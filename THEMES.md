# Rocketype Themes

This document provides an overview of all available themes in rocketype.

## How to Switch Themes

- **Keyboard shortcut**: Press `Ctrl+T` to cycle through themes
- **Command menu**: Press `Ctrl+P` and type `theme:` to see all available themes

---

## Dark Themes

### 1. Default Theme
**Name**: `default`  
**Description**: Uses your terminal's default color palette, adapting to your terminal configuration.

- Background: Terminal default
- Text: Green (correct), Red (incorrect), Yellow (cursor)
- Best for: Users who want the app to respect their terminal theme

---

### 2. Gruvbox Dark
**Name**: `gruvbox`  
**Description**: Warm, retro-inspired color scheme with earthy tones. Based on the popular Gruvbox palette by Pavel Pertsev.

**Colors**:
- Background: `#282828` - Dark warm gray
- Foreground: `#ebdbb2` - Light cream
- Correct text: `#ebdbb2` - Light cream
- Incorrect text: `#fb4934` - Bright red
- Cursor: `#fabd2f` - Warm yellow

**Best for**: Extended typing sessions, retro aesthetic fans

---

### 3. Kanagawa
**Name**: `kanagawa`  
**Description**: Inspired by traditional Japanese painting and "The Great Wave off Kanagawa" by Hokusai. Features deep, rich colors.

**Colors**:
- Background: `#1F1F28` - Deep navy
- Foreground: `#DCD7BA` - Soft beige
- Correct text: `#DCD7BA` - Soft beige
- Incorrect text: `#E82424` - Vibrant red
- Cursor: `#FF9E3B` - Warm orange

**Best for**: Users who prefer darker, cooler tones with excellent contrast

---

## Light Themes

### 4. Gruvbox Light
**Name**: `gruvbox-light`  
**Description**: Light variant of Gruvbox with warm, retro-inspired tones and a cream background.

**Colors**:
- Background: `#fbf1c7` - Light cream
- Foreground: `#3c3836` - Dark brown-gray
- Correct text: `#3c3836` - Dark brown-gray
- Incorrect text: `#cc241d` - Deep red
- Cursor: `#b57614` - Dark golden

**Best for**: Daytime use, well-lit environments, users who prefer warm light themes

---

### 5. Solarized Light
**Name**: `solarized-light`  
**Description**: Based on Ethan Schoonover's Solarized palette with scientifically balanced colors and optimal contrast ratios.

**Colors**:
- Background: `#fdf6e3` - Base3 (light beige)
- Foreground: `#657b83` - Base00 (gray-blue)
- Correct text: `#586e75` - Base01 (dark gray-blue)
- Incorrect text: `#dc322f` - Red
- Cursor: `#cb4b16` - Orange
- Title: `#b58900` - Yellow

**Best for**: Users who want scientifically designed colors for reduced eye strain

---

### 6. Catppuccin Latte
**Name**: `catppuccin-latte`  
**Description**: Light variant of the modern Catppuccin palette featuring soft pastel colors with a warm latte aesthetic.

**Colors**:
- Background: `#eff1f5` - Base (light lavender-gray)
- Foreground: `#4c4f69` - Text (dark slate)
- Correct text: `#4c4f69` - Text (dark slate)
- Incorrect text: `#d20f39` - Red
- Cursor: `#fe640b` - Peach
- Title: `#df8e1d` - Yellow

**Best for**: Modern aesthetic, users who prefer softer colors, Gen Z appeal

---

## Theme Comparison

| Theme | Type | Background | Best Use Case |
|-------|------|------------|---------------|
| **default** | Adaptive | Terminal default | Respect terminal settings |
| **gruvbox** | Dark | Warm dark gray | Long sessions, retro vibes |
| **kanagawa** | Dark | Deep navy | Cool dark preference |
| **gruvbox-light** | Light | Cream | Daytime, warm preference |
| **solarized-light** | Light | Light beige | Scientific color balance |
| **catppuccin-latte** | Light | Lavender-gray | Modern, soft aesthetics |

---

## Tips for Choosing a Theme

1. **Environment lighting**: 
   - Bright room → Light themes
   - Dim/dark room → Dark themes

2. **Time of day**:
   - Morning/afternoon → Light themes may be easier on the eyes
   - Evening/night → Dark themes reduce eye strain

3. **Personal preference**:
   - Warm colors → Gruvbox (light or dark)
   - Cool colors → Kanagawa, Solarized Light
   - Modern/pastel → Catppuccin Latte
   - Classic → Default

4. **Extended sessions**: All themes are designed for comfortable long-duration typing

---

## Adding Custom Themes

To add your own theme, edit `internal/theme.go`:

1. Define a new `Theme` struct with your colors
2. Add it to the `AvailableThemes()` function
3. Add a command in `internal/app.go`'s `initCommands()` function

Example:
```go
MyCustomTheme = Theme{
    Name:           "my-theme",
    Background:     tcell.NewRGBColor(255, 255, 255),
    Foreground:     tcell.NewRGBColor(0, 0, 0),
    // ... other colors
}
```
