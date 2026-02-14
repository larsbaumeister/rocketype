package internal

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
