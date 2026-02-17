package internal

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// WordSet represents a word list with its metadata.
type WordSet struct {
	Name  string   // Display name (filename without extension)
	Words []string // The list of words
	Path  string   // Full file path
}

// WordLibrary manages the collection of available word sets.
type WordLibrary struct {
	wordSets   []WordSet
	currentIdx int    // Index of currently selected word set
	wordsDir   string // Directory where word files are stored
	rand       *rand.Rand
}

// NewWordLibrary creates a new WordLibrary instance.
// It loads all .txt files from the specified directory.
//
// Parameters:
//   - wordsDir: directory path to search for .txt files with word lists
//
// Returns a WordLibrary (may be empty if no files found).
func NewWordLibrary(wordsDir string) *WordLibrary {
	wl := &WordLibrary{
		wordsDir:   wordsDir,
		wordSets:   make([]WordSet, 0),
		currentIdx: 0,
		rand:       rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	// Try to load word sets from directory
	_ = wl.loadWordSets()

	return wl
}

// loadWordSets reads all .txt files from the words directory.
// Each file should contain words (one per line or space-separated).
func (wl *WordLibrary) loadWordSets() error {
	// Check if directory exists
	if _, err := os.Stat(wl.wordsDir); os.IsNotExist(err) {
		return fmt.Errorf("words directory not found: %s", wl.wordsDir)
	}

	// Read all files in directory
	entries, err := os.ReadDir(wl.wordsDir)
	if err != nil {
		return fmt.Errorf("failed to read words directory: %w", err)
	}

	// Load each .txt file
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// Only process .txt files
		if !strings.HasSuffix(strings.ToLower(entry.Name()), ".txt") {
			continue
		}

		path := filepath.Join(wl.wordsDir, entry.Name())
		content, err := os.ReadFile(path)
		if err != nil {
			// Skip files that can't be read
			continue
		}

		// Parse words from file (support both newline and space-separated)
		text := string(content)
		words := make([]string, 0)

		// Split by both newlines and spaces
		lines := strings.Split(text, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			// Split each line by spaces in case words are space-separated
			lineWords := strings.Fields(line)
			words = append(words, lineWords...)
		}

		// Skip empty word sets
		if len(words) == 0 {
			continue
		}

		// Create word set
		name := strings.TrimSuffix(entry.Name(), ".txt")
		wl.wordSets = append(wl.wordSets, WordSet{
			Name:  name,
			Words: words,
			Path:  path,
		})
	}

	return nil
}

// GetCurrentWordSet returns the currently selected word set.
// Returns empty WordSet if none selected or library is empty.
func (wl *WordLibrary) GetCurrentWordSet() WordSet {
	if wl.currentIdx >= 0 && wl.currentIdx < len(wl.wordSets) {
		return wl.wordSets[wl.currentIdx]
	}
	return WordSet{}
}

// SelectByName selects a word set by its name.
// Returns false if no word set with that name is found.
func (wl *WordLibrary) SelectByName(name string) bool {
	for i, wordSet := range wl.wordSets {
		if wordSet.Name == name {
			wl.currentIdx = i
			return true
		}
	}
	return false
}

// GetAllWordSets returns a slice of all available word sets.
func (wl *WordLibrary) GetAllWordSets() []WordSet {
	return wl.wordSets
}

// Count returns the number of available word sets.
func (wl *WordLibrary) Count() int {
	return len(wl.wordSets)
}

// GenerateRandomWords generates a string of random words from the current word set.
// Words are separated by spaces and selected randomly with replacement.
//
// Parameters:
//   - count: number of words to generate
//
// Returns empty string if no word set is selected or word set is empty.
func (wl *WordLibrary) GenerateRandomWords(count int) string {
	wordSet := wl.GetCurrentWordSet()
	if len(wordSet.Words) == 0 {
		return ""
	}

	words := make([]string, count)
	for i := range count {
		words[i] = wordSet.Words[wl.rand.Intn(len(wordSet.Words))]
	}

	return strings.Join(words, " ")
}

// HasWordSets returns true if the library has at least one word set.
func (wl *WordLibrary) HasWordSets() bool {
	return len(wl.wordSets) > 0
}
