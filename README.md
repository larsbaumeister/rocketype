# rocketype

A fast, minimalist typing test application for the terminal, inspired by [monkeytype](https://monkeytype.com).

![rocketype demo](https://via.placeholder.com/800x400?text=rocketype+demo)

## Features

- 🎨 **Multiple Themes** - Choose from dark and light themes (Gruvbox, Kanagawa, Solarized, Catppuccin)
- ⚡ **Real-Time Feedback** - Instant visual feedback as you type
- 📊 **Detailed Statistics** - Track WPM, accuracy, and misspelled words
- 🎯 **Error Tracking** - Records mistakes even if corrected (true accuracy)
- 🔍 **Command Palette** - Quick access to all features with fuzzy search
- 🌈 **Beautiful UI** - Clean, distraction-free interface with full-screen backgrounds
- 📝 **Multi-line Support** - Practice with realistic text passages including newlines
- 📚 **Custom Texts** - Load your own practice texts from .txt files
- 🔌 **Stdin Support** - Pipe any text directly for instant practice
- 🌍 **Unicode Support** - Full support for international characters (ö, ä, ü, etc.)
- 📂 **Platform-Aware** - Uses appropriate default locations for each OS

## Installation

### From Source

```bash
git clone https://github.com/yourusername/rocketype.git
cd rocketype
go build -o rocketype ./cmd/rocketype
./rocketype
```

## Usage

Simply run the application to start typing immediately:

```bash
# Use default texts location
rocketype

# Use custom texts directory
rocketype --texts-dir ~/my-typing-texts

# Show default paths for your platform
rocketype --print-paths

# Show help
rocketype --help
```

### Default Text Locations

Rocketype uses platform-appropriate default directories:

- **Linux/BSD**: `~/.config/rocketype/texts/`
- **macOS**: `~/Library/Application Support/rocketype/texts/`
- **Windows**: `%APPDATA%\rocketype\texts\`

The directory is created automatically on first run. See [TEXTS.md](TEXTS.md) for detailed configuration.

### Practice with Custom Text via stdin

You can pipe any text directly into rocketype for instant practice:

```bash
# Practice with text from a file
cat myfile.txt | rocketype

# Practice with command output
echo "The quick brown fox jumps over the lazy dog" | rocketype

# Practice with clipboard content
pbpaste | rocketype  # macOS
xclip -o | rocketype  # Linux

# Practice with curl output
curl -s https://example.com/quote.txt | rocketype

# Practice typing git commit messages
git log --oneline -5 | rocketype
```

When stdin is provided, the piped text becomes the practice text with the name "stdin" visible in the title bar.

### Keyboard Shortcuts

**During Typing:**
- `Esc` or `Ctrl+C` - Quit application
- `Ctrl+P` - Open command palette
- `Ctrl+T` - Cycle through themes
- `Backspace` - Delete last character
- `Enter` - Type newline character

**In Results Screen:**
- `Enter` or `r` - Restart test
- `Ctrl+P` - Open command palette
- `Ctrl+T` - Change theme
- `Esc` or `Ctrl+C` - Quit application

**Command Palette:**
- `↑`/`↓` or `Ctrl+K`/`Ctrl+J` - Navigate commands
- `Enter` - Execute selected command
- `Esc` or `Ctrl+P` - Close palette
- Type to filter commands

## Themes

Rocketype includes 6 beautiful themes:

### Dark Themes
- **Default** - Respects your terminal's color scheme
- **Gruvbox** - Warm, retro-inspired with earthy tones
- **Kanagawa** - Inspired by traditional Japanese painting

### Light Themes
- **Gruvbox Light** - Warm cream background with earthy accents
- **Solarized Light** - Scientifically balanced colors for reduced eye strain
- **Catppuccin Latte** - Modern pastel aesthetic

Switch themes with `Ctrl+T` or via the command palette (`Ctrl+P` → `theme:`).

See [THEMES.md](THEMES.md) for detailed color information and screenshots.

## Custom Practice Texts

Rocketype supports loading custom typing texts from `.txt` files.

### Quick Start

1. Check your default texts directory:
   ```bash
   rocketype --print-paths
   ```

2. Add `.txt` files to that directory:
   ```bash
   # Linux example
   echo "Your custom text" > ~/.config/rocketype/texts/my-text.txt
   ```

3. Launch rocketype - it will automatically load your texts!

### Migration from Local Directory

If you have texts in a local `./texts` directory, use the migration script:

```bash
./migrate-texts.sh
```

### Detailed Configuration

For comprehensive information about:
- Platform-specific paths
- Adding custom texts
- File naming conventions
- Sharing text collections
- Advanced configuration

See **[TEXTS.md](TEXTS.md)** for complete documentation.

### Text Selection

- **Automatic random selection** - On startup, a random text is chosen
- **Command palette** - Press `Ctrl+P` and type `text:` to see all available texts
  - `text: random` - Select a random text
  - `text: [name]` - Select a specific text by name
- **Title bar** - Shows the currently active text name

### Example Text Files

The repository includes several example texts in the `texts/` directory:
- `pangrams.txt` - Classic pangrams for practicing all letters
- `hobbit.txt` - Opening from "The Hobbit" by J.R.R. Tolkien
- `tale-of-two-cities.txt` - Opening from Dickens
- `javascript-code.txt` - JavaScript code snippets for programming practice

On first run, these will be available in your platform's default texts directory.

### Text File Format

- Plain text files with `.txt` extension
- Can contain multiple lines (newlines are preserved)
- Filename (without extension) becomes the display name
- UTF-8 encoding recommended
- Empty files are ignored

**Tip:** Create specialized texts for different practice goals:
- `symbols.txt` - Practice special characters and punctuation
- `python-stdlib.txt` - Common Python standard library imports
- `git-commands.txt` - Frequently used git commands
- `medical-terms.txt` - Domain-specific vocabulary

## How It Works

### Visual Feedback

- **Gray text** - Characters you haven't typed yet
- **Green text** - Correctly typed characters
- **Red bold text** - Incorrectly typed characters
- **Yellow underline** - Current cursor position
- **Small red text above** - Shows what you actually typed when incorrect
- **Underscore `_`** - Represents a mistyped space
- **Return symbol `↵`** - Represents a newline

### Statistics

- **WPM (Words Per Minute)** - Calculated using the industry standard: 5 characters = 1 word
- **Accuracy** - Percentage of correctly typed characters
- **Misspelled Words** - Lists all words typed incorrectly, even if later corrected
  - Words are shown in the order they were first misspelled
  - Count shows how many times each word was mistyped

### Adding Custom Texts

Simply create `.txt` files in the `texts/` directory. The application automatically:
- Loads all `.txt` files on startup
- Generates commands in the palette for each text
- Displays the filename (without extension) as the text name
- Preserves line breaks and formatting

