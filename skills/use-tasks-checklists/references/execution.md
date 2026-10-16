# Execute and resume

Compare recorded progress and checked state with the working tree before acting.
Preserve unrelated work. Execute only the requested scope: the next task for a
next-task request, or dependency-ordered tasks for an all-tasks request. Existing
authorization remains subject to project review gates. Delegate independent
tasks only when requested or required and available; preserve dependency order.

For each task, implement the named change within approved scope, run its done
check and required project verification, and record concise evidence beside it:
command or inspection, result, revision or working-tree state, and material
limits. Retain meaningful expected failures for TDD; distinguish failures,
skipped checks, and checks not run. Link long output instead of embedding logs.

Record review findings and dispositions, decisions, proposed commit messages,
and actual commit references beside the affected task. Proposals are not commits.
When requested, keep short task-referenced or timestamped entries in
`implementation_notes.md` for unexpected behavior, unplanned decisions,
assumptions, and non-obvious discoveries; link rather than duplicate task evidence.

Mark `[x]` only when the done condition and required reviews have evidence.
Leave partial work unchecked. Record explicit user acceptance with its reference,
accepted scope, and remaining work. Passing checks and completed tasks imply
neither acceptance nor permission to commit, publish, or exceed approved scope.

Do not modify `spec.md` or `plan.md` during execution. If findings require a
requirement or approach change, report what was tried, record the blocker using
the project's question convention, and stop dependent work for a decision.
Separately authorized document revisions belong to their authoring workflow.

At review boundaries and session end, update task checks, evidence, blockers,
and the next action. Reflect proven progress in the relevant ROADMAP.md group
without duplicating tasks. Keep its feature/sub-feature headings and checklists;
leave dependency details and narrative in the records. End with one concrete
next action, including required review when applicable.
