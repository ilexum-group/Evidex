// Package utils provides utility functions for Evidex
package utils

import (
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/ilexum/evidex/internal/logger"
	"github.com/ilexum/evidex/internal/models"
)

// LogLevel represents the severity level of a log entry
type LogLevel string

const (
	LogLevelInfo    LogLevel = "INFO"
	LogLevelWarning LogLevel = "WARNING"
	LogLevelError   LogLevel = "ERROR"
)

var (
	logEntries []*models.LogEntry
	//nolint:unused
	logStartTime time.Time
)

// InitLogging initializes the logging system
func InitLogging() {
	logEntries = make([]*models.LogEntry, 0)
	logStartTime = time.Now()
}

// Log adds an entry to the log
func Log(level LogLevel, message, details, errMsg string) {
	entry := &models.LogEntry{
		Timestamp: time.Now(),
		Level:     string(level),
		Message:   message,
		Details:   details,
		Error:     errMsg,
	}
	logEntries = append(logEntries, entry)

	// Print to stderr for visibility
	fmt.Fprintf(os.Stderr, "[%s] %s: %s\n", entry.Level, entry.Timestamp.Format(time.RFC3339), message)
	if details != "" {
		fmt.Fprintf(os.Stderr, "  Details: %s\n", details)
	}
	if errMsg != "" {
		fmt.Fprintf(os.Stderr, "  Error: %s\n", errMsg)
	}
}

// LogInfo logs an informational message with RFC 5424 format
func LogInfo(message string, meta map[string]string) {
	logger.LogInfo(message, meta)
}

// LogWarn logs a warning message with RFC 5424 format
func LogWarn(message string, meta map[string]string) {
	logger.LogWarn(message, meta)
}

// LogError logs an error message with RFC 5424 format
func LogError(message string, meta map[string]string) {
	logger.LogError(message, meta)
}

// LogDebug logs a debug message with RFC 5424 format
func LogDebug(message string, meta map[string]string) {
	logger.LogDebug(message, meta)
}

// CalculateSHA256 calculates SHA-256 hash of a file
func CalculateSHA256(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer func() {
		_ = file.Close() //nolint:errcheck
	}()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("failed to hash file: %w", err)
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// CalculateSHA512 calculates SHA-512 hash of a file
func CalculateSHA512(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer func() {
		_ = file.Close() //nolint:errcheck
	}()

	hash := sha512.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("failed to hash file: %w", err)
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// CalculateMD5 calculates MD5 hash of a file (for compatibility only)
func CalculateMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer func() {
		_ = file.Close() //nolint:errcheck
	}()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("failed to hash file: %w", err)
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// VerifyHash verifies a file's hash against a known value
func VerifyHash(filePath, expectedHash string, algorithm string) (bool, error) {
	var actualHash string
	var err error

	switch algorithm {
	case "SHA-256":
		actualHash, err = CalculateSHA256(filePath)
	case "SHA-512":
		actualHash, err = CalculateSHA512(filePath)
	case "MD5":
		actualHash, err = CalculateMD5(filePath)
	default:
		return false, fmt.Errorf("unsupported hash algorithm: %s", algorithm)
	}

	if err != nil {
		return false, err
	}

	return actualHash == expectedHash, nil
}

// GetSystemContext gathers system information for chain of custody
func GetSystemContext() *models.SystemContext {
	hostname, _ := os.Hostname()
	currentUser, _ := GetCurrentUser()
	cwd, _ := os.Getwd()

	return &models.SystemContext{
		Hostname:         hostname,
		Username:         currentUser,
		OperatingSystem:  runtime.GOOS,
		OSVersion:        runtime.Version(),
		Architecture:     runtime.GOARCH,
		TimeZone:         time.Now().Format("MST"),
		LocalTime:        time.Now(),
		UTCTime:          time.Now().UTC(),
		BinaryVersion:    "1.0.0",
		WorkingDirectory: cwd,
	}
}

// GetCurrentUser retrieves the current executing user
func GetCurrentUser() (string, error) {
	switch runtime.GOOS {
	case "windows":
		return os.Getenv("USERNAME"), nil
	case "linux", "darwin":
		return os.Getenv("USER"), nil
	default:
		return "unknown", nil
	}
}

// IsReadOnly checks if a file can be opened in read-only mode
func IsReadOnly(filePath string) (bool, error) {
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0)
	if err != nil {
		if os.IsPermission(err) {
			return false, fmt.Errorf("permission denied")
		}
		return false, err
	}
	defer func() {
		_ = file.Close() //nolint:errcheck
	}()
	return true, nil
}

// GenerateEvidenceID creates a unique evidence package identifier
func GenerateEvidenceID() string {
	hostname, _ := os.Hostname()
	timestamp := time.Now().Format("20060102150405")
	return fmt.Sprintf("EVDX-%s-%s", hostname, timestamp)
}

// GetFileTypeFromExtension maps file extensions to MIME types
func GetFileTypeFromExtension(filePath string) string {
	ext := filepath.Ext(filePath)

	// Common multimedia types
	mimeTypes := map[string]string{
		// Images
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".bmp":  "image/bmp",
		".tiff": "image/tiff",
		".tif":  "image/tiff",
		".webp": "image/webp",
		".ico":  "image/x-icon",
		".svg":  "image/svg+xml",
		".raw":  "image/x-raw",
		".heic": "image/heic",
		".heif": "image/heif",

		// Videos
		".mp4":  "video/mp4",
		".avi":  "video/x-msvideo",
		".mov":  "video/quicktime",
		".mkv":  "video/x-matroska",
		".flv":  "video/x-flv",
		".wmv":  "video/x-ms-wmv",
		".webm": "video/webm",
		".m3u8": "application/vnd.apple.mpegurl",
		".ts":   "video/mp2t",
		".3gp":  "video/3gpp",
		".ogv":  "video/ogg",
	}

	if mimeType, exists := mimeTypes[ext]; exists {
		return mimeType
	}
	return "application/octet-stream"
}

// IsImageFile checks if a file is an image
func IsImageFile(filePath string) bool {
	mimeType := GetFileTypeFromExtension(filePath)
	return len(mimeType) > 0 && mimeType[:6] == "image/"
}

// IsVideoFile checks if a file is a video
func IsVideoFile(filePath string) bool {
	mimeType := GetFileTypeFromExtension(filePath)
	return len(mimeType) > 0 && mimeType[:6] == "video/"
}

// EnsureDirectory creates a directory if it doesn't exist
func EnsureDirectory(dirPath string) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return os.MkdirAll(dirPath, 0755)
	}
	return nil
}

// SafeCopyFile copies a file while maintaining read-only semantics
func SafeCopyFile(sourcePath, destinationPath string) error {
	// Verify source exists and is readable
	if _, err := os.Stat(sourcePath); err != nil {
		return fmt.Errorf("source file not accessible: %w", err)
	}

	// Open source
	src, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to open source: %w", err)
	}
	defer func() {
		_ = src.Close() //nolint:errcheck
	}()

	// Create destination
	dst, err := os.Create(destinationPath)
	if err != nil {
		return fmt.Errorf("failed to create destination: %w", err)
	}
	defer func() {
		_ = dst.Close() //nolint:errcheck
	}()

	// Copy contents
	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	return nil
}

// TrimPath removes leading separators for relative paths
func TrimPath(path string) string {
	if len(path) > 0 && (path[0] == '/' || path[0] == '\\') {
		return path[1:]
	}
	return path
}
