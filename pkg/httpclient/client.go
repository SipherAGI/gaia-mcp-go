package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// HeaderInterceptor is a function that can modify headers before a request is sent
type HeaderInterceptor func(req *http.Request) error

// Client represents our custom HTTP client with enhanced features
type Client struct {
	client             *http.Client        // The underlying HTTP client
	baseURL            string              // Base URL for all requests
	timeout            time.Duration       // Request timeout
	maxRetries         int                 // Maximum number of retry attempts
	retryDelay         time.Duration       // Delay between retries
	debug              bool                // Enable debug logging
	defaultHeaders     map[string]string   // Headers applied to every request
	headerInterceptors []HeaderInterceptor // Functions to modify headers before requests
}

// Config holds configuration options for creating a new HTTP client
type Config struct {
	BaseURL        string            // Base URL for the API
	Timeout        time.Duration     // Request timeout (default: 30 seconds)
	MaxRetries     int               // Maximum retry attempts (default: 3)
	RetryDelay     time.Duration     // Delay between retries (default: 1 second)
	Debug          bool              // Enable debug logging
	DefaultHeaders map[string]string // Headers to add to every request
}

// APIError represents an error returned by the API
type APIError struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

// Error implements the error interface for APIError
func (e *APIError) Error() string {
	return fmt.Sprintf("API Error %d: %s", e.StatusCode, e.Message)
}

