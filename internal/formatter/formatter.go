// Package formatter handles evidence package serialization and export
package formatter

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/evidex/internal/models"
	"github.com/evidex/internal/utils"
)

// PackageFormatter handles various output formats for evidence packages
type PackageFormatter struct {
	outputDir string
	package_  *models.EvidencePackage
}

// NewPackageFormatter creates a new formatter
func NewPackageFormatter(outputDir string, pkg *models.EvidencePackage) *PackageFormatter {
	return &PackageFormatter{
		outputDir: outputDir,
		package_:  pkg,
	}
}

// ExportJSON exports the evidence package as JSON metadata
func (pf *PackageFormatter) ExportJSON() error {
	metadataDir := filepath.Join(pf.outputDir, "metadata")
	if err := utils.EnsureDirectory(metadataDir); err != nil {
		return fmt.Errorf("failed to create metadata directory: %w", err)
	}

	// Create manifest JSON
	manifestFile := filepath.Join(metadataDir, "manifest.json")
	manifestData, err := json.MarshalIndent(pf.package_.Manifest, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal manifest: %w", err)
	}

	if err := os.WriteFile(manifestFile, manifestData, 0644); err != nil {
		return fmt.Errorf("failed to write manifest file: %w", err)
	}

	// Create acquisition log JSON
	logFile := filepath.Join(metadataDir, "acquisition_log.json")
	logData, err := json.MarshalIndent(pf.package_.AcquisitionLog, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal acquisition log: %w", err)
	}

	if err := os.WriteFile(logFile, logData, 0644); err != nil {
		return fmt.Errorf("failed to write acquisition log: %w", err)
	}

	// Create system context JSON
	contextFile := filepath.Join(metadataDir, "system_context.json")
	contextData, err := json.MarshalIndent(pf.package_.SystemContext, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal system context: %w", err)
	}

	if err := os.WriteFile(contextFile, contextData, 0644); err != nil {
		return fmt.Errorf("failed to write system context: %w", err)
	}

	// Create file metadata catalog
	catalogFile := filepath.Join(metadataDir, "file_catalog.json")
	catalogData, err := json.MarshalIndent(pf.package_.Files, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal file catalog: %w", err)
	}

	if err := os.WriteFile(catalogFile, catalogData, 0644); err != nil {
		return fmt.Errorf("failed to write file catalog: %w", err)
	}

	return nil
}

// ExportHashes creates a hashes file for integrity verification
func (pf *PackageFormatter) ExportHashes() error {
	hashesFile := filepath.Join(pf.outputDir, "HASHES.txt")
	hashesText := "EVIDEX FORENSIC HASHES\n"
	hashesText += fmt.Sprintf("Generated: %s\n", time.Now().Format(time.RFC3339))
	hashesText += fmt.Sprintf("Evidence ID: %s\n", pf.package_.Manifest.ID)
	hashesText += repeat("=", 80) + "\n\n"

	hashesText += fmt.Sprintf("File Count: %d\n", len(pf.package_.Files))
	hashesText += fmt.Sprintf("Total Size: %d bytes\n", pf.package_.Manifest.TotalSize)
	hashesText += "\n"

	hashesText += "FILES AND HASHES:\n"
	hashesText += repeat("-", 80) + "\n"

	for _, file := range pf.package_.Files {
		hashesText += fmt.Sprintf("\nFile: %s\n", file.SourcePath)
		hashesText += fmt.Sprintf("Size: %d bytes\n", file.FileSize)
		if file.Hashes != nil {
			hashesText += fmt.Sprintf("SHA-256: %s\n", file.Hashes.SHA256)
			if file.Hashes.SHA512 != "" {
				hashesText += fmt.Sprintf("SHA-512: %s\n", file.Hashes.SHA512)
			}
		}
	}

	if err := os.WriteFile(hashesFile, []byte(hashesText), 0644); err != nil {
		return fmt.Errorf("failed to write hashes file: %w", err)
	}

	return nil
}

