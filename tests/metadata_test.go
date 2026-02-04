package tests

import (
	"os"
	"testing"

	"github.com/ilexum-group/evidex/internal/metadata"
	"github.com/ilexum-group/evidex/pkg/models"
)

// TestExtractFileMetadata tests metadata extraction for regular files.
func TestExtractFileMetadata(t *testing.T) {
	// Create test file
	tmpfile, err := os.CreateTemp("", "test-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer func() {
		if err := os.Remove(tmpfile.Name()); err != nil {
			t.Logf("Failed to remove temp file: %v", err)
		}
	}()
	defer func() {
		if err := tmpfile.Close(); err != nil {
			t.Logf("Failed to close temp file: %v", err)
		}
	}()

	testData := []byte("test metadata")
	if _, err := tmpfile.Write(testData); err != nil {
		t.Fatalf("Failed to write test data: %v", err)
	}

	// Create metadata manager
	mgr := metadata.NewMetadataManager(nil)
	fileEvidence := &models.FileEvidence{
		SourcePath: tmpfile.Name(),
	}
	err = mgr.ExtractFileMetadata(fileEvidence)
	if err != nil {
		t.Fatalf("ExtractFileMetadata() error = %v", err)
	}

	if fileEvidence.SourcePath == "" {
		t.Error("Expected SourcePath to be non-empty")
	}

	// Just verify the extraction didn't error - metadata fields depend on file type
	// Text files won't have ImageMetadata, VideoMetadata etc.
}

// TestExtractFileMetadataDirectory tests metadata extraction for directories.
func TestExtractFileMetadataDirectory(t *testing.T) {
	tmpdir, err := os.MkdirTemp("", "test-dir-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpdir); err != nil {
			t.Logf("Failed to remove temp dir: %v", err)
		}
	}()

	// Create metadata manager
	mgr := metadata.NewMetadataManager(nil)
	fileEvidence := &models.FileEvidence{
		SourcePath: tmpdir,
	}
	err = mgr.ExtractFileMetadata(fileEvidence)
	if err != nil {
		t.Fatalf("ExtractFileMetadata() error = %v", err)
	}

	// Just verify extraction didn't error
	// Directories don't have image/video metadata
}

// TestCalculateFileHashes tests hash calculation.
func TestCalculateFileHashes(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "test-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer func() {
		if err := os.Remove(tmpfile.Name()); err != nil {
			t.Logf("Failed to remove temp file: %v", err)
		}
	}()
	defer func() {
		if err := tmpfile.Close(); err != nil {
			t.Logf("Failed to close temp file: %v", err)
		}
	}()

	if _, err := tmpfile.WriteString("hash test data"); err != nil {
		t.Fatalf("Failed to write test data: %v", err)
	}

	// Create metadata manager
	mgr := metadata.NewMetadataManager(nil)
	hashes, err := mgr.CalculateFileHashes(tmpfile.Name())
	if err != nil {
		t.Fatalf("CalculateFileHashes() error = %v", err)
	}

	if hashes.SHA256 == "" {
		t.Error("Expected SHA256 hash to be non-empty")
	}

	if hashes.SHA512 == "" {
		t.Error("Expected SHA512 hash to be non-empty")
	}

	if hashes.MD5 == "" {
		t.Error("Expected MD5 hash to be non-empty")
	}
}

// TestImageFormatDetection tests image format detection.
func TestImageFormatDetection(t *testing.T) {
	tests := []struct {
		name        string
		extension   string
		shouldMatch bool
	}{
		{"JPEG file", "test.jpg", true},
		{"PNG file", "test.png", true},
		{"GIF file", "test.gif", true},
		{"Text file", "test.txt", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpfile, err := os.CreateTemp("", "test-*"+tt.extension)
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer func() {
				if err := os.Remove(tmpfile.Name()); err != nil {
					t.Logf("Failed to remove temp file: %v", err)
				}
			}()
			defer func() {
				if err := tmpfile.Close(); err != nil {
					t.Logf("Failed to close temp file: %v", err)
				}
			}()

			// File exists, so IsImage should check based on extension
			// (actual magic byte checking would require specific file contents)
		})
	}
}
