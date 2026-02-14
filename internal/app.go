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
	renderer     *Renderer
	typingTest   *TypingTest
	inputHandler *InputHandler
	commandMenu  *CommandMenu
	textLibrary  *TextLibrary

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
//
// Returns an error if the screen cannot be created or initialized.
func NewApp(stdinText string) (*App, error) {
	screen, err := tcell.NewScreen()
	if err != nil {
		return nil, fmt.Errorf("failed to create screen: %w", err)
	}

	if err := screen.Init(); err != nil {
		return nil, fmt.Errorf("failed to initialize screen: %w", err)
	}

	// Load text library
	textLibrary := NewTextLibrary(DefaultTextsDir)

	// Select initial text
	var initialText TextSource
	if stdinText != "" {
		stdinSource := TextSource{
			Name:    "stdin",
			Content: stdinText,
			Path:    "",
		}
		textLibrary.AddText(stdinSource)
		textLibrary.SelectByName("stdin")
		initialText = stdinSource
	} else {
		initialText = textLibrary.SelectRandom()
	}

	// Create components
	renderer := NewRenderer(screen)
	typingTest := NewTypingTest(initialText.Content)
	commandMenu := NewCommandMenu()

	app := &App{
		renderer:    renderer,
		typingTest:  typingTest,
		commandMenu: commandMenu,
		textLibrary: textLibrary,
		theme:       DefaultTheme,
		screen:      screen,
		quit:        false,
		showResults: false,
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

// cycleTheme switches to the next theme.
func (a *App) cycleTheme() {
	a.theme = GetNextTheme(a.theme)
}

// restartTest resets the current typing test.
func (a *App) restartTest() {
	a.typingTest.Reset()
	a.showResults = false
}

// selectRandomText selects a random text and restarts the test.
func (a *App) selectRandomText() {
	text := a.textLibrary.SelectRandom()
	a.typingTest.SetSampleText(text.Content)
}

// selectTextByName selects a text by name and restarts the test.
func (a *App) selectTextByName(name string) {
	if a.textLibrary.SelectByName(name) {
		text := a.textLibrary.GetCurrentText()
		a.typingTest.SetSampleText(text.Content)
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
			},
		},
		{
			Name:        "theme: gruvbox",
			Description: "Switch to gruvbox theme (dark)",
			Action: func(app *App) {
				app.theme = GruvboxTheme
			},
		},
		{
			Name:        "theme: kanagawa",
			Description: "Switch to kanagawa theme (dark)",
			Action: func(app *App) {
				app.theme = KanagawaTheme
			},
		},
		{
			Name:        "theme: gruvbox-light",
			Description: "Switch to gruvbox light theme",
			Action: func(app *App) {
				app.theme = GruvboxLightTheme
			},
		},
		{
			Name:        "theme: solarized-light",
			Description: "Switch to solarized light theme",
			Action: func(app *App) {
				app.theme = SolarizedLightTheme
			},
		},
		{
			Name:        "theme: catppuccin-latte",
			Description: "Switch to catppuccin latte theme (light)",
			Action: func(app *App) {
				app.theme = CatppuccinLatteTheme
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
