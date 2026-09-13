# Authoring reference

Use the [Agent Skills specification](https://agentskills.io/specification) as
the format reference and its [authoring guide](https://agentskills.io/skill-creation/best-practices)
when deciding what instructions belong in a skill.

## Metadata

Write YAML frontmatter followed by Markdown instructions.

| Field | Constraint |
| --- | --- |
| `name` | Required; 1–64 lowercase alphanumeric characters or hyphens; match the directory; no leading, trailing, or consecutive hyphens. |
| `description` | Required; 1–1024 characters explaining the task and when it applies. |
| `license` | Optional license identifier or bundled license reference. |
| `compatibility` | Optional; 1–500 characters describing actual environment requirements. |
| `metadata` | Optional mapping of string keys to string values; quote numeric-looking values. |
| `allowed-tools` | Optional experimental space-separated tool string; support depends on the client. |

This demo uses `allowed-tools: Read` to demonstrate the experimental field.

## Resources and instructions

- `scripts/` holds executable helpers. Document dependencies, usage, and failures.
- `references/` holds guidance read when a task needs it.
- `assets/` holds templates and other output resources.

These folders are optional. This demo includes all three to exercise copying.
The standard also permits other directories, but this example intentionally
uses only those three conventions.

Keep the entrypoint below 500 lines and preferably 5,000 tokens. Link resources
directly from `SKILL.md` using paths relative to the skill root. Give the agent
a concrete task, expected result, and checks; load supporting detail as needed.

## Using this repository

The source extensions are defined by this repository's `docs/SPEC.md`:

- `status: draft` keeps a skill unpublished. Omitting `status` or setting
  `published` makes it installable, so mark work in progress explicitly.
- `tags` groups skills in the catalog. Omitting it means no tags.
- `x-claude` holds Claude-only fields. Its keys must not duplicate shared fields
  or redefine `name`, `description`, `status`, `tags`, or `x-claude`.

The CLI removes `status` and `tags` on install. It lifts `x-claude` into Claude's
frontmatter and drops it for every other agent. These are store extensions, not
standard fields. See the [Claude skills reference](https://code.claude.com/docs/en/skills)
for supported extension values.

Claude fields such as `model`, `hooks`, or `paths` placed at the top level are
copied to every agent. Put them under `x-claude` when other agents should not
see them. Unknown fields pass through unchanged. The CLI keeps these values but
does not check that an agent supports them. Every mapping, including nested
ones such as `metadata` and `hooks`, must have unique keys. A Claude skill
without `name` and `description` needs both added before it goes into a store.

The CLI copies regular files and preserves executable intent. Use real files;
symlinks are unsupported. The root `.skill-lock.json` belongs to the installer.
Keep generated reports outside the installed skill so they do not become local
installation content.

After installation, run `skills-ref validate` against the shared
`.agents/skills` copy if the official reference validator is available. The
source and Claude copy contain extensions that this strict validator rejects.
Check Claude's extra fields separately. Metadata validation does not replace
trying the skill on a representative request.
