// Package generic provides metadata extraction for generic file formats including executables, archives, and documents.
package generic

import (
	"debug/pe"
	"fmt"
	"os"
	"time"

	"github.com/ilexum-group/evidex/internal/utils"
	"github.com/ilexum-group/evidex/pkg/models"
)

// PEExtractor implements metadata extraction for Windows PE executables.
type PEExtractor struct{}

// NewPEExtractor creates a new PE extractor.
func NewPEExtractor() *PEExtractor {
	return &PEExtractor{}
}

// CanHandle checks if the file is a PE executable.
func (e *PEExtractor) CanHandle(filePath string) bool {
	return IsPE(filePath)
}

// Extract implements MetadataExtractor interface.
func (e *PEExtractor) Extract(filePath string, logCmd models.CommandLogger) (interface{}, error) {
	return e.ExtractDocument(filePath, logCmd), nil
}

// GetType returns the extractor type.
func (e *PEExtractor) GetType() string {
	return "application/x-msdownload"
}

// ExtractDocument extracts PE metadata.
func (e *PEExtractor) ExtractDocument(filePath string, logCmd models.CommandLogger) map[string]string {
	return ParsePE(filePath, logCmd)
}

// IsPE checks if a file is a Windows PE executable by magic bytes.
func IsPE(filePath string) bool {
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

	// PE files start with MZ header
	return len(magic) == 2 && magic[0] == 'M' && magic[1] == 'Z'
}

// ParsePE extracts metadata from Windows PE executables.
func ParsePE(filePath string, logCmd models.CommandLogger) map[string]string {
	metadata := make(map[string]string)
	metadata["Type"] = "Windows PE Executable"

	startTime := time.Now()
	file, err := pe.Open(filePath)
	endTime := time.Now()
	exitCode := 0
	if err != nil {
		exitCode = 1
	}
	if logCmd != nil {
		logCmd(utils.GenerateRandomID(), "pe.Open", []string{filePath}, startTime, endTime, exitCode, err, "", filePath)
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
			logCmd(utils.GenerateRandomID(), "pe.File.Close", []string{filePath}, startTime, endTime, exitCode, err, "", filePath)
		}
	}()

	// Determine architecture
	switch file.Machine {
	case pe.IMAGE_FILE_MACHINE_I386:
		metadata["Architecture"] = "x86 (32-bit)"
	case pe.IMAGE_FILE_MACHINE_AMD64:
		metadata["Architecture"] = "x86-64 (64-bit)"
	case pe.IMAGE_FILE_MACHINE_ARM:
		metadata["Architecture"] = "ARM"
	case pe.IMAGE_FILE_MACHINE_ARM64:
		metadata["Architecture"] = "ARM64"
	default:
		metadata["Architecture"] = fmt.Sprintf("Unknown (0x%X)", file.Machine)
	}

	// Get compilation timestamp
	if fileHeader, ok := file.OptionalHeader.(*pe.OptionalHeader32); ok {
		metadata["Subsystem"] = fmt.Sprintf("%d", fileHeader.Subsystem)
		metadata["ImageBase"] = fmt.Sprintf("0x%X", fileHeader.ImageBase)
		metadata["EntryPoint"] = fmt.Sprintf("0x%X", fileHeader.AddressOfEntryPoint)
	} else if fileHeader, ok := file.OptionalHeader.(*pe.OptionalHeader64); ok {
		metadata["Subsystem"] = fmt.Sprintf("%d", fileHeader.Subsystem)
		metadata["ImageBase"] = fmt.Sprintf("0x%X", fileHeader.ImageBase)
		metadata["EntryPoint"] = fmt.Sprintf("0x%X", fileHeader.AddressOfEntryPoint)
	}

	// Get timestamp (seconds since 1970-01-01)
	if file.TimeDateStamp > 0 {
		timestamp := time.Unix(int64(file.TimeDateStamp), 0)
		metadata["CompilationTime"] = timestamp.String()
	}

	// Count sections
	metadata["SectionCount"] = fmt.Sprintf("%d", len(file.Sections))

	// List section names
	var sections []string
	for _, section := range file.Sections {
		sections = append(sections, section.Name)
	}
	if len(sections) > 0 {
		metadata["Sections"] = fmt.Sprintf("%v", sections)
	}

	// Characteristics
	metadata["Characteristics"] = fmt.Sprintf("0x%X", file.Characteristics)

	return metadata
}
