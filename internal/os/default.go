package os

import (
	stdos "os"
	"runtime"
	"time"

	"github.com/ilexum-group/evidex/internal/utils"
	"github.com/ilexum-group/evidex/pkg/models"
)

// OS defines the interface for OS-specific metadata extraction.
type OS interface {
	// Hostname returns the host name reported by the kernel.
	Hostname() (string, error)

	// Getenv retrieves the value of the environment variable named by the key.
	Getenv(key string) string

	// Getwd returns a rooted path name corresponding to the current directory.
	Getwd() (string, error)

	// Stat returns a FileInfo describing the named file.
	Stat(name string) (stdos.FileInfo, error)

	// GetCurrentUser retrieves the current executing user based on the OS.
	GetCurrentUser() (string, error)

	// GetProcessID returns the current process ID.
	GetProcessID() int

	// SetLogger sets the logger function for OS operations.
	SetLogger(logger models.CommandLogger)

	// ExtractAdvancedTimes extracts creation, access, and change times.
	ExtractAdvancedTimes(filePath string) (accessedTime time.Time, modifiedTime time.Time, changeTime time.Time, createdTime time.Time)

	// ExtractOwnershipInfo extracts file owner and group information.
	ExtractOwnershipInfo(filePath string) (username string, groupname string)
}

// DefaultImpl is the base implementation of Default with default methods.
type DefaultImpl struct {
	logger models.CommandLogger
}

// NewDefault creates a new default OS implementation.
func NewDefault() *DefaultImpl {
	return &DefaultImpl{}
}

// Hostname returns the host name reported by the kernel.
func (d *DefaultImpl) Hostname() (string, error) {
	startTime := time.Now()
	hostname, err := stdos.Hostname()
	endTime := time.Now()
	exitCode := 0
	if err != nil {
		exitCode = 1
	}
	d.logger(utils.GenerateRandomID(), "os.Hostname", []string{}, startTime, endTime, exitCode, err, "", "")
	return hostname, err
}

// Getenv retrieves the value of the environment variable named by the key.
func (d *DefaultImpl) Getenv(key string) string {
	startTime := time.Now()
	value := stdos.Getenv(key)
	endTime := time.Now()
	d.logger(utils.GenerateRandomID(), "os.Getenv", []string{key}, startTime, endTime, 0, nil, "", "")
	return value
}

// Getwd returns a rooted path name corresponding to the current directory.
func (d *DefaultImpl) Getwd() (string, error) {
	startTime := time.Now()
	path, err := stdos.Getwd()
	endTime := time.Now()
	exitCode := 0
	if err != nil {
		exitCode = 1
	}
	d.logger(utils.GenerateRandomID(), "os.Getwd", []string{}, startTime, endTime, exitCode, err, "", "")
	return path, err
}

// Stat returns a FileInfo describing the named file.
func (d *DefaultImpl) Stat(name string) (stdos.FileInfo, error) {
	startTime := time.Now()
	fileInfo, err := stdos.Stat(name)
	endTime := time.Now()
	exitCode := 0
	if err != nil {
		exitCode = 1
	}
	d.logger(utils.GenerateRandomID(), "os.Stat", []string{name}, startTime, endTime, exitCode, err, "", name)
	return fileInfo, err
}

// GetCurrentUser retrieves the current executing user (default implementation).
func (d *DefaultImpl) GetCurrentUser() (string, error) {
	return "unknown", nil
}

// GetProcessID returns the current process ID.
func (d *DefaultImpl) GetProcessID() int {
	startTime := time.Now()
	pid := stdos.Getpid()
	endTime := time.Now()
	d.logger(utils.GenerateRandomID(), "os.GetProcessID", nil, startTime, endTime, 0, nil, "", "")
	return pid
}

// SetLogger sets the logger function for OS operations.
func (d *DefaultImpl) SetLogger(logger models.CommandLogger) {
	d.logger = logger
}

// ExtractAdvancedTimes extracts creation, access, and change times (default implementation).
func (d *DefaultImpl) ExtractAdvancedTimes(filePath string) (accessedTime time.Time, modifiedTime time.Time, changeTime time.Time, createdTime time.Time) {
	startTime := time.Now()
	fileInfo, err := stdos.Stat(filePath)
	if err != nil {
		end := time.Now()
		if d.logger != nil {
			d.logger(utils.GenerateRandomID(), "os.Stat", []string{filePath}, startTime, end, 1, err, "", filePath)
		}
		return time.Time{}, time.Time{}, time.Time{}, time.Time{}
	}

	// Default implementation: only modified time is reliably available
	mTime := fileInfo.ModTime()
	modifiedTime = mTime
	// Use modified time as fallback for other times
	accessedTime = mTime
	changeTime = mTime
	createdTime = mTime

	endTime := time.Now()
	if d.logger != nil {
		d.logger(utils.GenerateRandomID(), "ExtractAdvancedTimes", []string{filePath}, startTime, endTime, 0, nil, "", filePath)
	}

	return accessedTime, modifiedTime, changeTime, createdTime
}

// ExtractOwnershipInfo extracts file owner and group information (default implementation).
func (d *DefaultImpl) ExtractOwnershipInfo(filePath string) (username string, groupname string) {
	startTime := time.Now()
	username = "unknown"
	groupname = "unknown"

	endTime := time.Now()
	if d.logger != nil {
		d.logger(utils.GenerateRandomID(), "ExtractOwnershipInfo", []string{filePath, username, groupname}, startTime, endTime, 0, nil, "", filePath)
	}

	return username, groupname
}

// New returns the appropriate OS implementation based on the runtime OS.
func New() OS {
	// Initialize with the appropriate OS-specific implementation
	switch runtime.GOOS {
	case "windows":
		return NewWindows()
	case "linux":
		return NewLinux()
	case "darwin":
		return NewDarwin()
	case "freebsd", "openbsd", "netbsd", "dragonfly":
		return NewUnix()
	default:
		return NewDefault()
	}
}
