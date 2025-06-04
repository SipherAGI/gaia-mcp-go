package api

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"gaia-mcp-go/pkg/httpclient"
	"gaia-mcp-go/pkg/imageutil"
	"gaia-mcp-go/pkg/shared"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

// Package api provides a client interface for interacting with the Gaia API.
//
// The Gaia API allows you to create SD styles, generate images, and upload images
// with multipart upload support. This package handles authentication, request
// formatting, and error handling for all Gaia API operations.
//
// Example usage:
//
//	cfg := api.GaiaApiConfig{
//		BaseUrl: "https://api.gaia.com",
//		ApiKey:  "your-api-key",
//	}
//	client := api.NewGaiaApi(cfg)
//
//	// Create a new style
//	style, err := client.CreateStyle(ctx, imageUrls, "My Style", &description)

// GaiaApi defines the interface for interacting with the Gaia API.
//
// This interface provides methods for creating SD styles, generating images,
// and uploading multiple images with automatic chunking and concurrent upload
// support for large files.
type GaiaApi interface {
	// CreateStyle creates a new SD (Stable Diffusion) style from provided images.
	//
	// Parameters:
	//   - ctx: Context for request cancellation and timeout control
	//   - imageUrls: Slice of HTTP(S) URLs pointing to reference images
	//   - name: Human-readable name for the style
	//   - description: Optional description of the style (can be nil)
	//
	// Returns the created SdStyle with its unique ID and metadata, or an error
	// if the creation fails due to invalid URLs, network issues, or API errors.
	CreateStyle(ctx context.Context, imageUrls []string, name string, description *string) (SdStyle, error)

	// GenerateImages creates a new image generation task using the Gaia AGI system.
	//
	// Parameters:
	//   - ctx: Context for request cancellation and timeout control
	//   - req: GenerateImagesRequest containing generation parameters like prompt,
	//          style settings, dimensions, and other generation options
	//
	// Returns ImageGeneratedResponse containing the task ID and status information,
	// or an error if the generation request fails validation or submission.
	GenerateImages(ctx context.Context, req GenerateImagesRequest) (ImageGeneratedResponse, error)

	// UploadImages uploads multiple images concurrently using multipart upload.
	//
	// This method downloads images from the provided URLs, processes them,
	// and uploads them to Gaia's storage using chunked multipart uploads
	// for efficient handling of large files.
	//
	// Parameters:
	//   - ctx: Context for request cancellation and timeout control
	//   - imageUrls: Slice of HTTP(S) URLs pointing to images to upload
	//   - associatedResource: Metadata about the resource these images are associated with
	//
	// Returns a slice of UploadFile containing the uploaded file metadata,
	// or an error if any uploads fail. Partial failures are reported in the error.
	UploadImages(ctx context.Context, imageUrls []string, associatedResource shared.FileAssociatedResource) ([]UploadFile, error)
}

// GaiaApiConfig holds the configuration needed to create a Gaia API client.
//
// Both BaseUrl and ApiKey are required fields. The BaseUrl should include
// the protocol (http:// or https://) and should not end with a trailing slash.
type GaiaApiConfig struct {
	// BaseUrl is the base URL of the Gaia API server (e.g., "https://api.gaia.com")
	BaseUrl string
	// ApiKey is the authentication token for accessing the Gaia API
	ApiKey string
}

// gaiaApi is the concrete implementation of the GaiaApi interface.
//
// This struct contains an HTTP client configured with the appropriate
// base URL, authentication headers, and timeout settings for Gaia API calls.
type gaiaApi struct {
	client *httpclient.Client
}

// NewGaiaApi creates a new Gaia API client with the provided configuration.
//
// The client is configured with:
//   - Bearer token authentication using the provided API key
//   - 60-second timeout for all API requests
//   - Automatic request/response JSON marshaling
//
// Parameters:
//   - cfg: GaiaApiConfig containing the base URL and API key
//
// Returns a GaiaApi interface implementation ready for use.
func NewGaiaApi(cfg GaiaApiConfig) GaiaApi {
	client := httpclient.New(httpclient.Config{
		BaseURL: cfg.BaseUrl,
		DefaultHeaders: map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", cfg.ApiKey),
		},
		Timeout: 60 * time.Second, // 60 seconds timeout for calling the API
	})
	return &gaiaApi{client: client}
}

