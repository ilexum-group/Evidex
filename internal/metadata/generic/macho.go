// Package generic provides metadata extraction for generic file formats including executables, archives, and documents.
package generic

import (
	"debug/macho"
	"fmt"
	"os"
	"time"

	"github.com/ilexum-group/evidex/internal/utils"
	"github.com/ilexum-group/evidex/pkg/models"
)

// MachOExtractor implements metadata extraction for macOS Mach-O executables.
type MachOExtractor struct{}

// NewMachOExtractor creates a new Mach-O extractor.
func NewMachOExtractor() *MachOExtractor {
	return &MachOExtractor{}
}

// CanHandle checks if the file is a Mach-O executable.
func (e *MachOExtractor) CanHandle(filePath string) bool {
	return IsMachO(filePath)
}

// Extract implements MetadataExtractor interface.
func (e *MachOExtractor) Extract(filePath string, logCmd models.CommandLogger) (interface{}, error) {
	return e.ExtractDocument(filePath, logCmd), nil
}

// GetType returns the extractor type.
func (e *MachOExtractor) GetType() string {
	return "application/x-mach-binary"
}

// ExtractDocument extracts Mach-O metadata.
func (e *MachOExtractor) ExtractDocument(filePath string, logCmd models.CommandLogger) map[string]string {
	return ParseMachO(filePath, logCmd)
}

// IsMachO checks if a file is a macOS Mach-O executable by magic bytes.
func IsMachO(filePath string) bool {
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

	if len(magic) != 4 {
		return false
	}

	// Mach-O magic numbers
	return (magic[0] == 0xFE && magic[1] == 0xED && magic[2] == 0xFA && magic[3] == 0xCE) || // 32-bit
		(magic[0] == 0xFE && magic[1] == 0xED && magic[2] == 0xFA && magic[3] == 0xCF) || // 64-bit
		(magic[0] == 0xCE && magic[1] == 0xFA && magic[2] == 0xED && magic[3] == 0xFE) || // 32-bit reverse
		(magic[0] == 0xCF && magic[1] == 0xFA && magic[2] == 0xED && magic[3] == 0xFE) || // 64-bit reverse
		(magic[0] == 0xCA && magic[1] == 0xFE && magic[2] == 0xBA && magic[3] == 0xBE) // Universal/Fat binary
}

// ParseMachO extracts metadata from macOS Mach-O executables.
func ParseMachO(filePath string, logCmd models.CommandLogger) map[string]string {
	metadata := make(map[string]string)
	metadata["type"] = "macOS Mach-O Executable"

	startTime := time.Now()
	file, err := macho.Open(filePath)
	endTime := time.Now()
	exitCode := 0
	if err != nil {
		exitCode = 1
	}
	if logCmd != nil {
		logCmd(utils.GenerateRandomID(), "macho.Open", []string{filePath}, startTime, endTime, exitCode, err, "", filePath)
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
			logCmd(utils.GenerateRandomID(), "macho.File.Close", []string{filePath}, startTime, endTime, exitCode, err, "", filePath)
		}
	}()

	// CPU type
	switch file.Cpu {
	case macho.Cpu386:
		metadata["Architecture"] = "Intel 386"
	case macho.CpuAmd64:
		metadata["Architecture"] = "AMD64 (x86-64)"
	case macho.CpuArm:
		metadata["Architecture"] = "ARM"
	case macho.CpuArm64:
		metadata["Architecture"] = "ARM64"
	case macho.CpuPpc:
		metadata["Architecture"] = "PowerPC"
	case macho.CpuPpc64:
		metadata["Architecture"] = "PowerPC 64-bit"
	default:
		metadata["Architecture"] = fmt.Sprintf("Unknown (0x%X)", file.Cpu)
	}

	// File type
	switch file.Type {
	case macho.TypeObj:
		metadata["FileType"] = "Relocatable object file"
	case macho.TypeExec:
		metadata["FileType"] = "Executable"
	case macho.TypeDylib:
		metadata["FileType"] = "Dynamic library"
	case macho.TypeBundle:
		metadata["FileType"] = "Bundle"
	default:
		metadata["FileType"] = fmt.Sprintf("Unknown (%d)", file.Type)
	}

	// Flags
	metadata["Flags"] = fmt.Sprintf("0x%X", file.Flags)

	// Count segments
	var segments []string
	for _, load := range file.Loads {
		if seg, ok := load.(*macho.Segment); ok {
			segments = append(segments, seg.Name)
		}
	}

	metadata["SegmentCount"] = fmt.Sprintf("%d", len(segments))
	if len(segments) > 0 {
		metadata["Segments"] = fmt.Sprintf("%v", segments)
	}

	// Load command count
	metadata["LoadCommandCount"] = fmt.Sprintf("%d", len(file.Loads))

	return metadata
}
