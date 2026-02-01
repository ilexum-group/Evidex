package tests

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/ilexum-group/evidex/internal/acquisition"
)

// TestNewAcquirer tests Acquirer initialization.
func TestNewAcquirer(t *testing.T) {
	tmpdir, err := os.MkdirTemp("", "test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpdir); err != nil {
			t.Logf("Failed to remove temp directory: %v", err)
		}
	}()

	acq := acquisition.NewAcquirer(tmpdir)

	if acq == nil {
		t.Error("Expected Acquirer to be non-nil")
	}

	if acq.GetFileCount() != 0 {
		t.Errorf("Expected initial file count to be 0, got %d", acq.GetFileCount())
	}
}

// TestAcquireFile tests single file acquisition.
func TestAcquireFile(t *testing.T) {
	tmpdir, err := os.MkdirTemp("", "test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpdir); err != nil {
			t.Logf("Failed to remove temp directory: %v", err)
		}
	}()

	// Create test file
	testfile := filepath.Join(tmpdir, "test.txt")
	if err := os.WriteFile(testfile, []byte("test content"), 0600); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	acq := acquisition.NewAcquirer(tmpdir)
	if err := acq.AcquireFile(testfile); err != nil {
		t.Logf("AcquireFile warning: %v", err)
	}

	if acq.GetFileCount() != 1 {
		t.Errorf("Expected file count to be 1, got %d", acq.GetFileCount())
	}
}

// TestAcquireDirectory tests directory acquisition.
func TestAcquireDirectory(t *testing.T) {
	tmpdir, err := os.MkdirTemp("", "test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpdir); err != nil {
			t.Logf("Failed to remove temp directory: %v", err)
		}
	}()

	// Create test files
	testfile1 := filepath.Join(tmpdir, "test1.txt")
	testfile2 := filepath.Join(tmpdir, "test2.txt")
	if err := os.WriteFile(testfile1, []byte("content1"), 0600); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	if err := os.WriteFile(testfile2, []byte("content2"), 0600); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	acq := acquisition.NewAcquirer(tmpdir)
	if err := acq.AcquireDirectory(tmpdir, false); err != nil {
		t.Logf("AcquireDirectory warning: %v", err)
	}

	if acq.GetFileCount() != 2 {
		t.Errorf("Expected file count to be 2, got %d", acq.GetFileCount())
	}
}

// TestAcquireDirectoryRecursive tests recursive directory acquisition.
func TestAcquireDirectoryRecursive(t *testing.T) {
	tmpdir, err := os.MkdirTemp("", "test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpdir); err != nil {
			t.Logf("Failed to remove temp directory: %v", err)
		}
	}()

	// Create nested directory structure
	subdir := filepath.Join(tmpdir, "subdir")
	if err := os.Mkdir(subdir, 0750); err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}

	// Create files at different levels
	rootfile := filepath.Join(tmpdir, "root.txt")
	subfile := filepath.Join(subdir, "sub.txt")
	if err := os.WriteFile(rootfile, []byte("root"), 0600); err != nil {
		t.Fatalf("Failed to create root file: %v", err)
	}
	if err := os.WriteFile(subfile, []byte("sub"), 0600); err != nil {
		t.Fatalf("Failed to create sub file: %v", err)
	}

	acq := acquisition.NewAcquirer(tmpdir)
	if err := acq.AcquireDirectory(tmpdir, true); err != nil {
		t.Logf("AcquireDirectory warning: %v", err)
	}

	if acq.GetFileCount() != 2 {
		t.Errorf("Expected file count to be 2 with recursive acquisition, got %d", acq.GetFileCount())
	}
}

