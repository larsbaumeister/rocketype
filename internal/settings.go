package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Settings represents persistent user preferences that survive across sessions.
// These settings are preserved even when clearing session data.
type Settings struct {
	ThemeName string `json:"theme_name"` // Current theme preference

	// Mode settings
	Mode string `json:"mode"` // "text" or "words"

	// Word mode settings
	LimitType   string `json:"limit_type"`    // "time" or "words"
	TimeLimit   int    `json:"time_limit"`    // Time limit in seconds (default: 60)
	WordLimit   int    `json:"word_limit"`    // Word count limit (default: 50)
	LastWordSet string `json:"last_word_set"` // Last selected word set name
}

// SettingsManager handles saving and loading user settings.
type SettingsManager struct {
	settingsPath string
}

// NewSettingsManager creates a new settings manager.
// It uses the platform-appropriate config directory.
func NewSettingsManager() (*SettingsManager, error) {
	configDir, err := getConfigDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get config directory: %w", err)
	}

	settingsPath := filepath.Join(configDir, "settings.json")

	return &SettingsManager{
		settingsPath: settingsPath,
	}, nil
}

// SaveSettings saves user settings to disk.
func (sm *SettingsManager) SaveSettings(settings Settings) error {
	// Marshal to JSON
	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal settings: %w", err)
	}

	// Write to file
	if err := os.WriteFile(sm.settingsPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write settings file: %w", err)
	}

	return nil
}

// LoadSettings loads user settings from disk.
// Returns default settings if no settings file exists.
func (sm *SettingsManager) LoadSettings() (*Settings, error) {
	// Check if settings file exists
	if _, err := os.Stat(sm.settingsPath); os.IsNotExist(err) {
		// Return default settings
		return &Settings{
			ThemeName:   "default",
			Mode:        "text",
			LimitType:   "time",
			TimeLimit:   60,
			WordLimit:   50,
			LastWordSet: "",
		}, nil
	}

	// Read file
	data, err := os.ReadFile(sm.settingsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read settings file: %w", err)
	}

	// Unmarshal JSON
	var settings Settings
	if err := json.Unmarshal(data, &settings); err != nil {
		return nil, fmt.Errorf("failed to unmarshal settings: %w", err)
	}

	// Apply defaults for any missing fields
	if settings.Mode == "" {
		settings.Mode = "text"
	}
	if settings.LimitType == "" {
		settings.LimitType = "time"
	}
	if settings.TimeLimit == 0 {
		settings.TimeLimit = 60
	}
	if settings.WordLimit == 0 {
		settings.WordLimit = 50
	}

	return &settings, nil
}

// GetSettingsPath returns the path to the settings file.
func (sm *SettingsManager) GetSettingsPath() string {
	return sm.settingsPath
}
