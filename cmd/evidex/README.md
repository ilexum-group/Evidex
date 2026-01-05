# Evidex CLI Documentation

## Overview

The Evidex command-line interface provides comprehensive options for evidence acquisition with chain of custody compliance.

## Basic Invocation

```bash
evidex [options] <file|directory> [file|directory]...
```

## Command-Line Flags

### Required Flags

#### `-o, -output DIR`
**Output directory for evidence package** (REQUIRED)

The directory where the evidence package will be created. Must be writable.

```bash
evidex -o ./evidence image.jpg
evidex -output /path/to/evidence/storage image.jpg
```

If directory doesn't exist, it will be created automatically.

### Input Selection

#### `<file|directory>` (positional arguments)
**One or more files or directories to acquire**

Files and directories must exist and be readable.

```bash
# Single file
evidex -o ./evidence image.jpg

# Multiple files
evidex -o ./evidence file1.jpg file2.jpg file3.mp4

# Single directory
evidex -o ./evidence /path/to/folder

# Mixed input
evidex -o ./evidence image.jpg /path/to/videos/ photo.png
```

### Optional Flags

#### `-r, -recursive`
**Recursively process directories**

When specified, directories are traversed recursively to include all files in subdirectories.

```bash
# Non-recursive (only direct children)
evidex -o ./evidence /path/to/folder

# Recursive (include subdirectories)
evidex -o ./evidence -r /path/to/folder
```

Default: `false` (non-recursive)

#### `-hash ALGORITHM`
**Hash algorithm to use**

Specifies the primary cryptographic hash algorithm.

```bash
# SHA-256 (default, recommended)
evidex -o ./evidence image.jpg
evidex -o ./evidence -hash SHA256 image.jpg

# SHA-512 (stronger security)
evidex -o ./evidence -hash SHA512 image.jpg
```

Valid values: `SHA-256`, `SHA-512`, `SHA256`, `SHA512`

Default: `SHA-256`

#### `-desc DESCRIPTION`
**Evidence description for manifest**

Adds case-specific description or case number to the manifest.

```bash
evidex -o ./evidence -desc "Case: ABC-2025-001" image.jpg
evidex -o ./evidence -desc "Suspect Device - Serial #12345" file.jpg
evidex -o ./evidence -desc "Court Order #2025-123456" evidence.jpg
```

This text appears in `manifest.json` under `evidence_description`.

#### `-notes NOTES`
**Examiner notes for chain of custody**

Additional narrative notes from the examining forensic professional.

```bash
evidex -o ./evidence \
  -notes "Files collected from suspect's USB drive during execution of search warrant" \
  image.jpg

evidex -o ./evidence \
  -notes "Evidence acquired under chain of custody protocol version 3.2" \
  file.jpg
```

This text appears in `manifest.json` under `examiners_notes`.

### Export Format Flags

#### `-json`
**Export JSON metadata**

Controls JSON export of metadata files.

```bash
# Enable JSON (default)
evidex -o ./evidence -json image.jpg
evidex -o ./evidence image.jpg

# Disable JSON
evidex -o ./evidence -json=false image.jpg
```

Default: `true`

Creates:
- `metadata/manifest.json`
- `metadata/acquisition_log.json`
- `metadata/system_context.json`
- `metadata/file_catalog.json`

#### `-csv`
**Export CSV file manifest**

Creates CSV file for import into analysis tools.

```bash
# Enable CSV export
evidex -o ./evidence -csv image.jpg
```

Default: `false`

Creates: `metadata/file_manifest.csv`

Columns: `Filename,Source Path,Size,Modified Time,Hash,Verified,File Type,Owner`

#### `-hashes`
**Create HASHES.txt file**

Creates simple text file with hashes for verification.

```bash
# Enable (default)
evidex -o ./evidence image.jpg

# Disable
evidex -o ./evidence -hashes=false image.jpg
```

Default: `true`

Creates: `HASHES.txt` with per-file hashes for verification

#### `-report`
**Create integrity report**

