package video

import (
	"os"

	"github.com/ilexum-group/evidex/pkg/models"
)

// MOVExtractor implements metadata extraction for QuickTime MOV videos
type MOVExtractor struct{}

// NewMOVExtractor creates a new MOV extractor
func NewMOVExtractor() *MOVExtractor {
	return &MOVExtractor{}
}

// CanHandle checks if the file is a MOV video
func (e *MOVExtractor) CanHandle(filePath string) bool {
	return IsMOV(filePath)
}

// Extract implements MetadataExtractor interface
func (e *MOVExtractor) Extract(filePath string) (interface{}, error) {
	return e.ExtractVideo(filePath)
}

// GetType returns the extractor type
func (e *MOVExtractor) GetType() string {
	return "video/quicktime"
}

// ExtractVideo extracts MOV-specific metadata
func (e *MOVExtractor) ExtractVideo(filePath string) (*models.VideoMetadata, error) {
	metadata := &models.VideoMetadata{
		Format: "QuickTime",
	}
	return metadata, nil
}

// IsMOV checks if a file is a QuickTime MOV video by magic bytes
func IsMOV(filePath string) bool {
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

	// MOV signature: ftyp with qt subtype
	return magic[4] == 'f' && magic[5] == 't' && magic[6] == 'y' && magic[7] == 'p' &&
		magic[8] == 'q' && magic[9] == 't'
}
