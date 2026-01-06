//go:build linux

package os

import (
	"syscall"
	"time"
)

// hasStatBirthtime returns false on Linux as it typically doesn't support birth time
func hasStatBirthtime() bool {
	return false
}

// getStatBirthtime returns change time as fallback on Linux
func getStatBirthtime(stat *syscall.Stat_t) time.Time {
	// Linux doesn't have birth time, return change time as fallback
	return time.Unix(stat.Ctim.Sec, stat.Ctim.Nsec)
}

// getAccessTime extracts access time from stat on Linux
func getAccessTime(stat *syscall.Stat_t) time.Time {
	return time.Unix(stat.Atim.Sec, stat.Atim.Nsec)
}

// getModifiedTime extracts modified time from stat on Linux
func getModifiedTime(stat *syscall.Stat_t) time.Time {
	return time.Unix(stat.Mtim.Sec, stat.Mtim.Nsec)
}

// getChangeTime extracts change time from stat on Linux
func getChangeTime(stat *syscall.Stat_t) time.Time {
	return time.Unix(stat.Ctim.Sec, stat.Ctim.Nsec)
}
