//go:build !linux

package os

// NewLinux is a stub for non-Linux platforms.
// This function is only used during compilation on non-Linux systems
// and should never be called at runtime due to the New() factory logic.
func NewLinux() OS {
	return NewDefault()
}
