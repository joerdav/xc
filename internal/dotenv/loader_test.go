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

func TestLoad_WithLocal_OverridesBase(t *testing.T) {
	// Arrange
	tmpDir := t.TempDir()
	
	// Create .env with base values
	envFile := filepath.Join(tmpDir, ".env")
	if err := os.WriteFile(envFile, []byte("KEY=base\nONLY_BASE=base_value"), 0600); err != nil {
		t.Fatal(err)
	}
	
	// Create .env.local with overrides
	localFile := filepath.Join(tmpDir, ".env.local")
	if err := os.WriteFile(localFile, []byte("KEY=local\nONLY_LOCAL=local_value"), 0600); err != nil {
		t.Fatal(err)
	}
	
	// Act
	err := Load(tmpDir)
	
	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := os.Getenv("KEY"); got != "local" {
		t.Errorf("KEY = %q, want %q (should be overridden)", got, "local")
	}
	if got := os.Getenv("ONLY_BASE"); got != "base_value" {
		t.Errorf("ONLY_BASE = %q, want %q", got, "base_value")
	}
	if got := os.Getenv("ONLY_LOCAL"); got != "local_value" {
		t.Errorf("ONLY_LOCAL = %q, want %q", got, "local_value")
	}
	
	// Cleanup
	t.Cleanup(func() {
		os.Unsetenv("KEY")
		os.Unsetenv("ONLY_BASE")
		os.Unsetenv("ONLY_LOCAL")
	})
}

func TestLoad_WorldReadable_LogsWarningAndSkips(t *testing.T) {
	// This test only runs on Unix systems
	if os.Getenv("GOOS") == "windows" {
		t.Skip("skipping permission test on Windows")
	}
	
	// Arrange
	tmpDir := t.TempDir()
	envFile := filepath.Join(tmpDir, ".env")
	
	// Create world-readable .env file
	if err := os.WriteFile(envFile, []byte("SECRET=exposed"), 0644); err != nil {
		t.Fatal(err)
	}
	
	// Act
	err := Load(tmpDir)
	
	// Assert
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	
	// Verify SECRET was NOT loaded (file was skipped)
	if got := os.Getenv("SECRET"); got != "" {
		t.Errorf("SECRET should not be loaded from world-readable file, got %q", got)
	}
	
	// Cleanup
	t.Cleanup(func() {
		os.Unsetenv("SECRET")
	})
}
