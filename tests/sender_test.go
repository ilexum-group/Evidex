package tests

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ilexum-group/evidex/internal/sender"
	"github.com/ilexum-group/evidex/pkg/models"
)

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
		Files:        []*models.FileEvidence{},
		CustodyChain: &models.CustodyChainEntry{ID: "test-id"},
	}

	// Send package
	snd := sender.NewSender(server.URL, "test-token")
	err := snd.SendEvidencePackage(pkg)
	if err != nil {
		t.Fatalf("SendEvidencePackage() error = %v", err)
	}
}

// TestSendEvidencePackageInvalidEndpoint tests sending to invalid endpoint.
func TestSendEvidencePackageInvalidEndpoint(t *testing.T) {
	pkg := &models.EvidencePackage{}

	// Try to send to invalid URL
	snd := sender.NewSender("http://invalid-url-that-does-not-exist:9999", "token")
	err := snd.SendEvidencePackage(pkg)
	if err == nil {
		t.Error("Expected error for invalid endpoint")
	}
}

// TestSendEvidencePackageServerError tests handling server error response.
func TestSendEvidencePackageServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = io.WriteString(w, `{"error":"Internal server error"}`)
	}))
	defer server.Close()

	pkg := &models.EvidencePackage{}

	snd := sender.NewSender(server.URL, "token")
	err := snd.SendEvidencePackage(pkg)
	if err == nil {
		t.Error("Expected error for server error response")
	}
}

// TestSendEvidencePackageWithLogs tests sending package with logs.
func TestSendEvidencePackageWithLogs(t *testing.T) {
	// Create mock server
	var receivedLogs []models.LogEntry
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var pkg models.EvidencePackage
		_ = json.NewDecoder(r.Body).Decode(&pkg)
		receivedLogs = pkg.CustodyChain.LogEntries

		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, `{"status":"success"}`)
	}))
	defer server.Close()

	// Create test package with logs
	logs := []models.LogEntry{
		{Level: "INFO", Message: "Log entry 1"},
		{Level: "INFO", Message: "Log entry 2"},
	}
	pkg := &models.EvidencePackage{
		Files:        []*models.FileEvidence{},
		CustodyChain: &models.CustodyChainEntry{ID: "test-id", LogEntries: logs},
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

	pkg := &models.EvidencePackage{}
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
		Files: []*models.FileEvidence{},
		CustodyChain: &models.CustodyChainEntry{
			ID:           "test-id",
			AgentVersion: "1.0",
			LogEntries: []models.LogEntry{
				{Level: "INFO", Message: "Log 1"},
				{Level: "INFO", Message: "Log 2"},
				{Level: "INFO", Message: "Log 3"},
			},
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
		if received.CustodyChain.AgentVersion != "1.0" {
			t.Errorf("Expected version 1.0, got %s", received.CustodyChain.AgentVersion)
		}
		if len(received.CustodyChain.LogEntries) != 3 {
			t.Errorf("Expected 3 logs, got %d", len(received.CustodyChain.LogEntries))
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

	pkg := &models.EvidencePackage{
		Files:        []*models.FileEvidence{},
		CustodyChain: &models.CustodyChainEntry{ID: "test-id", AgentVersion: "1.0"},
	}

	// This should succeed
	snd := sender.NewSender(server.URL, "token")
	err := snd.SendEvidencePackage(pkg)
	if err != nil {
		t.Fatalf("Expected successful response, got error: %v", err)
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

		Files:        []*models.FileEvidence{},
		CustodyChain: &models.CustodyChainEntry{ID: "test-id", LogEntries: []models.LogEntry{{Level: "INFO", Message: "Log 1"}, {Level: "INFO", Message: "Log 2"}, {Level: "INFO", Message: "Log 3"}}},
	}

	snd := sender.NewSender(server.URL, "token")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = snd.SendEvidencePackage(pkg)
	}
}
