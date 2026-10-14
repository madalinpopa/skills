# Product features

Index the project's actual features. Requirements must fit
[MISSION.md](MISSION.md); delivery order and status belong in
[ROADMAP.md](ROADMAP.md). Do not invent features to populate the index.

## Feature index

Use sectioned bullet lists for the index, with a short outcome and links under
each feature. Keep long behavior descriptions in requirements.md.

Pending: list features agreed with the user. Give each the next unused
four-digit number, such as `0001`, a name, and a short intended outcome. The
number is the feature ID everywhere; it becomes the folder name
`NNNN-<slug>` when planning begins. Do not reuse or renumber IDs. Add links to
the feature folder when it exists:

- `requirements.md`: detailed behavior, scope, and acceptance criteria.
- `plan.md`: delivery steps and current implementation checkpoint.
- `validation.md`: verification evidence and user acceptance.

Use the naming rules and three templates in the
[feature guide](../features/README.md). A feature can stay in this index without
an empty folder until planning begins.

## Contract ownership

Detailed behavior lives in each feature's requirements.md: input rules, visible
results, state that must stay true, meaningful failures, and, where relevant,
what happens on repeated, conflicting, or partial runs. Acceptance criteria need concrete
scenarios and IDs such as `R-001`; cite them elsewhere as `NNNN/R-001`.
validation.md links those criteria to evidence.

Keep shared goals and constraints in MISSION.md and cross-feature choices in
[decision records](../decisions/README.md). Link these owners rather than
maintaining duplicate requirements here. Do not track progress here.
Changing an agreed contract requires an explicit user decision; update related
plans and validation expectations.
