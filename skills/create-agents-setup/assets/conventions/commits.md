# Commit conventions

Follow established repository conventions. Where none exist, use
`<type>(<scope>): <imperative summary>`; omit scope when it adds no information.
Choose `feat`, `fix`, `test`, `docs`, `refactor`, `build`, `ci`, or `chore`
according to the change. Describe one outcome, not a list of files.

Examples: `test(config): cover missing settings`,
`fix(config): preserve defaults when settings are absent`,
`docs(setup): explain local verification`.

Propose a message after each completed, verified change. A conversational step
with no file change needs no empty commit. Follow the project's test-review
workflow for commit boundaries.

Review the exact diff before committing; exclude unrelated work. Add a short
body only when the reason, tradeoff, or compatibility impact needs explanation.
Never add AI attribution; preserve genuine human attribution.

Never commit, push, open or merge PRs, close issues, or change labels without
explicit authorization. PRs also require completed review. Approval of a plan
or implementation does not authorize publication. Use the applicable
version-control and GitHub skills selected by AGENTS.md.
