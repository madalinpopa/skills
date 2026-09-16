# Agent instructions

<Project name> is <purpose and users in one sentence>.

## Project information

- [README.md](README.md): setup and project overview.
- [docs/SPEC.md](docs/SPEC.md): behavior, constraints, and acceptance criteria.
- <Source paths and their responsibilities; link existing architecture docs.>

Read only relevant spec sections and task files. Check code and diffs before
advice or edits; preserve unrelated work. Surface spec/code conflicts before
changing the contract.

## Commands

Run from <working directory>; use the project's commands.

- Focused tests: `<command>`
- Full tests: `<command>`
- Format and lint: `<commands>`
- Build: `<command>`

## Conventions

Load only for the current task:

| Task | Guidance |
| --- | --- |
| Behavior changes and test review | [TDD](docs/conventions/tdd.md) |
| Commit proposals and delivery | [Commits](docs/conventions/commits.md) |
| Code changes and review | [Maintainable code](docs/conventions/code.md) |
| Go work | [Idiomatic Go](docs/conventions/go.md) |

## Skills to use

| Situation | Skill |
| --- | --- |
| Changes, issue selection, progress review | `workflow-mentor` |
| Go guidance, edits, or review | `modern-go-guidelines:use-modern-go` |
| Repository history when `.jj` exists | `use-jj` |
| GitHub issues, PRs, or checks | `use-gh` |

Load the selected skill when needed. Resolve required missing skills with
`use-skills-cli`; report unavailable dependencies and continue independent work.

## How we work

The user writes code by default. Follow `workflow-mentor`: one action and its
reason, then review the attempt. Give code or edit implementation only when
requested. Answer direct questions directly.

For feature work, follow the [workflow](docs/conventions/workflow.md). Use
`docs/feature.md` for scope, acceptance, and the current checkpoint; verify
recorded progress against the checkout on resume.

Respect review gates. Propose a commit message after verified work; never
commit or publish without an explicit request. Never add AI attribution.
