package tests

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/ilexum-group/evidex/internal/models"
	"github.com/ilexum-group/evidex/internal/sender"
)

// createTestFile creates a temporary test file with specified size.
func createTestFile(t *testing.T, size int) string {
	tmpfile, err := os.CreateTemp("", "test-*.bin")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	// Write test data
	data := make([]byte, size)
	for i := range data {
		data[i] = byte(i % 256)
	}

	if _, err := tmpfile.Write(data); err != nil {
		t.Fatalf("Failed to write test data: %v", err)
	}

	if err := tmpfile.Close(); err != nil {
		t.Logf("Failed to close temp file: %v", err)
	}

	t.Cleanup(func() {
		_ = os.Remove(tmpfile.Name())
	})

	return tmpfile.Name()
}

// TestSendEvidencePackage tests sending evidence package to endpoint.
func TestSendEvidencePackage(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		// Verify authorization header
		auth := r.Header.Get("Authorization")
		if auth == "" {
			t.Error("Missing Authorization header")
		}

		// Verify content type
		ct := r.Header.Get("Content-Type")
		if ct != "application/json" {
			t.Errorf("Expected application/json, got %s", ct)
		}

		// Verify body contains valid JSON
		var pkg models.EvidencePackage
		if err := json.NewDecoder(r.Body).Decode(&pkg); err != nil {
			t.Logf("Failed to decode package: %v", err)
		}

		// Return success
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, `{"status":"success"}`)
	}))
	defer server.Close()

	// Create test package
	pkg := &models.EvidencePackage{
		Logs:      []string{"Log entry 1", "Log entry 2"},
		ServerURL: server.URL,
		Version:   "1.0",
	}

	// Send package
	snd := sender.NewSender(server.URL, "test-token")
	err := snd.SendEvidencePackage(pkg)
	if err != nil {
		t.Fatalf("SendEvidencePackage() error = %v", err)
	}
}

// TestSendEvidenceFile tests sending evidence file to endpoint.
func TestSendEvidenceFile(t *testing.T) {
	// Create test file
	testFilePath := createTestFile(t, 1024)

	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		// Verify authorization header
		if r.Header.Get("Authorization") == "" {
			t.Error("Missing authorization header")
		}

		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, `{"status":"success"}`)
	}))
	defer server.Close()

	// Send file
	meta := map[string]string{"filename": filepath.Base(testFilePath)}
	snd := sender.NewSender(server.URL, "test-token")
	err := snd.SendEvidenceFile(testFilePath, meta)
	if err != nil {
		t.Fatalf("SendEvidenceFile() error = %v", err)
	}
}

// TestSendEvidencePackageInvalidEndpoint tests sending to invalid endpoint.
func TestSendEvidencePackageInvalidEndpoint(t *testing.T) {
	pkg := &models.EvidencePackage{
		Version: "1.0",
	}

	// Try to send to invalid URL
	snd := sender.NewSender("http://invalid-url-that-does-not-exist:9999", "token")
	err := snd.SendEvidencePackage(pkg)
	if err == nil {
		t.Error("Expected error for invalid endpoint")
	}
}

// TestSendEvidenceFileNotFound tests sending non-existent file.
func TestSendEvidenceFileNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	meta := map[string]string{}
	snd := sender.NewSender(server.URL, "token")
	err := snd.SendEvidenceFile("/non/existent/file.txt", meta)
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

// TestSendEvidencePackageServerError tests handling server error response.
func TestSendEvidencePackageServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(w, `{"error":"Internal server error"}`)
	}))
	defer server.Close()

	pkg := &models.EvidencePackage{
		Version: "1.0",
	}

	snd := sender.NewSender(server.URL, "token")
	err := snd.SendEvidencePackage(pkg)
	if err == nil {
		t.Error("Expected error for server error response")
	}
}

// TestSendEvidencePackageWithLogs tests sending package with logs.
func TestSendEvidencePackageWithLogs(t *testing.T) {
	// Create mock server
	var receivedLogs []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var pkg models.EvidencePackage
		_ = json.NewDecoder(r.Body).Decode(&pkg)
		receivedLogs = pkg.Logs

		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, `{"status":"success"}`)
	}))
	defer server.Close()

	// Create test package with logs
	logs := []string{
		"<134>1 2026-01-06T10:00:00Z hostname evidex 1234 - [meta@1] Log entry 1",
		"<134>1 2026-01-06T10:00:01Z hostname evidex 1234 - [meta@1] Log entry 2",
	}
	pkg := &models.EvidencePackage{
		Logs:    logs,
		Version: "1.0",
	}

	snd := sender.NewSender(server.URL, "test-token")
	err := snd.SendEvidencePackage(pkg)
	if err != nil {
		t.Fatalf("SendEvidencePackage() error = %v", err)
	}

	if len(receivedLogs) != 2 {
		t.Errorf("Expected 2 logs received, got %d", len(receivedLogs))
	}
}

