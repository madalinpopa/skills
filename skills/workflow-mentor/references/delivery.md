# Prepare and open the PR

After every issue criterion is verified:

1. Review the full intended change against the confirmed base branch,
   including commits made during mentoring. Check scope and required checks.
2. Resolve the head branch or jj bookmark, repository, and remote. Verify
   the remote head matches the reviewed local commit using current remote
   evidence; local `git status` or `jj log` alone does not prove publication.
   Ask the user to push if needed, unless pushing is already authorized.
3. Check for an existing PR for that head. Reuse it rather than creating a
   duplicate. Prepare the title and body using the PR template linked from
   SKILL.md and the project's conventions. Describe the final behavior and
   actual verification. Include `Closes #N` only when the full issue is done.
4. Present the concrete draft for any required project review. Open or update
   the PR only with explicit user authorization, honoring prior approval.
   Use `use-gh`, an explicit base/head, and a body file. Do not let PR creation
   implicitly push or fork an unverified branch.
5. Report the PR URL, completed verification, and outstanding checks. If the
   command fails ambiguously, check remote state before retrying.

Keep the issue open for the PR's merge. Do not merge, close it manually, or
change labels without explicit authorization. Never add AI attribution to
commits or PRs.
