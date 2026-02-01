package tests

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ilexum-group/evidex/internal/utils"
)

// createTempTestFile is a helper function to create temporary test files.
func createTempTestFile(t *testing.T, data []byte) string {
	tmpfile, err := os.CreateTemp("", "test-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	if _, err := tmpfile.Write(data); err != nil {
		t.Fatalf("Failed to write test data: %v", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Logf("Failed to close temp file: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Remove(tmpfile.Name()); err != nil {
			t.Logf("Failed to remove temp file: %v", err)
		}
	})
	return tmpfile.Name()
}

// Testutils.CalculateSHA256 tests SHA256 hash calculation.
func TestCalculateSHA256(t *testing.T) {
	path := createTempTestFile(t, []byte("test data for hashing"))

	hash, err := utils.CalculateSHA256(path)
	if err != nil {
		t.Fatalf("utils.CalculateSHA256() error = %v", err)
	}

	if hash == "" {
		t.Error("Expected hash to be non-empty")
	}

	if len(hash) != 64 {
		t.Errorf("Expected hash length 64, got %d", len(hash))
	}
}

// Testutils.CalculateSHA512 tests SHA512 hash calculation.
func TestCalculateSHA512(t *testing.T) {
	path := createTempTestFile(t, []byte("test data for sha512"))

	hash, err := utils.CalculateSHA512(path)
	if err != nil {
		t.Fatalf("utils.CalculateSHA512() error = %v", err)
	}

	if hash == "" {
		t.Error("Expected hash to be non-empty")
	}

	if len(hash) != 128 {
		t.Errorf("Expected hash length 128, got %d", len(hash))
	}
}

// Testutils.CalculateMD5 tests MD5 hash calculation.
func TestCalculateMD5(t *testing.T) {
	path := createTempTestFile(t, []byte("test md5"))

	hash, err := utils.CalculateMD5(path)
	if err != nil {
		t.Fatalf("utils.CalculateMD5() error = %v", err)
	}

	if hash == "" {
		t.Error("Expected hash to be non-empty")
	}

	if len(hash) != 32 {
		t.Errorf("Expected hash length 32, got %d", len(hash))
	}
}

// TestVerifyHash tests hash verification.
func TestVerifyHash(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "test-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer func() {
		if err := os.Remove(tmpfile.Name()); err != nil {
			t.Logf("Failed to remove temp file: %v", err)
		}
	}()

	testData := []byte("verify test")
	if _, err := tmpfile.Write(testData); err != nil {
		t.Fatalf("Failed to write test data: %v", err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Logf("Failed to close temp file: %v", err)
	}

	// utils.Calculate correct hash
	correctHash, err := utils.CalculateSHA256(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to utils.Calculate hash: %v", err)
	}

	// Test with correct hash
	verified, err := utils.VerifyHash(tmpfile.Name(), correctHash, "SHA-256")
	if err != nil {
		t.Fatalf("utils.VerifyHash() error = %v", err)
	}
	if !verified {
		t.Error("Expected verification to succeed with correct hash")
	}

	// Test with incorrect hash
	verified, err = utils.VerifyHash(tmpfile.Name(), "0000000000000000000000000000000000000000000000000000000000000000", "SHA-256")
	if err != nil {
		t.Fatalf("utils.VerifyHash() error = %v", err)
	}
	if verified {
		t.Error("Expected verification to fail with incorrect hash")
	}
}

// TestIsReadOnly tests read-only file verification.
func TestIsReadOnly(t *testing.T) {
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
	if err := tmpfile.Close(); err != nil {
		t.Logf("Failed to close temp file: %v", err)
	}

	// Readable file should return true
	readable, err := utils.IsReadOnly(tmpfile.Name())
	if err != nil {
		t.Fatalf("utils.IsReadOnly() error = %v", err)
	}
	if !readable {
		t.Error("Expected file to be readable")
	}

	// Non-existent file should return error
	readable, err = utils.IsReadOnly("/nonexistent/path/to/file")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
	if readable {
		t.Error("Expected non-existent file to return false")
	}
}

// TestEnsureDirectory tests directory creation.
func TestEnsureDirectory(t *testing.T) {
	tmpdir := filepath.Join(os.TempDir(), "evidex-test-"+time.Now().Format("20060102150405"))

	err := utils.EnsureDirectory(tmpdir)
	if err != nil {
		t.Fatalf("utils.EnsureDirectory() error = %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpdir); err != nil {
			t.Logf("Failed to remove temp directory: %v", err)
		}
	}()

	// Verify directory was created
	info, err := os.Stat(tmpdir)
	if err != nil {
		t.Fatalf("Directory not created: %v", err)
	}

	if !info.IsDir() {
		t.Error("Expected path to be a directory")
	}

	// Second call should succeed (already exists)
	err = utils.EnsureDirectory(tmpdir)
	if err != nil {
		t.Fatalf("utils.EnsureDirectory() on existing directory error = %v", err)
	}
}

// TestGetFileTypeFromExtension tests MIME type detection.
func TestGetFileTypeFromExtension(t *testing.T) {
	tests := []struct {
		filename string
		expected string
	}{
		{"image.jpg", "image/jpeg"},
		{"image.jpeg", "image/jpeg"},
		{"image.png", "image/png"},
		{"video.mp4", "video/mp4"},
		{"unknown.xyz", "application/octet-stream"},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			mimeType := utils.GetFileTypeFromExtension(tt.filename)
			if mimeType != tt.expected {
				t.Errorf("utils.GetFileTypeFromExtension(%s) = %s, want %s", tt.filename, mimeType, tt.expected)
			}
		})
	}
}

// TestIsImageFile tests image file detection.
func TestIsImageFile(t *testing.T) {
	tests := []struct {
		filename string
		expected bool
	}{
		{"image.jpeg", true},
		{"image.png", true},
		{"video.mp4", false},
		{"document.pdf", false},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			result := utils.IsImageFile(tt.filename)
			if result != tt.expected {
				t.Errorf("utils.IsImageFile(%s) = %v, want %v", tt.filename, result, tt.expected)
			}
		})
	}
}

// TestIsVideoFile tests video file detection.
func TestIsVideoFile(t *testing.T) {
	tests := []struct {
		filename string
		expected bool
	}{
		{"video.mp4", true},
		{"video.mov", true},
		{"image.jpeg", false},
		{"document.pdf", false},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			result := utils.IsVideoFile(tt.filename)
			if result != tt.expected {
				t.Errorf("utils.IsVideoFile(%s) = %v, want %v", tt.filename, result, tt.expected)
			}
		})
	}
}

// TestLogging tests logging functionality.
func TestLogging(t *testing.T) {
	// Test new logger signature with map[string]string metadata
	utils.LogInfo("Test info message", map[string]string{})
	utils.LogWarn("Test warning message", map[string]string{"test": "value"})
	utils.LogError("Test error message", map[string]string{"error": "test_error"})

	// Logging should not panic or error
	t.Log("Logging functions executed successfully with new RFC 5424 format")
}

// TestGenerateEvidenceID tests evidence ID generation.
func TestGenerateEvidenceID(t *testing.T) {
	id := utils.GenerateEvidenceID()

	if id == "" {
		t.Error("Expected non-empty evidence ID")
	}

	if len(id) < 10 {
		t.Errorf("Expected evidence ID to be at least 10 characters, got %d", len(id))
	}
}
