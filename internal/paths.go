package internal

import (
	"os"
	"path/filepath"
	"runtime"
)

// GetDefaultTextsDir returns the platform-appropriate default directory for texts.
// It follows platform conventions:
//   - Linux: ~/.config/rocketype/texts (XDG Base Directory)
//   - macOS: ~/Library/Application Support/rocketype/texts (Apple guidelines)
//   - Windows: %APPDATA%\rocketype\texts (Windows standard)
//
// If the directory doesn't exist, it will be created.
func GetDefaultTextsDir() (string, error) {
	var baseDir string
	var err error

	switch runtime.GOOS {
	case "linux", "freebsd", "openbsd", "netbsd":
		// Linux/BSD: Use XDG_CONFIG_HOME or ~/.config
		configHome := os.Getenv("XDG_CONFIG_HOME")
		if configHome == "" {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			configHome = filepath.Join(homeDir, ".config")
		}
		baseDir = filepath.Join(configHome, "rocketype", "texts")

	case "darwin":
		// macOS: Use ~/Library/Application Support
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		baseDir = filepath.Join(homeDir, "Library", "Application Support", "rocketype", "texts")

	case "windows":
		// Windows: Use %APPDATA%
		appData := os.Getenv("APPDATA")
		if appData == "" {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			appData = filepath.Join(homeDir, "AppData", "Roaming")
		}
		baseDir = filepath.Join(appData, "rocketype", "texts")

	default:
		// Fallback for unknown platforms
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		baseDir = filepath.Join(homeDir, ".rocketype", "texts")
	}

	// Create directory if it doesn't exist
	if err = os.MkdirAll(baseDir, 0755); err != nil {
		return "", err
	}

	return baseDir, nil
}

// GetFallbackTextsDir returns the local ./texts directory as a fallback.
// This is useful for development or when the user has texts in the current directory.
func GetFallbackTextsDir() string {
	return "texts"
}

// GetDefaultWordsDir returns the platform-appropriate default directory for word lists.
// It follows the same platform conventions as GetDefaultTextsDir but uses "words" subdirectory.
//   - Linux: ~/.config/rocketype/words (XDG Base Directory)
//   - macOS: ~/Library/Application Support/rocketype/words (Apple guidelines)
//   - Windows: %APPDATA%\rocketype\words (Windows standard)
//
// If the directory doesn't exist, it will be created.
func GetDefaultWordsDir() (string, error) {
	var baseDir string
	var err error

	switch runtime.GOOS {
	case "linux", "freebsd", "openbsd", "netbsd":
		// Linux/BSD: Use XDG_CONFIG_HOME or ~/.config
		configHome := os.Getenv("XDG_CONFIG_HOME")
		if configHome == "" {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			configHome = filepath.Join(homeDir, ".config")
		}
		baseDir = filepath.Join(configHome, "rocketype", "words")

	case "darwin":
		// macOS: Use ~/Library/Application Support
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		baseDir = filepath.Join(homeDir, "Library", "Application Support", "rocketype", "words")

	case "windows":
		// Windows: Use %APPDATA%
		appData := os.Getenv("APPDATA")
		if appData == "" {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			appData = filepath.Join(homeDir, "AppData", "Roaming")
		}
		baseDir = filepath.Join(appData, "rocketype", "words")

	default:
		// Fallback for unknown platforms
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		baseDir = filepath.Join(homeDir, ".rocketype", "words")
	}

	// Create directory if it doesn't exist
	if err = os.MkdirAll(baseDir, 0755); err != nil {
		return "", err
	}

	return baseDir, nil
}

// GetFallbackWordsDir returns the local ./words directory as a fallback.
func GetFallbackWordsDir() string {
	return "words"
}

// EnsureTextsDir ensures the texts directory exists and is readable.
// If the directory doesn't exist, it creates it and optionally copies
// default texts if available.
func EnsureTextsDir(dir string) error {
	// Check if directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// Create the directory
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	// Verify it's readable
	if _, err := os.ReadDir(dir); err != nil {
		return err
	}

	return nil
}