// APIResponse represents a generic API response wrapper
type APIResponse[T any] struct {
	Success bool   `json:"success"`
	Data    T      `json:"data"`
	Message string `json:"message"`
	Count   int    `json:"count,omitempty"` // For pagination
	Total   int    `json:"total,omitempty"` // For pagination
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse[T any] struct {
	Data       []T `json:"data"`
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// New creates a new HTTP client with the provided configuration
func New(config Config) *Client {
	// Set default values if not provided
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}
	if config.MaxRetries == 0 {
		config.MaxRetries = 3
	}
	if config.RetryDelay == 0 {
		config.RetryDelay = 1 * time.Second
	}

	// Initialize default headers if nil
	if config.DefaultHeaders == nil {
		config.DefaultHeaders = make(map[string]string)
	}

	// Create the underlying HTTP client with timeout
	httpClient := &http.Client{
		Timeout: config.Timeout,
		// Add transport configuration for better performance
		Transport: &http.Transport{
			MaxIdleConns:        100,              // Maximum idle connections
			MaxIdleConnsPerHost: 10,               // Maximum idle connections per host
			IdleConnTimeout:     90 * time.Second, // How long to keep idle connections
		},
	}

	return &Client{
		client:             httpClient,
		baseURL:            config.BaseURL,
		timeout:            config.Timeout,
		maxRetries:         config.MaxRetries,
		retryDelay:         config.RetryDelay,
		debug:              config.Debug,
		defaultHeaders:     config.DefaultHeaders,
		headerInterceptors: make([]HeaderInterceptor, 0),
	}
}

// Generic HTTP Methods - Type-safe versions

// GetJSON performs a GET request and returns typed data
func GetJSON[T any](c *Client, ctx context.Context, endpoint string, headers map[string]string) (T, error) {
	var result T
	resp, err := c.doRequest(ctx, "GET", endpoint, nil, headers)
	if err != nil {
		return result, err
	}

	err = c.parseJSONResponse(resp, &result)
	return result, err
}

// PostJSON performs a POST request and returns typed data
func PostJSON[T any](c *Client, ctx context.Context, endpoint string, payload interface{}, headers map[string]string) (T, error) {
	var result T
	resp, err := c.doRequest(ctx, "POST", endpoint, payload, headers)
	if err != nil {
		return result, err
	}

	err = c.parseJSONResponse(resp, &result)
	return result, err
}

// PutJSON performs a PUT request and returns typed data
func PutJSON[T any](c *Client, ctx context.Context, endpoint string, payload interface{}, headers map[string]string) (T, error) {
	var result T
	resp, err := c.doRequest(ctx, "PUT", endpoint, payload, headers)
	if err != nil {
		return result, err
	}

	err = c.parseJSONResponse(resp, &result)
	return result, err
}

// DeleteJSON performs a DELETE request and returns typed data
func DeleteJSON[T any](c *Client, ctx context.Context, endpoint string, headers map[string]string) (T, error) {
	var result T
	resp, err := c.doRequest(ctx, "DELETE", endpoint, nil, headers)
	if err != nil {
		return result, err
	}

	err = c.parseJSONResponse(resp, &result)
	return result, err
}

// Generic Methods on Client struct

// GetJSON performs a type-safe GET request
func (c *Client) GetJSON(ctx context.Context, endpoint string, headers map[string]string) *TypedRequestBuilder {
	return &TypedRequestBuilder{
		client:   c,
		ctx:      ctx,
		method:   "GET",
		endpoint: endpoint,
		headers:  headers,
		payload:  nil,
	}
}

// PostJSON performs a type-safe POST request
func (c *Client) PostJSON(ctx context.Context, endpoint string, payload interface{}, headers map[string]string) *TypedRequestBuilder {
	return &TypedRequestBuilder{
		client:   c,
		ctx:      ctx,
		method:   "POST",
		endpoint: endpoint,
		headers:  headers,
		payload:  payload,
	}
}

// PutJSON performs a type-safe PUT request
func (c *Client) PutJSON(ctx context.Context, endpoint string, payload interface{}, headers map[string]string) *TypedRequestBuilder {
	return &TypedRequestBuilder{
		client:   c,
		ctx:      ctx,
		method:   "PUT",
		endpoint: endpoint,
		headers:  headers,
		payload:  payload,
	}
}

// DeleteJSON performs a type-safe DELETE request
func (c *Client) DeleteJSON(ctx context.Context, endpoint string, headers map[string]string) *TypedRequestBuilder {
	return &TypedRequestBuilder{
		client:   c,
		ctx:      ctx,
		method:   "DELETE",
		endpoint: endpoint,
		headers:  headers,
		payload:  nil,
	}
}

// TypedRequestBuilder allows for fluent, type-safe API calls
type TypedRequestBuilder struct {
	client   *Client
	ctx      context.Context
	method   string
	endpoint string
	headers  map[string]string
	payload  interface{}
}

// Into executes the request and unmarshals the response into the specified type
func (rb *TypedRequestBuilder) Into(target interface{}) error {
	resp, err := rb.client.doRequest(rb.ctx, rb.method, rb.endpoint, rb.payload, rb.headers)
	if err != nil {
		return err
	}
	return rb.client.parseJSONResponse(resp, target)
}

// As executes the request and returns the response as the specified type
func As[T any](rb *TypedRequestBuilder) (T, error) {
	var result T
	resp, err := rb.client.doRequest(rb.ctx, rb.method, rb.endpoint, rb.payload, rb.headers)
	if err != nil {
		return result, err
	}

	err = rb.client.parseJSONResponse(resp, &result)
	return result, err
}

// AsResponse executes the request and returns the response wrapped in APIResponse
func AsResponse[T any](rb *TypedRequestBuilder) (APIResponse[T], error) {
	var result APIResponse[T]
	resp, err := rb.client.doRequest(rb.ctx, rb.method, rb.endpoint, rb.payload, rb.headers)
	if err != nil {
		return result, err
	}

	err = rb.client.parseJSONResponse(resp, &result)
	return result, err
}

// AsPaginated executes the request and returns a paginated response
func AsPaginated[T any](rb *TypedRequestBuilder) (PaginatedResponse[T], error) {
	var result PaginatedResponse[T]
	resp, err := rb.client.doRequest(rb.ctx, rb.method, rb.endpoint, rb.payload, rb.headers)
	if err != nil {
		return result, err
	}

	err = rb.client.parseJSONResponse(resp, &result)
	return result, err
}

// SetDefaultHeader sets a default header that will be applied to all requests
func (c *Client) SetDefaultHeader(key, value string) {
	if c.defaultHeaders == nil {
		c.defaultHeaders = make(map[string]string)
	}
	c.defaultHeaders[key] = value
}

// SetDefaultHeaders sets multiple default headers at once
func (c *Client) SetDefaultHeaders(headers map[string]string) {
	if c.defaultHeaders == nil {
		c.defaultHeaders = make(map[string]string)
	}
	for key, value := range headers {
		c.defaultHeaders[key] = value
	}
}

// RemoveDefaultHeader removes a default header
func (c *Client) RemoveDefaultHeader(key string) {
	if c.defaultHeaders != nil {
		delete(c.defaultHeaders, key)
	}
}

// AddHeaderInterceptor adds a function that will be called to modify headers before each request
func (c *Client) AddHeaderInterceptor(interceptor HeaderInterceptor) {
	c.headerInterceptors = append(c.headerInterceptors, interceptor)
}

// SetBearerToken sets the Authorization header with a Bearer token
func (c *Client) SetBearerToken(token string) {
	c.SetDefaultHeader("Authorization", "Bearer "+token)
}

// SetAPIKey sets an API key header (common patterns)
func (c *Client) SetAPIKey(key, value string) {
	c.SetDefaultHeader(key, value)
}

// SetBasicAuth sets basic authentication header
func (c *Client) SetBasicAuth(username, password string) {
	// We'll add the interceptor to handle basic auth
	c.AddHeaderInterceptor(func(req *http.Request) error {
		req.SetBasicAuth(username, password)
		return nil
	})
}

// WithHeaders creates a temporary client with additional headers for a single request chain
func (c *Client) WithHeaders(headers map[string]string) *RequestBuilder {
	return &RequestBuilder{
		client:  c,
		headers: headers,
	}
}

// RequestBuilder allows for fluent API building with custom headers
type RequestBuilder struct {
	client  *Client
	headers map[string]string
}

// GET performs a GET request with the builder's headers
func (rb *RequestBuilder) GET(ctx context.Context, endpoint string) (*http.Response, error) {
	return rb.client.doRequest(ctx, "GET", endpoint, nil, rb.headers)
}

// POST performs a POST request with the builder's headers
func (rb *RequestBuilder) POST(ctx context.Context, endpoint string, payload interface{}) (*http.Response, error) {
	return rb.client.doRequest(ctx, "POST", endpoint, payload, rb.headers)
}

// PUT performs a PUT request with the builder's headers
func (rb *RequestBuilder) PUT(ctx context.Context, endpoint string, payload interface{}) (*http.Response, error) {
	return rb.client.doRequest(ctx, "PUT", endpoint, payload, rb.headers)
}

// DELETE performs a DELETE request with the builder's headers
func (rb *RequestBuilder) DELETE(ctx context.Context, endpoint string) (*http.Response, error) {
	return rb.client.doRequest(ctx, "DELETE", endpoint, nil, rb.headers)
}

// GetJSON performs a type-safe GET request with builder headers
func (rb *RequestBuilder) GetJSON(ctx context.Context, endpoint string) *TypedRequestBuilder {
	return &TypedRequestBuilder{
		client:   rb.client,
		ctx:      ctx,
		method:   "GET",
		endpoint: endpoint,
		headers:  rb.headers,
		payload:  nil,
	}
}

// PostJSON performs a type-safe POST request with builder headers
func (rb *RequestBuilder) PostJSON(ctx context.Context, endpoint string, payload interface{}) *TypedRequestBuilder {
	return &TypedRequestBuilder{
		client:   rb.client,
		ctx:      ctx,
		method:   "POST",
		endpoint: endpoint,
		headers:  rb.headers,
		payload:  payload,
	}
}

// GET performs a GET request to the specified endpoint
func (c *Client) GET(ctx context.Context, endpoint string, headers map[string]string) (*http.Response, error) {
	return c.doRequest(ctx, "GET", endpoint, nil, headers)
}

// POST performs a POST request with JSON payload
func (c *Client) POST(ctx context.Context, endpoint string, payload interface{}, headers map[string]string) (*http.Response, error) {
	return c.doRequest(ctx, "POST", endpoint, payload, headers)
}

// PUT performs a PUT request with JSON payload
func (c *Client) PUT(ctx context.Context, endpoint string, payload interface{}, headers map[string]string) (*http.Response, error) {
	return c.doRequest(ctx, "PUT", endpoint, payload, headers)
}

// DELETE performs a DELETE request
func (c *Client) DELETE(ctx context.Context, endpoint string, headers map[string]string) (*http.Response, error) {
	return c.doRequest(ctx, "DELETE", endpoint, nil, headers)
}

// doRequest is the core method that handles all HTTP requests with retry logic
func (c *Client) doRequest(ctx context.Context, method, endpoint string, payload interface{}, headers map[string]string) (*http.Response, error) {
	// Build the full URL
	url := c.baseURL + endpoint

	// Prepare the request body if payload is provided
	var body io.Reader
	if payload != nil {
		jsonData, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal payload: %w", err)
		}
		body = bytes.NewBuffer(jsonData)
	}

	// Retry logic with exponential backoff
	var lastErr error
	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		// Create a new request for each attempt
		req, err := http.NewRequestWithContext(ctx, method, url, body)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		// Apply headers in order: defaults -> custom -> interceptors
		if err := c.applyHeaders(req, headers); err != nil {
			return nil, fmt.Errorf("failed to apply headers: %w", err)
		}

		// Log the request if debug is enabled
		if c.debug {
			c.logRequest(req, method, url, attempt)
		}

		// Perform the request
		resp, err := c.client.Do(req)
		if err != nil {
			lastErr = err
			if attempt < c.maxRetries {
				// Wait before retrying
				time.Sleep(c.retryDelay * time.Duration(attempt+1)) // Exponential backoff
				continue
			}
			return nil, fmt.Errorf("request failed after %d attempts: %w", c.maxRetries+1, err)
		}

		// Check if we should retry based on status code
		if c.shouldRetry(resp.StatusCode) && attempt < c.maxRetries {
			resp.Body.Close() // Important: close the response body
			if c.debug {
				log.Printf("Retrying request due to status code %d", resp.StatusCode)
			}
			time.Sleep(c.retryDelay * time.Duration(attempt+1))
			continue
		}

		// Log successful response if debug is enabled
		if c.debug {
			log.Printf("Request completed with status %d", resp.StatusCode)
		}

		return resp, nil
	}

	return nil, lastErr
}

