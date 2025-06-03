package shared

import (
	"fmt"
	"time"
)

// Common time formats used in the application
const (
	// RFC3339 is the standard ISO 8601 format: "2006-01-02T15:04:05Z07:00"
	TimeFormatRFC3339 = time.RFC3339

	// RFC3339Nano includes nanoseconds: "2006-01-02T15:04:05.999999999Z07:00"
	TimeFormatRFC3339Nano = time.RFC3339Nano

	// Custom format if your API uses a different format
	TimeFormatCustom = "2006-01-02 15:04:05"
)

// ParseTimeString converts a time string to time.Time using multiple format attempts
// Returns zero time and error if parsing fails with all formats
func ParseTimeString(timeStr string) (time.Time, error) {
	if timeStr == "" {
		return time.Time{}, fmt.Errorf("empty time string")
	}

	// Try parsing with different formats in order of preference
	formats := []string{
		TimeFormatRFC3339Nano,
		TimeFormatRFC3339,
		TimeFormatCustom,
		"2006-01-02T15:04:05",  // Without timezone
		"2006-01-02 15:04:05Z", // Alternative format
	}

	for _, format := range formats {
		if t, err := time.Parse(format, timeStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse time string: %s", timeStr)
}

// ParseOptionalTimeString converts an optional time string (*string) to *time.Time
// Returns nil if input is nil, otherwise attempts to parse the string
func ParseOptionalTimeString(timeStr *string) (*time.Time, error) {
	if timeStr == nil {
		return nil, nil
	}

	t, err := ParseTimeString(*timeStr)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

// MustParseTimeString is like ParseTimeString but panics on error
// Use only when you're certain the time string is valid
func MustParseTimeString(timeStr string) time.Time {
	t, err := ParseTimeString(timeStr)
	if err != nil {
		panic(fmt.Sprintf("failed to parse time string '%s': %v", timeStr, err))
	}
	return t
}

// FormatTime converts time.Time to string using RFC3339 format
func FormatTime(t time.Time) string {
	return t.Format(TimeFormatRFC3339)
}

// FormatOptionalTime converts *time.Time to *string using RFC3339 format
// Returns nil if input is nil
func FormatOptionalTime(t *time.Time) *string {
	if t == nil {
		return nil
	}

	formatted := FormatTime(*t)
	return &formatted
}

// IsTimeExpired checks if a time string represents a time in the past
func IsTimeExpired(timeStr *string) bool {
	if timeStr == nil {
		return false
	}

	t, err := ParseTimeString(*timeStr)
	if err != nil {
		return false
	}

	return time.Now().After(t)
}

// GetTimeDuration calculates duration between two time strings
func GetTimeDuration(startTimeStr, endTimeStr string) (time.Duration, error) {
	startTime, err := ParseTimeString(startTimeStr)
	if err != nil {
		return 0, fmt.Errorf("invalid start time: %v", err)
	}

	endTime, err := ParseTimeString(endTimeStr)
	if err != nil {
		return 0, fmt.Errorf("invalid end time: %v", err)
	}

	return endTime.Sub(startTime), nil
}

// TimePtr creates a pointer to time.Time (helper for optional fields)
func TimePtr(t time.Time) *time.Time {
	return &t
}
