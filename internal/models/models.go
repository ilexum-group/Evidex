// Package models defines the data structures for Evidex forensic evidence acquisition.
// All models are designed to maintain chain of custody and forensic integrity.
package models

import (
	"time"
)

// EvdencePackage represents the complete forensic evidence package
type EvidencePackage struct {
	Manifest       *ChainOfCustodyManifest `json:"manifest"`
	Files          []*FileEvidence         `json:"files"`
	AcquisitionLog *AcquisitionLog         `json:"acquisition_log"`
	SystemContext  *SystemContext          `json:"system_context"`
	CreatedAt      time.Time               `json:"created_at"`
	Version        string                  `json:"version"`
	Logs           []string                `json:"logs"`       // RFC 5424 syslog entries
	ServerURL      string                  `json:"server_url"` // Remote endpoint for transmission
	AuthToken      string                  `json:"auth_token"` // Bearer token for server authentication
}

// ChainOfCustodyManifest records all actions and modifications to the evidence
type ChainOfCustodyManifest struct {
	ID                  string      `json:"id"`                   // Unique evidence package identifier
	CreatedAt           time.Time   `json:"created_at"`           // Package creation timestamp
	CreatedByUser       string      `json:"created_by_user"`      // User who initiated acquisition
	CreatedByHostname   string      `json:"created_by_hostname"`  // Hostname where binary ran
	Signature           string      `json:"signature"`            // HMAC signature of the package
	FileCount           int         `json:"file_count"`           // Total files in package
	TotalSize           int64       `json:"total_size"`           // Total size in bytes
	CompressionMethod   string      `json:"compression_method"`   // Compression algorithm used
	HashAlgorithm       string      `json:"hash_algorithm"`       // Primary hash algorithm (SHA-256, SHA-512)
	SecondaryAlgorithm  string      `json:"secondary_algorithm"`  // Secondary hash (optional)
	AcquisitionMethod   string      `json:"acquisition_method"`   // read-only, file copy, etc.
	EvidenceDescription string      `json:"evidence_description"` // Case-specific notes
	ExaminersNotes      string      `json:"examiners_notes"`      // Additional context
	Integrity           string      `json:"integrity"`            // Integrity status (verified, unverified)
	Custodians          []Custodian `json:"custodians"`           // Chain of custody history
	Hash                string      `json:"hash"`                 // Hash of the entire manifest
}

// Custodian tracks who had custody of evidence and when
type Custodian struct {
	Name      string    `json:"name"`      // Custodian name
	Role      string    `json:"role"`      // Role/title
	Signature string    `json:"signature"` // Digital signature
	Action    string    `json:"action"`    // received, transferred, archived
	Timestamp time.Time `json:"timestamp"` // When action occurred
	Location  string    `json:"location"`  // Physical or logical location
	Notes     string    `json:"notes"`     // Optional notes about custody action
}

// FileEvidence represents a single evidence file with complete metadata
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
	Verified        bool              `json:"verified"`                   // Hash verification status
	VerificationErr string            `json:"verification_err"`           // Verification error if any
	FileContent     []byte            `json:"file_content"`               // File content for transmission (read-only, never modifies source)
}

// FileHashes contains cryptographic hashes for integrity verification
type FileHashes struct {
	SHA256 string `json:"sha256"` // Primary hash
	SHA512 string `json:"sha512"` // Secondary hash (optional)
	MD5    string `json:"md5"`    // For compatibility (not recommended for evidence)
}

// ImageMetadata contains image-specific metadata
type ImageMetadata struct {
	Format       string            `json:"format"`         // Image format (JPEG, PNG, etc.)
	Width        int               `json:"width"`          // Image width in pixels
	Height       int               `json:"height"`         // Image height in pixels
	ColorSpace   string            `json:"color_space"`    // Color space (RGB, CMYK, etc.)
	BitDepth     int               `json:"bit_depth"`      // Bits per channel
	EXIF         map[string]string `json:"exif,omitempty"` // EXIF metadata
	XMP          map[string]string `json:"xmp,omitempty"`  // XMP metadata
	IPTC         map[string]string `json:"iptc,omitempty"` // IPTC metadata
	GPS          *GPSCoordinates   `json:"gps,omitempty"`  // GPS coordinates if present
	CameraModel  string            `json:"camera_model"`   // Camera model
	LensMake     string            `json:"lens_make"`      // Lens manufacturer
	ExposureTime string            `json:"exposure_time"`  // Exposure time
	FNumber      string            `json:"f_number"`       // F-number
	ISO          int               `json:"iso"`            // ISO sensitivity
	FocalLength  string            `json:"focal_length"`   // Focal length
	DateTime     time.Time         `json:"date_time"`      // Photo date/time
	Software     string            `json:"software"`       // Software used
	Copyright    string            `json:"copyright"`      // Copyright info
	Artist       string            `json:"artist"`         // Artist/photographer
	Description  string            `json:"description"`    // Image description
	Keywords     []string          `json:"keywords"`       // Keywords/tags
}

// GPSCoordinates represents GPS location data
type GPSCoordinates struct {
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Altitude  float64   `json:"altitude"`
	Timestamp time.Time `json:"timestamp"`
}

// VideoMetadata contains video-specific metadata
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

// AcquisitionLog records all operations during evidence acquisition
type AcquisitionLog struct {
	StartTime      time.Time  `json:"start_time"`      // Acquisition start
	EndTime        time.Time  `json:"end_time"`        // Acquisition end
	Duration       string     `json:"duration"`        // Total duration
	OperationCount int        `json:"operation_count"` // Total operations
	ErrorCount     int        `json:"error_count"`     // Total errors
	WarningCount   int        `json:"warning_count"`   // Total warnings
	Entries        []LogEntry `json:"entries"`         // Detailed log entries
}

// LogEntry represents a single log entry
type LogEntry struct {
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"` // INFO, WARNING, ERROR
	Message   string    `json:"message"`
	Details   string    `json:"details"`
	Error     string    `json:"error"`
}

// SystemContext records system information during acquisition
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

// VerificationReport provides integrity verification results
type VerificationReport struct {
	Timestamp      time.Time `json:"timestamp"`
	VerifiedFiles  int       `json:"verified_files"`
	FailedFiles    int       `json:"failed_files"`
	SkippedFiles   int       `json:"skipped_files"`
	VerificationOK bool      `json:"verification_ok"`
	Failures       []string  `json:"failures"`
	Details        string    `json:"details"`
}
