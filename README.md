# Evidex - Forensic Evidence Acquisition Binary

A portable, auditable forensic tool for multimedia file collection with complete chain of custody documentation. Designed for judicial proceedings with zero modifications to source evidence.

## Overview

**Evidex** is a forensic evidence acquisition binary that collects image and video files while guaranteeing:

- ✅ **Zero Modifications**: Source files are NEVER modified, deleted, or moved
- ✅ **Read-Only Access**: All operations use read-only file handles
- ✅ **Cryptographic Integrity**: SHA-256 and SHA-512 hashing for verification
- ✅ **Complete Chain of Custody**: Immutable custody documentation and audit trail
- ✅ **Portable**: Cross-platform (Linux, macOS, Windows, FreeBSD, OpenBSD)
- ✅ **Auditable**: Open source, no external dependencies
- ✅ **Legally Defensible**: Compliant with ISO 27037 and NIST standards

## Key Features

### Evidence Acquisition
- Acquire single files or entire directories
- Recursive directory traversal
- Batch processing of multiple paths
- Non-destructive collection (read-only)

### Metadata Extraction
- Filesystem metadata (timestamps, permissions, ownership)
- Image metadata (EXIF, XMP, IPTC, GPS coordinates)
- Video metadata (codec, duration, frame rate, bitrate, resolution)
- Automatic format detection (MIME types, magic bytes)

### Chain of Custody
- Unique evidence package identifier
- Complete custody history with timestamps
- Digital signature capability
- Examiner notes and case description
- Transferable custody records

### Cryptographic Verification
- Primary hash: SHA-256
- Secondary hash: SHA-512
- Hash-based integrity verification
- Collision-resistant algorithms
- Industry-standard implementation

### Output Formats
- **JSON**: Complete metadata and manifest
- **CSV**: File listing for analysis tools
- **TXT**: Human-readable hashes and reports
- **Original Files**: Unmodified copies in package

### Comprehensive Logging
- Detailed acquisition operation log
- System context capture
- Error and warning tracking
- UTC and local timestamps
- Timezone-aware recording

## Quick Start

### Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/evidex.git
cd evidex

# Build for your platform
make build

# Or build for all platforms
make build-all
```

### Basic Usage

```bash
# Acquire a single image file
./build/evidex -o ./evidence image.jpg

# Acquire multiple files
./build/evidex -o ./evidence photo1.jpg video1.mp4 photo2.jpg

# Recursively acquire all files in a directory
./build/evidex -o ./evidence -r /path/to/folder

# Acquire with case description
./build/evidex -o ./evidence -desc "Case: ABC-2025-001" image.jpg

# Use SHA-512 hashing
./build/evidex -o ./evidence -hash SHA-512 file.jpg
```

### Command-Line Options

```
-h, -help              Show help message
-v, -version           Show version information
-o, -output DIR        Output directory for evidence package (required)
-r, -recursive         Recursively process directories
-hash ALGORITHM        Hash algorithm: SHA256 (default), SHA512
-desc DESCRIPTION      Evidence description for manifest
-notes NOTES          Examiner notes for chain of custody
-json                  Export JSON metadata (default: true)
-csv                   Export CSV file manifest
-hashes                Create HASHES.txt file for verification
-report                Create integrity report
-verify                Verify hashes after acquisition (default: true)
```

## Output Package Structure

```
evidence_package/
├── files/
│   ├── image1.jpg
│   ├── image2.png
│   ├── video1.mp4
│   └── ...
│
├── metadata/
│   ├── manifest.json           # Chain of custody
│   ├── acquisition_log.json    # Detailed operations
│   ├── system_context.json     # System information
│   ├── file_catalog.json       # File metadata
│   └── file_manifest.csv       # CSV format
│
├── HASHES.txt                  # Hash verification
├── INTEGRITY_REPORT.txt        # Integrity status
├── README.txt                  # Documentation
└── PACKAGE_CONTENTS.txt        # Package structure
```

## Files Explained

### manifest.json
Master chain of custody record containing:
- Evidence ID (unique identifier)
- Creation date and user
- File count and total size
- Hash algorithms used
- Acquisition method
- Custodian history
- Evidence description and examiner notes

### file_catalog.json
Complete metadata for each acquired file:
- File path and name
- Size and timestamps
- Permissions and ownership
- Cryptographic hashes
- Image/video specific metadata
- Verification status

### acquisition_log.json
Detailed log of all acquisition operations:
- Start and end times
- Duration
- Operation count
- Error and warning counts
- Individual log entries with timestamps

### HASHES.txt
Simple hash verification file for independent verification:
- Evidence ID
- File count and total size
- Per-file hashes (SHA-256, SHA-512)
- Suitable for hash database lookup

## Building

### Prerequisites
- Go 1.25 or higher
- Make utility
- Git (for version control)

### Build Targets

```bash
# Current platform
make build

# Linux
make build-linux          # Build for Linux (amd64, arm64)
make build-linux-amd64    # Linux x86_64
make build-linux-arm64    # Linux ARM 64-bit

# macOS
make build-darwin         # Build for macOS (amd64, arm64)
make build-darwin-amd64   # macOS Intel
make build-darwin-arm64   # macOS Apple Silicon

# Windows
make build-windows        # Build for Windows (amd64, arm64)
make build-windows-amd64  # Windows x86_64
make build-windows-arm64  # Windows ARM 64-bit

# Other Unix
make build-freebsd-amd64  # FreeBSD
make build-openbsd-amd64  # OpenBSD

# All platforms
make build-all            # Build all supported platforms

