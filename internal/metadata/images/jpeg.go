package images

import (
	"image"
	_ "image/jpeg" // Register JPEG format
	"os"

	"github.com/evidex/internal/models"
)

// JPEGExtractor implements metadata extraction for JPEG images
type JPEGExtractor struct{}

// NewJPEGExtractor creates a new JPEG extractor
func NewJPEGExtractor() *JPEGExtractor {
	return &JPEGExtractor{}
}

// CanHandle checks if the file is a JPEG image
func (e *JPEGExtractor) CanHandle(filePath string) bool {
	return IsJPEG(filePath)
}

// Extract implements MetadataExtractor interface
func (e *JPEGExtractor) Extract(filePath string) (interface{}, error) {
	return e.ExtractImage(filePath)
}

// GetType returns the extractor type
func (e *JPEGExtractor) GetType() string {
	return "image/jpeg"
}

// ExtractImage extracts JPEG-specific metadata
func (e *JPEGExtractor) ExtractImage(filePath string) (*models.ImageMetadata, error) {
	metadata := &models.ImageMetadata{
		Format: "JPEG",
		EXIF:   make(map[string]string),
		XMP:    make(map[string]string),
		IPTC:   make(map[string]string),
	}

	// Get basic image dimensions
	file, err := os.Open(filePath)
	if err != nil {
		return metadata, err
	}
	defer func() {
		_ = file.Close()
	}()

	config, _, err := image.DecodeConfig(file)
	if err == nil {
		metadata.Width = config.Width
		metadata.Height = config.Height
	}

	return metadata, nil
}

// IsJPEG checks if a file is a JPEG image by magic bytes
func IsJPEG(filePath string) bool {
	file, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer func() {
		_ = file.Close()
	}()

	magic := make([]byte, 3)
	if _, err := file.Read(magic); err != nil {
		return false
	}

	// JPEG signature: FF D8 FF
	return len(magic) == 3 && magic[0] == 0xFF && magic[1] == 0xD8 && magic[2] == 0xFF
}
