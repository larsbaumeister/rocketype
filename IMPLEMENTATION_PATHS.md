# Platform-Specific Paths Implementation

## Summary

Successfully implemented platform-aware default text directories with command-line override support.

## Changes Made

### New Files

1. **`internal/paths.go`**
   - `GetDefaultTextsDir()` - Returns platform-appropriate default directory
   - `GetFallbackTextsDir()` - Returns local `./texts` fallback
   - `EnsureTextsDir()` - Creates directory if missing
   - Platform detection using `runtime.GOOS`

2. **`TEXTS.md`**
   - Comprehensive documentation for text configuration
   - Platform-specific path explanations
   - Migration guide
   - Best practices and examples

3. **`migrate-texts.sh`**
   - Bash script to help users migrate from old location
   - Safe migration with confirmation prompts
   - Color-coded output

### Modified Files

1. **`cmd/rocketype/main.go`**
   - Added `--texts-dir` flag for custom directory
   - Added `--print-paths` flag to show default locations
   - Enhanced `--help` output with platform info
   - Smart fallback logic: platform default → local fallback

2. **`internal/app.go`**
   - Updated `NewApp()` to accept `textsDir` parameter
   - Removed hardcoded `DefaultTextsDir` reference

3. **`internal/textlib.go`**
   - Removed `DefaultTextsDir` constant (now in paths.go)

4. **`README.md`**
   - Updated with platform-specific path information
   - Added command-line flag documentation
   - References to TEXTS.md for details

## Platform Defaults

| Platform | Default Path |
|----------|--------------|
| **Linux/BSD** | `~/.config/rocketype/texts/` |
| **macOS** | `~/Library/Application Support/rocketype/texts/` |
| **Windows** | `%APPDATA%\rocketype\texts\` |

All platforms also check `./texts/` as a fallback.

## Features

### 1. Automatic Directory Creation
The default directory is created automatically on first run with proper permissions (0755).

### 2. Intelligent Fallback
```
Priority order:
1. --texts-dir flag (if provided)
2. Platform default (if exists)
3. ./texts fallback
4. Built-in default text
```

### 3. Environment Variable Support
Linux respects `XDG_CONFIG_HOME`:
```bash
export XDG_CONFIG_HOME=~/my-config
# Uses: ~/my-config/rocketype/texts/
```

### 4. Helpful Commands

**Show paths:**
```bash
rocketype --print-paths
```

**Use custom directory:**
```bash
rocketype --texts-dir ~/my-texts
```

**Show help:**
```bash
rocketype --help
```

## Testing

All functionality tested on Linux:

✅ Default path detection (`/home/lars/.config/rocketype/texts`)  
✅ Directory auto-creation  
✅ `--print-paths` flag  
✅ `--help` flag with platform info  
✅ Custom directory via `--texts-dir`  
✅ Fallback to `./texts`  
✅ Migration script  

## Benefits

1. **Platform conventions** - Respects OS standards
2. **User-friendly** - Works out of the box
3. **Flexible** - Override with flag when needed
4. **Backwards compatible** - Still checks `./texts`
5. **Well documented** - Clear help and docs

## Migration Guide for Users

### From Local `./texts`

**Option 1: Use migration script**
```bash
./migrate-texts.sh
```

**Option 2: Manual copy**
```bash
# Linux
cp texts/*.txt ~/.config/rocketype/texts/

# macOS
cp texts/*.txt ~/Library/Application\ Support/rocketype/texts/

# Windows (PowerShell)
Copy-Item texts\*.txt $env:APPDATA\rocketype\texts\
```

**Option 3: Keep using local directory**
```bash
rocketype --texts-dir ./texts
```

## Future Enhancements

Potential improvements:
- Config file support (e.g., `~/.config/rocketype/config.toml`)
- Multiple text directories (like PATH)
- Text library sync across machines
- Download texts from URLs
- Built-in text marketplace/repository

## Code Quality

- ✅ Clean separation of concerns
- ✅ Well-documented functions
- ✅ Error handling
- ✅ Cross-platform compatible
- ✅ User-facing documentation
- ✅ Migration tools
