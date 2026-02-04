package tests

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/ilexum-group/evidex/internal/acquisition"
	"github.com/ilexum-group/evidex/internal/metadata"
	osWrapper "github.com/ilexum-group/evidex/internal/os"
	"github.com/ilexum-group/evidex/pkg/models"
)

// createTestAcquirer creates a test acquirer with all dependencies.
func createTestAcquirer(t *testing.T) *acquisition.Acquirer {
	t.Helper()

	// Create custody chain first
	custodyChain, err := models.NewCustodyChainEntry("evidex-test", "1.0.0")
	if err != nil {
		t.Fatalf("Failed to create custody chain: %v", err)
	}

	// Create OS wrapper
	osImpl := osWrapper.New()

	// Configure OS logger BEFORE using any OS operations
	osImpl.SetLogger(custodyChain.LogCommand)

	// Get system info (now logger is set)
	hostname, err := osImpl.Hostname()
	if err != nil {
		hostname = "test-host"
	}

	currentUser, err := osImpl.GetCurrentUser()
	if err != nil {
		currentUser = "test-user"
	}

	custodyChain.SetAgentHostname(hostname)
	custodyChain.SetAgentUser(currentUser)

	// Create metadata manager
	metadataMgr := metadata.NewMetadataManager(custodyChain.LogCommand)

	// Create acquirer
	return acquisition.NewAcquirer(custodyChain, osImpl, metadataMgr)
}

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

	acq := createTestAcquirer(t)

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

	acq := createTestAcquirer(t)
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

	acq := createTestAcquirer(t)
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

	acq := createTestAcquirer(t)
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

	acq := createTestAcquirer(t)
	if err := acq.AcquireFile(testfile); err != nil {
		t.Logf("AcquireFile warning: %v", err)
	}

	pkg := acq.GetEvidencePackage()

	if pkg == nil {
		t.Error("Expected EvidencePackage to be non-nil")
		return
	}

	if pkg.CustodyChain == nil {
		t.Error("Expected CustodyChain to be non-nil")
	}

	if pkg.CustodyChain.ID == "" {
		t.Error("Expected Evidence ID to be non-empty")
	}

	if len(pkg.Files) != 1 {
		t.Errorf("Expected 1 file in package, got %d", len(pkg.Files))
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

	acq := createTestAcquirer(t)

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

	acq := createTestAcquirer(t)

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
