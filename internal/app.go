package internal

import (
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
)

// App represents the main application controller.
// It coordinates between UI rendering (Renderer), business logic (TypingTest),
// input handling (InputHandler), and application state (themes, commands, text library).
//
// The refactored App follows clean architecture principles:
//   - Renderer handles all drawing logic
//   - TypingTest manages typing test business logic
//   - InputHandler routes keyboard events
//   - App coordinates between components and manages state
type App struct {
	// Components
	renderer        *Renderer
	typingTest      *TypingTest
	inputHandler    *InputHandler
	commandMenu     *CommandMenu
	textLibrary     *TextLibrary
	wordLibrary     *WordLibrary
	sessionManager  *SessionManager
	settingsManager *SettingsManager

	// State
	theme       Theme
	screen      tcell.Screen
	quit        bool
	showResults bool

	// Mode settings
	mode              string    // "text" or "words"
	limitType         string    // "time" or "words"
	timeLimit         int       // Time limit in seconds
	wordLimit         int       // Word count limit
	testStarted       time.Time // When test was started (for time limit)
	lastCheckPosition int       // Last cursor position when we checked for more words (optimization)
}

const (
	// defaultSampleText is the fallback text when no texts directory exists.
	defaultSampleText = "Roads go ever ever on,\nOver rock and under tree,\nBy caves where never sun has shone,\nBy streams that never find the sea;\nOver snow by winter sown,\nAnd through the merry flowers of June,\nOver grass and over stone,\nAnd under mountains in the moon."

	// Word mode constants
	initialWordCount        = 100 // Initial words generated when entering word mode
	wordGenerationChunk     = 50  // Number of words to generate when buffer runs low
	wordModeVisibleLines    = 3   // Number of lines visible in word mode (cursor + 2 below)
	wordModeLinesThreshold  = 3   // Minimum lines remaining before generating more words
	timerUpdateIntervalMS   = 100 // Timer update interval in milliseconds
	wordLimitMultiplier     = 2   // Multiplier for initial word generation in word limit mode
	lastCheckPositionOffset = 10  // Don't check for more words until cursor advances by this many characters
)

