// Package acquisition handles evidence file acquisition and chain of custody
package acquisition

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ilexum-group/evidex/internal/metadata"
	"github.com/ilexum-group/evidex/internal/models"
	"github.com/ilexum-group/evidex/internal/utils"
)

// Acquirer manages the evidence acquisition process
type Acquirer struct {
	outputDir        string
	files            []*models.FileEvidence
	acquisitionLog   *models.AcquisitionLog
	primaryHashAlg   string
	secondaryHashAlg string
	validateAccess   bool
}

// NewAcquirer creates a new evidence acquirer
func NewAcquirer(outputDir string) *Acquirer {
	return &Acquirer{
		outputDir:        outputDir,
		files:            make([]*models.FileEvidence, 0),
		acquisitionLog:   createAcquisitionLog(),
		primaryHashAlg:   "SHA-256",
		secondaryHashAlg: "SHA-512",
		validateAccess:   true,
	}
}

// createAcquisitionLog initializes the acquisition log
func createAcquisitionLog() *models.AcquisitionLog {
	return &models.AcquisitionLog{
		StartTime:      time.Now(),
		Entries:        make([]models.LogEntry, 0),
		OperationCount: 0,
		ErrorCount:     0,
		WarningCount:   0,
	}
}

// AcquireFile adds a single file to the evidence package
func (a *Acquirer) AcquireFile(filePath string) error {
	// Verify file exists
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		a.logError("AcquireFile", filePath, err)
		return fmt.Errorf("file not found or inaccessible: %s", filePath)
	}

	// Skip directories
	if fileInfo.IsDir() {
		errMsg := fmt.Sprintf("Skipping directory (not a file): %s", filePath)
		a.logWarning("AcquireFile", errMsg)
		return fmt.Errorf("is a directory")
	}

	// Verify read-only access
	if a.validateAccess {
		if _, err := utils.IsReadOnly(filePath); err != nil {
			errMsg := fmt.Sprintf("No read access to file: %s", filePath)
			a.logWarning("AcquireFile", errMsg)
			return fmt.Errorf("no read access")
		}
	}

	// Log acquisition start
	a.logInfo("AcquireFile", fmt.Sprintf("Starting acquisition of: %s", filePath))

	// Extract metadata
	evidence, err := metadata.ExtractFileMetadata(filePath)
	if err != nil {
		a.logError("ExtractMetadata", filePath, err)
		return fmt.Errorf("metadata extraction failed: %w", err)
	}

	// Calculate hashes
	hashes, err := metadata.CalculateFileHashes(filePath)
	if err != nil {
		a.logError("CalculateHashes", filePath, err)
		evidence.VerificationErr = err.Error()
		evidence.Verified = false
	} else {
		evidence.Hashes = hashes
		evidence.Verified = true
		a.logInfo("HashCalculation",
			fmt.Sprintf("SHA-256: %s", evidence.Hashes.SHA256))
	}

	// Store evidence
	a.files = append(a.files, evidence)
	a.acquisitionLog.OperationCount++

	a.logInfo("FileAcquired", fmt.Sprintf("Successfully acquired: %s (%d bytes)",
		filePath, evidence.FileSize))

	return nil
}

// AcquireDirectory recursively acquires all files in a directory
func (a *Acquirer) AcquireDirectory(dirPath string, recursive bool) error {
	a.logInfo("AcquireDirectory", fmt.Sprintf("Starting acquisition of directory: %s (recursive: %v)",
		dirPath, recursive))

	if recursive {
		return filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				a.logWarning("DirectoryWalk", fmt.Sprintf("Error accessing: %s (%v)", path, err))
				return nil // Continue walking
			}

			if !info.IsDir() {
				if err := a.AcquireFile(path); err != nil {
					// Log but continue
					a.acquisitionLog.ErrorCount++
				}
			}
			return nil
		})
	}

	// Non-recursive: only files directly in directory
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		a.logError("ReadDirectory", dirPath, err)
		return fmt.Errorf("failed to read directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			filePath := filepath.Join(dirPath, entry.Name())
			if err := a.AcquireFile(filePath); err != nil {
				a.acquisitionLog.ErrorCount++
			}
		}
	}

	return nil
}

// AcquireMultiple acquires multiple files from a list
func (a *Acquirer) AcquireMultiple(filePaths []string) error {
	for _, filePath := range filePaths {
		if err := a.AcquireFile(filePath); err != nil {
			a.logError("AcquireMultiple", filePath, err)
			a.acquisitionLog.ErrorCount++
		}
	}
	return nil
}