// ExportCSV exports file metadata in CSV format for analysis
func (pf *PackageFormatter) ExportCSV() error {
	csvFile := filepath.Join(pf.outputDir, "metadata", "file_manifest.csv")

	// Ensure metadata directory exists
	if err := utils.EnsureDirectory(filepath.Dir(csvFile)); err != nil {
		return fmt.Errorf("failed to create metadata directory: %w", err)
	}

	file, err := os.Create(csvFile)
	if err != nil {
		return fmt.Errorf("failed to create CSV file: %w", err)
	}
	defer func() {
		_ = file.Close() //nolint:errcheck
	}()

	// Write CSV header
	header := "Filename,Source Path,Size (bytes),Modified Time,Hash (SHA-256),Verified,File Type,Owner\n"
	if _, err := file.WriteString(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write file entries
	for _, f := range pf.package_.Files {
		line := fmt.Sprintf("%s,%s,%d,%s,%s,%v,%s,%s\n",
			f.Filename,
			f.SourcePath,
			f.FileSize,
			f.ModifiedTime.Format(time.RFC3339),
			f.Hashes.SHA256,
			f.Verified,
			f.FileType,
			f.Owner,
		)
		if _, err := file.WriteString(line); err != nil {
			return fmt.Errorf("failed to write CSV line: %w", err)
		}
	}

	return nil
}

// CreateReadme creates an informational README file
func (pf *PackageFormatter) CreateReadme() error {
	readmeFile := filepath.Join(pf.outputDir, "README.txt")

	readme := fmt.Sprintf(`FORENSIC EVIDENCE PACKAGE
Generated by: Evidex v%s
Evidence ID: %s
Created: %s
Created By: %s@%s

PACKAGE CONTENTS:
- files/               : Acquired evidence files (unchanged)
- metadata/           : File metadata and hashes
  - manifest.json     : Chain of custody manifest
  - acquisition_log.json : Detailed acquisition log
  - system_context.json  : System information
  - file_catalog.json    : Complete file metadata
  - file_manifest.csv    : CSV file listing
- HASHES.txt          : File hashes for verification
- README.txt          : This file

INTEGRITY VERIFICATION:
All files in this package have been acquired in read-only mode without modification.
Cryptographic hashes (SHA-256 and SHA-512) are provided for integrity verification.

CHAIN OF CUSTODY:
The manifest.json file contains the complete chain of custody history.

FILE COUNT: %d
TOTAL SIZE: %d bytes
COMPRESSION: %s

For more information, see manifest.json and acquisition_log.json
`,
		pf.package_.Version,
		pf.package_.Manifest.ID,
		pf.package_.CreatedAt.Format(time.RFC3339),
		pf.package_.SystemContext.Username,
		pf.package_.SystemContext.Hostname,
		pf.package_.Manifest.FileCount,
		pf.package_.Manifest.TotalSize,
		pf.package_.Manifest.CompressionMethod,
	)

	if err := os.WriteFile(readmeFile, []byte(readme), 0644); err != nil {
		return fmt.Errorf("failed to write README: %w", err)
	}

	return nil
}

// CalculatePackageHash calculates a hash of the entire package metadata
func (pf *PackageFormatter) CalculatePackageHash() (string, error) {
	hash := sha256.New()

	// Hash manifest
	manifestJSON, _ := json.Marshal(pf.package_.Manifest)
	if _, err := hash.Write(manifestJSON); err != nil {
		return "", fmt.Errorf("failed to hash manifest: %w", err)
	}

	// Hash all file hashes
	for _, file := range pf.package_.Files {
		if file.Hashes != nil {
			if _, err := hash.Write([]byte(file.Hashes.SHA256)); err != nil {
				return "", fmt.Errorf("failed to hash file hash: %w", err)
			}
		}
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// ExportIntegrityReport creates a comprehensive integrity report
func (pf *PackageFormatter) ExportIntegrityReport() error {
	reportFile := filepath.Join(pf.outputDir, "INTEGRITY_REPORT.txt")

	report := fmt.Sprintf(`FORENSIC EVIDENCE INTEGRITY REPORT
Generated: %s
Evidence ID: %s

SUMMARY:
Total Files Acquired: %d
Total Size: %d bytes
Verified Files: %d
Failed Verifications: %d
Warnings: %d
Errors: %d

ACQUISITION METHOD: %s
HASH ALGORITHM: %s
SECONDARY ALGORITHM: %s

ACQUISITION LOG:
Start Time: %s
End Time: %s
Duration: %s

SYSTEM CONTEXT:
Operating System: %s %s
Hostname: %s
User: %s
Architecture: %s

CHAIN OF CUSTODY:
Initiator: %s
Role: Examiner
Status: %s

`,
		time.Now().Format(time.RFC3339),
		pf.package_.Manifest.ID,
		pf.package_.Manifest.FileCount,
		pf.package_.Manifest.TotalSize,
		countVerified(pf.package_.Files),
		countFailed(pf.package_.Files),
		pf.package_.AcquisitionLog.WarningCount,
		pf.package_.AcquisitionLog.ErrorCount,
		pf.package_.Manifest.AcquisitionMethod,
		pf.package_.Manifest.HashAlgorithm,
		pf.package_.Manifest.SecondaryAlgorithm,
		pf.package_.AcquisitionLog.StartTime.Format(time.RFC3339),
		pf.package_.AcquisitionLog.EndTime.Format(time.RFC3339),
		pf.package_.AcquisitionLog.Duration,
		pf.package_.SystemContext.OperatingSystem,
		pf.package_.SystemContext.OSVersion,
		pf.package_.SystemContext.Hostname,
		pf.package_.SystemContext.Username,
		pf.package_.SystemContext.Architecture,
		pf.package_.Manifest.CreatedByUser,
		pf.package_.Manifest.Integrity,
	)

	if err := os.WriteFile(reportFile, []byte(report), 0644); err != nil {
		return fmt.Errorf("failed to write integrity report: %w", err)
	}

	return nil
}

// CompressPackage compresses the package directory
func (pf *PackageFormatter) CompressPackage(compressionFormat string) error {
	// This would require archive/tar and compress/gzip
	// For now, just create a simple archive list

	listFile := filepath.Join(pf.outputDir, "PACKAGE_CONTENTS.txt")

	list := "EVIDEX PACKAGE CONTENTS\n"
	list += fmt.Sprintf("Evidence ID: %s\n", pf.package_.Manifest.ID)
	list += fmt.Sprintf("Generated: %s\n", time.Now().Format(time.RFC3339))
	list += repeat("=", 80) + "\n\n"

	// List directory structure
	err := filepath.Walk(pf.outputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(pf.outputDir, path)
		if info.IsDir() {
			list += fmt.Sprintf("[DIR]  %s/\n", rel)
		} else {
			list += fmt.Sprintf("[FILE] %s (%d bytes)\n", rel, info.Size())
		}
		return nil
	})

	if err == nil {
		if err := os.WriteFile(listFile, []byte(list), 0644); err != nil {
			return fmt.Errorf("failed to write file list: %w", err)
		}
	}

	return nil
}

// Helper functions
func countVerified(files []*models.FileEvidence) int {
	count := 0
	for _, f := range files {
		if f.Verified {
			count++
		}
	}
	return count
}

func countFailed(files []*models.FileEvidence) int {
	count := 0
	for _, f := range files {
		if !f.Verified {
			count++
		}
	}
	return count
}

// repeat creates a string repeated count times
func repeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}
