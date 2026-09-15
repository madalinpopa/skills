# Phase 01: Local authoring

- **Feature**: [Self-contained authoring](FEATURE.md)
- **Depends on**: None
- **Outcome**: Published-by-default skills and local authoring helpers with conditional guidance.

## Current checkpoint

- **Status**: accepted
- **Approved scope**: User requested published defaults, merged authoring guidance,
  context/token optimization, and Bash automation; explicitly approved helper
  implementation after reviewing tests in the current conversation.
- **Working tree**: Isolated jj workspace /private/tmp/skills-self-contained-authoring-20260915.
  Test commit `862b4194`; implementation and documentation ready; no Go changes.
- **Last verification**: Seven helper tests and source validation passed;
  go fix, go test, and go vet passed with a per-command GOROOT override.
  Repository formatting is blocked by existing Go template placeholders.
- **Pending decision**: None; user approved proceeding after implementation review.
- **Next action**: Commit implementation, publish the focused PR, and merge after
  its checks pass, as authorized in this conversation.

## Scope and context

Read AGENTS.md, the feature overview, docs/SPEC.md Skill metadata, and the
current create-repo-skill source. Source reference material was inspected from
the installed Apache-2.0 authoring guidance; maintain its provenance/license.

## Acceptance criteria

- Initializer creates valid metadata with published default, explicit draft,
  requested resources, and correctly quoted descriptions.
- Invalid names/options and existing or symlinked destinations are rejected
  before replacing user content.
- Validator accepts repository/optional fields without changing source, and
  rejects malformed store metadata, duplicate keys, and unfinished instructions.
- UI metadata updates preserve existing interface/policy/dependencies; invalid
  values do not write files.
- Design, optimization, metadata, scripts, and validation references cover the
  source guidance without requiring another authoring skill or blanket reads.

## Change plan and review checkpoints

1. Tests: public Bash-command interface tests in the skill's tests directory.
   Proposed commit: `test(create-repo-skill): cover local authoring helpers`.
2. Test review: user explicitly approved helper implementation.
3. Implement local helpers and documentation; no Go code changes planned.
   Proposed commit: `feat(create-repo-skill): make authoring self-contained`.
4. Run helper tests, validate the skill, inspect links and scope, run repository
   verification, then present implementation results before publication.

## Verification evidence

| Check | Result |
| --- | --- |
| uv run --offline --with pyyaml python -m unittest discover -s skills/create-repo-skill/tests -v | Expected red: 7 cases, 13 failures, unimplemented local entrypoints |
| Same helper tests after implementation | Pass: 7 cases |
| Local validate-skill.sh on create-repo-skill | Pass |
| Bash syntax, local Markdown links/anchors, trailing whitespace | Pass |
| go fix ./... | Pass after unsetting inherited GOROOT per command; initial attempt failed from Go 1.27.0/1.27.1 mismatch |
| task format | Blocked by existing create-api-module Go template placeholders; no Go files changed |
| go test ./... and go vet ./... | Pass with per-command GOROOT override |
| Docker integration suite | Not run: no Go integration boundary changed; helpers use temporary directories |

The entrypoint is 584 words; detailed authoring topics are conditional reads.
This is a document count, not measured runtime token savings. Behavioral review
covers new skills, narrow metadata edits, explicit drafts, unrelated Go requests,
and failed prerequisites. No independent subagent evaluation was performed.

## Developer acceptance

- Test review: approved in current conversation.
- Implementation review: approved by the user's “proceed” in this conversation.
- Remaining work: authorized publication; use the PR state for delivery status.
