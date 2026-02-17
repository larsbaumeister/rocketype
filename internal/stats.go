package internal

import (
	"strings"
	"time"
)

const (
	// CharsPerWord represents the standard conversion factor for WPM calculation.
	// The industry standard is 5 characters = 1 word.
	CharsPerWord = 5.0
)

// WPMSnapshot represents a WPM measurement at a specific time.
type WPMSnapshot struct {
	Timestamp time.Time // When this measurement was taken
	WPM       float64   // Words per minute at this point
}

// keystrokeEvent records a single keystroke with its timestamp.
type keystrokeEvent struct {
	timestamp time.Time
	correct   bool
}

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

	// Instantaneous WPM tracking
	keystrokeEvents  []keystrokeEvent // Recent keystrokes with timestamps
	instantWindowSec float64          // Time window for instantaneous WPM (e.g., 3 seconds)

	// WPM timeline tracking
	wpmHistory          []WPMSnapshot // Historical WPM measurements
	lastSnapshotTime    time.Time     // Last time we took a snapshot
	snapshotIntervalSec float64       // Seconds between snapshots

	// Error tracking
	errorTimestamps []time.Time    // Timestamps of when errors occurred
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
		misspelledWords:     make(map[string]int),
		wordHadError:        make(map[int]bool),
		currentWordStart:    0,
		testComplete:        false,
		wpmHistory:          make([]WPMSnapshot, 0, 60),      // Pre-allocate for ~60 seconds
		errorTimestamps:     make([]time.Time, 0, 100),       // Pre-allocate for typical errors
		snapshotIntervalSec: 1.0,                             // Take snapshot every second
		keystrokeEvents:     make([]keystrokeEvent, 0, 1000), // Pre-allocate for typical keystrokes
		instantWindowSec:    3.0,                             // 3-second rolling window
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
// It also updates the WPM timeline at regular intervals.
//
// Parameters:
//   - correct: true if the typed character matches the expected character
func (s *Stats) RecordKeystroke(correct bool) {
	s.totalKeystrokes++
	if correct {
		s.correctKeystrokes++
	} else {
		// Record timestamp of error
		if !s.startTime.IsZero() {
			s.errorTimestamps = append(s.errorTimestamps, time.Now())
		}
	}

	// Record keystroke event with timestamp for instantaneous WPM
	if !s.startTime.IsZero() {
		s.keystrokeEvents = append(s.keystrokeEvents, keystrokeEvent{
			timestamp: time.Now(),
			correct:   correct,
		})
	}

	// Update WPM timeline
	s.updateWPMTimeline()
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
	// Trim all whitespace (spaces, tabs, newlines) from the word
	word = strings.TrimSpace(word)

	if word == "" {
		return
	}

	// If the word contains newlines or multiple spaces, it was incorrectly captured
	// Split it into individual words and record each separately
	if strings.ContainsAny(word, "\n\r\t") || strings.Contains(word, "  ") {
		// Split on whitespace and record each word separately
		words := strings.Fields(word)
		for _, w := range words {
			s.RecordMisspelledWord(w) // Recursive call with cleaned word
		}
		return
	}

	// Filter out suspicious single-character entries that aren't valid words
	// Keep "a", "A", "I", "i" as they are valid English words
	// Filter out other single characters/spaces that might be extraction errors
	if len(word) == 1 {
		validSingleChars := map[rune]bool{
			'a': true, 'A': true,
			'i': true, 'I': true,
		}
		if !validSingleChars[rune(word[0])] {
			// Skip single characters that aren't valid words
			return
		}
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

// GetStartTime returns the time when the test started.
// Returns zero time if test hasn't started yet.
func (s *Stats) GetStartTime() time.Time {
	return s.startTime
}

// GetTotalKeystrokes returns the total number of keystrokes recorded.
func (s *Stats) GetTotalKeystrokes() int {
	return s.totalKeystrokes
}

// GetCorrectKeystrokes returns the number of correct keystrokes recorded.
func (s *Stats) GetCorrectKeystrokes() int {
	return s.correctKeystrokes
}

// GetMisspelledWordsMap returns the map of misspelled words and their counts.
func (s *Stats) GetMisspelledWordsMap() map[string]int {
	// Return a copy to prevent external modification
	result := make(map[string]int, len(s.misspelledWords))
	for k, v := range s.misspelledWords {
		result[k] = v
	}
	return result
}

// GetWordErrorsMap returns the map of word positions that had errors.
// Returns as map[int]bool where key is word start position.
func (s *Stats) GetWordErrorsMap() map[int]bool {
	// Return a copy to prevent external modification
	result := make(map[int]bool, len(s.wordHadError))
	for k, v := range s.wordHadError {
		result[k] = v
	}
	return result
}

// getInstantaneousWPM calculates WPM based on keystrokes in the last N seconds.
// This gives a real-time measure of typing speed that drops to 0 when typing stops.
func (s *Stats) getInstantaneousWPM() float64 {
	if s.startTime.IsZero() || len(s.keystrokeEvents) == 0 {
		return 0
	}

	now := time.Now()
	cutoffTime := now.Add(-time.Duration(s.instantWindowSec * float64(time.Second)))

	// Count correct keystrokes in the rolling window
	correctInWindow := 0
	for i := len(s.keystrokeEvents) - 1; i >= 0; i-- {
		event := s.keystrokeEvents[i]
		if event.timestamp.Before(cutoffTime) {
			break // Events are chronological, so we can stop
		}
		if event.correct {
			correctInWindow++
		}
	}

	// Calculate WPM from keystrokes in window
	words := float64(correctInWindow) / CharsPerWord
	minutes := s.instantWindowSec / 60.0

	return words / minutes
}

// updateWPMTimeline takes a WPM snapshot if enough time has elapsed since the last snapshot.
// This is called automatically by RecordKeystroke.
func (s *Stats) updateWPMTimeline() {
	if s.startTime.IsZero() {
		return
	}

	now := time.Now()

	// Initialize last snapshot time if this is the first call
	if s.lastSnapshotTime.IsZero() {
		s.lastSnapshotTime = s.startTime
	}

	// Check if enough time has passed for a new snapshot
	elapsed := now.Sub(s.lastSnapshotTime).Seconds()
	if elapsed >= s.snapshotIntervalSec {
		// Calculate instantaneous WPM for the graph
		currentWPM := s.getInstantaneousWPM()

		// Add snapshot
		s.wpmHistory = append(s.wpmHistory, WPMSnapshot{
			Timestamp: now,
			WPM:       currentWPM,
		})

		s.lastSnapshotTime = now

		// Clean up old keystroke events to prevent unbounded growth
		// Keep events from the last 10 seconds for cleanup
		cleanupCutoff := now.Add(-10 * time.Second)
		firstValidIdx := 0
		for i, event := range s.keystrokeEvents {
			if !event.timestamp.Before(cleanupCutoff) {
				firstValidIdx = i
				break
			}
		}
		if firstValidIdx > 0 {
			s.keystrokeEvents = s.keystrokeEvents[firstValidIdx:]
		}
	}
}

// GetWPMHistory returns a copy of the WPM history for timeline display.
func (s *Stats) GetWPMHistory() []WPMSnapshot {
	// Return a copy to prevent external modification
	result := make([]WPMSnapshot, len(s.wpmHistory))
	copy(result, s.wpmHistory)
	return result
}

// GetErrorTimestamps returns a copy of the error timestamps for visualization.
func (s *Stats) GetErrorTimestamps() []time.Time {
	// Return a copy to prevent external modification
	result := make([]time.Time, len(s.errorTimestamps))
	copy(result, s.errorTimestamps)
	return result
}

// RestoreFromSession restores stats from saved session data.
// This allows resuming a typing test with accurate WPM and accuracy tracking.
func (s *Stats) RestoreFromSession(startTime time.Time, totalKeystrokes, correctKeystrokes int, misspelledWords map[string]int, misspelledOrder []string, wordHadError map[int]bool) {
	s.startTime = startTime
	s.totalKeystrokes = totalKeystrokes
	s.correctKeystrokes = correctKeystrokes
	s.misspelledWords = misspelledWords
	s.misspelledOrder = misspelledOrder
	s.wordHadError = wordHadError
	s.testComplete = false
}
