# HTTP Client Package

A powerful, type-safe HTTP client for Go that provides enhanced features like automatic retries, request/response logging, authentication helpers, and fluent API design.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Configuration](#configuration)
- [Basic Usage](#basic-usage)
- [Type-Safe Requests](#type-safe-requests)
- [Authentication](#authentication)
- [Error Handling](#error-handling)
- [Advanced Features](#advanced-features)
- [Best Practices](#best-practices)

## Features

- **Type-Safe API**: Generic methods for type-safe JSON requests and responses
- **Automatic Retries**: Configurable retry logic with exponential backoff
- **Authentication Support**: Built-in support for Bearer tokens, API keys, and Basic auth
- **Header Management**: Default headers, custom headers, and header interceptors
- **Request Builders**: Fluent API for building complex requests
- **Debug Logging**: Detailed request/response logging (with sensitive data protection)
- **Error Handling**: Custom error types with detailed API error information
- **Pagination Support**: Built-in support for paginated responses
- **Context Support**: Full context support for request cancellation and timeouts

## Installation

```bash
go get your-module/pkg/httpclient
```

## Quick Start

Here's a simple example to get you started:

```go
package main

import (
    "context"
    "fmt"
    "time"

    "your-module/pkg/httpclient"
)

// Define your data structures
type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

func main() {
    // Create a new HTTP client
    client := httpclient.New(httpclient.Config{
        BaseURL:    "https://api.example.com",
        Timeout:    30 * time.Second,
        MaxRetries: 3,
        Debug:      true,
    })
    defer client.Close()

    ctx := context.Background()

    // Make a type-safe GET request
    user, err := httpclient.As[User](
        client.GetJSON(ctx, "/users/1", nil),
    )
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    fmt.Printf("User: %+v\n", user)
}
```

## Configuration

The `Config` struct allows you to customize the HTTP client behavior:

```go
config := httpclient.Config{
    BaseURL:    "https://api.example.com",  // Base URL for all requests
    Timeout:    30 * time.Second,           // Request timeout (default: 30s)
    MaxRetries: 3,                          // Max retry attempts (default: 3)
    RetryDelay: 1 * time.Second,            // Delay between retries (default: 1s)
    Debug:      true,                       // Enable debug logging
    DefaultHeaders: map[string]string{      // Headers for all requests
        "X-App-Version": "1.0.0",
    },
}

client := httpclient.New(config)
```

## Basic Usage

### Standard HTTP Methods

```go
// GET request
resp, err := client.GET(ctx, "/users", nil)
if err != nil {
    // handle error
}
defer resp.Body.Close()

// POST request with payload
payload := map[string]interface{}{
    "name":  "John Doe",
    "email": "john@example.com",
}

resp, err = client.POST(ctx, "/users", payload, nil)
if err != nil {
    // handle error
}
defer resp.Body.Close()

// PUT request
resp, err = client.PUT(ctx, "/users/1", payload, nil)

// DELETE request
resp, err = client.DELETE(ctx, "/users/1", nil)
```

## Type-Safe Requests

The package provides type-safe methods that automatically handle JSON marshaling/unmarshaling:

### Using Generic Functions

```go
// GET request returning a single user
user, err := httpclient.GetJSON[User](client, ctx, "/users/1", nil)
if err != nil {
    // handle error
}

// POST request with payload, returning created user
newUser := User{Name: "Jane Doe", Email: "jane@example.com"}
createdUser, err := httpclient.PostJSON[User](client, ctx, "/users", newUser, nil)

// GET request returning a list of users
users, err := httpclient.GetJSON[[]User](client, ctx, "/users", nil)
```

### Using Fluent API

```go
// Type-safe GET request using fluent API
user, err := httpclient.As[User](
    client.GetJSON(ctx, "/users/1", nil),
)

// POST request with fluent API
createdUser, err := httpclient.As[User](
    client.PostJSON(ctx, "/users", newUser, nil),
)

// Using Into() method to unmarshal into existing variable
var user User
err := client.GetJSON(ctx, "/users/1", nil).Into(&user)
```

### Working with API Responses

Many APIs wrap their responses in a standard format. The package supports this pattern:

```go
// For APIs that return wrapped responses
type APIResponse[T any] struct {
    Success bool   `json:"success"`
    Data    T      `json:"data"`
    Message string `json:"message"`
    Count   int    `json:"count,omitempty"`
    Total   int    `json:"total,omitempty"`
}

// Get wrapped response
userResponse, err := httpclient.AsResponse[User](
    client.GetJSON(ctx, "/users/1", nil),
)
if err != nil {
    // handle error
}

if userResponse.Success {
    user := userResponse.Data
    fmt.Printf("User: %+v\n", user)
}
```

### Paginated Responses

For APIs that return paginated data:

```go
// Get paginated response
usersPage, err := httpclient.AsPaginated[User](
    client.GetJSON(ctx, "/users?page=1&per_page=10", nil),
)
if err != nil {
    // handle error
}

fmt.Printf("Page %d of %d\n", usersPage.Page, usersPage.TotalPages)
fmt.Printf("Users: %+v\n", usersPage.Data)
```

## Authentication

### Bearer Token Authentication

```go
client := httpclient.New(config)

// Set Bearer token for all requests
client.SetBearerToken("your-jwt-token-here")

// Now all requests will include the Authorization header
user, err := httpclient.As[User](
    client.GetJSON(ctx, "/protected/users/1", nil),
)
```

### API Key Authentication

```go
// Set API key header
client.SetAPIKey("X-API-Key", "your-api-key-here")

// Or set multiple API-related headers
client.SetDefaultHeaders(map[string]string{
    "X-API-Key":     "your-api-key",
    "X-API-Version": "v1",
})
```

### Basic Authentication

```go
// Set basic authentication
client.SetBasicAuth("username", "password")

// This will add the Authorization header to all requests
```

### Custom Authentication with Header Interceptors

For more complex authentication scenarios:

```go
// Add a custom header interceptor
client.AddHeaderInterceptor(func(req *http.Request) error {
    // Add custom authentication logic here
    token, err := getTokenFromSomeSource()
    if err != nil {
        return err
    }

    req.Header.Set("X-Custom-Auth", token)
    return nil
})
```

## Error Handling

The package provides detailed error information through the `APIError` type:

```go
user, err := httpclient.As[User](
    client.GetJSON(ctx, "/users/999", nil),
)
if err != nil {
    // Check if it's an API error
    var apiErr *httpclient.APIError
    if errors.As(err, &apiErr) {
        fmt.Printf("API Error %d: %s\n", apiErr.StatusCode, apiErr.Message)
        if apiErr.Details != "" {
            fmt.Printf("Details: %s\n", apiErr.Details)
        }
    } else {
        // Handle other types of errors (network, parsing, etc.)
        fmt.Printf("Other error: %v\n", err)
    }
}
```

## Advanced Features

### Request Builder with Custom Headers

```go
// Create a request builder with custom headers
builder := client.WithHeaders(map[string]string{
    "X-Request-ID": "unique-request-id",
    "X-User-Agent": "MyApp/1.0",
})

// Use the builder for requests
user, err := httpclient.As[User](
    builder.GetJSON(ctx, "/users/1"),
)

// Or make standard HTTP requests
resp, err := builder.GET(ctx, "/users/1")
```

### Managing Default Headers

```go
// Set individual default headers
client.SetDefaultHeader("X-App-Version", "1.2.0")

// Set multiple default headers
client.SetDefaultHeaders(map[string]string{
    "X-Client-Type":    "mobile",
    "X-Client-Version": "2.1.0",
})

// Remove a default header
client.RemoveDefaultHeader("X-App-Version")
```

### Debug Logging

When debug mode is enabled, the client logs detailed request/response information:

```go
client := httpclient.New(httpclient.Config{
    BaseURL: "https://api.example.com",
    Debug:   true, // Enable debug logging
})

// This will log:
// - Request method, URL, and attempt number
// - Request headers (sensitive headers are redacted)
// - Response status code
// - Retry attempts if they occur
```

### Custom Request Processing

You can add header interceptors to modify requests before they're sent:

```go
// Add timestamp to all requests
client.AddHeaderInterceptor(func(req *http.Request) error {
    req.Header.Set("X-Timestamp", time.Now().Format(time.RFC3339))
    return nil
})

// Add request signing
client.AddHeaderInterceptor(func(req *http.Request) error {
    signature, err := signRequest(req)
    if err != nil {
        return err
    }
    req.Header.Set("X-Signature", signature)
    return nil
})
```

## Best Practices

### 1. Always Use Context

```go
// Always pass a context, preferably with timeout
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

user, err := httpclient.As[User](
    client.GetJSON(ctx, "/users/1", nil),
)
```

### 2. Handle Errors Appropriately

```go
user, err := httpclient.As[User](
    client.GetJSON(ctx, "/users/1", nil),
)
if err != nil {
    var apiErr *httpclient.APIError
    if errors.As(err, &apiErr) {
        switch apiErr.StatusCode {
        case 404:
            // User not found - this might be expected
            return nil, ErrUserNotFound
        case 401:
            // Authentication error - might need to refresh token
            return nil, ErrUnauthorized
        default:
            // Other API errors
            return nil, fmt.Errorf("API error: %w", apiErr)
        }
    }
    // Network or parsing errors
    return nil, fmt.Errorf("request failed: %w", err)
}
```

### 3. Reuse Client Instances

```go
// Create one client instance and reuse it
var apiClient *httpclient.Client

func init() {
    apiClient = httpclient.New(httpclient.Config{
        BaseURL:    os.Getenv("API_BASE_URL"),
        Timeout:    30 * time.Second,
        MaxRetries: 3,
    })
}

// Use the same client instance across your application
func GetUser(ctx context.Context, userID int) (*User, error) {
    return httpclient.As[User](
        apiClient.GetJSON(ctx, fmt.Sprintf("/users/%d", userID), nil),
    )
}
```

### 4. Clean Up Resources

```go
func main() {
    client := httpclient.New(config)

    // Always close the client when done
    defer client.Close()

    // Your application logic here
}
```

### 5. Use Type Safety

```go
// Define clear types for your API responses
type CreateUserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

type CreateUserResponse struct {
    ID        int       `json:"id"`
    Name      string    `json:"name"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
}

// Use these types in your API calls
func CreateUser(ctx context.Context, req CreateUserRequest) (*CreateUserResponse, error) {
    response, err := httpclient.As[CreateUserResponse](
        apiClient.PostJSON(ctx, "/users", req, nil),
    )
    if err != nil {
        return nil, err
    }
    return &response, nil
}
```

This HTTP client package provides a robust, type-safe foundation for making HTTP requests in your Go applications. It handles common concerns like retries, authentication, and error handling while maintaining a clean, intuitive API.
