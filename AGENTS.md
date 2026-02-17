# AGENTS.md - Developer Guide for rocketype

This guide is for AI coding agents and developers working on the rocketype codebase.

## Project Overview

Rocketype is a terminal-based typing test application written in Go, inspired by monkeytype. It uses the tcell library for terminal manipulation and follows clean architecture principles with clear separation of concerns.

**Key Technologies:**
- Go 1.25.7
- github.com/gdamore/tcell/v2 for terminal UI
- Standard library for file I/O and JSON handling

## Build & Run Commands

### Building
```bash
# Build the binary
go build -o rocketype ./cmd/rocketype

# Build and run
go build -o rocketype ./cmd/rocketype && ./rocketype

# Build for specific platform
GOOS=linux GOARCH=amd64 go build -o rocketype ./cmd/rocketype
```

### Running
```bash
# Run directly without building
go run ./cmd/rocketype

# Run with custom texts directory
go run ./cmd/rocketype --texts-dir ~/my-texts

# Run with piped input
echo "custom text" | go run ./cmd/rocketype

# Show platform paths
go run ./cmd/rocketype --print-paths
```

### Testing
```bash
# Run all tests
go test ./...

# Run tests for a specific package
go test ./internal

# Run tests with verbose output
go test -v ./...

# Run a specific test function
go test -v ./internal -run TestFunctionName

# Run tests with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Linting & Formatting
```bash
# Format all code (must be run before committing)
gofmt -w .

# Check formatting without modifying
gofmt -d .

# Vet code for common issues
go vet ./...

# Install and run staticcheck (if available)
go install honnef.co/go/tools/cmd/staticcheck@latest
staticcheck ./...
```

## Architecture & Code Structure

### Directory Layout
```
rocketype/
├── cmd/rocketype/          # Entry point (main.go)
│   └── main.go            # CLI flags, stdin handling, app initialization
├── internal/              # Core application code (not importable by other projects)
│   ├── doc.go            # Package documentation and architecture overview
│   ├── app.go            # Main application controller & coordinator
│   ├── renderer.go       # UI rendering (presentation layer)
│   ├── typingtest.go     # Business logic for typing tests
│   ├── stats.go          # Statistics tracking (WPM, accuracy, errors)
│   ├── input_handler.go  # Input routing based on application mode
│   ├── theme.go          # Color scheme definitions
│   ├── command.go        # Command palette with fuzzy search
│   ├── textlib.go        # Text file loading and management
│   ├── wordlib.go        # Word list management for word mode
│   ├── session.go        # Session persistence (state saving/loading)
│   ├── settings.go       # User preferences (theme, mode settings)
│   └── paths.go          # Platform-specific path resolution
├── go.mod                 # Module definition
└── README.md             # User-facing documentation
```

### Key Design Principles

1. **Clean Architecture**: Business logic (TypingTest) has no dependencies on UI (Renderer)
2. **Separation of Concerns**: UI rendering, business logic, and input handling are separate
3. **Data-Oriented Rendering**: Renderer receives data structs, not business objects
4. **Theme-Driven UI**: All colors come from theme definitions (no hardcoded colors)
5. **Error Persistence**: Misspelled words tracked even if corrected via backspace

### Component Responsibilities

- **App** (app.go): Coordinates between components, manages state, event loop
- **Renderer** (renderer.go): Pure presentation - draws UI from data structures
- **TypingTest** (typingtest.go): Core typing logic - no UI knowledge
- **Stats** (stats.go): Accuracy, WPM calculation, error tracking
- **InputHandler** (input_handler.go): Routes keyboard events based on mode
- **CommandMenu** (command.go): Command palette with fuzzy filtering
- **TextLibrary** (textlib.go): Loads .txt files from platform directories
- **SessionManager** (session.go): Persists and restores application state
- **SettingsManager** (settings.go): Manages user preferences

## Code Style Guidelines

### Imports
- Standard library first
- Third-party packages second
- Local packages last
- Groups separated by blank lines
- Example:
```go
import (
    "fmt"
    "strings"
    "time"

    "github.com/gdamore/tcell/v2"
)
```

### Formatting
- **Required**: Run `gofmt -w .` before committing
- Use tabs for indentation (Go standard)
- Line length: No strict limit, but keep it reasonable (~120 chars preferred)
- Prefer early returns over nested conditionals

### Naming Conventions
- **Packages**: lowercase, single word when possible (`internal`, not `internal_lib`)
- **Types**: PascalCase (`TypingTest`, `TextLibrary`)
- **Functions/Methods**: PascalCase for exported, camelCase for unexported
- **Constants**: PascalCase or camelCase, use `const` blocks with iota when appropriate
- **Variables**: camelCase, descriptive names preferred over abbreviations
- **Private fields**: Start with lowercase (`cursorPos`, `userInput`)
- **Acronyms**: Keep case consistent (e.g., `WPM`, `GetDefaultTextsDir`)

### Types & Structs
- Document all exported types with godoc comments
- Group related fields together in structs
- Use struct tags for JSON serialization: `` `json:"field_name"` ``
- Initialize complex structs with constructor functions (e.g., `NewApp`, `NewStats`)
- Example:
```go
// Stats tracks typing test statistics including timing, accuracy, and error tracking.
// It maintains detailed information about keystrokes, misspelled words, and test progress.
type Stats struct {
    // Timing information
    startTime time.Time
    endTime   time.Time

    // Keystroke tracking
    totalKeystrokes   int
    correctKeystrokes int
}
```

### Comments & Documentation
- All exported functions, types, and methods must have godoc comments
- Comment starts with the name of the thing being documented
- Use `//` for single-line comments, not `/* */`
- Document parameters with `Parameters:` section when complex
- Document return values when non-obvious
- Example:
```go
// NewTextLibrary creates a new TextLibrary instance.
// It loads all .txt files from the specified directory, or uses the default
// embedded text if the directory doesn't exist or contains no files.
//
// Parameters:
//   - textsDir: directory path to search for .txt files
//
// Returns a TextLibrary with at least one text (the default if no files found).
func NewTextLibrary(textsDir string) *TextLibrary {
```

