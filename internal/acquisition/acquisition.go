// Package acquisition handles evidence file acquisition and chain of custody.
package acquisition

import (
	"fmt"
	"io/fs"
	stdos "os"
	"path/filepath"
	"time"

	"github.com/ilexum-group/evidex/internal/metadata"
	osWrapper "github.com/ilexum-group/evidex/internal/os"
	"github.com/ilexum-group/evidex/internal/utils"
	"github.com/ilexum-group/evidex/pkg/models"
)

// Acquirer manages the evidence acquisition process.
type Acquirer struct {
	files        []*models.FileEvidence
	custodyChain *models.CustodyChainEntry
	osImpl       osWrapper.OS
	metadataMgr  *metadata.MetadataManager
}

// NewAcquirer creates a new evidence acquirer.
func NewAcquirer(custodyChainEntry *models.CustodyChainEntry, os osWrapper.OS, metadataManager *metadata.MetadataManager) *Acquirer {
	return &Acquirer{
		files:        make([]*models.FileEvidence, 0),
		custodyChain: custodyChainEntry,
		osImpl:       os,
		metadataMgr:  metadataManager,
	}
}

// AcquireFile adds a single file to the evidence package.
func (a *Acquirer) AcquireFile(filePath string) error {
	a.custodyChain.LogInfo("AcquireFile", fmt.Sprintf("Starting acquisition of: %s", filePath))

	// Verify file exists using OS wrapper
	fileInfo, err := a.osImpl.Stat(filePath)
	if err != nil {
		a.custodyChain.LogError("AcquireFile", filePath, err)
		return fmt.Errorf("file not found or inaccessible: %w", err)
	}

	a.custodyChain.LogInfo("FileValidation", fmt.Sprintf("File found: %s (size: %d bytes)", filePath, fileInfo.Size()))

	// Skip directories
	if fileInfo.IsDir() {
		errMsg := fmt.Sprintf("Skipping directory (not a file): %s", filePath)
		a.custodyChain.LogWarning("AcquireFile", errMsg)
		return fmt.Errorf("is a directory")
	}

	// Verify read-only access
	if _, err := utils.IsReadOnly(filePath); err != nil {
		a.custodyChain.LogWarning("AcquireFile", fmt.Sprintf("No read access to file: %s", filePath))
		return fmt.Errorf("no read access: %w", err)
	}

	// Create evidence structure
	evidence, err := a.createFileEvidence(filePath, fileInfo)
	if err != nil {
		return fmt.Errorf("failed to create evidence: %w", err)
	}

	// Extract metadata
	if err := a.extractMetadata(filePath, evidence); err != nil {
		return fmt.Errorf("metadata extraction failed: %w", err)
	}

	// Calculate hashes
	if err := a.calculateHashes(filePath, evidence); err != nil {
		a.custodyChain.LogWarning("HashCalculation", fmt.Sprintf("Hash calculation failed for %s: %v", filePath, err))
	}

	// Store evidence
	a.files = append(a.files, evidence)

	a.custodyChain.LogInfo("FileAcquired", fmt.Sprintf("Successfully acquired: %s (%d bytes)",
		filePath, evidence.FileSize))

	return nil
}

// createFileEvidence creates the basic evidence structure for a file.
func (a *Acquirer) createFileEvidence(filePath string, fileInfo fs.FileInfo) (*models.FileEvidence, error) {
	evidence := &models.FileEvidence{
		SourcePath:      filePath,
		Filename:        fileInfo.Name(),
		FileSize:        fileInfo.Size(),
		ModifiedTime:    fileInfo.ModTime(),
		AcquisitionTime: time.Now(),
		FileType:        a.metadataMgr.GetFileTypeFromMimeType(filePath),
	}

	// Extract advanced timestamps and ownership (OS-specific)
	accessTime, modTime, changeTime, createTime := a.osImpl.ExtractAdvancedTimes(filePath)
	evidence.AccessedTime = accessTime
	evidence.ChangeTime = changeTime
	evidence.CreatedTime = createTime
	evidence.ModifiedTime = modTime

	username, groupname := a.osImpl.ExtractOwnershipInfo(filePath)
	evidence.Owner = username
	evidence.Group = groupname

	a.custodyChain.LogInfo("TimestampExtraction", fmt.Sprintf("Timestamps extracted for: %s", filePath))

	return evidence, nil
}

// extractMetadata extracts file-specific metadata.
func (a *Acquirer) extractMetadata(filePath string, evidence *models.FileEvidence) error {
	a.custodyChain.LogInfo("MetadataExtraction", fmt.Sprintf("Extracting metadata from: %s", filePath))

	err := a.metadataMgr.ExtractFileMetadata(evidence)
	if err != nil {
		a.custodyChain.LogError("ExtractMetadata", filePath, err)
		return err
	}

	return nil
}

