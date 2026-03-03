# Feature Specification: Responsive Results View

**Feature Branch**: `001-responsive-results-view`  
**Created**: Tue Mar 03 2026  
**Status**: Draft  
**Input**: User description: "make the results view responsive. it now displays many different things. depending on, if the user has a wide terminla or high termial, the layout shouls adjust"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Responsive Results With Optional Chart (Priority: P1)

When a user completes a typing test, the results view stacks all primary metrics in a single vertical column. If the terminal is wide enough, the WPM chart is displayed in a separate right-side column. If the terminal is not wide enough, the chart is hidden entirely.

**Why this priority**: The default single-column layout is the most reliable and readable across terminals, while the optional chart enhances insight when space permits.

**Independent Test**: Can be fully tested by running a typing test and resizing the terminal to validate single-column layout, two-column chart breakout, and chart hiding when narrow.

**Acceptance Scenarios**:
1. **Given** a standard terminal width, **When** the typing test ends, **Then** all results items appear stacked in a single column and the chart is hidden.
2. **Given** a wide terminal, **When** the typing test ends, **Then** the WPM chart appears in its own right-side column alongside the single-column results.
3. **Given** a terminal too narrow for a chart split, **When** the typing test ends, **Then** the chart is not shown and the single-column results remain intact.

---

### User Story 2 - Minimum Leaderboard Visibility (Priority: P2)

When terminal height is reduced, the results view must preserve at least three leaderboard rows, even if that means hiding lower-priority sections (such as the chart or extra spacing).

**Why this priority**: Leaderboard visibility is a core part of results feedback and should remain visible even in constrained terminals.

**Independent Test**: Can be tested by reducing terminal height and confirming that at least three leaderboard rows are still displayed.

**Acceptance Scenarios**:
1. **Given** a short terminal height, **When** the results screen is displayed, **Then** the leaderboard shows at least three rows.

---

### User Story 3 - Responsive Adaptation to Terminal Resize (Priority: P3)

If the user resizes the terminal window during or after the test (either wider or taller), the results view immediately adapts its layout to use available space optimally for visibility and readability.

**Why this priority**: Ensures the feature remains usable and dynamic, adapting in real time to user preferences or device changes.

**Independent Test**: Results view updates layout within 0.5 seconds after terminal is resized; can be tested by resizing window and observing layout change.

**Acceptance Scenarios**:
1. **Given** any terminal orientation, **When** the terminal is resized, **Then** results layout adapts accordingly, with result items always visible and readable.

---

### Edge Cases

- What happens when the terminal is smaller than 40 columns or 5 rows? Results view collapses to a single-column minimal mode with only essential metrics shown.
- How does the system handle missing result data (e.g., one metric unavailable)? Layout still displays all available items, skipping missing.
- User rapidly resizes terminal from wide to narrow and back: results view updates each time without layout breakage.
- When height is constrained, chart and low-priority sections may be suppressed to keep at least three leaderboard rows visible.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST detect terminal dimensions (width and height) and select results layout accordingly.
- **FR-002**: System MUST display all primary result metrics in a single vertical column by default.
- **FR-003**: System MUST display the WPM chart in a separate right-side column when the terminal is wide enough for a split layout.
- **FR-004**: System MUST hide the WPM chart entirely when the terminal is not wide enough for a split layout.
- **FR-005**: System MUST preserve at least three leaderboard rows even when terminal height is reduced.
- **FR-006**: System MUST adapt layout dynamically if the terminal is resized during the results view.
- **FR-007**: System MUST enter a minimal single-column mode if the terminal size falls below 40 columns or 5 rows.
- **FR-008**: All available result types MUST be displayed in each layout orientation; if some items are missing, placeholders or indicators may be shown.

### Key Entities

- **ResultsView**: Contains the set of result metrics, tracks layout orientation, and stores available space parameters (width, height).
- **ResultItem**: Represents a single result metric (e.g., WPM, accuracy). Attributes: label, value, display priority.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 95% of users report results are clearly readable and accessible after resizing the terminal.
- **SC-002**: Results view layout adapts within 0.5 seconds of terminal resize.
- **SC-003**: No user-reported cases where results view becomes unusable (information lost or layout broken) on supported terminals.
- **SC-004**: In constrained height scenarios, the leaderboard shows at least three rows.
- **SC-005**: The WPM chart appears only when the terminal is wide enough for a split layout; otherwise it is hidden.