// NewApp creates a new application instance and initializes all components.
//
// Parameters:
//   - stdinText: optional text from stdin (empty string if not provided)
//   - textsDir: directory path for text files
//   - restoreSession: whether to attempt to restore a saved session
//
// Returns an error if the screen cannot be created or initialized.
func NewApp(stdinText, textsDir string, restoreSession bool) (*App, error) {
	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, fmt.Errorf("failed to create screen: %w", err)
	}

	if err := screen.Init(); err != nil {
		return nil, fmt.Errorf("failed to initialize screen: %w", err)
	}

	// Initialize session manager
	sessionManager, err := NewSessionManager()
	if err != nil {
		return nil, fmt.Errorf("failed to create session manager: %w", err)
	}

	// Initialize settings manager
	settingsManager, err := NewSettingsManager()
	if err != nil {
		return nil, fmt.Errorf("failed to create settings manager: %w", err)
	}

	// Load settings (theme, etc.)
	settings, err := settingsManager.LoadSettings()
	if err != nil {
		// If settings load fails, use defaults
		settings = &Settings{ThemeName: "default"}
	}

	// Resolve theme from settings
	initialTheme := DefaultTheme
	for _, theme := range AvailableThemes() {
		if theme.Name == settings.ThemeName {
			initialTheme = theme
			break
		}
	}

	// Load text library
	textLibrary := NewTextLibrary(textsDir)

	// Load word library
	wordsDir, err := GetDefaultWordsDir()
	if err != nil {
		wordsDir = GetFallbackWordsDir()
	}
	wordLibrary := NewWordLibrary(wordsDir)

	// Try to restore session if requested and available (unless stdin is provided)
	var initialText TextSource
	var typingTest *TypingTest

	// stdin text takes precedence over session restoration, always text mode
	if stdinText != "" {
		stdinSource := TextSource{
			Name:    "stdin",
			Content: stdinText,
			Path:    "",
		}
		textLibrary.AddText(stdinSource)
		textLibrary.SelectByName("stdin")
		initialText = stdinSource
		typingTest = NewTypingTest(initialText.Content)
		settings.Mode = "text" // Force text mode for stdin
	} else if restoreSession && sessionManager.HasSession() {
		session, err := sessionManager.LoadSession()
		if err == nil && session != nil {
			// Restore from session
			initialText = TextSource{
				Name:    session.TextName,
				Content: session.TextContent,
				Path:    session.TextPath,
			}

			// Create typing test with restored state
			typingTest = NewTypingTest(session.TextContent)
			// Restore progress
			typingTest.RestoreState(session.UserInput, session.CursorPos)
			// Restore stats
			_ = typingTest.RestoreStatsFromSession(
				session.StartTime,
				session.TotalKeystrokes,
				session.CorrectKeystrokes,
				session.MisspelledWords,
				session.MisspelledOrder,
				session.WordHadError,
			)

			// Add to library if not already there
			textLibrary.AddText(initialText)
		} else {
			// Session loading failed, initialize based on mode
			if settings.Mode == "words" && wordLibrary.HasWordSets() {
				// Word mode - generate random words
				if settings.LastWordSet != "" {
					wordLibrary.SelectByName(settings.LastWordSet)
				}
				wordCount := initialWordCount
				if settings.LimitType == "words" {
					wordCount = settings.WordLimit * wordLimitMultiplier
				}
				content := wordLibrary.GenerateRandomWords(wordCount)
				initialText = TextSource{
					Name:    "Random Words",
					Content: content,
					Path:    "",
				}
				typingTest = NewTypingTest(content)
			} else {
				// Text mode - use random text
				settings.Mode = "text"
				initialText = textLibrary.SelectRandom()
				typingTest = NewTypingTest(initialText.Content)
			}
		}
	} else {
		// No stdin, no session - initialize based on mode
		if settings.Mode == "words" && wordLibrary.HasWordSets() {
			// Word mode - generate random words
			if settings.LastWordSet != "" {
				wordLibrary.SelectByName(settings.LastWordSet)
			}
			wordCount := initialWordCount
			if settings.LimitType == "words" {
				wordCount = settings.WordLimit * wordLimitMultiplier
			}
			content := wordLibrary.GenerateRandomWords(wordCount)
			initialText = TextSource{
				Name:    "Random Words",
				Content: content,
				Path:    "",
			}
			typingTest = NewTypingTest(content)
		} else {
			// Text mode - use random text
			settings.Mode = "text"
			initialText = textLibrary.SelectRandom()
			typingTest = NewTypingTest(initialText.Content)
		}
	}

	// Create components
	renderer := NewRenderer(screen)
	commandMenu := NewCommandMenu()

	app := &App{
		renderer:        renderer,
		typingTest:      typingTest,
		commandMenu:     commandMenu,
		textLibrary:     textLibrary,
		wordLibrary:     wordLibrary,
		sessionManager:  sessionManager,
		settingsManager: settingsManager,
		theme:           initialTheme,
		screen:          screen,
		quit:            false,
		showResults:     false,
		mode:            settings.Mode,
		limitType:       settings.LimitType,
		timeLimit:       settings.TimeLimit,
		wordLimit:       settings.WordLimit,
		testStarted:     time.Time{}, // Will be set when typing starts
	}

	// Initialize input handler with callbacks
	app.inputHandler = NewInputHandler(
		func() { app.quit = true },
		func() { app.toggleCommandMenu() },
		func() { app.cycleTheme() },
		func() { app.restartTest() },
		typingTest,
		commandMenu,
	)

	// Initialize commands
	app.initCommands()

	return app, nil
}

