# Roadmap

Record delivery order and status for features in [FEATURES.md](FEATURES.md),
guided by [MISSION.md](MISSION.md). This is the only place that tracks feature
status. Detailed execution belongs in numbered
[feature records](../features/README.md), not duplicate progress logs here.

## Delivery sequence

- [ ] Pending: agree on the first useful outcome and its dependencies.

For each agreed feature, use a checklist entry with short descriptive subtasks:

- [ ] `NNNN`: Pending: intended outcome; link to plan.md when it exists.
  - [ ] Pending: phase outcome and a link to its plan tasks or evidence.
  - [ ] User acceptance: Pending: link to the decision in validation.md.

- Record each entry's dependency and status as short bullets beneath it.
- Mark delivered milestones `[x]` with evidence; keep deferred or blocked ones
  `[ ]` with a reason. Detailed implementation tasks stay in plan.md.
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

Update status and ordering when the user makes a delivery decision or reviews
an outcome. Link the active feature's plan.md for its checkpoint and next action,
and validation.md for evidence and acceptance.
