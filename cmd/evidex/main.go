// Package main implements the Evidex CLI for forensic evidence acquisition
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/evidex/internal/acquisition"
	"github.com/evidex/internal/formatter"
	"github.com/evidex/internal/utils"
)

const (
	version = "1.0.0"
	usage   = `Evidex - Forensic Evidence Acquisition Binary
A portable, auditable forensic tool for evidence collection with chain of custody

Usage:
  evidex [options] <file|directory> [file|directory]...

Options:
  -h, -help              Show this help message
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

Examples:
  # Acquire a single image file
  evidex -o ./evidence image.jpg

  # Acquire multiple files
  evidex -o ./evidence photo1.jpg video1.mp4 photo2.jpg

  # Recursively acquire all files in a directory
  evidex -o ./evidence -r /path/to/folder

  # Acquire with custom description
  evidex -o ./evidence -desc "Case: ABC-2025-001" image.jpg

  # Acquire with SHA-512 hashing
  evidex -o ./evidence -hash SHA512 file.jpg

Chain of Custody:
  All evidence is acquired in read-only mode without any modifications.
  Every operation is logged and included in the forensic manifest.
  Cryptographic hashes are calculated for integrity verification.
`
)

// Config holds CLI configuration
type Config struct {
	OutputDir     string
	FilePaths     []string
	Recursive     bool
	HashAlgorithm string
	Description   string
	Notes         string
	ExportJSON    bool
	ExportCSV     bool
	ExportHashes  bool
	CreateReport  bool
	VerifyAfter   bool
}

func main() {
	// Initialize logging
	utils.InitLogging()

	cfg := parseFlags()

	// Validate required flags
	if err := validateConfig(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		fmt.Fprintf(os.Stderr, "\nUse -help for usage information\n")
		os.Exit(1)
	}

	utils.LogInfo("Starting Evidex", map[string]string{"version": version})
	utils.LogInfo("Output directory", map[string]string{"path": cfg.OutputDir})

	// Create output directory
	if err := utils.EnsureDirectory(cfg.OutputDir); err != nil {
		utils.LogError("Failed to create output directory", map[string]string{"path": cfg.OutputDir, "error": err.Error()})
		os.Exit(1)
	}

	// Create acquirer
	acquirer := acquisition.NewAcquirer(cfg.OutputDir)

	// Set hash algorithm
	if cfg.HashAlgorithm != "" {
		acquirer.SetHashAlgorithm(cfg.HashAlgorithm)
		utils.LogInfo("Hash algorithm set", map[string]string{"algorithm": cfg.HashAlgorithm})
	}

	// Process input files/directories
	for _, filePath := range cfg.FilePaths {
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			utils.LogError("File not found", map[string]string{"path": filePath, "error": err.Error()})
			continue
		}

		if fileInfo.IsDir() {
			utils.LogInfo("Processing directory", map[string]string{"path": filePath, "recursive": fmt.Sprintf("%v", cfg.Recursive)})
			if err := acquirer.AcquireDirectory(filePath, cfg.Recursive); err != nil {
				utils.LogError("Directory acquisition failed", map[string]string{"path": filePath, "error": err.Error()})
			}
		} else {
			utils.LogInfo("Processing file", map[string]string{"path": filePath})
			if err := acquirer.AcquireFile(filePath); err != nil {
				utils.LogError("File acquisition failed", map[string]string{"path": filePath, "error": err.Error()})
			}
		}
	}

	// Copy files to package
	utils.LogInfo("Copying files to package", map[string]string{})
	if err := acquirer.CopyFilesToPackage(); err != nil {
		utils.LogError("Failed to copy files", map[string]string{"error": err.Error()})
		os.Exit(1)
	}

	// Build evidence package
	utils.LogInfo("Building evidence package", map[string]string{})
	pkg := acquirer.GetEvidencePackage()

	// Set manifest properties
	if cfg.Description != "" {
		pkg.Manifest.EvidenceDescription = cfg.Description
	}
	if cfg.Notes != "" {
		pkg.Manifest.ExaminersNotes = cfg.Notes
	}

	// Export formats
	pkgFormatter := formatter.NewPackageFormatter(cfg.OutputDir, pkg)

	if cfg.ExportJSON {
		utils.LogInfo("Exporting JSON metadata", map[string]string{})
		if err := pkgFormatter.ExportJSON(); err != nil {
			utils.LogError("JSON export failed", map[string]string{"error": err.Error()})
		}
	}

	if cfg.ExportCSV {
		utils.LogInfo("Exporting CSV manifest", map[string]string{})
		if err := pkgFormatter.ExportCSV(); err != nil {
			utils.LogError("CSV export failed", map[string]string{"error": err.Error()})
		}
	}

	if cfg.ExportHashes {
		utils.LogInfo("Exporting hash file", map[string]string{})
		if err := pkgFormatter.ExportHashes(); err != nil {
			utils.LogError("Hashes export failed", map[string]string{"error": err.Error()})
		}
	}

	if cfg.CreateReport {
		utils.LogInfo("Creating integrity report", map[string]string{})
		if err := pkgFormatter.ExportIntegrityReport(); err != nil {
			utils.LogError("Integrity report failed", map[string]string{"error": err.Error()})
		}
	}

	// Create README
	utils.LogInfo("Creating package documentation", map[string]string{})
	if err := pkgFormatter.CreateReadme(); err != nil {
		utils.LogError("README creation failed", map[string]string{"error": err.Error()})
	}

	// Package compression info
	if err := pkgFormatter.CompressPackage("tar.gz"); err != nil {
		utils.LogWarn("Compression metadata failed", map[string]string{"error": err.Error()})
	}

	// Summary
	fmt.Printf("\n")
	fmt.Printf("================================\n")
	fmt.Printf("FORENSIC EVIDENCE PACKAGE CREATED\n")
	fmt.Printf("================================\n")
	fmt.Printf("Evidence ID: %s\n", pkg.Manifest.ID)
	fmt.Printf("Files Acquired: %d\n", pkg.Manifest.FileCount)
	fmt.Printf("Total Size: %d bytes\n", pkg.Manifest.TotalSize)
	fmt.Printf("Output Directory: %s\n", cfg.OutputDir)
	fmt.Printf("Hash Algorithm: %s\n", pkg.Manifest.HashAlgorithm)
	fmt.Printf("Acquisition Status: %s\n", pkg.Manifest.Integrity)
	fmt.Printf("\nPackage Contents:\n")
	fmt.Printf("  - files/           (Acquired evidence files)\n")
	fmt.Printf("  - metadata/        (File metadata and hashes)\n")
	fmt.Printf("  - HASHES.txt       (Hash verification)\n")
	fmt.Printf("  - README.txt       (Package information)\n")
	fmt.Printf("  - INTEGRITY_REPORT.txt (Verification report)\n")
	fmt.Printf("\nTo verify integrity:\n")
	fmt.Printf("  Compare file hashes with HASHES.txt\n")
	fmt.Printf("  Review manifest.json for chain of custody\n")
	fmt.Printf("\n")

	utils.LogInfo("Evidex acquisition completed", map[string]string{"status": "successfully"})
}

