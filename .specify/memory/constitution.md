<!--
Sync Impact Report

Version change: initial → 1.0.0
Modified principles: Initial formalization
Added sections: Platform/Tech Constraints, Development Workflow
Removed sections: None
Templates requiring updates:
- .specify/templates/plan-template.md   ✅ (principles now testable/gated)
- .specify/templates/spec-template.md   ✅ (requirements must align with principles)
- .specify/templates/tasks-template.md  ✅ (task categories reflect test, theme, dependency, CLI)
- README.md                            ✅ (matches terminal, theme, minimal, simplicity focus)
Follow-up TODOs:
- TODO(RATIFICATION_DATE): explanation, fill once project history clarified

No unexplained placeholders remain.
-->

# Rocketype Constitution

## Core Principles

### I. Terminal-First Experience

All features MUST be accessible through a standard terminal interface, with no graphical dependencies beyond terminal emulation. UX, commands, and visual feedback are optimized for keyboard interaction.

Rationale: Ensures portability, accessibility, and minimalism.

---

### II. Minimal Dependencies

Only essential external packages (e.g., tcell for terminal handling, Go standard library) are permitted. New dependencies require explicit rationale and code review.

Rationale: Guarantees maintainability, security, and ease of distribution.

---

### III. Multi-Theme Flexibility

The application MUST support multiple selectable themes, allowing runtime theme changes. Themes are defined centrally, not hardcoded in logic.

Rationale: Provides user personalization and accessibility for different environments.

---

### IV. Simplicity & Maintainability

Code MUST remain simple, readable, and easy to maintain; documentation is mandatory for all exported functions and key sections. Avoid feature bloat and premature optimization.

Rationale: Enables long-term stability and community contribution.

---

## Platform and Technology Constraints

- Language: Go 1.x
- Terminal UI: Only tcell v2 and standard library allowed
- Multi-OS: Linux, macOS, Windows supported
- Text file input: Plain UTF-8 `.txt` files and piped stdin practice required
- Theme architecture: All themes defined in a dedicated module; no hardcoded colors
- Custom texts: Loaded dynamically, platform directories auto-managed

---

## Development Workflow

- All code changes via Pull Request (PR), reviewed for compliance with all principles.
- `gofmt` and `go vet` must pass before acceptance.
- Linting is required for every proposed change.
- Every PR must document principle compliance and justify complexity.
- CI/CD system should gate merges on lint and principle check.
- New dependencies must be justified in the PR and require approval.

---

## Governance

- This constitution supersedes all workflow conventions and is mandatory for all contributors.
- Amendment procedure: All changes require PR, full rationale, clear diff of principles or sections changed, and full team or maintainer approval.
- Versioning: Amendments that remove or significantly alter principles trigger MAJOR version bumps; new principles or material expansion is MINOR; refinements/clarifications are PATCH.
- Compliance reviews: All PRs and releases must explicitly document principle compliance; reviewers will block merges if requirements are not met.
- The authoritative runtime guidance is found in `README.md`; it must be manually synced with amendments to this constitution.

---

**Version**: 1.0.0 | **Ratified**: TODO(RATIFICATION_DATE): explanation, fill once project history clarified | **Last Amended**: 2026-03-03
