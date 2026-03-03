# Tasks: Responsive Results View

## Phase 1: Setup
- [X] T001 Create directory structure in internal/ and specs/001-responsive-results-view
- [X] T002 Initialize Go module and tcell dependency in go.mod
- [X] T003 Add initial theme definitions to internal/theme.go

## Phase 2: Foundational
- [ ] T004 Scaffold ResultsView data model in internal/renderer.go
- [ ] T005 Scaffold ResultItem entity in internal/renderer.go
- [X] T006 Prepare layout config and dimension detection logic in internal/renderer.go

## Phase 3: User Story 1 (Priority: P1) - Wide Terminal Optimized Results
- [X] T007 [P] [US1] Implement one-column vertical stack layout for all result metrics in internal/renderer.go
- [X] T008 [P] [US1] Add logic for chart breakout to the right only if terminal is wide enough in internal/renderer.go
- [X] T009 [US1] Implement hiding of chart if not enough width in internal/renderer.go
- [X] T010 [US1] Integrate WPM chart stats source in internal/stats.go

## Phase 4: User Story 2 (Priority: P2) - Tall Terminal Optimized Results
- [X] T011 [US2] Implement vertical stacking for tall terminals in internal/renderer.go
- [X] T012 [US2] Integrate edge case handling for very small terminals in internal/renderer.go

## Phase 5: User Story 3 (Priority: P3) - Responsive Adaptation
- [X] T013 [US3] Implement runtime layout adaptation on terminal resize in internal/renderer.go
- [X] T014 [P] [US3] Ensure all layout transitions are seamless and readable in internal/renderer.go

## Final Phase: Polish & Cross-Cutting Concerns
- [X] T015 Add comments to all new exported functions per constitution in internal/renderer.go, internal/stats.go
- [X] T016 Ensure theme-driven rendering (no hardcoded colors) in internal/theme.go
- [X] T017 Add manual test walkthrough instructions to specs/001-responsive-results-view/quickstart.md
- [X] T018 Run gofmt, go vet, and lint checks

## Dependencies
- US1 must be complete before US2 and US3
- US2 and US3 can be developed/tested in parallel once US1 is implemented

## Parallel Execution Examples
- T007 and T008 ([US1]) can be implemented in parallel (separate conditional logic)
- T013 ([US3]) can be implemented in parallel with polish phase tasks

## Implementation Strategy
- Focus MVP on User Story 1: One-column layout with responsive breakout chart logic for wide terminals
- Add vertical stacking and real-time resize responsiveness as incremental enhancements
- Deliver each story as an independently testable increment
