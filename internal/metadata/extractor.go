// Package metadata provides a polymorphic interface for metadata extraction
package metadata

import (
	"github.com/ilexum-group/evidex/pkg/models"
)

// MetadataExtractor defines the interface for all metadata extractors.
// This allows for polymorphic behavior across different file types.
type MetadataExtractor interface {
	// CanHandle determines if this extractor can process the given file.
	CanHandle(filePath string) bool

	// Extract performs the metadata extraction and returns structured data.
	Extract(filePath string, commandLogger models.CommandLogger) (interface{}, error)

	// GetType returns the type of metadata this extractor produces.
	GetType() string
}

// ImageExtractor is the interface for image metadata extraction.
type ImageExtractor interface {
	MetadataExtractor
	ExtractImage(filePath string, commandLogger models.CommandLogger) (*models.ImageMetadata, error)
}

// VideoExtractor is the interface for video metadata extraction.
type VideoExtractor interface {
	MetadataExtractor
	ExtractVideo(filePath string, commandLogger models.CommandLogger) (*models.VideoMetadata, error)
}

// DocumentExtractor is the interface for document metadata extraction.
type DocumentExtractor interface {
	MetadataExtractor
	ExtractDocument(filePath string, commandLogger models.CommandLogger) map[string]string
}

// ExtractorRegistry manages all registered extractors.
type ExtractorRegistry struct {
	extractors []MetadataExtractor
}

// NewExtractorRegistry creates a new registry with all extractors.
func NewExtractorRegistry() *ExtractorRegistry {
	return &ExtractorRegistry{
		extractors: make([]MetadataExtractor, 0),
	}
}

// Register adds an extractor to the registry.
func (r *ExtractorRegistry) Register(extractor MetadataExtractor) {
	r.extractors = append(r.extractors, extractor)
}

// FindExtractor finds the appropriate extractor for a file.
func (r *ExtractorRegistry) FindExtractor(filePath string) MetadataExtractor {
	for _, extractor := range r.extractors {
		if extractor.CanHandle(filePath) {
			return extractor
		}
	}
	return nil
}
