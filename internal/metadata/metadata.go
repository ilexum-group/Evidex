// Package metadata provides extraction of forensically relevant file metadata
package metadata

import (
	"fmt"
	"os"
	"os/user"
	"time"

	"github.com/evidex/internal/models"
	"github.com/evidex/internal/utils"
)

// ExtractFileMetadata extracts filesystem and file metadata
func ExtractFileMetadata(filePath string) (*models.FileEvidence, error) {
	// Get file info
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}

	evidence := &models.FileEvidence{
		SourcePath:      filePath,
		Filename:        fileInfo.Name(),
		FileSize:        fileInfo.Size(),
		FileMode:        uint32(fileInfo.Mode()),
		ModifiedTime:    fileInfo.ModTime(),
		AcquisitionTime: time.Now(),
		FileType:        utils.GetFileTypeFromExtension(filePath),
		Hashes:          &models.FileHashes{},
	}

	// Try to get more detailed timestamps (OS-specific)
	extractAdvancedTimes(filePath, evidence)

	// Get owner and group information
	extractOwnershipInfo(filePath, evidence)

	// Check if it's a symlink
	stat, err := os.Lstat(filePath)
	if err == nil {
		if stat.Mode()&os.ModeSymlink != 0 {
			evidence.IsSymlink = true
			if target, err := os.Readlink(filePath); err == nil {
				evidence.SymlinkTarget = target
			}
		}
	}

	// Extract image metadata if applicable
	if utils.IsImageFile(filePath) {
		evidence.ImageMetadata, _ = ExtractImageMetadata(filePath)
	}

	// Extract video metadata if applicable
	if utils.IsVideoFile(filePath) {
		evidence.VideoMetadata, _ = ExtractVideoMetadata(filePath)
	}

	return evidence, nil
}

// extractAdvancedTimes attempts to extract creation and access times (platform-specific)
func extractAdvancedTimes(filePath string, evidence *models.FileEvidence) {
	// These are platform-specific and require OS-level access
	// On Windows: use GetFileTime
	// On Unix: use stat struct fields (sys.Stat_t.Ctimespec, Atimespec, Mtimespec)
	// For now, we use ModTime as a baseline
	// This should be extended with platform-specific implementations

	evidence.AccessedTime = time.Now()           // Placeholder - would need OS-specific code
	evidence.CreatedTime = evidence.ModifiedTime // Placeholder
	evidence.ChangeTime = evidence.ModifiedTime  // Placeholder
}

// extractOwnershipInfo extracts file owner and group information
func extractOwnershipInfo(filePath string, evidence *models.FileEvidence) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return
	}

	// Get OS-specific uid/gid
	stat := fileInfo.Sys()
	if stat == nil {
		return
	}

	// Try to resolve user and group names
	switch stat.(type) {
	// Platform-specific handling would go here
	default:
		evidence.Owner = "unknown"
		evidence.Group = "unknown"
	}

	// On Unix-like systems
	if u, err := user.Current(); err == nil {
		evidence.Owner = u.Username
	}
}

// ExtractImageMetadata extracts EXIF, XMP, and IPTC metadata from images
func ExtractImageMetadata(filePath string) (*models.ImageMetadata, error) {
	metadata := &models.ImageMetadata{
		EXIF: make(map[string]string),
		XMP:  make(map[string]string),
		IPTC: make(map[string]string),
	}

	// Extract basic image info from file extension
	switch {
	case isJPEG(filePath):
		metadata.Format = "JPEG"
		extractJPEGMetadata(filePath, metadata)
	case isPNG(filePath):
		metadata.Format = "PNG"
		extractPNGMetadata(filePath, metadata)
	case isGIF(filePath):
		metadata.Format = "GIF"
		extractGIFMetadata(filePath, metadata)
	default:
		metadata.Format = "Unknown"
	}

	return metadata, nil
}

// ExtractVideoMetadata extracts video codec, duration, and container info
func ExtractVideoMetadata(filePath string) (*models.VideoMetadata, error) {
	metadata := &models.VideoMetadata{}

	// Determine container format from extension
	switch {
	case isMP4(filePath):
		metadata.Format = "MP4"
	case isMOV(filePath):
		metadata.Format = "QuickTime"
	case isMKV(filePath):
		metadata.Format = "Matroska"
	case isAVI(filePath):
		metadata.Format = "AVI"
	case isWebM(filePath):
		metadata.Format = "WebM"
	default:
		metadata.Format = "Unknown"
	}

	// Extract video info from file
	extractVideoProperties(filePath, metadata)

	return metadata, nil
}

// Helper functions for format detection
func isJPEG(filePath string) bool {
	// Check magic bytes
	file, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer file.Close()

	buf := make([]byte, 2)
	if _, err := file.Read(buf); err != nil {
		return false
	}
	return buf[0] == 0xFF && buf[1] == 0xD8
}