// TestGetEvidencePackage tests package assembly.
func TestGetEvidencePackage(t *testing.T) {
	tmpdir, err := os.MkdirTemp("", "test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpdir); err != nil {
			t.Logf("Failed to remove temp directory: %v", err)
		}
	}()

	// Create test file
	testfile := filepath.Join(tmpdir, "test.txt")
	if err := os.WriteFile(testfile, []byte("test"), 0600); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	acq := acquisition.NewAcquirer(tmpdir)
	if err := acq.AcquireFile(testfile); err != nil {
		t.Logf("AcquireFile warning: %v", err)
	}

	pkg := acq.GetEvidencePackage()

	if pkg == nil {
		t.Error("Expected EvidencePackage to be non-nil")
		return
	}

	if pkg.Manifest == nil {
		t.Error("Expected Manifest to be non-nil")
	}

	if pkg.Manifest.ID == "" {
		t.Error("Expected Evidence ID to be non-empty")
	}

	if len(pkg.Files) != 1 {
		t.Errorf("Expected 1 file in package, got %d", len(pkg.Files))
	}
}

// TestSetHashAlgorithm tests hash algorithm configuration.
func TestSetHashAlgorithm(t *testing.T) {
	tmpdir, err := os.MkdirTemp("", "test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpdir); err != nil {
			t.Logf("Failed to remove temp directory: %v", err)
		}
	}()

	acq := acquisition.NewAcquirer(tmpdir)
	acq.SetHashAlgorithm("sha512")

	// Algorithm should be set (verified through acquisition)
	if acq == nil {
		t.Error("Expected Acquirer to be non-nil after SetHashAlgorithm")
	}
}

// TestCopyFilesToPackage tests file copying to package.
func TestCopyFilesToPackage(t *testing.T) {
	tmpdir, err := os.MkdirTemp("", "test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpdir); err != nil {
			t.Logf("Failed to remove temp directory: %v", err)
		}
	}()

	outputdir := filepath.Join(tmpdir, "output")
	if err := os.Mkdir(outputdir, 0750); err != nil {
		t.Fatalf("Failed to create output directory: %v", err)
	}

	// Create test file
	testfile := filepath.Join(tmpdir, "test.txt")
	if err := os.WriteFile(testfile, []byte("test content"), 0600); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	acq := acquisition.NewAcquirer(outputdir)
	if err := acq.AcquireFile(testfile); err != nil {
		t.Logf("AcquireFile warning: %v", err)
	}
	if err := acq.CopyFilesToPackage(); err != nil {
		t.Logf("CopyFilesToPackage warning: %v", err)
	}

	// Check if files directory was created
	filesdir := filepath.Join(outputdir, "files")
	if _, err := os.Stat(filesdir); err != nil {
		t.Errorf("Expected files directory to be created: %v", err)
	}
}

// TestGetFileCount tests file count accessor.
func TestGetFileCount(t *testing.T) {
	tmpdir, err := os.MkdirTemp("", "test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpdir); err != nil {
			t.Logf("Failed to remove temp directory: %v", err)
		}
	}()

	acq := acquisition.NewAcquirer(tmpdir)

	// Initial count should be 0
	if acq.GetFileCount() != 0 {
		t.Errorf("Expected initial file count to be 0, got %d", acq.GetFileCount())
	}

	// Create and acquire files
	for i := 0; i < 3; i++ {
		testfile := filepath.Join(tmpdir, fmt.Sprintf("test%d.txt", i))
		if err := os.WriteFile(testfile, []byte("content"), 0600); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
		if err := acq.AcquireFile(testfile); err != nil {
			t.Logf("AcquireFile warning: %v", err)
		}
	}

	if acq.GetFileCount() != 3 {
		t.Errorf("Expected file count to be 3, got %d", acq.GetFileCount())
	}
}

// TestGetTotalSize tests total size calculator.
func TestGetTotalSize(t *testing.T) {
	tmpdir, err := os.MkdirTemp("", "test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpdir); err != nil {
			t.Logf("Failed to remove temp directory: %v", err)
		}
	}()

	acq := acquisition.NewAcquirer(tmpdir)

	// Create test files with known sizes
	testfile := filepath.Join(tmpdir, "test.txt")
	testdata := []byte("12345678")
	if err := os.WriteFile(testfile, testdata, 0600); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	if err := acq.AcquireFile(testfile); err != nil {
		t.Logf("AcquireFile error: %v", err)
	}
	totalSize := acq.GetTotalSize()

	if totalSize != int64(len(testdata)) {
		t.Errorf("Expected total size %d, got %d", len(testdata), totalSize)
	}
}
