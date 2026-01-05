# Evidex Internal Modules Documentation

## Overview

The `internal/` directory contains the core forensic acquisition logic, separated into specialized modules:

```
internal/
├── acquisition/    # Evidence acquisition engine
├── formatter/      # Output format generation
├── metadata/       # Metadata extraction
├── models/         # Data structures
└── utils/          # Utility functions
```

---

## Module: Models (`internal/models/`)

### Purpose
Defines all data structures used throughout Evidex for type safety and consistency.

### Key Types

#### `EvidencePackage`
Complete forensic evidence package containing:
- Manifest (chain of custody)
- Files (acquired evidence with metadata)
- Acquisition log (operations log)
- System context (system information)
- Timestamps and version

```go
type EvidencePackage struct {
    Manifest       *ChainOfCustodyManifest
    Files          []*FileEvidence
    AcquisitionLog *AcquisitionLog
    SystemContext  *SystemContext
    CreatedAt      time.Time
    Version        string
}
```

#### `ChainOfCustodyManifest`
Master record for legal admissibility:
- Evidence ID and creation date
- User and hostname information
- File statistics
- Hash algorithm configuration
- Acquisition method documentation
- Custodian history
- Integrity status

```go
type ChainOfCustodyManifest struct {
    ID                  string        // Unique evidence package ID
    CreatedAt           time.Time     // Package creation
    CreatedByUser       string        // Examiner name
    CreatedByHostname   string        // System hostname
    FileCount           int           // Total files
    TotalSize           int64         // Total bytes
    HashAlgorithm       string        // SHA-256, SHA-512
    AcquisitionMethod   string        // read-only
    Custodians          []Custodian   // Custody history
}
```

#### `FileEvidence`
Per-file metadata and acquisition record:
- Source and relative paths
- File size and timestamps
- Permissions and ownership
- Cryptographic hashes
- Image-specific metadata (EXIF, GPS)
- Video-specific metadata (codec, duration)
- Verification status

#### `ImageMetadata`
Image-specific data extraction:
- Format detection (JPEG, PNG, etc.)
- Dimensions and color space
- EXIF data (camera, lens, exposure)
- XMP data (copyright, keywords)
- GPS coordinates
- Timestamps

#### `VideoMetadata`
Video file information:
- Container format (MP4, MKV, etc.)
- Codec information (video, audio)
- Duration and frame rate
- Bitrate information
- Resolution
- Creation time and software

#### `AcquisitionLog`
Detailed operation log:
- Start and end times
- Operation count and error tracking
- Individual log entries
- Each entry: timestamp, level, message, error

#### `SystemContext`
System information snapshot:
- Hostname and username
- Operating system and version
- Architecture and timezone
- Local and UTC times
- Binary version and path

---

## Module: Utils (`internal/utils/`)

### Purpose
Utility functions for logging, hashing, file operations, and system information.

### Key Functions

#### Hashing Functions

```go
// Calculate SHA-256 hash of file
CalculateSHA256(filePath string) (string, error)

// Calculate SHA-512 hash of file
CalculateSHA512(filePath string) (string, error)

// Calculate MD5 hash (legacy compatibility only)
CalculateMD5(filePath string) (string, error)

// Verify file hash against known value
VerifyHash(filePath, expectedHash, algorithm string) (bool, error)
```

#### Logging Functions

```go
// Initialize logging system
InitLogging()

// Log at various levels
Log(level LogLevel, message, details, errMsg string)
LogInfo(message, details string)
LogWarning(message, details string)
LogError(message, details string, err error)

// Retrieve all log entries
GetLogEntries() []*models.LogEntry
```

#### File Operations

```go
// Check if file can be opened read-only
IsReadOnly(filePath string) (bool, error)

// Safely copy file without modifying source
SafeCopyFile(sourcePath, destinationPath string) error

// Ensure directory exists
EnsureDirectory(dirPath string) error
```

#### System Information

```go
// Get current system context
GetSystemContext() *models.SystemContext

// Get current executing user
GetCurrentUser() (string, error)

// Generate unique evidence ID
GenerateEvidenceID() string
```

#### File Type Detection

```go
// Get MIME type from file extension
GetFileTypeFromExtension(filePath string) string

// Check if file is image
IsImageFile(filePath string) bool

// Check if file is video
IsVideoFile(filePath string) bool
```

