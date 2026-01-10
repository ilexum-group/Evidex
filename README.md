# Evidex - Forensic Metadata Acquisition & Analysis Tool

A portable, auditable forensic tool for acquiring comprehensive metadata from all file types with complete chain of custody documentation and remote transmission for analysis. Designed for judicial proceedings with zero modifications to source evidence.

## Overview

**Evidex** is a forensic evidence acquisition tool that collects files and their metadata for transmission to remote analysis servers while guaranteeing:

- ✅ **Complete Evidence Package**: Transmits file content and comprehensive metadata to remote servers
- ✅ **Zero Modifications**: Source files are NEVER modified, deleted, or moved
- ✅ **Read-Only Access**: All operations use read-only file handles
- ✅ **Cryptographic Integrity**: SHA-256 and SHA-512 hashing for verification
- ✅ **Complete Chain of Custody**: Immutable custody documentation and audit trail
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

#### Production Mode (Recommended)

In production mode, evidence is transmitted directly to remote servers **without writing to local disk**:

```bash
# Acquire and transmit evidence to analysis server
./build/evidex -send -server https://analysis.server.com/api/evidence -token YOUR_TOKEN image.jpg

# Process multiple files with case ID
./build/evidex -send -server https://analysis.server.com/api/evidence -token YOUR_TOKEN -case-id CASE-2026-001 file1.jpg file2.mp4

# Process entire directory recursively
./build/evidex -send -server https://analysis.server.com/api/evidence -token YOUR_TOKEN -case-id CASE-2026-001 -r /path/to/folder

# With case ID and examiner notes
./build/evidex -send -server https://analysis.server.com/api/evidence -token YOUR_TOKEN -case-id CASE-2026-001 -notes "Initial acquisition" evidence.jpg
```

**Production Mode Benefits**:
- No local storage required (uses temporary directory, auto-cleaned)
- Direct transmission to analysis server
- Minimal forensic footprint on acquisition system
- Ideal for field operations

#### Debug Mode (Local Storage Only)

In debug mode, evidence is saved locally only (no transmission):

```bash
# Save evidence package locally
./build/evidex -o ./evidence image.jpg

# Process multiple files with local storage
./build/evidex -o ./evidence photo1.jpg video1.mp4 document.pdf

# Recursively process directory with local storage
./build/evidex -o ./evidence -r /path/to/folder

# Use SHA-512 hashing
./build/evidex -o ./evidence -hash SHA-512 file.jpg
```

**Debug Mode Benefits**:
- Local copy for validation
- Offline operation support
- Evidence package inspection
- Testing and development

**Note**: Evidex acquires complete files and comprehensive metadata (EXIF, file properties, codecs, etc.) for forensic analysis. Source files are **never modified** - only read in a non-destructive manner.

**Important**: The `-send` and `-o` flags are mutually exclusive. Use `-send` for production (remote transmission) or `-o` for debug (local storage only).

### Command-Line Options

```
-h, -help              Show help message
-v, -version           Show version information
-o, -output DIR        Output directory for evidence package (required for debug mode, optional with -send)
-r, -recursive         Recursively process directories
-hash ALGORITHM        Hash algorithm: SHA256 (default), SHA512
-case-id ID            Case identifier for correlation (e.g., CASE-2026-001)
-notes NOTES          Examiner notes for chain of custody
-json                  Export JSON metadata (default: true)
-csv                   Export CSV file manifest
-hashes                Create HASHES.txt file for verification
-report                Create integrity report
-verify                Verify hashes after acquisition (default: true)
-send                  Send evidence package to remote server (production mode)
-server URL            Remote server endpoint for transmission (required with -send)
-token TOKEN           Authentication token for server communication
```

**Mode Selection**:
- **Production Mode**: Use `-send -server URL -token TOKEN` (transmits to server, no local storage)
- **Debug Mode**: Use `-o DIR` (saves locally only, no transmission)
- **Note**: `-send` and `-o` are mutually exclusive

### Remote Transmission

Evidex transmits complete evidence packages (files + metadata) directly to remote analysis server endpoints:

```bash
# Production mode: Direct transmission without local storage
./build/evidex -send -server https://evidence-analysis.server.com:8443/api/evidence -token YOUR_TOKEN image.jpg

# The tool will:
# 1. Read files in read-only mode (never modifies source)
# 2. Extract comprehensive metadata from all file types
# 2. Create metadata package with complete chain of custody
# 3. Include RFC 5424 structured logging with metadata details
# 4. Intelligently chunk large metadata sets
# 5. Transmit securely to analysis endpoint
# 6. Verify successful transmission
# 7. Source files remain unmodified at origin
```

**Evidence Package Transmitted**:
- **Complete file content** (binary data, base64-encoded in JSON)
- File system metadata (permissions, timestamps, ownership, size)
- Image metadata (EXIF, XMP, IPTC, GPS, color profiles)
- Video metadata (codecs, resolution, bitrate, duration, frame rate)
- Document metadata (author, creation date, modification history)
- Archive metadata (compression type, member information)
- Cryptographic hashes (SHA-256, SHA-512)
- Chain of custody information
- RFC 5424 structured acquisition logs

### Comprehensive Logging

All operations are logged using RFC 5424 syslog format with structured data:

- **Log Output**: Real-time console output + in-memory buffer
- **Format**: RFC 5424 compliant with metadata
- **Transmission**: Logs included in evidence package for server
- **Levels**: Debug, Info, Warning, Error
- **Metadata**: Rich contextual information in each log entry

Example log entry:
```
<134>1 2026-01-06T10:30:45.123456Z hostname evidex 1234 - [meta@1 file_path="/path/to/file" file_size="1024"] Successfully acquired file
```

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
