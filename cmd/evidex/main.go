// Package main implements the Evidex CLI for forensic evidence acquisition
package main

import (
	"fmt"
	"os"

	"github.com/ilexum-group/evidex/internal/acquisition"
	"github.com/ilexum-group/evidex/internal/config"
	"github.com/ilexum-group/evidex/internal/formatter"
	"github.com/ilexum-group/evidex/internal/models"
	"github.com/ilexum-group/evidex/internal/sender"
	"github.com/ilexum-group/evidex/internal/utils"
)

const (
	version = "1.0.0"
)

func main() {
	// Initialize logging
	utils.InitLogging()

	cfg := config.ParseFlags()

	// Validate required flags
	if err := config.ValidateConfig(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		fmt.Fprintf(os.Stderr, "\nUse -help for usage information\n")
		os.Exit(1)
	}

	utils.LogInfo("Starting Evidex", map[string]string{"version": version})

	// Determine mode and setup
	if cfg.SendToServer {
		utils.LogInfo("Production mode: direct transmission (no disk writes)", map[string]string{})
		// Production mode: no working directory needed
		runProductionMode(cfg)
	} else {
		utils.LogInfo("Debug mode: local storage", map[string]string{"path": cfg.OutputDir})
		// Debug mode: use output directory
		runDebugMode(cfg)
	}

	utils.LogInfo("Evidex acquisition completed", map[string]string{"status": "successfully"})
}

// runProductionMode runs in production mode (no disk writes, direct transmission)
func runProductionMode(cfg *config.Config) {
	// Create acquirer without working directory (memory only)
	acquirer := acquisition.NewAcquirer("")

	// Set hash algorithm
	if cfg.HashAlgorithm != "" {
		acquirer.SetHashAlgorithm(cfg.HashAlgorithm)
		utils.LogInfo("Hash algorithm set", map[string]string{"algorithm": cfg.HashAlgorithm})
	}

	// Process input files (metadata only, no copying)
	processFiles(acquirer, cfg)

	// Build evidence package (in memory)
	utils.LogInfo("Building evidence package", map[string]string{})
	pkg := acquirer.GetEvidencePackage()

	// Set manifest properties
	if cfg.CaseID != "" {
		pkg.Manifest.CaseID = cfg.CaseID
	}
	if cfg.Notes != "" {
		pkg.Manifest.ExaminersNotes = cfg.Notes
	}

	// Print summary
	fmt.Printf("\n")
	fmt.Printf("================================\n")
	fmt.Printf("FORENSIC EVIDENCE PACKAGE\n")
	fmt.Printf("================================\n")
	fmt.Printf("Evidence ID: %s\n", pkg.Manifest.ID)
	fmt.Printf("Case ID: %s\n", pkg.Manifest.CaseID)
	fmt.Printf("Files Acquired: %d\n", pkg.Manifest.FileCount)
	fmt.Printf("Total Size: %d bytes\n", pkg.Manifest.TotalSize)
	fmt.Printf("Mode: Production (no local storage)\n")
	fmt.Printf("Hash Algorithm: %s\n", pkg.Manifest.HashAlgorithm)
	fmt.Printf("Acquisition Status: %s\n", pkg.Manifest.Integrity)

	// Send to remote server
	fmt.Printf("\n")
	fmt.Printf("================================\n")
	fmt.Printf("SENDING EVIDENCE TO SERVER\n")
	fmt.Printf("================================\n")
	fmt.Printf("Server URL: %s\n", cfg.ServerURL)
	utils.LogInfo("Sending evidence package to server", map[string]string{"url": cfg.ServerURL})

	snd := sender.NewSender(cfg.ServerURL, cfg.AuthToken)
	if err := snd.SendEvidencePackage(pkg); err != nil {
		utils.LogError("Failed to send evidence package", map[string]string{"error": err.Error()})
		fmt.Printf("\n❌ Failed to send: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\n✅ Evidence package sent successfully!\n")
	utils.LogInfo("Evidence package sent successfully", map[string]string{"status": "success"})
}

// runDebugMode runs in debug mode (local storage)
func runDebugMode(cfg *config.Config) {
	// Create working directory
	if err := utils.EnsureDirectory(cfg.OutputDir); err != nil {
		utils.LogError("Failed to create output directory", map[string]string{"path": cfg.OutputDir, "error": err.Error()})
		os.Exit(1)
	}

	// Create acquirer
	acquirer := createAcquirer(cfg.OutputDir, cfg)

	// Process input files
	processFiles(acquirer, cfg)

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
	if cfg.CaseID != "" {
		pkg.Manifest.CaseID = cfg.CaseID
	}
	if cfg.Notes != "" {
		pkg.Manifest.ExaminersNotes = cfg.Notes
	}

	// Export formats
	exportFormats(cfg.OutputDir, pkg, cfg)

	// Print summary
	fmt.Printf("\n")
	fmt.Printf("================================\n")
	fmt.Printf("FORENSIC EVIDENCE PACKAGE CREATED\n")
	fmt.Printf("================================\n")
	fmt.Printf("Evidence ID: %s\n", pkg.Manifest.ID)
	fmt.Printf("Files Acquired: %d\n", pkg.Manifest.FileCount)
	fmt.Printf("Total Size: %d bytes\n", pkg.Manifest.TotalSize)
	fmt.Printf("Output Directory: %s\n", cfg.OutputDir)
	fmt.Printf("Mode: Debug (local storage)\n")
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
}

// createAcquirer creates and configures the acquirer
func createAcquirer(workingDir string, cfg *config.Config) *acquisition.Acquirer {
	acquirer := acquisition.NewAcquirer(workingDir)

	// Set hash algorithm
	if cfg.HashAlgorithm != "" {
		acquirer.SetHashAlgorithm(cfg.HashAlgorithm)
		utils.LogInfo("Hash algorithm set", map[string]string{"algorithm": cfg.HashAlgorithm})
	}

	return acquirer
}

// processFiles processes all input files and directories
func processFiles(acquirer *acquisition.Acquirer, cfg *config.Config) {
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
}

// exportFormats exports all requested formats
func exportFormats(workingDir string, pkg *models.EvidencePackage, cfg *config.Config) {
	pkgFormatter := formatter.NewPackageFormatter(workingDir, pkg)

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
}