// Run starts the main application event loop.
// This method blocks until the application is quit.
func (a *App) Run() error {
	defer a.screen.Fini()

	a.draw()

	// Start a ticker for updating the display (e.g., timer countdown)
	ticker := time.NewTicker(timerUpdateIntervalMS * time.Millisecond)
	defer ticker.Stop()

	// Channel for screen events
	eventChan := make(chan tcell.Event)
	quitEventLoop := make(chan struct{})

	// Start event polling goroutine
	go func() {
		defer close(eventChan)
		for {
			select {
			case <-quitEventLoop:
				return
			default:
				ev := a.screen.PollEvent()
				select {
				case eventChan <- ev:
				case <-quitEventLoop:
					return
				}
			}
		}
	}()

	for !a.quit {
		select {
		case ev, ok := <-eventChan:
			if !ok {
				// Event channel closed, exit
				return nil
			}
			switch ev := ev.(type) {
			case *tcell.EventResize:
				a.screen.Sync()
				a.draw()

			case *tcell.EventKey:
				a.handleKey(ev)
				a.draw()
			}

		case <-ticker.C:
			// Periodic redraw for timer updates
			if a.mode == "words" && a.limitType == "time" && !a.testStarted.IsZero() {
				a.draw()
			}
		}
	}

	// Signal the event polling goroutine to stop
	close(quitEventLoop)

	// Handle session and settings on quit
	if a.typingTest.IsFinished() {
		// Test is finished - clear any saved session
		_ = a.sessionManager.ClearSession()
	} else if a.typingTest.GetCursorPos() > 0 {
		// Test in progress - save session with stats
		currentText := a.textLibrary.GetCurrentText()
		stats := a.typingTest.GetStats()

		session := Session{
			TextName:          currentText.Name,
			TextContent:       a.typingTest.GetSampleText(),
			TextPath:          currentText.Path,
			UserInput:         a.typingTest.GetUserInput(),
			CursorPos:         a.typingTest.GetCursorPos(),
			StartTime:         a.typingTest.GetStatsStartTime(),
			TotalKeystrokes:   a.typingTest.GetTotalKeystrokes(),
			CorrectKeystrokes: a.typingTest.GetCorrectKeystrokes(),
			MisspelledWords:   a.typingTest.GetMisspelledWordsMap(),
			MisspelledOrder:   stats.GetMisspelledWords(),
			WordHadError:      a.typingTest.GetWordErrorsMap(),
		}
		err := a.sessionManager.SaveSession(session)
		if err != nil {
			// Log error but don't fail the quit
			_ = err
		}
	}

	// Always save settings (theme preference and mode settings persist)
	currentWordSet := ""
	if a.mode == "words" {
		wordSet := a.wordLibrary.GetCurrentWordSet()
		currentWordSet = wordSet.Name
	}

	settings := Settings{
		ThemeName:   a.theme.Name,
		Mode:        a.mode,
		LimitType:   a.limitType,
		TimeLimit:   a.timeLimit,
		WordLimit:   a.wordLimit,
		LastWordSet: currentWordSet,
	}
	_ = a.settingsManager.SaveSettings(settings)

	return nil
}

// handleKey routes keyboard events to the input handler with current mode.
func (a *App) handleKey(ev *tcell.EventKey) {
	mode := a.getCurrentMode()

	// Special case: command menu execution needs app context
	if mode == ModeCommandMenu && ev.Key() == tcell.KeyEnter {
		a.commandMenu.ExecuteSelected(a)
		a.commandMenu.Hide()
		return
	}

	a.inputHandler.HandleKey(ev, mode)

	// Track test start time for word mode limits
	if mode == ModeTyping && a.mode == "words" && a.testStarted.IsZero() && a.typingTest.GetCursorPos() > 0 {
		a.testStarted = time.Now()
	}

	// Dynamically extend text in word mode if needed
	if mode == ModeTyping && a.mode == "words" {
		a.ensureEnoughWords()
	}

	// Check limits in word mode
	if a.mode == "words" && !a.typingTest.IsFinished() {
		limitReached := false

		if a.limitType == "time" && !a.testStarted.IsZero() {
			elapsed := time.Since(a.testStarted).Seconds()
			if elapsed >= float64(a.timeLimit) {
				limitReached = true
			}
		} else if a.limitType == "words" {
			// Count words typed by splitting user input
			userInput := a.typingTest.GetUserInput()
			wordCount := len(strings.Fields(userInput))
			if wordCount >= a.wordLimit {
				limitReached = true
			}
		}

		if limitReached {
			// Mark test as finished and show results
			a.typingTest.MarkFinished()
			a.showResults = true
		}
	}

	// Update results state after input
	if a.typingTest.IsFinished() {
		a.showResults = true
	}
}

