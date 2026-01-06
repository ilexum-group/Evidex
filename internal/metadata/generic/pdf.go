package generic

import (
	"bufio"
	"os"
	"strings"
)

// PDFExtractor implements metadata extraction for PDF documents
type PDFExtractor struct{}

// NewPDFExtractor creates a new PDF extractor
func NewPDFExtractor() *PDFExtractor {
	return &PDFExtractor{}
}

// CanHandle checks if the file is a PDF document
func (e *PDFExtractor) CanHandle(filePath string) bool {
	return IsPDF(filePath)
}

// Extract implements MetadataExtractor interface
func (e *PDFExtractor) Extract(filePath string) (interface{}, error) {
	return e.ExtractDocument(filePath), nil
}

// GetType returns the extractor type
func (e *PDFExtractor) GetType() string {
	return "application/pdf"
}

// ExtractDocument extracts PDF metadata
func (e *PDFExtractor) ExtractDocument(filePath string) map[string]string {
	return ParsePDF(filePath)
}

// IsPDF checks if a file is a PDF by magic bytes
func IsPDF(filePath string) bool {
	file, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer func() {
		_ = file.Close()
	}()

	magic := make([]byte, 5)
	if _, err := file.Read(magic); err != nil {
		return false
	}

	// PDF signature: %PDF-
	return len(magic) == 5 && string(magic) == "%PDF-"
}

// ParsePDF extracts metadata from PDF files
func ParsePDF(filePath string) map[string]string {
	metadata := make(map[string]string)
	metadata["Type"] = "PDF"

	file, err := os.Open(filePath)
	if err != nil {
		return metadata
	}
	defer func() {
		_ = file.Close()
	}()

	scanner := bufio.NewScanner(file)
	infoDict := false

	for scanner.Scan() {
		line := scanner.Text()

		// Look for PDF version
		if strings.HasPrefix(line, "%PDF-") {
			metadata["Version"] = strings.TrimPrefix(line, "%PDF-")
		}

		// Look for Info dictionary
		if strings.Contains(line, "/Info") {
			infoDict = true
		}

		// Extract metadata fields
		if infoDict {
			if strings.Contains(line, "/Title") {
				metadata["Title"] = extractPDFValue(line, "/Title")
			}
			if strings.Contains(line, "/Author") {
				metadata["Author"] = extractPDFValue(line, "/Author")
			}
			if strings.Contains(line, "/Subject") {
				metadata["Subject"] = extractPDFValue(line, "/Subject")
			}
			if strings.Contains(line, "/Keywords") {
				metadata["Keywords"] = extractPDFValue(line, "/Keywords")
			}
			if strings.Contains(line, "/Creator") {
				metadata["Creator"] = extractPDFValue(line, "/Creator")
			}
			if strings.Contains(line, "/Producer") {
				metadata["Producer"] = extractPDFValue(line, "/Producer")
			}
			if strings.Contains(line, "/CreationDate") {
				metadata["CreationDate"] = extractPDFValue(line, "/CreationDate")
			}
			if strings.Contains(line, "/ModDate") {
				metadata["ModificationDate"] = extractPDFValue(line, "/ModDate")
			}

			// End of Info dictionary
			if strings.Contains(line, ">>") && strings.Count(line, ">>") > strings.Count(line, "<<") {
				infoDict = false
			}
		}
	}

	return metadata
}

// extractPDFValue extracts a value from a PDF metadata line
func extractPDFValue(line, key string) string {
	idx := strings.Index(line, key)
	if idx == -1 {
		return ""
	}

	rest := line[idx+len(key):]
	rest = strings.TrimSpace(rest)

	// Handle parentheses-enclosed strings
	if strings.HasPrefix(rest, "(") {
		endIdx := strings.Index(rest, ")")
		if endIdx != -1 {
			return rest[1:endIdx]
		}
	}

	// Handle angle bracket-enclosed hex strings
	if strings.HasPrefix(rest, "<") {
		endIdx := strings.Index(rest, ">")
		if endIdx != -1 {
			return rest[1:endIdx]
		}
	}

	// Handle other formats
	parts := strings.Fields(rest)
	if len(parts) > 0 {
		return parts[0]
	}

	return ""
}
