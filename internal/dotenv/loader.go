package dotenv

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// Load loads .env files from the specified directory.
// Loads .env first, then .env.local (which overrides .env values).
// If neither file exists, no error is returned.
func Load(dir string) error {
	// Load .env
	envPath := filepath.Join(dir, ".env")
	if _, err := os.Stat(envPath); !errors.Is(err, os.ErrNotExist) {
		if err := godotenv.Load(envPath); err != nil {
			return err
		}
	}
	
	// Load .env.local (overrides .env)
	localPath := filepath.Join(dir, ".env.local")
	if _, err := os.Stat(localPath); !errors.Is(err, os.ErrNotExist) {
		// Use Overload to override existing vars from .env
		if err := godotenv.Overload(localPath); err != nil {
			return err
		}
	}
	
	return nil
}
