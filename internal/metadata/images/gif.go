package images

import (
	"image"
	_ "image/gif" // Register GIF format
	"os"

	"github.com/ilexum/evidex/internal/models"
)

// GIFExtractor implements metadata extraction for GIF images
type GIFExtractor struct{}

// NewGIFExtractor creates a new GIF extractor
func NewGIFExtractor() *GIFExtractor {
	return &GIFExtractor{}
}

// CanHandle checks if the file is a GIF image
func (e *GIFExtractor) CanHandle(filePath string) bool {
	return IsGIF(filePath)
}

// Extract implements MetadataExtractor interface
func (e *GIFExtractor) Extract(filePath string) (interface{}, error) {
	return e.ExtractImage(filePath)
}

// GetType returns the extractor type
func (e *GIFExtractor) GetType() string {
	return "image/gif"
}

// ExtractImage extracts GIF-specific metadata
func (e *GIFExtractor) ExtractImage(filePath string) (*models.ImageMetadata, error) {
	metadata := &models.ImageMetadata{Format: "GIF"}

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

// IsGIF checks if a file is a GIF image by magic bytes
func IsGIF(filePath string) bool {
	file, err := os.Open(filePath)
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