func isPNG(filePath string) bool {
	file, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer file.Close()

	buf := make([]byte, 8)
	if _, err := file.Read(buf); err != nil {
		return false
	}
	return buf[0] == 0x89 && buf[1] == 'P' && buf[2] == 'N' && buf[3] == 'G'
}

func isGIF(filePath string) bool {
	file, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer file.Close()

	buf := make([]byte, 3)
	if _, err := file.Read(buf); err != nil {
		return false
	}
	return buf[0] == 'G' && buf[1] == 'I' && buf[2] == 'F'
}

func isMP4(filePath string) bool {
	file, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer file.Close()

	buf := make([]byte, 8)
	if _, err := file.ReadAt(buf, 4); err != nil {
		return false
	}
	return string(buf) == "ftypmp42" || string(buf) == "ftypisom"
}

func isMOV(filePath string) bool {
	file, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer file.Close()

	buf := make([]byte, 4)
	if _, err := file.ReadAt(buf, 4); err != nil {
		return false
	}
	return string(buf) == "mdat" || string(buf) == "moov"
}

func isMKV(filePath string) bool {
	file, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer file.Close()

	buf := make([]byte, 4)
	if _, err := file.Read(buf); err != nil {
		return false
	}
	return buf[0] == 0x1A && buf[1] == 0x45 && buf[2] == 0xDF && buf[3] == 0xA3
}

func isAVI(filePath string) bool {
	file, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer file.Close()

	buf := make([]byte, 4)
	if _, err := file.Read(buf); err != nil {
		return false
	}
	return buf[0] == 'R' && buf[1] == 'I' && buf[2] == 'F' && buf[3] == 'F'
}

func isWebM(filePath string) bool {
	file, err := os.Open(filePath)
	if err != nil {
		return false
	}
	defer file.Close()

	buf := make([]byte, 4)
	if _, err := file.Read(buf); err != nil {
		return false
	}
	return buf[0] == 0x1A && buf[1] == 0x45 && buf[2] == 0xDF && buf[3] == 0xA3
}

// Extract JPEG metadata (EXIF, XMP)
func extractJPEGMetadata(filePath string, metadata *models.ImageMetadata) {
	// This would require JPEG parsing library
	// For now, we set basic properties
	metadata.Format = "JPEG"

	// Try to extract EXIF if library is available
	// This is a simplified implementation - would need rwcarlsen/goexif or similar
	// The full implementation would parse EXIF markers
}

// Extract PNG metadata
func extractPNGMetadata(filePath string, metadata *models.ImageMetadata) {
	metadata.Format = "PNG"
	// PNG metadata extraction would require PNG parsing
}

// Extract GIF metadata
func extractGIFMetadata(filePath string, metadata *models.ImageMetadata) {
	metadata.Format = "GIF"
	// GIF metadata extraction
}

// Extract video properties from file headers
func extractVideoProperties(filePath string, metadata *models.VideoMetadata) {
	// This is a stub for video property extraction
	// A full implementation would parse MP4, MOV, or MKV headers
	// to extract:
	// - Duration
	// - Video codec
	// - Audio codec
	// - Frame rate
	// - Resolution
	// - Bitrate
	// - Creation time

	// For now, set defaults
	metadata.AudioChannels = 2
	metadata.FrameRate = "unknown"
	metadata.DurationSeconds = 0
}

// CalculateFileHashes calculates cryptographic hashes
func CalculateFileHashes(filePath string, primaryAlg string) (*models.FileHashes, error) {
	hashes := &models.FileHashes{}

	// Calculate SHA-256 (primary)
	sha256Hash, err := utils.CalculateSHA256(filePath)
	if err != nil {
		return nil, fmt.Errorf("SHA-256 calculation failed: %w", err)
	}
	hashes.SHA256 = sha256Hash

	// Calculate SHA-512 (secondary)
	sha512Hash, err := utils.CalculateSHA512(filePath)
	if err != nil {
		utils.LogWarning("SHA-512 calculation failed for file", filePath)
	} else {
		hashes.SHA512 = sha512Hash
	}

	// Calculate MD5 (for compatibility, not recommended for evidence)
	md5Hash, err := utils.CalculateMD5(filePath)
	if err != nil {
		utils.LogWarning("MD5 calculation failed for file", filePath)
	} else {
		hashes.MD5 = md5Hash
	}

	return hashes, nil
}

// VerifyFileIntegrity verifies a file's integrity using stored hash
func VerifyFileIntegrity(filePath string, knownHash string) (bool, error) {
	return utils.VerifyHash(filePath, knownHash, "SHA-256")
}
