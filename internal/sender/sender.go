// Package sender handles sending evidence packages to remote endpoints
package sender

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/ilexum/evidex/internal/models"
	"github.com/ilexum/evidex/internal/utils"
)

const (
	// ChunkSize defines the size of each chunk for large files (64 MB)
	ChunkSize = 64 * 1024 * 1024
	// MaxPayloadSize defines the maximum JSON payload size (100 MB)
	MaxPayloadSize = 100 * 1024 * 1024
)

// Sender encapsulates configuration for sending evidence to a remote server
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

// WithHTTPClient overrides the HTTP client (useful for tests and custom transports).
func (s *Sender) WithHTTPClient(client *http.Client) *Sender {
	if client != nil {
		s.httpClient = client
	}
	return s
}

// SendEvidencePackage sends the evidence package to a remote server.
// Implements intelligent chunking for large evidence files.
func (s *Sender) SendEvidencePackage(pkg *models.EvidencePackage) error {
	if pkg == nil {
		return fmt.Errorf("evidence package is nil")
	}
	if s.serverURL == "" {
		return fmt.Errorf("server URL is required")
	}

	utils.LogInfo("Preparing to send evidence package to server", map[string]string{
		"url":        s.serverURL,
		"file_count": fmt.Sprintf("%d", len(pkg.Files)),
	})

	// Set authentication headers
	pkg.ServerURL = s.serverURL
	pkg.AuthToken = s.authToken

	// Load file contents for transmission
	utils.LogInfo("Loading file contents for transmission", map[string]string{
		"file_count": fmt.Sprintf("%d", len(pkg.Files)),
	})

	for i, fileEvidence := range pkg.Files {
		// Read file content
		content, err := os.ReadFile(fileEvidence.SourcePath)
		if err != nil {
			utils.LogWarn("Failed to read file content for transmission", map[string]string{
				"path":  fileEvidence.SourcePath,
				"error": err.Error(),
			})
			continue
		}

		pkg.Files[i].FileContent = content
		utils.LogDebug("Loaded file content", map[string]string{
			"file": fileEvidence.Filename,
			"size": fmt.Sprintf("%d bytes", len(content)),
		})
	}

	// Strategy: Send as JSON with intelligent chunking for large files
	jsonData, err := json.Marshal(pkg)
	if err != nil {
		utils.LogError("Failed to marshal evidence package", map[string]string{"error": err.Error()})
		return fmt.Errorf("failed to marshal evidence package: %w", err)
	}

	contentLength := len(jsonData)
	utils.LogDebug("Sending evidence package payload", map[string]string{
		"content_length": fmt.Sprintf("%d bytes", contentLength),
		"size_mb":        fmt.Sprintf("%.2f MB", float64(contentLength)/1024/1024),
	})

	return s.sendHTTPRequest(bytes.NewBuffer(jsonData), contentLength, "application/json")
}

// SendEvidenceFile sends an evidence file separately (for large files)
func (s *Sender) SendEvidenceFile(filePath string, metadata map[string]string) error {
	if s.serverURL == "" {
		return fmt.Errorf("server URL is required")
	}

	file, err := os.Open(filePath)
	if err != nil {
		utils.LogError("Failed to open evidence file", map[string]string{"error": err.Error()})
		return fmt.Errorf("failed to open evidence file: %w", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			utils.LogError("Failed to close evidence file", map[string]string{"error": err.Error()})
		}
	}()

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %w", err)
	}

	fileSize := fileInfo.Size()
	utils.LogInfo("Starting evidence file transmission", map[string]string{
		"file_path": filePath,
		"file_size": fmt.Sprintf("%d bytes", fileSize),
		"size_mb":   fmt.Sprintf("%.2f MB", float64(fileSize)/1024/1024),
	})

	// For files larger than ChunkSize, chunk them
	if fileSize > ChunkSize {
		return s.sendFileChunked(file, filePath, fileSize, metadata)
	}

	// For smaller files, send directly
	return s.sendHTTPRequest(file, int(fileSize), "application/octet-stream")
}

