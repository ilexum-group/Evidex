//go:build !windows

package os

// NewWindows is a stub for non-Windows platforms
// This function is only used during compilation on non-Windows systems
// and should never be called at runtime due to the New() factory logic
func NewWindows() OS {
	return NewDefault()
}
