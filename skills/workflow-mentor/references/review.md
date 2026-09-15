# Review a mentoring step

## Establish evidence

1. Resolve the task and its review boundary. Inspect all relevant changes:
   untracked files, unstaged/staged changes, and commits since the agreed
   base. A clean working tree does not mean there is no work to review.
   Use the repository's version-control skill; preserve unrelated edits.
2. Read tests and the changed behavior, following callers and boundary code
   as needed. Check the task's observable acceptance criteria, not just the
   patch's apparent intent. Apply Go guidance to the relevant changed files.
3. Use recorded red-run evidence. If missing, establish that the test detects
   the regression in an isolated temporary copy when practical, or ask for
   the missing evidence. Never revert the user's work to manufacture red.
4. Run project checks in the documented order, narrowest first where allowed.
   Distinguish failures caused by the change from environment failures.
   Do not run auto-fix commands or modify implementation files during review
   without authorization; guide the user through required mutating checks.
5. Decide what this boundary proves. A test-only task can pass review with
   the expected red result; the behavior and issue remain incomplete until
   implementation and final checks pass. Do not infer success from unavailable
   checks or skip required project gates.

Prioritize wrong behavior, regression coverage, error handling, and applicable
security boundaries. Request clarity improvements with a concrete reason.
Avoid speculative abstractions, performance work without evidence, or tests
that duplicate existing coverage.

## Respond

- **Verdict:** verified for this boundary, needs correction, or blocked by
  missing evidence. Summarize checks and any limits on the conclusion.
- **Feedback:** identify useful choices worth retaining, then actionable
  findings with location, consequence, and blocking/non-blocking status.
  No praise quota, invented findings, or taste-based blockers.
- **Next action:** one correction or next step in Goal / Do / Why / Done when
  form. Report all material blockers found, but guide one correction at a time.

Keep unrelated suggestions outside the task. Surface a serious correctness
or security discovery separately; do not silently edit it or expand scope.
Propose a commit message at a commit boundary in the project's style.

## Record completion

Only tick a checkbox when its full outcome passes review, not after every
conversational step. For an authorized issue update, run from the repository:

```sh
"<skill-dir>/scripts/issue-task.sh" <N> --check "<unique task fragment>" -R <owner/repo>
```

Exit `0` means updated or already checked; `2` is usage; `3` means zero or
multiple matches. Inspect the current body and refine the fragment. Other
failures may occur after a remote write: re-read before retrying and report
whether the update persisted. The helper replaces the body; avoid concurrent
issue edits and verify unrelated content was preserved.

When all criteria are verified, summarize what was learned and what passed,
then use the delivery reference linked from SKILL.md. Leave the issue open.
