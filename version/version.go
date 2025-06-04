package version

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

var (
	// These will be set during build time using ldflags
	Version   = "0.1.0-dev"
	GitCommit = "unknown"
	BuildDate = "unknown"
)

// SemVer represents a semantic version
type SemVer struct {
	Major      int    `json:"major"`
	Minor      int    `json:"minor"`
	Patch      int    `json:"patch"`
	PreRelease string `json:"preRelease,omitempty"`
	Build      string `json:"build,omitempty"`
}

// Info contains comprehensive version information
type Info struct {
	Version   SemVer `json:"version"`
	Raw       string `json:"raw"`
	GitCommit string `json:"gitCommit"`
	BuildDate string `json:"buildDate"`
	GoVersion string `json:"goVersion"`
	Platform  string `json:"platform"`
}

// ParseSemVer parses a semantic version string
func ParseSemVer(v string) (SemVer, error) {
	// Remove 'v' prefix if present
	v = strings.TrimPrefix(v, "v")

	var semver SemVer

	// Split by + for build metadata
	parts := strings.Split(v, "+")
	if len(parts) > 1 {
		semver.Build = parts[1]
	}

	// Split by - for pre-release
	versionPart := parts[0]
	releaseParts := strings.Split(versionPart, "-")
	if len(releaseParts) > 1 {
		semver.PreRelease = strings.Join(releaseParts[1:], "-")
	}

	// Parse major.minor.patch
	versionNumbers := strings.Split(releaseParts[0], ".")
	if len(versionNumbers) < 3 {
		return semver, fmt.Errorf("invalid semantic version: %s", v)
	}

	var err error
	semver.Major, err = strconv.Atoi(versionNumbers[0])
	if err != nil {
		return semver, fmt.Errorf("invalid major version: %s", versionNumbers[0])
	}

	semver.Minor, err = strconv.Atoi(versionNumbers[1])
	if err != nil {
		return semver, fmt.Errorf("invalid minor version: %s", versionNumbers[1])
	}

	semver.Patch, err = strconv.Atoi(versionNumbers[2])
	if err != nil {
		return semver, fmt.Errorf("invalid patch version: %s", versionNumbers[2])
	}

	return semver, nil
}

// String returns the semantic version as a string
func (s SemVer) String() string {
	version := fmt.Sprintf("%d.%d.%d", s.Major, s.Minor, s.Patch)

	if s.PreRelease != "" {
		version += "-" + s.PreRelease
	}

	if s.Build != "" {
		version += "+" + s.Build
	}

	return version
}

// IsPreRelease returns true if this is a pre-release version
func (s SemVer) IsPreRelease() bool {
	return s.PreRelease != ""
}

// IsStable returns true if this is a stable release (>= 1.0.0 and no pre-release)
func (s SemVer) IsStable() bool {
	return s.Major >= 1 && !s.IsPreRelease()
}

// Get returns comprehensive version information
func Get() Info {
	semver, err := ParseSemVer(Version)
	if err != nil {
		// Fallback for invalid version strings
		semver = SemVer{Major: 0, Minor: 1, Patch: 0, PreRelease: "dev"}
	}

	return Info{
		Version:   semver,
		Raw:       Version,
		GitCommit: GitCommit,
		BuildDate: BuildDate,
		GoVersion: runtime.Version(),
		Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}

// String returns a formatted version string
func (i Info) String() string {
	return fmt.Sprintf("gaia-mcp-server v%s (%s) built on %s for %s with %s",
		i.Version.String(), i.GitCommit, i.BuildDate, i.Platform, i.GoVersion)
}

// Short returns just the version number
func (i Info) Short() string {
	return i.Version.String()
}

// IsStable returns true if this is a stable release
func (i Info) IsStable() bool {
	return i.Version.IsStable()
}
