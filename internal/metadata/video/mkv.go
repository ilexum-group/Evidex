package video

import (
	"os"

	"github.com/ilexum-group/evidex/pkg/models"
)

// MKVExtractor implements metadata extraction for Matroska MKV videos
type MKVExtractor struct{}

// NewMKVExtractor creates a new MKV extractor
func NewMKVExtractor() *MKVExtractor {
	return &MKVExtractor{}
}

// CanHandle checks if the file is an MKV video
func (e *MKVExtractor) CanHandle(filePath string) bool {
	return IsMKV(filePath)
}

// Extract implements MetadataExtractor interface
func (e *MKVExtractor) Extract(filePath string) (interface{}, error) {
	return e.ExtractVideo(filePath)
}

// GetType returns the extractor type
func (e *MKVExtractor) GetType() string {
	return "video/x-matroska"
}

// ExtractVideo extracts MKV-specific metadata
func (e *MKVExtractor) ExtractVideo(filePath string) (*models.VideoMetadata, error) {
	metadata := &models.VideoMetadata{
		Format: "Matroska",
	}
	return metadata, nil
}

// IsMKV checks if a file is a Matroska MKV video by magic bytes
func IsMKV(filePath string) bool {
	file, err := os.Open(filePath)
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

	// MKV/Matroska signature: 1A 45 DF A3
	return len(magic) == 4 &&
		magic[0] == 0x1A && magic[1] == 0x45 &&
		magic[2] == 0xDF && magic[3] == 0xA3
}
