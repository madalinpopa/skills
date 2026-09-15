# Set up project instructions

Use this only on first adoption or when project instructions are missing.
Skip it during ordinary questions and reviews. Read existing instructions
before proposing changes; they take precedence over this skill's defaults.

Inspect `AGENTS.md` and `CLAUDE.md` at the repository root. Stop setup on a
symlink or non-regular target rather than copying through it. If either file
is absent, use the bundled helper:

```sh
"<skill-dir>/scripts/init-agent-docs.sh" <repository-root>
```

The helper copies missing templates and keeps existing files. Fill created
AGENTS.md placeholders from the repository: purpose, applicable spec, and
build/test/format/lint commands. Remove inapplicable lines; ask only for
essential information the repository cannot supply. Show the completed files.

If AGENTS.md exists without this workflow, propose a compact mentoring trigger
and any needed project-specific limits. Add it only within the user's
authorization; do not append the whole template or replace project rules.
If CLAUDE.md does not import `@AGENTS.md`, propose that line without rewriting
the file. Skip once both files are present and the workflow is configured.

Report created, kept, or edited paths and their uncommitted state. On failure,
inspect both paths and report partial creation; the helper is not transactional.
