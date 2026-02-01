package config

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// Config holds CLI configuration
type Config struct {
	OutputDir     string
	FilePaths     []string
	Recursive     bool
	HashAlgorithm string
	CaseID        string
	Notes         string
	ExportJSON    bool
	ExportCSV     bool
	ExportHashes  bool
	CreateReport  bool
	VerifyAfter   bool
	SendToServer  bool
	ServerURL     string
	AuthToken     string
}

const usage = `Evidex - Forensic Evidence Acquisition Binary
A portable, auditable forensic tool for evidence collection with chain of custody

Usage:
  evidex [options] <file|directory> [file|directory]...

Options:
  -h, -help              Show this help message
  -v, -version           Show version information
  -o, -output DIR        Output directory for evidence package (required for debug mode, optional for production mode with -send)
  -r, -recursive         Recursively process directories
  -hash ALGORITHM        Hash algorithm: SHA256 (default), SHA512
  -desc DESCRIPTION      Evidence description for manifest
  -notes NOTES          Examiner notes for chain of custody
  -json                  Export JSON metadata (default: true)
  -csv                   Export CSV file manifest
  -hashes                Create HASHES.txt file for verification
  -report                Create integrity report
  -verify                Verify hashes after acquisition (default: true)
  -send                  Send evidence package to remote server
  -server URL            Remote server endpoint URL
  -token TOKEN           Authentication token for remote server

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

  # Production mode: Send to server without local storage
  evidex -send -server https://server.com/api/evidence -token TOKEN image.jpg

  # Debug mode: Save locally only (no transmission)
  evidex -o ./evidence image.jpg

Chain of Custody:
  All evidence is acquired in read-only mode without any modifications.
  Every operation is logged and included in the forensic manifest.
  Cryptographic hashes are calculated for integrity verification.
`

// ParseFlags parses command-line flags
func ParseFlags() *Config {
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
	fs.StringVar(&cfg.CaseID, "case-id", "", "Case identifier for correlation")
	fs.StringVar(&cfg.Notes, "notes", "", "Examiner notes")
	fs.BoolVar(&cfg.ExportJSON, "json", true, "Export JSON metadata")
	fs.BoolVar(&cfg.ExportCSV, "csv", false, "Export CSV manifest")
	fs.BoolVar(&cfg.ExportHashes, "hashes", true, "Export hash file")
	fs.BoolVar(&cfg.CreateReport, "report", true, "Create integrity report")
	fs.BoolVar(&cfg.VerifyAfter, "verify", true, "Verify hashes")
	fs.BoolVar(&cfg.SendToServer, "send", false, "Send evidence package to remote server")
	fs.StringVar(&cfg.ServerURL, "server", "", "Remote server endpoint URL")
	fs.StringVar(&cfg.AuthToken, "token", "", "Authentication token for remote server")

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
		fmt.Printf("Evidex v1.0.0\n")
		os.Exit(0)
	}

	// Remaining arguments are file paths
	cfg.FilePaths = fs.Args()

	return cfg
}

// ValidateConfig validates the configuration
func ValidateConfig(cfg *Config) error {
	// -send and -o are mutually exclusive
	if cfg.SendToServer && cfg.OutputDir != "" {
		return fmt.Errorf("cannot use -send and -o together: use -send for production mode (remote transmission) or -o for debug mode (local storage only)")
	}

	// Output directory is required unless sending to server (production mode)
	if cfg.OutputDir == "" && !cfg.SendToServer {
		return fmt.Errorf("output directory (-o) is required in debug mode (or use -send for production mode)")
	}

	// When sending to server, validate server URL
	if cfg.SendToServer && cfg.ServerURL == "" {
		return fmt.Errorf("server URL (-server) is required when using -send")
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