// applyHeaders applies headers in the correct order: defaults -> custom -> interceptors
func (c *Client) applyHeaders(req *http.Request, customHeaders map[string]string) error {
	// Step 1: Set standard headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "gaia-mcp-go/1.0")

	// Step 2: Apply default headers (can override standard headers)
	for key, value := range c.defaultHeaders {
		req.Header.Set(key, value)
	}

	// Step 3: Apply custom headers for this request (can override defaults)
	for key, value := range customHeaders {
		req.Header.Set(key, value)
	}

	// Step 4: Apply header interceptors (can override everything)
	for _, interceptor := range c.headerInterceptors {
		if err := interceptor(req); err != nil {
			return fmt.Errorf("header interceptor failed: %w", err)
		}
	}

	return nil
}

// logRequest logs request details when debug is enabled
func (c *Client) logRequest(req *http.Request, method, url string, attempt int) {
	log.Printf("Making %s request to %s (attempt %d/%d)", method, url, attempt+1, c.maxRetries+1)

	// Log important headers (but hide sensitive ones)
	for key, values := range req.Header {
		if c.isSensitiveHeader(key) {
			log.Printf("  %s: [REDACTED]", key)
		} else {
			log.Printf("  %s: %s", key, strings.Join(values, ", "))
		}
	}
}

