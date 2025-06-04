package imageutil

import (
	"context"
	"gaia-mcp-go/internal/testutil"
	"image"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestDefaultConfig tests the default configuration
func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	assert.Equal(t, 1024, config.MaxWidth)
	assert.Equal(t, 1024, config.MaxHeight)
	assert.Equal(t, 30*time.Second, config.Timeout)
	assert.Equal(t, 90, config.JPEGQuality)
	assert.Equal(t, "Gaia-MCP-Go/1.0", config.UserAgent)
}

// TestNewProcessor tests processor creation
func TestNewProcessor(t *testing.T) {
	config := ProcessorConfig{
		MaxWidth:    512,
		MaxHeight:   512,
		Timeout:     10 * time.Second,
		JPEGQuality: 85,
		UserAgent:   "Test/1.0",
	}

	processor := NewProcessor(config)

	assert.NotNil(t, processor)
	assert.Equal(t, config, processor.config)
	assert.NotNil(t, processor.client)
	assert.Equal(t, config.Timeout, processor.client.Timeout)
}

// TestNewDefaultProcessor tests the default processor
func TestNewDefaultProcessor(t *testing.T) {
	processor := NewDefaultProcessor()

	assert.NotNil(t, processor)
	assert.Equal(t, DefaultConfig(), processor.config)
}

// TestProcessImageFromURL tests the complete image processing pipeline
func TestProcessImageFromURL(t *testing.T) {
	// Create a test server with mock image
	testServer := testutil.NewTestServer()
	defer testServer.Close()

	// Create a simple 1x1 PNG image
	mockImageData := testutil.CreateMockImage()

	t.Run("Successful image processing", func(t *testing.T) {
		// Setup mock response
		testServer.AddResponse("GET", "/test-image.png", testutil.MockResponse{
			StatusCode: http.StatusOK,
			Body:       mockImageData,
			Headers: map[string]string{
				"Content-Type": "image/png",
			},
		})

		processor := NewDefaultProcessor()
		ctx := context.Background()

		result, err := processor.ProcessImageFromURL(ctx, testServer.URL+"/test-image.png")

		assert.NoError(t, err)
		assert.NotEmpty(t, result)

		// Check that result is a valid base64 data URL
		assert.True(t, strings.HasPrefix(result, "data:image/"), "Result should be a data URL")
		assert.Contains(t, result, "base64,", "Result should contain base64 data")
	})

	t.Run("HTTP error handling", func(t *testing.T) {
		// Setup mock 404 response
		testServer.AddResponse("GET", "/not-found.png", testutil.MockResponse{
			StatusCode: http.StatusNotFound,
			Body:       "Not Found",
		})

		processor := NewDefaultProcessor()
		ctx := context.Background()

		result, err := processor.ProcessImageFromURL(ctx, testServer.URL+"/not-found.png")

		assert.Error(t, err)
		assert.Empty(t, result)
		assert.Contains(t, err.Error(), "downloading image")
	})

	t.Run("Invalid URL", func(t *testing.T) {
		processor := NewDefaultProcessor()
		ctx := context.Background()

		result, err := processor.ProcessImageFromURL(ctx, "invalid-url")

		assert.Error(t, err)
		assert.Empty(t, result)
	})

	t.Run("Context timeout", func(t *testing.T) {
		// Setup slow response
		testServer.AddResponse("GET", "/slow-image.png", testutil.MockResponse{
			StatusCode: http.StatusOK,
			Body:       mockImageData,
			Delay:      2 * time.Second,
		})

		processor := NewDefaultProcessor()
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
		defer cancel()

		result, err := processor.ProcessImageFromURL(ctx, testServer.URL+"/slow-image.png")

		assert.Error(t, err)
		assert.Empty(t, result)
	})
}

