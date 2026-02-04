package tests

import (
	"os"
	"testing"

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

// TestVerifyHash tests hash verification functionality (removed - function doesn't exist in utils).

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

// TestGenerateEvidenceID tests evidence ID generation.
func TestGenerateEvidenceID(t *testing.T) {
	id := utils.GenerateEvidenceID("testhost")

	if id == "" {
		t.Error("Expected non-empty evidence ID")
	}

	if len(id) < 10 {
		t.Errorf("Expected evidence ID to be at least 10 characters, got %d", len(id))
	}
}
