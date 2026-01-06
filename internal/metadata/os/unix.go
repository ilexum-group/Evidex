//go:build !windows

package os

import (
	"os"
	"os/user"
	"strconv"
	"syscall"

	"github.com/evidex/internal/models"
)

// UnixExtractor implements OS-specific metadata extraction for Unix-like systems
type UnixExtractor struct{}

// newExtractor creates a new UnixExtractor instance for non-Windows platforms
func newExtractor() OSMetadataExtractor {
	return &UnixExtractor{}
}

// GetOSName returns the name of the operating system
func (u *UnixExtractor) GetOSName() string {
	return "Unix"
}

// ExtractAdvancedTimes extracts creation, access, and change times on Unix/Linux/macOS
func (u *UnixExtractor) ExtractAdvancedTimes(filePath string, evidence *models.FileEvidence) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		evidence.AccessedTime = evidence.ModifiedTime
		evidence.CreatedTime = evidence.ModifiedTime
		evidence.ChangeTime = evidence.ModifiedTime
		return
	}

	stat, ok := fileInfo.Sys().(*syscall.Stat_t)
	if !ok {
		evidence.AccessedTime = evidence.ModifiedTime
		evidence.CreatedTime = evidence.ModifiedTime
		evidence.ChangeTime = evidence.ModifiedTime
		return
	}

	// Access time
	evidence.AccessedTime = getAccessTime(stat)

	// Modified time (already set, but we can update it from stat for consistency)
	evidence.ModifiedTime = getModifiedTime(stat)

	// Change time (inode change time)
	evidence.ChangeTime = getChangeTime(stat)

	// Birth time (creation time) - platform specific
	// On Linux: Ctim is inode change time, not birth time (birth time may not be available)
	// On macOS: Birthtim is available
	// On BSD: Birthtim is available
	// For maximum compatibility, we use Ctim as fallback
	evidence.CreatedTime = evidence.ChangeTime

	// Try to get birth time if available (macOS/BSD)
	if hasStatBirthtime() {
		evidence.CreatedTime = getStatBirthtime(stat)
	}
}

// ExtractOwnershipInfo extracts file owner and group information on Unix/Linux/macOS
func (u *UnixExtractor) ExtractOwnershipInfo(filePath string, evidence *models.FileEvidence) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		evidence.Owner = "unknown"
		evidence.Group = "unknown"
		return
	}

	stat, ok := fileInfo.Sys().(*syscall.Stat_t)
	if !ok {
		evidence.Owner = "unknown"
		evidence.Group = "unknown"
		return
	}

	// Get UID and GID
	uid := stat.Uid
	gid := stat.Gid

	// Try to resolve user name from UID
	if u, err := user.LookupId(strconv.Itoa(int(uid))); err == nil {
		evidence.Owner = u.Username
	} else {
		evidence.Owner = strconv.Itoa(int(uid))
	}

	// Try to resolve group name from GID
	if g, err := user.LookupGroupId(strconv.Itoa(int(gid))); err == nil {
		evidence.Group = g.Name
	} else {
		evidence.Group = strconv.Itoa(int(gid))
	}
}
