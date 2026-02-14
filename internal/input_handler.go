package internal

import "github.com/gdamore/tcell/v2"

// AppMode represents the current mode of the application.
type AppMode int

const (
	// ModeTyping is the normal typing test mode.
	ModeTyping AppMode = iota
	// ModeResults is the results screen mode.
	ModeResults
	// ModeCommandMenu is when the command menu is visible.
	ModeCommandMenu
)

// InputHandler handles keyboard input routing based on application mode.
// It separates input handling logic from the main application controller.
type InputHandler struct {
	// Callbacks for different actions
	onQuit              func()
	onToggleCommandMenu func()
	onCycleTheme        func()
	onRestartTest       func()

	// Mode-specific handlers
	typingHandler      *TypingInputHandler
	resultsHandler     *ResultsInputHandler
	commandMenuHandler *CommandMenuInputHandler
}

// NewInputHandler creates a new input handler with the given callbacks.
func NewInputHandler(
	onQuit func(),
	onToggleCommandMenu func(),
	onCycleTheme func(),
	onRestartTest func(),
	typingTest *TypingTest,
	commandMenu *CommandMenu,
) *InputHandler {
	return &InputHandler{
		onQuit:              onQuit,
		onToggleCommandMenu: onToggleCommandMenu,
		onCycleTheme:        onCycleTheme,
		onRestartTest:       onRestartTest,
		typingHandler:       NewTypingInputHandler(typingTest),
		resultsHandler:      NewResultsInputHandler(),
		commandMenuHandler:  NewCommandMenuInputHandler(commandMenu),
	}
}

// HandleKey routes keyboard events to the appropriate handler based on mode.
func (h *InputHandler) HandleKey(ev *tcell.EventKey, mode AppMode) {
	switch mode {
	case ModeCommandMenu:
		h.handleCommandMenuKey(ev)
	case ModeResults:
		h.handleResultsKey(ev)
	case ModeTyping:
		h.handleTypingKey(ev)
	}
}

// handleTypingKey processes input during typing mode.
func (h *InputHandler) handleTypingKey(ev *tcell.EventKey) {
	switch ev.Key() {
	case tcell.KeyEscape, tcell.KeyCtrlC:
		h.onQuit()
	case tcell.KeyCtrlP:
		h.onToggleCommandMenu()
	case tcell.KeyCtrlT:
		h.onCycleTheme()
	case tcell.KeyBackspace, tcell.KeyBackspace2:
		h.typingHandler.HandleBackspace()
	case tcell.KeyEnter:
		h.typingHandler.HandleEnter()
	case tcell.KeyRune:
		h.typingHandler.HandleRune(ev.Rune())
	}
}

// handleResultsKey processes input during results screen mode.
func (h *InputHandler) handleResultsKey(ev *tcell.EventKey) {
	switch ev.Key() {
	case tcell.KeyEscape, tcell.KeyCtrlC:
		h.onQuit()
	case tcell.KeyCtrlP:
		h.onToggleCommandMenu()
	case tcell.KeyCtrlT:
		h.onCycleTheme()
	case tcell.KeyEnter, tcell.KeyRune:
		if ev.Rune() == 'r' || ev.Key() == tcell.KeyEnter {
			h.onRestartTest()
		}
	}
}

// handleCommandMenuKey processes input when command menu is visible.
func (h *InputHandler) handleCommandMenuKey(ev *tcell.EventKey) {
	switch ev.Key() {
	case tcell.KeyEscape, tcell.KeyCtrlC, tcell.KeyCtrlP:
		h.onToggleCommandMenu()
	case tcell.KeyUp, tcell.KeyCtrlK:
		h.commandMenuHandler.HandleMoveUp()
	case tcell.KeyDown, tcell.KeyCtrlJ:
		h.commandMenuHandler.HandleMoveDown()
	case tcell.KeyEnter:
		h.commandMenuHandler.HandleExecute()
		h.onToggleCommandMenu() // Close menu after execution
	case tcell.KeyBackspace, tcell.KeyBackspace2:
		h.commandMenuHandler.HandleBackspace()
	case tcell.KeyRune:
		h.commandMenuHandler.HandleRune(ev.Rune())
	}
}

// TypingInputHandler handles input during typing mode.
type TypingInputHandler struct {
	test *TypingTest
}

// NewTypingInputHandler creates a new typing input handler.
func NewTypingInputHandler(test *TypingTest) *TypingInputHandler {
	return &TypingInputHandler{test: test}
}

// HandleRune handles a regular character input.
func (h *TypingInputHandler) HandleRune(r rune) {
	h.test.TypeCharacter(r)
}

// HandleEnter handles the Enter key (newline).
func (h *TypingInputHandler) HandleEnter() {
	h.test.TypeNewline()
}

// HandleBackspace handles the Backspace key.
func (h *TypingInputHandler) HandleBackspace() {
	h.test.Backspace()
}

// ResultsInputHandler handles input during results screen mode.
type ResultsInputHandler struct{}

// NewResultsInputHandler creates a new results input handler.
func NewResultsInputHandler() *ResultsInputHandler {
	return &ResultsInputHandler{}
}

// CommandMenuInputHandler handles input when command menu is visible.
type CommandMenuInputHandler struct {
	menu *CommandMenu
}

// NewCommandMenuInputHandler creates a new command menu input handler.
func NewCommandMenuInputHandler(menu *CommandMenu) *CommandMenuInputHandler {
	return &CommandMenuInputHandler{menu: menu}
}

// HandleMoveUp moves selection up in the command menu.
func (h *CommandMenuInputHandler) HandleMoveUp() {
	h.menu.MoveUp()
}

// HandleMoveDown moves selection down in the command menu.
func (h *CommandMenuInputHandler) HandleMoveDown() {
	h.menu.MoveDown()
}

// HandleExecute executes the selected command.
func (h *CommandMenuInputHandler) HandleExecute() {
	// Note: Actual execution needs to be done by the app
	// This just signals that execute was requested
}

// HandleBackspace handles backspace in the filter input.
func (h *CommandMenuInputHandler) HandleBackspace() {
	h.menu.Backspace()
}

// HandleRune handles character input for filtering.
func (h *CommandMenuInputHandler) HandleRune(r rune) {
	h.menu.AddChar(r)
}