// TestSendEvidenceFileWithMetadata tests sending file with metadata.
func TestSendEvidenceFileWithMetadata(t *testing.T) {
	testFilePath := createTestFile(t, 512)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, `{"status":"success"}`)
	}))
	defer server.Close()

	meta := map[string]string{
		"file_size": "512",
		"file_hash": "abc123",
	}

	snd := sender.NewSender(server.URL, "token")
	err := snd.SendEvidenceFile(testFilePath, meta)
	if err != nil {
		t.Fatalf("SendEvidenceFile() with metadata error = %v", err)
	}
}

// TestSendHTTPRequest tests basic HTTP request functionality.
func TestSendHTTPRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify headers
		if r.Header.Get("Content-Type") != "application/json" {
			t.Error("Invalid Content-Type header")
		}

		if r.Header.Get("User-Agent") != "Evidex-Agent/1.0" {
			t.Error("Invalid User-Agent header")
		}

		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, `{"result":"ok"}`)
	}))
	defer server.Close()

	pkg := &models.EvidencePackage{
		Version: "1.0",
	}
	snd := sender.NewSender(server.URL, "token")
	err := snd.SendEvidencePackage(pkg)
	if err != nil {
		t.Fatalf("Error sending request: %v", err)
	}
}

// TestSendEvidencePackageWithComplexPayload tests sending complex package.
func TestSendEvidencePackageWithComplexPayload(t *testing.T) {
	receivedPackage := make(chan *models.EvidencePackage, 1)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var pkg models.EvidencePackage
		if err := json.NewDecoder(r.Body).Decode(&pkg); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		receivedPackage <- &pkg
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, `{"status":"success"}`)
	}))
	defer server.Close()

	// Create complex package
	pkg := &models.EvidencePackage{
		Version: "1.0",
		Logs: []string{
			"Log 1",
			"Log 2",
			"Log 3",
		},
	}

	snd := sender.NewSender(server.URL, "test-token")
	err := snd.SendEvidencePackage(pkg)
	if err != nil {
		t.Fatalf("SendEvidencePackage() error = %v", err)
	}

	// Verify received package
	select {
	case received := <-receivedPackage:
		if received.Version != "1.0" {
			t.Errorf("Expected version 1.0, got %s", received.Version)
		}
		if len(received.Logs) != 3 {
			t.Errorf("Expected 3 logs, got %d", len(received.Logs))
		}
	case <-make(chan struct{}):
		t.Error("Failed to receive package")
	}
}

// TestSendEvidencePackageResponseValidation tests response validation.
func TestSendEvidencePackageResponseValidation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Return success only on correct request
		if r.Method == "POST" && r.Header.Get("Authorization") != "" {
			w.WriteHeader(http.StatusOK)
			_, _ = io.WriteString(w, `{"status":"received"}`)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	}))
	defer server.Close()

	pkg := &models.EvidencePackage{Version: "1.0"}

	// This should succeed
	snd := sender.NewSender(server.URL, "token")
	err := snd.SendEvidencePackage(pkg)
	if err != nil {
		t.Fatalf("Expected successful response, got error: %v", err)
	}
}

// TestChunkSizeConstant removed - cannot access private constants from external package.

// TestSendLargeFile tests sending a large file.
func TestSendLargeFile(t *testing.T) {
	// Create a smaller test file (actual chunking would be >64MB)
	testFilePath := createTestFile(t, 100*1024) // 100KB

	requestCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++

		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		// Verify headers
		if r.Header.Get("Authorization") == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	meta := map[string]string{"chunk_test": "true"}
	snd := sender.NewSender(server.URL, "token")
	err := snd.SendEvidenceFile(testFilePath, meta)
	if err != nil {
		t.Fatalf("SendEvidenceFile() error = %v", err)
	}

	if requestCount == 0 {
		t.Error("Expected at least one HTTP request")
	}
}

// TestSendEvidenceFileResponseParsing tests response parsing.
func TestSendEvidenceFileResponseParsing(t *testing.T) {
	testFilePath := createTestFile(t, 256)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		response := map[string]interface{}{
			"status":      "success",
			"file_id":     "12345",
			"size_bytes":  256,
			"received_at": "2026-01-06T10:00:00Z",
		}

		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	meta := map[string]string{}
	snd := sender.NewSender(server.URL, "token")
	err := snd.SendEvidenceFile(testFilePath, meta)
	if err != nil {
		t.Fatalf("SendEvidenceFile() error = %v", err)
	}
}

// BenchmarkSendEvidencePackage benchmarks package sending.
func BenchmarkSendEvidencePackage(b *testing.B) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, `{"status":"ok"}`)
	}))
	defer server.Close()

	pkg := &models.EvidencePackage{
		Version: "1.0",
		Logs:    []string{"Log 1", "Log 2", "Log 3"},
	}

	snd := sender.NewSender(server.URL, "token")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = snd.SendEvidencePackage(pkg)
	}
}
