package internal

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

// Renderer handles all UI rendering logic for the rocketype application.
// It separates drawing concerns from business logic, providing a clean interface
// for rendering the typing test UI, command menu, and results screen.
type Renderer struct {
	screen tcell.Screen
}

// NewRenderer creates a new Renderer instance with the given screen.
func NewRenderer(screen tcell.Screen) *Renderer {
	return &Renderer{
		screen: screen,
	}
}

// Clear clears the entire screen.
func (r *Renderer) Clear() {
	r.screen.Clear()
}

// Show updates the physical screen with buffered changes.
func (r *Renderer) Show() {
	r.screen.Show()
}

// Size returns the current screen dimensions.
func (r *Renderer) Size() (width, height int) {
	return r.screen.Size()
}

// FillBackground fills the entire screen with the given background color.
func (r *Renderer) FillBackground(bg tcell.Color) {
	width, height := r.screen.Size()
	style := tcell.StyleDefault.Background(bg)
	for y := range height {
		for x := range width {
			r.screen.SetContent(x, y, ' ', nil, style)
		}
	}
}

// DrawText renders a string at the specified coordinates with the given colors.
func (r *Renderer) DrawText(x, y int, text string, fg, bg tcell.Color) {
	style := tcell.StyleDefault.Foreground(fg).Background(bg)
	for i, ch := range text {
		r.screen.SetContent(x+i, y, ch, nil, style)
	}
}

// DrawTitle renders the title bar with theme and text information.
func (r *Renderer) DrawTitle(themeName, textName string, theme Theme) {
	width, _ := r.screen.Size()
	title := fmt.Sprintf("rocketype [%s] - %s", themeName, textName)
	x := width/2 - len(title)/2
	r.DrawText(x, 2, title, theme.Title, theme.Background)
}

// DrawHelpText renders the help text at the bottom of the screen.
func (r *Renderer) DrawHelpText(theme Theme) {
	width, height := r.screen.Size()
	help := "Esc/Ctrl+C: quit  |  Ctrl+P: command menu  |  Ctrl+T: change theme"
	x := width/2 - len(help)/2
	r.DrawText(x, height-2, help, theme.Help, theme.Background)
}

// DrawStats renders the live statistics (WPM and accuracy) at the bottom.
func (r *Renderer) DrawStats(wpm, accuracy float64, theme Theme) {
	width, height := r.screen.Size()
	statsText := fmt.Sprintf("WPM: %.0f  |  Accuracy: %.1f%%", wpm, accuracy)
	x := width/2 - len(statsText)/2
	r.DrawText(x, height-3, statsText, theme.Help, theme.Background)
}

// TypingViewData contains all data needed to render the typing test view.
type TypingViewData struct {
	SampleText  string
	SampleRunes []rune // Cached rune slice to avoid repeated conversions
	UserInput   string
	UserRunes   []rune // Cached rune slice to avoid repeated conversions
	CursorPos   int
	ScrollLine  int // Which wrapped line should be at the top of the viewport
	Theme       Theme
}

// DrawTypingView renders the main typing test interface with wrapped text and visual feedback.
func (r *Renderer) DrawTypingView(data TypingViewData) {
	width, height := r.screen.Size()

	// Calculate available space for text
	maxWidth := width - 8
	if maxWidth < 20 {
		maxWidth = width
	}

	// Wrap text to fit screen width
	lines := wrapText(data.SampleText, maxWidth)

	// Calculate available height for text lines
	availableHeight := height - 8
	maxVisibleLines := availableHeight / 2 // 2 screen rows per text line

	// Adjust scroll position if needed
	scrollLine := data.ScrollLine
	if scrollLine < 0 {
		scrollLine = 0
	}
	if scrollLine > len(lines)-1 {
		scrollLine = len(lines) - 1
	}

	// Start drawing from top of available area
	startY := 5

	// Center horizontally by finding the longest line
	maxLineLen := 0
	for _, line := range lines {
		lineLen := 0
		for _, ch := range line {
			if ch != '\n' {
				lineLen++
			}
		}
		if lineLen > maxLineLen {
			maxLineLen = lineLen
		}
	}

	startX := (width - maxLineLen) / 2
	if startX < 0 {
		startX = 2
	}

	r.drawTypingText(lines, startX, startY, height, scrollLine, maxVisibleLines, data)
}

