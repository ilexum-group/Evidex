// Package video provides metadata extraction for video file formats including MP4, AVI, MKV, MOV, and WebM.
package video

import (
	"os"
	"time"

	"github.com/ilexum-group/evidex/internal/utils"
	"github.com/ilexum-group/evidex/pkg/models"
)

// MOVExtractor implements metadata extraction for QuickTime MOV videos.
type MOVExtractor struct{}

// NewMOVExtractor creates a new MOV extractor.
func NewMOVExtractor() *MOVExtractor {
	return &MOVExtractor{}
}

// CanHandle checks if the file is a MOV video.
func (e *MOVExtractor) CanHandle(filePath string) bool {
	return IsMOV(filePath)
}

// Extract implements MetadataExtractor interface.
func (e *MOVExtractor) Extract(filePath string, logCmd models.CommandLogger) (interface{}, error) {
	return e.ExtractVideo(filePath, logCmd)
}

// GetType returns the extractor type.
func (e *MOVExtractor) GetType() string {
	return "video/quicktime"
}

// ExtractVideo extracts MOV-specific metadata.
func (e *MOVExtractor) ExtractVideo(filePath string, logCmd models.CommandLogger) (*models.VideoMetadata, error) {
	if logCmd != nil {
		startTime := time.Now()
		logCmd(utils.GenerateRandomID(), "video.ExtractMOV", []string{filePath}, startTime, time.Now(), 0, nil, "", filePath)
	}

	metadata := &models.VideoMetadata{
		Format: "QuickTime",
	}
	return metadata, nil
}

// IsMOV checks if a file is a QuickTime MOV video by magic bytes.
func IsMOV(filePath string) bool {
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

	// MOV signature: ftyp with qt subtype
	return magic[4] == 'f' && magic[5] == 't' && magic[6] == 'y' && magic[7] == 'p' &&
		magic[8] == 'q' && magic[9] == 't'
}
