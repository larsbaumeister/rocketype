# Local Leaderboard Feature

## Overview
The local leaderboard stores the top 10 results per challenge, where a challenge is either a specific text (text mode) or a word set (word mode). Results are stored only on the local machine and are never transmitted.

Each entry stores:
- Username (from OS)
- Real name (from OS, if available)
- WPM
- Accuracy
- Timestamp

## Storage
Leaderboard data is stored in `leaderboard.json` inside the Rocketype config directory. It is written atomically using a temporary file and rename strategy, with a `.bak` backup retained during updates.

If the file is missing, an empty leaderboard is created. If the file is corrupt, it is reset to empty and the app continues without crashing.

## Sorting and Ranking
Entries are sorted and trimmed to the top 10:
1. Highest WPM
2. Highest accuracy
3. Earliest timestamp

Duplicate usernames are allowed and treated as separate entries. If the real name is missing, the UI falls back to showing the username in the Name column.

## Rendering
The results screen includes a leaderboard table that shows up to 10 entries. If fewer than 10 entries exist, only available entries are shown. Unicode names are supported and truncated safely for display.

## Quickstart Validation Notes
- Complete a typing test in text mode and verify the leaderboard updates.
- Restart the app and confirm the leaderboard persists for the same text.
- Switch to word mode and verify the leaderboard is separate.
- Test with missing real name and duplicate usernames.
