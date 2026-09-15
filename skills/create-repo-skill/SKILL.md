---
name: create-repo-skill
description: Creates or updates catalog skills and authoring guidance in this repository or its forks, following local naming, metadata, and packaging rules. Not for generic skill creation elsewhere or Go CLI implementation.
status: published
tags: [skills, authoring, repository]
---

# Create a repository skill

## Prepare

Use the available `skill-creator` first for design and validation; reuse its
instructions if already loaded for this task. If unavailable, report the
missing dependency before editing.

Identify the current checkout by `AGENTS.md`, `docs/SPEC.md`, and `skills/`.
Follow its guidance and inspect the diff. Author in this checkout, including
forks, rather than the installed skill or configured store. Do not change
`store.repo` or `store.branch` merely because the checkout is a fork.

## Repository essentials

- Put each skill and all its resources in `skills/<skill-name>/`; installation
  targets such as `.agents/` and `.claude/` are not source directories. Pass
  the checkout's absolute `skills/` path to Skill Creator's initializer.
- Start new skills with explicit `status: draft`. Preserve existing status
  unless a change is requested. Keep `status` and optional `tags` at the top
  level; they are store fields stripped on install.
- Keep resources self-contained, use relative links and regular files, and
  avoid source symlinks or a root `.skill-lock.json`. The CLI preserves body
  bytes, so shared instructions must work across agents.
- Apply Skill Creator's guidance without duplicating its manual. Check the
  catalog for overlap when introducing or expanding a skill's scope.

## Load details only when needed

| Task | Read |
| --- | --- |
| Choose or change a name/category | [Naming conventions](references/naming.md) |
| Create or update skill coordination | [Workflow skills](references/naming.md#workflow-skills) |
| Create a skill, change metadata/packaging, or publish | Relevant sections of [authoring rules](references/authoring.md) and the checkout's `docs/SPEC.md` |
| Change agent settings | [Agent-specific options](references/authoring.md#agent-specific-options) and current docs for the affected agent |
| Resolve validator incompatibility | [Validation details](references/authoring.md#validation-details) |

A wording-only edit needs no supporting reference unless it affects one of
these areas.

## Verify and report

Use Skill Creator's validator on a temporary shared copy with `status`, `tags`,
and `x-claude` removed; retain those fields in source. Check the source's
folder/name match, YAML types and unique keys, and relative links; run changed
scripts. Review representative requests, including a nearby non-trigger.
For workflows, verify handoffs and failure stops; for agent settings, inspect
both installed representations using the authoring rules.

Report changes, checks, limitations, and a proposed commit message. Commit or
publish only when explicitly requested; honor authorization already given.
