// Package main is the entry point for the rocketype terminal typing test application.
//
// Rocketype is a terminal-based typing test inspired by monkeytype, built using
// the tcell library for direct terminal manipulation. It provides real-time typing
// feedback, accuracy tracking, and customizable themes.
//
// Usage:
//
//	rocketype                                    # Start with random text from default location
//	rocketype --texts-dir ~/my-texts             # Use custom texts directory
//	cat myfile.txt | rocketype                   # Practice with custom text via stdin
//	echo "custom text" | rocketype               # Practice with inline text
//
// Default text locations:
//   - Linux: ~/.config/rocketype/texts
//   - macOS: ~/Library/Application Support/rocketype/texts
//   - Windows: %APPDATA%\rocketype\texts
//
// The application starts immediately in typing mode. Use Ctrl+P to access the
// command palette, Ctrl+T to cycle themes, and Ctrl+C or Esc to quit.
//
// When text is piped via stdin, it becomes available as the "stdin" text source
// and is automatically selected as the active practice text.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"baumeister.de/rocketype/internal"
)

func main() {
	// Define command-line flags
	textsDir := flag.String("texts-dir", "", "Path to texts directory (overrides platform default)")
	printPaths := flag.Bool("print-paths", false, "Print default paths and exit")
	restoreSession := flag.Bool("restore-session", true, "Restore previous session on startup (default: true)")

	// Custom usage message
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Rocketype - A terminal-based typing test application\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nDefault text locations:\n")
		defaultDir, _ := internal.GetDefaultTextsDir()
		fmt.Fprintf(os.Stderr, "  Platform default: %s\n", defaultDir)
		fmt.Fprintf(os.Stderr, "  Local fallback:   %s\n", internal.GetFallbackTextsDir())
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s                           # Use default texts location\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s --texts-dir ~/my-texts   # Use custom directory\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  cat file.txt | %s           # Practice with piped text\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s --restore-session=false  # Start fresh, ignore saved session\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nKeyboard shortcuts:\n")
		fmt.Fprintf(os.Stderr, "  Ctrl+P     - Open command menu\n")
		fmt.Fprintf(os.Stderr, "  Ctrl+T     - Cycle themes\n")
		fmt.Fprintf(os.Stderr, "  Ctrl+C/Esc - Quit\n")
	}

	flag.Parse()

	// If user wants to see paths, print and exit
	if *printPaths {
		defaultDir, err := internal.GetDefaultTextsDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting default path: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Default texts directory: %s\n", defaultDir)
		fmt.Printf("Fallback directory: %s\n", internal.GetFallbackTextsDir())
		os.Exit(0)
	}

	// Determine which texts directory to use
	var finalTextsDir string
	if *textsDir != "" {
		// User specified a custom directory
		finalTextsDir = *textsDir
	} else {
		// Try platform default first
		defaultDir, err := internal.GetDefaultTextsDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Could not determine default texts directory: %v\n", err)
			fmt.Fprintf(os.Stderr, "Falling back to: %s\n", internal.GetFallbackTextsDir())
			finalTextsDir = internal.GetFallbackTextsDir()
		} else {
			finalTextsDir = defaultDir

			// Check if platform default exists and has files, otherwise try fallback
			if _, err := os.Stat(finalTextsDir); os.IsNotExist(err) {
				// Check if fallback exists
				fallback := internal.GetFallbackTextsDir()
				if _, err := os.Stat(fallback); err == nil {
					finalTextsDir = fallback
				}
			}
		}
	}

	// Check if input is being piped via stdin
	var stdinText string
	stat, err := os.Stdin.Stat()
	if err == nil && (stat.Mode()&os.ModeCharDevice) == 0 {
		// stdin is being piped, read the content
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
			os.Exit(1)
		}
		stdinText = strings.TrimSpace(string(data))
		if stdinText == "" {
			fmt.Fprintf(os.Stderr, "Error: stdin is empty\n")
			os.Exit(1)
		}
	}

	// Create and initialize the application
	app, err := internal.NewApp(stdinText, finalTextsDir, *restoreSession)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating app: %v\n", err)
		os.Exit(1)
	}

	// Run the main event loop (blocks until quit)
	if err := app.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running app: %v\n", err)
		os.Exit(1)
	}
}
