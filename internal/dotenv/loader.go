package dotenv

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// Load loads .env files from the specified directory.
// If .env does not exist, no error is returned.
func Load(dir string) error {
	envPath := filepath.Join(dir, ".env")
	
	// Check if file exists
	if _, err := os.Stat(envPath); errors.Is(err, os.ErrNotExist) {
		return nil // File not found is OK
	}
	
	// Load the .env file
	return godotenv.Load(envPath)
}
