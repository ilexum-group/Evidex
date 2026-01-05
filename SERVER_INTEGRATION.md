# Evidex Server Integration Guide

## Overview

Evidex now supports remote transmission of evidence packages to a central server endpoint, similar to Tracium's architecture.

## Architecture Components

### 1. Enhanced Logging (`internal/logger/`)

RFC 5424 compliant syslog format with structured data:

```go
// Example usage
logger.LogInfo("File acquired", map[string]string{
    "file_path": "/path/to/file.jpg",
    "file_size": "1024",
    "hash": "abc123...",
})
```

**Features**:
- Structured metadata in each log entry
- In-memory log buffering for transmission
- Thread-safe concurrent logging
- Multiple severity levels (Debug, Info, Warning, Error)

### 2. Server Sender (`internal/sender/`)

Handles transmission of evidence packages to remote endpoints:

```go
// Send complete evidence package
sender.SendEvidencePackage("https://server.com/api/evidence", authToken, evidencePackage)

// Send individual large files with chunking
sender.SendEvidenceFile("https://server.com/api/evidence", authToken, filePath, metadata)
```

**Features**:
- Intelligent chunking for large files (64 MB chunks)
- JSON payload serialization
- HTTP/HTTPS transmission with authentication
- Automatic retry with logging
- Progress tracking and reporting

### 3. Updated Models

Evidence package model includes:
- `Logs []string` - All acquisition logs
- `ServerURL string` - Transmission endpoint
- `AuthToken string` - Authentication token

## Implementation in Main

### 1. Initialize Logger

```go
import "github.com/evidex/internal/logger"

func main() {
    // Initialize global logger
    if err := logger.InitDefaultLogger(); err != nil {
        fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
        os.Exit(1)
    }
    
    // All logging now uses RFC 5424 format
    logger.LogInfo("Starting acquisition", map[string]string{
        "version": "1.0.0",
    })
}
```

### 2. Add CLI Flags

```go
// Add to parseFlags()
serverURL := flag.String("server", "", "Remote server endpoint")
authToken := flag.String("token", "", "Authentication token for server")

cfg.ServerURL = *serverURL
cfg.AuthToken = *authToken
```

### 3. Transmit After Acquisition

```go
import "github.com/evidex/internal/sender"

// After creating evidence package
pkg := acquirer.GetEvidencePackage()

// Attach logs to package
pkg.Logs = logger.GetLogs()

// Transmit if server is configured
if cfg.ServerURL != "" {
    logger.LogInfo("Transmitting evidence to server", map[string]string{
        "server": cfg.ServerURL,
        "package_id": pkg.Manifest.ID,
        "file_count": fmt.Sprintf("%d", len(pkg.Files)),
    })
    
    if err := sender.SendEvidencePackage(cfg.ServerURL, cfg.AuthToken, pkg); err != nil {
        logger.LogError("Failed to transmit package", map[string]string{
            "error": err.Error(),
        })
    } else {
        logger.LogInfo("Evidence transmitted successfully", map[string]string{
            "package_id": pkg.Manifest.ID,
        })
    }
}
```

## Log Format

RFC 5424 format with structured data:

```
<PRI>1 TIMESTAMP HOSTNAME APP-NAME PROCID MSGID [META@1 key="value" ...] MESSAGE
```

Example:
```
<134>1 2026-01-06T10:30:45.123Z hostname evidex 1234 - [meta@1 file_path="/evidence/file.jpg" size="2048" hash="abc123"] File successfully acquired
```

**Priority Calculation**:
- Facility: 16 (local0) = 16 × 8 = 128
- Severity: 6 (Info) = 6
- Priority = 128 + 6 = 134

## Server Expectations

The remote server should:

### Endpoint Requirements
- **Path**: `/api/evidence` (configurable)
- **Method**: `POST`
- **Content-Type**: `application/json`
- **Headers**:
  - `Authorization: Bearer <token>`
  - `User-Agent: Evidex-Agent/1.0`

### Request Body

```json
{
  "version": "1.0.0",
  "created_at": "2026-01-06T10:30:45Z",
  "manifest": {
    "id": "EV-2026-01-06-ABC123",
    "created_at": "2026-01-06T10:30:45Z",
    "created_by_user": "investigator",
    "created_by_hostname": "forensic-workstation",
    "file_count": 5,
    "total_size": 1048576,
    "hash_algorithm": "SHA-256",
    "integrity": "VERIFIED"
  },
  "files": [
    {
      "source_path": "/mnt/evidence/image.jpg",
      "file_size": 1024,
      "hashes": {
        "sha256": "abc123...",
        "sha512": "def456..."
      },
      "verified": true
    }
  ],
  "logs": [
    "<134>1 2026-01-06T10:30:45.123Z hostname evidex 1234 - [meta@1 ...] File acquired",
    "<134>1 2026-01-06T10:30:46.456Z hostname evidex 1234 - [meta@1 ...] Hash verified"
  ]
}
```

### Chunked File Transmission

For files > 64 MB:

```json
{
  "type": "evidence_chunk",
  "file_path": "/evidence/large_video.mp4",
  "chunk_num": 1,
  "total_chunks": 100,
  "chunk_size": 67108864,
  "file_size": 6710886400,
  "metadata": { ... },
  "data": "<base64-encoded-chunk>"
}
```

## Usage Examples

### Acquire locally only
```bash
evidex -o ./evidence -desc "Case: 2026-001" image.jpg
```

### Acquire and transmit to server
```bash
evidex -o ./evidence \
  -desc "Case: 2026-001" \
  -server https://evidence.company.com:8443/api/evidence \
  -token "secret-token-123" \
  image.jpg
```

### Batch transmission of multiple files
```bash
evidex -o ./evidence \
  -r \
  -server https://evidence.company.com:8443/api/evidence \
  -token "secret-token-123" \
  /mnt/devices/phone
```

## Security Considerations

### Authentication
- Use HTTPS for all transmissions
- Store tokens securely (environment variables, credential managers)
- Rotate tokens regularly

### Encryption
- HTTPS TLS 1.3 minimum recommended
- Consider end-to-end encryption for sensitive evidence
- Verify server certificate validity

### Logging
- Logs contain metadata but not sensitive data
- Log transmission can be disabled
- Logs are immutable once included in package

### Package Integrity
- SHA-256 hashes verified locally before transmission
- Server should verify hashes on receipt
- Evidence package is read-only after creation

## Testing

```bash
# Test logging
go test ./internal/logger -v

# Test sender
go test ./internal/sender -v

# Integration test
./build/evidex -o ./test-evidence \
  -desc "Test case" \
  test-file.jpg
```

## Troubleshooting

### Connection refused
- Verify server URL is correct
- Check firewall/network connectivity
- Verify HTTPS certificate is valid

### Authentication failed (401)
- Verify token is correct
- Check token hasn't expired
- Verify token format (Bearer prefix)

### Transmission timeout
- Increase client timeout in sender.go
- Check network bandwidth
- Consider chunking for large files

### Log format issues
- Verify RFC 5424 compliance
- Check structured data escaping
- Review syslog server configuration

## Future Enhancements

- [ ] Log rotation and archival
- [ ] Signature verification on server
- [ ] Batch evidence transmission
- [ ] Real-time progress webhooks
- [ ] Evidence retrieval/verification endpoint
- [ ] Distributed server clusters
