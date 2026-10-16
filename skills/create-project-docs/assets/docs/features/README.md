# Feature and codebase work records

Use records for real features, fixes, refactoring, migrations, tooling, and
planned documentation work. For disposable experiments, follow
[PROTOTYPE.md](../conventions/PROTOTYPE.md).

## Ownership and routes

- Use `NNNN-<slug>/spec.md`, `plan.md`, and `tasks.md`. Choose the next unused
  four-digit ID across feature folders and [ROADMAP.md](../specs/ROADMAP.md);
  preserve existing IDs and locations. Do not add competing phase/validation records.
- For requirements authoring, use `use-task-spec`; `spec.md` owns outcomes,
  constraints, and acceptance criteria.
- For technical planning, use `use-task-plan`; `plan.md` owns approach and sequencing.
- For checklist creation, execution, or resumption, use `use-tasks-checklists`;
  `tasks.md` owns progress, evidence, and the next action. On resume, start with
  the active phase and nearby task evidence, then relevant spec/plan sections.
- Completion does not imply acceptance, commit permission, or additional scope.
- Keep actual unresolved questions in [the question log](../QUESTIONS.md),
  following [its guide](../conventions/QUESTIONS.md).

Skills are installed through `use-skills-cli`; report missing prerequisites or
local conflicts instead of replacing installed content manually. Keep project
facts here and reusable workflow instructions in their owning skills.
