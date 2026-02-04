//go:build !freebsd && !openbsd && !netbsd && !dragonfly

package os

// NewUnix is a stub for non-BSD platforms.
// This function is only used during compilation on non-BSD systems
// and should never be called at runtime due to the New() factory logic.
func NewUnix() OS {
	return NewDefault()
}
