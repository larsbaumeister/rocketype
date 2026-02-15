package internal

import (
	"fmt"

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
	sessionManager  *SessionManager
	settingsManager *SettingsManager

	// State
	theme       Theme
	screen      tcell.Screen
	quit        bool
	showResults bool
}

const (
	// defaultSampleText is the fallback text when no texts directory exists.
	defaultSampleText = "Roads go ever ever on,\nOver rock and under tree,\nBy caves where never sun has shone,\nBy streams that never find the sea;\nOver snow by winter sown,\nAnd through the merry flowers of June,\nOver grass and over stone,\nAnd under mountains in the moon."
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

	// Try to restore session if requested and available (unless stdin is provided)
	var initialText TextSource
	var typingTest *TypingTest

	// stdin text takes precedence over session restoration
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
			// Session loading failed, use random text
			initialText = textLibrary.SelectRandom()
			typingTest = NewTypingTest(initialText.Content)
		}
	} else {
		// No stdin, no session - use random text
		initialText = textLibrary.SelectRandom()
		typingTest = NewTypingTest(initialText.Content)
	}

	// Create components
	renderer := NewRenderer(screen)
	commandMenu := NewCommandMenu()

	app := &App{
		renderer:        renderer,
		typingTest:      typingTest,
		commandMenu:     commandMenu,
		textLibrary:     textLibrary,
		sessionManager:  sessionManager,
		settingsManager: settingsManager,
		theme:           initialTheme,
		screen:          screen,
		quit:            false,
		showResults:     false,
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

	for !a.quit {
		ev := a.screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			a.screen.Sync()
			a.draw()

		case *tcell.EventKey:
			a.handleKey(ev)
			a.draw()
		}
	}

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

	// Always save settings (theme preference persists regardless of session state)
	settings := Settings{
		ThemeName: a.theme.Name,
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

	// Draw title
	currentText := a.textLibrary.GetCurrentText()
	a.renderer.DrawTitle(a.theme.Name, currentText.Name, a.theme)

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

	// Get cached rune slices (no conversion needed!)
	sampleRunes := a.typingTest.GetSampleRunes()
	cursorPos := a.typingTest.GetCursorPos()

	// Calculate which line the cursor is on (use string for wrapping)
	sampleText := a.typingTest.GetSampleText()
	cursorLine := CalculateCursorLine(sampleText, cursorPos, maxWidth)

	// Calculate total wrapped lines
	lines := wrapText(sampleText, maxWidth)
	totalLines := len(lines)

	// Calculate scroll position to keep cursor visible
	scrollLine := CalculateScrollLine(cursorLine, maxVisibleLines, totalLines)

	// Draw typing view with cached rune slices
	viewData := TypingViewData{
		SampleText:  sampleText,
		SampleRunes: sampleRunes,
		UserInput:   a.typingTest.GetUserInput(),
		UserRunes:   a.typingTest.GetUserRunes(),
		CursorPos:   cursorPos,
		ScrollLine:  scrollLine,
		Theme:       a.theme,
	}
	a.renderer.DrawTypingView(viewData)

	// Draw stats
	stats := a.typingTest.GetStats()
	a.renderer.DrawStats(stats.GetWPM(), stats.GetAccuracy(), a.theme)

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
		ThemeName: a.theme.Name,
	}
	_ = a.settingsManager.SaveSettings(settings)
}

// restartTest resets the current typing test.
func (a *App) restartTest() {
	a.typingTest.Reset()
	a.showResults = false
	// Clear saved session when explicitly restarting
	_ = a.sessionManager.ClearSession()
}

// clearSession clears the saved session file and resets text/progress to defaults.
// Non-text-related settings like theme are preserved.
func (a *App) clearSession() {
	if err := a.sessionManager.ClearSession(); err != nil {
		// In a real app you might want to show an error message
		_ = err
	}

	// Reset to a fresh test with random text
	// Theme and other preferences are preserved
	text := a.textLibrary.SelectRandom()
	a.typingTest.SetSampleText(text.Content)
	a.showResults = false
}

// selectRandomText selects a random text and restarts the test.
func (a *App) selectRandomText() {
	text := a.textLibrary.SelectRandom()
	a.typingTest.SetSampleText(text.Content)
	// Clear saved session when selecting new text
	_ = a.sessionManager.ClearSession()
}

// selectTextByName selects a text by name and restarts the test.
func (a *App) selectTextByName(name string) {
	if a.textLibrary.SelectByName(name) {
		text := a.textLibrary.GetCurrentText()
		a.typingTest.SetSampleText(text.Content)
		// Clear saved session when selecting new text
		_ = a.sessionManager.ClearSession()
	}
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

	a.commandMenu.SetCommands(commands)
}
