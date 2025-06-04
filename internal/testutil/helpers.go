package testutil

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestServer represents a test HTTP server for mocking external APIs
type TestServer struct {
	*httptest.Server
	responses map[string]MockResponse
}

// MockResponse represents a mock HTTP response
type MockResponse struct {
	StatusCode int
	Body       interface{}
	Headers    map[string]string
	Delay      time.Duration
}

// NewTestServer creates a new test server with predefined responses
func NewTestServer() *TestServer {
	ts := &TestServer{
		responses: make(map[string]MockResponse),
	}

	// Create the actual HTTP server
	ts.Server = httptest.NewServer(http.HandlerFunc(ts.handler))
	return ts
}

// AddResponse adds a mock response for a specific endpoint
func (ts *TestServer) AddResponse(method, path string, response MockResponse) {
	key := fmt.Sprintf("%s:%s", method, path)
	ts.responses[key] = response
}

// handler handles incoming requests and returns mock responses
func (ts *TestServer) handler(w http.ResponseWriter, r *http.Request) {
	key := fmt.Sprintf("%s:%s", r.Method, r.URL.Path)

	// Check if we have a mock response for this endpoint
	response, exists := ts.responses[key]
	if !exists {
		// Default response for unmocked endpoints
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `{"error": "Mock not found for %s"}`, key)
		return
	}

	// Add artificial delay if specified
	if response.Delay > 0 {
		time.Sleep(response.Delay)
	}

	// Set custom headers
	for key, value := range response.Headers {
		w.Header().Set(key, value)
	}

	// Set status code
	w.WriteHeader(response.StatusCode)

	// Write response body
	switch body := response.Body.(type) {
	case string:
		w.Write([]byte(body))
	case []byte:
		w.Write(body)
	default:
		// JSON encode other types
		jsonData, err := json.Marshal(body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"error": "Failed to marshal response: %s"}`, err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	}
}

// AssertJSONRequest validates that a request contains the expected JSON data
func AssertJSONRequest(t *testing.T, r *http.Request, expected interface{}) {
	t.Helper()

	body, err := io.ReadAll(r.Body)
	require.NoError(t, err, "Failed to read request body")

	// Parse the actual JSON
	var actual interface{}
	err = json.Unmarshal(body, &actual)
	require.NoError(t, err, "Failed to parse request JSON")

	// Compare with expected
	expectedJSON, err := json.Marshal(expected)
	require.NoError(t, err, "Failed to marshal expected data")

	var expectedParsed interface{}
	err = json.Unmarshal(expectedJSON, &expectedParsed)
	require.NoError(t, err, "Failed to parse expected JSON")

	assert.Equal(t, expectedParsed, actual, "Request JSON does not match expected")
}

// CreateTestContext creates a context with timeout for testing
func CreateTestContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 10*time.Second)
}

// LoadTestData loads test data from a file in the testdata directory
func LoadTestData(t *testing.T, filename string) []byte {
	t.Helper()

	// Find the testdata directory (walk up from current directory)
	path, err := findTestDataPath(filename)
	require.NoError(t, err, "Failed to find testdata file: %s", filename)

	data, err := os.ReadFile(path)
	require.NoError(t, err, "Failed to read testdata file: %s", filename)

	return data
}

// findTestDataPath searches for the testdata directory and file
func findTestDataPath(filename string) (string, error) {
	// Start from current directory and walk up
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		testdataPath := filepath.Join(dir, "testdata", filename)
		if _, err := os.Stat(testdataPath); err == nil {
			return testdataPath, nil
		}

		// Move up one directory
		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached root directory
			break
		}
		dir = parent
	}

	return "", fmt.Errorf("testdata file not found: %s", filename)
}

// AssertNoError is a helper that fails the test if err is not nil
func AssertNoError(t *testing.T, err error, msgAndArgs ...interface{}) {
	t.Helper()
	if err != nil {
		if len(msgAndArgs) > 0 {
			t.Fatalf("Expected no error, but got: %v. Message: %v", err, msgAndArgs[0])
		} else {
			t.Fatalf("Expected no error, but got: %v", err)
		}
	}
}

// AssertError is a helper that fails the test if err is nil
func AssertError(t *testing.T, err error, msgAndArgs ...interface{}) {
	t.Helper()
	if err == nil {
		if len(msgAndArgs) > 0 {
			t.Fatalf("Expected an error, but got nil. Message: %v", msgAndArgs[0])
		} else {
			t.Fatal("Expected an error, but got nil")
		}
	}
}

// AssertContains checks if a string contains a substring
func AssertContains(t *testing.T, haystack, needle string, msgAndArgs ...interface{}) {
	t.Helper()
	if !strings.Contains(haystack, needle) {
		if len(msgAndArgs) > 0 {
			t.Fatalf("Expected '%s' to contain '%s'. Message: %v", haystack, needle, msgAndArgs[0])
		} else {
			t.Fatalf("Expected '%s' to contain '%s'", haystack, needle)
		}
	}
}

// CreateMockImage creates a simple mock image data for testing
func CreateMockImage() []byte {
	// This is a minimal valid PNG image (1x1 pixel, transparent)
	return []byte{
		0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x00, 0x00, 0x00, 0x0D,
		0x49, 0x48, 0x44, 0x52, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x06, 0x00, 0x00, 0x00, 0x1F, 0x15, 0xC4, 0x89, 0x00, 0x00, 0x00,
		0x0A, 0x49, 0x44, 0x41, 0x54, 0x78, 0x9C, 0x63, 0x00, 0x01, 0x00, 0x00,
		0x05, 0x00, 0x01, 0x0D, 0x0A, 0x2D, 0xB4, 0x00, 0x00, 0x00, 0x00, 0x49,
		0x45, 0x4E, 0x44, 0xAE, 0x42, 0x60, 0x82,
	}
}

// TestConfig represents common test configuration
type TestConfig struct {
	APIKey  string
	BaseURL string
	Timeout time.Duration
}

// DefaultTestConfig returns a default test configuration
func DefaultTestConfig() TestConfig {
	return TestConfig{
		APIKey:  "test-api-key-12345",
		BaseURL: "https://api.test.com",
		Timeout: 5 * time.Second,
	}
}
