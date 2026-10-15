# Commit conventions

## Messages

- Use [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/):
  `<type>(<scope>): <imperative summary>`. Omit unhelpful scope; otherwise reuse
  an established scope. Describe the concrete outcome, e.g.
  `fix(config): reject empty store paths`.
- Types: `feat` adds behavior; `fix` corrects it; `refactor` preserves it; `test`
  changes tests/support only; `docs` documents; `build` changes build tooling or
  dependencies; `ci` changes CI. Use `chore` only when no clearer type fits.
- Mark breaking changes with `!` or a `BREAKING CHANGE:` footer; explain impact
  and migration. Add a body for useful rationale/tradeoffs. Keep command logs and
  review transcripts in validation records when present.
- Never add AI attribution; preserve genuine human attribution.

## Scope and approval

- Keep each commit to one reviewable outcome with necessary tests/docs; respect
  project rules requiring separate test/implementation commits.
- Review the exact diff, exclude unrelated work, and propose a message matching
  final scope. Follow [TDD.md](TDD.md#workflow) verification/review and use the
  repository's chosen version-control tool.
- Commit only with approval for that reviewed diff. Plan/work approval is not
  commit approval; later edits require renewed review and approval.
- Push, open/merge pull requests, or publish only on explicit user request.
