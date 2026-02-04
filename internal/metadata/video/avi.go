// Package video provides metadata extraction for video file formats including MP4, AVI, MKV, MOV, and WebM.
package video

import (
	"os"
	"time"

	"github.com/ilexum-group/evidex/internal/utils"
	"github.com/ilexum-group/evidex/pkg/models"
)

// AVIExtractor implements metadata extraction for AVI videos.
type AVIExtractor struct{}

// NewAVIExtractor creates a new AVI extractor.
func NewAVIExtractor() *AVIExtractor {
	return &AVIExtractor{}
}

// CanHandle checks if the file is an AVI video.
func (e *AVIExtractor) CanHandle(filePath string) bool {
	return IsAVI(filePath)
}

// Extract implements MetadataExtractor interface.
func (e *AVIExtractor) Extract(filePath string, logCmd models.CommandLogger) (interface{}, error) {
	return e.ExtractVideo(filePath, logCmd)
}

// GetType returns the extractor type.
func (e *AVIExtractor) GetType() string {
	return "video/x-msvideo"
}

// ExtractVideo extracts AVI-specific metadata.
func (e *AVIExtractor) ExtractVideo(filePath string, logCmd models.CommandLogger) (*models.VideoMetadata, error) {
	if logCmd != nil {
		startTime := time.Now()
		logCmd(utils.GenerateRandomID(), "video.ExtractAVI", []string{filePath}, startTime, time.Now(), 0, nil, "", filePath)
	}

	metadata := &models.VideoMetadata{
		Format: "AVI",
	}
	return metadata, nil
}

// IsAVI checks if a file is an AVI video by magic bytes.
func IsAVI(filePath string) bool {
	file, err := os.Open(filePath) // #nosec G304 - filepath is user-provided evidence path in forensic tool
	if err != nil {
		return false
	}
	defer func() {
		_ = file.Close()
	}()

	magic := make([]byte, 12)
	if _, err := file.Read(magic); err != nil {
		return false
	}

	if len(magic) < 12 {
		return false
	}

	// AVI signature: RIFF....AVI
	return magic[0] == 'R' && magic[1] == 'I' && magic[2] == 'F' && magic[3] == 'F' &&
		magic[8] == 'A' && magic[9] == 'V' && magic[10] == 'I'
}
