package video

import (
	"os"

	"github.com/ilexum-group/evidex/internal/models"
)

// MP4Extractor implements metadata extraction for MP4 videos
type MP4Extractor struct{}

// NewMP4Extractor creates a new MP4 extractor
func NewMP4Extractor() *MP4Extractor {
	return &MP4Extractor{}
}

// CanHandle checks if the file is an MP4 video
func (e *MP4Extractor) CanHandle(filePath string) bool {
	return IsMP4(filePath)
}

// Extract implements MetadataExtractor interface
func (e *MP4Extractor) Extract(filePath string) (interface{}, error) {
	return e.ExtractVideo(filePath)
}

// GetType returns the extractor type
func (e *MP4Extractor) GetType() string {
	return "video/mp4"
}

// ExtractVideo extracts MP4-specific metadata
func (e *MP4Extractor) ExtractVideo(filePath string) (*models.VideoMetadata, error) {
	metadata := &models.VideoMetadata{
		Format: "MP4",
	}

	return metadata, nil
}

// IsMP4 checks if a file is an MP4 video by magic bytes
func IsMP4(filePath string) bool {
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

	// MP4 signature: ftyp at offset 4
	return magic[4] == 'f' && magic[5] == 't' && magic[6] == 'y' && magic[7] == 'p'
}
