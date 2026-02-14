package internal

import "strings"

// Command represents an executable action in the command palette.
// Commands can be filtered by name or description and executed with a single keystroke.
type Command struct {
	Name        string     // Display name shown in the command palette
	Description string     // Descriptive text explaining what the command does
	Action      func(*App) // Function to execute when the command is selected
}

// CommandMenu manages the command palette overlay, including visibility,
// filtering, selection, and command execution. It provides a keyboard-driven
// interface for accessing application features.
type CommandMenu struct {
	visible  bool      // Whether the command menu is currently displayed
	filter   string    // Current filter text for searching commands
	selected int       // Index of currently selected command in filtered list
	commands []Command // All available commands
}

// NewCommandMenu creates a new CommandMenu instance with default values.
// The menu starts hidden with no filter applied and no commands registered.
// Use SetCommands to populate the available commands.
func NewCommandMenu() *CommandMenu {
	return &CommandMenu{
		visible:  false,
		filter:   "",
		selected: 0,
		commands: []Command{},
	}
}

// Show displays the command menu and resets the filter and selection.
// This ensures a clean slate each time the menu is opened.
func (cm *CommandMenu) Show() {
	cm.visible = true
	cm.filter = ""
	cm.selected = 0
}

// Hide closes the command menu and clears any active filter and selection state.
// The command list is preserved for the next time the menu is opened.
func (cm *CommandMenu) Hide() {
	cm.visible = false
	cm.filter = ""
	cm.selected = 0
}

// IsVisible returns whether the command menu is currently displayed.
func (cm *CommandMenu) IsVisible() bool {
	return cm.visible
}

// SetCommands replaces the available commands with the provided list.
// This should typically be called once during application initialization.
//
// Parameters:
//   - commands: slice of Command structs to make available in the menu
func (cm *CommandMenu) SetCommands(commands []Command) {
	cm.commands = commands
}

// AddChar appends a character to the current filter string.
// Resets the selection to the first item in the newly filtered list.
//
// Parameters:
//   - ch: the character to add to the filter
func (cm *CommandMenu) AddChar(ch rune) {
	cm.filter += string(ch)
	cm.selected = 0 // Reset selection when filter changes
}

// Backspace removes the last character from the filter string.
// If the filter is already empty, this is a no-op.
// Resets selection to the first item after modifying the filter.
func (cm *CommandMenu) Backspace() {
	if len(cm.filter) > 0 {
		cm.filter = cm.filter[:len(cm.filter)-1]
		cm.selected = 0
	}
}

// GetFilter returns the current filter string being applied to commands.
func (cm *CommandMenu) GetFilter() string {
	return cm.filter
}

// GetFilteredCommands returns commands that match the current filter.
// Matching is case-insensitive and searches both command names and descriptions.
// If no filter is applied, returns all commands.
//
// Returns a slice of matching Command structs in their original order.
func (cm *CommandMenu) GetFilteredCommands() []Command {
	if cm.filter == "" {
		return cm.commands
	}

	filter := strings.ToLower(cm.filter)
	var filtered []Command

	for _, cmd := range cm.commands {
		nameMatch := strings.Contains(strings.ToLower(cmd.Name), filter)
		descMatch := strings.Contains(strings.ToLower(cmd.Description), filter)
		if nameMatch || descMatch {
			filtered = append(filtered, cmd)
		}
	}

	return filtered
}

// MoveUp moves the selection cursor up by one position.
// Does nothing if already at the first item.
func (cm *CommandMenu) MoveUp() {
	if cm.selected > 0 {
		cm.selected--
	}
}

// MoveDown moves the selection cursor down by one position.
// Does nothing if already at the last item in the filtered list.
func (cm *CommandMenu) MoveDown() {
	filtered := cm.GetFilteredCommands()
	if cm.selected < len(filtered)-1 {
		cm.selected++
	}
}

// GetSelected returns the index of the currently selected command
// within the filtered command list.
func (cm *CommandMenu) GetSelected() int {
	return cm.selected
}

// ExecuteSelected executes the currently selected command and closes the menu.
// If no commands match the filter or selection is invalid, this is a no-op.
// The menu is automatically hidden after successful command execution.
//
// Parameters:
//   - app: the App instance to pass to the command's action function
func (cm *CommandMenu) ExecuteSelected(app *App) {
	filtered := cm.GetFilteredCommands()
	if len(filtered) > 0 && cm.selected < len(filtered) {
		filtered[cm.selected].Action(app)
		cm.Hide()
	}
}
