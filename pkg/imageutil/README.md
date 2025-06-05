# ImageUtil Package

The `imageutil` package provides robust image processing utilities for downloading, resizing, and encoding images. It's designed to be reusable across different tools and functions in the Gaia MCP Go application.

## Features

- **Download images** from URLs with configurable timeouts
- **Resize images** while maintaining aspect ratio
- **Base64 encoding** with data URL format
- **Multiple image formats** support (JPEG, PNG, etc.)
- **Configurable processing** options
- **Error handling** with proper context propagation

## Quick Start

### Basic Usage

```go
import (
    "context"
    "gaia-mcp-go/pkg/imageutil"
)

// Process an image with default settings (1024x1024 max)
base64Image, err := imageutil.ProcessImageQuick(ctx, "https://example.com/image.jpg")
if err != nil {
    log.Fatal(err)
}
```

### Custom Configuration

```go
// Create a processor with custom settings
processor := imageutil.NewQuickProcessor().
    WithMaxSize(512, 512).
    WithTimeout(10 * time.Second).
    WithJPEGQuality(85).
    Build()

base64Image, err := processor.ProcessImageFromURL(ctx, imageURL)
```

### Advanced Usage

```go
// Create processor with full configuration
config := imageutil.ProcessorConfig{
    MaxWidth:    2048,
    MaxHeight:   2048,
    Timeout:     60 * time.Second,
    JPEGQuality: 95,
    UserAgent:   "MyApp/1.0",
}

processor := imageutil.NewProcessor(config)
base64Image, err := processor.ProcessImageFromURL(ctx, imageURL)
```

## Convenience Functions

### Predefined Sizes

```go
// Create thumbnail (256x256)
thumbnail, err := imageutil.ProcessImageThumbnail(ctx, imageURL)

// Create large image (2048x2048)
largeImage, err := imageutil.ProcessImageLarge(ctx, imageURL)

// Custom size
customImage, err := imageutil.ProcessImageWithSize(ctx, imageURL, 800, 600)
```

### MCP-Specific Functions

For MCP (Model Context Protocol) applications that need pure base64 data without data URL prefixes:

```go
// Process image for MCP with optimized settings (512x512, quality 70)
// Recommended for most MCP use cases to avoid size limits
base64Data, mimeType, err := imageutil.ProcessImageQuickForMCP(ctx, imageURL)

// Process image for MCP with default settings (1024x1024, quality 90)
// May exceed MCP size limits for large images
base64Data, mimeType, err := imageutil.ProcessImageDefaultForMCP(ctx, imageURL)

// Create thumbnail for MCP (256x256)
base64Data, mimeType, err := imageutil.ProcessImageThumbnailForMCP(ctx, imageURL)

// Create large image for MCP (2048x2048)
base64Data, mimeType, err := imageutil.ProcessImageLargeForMCP(ctx, imageURL)

// Custom size for MCP
base64Data, mimeType, err := imageutil.ProcessImageWithSizeForMCP(ctx, imageURL, 800, 600)

// Use in MCP tool result
return mcp.NewToolResultImage("Description", base64Data, mimeType)
```

**Note**: `ProcessImageQuickForMCP` is now optimized for MCP size constraints (512x512 max, 70% quality) to prevent "result exceeds maximum length" errors in Claude Desktop and other MCP clients.

### Utility Functions

```go
// Get image dimensions without processing
width, height, err := imageutil.GetImageDimensions(ctx, imageURL)

// Validate if URL points to a valid image
err := imageutil.ValidateImageURL(ctx, imageURL)
```

## Configuration Options

| Option        | Description                       | Default           |
| ------------- | --------------------------------- | ----------------- |
| `MaxWidth`    | Maximum width for resized images  | 1024              |
| `MaxHeight`   | Maximum height for resized images | 1024              |
| `Timeout`     | HTTP request timeout              | 30 seconds        |
| `JPEGQuality` | JPEG compression quality (1-100)  | 90                |
| `UserAgent`   | User agent for HTTP requests      | "Gaia-MCP-Go/1.0" |

## Integration Examples

### In Tools

```go
type MyTool struct {
    imageProcessor *imageutil.Processor
}

func NewMyTool() *MyTool {
    return &MyTool{
        imageProcessor: imageutil.NewDefaultProcessor(),
    }
}

func (t *MyTool) ProcessImage(ctx context.Context, imageURL string) (string, error) {
    return t.imageProcessor.ProcessImageFromURL(ctx, imageURL)
}
```

### Custom Processors

```go
// Different tools can use different configurations
thumbnailProcessor := imageutil.NewQuickProcessor().
    WithMaxSize(256, 256).
    Build()

highQualityProcessor := imageutil.NewQuickProcessor().
    WithMaxSize(2048, 2048).
    WithJPEGQuality(95).
    Build()
```

## Error Handling

The package provides comprehensive error handling:

- **Network errors**: Connection timeouts, DNS failures
- **HTTP errors**: Non-200 status codes, server errors
- **Image errors**: Invalid formats, corrupted data
- **Processing errors**: Encoding failures, memory issues

All errors are wrapped with context for easy debugging:

```go
base64Image, err := processor.ProcessImageFromURL(ctx, imageURL)
if err != nil {
    // Error will contain context about which step failed
    log.Printf("Image processing failed: %v", err)
}
```

## Performance Considerations

- **Memory usage**: Large images are processed in memory
- **Network bandwidth**: Images are downloaded completely before processing
- **CPU usage**: High-quality resizing uses more CPU
- **Caching**: Consider implementing caching for frequently accessed images

## Thread Safety

All functions and methods in this package are thread-safe and can be used concurrently across multiple goroutines.

## Output Format

The package returns base64-encoded images with data URL format:

```
data:image/jpeg;base64,/9j/4AAQSkZJRgABAQAA...
```

This format can be used directly in:

- HTML `<img>` tags
- CSS `background-image` properties
- JSON responses
- Database storage
