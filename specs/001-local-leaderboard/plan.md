# Implementation Plan: Local Leaderboard Feature

**Feature branch:** 001-local-leaderboard

---
## Overview
Implement a robust, privacy-preserving local leaderboard for Rocketype that:
- Tracks top 10 scores per mode/text
- Displays username and (if present) real name, sourced from OS
- Persists data locally in user Rocketype directory
- Renders a leaderboard table on results screen
- Handles edge cases (duplicates, ties, missing fields, data corruption, <10 results, unicode)

---
## Architecture & Key Components
1. **Leaderboard Data Model**
   - Struct: Score, username, name, timestamp, mode/text
   - File: leaderboard.json (JSON array, atomic write)
2. **User Info Extraction**
   - Use os/user for username and full name (cross-platform)
   - Fallback to username if name missing
3. **Persistence & Error Handling**
   - Atomic file writes (temp file, rename)
   - Backup old file before overwrite
   - Gracefully degrade on read/corruption (restore, reset, partial render)
4. **UI Rendering (Results Page)**
   - Render top 10 table with username, name, score, date
   - Distinct handling for ties, missing fields, <10 results
   - Unicode rendering compatible (tcell, Go stdlib)
5. **Edge Cases**
   - Edge cases: duplicate names, ties, corrupted file, missing fields, non-ASCII
   - Backwards compatibility: no errors if old file layout

---
## Development Sequence
1. Define leaderboard struct, contract; validate with test data
2. Implement data loading, saving, error recovery for leaderboard.json
3. Integrate user info extraction (os/user): test on all platforms
4. Connect stats-logic for updating leaderboard after results
5. Render leaderboard as table on results screen (UI, tcell)
6. Comprehensive edge-case handling

---
## Edge Cases
- Duplicate usernames: allow, differentiate by timestamp
- Ties: sort by score, then by earliest timestamp
- File corruption: reset, restore, retain best effort render
- <10 results: show partial leaderboard
- Unicode/emoji: test with names and scores, ensure correct display
- Missing fields: fallback to username, show placeholder
- Data migration: handle absent/older leaderboard.json gracefully

---
## References & Resources
- Spec: specs/001-local-leaderboard/spec.md
- Research: specs/001-local-leaderboard/research.md
- Go docs: os/user, encoding/json, atomic file write, unicode

---
## Responsible Practices
- Privacy: local-only, no extra collection
- Minimal dependencies: Go stdlib exclusively
- Platform support: Linux, macOS, Windows
- Maintainability: clear separation of business/UI/persistence per constitution

---
## Completion
When all steps pass, feature branch is ready for review/integration.
