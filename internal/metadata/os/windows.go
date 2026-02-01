//go:build windows

package os

import (
	"os"
	"syscall"
	"time"

	"github.com/ilexum/evidex/internal/models"
	"golang.org/x/sys/windows"
)

// WindowsExtractor implements OS-specific metadata extraction for Windows
type WindowsExtractor struct{}

// newExtractor creates a new WindowsExtractor instance
func newExtractor() OSMetadataExtractor {
	return &WindowsExtractor{}
}

// GetOSName returns the name of the operating system
func (w *WindowsExtractor) GetOSName() string {
	return "Windows"
}

// ExtractAdvancedTimes extracts creation, access, and change times on Windows
func (w *WindowsExtractor) ExtractAdvancedTimes(filePath string, evidence *models.FileEvidence) {
	// Open file handle
	pathPtr, err := syscall.UTF16PtrFromString(filePath)
	if err != nil {
		evidence.AccessedTime = evidence.ModifiedTime
		evidence.CreatedTime = evidence.ModifiedTime
		evidence.ChangeTime = evidence.ModifiedTime
		return
	}

	handle, err := syscall.CreateFile(
		pathPtr,
		syscall.GENERIC_READ,
		syscall.FILE_SHARE_READ,
		nil,
		syscall.OPEN_EXISTING,
		syscall.FILE_FLAG_BACKUP_SEMANTICS,
		0,
	)
	if err != nil {
		evidence.AccessedTime = evidence.ModifiedTime
		evidence.CreatedTime = evidence.ModifiedTime
		evidence.ChangeTime = evidence.ModifiedTime
		return
	}
	defer func() {
		_ = syscall.CloseHandle(handle)
	}()

	// Get file times
	var creationTime, accessTime, writeTime windows.Filetime
	err = windows.GetFileTime(windows.Handle(handle), &creationTime, &accessTime, &writeTime)
	if err != nil {
		evidence.AccessedTime = evidence.ModifiedTime
		evidence.CreatedTime = evidence.ModifiedTime
		evidence.ChangeTime = evidence.ModifiedTime
		return
	}

	// Convert Windows FILETIME to Go time.Time
	evidence.CreatedTime = time.Unix(0, creationTime.Nanoseconds())
	evidence.AccessedTime = time.Unix(0, accessTime.Nanoseconds())
	evidence.ModifiedTime = time.Unix(0, writeTime.Nanoseconds())
	// On Windows, change time is same as modified time for files
	evidence.ChangeTime = evidence.ModifiedTime
}

// ExtractOwnershipInfo extracts file owner and group information on Windows
func (w *WindowsExtractor) ExtractOwnershipInfo(filePath string, evidence *models.FileEvidence) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		evidence.Owner = "unknown"
		evidence.Group = "unknown"
		return
	}

	// Get security descriptor for the file
	sd, err := windows.GetNamedSecurityInfo(
		filePath,
		windows.SE_FILE_OBJECT,
		windows.OWNER_SECURITY_INFORMATION|windows.GROUP_SECURITY_INFORMATION,
	)
	if err != nil {
		evidence.Owner = "unknown"
		evidence.Group = "unknown"
		return
	}

	// Get owner SID
	ownerSID, _, err := sd.Owner()
	if err != nil {
		evidence.Owner = "unknown"
	} else {
		// Try to lookup account name from SID
		account, domain, _, err := ownerSID.LookupAccount("")
		if err != nil {
			// If lookup fails, use SID string
			evidence.Owner = ownerSID.String()
		} else {
			if domain != "" {
				evidence.Owner = domain + "\\" + account
			} else {
				evidence.Owner = account
			}
		}
	}

	// Get group SID
	groupSID, _, err := sd.Group()
	if err != nil {
		evidence.Group = "unknown"
	} else {
		// Try to lookup group name from SID
		account, domain, _, err := groupSID.LookupAccount("")
		if err != nil {
			// If lookup fails, use SID string
			evidence.Group = groupSID.String()
		} else {
			if domain != "" {
				evidence.Group = domain + "\\" + account
			} else {
				evidence.Group = account
			}
		}
	}

	_ = fileInfo // Mark as used
}
