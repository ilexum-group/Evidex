// Package main implements the Evidex CLI for forensic evidence acquisition.
package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/ilexum-group/evidex/internal/acquisition"
	"github.com/ilexum-group/evidex/internal/config"
	"github.com/ilexum-group/evidex/internal/logger"
	"github.com/ilexum-group/evidex/internal/metadata"
	osWrapper "github.com/ilexum-group/evidex/internal/os"
	"github.com/ilexum-group/evidex/internal/sender"
	"github.com/ilexum-group/evidex/pkg/models"
)

const (
	applicationName = "Evidex"
)

// version is set at build time via -ldflags.
var version = "placeholder"

func main() {
	// Parse configuration first
	cfg := config.ParseFlags()

	// Validate required flags before initialization
	if err := config.ValidateConfig(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Configuration validation failed: %v\n", err)
		fmt.Fprintln(os.Stderr, "Use --help for usage information")
		os.Exit(1)
	}

	// Create custody chain entry
	custodyChain, err := models.NewCustodyChainEntry(applicationName, version)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to create custody chain", map[string]string{"error": err.Error()})
		os.Exit(1)
	}

	// Initialize OS wrapper
	osImpl := osWrapper.New()

	// Configure OS wrapper with custody chain logger
	osImpl.SetLogger(custodyChain.LogCommand)

	// Get system information for initialization
	hostname, err := osImpl.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	currentUser, err := osImpl.GetCurrentUser()
	if err != nil {
		currentUser = "unknown"
	}

	processID := osImpl.GetProcessID()

	// Set hostname and user in custody chain
	custodyChain.SetAgentHostname(hostname)
	custodyChain.SetAgentUser(currentUser)

	// Initialize logger module (prioritize this over utils)
	if err := logger.InitDefaultLogger(applicationName, hostname, strconv.Itoa(processID)); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	logger.LogInfo("Starting Evidex", map[string]string{
		"version":  version,
		"hostname": hostname,
		"user":     currentUser,
		"pid":      strconv.Itoa(processID),
	})

	// Initialize metadata manager with custody chain logger
	metadataMgr := metadata.NewMetadataManager(custodyChain.LogCommand)

	// Create acquirer with dependency injection
	acquirer := acquisition.NewAcquirer(custodyChain, osImpl, metadataMgr)

	// Process input files (metadata only)
	processFiles(acquirer, cfg)

	// Build evidence package (in memory)
	logger.LogInfo("Building evidence package", map[string]string{})
	pkg := acquirer.GetEvidencePackage()

	// Set Case ID if provided
	pkg.CaseID = cfg.CaseID

	// Log summary
	logger.LogInfo("Forensic evidence package built", map[string]string{
		"evidence_id":    pkg.CustodyChain.ID,
		"files_acquired": fmt.Sprintf("%d", pkg.CustodyChain.ItemCount),
		"total_size":     fmt.Sprintf("%d bytes", pkg.CustodyChain.TotalSizeBytes),
		"hash_algorithm": "MD5, SHA1, SHA256, SHA512",
	})

	// Send to remote server
	logger.LogInfo("Sending evidence to server", map[string]string{"server_url": cfg.ServerURL})

	snd := sender.NewSender(cfg.ServerURL, cfg.AuthToken)
	if err := snd.SendEvidencePackage(pkg); err != nil {
		logger.LogError("Failed to send evidence package", map[string]string{"error": err.Error()})
		os.Exit(1)
	}

	logger.LogInfo("Evidence package sent successfully", map[string]string{"status": "success"})
	logger.LogInfo("Evidex acquisition completed", map[string]string{"status": "successfully"})
}

// processFiles processes all input files and directories.
func processFiles(acquirer *acquisition.Acquirer, cfg *config.Config) {
	for _, filePath := range cfg.FilePaths {
		fileInfo, err := os.Stat(filePath)
		if err != nil {
			logger.LogError("File not found", map[string]string{"path": filePath, "error": err.Error()})
			continue
		}

		if fileInfo.IsDir() {
			logger.LogInfo("Processing directory", map[string]string{"path": filePath, "recursive": fmt.Sprintf("%v", cfg.Recursive)})
			if err := acquirer.AcquireDirectory(filePath, cfg.Recursive); err != nil {
				logger.LogError("Directory acquisition failed", map[string]string{"path": filePath, "error": err.Error()})
			}
		} else {
			logger.LogInfo("Processing file", map[string]string{"path": filePath})
			if err := acquirer.AcquireFile(filePath); err != nil {
				logger.LogError("File acquisition failed", map[string]string{"path": filePath, "error": err.Error()})
			}
		}
	}
}