// CreateStyle creates a new SD style from reference images.
//
// This method formats the provided image URLs into the expected API payload
// format, with each image assigned a default weight of 0.5. The style can
// be used later for image generation tasks.
//
// The method handles:
//   - Image URL validation and formatting
//   - Optional description parameter
//   - JSON marshaling/unmarshaling
//   - Error processing and wrapping
//
// Parameters:
//   - ctx: Request context for cancellation and timeout
//   - imageUrls: URLs of reference images (must be HTTP/HTTPS)
//   - name: Display name for the style
//   - description: Optional style description (pass nil if not needed)
//
// Returns the created SdStyle containing the style ID and metadata,
// or an error if creation fails.
func (a *gaiaApi) CreateStyle(ctx context.Context, imageUrls []string, name string, description *string) (SdStyle, error) {
	// Formatting imageUrls to be an array of images
	images := make([]map[string]interface{}, len(imageUrls))
	for i, imageUrl := range imageUrls {
		images[i] = map[string]interface{}{
			"url":    imageUrl,
			"weight": 0.5,
		}
	}

	payload := map[string]interface{}{
		"images": images,
		"name":   name,
	}

	if description != nil {
		payload["description"] = *description
	}

	// Use the type-safe As[T] function - cleaner and more idiomatic
	sdStyle, err := httpclient.As[SdStyle](
		a.client.PostJSON(ctx, "/api/sd-styles", payload, map[string]string{}),
	)
	if err != nil {
		return SdStyle{}, ProcessError(err)
	}

	return sdStyle, nil
}

// GenerateImages submits an image generation request to the Gaia AGI system.
//
// This method calls the agi-tasks/create-task endpoint to start a new
// image generation task. The task runs asynchronously, and you can
// check its status using the returned task information.
//
// Parameters:
//   - ctx: Request context for cancellation and timeout
//   - req: Complete generation request with prompt, style, dimensions, etc.
//
// Returns ImageGeneratedResponse with task ID and initial status,
// or an error if the request fails validation or submission.
func (a *gaiaApi) GenerateImages(ctx context.Context, req GenerateImagesRequest) (ImageGeneratedResponse, error) {
	// Use the type-safe As[T] function - cleaner and more idiomatic
	imageGeneratedResponse, err := httpclient.As[ImageGeneratedResponse](
		a.client.PostJSON(ctx, "/api/recipe/agi-tasks/create-task", req, map[string]string{}),
	)
	if err != nil {
		return ImageGeneratedResponse{}, ProcessError(err)
	}

	return imageGeneratedResponse, nil
}

// UploadImages handles concurrent multipart upload of multiple images.
//
// This method performs the following steps for each image:
//  1. Downloads and validates the image from the provided URL
//  2. Processes the image to extract dimensions and convert to bytes
//  3. Initializes a multipart upload session with the API
//  4. Uploads image data in chunks concurrently for better performance
//  5. Completes the multipart upload process
//
// The method uses goroutines for concurrent chunk uploads within each image
// and continues processing other images even if some fail. All failures
// are collected and reported in the final error.
//
// Parameters:
//   - ctx: Request context for cancellation and timeout control
//   - imageUrls: Slice of HTTP/HTTPS URLs pointing to images
//   - associatedResource: Metadata linking uploads to a specific resource
//
// Returns a slice of successfully uploaded files, or an error containing
// details of any failures. If some uploads succeed and others fail,
// only the error is returned with failure details.
func (a *gaiaApi) UploadImages(ctx context.Context, imageUrls []string, associatedResource shared.FileAssociatedResource) ([]UploadFile, error) {
	var uploadedFiles []UploadFile
	var failedFiles []map[string]string

	for _, imageUrl := range imageUrls {
		if !strings.HasPrefix(imageUrl, "http") {
			failedFiles = append(failedFiles, map[string]string{
				"url":   imageUrl,
				"error": "URL must start with http:// or https://",
			})
			continue
		}

		// process the image
		imageData, _, w, h, err := a.processImage(ctx, imageUrl)
		if err != nil {
			failedFiles = append(failedFiles, map[string]string{
				"url":   imageUrl,
				"error": err.Error(),
			})
			continue
		}

		// Initialize the upload file
		initUploadResponse, err := a.initUploadImage(ctx, imageData, w, h, associatedResource)
		if err != nil {
			failedFiles = append(failedFiles, map[string]string{
				"url":   imageUrl,
				"error": err.Error(),
			})
			continue
		}

		// Upload chunk concurrently
		var wg sync.WaitGroup
		uploadParts := make([]*UploadPart, len(initUploadResponse.UploadUrls))
		uploadErrs := make([]error, len(initUploadResponse.UploadUrls))

		for i, url := range initUploadResponse.UploadUrls {
			wg.Add(1)
			go func(i int, url string) {
				defer wg.Done()

				// Calculate the chunk boundaries
				start := i * shared.UPLOAD_CHUNK_SIZE
				end := start + shared.UPLOAD_CHUNK_SIZE
				if end > len(imageData) {
					end = len(imageData)
				}

				chunk := imageData[start:end]
				partNumber := i + 1

				// Upload the chunk
				part, err := a.uploadChunk(ctx, chunk, url, partNumber)
				if err != nil {
					uploadErrs[i] = err
					return
				}

				uploadParts[i] = part
			}(i, url)
		}

		wg.Wait()

		// Check for errors
		var hasErrors bool
		for _, err := range uploadErrs {
			if err != nil {
				hasErrors = true
				break
			}
		}

		if hasErrors {
			failedFiles = append(failedFiles, map[string]string{
				"url":   imageUrl,
				"error": "Failed to upload some chunks",
			})
			continue
		}

		// Convert to slice without nil pointers
		var parts []UploadPart
		for _, part := range uploadParts {
			if part != nil {
				parts = append(parts, *part)
			}
		}

		// Complete the upload
		if err := a.completeUpload(ctx, initUploadResponse.Key, initUploadResponse.UploadId, parts); err != nil {
			failedFiles = append(failedFiles, map[string]string{
				"url":   imageUrl,
				"error": err.Error(),
			})
			continue
		}

		uploadedFiles = append(uploadedFiles, initUploadResponse.File)
	}

	if len(failedFiles) > 0 {
		return nil, fmt.Errorf("failed to upload some files: %v", failedFiles)
	}

	return uploadedFiles, nil
}

