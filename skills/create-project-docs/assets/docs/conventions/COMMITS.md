# Commit conventions

Use `<type>(<scope>): <imperative summary>`; omit scope when it adds no useful
information. Types include `feat`, `fix`, `test`, `docs`, `refactor`, `build`,
`ci`, and `chore`. Describe one reviewable outcome.

Propose a message after each verified change. For behavior changes, keep the
failing tests and implementation separately reviewable for two commits, following
[TDD.md](TDD.md). Documentation or configuration work that gains nothing from
a failing test may use one reviewed commit; classify by effect, not extension.

Review the exact diff and exclude unrelated work. Add a body when the reason,
tradeoff, or compatibility impact needs explanation. Never add AI attribution;
preserve genuine human attribution.

Never create commits, push, open or merge pull requests, or publish without an
explicit user request. Approval of a plan or implementation does not authorize
these actions. Use the repository's chosen version-control tool.
