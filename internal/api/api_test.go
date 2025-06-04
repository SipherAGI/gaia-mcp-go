package api

import (
	"context"
	"gaia-mcp-go/internal/testutil"
	"gaia-mcp-go/pkg/shared"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewGaiaApi(t *testing.T) {
	tests := []struct {
		name   string
		config GaiaApiConfig
	}{
		{
			name: "Valid configuration",
			config: GaiaApiConfig{
				BaseUrl: "https://api.gaia.com",
				ApiKey:  "test-api-key",
			},
		},
		{
			name: "Different base URL",
			config: GaiaApiConfig{
				BaseUrl: "https://staging.gaia.com",
				ApiKey:  "staging-key",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewGaiaApi(tt.config)
			assert.NotNil(t, client)

			// Verify it implements the interface
			var _ GaiaApi = client
		})
	}
}

func TestGaiaApi_CreateStyle(t *testing.T) {
	tests := []struct {
		name          string
		imageUrls     []string
		styleName     string
		description   *string
		mockResponse  testutil.MockResponse
		expectedStyle SdStyle
		expectedError string
	}{
		{
			name:        "Successful style creation",
			imageUrls:   []string{"https://example.com/image1.jpg", "https://example.com/image2.jpg"},
			styleName:   "Test Style",
			description: func() *string { s := "A test style"; return &s }(),
			mockResponse: testutil.MockResponse{
				StatusCode: 200,
				Body: SdStyle{
					Id:          "style-123",
					Name:        "Test Style",
					Description: "A test style",
					CreatedAt:   "2023-01-01T00:00:00Z",
				},
			},
			expectedStyle: SdStyle{
				Id:          "style-123",
				Name:        "Test Style",
				Description: "A test style",
			},
		},
		{
			name:        "Style creation without description",
			imageUrls:   []string{"https://example.com/image1.jpg"},
			styleName:   "Simple Style",
			description: nil,
			mockResponse: testutil.MockResponse{
				StatusCode: 200,
				Body: SdStyle{
					Id:        "style-456",
					Name:      "Simple Style",
					CreatedAt: "2023-01-01T00:00:00Z",
				},
			},
			expectedStyle: SdStyle{
				Id:   "style-456",
				Name: "Simple Style",
			},
		},
		{
			name:      "API error response",
			imageUrls: []string{"https://example.com/image1.jpg"},
			styleName: "Failed Style",
			mockResponse: testutil.MockResponse{
				StatusCode: 400,
				Body: map[string]interface{}{
					"error": "Invalid image URL",
					"code":  "INVALID_URL",
				},
			},
			expectedError: "Invalid image URL",
		},
		{
			name:      "Network timeout",
			imageUrls: []string{"https://example.com/image1.jpg"},
			styleName: "Timeout Style",
			mockResponse: testutil.MockResponse{
				StatusCode: 0, // Indicates network error
				Delay:      2 * time.Second,
			},
			expectedError: "context deadline exceeded",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test server
			server := testutil.NewTestServer()
			defer server.Close()

			// Configure mock response
			server.AddResponse("POST", "/api/sd-styles", tt.mockResponse)

			// Create API client
			client := NewGaiaApi(GaiaApiConfig{
				BaseUrl: server.URL,
				ApiKey:  "test-key",
			})

			// Create context with timeout for timeout tests
			ctx := context.Background()
			if strings.Contains(tt.name, "timeout") {
				var cancel context.CancelFunc
				ctx, cancel = context.WithTimeout(ctx, 1*time.Second)
				defer cancel()
			}

			// Execute test
			style, err := client.CreateStyle(ctx, tt.imageUrls, tt.styleName, tt.description)

			// Verify results
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Empty(t, style.Id)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStyle.Id, style.Id)
				assert.Equal(t, tt.expectedStyle.Name, style.Name)
				if tt.expectedStyle.Description != "" {
					assert.Equal(t, tt.expectedStyle.Description, style.Description)
				}
			}
		})
	}
}

func TestGaiaApi_GenerateImages(t *testing.T) {
	tests := []struct {
		name             string
		request          GenerateImagesRequest
		mockResponse     testutil.MockResponse
		expectedResponse ImageGeneratedResponse
		expectedError    string
	}{
		{
			name: "Successful image generation",
			request: GenerateImagesRequest{
				RecipeId: shared.RecipeIdImageGeneratorSimple,
				Params: map[string]interface{}{
					"prompt":      "A beautiful sunset",
					"styleId":     "style-123",
					"aspectRatio": shared.AspectRatio16_9,
				},
			},
			mockResponse: testutil.MockResponse{
				StatusCode: 200,
				Body: ImageGeneratedResponse{
					Success: true,
					Images:  []string{"image-url-1", "image-url-2"},
				},
			},
			expectedResponse: ImageGeneratedResponse{
				Success: true,
				Images:  []string{"image-url-1", "image-url-2"},
			},
		},
		{
			name: "Invalid request parameters",
			request: GenerateImagesRequest{
				RecipeId: shared.RecipeIdImageGeneratorSimple,
				Params: map[string]interface{}{
					"prompt": "", // Empty prompt should cause error
				},
			},
			mockResponse: testutil.MockResponse{
				StatusCode: 400,
				Body: map[string]interface{}{
					"error": "Prompt cannot be empty",
					"code":  "INVALID_PROMPT",
				},
			},
			expectedError: "Prompt cannot be empty",
		},
		{
			name: "API rate limit exceeded",
			request: GenerateImagesRequest{
				RecipeId: shared.RecipeIdImageGeneratorSimple,
				Params: map[string]interface{}{
					"prompt": "Test prompt",
				},
			},
			mockResponse: testutil.MockResponse{
				StatusCode: 429,
				Body: map[string]interface{}{
					"error": "Rate limit exceeded",
					"code":  "RATE_LIMIT",
				},
			},
			expectedError: "Rate limit exceeded",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test server
			server := testutil.NewTestServer()
			defer server.Close()

			// Configure mock response
			server.AddResponse("POST", "/api/recipe/agi-tasks/create-task", tt.mockResponse)

			// Create API client
			client := NewGaiaApi(GaiaApiConfig{
				BaseUrl: server.URL,
				ApiKey:  "test-key",
			})

			// Execute test
			ctx := context.Background()
			response, err := client.GenerateImages(ctx, tt.request)

			// Verify results
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.False(t, response.Success)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResponse.Success, response.Success)
				assert.Equal(t, tt.expectedResponse.Images, response.Images)
			}
		})
	}
}

