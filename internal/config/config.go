// Package config provides configuration management for the Evidex application.
package config

import (
	"flag"
	"fmt"
	"os"
)

// Config holds CLI configuration.
type Config struct {
	FilePaths []string
	Recursive bool
	CaseID    string
	ServerURL string
	AuthToken string
}

const usage = `Evidex - Forensic Evidence Acquisition Binary
A portable, auditable forensic tool for evidence collection with chain of custody

Usage:
  evidex [options] <file|directory> [file|directory]...

Options:
  -h, --help              Show this help message
  -v, --version           Show version information
  -r, --recursive         Recursively process directories
  --case-id ID            Case identifier for correlation (required)
  --server URL            Remote server endpoint URL (required)
  --token TOKEN           Authentication token for remote server (required)

Examples:
  # Acquire a single file and send to server
  evidex --server https://server.com/api/evidence --token TOKEN --case-id "ABC-2025-001" image.jpg

  # Acquire multiple files
  evidex --server https://server.com/api/evidence --token TOKEN --case-id "ABC-2025-001" photo1.jpg video1.mp4

  # Recursively acquire all files in a directory
  evidex --server https://server.com/api/evidence --token TOKEN --case-id "ABC-2025-001" -r /path/to/folder

Chain of Custody:
  All evidence is acquired in read-only mode without any modifications.
  Every operation is logged and included in the forensic manifest.
  All files are hashed with MD5, SHA1, and SHA256 for integrity verification.
  Evidence is transmitted directly to the server without local storage.
`

// ParseFlags parses command-line flags.
func ParseFlags() *Config {
	cfg := &Config{}

	fs := flag.NewFlagSet("evidex", flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprint(os.Stderr, usage)
	}

	fs.BoolVar(&cfg.Recursive, "r", false, "Recursively process directories")
	fs.BoolVar(&cfg.Recursive, "recursive", false, "Recursively process directories")
	fs.StringVar(&cfg.CaseID, "case-id", "", "Case identifier for correlation")
	fs.StringVar(&cfg.ServerURL, "server", "", "Remote server endpoint URL")
	fs.StringVar(&cfg.AuthToken, "token", "", "Authentication token for remote server")

	// Handle help and version
	helpFlag := fs.Bool("help", false, "Show help message")
	fs.Bool("v", false, "Show version")
	fs.Bool("h", false, "Show help message")
	versionFlag := fs.Bool("version", false, "Show version")

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

// ValidateConfig validates the configuration.
func ValidateConfig(cfg *Config) error {
	// Server URL is required
	if cfg.ServerURL == "" {
		return fmt.Errorf("server URL (--server) is required for evidence transmission")
	}

	// Case ID is required
	if cfg.CaseID == "" {
		return fmt.Errorf("case ID (--case-id) is required for evidence correlation")
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

	return nil
}
