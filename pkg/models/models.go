// Package models defines the data structures for Evidex forensic evidence acquisition.
// All models are designed to maintain chain of custody and forensic integrity.
package models

import (
	"time"
)

// EvidencePackage represents the complete forensic evidence package.
type EvidencePackage struct {
	CaseID       string             `json:"case_id,omitempty"` // Case ID for grouping evidence packages
	Files        []*FileEvidence    `json:"files"`
	CustodyChain *CustodyChainEntry `json:"custody_chain"` // Custody Chain - Complete digital evidence custody tracking
}

// FileEvidence represents a single evidence file with complete metadata.
type FileEvidence struct {
	SourcePath      string            `json:"source_path"`                // Original file path
	RelativePath    string            `json:"relative_path"`              // Path within evidence package
	Filename        string            `json:"filename"`                   // Original filename
	FileSize        int64             `json:"file_size"`                  // Size in bytes
	FileMode        uint32            `json:"file_mode"`                  // Unix permissions
	AccessedTime    time.Time         `json:"accessed_time"`              // Access time
	ModifiedTime    time.Time         `json:"modified_time"`              // Modification time
	CreatedTime     time.Time         `json:"created_time"`               // Creation/birth time
	ChangeTime      time.Time         `json:"change_time"`                // Metadata change time
	Owner           string            `json:"owner"`                      // File owner
	Group           string            `json:"group"`                      // File group
	FileType        string            `json:"file_type"`                  // MIME type
	IsSymlink       bool              `json:"is_symlink"`                 // Is symbolic link
	SymlinkTarget   string            `json:"symlink_target"`             // Symlink target if applicable
	Hashes          *FileHashes       `json:"hashes"`                     // Cryptographic hashes
	ImageMetadata   *ImageMetadata    `json:"image_metadata"`             // EXIF/XMP for images
	VideoMetadata   *VideoMetadata    `json:"video_metadata"`             // Codec info for videos
	GenericMetadata map[string]string `json:"generic_metadata,omitempty"` // Other file types
	AcquisitionTime time.Time         `json:"acquisition_time"`           // When file was acquired
	FileContent     []byte            `json:"file_content"`               // File content for transmission (read-only, never modifies source)
}

// FileHashes contains cryptographic hashes for integrity verification.
type FileHashes struct {
	MD5    string `json:"md5"`    // MD5 hash (128-bit) - for compatibility with legacy systems
	SHA1   string `json:"sha1"`   // SHA1 hash (160-bit) - for legacy compatibility
	SHA256 string `json:"sha256"` // SHA256 hash (256-bit) - primary hash for evidence integrity
	SHA512 string `json:"sha512"` // SHA512 hash (512-bit) - secondary hash (optional)
}

// ImageMetadata contains image-specific metadata.
// Reduced to the fields we currently populate.
type ImageMetadata struct {
	Format string `json:"format"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

// GPSCoordinates represents GPS location data.
type GPSCoordinates struct {
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Altitude  float64   `json:"altitude"`
	Timestamp time.Time `json:"timestamp"`
}

// VideoMetadata contains video-specific metadata.
type VideoMetadata struct {
	Format          string    `json:"format"`            // Container format (MP4, MOV, AVI, etc.)
	Duration        string    `json:"duration"`          // Duration (HH:MM:SS)
	DurationSeconds float64   `json:"duration_seconds"`  // Duration in seconds
	VideoCodec      string    `json:"video_codec"`       // Video codec (h264, h265, etc.)
	AudioCodec      string    `json:"audio_codec"`       // Audio codec (aac, mp3, etc.)
	Width           int       `json:"width"`             // Video width in pixels
	Height          int       `json:"height"`            // Video height in pixels
	FrameRate       string    `json:"frame_rate"`        // Frame rate (fps)
	BitRate         string    `json:"bit_rate"`          // Overall bitrate
	VideoBitRate    string    `json:"video_bit_rate"`    // Video bitrate
	AudioBitRate    string    `json:"audio_bit_rate"`    // Audio bitrate
	AudioChannels   int       `json:"audio_channels"`    // Number of audio channels
	AudioSampleRate string    `json:"audio_sample_rate"` // Audio sample rate
	Encoder         string    `json:"encoder"`           // Encoder software
	CreationTime    time.Time `json:"creation_time"`     // Media creation time
	Title           string    `json:"title"`             // Video title
	Description     string    `json:"description"`       // Video description
	Copyright       string    `json:"copyright"`         // Copyright notice
}

// SystemContext records system information during acquisition.
type SystemContext struct {
	Hostname         string    `json:"hostname"`          // System hostname
	Username         string    `json:"username"`          // User running acquisition
	OperatingSystem  string    `json:"operating_system"`  // OS name
	OSVersion        string    `json:"os_version"`        // OS version
	Architecture     string    `json:"architecture"`      // CPU architecture
	TimeZone         string    `json:"timezone"`          // System timezone
	LocalTime        time.Time `json:"local_time"`        // Local acquisition time
	UTCTime          time.Time `json:"utc_time"`          // UTC acquisition time
	BinaryVersion    string    `json:"binary_version"`    // Evidex version
	BinaryPath       string    `json:"binary_path"`       // Path to binary
	BinaryHash       string    `json:"binary_hash"`       // Hash of binary itself
	WorkingDirectory string    `json:"working_directory"` // CWD during execution
}