// TestProcessImageFromURLPure tests the pure base64 processing
func TestProcessImageFromURLPure(t *testing.T) {
	testServer := testutil.NewTestServer()
	defer testServer.Close()

	mockImageData := testutil.CreateMockImage()

	testServer.AddResponse("GET", "/test-image.png", testutil.MockResponse{
		StatusCode: http.StatusOK,
		Body:       mockImageData,
		Headers: map[string]string{
			"Content-Type": "image/png",
		},
	})

	processor := NewDefaultProcessor()
	ctx := context.Background()

	base64Data, mimeType, err := processor.ProcessImageFromURLForMCP(ctx, testServer.URL+"/test-image.png")

	assert.NoError(t, err)
	assert.NotEmpty(t, base64Data)
	assert.Equal(t, "image/png", mimeType)

	// Verify base64 data doesn't contain data URL prefix
	assert.False(t, strings.HasPrefix(base64Data, "data:"), "Pure base64 should not have data URL prefix")
}

// TestResizeImage tests the image resizing functionality
func TestResizeImage(t *testing.T) {
	processor := NewProcessor(ProcessorConfig{
		MaxWidth:  100,
		MaxHeight: 100,
	})

	t.Run("Resize large image", func(t *testing.T) {
		// Create a 200x200 test image
		src := image.NewRGBA(image.Rect(0, 0, 200, 200))

		resized := processor.ResizeImage(src)
		bounds := resized.Bounds()

		// Image should be resized to fit within 100x100
		assert.True(t, bounds.Dx() <= 100, "Width should be <= 100")
		assert.True(t, bounds.Dy() <= 100, "Height should be <= 100")
	})

	t.Run("Don't resize small image", func(t *testing.T) {
		// Create a 50x50 test image
		src := image.NewRGBA(image.Rect(0, 0, 50, 50))

		resized := processor.ResizeImage(src)
		bounds := resized.Bounds()

		// Image should remain the same size
		assert.Equal(t, 50, bounds.Dx())
		assert.Equal(t, 50, bounds.Dy())
	})

	t.Run("Maintain aspect ratio", func(t *testing.T) {
		// Create a 200x100 test image (2:1 aspect ratio)
		src := image.NewRGBA(image.Rect(0, 0, 200, 100))

		resized := processor.ResizeImage(src)
		bounds := resized.Bounds()

		// Should maintain 2:1 aspect ratio
		aspectRatio := float64(bounds.Dx()) / float64(bounds.Dy())
		assert.InDelta(t, 2.0, aspectRatio, 0.1, "Aspect ratio should be maintained")
	})
}

// TestEncodeImageToBase64 tests the base64 encoding
func TestEncodeImageToBase64(t *testing.T) {
	processor := NewDefaultProcessor()

	// Create a simple test image
	testImage := image.NewRGBA(image.Rect(0, 0, 10, 10))

	t.Run("Encode PNG", func(t *testing.T) {
		result, err := processor.EncodeImageToBase64(testImage, "png")

		assert.NoError(t, err)
		assert.NotEmpty(t, result)
		assert.True(t, strings.HasPrefix(result, "data:image/png;base64,"))
	})

	t.Run("Encode JPEG", func(t *testing.T) {
		result, err := processor.EncodeImageToBase64(testImage, "jpeg")

		assert.NoError(t, err)
		assert.NotEmpty(t, result)
		assert.True(t, strings.HasPrefix(result, "data:image/jpeg;base64,"))
	})

	t.Run("Encode unknown format defaults to PNG", func(t *testing.T) {
		result, err := processor.EncodeImageToBase64(testImage, "unknown")

		assert.NoError(t, err)
		assert.NotEmpty(t, result)
		assert.True(t, strings.HasPrefix(result, "data:image/png;base64,"))
	})
}

