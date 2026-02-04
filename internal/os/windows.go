//go:build windows

package os

import (
	"syscall"
	"time"

	"github.com/ilexum-group/evidex/internal/utils"
	"golang.org/x/sys/windows"
)

// Windows implements OS-specific metadata extraction for Windows.
type Windows struct {
	*DefaultImpl // Embed default implementation
}

// NewWindows creates a new Windows-specific OS implementation.
func NewWindows() OS {
	return &Windows{
		DefaultImpl: NewDefault(),
	}
}

// ExtractAdvancedTimes extracts creation, access, and change times on Windows.
func (w *Windows) ExtractAdvancedTimes(filePath string) (accessedTime time.Time, modifiedTime time.Time, changeTime time.Time, createdTime time.Time) {
	startTime := time.Now()

	// Open file handle
	pathPtr, err := syscall.UTF16PtrFromString(filePath)
	if err != nil {
		end := time.Now()
		if w.logger != nil {
			w.logger(utils.GenerateRandomID(), "syscall.UTF16PtrFromString", []string{filePath}, startTime, end, 1, err, "", filePath)
		}
		return time.Time{}, time.Time{}, time.Time{}, time.Time{}
	}

	handleStart := time.Now()
	handle, err := syscall.CreateFile(
		pathPtr,
		syscall.GENERIC_READ,
		syscall.FILE_SHARE_READ,
		nil,
		syscall.OPEN_EXISTING,
		syscall.FILE_FLAG_BACKUP_SEMANTICS,
		0,
	)
	handleEnd := time.Now()
	if err != nil {
		if w.logger != nil {
			w.logger(utils.GenerateRandomID(), "syscall.CreateFile", []string{filePath, "GENERIC_READ"}, handleStart, handleEnd, 1, err, "", filePath)
		}
		return time.Time{}, time.Time{}, time.Time{}, time.Time{}
	}
	if w.logger != nil {
		w.logger(utils.GenerateRandomID(), "syscall.CreateFile", []string{filePath, "GENERIC_READ"}, handleStart, handleEnd, 0, nil, "", filePath)
	}

	defer func() {
		closeStart := time.Now()
		closeErr := syscall.CloseHandle(handle)
		closeEnd := time.Now()
		exitCode := 0
		if closeErr != nil {
			exitCode = 1
		}
		if w.logger != nil {
			w.logger(utils.GenerateRandomID(), "syscall.CloseHandle", []string{filePath}, closeStart, closeEnd, exitCode, closeErr, "", filePath)
		}
	}()

	// Get file times
	var creationTime, accessTime, writeTime windows.Filetime
	getTimeStart := time.Now()
	err = windows.GetFileTime(windows.Handle(handle), &creationTime, &accessTime, &writeTime)
	getTimeEnd := time.Now()
	if err != nil {
		if w.logger != nil {
			w.logger(utils.GenerateRandomID(), "windows.GetFileTime", []string{filePath}, getTimeStart, getTimeEnd, 1, err, "", filePath)
		}
		return time.Time{}, time.Time{}, time.Time{}, time.Time{}
	}
	if w.logger != nil {
		w.logger(utils.GenerateRandomID(), "windows.GetFileTime", []string{filePath}, getTimeStart, getTimeEnd, 0, nil, "", filePath)
	}

	// Convert Windows FILETIME to Go time.Time
	createdTime = time.Unix(0, creationTime.Nanoseconds())
	accessedTime = time.Unix(0, accessTime.Nanoseconds())
	modifiedTime = time.Unix(0, writeTime.Nanoseconds())
	// On Windows, change time is same as modified time for files
	changeTime = modifiedTime

	endTime := time.Now()
	if w.logger != nil {
		w.logger(utils.GenerateRandomID(), "ExtractAdvancedTimes", []string{filePath}, startTime, endTime, 0, nil, "", filePath)
	}

	return accessedTime, modifiedTime, changeTime, createdTime
}

