package version

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestParseSemVer tests the semantic version parsing functionality
func TestParseSemVer(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected SemVer
		hasError bool
	}{
		{
			name:  "Valid basic version",
			input: "1.2.3",
			expected: SemVer{
				Major: 1,
				Minor: 2,
				Patch: 3,
			},
		},
		{
			name:  "Valid version with v prefix",
			input: "v1.2.3",
			expected: SemVer{
				Major: 1,
				Minor: 2,
				Patch: 3,
			},
		},
		{
			name:  "Valid version with pre-release",
			input: "1.2.3-alpha.1",
			expected: SemVer{
				Major:      1,
				Minor:      2,
				Patch:      3,
				PreRelease: "alpha.1",
			},
		},
		{
			name:  "Valid version with build metadata",
			input: "1.2.3+build.123",
			expected: SemVer{
				Major: 1,
				Minor: 2,
				Patch: 3,
				Build: "build.123",
			},
		},
		{
			name:  "Valid version with pre-release and build",
			input: "1.2.3-beta.2+build.456",
			expected: SemVer{
				Major:      1,
				Minor:      2,
				Patch:      3,
				PreRelease: "beta.2",
				Build:      "build.456",
			},
		},
		{
			name:  "Valid dev version",
			input: "0.1.0-dev",
			expected: SemVer{
				Major:      0,
				Minor:      1,
				Patch:      0,
				PreRelease: "dev",
			},
		},
		{
			name:     "Invalid version - missing patch",
			input:    "1.2",
			hasError: true,
		},
		{
			name:     "Invalid version - non-numeric major",
			input:    "a.2.3",
			hasError: true,
		},
		{
			name:     "Invalid version - non-numeric minor",
			input:    "1.b.3",
			hasError: true,
		},
		{
			name:     "Invalid version - non-numeric patch",
			input:    "1.2.c",
			hasError: true,
		},
		{
			name:     "Empty string",
			input:    "",
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseSemVer(tt.input)

			if tt.hasError {
				assert.Error(t, err, "Expected an error for input: %s", tt.input)
				return
			}

			require.NoError(t, err, "Unexpected error for input: %s", tt.input)
			assert.Equal(t, tt.expected, result, "Mismatch for input: %s", tt.input)
		})
	}
}

