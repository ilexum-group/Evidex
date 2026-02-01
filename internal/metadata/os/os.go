package os

import (
	"github.com/ilexum/evidex/internal/models"
)

// OSMetadataExtractor defines the interface for OS-specific metadata extraction
type OSMetadataExtractor interface {
	// ExtractAdvancedTimes extracts creation, access, and change times
	ExtractAdvancedTimes(filePath string, evidence *models.FileEvidence)

	// ExtractOwnershipInfo extracts file owner and group information
	ExtractOwnershipInfo(filePath string, evidence *models.FileEvidence)

	// GetOSName returns the name of the operating system
	GetOSName() string
}

// Current holds the OS-specific metadata extractor for the current platform
var Current OSMetadataExtractor

func init() {
	// This will be set by the platform-specific init() functions
	Current = newExtractor()
}
