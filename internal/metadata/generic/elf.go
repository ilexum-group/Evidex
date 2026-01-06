package generic

import (
	"debug/elf"
	"fmt"
	"os"
)

// ELFExtractor implements metadata extraction for Linux ELF executables
type ELFExtractor struct{}

// NewELFExtractor creates a new ELF extractor
func NewELFExtractor() *ELFExtractor {
	return &ELFExtractor{}
}

// CanHandle checks if the file is an ELF executable
func (e *ELFExtractor) CanHandle(filePath string) bool {
	return IsELF(filePath)
}

// Extract implements MetadataExtractor interface
func (e *ELFExtractor) Extract(filePath string) (interface{}, error) {
	return e.ExtractDocument(filePath), nil
}

// GetType returns the extractor type
func (e *ELFExtractor) GetType() string {
	return "application/x-executable"
}

// ExtractDocument extracts ELF metadata
func (e *ELFExtractor) ExtractDocument(filePath string) map[string]string {
	return ParseELF(filePath)
}

// IsELF checks if a file is a Linux ELF executable by magic bytes
func IsELF(filePath string) bool {
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

	// ELF signature: 7F 45 4C 46
	return len(magic) == 4 && magic[0] == 0x7F &&
		magic[1] == 'E' && magic[2] == 'L' && magic[3] == 'F'
}

// ParseELF extracts metadata from Linux ELF executables
func ParseELF(filePath string) map[string]string {
	metadata := make(map[string]string)
	metadata["Type"] = "ELF Executable"

	file, err := elf.Open(filePath)
	if err != nil {
		return metadata
	}
	defer func() {
		_ = file.Close()
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
