# Evidex Deployment

## Requirements

- Go 1.25 or higher
- Make utility
- Git

## Building

```bash
# Clone repository
git clone https://github.com/yourusername/evidex.git
cd evidex

# Build for current platform
make build

# Build for all platforms
make build-all
```

## Execution

### Production Mode (Remote Transmission)
```bash
# Acquire and transmit evidence
./build/evidex -send -server https://analysis.server.com/api/evidence -token YOUR_TOKEN file.jpg

# Process directory recursively
./build/evidex -send -server https://analysis.server.com/api/evidence -token YOUR_TOKEN -r /path/to/folder
```

### Debug Mode (Local Storage)
```bash
# Save evidence package locally
./build/evidex -o ./evidence file.jpg

# Process directory with local storage
./build/evidex -o ./evidence -r /path/to/folder
```

### Key Flags
- `-send` - Transmit to remote server (production mode)
- `-server URL` - Analysis server endpoint
- `-token TOKEN` - Authentication token
- `-o DIR` - Output directory (debug mode)
- `-r` - Recursive directory processing
- `-hash ALGORITHM` - SHA-256 (default) or SHA-512
- `-case-id ID` - Case identifier
- `-notes TEXT` - Examiner notes