// getCurrentMode determines the current application mode.
func (a *App) getCurrentMode() AppMode {
	if a.commandMenu.IsVisible() {
		return ModeCommandMenu
	}
	if a.showResults {
		return ModeResults
	}
	return ModeTyping
}

// draw renders the entire UI using the Renderer.
func (a *App) draw() {
	a.renderer.Clear()
	a.renderer.FillBackground(a.theme.Background)

	// Draw title with mode information
	var textName string
	var modeInfo string

	if a.mode == "words" {
		wordSet := a.wordLibrary.GetCurrentWordSet()
		textName = wordSet.Name
		if a.limitType == "time" {
			modeInfo = fmt.Sprintf("words mode, %ds", a.timeLimit)
		} else {
			modeInfo = fmt.Sprintf("words mode, %d words", a.wordLimit)
		}
	} else {
		currentText := a.textLibrary.GetCurrentText()
		textName = currentText.Name
		modeInfo = ""
	}

	a.renderer.DrawTitle(a.theme.Name, textName, modeInfo, a.theme)

	// Draw main content
	if a.showResults {
		a.drawResultsScreen()
	} else {
		a.drawTypingScreen()
	}

	// Draw overlays (always on top)
	if a.commandMenu.IsVisible() {
		a.drawCommandMenuOverlay()
	}

	a.renderer.Show()
}

// drawTypingScreen renders the typing test interface.
func (a *App) drawTypingScreen() {
	width, height := a.screen.Size()

	// Calculate text wrapping parameters
	maxWidth := width - 8
	if maxWidth < 20 {
		maxWidth = width
	}

	// Calculate available height and visible lines
	availableHeight := height - 8
	maxVisibleLines := availableHeight / 2 // 2 screen rows per text line

	// In word mode, only show 2 lines below cursor
	if a.mode == "words" {
		maxVisibleLines = wordModeVisibleLines // cursor line + 2 lines below
	}

	// Get cached rune slices (no conversion needed!)
	sampleRunes := a.typingTest.GetSampleRunes()
	cursorPos := a.typingTest.GetCursorPos()

	// Calculate which line the cursor is on (use string for wrapping)
	sampleText := a.typingTest.GetSampleText()
	cursorLine := CalculateCursorLine(sampleText, cursorPos, maxWidth)

	// Calculate total wrapped lines
	lines := wrapText(sampleText, maxWidth)
	totalLines := len(lines)

	// Calculate scroll position
	var scrollLine int
	if a.mode == "words" {
		// In word mode, scroll so cursor is on the first visible line
		scrollLine = cursorLine
		if scrollLine < 0 {
			scrollLine = 0
		}
	} else {
		// In text mode, use standard scroll calculation
		scrollLine = CalculateScrollLine(cursorLine, maxVisibleLines, totalLines)
	}

	// Draw typing view with cached rune slices
	viewData := TypingViewData{
		SampleText:  sampleText,
		SampleRunes: sampleRunes,
		UserInput:   a.typingTest.GetUserInput(),
		UserRunes:   a.typingTest.GetUserRunes(),
		CursorPos:   cursorPos,
		ScrollLine:  scrollLine,
		Theme:       a.theme,
		WordMode:    a.mode == "words",
	}
	a.renderer.DrawTypingView(viewData)

	// Draw stats
	stats := a.typingTest.GetStats()
	a.renderer.DrawStats(stats.GetWPM(), stats.GetAccuracy(), a.theme)

	// Draw progress for word mode
	if a.mode == "words" && !a.testStarted.IsZero() {
		var progressText string
		if a.limitType == "time" {
			elapsed := time.Since(a.testStarted).Seconds()
			remaining := float64(a.timeLimit) - elapsed
			if remaining < 0 {
				remaining = 0
			}
			progressText = fmt.Sprintf("Time: %.1fs", remaining)
		} else {
			// Count words typed
			wordsTyped := len(strings.Fields(a.typingTest.GetUserInput()))
			progressText = fmt.Sprintf("Words: %d / %d", wordsTyped, a.wordLimit)
		}
		a.renderer.DrawProgress(progressText, a.theme)
	}

	// Draw help text
	a.renderer.DrawHelpText(a.theme)
}

