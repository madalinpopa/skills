# Commit conventions

Use `<type>(<scope>): <imperative summary>`; omit scope when it adds no useful
information. Types include `feat`, `fix`, `test`, `docs`, `refactor`, `build`,
`ci`, and `chore`. Describe one reviewable outcome.

Propose one commit message per slice, after its adversarial review is clean,
following [TDD.md](TDD.md#workflow). The commit holds the slice's tests and
implementation together.

Review the exact diff and exclude unrelated work. Add a body when the reason,
tradeoff, or compatibility impact needs explanation. Never add AI attribution;
preserve genuine human attribution.

Create a commit only when the user approves that commit; approving a plan or
other work does not count. Never push, open or merge pull requests, or publish
without an explicit user request. Use the repository's chosen version-control tool.