Creates summary integrity verification report.

```bash
# Enable (default)
evidex -o ./evidence image.jpg

# Disable
evidex -o ./evidence -report=false image.jpg
```

Default: `true`

Creates: `INTEGRITY_REPORT.txt` with acquisition summary and statistics

### Verification Flags

#### `-verify`
**Verify hashes after acquisition**

Controls hash verification during acquisition.

```bash
# Enable (default)
evidex -o ./evidence image.jpg

# Disable verification
evidex -o ./evidence -verify=false image.jpg
```

Default: `true`

When enabled, files are re-hashed after copying to verify integrity.

### Information Flags

#### `-h, -help`
**Show help message**

```bash
evidex -help
evidex -h
```

Displays usage information and available options.

#### `-v, -version`
**Show version information**

```bash
evidex -v
evidex -version
```

Displays Evidex version number.

## Examples

### Basic Image Acquisition

```bash
# Single image
evidex -o ./evidence image.jpg

# Multiple images
evidex -o ./evidence photo1.jpg photo2.jpg photo3.jpg
```

### Recursive Directory Acquisition

```bash
# All files in directory and subdirectories
evidex -o ./evidence -r /mnt/usb/photos

# With case description
evidex -o ./evidence -r -desc "Case XYZ-2025-001" /mnt/usb
```

### Enhanced Hash Security

```bash
# Use SHA-512 for maximum security
evidex -o ./evidence -hash SHA512 image.jpg

# For cold case archives (long-term preservation)
evidex -o ./evidence -r -hash SHA512 -desc "Cold Case Archive" /evidence/folder
```

### Full Forensic Acquisition

```bash
# Complete acquisition with all exports and notes
evidex -o ./evidence \
  -r \
  -hash SHA512 \
  -csv \
  -report \
  -desc "Case: Criminal Investigation ABC-2025-001" \
  -notes "Evidence collected per warrant #2025-456789 on 2025-12-25 by Agent Smith" \
  /mnt/evidence/drive1
```

### Batch Processing

```bash
# Acquire multiple separate evidence sources
evidex -o ./evidence_1 -desc "Device A" /mnt/device_a
evidex -o ./evidence_2 -desc "Device B" /mnt/device_b
evidex -o ./evidence_3 -desc "USB Drive" /mnt/usb_drive
```

### High-Volume Acquisition

```bash
# Large directory structure with SHA-512
evidex -o ./archive \
  -r \
  -hash SHA512 \
  -csv \
  -report \
  /large/evidence/collection
```

## Output and Feedback

### Console Output

During acquisition, Evidex provides:

```
[INFO] Starting Evidex: Version 1.0.0
[INFO] Output directory: ./evidence
[INFO] Hash algorithm: SHA-256
[INFO] [AcquireFile] Starting acquisition of: image.jpg
[INFO] [ExtractMetadata] Extracted metadata
[INFO] [HashCalculation] SHA-256: a1b2c3d4e5f6...
[INFO] [FileCopied] Successfully acquired: image.jpg (5242880 bytes)
```

### Summary Report

After completion:

```
================================
FORENSIC EVIDENCE PACKAGE CREATED
================================
Evidence ID: EVDX-HOSTNAME-20260105181234
Files Acquired: 15
Total Size: 2147483648 bytes
Output Directory: ./evidence
Hash Algorithm: SHA-256
Acquisition Status: verified

Package Contents:
  - files/           (Acquired evidence files)
  - metadata/        (File metadata and hashes)
  - HASHES.txt       (Hash verification)
  - README.txt       (Package information)
  - INTEGRITY_REPORT.txt (Verification report)

To verify integrity:
  Compare file hashes with HASHES.txt
  Review manifest.json for chain of custody
```

## Exit Codes

- **0**: Successful completion
- **1**: Fatal error (invalid flags, missing files, output creation failure)
- **2**: Acquisition completed with errors (some files failed)

## Best Practices