// TestSemVerString tests the String() method of SemVer
func TestSemVerString(t *testing.T) {
	tests := []struct {
		name     string
		semver   SemVer
		expected string
	}{
		{
			name: "Basic version",
			semver: SemVer{
				Major: 1,
				Minor: 2,
				Patch: 3,
			},
			expected: "1.2.3",
		},
		{
			name: "Version with pre-release",
			semver: SemVer{
				Major:      1,
				Minor:      2,
				Patch:      3,
				PreRelease: "alpha.1",
			},
			expected: "1.2.3-alpha.1",
		},
		{
			name: "Version with build",
			semver: SemVer{
				Major: 1,
				Minor: 2,
				Patch: 3,
				Build: "build.123",
			},
			expected: "1.2.3+build.123",
		},
		{
			name: "Version with pre-release and build",
			semver: SemVer{
				Major:      1,
				Minor:      2,
				Patch:      3,
				PreRelease: "beta.2",
				Build:      "build.456",
			},
			expected: "1.2.3-beta.2+build.456",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.semver.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestSemVerIsPreRelease tests the IsPreRelease() method
func TestSemVerIsPreRelease(t *testing.T) {
	tests := []struct {
		name     string
		semver   SemVer
		expected bool
	}{
		{
			name: "Stable version",
			semver: SemVer{
				Major: 1,
				Minor: 2,
				Patch: 3,
			},
			expected: false,
		},
		{
			name: "Pre-release version",
			semver: SemVer{
				Major:      1,
				Minor:      2,
				Patch:      3,
				PreRelease: "alpha.1",
			},
			expected: true,
		},
		{
			name: "Dev version",
			semver: SemVer{
				Major:      0,
				Minor:      1,
				Patch:      0,
				PreRelease: "dev",
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.semver.IsPreRelease()
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestSemVerIsStable tests the IsStable() method
func TestSemVerIsStable(t *testing.T) {
	tests := []struct {
		name     string
		semver   SemVer
		expected bool
	}{
		{
			name: "Stable version 1.0.0",
			semver: SemVer{
				Major: 1,
				Minor: 0,
				Patch: 0,
			},
			expected: true,
		},
		{
			name: "Stable version 2.5.3",
			semver: SemVer{
				Major: 2,
				Minor: 5,
				Patch: 3,
			},
			expected: true,
		},
		{
			name: "Pre-release version",
			semver: SemVer{
				Major:      1,
				Minor:      0,
				Patch:      0,
				PreRelease: "alpha.1",
			},
			expected: false,
		},
		{
			name: "Version 0.x.x (not stable)",
			semver: SemVer{
				Major: 0,
				Minor: 5,
				Patch: 3,
			},
			expected: false,
		},
		{
			name: "Dev version",
			semver: SemVer{
				Major:      0,
				Minor:      1,
				Patch:      0,
				PreRelease: "dev",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.semver.IsStable()
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestGet tests the Get() function
func TestGet(t *testing.T) {
	// Store original values
	originalVersion := Version
	originalGitCommit := GitCommit
	originalBuildDate := BuildDate

	// Restore original values after test
	defer func() {
		Version = originalVersion
		GitCommit = originalGitCommit
		BuildDate = originalBuildDate
	}()

	t.Run("Test with valid version", func(t *testing.T) {
		// Set test values
		Version = "1.2.3"
		GitCommit = "abc123"
		BuildDate = "2024-01-01_12:00:00"

		info := Get()

		// Check basic fields
		assert.Equal(t, "1.2.3", info.Raw)
		assert.Equal(t, "abc123", info.GitCommit)
		assert.Equal(t, "2024-01-01_12:00:00", info.BuildDate)
		assert.Equal(t, runtime.Version(), info.GoVersion)
		assert.Contains(t, info.Platform, runtime.GOOS)
		assert.Contains(t, info.Platform, runtime.GOARCH)

		// Check parsed version
		assert.Equal(t, 1, info.Version.Major)
		assert.Equal(t, 2, info.Version.Minor)
		assert.Equal(t, 3, info.Version.Patch)
	})

	t.Run("Test with invalid version fallback", func(t *testing.T) {
		// Set invalid version
		Version = "invalid-version"
		GitCommit = "def456"
		BuildDate = "2024-01-01_12:00:00"

		info := Get()

		// Check that fallback version is used
		assert.Equal(t, "invalid-version", info.Raw)
		assert.Equal(t, 0, info.Version.Major)
		assert.Equal(t, 1, info.Version.Minor)
		assert.Equal(t, 0, info.Version.Patch)
		assert.Equal(t, "dev", info.Version.PreRelease)
	})
}

// TestInfoString tests the String() method of Info
func TestInfoString(t *testing.T) {
	// Store original values
	originalVersion := Version
	originalGitCommit := GitCommit
	originalBuildDate := BuildDate

	// Restore original values after test
	defer func() {
		Version = originalVersion
		GitCommit = originalGitCommit
		BuildDate = originalBuildDate
	}()

	t.Run("Test info string format", func(t *testing.T) {
		Version = "1.2.3"
		GitCommit = "abc123"
		BuildDate = "2024-01-01_12:00:00"

		info := Get()
		result := info.String()

		// Check that the string contains expected components
		assert.Contains(t, result, "gaia-mcp-server")
		assert.Contains(t, result, "v1.2.3")
		assert.Contains(t, result, "abc123")
		assert.Contains(t, result, "2024-01-01_12:00:00")
		assert.Contains(t, result, runtime.Version())
		assert.Contains(t, result, runtime.GOOS)
		assert.Contains(t, result, runtime.GOARCH)
	})
}

// TestInfoShort tests the Short() method of Info
func TestInfoShort(t *testing.T) {
	// Store original values
	originalVersion := Version

	// Restore original values after test
	defer func() {
		Version = originalVersion
	}()

	t.Run("Test short version format", func(t *testing.T) {
		Version = "1.2.3-beta.1"

		info := Get()
		result := info.Short()

		assert.Equal(t, "1.2.3-beta.1", result)
	})
}

// TestInfoIsStable tests the IsStable() method of Info
func TestInfoIsStable(t *testing.T) {
	// Store original values
	originalVersion := Version

	// Restore original values after test
	defer func() {
		Version = originalVersion
	}()

	tests := []struct {
		name     string
		version  string
		expected bool
	}{
		{
			name:     "Stable version",
			version:  "1.0.0",
			expected: true,
		},
		{
			name:     "Pre-release version",
			version:  "1.0.0-alpha.1",
			expected: false,
		},
		{
			name:     "Dev version",
			version:  "0.1.0-dev",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Version = tt.version

			info := Get()
			result := info.IsStable()

			assert.Equal(t, tt.expected, result)
		})
	}
}

// Benchmark tests
func BenchmarkParseSemVer(b *testing.B) {
	version := "1.2.3-alpha.1+build.123"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ParseSemVer(version)
	}
}

func BenchmarkSemVerString(b *testing.B) {
	semver := SemVer{
		Major:      1,
		Minor:      2,
		Patch:      3,
		PreRelease: "alpha.1",
		Build:      "build.123",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = semver.String()
	}
}

func BenchmarkGet(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = Get()
	}
}