// processImage downloads, processes, and extracts metadata from an image URL.
//
// This method performs several operations on the source image:
//  1. Downloads the image from the provided URL using the imageutil package
//  2. Processes it without resizing to maintain original quality
//  3. Converts the base64-encoded image data to raw bytes
//  4. Extracts image dimensions (width and height)
//
// The method is used internally by UploadImages to prepare image data
// for multipart upload. It handles various image formats and ensures
// the data is ready for chunked upload operations.
//
// Parameters:
//   - ctx: Request context for cancellation and timeout control
//   - imageUrl: HTTP/HTTPS URL pointing to the image to process
//
// Returns:
//   - imageData: Raw image bytes ready for upload
//   - mimeType: MIME type of the processed image (e.g., "image/png")
//   - w: Image width in pixels
//   - h: Image height in pixels
//   - err: Error if download, processing, or dimension extraction fails
func (a *gaiaApi) processImage(ctx context.Context, imageUrl string) (imageData []byte, mimeType string, w, h int, err error) {
	var base64Data string

	// Fetch the image
	base64Data, mimeType, err = imageutil.ProcessImageNoResizeForMCP(ctx, imageUrl)
	if err != nil {
		return nil, "", 0, 0, fmt.Errorf("failed to process image: %w", err)
	}

	// Convert base64 data to bytes
	imageData, err = base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return nil, "", 0, 0, fmt.Errorf("failed to decode base64 data: %w", err)
	}

	// Get dimensions of the image
	w, h, err = imageutil.GetImageDimensions(ctx, imageUrl)
	if err != nil {
		return nil, "", 0, 0, fmt.Errorf("failed to get image dimensions: %w", err)
	}

	return imageData, mimeType, w, h, nil
}

