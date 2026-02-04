# Evidex Deployment

## Requirements

- Go 1.25 or higher
- Make utility
- Git

## Architecture

Evidex uses a **dependency injection pattern** with the following components:

**Core Modules**:
- **OS Wrapper**: Platform-specific operations with automatic command logging
- **Metadata Manager**: File analysis, hash calculation (MD5, SHA1, SHA256, SHA512)
- **Logger**: RFC 5424 structured logging with metadata
- **Acquirer**: Evidence acquisition orchestration
- **Custody Chain**: Immutable chain of custody tracking

**Key Features**:
- Automatic command logging via OS wrapper
- Cross-platform compilation with build tags
- Native `time.Time` handling for timestamps
- Dependency injection for modularity and testing

## Building

### Cross-Platform Compilation

Evidex supports cross-compilation using Go build tags:

```bash
# Clone repository
git clone https://github.com/yourusername/evidex.git
cd evidex

# Build for current platform
make build

# Build for specific platforms
make build-linux         # Linux (amd64, arm64)
make build-darwin        # macOS (amd64, arm64)
make build-windows       # Windows (amd64, arm64)
make build-freebsd-amd64 # FreeBSD
make build-openbsd-amd64 # OpenBSD

# Build for all platforms
make build-all
```

**Build Tags**:
- `//go:build linux` - Linux-specific code
- `//go:build windows` - Windows-specific code
- `//go:build darwin` - macOS-specific code
- Stub files (`*_stub.go`) enable cross-compilation

**Benefits**:
- Compile for any platform from any platform
- No missing platform-specific code errors
- Graceful fallback to default implementations

## Initialization Sequence

When Evidex starts, it follows this initialization sequence:

1. **Parse Configuration**: Command-line flags and validation
2. **Create OS Wrapper**: Platform-specific implementation (Windows/Linux/macOS/BSD)
3. **Get System Info**: Hostname, username, process ID
4. **Initialize Logger**: RFC 5424 structured logging
5. **Create Custody Chain**: Evidence tracking with unique ID
6. **Enable Auto-Logging**: OS operations automatically logged
7. **Create Metadata Manager**: File analysis and hash calculation
8. **Create Acquirer**: Evidence acquisition with all dependencies
9. **Execute Acquisition**: Process files/directories

This ensures all components are properly initialized before evidence collection begins.

## Execution

Evidex operates in memory-only mode, transmitting evidence directly to the server:

```bash
# Acquire and transmit evidence
./build/evidex --server https://analysis.server.com/api/evidence --token YOUR_TOKEN --case-id CASE-001 file.jpg

# Process directory recursively
./build/evidex --server https://analysis.server.com/api/evidence --token YOUR_TOKEN --case-id CASE-001 -r /path/to/folder
```

### Key Flags
- `--server URL` - Analysis server endpoint (required)
- `--token TOKEN` - Authentication token (required)
- `--case-id ID` - Case identifier (required)
- `-r, --recursive` - Recursive directory processing

### Features
- **Memory-only processing**: No local storage
- **Automatic hashing**: MD5, SHA1, SHA256 calculated for all files
- **Command logging**: All OS and file operations logged
- **Direct transmission**: Evidence sent directly to server
