# Implementation Plan: Responsive Results View

**Branch**: `001-responsive-results-view` | **Date**: Tue Mar 03 2026 | **Spec**: [/specs/001-responsive-results-view/spec.md]
**Input**: Feature specification from `/specs/001-responsive-results-view/spec.md`

## Summary

User intent:
- Preferred layout: All result items should be stacked in a single column below each other for maximum clarity on any terminal size.
- Exception: If the terminal is wide enough, break out the chart displaying WPM statistics into its *own* column to the right of the other items.
- Fallback: If the terminal is not wide enough to fit the chart and the results column side-by-side, do not show the chart.
- Height constraint: Leaderboard must always show at least three rows; lower-priority sections can be suppressed to preserve this.

Technical note:
- The design must enable seamless switching between these states based on live terminal dimension checks, with no graphical dependencies beyond tcell.
- All state transitions and layouts must remain accessible and readable in standard (keyboard-focused, terminal-first) workflows.

## Technical Context

**Language/Version**: Go 1.x (latest stable; required by constitution)
**Primary Dependencies**: tcell v2 library, Go standard library only
**Storage**: N/A (no persistent results for this view)
**Testing**: go test, manual TUI walkthrough
**Target Platform**: Terminal emulator (Linux, macOS, Windows)
**Project Type**: Terminal application (TUI)
**Performance Goals**: Layout must update within 0.5 seconds of terminal resize; all metrics readable in all valid layouts
**Constraints**: No graphical dependencies, no external packages other than tcell v2; must use theme definitions; must not hardcode colors.
**Scale/Scope**: Results view supports any terminal with 40+ columns and 5+ rows; chart breakout triggers only on sufficient width (determined experimentally, e.g. ≥80 columns for chart split)

**Layout algorithm:**
- Always prefer a vertical single-column stack for all metrics as default
- If the terminal width is sufficient, break out the WPM chart to its own column on the right; otherwise, do not display chart
- Adapt layout on-the-fly based on width/height; trigger layout changes immediately as dimensions change
- Ensure the leaderboard retains at least three visible rows; suppress lower-priority content when height is constrained
- If not enough width for chart split, hide chart cleanly (no placeholder or empty space)
- All logic and rendering driven by data-oriented model, using central theme definitions

**Unknowns / NEEDS CLARIFICATION:**
- Exact column width threshold for chart breakout: [Recommend ≥80 columns; NEEDS CLARIFICATION if another value desired]
- Are any interactive controls present in the results view (e.g., open/collapse sections), or strictly read-only? [Assume read-only]
- Should chart appearance (colors/sizing) adhere to theme, or are there additional display requirements? [Assume theme-driven]
- Which sections are lowest priority to hide when enforcing minimum leaderboard rows? [Assume chart first, then misspelled words]

**Dependencies:**
- tcell v2 must remain available, compatible with all target OS terminals
- Rendering logic must avoid hardcoded layout, use measured widths

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- ✅ Terminal-First Experience: Only terminal UI and tcell used
- ✅ Minimal Dependencies: Go + tcell only
- ✅ Multi-Theme Flexibility: All colors/settings driven by theme module
- ✅ Simplicity & Maintainability: Requirement remains simple and clear

## Project Structure

### Documentation (this feature)
```
specs/001-responsive-results-view/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)
```
internal/
├── renderer.go          # Implements layout switching logic and rendering
├── typingtest.go        # Provides result data
├── stats.go             # Sources metrics for view and chart
```

## Complexity Tracking

No current Constitution violations or rejected simpler alternatives.
