// Package internal contains the core implementation of the rocketype typing test application.
//
// Architecture Overview:
//
// The package follows clean architecture principles with clear separation of concerns:
//
// Core Components:
//
//   - App (app.go): Application controller that coordinates between components.
//     Manages state, delegates to specialized components, and handles the event loop.
//
//   - Renderer (renderer.go): UI rendering layer responsible for all drawing logic.
//     Accepts data structures and renders them to the screen using tcell.
//     No business logic - pure presentation.
//
//   - TypingTest (typingtest.go): Business logic for the typing test.
//     Handles character input, cursor movement, word tracking, and completion detection.
//     No knowledge of UI or rendering.
//
//   - InputHandler (input_handler.go): Input routing based on application mode.
//     Separates keyboard event handling from business logic and UI concerns.
//
//   - Stats (stats.go): Statistics tracking including WPM, accuracy, and misspelled words.
//     Tracks errors even after correction for authentic accuracy metrics.
//
//   - Theme (theme.go): Color scheme definitions and theme management.
//     All colors are defined in theme structs (no hardcoded colors in rendering code).
//
//   - CommandMenu (command.go): Command palette with fuzzy filtering for accessing features.
//
//   - TextLibrary (textlib.go): Manages loading and selection of practice texts from files.
//
// Design Principles:
//
// 1. Separation of Concerns: UI rendering, business logic, and input handling are separate.
//
// 2. Clean Architecture: Business logic (TypingTest) has no dependencies on UI (Renderer).
//
//  3. Data-Oriented Rendering: Renderer receives data structs (TypingViewData, ResultsData)
//     and has no knowledge of business logic.
//
//  4. Theme-Driven Rendering: All colors come from theme definitions, making themes
//     easy to add without modifying rendering code.
//
// 5. Error Persistence: Misspelled words are tracked even if corrected via backspace.
//
// 6. Testability: Components are decoupled and can be tested independently.
//
// Component Interactions:
//
//	User Input → InputHandler → TypingTest (updates state)
//	                          ↓
//	App (coordinates) → Renderer (draws UI)
//	                          ↑
//	TypingTest (provides data) → Stats, Theme
//
// Extension Points:
//
// To add a new theme:
//  1. Define a new Theme variable in theme.go
//  2. Add it to the AvailableThemes() function
//  3. Add a command in app.initCommands() to select it
//
// To add a new command:
//  1. Add a Command struct to the slice in app.initCommands()
//  2. Define the Action function inline or as a method
//
// To add custom practice texts:
//  1. Create .txt files in the texts/ directory
//  2. The TextLibrary automatically loads them on startup
//
// To extend rendering:
//  1. Add methods to Renderer with appropriate data structures
//  2. Call from App.draw() or related methods
package internal
