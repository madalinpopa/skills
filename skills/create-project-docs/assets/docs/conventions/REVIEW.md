# Adversarial review

Default to one independent, read-only reviewer via the available subagent tool.
Provide the diff, approved scope, requirements, test evidence, and necessary code
context; omit implementer conclusions.

Reviewer instructions:

- Challenge correctness, regressions, boundary handling, security, and test adequacy.
  Inspect relevant callers, contracts, and existing tests beyond the diff; check
  whether tests fail for broken behavior.
- Verify findings with a reproducer, failing check, or clear code path. Report
  actionable issues introduced/exposed by this slice, each root cause once;
  distinguish questions from confirmed defects.
- Order findings by severity based on impact/reachability. Include location,
  trigger, impact, and evidence; distinguish blockers from optional improvements.
- Exclude style preferences; maintainability findings need concrete impact.
  Do not block on polish or equally valid designs.
- Return concise findings or “No actionable findings”, with material limitations,
  unrun checks, and relevant context not inspected.

Verify each finding; mark it fixed, rejected with evidence, or awaiting a user
decision. Do not repeat rejected findings without new evidence. Record findings,
dispositions, and checked state beside the affected task in tasks.md when present;
avoid repeated review logs in replies.

Repeat review after fixes until no blocking findings remain; later code edits
require renewed review. Optional improvements do not block completion. If review
repeats a disputed finding without new evidence or stalls on blockers, stop and
request a user decision with evidence and remaining disagreement. Never claim
clean review with unresolved blockers; clean review does not prove defect-free code.

Review authorizes neither new scope nor dependency changes. Record unresolved
questions in [the question log](../QUESTIONS.md), following the
[question guide](QUESTIONS.md), and link them from affected requirements/plans.
Ask with the available question tool
before dependent work. If delegation is unavailable, report the blocker and request
a decision; do not substitute self-review or claim clean review.