### Error Handling
- Always check errors, don't ignore them
- Wrap errors with context using `fmt.Errorf` with `%w` verb
- Return errors up the stack rather than logging and continuing
- Use early returns to reduce nesting
- Example:
```go
data, err := os.ReadFile(path)
if err != nil {
    return fmt.Errorf("failed to read settings file: %w", err)
}
```

### Constants
- Define constants at package level or within functions if scoped
- Use `const` blocks for related constants
- Document constant blocks, especially when using iota
- Example:
```go
const (
    // CharsPerWord represents the standard conversion factor for WPM calculation.
    // The industry standard is 5 characters = 1 word.
    CharsPerWord = 5.0
)
```

### Functions
- Keep functions focused and single-purpose
- Prefer small functions over large ones
- Use named return values only when it improves clarity
- Constructor functions should be named `New<TypeName>`
- Getters don't use "Get" prefix unless necessary for clarity

## Common Patterns

### Adding a New Theme
1. Define a new `Theme` variable in `theme.go`
2. Add it to the `AvailableThemes()` function
3. Add a command in `app.initCommands()` to select it

### Adding a New Command
1. Add a `Command` struct to the slice in `app.initCommands()`
2. Define the `Action` function inline or as a method
3. Provide clear `Name` and `Description` for fuzzy search

### Platform-Specific Paths
- Use `runtime.GOOS` to detect platform
- Follow conventions: XDG on Linux, ~/Library on macOS, %APPDATA% on Windows
- Always create directories with `os.MkdirAll` and check errors
- Use `filepath.Join` for cross-platform path construction

### Working with Unicode
- Use `[]rune` for text that needs character-level manipulation
- Cache rune slices when accessed frequently
- Be mindful of byte vs. rune indexing

## Testing Guidelines

- Test files should be named `*_test.go` (though none exist yet)
- Place tests in the same package as the code being tested
- Use table-driven tests when testing multiple cases
- Test exported functions; unexported functions indirectly via exports

## Git Workflow

- Commit messages should be concise and descriptive
- Focus on "why" rather than "what"
- Reference internal/doc.go for architecture decisions

## File References

When discussing code locations, use the format: `path/to/file.go:line`
Example: "Stats tracking happens in internal/stats.go:42"
