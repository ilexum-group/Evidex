//go:build darwin

package os

import (
	"syscall"
	"time"
)

// hasStatBirthtime returns true on macOS as it supports birth time
func hasStatBirthtime() bool {
	return true
}

// getStatBirthtime extracts birth time from stat on macOS
func getStatBirthtime(stat *syscall.Stat_t) time.Time {
	return time.Unix(stat.Birthtimespec.Sec, stat.Birthtimespec.Nsec)
}

// getAccessTime extracts access time from stat on macOS
func getAccessTime(stat *syscall.Stat_t) time.Time {
	return time.Unix(stat.Atimespec.Sec, stat.Atimespec.Nsec)
}

// getModifiedTime extracts modified time from stat on macOS
func getModifiedTime(stat *syscall.Stat_t) time.Time {
	return time.Unix(stat.Mtimespec.Sec, stat.Mtimespec.Nsec)
}

// getChangeTime extracts change time from stat on macOS
func getChangeTime(stat *syscall.Stat_t) time.Time {
	return time.Unix(stat.Ctimespec.Sec, stat.Ctimespec.Nsec)
}
