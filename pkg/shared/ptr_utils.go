package shared

// Pointer utility functions for working with optional fields in structs
// These helpers make it easier to create pointers to basic types

// StringPtr creates a pointer to string (helper for optional fields)
func StringPtr(s string) *string {
	return &s
}

// BoolPtr creates a pointer to bool (helper for optional fields)
func BoolPtr(b bool) *bool {
	return &b
}

// IntPtr creates a pointer to int (helper for optional fields)
func IntPtr(i int) *int {
	return &i
}

// Int64Ptr creates a pointer to int64 (helper for optional fields)
func Int64Ptr(i int64) *int64 {
	return &i
}

// Float64Ptr creates a pointer to float64 (helper for optional fields)
func Float64Ptr(f float64) *float64 {
	return &f
}

// Utility functions to safely dereference pointers with default values

// StringValue safely dereferences a string pointer, returning empty string if nil
func StringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// BoolValue safely dereferences a bool pointer, returning false if nil
func BoolValue(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

// IntValue safely dereferences an int pointer, returning 0 if nil
func IntValue(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}

// Int64Value safely dereferences an int64 pointer, returning 0 if nil
func Int64Value(i *int64) int64 {
	if i == nil {
		return 0
	}
	return *i
}

// Float64Value safely dereferences a float64 pointer, returning 0.0 if nil
func Float64Value(f *float64) float64 {
	if f == nil {
		return 0.0
	}
	return *f
}
