# Commit conventions

## Message format

- Follow [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/):
  `<type>(<scope>): <imperative summary>`. Omit scope when it adds no useful
  information; otherwise reuse an established scope for the affected area.
- Choose the type by the outcome: `feat` adds behavior, `fix` corrects behavior,
  `refactor` preserves behavior, and `test` changes only tests or test support.
  Use `docs` for documentation, `build` for build tooling and dependencies,
  and `ci` for CI workflows. Reserve `chore` for maintenance with no clearer type.
- Describe the concrete change, such as `fix(config): reject empty store paths`,
  rather than a vague summary such as `fix(config): update validation`.
- Mark breaking changes with `!` before the colon or a `BREAKING CHANGE:` footer.
  Explain the compatibility impact and any required migration.
- Add a body when the reason, tradeoff, or migration needs explanation. Keep
  command logs and review transcripts in validation records when present.
- Never add AI attribution; preserve genuine human attribution.

## Commit scope

- Keep each commit focused on one coherent, reviewable outcome. Include its
  necessary tests and documentation, respecting project policy when tests and
  implementation require separate commits.
- Review the exact diff and exclude unrelated work. Describe the final changes,
  including any scope adjustments, in the proposed message.
- Follow [TDD.md](TDD.md#workflow) for verification and review before proposing
  the commit. Use the repository's chosen version-control tool.

## Approval

- Create a commit only when the user approves that commit; approving a plan or
  other work does not count. Approval applies to the reviewed diff; later edits
  require renewed review and approval before committing.
- Never push, open or merge pull requests, or publish without an explicit user
  request.
