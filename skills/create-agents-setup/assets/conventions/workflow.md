# Mentoring and feature work

`workflow-mentor` owns conversational steps, reviewing attempts, issue selection,
and delivery. Use its current-stage reference when needed; do not copy its
procedures into project instructions. Project TDD and commit documents own
engineering gates and publication limits.

## Feature record

Use `docs/feature.md` for the active feature unless the project already has a
feature-document convention. Reuse existing `docs/features/<name>/FEATURE.md`
and phase files where present; never maintain a competing progress record.
Create a feature record only for requested feature planning, not during setup.
Use the project's FEATURE and PHASE templates when available.

The record holds the problem, agreed scope, non-goals, acceptance criteria,
decisions, relevant spec links, and next action. Present the plan before
implementation and preserve the user's explicit decisions. `docs/SPEC.md`
defines product behavior; the feature record defines this work's agreement.
Engineering conventions stay in AGENTS.md and its linked documents.

On resume, read the overview, current checkpoint or active phase, and linked
spec sections. Load historical phases only for a specific dependency.
Use workflow-mentor's checkpoint procedure at review boundaries and session
end, writing into this existing record. Link the related issue rather than
copying its history. Acceptance criteria need evidence and developer review;
checked tasks alone do not establish completion.