---

## Module: Metadata (`internal/metadata/`)

### Purpose
Extraction of forensically relevant metadata from files.

### Key Functions

#### File Metadata Extraction

```go
// Extract complete file metadata
ExtractFileMetadata(filePath string) (*models.FileEvidence, error)
```

This function:
1. Gets filesystem stats (size, permissions)
2. Extracts timestamps (modified, accessed, created)
3. Gets ownership information
4. Detects symlinks
5. Calls format-specific extractors
6. Returns complete FileEvidence structure

#### Image Metadata

```go
// Extract image metadata (EXIF, XMP, IPTC)
ExtractImageMetadata(filePath string) (*models.ImageMetadata, error)
```

Supported formats:
- JPEG: Full EXIF extraction with GPS
- PNG: Metadata chunks
- GIF: Animation and timing info
- Others: Basic format detection

#### Video Metadata

```go
// Extract video codec and duration info
ExtractVideoMetadata(filePath string) (*models.VideoMetadata, error)
```

Extracts:
- Container format identification
- Video codec and bitrate
- Audio codec and sample rate
- Frame rate and duration
- Resolution

#### Hash Calculation

```go
// Calculate all hashes for a file
CalculateFileHashes(filePath string, primaryAlg string) (*models.FileHashes, error)

// Verify file integrity
VerifyFileIntegrity(filePath, knownHash string) (bool, error)
```

#### Format Detection

Internal functions for file type detection:
- `isJPEG()`: JPEG magic bytes check
- `isPNG()`: PNG signature validation
- `isGIF()`: GIF header detection
- `isMP4()`: MP4 ftyp box detection
- `isMOV()`: QuickTime box detection
- `isMKV()`: EBML header detection
- `isAVI()`: RIFF/AVI signature
- `isWebM()`: WebM EBML detection

---

## Module: Acquisition (`internal/acquisition/`)

### Purpose
Core evidence acquisition engine managing the entire collection process.

### Type: Acquirer

Main acquisition coordinator:

```go
type Acquirer struct {
    outputDir        string                     // Output directory
    files            []*models.FileEvidence    // Collected files
    acquisitionLog   *models.AcquisitionLog    // Operation log
    primaryHashAlg   string                    // SHA-256/512
    secondaryHashAlg string                    // SHA-512
    validateAccess   bool                      // Verify read-only
}
```

### Key Methods

#### Acquisition

```go
// Acquire single file
AcquireFile(filePath string) error
    - Verify file exists and readable
    - Extract metadata
    - Calculate hashes
    - Log operation
    - Store FileEvidence

// Acquire directory (recursive/non-recursive)
AcquireDirectory(dirPath string, recursive bool) error
    - Enumerate directory
    - Call AcquireFile for each file
    - Optionally recurse subdirectories

// Acquire multiple files from list
AcquireMultiple(filePaths []string) error
    - Process each path
    - Track errors
    - Continue on failures
```

#### Package Building

```go
// Copy acquired files to output directory
CopyFilesToPackage() error
    - Create files/ subdirectory
    - Copy each file safely
    - Preserve original

// Build complete evidence package
GetEvidencePackage() *models.EvidencePackage
    - Create manifest
    - Get system context
    - Compile log entries
    - Return complete package
```

#### Information

```go
// Get count of acquired files
GetFileCount() int

// Get total size of acquired files
GetTotalSize() int64

// Set hash algorithm
SetHashAlgorithm(algorithm string)
```

### Acquisition Workflow

1. **File Validation**
   - Check file exists
   - Verify read-only access
   - Skip directories (files only)

2. **Metadata Extraction**
   - Filesystem metadata (timestamps, permissions)
   - Format-specific metadata (EXIF for images, codec for video)
   - Ownership information

3. **Hash Calculation**
   - Primary hash (SHA-256 or SHA-512)
   - Secondary hash (SHA-512)
   - Optional MD5 for compatibility

4. **File Storage**
   - Create FileEvidence structure
   - Add to acquisition list
   - Log operation

5. **Package Assembly**
   - Copy files to output directory
   - Create manifest with metadata
   - Compile acquisition log
   - Capture system context

---

## Module: Formatter (`internal/formatter/`)

### Purpose
Export evidence packages in various formats for analysis and verification.

### Type: PackageFormatter

