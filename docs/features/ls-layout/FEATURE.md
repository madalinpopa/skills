# Feature: Readable `ls` blocks

`skills ls` aligns its columns with `text/tabwriter`, so the description
column is padded to the longest description in the list. Real skill
descriptions run to several hundred characters, so the tags column lands far
off screen and every row wraps into unreadable blocks. The point of `ls` is to
read each skill and decide what to install, so the outcome is one block per
skill with the full description wrapped to the terminal width.

- **Active phase**: [Phase 01: one wrapped block per skill](phase-01-wrapped-blocks.md)
- **Feature acceptance**: Pending

## Scope

- **Included**: Both `ls` views (store, `--local`, `--global`) print one block
  per skill: a header line with the name and the trailing column (tags, or
  agents), then the full description wrapped and indented below, with a blank
  line between skills.
- **Non-goals**: Cutting descriptions; a `--json` or machine-readable format;
  display-width handling for wide (CJK) glyphs; changes to other commands.
- **Constraints**: One block per skill, never one per file. `cmd` owns
  rendering; domain packages do not print. Terminal detection follows the
  existing `isTerminal` helper used for color.
- **Spec references**: [Commands, `ls` examples](../../SPEC.md#commands),
  [Output](../../SPEC.md#output).

## Agreed layout

```
$ skills ls
  create-repo-skill   skills, authoring, repository
      Creates or updates skills in this repository or its forks using the
      checkout's naming, metadata, and packaging rules. Use when authoring a
      skill catalog with this repository's layout and conventions or maintaining
      its skill authoring guidance; not for generic skill creation elsewhere or
      Go CLI implementation.

  demo   demo, example, testing
      Demonstrates Agent Skills packaging with a bundled script, reference, and
      output template. Use when exploring this demo, learning the skill folder
      layout, or smoke-testing a skills CLI installation.

$ skills ls --local
  create-repo-skill   agents, claude
      Creates or updates skills in this repository or its forks using the
      ...
```

- Header: two-space indent, name, three spaces, tags or agents. Empty tags
  leave the header as the name alone.
- Description: six-space indent, word-wrapped so no line exceeds the width.
  A word longer than the width stays whole on its own line.
- Width: the terminal width when stdout is a terminal, otherwise 80.
- One blank line between blocks, none after the last.
- Installed views keep their attention lines (`! damaged lock: ...`,
  `! description unavailable`) as the description block.

## Acceptance criteria

- [ ] In a terminal, every description line fits the terminal width and the
  whole description is printed. Evidence: phase 01, step 1 tests.
- [ ] Piped output uses the same blocks wrapped at 80 columns.
  Evidence: phase 01, step 1 tests.
- [ ] `--local` and `--global` print the agents on the header line and the
  installed description below, including attention lines and exit code 3.
  Evidence: phase 01, step 1 tests.
- [ ] `docs/SPEC.md` shows the block layout in the `ls` examples.
  Evidence: phase 01 implementation review.

## Agreed decisions

| Decision | Reason | Approval or pending question |
| --- | --- | --- |
| Layout A: name and tags header, wrapped description below | Full description is the point of `ls`; header keeps name and tags scannable | Developer choice, 2026-09-14 |
| Piped output uses the same blocks wrapped at 80 columns | One code path and one documented shape | Developer choice, 2026-09-14 |
| Read the width with `golang.org/x/term` | Small Go-team package, works on POSIX and Windows; `x/sys` is already an indirect dependency | Pending developer approval |
| No maximum line width on wide terminals | Simplest rule; can be revisited if long lines read badly | Proposal |

## Phase roadmap

| Phase | Outcome | Depends on | Status | Phase file |
| --- | --- | --- | --- | --- |
| 01 | One wrapped block per skill in every `ls` view | None | planned | [phase-01-wrapped-blocks.md](phase-01-wrapped-blocks.md) |
