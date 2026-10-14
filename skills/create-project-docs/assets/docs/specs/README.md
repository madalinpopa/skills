# Writing specs

## Working rules

- Start spec files empty; add a title and sections only with real content.
- Group related facts under descriptive headings; use short bullet lists.
- Keep one point per bullet. Avoid long paragraphs and prose-heavy tables.
- Use tables only for compact comparisons or mappings.
- Record established facts and user decisions; never invent scope to fill space.
- Omit unknown topics, empty headings, placeholders, and generic questions.
- For requested spec work, use the question tool with concrete proposed options.
- Ask about the requested scope, reuse earlier answers, and wait for decisions.
- If no question tool is available, explain and ask option-based questions in chat.
- Link shared facts to their owner instead of duplicating them.

## MISSION.md

- Problem and users: describe the actual problem, affected users, and current
  alternatives. Group these together when their context is shared.
- Intended outcome: state what improves for users and how success can be observed.
- Scope and non-goals: list included outcomes, deliberate exclusions, and real
  constraints. Do not turn undecided work into a non-goal.
- Open questions: record only unresolved product questions raised during work.
- Keep detailed feature behavior in feature requirements; tasks belong in ROADMAP.md.

## ROADMAP.md

- Group features under headings; nest related sub-features beneath them.
- Use stable four-digit feature IDs; check the roadmap and feature folders first.
- List only identified tasks and subtasks as short checklists within each group.
- Order work by dependencies and useful delivery; keep dependency details in plans.
- Link to relevant plans or evidence directly from tasks when available.
- Mark tasks `[x]` only with evidence; task completion does not mean user acceptance.
- Include no status tables, narrative, acceptance criteria, or workflow guidance.
- Leave the file empty when no feature or task is clear.

## STACK.md

- Technology choices: list actual languages, frameworks, storage, services,
  architecture boundaries, build tools, and deployment choices when known.
- For each choice, record the selected version, repository evidence or user
  decision, official documentation link, and verification date.
- Compare the selected and latest released versions; do not upgrade automatically.
- Verification commands: group tests, formatting, lint, type checks, build, and
  integration commands. Include working directories and required opt-in gates.
- Copy exact commands from manifests or scripts; discovery does not prove a run.
- Local setup: list prerequisites, development commands, and configuration names.
- Never record secrets or choose technologies merely to complete the document.

## Feature and decision records

- Follow the [feature guide](../features/README.md) for behavior, plans, and evidence.
- Follow the [decision guide](../decisions/README.md) for significant agreed choices.
- Add sections progressively using these same concise writing rules.
