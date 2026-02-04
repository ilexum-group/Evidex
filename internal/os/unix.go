//go:build freebsd || openbsd || netbsd || dragonfly

package os

import (
	"os"
	"os/user"
	"strconv"
	"syscall"
	"time"

	"github.com/ilexum-group/evidex/internal/utils"
)

// Unix implements OS-specific metadata extraction for Unix-like systems
type Unix struct {
	*DefaultImpl // Embed default implementation
}

// NewUnix creates a new Unix-specific OS implementation
func NewUnix() OS {
	return &Unix{
		DefaultImpl: NewDefault(),
	}
}

// ExtractAdvancedTimes extracts creation, access, and change times on Unix
func (u *Unix) ExtractAdvancedTimes(filePath string) (accessedTime time.Time, modifiedTime time.Time, changeTime time.Time, createdTime time.Time) {
	startTime := time.Now()
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		end := time.Now()
		if u.logger != nil {
			u.logger(utils.GenerateRandomID(), "os.Stat", []string{filePath}, startTime, end, 1, err, "", filePath)
		}
		return time.Time{}, time.Time{}, time.Time{}, time.Time{}
	}

	stat, ok := fileInfo.Sys().(*syscall.Stat_t)
	if !ok {
		end := time.Now()
		if u.logger != nil {
			u.logger(utils.GenerateRandomID(), "syscall.Stat_t.Type", []string{filePath}, startTime, end, 1, nil, "", filePath)
		}
		return time.Time{}, time.Time{}, time.Time{}, time.Time{}
	}

	// Extract times
	accessedTime = getAccessTime(stat)
	modifiedTime = getModifiedTime(stat)
	changeTime = getChangeTime(stat)

	// Try to get birth time if available (macOS/BSD)
	if hasStatBirthtime() {
		createdTime = getStatBirthtime(stat)
	} else {
		// Use change time as fallback
		createdTime = changeTime
	}

	endTime := time.Now()
	if u.logger != nil {
		u.logger(utils.GenerateRandomID(), "ExtractAdvancedTimes", []string{filePath}, startTime, endTime, 0, nil, "", filePath)
	}

	return accessedTime, modifiedTime, changeTime, createdTime
}

// ExtractOwnershipInfo extracts file owner and group information on Unix
func (u *Unix) ExtractOwnershipInfo(filePath string) (username string, groupname string) {
	startTime := time.Now()
	username = "unknown"
	groupname = "unknown"

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		end := time.Now()
		if u.logger != nil {
			u.logger(utils.GenerateRandomID(), "os.Stat", []string{filePath}, startTime, end, 1, err, "", filePath)
		}
		return username, groupname
	}

	stat, ok := fileInfo.Sys().(*syscall.Stat_t)
	if !ok {
		end := time.Now()
		if u.logger != nil {
			u.logger(utils.GenerateRandomID(), "syscall.Stat_t.Type", []string{filePath}, startTime, end, 1, nil, "", filePath)
		}
		return username, groupname
	}

	// Get UID and GID
	uid := stat.Uid
	gid := stat.Gid

	// Try to resolve user name from UID
	if usr, err := user.LookupId(strconv.Itoa(int(uid))); err == nil {
		username = usr.Username
	} else {
		username = strconv.Itoa(int(uid))
	}

	// Try to resolve group name from GID
	if grp, err := user.LookupGroupId(strconv.Itoa(int(gid))); err == nil {
		groupname = grp.Name
	} else {
		groupname = strconv.Itoa(int(gid))
	}

	endTime := time.Now()
	if u.logger != nil {
		u.logger(utils.GenerateRandomID(), "ExtractOwnershipInfo", []string{filePath, username, groupname}, startTime, endTime, 0, nil, "", filePath)
	}

	return username, groupname
}

// GetCurrentUser retrieves the current executing user on Unix systems
func (u *Unix) GetCurrentUser() (string, error) {
	startTime := time.Now()

	currentStart := time.Now()
	usr, err := user.Current()
	currentEnd := time.Now()
	if err != nil {
		if u.logger != nil {
			u.logger(utils.GenerateRandomID(), "user.Current", []string{}, currentStart, currentEnd, 1, err, "", "")
		}
		return "unknown", err
	}
	if u.logger != nil {
		u.logger(utils.GenerateRandomID(), "user.Current", []string{usr.Username, usr.Uid}, currentStart, currentEnd, 0, nil, "", "")
	}

	endTime := time.Now()
	if u.logger != nil {
		u.logger(utils.GenerateRandomID(), "GetCurrentUser", []string{usr.Username}, startTime, endTime, 0, nil, "", "")
	}

	return usr.Username, nil
}