# Release archives
make release              # Create compressed packages
```

### Other Targets

```bash
make deps      # Download dependencies
make test      # Run tests
make clean     # Clean build artifacts
make fmt       # Format code
make lint      # Lint code
make vendor    # Update vendor directory
make help      # Show all targets
```

## Forensic Standards Compliance

Evidex implements forensic best practices aligned with:

- **ISO 27037**: Digital evidence identification and preservation
- **NIST SP 800-86**: Digital forensics guidelines
- **Chain of Custody Requirements**: Legal evidence handling
- **Daubert Standard**: Expert testimony standards
- **Non-Repudiation**: Evidence integrity through hashing

## Supported File Types

### Images
- JPEG (.jpg, .jpeg)
- PNG (.png)
- GIF (.gif)
- BMP (.bmp)
- TIFF (.tiff, .tif)
- WebP (.webp)
- RAW (.raw)
- HEIC/HEIF (.heic, .heif)

### Videos
- MP4 (.mp4)
- MOV (.mov) - QuickTime
- MKV (.mkv) - Matroska
- AVI (.avi)
- FLV (.flv)
- WMV (.wmv)
- WebM (.webm)
- 3GP (.3gp)
- OGG (.ogv)
- MPEG-TS (.ts)

## Use Cases

1. **Criminal Investigation**: Collect evidence from digital devices
2. **Civil Litigation**: Preserve digital evidence for court proceedings
3. **Incident Response**: Acquire evidence during security investigations
4. **Digital Forensics**: Professional evidence collection
5. **E-Discovery**: Systematic file collection for legal review
6. **Intellectual Property**: Image/video ownership documentation

## Legal Considerations

### Admissibility in Court

Evidex evidence packages are designed for legal admissibility by:

1. **Clear Chain of Custody**: Complete custody history and transitions
2. **Integrity Verification**: Cryptographic hashes prove no tampering
3. **Method Transparency**: Non-destructive, well-documented procedures
4. **Expert Testimony**: Acquisition method is technically defensible
5. **Reproducibility**: Same input produces same results

### Important Disclaimers

- **Legal Advice**: Consult with legal counsel before using in proceedings
- **Jurisdiction**: Requirements vary by jurisdiction
- **Authentication**: Ensure proper authentication of chain of custody
- **Documentation**: Maintain detailed records of all acquisition activities
- **Preservation**: Preserve original evidence alongside package

## Security Considerations

### File Access
- Read-only access to all source files
- No modifications to source data
- Separate isolated output directory
- File permissions preserved in metadata

### Cryptographic Integrity
- SHA-256: Industry-standard for forensics
- SHA-512: Stronger security margin for long-term storage
- No password protection (files are plaintext JSON)
- Consider encrypting package for transport/storage

### Source Code Auditing
- Open source for review
- No obfuscation or hidden behavior
- Simple, straightforward Go code
- Reproducible builds possible

### Chain of Custody
- Immutable manifest after creation
- Signature capability for authentication
- No post-acquisition modifications
- Complete audit trail

## Performance

- **Speed**: Depends on file size and storage speed (typically limited by disk I/O)
- **Memory**: Minimal memory footprint (streaming hash calculation)
- **Disk Space**: Requires space for original files + metadata
- **Scaling**: Handles thousands of files efficiently

## Troubleshooting

### "Permission denied"
- Ensure read access to source files
- Run with appropriate privileges
- Check file ownership and permissions

### "Hash mismatch"
- Verify source files haven't changed since acquisition
- Re-run acquisition to regenerate hashes
- Compare against original HASHES.txt

### "Output directory not found"
- Ensure parent directory exists
- Create output directory manually
- Check write permissions to output location

### "Cannot open file"
- Verify file path is correct
- Check for symbolic links or shortcuts
- Ensure sufficient disk space

## Configuration

No configuration file is needed. All options are command-line flags.

Environment variables used:
- None (self-contained binary)

## Requirements

- Go 1.25+ (for building)
- 100 MB+ disk space (for output)
- Read access to evidence files
- Write access to output directory

## Dependencies

Minimal external dependencies:
- `rwcarlsen/goexif` (EXIF parsing)
- `golang.org/x/crypto` (for optional future features)

No network dependencies - fully offline capable.

## License

Licensed under the Apache License, Version 2.0.

```
Copyright 2026 Evidex Contributors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```

## Documentation

- [ARCHITECTURE.md](ARCHITECTURE.md) - Comprehensive technical design
- [cmd/evidex/README.md](cmd/evidex/README.md) - CLI documentation
- [internal/README.md](internal/README.md) - Module documentation

## Contributing

Contributions welcome! Please ensure:

1. Code follows Go conventions
2. Tests pass (`make test`)
3. Code is formatted (`make fmt`)
4. Linting passes (`make lint`)
5. Documentation is updated
6. Forensic integrity is maintained

## Author

Designed and implemented as a professional forensic evidence acquisition tool.

## Support

For issues, questions, or forensic consultation:
- Review ARCHITECTURE.md for design details
- Check troubleshooting section
- Examine acquisition_log.json in generated packages
- Review source code for implementation details

## Disclaimer

**This tool is provided "as-is" for legitimate forensic purposes only.**

Users are responsible for:
- Compliance with applicable laws
- Proper chain of custody procedures
- Legal consultation before use in proceedings
- Preservation of original evidence
- Proper storage and access control of evidence packages

**Do not use for unauthorized access to files or systems.**

---

**Version**: 1.0.0  
**Status**: Production Ready  
**Last Updated**: 2026-01-05
