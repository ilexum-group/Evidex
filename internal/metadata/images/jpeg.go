// Package images provides metadata extraction for image file formats including JPEG, PNG, and GIF.
package images

import (
	_ "image/jpeg" // Register JPEG format
	"os"

	"github.com/ilexum-group/evidex/pkg/models"
)

// JPEGExtractor implements metadata extraction for JPEG images.
type JPEGExtractor struct{}

// NewJPEGExtractor creates a new JPEG extractor.
func NewJPEGExtractor() *JPEGExtractor {
	return &JPEGExtractor{}
}

// CanHandle checks if the file is a JPEG image.
func (e *JPEGExtractor) CanHandle(filePath string) bool {
	return IsJPEG(filePath)
}

// Extract implements MetadataExtractor interface.
func (e *JPEGExtractor) Extract(filePath string, logCmd models.CommandLogger) (interface{}, error) {
	return e.ExtractImage(filePath, logCmd)
}

// GetType returns the extractor type.
func (e *JPEGExtractor) GetType() string {
	return "image/jpeg"
}

// ExtractImage extracts JPEG-specific metadata.
func (e *JPEGExtractor) ExtractImage(filePath string, logCmd models.CommandLogger) (*models.ImageMetadata, error) {
	metadata := &models.ImageMetadata{Format: "JPEG"}

	// Get basic image dimensions
	width, height, err := extractImageDimensions(filePath, "JPEG", logCmd)
	metadata.Width = width
	metadata.Height = height

	return metadata, err
}

// IsJPEG checks if a file is a JPEG image by magic bytes.
func IsJPEG(filePath string) bool {
	file, err := os.Open(filePath) // #nosec G304 - filepath is user-provided evidence path in forensic tool
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