// drawResultsScreen renders the results screen.
func (a *App) drawResultsScreen() {
	stats := a.typingTest.GetStats()
	misspelledWords := stats.GetMisspelledWords()

	// Build word counts map
	wordCounts := make(map[string]int)
	for _, word := range misspelledWords {
		wordCounts[word] = stats.GetMisspelledWordCount(word)
	}

	resultsData := ResultsData{
		WPM:             stats.GetWPM(),
		Accuracy:        stats.GetAccuracy(),
		MisspelledWords: misspelledWords,
		WordCounts:      wordCounts,
		Theme:           a.theme,
	}
	a.renderer.DrawResults(resultsData)
}

// drawCommandMenuOverlay renders the command menu.
func (a *App) drawCommandMenuOverlay() {
	menuData := CommandMenuData{
		Filter:           a.commandMenu.GetFilter(),
		FilteredCommands: a.commandMenu.GetFilteredCommands(),
		Selected:         a.commandMenu.GetSelected(),
		Theme:            a.theme,
	}
	a.renderer.DrawCommandMenu(menuData)
}

// toggleCommandMenu toggles the command menu visibility.
func (a *App) toggleCommandMenu() {
	if a.commandMenu.IsVisible() {
		a.commandMenu.Hide()
	} else {
		a.commandMenu.Show()
	}
}

// cycleTheme switches to the next theme and saves the preference.
func (a *App) cycleTheme() {
	a.theme = GetNextTheme(a.theme)
	a.saveThemePreference()
}

// saveThemePreference saves the current theme to settings.
func (a *App) saveThemePreference() {
	settings := Settings{
		ThemeName:   a.theme.Name,
		Mode:        a.mode,
		LimitType:   a.limitType,
		TimeLimit:   a.timeLimit,
		WordLimit:   a.wordLimit,
		LastWordSet: a.getLastWordSet(),
	}
	_ = a.settingsManager.SaveSettings(settings)
}

// saveAllSettings saves all current settings including theme, mode, and limits.
func (a *App) saveAllSettings() {
	settings := Settings{
		ThemeName:   a.theme.Name,
		Mode:        a.mode,
		LimitType:   a.limitType,
		TimeLimit:   a.timeLimit,
		WordLimit:   a.wordLimit,
		LastWordSet: a.getLastWordSet(),
	}
	_ = a.settingsManager.SaveSettings(settings)
}

// getLastWordSet returns the current word set name or empty string.
func (a *App) getLastWordSet() string {
	if a.mode == "words" {
		wordSet := a.wordLibrary.GetCurrentWordSet()
		return wordSet.Name
	}
	return ""
}

// restartTest resets the current typing test.
func (a *App) restartTest() {
	a.typingTest.Reset()
	a.showResults = false
	a.testStarted = time.Time{} // Reset timer for word mode
	// Clear saved session when explicitly restarting
	_ = a.sessionManager.ClearSession()
}

// clearSession clears the saved session file and resets text/progress to defaults.
// Non-text-related settings like theme and mode are preserved.
func (a *App) clearSession() {
	if err := a.sessionManager.ClearSession(); err != nil {
		_ = err
	}

	// Reset based on current mode
	if a.mode == "words" && a.wordLibrary.HasWordSets() {
		// Generate new random words
		wordCount := initialWordCount
		if a.limitType == "words" {
			wordCount = a.wordLimit * wordLimitMultiplier
		}
		content := a.wordLibrary.GenerateRandomWords(wordCount)
		a.typingTest.SetSampleText(content)
	} else {
		// Select random text
		text := a.textLibrary.SelectRandom()
		a.typingTest.SetSampleText(text.Content)
	}

	a.showResults = false
	a.testStarted = time.Time{} // Reset timer
}

