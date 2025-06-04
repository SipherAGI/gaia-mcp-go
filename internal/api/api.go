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

const (
	CHUNK_SIZE = 1024 * 1024 * 10 // 10MB chunks
)

type GaiaApi interface {
	CreateStyle(ctx context.Context, imageUrls []string, name string, description *string) (SdStyle, error)
	GenerateImages(ctx context.Context, req GenerateImagesRequest) (ImageGeneratedResponse, error)
	UploadImages(ctx context.Context, imageUrls []string, associatedResource shared.FileAssociatedResource) ([]UploadFile, error)
}

type GaiaApiConfig struct {
	BaseUrl string
	ApiKey  string
}

type gaiaApi struct {
	client *httpclient.Client
}

func NewGaiaApi(cfg GaiaApiConfig) GaiaApi {
	client := httpclient.New(httpclient.Config{
		BaseURL: cfg.BaseUrl,
		DefaultHeaders: map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", cfg.ApiKey),
		},
	})
	return &gaiaApi{client: client}
}

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
		return SdStyle{}, err
	}

	return sdStyle, nil
}

func (a *gaiaApi) GenerateImages(ctx context.Context, req GenerateImagesRequest) (ImageGeneratedResponse, error) {
	// Use the type-safe As[T] function - cleaner and more idiomatic
	imageGeneratedResponse, err := httpclient.As[ImageGeneratedResponse](
		a.client.PostJSON(ctx, "/api/recipe/agi-tasks/create-task", req, map[string]string{}),
	)
	if err != nil {
		return ImageGeneratedResponse{}, err
	}

	return imageGeneratedResponse, nil
}

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
				start := i * CHUNK_SIZE
				end := start + CHUNK_SIZE
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
		"chunkSize":          CHUNK_SIZE,
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
