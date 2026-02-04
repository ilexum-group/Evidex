// Package generic provides metadata extraction for generic file formats including executables, archives, and documents.
package generic

import (
	"bufio"
	"os"
	"time"
	"unicode/utf8"

	"github.com/ilexum-group/evidex/internal/utils"
	"github.com/ilexum-group/evidex/pkg/models"
)

// TextExtractor implements metadata extraction for text files.
type TextExtractor struct{}

// NewTextExtractor creates a new text extractor.
func NewTextExtractor() *TextExtractor {
	return &TextExtractor{}
}

// CanHandle checks if the file is a text file.
func (e *TextExtractor) CanHandle(_ string) bool {
	// This is a fallback extractor, it can handle any file
	// but should be registered last in the registry
	return true
}

// Extract implements MetadataExtractor interface.
func (e *TextExtractor) Extract(filePath string, logCmd models.CommandLogger) (interface{}, error) {
	return e.ExtractDocument(filePath, logCmd), nil
}

// GetType returns the extractor type.
func (e *TextExtractor) GetType() string {
	return "text/plain"
}

// ExtractDocument extracts text file metadata.
func (e *TextExtractor) ExtractDocument(filePath string, logCmd models.CommandLogger) map[string]string {
	return ProbeText(filePath, logCmd)
}

// ProbeText attempts to detect if a file is text and count lines/chars.
func ProbeText(filePath string, logCmd models.CommandLogger) map[string]string {
	metadata := make(map[string]string)

	startTime := time.Now()
	file, err := os.Open(filePath) // #nosec G304 - filepath is user-provided evidence path in forensic tool
	endTime := time.Now()
	exitCode := 0
	if err != nil {
		exitCode = 1
	}
	if logCmd != nil {
		logCmd(utils.GenerateRandomID(), "os.Open", []string{filePath, "text probe"}, startTime, endTime, exitCode, err, "", filePath)
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
			logCmd(utils.GenerateRandomID(), "file.Close", []string{filePath, "text probe"}, startTime, endTime, exitCode, err, "", filePath)
		}
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