// selectRandomText selects a random text and restarts the test.
func (a *App) selectRandomText() {
	text := a.textLibrary.SelectRandom()
	a.typingTest.SetSampleText(text.Content)
	a.mode = "text"
	a.testStarted = time.Time{}
	// Clear saved session when selecting new text
	_ = a.sessionManager.ClearSession()
	a.saveAllSettings()
}

// selectTextByName selects a text by name and restarts the test.
func (a *App) selectTextByName(name string) {
	if a.textLibrary.SelectByName(name) {
		text := a.textLibrary.GetCurrentText()
		a.typingTest.SetSampleText(text.Content)
		a.mode = "text"
		a.testStarted = time.Time{}
		// Clear saved session when selecting new text
		_ = a.sessionManager.ClearSession()
		a.saveAllSettings()
	}
}

// selectWordSet selects a word set and generates random words.
func (a *App) selectWordSet(name string) {
	if a.wordLibrary.SelectByName(name) {
		a.mode = "words"
		// Start with a reasonable initial amount of words
		// We'll dynamically generate more as the user types
		content := a.wordLibrary.GenerateRandomWords(initialWordCount)
		a.typingTest.SetSampleText(content)
		a.testStarted = time.Time{}
		a.lastCheckPosition = 0 // Reset check position
		_ = a.sessionManager.ClearSession()
		a.saveAllSettings()
	}
}

// ensureEnoughWords checks if there's enough text ahead of the cursor and generates more if needed.
// This ensures the user always has at least 2 lines of text visible below the cursor.
// Optimized to only check periodically (not on every keystroke) for performance.
func (a *App) ensureEnoughWords() {
	cursorPos := a.typingTest.GetCursorPos()

	// Performance optimization: only check when cursor advances significantly
	if cursorPos < a.lastCheckPosition+lastCheckPositionOffset {
		return
	}
	a.lastCheckPosition = cursorPos

	width, _ := a.screen.Size()
	maxWidth := width - 8
	if maxWidth < 20 {
		maxWidth = width
	}

	sampleText := a.typingTest.GetSampleText()

	// Calculate how much text remains after cursor
	remainingText := ""
	sampleRunes := []rune(sampleText)
	if cursorPos < len(sampleRunes) {
		remainingText = string(sampleRunes[cursorPos:])
	}

	// Wrap remaining text to see how many lines are left
	remainingLines := wrapText(remainingText, maxWidth)

	// If less than threshold lines remaining, generate more words
	if len(remainingLines) < wordModeLinesThreshold {
		// Generate a chunk of new words
		newWords := a.wordLibrary.GenerateRandomWords(wordGenerationChunk)
		if newWords != "" {
			// Append new words to existing text
			updatedText := sampleText + " " + newWords
			a.typingTest.SetSampleText(updatedText)
		}
	}
}

// setTimeLimit sets the time limit in seconds and switches to time-based limit.
func (a *App) setTimeLimit(seconds int) {
	a.timeLimit = seconds
	a.limitType = "time"
	a.saveAllSettings()
}

// setWordLimit sets the word count limit and switches to word-based limit.
func (a *App) setWordLimit(words int) {
	a.wordLimit = words
	a.limitType = "words"
	// If already in word mode, regenerate text with appropriate word count
	if a.mode == "words" {
		wordCount := words * wordLimitMultiplier
		content := a.wordLibrary.GenerateRandomWords(wordCount)
		a.typingTest.SetSampleText(content)
		a.testStarted = time.Time{}
		a.lastCheckPosition = 0 // Reset check position
	}
	a.saveAllSettings()
}

