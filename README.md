# Evidex

## Description

Evidex is a forensic evidence acquisition tool that collects files and comprehensive metadata for transmission to remote analysis servers. It operates in read-only mode, ensuring zero modifications to source evidence while maintaining complete chain of custody documentation.

## Purpose

Evidex acquires forensic evidence from files and directories, extracting detailed metadata from images (EXIF, GPS, XMP), videos (codecs, duration, resolution), and documents. All evidence is packaged with cryptographic hashes (MD5, SHA1, SHA256) and transmitted securely to remote analysis servers for processing. All operations are performed in-memory without local storage.

## Problem It Solves

Digital forensic investigations require collecting evidence without altering source files. Evidex provides a portable, auditable solution that acquires complete file metadata and content while guaranteeing forensic integrity through read-only operations, cryptographic verification, and comprehensive logging. It enables secure transmission of evidence packages to centralized analysis platforms, maintaining legally defensible chain of custody throughout the process.

- ✅ **Complete Evidence Package**: Transmits file content and comprehensive metadata to remote servers
- ✅ **Zero Modifications**: Source files are NEVER modified, deleted, or moved
- ✅ **Read-Only Access**: All operations use read-only file handles
- ✅ **Cryptographic Integrity**: MD5, SHA1, and SHA256 hashing for verification
- ✅ **Complete Chain of Custody**: Immutable custody documentation with command execution history
- ✅ **Memory-Only Processing**: No local storage, direct transmission to server
- ✅ **Remote Analysis Ready**: Transmit metadata packages to analysis servers via secure HTTP endpoints
- ✅ **Portable**: Cross-platform (Linux, macOS, Windows, FreeBSD, OpenBSD)
- ✅ **Auditable**: Open source, no external dependencies
- ✅ **Legally Defensible**: Compliant with ISO 27037 and NIST standards

## Key Features

**Purpose**: Evidex is exclusively designed for metadata extraction and acquisition from all file types for forensic analysis and transmission to remote analysis platforms.

| Feature Category | Description |
|-----------------|-------------|
| **Metadata Extraction** | Comprehensive metadata from all file types, filesystem metadata (timestamps, permissions, ownership), image metadata (EXIF, XMP, IPTC, GPS coordinates), video metadata (codec, duration, frame rate, bitrate, resolution), document metadata (author, creation date, modification history), and automatic format detection (MIME types, magic bytes) |
| **Evidence Acquisition** | Non-destructive evidence collection from single files or entire directories, recursive directory traversal, batch processing of multiple paths, read-only operations guarantee no source modification |
| **Remote Analysis Transmission** | Direct transmission of complete evidence packages (file content + metadata) to analysis servers, secure HTTP/HTTPS endpoints with Bearer token authentication, RFC 5424 structured logging, intelligent chunking for large files, and transfer verification |
| **Chain of Custody** | Unique evidence package identifier, complete custody history with timestamps, digital signature capability, examiner notes and case description, and transferable custody records |
| **Cryptographic Verification** | Primary hash: SHA-256, secondary hash: SHA-512, hash-based integrity verification, collision-resistant algorithms, and industry-standard implementation |
| **Output Formats** | **JSON** (complete metadata manifest), **CSV** (file listing for analysis tools), **TXT** (human-readable hashes and reports), and **RFC 5424 Logs** (structured syslog format for analysis servers) |

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

Evidex operates in memory-only mode, transmitting evidence directly to remote servers without local storage:

```bash
# Acquire and transmit evidence to analysis server
./build/evidex --server https://analysis.server.com/api/evidence --token YOUR_TOKEN --case-id CASE-2026-001 image.jpg

# Process multiple files with case ID
./build/evidex --server https://analysis.server.com/api/evidence --token YOUR_TOKEN --case-id CASE-2026-001 file1.jpg file2.mp4

# Process entire directory recursively
./build/evidex --server https://analysis.server.com/api/evidence --token YOUR_TOKEN --case-id CASE-2026-001 -r /path/to/folder
```

**Key Benefits**:
- No local storage required (memory-only processing)
- Direct transmission to analysis server
- Minimal forensic footprint on acquisition system
- Complete command execution logging
- Ideal for field operations
- Evidence package inspection
- Testing and development