// sendFileChunked sends a file in chunks for large evidence files
func (s *Sender) sendFileChunked(file *os.File, filePath string, fileSize int64, metadata map[string]string) error {
	totalChunks := (fileSize + ChunkSize - 1) / ChunkSize

	utils.LogInfo("File exceeds chunk size, using chunked transfer", map[string]string{
		"file_path":     filePath,
		"chunk_size_mb": fmt.Sprintf("%.2f MB", float64(ChunkSize)/1024/1024),
		"total_chunks":  fmt.Sprintf("%d", totalChunks),
	})

	chunkBuffer := make([]byte, ChunkSize)
	chunkNum := 0

	for {
		bytesRead, err := file.Read(chunkBuffer)
		if err != nil && err != io.EOF {
			utils.LogError("Failed to read chunk", map[string]string{"error": err.Error()})
			return fmt.Errorf("failed to read chunk: %w", err)
		}

		if bytesRead == 0 {
			break
		}

		chunkNum++
		chunkData := chunkBuffer[:bytesRead]

		// Prepare chunk metadata
		chunkPayload := map[string]interface{}{
			"type":         "evidence_chunk",
			"file_path":    filePath,
			"chunk_num":    chunkNum,
			"total_chunks": totalChunks,
			"chunk_size":   bytesRead,
			"file_size":    fileSize,
			"metadata":     metadata,
			"data":         string(chunkData),
		}

		payloadJSON, err := json.Marshal(chunkPayload)
		if err != nil {
			utils.LogError("Failed to marshal chunk", map[string]string{"error": err.Error()})
			return fmt.Errorf("failed to marshal chunk: %w", err)
		}

		utils.LogDebug("Sending chunk", map[string]string{
			"chunk":      fmt.Sprintf("%d/%d", chunkNum, totalChunks),
			"size_bytes": fmt.Sprintf("%d", bytesRead),
			"size_mb":    fmt.Sprintf("%.2f", float64(bytesRead)/1024/1024),
			"progress":   fmt.Sprintf("%.1f%%", float64(chunkNum*100)/float64(totalChunks)),
		})

		if err := s.sendHTTPRequest(bytes.NewBuffer(payloadJSON), len(payloadJSON), "application/json"); err != nil {
			utils.LogError("Failed to send chunk", map[string]string{
				"chunk": fmt.Sprintf("%d/%d", chunkNum, totalChunks),
				"error": err.Error(),
			})
			return err
		}

		utils.LogInfo("Chunk sent successfully", map[string]string{
			"chunk":    fmt.Sprintf("%d/%d", chunkNum, totalChunks),
			"progress": fmt.Sprintf("%.1f%%", float64(chunkNum*100)/float64(totalChunks)),
		})
	}

	return nil
}

// sendHTTPRequest performs the actual HTTP POST request
func (s *Sender) sendHTTPRequest(body io.Reader, contentLength int, contentType string) error {
	if s.serverURL == "" {
		return fmt.Errorf("server URL is required")
	}

	req, err := http.NewRequest("POST", s.serverURL, body)
	if err != nil {
		utils.LogError("Failed to create HTTP request", map[string]string{"error": err.Error()})
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("User-Agent", "Evidex-Agent/1.0")
	if s.authToken != "" {
		req.Header.Set("Authorization", "Bearer "+s.authToken)
	}
	req.ContentLength = int64(contentLength)

	utils.LogDebug("Sending HTTP request", map[string]string{
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
		utils.LogError("Failed to send HTTP request", map[string]string{"error": err.Error()})
		return fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			utils.LogError("Failed to close response body", map[string]string{"error": err.Error()})
		}
	}()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		utils.LogWarn("Server returned non-OK status", map[string]string{
			"status_code": fmt.Sprintf("%d", resp.StatusCode),
		})
		return fmt.Errorf("server returned status %d", resp.StatusCode)
	}

	utils.LogInfo("Evidence successfully transmitted to server", map[string]string{
		"status_code": fmt.Sprintf("%d", resp.StatusCode),
	})

	return nil
}

// SendEvidencePackage is a convenience wrapper that constructs a Sender on the fly.
// Deprecated: Prefer creating a Sender with NewSender for reuse and configurability.
func SendEvidencePackage(serverURL, authToken string, pkg *models.EvidencePackage) error {
	return NewSender(serverURL, authToken).SendEvidencePackage(pkg)
}

// SendEvidenceFile is a convenience wrapper that constructs a Sender on the fly.
// Deprecated: Prefer creating a Sender with NewSender for reuse and configurability.
func SendEvidenceFile(serverURL, authToken, filePath string, metadata map[string]string) error {
	return NewSender(serverURL, authToken).SendEvidenceFile(filePath, metadata)
}
