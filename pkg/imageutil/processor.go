package imageutil

import (
	"context"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"strings"
	"time"

	"golang.org/x/image/draw"
)

// ProcessorConfig holds configuration for image processing
type ProcessorConfig struct {
	// MaxWidth is the maximum width for resized images
	MaxWidth int
	// MaxHeight is the maximum height for resized images
	MaxHeight int
	// Timeout for HTTP requests when downloading images
	Timeout time.Duration
	// Quality for JPEG encoding (1-100)
	JPEGQuality int
	// UserAgent for HTTP requests
	UserAgent string
}

// DefaultConfig returns a sensible default configuration
func DefaultConfig() ProcessorConfig {
	return ProcessorConfig{
		MaxWidth:    1024,
		MaxHeight:   1024,
		Timeout:     30 * time.Second,
		JPEGQuality: 90,
		UserAgent:   "Gaia-MCP-Go/1.0",
	}
}

// Processor handles image processing operations
type Processor struct {
	config ProcessorConfig
	client *http.Client
}

// NewProcessor creates a new image processor with the given configuration
func NewProcessor(config ProcessorConfig) *Processor {
	return &Processor{
		config: config,
		client: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

// NewDefaultProcessor creates a new image processor with default configuration
func NewDefaultProcessor() *Processor {
	return NewProcessor(DefaultConfig())
}

// ProcessImageFromURL downloads an image from URL, resizes it, and returns base64 encoded string
func (p *Processor) ProcessImageFromURL(ctx context.Context, imageURL string) (string, error) {
	// Step 1: Download the image
	img, format, err := p.downloadImage(ctx, imageURL)
	if err != nil {
		return "", fmt.Errorf("downloading image: %w", err)
	}

	// Step 2: Resize the image
	resizedImg := p.resizeImage(img)

	// Step 3: Encode to base64
	base64Str, err := p.encodeImageToBase64(resizedImg, format)
	if err != nil {
		return "", fmt.Errorf("encoding image to base64: %w", err)
	}

	return base64Str, nil
}

// ProcessImageFromURLForMCP downloads an image from URL, resizes it, and returns pure base64 data and MIME type for MCP
func (p *Processor) ProcessImageFromURLForMCP(ctx context.Context, imageURL string) (base64Data string, mimeType string, err error) {
	// Step 1: Download the image
	img, format, err := p.downloadImage(ctx, imageURL)
	if err != nil {
		return "", "", fmt.Errorf("downloading image: %w", err)
	}

	// Step 2: Resize the image
	resizedImg := p.resizeImage(img)

	// Step 3: Encode to base64 (pure base64, no data URL prefix)
	base64Data, mimeType, err = p.encodeImageToBase64Pure(resizedImg, format)
	if err != nil {
		return "", "", fmt.Errorf("encoding image to base64: %w", err)
	}

	return base64Data, mimeType, nil
}

// DownloadImage downloads an image from URL and returns the decoded image
func (p *Processor) DownloadImage(ctx context.Context, url string) (image.Image, string, error) {
	return p.downloadImage(ctx, url)
}

// ResizeImage resizes an image to fit within the configured dimensions
func (p *Processor) ResizeImage(img image.Image) image.Image {
	return p.resizeImage(img)
}

// EncodeImageToBase64 encodes an image to base64 string with data URL prefix
func (p *Processor) EncodeImageToBase64(img image.Image, format string) (string, error) {
	return p.encodeImageToBase64(img, format)
}

// EncodeImageToBase64Pure encodes an image to pure base64 string and returns MIME type
func (p *Processor) EncodeImageToBase64Pure(img image.Image, format string) (string, string, error) {
	return p.encodeImageToBase64Pure(img, format)
}

// downloadImage downloads an image from the given URL and returns the decoded image
func (p *Processor) downloadImage(ctx context.Context, url string) (image.Image, string, error) {
	// Create request with context
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, "", fmt.Errorf("creating request: %w", err)
	}

	// Set user agent to avoid blocking
	req.Header.Set("User-Agent", p.config.UserAgent)

	// Download the image
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("downloading image: %w", err)
	}
	defer resp.Body.Close()

	// Check if the response is successful
	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read response body
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("reading response body: %w", err)
	}

	// Decode the image
	img, format, err := image.Decode(strings.NewReader(string(data)))
	if err != nil {
		return nil, "", fmt.Errorf("decoding image: %w", err)
	}

	return img, format, nil
}

// resizeImage resizes an image to fit within maxWidth x maxHeight while maintaining aspect ratio
func (p *Processor) resizeImage(src image.Image) image.Image {
	srcBounds := src.Bounds()
	srcWidth := srcBounds.Dx()
	srcHeight := srcBounds.Dy()

	maxWidth := p.config.MaxWidth
	maxHeight := p.config.MaxHeight

	// If image is already smaller than max dimensions, return as is
	if srcWidth <= maxWidth && srcHeight <= maxHeight {
		return src
	}

	// Calculate scaling factor to maintain aspect ratio
	scaleX := float64(maxWidth) / float64(srcWidth)
	scaleY := float64(maxHeight) / float64(srcHeight)
	scale := scaleX
	if scaleY < scaleX {
		scale = scaleY
	}

	// Calculate new dimensions
	newWidth := int(float64(srcWidth) * scale)
	newHeight := int(float64(srcHeight) * scale)

	// Create new image
	dst := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	// Resize using high-quality scaling
	draw.BiLinear.Scale(dst, dst.Bounds(), src, srcBounds, draw.Over, nil)

	return dst
}

// encodeImageToBase64Pure encodes an image to pure base64 string without data URL prefix
func (p *Processor) encodeImageToBase64Pure(img image.Image, format string) (base64Data string, mimeType string, err error) {
	var buf strings.Builder

	// Determine the MIME type and encoding based on original format
	switch format {
	case "jpeg", "jpg":
		mimeType = "image/jpeg"
		if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: p.config.JPEGQuality}); err != nil {
			return "", "", fmt.Errorf("encoding JPEG: %w", err)
		}
	case "png":
		mimeType = "image/png"
		if err := png.Encode(&buf, img); err != nil {
			return "", "", fmt.Errorf("encoding PNG: %w", err)
		}
	default:
		// Default to PNG for unknown formats
		mimeType = "image/png"
		if err := png.Encode(&buf, img); err != nil {
			return "", "", fmt.Errorf("encoding PNG: %w", err)
		}
	}

	// Encode to base64 (pure base64, no prefix)
	encoded := base64.StdEncoding.EncodeToString([]byte(buf.String()))

	return encoded, mimeType, nil
}

// encodeImageToBase64 encodes an image to base64 string with data URL prefix
func (p *Processor) encodeImageToBase64(img image.Image, format string) (string, error) {
	// Use the pure base64 function and add the data URL prefix
	encoded, mimeType, err := p.encodeImageToBase64Pure(img, format)
	if err != nil {
		return "", err
	}

	// Return with data URL prefix
	return fmt.Sprintf("data:%s;base64,%s", mimeType, encoded), nil
}
