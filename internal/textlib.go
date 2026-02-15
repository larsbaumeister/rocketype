package internal

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// TextSource represents a typing test text with its metadata.
type TextSource struct {
	Name    string // Display name (filename without extension)
	Content string // The actual text content
	Path    string // Full file path
}

// TextLibrary manages the collection of available typing test texts.
type TextLibrary struct {
	texts       []TextSource
	currentIdx  int    // Index of currently selected text
	textsDir    string // Directory where text files are stored
	defaultText TextSource
	rand        *rand.Rand
}

// NewTextLibrary creates a new TextLibrary instance.
// It loads all .txt files from the specified directory, or uses the default
// embedded text if the directory doesn't exist or contains no files.
//
// Parameters:
//   - textsDir: directory path to search for .txt files
//
// Returns a TextLibrary with at least one text (the default if no files found).
func NewTextLibrary(textsDir string) *TextLibrary {
	tl := &TextLibrary{
		textsDir: textsDir,
		defaultText: TextSource{
			Name:    "Default (Tolkien)",
			Content: defaultSampleText,
			Path:    "",
		},
		texts:      make([]TextSource, 0),
		currentIdx: 0,
		rand:       rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	// Try to load texts from directory
	if err := tl.loadTexts(); err != nil {
		// If loading fails, use default text
		tl.texts = []TextSource{tl.defaultText}
	}

	// If no texts were loaded, add default
	if len(tl.texts) == 0 {
		tl.texts = []TextSource{tl.defaultText}
	}

	return tl
}

// loadTexts reads all .txt files from the texts directory.
func (tl *TextLibrary) loadTexts() error {
	// Check if directory exists
	if _, err := os.Stat(tl.textsDir); os.IsNotExist(err) {
		return fmt.Errorf("texts directory not found: %s", tl.textsDir)
	}

	// Read all files in directory
	entries, err := os.ReadDir(tl.textsDir)
	if err != nil {
		return fmt.Errorf("failed to read texts directory: %w", err)
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

		path := filepath.Join(tl.textsDir, entry.Name())
		content, err := os.ReadFile(path)
		if err != nil {
			// Skip files that can't be read
			continue
		}

		// Skip empty files
		text := strings.TrimSpace(string(content))
		if text == "" {
			continue
		}

		// Create text source
		name := strings.TrimSuffix(entry.Name(), ".txt")
		tl.texts = append(tl.texts, TextSource{
			Name:    name,
			Content: text,
			Path:    path,
		})
	}

	return nil
}

// GetCurrentText returns the currently selected text.
func (tl *TextLibrary) GetCurrentText() TextSource {
	if tl.currentIdx >= 0 && tl.currentIdx < len(tl.texts) {
		return tl.texts[tl.currentIdx]
	}
	return tl.defaultText
}

// SelectRandom selects a random text from the library.
func (tl *TextLibrary) SelectRandom() TextSource {
	if len(tl.texts) == 0 {
		return tl.defaultText
	}
	tl.currentIdx = tl.rand.Intn(len(tl.texts))
	return tl.GetCurrentText()
}

// SelectByIndex selects a text by its index in the library.
// Returns false if the index is out of bounds.
func (tl *TextLibrary) SelectByIndex(index int) bool {
	if index >= 0 && index < len(tl.texts) {
		tl.currentIdx = index
		return true
	}
	return false
}

// SelectByName selects a text by its name.
// Returns false if no text with that name is found.
func (tl *TextLibrary) SelectByName(name string) bool {
	for i, text := range tl.texts {
		if text.Name == name {
			tl.currentIdx = i
			return true
		}
	}
	return false
}

// GetAllTexts returns a slice of all available texts.
func (tl *TextLibrary) GetAllTexts() []TextSource {
	return tl.texts
}

// Count returns the number of available texts.
func (tl *TextLibrary) Count() int {
	return len(tl.texts)
}

// GetCurrentIndex returns the index of the currently selected text.
func (tl *TextLibrary) GetCurrentIndex() int {
	return tl.currentIdx
}

// AddText adds a new text to the library.
// This is useful for dynamically adding texts like stdin input.
func (tl *TextLibrary) AddText(text TextSource) {
	tl.texts = append(tl.texts, text)
}
