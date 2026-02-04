// Package images provides metadata extraction for image file formats including JPEG, PNG, and GIF.
package images

import (
	_ "image/gif" // Register GIF format
	"os"

	"github.com/ilexum-group/evidex/pkg/models"
)

// GIFExtractor implements metadata extraction for GIF images.
type GIFExtractor struct{}

// NewGIFExtractor creates a new GIF extractor.
func NewGIFExtractor() *GIFExtractor {
	return &GIFExtractor{}
}

// CanHandle checks if the file is a GIF image.
func (e *GIFExtractor) CanHandle(filePath string) bool {
	return IsGIF(filePath)
}

// Extract implements MetadataExtractor interface.
func (e *GIFExtractor) Extract(filePath string, logCmd models.CommandLogger) (interface{}, error) {
	return e.ExtractImage(filePath, logCmd)
}

// GetType returns the extractor type.
func (e *GIFExtractor) GetType() string {
	return "image/gif"
}

// ExtractImage extracts GIF-specific metadata.
func (e *GIFExtractor) ExtractImage(filePath string, logCmd models.CommandLogger) (*models.ImageMetadata, error) {
	metadata := &models.ImageMetadata{Format: "GIF"}

	// Get basic image dimensions
	width, height, err := extractImageDimensions(filePath, "GIF", logCmd)
	metadata.Width = width
	metadata.Height = height

	return metadata, err
}

// IsGIF checks if a file is a GIF image by magic bytes.
func IsGIF(filePath string) bool {
	file, err := os.Open(filePath) // #nosec G304 - filepath is user-provided evidence path in forensic tool
	if err != nil {
		return false
	}
	defer func() {
		_ = file.Close()
	}()

	magic := make([]byte, 6)
	if _, err := file.Read(magic); err != nil {
		return false
	}

	// GIF signature: "GIF87a" or "GIF89a"
	if len(magic) < 6 {
		return false
	}

	return (magic[0] == 'G' && magic[1] == 'I' && magic[2] == 'F' &&
		magic[3] == '8' && (magic[4] == '7' || magic[4] == '9') && magic[5] == 'a')
}
