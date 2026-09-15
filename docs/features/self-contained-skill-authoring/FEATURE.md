# Feature: Self-contained repository skill authoring

Make create-repo-skill the complete local authoring workflow, with published
new skills, local helpers, and progressively loaded design/optimization guidance.

- **Active phase**: [01: Local authoring](phase-01-local-authoring.md)
- **Feature acceptance**: Approved in the current conversation after implementation review.

## Scope

- Include self-contained design, resource selection, creation, UI metadata,
  validation, conditional evaluation, and context/token optimization guidance.
- Default new skills to published; preserve existing status and permit explicit drafts.
- Provide Bash initialization/validation/metadata commands with Python/PyYAML
  for structured YAML. Preserve existing content and metadata on errors/updates.
- Update root AGENTS.md to remove the external skill prerequisite.
- Non-goals: change Go CLI behavior, rewrite other skills, add automatic
  publication, or install runtime dependencies silently.
- Spec: [Skill metadata](../../SPEC.md#skill-metadata).

## Acceptance criteria

- [x] Ordinary authoring needs no external authoring skill installation.
- [x] New skills default to published and explicit draft remains available.
- [x] Helpers pass the phase's behavioral tests and preserve existing files.
- [x] Guidance covers all source authoring topics with relevant-only loading.
- [x] Documentation, helper checks, and repository verification are complete;
  the existing formatter limitation is recorded in the phase evidence.

## Agreed decisions

| Decision | Reason | Approval |
| --- | --- | --- |
| Published default; merge source guidance and optimization practice | User request | Current conversation, 2026-09-15 |
| Bash entrypoints with structured YAML parser | Reliable repetitive operations | Helper implementation approved after test review |
| Keep source attribution in NOTICE/LICENSE only | Retain source licensing without runtime dependency | Implementation detail |

## Phase roadmap

| Phase | Outcome | Depends on | Status | Phase file |
| --- | --- | --- | --- | --- |
| 01 | Self-contained authoring workflow and helpers | None | accepted | [Phase 01](phase-01-local-authoring.md) |
