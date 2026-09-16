# Set up project instructions

Inspect root `AGENTS.md` and `CLAUDE.md` once when orienting. If either is
missing, use `create-agents-setup`. It owns instruction templates, conventions,
existing-file handling, and verification. Resolve a missing setup skill with
`use-skills-cli`, then invoke it or read its installed `SKILL.md` and follow it.
Do not keep a second setup template or copy files from another checkout.

If the dependency cannot be resolved, report the exact blocker. Continue direct
answers and independent read-only mentoring; pause changes requiring setup.
Setup does not authorize implementation or publication. After success, read
the resulting instructions and resume the same issue and checkpoint.

If both files exist, preserve their rules. Empty files or an absent CLAUDE
import need targeted repair through `create-agents-setup` within authorized
setup scope. Do not reinitialize files during ordinary questions or reviews.
If instructions do not select mentoring, follow the user's current mentoring
request; propose a compact persistent trigger only when adopting the workflow.
Surface conflicting instructions instead of silently replacing them.
