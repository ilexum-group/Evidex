// Package images provides metadata extraction for image file formats including JPEG, PNG, and GIF.
package images

import (
	_ "image/png" // Register PNG format
	"os"

	"github.com/ilexum-group/evidex/pkg/models"
)

// PNGExtractor implements metadata extraction for PNG images.
type PNGExtractor struct{}

// NewPNGExtractor creates a new PNG extractor.
func NewPNGExtractor() *PNGExtractor {
	return &PNGExtractor{}
}

// CanHandle checks if the file is a PNG image.
func (e *PNGExtractor) CanHandle(filePath string) bool {
	return IsPNG(filePath)
}

// Extract implements MetadataExtractor interface.
func (e *PNGExtractor) Extract(filePath string, logCmd models.CommandLogger) (interface{}, error) {
	return e.ExtractImage(filePath, logCmd)
}

// GetType returns the extractor type.
func (e *PNGExtractor) GetType() string {
	return "image/png"
}

// ExtractImage extracts PNG-specific metadata.
func (e *PNGExtractor) ExtractImage(filePath string, logCmd models.CommandLogger) (*models.ImageMetadata, error) {
	metadata := &models.ImageMetadata{Format: "PNG"}

	// Get basic image dimensions
	width, height, err := extractImageDimensions(filePath, "PNG", logCmd)
	metadata.Width = width
	metadata.Height = height

	return metadata, err
}

// IsPNG checks if a file is a PNG image by magic bytes.
func IsPNG(filePath string) bool {
	file, err := os.Open(filePath) // #nosec G304 - filepath is user-provided evidence path in forensic tool
	if err != nil {
		return false
	}
	defer func() {
		_ = file.Close()
	}()

	magic := make([]byte, 8)
	if _, err := file.Read(magic); err != nil {
		return false
	}

	// PNG signature: 89 50 4E 47 0D 0A 1A 0A
	return len(magic) == 8 &&
		magic[0] == 0x89 && magic[1] == 0x50 &&
		magic[2] == 0x4E && magic[3] == 0x47 &&
		magic[4] == 0x0D && magic[5] == 0x0A &&
		magic[6] == 0x1A && magic[7] == 0x0A
}
