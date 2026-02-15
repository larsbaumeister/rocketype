package internal

import (
	"fmt"
	"time"
)

// TypingTest manages the business logic of a typing test session.
// It handles user input tracking, cursor position, word boundaries,
// and coordinates with Stats for accuracy and error tracking.
//
// This component is purely about the typing test logic and has no
// knowledge of rendering or UI concerns.
type TypingTest struct {
	sampleText  string // The reference text the user types
	sampleRunes []rune // Cached rune slice of sampleText for efficient Unicode handling
	userInput   string // What the user has typed so far
	userRunes   []rune // Cached rune slice of userInput for efficient Unicode handling
	cursorPos   int    // Current position in sampleText (in runes, not bytes)
	wordStart   int    // Index where current word starts (in runes, not bytes)
	stats       *Stats // Statistics tracker
	finished    bool   // Whether the test is complete
}

// NewTypingTest creates a new typing test with the given sample text.
func NewTypingTest(sampleText string) *TypingTest {
	return &TypingTest{
		sampleText:  sampleText,
		sampleRunes: []rune(sampleText),
		userInput:   "",
		userRunes:   []rune{},
		cursorPos:   0,
		wordStart:   0,
		stats:       NewStats(),
		finished:    false,
	}
}

// GetSampleText returns the reference text.
func (t *TypingTest) GetSampleText() string {
	return t.sampleText
}

// GetSampleRunes returns the reference text as a rune slice (cached).
func (t *TypingTest) GetSampleRunes() []rune {
	return t.sampleRunes
}

// GetUserInput returns what the user has typed so far.
func (t *TypingTest) GetUserInput() string {
	return t.userInput
}

// GetUserRunes returns what the user has typed as a rune slice (cached).
func (t *TypingTest) GetUserRunes() []rune {
	return t.userRunes
}

// GetCursorPos returns the current cursor position.
func (t *TypingTest) GetCursorPos() int {
	return t.cursorPos
}

// GetStats returns the statistics tracker.
func (t *TypingTest) GetStats() *Stats {
	return t.stats
}

// IsFinished returns whether the test is complete.
func (t *TypingTest) IsFinished() bool {
	return t.finished
}

// MarkFinished marks the test as complete and finalizes stats.
// This should be called when ending the test early (e.g., time/word limit reached in word mode).
func (t *TypingTest) MarkFinished() {
	if !t.finished {
		t.stats.Finish()
		t.finished = true
	}
}

// SetSampleText updates the sample text and resets the test.
func (t *TypingTest) SetSampleText(text string) {
	t.sampleText = text
	t.sampleRunes = []rune(text)
	t.Reset()
}

// Reset resets the test to initial state, keeping the same sample text.
func (t *TypingTest) Reset() {
	t.userInput = ""
	t.userRunes = []rune{}
	t.cursorPos = 0
	t.wordStart = 0
	t.stats = NewStats()
	t.finished = false
}

// RestoreState restores the test state from a saved session.
// This allows resuming a typing test from where the user left off.
func (t *TypingTest) RestoreState(userInput string, cursorPos int) {
	t.userInput = userInput
	t.userRunes = []rune(userInput)
	t.cursorPos = cursorPos

	// Find the start of the current word by looking backwards for a space or newline
	t.wordStart = 0
	for i := cursorPos - 1; i >= 0; i-- {
		if t.sampleRunes[i] == ' ' || t.sampleRunes[i] == '\n' {
			t.wordStart = i + 1
			break
		}
	}

	t.finished = false
	t.stats = NewStats()
	// Stats will start when user types next character
}

// GetStatsStartTime returns the start time from stats as an RFC3339 string.
// Returns empty string if test hasn't started.
func (t *TypingTest) GetStatsStartTime() string {
	startTime := t.stats.GetStartTime()
	if startTime.IsZero() {
		return ""
	}
	return startTime.Format(time.RFC3339)
}

// GetTotalKeystrokes returns the total keystrokes from stats.
func (t *TypingTest) GetTotalKeystrokes() int {
	return t.stats.GetTotalKeystrokes()
}

// GetCorrectKeystrokes returns the correct keystrokes from stats.
func (t *TypingTest) GetCorrectKeystrokes() int {
	return t.stats.GetCorrectKeystrokes()
}

// GetMisspelledWordsMap returns the misspelled words map from stats.
func (t *TypingTest) GetMisspelledWordsMap() map[string]int {
	return t.stats.GetMisspelledWordsMap()
}

