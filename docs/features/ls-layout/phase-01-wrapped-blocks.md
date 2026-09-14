# Phase 01: One wrapped block per skill

- **Feature**: [Feature agreement](FEATURE.md)
- **Depends on**: None
- **Outcome**: Every `ls` view prints one block per skill: name and tags (or
  agents) on the header line, the full description wrapped below.

## Current checkpoint

- **Status**: awaiting test review
- **Current step**: Step 1, test review is next
- **Approved scope**: Layout and piped behavior agreed on 2026-09-14; go-ahead
  for the step 1 tests given on 2026-09-14 ("Then proceed with the tests")
- **Working tree**: on top of `lutuzxmm` (plan docs): new
  `cmd/ls_blocks_test.go`; adapted `cmd/ls_test.go` and `cmd/damaged_test.go`
- **Last verification**: `go test ./cmd/` fails to build (`block` and
  `printBlocks` undefined); with the internal test set aside, the three
  adapted command tests fail on the line shape
- **Blocker or pending decision**: Developer review of the tests; approval of
  the `golang.org/x/term` dependency
- **Next action**: On approval, commit the tests and implement step 1

## Scope and required context

- **Included**: Width detection for stdout; the block renderer; the spec
  examples update.
- **Non-goals**: Cutting text; any new flag; other commands' output.
- **Expected affected areas**:
  - `cmd/ls.go`: `printRows` becomes a block renderer taking a width.
  - `cmd/render.go`: terminal width helper next to `isTerminal`.
  - `go.mod`: `golang.org/x/term` becomes a direct dependency.
  - `docs/SPEC.md`: `ls` examples in Commands.

| Reference | Why it is needed |
| --- | --- |
| `docs/SPEC.md` Commands, `ls` examples | Current examples to replace with the block layout |
| `docs/SPEC.md` Output | "One line per skill" wording to reconcile with one block per skill |
| `cmd/ls.go` `printRows`, `listInstalled` | Current tabwriter layout and the attention lines to keep |
| `cmd/render.go` `isTerminal` | Existing terminal detection to reuse for width |
| `cmd/ls_test.go` | Existing assertions on row contents to adapt to blocks |
| `golang.org/x/term` `GetSize` docs | Verify the current signature before use |

## Acceptance criteria

- [ ] Given a long description and width 80, `skills ls` prints the header
  `  <name>   <tags>`, then the whole description on lines of at most 80
  characters indented by six spaces, in original word order.
- [ ] Given two skills, exactly one blank line separates the blocks and none
  follows the last block.
- [ ] Given empty tags, the header is the name alone with no trailing spaces.
- [ ] Given a word longer than the width, it is printed whole on its own line.
- [ ] Given `--local` or `--global`, the header shows the agents and the
  description block shows the installed description; a damaged lock shows the
  attention line and the command exits 3.
- [ ] Given a multi-byte description, wrapping counts runes and keeps valid
  UTF-8.

## Change plan and review checkpoints

### Step 1: Render one wrapped block per skill

- **Change and reason**: Replace the tabwriter rows with a block renderer.
  `printRows` takes the output width and, for each row, writes the header line
  and the wrapped description. The width comes from `term.GetSize` on stdout
  when it is a terminal, else 80. The store and installed views share it.
- **Test cases**:
  - `cmd` internal table test of the block renderer with a fixed width: long
    description wraps under 80 with a six-space indent; blank line between
    blocks and none after the last; empty tags header; over-long word kept
    whole; multi-byte text wrapped by runes.
  - `cmd_test` `TestLs_listsPublishedSkills` and `TestLs_installedScopes`
    adapted to the block shape: header line holds name and tags or agents, the
    next line holds the description.
  - Expected failure: the renderer does not exist yet (compile error in the
    internal test); the adapted command tests fail on the line shape.
- **Test commit proposal**: `test(ls): cover one wrapped block per skill`
- **Implementation commit proposal**: `feat(ls): print one wrapped block per skill`

1. Write the tests above and run `go test ./cmd/...`; confirm the failure is the
   missing behavior.
2. **Stop for test review.**
3. After approval, add `golang.org/x/term`, the width helper, the block
   renderer, and the spec examples.
4. Run the `AGENTS.md` verification, including `task lint`.
5. **Stop for implementation review.**

## Verification evidence

| Step / criterion | Command or inspection | Result and relevant evidence | Tested state |
| --- | --- | --- | --- |
| Not run | | | |

## Developer acceptance

- **Decision**: Pending
- **Reference and scope**: None
- **Remaining work**: Step 1
