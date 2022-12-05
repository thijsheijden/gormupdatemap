package gormupdatemap

// strPtr returns a pointer to the given string s. Used for inline pointers
func strPtr(s string) *string {
	return &s
}

// boolPtr returns a pointer to the given bool b. Used for inline pointers
func boolPtr(b bool) *bool {
	return &b
}

// float64Ptr returns a pointer to the given float64 f. Used for inline pointers
func float64Ptr(f float64) *float64 {
	return &f
}
