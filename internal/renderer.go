package internal

import (
	"fmt"
	"time"

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
func (r *Renderer) DrawTitle(themeName, textName, modeInfo string, theme Theme) {
	width, _ := r.screen.Size()
	var title string
	if modeInfo != "" {
		title = fmt.Sprintf("rocketype [%s] - %s (%s)", themeName, textName, modeInfo)
	} else {
		title = fmt.Sprintf("rocketype [%s] - %s", themeName, textName)
	}
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

// DrawProgress renders progress information (timer or word count) above stats.
func (r *Renderer) DrawProgress(progressText string, theme Theme) {
	width, height := r.screen.Size()
	x := width/2 - len(progressText)/2
	r.DrawText(x, height-4, progressText, theme.Help, theme.Background)
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
	WordMode    bool // True if in word mode (shows only 2 lines below cursor)
}

// DrawTypingView renders the main typing test interface with wrapped text and visual feedback.
func (r *Renderer) DrawTypingView(data TypingViewData) {
	width, height := r.screen.Size()

	// Calculate available space for text
	// IMPORTANT: This maxWidth calculation must match app.go's calculations
	// for cursor line and text wrapping to work correctly!
	maxWidth := width - 20
	if maxWidth < 20 {
		maxWidth = width
	}

	// Wrap text to fit screen width
	lines := wrapText(data.SampleText, maxWidth)

	// Calculate available height for text lines
	availableHeight := height - 8
	maxVisibleLines := availableHeight / 2 // 2 screen rows per text line

	// In word mode, only show 3 lines (cursor line + 2 below)
	// This constant should match wordModeVisibleLines in app.go
	const wordModeVisibleLines = 3
	if data.WordMode {
		maxVisibleLines = wordModeVisibleLines
	}

	// Adjust scroll position if needed
	scrollLine := data.ScrollLine
	if scrollLine < 0 {
		scrollLine = 0
	}
	if scrollLine > len(lines)-1 {
		scrollLine = len(lines) - 1
	}

	// Calculate how many lines will actually be rendered
	endLine := scrollLine + maxVisibleLines
	if endLine > len(lines) {
		endLine = len(lines)
	}
	visibleLineCount := endLine - scrollLine

	// Calculate vertical centering
	// Each line takes 2 rows (text + space), calculate total height needed
	contentHeight := visibleLineCount * 2
	// Center the content vertically in available space
	startY := (height - contentHeight) / 2
	if startY < 4 {
		startY = 4 // Keep minimum spacing from top
	}

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
	ScrollOffset     int
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
	WPMHistory      []WPMSnapshot // Timeline of WPM measurements
	ErrorTimestamps []time.Time   // Timestamps when errors occurred
	Theme           Theme
}

// DrawResults renders the results screen overlay.
func (r *Renderer) DrawResults(data ResultsData) {
	width, height := r.screen.Size()

	// Make box larger to accommodate taller graph
	boxWidth := min(width*4/5, 80)
	boxHeight := min(height*4/5, 45)
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

	// Calculate the visible window based on scroll offset
	startIdx := data.ScrollOffset
	endIdx := min(startIdx+maxCommands, len(data.FilteredCommands))

	// Draw scroll indicators if needed
	if startIdx > 0 {
		// Show "more above" indicator
		r.DrawText(menuX+menuWidth-3, menuY+3, "▲", data.Theme.Border, data.Theme.Background)
	}
	if endIdx < len(data.FilteredCommands) {
		// Show "more below" indicator
		r.DrawText(menuX+menuWidth-3, menuY+menuHeight-2, "▼", data.Theme.Border, data.Theme.Background)
	}

	// Draw visible commands
	for i := startIdx; i < endIdx; i++ {
		cmd := data.FilteredCommands[i]
		displayIdx := i - startIdx
		y := startY + displayIdx

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

// drawResultsContent draws the statistics, WPM timeline graph, and misspelled words in the results screen.
func (r *Renderer) drawResultsContent(boxX, boxY, boxWidth, boxHeight int, data ResultsData) {
	currentY := boxY + 2

	// Draw stats
	wpmText := fmt.Sprintf("WPM: %.1f", data.WPM)
	r.DrawText(boxX+4, currentY, wpmText, data.Theme.Foreground, data.Theme.Background)
	currentY++

	accuracyText := fmt.Sprintf("Accuracy: %.1f%%", data.Accuracy)
	r.DrawText(boxX+4, currentY, accuracyText, data.Theme.Foreground, data.Theme.Background)
	currentY += 2

	// Draw WPM timeline graph if we have history
	if len(data.WPMHistory) > 1 {
		graphHeight := 17 // Increased to accommodate X-axis labels
		graphWidth := boxWidth - 8
		r.drawWPMGraph(boxX+4, currentY, graphWidth, graphHeight, data.WPMHistory, data.ErrorTimestamps, data.Theme)
		currentY += graphHeight + 2
	}

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

		// Build comma-separated list with counts
		var wordList []string
		for _, word := range data.MisspelledWords {
			count := data.WordCounts[word]
			if count > 1 {
				wordList = append(wordList, fmt.Sprintf("%s (x%d)", word, count))
			} else {
				wordList = append(wordList, word)
			}
		}

		// Calculate available width and height for wrapping
		contentWidth := boxWidth - 12                        // Leave margin on both sides
		availableHeight := boxHeight - (currentY - boxY) - 3 // Space until help text

		// Wrap the text to fit width
		currentLine := ""
		linesDrawn := 0

		for i, word := range wordList {
			var testLine string
			if currentLine == "" {
				testLine = word
			} else {
				testLine = currentLine + ", " + word
			}

			// Check if adding this word exceeds width
			if len(testLine) > contentWidth {
				// Draw current line and start new one
				if linesDrawn < availableHeight {
					r.DrawText(boxX+6, currentY, currentLine, data.Theme.TextIncorrect, data.Theme.Background)
					currentY++
					linesDrawn++
				}
				currentLine = word
			} else {
				currentLine = testLine
			}

			// If this is the last word, draw the remaining line
			if i == len(wordList)-1 {
				if linesDrawn < availableHeight {
					r.DrawText(boxX+6, currentY, currentLine, data.Theme.TextIncorrect, data.Theme.Background)
				} else {
					// Too many words, show truncation message
					moreText := fmt.Sprintf("... and more")
					r.DrawText(boxX+6, currentY, moreText, data.Theme.MenuDimText, data.Theme.Background)
				}
			}
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
// It keeps the cursor positioned with space above and below for context.
func CalculateScrollLine(cursorLine, maxVisibleLines, totalLines int) int {
	// If all text fits on screen, don't scroll
	if totalLines <= maxVisibleLines {
		return 0
	}

	// Desired buffer: keep at least 1 line visible below cursor
	const minBufferBelow = 1

	// Desired position: top third of viewport (gives more context below)
	desiredCursorPosition := maxVisibleLines / 3

	// But ensure we leave room for the buffer below
	maxCursorPosition := maxVisibleLines - minBufferBelow - 1
	if desiredCursorPosition > maxCursorPosition {
		desiredCursorPosition = maxCursorPosition
	}

	// Calculate scroll line based on desired cursor position
	scrollLine := cursorLine - desiredCursorPosition

	// Clamp to valid range [0, maxScroll]
	if scrollLine < 0 {
		scrollLine = 0
	}

	maxScroll := totalLines - maxVisibleLines
	if scrollLine > maxScroll {
		scrollLine = maxScroll
	}

	return scrollLine
}

// drawWPMGraph renders a timeline graph of WPM changes over time.
// The graph uses ASCII characters to draw a line chart showing typing speed progression.
//
// Parameters:
//   - x, y: top-left position of the graph
//   - width, height: dimensions of the graph area
//   - history: slice of WPM snapshots to plot
//   - errorTimestamps: timestamps when typing errors occurred
//   - theme: color theme for rendering
func (r *Renderer) drawWPMGraph(x, y, width, height int, history []WPMSnapshot, errorTimestamps []time.Time, theme Theme) {
	if len(history) < 2 || width < 10 || height < 3 {
		return
	}

	// Find max WPM for scaling
	maxWPM := history[0].WPM
	for _, snapshot := range history {
		if snapshot.WPM > maxWPM {
			maxWPM = snapshot.WPM
		}
	}

	// Always start at 0 and round up to nearest 25 WPM increment
	minWPM := 0.0
	maxWPM = float64(int(maxWPM/25)+1) * 25
	if maxWPM < 25 {
		maxWPM = 25
	}

	// Draw title
	title := "WPM Timeline"
	titleX := x + (width-len(title))/2
	r.DrawText(titleX, y, title, theme.Title, theme.Background)
	y++ // Move down after title

	graphHeight := height - 3 // Reserve space for title, Y-axis labels, and X-axis
	graphWidth := width - 7   // Reserve space for Y-axis labels (4 chars + 3 space)

	// Draw Y-axis labels at every 25 WPM increment
	numLabels := int(maxWPM/25) + 1 // Number of labels from 0 to maxWPM
	for i := 0; i < numLabels; i++ {
		wpmValue := float64(i) * 25
		label := fmt.Sprintf("%4.0f", wpmValue) // Right-align with width of 4

		// Calculate Y position for this label (inverted)
		normalized := (wpmValue - minWPM) / (maxWPM - minWPM)
		labelY := y + graphHeight - 1 - int(normalized*float64(graphHeight-1))

		r.DrawText(x, labelY, label, theme.Help, theme.Background)
	}

	// Starting position for graph content
	graphX := x + 6 // Y-axis labels are 4 chars + 2 space padding
	graphY := y

	// Initialize graph area with spaces
	graphStyle := tcell.StyleDefault.Foreground(theme.TextDefault).Background(theme.Background)
	for gy := 0; gy < graphHeight; gy++ {
		for gx := 0; gx < graphWidth; gx++ {
			r.screen.SetContent(graphX+gx, graphY+gy, ' ', nil, graphStyle)
		}
	}

	// Calculate points for the graph
	points := make([]int, graphWidth)
	for i := range points {
		// Map column to history index
		historyIdx := int(float64(i) / float64(graphWidth-1) * float64(len(history)-1))
		if historyIdx >= len(history) {
			historyIdx = len(history) - 1
		}

		wpm := history[historyIdx].WPM

		// Scale WPM to graph height (inverted because Y increases downward)
		normalized := (wpm - minWPM) / (maxWPM - minWPM)
		if normalized < 0 {
			normalized = 0
		}
		if normalized > 1 {
			normalized = 1
		}

		// Convert to screen coordinates (invert Y)
		points[i] = graphHeight - 1 - int(normalized*float64(graphHeight-1))
	}

	// Draw graph using braille characters for smooth lines
	lineStyle := tcell.StyleDefault.Foreground(theme.TextCorrect).Background(theme.Background).Bold(true)

	// Helper function for absolute value
	abs := func(n int) int {
		if n < 0 {
			return -n
		}
		return n
	}

	// Braille characters are 2 dots wide by 4 dots tall
	// We need braille grid that matches our graphWidth and graphHeight
	brailleWidth := graphWidth   // Number of braille characters horizontally
	brailleHeight := graphHeight // Number of braille characters vertically
	brailleGrid := make([][]uint8, brailleHeight)
	for i := range brailleGrid {
		brailleGrid[i] = make([]uint8, brailleWidth)
	}

	// Map each point to braille grid with 2x4 sub-pixel precision
	for i := 0; i < len(points)-1; i++ {
		// Points array has one entry per character column
		// Scale to 2x width (braille horizontal resolution) and 4x height (braille vertical resolution)
		x1, y1 := i*2, points[i]*4
		x2, y2 := (i+1)*2, points[i+1]*4

		// Draw line segment between consecutive points using Bresenham's algorithm
		dx := abs(x2 - x1)
		dy := abs(y2 - y1)
		sx := 1
		if x1 > x2 {
			sx = -1
		}
		sy := 1
		if y1 > y2 {
			sy = -1
		}
		err := dx - dy

		x, y := x1, y1
		for {
			// Convert sub-pixel coordinates to braille cell and dot position
			cellX := x / 2
			cellY := y / 4
			dotX := x % 2 // 0 or 1 (left or right dot in braille cell)
			dotY := y % 4 // 0-3 (which row of dots in braille cell)

			// Set braille dot if within grid bounds
			if cellX >= 0 && cellX < brailleWidth && cellY >= 0 && cellY < brailleHeight {
				// Braille dot pattern: dots are numbered 0-7
				// Left column: 0,1,2,3 (top to bottom), Right column: 4,5,6,7 (top to bottom)
				dotIndex := dotX*4 + dotY
				brailleGrid[cellY][cellX] |= (1 << dotIndex)
			}

			if x == x2 && y == y2 {
				break
			}

			e2 := 2 * err
			if e2 > -dy {
				err -= dy
				x += sx
			}
			if e2 < dx {
				err += dx
				y += sy
			}
		}
	}

	// Convert braille grid to unicode braille characters and render
	brailleBase := rune(0x2800) // Unicode braille pattern base
	for cellY := 0; cellY < brailleHeight; cellY++ {
		for cellX := 0; cellX < brailleWidth; cellX++ {
			if brailleGrid[cellY][cellX] != 0 {
				brailleChar := brailleBase + rune(brailleGrid[cellY][cellX])
				r.screen.SetContent(graphX+cellX, graphY+cellY, brailleChar, nil, lineStyle)
			}
		}
	}

	// Draw error markers
	if len(history) > 0 && len(errorTimestamps) > 0 {
		startTime := history[0].Timestamp
		endTime := history[len(history)-1].Timestamp
		totalDuration := endTime.Sub(startTime).Seconds()

		if totalDuration > 0 {
			errorStyle := tcell.StyleDefault.Foreground(theme.TextIncorrect).Background(theme.Background)

			for _, errorTime := range errorTimestamps {
				// Calculate time offset from start
				errorOffset := errorTime.Sub(startTime).Seconds()

				// Skip errors outside the graph range
				if errorOffset < 0 || errorOffset > totalDuration {
					continue
				}

				// Calculate X position for this error
				normalized := errorOffset / totalDuration
				errorX := graphX + int(normalized*float64(graphWidth-1))

				// Draw error marker at the bottom of the graph
				r.screen.SetContent(errorX, graphY+graphHeight-1, '×', nil, errorStyle)
			}
		}
	}

	// Draw X-axis with time labels
	if len(history) > 0 {
		// Calculate total duration in seconds
		totalDuration := history[len(history)-1].Timestamp.Sub(history[0].Timestamp).Seconds()

		// Draw X-axis time labels
		xAxisY := graphY + graphHeight + 1

		// Determine time interval for labels based on duration
		var interval float64
		if totalDuration <= 30 {
			interval = 5 // Every 5 seconds for short tests
		} else if totalDuration <= 60 {
			interval = 10 // Every 10 seconds for medium tests
		} else if totalDuration <= 120 {
			interval = 15 // Every 15 seconds for longer tests
		} else {
			interval = 30 // Every 30 seconds for very long tests
		}

		// Draw time labels
		for t := 0.0; t <= totalDuration; t += interval {
			// Calculate X position for this time
			normalized := t / totalDuration
			labelX := graphX + int(normalized*float64(graphWidth-1))

			// Format time label
			var timeLabel string
			if t >= 60 {
				minutes := int(t / 60)
				seconds := int(t) % 60
				if seconds == 0 {
					timeLabel = fmt.Sprintf("%dm", minutes)
				} else {
					timeLabel = fmt.Sprintf("%dm%ds", minutes, seconds)
				}
			} else {
				timeLabel = fmt.Sprintf("%.0fs", t)
			}

			// Center the label on the tick mark
			labelX -= len(timeLabel) / 2
			if labelX < graphX {
				labelX = graphX
			}
			if labelX+len(timeLabel) > graphX+graphWidth {
				labelX = graphX + graphWidth - len(timeLabel)
			}

			r.DrawText(labelX, xAxisY, timeLabel, theme.Help, theme.Background)
		}
	}
}
