# Repository skill authoring

Use the [Agent Skills specification](https://agentskills.io/specification) as
the base format. Read `docs/SPEC.md` in the target `madalinpopa/skills` checkout,
especially "Skill metadata" and "Agent differences", for the CLI contract.
For optional agent settings, consult the current
[OpenAI skill documentation](https://learn.chatgpt.com/docs/build-skills#optional-metadata)
and [Anthropic Claude Code reference](https://code.claude.com/docs/en/skills#frontmatter-reference).

## Contents

- [Folder and scope](#folder-and-scope)
- [Frontmatter](#frontmatter)
- [Agent-specific options](#agent-specific-options)
- [Validate before review](#validate-before-review)

## Folder and scope

One directory per skill, with one shared `SKILL.md` and only the resources it
needs:

```
skills/
  review-code/
    SKILL.md                shared by every agent
    agents/openai.yaml      optional, Codex settings
    references/             optional, guidance read on demand
    scripts/                optional, executable helpers
    assets/                 optional, templates or output resources
```

All skill resources belong inside `skills/<skill-name>/`. Agent installation
directories such as `.agents/skills/`, `.claude/skills/`, and home directories
are destinations for the CLI, not source locations for this library. Keep
resources self-contained, use relative links, and use regular files and
directories; source symlinks and a root `.skill-lock.json` are unsupported.

Keep each skill focused on one reusable job with clear inputs and an expected
result. Inspect existing skills before adding overlapping instructions. Keep
`SKILL.md` short and below 500 lines; move substantial conditional guidance to
linked references and say when to read them. Add scripts only for concrete
automation needs. Avoid empty scaffolding and generic advice the agent already
knows. Shared instructions must work on both agents; the CLI preserves body
bytes, so it cannot translate Claude-specific substitutions or commands.

## Frontmatter

Begin `SKILL.md` with YAML frontmatter, followed by the instructions:

```markdown
---
name: review-code
description: Reviews code changes for correctness and maintainability. Use when
  asked to review a diff or pull request and report actionable findings.
status: draft
tags: [review, code]
---

# Review code

Read the changed code and its callers. Report actionable findings with file
locations and observable consequences.
```

| Field | Authoring rule |
| --- | --- |
| `name` | Required; matches the folder and chosen category prefix. |
| `description` | Required, nonempty, at most 1,024 characters; say what the skill does and when it applies. |
| `status` | Store field; explicitly use `draft` for new work and `published` when ready. Omission means `published`; other values are invalid. |
| `tags` | Optional store field; use a short list of strings for useful subjects, for example `[go, review]`. |
| `license`, `compatibility`, `metadata`, `allowed-tools` | Optional standard fields; include only when needed and follow the base specification's types and limits. Author and version belong in `metadata`. |
| `x-claude` | Optional store mapping for Claude-specific frontmatter; see below. |

`status` and `tags` belong at the top level, not under `metadata`. They are CLI
extensions, not part of the base standard, and are stripped on install. A draft
is hidden from store listings and cannot be installed. Publishing content needs
no CLI release: after the published skill reaches the configured store branch,
`skills sync` makes it available. The CLI reads committed store content, so
`skills ls` does not validate uncommitted authoring edits.

## Agent-specific options

Keep standard fields at the top level and isolate optional vendor settings:

| Agent | Source location | What the CLI writes |
| --- | --- | --- |
| Claude Code | `x-claude` in `SKILL.md` | Lifts its keys to top-level frontmatter for Claude; drops the mapping for other agents. |
| Codex | `agents/openai.yaml` | Copies the file to shared agent targets; omits it for Claude. There is no `x-codex` transform. |

For example, a Claude autocomplete hint belongs under `x-claude`:

```yaml
x-claude:
  argument-hint: "[file-or-diff]"
```

Other optional Claude settings include `model`, `context`, `agent`, `hooks`,
and invocation controls. Use only settings the workflow needs, with values
checked against the current Claude reference. `x-claude` cannot contain `name`,
`description`, `status`, `tags`, or another `x-claude`, and cannot duplicate a
top-level key. All YAML mappings must have unique keys. A top-level
`allowed-tools` uses the standard space-separated string format; Claude-only
tool grants can go under `x-claude` instead, without duplicating that key.

For Codex, optional `agents/openai.yaml` settings cover `interface` display
metadata, `policy` invocation controls, and `dependencies.tools` MCP
requirements. Follow the OpenAI documentation and skill-creator's sidecar
guidance; add only the settings needed for the skill.

Keep automatic discovery enabled unless explicit-only invocation is requested.
For that mode, use `x-claude.disable-model-invocation: true` for Claude and
`policy.allow_implicit_invocation: false` in `agents/openai.yaml` for Codex.
These settings control invocation, not authorization for the skill's actions.

Unknown top-level fields pass through unchanged; this does not mean an agent
supports them. See the "Field kinds" section of the target repository's
`docs/SPEC.md` for the full contract.

## Validate before review

Use skill-creator's validation workflow, then check this repository's additional
metadata and transformation rules. A validator for the base format may reject
the store-only keys: validate the shared transformed copy in a temporary
directory without removing those fields from the source. Check Claude extras
against the Claude reference separately. If a validator lacks a dependency or
rejects a currently supported standard field, report the limitation and check
that field against the current specification; do not weaken source metadata
just to satisfy an outdated validator. The CLI is not a full Agent Skills
validator, and there is no `skills check` command.

Check the folder/name match, YAML types, relative links, and any agent options.
Exercise added scripts and try representative prompts, including a nearby task
that should not select the skill. Review the rendered metadata for both Claude
and Codex when adding agent settings. Report checks and limitations, leave new
skills as drafts until ready for publication, and propose a commit message for
review. Commit or publish only when explicitly requested.