### Case Management
```bash
# Use case number in description
evidex -o /evidence/CASE-2025-001234 \
  -desc "CASE 2025-001234: Evidence Acquisition" \
  /evidence/source
```

### Examiner Identification
```bash
# Include examiner information in notes
evidex -o ./evidence \
  -notes "Acquired by: John Doe, Badge #12345, 2026-01-05 18:00 UTC" \
  image.jpg
```

### Security Level Selection
```bash
# Standard: SHA-256 (daily investigations)
evidex -o ./evidence -hash SHA256 image.jpg

# Enhanced: SHA-512 (high-profile cases)
evidex -o ./evidence -hash SHA512 image.jpg

# Maximum: SHA-512 + CSV + Report (long-term archives)
evidex -o ./evidence -hash SHA512 -csv -report image.jpg
```

### Batch Processing Workflow
```bash
# Process multiple devices in sequence
for device in /mnt/device_*; do
  name=$(basename $device)
  evidex -o ./evidence_${name} \
    -r \
    -desc "Device: ${name}" \
    $device
done
```

## Legal Compliance Tips

1. **Always use `-desc`** to document case information
2. **Always use `-notes`** to document chain of custody details
3. **Archive with SHA-512** (`-hash SHA512`) for long-term evidence
4. **Export CSV** (`-csv`) for compatibility with legal analysis tools
5. **Create report** (`-report`) for court documentation
6. **Preserve output** in secure storage with access controls

## Troubleshooting

### "File not found"
```bash
# Verify file exists
ls -la /path/to/file.jpg

# Check correct path
evidex -o ./evidence /path/to/file.jpg
```

### "Permission denied"
```bash
# Verify read access
test -r /path/to/file.jpg && echo "readable" || echo "not readable"

# May need elevated privileges
sudo evidex -o ./evidence /path/to/file.jpg
```

### "Output directory creation failed"
```bash
# Create parent directory first
mkdir -p /path/to/evidence
evidex -o /path/to/evidence/package file.jpg

# Or use current directory
evidex -o ./evidence file.jpg
```

### "Invalid hash algorithm"
```bash
# Use supported algorithms
evidex -o ./evidence -hash SHA256 file.jpg    # Correct
evidex -o ./evidence -hash SHA512 file.jpg    # Correct
evidex -o ./evidence -hash SHA-256 file.jpg   # Also works

# These are invalid
evidex -o ./evidence -hash SHA1 file.jpg      # Not supported
evidex -o ./evidence -hash MD5 file.jpg       # Not supported
```

## Environment Variables

Currently, Evidex does not use environment variables. All configuration is through command-line flags.

Future versions may support:
- `EVIDEX_OUTPUT_DIR`
- `EVIDEX_HASH_ALGORITHM`
- `EVIDEX_EXAMINER_NAME`
- `EVIDEX_CASE_NUMBER`

## Version Information

Current version: **1.0.0**

Check version:
```bash
evidex -v
evidex -version
```

## Command Reference Summary

| Flag | Type | Default | Purpose |
|------|------|---------|---------|
| `-o, -output` | string | *required* | Output directory |
| `-r, -recursive` | bool | false | Recurse directories |
| `-hash` | string | SHA-256 | Hash algorithm |
| `-desc` | string | "" | Case description |
| `-notes` | string | "" | Examiner notes |
| `-json` | bool | true | Export JSON |
| `-csv` | bool | false | Export CSV |
| `-hashes` | bool | true | Export hashes |
| `-report` | bool | true | Create report |
| `-verify` | bool | true | Verify hashes |
| `-h, -help` | bool | false | Show help |
| `-v, -version` | bool | false | Show version |

## Legal Notice

**Evidex is intended for legitimate forensic investigations only.**

Users must:
- Comply with all applicable laws
- Maintain proper chain of custody
- Obtain necessary legal authorization
- Preserve original evidence
- Consult with legal counsel
- Follow organizational policies

**Unauthorized access to files or systems is illegal.**

---

**Documentation Version**: 1.0.0  
**Last Updated**: 2026-01-05
