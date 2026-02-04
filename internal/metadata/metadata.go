// Package metadata provides extraction of forensically relevant file metadata.
// This package demonstrates polymorphism and interface-based design for extensible metadata extraction
package metadata

import (
	"fmt"

	"github.com/ilexum-group/evidex/internal/metadata/generic"
	"github.com/ilexum-group/evidex/internal/metadata/images"
	"github.com/ilexum-group/evidex/internal/metadata/video"
	"github.com/ilexum-group/evidex/internal/utils"
	"github.com/ilexum-group/evidex/pkg/models"
)

// MetadataManager manages metadata extraction with polymorphic extractors.
type MetadataManager struct {
	registry      *ExtractorRegistry
	commandLogger models.CommandLogger
}

// NewMetadataManager creates a new metadata manager with all registered extractors.
// This demonstrates polymorphism: each extractor implements the MetadataExtractor interface.
func NewMetadataManager(commandLogger models.CommandLogger) *MetadataManager {
	registry := NewExtractorRegistry()

	// Register image extractors
	registry.Register(images.NewJPEGExtractor())
	registry.Register(images.NewPNGExtractor())
	registry.Register(images.NewGIFExtractor())

	// Register video extractors
	registry.Register(video.NewMP4Extractor())
	registry.Register(video.NewMOVExtractor())
	registry.Register(video.NewMKVExtractor())
	registry.Register(video.NewAVIExtractor())
	registry.Register(video.NewWebMExtractor())

	// Register document/archive/executable extractors
	registry.Register(generic.NewPDFExtractor())
	registry.Register(generic.NewZIPExtractor())
	registry.Register(generic.NewGZIPExtractor())
	registry.Register(generic.NewTARExtractor())
	registry.Register(generic.NewPEExtractor())
	registry.Register(generic.NewELFExtractor())
	registry.Register(generic.NewMachOExtractor())

	// Text extractor as fallback (register last)
	registry.Register(generic.NewTextExtractor())

	return &MetadataManager{
		registry:      registry,
		commandLogger: commandLogger,
	}
}

// ExtractFileMetadata extracts filesystem and file metadata using polymorphic extractors.
func (m *MetadataManager) ExtractFileMetadata(fileEvidence *models.FileEvidence) error {
	// Use polymorphic extraction: find the right extractor and extract metadata
	// This is the key polymorphism pattern - the registry dispatches to the correct extractor
	extractor := m.registry.FindExtractor(fileEvidence.SourcePath)
	if extractor != nil {
		result, err := extractor.Extract(fileEvidence.SourcePath, m.commandLogger)
		if err == nil && result != nil {
			// Type assertion based on the result type
			switch v := result.(type) {
			case *models.ImageMetadata:
				fileEvidence.ImageMetadata = v
			case *models.VideoMetadata:
				fileEvidence.VideoMetadata = v
			case map[string]string:
				fileEvidence.GenericMetadata = v
			}
		}
	}

	return nil
}

// CalculateFileHashes calculates all cryptographic hashes for a file.
func (m *MetadataManager) CalculateFileHashes(filePath string) (*models.FileHashes, error) {
	hashes := &models.FileHashes{}

	// Calculate MD5
	md5Hash, err := utils.CalculateMD5(filePath)
	if err != nil {
		return nil, fmt.Errorf("MD5 calculation failed: %w", err)
	}
	hashes.MD5 = md5Hash

	// Calculate SHA1
	sha1Hash, err := utils.CalculateSHA1(filePath)
	if err != nil {
		return nil, fmt.Errorf("SHA1 calculation failed: %w", err)
	}
	hashes.SHA1 = sha1Hash

	// Calculate SHA256
	sha256Hash, err := utils.CalculateSHA256(filePath)
	if err != nil {
		return nil, fmt.Errorf("SHA256 calculation failed: %w", err)
	}
	hashes.SHA256 = sha256Hash

	// Calculate SHA512
	sha512Hash, err := utils.CalculateSHA512(filePath)
	if err != nil {
		return nil, fmt.Errorf("SHA512 calculation failed: %w", err)
	}
	hashes.SHA512 = sha512Hash

	return hashes, nil
}

// GetFileTypeFromMimeType determines the file type based on mimetype using registered extractors.
func (m *MetadataManager) GetFileTypeFromMimeType(filePath string) string {
	// Find the appropriate extractor for this file
	extractor := m.registry.FindExtractor(filePath)
	if extractor != nil {
		// Return the mimetype from the extractor
		return extractor.GetType()
	}

	// Default fallback for unknown types
	return "application/octet-stream"
}
