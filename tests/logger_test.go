package tests

import (
	"testing"

	"github.com/ilexum-group/evidex/internal/logger"
)

// TestNewLogger tests logger creation.
func TestNewLogger(t *testing.T) {
	appName := "test-app"
	hostname := "test-host"
	processID := "12345"

	log, err := logger.NewLogger(appName, hostname, processID)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	if log == nil {
		t.Fatal("Expected logger to be non-nil")
	}
}

// TestLogInfo tests info level logging.
func TestLogInfo(t *testing.T) {
	log, err := logger.NewLogger("test-app", "test-host", "12345")
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	meta := map[string]string{"key": "value"}
	// Just verify it doesn't panic - logger writes to stdout
	log.LogInfo("Test info message", meta)

	// Test passes if no panic occurred
}

// TestLogWarn tests warning level logging.
func TestLogWarn(t *testing.T) {
	log, err := logger.NewLogger("test-app", "test-host", "12345")
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	meta := map[string]string{"severity": "high"}
	// Just verify it doesn't panic - logger writes to stdout
	log.LogWarn("Test warning message", meta)

	// Test passes if no panic occurred
}

// TestLogError tests error level logging.
func TestLogError(t *testing.T) {
	log, err := logger.NewLogger("test-app", "test-host", "12345")
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	meta := map[string]string{"error_code": "500"}
	// Just verify it doesn't panic - logger writes to stdout
	log.LogError("Test error message", meta)

	// Test passes if no panic occurred
}

// TestLogDebug tests debug level logging.
func TestLogDebug(t *testing.T) {
	log, err := logger.NewLogger("test-app", "test-host", "12345")
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	meta := map[string]string{"component": "acquisition"}
	// Just verify it doesn't panic - logger writes to stdout
	log.LogDebug("Test debug message", meta)

	// Test passes if no panic occurred
}

// TestInitDefaultLogger tests global logger initialization.
func TestInitDefaultLogger(t *testing.T) {
	err := logger.InitDefaultLogger("test-app", "test-host", "12345")
	if err != nil {
		t.Fatalf("Failed to initialize default logger: %v", err)
	}

	// Test global convenience functions
	logger.LogInfo("Global logger test", map[string]string{"test": "value"})
	logger.LogWarn("Global warning test", map[string]string{})
	logger.LogError("Global error test", map[string]string{})
	logger.LogDebug("Global debug test", map[string]string{})

	// Test passes if no panics occurred
}

// TestEmptyMetadata tests logging with empty metadata.
func TestEmptyMetadata(t *testing.T) {
	log, err := logger.NewLogger("test-app", "test-host", "12345")
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// Should handle empty metadata gracefully
	log.LogInfo("Message with no metadata", map[string]string{})

	// Test passes if no panic occurred
}