// parseFlags parses command-line flags
func parseFlags() *Config {
	cfg := &Config{
		Recursive:     false,
		HashAlgorithm: "SHA-256",
		ExportJSON:    true,
		ExportCSV:     false,
		ExportHashes:  true,
		CreateReport:  true,
		VerifyAfter:   true,
	}

	fs := flag.NewFlagSet("evidex", flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprint(os.Stderr, usage)
	}

	fs.StringVar(&cfg.OutputDir, "o", "", "Output directory for evidence package")
	fs.StringVar(&cfg.OutputDir, "output", "", "Output directory for evidence package")
	fs.BoolVar(&cfg.Recursive, "r", false, "Recursively process directories")
	fs.BoolVar(&cfg.Recursive, "recursive", false, "Recursively process directories")
	fs.StringVar(&cfg.HashAlgorithm, "hash", "SHA-256", "Hash algorithm")
	fs.StringVar(&cfg.Description, "desc", "", "Evidence description")
	fs.StringVar(&cfg.Notes, "notes", "", "Examiner notes")
	fs.BoolVar(&cfg.ExportJSON, "json", true, "Export JSON metadata")
	fs.BoolVar(&cfg.ExportCSV, "csv", false, "Export CSV manifest")
	fs.BoolVar(&cfg.ExportHashes, "hashes", true, "Export hash file")
	fs.BoolVar(&cfg.CreateReport, "report", true, "Create integrity report")
	fs.BoolVar(&cfg.VerifyAfter, "verify", true, "Verify hashes")

	// Handle help and version
	helpFlag := fs.Bool("help", false, "Show help message")
	versionFlag := fs.Bool("v", false, "Show version")
	_ = fs.Bool("h", false, "Show help message")
	_ = fs.Bool("version", false, "Show version")

	if err := fs.Parse(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		os.Exit(1)
	}

	if *helpFlag {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(0)
	}

	if *versionFlag {
		fmt.Printf("Evidex v%s\n", version)
		os.Exit(0)
	}

	// Remaining arguments are file paths
	cfg.FilePaths = fs.Args()

	return cfg
}

// validateConfig validates the configuration
func validateConfig(cfg *Config) error {
	if cfg.OutputDir == "" {
		return fmt.Errorf("output directory (-o) is required")
	}

	if len(cfg.FilePaths) == 0 {
		return fmt.Errorf("at least one file or directory path is required")
	}

	// Validate file paths exist
	for _, filePath := range cfg.FilePaths {
		if _, err := os.Stat(filePath); err != nil {
			return fmt.Errorf("file not found: %s", filePath)
		}
	}

	// Validate hash algorithm
	validAlgs := []string{"SHA-256", "SHA-512", "SHA256", "SHA512"}
	validAlg := false
	for _, alg := range validAlgs {
		if strings.ToUpper(cfg.HashAlgorithm) == alg {
			validAlg = true
			cfg.HashAlgorithm = alg
			break
		}
	}
	if !validAlg {
		return fmt.Errorf("invalid hash algorithm: %s", cfg.HashAlgorithm)
	}

	return nil
}
