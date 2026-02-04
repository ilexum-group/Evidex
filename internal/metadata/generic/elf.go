// Package generic provides metadata extraction for generic file formats including executables, archives, and documents.
package generic

import (
	"debug/elf"
	"fmt"
	"os"
	"time"

	"github.com/ilexum-group/evidex/internal/utils"
	"github.com/ilexum-group/evidex/pkg/models"
)

// ELFExtractor implements metadata extraction for Linux ELF executables.
type ELFExtractor struct{}

// NewELFExtractor creates a new ELF extractor.
func NewELFExtractor() *ELFExtractor {
	return &ELFExtractor{}
}

// CanHandle checks if the file is an ELF executable.
func (e *ELFExtractor) CanHandle(filePath string) bool {
	return IsELF(filePath)
}

// Extract implements MetadataExtractor interface.
func (e *ELFExtractor) Extract(filePath string, logCmd models.CommandLogger) (interface{}, error) {
	return e.ExtractDocument(filePath, logCmd), nil
}

// GetType returns the extractor type.
func (e *ELFExtractor) GetType() string {
	return "application/x-executable"
}

// ExtractDocument extracts ELF metadata.
func (e *ELFExtractor) ExtractDocument(filePath string, logCmd models.CommandLogger) map[string]string {
	return ParseELF(filePath, logCmd)
}

// IsELF checks if a file is a Linux ELF executable by magic bytes.
func IsELF(filePath string) bool {
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

	// ELF signature: 7F 45 4C 46
	return len(magic) == 4 && magic[0] == 0x7F &&
		magic[1] == 'E' && magic[2] == 'L' && magic[3] == 'F'
}

// ParseELF extracts metadata from Linux ELF executables.
func ParseELF(filePath string, logCmd models.CommandLogger) map[string]string {
	metadata := make(map[string]string)
	metadata["type"] = "Linux ELF Executable"

	startTime := time.Now()
	file, err := elf.Open(filePath)
	endTime := time.Now()
	exitCode := 0
	if err != nil {
		exitCode = 1
	}
	if logCmd != nil {
		logCmd(utils.GenerateRandomID(), "elf.Open", []string{filePath}, startTime, endTime, exitCode, err, "", filePath)
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
			logCmd(utils.GenerateRandomID(), "elf.File.Close", []string{filePath}, startTime, endTime, exitCode, err, "", filePath)
		}
	}()

	// Class (32-bit or 64-bit)
	switch file.Class {
	case elf.ELFCLASS32:
		metadata["Class"] = "32-bit"
	case elf.ELFCLASS64:
		metadata["Class"] = "64-bit"
	default:
		metadata["Class"] = "Unknown"
	}

	// Data encoding (endianness)
	switch file.Data {
	case elf.ELFDATA2LSB:
		metadata["Endianness"] = "Little-endian"
	case elf.ELFDATA2MSB:
		metadata["Endianness"] = "Big-endian"
	default:
		metadata["Endianness"] = "Unknown"
	}

	// OS/ABI
	switch file.OSABI {
	case elf.ELFOSABI_NONE:
		metadata["OS/ABI"] = "UNIX System V"
	case elf.ELFOSABI_LINUX:
		metadata["OS/ABI"] = "Linux"
	case elf.ELFOSABI_FREEBSD:
		metadata["OS/ABI"] = "FreeBSD"
	case elf.ELFOSABI_NETBSD:
		metadata["OS/ABI"] = "NetBSD"
	case elf.ELFOSABI_OPENBSD:
		metadata["OS/ABI"] = "OpenBSD"
	case elf.ELFOSABI_SOLARIS:
		metadata["OS/ABI"] = "Solaris"
	default:
		metadata["OS/ABI"] = fmt.Sprintf("Unknown (%d)", file.OSABI)
	}

	// File type
	switch file.Type {
	case elf.ET_NONE:
		metadata["FileType"] = "None"
	case elf.ET_REL:
		metadata["FileType"] = "Relocatable"
	case elf.ET_EXEC:
		metadata["FileType"] = "Executable"
	case elf.ET_DYN:
		metadata["FileType"] = "Shared object"
	case elf.ET_CORE:
		metadata["FileType"] = "Core file"
	default:
		metadata["FileType"] = fmt.Sprintf("Unknown (%d)", file.Type)
	}

	// Machine architecture
	switch file.Machine {
	case elf.EM_386:
		metadata["Architecture"] = "Intel 80386"
	case elf.EM_X86_64:
		metadata["Architecture"] = "AMD x86-64"
	case elf.EM_ARM:
		metadata["Architecture"] = "ARM"
	case elf.EM_AARCH64:
		metadata["Architecture"] = "AArch64"
	case elf.EM_PPC:
		metadata["Architecture"] = "PowerPC"
	case elf.EM_PPC64:
		metadata["Architecture"] = "PowerPC 64-bit"
	case elf.EM_MIPS:
		metadata["Architecture"] = "MIPS"
	case elf.EM_RISCV:
		metadata["Architecture"] = "RISC-V"
	default:
		metadata["Architecture"] = fmt.Sprintf("Unknown (0x%X)", file.Machine)
	}

	// Entry point
	metadata["EntryPoint"] = fmt.Sprintf("0x%X", file.Entry)

	// Count sections
	metadata["SectionCount"] = fmt.Sprintf("%d", len(file.Sections))

	// List section names
	var sections []string
	for _, section := range file.Sections {
		if section.Name != "" {
			sections = append(sections, section.Name)
		}
	}
	if len(sections) > 0 {
		metadata["Sections"] = fmt.Sprintf("%v", sections)
	}

	return metadata
}
