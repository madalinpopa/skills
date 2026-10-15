# Commits and pull requests

Each behavioral change has two commits: failing tests, then implementation.
Config, docs, and release changes that gain nothing from failing tests may use
one commit; developer review before committing still applies.

Propose a commit message for every requested change; never create Git or
Jujutsu commits unless asked. Keep PRs focused on one piece of work; open them
only after full developer review and an explicit request.
On “ready to create PR,” review the latest commits between `main` and the
working-copy revision and provide a PR title. Use `use-jj` for this repository's
version-control operations.

Never add AI attribution (including Co-Authored-By trailers or generated-by
signatures) to commits or PRs; preserve genuine human attribution.
