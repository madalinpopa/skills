# Repository skill authoring

Read the sections relevant to the change. The current checkout's
`docs/SPEC.md` is the CLI contract; consult its "Skill metadata" section for
metadata and transformations, and "Agent differences" for agent integration.
Use the [Agent Skills specification](https://agentskills.io/specification)
for the base format. Routine body edits do not require these full references.

## Contents

- [Metadata](#metadata)
- [Packaging and publication](#packaging-and-publication)
- [Agent-specific options](#agent-specific-options)
- [Validation details](#validation-details)

## Metadata

| Field | Repository rule |
| --- | --- |
| `name` | Required, nonempty, matches the immediate skill folder and naming conventions. |
| `description` | Required, nonempty, at most 1,024 characters; identifies capability and trigger. |
| `status` | New skills explicitly use `draft`; preserve existing status unless requested. Omission means published; only `draft` and `published` are valid values. |
| `tags` | Optional list of strings; null and non-string entries are invalid. |
| `x-claude` | Optional mapping for Claude-only frontmatter; see below. |
| Standard optional fields | Check `license`, `compatibility`, `metadata`, and `allowed-tools` against the base specification when used. Author/version belong under `metadata`; top-level `allowed-tools` is a space-separated string. |

Every YAML mapping, including nested mappings, must have unique keys. Unknown
fields pass through unchanged; this does not establish agent support. Preserve
supported fields rather than restricting source to a validator's allowlist.

## Packaging and publication

The CLI discovers immediate directories under `skills/`, each containing
`SKILL.md` and its resources. Do not nest category directories. Keep the body
below 500 lines, references directly discoverable from it, and a contents
section in references longer than 100 lines. Avoid agent-specific substitutions
in shared body text: installation transforms metadata, not instructions.

`status` and `tags` are top-level store extensions stripped on install. Drafts
are hidden from store listings and cannot be installed. Publication needs
committed, published content on the configured store branch, followed by
`skills sync` and install/update in consuming projects; no CLI release is
needed. `skills ls` does not validate uncommitted authoring edits.

## Agent-specific options

Consult current docs only for settings being introduced or changed:
[OpenAI optional metadata](https://learn.chatgpt.com/docs/build-skills#optional-metadata)
and [Claude frontmatter](https://code.claude.com/docs/en/skills#frontmatter-reference).

| Agent | Source | Install behavior |
| --- | --- | --- |
| Claude Code | `x-claude` in `SKILL.md` | Lift its fields to top level for Claude; drop the mapping for other agents. |
| Codex | `agents/openai.yaml` | Copy to shared targets; omit for Claude. No `x-codex` transform. |

`x-claude` requires unique string keys. It cannot contain `name`, `description`,
`status`, `tags`, or `x-claude`, or collide with top-level fields. Standard
fields stay at top level. Put Claude-only settings such as `argument-hint`,
`model`, `context`, or hooks under `x-claude`; follow vendor-supported types.
Claude-only tool grants may use `x-claude.allowed-tools` when absent at top level.

Use Skill Creator's sidecar guidance for Codex `interface`, `policy`, and
`dependencies.tools`. Preserve automatic discovery unless explicit-only use
is requested: that mode uses `x-claude.disable-model-invocation: true` for
Claude and `policy.allow_implicit_invocation: false` in `agents/openai.yaml`
for Codex. Invocation settings grant no additional action permissions.

## Validation details

The base validator may reject store-only keys. Validate a temporary shared
representation with `status`, `tags`, and `x-claude` removed, while checking
source metadata separately against the CLI contract. Preserve YAML values,
expand aliases whose anchors are removed, and leave body bytes unchanged.

When changing agent settings, also inspect the Claude representation with
`x-claude` lifted and the shared representation with it dropped; verify the
sidecar's inclusion/exclusion. Check vendor fields against the affected
agent's documentation rather than a base validator's field allowlist.

Report missing validator dependencies or unsupported standard fields; do not
weaken source metadata to make an outdated validator pass. The CLI is not a
full Agent Skills validator and has no `skills check` command.
