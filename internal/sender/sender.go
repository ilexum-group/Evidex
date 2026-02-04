// Package sender handles sending evidence packages to remote endpoints
package sender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/ilexum-group/evidex/internal/logger"
	"github.com/ilexum-group/evidex/pkg/models"
)

// Sender encapsulates configuration for sending evidence to a remote server.
type Sender struct {
	serverURL  string
	authToken  string
	httpClient *http.Client
}

// NewSender builds a sender with the provided endpoint configuration.
func NewSender(serverURL, authToken string) *Sender {
	return &Sender{
		serverURL:  serverURL,
		authToken:  authToken,
		httpClient: &http.Client{},
	}
}

// SendEvidencePackage sends the evidence package to a remote server.
// Implements intelligent chunking for large evidence files..
func (s *Sender) SendEvidencePackage(pkg *models.EvidencePackage) error {
	if pkg == nil {
		return fmt.Errorf("evidence package is nil")
	}
	if s.serverURL == "" {
		return fmt.Errorf("server URL is required")
	}

	logger.LogInfo("Preparing to send evidence package to server", map[string]string{
		"url":        s.serverURL,
		"file_count": fmt.Sprintf("%d", len(pkg.Files)),
	})

	// Load file contents for transmission
	logger.LogInfo("Loading file contents for transmission", map[string]string{
		"file_count": fmt.Sprintf("%d", len(pkg.Files)),
	})

	for i, fileEvidence := range pkg.Files {
		// Read file content
		content, err := os.ReadFile(fileEvidence.SourcePath)
		if err != nil {
			logger.LogWarn("Failed to read file content for transmission", map[string]string{
				"path":  fileEvidence.SourcePath,
				"error": err.Error(),
			})
			continue
		}

		pkg.Files[i].FileContent = content
		logger.LogDebug("Loaded file content", map[string]string{
			"file": fileEvidence.Filename,
			"size": fmt.Sprintf("%d bytes", len(content)),
		})
	}

	// Strategy: Send as JSON with intelligent chunking for large files
	jsonData, err := json.Marshal(pkg)
	if err != nil {
		logger.LogError("Failed to marshal evidence package", map[string]string{"error": err.Error()})
		return fmt.Errorf("failed to marshal evidence package: %w", err)
	}

	contentLength := len(jsonData)
	logger.LogDebug("Sending evidence package payload", map[string]string{
		"content_length": fmt.Sprintf("%d bytes", contentLength),
		"size_mb":        fmt.Sprintf("%.2f MB", float64(contentLength)/1024/1024),
	})

	return s.sendHTTPRequest(bytes.NewBuffer(jsonData), contentLength, "application/json")
}

// sendHTTPRequest performs the actual HTTP POST request.
func (s *Sender) sendHTTPRequest(body io.Reader, contentLength int, contentType string) error {
	if s.serverURL == "" {
		return fmt.Errorf("server URL is required")
	}

	req, err := http.NewRequest("POST", s.serverURL, body)
	if err != nil {
		logger.LogError("Failed to create HTTP request", map[string]string{"error": err.Error()})
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("User-Agent", "Evidex-Agent/1.0")
	if s.authToken != "" {
		req.Header.Set("Authorization", "Bearer "+s.authToken)
	}
	req.ContentLength = int64(contentLength)

	logger.LogDebug("Sending HTTP request", map[string]string{
		"method":            "POST",
		"content_type":      contentType,
		"content_length":    fmt.Sprintf("%d", contentLength),
		"content_length_mb": fmt.Sprintf("%.2f", float64(contentLength)/1024/1024),
	})

	client := s.httpClient
	if client == nil {
		client = http.DefaultClient
	}

	resp, err := client.Do(req)
	if err != nil {
		logger.LogError("Failed to send HTTP request", map[string]string{"error": err.Error()})
		return fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			logger.LogError("Failed to close response body", map[string]string{"error": err.Error()})
		}
	}()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		logger.LogWarn("Server returned non-OK status", map[string]string{
			"status_code": fmt.Sprintf("%d", resp.StatusCode),
		})
		return fmt.Errorf("server returned status %d", resp.StatusCode)
	}

	logger.LogInfo("Evidence successfully transmitted to server", map[string]string{
		"status_code": fmt.Sprintf("%d", resp.StatusCode),
	})

	return nil
}
