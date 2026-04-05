package dotenv

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
)

// Load loads .env files from the specified directory.
// Loads .env first, then .env.local (which overrides .env values).
// If neither file exists, no error is returned.
// Files with world-readable permissions are skipped with a warning.
func Load(dir string) error {
	// Load .env
	envPath := filepath.Join(dir, ".env")
	if err := loadFile(envPath, false); err != nil {
		return err
	}
	
	// Load .env.local (overrides .env)
	localPath := filepath.Join(dir, ".env.local")
	if err := loadFile(localPath, true); err != nil {
		return err
	}
	
	return nil
}

// LoadFile loads a single environment file (for custom --env-file flag).
func LoadFile(path string) error {
	return loadFile(path, false)
}

// loadFile loads a single env file with security checks.
// If override is true, uses Overload instead of Load.
func loadFile(path string, override bool) error {
	info, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil // File not found is OK
	}
	if err != nil {
		return err
	}
	
	// Security check: Skip world-readable or group-readable files (Unix only)
	if runtime.GOOS != "windows" {
		perm := info.Mode().Perm()
		// Check if world-readable (others can read: 0004) or group-readable (0040)
		if perm&0044 != 0 {
			log.Printf("warning: %s is world/group readable (permissions: %o), skipping for security", path, perm)
			return nil
		}
	}
	
	// Load the file
	if override {
		return godotenv.Overload(path)
	}
	return godotenv.Load(path)
}
