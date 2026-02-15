# Text Files Configuration

Rocketype loads practice texts from `.txt` files. This document explains where these files are stored and how to customize them.

## Default Text Locations

Rocketype uses **platform-appropriate** default locations for storing text files:

### Linux / BSD
```
~/.config/rocketype/texts/
```
Follows the [XDG Base Directory Specification](https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html).

You can override this with the `XDG_CONFIG_HOME` environment variable:
```bash
export XDG_CONFIG_HOME=~/my-config
# Rocketype will use: ~/my-config/rocketype/texts/
```

### macOS
```
~/Library/Application Support/rocketype/texts/
```
Follows Apple's [File System Programming Guide](https://developer.apple.com/library/archive/documentation/FileManagement/Conceptual/FileSystemProgrammingGuide/FileSystemOverview/FileSystemOverview.html).

### Windows
```
%APPDATA%\rocketype\texts\
```
Typically expands to:
```
C:\Users\YourName\AppData\Roaming\rocketype\texts\
```

## Checking Your Default Location

Run this command to see where rocketype looks for texts on your system:

```bash
rocketype --print-paths
```

Output example:
```
Default texts directory: /home/username/.config/rocketype/texts
Fallback directory: texts
```

## Using a Custom Directory

You can override the default location with the `--texts-dir` flag:

```bash
# Use texts from a specific directory
rocketype --texts-dir ~/my-typing-texts

# Use texts from current directory
rocketype --texts-dir ./texts

# Use absolute path
rocketype --texts-dir /usr/share/rocketype-texts
```

## Local Fallback

If the platform default directory doesn't exist or is empty, rocketype will automatically fall back to a `./texts/` directory in the current working directory. This is useful for:

- Development
- Portable installations
- Running from a USB drive

## Adding Your Own Texts

### 1. Find Your Texts Directory

```bash
rocketype --print-paths
```

### 2. Add Text Files

Create `.txt` files in the texts directory. Each file becomes a selectable practice text.

**Example**: Create `~/.config/rocketype/texts/shakespeare.txt`
```
To be, or not to be, that is the question:
Whether 'tis nobler in the mind to suffer
The slings and arrows of outrageous fortune,
Or to take arms against a sea of troubles,
And by opposing end them.
```

### 3. File Naming

- Files must have `.txt` extension
- Filename (without `.txt`) becomes the display name
- Example: `my-practice-text.txt` → Shows as "my-practice-text"

### 4. Access Your Texts

Press `Ctrl+P` in rocketype, then type `text:` to see all available texts, including your custom ones!

## Text File Guidelines

### Recommended Formats

**Short practice** (50-200 words):
- Good for warm-up
- Quick accuracy checks
- Specific skill practice (numbers, punctuation, etc.)

**Medium practice** (200-500 words):
- Standard typing test length
- Good for WPM measurements
- Paragraph-level flow

**Long practice** (500+ words):
- Endurance training
- Sustained accuracy practice
- Real-world typing simulation

### Best Practices

1. **Use plain text**: No formatting, just text
2. **UTF-8 encoding**: Supports international characters (ö, ä, ü, etc.)
3. **Multiline support**: Newlines are preserved and must be typed
4. **Mixed content**: Combine letters, numbers, and punctuation for comprehensive practice

### Example Text Structures

**Code practice**:
```javascript
function fibonacci(n) {
    if (n <= 1) return n;
    return fibonacci(n - 1) + fibonacci(n - 2);
}
```

**Multilingual**:
```
Übung macht den Meister.
Practice makes perfect.
La pratique rend parfait.
```

**Special characters**:
```
Email: user@example.com
Path: /usr/local/bin/app
Math: (a + b) * c = ac + bc
Symbols: $100, 50%, #hashtag, @mention
```

## Directory Structure Example

```
~/.config/rocketype/
└── texts/
    ├── quotes.txt              # Famous quotes
    ├── programming-go.txt      # Go code samples
    ├── programming-python.txt  # Python code samples
    ├── german-practice.txt     # German umlauts
    ├── numbers-practice.txt    # Number practice
    ├── punctuation.txt         # Punctuation drills
    ├── literature-hobbit.txt   # Book excerpts
    └── custom-words.txt        # Personal word lists
```

## Migrating from Old Location

If you were using the old `./texts/` directory, you can migrate your files:

### Linux/macOS
```bash
mkdir -p ~/.config/rocketype/texts  # or appropriate macOS path
cp texts/*.txt ~/.config/rocketype/texts/
```

### Windows (PowerShell)
```powershell
New-Item -Path "$env:APPDATA\rocketype\texts" -ItemType Directory -Force
Copy-Item -Path "texts\*.txt" -Destination "$env:APPDATA\rocketype\texts\"
```

## Troubleshooting

### "No texts found" Error

**Cause**: The texts directory is empty or doesn't contain `.txt` files.

**Solution**: Add at least one `.txt` file to your texts directory, or rocketype will use the built-in default text.

### Permission Denied

**Cause**: No write permissions for the default directory.

**Solution**: Use a custom directory where you have write access:
```bash
rocketype --texts-dir ~/my-texts
```

### Wrong Directory

**Cause**: Multiple rocketype installations or incorrect configuration.

**Solution**: Check your actual paths:
```bash
rocketype --print-paths
```

### Texts Don't Appear

**Checklist**:
1. Files have `.txt` extension
2. Files contain text (not empty)
3. Files are in the correct directory
4. Files have read permissions
5. Restart rocketype after adding files

## Sharing Text Collections

You can share your text collections by sharing the directory or individual files:

```bash
# Create a shareable archive
cd ~/.config/rocketype
tar czf rocketype-texts.tar.gz texts/

# Share or backup
cp rocketype-texts.tar.gz ~/Dropbox/

# Restore on another machine
tar xzf rocketype-texts.tar.gz -C ~/.config/rocketype/
```

## Advanced: Version Control

Keep your text collection in git:

```bash
cd ~/.config/rocketype/texts
git init
git add *.txt
git commit -m "Initial text collection"
git remote add origin https://github.com/yourusername/rocketype-texts
git push -u origin main
```

Now you can sync across machines and track changes!

---

## Quick Reference

| Action | Command |
|--------|---------|
| Show default path | `rocketype --print-paths` |
| Use custom directory | `rocketype --texts-dir ~/my-texts` |
| Add new text | Create `.txt` file in texts directory |
| Select text | Press `Ctrl+P` → type `text:` |
| Reset to default | Delete custom texts, use built-in default |

For more help, see the main [README.md](./README.md) or run `rocketype --help`.