// isSensitiveHeader checks if a header contains sensitive information
func (c *Client) isSensitiveHeader(key string) bool {
	sensitiveHeaders := []string{
		"authorization",
		"cookie",
		"x-api-key",
		"x-auth-token",
	}

	keyLower := strings.ToLower(key)
	for _, sensitive := range sensitiveHeaders {
		if keyLower == sensitive || strings.Contains(keyLower, sensitive) {
			return true
		}
	}
	return false
}

// shouldRetry determines if a request should be retried based on the status code
func (c *Client) shouldRetry(statusCode int) bool {
	// Retry on server errors (5xx) and specific client errors
	switch statusCode {
	case http.StatusTooManyRequests, // 429
		http.StatusInternalServerError, // 500
		http.StatusBadGateway,          // 502
		http.StatusServiceUnavailable,  // 503
		http.StatusGatewayTimeout:      // 504
		return true
	}
	return false
}

// parseJSONResponse is the internal method used by generic functions
func (c *Client) parseJSONResponse(resp *http.Response, target interface{}) error {
	defer resp.Body.Close() // Always close the response body

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Check if the response indicates an error
	if resp.StatusCode >= 400 {
		// Try to parse as API error first
		var apiErr APIError
		if json.Unmarshal(body, &apiErr) == nil && apiErr.Message != "" {
			apiErr.StatusCode = resp.StatusCode
			return &apiErr
		}

		// If we can't parse as API error, create a generic error
		return &APIError{
			StatusCode: resp.StatusCode,
			Message:    string(body),
		}
	}

	// Parse successful response
	if err := json.Unmarshal(body, target); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return nil
}

// ParseJSONResponse is a helper method to parse JSON responses (legacy method)
func (c *Client) ParseJSONResponse(resp *http.Response, target interface{}) error {
	return c.parseJSONResponse(resp, target)
}

// Close closes the HTTP client and cleans up resources
func (c *Client) Close() {
	// Close idle connections
	c.client.CloseIdleConnections()
}
