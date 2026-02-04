// Package generic provides metadata extraction for generic file formats including executables, archives, and documents.
package generic

import (
	"compress/gzip"
	"fmt"
	"os"
	"time"

	"github.com/ilexum-group/evidex/internal/utils"
	"github.com/ilexum-group/evidex/pkg/models"
)

// GZIPExtractor implements metadata extraction for GZIP archives.
type GZIPExtractor struct{}

// NewGZIPExtractor creates a new GZIP extractor.
func NewGZIPExtractor() *GZIPExtractor {
	return &GZIPExtractor{}
}

// CanHandle checks if the file is a GZIP archive.
func (e *GZIPExtractor) CanHandle(filePath string) bool {
	return IsGZIP(filePath)
}

// Extract implements MetadataExtractor interface.
func (e *GZIPExtractor) Extract(filePath string, logCmd models.CommandLogger) (interface{}, error) {
	return e.ExtractDocument(filePath, logCmd), nil
}

// GetType returns the extractor type.
func (e *GZIPExtractor) GetType() string {
	return "application/gzip"
}

// ExtractDocument extracts GZIP metadata.
func (e *GZIPExtractor) ExtractDocument(filePath string, logCmd models.CommandLogger) map[string]string {
	return ParseGZIP(filePath, logCmd)
}

// IsGZIP checks if a file is a GZIP archive by magic bytes.
func IsGZIP(filePath string) bool {
	file, err := os.Open(filePath) // #nosec G304 - filepath is user-provided evidence path in forensic tool
	if err != nil {
		return false
	}
	defer func() {
		_ = file.Close()
	}()

	magic := make([]byte, 2)
	if _, err := file.Read(magic); err != nil {
		return false
	}

	// GZIP signature: 1F 8B
	return len(magic) == 2 && magic[0] == 0x1F && magic[1] == 0x8B
}

// ParseGZIP extracts metadata from GZIP archives.
func ParseGZIP(filePath string, logCmd models.CommandLogger) map[string]string {
	metadata := make(map[string]string)
	metadata["Type"] = "GZIP Archive"

	startTime := time.Now()
	file, err := os.Open(filePath) // #nosec G304 - filepath is user-provided evidence path in forensic tool
	endTime := time.Now()
	exitCode := 0
	if err != nil {
		exitCode = 1
	}
	if logCmd != nil {
		logCmd(utils.GenerateRandomID(), "os.Open", []string{filePath, "GZIP parse"}, startTime, endTime, exitCode, err, "", filePath)
	}
	if err != nil {
		return metadata
	}
	defer func() {
		startTime := time.Now()
		err := file.Close()
		endTime := time.Now()
		exitCode := 0
		if err != nil {
			exitCode = 1
		}
		if logCmd != nil {
			logCmd(utils.GenerateRandomID(), "file.Close", []string{filePath, "GZIP parse"}, startTime, endTime, exitCode, err, "", filePath)
		}
	}()

	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return metadata
	}
	defer func() {
		startTime := time.Now()
		err := gzReader.Close()
		endTime := time.Now()
		exitCode := 0
		if err != nil {
			exitCode = 1
		}
		if logCmd != nil {
			logCmd(utils.GenerateRandomID(), "gzip.Reader.Close", []string{filePath}, startTime, endTime, exitCode, err, "", filePath)
		}
	}()

	// Extract GZIP header information
	if gzReader.Name != "" {
		metadata["OriginalFilename"] = gzReader.Name
	}

	if gzReader.Comment != "" {
		metadata["Comment"] = gzReader.Comment
	}

	if !gzReader.ModTime.IsZero() {
		metadata["ModificationTime"] = gzReader.ModTime.String()
	}

	// OS type
	switch gzReader.OS {
	case 0:
		metadata["OS"] = "FAT filesystem (MS-DOS, OS/2, NT/Win32)"
	case 1:
		metadata["OS"] = "Amiga"
	case 2:
		metadata["OS"] = "VMS (or OpenVMS)"
	case 3:
		metadata["OS"] = "Unix"
	case 4:
		metadata["OS"] = "VM/CMS"
	case 5:
		metadata["OS"] = "Atari TOS"
	case 6:
		metadata["OS"] = "HPFS filesystem (OS/2, NT)"
	case 7:
		metadata["OS"] = "Macintosh"
	case 8:
		metadata["OS"] = "Z-System"
	case 9:
		metadata["OS"] = "CP/M"
	case 10:
		metadata["OS"] = "TOPS-20"
	case 11:
		metadata["OS"] = "NTFS filesystem (NT)"
	case 12:
		metadata["OS"] = "QDOS"
	case 13:
		metadata["OS"] = "Acorn RISCOS"
	default:
		metadata["OS"] = fmt.Sprintf("Unknown (%d)", gzReader.OS)
	}

	return metadata
}
