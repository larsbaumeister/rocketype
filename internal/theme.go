package internal

import "github.com/gdamore/tcell/v2"

// Theme defines the complete color scheme for the application.
// All visual elements should reference colors from the active theme to ensure
// consistency and allow for easy theme switching.
//
// Design principles:
//   - No hardcoded colors in rendering code
//   - Full-screen background color support
//   - All UI elements have defined colors
//   - Themes should be visually distinct and accessible
type Theme struct {
	// Name is the unique identifier for the theme
	Name string

	// Base colors
	Background tcell.Color // Full-screen background color
	Foreground tcell.Color // Default text color for general UI elements

	// Text state colors for the typing area
	TextDefault   tcell.Color // Untyped characters (not yet reached)
	TextCorrect   tcell.Color // Correctly typed characters
	TextIncorrect tcell.Color // Incorrectly typed characters
	TextCursor    tcell.Color // Character at cursor position

	// UI component colors
	Title  tcell.Color // Title and header text
	Border tcell.Color // Borders and separators
	Help   tcell.Color // Help text and secondary information

	// Command menu colors
	MenuSelectedBg tcell.Color // Background of selected menu item
	MenuSelectedFg tcell.Color // Foreground of selected menu item
	MenuDimText    tcell.Color // Dimmed/disabled text in menus
}

var (
	// DefaultTheme uses the terminal's default color palette, respecting user preferences.
	// This theme adapts to the user's terminal configuration and works well in both
	// light and dark terminal themes.
	DefaultTheme = Theme{
		Name:           "default",
		Background:     tcell.ColorDefault,
		Foreground:     tcell.ColorDefault,
		TextDefault:    tcell.ColorDefault,
		TextCorrect:    tcell.ColorGreen,
		TextIncorrect:  tcell.ColorRed,
		TextCursor:     tcell.ColorYellow,
		Title:          tcell.ColorYellow,
		Border:         tcell.ColorDarkGray,
		Help:           tcell.ColorDarkGray,
		MenuSelectedBg: tcell.ColorDarkGray,
		MenuSelectedFg: tcell.ColorWhite,
		MenuDimText:    tcell.ColorGray,
	}

	// GruvboxTheme is a warm, retro-inspired color scheme with earthy tones.
	// Based on the popular Gruvbox color scheme by Pavel Pertsev.
	// Provides excellent contrast and is easy on the eyes for extended typing sessions.
	GruvboxTheme = Theme{
		Name:           "gruvbox",
		Background:     tcell.NewRGBColor(40, 40, 40),    // #282828 - Dark warm gray
		Foreground:     tcell.NewRGBColor(235, 219, 178), // #ebdbb2 - Light cream
		TextDefault:    tcell.NewRGBColor(146, 131, 116), // #928374 - Medium gray
		TextCorrect:    tcell.NewRGBColor(235, 219, 178), // #ebdbb2 - Light cream
		TextIncorrect:  tcell.NewRGBColor(251, 73, 52),   // #fb4934 - Bright red
		TextCursor:     tcell.NewRGBColor(250, 189, 47),  // #fabd2f - Warm yellow
		Title:          tcell.NewRGBColor(250, 189, 47),  // #fabd2f - Warm yellow
		Border:         tcell.NewRGBColor(80, 73, 69),    // #504945 - Dark brown-gray
		Help:           tcell.NewRGBColor(146, 131, 116), // #928374 - Medium gray
		MenuSelectedBg: tcell.NewRGBColor(60, 56, 54),    // #3c3836 - Slightly lighter than bg
		MenuSelectedFg: tcell.NewRGBColor(235, 219, 178), // #ebdbb2 - Light cream
		MenuDimText:    tcell.NewRGBColor(146, 131, 116), // #928374 - Medium gray
	}

	// KanagawaTheme is inspired by traditional Japanese painting and wave aesthetics.
	// Features deep, rich colors with excellent contrast. Named after "The Great Wave
	// off Kanagawa" by Hokusai.
	KanagawaTheme = Theme{
		Name:           "kanagawa",
		Background:     tcell.NewRGBColor(31, 31, 40),    // #1F1F28 - Deep navy
		Foreground:     tcell.NewRGBColor(220, 215, 186), // #DCD7BA - Soft beige
		TextDefault:    tcell.NewRGBColor(114, 113, 105), // #727169 - Muted gray
		TextCorrect:    tcell.NewRGBColor(220, 215, 186), // #DCD7BA - Soft beige
		TextIncorrect:  tcell.NewRGBColor(232, 36, 36),   // #E82424 - Vibrant red
		TextCursor:     tcell.NewRGBColor(255, 158, 59),  // #FF9E3B - Warm orange
		Title:          tcell.NewRGBColor(255, 158, 59),  // #FF9E3B - Warm orange
		Border:         tcell.NewRGBColor(84, 84, 109),   // #54546D - Cool gray
		Help:           tcell.NewRGBColor(114, 113, 105), // #727169 - Muted gray
		MenuSelectedBg: tcell.NewRGBColor(54, 54, 68),    // #363644 - Slightly lighter than bg
		MenuSelectedFg: tcell.NewRGBColor(220, 215, 186), // #DCD7BA - Soft beige
		MenuDimText:    tcell.NewRGBColor(114, 113, 105), // #727169 - Muted gray
	}
)

// AvailableThemes returns all available themes in the order they appear in the theme cycle.
// This function is the single source of truth for theme ordering and can be extended
// by adding new themes to the returned slice.
func AvailableThemes() []Theme {
	return []Theme{
		DefaultTheme,
		GruvboxTheme,
		KanagawaTheme,
	}
}

// GetNextTheme returns the next theme in the rotation cycle.
// When the last theme is reached, it wraps around to the first theme.
// If the current theme is not found, returns DefaultTheme as a safe fallback.
//
// Parameters:
//   - current: the currently active theme
//
// Returns the next theme in sequence, or DefaultTheme if current is not found.
func GetNextTheme(current Theme) Theme {
	themes := AvailableThemes()
	for i, theme := range themes {
		if theme.Name == current.Name {
			return themes[(i+1)%len(themes)]
		}
	}
	return DefaultTheme
}