// drawTypingText renders each character of the typing test with appropriate styling.
func (r *Renderer) drawTypingText(lines []string, startX, startY, height, scrollLine, maxVisibleLines int, data TypingViewData) {
	currentY := startY
	charIndex := 0

	// Use the cached rune slices from data (no conversion needed!)
	sampleRunes := data.SampleRunes
	userRunes := data.UserRunes

	// Calculate the end line to render
	endLine := scrollLine + maxVisibleLines
	if endLine > len(lines) {
		endLine = len(lines)
	}

	// Skip characters before scrollLine
	for lineIdx := 0; lineIdx < scrollLine && lineIdx < len(lines); lineIdx++ {
		charIndex += len([]rune(lines[lineIdx]))
	}

	// Render only visible lines
	for lineIdx := scrollLine; lineIdx < endLine; lineIdx++ {
		line := lines[lineIdx]
		currentX := startX

		for _, ch := range line {
			if charIndex >= len(sampleRunes) {
				break
			}

			if currentY >= height-4 {
				break
			}

			style, displayChar := r.getCharStyle(charIndex, ch, sampleRunes, userRunes, data)

			// Draw mistyped character above if incorrect
			if charIndex < len(userRunes) && userRunes[charIndex] != ch {
				r.drawMistypedChar(currentX, currentY-1, userRunes[charIndex], data.Theme)
			}

			// Draw the character
			if ch != '\n' {
				r.screen.SetContent(currentX, currentY, displayChar, nil, style)
				currentX++
			} else {
				r.screen.SetContent(currentX, currentY, displayChar, nil, style)
			}

			charIndex++
		}

		currentY += 2
	}
}

// getCharStyle determines the style and display character for a given position.
func (r *Renderer) getCharStyle(charIndex int, ch rune, sampleRunes, userRunes []rune, data TypingViewData) (tcell.Style, rune) {
	displayChar := ch
	var style tcell.Style

	if charIndex < len(userRunes) {
		// Already typed
		if userRunes[charIndex] == ch {
			// Correct
			style = tcell.StyleDefault.Foreground(data.Theme.TextCorrect).Background(data.Theme.Background)
		} else {
			// Incorrect
			style = tcell.StyleDefault.Foreground(data.Theme.TextIncorrect).Background(data.Theme.Background).Bold(true)
			if ch == ' ' {
				displayChar = '_'
			} else if ch == '\n' {
				displayChar = '↵'
			}
		}
	} else if charIndex == data.CursorPos {
		// Cursor position
		style = tcell.StyleDefault.Foreground(data.Theme.TextCursor).Background(data.Theme.Background).Underline(true).Bold(true)
		if ch == '\n' {
			displayChar = '↵'
		}
	} else {
		// Not yet typed
		style = tcell.StyleDefault.Foreground(data.Theme.TextDefault).Background(data.Theme.Background)
		if ch == '\n' {
			displayChar = '↵'
		}
	}

	return style, displayChar
}

// drawMistypedChar renders a mistyped character above the expected character.
func (r *Renderer) drawMistypedChar(x, y int, mistypedChar rune, theme Theme) {
	if mistypedChar == ' ' {
		mistypedChar = '_'
	} else if mistypedChar == '\n' {
		mistypedChar = '↵'
	}
	style := tcell.StyleDefault.Foreground(theme.TextIncorrect).Background(theme.Background).Dim(true)
	r.screen.SetContent(x, y, mistypedChar, nil, style)
}

// CommandMenuData contains all data needed to render the command menu.
type CommandMenuData struct {
	Filter           string
	FilteredCommands []Command
	Selected         int
	Theme            Theme
}

// DrawCommandMenu renders the command palette overlay.
func (r *Renderer) DrawCommandMenu(data CommandMenuData) {
	width, height := r.screen.Size()

	menuWidth := min(width*2/3, 60)
	menuHeight := min(height*2/3, 15)
	menuX := (width - menuWidth) / 2
	menuY := (height - menuHeight) / 2

	r.drawBox(menuX, menuY, menuWidth, menuHeight, data.Theme)
	r.drawBoxTitle(menuX, menuY, menuWidth, " command menu ", data.Theme)
	r.drawFilterInput(menuX, menuY, data.Filter, data.Theme)
	r.drawCommandList(menuX, menuY, menuWidth, menuHeight, data)
}

// ResultsData contains all data needed to render the results screen.
type ResultsData struct {
	WPM             float64
	Accuracy        float64
	MisspelledWords []string
	WordCounts      map[string]int
	Theme           Theme
}

// DrawResults renders the results screen overlay.
func (r *Renderer) DrawResults(data ResultsData) {
	width, height := r.screen.Size()

	boxWidth := min(width*3/4, 70)
	boxHeight := min(height*3/4, 25)
	boxX := (width - boxWidth) / 2
	boxY := (height - boxHeight) / 2

	r.drawBox(boxX, boxY, boxWidth, boxHeight, data.Theme)
	r.drawBoxTitle(boxX, boxY, boxWidth, " Test Complete! ", data.Theme)
	r.drawResultsContent(boxX, boxY, boxWidth, boxHeight, data)
}

