package tests

import (
	"os"
	"testing"

	"github.com/ilexum-group/evidex/internal/config"
)

// Testconfig.ParseFlags tests the flag parsing functionality.
func TestParseFlags(t *testing.T) {
	// Save original args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// Set test args
	os.Args = []string{"evidex", "-server", "https://test.com/api/evidence", "-token", "test-token", "/path/to/file"}

	cfg := config.ParseFlags()

	if cfg.ServerURL == "" {
		t.Error("Expected server URL to be set")
	}

	if len(cfg.FilePaths) == 0 {
		t.Error("Expected at least one path to be parsed")
	}
}

// TestFlagDefaults tests that flag defaults are set correctly.
func TestFlagDefaults(t *testing.T) {
	// Save original args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"evidex", "-server", "https://test.com", "-token", "test-token", "/path"}

	cfg := config.ParseFlags()

	if cfg.Recursive != false {
		t.Error("Expected Recursive to default to false")
	}
}

// Testconfig.ValidateConfig tests input validation.
func TestValidateConfig(t *testing.T) {
	// Create temporary test file
	tmpfile, err := os.CreateTemp("", "test-*.txt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer func() {
		if err := os.Remove(tmpfile.Name()); err != nil {
			t.Logf("Failed to remove temp file: %v", err)
		}
	}()
	if err := tmpfile.Close(); err != nil {
		t.Logf("Failed to close temp file: %v", err)
	}

	tests := []struct {
		name    string
		config  *config.Config
		wantErr bool
	}{
		{
			name: "Valid minimal config",
			config: &config.Config{
				FilePaths: []string{tmpfile.Name()},
				ServerURL: "http://test.com",
			},
			wantErr: false,
		},
		{
			name: "Missing server URL",
			config: &config.Config{
				FilePaths: []string{tmpfile.Name()},
			},
			wantErr: true,
		},
		{
			name: "No paths provided",
			config: &config.Config{
				FilePaths: []string{},
				ServerURL: "http://test.com",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := config.ValidateConfig(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("config.ValidateConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