// initUploadImage initializes a multipart upload session for a single image.
//
// This method creates the necessary setup for uploading large images using
// chunked multipart upload. It sends image metadata including dimensions,
// file size, and associated resource information to the Gaia API to receive
// presigned upload URLs for each chunk.
//
// The method:
//  1. Constructs the upload initialization payload with file metadata
//  2. Includes image dimensions, MIME type, and calculated file size
//  3. Specifies the chunk size for multipart upload
//  4. Associates the upload with a specific resource (e.g., style, task)
//  5. Returns presigned URLs for each chunk and upload tracking information
//
// Parameters:
//   - ctx: Request context for cancellation and timeout control
//   - imageData: Raw image bytes to be uploaded
//   - w: Image width in pixels (for metadata)
//   - h: Image height in pixels (for metadata)
//   - associatedResource: Resource metadata linking this upload to a specific entity
//
// Returns:
//   - InitUploadResponse: Contains upload ID, presigned URLs, and file metadata
//   - error: Error if initialization request fails or returns invalid response
func (a *gaiaApi) initUploadImage(
	ctx context.Context,
	imageData []byte,
	w, h int,
	associatedResource shared.FileAssociatedResource,
) (*InitUploadResponse, error) {
	// Prepare the request payload
	payload := map[string]interface{}{
		"files": []map[string]interface{}{
			{
				"filename": fmt.Sprintf("image_%d.png", time.Now().Unix()),
				"mimetype": "image/png",
				"metadata": map[string]int{
					"width":  w,
					"height": h,
				},
				"fileSize": len(imageData),
			},
		},
		"associatedResource": associatedResource,
		"chunkSize":          shared.UPLOAD_CHUNK_SIZE,
	}

	// Send the request - the API returns an array of InitUploadResponse
	initUploadResponses, err := httpclient.As[[]InitUploadResponse](
		a.client.PostJSON(ctx, "/api/upload/initialize", payload, map[string]string{}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	// Since we're uploading one file, we expect exactly one response
	if len(initUploadResponses) < 1 {
		return nil, fmt.Errorf("no upload responses received")
	}

	return &initUploadResponses[0], nil
}

// uploadChunk uploads a single data chunk directly to a presigned S3 URL.
//
// This method handles the actual upload of individual chunks in a multipart
// upload process. It bypasses the regular API client to send data directly
// to AWS S3 using presigned URLs for better performance and reduced server load.
//
// The method:
//  1. Creates a direct HTTP PUT request to the presigned S3 URL
//  2. Sets appropriate headers for S3 compatibility (Content-Type, Content-Length)
//  3. Uses a dedicated HTTP client with extended timeout for large chunks
//  4. Validates the upload response and extracts the required ETag
//  5. Returns upload part information needed for multipart completion
//
// Parameters:
//   - ctx: Request context for cancellation and timeout control
//   - chunk: Raw data bytes for this specific chunk
//   - url: Presigned S3 URL for uploading this chunk
//   - partNumber: Sequential part number (1-based) for this chunk
//
// Returns:
//   - UploadPart: Contains ETag and part number required for upload completion
//   - error: Error if HTTP request fails, upload is rejected, or ETag is missing
func (a *gaiaApi) uploadChunk(ctx context.Context, chunk []byte, url string, partNumber int) (*UploadPart, error) {
	// Create a direct HTTP request to the presigned S3 URL
	// Don't use a.client.PUT() because it prepends the base URL
	req, err := http.NewRequestWithContext(ctx, "PUT", url, bytes.NewReader(chunk))
	if err != nil {
		return nil, fmt.Errorf("failed to create request for chunk %d: %w", partNumber, err)
	}

	// Set required headers for S3 upload
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Content-Length", fmt.Sprintf("%d", len(chunk)))

	// Create a new HTTP client for direct S3 uploads
	httpClient := &http.Client{
		Timeout: 60 * time.Second, // Longer timeout for large uploads
	}

	// Execute the request
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to upload chunk %d: %w", partNumber, err)
	}
	defer resp.Body.Close()

	// Check for successful upload (S3 returns 200 for successful chunk uploads)
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("chunk %d upload failed with status %d: %s", partNumber, resp.StatusCode, string(body))
	}

	// Extract ETag from response headers (required for multipart upload completion)
	etag := resp.Header.Get("ETag")
	if etag == "" {
		return nil, fmt.Errorf("missing ETag in response for chunk %d", partNumber)
	}

	uploadPart := &UploadPart{
		ETag:       etag,
		PartNumber: partNumber,
	}

	return uploadPart, nil
}

// completeUpload finalizes a multipart upload by combining all uploaded chunks.
//
// This method notifies the Gaia API that all chunks have been successfully
// uploaded and requests the finalization of the multipart upload process.
// The API will then combine all parts into a single file and make it
// available for use.
//
// The method:
//  1. Constructs completion payload with upload ID, key, and all part information
//  2. Sends the completion request to the API endpoint
//  3. Validates the response to ensure successful completion
//  4. Handles both 200 OK and 201 Created as successful completion statuses
//
// Parameters:
//   - ctx: Request context for cancellation and timeout control
//   - key: Unique identifier for the upload session
//   - uploadId: Multipart upload ID from the initialization response
//   - parts: Slice of UploadPart containing ETag and part number for each chunk
//
// Returns:
//   - error: Error if completion request fails or server rejects the completion
func (a *gaiaApi) completeUpload(ctx context.Context, key, uploadId string, parts []UploadPart) error {
	payload := []map[string]interface{}{
		{
			"key":      key,
			"uploadId": uploadId,
			"parts":    parts,
		},
	}

	// Send the request
	res, err := a.client.POST(ctx, "/api/upload/complete", payload, map[string]string{})
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer res.Body.Close()

	// Read the response body for proper error handling
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for successful completion
	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to complete upload (status %d): %s", res.StatusCode, string(body))
	}

	return nil
}
