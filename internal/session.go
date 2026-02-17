package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// Session represents a saved typing session that can be resumed.
// This contains only ephemeral progress data that gets cleared when
// you finish a test, restart, or select a new text.
type Session struct {
	// Text information
	TextName    string `json:"text_name"`    // Name of the text being typed
	TextContent string `json:"text_content"` // Full text content
	TextPath    string `json:"text_path"`    // Path to text file (if from file)

	// Progress information
	UserInput string `json:"user_input"` // What the user has typed so far
	CursorPos int    `json:"cursor_pos"` // Current cursor position (in runes)

	// Statistics (for accurate WPM and accuracy restoration)
	StartTime         string         `json:"start_time"`         // ISO timestamp when test started
	TotalKeystrokes   int            `json:"total_keystrokes"`   // Total keys pressed
	CorrectKeystrokes int            `json:"correct_keystrokes"` // Correct keys pressed
	MisspelledWords   map[string]int `json:"misspelled_words"`   // Words misspelled and their counts
	MisspelledOrder   []string       `json:"misspelled_order"`   // Order words were first misspelled
	WordHadError      map[string]int `json:"word_had_error"`     // Word positions that had errors (as string keys for JSON)

	// Metadata
	SavedAt string `json:"saved_at"` // ISO timestamp when saved

	// Note: Theme is stored separately in settings.json, not here.
}

// SessionManager handles saving and loading typing sessions.
type SessionManager struct {
	sessionPath string
}

// NewSessionManager creates a new session manager.
// It uses the platform-appropriate config directory.
func NewSessionManager() (*SessionManager, error) {
	configDir, err := getConfigDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get config directory: %w", err)
	}

	sessionPath := filepath.Join(configDir, "session.json")

	return &SessionManager{
		sessionPath: sessionPath,
	}, nil
}

// SaveSession saves the current typing session to disk.
func (sm *SessionManager) SaveSession(session Session) error {
	// Normalize whitespace before saving to ensure consistency
	session.TextContent = NormalizeWhitespace(session.TextContent)

	// Add timestamp
	session.SavedAt = time.Now().Format(time.RFC3339)

	// Marshal to JSON
	data, err := json.MarshalIndent(session, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	// Write to file
	if err := os.WriteFile(sm.sessionPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write session file: %w", err)
	}

	return nil
}

// LoadSession loads a saved session from disk.
// Returns nil if no session file exists.
func (sm *SessionManager) LoadSession() (*Session, error) {
	// Check if session file exists
	if _, err := os.Stat(sm.sessionPath); os.IsNotExist(err) {
		return nil, nil // No session to load
	}

	// Read file
	data, err := os.ReadFile(sm.sessionPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read session file: %w", err)
	}

	// Unmarshal JSON
	var session Session
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session: %w", err)
	}

	// Normalize whitespace in loaded content to ensure compatibility
	session.TextContent = NormalizeWhitespace(session.TextContent)

	return &session, nil
}

// HasSession checks if a saved session exists.
func (sm *SessionManager) HasSession() bool {
	_, err := os.Stat(sm.sessionPath)
	return err == nil
}

// ClearSession deletes the saved session file.
func (sm *SessionManager) ClearSession() error {
	if _, err := os.Stat(sm.sessionPath); os.IsNotExist(err) {
		return nil // Already doesn't exist
	}

	if err := os.Remove(sm.sessionPath); err != nil {
		return fmt.Errorf("failed to remove session file: %w", err)
	}

	return nil
}

// GetSessionPath returns the path to the session file.
func (sm *SessionManager) GetSessionPath() string {
	return sm.sessionPath
}

// getConfigDir returns the platform-appropriate config directory.
// This is similar to GetDefaultTextsDir but returns the base config dir.
func getConfigDir() (string, error) {
	var configDir string

	switch runtime.GOOS {
	case "windows":
		appData := os.Getenv("APPDATA")
		if appData == "" {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			appData = filepath.Join(homeDir, "AppData", "Roaming")
		}
		configDir = filepath.Join(appData, "rocketype")

	case "darwin":
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		configDir = filepath.Join(homeDir, "Library", "Application Support", "rocketype")

	default: // Linux and other Unix-like
		xdgConfig := os.Getenv("XDG_CONFIG_HOME")
		if xdgConfig == "" {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return "", err
			}
			xdgConfig = filepath.Join(homeDir, ".config")
		}
		configDir = filepath.Join(xdgConfig, "rocketype")
	}

	// Ensure directory exists
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", err
	}

	return configDir, nil
}

// CreateSessionFromApp creates a Session from the current app state.
// Note: Theme is not included as it's stored separately in settings.
func CreateSessionFromApp(app *App) Session {
	currentText := app.textLibrary.GetCurrentText()

	return Session{
		TextName:    currentText.Name,
		TextContent: app.typingTest.GetSampleText(),
		TextPath:    currentText.Path,
		UserInput:   app.typingTest.GetUserInput(),
		CursorPos:   app.typingTest.GetCursorPos(),
	}
}
