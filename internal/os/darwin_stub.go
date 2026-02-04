//go:build !darwin

// Package os provides operating system-specific metadata and forensic information collection.
package os

// NewDarwin is a stub for non-Darwin (macOS) platforms.
// This function is only used during compilation on non-Darwin systems
// and should never be called at runtime due to the New() factory logic.
func NewDarwin() OS {
	return NewDefault()
}
