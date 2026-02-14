// Package main is the entry point for the rocketype terminal typing test application.
//
// Rocketype is a terminal-based typing test inspired by monkeytype, built using
// the tcell library for direct terminal manipulation. It provides real-time typing
// feedback, accuracy tracking, and customizable themes.
//
// Usage:
//
//	rocketype                                    # Start with random text from texts/ directory
//	cat myfile.txt | rocketype                  # Practice with custom text via stdin
//	echo "custom text" | rocketype              # Practice with inline text
//
// The application starts immediately in typing mode. Use Ctrl+P to access the
// command palette, Ctrl+T to cycle themes, and Ctrl+C or Esc to quit.
//
// When text is piped via stdin, it becomes available as the "stdin" text source
// and is automatically selected as the active practice text.
package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"baumeister.de/rocketype/internal"
)

func main() {
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
	app, err := internal.NewApp(stdinText)
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