// drawBox draws a bordered box at the specified position.
func (r *Renderer) drawBox(x, y, width, height int, theme Theme) {
	borderStyle := tcell.StyleDefault.Foreground(theme.Border).Background(theme.Background)
	bgStyle := tcell.StyleDefault.Background(theme.Background)

	// Top border
	r.screen.SetContent(x, y, '┌', nil, borderStyle)
	for i := x + 1; i < x+width-1; i++ {
		r.screen.SetContent(i, y, '─', nil, borderStyle)
	}
	r.screen.SetContent(x+width-1, y, '┐', nil, borderStyle)

	// Side borders and fill
	for j := y + 1; j < y+height-1; j++ {
		r.screen.SetContent(x, j, '│', nil, borderStyle)
		for i := x + 1; i < x+width-1; i++ {
			r.screen.SetContent(i, j, ' ', nil, bgStyle)
		}
		r.screen.SetContent(x+width-1, j, '│', nil, borderStyle)
	}

	// Bottom border
	r.screen.SetContent(x, y+height-1, '└', nil, borderStyle)
	for i := x + 1; i < x+width-1; i++ {
		r.screen.SetContent(i, y+height-1, '─', nil, borderStyle)
	}
	r.screen.SetContent(x+width-1, y+height-1, '┘', nil, borderStyle)
}

// drawBoxTitle draws a centered title on the top border of a box.
func (r *Renderer) drawBoxTitle(boxX, boxY, boxWidth int, title string, theme Theme) {
	titleX := boxX + (boxWidth-len(title))/2
	r.DrawText(titleX, boxY, title, theme.Title, theme.Background)
}

// drawFilterInput draws the filter input line in the command menu.
func (r *Renderer) drawFilterInput(menuX, menuY int, filter string, theme Theme) {
	filterPrompt := "> " + filter
	r.DrawText(menuX+2, menuY+2, filterPrompt, theme.Foreground, theme.Background)

	// Draw cursor
	cursorX := menuX + 2 + len(filterPrompt)
	cursorStyle := tcell.StyleDefault.Foreground(theme.TextCursor).Background(theme.Background)
	r.screen.SetContent(cursorX, menuY+2, '▏', nil, cursorStyle)
}

// drawCommandList draws the list of filtered commands in the command menu.
func (r *Renderer) drawCommandList(menuX, menuY, menuWidth, menuHeight int, data CommandMenuData) {
	// Draw separator
	borderStyle := tcell.StyleDefault.Foreground(data.Theme.Border).Background(data.Theme.Background)
	for x := menuX + 1; x < menuX+menuWidth-1; x++ {
		r.screen.SetContent(x, menuY+3, '─', nil, borderStyle)
	}

	maxCommands := menuHeight - 5
	startY := menuY + 4

	if len(data.FilteredCommands) == 0 {
		r.drawNoResults(menuX, menuWidth, startY, data.Theme)
		return
	}

	for i := 0; i < maxCommands && i < len(data.FilteredCommands); i++ {
		cmd := data.FilteredCommands[i]
		y := startY + i

		var style tcell.Style
		if i == data.Selected {
			style = tcell.StyleDefault.Foreground(data.Theme.MenuSelectedFg).Background(data.Theme.MenuSelectedBg).Bold(true)
		} else {
			style = tcell.StyleDefault.Foreground(data.Theme.Foreground).Background(data.Theme.Background)
		}

		// Clear line with style
		for x := menuX + 2; x < menuX+menuWidth-2; x++ {
			r.screen.SetContent(x, y, ' ', nil, style)
		}

		// Draw command name (truncated if needed)
		maxLen := menuWidth - 4
		displayText := cmd.Name
		if len(displayText) > maxLen {
			displayText = displayText[:maxLen-3] + "..."
		}

		for j, ch := range displayText {
			r.screen.SetContent(menuX+2+j, y, ch, nil, style)
		}
	}
}

// drawNoResults draws the "no matching commands" message.
func (r *Renderer) drawNoResults(menuX, menuWidth, startY int, theme Theme) {
	noResults := "no matching commands"
	noResultsX := menuX + (menuWidth-len(noResults))/2
	r.DrawText(noResultsX, startY+2, noResults, theme.MenuDimText, theme.Background)
}

