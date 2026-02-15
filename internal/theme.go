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

	// === Fancy Dark Themes ===

	// CyberpunkTheme features neon colors on a dark background, inspired by cyberpunk
	// aesthetics and synthwave art. Vibrant pinks, purples, and cyans create a futuristic feel.
	CyberpunkTheme = Theme{
		Name:           "cyberpunk",
		Background:     tcell.NewRGBColor(16, 16, 28),    // #10101c - Deep dark blue
		Foreground:     tcell.NewRGBColor(230, 230, 250), // #e6e6fa - Light lavender
		TextDefault:    tcell.NewRGBColor(150, 150, 200), // #9696c8 - Brighter purple (better contrast)
		TextCorrect:    tcell.NewRGBColor(0, 255, 255),   // #00ffff - Cyan (neon)
		TextIncorrect:  tcell.NewRGBColor(255, 0, 128),   // #ff0080 - Hot pink
		TextCursor:     tcell.NewRGBColor(255, 0, 255),   // #ff00ff - Magenta (neon)
		Title:          tcell.NewRGBColor(255, 0, 255),   // #ff00ff - Magenta (neon)
		Border:         tcell.NewRGBColor(138, 43, 226),  // #8a2be2 - Blue violet
		Help:           tcell.NewRGBColor(147, 112, 219), // #9370db - Medium purple
		MenuSelectedBg: tcell.NewRGBColor(75, 0, 130),    // #4b0082 - Indigo
		MenuSelectedFg: tcell.NewRGBColor(0, 255, 255),   // #00ffff - Cyan
		MenuDimText:    tcell.NewRGBColor(120, 120, 170), // #7878aa - Medium purple
	}

	// MidnightTheme is an elegant dark theme with deep blues and silvers,
	// reminiscent of a clear night sky. Sophisticated and easy on the eyes.
	MidnightTheme = Theme{
		Name:           "midnight",
		Background:     tcell.NewRGBColor(13, 17, 23),    // #0d1117 - Deep midnight blue
		Foreground:     tcell.NewRGBColor(201, 209, 217), // #c9d1d9 - Cool gray
		TextDefault:    tcell.NewRGBColor(139, 148, 158), // #8b949e - Brighter gray-blue (better contrast)
		TextCorrect:    tcell.NewRGBColor(125, 174, 255), // #7dafff - Soft blue
		TextIncorrect:  tcell.NewRGBColor(242, 105, 125), // #f2697d - Coral red
		TextCursor:     tcell.NewRGBColor(121, 192, 255), // #79c0ff - Bright sky blue
		Title:          tcell.NewRGBColor(121, 192, 255), // #79c0ff - Bright sky blue
		Border:         tcell.NewRGBColor(33, 38, 45),    // #21262d - Dark border
		Help:           tcell.NewRGBColor(110, 118, 129), // #6e7681 - Muted gray-blue
		MenuSelectedBg: tcell.NewRGBColor(33, 38, 45),    // #21262d - Dark selection
		MenuSelectedFg: tcell.NewRGBColor(201, 209, 217), // #c9d1d9 - Cool gray
		MenuDimText:    tcell.NewRGBColor(110, 118, 129), // #6e7681 - Muted gray-blue
	}

	// OceanDeepTheme features deep ocean blues and aqua highlights,
	// creating a calming underwater atmosphere. Rich and immersive.
	OceanDeepTheme = Theme{
		Name:           "ocean-deep",
		Background:     tcell.NewRGBColor(10, 25, 47),    // #0a192f - Deep ocean blue
		Foreground:     tcell.NewRGBColor(204, 214, 246), // #ccd6f6 - Light periwinkle
		TextDefault:    tcell.NewRGBColor(136, 146, 176), // #8892b0 - Brighter slate (better contrast)
		TextCorrect:    tcell.NewRGBColor(128, 203, 196), // #80cbc4 - Aqua
		TextIncorrect:  tcell.NewRGBColor(255, 107, 107), // #ff6b6b - Coral
		TextCursor:     tcell.NewRGBColor(100, 255, 218), // #64ffda - Bright aqua
		Title:          tcell.NewRGBColor(100, 255, 218), // #64ffda - Bright aqua
		Border:         tcell.NewRGBColor(23, 42, 69),    // #172a45 - Dark teal
		Help:           tcell.NewRGBColor(136, 146, 176), // #8892b0 - Cool gray
		MenuSelectedBg: tcell.NewRGBColor(17, 34, 64),    // #112240 - Dark selection
		MenuSelectedFg: tcell.NewRGBColor(100, 255, 218), // #64ffda - Bright aqua
		MenuDimText:    tcell.NewRGBColor(100, 115, 146), // #647392 - Muted slate
	}

	// DraculaTheme is inspired by the popular Dracula color scheme,
	// featuring a dark purple background with vibrant accent colors.
	DraculaTheme = Theme{
		Name:           "dracula",
		Background:     tcell.NewRGBColor(40, 42, 54),    // #282a36 - Dark purple-gray
		Foreground:     tcell.NewRGBColor(248, 248, 242), // #f8f8f2 - Off-white
		TextDefault:    tcell.NewRGBColor(139, 157, 195), // #8b9dc3 - Brighter blue (better contrast)
		TextCorrect:    tcell.NewRGBColor(80, 250, 123),  // #50fa7b - Bright green
		TextIncorrect:  tcell.NewRGBColor(255, 85, 85),   // #ff5555 - Bright red
		TextCursor:     tcell.NewRGBColor(241, 250, 140), // #f1fa8c - Bright yellow
		Title:          tcell.NewRGBColor(189, 147, 249), // #bd93f9 - Bright purple
		Border:         tcell.NewRGBColor(68, 71, 90),    // #44475a - Lighter purple-gray
		Help:           tcell.NewRGBColor(98, 114, 164),  // #6272a4 - Muted blue
		MenuSelectedBg: tcell.NewRGBColor(68, 71, 90),    // #44475a - Lighter purple-gray
		MenuSelectedFg: tcell.NewRGBColor(248, 248, 242), // #f8f8f2 - Off-white
		MenuDimText:    tcell.NewRGBColor(98, 114, 164),  // #6272a4 - Muted blue
	}

	// === Pastel Light Themes ===

	// LavenderDreamTheme features soft lavender and purple tones,
	// creating a gentle, dreamy aesthetic. Calming and elegant.
	LavenderDreamTheme = Theme{
		Name:           "lavender-dream",
		Background:     tcell.NewRGBColor(245, 243, 250), // #f5f3fa - Very light lavender
		Foreground:     tcell.NewRGBColor(80, 73, 90),    // #50495a - Dark purple-gray
		TextDefault:    tcell.NewRGBColor(140, 120, 160), // #8c78a0 - Darker purple (better contrast)
		TextCorrect:    tcell.NewRGBColor(100, 88, 120),  // #645878 - Medium purple
		TextIncorrect:  tcell.NewRGBColor(200, 84, 120),  // #c85478 - Mauve-pink
		TextCursor:     tcell.NewRGBColor(145, 106, 190), // #916abe - Soft violet
		Title:          tcell.NewRGBColor(145, 106, 190), // #916abe - Soft violet
		Border:         tcell.NewRGBColor(225, 218, 235), // #e1daeb - Light lavender
		Help:           tcell.NewRGBColor(150, 140, 165), // #968ca5 - Muted lavender
		MenuSelectedBg: tcell.NewRGBColor(230, 223, 240), // #e6dff0 - Pale lavender
		MenuSelectedFg: tcell.NewRGBColor(80, 73, 90),    // #50495a - Dark purple-gray
		MenuDimText:    tcell.NewRGBColor(180, 167, 194), // #b4a7c2 - Soft purple
	}

	// MintFreshTheme features refreshing mint green tones with soft accents,
	// creating a clean and energizing atmosphere.
	MintFreshTheme = Theme{
		Name:           "mint-fresh",
		Background:     tcell.NewRGBColor(243, 250, 246), // #f3faf6 - Very light mint
		Foreground:     tcell.NewRGBColor(60, 80, 70),    // #3c5046 - Dark teal-green
		TextDefault:    tcell.NewRGBColor(100, 140, 120), // #648c78 - Darker mint (better contrast)
		TextCorrect:    tcell.NewRGBColor(70, 130, 100),  // #468264 - Medium green
		TextIncorrect:  tcell.NewRGBColor(200, 90, 100),  // #c85a64 - Soft coral
		TextCursor:     tcell.NewRGBColor(80, 180, 140),  // #50b48c - Bright mint
		Title:          tcell.NewRGBColor(80, 180, 140),  // #50b48c - Bright mint
		Border:         tcell.NewRGBColor(220, 235, 228), // #dcebe4 - Light mint
		Help:           tcell.NewRGBColor(140, 165, 150), // #8ca596 - Muted mint
		MenuSelectedBg: tcell.NewRGBColor(225, 240, 233), // #e1f0e9 - Pale mint
		MenuSelectedFg: tcell.NewRGBColor(60, 80, 70),    // #3c5046 - Dark teal-green
		MenuDimText:    tcell.NewRGBColor(165, 195, 180), // #a5c3b4 - Soft mint
	}

	// PeachSoftTheme features warm peach and coral tones,
	// creating a cozy and inviting atmosphere.
	PeachSoftTheme = Theme{
		Name:           "peach-soft",
		Background:     tcell.NewRGBColor(255, 248, 242), // #fff8f2 - Very light peach
		Foreground:     tcell.NewRGBColor(90, 70, 60),    // #5a463c - Dark brown
		TextDefault:    tcell.NewRGBColor(160, 130, 110), // #a0826e - Darker tan (better contrast)
		TextCorrect:    tcell.NewRGBColor(140, 100, 80),  // #8c6450 - Medium brown
		TextIncorrect:  tcell.NewRGBColor(220, 100, 100), // #dc6464 - Soft red
		TextCursor:     tcell.NewRGBColor(255, 160, 120), // #ffa078 - Bright peach
		Title:          tcell.NewRGBColor(255, 160, 120), // #ffa078 - Bright peach
		Border:         tcell.NewRGBColor(245, 225, 210), // #f5e1d2 - Light peach
		Help:           tcell.NewRGBColor(170, 145, 130), // #aa9182 - Muted tan
		MenuSelectedBg: tcell.NewRGBColor(250, 235, 220), // #faebdc - Pale peach
		MenuSelectedFg: tcell.NewRGBColor(90, 70, 60),    // #5a463c - Dark brown
		MenuDimText:    tcell.NewRGBColor(210, 180, 165), // #d2b4a5 - Soft tan
	}

	// RosewaterTheme features soft pink and rose tones,
	// creating a gentle and romantic aesthetic.
	RosewaterTheme = Theme{
		Name:           "rosewater",
		Background:     tcell.NewRGBColor(255, 245, 248), // #fff5f8 - Very light rose
		Foreground:     tcell.NewRGBColor(80, 60, 70),    // #503c46 - Dark mauve
		TextDefault:    tcell.NewRGBColor(160, 130, 145), // #a08291 - Darker rose (better contrast)
		TextCorrect:    tcell.NewRGBColor(120, 80, 100),  // #785064 - Medium mauve
		TextIncorrect:  tcell.NewRGBColor(220, 90, 120),  // #dc5a78 - Rose red
		TextCursor:     tcell.NewRGBColor(255, 140, 170), // #ff8caa - Bright rose
		Title:          tcell.NewRGBColor(255, 140, 170), // #ff8caa - Bright rose
		Border:         tcell.NewRGBColor(245, 220, 230), // #f5dce6 - Light rose
		Help:           tcell.NewRGBColor(170, 145, 155), // #aa919b - Muted rose
		MenuSelectedBg: tcell.NewRGBColor(250, 230, 238), // #fae6ee - Pale rose
		MenuSelectedFg: tcell.NewRGBColor(80, 60, 70),    // #503c46 - Dark mauve
		MenuDimText:    tcell.NewRGBColor(210, 180, 190), // #d2b4be - Soft rose
	}

	// === High Contrast Themes ===

	// HighContrastDarkTheme features pure white text on pure black background,
	// providing maximum contrast for accessibility and reduced eye strain in dark environments.
	HighContrastDarkTheme = Theme{
		Name:           "high-contrast-dark",
		Background:     tcell.NewRGBColor(0, 0, 0),       // #000000 - Pure black
		Foreground:     tcell.NewRGBColor(255, 255, 255), // #ffffff - Pure white
		TextDefault:    tcell.NewRGBColor(140, 140, 140), // #8c8c8c - Medium gray
		TextCorrect:    tcell.NewRGBColor(255, 255, 255), // #ffffff - Pure white
		TextIncorrect:  tcell.NewRGBColor(255, 50, 50),   // #ff3232 - Bright red
		TextCursor:     tcell.NewRGBColor(255, 255, 0),   // #ffff00 - Pure yellow
		Title:          tcell.NewRGBColor(255, 255, 0),   // #ffff00 - Pure yellow
		Border:         tcell.NewRGBColor(80, 80, 80),    // #505050 - Dark gray
		Help:           tcell.NewRGBColor(180, 180, 180), // #b4b4b4 - Light gray
		MenuSelectedBg: tcell.NewRGBColor(255, 255, 255), // #ffffff - White (inverted)
		MenuSelectedFg: tcell.NewRGBColor(0, 0, 0),       // #000000 - Black (inverted)
		MenuDimText:    tcell.NewRGBColor(140, 140, 140), // #8c8c8c - Medium gray
	}

	// HighContrastLightTheme features pure black text on pure white background,
	// providing maximum contrast for accessibility in bright environments.
	HighContrastLightTheme = Theme{
		Name:           "high-contrast-light",
		Background:     tcell.NewRGBColor(255, 255, 255), // #ffffff - Pure white
		Foreground:     tcell.NewRGBColor(0, 0, 0),       // #000000 - Pure black
		TextDefault:    tcell.NewRGBColor(140, 140, 140), // #8c8c8c - Medium gray
		TextCorrect:    tcell.NewRGBColor(0, 0, 0),       // #000000 - Pure black
		TextIncorrect:  tcell.NewRGBColor(200, 0, 0),     // #c80000 - Dark red
		TextCursor:     tcell.NewRGBColor(0, 0, 200),     // #0000c8 - Dark blue
		Title:          tcell.NewRGBColor(0, 0, 200),     // #0000c8 - Dark blue
		Border:         tcell.NewRGBColor(200, 200, 200), // #c8c8c8 - Light gray
		Help:           tcell.NewRGBColor(100, 100, 100), // #646464 - Dark gray
		MenuSelectedBg: tcell.NewRGBColor(0, 0, 0),       // #000000 - Black (inverted)
		MenuSelectedFg: tcell.NewRGBColor(255, 255, 255), // #ffffff - White (inverted)
		MenuDimText:    tcell.NewRGBColor(140, 140, 140), // #8c8c8c - Medium gray
	}

	// HighVisibilityTheme features bright yellow background with black text,
	// designed for maximum visibility and attention. Bold and energetic.
	HighVisibilityTheme = Theme{
		Name:           "high-visibility",
		Background:     tcell.NewRGBColor(255, 255, 0), // #ffff00 - Bright yellow
		Foreground:     tcell.NewRGBColor(0, 0, 0),     // #000000 - Pure black
		TextDefault:    tcell.NewRGBColor(100, 100, 0), // #646400 - Dark yellow-green
		TextCorrect:    tcell.NewRGBColor(0, 0, 0),     // #000000 - Pure black
		TextIncorrect:  tcell.NewRGBColor(200, 0, 0),   // #c80000 - Dark red
		TextCursor:     tcell.NewRGBColor(0, 0, 200),   // #0000c8 - Dark blue
		Title:          tcell.NewRGBColor(0, 0, 200),   // #0000c8 - Dark blue
		Border:         tcell.NewRGBColor(200, 200, 0), // #c8c800 - Dark yellow
		Help:           tcell.NewRGBColor(80, 80, 0),   // #505000 - Olive
		MenuSelectedBg: tcell.NewRGBColor(0, 0, 0),     // #000000 - Black (inverted)
		MenuSelectedFg: tcell.NewRGBColor(255, 255, 0), // #ffff00 - Yellow (inverted)
		MenuDimText:    tcell.NewRGBColor(100, 100, 0), // #646400 - Dark yellow-green
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
		// Fancy dark themes
		CyberpunkTheme,
		MidnightTheme,
		OceanDeepTheme,
		DraculaTheme,
		// Pastel light themes
		LavenderDreamTheme,
		MintFreshTheme,
		PeachSoftTheme,
		RosewaterTheme,
		// High contrast themes
		HighContrastDarkTheme,
		HighContrastLightTheme,
		HighVisibilityTheme,
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
