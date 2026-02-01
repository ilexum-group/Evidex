// Package metadata provides extraction of forensically relevant file metadata
// This package demonstrates polymorphism and interface-based design for extensible metadata extraction
package metadata

import (
	"fmt"
	"os"
	"time"

	"github.com/ilexum-group/evidex/internal/metadata/generic"
	"github.com/ilexum-group/evidex/internal/metadata/images"
	osmetadata "github.com/ilexum-group/evidex/internal/metadata/os"
	"github.com/ilexum-group/evidex/internal/metadata/video"
	"github.com/ilexum-group/evidex/internal/models"
	"github.com/ilexum-group/evidex/internal/utils"
)

// Global extractor registry initialized once
var globalRegistry *ExtractorRegistry

// init initializes the global extractor registry with all available extractors
// This demonstrates polymorphism: each extractor implements the MetadataExtractor interface
func init() {
	globalRegistry = NewExtractorRegistry()

	// Register image extractors
	globalRegistry.Register(images.NewJPEGExtractor())
	globalRegistry.Register(images.NewPNGExtractor())
	globalRegistry.Register(images.NewGIFExtractor())

	// Register video extractors
	globalRegistry.Register(video.NewMP4Extractor())
	globalRegistry.Register(video.NewMOVExtractor())
	globalRegistry.Register(video.NewMKVExtractor())
	globalRegistry.Register(video.NewAVIExtractor())
	globalRegistry.Register(video.NewWebMExtractor())

	// Register document/archive/executable extractors
	globalRegistry.Register(generic.NewPDFExtractor())
	globalRegistry.Register(generic.NewZIPExtractor())
	globalRegistry.Register(generic.NewGZIPExtractor())
	globalRegistry.Register(generic.NewTARExtractor())
	globalRegistry.Register(generic.NewPEExtractor())
	globalRegistry.Register(generic.NewELFExtractor())
	globalRegistry.Register(generic.NewMachOExtractor())

	// Text extractor as fallback (register last)
	globalRegistry.Register(generic.NewTextExtractor())
}

// ExtractFileMetadata extracts filesystem and file metadata using polymorphic extractors
func ExtractFileMetadata(filePath string) (*models.FileEvidence, error) {
	// Get file info
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}

	evidence := &models.FileEvidence{
		SourcePath:      filePath,
		Filename:        fileInfo.Name(),
		FileSize:        fileInfo.Size(),
		FileMode:        uint32(fileInfo.Mode()),
		ModifiedTime:    fileInfo.ModTime(),
		AcquisitionTime: time.Now(),
		FileType:        utils.GetFileTypeFromExtension(filePath),
		Hashes:          &models.FileHashes{},
	}

	// Try to get more detailed timestamps (OS-specific)
	// Uses polymorphic OS metadata extractor
	osmetadata.Current.ExtractAdvancedTimes(filePath, evidence)

	// Get owner and group information (OS-specific)
	osmetadata.Current.ExtractOwnershipInfo(filePath, evidence)

	// Check if it's a symlink
	stat, err := os.Lstat(filePath)
	if err == nil {
		if stat.Mode()&os.ModeSymlink != 0 {
			evidence.IsSymlink = true
			if target, err := os.Readlink(filePath); err == nil {
				evidence.SymlinkTarget = target
			}
		}
	}

	// Use polymorphic extraction: find the right extractor and extract metadata
	// This is the key polymorphism pattern - the registry dispatches to the correct extractor
	extractor := globalRegistry.FindExtractor(filePath)
	if extractor != nil {
		result, err := extractor.Extract(filePath)
		if err == nil && result != nil {
			// Type assertion based on the result type
			switch v := result.(type) {
			case *models.ImageMetadata:
				evidence.ImageMetadata = v
			case *models.VideoMetadata:
				evidence.VideoMetadata = v
			case map[string]string:
				evidence.GenericMetadata = v
			}
		}
	}

	return evidence, nil
}

// CalculateFileHashes computes cryptographic hashes of a file
func CalculateFileHashes(filePath string) (*models.FileHashes, error) {
	hashes := &models.FileHashes{}

	// Calculate SHA-256
	sha256Hash, err := utils.CalculateSHA256(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate SHA-256: %w", err)
	}
	hashes.SHA256 = sha256Hash

	// Calculate SHA-512
	sha512Hash, err := utils.CalculateSHA512(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate SHA-512: %w", err)
	}
	hashes.SHA512 = sha512Hash

	// Calculate MD5
	md5Hash, err := utils.CalculateMD5(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate MD5: %w", err)
	}
	hashes.MD5 = md5Hash

	return hashes, nil
}

// VerifyFileIntegrity verifies a file against known hashes
func VerifyFileIntegrity(filePath string, expectedHashes *models.FileHashes) error {
	actualHashes, err := CalculateFileHashes(filePath)
	if err != nil {
		return fmt.Errorf("failed to calculate hashes: %w", err)
	}

	if expectedHashes.SHA256 != "" && actualHashes.SHA256 != expectedHashes.SHA256 {
		return fmt.Errorf("SHA-256 mismatch: expected %s, got %s", expectedHashes.SHA256, actualHashes.SHA256)
	}

	if expectedHashes.SHA512 != "" && actualHashes.SHA512 != expectedHashes.SHA512 {
		return fmt.Errorf("SHA-512 mismatch: expected %s, got %s", expectedHashes.SHA512, actualHashes.SHA512)
	}

	if expectedHashes.MD5 != "" && actualHashes.MD5 != expectedHashes.MD5 {
		return fmt.Errorf("MD5 mismatch: expected %s, got %s", expectedHashes.MD5, actualHashes.MD5)
	}

	return nil
}
