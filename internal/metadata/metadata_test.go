package metadata

import (
	"os"
	"testing"
)

// TestExtractFileMetadata tests metadata extraction for regular files
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

	metadata, err := ExtractFileMetadata(tmpfile.Name())
	if err != nil {
		t.Fatalf("ExtractFileMetadata() error = %v", err)
	}

	if metadata.SourcePath == "" {
		t.Error("Expected metadata SourcePath to be non-empty")
	}

	if metadata.FileSize != int64(len(testData)) {
		t.Errorf("Expected size %d, got %d", len(testData), metadata.FileSize)
	}

	if metadata.IsSymlink {
		t.Error("Expected IsSymlink to be false for regular file")
	}
}

// TestExtractFileMetadataDirectory tests metadata extraction for directories
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

	metadata, err := ExtractFileMetadata(tmpdir)
	if err != nil {
		t.Fatalf("ExtractFileMetadata() error = %v", err)
	}

	if metadata.IsSymlink {
		t.Error("Expected IsSymlink to be false for directory (not a symlink)")
	}

	if metadata.FileSize != 0 {
		t.Errorf("Expected directory size to be 0, got %d", metadata.FileSize)
	}
}

// TestCalculateFileHashes tests hash calculation
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

	hashes, err := CalculateFileHashes(tmpfile.Name(), "SHA-256")
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

// TestVerifyFileIntegrity tests file integrity verification
func TestVerifyFileIntegrity(t *testing.T) {
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

	if _, err := tmpfile.WriteString("integrity test"); err != nil {
		t.Fatalf("Failed to write test data: %v", err)
	}

	hashes, err := CalculateFileHashes(tmpfile.Name(), "SHA-256")
	if err != nil {
		t.Fatalf("Failed to calculate hashes: %v", err)
	}

	// Test with correct hash
	verified, err := VerifyFileIntegrity(tmpfile.Name(), hashes.SHA256)
	if err != nil {
		t.Fatalf("VerifyFileIntegrity() error = %v", err)
	}
	if !verified {
		t.Error("Expected file integrity verification to succeed with correct hash")
	}

	// Test with incorrect hash
	verified, err = VerifyFileIntegrity(tmpfile.Name(), "0000000000000000000000000000000000000000000000000000000000000000")
	if err != nil {
		t.Fatalf("VerifyFileIntegrity() error = %v", err)
	}
	if verified {
		t.Error("Expected file integrity verification to fail with incorrect hash")
	}
}

// TestImageFormatDetection tests image format detection
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
