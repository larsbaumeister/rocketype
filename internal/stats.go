package internal

import (
	"time"
)

const (
	// CharsPerWord represents the standard conversion factor for WPM calculation.
	// The industry standard is 5 characters = 1 word.
	CharsPerWord = 5.0
)

// Stats tracks typing test statistics including timing, accuracy, and error tracking.
// It maintains detailed information about keystrokes, misspelled words, and test progress.
//
// The Stats type is designed to track errors even after correction - if a word is
// typed incorrectly and then backspaced and fixed, it still counts as misspelled
// in the final results.
type Stats struct {
	// Timing information
	startTime time.Time
	endTime   time.Time

	// Keystroke tracking
	totalKeystrokes   int
	correctKeystrokes int

	// Word-level error tracking
	misspelledWords map[string]int // Maps word to count of times misspelled
	misspelledOrder []string       // Maintains insertion order of misspelled words

	// Current word tracking for real-time error detection
	currentWordStart int          // Index where current word starts
	wordHadError     map[int]bool // Maps word start position to error flag

	// Test state
	testComplete bool
}

// NewStats creates a new Stats instance with all fields properly initialized.
// Returns a pointer to a Stats struct ready for tracking typing test metrics.
func NewStats() *Stats {
	return &Stats{
		misspelledWords:  make(map[string]int),
		wordHadError:     make(map[int]bool),
		currentWordStart: 0,
		testComplete:     false,
	}
}

// Start begins timing the typing test.
// This method is idempotent - calling it multiple times has no effect after the first call.
// The start time is recorded on the first invocation only.
func (s *Stats) Start() {
	if s.startTime.IsZero() {
		s.startTime = time.Now()
	}
}

// Finish marks the typing test as complete and records the end time.
// This should be called when the user has typed all characters in the sample text.
func (s *Stats) Finish() {
	s.endTime = time.Now()
	s.testComplete = true
}

// IsComplete returns whether the typing test has finished.
func (s *Stats) IsComplete() bool {
	return s.testComplete
}

// RecordKeystroke records a single keystroke and tracks whether it was correct.
// This method updates both total keystroke count and correct keystroke count.
//
// Parameters:
//   - correct: true if the typed character matches the expected character
func (s *Stats) RecordKeystroke(correct bool) {
	s.totalKeystrokes++
	if correct {
		s.correctKeystrokes++
	}
}

// MarkCurrentWordAsError marks that the word starting at the given position has an error.
// This flag persists even if the user backspaces and corrects the error, ensuring that
// corrections don't hide mistakes in the final statistics.
//
// Parameters:
//   - wordStart: the character index where the word begins in the sample text
func (s *Stats) MarkCurrentWordAsError(wordStart int) {
	s.wordHadError[wordStart] = true
}

// WordHadError checks whether the word at the given position ever had an error.
// Returns true even if the error was subsequently corrected via backspace.
//
// Parameters:
//   - wordStart: the character index where the word begins in the sample text
func (s *Stats) WordHadError(wordStart int) bool {
	return s.wordHadError[wordStart]
}

// RecordMisspelledWord records a word that was misspelled during the test.
// If the word was already misspelled, increments its count. Empty strings are ignored.
// The first occurrence of each misspelled word is tracked for maintaining display order.
//
// Parameters:
//   - word: the word from the sample text that was typed incorrectly
func (s *Stats) RecordMisspelledWord(word string) {
	if word == "" {
		return
	}

	// Track first occurrence order for consistent display
	if s.misspelledWords[word] == 0 {
		s.misspelledOrder = append(s.misspelledOrder, word)
	}

	s.misspelledWords[word]++
}

// SetCurrentWordStart updates the index where the current word begins.
// This is used for tracking word boundaries as the user types.
//
// Parameters:
//   - index: the character position in the sample text where the current word starts
func (s *Stats) SetCurrentWordStart(index int) {
	s.currentWordStart = index
}

// GetCurrentWord extracts the word currently being typed from the sample text.
// The word is determined by the current word start position and extends to either
// the cursor position or the next space character, whichever comes first.
//
// Parameters:
//   - sampleText: the reference text the user is typing
//   - userInput: the text the user has typed so far (currently unused but kept for API consistency)
//   - cursorPos: the current position in the text
//
// Returns an empty string if the word start is beyond the sample text bounds.
func (s *Stats) GetCurrentWord(sampleText, userInput string, cursorPos int) string {
	wordStart := s.currentWordStart
	wordEnd := cursorPos

	// Find the end of the current word in sample text
	for wordEnd < len(sampleText) && sampleText[wordEnd] != ' ' {
		wordEnd++
	}

	if wordStart >= len(sampleText) {
		return ""
	}

	return sampleText[wordStart:wordEnd]
}

// CheckCurrentWordForErrors checks if any character in the specified word range
// was typed incorrectly. This method compares user input against the sample text
// character by character within the given word boundaries.
//
// Parameters:
//   - sampleText: the reference text the user is typing
//   - userInput: the text the user has typed so far
//   - wordStart: the starting index of the word to check
//   - wordEnd: the ending index of the word to check
//
// Returns true if any character mismatch is found, false otherwise.
func (s *Stats) CheckCurrentWordForErrors(sampleText, userInput string, wordStart, wordEnd int) bool {
	if wordStart >= len(userInput) {
		return false
	}

	// Check if any character in the word was typed incorrectly
	for i := wordStart; i < wordEnd && i < len(userInput) && i < len(sampleText); i++ {
		if userInput[i] != sampleText[i] {
			return true
		}
	}
	return false
}

// GetWPM calculates the typing speed in words per minute (WPM).
// Uses the industry standard of 5 characters = 1 word. Only correct keystrokes
// contribute to the WPM calculation.
//
// The calculation uses elapsed time from start to either:
//   - The current time (if test is ongoing)
//   - The end time (if test is complete)
//
// Returns 0 if:
//   - The test hasn't started
//   - Less than 1 second has elapsed
//   - No time has passed (edge case)
func (s *Stats) GetWPM() float64 {
	if s.startTime.IsZero() {
		return 0
	}

	var duration time.Duration
	if s.testComplete {
		duration = s.endTime.Sub(s.startTime)
	} else {
		duration = time.Since(s.startTime)
	}

	if duration.Seconds() < 1 {
		return 0
	}

	words := float64(s.correctKeystrokes) / CharsPerWord
	minutes := duration.Minutes()

	if minutes == 0 {
		return 0
	}

	return words / minutes
}

// GetAccuracy calculates typing accuracy as a percentage.
// Accuracy is the ratio of correct keystrokes to total keystrokes.
//
// Returns 100.0 if no keystrokes have been recorded yet.
func (s *Stats) GetAccuracy() float64 {
	if s.totalKeystrokes == 0 {
		return 100.0
	}
	return (float64(s.correctKeystrokes) / float64(s.totalKeystrokes)) * 100.0
}

// GetMisspelledWords returns misspelled words in insertion order.
// The order represents the sequence in which words were first typed incorrectly.
func (s *Stats) GetMisspelledWords() []string {
	return s.misspelledOrder
}

// GetMisspelledWordCount returns the count of how many times a specific word was misspelled.
// Returns 0 if the word was never misspelled.
func (s *Stats) GetMisspelledWordCount(word string) int {
	return s.misspelledWords[word]
}
