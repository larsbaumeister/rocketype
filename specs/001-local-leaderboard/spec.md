# Feature Specification: Local Leaderboard

**Feature Branch**: `001-local-leaderboard`  
**Created**: 2026-03-03  
**Status**: Draft  
**Input**: User description: "i want to extend this app with a leaderboard. assume the user of that is currently typing is the os user. copy both the unix username and if present the real name from the os current user to the leaderboard. the leaderboard is local to a specific text or a specific word mode. there is only the local to challenge leaderboard not global (makes no sense). the leaderboard only containes the top 10 no more. it is rendered as a table in the results page that shows when a typing challange is completed. once reaching the result, the result is automatically added to the leaderboard."

## User Scenarios *(mandatory)*

### User Story 1 - Submit and See Result (Priority: P1)

After finishing a typing challenge, the user's result (with WPM, accuracy, and the current OS username/real name) is automatically added to the per-challenge leaderboard and the top 10 scores are displayed as a table on the results page.

**Why this priority**: Enables instant feedback and competition; core value for user motivation.

**Acceptance Scenarios**:
1. **Given** a new user completes a challenge, **When** the results screen shows, **Then** the leaderboard includes this result and lists the top 10, if available.
2. **Given** an existing leaderboard of 10 results, **When** a new user gets a better score, **Then** their score is correctly inserted and the lowest is dropped.

---

### User Story 2 - Handle User Data Edge Cases (Priority: P2)

If the user's real name is unavailable, the leaderboard displays only the username. Duplicate usernames/real names are supported and results are tied to scores, not only name.

**Why this priority**: Ensures robustness and inclusivity.

**Acceptance Scenarios**:
1. **Given** a user with no real name, **When** they complete a challenge, **Then** the leaderboard shows only their username.
2. **Given** two users with the same chosen username (from OS), **When** both post scores, **Then** both are shown as separate lines with their respective scores.

---

### User Story 3 - Persist and Reload Leaderboard Data (Priority: P3)

Leaderboard entries for each text/word-mode are saved locally and persist across application restarts.

**Why this priority**: User progress and competition history drives continued usage.

**Acceptance Scenarios**:
1. **Given** a leaderboard with 10 results, **When** the app is closed and reopened, **Then** the leaderboard for that challenge is unchanged.

---

### Edge Cases
- What happens if multiple users get the same score? The one who achieved it first (earlier timestamp) stays ranked higher.
- How does system handle non-ASCII usernames/real names? Must display and store Unicode safely.
- What if the local leaderboard file/record is missing/corrupted? App initializes an empty leaderboard and does not crash; corruption is handled gracefully with an error message or silent recovery.
- How does the leaderboard behave if less than 10 results are available? Only available entries are displayed, up to 10.
- What if the OS does not allow real name lookup? Only the username is shown, no error is displayed.

---

## Requirements *(mandatory)*

### Functional Requirements
- **FR-001**: System MUST automatically add the current user's typing challenge results (including WPM, accuracy, timestamp, UNIX username, and real name if available) to the appropriate leaderboard after each challenge completion.
- **FR-002**: System MUST store and display at most the top 10 results (per text or word-mode), ranked by highest WPM (accuracy used to break ties, then earliest timestamp if still tied).
- **FR-003**: System MUST persist the leaderboards locally on disk so that results are never lost after app restart.
- **FR-004**: System MUST render the leaderboard as a table in the results UI after every typing challenge.
- **FR-005**: System MUST fall back gracefully if the user's real name is unavailable, displaying only the username; Unicode is supported for all display/info.
- **FR-006**: System MUST function entirely locally, with no network or remote/global sharing of leaderboard data.
- **FR-007**: System MUST accept and display multiple entries for the same username/real name if scores are different and within top 10.
- **FR-008**: System MUST initialize and handle missing/corrupted leaderboard files by resetting them to empty (with recovery notice in logs/UI, no crash or data leak).

### Key Entities
- **Leaderboard**: Represents the top 10 results for a given text or word mode. Contains a sorted list of LeaderboardEntry, is loaded and saved from local storage using a unique key per text or mode.
- **LeaderboardEntry**: Record of a completed challenge for a specific user; includes UNIX username, real name (if any), WPM, accuracy, timestamp of completion, and challenge key (text/mode).

## Success Criteria *(mandatory)*

### Measurable Outcomes
- **SC-001**: 100% of completed typing challenges are reflected in the displayed top 10 leaderboard (if in top 10).
- **SC-002**: Leaderboard UI/table renders correctly in all tested results screens (correct columns, user display, sorted order, Unicode support).
- **SC-003**: After restarting the application, leaderboard for a given text/mode persists without data loss or corruption.
- **SC-004**: System never makes any network requests or attempts remote sync for leaderboards.
- **SC-005**: No leaderboard data is lost or incorrectly displayed during file corruption or low disk space events (system logs recovery and recovers gracefully, never crashes).
- **SC-006**: Users with duplicate usernames and/or no real names are fully supported in all leaderboard and UI operations.

---

### Assumptions
- The OS username/real name are available via the Go `os/user` package or standard equivalent.
- No global leaderboard; all persistence is device-local and per challenge.
- Scores are only entered for actual users after typing a challenge (no artificial/manual edits).
- Only the top 10 scores per challenge are kept; others are discarded.
- There is only one table per text or word mode (no "all time" or "daily" split).
- All necessary filesystem permissions for reading/writing local leaderboard data are handled gracefully if unavailable.
