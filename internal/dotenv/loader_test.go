package dotenv

import (
	"testing"
)

func TestLoad_FileNotFound_NoError(t *testing.T) {
	tmpDir := t.TempDir()
	err := Load(tmpDir)
	if err != nil {
		t.Errorf("expected no error when .env not found, got %v", err)
	}
}