// TestEncodeImageToBase64Pure tests the pure base64 encoding
func TestEncodeImageToBase64Pure(t *testing.T) {
	processor := NewDefaultProcessor()

	// Create a simple test image
	testImage := image.NewRGBA(image.Rect(0, 0, 10, 10))

	t.Run("Encode PNG pure", func(t *testing.T) {
		base64Data, mimeType, err := processor.EncodeImageToBase64Pure(testImage, "png")

		assert.NoError(t, err)
		assert.NotEmpty(t, base64Data)
		assert.Equal(t, "image/png", mimeType)
		assert.False(t, strings.HasPrefix(base64Data, "data:"))
	})

	t.Run("Encode JPEG pure", func(t *testing.T) {
		base64Data, mimeType, err := processor.EncodeImageToBase64Pure(testImage, "jpeg")

		assert.NoError(t, err)
		assert.NotEmpty(t, base64Data)
		assert.Equal(t, "image/jpeg", mimeType)
		assert.False(t, strings.HasPrefix(base64Data, "data:"))
	})
}

// TestDownloadImage tests the image download functionality
func TestDownloadImage(t *testing.T) {
	testServer := testutil.NewTestServer()
	defer testServer.Close()

	mockImageData := testutil.CreateMockImage()

	processor := NewDefaultProcessor()

	t.Run("Successful download", func(t *testing.T) {
		testServer.AddResponse("GET", "/test-image.png", testutil.MockResponse{
			StatusCode: http.StatusOK,
			Body:       mockImageData,
			Headers: map[string]string{
				"Content-Type": "image/png",
			},
		})

		ctx := context.Background()
		img, format, err := processor.DownloadImage(ctx, testServer.URL+"/test-image.png")

		assert.NoError(t, err)
		assert.NotNil(t, img)
		assert.Equal(t, "png", format)
	})

	t.Run("HTTP error", func(t *testing.T) {
		testServer.AddResponse("GET", "/error.png", testutil.MockResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Internal Server Error",
		})

		ctx := context.Background()
		img, format, err := processor.DownloadImage(ctx, testServer.URL+"/error.png")

		assert.Error(t, err)
		assert.Nil(t, img)
		assert.Empty(t, format)
	})

	t.Run("Invalid image data", func(t *testing.T) {
		testServer.AddResponse("GET", "/invalid.png", testutil.MockResponse{
			StatusCode: http.StatusOK,
			Body:       "not an image",
			Headers: map[string]string{
				"Content-Type": "image/png",
			},
		})

		ctx := context.Background()
		img, format, err := processor.DownloadImage(ctx, testServer.URL+"/invalid.png")

		assert.Error(t, err)
		assert.Nil(t, img)
		assert.Empty(t, format)
	})
}

// TestGetImageDimensions tests dimension extraction
func TestGetImageDimensions(t *testing.T) {
	testServer := testutil.NewTestServer()
	defer testServer.Close()

	// Create a simple PNG image that we know the dimensions of
	mockImageData := testutil.CreateMockImage()

	testServer.AddResponse("GET", "/test-image.png", testutil.MockResponse{
		StatusCode: http.StatusOK,
		Body:       mockImageData,
		Headers: map[string]string{
			"Content-Type": "image/png",
		},
	})

	ctx := context.Background()
	width, height, err := GetImageDimensions(ctx, testServer.URL+"/test-image.png")

	assert.NoError(t, err)
	assert.Equal(t, 1, width) // Our mock image is 1x1
	assert.Equal(t, 1, height)
}

// Benchmark tests
func BenchmarkProcessImageFromURL(b *testing.B) {
	testServer := testutil.NewTestServer()
	defer testServer.Close()

	mockImageData := testutil.CreateMockImage()
	testServer.AddResponse("GET", "/bench-image.png", testutil.MockResponse{
		StatusCode: http.StatusOK,
		Body:       mockImageData,
		Headers: map[string]string{
			"Content-Type": "image/png",
		},
	})

	processor := NewDefaultProcessor()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = processor.ProcessImageFromURL(ctx, testServer.URL+"/bench-image.png")
	}
}

func BenchmarkResizeImage(b *testing.B) {
	processor := NewDefaultProcessor()
	// Create a larger test image for more realistic benchmarking
	src := image.NewRGBA(image.Rect(0, 0, 500, 500))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = processor.ResizeImage(src)
	}
}