// ExtractOwnershipInfo extracts file owner and group information on Windows.
func (w *Windows) ExtractOwnershipInfo(filePath string) (username string, groupname string) {
	startTime := time.Now()
	username = "unknown"
	groupname = "unknown"

	// Get security descriptor for the file
	sdStart := time.Now()
	sd, err := windows.GetNamedSecurityInfo(
		filePath,
		windows.SE_FILE_OBJECT,
		windows.OWNER_SECURITY_INFORMATION|windows.GROUP_SECURITY_INFORMATION,
	)
	sdEnd := time.Now()
	if err != nil {
		if w.logger != nil {
			w.logger(utils.GenerateRandomID(), "windows.GetNamedSecurityInfo", []string{filePath}, sdStart, sdEnd, 1, err, "", filePath)
		}
		return username, groupname
	}
	if w.logger != nil {
		w.logger(utils.GenerateRandomID(), "windows.GetNamedSecurityInfo", []string{filePath}, sdStart, sdEnd, 0, nil, "", filePath)
	}

	// Get owner SID
	username = w.lookupSIDWithLogging(sd.Owner, "sd.Owner", filePath)

	// Get group SID
	groupname = w.lookupSIDWithLogging(sd.Group, "sd.Group", filePath)

	endTime := time.Now()
	if w.logger != nil {
		w.logger(utils.GenerateRandomID(), "ExtractOwnershipInfo", []string{filePath, username, groupname}, startTime, endTime, 0, nil, "", filePath)
	}

	return username, groupname
}

// GetCurrentUser retrieves the current executing user on Windows.
func (w *Windows) GetCurrentUser() (string, error) {
	startTime := time.Now()

	// Get current user token
	var token windows.Token
	tokenStart := time.Now()
	err := windows.OpenProcessToken(windows.CurrentProcess(), windows.TOKEN_QUERY, &token)
	tokenEnd := time.Now()
	if err != nil {
		if w.logger != nil {
			w.logger(utils.GenerateRandomID(), "windows.OpenProcessToken", []string{}, tokenStart, tokenEnd, 1, err, "", "")
		}
		return "unknown", err
	}
	if w.logger != nil {
		w.logger(utils.GenerateRandomID(), "windows.OpenProcessToken", []string{}, tokenStart, tokenEnd, 0, nil, "", "")
	}
	defer func() {
		_ = token.Close()
	}()

	// Get token user
	userStart := time.Now()
	tokenUser, err := token.GetTokenUser()
	userEnd := time.Now()
	if err != nil {
		if w.logger != nil {
			w.logger(utils.GenerateRandomID(), "token.GetTokenUser", []string{}, userStart, userEnd, 1, err, "", "")
		}
		return "unknown", err
	}
	if w.logger != nil {
		w.logger(utils.GenerateRandomID(), "token.GetTokenUser", []string{}, userStart, userEnd, 0, nil, "", "")
	}

	// Lookup account from SID
	lookupStart := time.Now()
	account, domain, _, err := tokenUser.User.Sid.LookupAccount("")
	lookupEnd := time.Now()
	if err != nil {
		if w.logger != nil {
			w.logger(utils.GenerateRandomID(), "Sid.LookupAccount", []string{}, lookupStart, lookupEnd, 1, err, "", "")
		}
		return tokenUser.User.Sid.String(), nil
	}
	if w.logger != nil {
		w.logger(utils.GenerateRandomID(), "Sid.LookupAccount", []string{account, domain}, lookupStart, lookupEnd, 0, nil, "", "")
	}

	var username string
	if domain != "" {
		username = domain + "\\" + account
	} else {
		username = account
	}

	endTime := time.Now()
	if w.logger != nil {
		w.logger(utils.GenerateRandomID(), "GetCurrentUser", []string{username}, startTime, endTime, 0, nil, "", "")
	}

	return username, nil
}

// lookupSIDWithLogging is a helper function that looks up a SID and logs the operation.
// This eliminates code duplication between owner and group lookups..
func (w *Windows) lookupSIDWithLogging(sidGetter func() (*windows.SID, bool, error), operation, filePath string) string {
	startTime := time.Now()
	sid, _, err := sidGetter()
	endTime := time.Now()

	if err != nil {
		if w.logger != nil {
			w.logger(utils.GenerateRandomID(), operation, []string{filePath}, startTime, endTime, 1, err, "", filePath)
		}
		return ""
	}

	if w.logger != nil {
		w.logger(utils.GenerateRandomID(), operation, []string{filePath}, startTime, endTime, 0, nil, "", filePath)
	}

	// Try to lookup account name from SID
	lookupStart := time.Now()
	account, domain, _, lookupErr := sid.LookupAccount("")
	lookupEnd := time.Now()

	if lookupErr != nil {
		if w.logger != nil {
			w.logger(utils.GenerateRandomID(), "SID.LookupAccount", []string{filePath}, lookupStart, lookupEnd, 1, lookupErr, "", filePath)
		}
		// If lookup fails, use SID string
		return sid.String()
	}

	if w.logger != nil {
		w.logger(utils.GenerateRandomID(), "SID.LookupAccount", []string{filePath, account, domain}, lookupStart, lookupEnd, 0, nil, "", filePath)
	}

	if domain != "" {
		return domain + "\\" + account
	}
	return account
}
