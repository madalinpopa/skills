# Choose or write an issue

For a named issue, read its body, acceptance criteria, and relevant discussion
directly. Reuse an existing feature/phase plan if the project requires one;
avoid duplicating it in the issue.

## Choose work

Search open issues for the requested outcome. For "what should I work on",
start with a compact list of numbers, titles, and labels. Inspect only likely
candidates and check their blockers and current implementation. A limited
list is not an exhaustive search; expand it before declaring no match.

Recommend one issue with a short reason: useful now, unblocked, small enough
to finish, and appropriate for the user's experience. Show alternatives only
when a real tradeoff needs a decision. Reuse an explicitly selected issue
without asking the user to select it again. Do not close stale or apparently
completed issues as part of discovery.

## Scope and draft

Use the issue template linked from SKILL.md. Follow project templates and
planning gates when they take precedence. Define:

- One observable outcome, why it matters, and how to verify it.
- Scope and relevant non-goals; tasks say what changes and where, not the
  implementation. For bugs, include reproduction and expected/actual behavior.
- At most five planned commits per implementation issue. Count separate red
  and green commits when required. Each checkbox represents one reviewable
  commit outcome; mentoring steps may be smaller.
- A useful learning objective from this work. Do not increase scope merely
  to make the lesson harder.

Split larger work into a parent and focused sub-issues. Order dependencies
explicitly, prefer slices that deliver working behavior, and begin with an
unblocked sub-issue. An intentionally failing test commit may precede its
implementation under the project's policy.

If the requested outcome is ambiguous, ask for the missing behavior before
drafting tasks. If new scope appears during mentoring, propose a separate
issue instead of silently expanding the current one.

## Publish

Show the concrete draft and obtain approval if not already authorized. Use
`use-gh`, a body file, and explicit repository context. Check installed CLI
support before using sub-issue flags; use supported GitHub APIs when needed.
Honor project label rules and report created URLs. If a multi-issue operation
partially succeeds, report the created issues and remaining links; inspect
existing state before retrying to avoid duplicates.
