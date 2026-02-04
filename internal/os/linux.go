//go:build linux

package os

import (
	"os"
	"os/user"
	"strconv"
	"syscall"
	"time"

	"github.com/ilexum-group/evidex/internal/utils"
)

// Linux is the Linux-specific implementation
type Linux struct {
	*DefaultImpl // Embed default implementation
}

// NewLinux creates a new Linux-specific OS implementation
func NewLinux() OS {
	return &Linux{
		DefaultImpl: NewDefault(),
	}
}

// ExtractAdvancedTimes extracts creation, access, and change times on Linux
func (l *Linux) ExtractAdvancedTimes(filePath string) (accessedTime time.Time, modifiedTime time.Time, changeTime time.Time, createdTime time.Time) {
	startTime := time.Now()
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		end := time.Now()
		if l.logger != nil {
			l.logger(utils.GenerateRandomID(), "os.Stat", []string{filePath}, startTime, end, 1, err, "", filePath)
		}
		return time.Time{}, time.Time{}, time.Time{}, time.Time{}
	}

	stat, ok := fileInfo.Sys().(*syscall.Stat_t)
	if !ok {
		end := time.Now()
		if l.logger != nil {
			l.logger(utils.GenerateRandomID(), "syscall.Stat_t.Type", []string{filePath}, startTime, end, 1, nil, "", filePath)
		}
		return time.Time{}, time.Time{}, time.Time{}, time.Time{}
	}

	// Extract all times
	accessedTime = time.Unix(stat.Atim.Sec, stat.Atim.Nsec)
	modifiedTime = time.Unix(stat.Mtim.Sec, stat.Mtim.Nsec)
	changeTime = time.Unix(stat.Ctim.Sec, stat.Ctim.Nsec)
	// Linux doesn't have birth time, use change time as fallback
	createdTime = changeTime

	endTime := time.Now()
	if l.logger != nil {
		l.logger(utils.GenerateRandomID(), "ExtractAdvancedTimes", []string{filePath}, startTime, endTime, 0, nil, "", filePath)
	}

	return accessedTime, modifiedTime, changeTime, createdTime
}

// ExtractOwnershipInfo extracts file owner and group information on Linux
func (l *Linux) ExtractOwnershipInfo(filePath string) (username string, groupname string) {
	startTime := time.Now()
	username = "unknown"
	groupname = "unknown"

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		end := time.Now()
		if l.logger != nil {
			l.logger(utils.GenerateRandomID(), "os.Stat", []string{filePath}, startTime, end, 1, err, "", filePath)
		}
		return username, groupname
	}

	stat, ok := fileInfo.Sys().(*syscall.Stat_t)
	if !ok {
		end := time.Now()
		if l.logger != nil {
			l.logger(utils.GenerateRandomID(), "syscall.Stat_t.Type", []string{filePath}, startTime, end, 1, nil, "", filePath)
		}
		return username, groupname
	}

	// Get UID and GID
	uid := stat.Uid
	gid := stat.Gid

	// Try to resolve user name from UID
	if u, err := user.LookupId(strconv.Itoa(int(uid))); err == nil {
		username = u.Username
	} else {
		username = strconv.Itoa(int(uid))
	}

	// Try to resolve group name from GID
	if g, err := user.LookupGroupId(strconv.Itoa(int(gid))); err == nil {
		groupname = g.Name
	} else {
		groupname = strconv.Itoa(int(gid))
	}

	endTime := time.Now()
	if l.logger != nil {
		l.logger(utils.GenerateRandomID(), "ExtractOwnershipInfo", []string{filePath, username, groupname}, startTime, endTime, 0, nil, "", filePath)
	}

	return username, groupname
}

// GetCurrentUser retrieves the current executing user on Linux
func (l *Linux) GetCurrentUser() (string, error) {
	startTime := time.Now()

	currentStart := time.Now()
	usr, err := user.Current()
	currentEnd := time.Now()
	if err != nil {
		if l.logger != nil {
			l.logger(utils.GenerateRandomID(), "user.Current", []string{}, currentStart, currentEnd, 1, err, "", "")
		}
		return "unknown", err
	}
	if l.logger != nil {
		l.logger(utils.GenerateRandomID(), "user.Current", []string{usr.Username, usr.Uid}, currentStart, currentEnd, 0, nil, "", "")
	}

	endTime := time.Now()
	if l.logger != nil {
		l.logger(utils.GenerateRandomID(), "GetCurrentUser", []string{usr.Username}, startTime, endTime, 0, nil, "", "")
	}

	return usr.Username, nil
}

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
