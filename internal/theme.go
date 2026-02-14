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

	// GruvboxLightTheme is the light variant of Gruvbox with warm, retro-inspired tones.
	// Features a cream background with earthy accents. Perfect for daytime use.
	GruvboxLightTheme = Theme{
		Name:           "gruvbox-light",
		Background:     tcell.NewRGBColor(251, 241, 199), // #fbf1c7 - Light cream
		Foreground:     tcell.NewRGBColor(60, 56, 54),    // #3c3836 - Dark brown-gray
		TextDefault:    tcell.NewRGBColor(189, 174, 147), // #bdae93 - Medium tan
		TextCorrect:    tcell.NewRGBColor(60, 56, 54),    // #3c3836 - Dark brown-gray
		TextIncorrect:  tcell.NewRGBColor(204, 36, 29),   // #cc241d - Deep red
		TextCursor:     tcell.NewRGBColor(181, 118, 20),  // #b57614 - Dark golden
		Title:          tcell.NewRGBColor(181, 118, 20),  // #b57614 - Dark golden
		Border:         tcell.NewRGBColor(213, 196, 161), // #d5c4a1 - Light tan
		Help:           tcell.NewRGBColor(146, 131, 116), // #928374 - Medium gray-brown
		MenuSelectedBg: tcell.NewRGBColor(235, 219, 178), // #ebdbb2 - Slightly darker cream
		MenuSelectedFg: tcell.NewRGBColor(60, 56, 54),    // #3c3836 - Dark brown-gray
		MenuDimText:    tcell.NewRGBColor(189, 174, 147), // #bdae93 - Medium tan
	}

	// SolarizedLightTheme is based on the precise color relationships of the Solarized
	// palette by Ethan Schoonover. Features carefully balanced colors with optimal
	// contrast ratios for reduced eye strain.
	SolarizedLightTheme = Theme{
		Name:           "solarized-light",
		Background:     tcell.NewRGBColor(253, 246, 227), // #fdf6e3 - Base3 (light background)
		Foreground:     tcell.NewRGBColor(101, 123, 131), // #657b83 - Base00 (body text)
		TextDefault:    tcell.NewRGBColor(147, 161, 161), // #93a1a1 - Base1 (optional emphasized)
		TextCorrect:    tcell.NewRGBColor(88, 110, 117),  // #586e75 - Base01 (emphasized)
		TextIncorrect:  tcell.NewRGBColor(220, 50, 47),   // #dc322f - Red
		TextCursor:     tcell.NewRGBColor(203, 75, 22),   // #cb4b16 - Orange
		Title:          tcell.NewRGBColor(181, 137, 0),   // #b58900 - Yellow
		Border:         tcell.NewRGBColor(238, 232, 213), // #eee8d5 - Base2 (background highlights)
		Help:           tcell.NewRGBColor(147, 161, 161), // #93a1a1 - Base1
		MenuSelectedBg: tcell.NewRGBColor(238, 232, 213), // #eee8d5 - Base2
		MenuSelectedFg: tcell.NewRGBColor(88, 110, 117),  // #586e75 - Base01
		MenuDimText:    tcell.NewRGBColor(147, 161, 161), // #93a1a1 - Base1
	}

	// CatppuccinLatteTheme is the light variant of Catppuccin, featuring pastel colors
	// with a soft, warm latte aesthetic. Modern and pleasing to the eye.
	CatppuccinLatteTheme = Theme{
		Name:           "catppuccin-latte",
		Background:     tcell.NewRGBColor(239, 241, 245), // #eff1f5 - Base
		Foreground:     tcell.NewRGBColor(76, 79, 105),   // #4c4f69 - Text
		TextDefault:    tcell.NewRGBColor(156, 160, 176), // #9ca0b0 - Overlay0
		TextCorrect:    tcell.NewRGBColor(76, 79, 105),   // #4c4f69 - Text
		TextIncorrect:  tcell.NewRGBColor(210, 15, 57),   // #d20f39 - Red
		TextCursor:     tcell.NewRGBColor(254, 100, 11),  // #fe640b - Peach
		Title:          tcell.NewRGBColor(223, 142, 29),  // #df8e1d - Yellow
		Border:         tcell.NewRGBColor(220, 224, 232), // #dce0e8 - Mantle
		Help:           tcell.NewRGBColor(140, 143, 161), // #8c8fa1 - Subtext0
		MenuSelectedBg: tcell.NewRGBColor(204, 208, 218), // #ccd0da - Surface0
		MenuSelectedFg: tcell.NewRGBColor(76, 79, 105),   // #4c4f69 - Text
		MenuDimText:    tcell.NewRGBColor(156, 160, 176), // #9ca0b0 - Overlay0
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
		GruvboxLightTheme,
		SolarizedLightTheme,
		CatppuccinLatteTheme,
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
