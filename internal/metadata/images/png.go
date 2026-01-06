package images

import (
	"image"
	_ "image/png" // Register PNG format
	"os"

	"github.com/evidex/internal/models"
)

// PNGExtractor implements metadata extraction for PNG images
type PNGExtractor struct{}

// NewPNGExtractor creates a new PNG extractor
func NewPNGExtractor() *PNGExtractor {
	return &PNGExtractor{}
}

// CanHandle checks if the file is a PNG image
func (e *PNGExtractor) CanHandle(filePath string) bool {
	return IsPNG(filePath)
}

// Extract implements MetadataExtractor interface
func (e *PNGExtractor) Extract(filePath string) (interface{}, error) {
	return e.ExtractImage(filePath)
}

// GetType returns the extractor type
func (e *PNGExtractor) GetType() string {
	return "image/png"
}

// ExtractImage extracts PNG-specific metadata
func (e *PNGExtractor) ExtractImage(filePath string) (*models.ImageMetadata, error) {
	metadata := &models.ImageMetadata{
		Format: "PNG",
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

// IsPNG checks if a file is a PNG image by magic bytes
func IsPNG(filePath string) bool {
	file, err := os.Open(filePath)
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
