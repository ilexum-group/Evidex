// Package generic provides metadata extraction for generic file formats including executables, archives, and documents.
package generic

import (
	"archive/tar"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/ilexum-group/evidex/internal/utils"
	"github.com/ilexum-group/evidex/pkg/models"
)

// TARExtractor implements metadata extraction for TAR archives.
type TARExtractor struct{}

// NewTARExtractor creates a new TAR extractor.
func NewTARExtractor() *TARExtractor {
	return &TARExtractor{}
}

// CanHandle checks if the file is a TAR archive.
func (e *TARExtractor) CanHandle(filePath string) bool {
	return IsTAR(filePath)
}

// Extract implements MetadataExtractor interface.
func (e *TARExtractor) Extract(filePath string, logCmd models.CommandLogger) (interface{}, error) {
	return e.ExtractDocument(filePath, logCmd), nil
}

// GetType returns the extractor type.
func (e *TARExtractor) GetType() string {
	return "application/x-tar"
}

// ExtractDocument extracts TAR metadata.
func (e *TARExtractor) ExtractDocument(filePath string, logCmd models.CommandLogger) map[string]string {
	return ParseTAR(filePath, logCmd)
}

// IsTAR checks if a file is a TAR archive by magic bytes.
func IsTAR(filePath string) bool {
	file, err := os.Open(filePath) // #nosec G304 - filepath is user-provided evidence path in forensic tool
	if err != nil {
		return false
	}
	defer func() {
		_ = file.Close()
	}()

	// TAR signature is at offset 257: "ustar"
	magic := make([]byte, 262)
	if n, err := file.Read(magic); err != nil || n < 262 {
		return false
	}

	return string(magic[257:262]) == "ustar"
}

// ParseTAR extracts metadata from TAR archives.
func ParseTAR(filePath string, logCmd models.CommandLogger) map[string]string {
	metadata := make(map[string]string)
	metadata["type"] = "TAR Archive"

	startTime := time.Now()
	file, err := os.Open(filePath) // #nosec G304 - filepath is user-provided evidence path in forensic tool
	endTime := time.Now()
	exitCode := 0
	if err != nil {
		exitCode = 1
	}
	if logCmd != nil {
		logCmd(utils.GenerateRandomID(), "os.Open", []string{filePath, "TAR parse"}, startTime, endTime, exitCode, err, "", filePath)
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
			logCmd(utils.GenerateRandomID(), "file.Close", []string{filePath, "TAR parse"}, startTime, endTime, exitCode, err, "", filePath)
		}
	}()

	tarReader := tar.NewReader(file)

	fileCount := 0
	dirCount := 0
	var totalSize int64

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			break
		}

		switch header.Typeflag {
		case tar.TypeDir:
			dirCount++
		case tar.TypeReg:
			fileCount++
			totalSize += header.Size
		}

		// Record the first owner/group found
		if _, exists := metadata["Owner"]; !exists && header.Uname != "" {
			metadata["Owner"] = header.Uname
		}
		if _, exists := metadata["Group"]; !exists && header.Gname != "" {
			metadata["Group"] = header.Gname
		}
	}

	metadata["FileCount"] = fmt.Sprintf("%d", fileCount)
	metadata["DirectoryCount"] = fmt.Sprintf("%d", dirCount)
	metadata["TotalSize"] = fmt.Sprintf("%d bytes", totalSize)

	return metadata
}
