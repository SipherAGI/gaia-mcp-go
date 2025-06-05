package imageutil

import (
	"context"
	"fmt"
	"image"
	"time"
)

// QuickProcessConfig provides a simple configuration builder
type QuickProcessConfig struct {
	processor *Processor
}

// NewQuickProcessor creates a processor with simplified configuration
func NewQuickProcessor() *QuickProcessConfig {
	return &QuickProcessConfig{
		processor: NewDefaultProcessor(),
	}
}

// WithMaxSize sets the maximum dimensions for image resizing
func (q *QuickProcessConfig) WithMaxSize(width, height int) *QuickProcessConfig {
	config := q.processor.config
	config.MaxWidth = width
	config.MaxHeight = height
	q.processor = NewProcessor(config)
	return q
}

// WithTimeout sets the HTTP timeout for downloads
func (q *QuickProcessConfig) WithTimeout(timeout time.Duration) *QuickProcessConfig {
	config := q.processor.config
	config.Timeout = timeout
	q.processor = NewProcessor(config)
	return q
}

// WithJPEGQuality sets the JPEG compression quality (1-100)
func (q *QuickProcessConfig) WithJPEGQuality(quality int) *QuickProcessConfig {
	config := q.processor.config
	config.JPEGQuality = quality
	q.processor = NewProcessor(config)
	return q
}

// Build returns the configured processor
func (q *QuickProcessConfig) Build() *Processor {
	return q.processor
}

// ProcessImage is a convenience function for simple image processing
func (q *QuickProcessConfig) ProcessImage(ctx context.Context, imageURL string) (string, error) {
	return q.processor.ProcessImageFromURL(ctx, imageURL)
}

// Convenience functions for common use cases

// ProcessImageQuick processes an image with default settings (1024x1024, 30s timeout)
func ProcessImageQuick(ctx context.Context, imageURL string) (string, error) {
	processor := NewDefaultProcessor()
	return processor.ProcessImageFromURL(ctx, imageURL)
}

// ProcessImageWithSize processes an image with custom maximum dimensions
func ProcessImageWithSize(ctx context.Context, imageURL string, maxWidth, maxHeight int) (string, error) {
	config := DefaultConfig()
	config.MaxWidth = maxWidth
	config.MaxHeight = maxHeight

	processor := NewProcessor(config)
	return processor.ProcessImageFromURL(ctx, imageURL)
}

// ProcessImageThumbnail creates a small thumbnail (256x256)
func ProcessImageThumbnail(ctx context.Context, imageURL string) (string, error) {
	return ProcessImageWithSize(ctx, imageURL, 256, 256)
}

// ProcessImageLarge processes an image for large display (2048x2048)
func ProcessImageLarge(ctx context.Context, imageURL string) (string, error) {
	return ProcessImageWithSize(ctx, imageURL, 2048, 2048)
}

// GetImageDimensions returns the dimensions of an image without processing it
func GetImageDimensions(ctx context.Context, imageURL string) (width int, height int, err error) {
	processor := NewDefaultProcessor()
	img, _, err := processor.DownloadImage(ctx, imageURL)
	if err != nil {
		return 0, 0, fmt.Errorf("downloading image: %w", err)
	}

	bounds := img.Bounds()
	return bounds.Dx(), bounds.Dy(), nil
}

// ValidateImageURL checks if a URL points to a valid image
func ValidateImageURL(ctx context.Context, imageURL string) error {
	processor := NewDefaultProcessor()
	_, _, err := processor.DownloadImage(ctx, imageURL)
	if err != nil {
		return fmt.Errorf("invalid image URL: %w", err)
	}
	return nil
}

