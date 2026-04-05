package dotenv

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_FileNotFound_NoError(t *testing.T) {
	tmpDir := t.TempDir()
	err := Load(tmpDir)
	if err != nil {
		t.Errorf("expected no error when .env not found, got %v", err)
	}
}

func TestLoad_ValidEnv_LoadsVariables(t *testing.T) {
	// Arrange
	tmpDir := t.TempDir()
	envFile := filepath.Join(tmpDir, ".env")
	content := "TEST_KEY=test_value\nANOTHER=value2"
	if err := os.WriteFile(envFile, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	
	// Act
	err := Load(tmpDir)
	
	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := os.Getenv("TEST_KEY"); got != "test_value" {
		t.Errorf("TEST_KEY = %q, want %q", got, "test_value")
	}
	if got := os.Getenv("ANOTHER"); got != "value2" {
		t.Errorf("ANOTHER = %q, want %q", got, "value2")
	}
	
	// Cleanup
	t.Cleanup(func() {
		os.Unsetenv("TEST_KEY")
		os.Unsetenv("ANOTHER")
	})
}
