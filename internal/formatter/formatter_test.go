package formatter

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/evidex/internal/models"
)

// TestNewPackageFormatter tests PackageFormatter initialization
func TestNewPackageFormatter(t *testing.T) {
	tmpdir, err := os.MkdirTemp("", "test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpdir); err != nil {
			t.Logf("Failed to remove temp directory: %v", err)
		}
	}()

	pkg := &models.EvidencePackage{
		Manifest: &models.ChainOfCustodyManifest{
			ID: "TEST-001",
		},
	}

	formatter := NewPackageFormatter(tmpdir, pkg)

	if formatter == nil {
		t.Error("Expected PackageFormatter to be non-nil")
	}
}

// TestExportJSON tests JSON export functionality
func TestExportJSON(t *testing.T) {
	tmpdir, err := os.MkdirTemp("", "test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpdir); err != nil {
			t.Logf("Failed to remove temp directory: %v", err)
		}
	}()

	now := time.Now()
	pkg := &models.EvidencePackage{
		Manifest: &models.ChainOfCustodyManifest{
			ID:                "TEST-001",
			CreatedAt:         now,
			CreatedByUser:     "test-user",
			CreatedByHostname: "test-host",
			FileCount:         1,
			TotalSize:         100,
			HashAlgorithm:     "sha256",
			AcquisitionMethod: "file",
		},
		Files: []*models.FileEvidence{
			{
				SourcePath:   "/test/file.txt",
				FileSize:     100,
				ModifiedTime: now,
				Hashes: &models.FileHashes{
					SHA256: "abc123",
				},
			},
		},
		AcquisitionLog: &models.AcquisitionLog{
			StartTime: now,
			EndTime:   now.Add(time.Second),
		},
	}

	formatter := NewPackageFormatter(tmpdir, pkg)
	err = formatter.ExportJSON()

	if err != nil {
		t.Fatalf("ExportJSON() error = %v", err)
	}

	// Verify manifest.json was created
	manifestPath := filepath.Join(tmpdir, "metadata", "manifest.json")
	if _, err := os.Stat(manifestPath); err != nil {
		t.Errorf("Expected manifest.json to be created: %v", err)
	}

	// Verify it's valid JSON
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatalf("Failed to read manifest.json: %v", err)
	}

	var manifest models.ChainOfCustodyManifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		t.Errorf("manifest.json is not valid JSON: %v", err)
	}
}

// TestExportCSV tests CSV export functionality
func TestExportCSV(t *testing.T) {
	tmpdir, err := os.MkdirTemp("", "test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpdir); err != nil {
			t.Logf("Failed to remove temp directory: %v", err)
		}
	}()

	now := time.Now()
	pkg := &models.EvidencePackage{
		Manifest: &models.ChainOfCustodyManifest{
			ID: "TEST-001",
		},
		Files: []*models.FileEvidence{
			{
				SourcePath:   "/test/file.txt",
				FileSize:     100,
				ModifiedTime: now,
				Hashes: &models.FileHashes{
					SHA256: "abc123",
				},
				Verified: true,
			},
		},
	}

	formatter := NewPackageFormatter(tmpdir, pkg)
	err = formatter.ExportCSV()

	if err != nil {
		t.Fatalf("ExportCSV() error = %v", err)
	}

	// Verify CSV file was created
	csvPath := filepath.Join(tmpdir, "metadata", "file_manifest.csv")
	if _, err := os.Stat(csvPath); err != nil {
		t.Errorf("Expected file_manifest.csv to be created: %v", err)
	}
}

// TestExportHashes tests hash file export
func TestExportHashes(t *testing.T) {
	tmpdir, err := os.MkdirTemp("", "test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpdir); err != nil {
			t.Logf("Failed to remove temp directory: %v", err)
		}
	}()

	pkg := &models.EvidencePackage{
		Manifest: &models.ChainOfCustodyManifest{
			ID: "TEST-001",
		},
		Files: []*models.FileEvidence{
			{
				SourcePath: "/test/file.txt",
				Hashes: &models.FileHashes{
					SHA256: "abc123",
					SHA512: "def456",
				},
			},
		},
	}

	formatter := NewPackageFormatter(tmpdir, pkg)
	err = formatter.ExportHashes()

	if err != nil {
		t.Fatalf("ExportHashes() error = %v", err)
	}

	// Verify HASHES.txt was created
	hashesPath := filepath.Join(tmpdir, "HASHES.txt")
	if _, err := os.Stat(hashesPath); err != nil {
		t.Errorf("Expected HASHES.txt to be created: %v", err)
	}
}

