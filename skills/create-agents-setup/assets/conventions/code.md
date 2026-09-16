# Clean and maintainable code

Keep changes within the agreed outcome and existing package responsibilities.
Verify current official documentation for external APIs being changed.

- Choose descriptive names and direct control flow. Keep functions and modules
  cohesive; extract a helper when it clarifies a responsibility or real reuse.
- Add abstractions and dependencies for a demonstrated need. Avoid speculative
  extension points, generic frameworks, and unrelated cleanup.
- Validate external input and filesystem or service results at boundaries.
  Trust established internal invariants; avoid repeated defensive checks.
- Handle errors explicitly with useful context. Keep ownership of resources,
  cleanup, and side effects visible. Do not hide failures or log secrets.
- Comments explain constraints and reasons the code cannot express. Update
  misleading comments; avoid narrating obvious statements.
- Follow the project's formatter and linter for mechanical style. Use tests
  to protect behavior during refactoring; do not rewrite working code solely
  to match a personal preference.

Explain what changed, why, and how it was verified. Distinguish evidence from
assumptions and surface unresolved contract or scope decisions.