func TestGaiaApi_UploadImages(t *testing.T) {
	tests := []struct {
		name               string
		imageUrls          []string
		associatedResource shared.FileAssociatedResource
		expectedError      string
	}{
		{
			name:               "Image download failure",
			imageUrls:          []string{"https://example.com/nonexistent.jpg"},
			associatedResource: shared.FileAssociatedResourceStyle,
			expectedError:      "failed to upload some files",
		},
		{
			name:               "Empty image URLs",
			imageUrls:          []string{},
			associatedResource: shared.FileAssociatedResourceStyle,
			expectedError:      "", // Should return empty slice, no error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test server
			server := testutil.NewTestServer()
			defer server.Close()

			// Create API client
			client := NewGaiaApi(GaiaApiConfig{
				BaseUrl: server.URL,
				ApiKey:  "test-key",
			})

			// Execute test
			ctx := context.Background()
			files, err := client.UploadImages(ctx, tt.imageUrls, tt.associatedResource)

			// Verify results
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Empty(t, files) // Empty URLs should return empty files
			}
		})
	}
}

// Benchmark tests for performance monitoring
func BenchmarkGaiaApi_CreateStyle(b *testing.B) {
	server := testutil.NewTestServer()
	defer server.Close()

	server.AddResponse("POST", "/api/sd-styles", testutil.MockResponse{
		StatusCode: 200,
		Body: SdStyle{
			Id:   "benchmark-style",
			Name: "Benchmark Style",
		},
	})

	client := NewGaiaApi(GaiaApiConfig{
		BaseUrl: server.URL,
		ApiKey:  "benchmark-key",
	})

	imageUrls := []string{"https://example.com/image1.jpg", "https://example.com/image2.jpg"}
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := client.CreateStyle(ctx, imageUrls, "Benchmark Style", nil)
		if err != nil {
			b.Fatalf("Benchmark failed: %v", err)
		}
	}
}

func BenchmarkGaiaApi_GenerateImages(b *testing.B) {
	server := testutil.NewTestServer()
	defer server.Close()

	server.AddResponse("POST", "/api/recipe/agi-tasks/create-task", testutil.MockResponse{
		StatusCode: 200,
		Body: ImageGeneratedResponse{
			Success: true,
			Images:  []string{"benchmark-image"},
		},
	})

	client := NewGaiaApi(GaiaApiConfig{
		BaseUrl: server.URL,
		ApiKey:  "benchmark-key",
	})

	request := GenerateImagesRequest{
		RecipeId: shared.RecipeIdImageGeneratorSimple,
		Params: map[string]interface{}{
			"prompt":      "Benchmark prompt",
			"styleId":     "style-123",
			"aspectRatio": shared.AspectRatio16_9,
		},
	}
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := client.GenerateImages(ctx, request)
		if err != nil {
			b.Fatalf("Benchmark failed: %v", err)
		}
	}
}

// Test helper functions
func TestGaiaApiConfig_Validation(t *testing.T) {
	tests := []struct {
		name   string
		config GaiaApiConfig
		valid  bool
	}{
		{
			name: "Valid config",
			config: GaiaApiConfig{
				BaseUrl: "https://api.gaia.com",
				ApiKey:  "valid-key",
			},
			valid: true,
		},
		{
			name: "Empty base URL",
			config: GaiaApiConfig{
				BaseUrl: "",
				ApiKey:  "valid-key",
			},
			valid: false,
		},
		{
			name: "Empty API key",
			config: GaiaApiConfig{
				BaseUrl: "https://api.gaia.com",
				ApiKey:  "",
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewGaiaApi(tt.config)

			// Basic validation - client should always be created
			// but functionality depends on valid config
			assert.NotNil(t, client)

			if !tt.valid {
				// Test that invalid configs lead to errors in actual usage
				ctx := context.Background()
				_, err := client.CreateStyle(ctx, []string{"https://example.com/test.jpg"}, "Test", nil)
				assert.Error(t, err)
			}
		})
	}
}
