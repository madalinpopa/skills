# Feature records

Create one folder per real feature, named `NNNN-<slug>`. Assign the next unused
four-digit number, checking [ROADMAP.md](../specs/ROADMAP.md) and existing feature
folders. `NNNN` is the feature's only ID everywhere. Preserve existing numbering
and names; never reuse or renumber IDs. A feature may appear in ROADMAP.md before
it needs a folder.

```text
0001-feature-x/
  plan.md
  requirements.md
  validation.md
```

Use `validation.md` with one dot before the extension. Create these records
only when a real feature is being planned, not as empty examples during setup.

## Document ownership

- `requirements.md`: feature scope, behavior, constraints, and acceptance criteria.
- `plan.md`: implementation approach, delivery steps, approvals, and current checkpoint.
- `validation.md`: verification evidence and the user's acceptance decision.

[ROADMAP.md](../specs/ROADMAP.md) groups features and sub-features with
task/subtask checklists, without a separate feature index. Link
tasks to plan.md for execution details and validation.md for evidence when those
records exist. Keep acceptance decisions in validation.md, dependencies and the
checkpoint in plan.md, and behavior in requirements.md. Do not add status fields,
acceptance prose, or workflow notes to the roadmap. Link shared facts rather
than duplicating them. Use roadmap order for delivery order.

## Writing structure

- Add sections only when there is real content, following the
  [spec writing guide](../specs/README.md). Do not seed empty outlines.
- Group related information logically: scope, behavior, and acceptance;
  approach, delivery, and checkpoint.
- Write short bullets under sections. Split substantial behavior contracts into
  subsections for inputs, outcomes, state, and failures when needed.
- Use tables only for compact comparisons or mappings; keep requirements,
  rationale, tasks, and verification evidence in sectioned lists.
- Use `[ ]` task and subtask checklists for all plans and delivery sequences,
  including future work. Give every task a short outcome description and each
  subtask a concrete action. Group related tasks into named phases.
- Mark a subtask `[x]` when its action is complete and link relevant evidence.
  Mark its parent `[x]` only after all required subtasks and its user review are
  complete. Leave blocked, deferred, or unverified work unchecked with a reason.
- Update checklists at review boundaries and session end. Task completion does
  not imply feature acceptance or permission to commit.

Start records empty; add the actual feature number, title, and relevant sections
when established facts and user answers are available. Use the section guidance
below; do not copy it into the records or prefill empty headings.
For requested spec work, always use the available question tool with proposed
options as described in [AGENTS.md](../../AGENTS.md); ask only about the requested
scope and wait for answers before filling dependent sections. In plan.md, detail
the tasks for agreed phases; add future work only when identified.
Each commit has one reviewable outcome and follows
[TDD.md](../conventions/TDD.md). Keep the three files as the feature's records
rather than creating a competing FEATURE.md or separate phase files.

## Plan phases and commits

- Group related implementation tasks into phases with a clear outcome and
  dependencies. Link tasks to the small delivery slices in ROADMAP.md; a phase
  may group related slices. Keep phases as sections in plan.md.
- Plan one to five commits per phase. Split a phase needing more than five into
  smaller related phases; do not add filler commits to meet a count.
- Give each planned commit a checklist task, a proposed message, affected areas,
  relevant acceptance criteria, and detailed action subtasks for tests,
  implementation, verification, and review as applicable.
- Describe concrete work, such as the behavior to add, the code boundary to
  change, and the scenario to verify; replace generic workflow labels with
  feature-specific steps when writing a real plan.
- Record the planned commit count per phase and respect project policy on whether
  tests and implementation share a commit or need separate commits.
- Mark a phase complete only after all required tasks and reviews are complete.
  Track whether each commit is merely planned, ready for approval, or actually
  committed; checked work alone does not claim a commit exists.

## Split oversized features

If a feature is too large to implement and review in small slices, split it into
sub-features before implementation. Give each a clear, independently verifiable
outcome, bounded requirements, and its own review checkpoints.

Give each sub-feature its own number in ROADMAP.md and a sibling folder
containing `plan.md`, `requirements.md`, and `validation.md`. Link the parent
and sub-features in both directions and record their dependencies and delivery
order through its grouping and task order in ROADMAP.md; put dependency details
in plan.md.
Keep the parent focused on the overall outcome and shared constraints; link to
the sub-feature owners instead of duplicating detailed requirements or evidence.

Present the proposed split for user review. Implement one sub-feature at a time
in small slices. Accept the parent only after its sub-features and any shared
end-to-end acceptance criteria have evidence and user acceptance.

## Writing requirements.md

- Agreement: record the actual user decision and approved scope.
- Problem and scope: group the actor, outcome, included work, non-goals, and
  constraints; link mission, stack, and decision owners.
- Behavior contract: group input rules, preconditions, visible results, side
  effects, state invariants, meaningful failures, and recovery.
- Repeated, conflicting, or partial execution: describe only relevant cases.
  Do not imply retries or rollback without an agreed requirement.
- Acceptance criteria: list concrete scenarios with stable `R-001` IDs.
- Open questions: list actual unresolved requirements raised during work.
- Related records: link plan.md and validation.md when they exist.

## Writing plan.md

- Related records: link requirements and evidence.
- Approach and boundaries: list the smallest approach, affected areas,
  dependencies, and agreed decisions.
- Delivery phases: group actual agreed outcomes into phases of one to five
  planned commits; follow the phase and commit rules above.
- Tasks: use concrete checklists with proposed messages, affected areas,
  acceptance links, tests, implementation, verification, and review steps.
- Future work: list only identified later work; omit speculative filler.
- Current checkpoint: group the active phase/commit/task, approved scope,
  working-tree changes, verification, actual blockers, and exact next action.

## Writing validation.md

- Related records: link requirements and the execution plan.
- Verification evidence: group entries by criterion or task; list the scenario,
  command or inspection, expected failure, final result, code checked, and limits.
- Distinguish change failures, environment problems, skipped checks, and checks
  not run. Passing commands alone do not prove every criterion.
- User acceptance: record the actual decision, reference, accepted scope, and
  remaining work. Do not imply acceptance from passing checks.

On resume, read requirements.md and the checkpoint in plan.md, then relevant
validation evidence and linked contracts. Compare the checkpoint with the
working tree. Update plan.md and validation.md at review boundaries and session
end, and update the corresponding roadmap checklists from that evidence. Replace
the checkpoint snapshot rather than appending a diary. Load historical records
only for a dependency.
