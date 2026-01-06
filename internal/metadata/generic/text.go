package generic

import (
	"bufio"
	"os"
	"unicode/utf8"
)

// TextExtractor implements metadata extraction for text files
type TextExtractor struct{}

// NewTextExtractor creates a new text extractor
func NewTextExtractor() *TextExtractor {
	return &TextExtractor{}
}

// CanHandle checks if the file is a text file
func (e *TextExtractor) CanHandle(filePath string) bool {
	// This is a fallback extractor, it can handle any file
	// but should be registered last in the registry
	return true
}

// Extract implements MetadataExtractor interface
func (e *TextExtractor) Extract(filePath string) (interface{}, error) {
	return e.ExtractDocument(filePath), nil
}

// GetType returns the extractor type
func (e *TextExtractor) GetType() string {
	return "text/plain"
}

// ExtractDocument extracts text file metadata
func (e *TextExtractor) ExtractDocument(filePath string) map[string]string {
	return ProbeText(filePath)
}

// ProbeText attempts to detect if a file is text and count lines/chars
func ProbeText(filePath string) map[string]string {
	metadata := make(map[string]string)

	file, err := os.Open(filePath)
	if err != nil {
		return metadata
	}
	defer func() {
		_ = file.Close()
	}()

	// Read first few bytes to determine if it's text
	header := make([]byte, 512)
	n, err := file.Read(header)
	if err != nil && n == 0 {
		return metadata
	}

	// Check if content is valid UTF-8
	if !utf8.Valid(header[:n]) {
		metadata["Type"] = "Binary"
		return metadata
	}

	// Count lines and words
	if _, err := file.Seek(0, 0); err != nil {
		return metadata
	}
	scanner := bufio.NewScanner(file)
	lineCount := 0
	wordCount := 0
	charCount := 0

	for scanner.Scan() {
		lineCount++
		line := scanner.Text()
		charCount += len(line)

		// Simple word count
		inWord := false
		for _, r := range line {
			if r == ' ' || r == '\t' || r == '\n' {
				inWord = false
			} else if !inWord {
				wordCount++
				inWord = true
			}
		}
	}

	metadata["Type"] = "Text"
	metadata["Lines"] = string(rune(lineCount + '0'))
	metadata["Words"] = string(rune(wordCount + '0'))
	metadata["Characters"] = string(rune(charCount + '0'))

	return metadata
}
