# Task skill setup fallback

Use this only when `use-skills-cli` is unavailable to load. Run from the target
project root, never from this skill's directory.

1. Run `skills version`. If unavailable, report the prerequisite and stop:
   `brew install --cask madalinpopa/tap/skills` on macOS, or
   `go install github.com/madalinpopa/skills@latest`. Do not install the CLI silently.
2. Select `codex` for Codex/Gemini or `claude` for Claude Code. Inspect configured
   project targets and `skills ls --local`; defaults are `.agents/skills/` and
   `.claude/skills/`. A skill's `.skill-lock.json` identifies its installation.
3. Run `skills sync`, then `skills ls`. Confirm `use-task-spec`, `use-task-plan`,
   and `use-tasks-checklists` are published and available. Sync never updates
   installed files. Stop on dirty/divergent store errors; never reset the store.
4. For each required skill, run `skills install <name> --agent <agent>` when its
   lock is absent, otherwise `skills update <name> --agent <agent>` to compare
   and apply the store version. Use project scope; never add `--force` or `--global`.
5. Exit 0 succeeds; exit 1 is a runtime/availability error; exit 2 is invalid
   input (correct and retry once); exit 3 needs attention for local edits,
   conflicts, or unavailable content. Report errors and the printed diff hint;
   stop dependent setup instead of manually replacing installed files.

If a successfully installed skill is not yet discoverable, read its installed
SKILL.md and follow it directly, loading references only as needed. Report the
agent, destination, and whether the skill was loaded this way. Do not change
ignore rules or commit installed files as part of setup. Unpublished source
must reach the configured store branch before these installations can succeed;
do not commit or push automatically to make it available.
