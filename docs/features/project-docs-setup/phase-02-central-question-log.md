# Phase 02: Central question log

- **Feature**: [Project documentation setup](FEATURE.md)
- **Depends on**: Existing initializer and question templates
- **Outcome**: Initialize both question documents with the same preservation and
  preflight checks as the other project documents.

## Current checkpoint

- **Status**: awaiting implementation review
- **Approved scope**: User's “proceed” approves the existing failing tests and
  adding the two question paths to the initializer's file list.
- **Working tree**: Initializer fix and phase checkpoint documentation only.
- **Last verification**: Full skill checks passed: 11 skills validated, 14 tests
  passed; ShellCheck and Bash syntax checks passed.
- **Pending decision**: Implementation review after verification.
- **Next action**: Developer reviews the implementation diff and results.

## Scope and acceptance criteria

- [x] Create `docs/QUESTIONS.md` and `docs/conventions/QUESTIONS.md` from existing
  assets, preserving existing bytes on rerun.
- [x] Reject symlinks and directory collisions at either path before writing.
- [x] Preserve existing-layout support and other initializer behavior.

Existing tests in
[test_initializer.py](../../../skills/create-project-docs/tests/test_initializer.py)
cover these behaviors; no new tests or template changes are necessary.

## Verification evidence

| Check | Result |
| --- | --- |
| Before fix: `python3 -B -m unittest discover -s skills/create-project-docs/tests` | Expected failure: 7 tests, 9 failures caused by the missing file-list entries |
| After fix: same helper test command | Passed: 7 tests |
| `shellcheck skills/create-project-docs/scripts/init-project-docs.sh` and `bash -n skills/create-project-docs/scripts/init-project-docs.sh` | Passed |
| Installed authoring validator, via offline uv with PyYAML | Passed: valid create-project-docs skill |
| `task test:skills` | Passed: 11 skills validated, 14 tests passed; initial sandbox attempt failed on uv cache lock, permitted retry passed |

Results cover the one-line initializer fix. No Go files or external libraries
changed; Go and integration checks are unrelated and were not run. Preparation
completed with `skills sync` (sandbox denial, permitted retry passed) and
`skills update create-repo-skill --agent codex` (0 changed).

## Developer acceptance

- **Decision**: Pending implementation review; no commit or publication authorized.
- **Implementation commit proposal**: `fix(create-project-docs): initialize question documents`