// ResizeImageToExactSize resizes an image to exact dimensions (may distort aspect ratio)
func ResizeImageToExactSize(img image.Image, width, height int) image.Image {
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	// Use bilinear scaling for smooth resizing
	// Note: This will distort the image if aspect ratios don't match
	srcBounds := img.Bounds()
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Map destination coordinates to source coordinates
			srcX := int(float64(x) * float64(srcBounds.Dx()) / float64(width))
			srcY := int(float64(y) * float64(srcBounds.Dy()) / float64(height))

			// Ensure we don't go out of bounds
			if srcX >= srcBounds.Dx() {
				srcX = srcBounds.Dx() - 1
			}
			if srcY >= srcBounds.Dy() {
				srcY = srcBounds.Dy() - 1
			}

			dst.Set(x, y, img.At(srcBounds.Min.X+srcX, srcBounds.Min.Y+srcY))
		}
	}
	return dst
}

// MCP-specific convenience functions that return pure base64 and MIME type

// ProcessImageDefaultForMCP processes an image with default settings and returns data suitable for MCP
// Uses original default dimensions (1024x1024) and high quality (90) - may exceed MCP size limits for large images
func ProcessImageDefaultForMCP(ctx context.Context, imageURL string) (base64Data string, mimeType string, err error) {
	processor := NewDefaultProcessor()
	return processor.ProcessImageFromURLForMCP(ctx, imageURL)
}

// ProcessImageQuickForMCP processes an image with MCP-optimized settings and returns data suitable for MCP
// Uses smaller dimensions (512x512) and moderate quality (70) to stay under MCP size limits
func ProcessImageQuickForMCP(ctx context.Context, imageURL string) (base64Data string, mimeType string, err error) {
	// Use MCP-optimized configuration to avoid size limit errors
	config := DefaultConfig()
	config.MaxWidth = 512
	config.MaxHeight = 512
	config.JPEGQuality = 70 // Lower quality for smaller file size while maintaining visual quality

	processor := NewProcessor(config)
	return processor.ProcessImageFromURLForMCP(ctx, imageURL)
}

// ProcessImageWithSizeForMCP processes an image with custom dimensions and returns data suitable for MCP
func ProcessImageWithSizeForMCP(ctx context.Context, imageURL string, maxWidth, maxHeight int) (base64Data string, mimeType string, err error) {
	config := DefaultConfig()
	config.MaxWidth = maxWidth
	config.MaxHeight = maxHeight

	processor := NewProcessor(config)
	return processor.ProcessImageFromURLForMCP(ctx, imageURL)
}

// ProcessImageThumbnailForMCP creates a small thumbnail and returns data suitable for MCP
func ProcessImageThumbnailForMCP(ctx context.Context, imageURL string) (base64Data string, mimeType string, err error) {
	return ProcessImageWithSizeForMCP(ctx, imageURL, 256, 256)
}

// ProcessImageLargeForMCP processes an image for large display and returns data suitable for MCP
func ProcessImageLargeForMCP(ctx context.Context, imageURL string) (base64Data string, mimeType string, err error) {
	return ProcessImageWithSizeForMCP(ctx, imageURL, 2048, 2048)
}

// ProcessImageNoResize processes an image without any size constraints
// This preserves the original image dimensions while still performing format conversion and compression
func ProcessImageNoResize(ctx context.Context, imageURL string) (string, error) {
	config := DefaultConfig()
	// Set very large max dimensions to effectively disable resizing
	// Using a large number that's unlikely to be exceeded by normal images
	config.MaxWidth = 100000  // 100k pixels width
	config.MaxHeight = 100000 // 100k pixels height

	processor := NewProcessor(config)
	return processor.ProcessImageFromURL(ctx, imageURL)
}

// ProcessImageNoResizeForMCP processes an image without resizing and returns data suitable for MCP
// This is useful when you want to preserve original image dimensions for MCP responses
func ProcessImageNoResizeForMCP(ctx context.Context, imageURL string) (base64Data string, mimeType string, err error) {
	config := DefaultConfig()
	// Set very large max dimensions to effectively disable resizing
	config.MaxWidth = 100000  // 100k pixels width
	config.MaxHeight = 100000 // 100k pixels height

	processor := NewProcessor(config)
	return processor.ProcessImageFromURLForMCP(ctx, imageURL)
}
