# Agent instructions

## Project

- `skills` is a Go CLI and a catalog of reusable agent skills.
- [docs/SPEC.md](docs/SPEC.md) is the CLI behavior source of truth.
- Inspect relevant code, specs, and diffs before changes; preserve unrelated work.
- Check current official docs for every external library touched.

## Review gates

- Behavioral changes: write necessary tests, run them, and **stop for test
  review** with the diff, failure evidence, and a proposed test commit message.
- Implement only after explicit developer approval. Then verify and **stop for
  implementation review** with changes, results, and a proposed commit message.
- Config, docs, and release changes that gain nothing from failing tests may
  omit the test gate; implementation review still applies.
- Plan approval never bypasses test review or authorizes additional scope.

## Version control

- This repository uses jj-vcs. Invoke `use-jj` for version-control work.
- Never commit unless explicitly asked. Open PRs only after full developer
  review and an explicit request. Never add AI attribution; preserve human attribution.
- Propose a commit message for every requested change. Read
  [COMMITS.md](docs/conventions/COMMITS.md) for delivery rules and PR preparation.

## Read for the current task

- Before creating or editing Go files: read
  [GOLANG.md](docs/conventions/GOLANG.md) and invoke
  `modern-go-guidelines:use-modern-go` as directed there.
- Before choosing or writing tests, changing behavior, or verifying work: read
  [TDD.md](docs/conventions/TDD.md). Test each behavior at the lowest useful
  layer; avoid duplicate coverage. Default tests need neither network nor Docker.
- Before repository skill creation or updates: read
  [SKILLS.md](docs/conventions/SKILLS.md), run `skills sync`, prepare
  `create-repo-skill` once with the CLI, then use the installed skill.
- For planned phases/commits or feature planning: use
  [FEATURE.md](docs/templates/FEATURE.md) and [PHASE.md](docs/templates/PHASE.md).
  Reuse an existing feature or create `docs/features/<feature-name>/FEATURE.md`;
  detail only the active phase in `phase-NN-<outcome>.md`.
- Present plans for review before changes; honor existing approval within scope.
  Update phase checkpoints and the feature roadmap at review boundaries and
  session end. On resume, read the overview and active phase and verify recorded
  state against the working tree.
