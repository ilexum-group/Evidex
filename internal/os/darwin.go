//go:build darwin

package os

import (
	"os"
	"os/user"
	"strconv"
	"syscall"
	"time"

	"github.com/ilexum-group/evidex/internal/utils"
)

// Darwin is the macOS (Darwin)-specific implementation
type Darwin struct {
	*DefaultImpl // Embed default implementation
}

// NewDarwin creates a new macOS-specific OS implementation
func NewDarwin() OS {
	return &Darwin{
		DefaultImpl: NewDefault(),
	}
}

// ExtractAdvancedTimes extracts creation, access, and change times on macOS
func (d *Darwin) ExtractAdvancedTimes(filePath string) (accessedTime time.Time, modifiedTime time.Time, changeTime time.Time, createdTime time.Time) {
	startTime := time.Now()
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		end := time.Now()
		if d.logger != nil {
			d.logger(utils.GenerateRandomID(), "os.Stat", []string{filePath}, startTime, end, 1, err, "", filePath)
		}
		return time.Time{}, time.Time{}, time.Time{}, time.Time{}
	}

	stat, ok := fileInfo.Sys().(*syscall.Stat_t)
	if !ok {
		end := time.Now()
		if d.logger != nil {
			d.logger(utils.GenerateRandomID(), "syscall.Stat_t.Type", []string{filePath}, startTime, end, 1, nil, "", filePath)
		}
		return time.Time{}, time.Time{}, time.Time{}, time.Time{}
	}

	// Extract all times
	accessedTime = time.Unix(stat.Atimespec.Sec, stat.Atimespec.Nsec)
	modifiedTime = time.Unix(stat.Mtimespec.Sec, stat.Mtimespec.Nsec)
	changeTime = time.Unix(stat.Ctimespec.Sec, stat.Ctimespec.Nsec)
	createdTime = time.Unix(stat.Birthtimespec.Sec, stat.Birthtimespec.Nsec)

	endTime := time.Now()
	if d.logger != nil {
		d.logger(utils.GenerateRandomID(), "ExtractAdvancedTimes", []string{filePath}, startTime, endTime, 0, nil, "", filePath)
	}

	return accessedTime, modifiedTime, changeTime, createdTime
}

// ExtractOwnershipInfo extracts file owner and group information on macOS
func (d *Darwin) ExtractOwnershipInfo(filePath string) (username string, groupname string) {
	startTime := time.Now()
	username = "unknown"
	groupname = "unknown"

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		end := time.Now()
		if d.logger != nil {
			d.logger(utils.GenerateRandomID(), "os.Stat", []string{filePath}, startTime, end, 1, err, "", filePath)
		}
		return username, groupname
	}

	stat, ok := fileInfo.Sys().(*syscall.Stat_t)
	if !ok {
		end := time.Now()
		if d.logger != nil {
			d.logger(utils.GenerateRandomID(), "syscall.Stat_t.Type", []string{filePath}, startTime, end, 1, nil, "", filePath)
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
	if d.logger != nil {
		d.logger(utils.GenerateRandomID(), "ExtractOwnershipInfo", []string{filePath, username, groupname}, startTime, endTime, 0, nil, "", filePath)
	}

	return username, groupname
}

// GetCurrentUser retrieves the current executing user on macOS
func (d *Darwin) GetCurrentUser() (string, error) {
	startTime := time.Now()

	currentStart := time.Now()
	usr, err := user.Current()
	currentEnd := time.Now()
	if err != nil {
		if d.logger != nil {
			d.logger(utils.GenerateRandomID(), "user.Current", []string{}, currentStart, currentEnd, 1, err, "", "")
		}
		return "unknown", err
	}
	if d.logger != nil {
		d.logger(utils.GenerateRandomID(), "user.Current", []string{usr.Username, usr.Uid}, currentStart, currentEnd, 0, nil, "", "")
	}

	endTime := time.Now()
	if d.logger != nil {
		d.logger(utils.GenerateRandomID(), "GetCurrentUser", []string{usr.Username}, startTime, endTime, 0, nil, "", "")
	}

	return usr.Username, nil
}

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