// TestExportIntegrityReport tests integrity report export
func TestExportIntegrityReport(t *testing.T) {
	tmpdir, err := os.MkdirTemp("", "test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpdir); err != nil {
			t.Logf("Failed to remove temp directory: %v", err)
		}
	}()

	now := time.Now()
	pkg := &models.EvidencePackage{
		Manifest: &models.ChainOfCustodyManifest{
			ID:                 "TEST-001",
			CreatedByUser:      "test-user",
			CreatedByHostname:  "test-host",
			FileCount:          1,
			TotalSize:          100,
			HashAlgorithm:      "SHA-256",
			SecondaryAlgorithm: "",
			CreatedAt:          now,
			AcquisitionMethod:  "file",
			Integrity:          "VERIFIED",
		},
		Files: []*models.FileEvidence{
			{
				SourcePath: "/test/file.txt",
				Hashes: &models.FileHashes{
					SHA256: "abc123",
				},
				Verified: true,
			},
		},
		AcquisitionLog: &models.AcquisitionLog{
			StartTime:    now,
			EndTime:      now.Add(time.Second),
			Duration:     "1s",
			WarningCount: 0,
			ErrorCount:   0,
		},
		SystemContext: &models.SystemContext{
			Hostname:        "test-host",
			Username:        "test-user",
			OperatingSystem: "linux",
			OSVersion:       "5.4.0",
			Architecture:    "x86_64",
		},
	}

	formatter := NewPackageFormatter(tmpdir, pkg)
	err = formatter.ExportIntegrityReport()

	if err != nil {
		t.Fatalf("ExportIntegrityReport() error = %v", err)
	}

	// Verify report was created
	reportPath := filepath.Join(tmpdir, "INTEGRITY_REPORT.txt")
	if _, err := os.Stat(reportPath); err != nil {
		t.Errorf("Expected INTEGRITY_REPORT.txt to be created: %v", err)
	}
}

// TestCreateReadme tests README generation
func TestCreateReadme(t *testing.T) {
	tmpdir, err := os.MkdirTemp("", "test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpdir); err != nil {
			t.Logf("Failed to remove temp directory: %v", err)
		}
	}()

	now := time.Now()
	pkg := &models.EvidencePackage{
		Version:   "1.0.0",
		CreatedAt: now,
		Manifest: &models.ChainOfCustodyManifest{
			ID:                "TEST-001",
			CreatedAt:         now,
			FileCount:         1,
			TotalSize:         100,
			CompressionMethod: "tar.gz",
		},
		SystemContext: &models.SystemContext{
			Username: "test-user",
			Hostname: "test-host",
		},
	}

	formatter := NewPackageFormatter(tmpdir, pkg)
	err = formatter.CreateReadme()

	if err != nil {
		t.Fatalf("CreateReadme() error = %v", err)
	}

	// Verify README was created
	readmePath := filepath.Join(tmpdir, "README.txt")
	if _, err := os.Stat(readmePath); err != nil {
		t.Errorf("Expected README.txt to be created: %v", err)
	}
}

// TestCalculatePackageHash tests package hash calculation
func TestCalculatePackageHash(t *testing.T) {
	pkg := &models.EvidencePackage{
		Manifest: &models.ChainOfCustodyManifest{
			ID: "TEST-001",
		},
	}

	formatter := NewPackageFormatter("", pkg)
	hash, err := formatter.CalculatePackageHash()

	if err != nil {
		t.Fatalf("CalculatePackageHash() error = %v", err)
	}

	if hash == "" {
		t.Error("Expected package hash to be non-empty")
	}

	// SHA256 hash should be 64 characters
	if len(hash) != 64 {
		t.Errorf("Expected hash length 64, got %d", len(hash))
	}
}

// TestCompressPackage tests package compression listing
func TestCompressPackage(t *testing.T) {
	tmpdir, err := os.MkdirTemp("", "test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpdir); err != nil {
			t.Logf("Failed to remove temp directory: %v", err)
		}
	}()

	pkg := &models.EvidencePackage{
		Manifest: &models.ChainOfCustodyManifest{
			ID: "TEST-001",
		},
		Files: []*models.FileEvidence{
			{
				SourcePath: "/test/file.txt",
				FileSize:   100,
			},
		},
	}

	formatter := NewPackageFormatter(tmpdir, pkg)
	err = formatter.CompressPackage("tar.gz")

	if err != nil {
		t.Fatalf("CompressPackage() error = %v", err)
	}

	// Verify package contents file was created
	contentsPath := filepath.Join(tmpdir, "PACKAGE_CONTENTS.txt")
	if _, err := os.Stat(contentsPath); err != nil {
		t.Errorf("Expected PACKAGE_CONTENTS.txt to be created: %v", err)
	}
}