// calculateHashes computes cryptographic hashes for a file.
func (a *Acquirer) calculateHashes(filePath string, evidence *models.FileEvidence) error {
	a.custodyChain.LogInfo("HashCalculation", fmt.Sprintf("Calculating hashes (MD5, SHA1, SHA256) for: %s", filePath))

	hashes, err := a.metadataMgr.CalculateFileHashes(filePath)
	if err != nil {
		a.custodyChain.LogError("CalculateHashes", filePath, err)
		return err
	}

	evidence.Hashes = hashes

	// Log hash preview
	a.custodyChain.LogInfo("HashCalculation",
		fmt.Sprintf("Hashes calculated - MD5: %s, SHA1: %s, SHA256: %s",
			truncateHash(hashes.MD5, 8),
			truncateHash(hashes.SHA1, 8),
			truncateHash(hashes.SHA256, 8)))

	return nil
}

// truncateHash safely truncates a hash string to the specified length.
func truncateHash(hash string, length int) string {
	if len(hash) > length {
		return hash[:length]
	}
	return hash
}

// AcquireDirectory recursively acquires all files in a directory.
func (a *Acquirer) AcquireDirectory(dirPath string, recursive bool) error {
	a.custodyChain.LogInfo("AcquireDirectory", fmt.Sprintf("Starting directory acquisition: %s (recursive: %v)", dirPath, recursive))

	if recursive {
		return a.acquireDirectoryRecursive(dirPath)
	}
	return a.acquireDirectoryFlat(dirPath)
}

// acquireDirectoryRecursive recursively walks and acquires all files in a directory tree.
func (a *Acquirer) acquireDirectoryRecursive(dirPath string) error {
	a.custodyChain.LogInfo("DirectoryWalk", fmt.Sprintf("Walking directory tree: %s", dirPath))

	walkErr := filepath.Walk(dirPath, func(path string, info stdos.FileInfo, err error) error {
		if err != nil {
			a.custodyChain.LogWarning("DirectoryWalk", fmt.Sprintf("Error accessing: %s (%v)", path, err))
			return nil // Continue walking
		}

		if !info.IsDir() {
			if err := a.AcquireFile(path); err != nil {
				a.custodyChain.LogError("AcquireFile", path, err)
			}
		}
		return nil
	})

	return walkErr
}

// acquireDirectoryFlat acquires only files directly in the directory (non-recursive).
func (a *Acquirer) acquireDirectoryFlat(dirPath string) error {
	a.custodyChain.LogInfo("DirectoryScan", fmt.Sprintf("Scanning directory (non-recursive): %s", dirPath))

	entries, err := stdos.ReadDir(dirPath)
	if err != nil {
		a.custodyChain.LogError("ReadDirectory", dirPath, err)
		return fmt.Errorf("failed to read directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			filePath := filepath.Join(dirPath, entry.Name())
			if err := a.AcquireFile(filePath); err != nil {
				a.custodyChain.LogError("AcquireFile", filePath, err)
			}
		}
	}

	a.custodyChain.LogInfo("DirectoryComplete", fmt.Sprintf("Directory acquisition completed: %s", dirPath))
	return nil
}

// AcquireMultiple acquires multiple files from a list.
func (a *Acquirer) AcquireMultiple(filePaths []string) error {
	a.custodyChain.LogInfo("AcquireMultiple", fmt.Sprintf("Starting acquisition of %d files", len(filePaths)))
	for _, filePath := range filePaths {
		if err := a.AcquireFile(filePath); err != nil {
			a.custodyChain.LogError("AcquireMultiple", filePath, err)
		}
	}
	a.custodyChain.LogInfo("AcquireMultiple", fmt.Sprintf("Completed acquisition of %d files", len(filePaths)))
	return nil
}

// GetEvidencePackage builds the complete evidence package.
func (a *Acquirer) GetEvidencePackage() *models.EvidencePackage {
	a.custodyChain.LogInfo("PackageCreation", "Building evidence package with custody chain")
	a.custodyChain.EndTimestamp = time.Now()
	a.custodyChain.Duration = a.custodyChain.EndTimestamp.Sub(a.custodyChain.StartTimestamp).String()

	totalSize := int64(0)
	for _, f := range a.files {
		totalSize += f.FileSize
	}
	a.custodyChain.TotalSizeBytes = totalSize
	a.custodyChain.ItemCount = len(a.files)

	a.custodyChain.LogInfo("PackageStats", fmt.Sprintf("Total files: %d, Total size: %d bytes, Duration: %s",
		len(a.files), a.GetTotalSize(), a.custodyChain.Duration))

	return &models.EvidencePackage{
		Files:        a.files,
		CustodyChain: a.custodyChain,
	}
}

// GetFileCount returns the number of acquired files.
func (a *Acquirer) GetFileCount() int {
	return len(a.files)
}

// GetTotalSize returns the total size of acquired files.
func (a *Acquirer) GetTotalSize() int64 {
	total := int64(0)
	for _, f := range a.files {
		total += f.FileSize
	}
	return total
}
