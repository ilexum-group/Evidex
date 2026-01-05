# Evidex - Forensic Evidence Acquisition Architecture

## Executive Summary

**Evidex** is a portable, auditable forensic evidence acquisition binary designed for collecting multimedia files (images and videos) with complete chain of custody documentation. The tool operates on read-only principles, making zero modifications to source files while maintaining forensic integrity through cryptographic hashing and comprehensive logging.

---

## Table of Contents

1. [Design Principles](#design-principles)
2. [System Architecture](#system-architecture)
3. [Data Flow](#data-flow)
4. [Package Structure](#package-structure)
5. [Chain of Custody Model](#chain-of-custody-model)
6. [Evidence Package Format](#evidence-package-format)
7. [Security Considerations](#security-considerations)
8. [Legal and Forensic Standards](#legal-and-forensic-standards)

---

## Design Principles

### Fundamental Axioms

1. **Non-Modification**: Source files are NEVER modified, deleted, or moved
2. **Read-Only Access**: All file operations use read-only file handles
3. **Complete Logging**: Every operation is logged with timestamps and context
4. **Reproducibility**: Identical inputs produce identical outputs (deterministic)
5. **Auditability**: All acquisition steps can be independently verified
6. **Portability**: Binary runs on Linux, macOS, Windows, FreeBSD, OpenBSD
7. **No External Dependencies**: Zero reliance on remote services or network calls
8. **Transparency**: Source code is open and auditable

### Forensic Integrity Commitments

- **Hash Verification**: Primary (SHA-256) + Secondary (SHA-512) hashes
- **Metadata Preservation**: Complete filesystem and file metadata captured
- **Temporal Accuracy**: UTC and local timestamps with timezone context
- **System Context**: Operating system, user, and hardware information recorded
- **Evidence ID**: Unique, traceable identifier for entire package

---

## System Architecture

### High-Level Components

```
┌─────────────────────────────────────────────────────────────┐
│                          EVIDEX CLI                          │
│  (cmd/evidex/main.go - Command-line interface & flags)      │
└────────────────┬────────────────────────────────────────────┘
                 │
┌────────────────▼────────────────────────────────────────────┐
│                   ACQUISITION LAYER                         │
│  (internal/acquisition/acquisition.go)                       │
│  - File enumeration and filtering                           │
│  - Read-only access verification                            │
│  - Metadata extraction                                      │
│  - Hash calculation & verification                          │
│  - Package assembly                                         │
└────────────────┬────────────────────────────────────────────┘
                 │
    ┌────────────┼────────────┬──────────────┐
    │            │            │              │
┌───▼───┐  ┌────▼─────┐ ┌───▼──────┐ ┌────▼──────┐
│Metadata Hashing   │ Format    │ Package
│Extractor Utility  │ Exporter  │ Verifier
└─────────┘ └────────┘ └──────────┘ └───────────┘
    │            │            │              │
    └────────────┼────────────┴──────────────┘
                 │
    ┌────────────▼────────────────────────────┐
    │   DATA MODELS & STRUCTURES               │
    │   (internal/models/models.go)            │
    │   - EvidencePackage                     │
    │   - ChainOfCustodyManifest              │
    │   - FileEvidence                        │
    │   - AcquisitionLog                      │
    │   - SystemContext                       │
    └─────────────────────────────────────────┘
                 │
    ┌────────────▼────────────────────────────┐
    │   OUTPUT FORMATS                        │
    │   - JSON (metadata)                     │
    │   - CSV (file listing)                  │
    │   - Text reports (hashes, integrity)    │
    │   - Original files (unmodified)         │
    └─────────────────────────────────────────┘
```

### Module Breakdown

#### 1. **CLI Interface** (`cmd/evidex/main.go`)
- Flag parsing and validation
- Input path processing (files/directories)
- Recursive directory traversal
- User feedback and progress reporting
- Package summary and verification info

#### 2. **Acquisition Engine** (`internal/acquisition/acquisition.go`)
- `Acquirer` type: Manages entire acquisition process
- `AcquireFile()`: Single file acquisition
- `AcquireDirectory()`: Directory traversal (recursive/non-recursive)
- `AcquireMultiple()`: Batch file processing
- `CopyFilesToPackage()`: Safe file copying to output directory
- `GetEvidencePackage()`: Build final evidence package

#### 3. **Metadata Extractor** (`internal/metadata/metadata.go`)
- Filesystem metadata (permissions, timestamps, ownership)
- Image metadata (EXIF, XMP, IPTC, GPS)
- Video metadata (codec, duration, frame rate, bitrate)
- Format detection (MIME types, magic bytes)
- Hash calculation (SHA-256, SHA-512, MD5)

#### 4. **Utility Functions** (`internal/utils/utils.go`)
- Hash calculation and verification
- System context gathering
- File access validation
- Logging and reporting
- Directory and file operations

#### 5. **Package Formatter** (`internal/formatter/formatter.go`)
- JSON metadata export
- CSV manifest generation
- Hash file creation
- Integrity report generation
- Package documentation

#### 6. **Data Models** (`internal/models/models.go`)
- `EvidencePackage`: Complete evidence container
- `ChainOfCustodyManifest`: Custody history and metadata
- `FileEvidence`: Per-file evidence record with all metadata
- `AcquisitionLog`: Detailed operation log
- `SystemContext`: System information snapshot

---

## Data Flow

### Acquisition Flow Diagram

```
User Input (Files/Directories)
    ↓
[CLI Parameter Validation]
    ↓
[Output Directory Setup]
    ↓
For Each Input Path:
    ├─ If File → AcquireFile()
    │   ├─ Verify file exists & accessible (read-only)
    │   ├─ Extract filesystem metadata
    │   ├─ Extract file-specific metadata
    │   │   ├─ Image: EXIF/XMP/IPTC parsing
    │   │   └─ Video: Codec/duration/frame rate detection
    │   ├─ Calculate hashes
    │   │   ├─ SHA-256 (primary)
    │   │   ├─ SHA-512 (secondary)
    │   │   └─ MD5 (legacy compatibility)
    │   └─ Store FileEvidence in acquisition package
    │
    └─ If Directory → AcquireDirectory()
        ├─ Enumerate directory contents
        ├─ For each file: AcquireFile()
        └─ Optionally recurse into subdirectories
    ↓
[Copy Files to Package Output]
    ├─ Create "files/" subdirectory
    ├─ Copy each file (preserve original, make copy)
    └─ Store relative path in metadata
    ↓
[Build Evidence Package]
    ├─ Create ChainOfCustodyManifest
    ├─ Record system context
    ├─ Compile acquisition log
    └─ Package all metadata
    ↓
[Export Formats]
    ├─ JSON metadata (manifest, files catalog, log)
    ├─ CSV file listing
    ├─ HASHES.txt verification file
    ├─ INTEGRITY_REPORT.txt
    └─ README.txt documentation
    ↓
[Final Package Structure]
    ├─ files/                    (Original files - unmodified)
    ├─ metadata/
    │   ├─ manifest.json        (Chain of custody)
    │   ├─ acquisition_log.json  (Detailed operations)
    │   ├─ system_context.json   (System information)
    │   ├─ file_catalog.json     (All file metadata)
    │   └─ file_manifest.csv     (CSV listing)
    ├─ HASHES.txt               (Hash verification)
    ├─ INTEGRITY_REPORT.txt      (Integrity status)
    ├─ README.txt                (Documentation)
    └─ PACKAGE_CONTENTS.txt      (Package structure)
```

---

## Package Structure

### Output Directory Layout

```
evidence_package/
├── files/
│   ├── image1.jpg
│   ├── image2.png
│   ├── video1.mp4
│   └── ...
│
├── metadata/
│   ├── manifest.json           # Chain of custody master record
│   ├── acquisition_log.json    # Detailed operation log
│   ├── system_context.json     # System information snapshot
│   ├── file_catalog.json       # Complete file metadata
│   └── file_manifest.csv       # CSV format for analysis tools
│
├── HASHES.txt                  # Hash verification file
├── INTEGRITY_REPORT.txt        # Integrity verification report
├── README.txt                  # Package documentation
└── PACKAGE_CONTENTS.txt        # Directory structure listing
```

### Key Files Explained

#### **manifest.json**
Chain of custody master record containing:
- Evidence ID (unique identifier)
- Creation timestamp and user
- File statistics (count, total size)
- Hash algorithm information
- Acquisition method
- Custodian history (who had custody when)
- Evidence description and notes
- Integrity status

Example:
```json
{
  "id": "EVDX-HOSTNAME-20260105181234",
  "created_at": "2026-01-05T18:12:34Z",
  "created_by_user": "investigator",
  "created_by_hostname": "forensic-workstation",
  "file_count": 15,
  "total_size": 2147483648,
  "hash_algorithm": "SHA-256",
  "secondary_algorithm": "SHA-512",
  "acquisition_method": "read-only",
  "integrity": "verified",
  "custodians": [
    {
      "name": "investigator",
      "role": "Examiner",
      "action": "created",
      "timestamp": "2026-01-05T18:12:34Z",
      "location": "Lab-1",
      "notes": "Initial evidence acquisition"
    }
  ]
}
```

#### **file_catalog.json**
Complete metadata for each file:
```json
[
  {
    "source_path": "/path/to/image.jpg",
    "filename": "image.jpg",
    "file_size": 5242880,
    "modified_time": "2025-12-25T10:30:00Z",
    "created_time": "2025-12-25T10:30:00Z",
    "hashes": {
      "sha256": "a1b2c3d4e5f6...",
      "sha512": "f6e5d4c3b2a1..."
    },
    "image_metadata": {
      "format": "JPEG",
      "width": 4000,
      "height": 3000,
      "camera_model": "Canon EOS 5D Mark IV",
      "date_time": "2025-12-25T10:30:00Z",
      "gps": {
        "latitude": 40.7128,
        "longitude": -74.0060,
        "altitude": 10.5
      }
    },
    "verified": true
  }
]
```

#### **HASHES.txt**
Simple hash verification file for auditing:
```
EVIDEX FORENSIC HASHES
Generated: 2026-01-05T18:12:34Z
Evidence ID: EVDX-HOSTNAME-20260105181234
================================================================================

File Count: 15
Total Size: 2147483648 bytes

FILES AND HASHES:
--------------------------------------------------------------------------------

File: /path/to/image.jpg
Size: 5242880 bytes
SHA-256: a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6
SHA-512: f6e5d4c3b2a1z9y8x7w6v5u4t3s2r1q0...

[... more files ...]
```

---

## Chain of Custody Model

### Conceptual Design

```
┌─────────────────────────────────────────────────────────────┐
│                CHAIN OF CUSTODY (COC)                       │
│                                                              │
│  Initial Evidence Collection:                              │
│  ┌────────────────────────────────────────────────────┐   │
│  │ Examiner: investigator                             │   │
│  │ Time: 2026-01-05T18:12:34Z                        │   │
│  │ Location: Lab-1                                   │   │
│  │ Action: created                                    │   │
│  │ Notes: Initial evidence acquisition               │   │
│  └────────────────────────────────────────────────────┘   │
│                      ↓                                      │
│  Evidence Transfer:                                        │
│  ┌────────────────────────────────────────────────────┐   │
│  │ From: investigator                                │   │
│  │ To: evidence_custodian                            │   │
│  │ Time: 2026-01-06T10:00:00Z                        │   │
│  │ Location: Secure Storage                          │   │
│  │ Action: transferred                                │   │
│  │ Signature: [digital signature]                    │   │
│  └────────────────────────────────────────────────────┘   │
│                      ↓                                      │
│  Further Transfers/Custody Changes...                     │
│                                                              │
│  Every transaction is immutable and timestamped            │
└─────────────────────────────────────────────────────────────┘
```

### Custody Record Structure

```go
type Custodian struct {
    Name       string    // Full name of custodian
    Role       string    // Title/role (Examiner, Analyst, etc.)
    Signature  string    // Digital signature (HMAC/RSA)
    Action     string    // created, transferred, archived, verified
    Timestamp  time.Time // When action occurred
    Location   string    // Physical or logical location
    Notes      string    // Free-form notes about custody action
}
```

### Legal Admissibility

The chain of custody manifest is designed to satisfy legal requirements:

1. **Identification**: Clear identification of evidence (Evidence ID)
2. **Handling Documentation**: Complete record of who had custody when
3. **Authentication**: Each custodian can sign/authenticate
4. **Integrity Verification**: Hash verification proves no tampering
5. **Timeline**: Complete temporal record from creation to present
6. **Context**: System information and acquisition method documented

---

## Evidence Package Format

### JSON Metadata Structure

The complete evidence package in JSON includes:

```json
{
  "manifest": { /* ChainOfCustodyManifest */ },
  "files": [ /* Array of FileEvidence */ ],
  "acquisition_log": { /* AcquisitionLog */ },
  "system_context": { /* SystemContext */ },
  "created_at": "2026-01-05T18:12:34Z",
  "version": "1.0.0"
}
```

### File Evidence Entry

```json
{
  "source_path": "/original/path/to/file.jpg",
  "relative_path": "file.jpg",
  "filename": "file.jpg",
  "file_size": 5242880,
  "file_mode": 33188,
  "accessed_time": "2025-12-25T10:30:00Z",
  "modified_time": "2025-12-25T10:30:00Z",
  "created_time": "2025-12-25T10:30:00Z",
  "change_time": "2025-12-25T10:30:00Z",
  "owner": "user",
  "group": "group",
  "file_type": "image/jpeg",
  "is_symlink": false,
  "hashes": {
    "sha256": "a1b2c3d4e5f6...",
    "sha512": "f6e5d4c3b2a1..."
  },
  "image_metadata": {
    "format": "JPEG",
    "width": 4000,
    "height": 3000,
    "color_space": "RGB",
    "camera_model": "Canon EOS 5D Mark IV",
    "lens_make": "Canon",
    "exposure_time": "1/125",
    "f_number": "8.0",
    "iso": 100,
    "focal_length": "50mm",
    "date_time": "2025-12-25T10:30:00Z",
    "copyright": "© 2025 Photographer",
    "gps": {
      "latitude": 40.7128,
      "longitude": -74.0060,
      "altitude": 10.5
    }
  },
  "acquisition_time": "2026-01-05T18:12:34Z",
  "verified": true,
  "verification_err": ""
}
```

---

## Security Considerations

### Cryptographic Hash Security

**SHA-256 (Primary)**
- Industry standard for forensics
- Collision-resistant for practical purposes
- Accepted in legal proceedings
- Recommended by NIST

**SHA-512 (Secondary)**
- Stronger security margin
- Higher output size (512 bits vs 256)
- Future-proofing against advances in cryptanalysis
- Recommended for long-term evidence preservation

**MD5 (Legacy)**
- Included for compatibility with legacy systems
- NOT recommended for new evidence
- Cryptographically broken (collision attacks possible)
- Kept for reference only

### File Access Security

1. **Read-Only Mode**: All file operations use `os.O_RDONLY` flag
2. **No Write Operations**: Source files are NEVER modified
3. **No Delete Operations**: Source files are NEVER removed
4. **Separate Copy**: Files are copied to isolated output directory
5. **Access Logging**: Every file access is logged with timestamp

### System Integrity

1. **Source Code Auditability**: Go source code is human-readable and reviewable
2. **No Obfuscation**: Binary compiled from unmodified source
3. **Deterministic Build**: Same source produces same binary (reproducible builds possible)
4. **Minimal Dependencies**: Only essential external packages
5. **Platform Transparency**: Cross-platform code is transparent and auditable

### Chain of Custody Security

1. **Immutable Manifest**: Manifest cannot be modified after creation
2. **Signature Scheme**: HMAC signatures on custodian records
3. **Timestamp Validation**: All times are UTC with timezone context
4. **Custody Verification**: Each transition can be independently verified

---

## Legal and Forensic Standards

### Applicable Standards

1. **NIST SP 800-86**: Guide to Integrating Forensic Techniques into Incident Response
2. **ISO 27037**: Guidelines for identification, collection, acquisition and preservation of digital evidence
3. **NCFS (National Commission for Forensic Science)**: Standards for digital evidence
4. **Daubert Standard**: Expert testimony standards (US Courts)
5. **Chain of Custody Requirements**: Court-accepted evidence handling procedures

### Evidex Compliance

✅ **Read-Only Acquisition**
- Evidence is never modified during collection
- Non-destructive to source media
- Complies with ISO 27037

✅ **Complete Metadata**
- All relevant file metadata captured
- Timestamps preserved with timezone info
- System context recorded

✅ **Cryptographic Integrity**
- Industry-standard hash algorithms
- Multiple hash values for verification
- Documented algorithm choices

✅ **Chain of Custody**
- Complete custody history maintained
- Timestamped with audit trail
- Supports digital signatures

✅ **Reproducibility**
- Deterministic results
- Documented processes
- Auditable source code

✅ **Expert Testimony Ready**
- All acquisition methods documented
- Technical approach is defensible
- Results can be independently verified

### Evidentiary Challenges & Responses

**Challenge**: "How do we know the files weren't modified?"
- **Response**: Cryptographic hashes provide mathematical proof of integrity
- **Evidence**: Original SHA-256/SHA-512 hashes can be independently verified
- **Documentation**: Acquisition log shows no modifications were made

**Challenge**: "What about timestamps?"
- **Response**: File modification times are captured from source filesystem
- **Evidence**: Timestamps are recorded with timezone and system context
- **Documentation**: UTC and local time both recorded for clarity

**Challenge**: "Is the tool trustworthy?"
- **Response**: Source code is open and auditable
- **Evidence**: Go source code available for review before compilation
- **Documentation**: Build process is documented and reproducible

---

## Workflow Summary

1. **User Initiates**: Specifies input files/directories and output location
2. **System Context**: Binary captures system information (OS, user, hostname)
3. **File Enumeration**: Identifies all files to be acquired
4. **Access Verification**: Confirms read-only access to each file
5. **Metadata Extraction**: Collects filesystem and file-specific metadata
6. **Hashing**: Calculates SHA-256/SHA-512 for each file
7. **File Copy**: Copies files to evidence package (original remains untouched)
8. **Package Assembly**: Creates manifest with chain of custody
9. **Export Formats**: Generates JSON, CSV, and text reports
10. **Verification**: Ensures all operations completed successfully
11. **Documentation**: Creates README and integrity report

---

## Deployment Considerations

### Portable Deployment

- Single binary (no dependencies)
- Works on air-gapped systems
- No installation required
- Can be executed from USB or network share

### Forensic Workstation Setup

1. Execute binary from read-only media
2. Direct output to separate evidence storage
3. Maintain custody documentation
4. Archive evidence package securely

### Integration with Forensic Suites

- JSON export compatible with DFIR tooling
- CSV format importable to analysis platforms
- Hash verification compatible with hash databases
- Chain of custody compatible with case management systems

---

## Future Enhancements

1. **Digital Signatures**: RSA or ECDSA signatures on manifests
2. **Video Analysis**: Detailed codec analysis with ffprobe
3. **Audio Metadata**: Extraction of audio-specific metadata
4. **Database Evidence**: SQLite/MySQL database acquisition
5. **Network Artifacts**: Browser history, cache, temporary files
6. **Encryption Handling**: Encrypted file handling and reporting
7. **Cloud Integration**: Optional upload to secure evidence storage
8. **Web Dashboard**: GUI for evidence acquisition
9. **Batch Processing**: Automation for multiple cases
10. **Blockchain Timestamp**: Immutable timestamp via blockchain

---

## Conclusion

Evidex implements forensic best practices through a carefully designed acquisition architecture. By maintaining read-only access, comprehensive logging, cryptographic verification, and complete chain of custody documentation, Evidex creates evidence packages suitable for judicial proceedings while remaining auditable and portable across platforms.

The tool prioritizes transparency, reproducibility, and legal defensibility—core requirements for any forensic acquisition tool used in professional investigations.
