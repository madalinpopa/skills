---
name: create-repo-skill
description: Creates or updates skills in the madalinpopa/skills repository using its naming, metadata, and packaging rules. Use when authoring this repository's skill catalog or maintaining its skill authoring guidance; not for generic skill creation elsewhere or Go CLI implementation.
status: published
tags: [skills, authoring, repository]
---

# Create a repository skill

## Start with skill-creator

Invoke the available `skill-creator` skill first (`$skill-creator` in Codex)
and read its instructions. Use it for skill design and validation, then apply
this repository's rules. If it is unavailable, report the missing dependency
before creating or editing skill files.

Work in the target `madalinpopa/skills` checkout. Resolve `AGENTS.md`,
`docs/SPEC.md`, and the destination `skills/` from that repository root, not
from an installed copy of this skill. Follow the target repository's agent
instructions and inspect the current diff before editing.

## Choose the outcome and name

Inspect existing skills for overlap. Define the requested job, its inputs,
expected result, and when it should be selected. Keep ordinary skills focused
on one reusable job.

Read [references/naming.md](references/naming.md) when choosing or changing a
name or category. It contains the category legend and the `workflow-*`
convention for skills that coordinate other skills. For a workflow, define
component availability, handoffs, and completion before writing its steps.

## Write the skill

Read [references/authoring.md](references/authoring.md) before creating or
editing skill files. It covers the base specification, CLI frontmatter,
Claude and Codex options, packaging, and review checks.

Create `skills/<skill-name>/SKILL.md` and every supporting resource inside that
folder. When using skill-creator's initializer, pass the target repository's
absolute `skills/` path as its output directory. Start new skills with explicit
`status: draft`; preserve an existing skill's status unless a change is requested.

Keep the entrypoint short and load detailed references only when needed. Add
scripts, assets, or agent settings only for a concrete requirement. These
instructions extend skill-creator for this repository; do not copy its general
authoring manual into each skill.

## Verify and report

Follow the authoring reference's validation checks and exercise any changed
scripts. Check representative requests against the description, including a
nearby request that should not trigger the skill. For a workflow, check that
each handoff supplies what the next skill needs and that a missing dependency
or failed step stops dependent work.

Report changed files, checks performed, and any limitations. Propose a commit
message. Leave publication and commits to the user's explicit request.
