package logger

import (
	"strings"
	"testing"
)

// TestNewLogger tests logger creation
func TestNewLogger(t *testing.T) {
	appName := "test-app"
	logger, err := NewLogger(appName)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	if logger == nil {
		t.Fatal("Expected logger to be non-nil")
	}

	if logger.appName != appName {
		t.Errorf("Expected app name %s, got %s", appName, logger.appName)
	}

	if logger.hostname == "" {
		t.Error("Expected hostname to be non-empty")
	}
}

// TestLogInfo tests info level logging
func TestLogInfo(t *testing.T) {
	logger, err := NewLogger("test-app")
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	logger.ClearLogs()

	meta := map[string]string{"key": "value"}
	logger.LogInfo("Test info message", meta)

	logs := logger.GetLogs()
	if len(logs) != 1 {
		t.Fatalf("Expected 1 log entry, got %d", len(logs))
	}

	if !strings.Contains(logs[0], "Test info message") {
		t.Errorf("Expected message to be in log: %s", logs[0])
	}
}

// TestLogWarn tests warning level logging
func TestLogWarn(t *testing.T) {
	logger, err := NewLogger("test-app")
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	logger.ClearLogs()

	meta := map[string]string{"severity": "high"}
	logger.LogWarn("Test warning message", meta)

	logs := logger.GetLogs()
	if len(logs) != 1 {
		t.Fatalf("Expected 1 log entry, got %d", len(logs))
	}

	if !strings.Contains(logs[0], "Test warning message") {
		t.Errorf("Expected message to be in log: %s", logs[0])
	}
}

// TestLogError tests error level logging
func TestLogError(t *testing.T) {
	logger, err := NewLogger("test-app")
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	logger.ClearLogs()

	meta := map[string]string{"error_code": "500"}
	logger.LogError("Test error message", meta)

	logs := logger.GetLogs()
	if len(logs) != 1 {
		t.Fatalf("Expected 1 log entry, got %d", len(logs))
	}

	if !strings.Contains(logs[0], "Test error message") {
		t.Errorf("Expected message to be in log: %s", logs[0])
	}
}

// TestLogDebug tests debug level logging
func TestLogDebug(t *testing.T) {
	logger, err := NewLogger("test-app")
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	logger.ClearLogs()

	meta := map[string]string{"component": "acquisition"}
	logger.LogDebug("Test debug message", meta)

	logs := logger.GetLogs()
	if len(logs) != 1 {
		t.Fatalf("Expected 1 log entry, got %d", len(logs))
	}

	if !strings.Contains(logs[0], "Test debug message") {
		t.Errorf("Expected message to be in log: %s", logs[0])
	}
}

// TestRFC5424Format tests RFC 5424 syslog format
func TestRFC5424Format(t *testing.T) {
	logger, err := NewLogger("evidex")
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	logger.ClearLogs()

	logger.LogInfo("Test RFC 5424 format", map[string]string{"field": "value"})

	logs := logger.GetLogs()
	if len(logs) == 0 {
		t.Fatal("Expected at least one log entry")
	}

	logEntry := logs[0]

	// RFC 5424 starts with <PRI>VERSION
	if !strings.HasPrefix(logEntry, "<") {
		t.Errorf("Expected log to start with <, got: %s", logEntry)
	}

	// Check for version
	if !strings.Contains(logEntry, ">1 ") {
		t.Errorf("Expected RFC 5424 version 1, got: %s", logEntry)
	}

	// Check for hostname
	if !strings.Contains(logEntry, logger.hostname) {
		t.Errorf("Expected hostname in log, got: %s", logEntry)
	}

	// Check for app name
	if !strings.Contains(logEntry, "evidex") {
		t.Errorf("Expected app name 'evidex' in log, got: %s", logEntry)
	}
}

// TestGetLogs tests log retrieval
func TestGetLogs(t *testing.T) {
	logger, err := NewLogger("test-app")
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	logger.ClearLogs()

	// Log multiple messages
	logger.LogInfo("Message 1", map[string]string{})
	logger.LogInfo("Message 2", map[string]string{})
	logger.LogInfo("Message 3", map[string]string{})

	logs := logger.GetLogs()
	if len(logs) != 3 {
		t.Errorf("Expected 3 logs, got %d", len(logs))
	}

	for i, log := range logs {
		if log == "" {
			t.Errorf("Log entry %d is empty", i)
		}
	}
}

// TestClearLogs tests log clearing
func TestClearLogs(t *testing.T) {
	logger, err := NewLogger("test-app")
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	logger.LogInfo("Message 1", map[string]string{})
	logger.LogInfo("Message 2", map[string]string{})

	logs := logger.GetLogs()
	if len(logs) != 2 {
		t.Errorf("Expected 2 logs before clear, got %d", len(logs))
	}

	logger.ClearLogs()

	logs = logger.GetLogs()
	if len(logs) != 0 {
		t.Errorf("Expected 0 logs after clear, got %d", len(logs))
	}
}

// TestConcurrentLogging tests thread-safe logging
func TestConcurrentLogging(t *testing.T) {
	logger, err := NewLogger("concurrent-test")
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	logger.ClearLogs()

	// Log from multiple goroutines
	done := make(chan bool, 5)
	for i := 0; i < 5; i++ {
		go func(n int) {
			meta := map[string]string{"goroutine": string(rune(n))}
			logger.LogInfo("Concurrent message", meta)
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 5; i++ {
		<-done
	}

	logs := logger.GetLogs()
	if len(logs) != 5 {
		t.Errorf("Expected 5 concurrent logs, got %d", len(logs))
	}
}

// TestEmptyMetadata tests logging with empty metadata
func TestEmptyMetadata(t *testing.T) {
	logger, err := NewLogger("test-app")
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	logger.ClearLogs()

	logger.LogInfo("Message with no metadata", map[string]string{})

	logs := logger.GetLogs()
	if len(logs) != 1 {
		t.Fatalf("Expected 1 log entry, got %d", len(logs))
	}

	// Should have "-" for empty structured data
	if !strings.Contains(logs[0], " - ") {
		t.Errorf("Expected empty structured data marker, got: %s", logs[0])
	}
}