// GetWordErrorsMap returns the word errors map as map[string]int for JSON serialization.
// Converts map[int]bool to map[string]int.
func (t *TypingTest) GetWordErrorsMap() map[string]int {
	wordErrors := t.stats.GetWordErrorsMap()
	result := make(map[string]int, len(wordErrors))
	for pos, hadError := range wordErrors {
		if hadError {
			result[fmt.Sprintf("%d", pos)] = 1
		}
	}
	return result
}

// RestoreStatsFromSession restores stats from saved session data.
func (t *TypingTest) RestoreStatsFromSession(startTimeStr string, totalKeystrokes, correctKeystrokes int, misspelledWords map[string]int, misspelledOrder []string, wordHadErrorMap map[string]int) error {
	// Parse start time
	var startTime time.Time
	var err error
	if startTimeStr != "" {
		startTime, err = time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			return fmt.Errorf("failed to parse start time: %w", err)
		}
	}

	// Convert wordHadErrorMap from map[string]int to map[int]bool
	wordHadError := make(map[int]bool)
	for posStr, val := range wordHadErrorMap {
		var pos int
		_, err := fmt.Sscanf(posStr, "%d", &pos)
		if err == nil && val != 0 {
			wordHadError[pos] = true
		}
	}

	// Restore stats
	t.stats.RestoreFromSession(startTime, totalKeystrokes, correctKeystrokes, misspelledWords, misspelledOrder, wordHadError)
	return nil
}

// TypeCharacter handles typing a regular character.
// Returns true if the character was processed (test not yet finished).
func (t *TypingTest) TypeCharacter(typedChar rune) bool {
	if t.cursorPos >= len(t.sampleRunes) {
		return false
	}

	t.stats.Start()

	expectedChar := t.sampleRunes[t.cursorPos]
	correct := expectedChar == typedChar

	// Record keystroke
	t.stats.RecordKeystroke(correct)

	// Mark word as having error if incorrect
	if !correct {
		t.stats.MarkCurrentWordAsError(t.wordStart)
	}

	t.userInput += string(typedChar)
	t.userRunes = append(t.userRunes, typedChar)
	t.cursorPos++

	// Handle word boundary (space completes a word)
	if typedChar == ' ' && t.wordStart < t.cursorPos-1 {
		t.finishWord(t.cursorPos - 1)
		t.wordStart = t.cursorPos
	}

	t.checkCompletion()
	return true
}

// TypeNewline handles typing a newline character (Enter key).
// Returns true if the character was processed (test not yet finished).
func (t *TypingTest) TypeNewline() bool {
	if t.cursorPos >= len(t.sampleRunes) {
		return false
	}

	t.stats.Start()

	expectedChar := t.sampleRunes[t.cursorPos]
	typedChar := '\n'
	correct := expectedChar == typedChar

	// Record keystroke
	t.stats.RecordKeystroke(correct)

	// Mark word as having error if incorrect
	if !correct {
		t.stats.MarkCurrentWordAsError(t.wordStart)
	}

	t.userInput += "\n"
	t.userRunes = append(t.userRunes, '\n')
	t.cursorPos++

	// Newline acts as word boundary
	if t.wordStart < t.cursorPos-1 {
		t.finishWord(t.cursorPos - 1)
	}
	t.wordStart = t.cursorPos

	t.checkCompletion()
	return true
}

// Backspace handles the backspace key, removing the last typed character.
func (t *TypingTest) Backspace() {
	if t.cursorPos <= 0 {
		return
	}

	t.cursorPos--

	// Remove last rune from both string and rune slice
	if len(t.userRunes) > 0 {
		t.userRunes = t.userRunes[:len(t.userRunes)-1]
		t.userInput = string(t.userRunes)
	}

	// Update word start if we backspaced into previous word
	if t.cursorPos < len(t.sampleRunes) && t.sampleRunes[t.cursorPos] == ' ' {
		// Find the start of the word we backspaced into
		for t.wordStart > 0 && t.sampleRunes[t.wordStart-1] != ' ' {
			t.wordStart--
		}
	}
}

// finishWord records a word as misspelled if it had any errors.
func (t *TypingTest) finishWord(wordEnd int) {
	if t.stats.WordHadError(t.wordStart) {
		word := string(t.sampleRunes[t.wordStart:wordEnd])
		t.stats.RecordMisspelledWord(word)
	}
}

// checkCompletion checks if the test is complete and finalizes stats.
func (t *TypingTest) checkCompletion() {
	if t.cursorPos >= len(t.sampleRunes) {
		// Record last word if it had errors
		if t.wordStart < len(t.sampleRunes) {
			if t.stats.WordHadError(t.wordStart) {
				word := string(t.sampleRunes[t.wordStart:])
				t.stats.RecordMisspelledWord(word)
			}
		}
		t.stats.Finish()
		t.finished = true
	}
}
