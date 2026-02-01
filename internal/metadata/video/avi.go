package video

import (
	"os"

	"github.com/ilexum/evidex/internal/models"
)

// AVIExtractor implements metadata extraction for AVI videos
type AVIExtractor struct{}

// NewAVIExtractor creates a new AVI extractor
func NewAVIExtractor() *AVIExtractor {
	return &AVIExtractor{}
}

// CanHandle checks if the file is an AVI video
func (e *AVIExtractor) CanHandle(filePath string) bool {
	return IsAVI(filePath)
}

// Extract implements MetadataExtractor interface
func (e *AVIExtractor) Extract(filePath string) (interface{}, error) {
	return e.ExtractVideo(filePath)
}

// GetType returns the extractor type
func (e *AVIExtractor) GetType() string {
	return "video/x-msvideo"
}

// ExtractVideo extracts AVI-specific metadata
func (e *AVIExtractor) ExtractVideo(filePath string) (*models.VideoMetadata, error) {
	metadata := &models.VideoMetadata{
		Format: "AVI",
	}
	return metadata, nil
}

// IsAVI checks if a file is an AVI video by magic bytes
func IsAVI(filePath string) bool {
	file, err := os.Open(filePath)
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
