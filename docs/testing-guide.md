# Testing Guide for gaia-mcp-go

This guide provides comprehensive information about testing in the gaia-mcp-go project, including how to run tests, write new tests, and understand our testing strategy.

## Table of Contents

1. [Testing Philosophy](#testing-philosophy)
2. [Test Structure](#test-structure)
3. [Running Tests](#running-tests)
4. [Current Test Coverage](#current-test-coverage)
5. [Writing Tests](#writing-tests)
6. [Test Utilities](#test-utilities)
7. [Mocking and Test Data](#mocking-and-test-data)
8. [Best Practices](#best-practices)
9. [Troubleshooting](#troubleshooting)
10. [Learning Resources](#learning-resources)

## Testing Philosophy

Our testing strategy focuses on:

- **Reliability**: Tests should be deterministic and not flaky
- **Maintainability**: Tests should be easy to understand and modify
- **Coverage**: Critical paths and edge cases should be tested
- **Performance**: Include benchmark tests for performance-critical code
- **Documentation**: Tests serve as living documentation of expected behavior

## Test Structure

### File Organization

Tests are organized following Go conventions:

```
gaia-mcp-go/
├── internal/
│   ├── api/
│   │   ├── api.go
│   │   └── api_test.go          # API client tests
│   └── testutil/
│       └── helpers.go           # Test utilities and helpers
├── pkg/
│   ├── imageutil/
│   │   ├── processor.go
│   │   └── processor_test.go    # Image processing tests
│   └── shared/
│       ├── type.go
│       └── type_test.go         # Shared types tests
├── version/
│   ├── version.go
│   └── version_test.go          # Version parsing tests
└── testdata/
    └── mock_responses.json      # Mock API responses
```

### Naming Conventions

- Test files: `*_test.go`
- Test functions: `TestFunctionName`
- Benchmark functions: `BenchmarkFunctionName`
- Test helpers: `testHelperName` (lowercase, not exported)

## Running Tests

### Basic Commands

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test ./... -v

# Run tests for a specific package
go test ./pkg/imageutil

# Run a specific test
go test ./pkg/shared -run TestRecipeTaskStatus

# Run tests with coverage
go test ./... -cover

# Generate detailed coverage report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

### Using Makefile

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run benchmarks
make benchmark

# Run tests in watch mode (if available)
make test-watch
```

### Test Flags

```bash
# Run tests multiple times to check for flakiness
go test ./... -count=10

# Run tests with race detection
go test ./... -race

# Run only short tests (skip long-running tests)
go test ./... -short

# Set timeout for tests
go test ./... -timeout=30s
```

## Current Test Coverage

### Tested Packages

#### ✅ `internal/api` - API Client Tests

- **Coverage**: Core API functionality
- **Test Count**: 12 test functions
- **Features Tested**:
  - Client initialization and configuration
  - Style creation (with/without description)
  - Image generation requests
  - File upload workflows
  - Error handling and timeouts
  - HTTP mocking and response validation
  - Configuration validation

#### ✅ `pkg/imageutil` - Image Processing Tests

- **Coverage**: Complete image processing pipeline
- **Test Count**: 8 test functions
- **Features Tested**:
  - Image downloading from URLs
  - Image resizing with aspect ratio preservation
  - Base64 encoding (data URL and pure formats)
  - Error handling for invalid images/URLs
  - HTTP timeout scenarios
  - Image format detection and conversion

#### ✅ `pkg/shared` - Shared Types Tests

- **Coverage**: All constants and type mappings
- **Test Count**: 14 test functions
- **Features Tested**:
  - All enum constants (RecipeTaskStatus, RecipeType, etc.)
  - Type mapping functionality
  - Edge cases with invalid values
  - String conversion methods

#### ✅ `version` - Version Parsing Tests

- **Coverage**: Semantic version handling
- **Test Count**: 8 test functions
- **Features Tested**:
  - Semantic version parsing (basic, v-prefix, pre-release)
  - Version validation and error handling
  - Version info generation
  - Stability detection (IsStable, IsPreRelease)

### Packages Needing Tests

#### ❌ `internal/tools` - MCP Tools

- **Priority**: High
- **Suggested Tests**:
  - Tool registration and discovery
  - Tool execution with various parameters
  - Error handling for invalid inputs
  - Integration with API client

#### ❌ `pkg/httpclient` - HTTP Client

- **Priority**: Medium
- **Suggested Tests**:
  - HTTP client configuration
  - Request/response handling
  - Timeout and retry logic
  - Authentication headers

#### ❌ `cmd/` - CLI Commands

- **Priority**: Medium
- **Suggested Tests**:
  - Command parsing and validation
  - Flag handling
  - Output formatting
  - Error scenarios

## Writing Tests

### Basic Test Structure

```go
func TestFunctionName(t *testing.T) {
    // Arrange
    input := "test input"
    expected := "expected output"

    // Act
    result := FunctionToTest(input)

    // Assert
    assert.Equal(t, expected, result)
}
```

### Table-Driven Tests

```go
func TestFunction(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected string
        wantErr  bool
    }{
        {
            name:     "Valid input",
            input:    "valid",
            expected: "processed",
            wantErr:  false,
        },
        {
            name:    "Invalid input",
            input:   "",
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := FunctionToTest(tt.input)

            if tt.wantErr {
                assert.Error(t, err)
                return
            }

            assert.NoError(t, err)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

### Testing with Context

```go
func TestWithContext(t *testing.T) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    result, err := FunctionWithContext(ctx, "input")
    assert.NoError(t, err)
    assert.NotEmpty(t, result)
}
```

### Benchmark Tests

```go
func BenchmarkFunction(b *testing.B) {
    input := "benchmark input"

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _ = FunctionToTest(input)
    }
}
```

## Test Utilities

### Available Helpers

The `internal/testutil` package provides:

#### TestServer

```go
// Create a mock HTTP server
server := testutil.NewTestServer()
defer server.Close()

// Add mock responses
server.AddResponse("POST", "/api/endpoint", testutil.MockResponse{
    StatusCode: 200,
    Body:       expectedResponse,
    Headers:    map[string]string{"Content-Type": "application/json"},
})
```

#### Test Context

```go
// Create context with timeout
ctx, cancel := testutil.CreateTestContext()
defer cancel()
```

#### Assertion Helpers

```go
// Custom assertion helpers
testutil.AssertNoError(t, err, "Operation should succeed")
testutil.AssertError(t, err, "Operation should fail")
testutil.AssertContains(t, haystack, needle, "Should contain substring")
```

#### JSON Validation

```go
// Validate JSON request bodies
testutil.AssertJSONRequest(t, request, expectedPayload)
```

### Mock Image Creation

```go
// Create mock image data for testing
imageData := testutil.CreateMockImageBytes("jpeg")
```

## Mocking and Test Data

### HTTP Mocking

Use the `TestServer` for mocking HTTP endpoints:

```go
server := testutil.NewTestServer()
defer server.Close()

// Mock successful response
server.AddResponse("GET", "/api/data", testutil.MockResponse{
    StatusCode: 200,
    Body:       map[string]string{"status": "success"},
})

// Mock error response
server.AddResponse("POST", "/api/create", testutil.MockResponse{
    StatusCode: 400,
    Body:       map[string]string{"error": "Bad request"},
})
```

### Test Data Files

Store test data in the `testdata/` directory:

```go
// Load test data
data := testutil.LoadTestData(t, "mock_responses.json")
```

### Environment Variables

```go
// Set test environment variables
t.Setenv("API_KEY", "test-key")
t.Setenv("BASE_URL", "http://localhost:8080")
```

## Best Practices

### Test Organization

1. **Group related tests** using subtests with `t.Run()`
2. **Use descriptive test names** that explain the scenario
3. **Follow AAA pattern**: Arrange, Act, Assert
4. **Test both success and failure cases**
5. **Include edge cases and boundary conditions**

### Error Testing

```go
// Test specific error types
var targetErr *CustomError
assert.ErrorAs(t, err, &targetErr)

// Test error messages
assert.Contains(t, err.Error(), "expected error text")

// Test error wrapping
assert.ErrorIs(t, err, ErrExpectedType)
```

### Performance Testing

```go
func BenchmarkCriticalFunction(b *testing.B) {
    // Setup
    input := setupBenchmarkData()

    b.ResetTimer()
    b.ReportAllocs() // Report memory allocations

    for i := 0; i < b.N; i++ {
        _ = CriticalFunction(input)
    }
}
```

### Test Cleanup

```go
func TestWithCleanup(t *testing.T) {
    // Setup
    resource := createResource()

    // Ensure cleanup happens
    t.Cleanup(func() {
        resource.Close()
    })

    // Test logic
    // ...
}
```

### Parallel Tests

```go
func TestParallel(t *testing.T) {
    tests := []struct{
        name string
        // test cases
    }{
        // test data
    }

    for _, tt := range tests {
        tt := tt // Capture range variable
        t.Run(tt.name, func(t *testing.T) {
            t.Parallel() // Run in parallel
            // test logic
        })
    }
}
```

## Troubleshooting

### Common Issues

#### Flaky Tests

- Use deterministic test data
- Avoid time-dependent assertions
- Use proper synchronization for concurrent code
- Set appropriate timeouts

#### Slow Tests

- Use `t.Skip()` for integration tests in unit test runs
- Implement proper mocking to avoid external dependencies
- Use `testing.Short()` to skip long-running tests

#### Memory Leaks

- Use `t.Cleanup()` for resource cleanup
- Check for goroutine leaks with `goleak`
- Profile memory usage in benchmarks

### Debugging Tests

```bash
# Run with verbose output
go test -v ./...

# Run specific test with debugging
go test -v -run TestSpecificFunction ./pkg/package

# Use delve debugger
dlv test ./pkg/package -- -test.run TestSpecificFunction
```

### Test Coverage Analysis

```bash
# Generate coverage profile
go test -coverprofile=coverage.out ./...

# View coverage by function
go tool cover -func=coverage.out

# Generate HTML coverage report
go tool cover -html=coverage.out -o coverage.html
```

## Learning Resources

### Go Testing Documentation

- [Go Testing Package](https://pkg.go.dev/testing)
- [Go Blog: Using Subtests and Sub-benchmarks](https://go.dev/blog/subtests)
- [Go Blog: Advanced Testing](https://go.dev/blog/advanced-go-testing)

### Testing Libraries

- [Testify](https://github.com/stretchr/testify) - Assertion library
- [GoMock](https://github.com/golang/mock) - Mock generation
- [Gock](https://github.com/h2non/gock) - HTTP mocking
- [GoLeak](https://github.com/uber-go/goleak) - Goroutine leak detection

### Best Practices

- [Go Testing Best Practices](https://go.dev/doc/tutorial/add-a-test)
- [Table Driven Tests](https://go.dev/wiki/TableDrivenTests)
- [Testing in Go](https://go.dev/doc/code#Testing)

---

## Summary

This testing guide provides a comprehensive overview of testing in the gaia-mcp-go project. We currently have robust test coverage for:

- ✅ **API Client** (`internal/api`) - 12 tests covering all major functionality
- ✅ **Image Processing** (`pkg/imageutil`) - 8 tests covering the complete pipeline
- ✅ **Shared Types** (`pkg/shared`) - 14 tests covering all constants and mappings
- ✅ **Version Handling** (`version`) - 8 tests covering semantic version parsing

**Next Steps for Complete Coverage:**

1. Add tests for MCP tools (`internal/tools`)
2. Add tests for HTTP client utilities (`pkg/httpclient`)
3. Add tests for CLI commands (`cmd/`)

The testing infrastructure is well-established with comprehensive utilities, mocking capabilities, and clear patterns to follow. All tests are currently passing and provide a solid foundation for continued development.
