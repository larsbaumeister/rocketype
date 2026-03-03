package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
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
	configDir, err := GetConfigDir()
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

// SaveLeaderboard writes the leaderboard map to disk atomically.
func SaveLeaderboard(leaderboards map[string][]LeaderboardEntry) error {
	path, err := GetLeaderboardPath()
	if err != nil {
		return fmt.Errorf("failed to resolve leaderboard path: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("failed to create leaderboard directory: %w", err)
	}

	data, err := json.MarshalIndent(leaderboards, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal leaderboard: %w", err)
	}

	return writeFileAtomic(path, data, 0644)
}

// LoadLeaderboard reads leaderboard data from disk.
// Returns an empty map if the file does not exist.
func LoadLeaderboard() (map[string][]LeaderboardEntry, error) {
	path, err := GetLeaderboardPath()
	if err != nil {
		return nil, fmt.Errorf("failed to resolve leaderboard path: %w", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return map[string][]LeaderboardEntry{}, nil
		}
		return nil, fmt.Errorf("failed to read leaderboard file: %w", err)
	}

	if len(data) == 0 {
		return map[string][]LeaderboardEntry{}, nil
	}

	var leaderboards map[string][]LeaderboardEntry
	if err := json.Unmarshal(data, &leaderboards); err != nil {
		if resetErr := ResetLeaderboard(); resetErr != nil {
			return nil, fmt.Errorf("leaderboard corrupt and reset failed: %w", resetErr)
		}
		return map[string][]LeaderboardEntry{}, fmt.Errorf("leaderboard corrupt, reset to empty: %w", err)
	}

	if leaderboards == nil {
		leaderboards = map[string][]LeaderboardEntry{}
	}

	return leaderboards, nil
}

// ResetLeaderboard resets the leaderboard file to an empty map.
func ResetLeaderboard() error {
	path, err := GetLeaderboardPath()
	if err != nil {
		return fmt.Errorf("failed to resolve leaderboard path: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("failed to create leaderboard directory: %w", err)
	}

	empty := map[string][]LeaderboardEntry{}
	data, err := json.MarshalIndent(empty, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal empty leaderboard: %w", err)
	}

	return writeFileAtomic(path, data, 0644)
}

func writeFileAtomic(path string, data []byte, perm os.FileMode) error {
	dir := filepath.Dir(path)
	base := filepath.Base(path)
	tmpFile, err := os.CreateTemp(dir, base+".tmp")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer func() {
		_ = os.Remove(tmpFile.Name())
	}()

	if _, err := tmpFile.Write(data); err != nil {
		_ = tmpFile.Close()
		return fmt.Errorf("failed to write temp file: %w", err)
	}

	if err := tmpFile.Sync(); err != nil {
		_ = tmpFile.Close()
		return fmt.Errorf("failed to sync temp file: %w", err)
	}

	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("failed to close temp file: %w", err)
	}

	backupPath := path + ".bak"
	if _, err := os.Stat(path); err == nil {
		_ = os.Remove(backupPath)
		if err := os.Rename(path, backupPath); err != nil {
			return fmt.Errorf("failed to backup leaderboard file: %w", err)
		}
	}

	if err := os.Rename(tmpFile.Name(), path); err != nil {
		if _, statErr := os.Stat(backupPath); statErr == nil {
			_ = os.Rename(backupPath, path)
		}
		return fmt.Errorf("failed to replace leaderboard file: %w", err)
	}

	return nil
}
