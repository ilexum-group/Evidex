package video

import (
	"os"
	"strings"

	"github.com/ilexum-group/evidex/internal/models"
)

// WebMExtractor implements metadata extraction for WebM videos
type WebMExtractor struct{}

// NewWebMExtractor creates a new WebM extractor
func NewWebMExtractor() *WebMExtractor {
	return &WebMExtractor{}
}

// CanHandle checks if the file is a WebM video
func (e *WebMExtractor) CanHandle(filePath string) bool {
	return IsWebM(filePath)
}

// Extract implements MetadataExtractor interface
func (e *WebMExtractor) Extract(filePath string) (interface{}, error) {
	return e.ExtractVideo(filePath)
}

// GetType returns the extractor type
func (e *WebMExtractor) GetType() string {
	return "video/webm"
}

// ExtractVideo extracts WebM-specific metadata
func (e *WebMExtractor) ExtractVideo(filePath string) (*models.VideoMetadata, error) {
	metadata := &models.VideoMetadata{
		Format: "WebM",
	}
	return metadata, nil
}

// IsWebM checks if a file is a WebM video by magic bytes
func IsWebM(filePath string) bool {
	file, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer func() {
		_ = file.Close()
	}()

	magic := make([]byte, 4)
	if _, err := file.Read(magic); err != nil {
		return false
	}

	// WebM uses Matroska container, signature: 1A 45 DF A3
	// Need to check for "webm" string in DocType element
	if len(magic) != 4 || magic[0] != 0x1A || magic[1] != 0x45 ||
		magic[2] != 0xDF || magic[3] != 0xA3 {
		return false
	}

	// Read more to find DocType
	if _, err := file.Seek(0, 0); err != nil {
		return false
	}
	header := make([]byte, 256)
	n, err := file.Read(header)
	if err != nil || n < 20 {
		return false
	}

	// Look for "webm" string in header
	headerStr := string(header[:n])
	return strings.Contains(headerStr, "webm")
}