// drawResultsContent draws the statistics and misspelled words in the results screen.
func (r *Renderer) drawResultsContent(boxX, boxY, boxWidth, boxHeight int, data ResultsData) {
	currentY := boxY + 2

	// Draw stats
	wpmText := fmt.Sprintf("WPM: %.1f", data.WPM)
	r.DrawText(boxX+4, currentY, wpmText, data.Theme.Foreground, data.Theme.Background)
	currentY++

	accuracyText := fmt.Sprintf("Accuracy: %.1f%%", data.Accuracy)
	r.DrawText(boxX+4, currentY, accuracyText, data.Theme.Foreground, data.Theme.Background)
	currentY += 2

	// Draw separator
	borderStyle := tcell.StyleDefault.Foreground(data.Theme.Border).Background(data.Theme.Background)
	for x := boxX + 2; x < boxX+boxWidth-2; x++ {
		r.screen.SetContent(x, currentY, '─', nil, borderStyle)
	}
	currentY += 2

	// Draw misspelled words
	if len(data.MisspelledWords) == 0 {
		perfectText := "Perfect! No mistakes!"
		perfectX := boxX + (boxWidth-len(perfectText))/2
		r.DrawText(perfectX, currentY, perfectText, data.Theme.TextCorrect, data.Theme.Background)
	} else {
		header := "Misspelled Words:"
		r.DrawText(boxX+4, currentY, header, data.Theme.Title, data.Theme.Background)
		currentY += 2

		maxWords := boxHeight - 12
		for i, word := range data.MisspelledWords {
			if i >= maxWords {
				moreText := fmt.Sprintf("... and %d more", len(data.MisspelledWords)-i)
				r.DrawText(boxX+6, currentY, moreText, data.Theme.MenuDimText, data.Theme.Background)
				break
			}

			count := data.WordCounts[word]
			wordText := fmt.Sprintf("• %s", word)
			if count > 1 {
				wordText = fmt.Sprintf("• %s (x%d)", word, count)
			}
			r.DrawText(boxX+6, currentY, wordText, data.Theme.TextIncorrect, data.Theme.Background)
			currentY++
		}
	}

	// Draw help text
	helpText := "Press Enter or 'r' to restart  |  Esc to quit"
	helpX := boxX + (boxWidth-len(helpText))/2
	r.DrawText(helpX, boxY+boxHeight-2, helpText, data.Theme.Help, data.Theme.Background)
}

// wrapText breaks text into lines that fit within maxWidth characters.
// Respects explicit newlines and attempts to break at word boundaries.
func wrapText(text string, maxWidth int) []string {
	var lines []string
	var currentLine []rune

	for _, ch := range text {
		if ch == '\n' {
			lines = append(lines, string(currentLine)+string(ch))
			currentLine = []rune{}
		} else if len(currentLine) >= maxWidth {
			// Auto-wrap at maxWidth - try to break at last space
			breakPoint := len(currentLine)
			for i := len(currentLine) - 1; i >= 0; i-- {
				if currentLine[i] == ' ' {
					breakPoint = i + 1
					break
				}
			}

			lines = append(lines, string(currentLine[:breakPoint]))
			currentLine = currentLine[breakPoint:]
			currentLine = append(currentLine, ch)
		} else {
			currentLine = append(currentLine, ch)
		}
	}

	if len(currentLine) > 0 {
		lines = append(lines, string(currentLine))
	}

	return lines
}

// CalculateCursorLine determines which wrapped line the cursor is on.
// Returns the line index (0-based) within the wrapped lines.
func CalculateCursorLine(text string, cursorPos int, maxWidth int) int {
	if cursorPos < 0 {
		return 0
	}

	lines := wrapText(text, maxWidth)
	charCount := 0

	for lineIdx, line := range lines {
		lineLen := len([]rune(line))
		if charCount+lineLen > cursorPos {
			return lineIdx
		}
		charCount += lineLen
	}

	// Cursor is at or beyond the end
	if len(lines) > 0 {
		return len(lines) - 1
	}
	return 0
}

// CalculateScrollLine calculates the optimal scroll line to keep the cursor visible.
// It tries to keep the cursor in the middle of the viewport when possible.
func CalculateScrollLine(cursorLine, maxVisibleLines, totalLines int) int {
	// If all text fits on screen, don't scroll
	if totalLines <= maxVisibleLines {
		return 0
	}

	// Try to keep cursor in the middle third of the viewport
	desiredCursorPosition := maxVisibleLines / 3

	scrollLine := cursorLine - desiredCursorPosition
	if scrollLine < 0 {
		scrollLine = 0
	}

	// Don't scroll past the end
	maxScroll := totalLines - maxVisibleLines
	if scrollLine > maxScroll {
		scrollLine = maxScroll
	}

	return scrollLine
}
