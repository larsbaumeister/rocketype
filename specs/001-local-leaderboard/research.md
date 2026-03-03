# Local Leaderboard Research Summary

## 1. Go Cross-platform Extraction of Username and Real Name

- Use [`os/user`](https://pkg.go.dev/os/user) package: `user.Current()` retrieves local username and (if available) real name.
- Works cross-platform: Linux, macOS, Windows; not supported for mobile/web.
- Example:
  ```go
  import "os/user"
  u, err := user.Current()
  if err == nil {
      fmt.Println(u.Username, u.Name)
  }
  ```
- Always fall back to username if name missing/empty.
- Respect privacy: never send/display outside local scope.
- Reference: [Stack Overflow](https://stackoverflow.com/q/11795190)

---

## 2. Recommended File Storage Format for Local Leaderboard

- Most robust Go standard: JSON via [`encoding/json`](https://pkg.go.dev/encoding/json).
- Human-readable, structured, supports Unicode; easy for error recovery.
- Write atomically: to temp file then rename, to prevent data loss.
- Example format:
  ```json
  [
    {"score": 112, "username": "alice", "name": "Alice Smith"},
    {"score": 110, "username": "bob", "name": "Robert Brown"}
  ]
  ```
- Use OS-specific Rocketype directory for storage.
- Reference: [Go Blog: JSON](https://blog.golang.org/json)

---

## 3. Handling File Corruption

- Always check for read/unmarshal errors; if file corrupt, reset or recover with backup.
- Atomic update recommendation: write new file, then use [`os.Rename`](https://github.com/golang/go/wiki/AtomicFileWrite) to replace; backup original before overwrite.
- If less than 10 results, display what exists and handle gracefully.

---

## 4. Unicode Handling for Names/Scores

- Go strings are UTF-8 encoded by default; robust for all international text.
- Use standard library: `utf8`, `strings`, `unicode`.
- Example: `utf8.RuneCountInString(str)` for length.
- Handles emojis, diacritics, other scripts reliably.
- Reference: [Unicode in Go](https://www.slingacademy.com/article/exploring-unicode-and-utf-8-in-go-strings/), [Go Cookbook Unicode](https://go-cookbook.com/snippets/strings/working-with-unicode-strings)

---

## 5. Privacy, Minimal Dependencies

- Do not sync, log, or collect extra info; only what’s needed for scores.
- Store leaderboard in user-only Rocketype dir, as plain JSON.
- Use Go stdlib only (`os/user`, `encoding/json`, `os`).
- Reference: [Go Privacy](https://stackoverflow.com/q/40384257)

---

## Key Actionable Notes
- Use os/user for name and username. JSON local file (atomic write) for leaderboard. Always check, recover, and respect privacy. Unicode is fully supported — no extra libs required.
- All sources and practices are current as of Feb 2026.
