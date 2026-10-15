# Repository skill authoring

When the user asks to create or update a repository skill, prepare
`create-repo-skill` with this repository's `skills` CLI once at the start of
the task, then use the installed skill:

1. Run `skills sync` from the repository root to refresh the configured
   store. If the CLI is unavailable, report it and follow the README's CLI
   installation instructions before continuing.
2. Select the current agent: `codex` or `claude`. Check its configured project
   target for `create-repo-skill/.skill-lock.json` (default roots are
   `.agents/skills/` and `.claude/skills/`, respectively).
3. If the lock is absent, run
   `skills install create-repo-skill --agent <agent>`; otherwise run
   `skills update create-repo-skill --agent <agent>`, replacing `<agent>`
   with the selected name. Use project scope. Report command failures or local
   conflicts before proceeding; do not add `--force` or replace files manually.
4. After success, invoke `create-repo-skill` (`$create-repo-skill` in Codex,
   `/create-repo-skill` in Claude). If it is not yet discoverable in the current
   session, read and follow the installed `SKILL.md` at the configured target.
   It supplies self-contained authoring rules and local helpers. Do not
   repeat this setup when the skill refers back to this file in the same task.

The CLI installs published, committed store content. If the skill is absent
after sync, report that its `status: published` source must reach the configured
store branch first; do not commit or push automatically. Author changes under
`skills/`; `.agents/` and `.claude/` are ignored local installation targets.

## Verification

Use the installed `create-repo-skill` validation guidance and run relevant
skill tests. `task test:skills` validates every catalog skill and runs its
Python tests; use it when changing shared authoring or validation machinery.
Run commands from the repository root. Inspect the final diff and report
failures and skipped checks.
