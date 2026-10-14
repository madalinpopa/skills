# Roadmap

Record delivery order and status for features in [FEATURES.md](FEATURES.md),
guided by [MISSION.md](MISSION.md). This is the only place that tracks feature
status. Detailed execution belongs in numbered
[feature records](../features/README.md), not duplicate progress logs here.

## Delivery sequence

Pending: agree on the first useful outcome and its dependencies. For each real
entry, record the feature number (`NNNN`), outcome, dependency, status, and link
to the feature's plan.md when it exists. Do not invent deadlines or tasks to
populate a table.

Statuses:

- `proposed`: listed but not yet agreed for planning.
- `planned`: requirements and plan agreed; work not started.
- `active`: work in progress.
- `awaiting review`: checks done; waiting for the user's acceptance decision.
- `accepted`: the user accepted it, as recorded in validation.md.
- `deferred`: set aside by user decision.

Passing checks alone do not make a feature accepted.

## Later work

Pending: outline future outcomes briefly; detail only the active work. Keep
excluded product scope in MISSION.md rather than treating it as an obligation.

## Review boundaries

Update status and ordering when the user makes a delivery decision or reviews
an outcome. Link the active feature's plan.md for its checkpoint and next action,
and validation.md for evidence and acceptance.