// initCommands initializes the command palette with all available commands.
func (a *App) initCommands() {
	commands := []Command{
		{
			Name:        "theme: default",
			Description: "Switch to default terminal theme",
			Action: func(app *App) {
				app.theme = DefaultTheme
				app.saveThemePreference()
			},
		},
		{
			Name:        "theme: gruvbox",
			Description: "Switch to gruvbox theme (dark)",
			Action: func(app *App) {
				app.theme = GruvboxTheme
				app.saveThemePreference()
			},
		},
		{
			Name:        "theme: kanagawa",
			Description: "Switch to kanagawa theme (dark)",
			Action: func(app *App) {
				app.theme = KanagawaTheme
				app.saveThemePreference()
			},
		},
		{
			Name:        "theme: gruvbox-light",
			Description: "Switch to gruvbox light theme",
			Action: func(app *App) {
				app.theme = GruvboxLightTheme
				app.saveThemePreference()
			},
		},
		{
			Name:        "theme: solarized-light",
			Description: "Switch to solarized light theme",
			Action: func(app *App) {
				app.theme = SolarizedLightTheme
				app.saveThemePreference()
			},
		},
		{
			Name:        "theme: catppuccin-latte",
			Description: "Switch to catppuccin latte theme (light)",
			Action: func(app *App) {
				app.theme = CatppuccinLatteTheme
				app.saveThemePreference()
			},
		},
		{
			Name:        "text: random",
			Description: "Select a random text",
			Action: func(app *App) {
				app.selectRandomText()
			},
		},
		{
			Name:        "restart test",
			Description: "Restart the typing test with current text",
			Action: func(app *App) {
				app.restartTest()
			},
		},
		{
			Name:        "clear session",
			Description: "Clear saved session and start fresh",
			Action: func(app *App) {
				app.clearSession()
			},
		},
		{
			Name:        "quit",
			Description: "Exit the application",
			Action: func(app *App) {
				app.quit = true
			},
		},
	}

	// Add commands for each available text
	for _, text := range a.textLibrary.GetAllTexts() {
		textName := text.Name
		commands = append(commands, Command{
			Name:        fmt.Sprintf("text: %s", textName),
			Description: fmt.Sprintf("Practice with '%s'", textName),
			Action: func(app *App) {
				app.selectTextByName(textName)
			},
		})
	}

	// Add commands for each available word set
	for _, wordSet := range a.wordLibrary.GetAllWordSets() {
		wordSetName := wordSet.Name
		commands = append(commands, Command{
			Name:        fmt.Sprintf("words: %s", wordSetName),
			Description: fmt.Sprintf("Practice random words from '%s'", wordSetName),
			Action: func(app *App) {
				app.selectWordSet(wordSetName)
			},
		})
	}

	// Add time limit commands (automatically switches to time-based limit)
	commands = append(commands, Command{
		Name:        "limit: 30 seconds",
		Description: "Set time limit to 30 seconds",
		Action: func(app *App) {
			app.setTimeLimit(30)
		},
	})
	commands = append(commands, Command{
		Name:        "limit: 60 seconds",
		Description: "Set time limit to 60 seconds",
		Action: func(app *App) {
			app.setTimeLimit(60)
		},
	})
	commands = append(commands, Command{
		Name:        "limit: 120 seconds",
		Description: "Set time limit to 120 seconds",
		Action: func(app *App) {
			app.setTimeLimit(120)
		},
	})

	// Add word limit commands (automatically switches to word-based limit)
	commands = append(commands, Command{
		Name:        "limit: 50 words",
		Description: "Set word limit to 50 words",
		Action: func(app *App) {
			app.setWordLimit(50)
		},
	})
	commands = append(commands, Command{
		Name:        "limit: 100 words",
		Description: "Set word limit to 100 words",
		Action: func(app *App) {
			app.setWordLimit(100)
		},
	})
	commands = append(commands, Command{
		Name:        "limit: 200 words",
		Description: "Set word limit to 200 words",
		Action: func(app *App) {
			app.setWordLimit(200)
		},
	})

	a.commandMenu.SetCommands(commands)
}