```go
type PackageFormatter struct {
    outputDir string                    // Output directory
    package_  *models.EvidencePackage  // Evidence package
}
```

### Export Functions

#### JSON Export

```go
// Export JSON metadata
ExportJSON() error
```

Creates:
- `metadata/manifest.json` - Chain of custody
- `metadata/acquisition_log.json` - Operation log
- `metadata/system_context.json` - System info
- `metadata/file_catalog.json` - Complete file metadata

#### CSV Export

```go
// Export file listing as CSV
ExportCSV() error
```

Creates:
- `metadata/file_manifest.csv`
- Columns: filename, path, size, hash, verified status
- Compatible with analysis tools

#### Hash File

```go
// Create hash verification file
ExportHashes() error
```

Creates:
- `HASHES.txt`
- Per-file SHA-256 and SHA-512
- Simple text format for verification

#### Reports

```go
// Create integrity verification report
ExportIntegrityReport() error
```

Creates:
- `INTEGRITY_REPORT.txt`
- Summary statistics
- File counts and sizes
- Verification status

#### Documentation

```go
// Create package README
CreateReadme() error
```

Creates:
- `README.txt`
- Package contents description
- Verification instructions
- Chain of custody information

#### Package Verification

```go
// Calculate hash of entire manifest
CalculatePackageHash() (string, error)

// List package contents
CompressPackage(compressionFormat string) error
```

---

## Data Flow Through Modules

```
User Input (CLI)
    ↓
    └─→ [models] Define data structures
        ↓
        └─→ [acquisition] Create Acquirer
            ├─→ [utils] Validate file access
            ├─→ [metadata] Extract metadata
            │   ├─→ [utils] Calculate hashes
            │   └─→ [metadata] Extract format-specific data
            ├─→ [utils] Copy files safely
            └─→ Build EvidencePackage
                ↓
                └─→ [formatter] Export formats
                    ├─→ JSON files
                    ├─→ CSV manifest
                    ├─→ Hash files
                    ├─→ Reports
                    └─→ Documentation
```

---

## Error Handling

### Acquisition Errors
- Non-blocking: Errors in individual files don't stop acquisition
- Logged: All errors recorded in acquisition log
- Tracked: Error count maintained in manifest
- Reported: Summary shown to user

### Critical Errors
- Output directory creation failure
- Package building failure
- Export format failure

These stop execution with error reporting.

### Validation Errors
- Missing required flags
- Invalid file paths
- Invalid hash algorithms
- Permission issues

These prevent execution and show usage help.

---

## Testing

Each module should have corresponding test coverage:

```
internal/
├── acquisition/
│   ├── acquisition.go
│   └── acquisition_test.go
├── formatter/
│   ├── formatter.go
│   └── formatter_test.go
├── metadata/
│   ├── metadata.go
│   └── metadata_test.go
├── models/
│   └── models.go              # No tests (data structures)
└── utils/
    ├── utils.go
    └── utils_test.go
```

---

## Performance Considerations

### Memory
- Streaming hash calculation (constant memory)
- Log entries accumulated (linear with operations)
- File metadata only (not file contents)

### Disk I/O
- Sequential file reading (for hashing)
- File copying (bandwidth limited)
- Metadata writing (small files)

### Scaling
- Linear time with file count
- Linear time with total file size
- Handles thousands of files efficiently

---

## Security Notes

### Input Validation
- File paths validated before access
- Directory traversal prevention
- Symlink detection and handling

### Access Control
- Read-only file access
- No privilege escalation
- User context captured
- File permissions preserved

### Cryptographic Security
- NIST-approved algorithms
- Properly seeded random generation
- No password storage
- No security assumptions about file permissions

---

## Future Enhancement Points

1. **Parallel Processing**: Concurrent file hashing
2. **Database Support**: SQLite/MySQL acquisition
3. **Network Evidence**: Browser cache, temp files
4. **Digital Signatures**: HMAC/RSA on manifests
5. **Compression**: Tar.gz package creation
6. **Encryption**: Optional AES encryption
7. **Upload**: Secure server transmission
8. **GUI**: Web-based interface

---

## Conclusion

The modular design of Evidex allows independent testing, maintenance, and enhancement of each component while maintaining forensic integrity throughout the acquisition pipeline.

Each module has a single, well-defined responsibility, making the codebase auditable and trustworthy for forensic proceedings.
