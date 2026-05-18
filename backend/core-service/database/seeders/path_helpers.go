package seeders

import (
	"os"
	"path/filepath"
	"runtime"

	"core-service/pkg/config"
)

// getBasePath returns the base path for locating resources
// Priority: config.App.BasePath > executable directory > current working directory
func getBasePath() string {
	cfg := config.Get()
	if cfg != nil && cfg.App.BasePath != "" {
		return cfg.App.BasePath
	}

	// Try to get executable directory
	exe, err := os.Executable()
	if err == nil {
		return filepath.Dir(exe)
	}

	// Fallback to working directory
	wd, err := os.Getwd()
	if err == nil {
		return wd
	}

	// Last resort: return current directory
	return "."
}

// getProjectRoot returns the project root directory (where go.mod is located)
// Useful for development when files are relative to project root
func getProjectRoot() string {
	// Try to get from current working directory
	if wd, err := os.Getwd(); err == nil {
		// Check if go.mod exists in current directory or parent
		if _, err := os.Stat(filepath.Join(wd, "go.mod")); err == nil {
			return wd
		}
		// Check parent directories (up to 3 levels)
		for i := 0; i < 3; i++ {
			parent := filepath.Dir(wd)
			if parent == wd {
				break
			}
			if _, err := os.Stat(filepath.Join(parent, "go.mod")); err == nil {
				return parent
			}
			wd = parent
		}
	}

	// Try executable path
	if exe, err := os.Executable(); err == nil {
		exeDir := filepath.Dir(exe)
		if _, err := os.Stat(filepath.Join(exeDir, "go.mod")); err == nil {
			return exeDir
		}
	}

	return "."
}

// fileExists checks if a file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// getCallerDir returns the directory of the caller function
// Useful for locating files relative to the source file
func getCallerDir() string {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		return "."
	}
	return filepath.Dir(file)
}