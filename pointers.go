package unstructured

// String returns a pointer to the given string value.
// This is useful when you need to pass optional string values to API requests.
func String(s string) *string {
	return &s
}

// Bool returns a pointer to the given boolean value.
// This is useful when you need to pass optional boolean values to API requests.
func Bool(b bool) *bool {
	return &b
}

// Int returns a pointer to the given integer value.
// This is useful when you need to pass optional integer values to API requests.
func Int(i int) *int {
	return &i
}

// ToString converts a string pointer to a string value.
// If the pointer is nil, it returns an empty string.
func ToString(p *string) string {
	if p == nil {
		return ""
	}

	return *p
}

// ToBool converts a boolean pointer to a boolean value.
// If the pointer is nil, it returns false.
func ToBool(p *bool) bool {
	if p == nil {
		return false
	}

	return *p
}

// ToInt converts an integer pointer to an integer value.
// If the pointer is nil, it returns 0.
func ToInt(p *int) int {
	if p == nil {
		return 0
	}

	return *p
}
