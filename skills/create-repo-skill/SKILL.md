---
name: create-repo-skill
description: Creates or updates catalog skills and authoring guidance in this repository or its forks, with local scaffolding, validation, and context-efficient design. Not for skill creation elsewhere or Go CLI implementation.
status: published
tags: [skills, authoring, repository]
---

# Create a repository skill

Create useful, self-contained skills that supply decisions and procedures an
agent cannot reliably infer. Keep the normal execution path short and load
conditional detail only when needed.

## Prepare

Identify the current checkout by `AGENTS.md`, `docs/SPEC.md`, and `skills/`.
Follow its instructions and inspect the diff and relevant spec. Author here,
including forks, rather than in an installed skill or configured store. Do
not change `store.repo` or `store.branch` merely because this is a fork.

Resolve the requested outcome, inputs, and scope from the request and nearby
skills. Ask only for missing information that changes the result; preserve
explicit choices and existing authorization. A skill must not expand the task
or grant permission to commit, publish, deploy, or contact others.

## Create or update

1. For a new or substantially redesigned skill, read
   [design](references/design.md) and [optimization](references/optimization.md).
   For narrow edits, use the essentials below and read only affected topics.
2. Author in `skills/<skill-name>/`, with all supporting resources inside it.
   New skills default to explicit `status: published`; use `draft` only when
   requested. Preserve existing status unless asked to change it. Published
   metadata does not authorize a commit or remote publication.
3. For a new skill, use the local initializer when it saves repetitive work:

   ```sh
   "<skill-dir>/scripts/init-skill.sh" <name> --repo <checkout> --description "<capability and trigger>"
   ```

   `<skill-dir>` is this skill's installed directory. Do not initialize an
   existing skill again. Optional resources, metadata commands, and runtime
   requirements are in [scripts](references/scripts.md).
4. Write only guidance that changes decisions or improves execution. Keep
   shared purpose, required inputs/defaults, essential constraints, commands,
   verification, and concise reporting in the entrypoint. Put substantial
   conditional guidance in directly linked references with explicit triggers.
   Copy templates as assets; do not routinely read or regenerate them.
5. Prefer Bash helpers for repeated file creation, copying, transformations,
   and multi-command verification. Run existing helpers rather than rewriting
   their logic. Use a proper parser for structured formats; do not parse YAML
   with shell substitutions. Add automation only for concrete repeated work.

## Repository essentials

Keep descriptions concise and front-load capability and trigger; include only
useful exclusions. Separate advice from actions and required inputs from
optional defaults. Do not require every reference, a full file tree, or full
successful command logs on ordinary runs. Preserve meaningful checks, failure
reporting, partial-change visibility, and authorization boundaries.

Use relative links and regular files. Source symlinks and a root
`.skill-lock.json` are unsupported. `.agents/` and `.claude/` are installation
targets, not source. The CLI preserves body bytes; shared instructions must
work across agents.

## Load details by task

| Task | Reference |
| --- | --- |
| Choose a name/category or coordinate skills | [Naming and workflows](references/naming.md) |
| Change metadata, packaging, or publication | [Repository authoring](references/authoring.md) |
| Add or edit Codex UI metadata or invocation policy | [Agent metadata](references/openai-yaml.md) |
| Assess behavior, scripts, or context savings | [Validation](references/validation.md) |

## Verify and report

Run `"<skill-dir>/scripts/validate-skill.sh" <skill-folder>`, exercise changed
scripts, and check realistic triggering and non-triggering requests. Format
validation does not prove good agent behavior. For complex changes, follow the
conditional evaluation procedure in the validation reference.

Report the result, changed files, checks and limitations, and proposed commit
message. Use concise counts when comparing context; distinguish words from
measured tokens. Commit or publish only when explicitly requested, honoring
existing authorization.
