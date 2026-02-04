// Package video provides metadata extraction for video file formats including MP4, AVI, MKV, MOV, and WebM.
package video

import (
	"os"
	"strings"
	"time"

	"github.com/ilexum-group/evidex/internal/utils"
	"github.com/ilexum-group/evidex/pkg/models"
)

// WebMExtractor implements metadata extraction for WebM videos.
type WebMExtractor struct{}

// NewWebMExtractor creates a new WebM extractor.
func NewWebMExtractor() *WebMExtractor {
	return &WebMExtractor{}
}

// CanHandle checks if the file is a WebM video.
func (e *WebMExtractor) CanHandle(filePath string) bool {
	return IsWebM(filePath)
}

// Extract implements MetadataExtractor interface.
func (e *WebMExtractor) Extract(filePath string, logCmd models.CommandLogger) (interface{}, error) {
	return e.ExtractVideo(filePath, logCmd)
}

// GetType returns the extractor type.
func (e *WebMExtractor) GetType() string {
	return "video/webm"
}

// ExtractVideo extracts WebM-specific metadata.
func (e *WebMExtractor) ExtractVideo(filePath string, logCmd models.CommandLogger) (*models.VideoMetadata, error) {
	if logCmd != nil {
		startTime := time.Now()
		logCmd(utils.GenerateRandomID(), "video.ExtractWebM", []string{filePath}, startTime, time.Now(), 0, nil, "", filePath)
	}

	metadata := &models.VideoMetadata{
		Format: "WebM",
	}
	return metadata, nil
}

// IsWebM checks if a file is a WebM video by magic bytes.
func IsWebM(filePath string) bool {
	file, err := os.Open(filePath) // #nosec G304 - filepath is user-provided evidence path in forensic tool
	if err != nil {
		return false
	}
	defer func() {
		_ = file.Close()
	}()

	magic := make([]byte, 4)
	if _, err := file.Read(magic); err != nil {
		return false
	}

	// WebM uses Matroska container, signature: 1A 45 DF A3
	// Need to check for "webm" string in DocType element
	if len(magic) != 4 || magic[0] != 0x1A || magic[1] != 0x45 ||
		magic[2] != 0xDF || magic[3] != 0xA3 {
		return false
	}

	// Read more to find DocType
	if _, err := file.Seek(0, 0); err != nil {
		return false
	}
	header := make([]byte, 256)
	n, err := file.Read(header)
	if err != nil || n < 20 {
		return false
	}

	// Look for "webm" string in header
	headerStr := string(header[:n])
	return strings.Contains(headerStr, "webm")
}
