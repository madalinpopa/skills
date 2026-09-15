# Reviewing

Read this when the user says a task is done or asks for a review. It covers
what to look at, how to phrase findings, and what "good" means here.

## Before writing

1. Read the whole diff for the task, not the files. Note anything outside
   the task's scope. Read the tests first: they should state the behavior
   the task asked for. If the user did not report a red run, ask how the
   test failed before the code existed.
2. Run the project's own checks, narrowest first: the package tests, then
   format, vet, lint, then the full suite when the task warrants it. Use
   the commands the repository documents, not assumed defaults.
3. For Go, run `use-modern-go` on each changed file and keep only findings
   that apply to this diff.
4. Decide the verdict first, then write. A review that wanders is not read.

## The three parts

**What is good.** Two or three specific things, each naming the code and
the principle it serves. Praise is information: it tells the user what to
keep doing. Skip generic praise.

**What to improve.** Ordered by severity, each with a label, the location,
the problem, and the why. Use the labels from Conventional Comments so the
weight is clear:

- `issue (blocking)`: wrong behavior, missing error handling, a security
  gap, a behavior change with no test, a test that does not test the
  change. The task is not done.
- `suggestion`: clearer or simpler code that the user should apply now.
- `nitpick (non-blocking)`: style or naming; note it, do not block.
- `question`: something you do not understand yet. Ask, do not assume.
- `thought (non-blocking)`: a lesson for later, not this task.

Comment on the code, never the person. Base every point on a principle or a
project rule, not taste. State the problem and let the user propose the fix;
give the fix only when the hint ladder in `mentoring.md` is exhausted or
they ask. One blocking issue at a time is enough to send back.

**Verdict.** "Task complete" or one correction step in Goal, Do, Why, Done
when form. Then the ticked task and the next step, or the PR handoff.

## What good means here

In this order, from Beck's rules of simple design:

1. It passes the tests, and the tests exercise the change.
2. It reveals intent: names say what, control flow is direct, a reader
   needs no comment to follow it.
3. It has no duplication the reader would trip over. A little copying beats
   a new dependency or a premature abstraction.
4. It has the fewest elements: no unused parameter, flag, interface, or
   configuration.

Correctness before performance. Do not ask for optimization without a
measurement that shows a problem in this task. Ask instead about error
paths, boundaries, empty and invalid input, and concurrent access.

KISS and YAGNI apply to added complexity, not to care. Challenge every
abstraction with "which current requirement needs this?" and name the
carrying cost when the answer is none. Never accept sloppy code because it
is small: it still needs tests, handled errors, and idiomatic shape.

Security is part of correctness. Look for unvalidated external input,
secrets in code or logs, paths built from user input, shell strings built
from arguments, and permissions wider than the task needs.

## Scope discipline

Review the task, not the repository. Unrelated code gets a comment only for
a correctness or security bug, and even then as `thought` or a proposed new
issue, never as a blocker for this task. Refactors the user did not ask for
and were not in the task go back out or into their own issue.

## Commits

Each task is one commit with an imperative subject in the project's style
and a body that says why when the subject is not enough. If the commit is
hard to summarize in one line, it is two commits. Ask the user to split it
before ticking the task.

## Sources

- Google, the standard of code review: <https://google.github.io/eng-practices/review/reviewer/standard.html>
- Google, how to write review comments: <https://google.github.io/eng-practices/review/reviewer/comments.html>
- Conventional Comments: <https://conventionalcomments.org/>
- Fowler, Beck's design rules: <https://martinfowler.com/bliki/BeckDesignRules.html>
- Fowler, YAGNI: <https://martinfowler.com/bliki/Yagni.html>
- Go proverbs: <https://go-proverbs.github.io/>
