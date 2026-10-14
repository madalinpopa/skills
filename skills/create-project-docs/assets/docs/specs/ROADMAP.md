# Roadmap

Record delivery order and status for features in [FEATURES.md](FEATURES.md),
guided by [MISSION.md](MISSION.md). This is the only place that tracks feature
status. Detailed execution belongs in numbered
[feature records](../features/README.md), not duplicate progress logs here.

## Delivery sequence

- [ ] Pending: agree on the first useful outcome and its dependencies.

For each agreed feature, list small delivery slices as checklist subtasks:

- [ ] `NNNN`: Pending: intended outcome; link to plan.md when it exists.
  - [ ] Slice 1: Pending: smallest useful, independently verifiable outcome.
    - Scope and acceptance: Pending: bounded behavior and criterion links.
    - Depends on: Pending: prerequisites or earlier slices; use None when known.
    - Plan and evidence: Pending: links to phase/tasks and validation.md.
  - [ ] Slice 2: Pending: next bounded outcome, only when needed.
    - Scope and acceptance: Pending: behavior and criterion links.
    - Depends on: Pending.
    - Plan and evidence: Pending.
  - [ ] User acceptance: Pending: link to the feature decision in validation.md.

- Record each entry's dependency and status as short bullets beneath it.
- Split broad features before scheduling implementation. Each slice must have
  one useful outcome that can be implemented, verified, and reviewed on its own;
  avoid slices defined only by technical layers with no observable outcome.
- Order slices by dependencies and early useful delivery. Keep their scope small
  enough for one focused review; split a slice again if it bundles unrelated work.
- Link each slice to the corresponding phase and tasks in plan.md. A phase may
  group related slices; its one-to-five-commit limit does not set slice size.
- Mark slices `[x]` after their verification and user review, with linked evidence;
  keep deferred or blocked slices `[ ]` with a reason. Detailed implementation
  tasks and commit proposals stay in plan.md.
- Use separate numbered sub-features when a slice needs its own requirements
  and acceptance record, following the feature guide; otherwise keep it under
  the existing feature ID.
- Mark a feature `[x]` only when its status is `accepted`; completed checks or
  implementation tasks alone do not mean the user accepted it.
- Replace placeholders only with real agreed work. Do not invent deadlines or
  tasks to fill the checklist.

Statuses:

- `proposed`: listed but not yet agreed for planning.
- `planned`: requirements and plan agreed; work not started.
- `active`: work in progress.
- `awaiting review`: checks done; waiting for the user's acceptance decision.
- `accepted`: the user accepted it, as recorded in validation.md.
- `deferred`: set aside by user decision.

Passing checks alone do not make a feature accepted.

## Later work

- [ ] Pending: outline a future outcome briefly; link detailed agreed phases in plan.md.

Keep excluded product scope in MISSION.md rather than treating it as an obligation.

## Review boundaries

Update slice checklists, feature status, and ordering when the user makes a
delivery decision or reviews an outcome. Link the active feature's plan.md for its checkpoint and next action,
and validation.md for evidence and acceptance.