**Note**: Evidex acquires complete files and comprehensive metadata (EXIF, file properties, codecs, etc.) for forensic analysis. Source files are **never modified** - only read in a non-destructive manner.

**Important**: The `-send` and `-o` flags are mutually exclusive. Use `-send` for production (remote transmission) or `-o` for debug (local storage only).

### Command-Line Options

```
-h, --help              Show help message
-v, --version           Show version information
-r, --recursive         Recursively process directories
--case-id ID            Case identifier for correlation (required)
--server URL            Remote server endpoint for transmission (required)
--token TOKEN           Authentication token for server communication (required)
```

**Required Options**:
- `--server URL` - Analysis server endpoint URL
- `--token TOKEN` - Authentication token
- `--case-id ID` - Case identifier for correlation
- At least one file or directory path

**Notes**:
- All hashes (MD5, SHA1, SHA256) are calculated automatically
- All file operations and commands are logged in the custody chain
- Evidence is transmitted directly to server without local storage

### Remote Transmission

Evidex transmits complete evidence packages directly to analysis servers:

```bash
# Direct transmission without local storage
./build/evidex --server https://analysis.server.com/api/evidence --token YOUR_TOKEN --case-id CASE-2026-001 image.jpg
```

**Transmitted Data**:
- Complete file content (base64-encoded)
- All metadata (filesystem, EXIF, video codecs, etc.)
- Cryptographic hashes (MD5, SHA1, SHA256, SHA512)
- Chain of custody with command history
- RFC 5424 structured logs

### Logging

RFC 5424 structured logging with metadata:
```
<134>1 2026-01-06T10:30:45Z hostname evidex 1234 - [meta@1 file_path="/path/to/file"] Successfully acquired file
```

## Building

### Prerequisites
- Go 1.25 or higher
- Make utility

### Build Commands

```bash
# Current platform
make build

# Specific platforms
make build-linux          # Linux (amd64, arm64)
make build-darwin         # macOS (amd64, arm64)
make build-windows        # Windows (amd64, arm64)
make build-freebsd-amd64  # FreeBSD
make build-openbsd-amd64  # OpenBSD

# All platforms
make build-all

# Other targets
make test      # Run tests
make clean     # Clean build artifacts
make fmt       # Format code
make vendor    # Update vendor directory
```

## Development

### Dependency Management

Evidex uses **Go modules with vendoring** for reproducible builds and offline compilation:

```bash
# Update vendor directory (after adding/updating dependencies)
go mod vendor

# Verify vendor integrity
go mod verify

# All builds automatically use vendored dependencies
go build -mod=vendor ./cmd/evidex
go test -mod=vendor ./tests/...
```

**Why Vendoring?**
- **Reproducibility**: Exact same dependencies across all builds
- **Offline builds**: No internet required for compilation
- **Security**: Audit and review all dependencies in version control
- **Forensic integrity**: Self-contained tool with known dependencies

### Running Tests

```bash
# Run all tests
make test

# Run tests with race detection
make test-race

# Run tests with coverage
make test-coverage

# Run with vendored dependencies
go test -mod=vendor -v ./tests/...
```

### Code Quality

```bash
# Format code
make fmt

# Run linter (requires golangci-lint)
make lint

# Run all quality checks
go fmt ./...
go vet -mod=vendor ./...
golangci-lint run ./...
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

## Configuration

No configuration file is needed. All options are command-line flags.

Environment variables:
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

## Digital Evidence Custody Chain

Evidex implements a comprehensive digital evidence custody chain:

**Features:**
- MD5, SHA1, SHA256 hash algorithms for integrity verification
- Complete command logging during collection
- Custody transfer tracking with timestamps
- Agent identification (hostname, username, version)
- Automatic timeline generation from file timestamps
- Integration with Processor for analysis

**Usage:**
The custody chain is automatically created and maintained during evidence collection. All operations are logged and included in the transmitted evidence package.

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

For issues or questions, open an issue on GitHub.

## Disclaimer

**This tool is provided "as-is" for legitimate forensic purposes only.**

Users are responsible for:
- Compliance with applicable laws
- Proper chain of custody procedures
- Legal consultation before use in proceedings
- Preservation of original evidence
- Proper storage and access control of evidence packages

**Do not use for unauthorized access to files or systems.**