// CopyFilesToPackage copies acquired files to the output directory
func (a *Acquirer) CopyFilesToPackage() error {
	filesDir := filepath.Join(a.outputDir, "files")
	if err := utils.EnsureDirectory(filesDir); err != nil {
		a.logError("CreateFilesDirectory", filesDir, err)
		return fmt.Errorf("failed to create files directory: %w", err)
	}

	for _, evidence := range a.files {
		// Create destination path structure
		relativePath := evidence.Filename
		destPath := filepath.Join(filesDir, relativePath)

		// Ensure destination directory exists
		destDir := filepath.Dir(destPath)
		if err := utils.EnsureDirectory(destDir); err != nil {
			a.logWarning("CreateDestinationDir", fmt.Sprintf("Failed to create: %s", destDir))
			continue
		}

		// Copy file (non-destructive, read-only)
		if err := utils.SafeCopyFile(evidence.SourcePath, destPath); err != nil {
			a.logError("CopyFile", fmt.Sprintf("%s -> %s", evidence.SourcePath, destPath), err)
			a.acquisitionLog.ErrorCount++
			continue
		}

		// Store relative path in evidence
		evidence.RelativePath = relativePath
		a.logInfo("FileCopied", fmt.Sprintf("Copied to: %s", relativePath))
	}

	return nil
}

// GetEvidencePackage builds the complete evidence package
func (a *Acquirer) GetEvidencePackage() *models.EvidencePackage {
	a.acquisitionLog.EndTime = time.Now()
	a.acquisitionLog.Duration = a.acquisitionLog.EndTime.Sub(a.acquisitionLog.StartTime).String()

	manifest := createManifest(a.files, a.acquisitionLog)
	systemContext := utils.GetSystemContext()

	return &models.EvidencePackage{
		Manifest:       manifest,
		Files:          a.files,
		AcquisitionLog: a.acquisitionLog,
		SystemContext:  systemContext,
		CreatedAt:      time.Now(),
		Version:        "1.0.0",
	}
}

// createManifest creates the chain of custody manifest
func createManifest(files []*models.FileEvidence, log *models.AcquisitionLog) *models.ChainOfCustodyManifest {
	hostname, _ := os.Hostname()
	username, _ := utils.GetCurrentUser()

	totalSize := int64(0)
	for _, f := range files {
		totalSize += f.FileSize
	}

	manifest := &models.ChainOfCustodyManifest{
		ID:                 utils.GenerateEvidenceID(),
		CreatedAt:          time.Now(),
		CreatedByUser:      username,
		CreatedByHostname:  hostname,
		FileCount:          len(files),
		TotalSize:          totalSize,
		CompressionMethod:  "none",
		HashAlgorithm:      "SHA-256",
		SecondaryAlgorithm: "SHA-512",
		AcquisitionMethod:  "read-only",
		Integrity:          "verified",
		Custodians: []models.Custodian{
			{
				Name:      username,
				Role:      "Examiner",
				Action:    "created",
				Timestamp: time.Now(),
				Location:  hostname,
				Notes:     "Initial evidence acquisition",
			},
		},
	}

	return manifest
}

// Logging helpers
func (a *Acquirer) logInfo(operation, message string) {
	utils.LogInfo(message, map[string]string{
		"operation": operation,
	})
	entry := models.LogEntry{
		Timestamp: time.Now(),
		Level:     "INFO",
		Message:   message,
		Details:   operation,
	}
	a.acquisitionLog.Entries = append(a.acquisitionLog.Entries, entry)
}

func (a *Acquirer) logWarning(operation, message string) {
	utils.LogWarn(message, map[string]string{
		"operation": operation,
	})
	a.acquisitionLog.WarningCount++
	entry := models.LogEntry{
		Timestamp: time.Now(),
		Level:     "WARNING",
		Message:   message,
		Details:   operation,
	}
	a.acquisitionLog.Entries = append(a.acquisitionLog.Entries, entry)
}

func (a *Acquirer) logError(operation, message string, err error) {
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}
	utils.LogError(message, map[string]string{
		"operation": operation,
		"error":     errMsg,
	})
	a.acquisitionLog.ErrorCount++
	entry := models.LogEntry{
		Timestamp: time.Now(),
		Level:     "ERROR",
		Message:   message,
		Details:   operation,
		Error:     errMsg,
	}
	a.acquisitionLog.Entries = append(a.acquisitionLog.Entries, entry)
}

// SetHashAlgorithm sets the primary hash algorithm
func (a *Acquirer) SetHashAlgorithm(algorithm string) {
	a.primaryHashAlg = algorithm
}

// GetFileCount returns the number of acquired files
func (a *Acquirer) GetFileCount() int {
	return len(a.files)
}

// GetTotalSize returns the total size of acquired files
func (a *Acquirer) GetTotalSize() int64 {
	total := int64(0)
	for _, f := range a.files {
		total += f.FileSize
	}
	return total
}
