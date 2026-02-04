// Package generic provides metadata extraction for generic file formats including executables, archives, and documents.
package generic

import (
	"archive/zip"
	"fmt"
	"os"
	"time"

	"github.com/ilexum-group/evidex/internal/utils"
	"github.com/ilexum-group/evidex/pkg/models"
)

// ZIPExtractor implements metadata extraction for ZIP archives.
type ZIPExtractor struct{}

// NewZIPExtractor creates a new ZIP extractor.
func NewZIPExtractor() *ZIPExtractor {
	return &ZIPExtractor{}
}

// CanHandle checks if the file is a ZIP archive.
func (e *ZIPExtractor) CanHandle(filePath string) bool {
	return IsZIP(filePath)
}

// Extract implements MetadataExtractor interface.
func (e *ZIPExtractor) Extract(filePath string, logCmd models.CommandLogger) (interface{}, error) {
	return e.ExtractDocument(filePath, logCmd), nil
}

// GetType returns the extractor type.
func (e *ZIPExtractor) GetType() string {
	return "application/zip"
}

// ExtractDocument extracts ZIP metadata.
func (e *ZIPExtractor) ExtractDocument(filePath string, logCmd models.CommandLogger) map[string]string {
	return ParseZIP(filePath, logCmd)
}

// IsZIP checks if a file is a ZIP archive by magic bytes.
func IsZIP(filePath string) bool {
	file, err := os.Open(filePath) // #nosec G304 - filepath is user-provided evidence path in forensic tool
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

	// ZIP signature: PK\x03\x04 or PK\x05\x06 (empty archive) or PK\x07\x08 (spanned archive)
	return len(magic) == 4 && magic[0] == 'P' && magic[1] == 'K' &&
		(magic[2] == 0x03 && magic[3] == 0x04 ||
			magic[2] == 0x05 && magic[3] == 0x06 ||
			magic[2] == 0x07 && magic[3] == 0x08)
}

// ParseZIP extracts metadata from ZIP archives.
func ParseZIP(filePath string, logCmd models.CommandLogger) map[string]string {
	metadata := make(map[string]string)
	metadata["Type"] = "ZIP Archive"

	startTime := time.Now()
	reader, err := zip.OpenReader(filePath)
	endTime := time.Now()
	exitCode := 0
	if err != nil {
		exitCode = 1
	}
	if logCmd != nil {
		logCmd(utils.GenerateRandomID(), "zip.OpenReader", []string{filePath}, startTime, endTime, exitCode, err, "", filePath)
	}
	if err != nil {
		return metadata
	}
	defer func() {
		startTime := time.Now()
		err := reader.Close()
		endTime := time.Now()
		exitCode := 0
		if err != nil {
			exitCode = 1
		}
		if logCmd != nil {
			logCmd(utils.GenerateRandomID(), "zip.Reader.Close", []string{filePath}, startTime, endTime, exitCode, err, "", filePath)
		}
	}()

	metadata["FileCount"] = fmt.Sprintf("%d", len(reader.File))

	var totalUncompressed, totalCompressed int64
	hasEncrypted := false

	for _, file := range reader.File {
		totalUncompressed += int64(file.UncompressedSize64)
		totalCompressed += int64(file.CompressedSize64)

		// Record compression method
		if _, exists := metadata["CompressionMethod"]; !exists {
			switch file.Method {
			case zip.Store:
				metadata["CompressionMethod"] = "Store (no compression)"
			case zip.Deflate:
				metadata["CompressionMethod"] = "Deflate"
			default:
				metadata["CompressionMethod"] = fmt.Sprintf("Method %d", file.Method)
			}
		}
	}

	metadata["TotalUncompressedSize"] = fmt.Sprintf("%d bytes", totalUncompressed)
	metadata["TotalCompressedSize"] = fmt.Sprintf("%d bytes", totalCompressed)

	if totalUncompressed > 0 {
		compressionRatio := float64(totalCompressed) / float64(totalUncompressed) * 100
		metadata["CompressionRatio"] = fmt.Sprintf("%.2f%%", compressionRatio)
	}

	if hasEncrypted {
		metadata["Encrypted"] = "Yes"
	}

	if len(reader.Comment) > 0 {
		metadata["Comment"] = reader.Comment
	}

	return metadata
}
