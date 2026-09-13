# Names and category legend

Use `<action>-<subject>`: `create-pr`, `create-issue`, `create-spec`, `use-jj`,
`use-gh`, `review-tests`, or `review-code`. Add a qualifier only when it
distinguishes a real specialization, such as `review-go-code` or
`create-github-issue`. Choose the primary outcome as the prefix; the tool used
to achieve it usually belongs in the description or tags.

Keep names to 1–64 lowercase ASCII letters, digits, and hyphens, with no leading,
trailing, or repeated hyphens. The folder name and frontmatter `name` must match.
For Claude compatibility, avoid names containing `claude` or `anthropic` and
the reserved name `synced`. These are conventions for new skills; renaming an
existing skill is a separate change because installed names are its identity.

This is a practical vocabulary for this library, not a required set of skills
or a ranking of usage. Start with `create`, `use`, and `review`; use the other
prefixes when a distinct, reusable workflow needs them.

| Prefix | Primary outcome | Examples |
| --- | --- | --- |
| `create-` | Produce a new artifact | `create-pr`, `create-issue`, `create-spec` |
| `use-` | Apply tool-specific working conventions | `use-jj`, `use-gh` |
| `review-` | Assess existing work and report findings | `review-code`, `review-tests` |
| `plan-` | Define scope, decisions, and implementation steps | `plan-feature`, `plan-migration` |
| `investigate-` | Gather evidence and identify causes | `investigate-ci`, `investigate-performance` |
| `fix-` | Correct a demonstrated defect | `fix-ci`, `fix-flaky-tests` |
| `refactor-` | Improve structure while preserving behavior | `refactor-go-package` |
| `test-` | Exercise behavior and report verification results | `test-api`, `test-ui` |
| `document-` | Explain or maintain documentation for existing behavior | `document-api`, `document-architecture` |
| `configure-` | Set up or adjust a tool or environment | `configure-lint`, `configure-ci` |
| `migrate-` | Move between versions, formats, or systems | `migrate-config`, `migrate-database` |
| `release-` | Prepare or carry out a release workflow | `release-cli`, `release-package` |
| `workflow-` | Coordinate other skills toward a larger outcome | `workflow-deliver-feature`, `workflow-prepare-release` |

Choose one prefix for the primary outcome. Use `workflow-<outcome>` when the
skill coordinates other skills; keep its component skills independently useful.
`create-tests` writes tests; `test-api` runs a verification workflow;
`review-tests` assesses test quality. `use-gh` supplies GitHub CLI conventions;
`create-pr` owns the pull-request outcome. A name never grants permission to
commit, publish, deploy, or perform other external actions.

Categories sort by name in a flat catalog. Create `skills/create-pr/`, not
`skills/create/pr/`: the CLI discovers only immediate skill directories.
Use tags for other dimensions such as `go`, `github`, or `testing`.

## Workflow skills

Use `workflow-<outcome>` for a skill that invokes other skills to complete a
larger task. For example, `workflow-deliver-feature` could coordinate
`plan-feature`, `create-tests`, `review-code`, and `create-pr` once those skills
exist and are available. These are illustrative names, not installed dependencies.

Keep the workflow focused on coordination:

- Name the required and optional skills and explain when each is used. Check
  availability before a dependent step; report missing requirements instead of
  claiming a skill ran. Use a fallback only when the workflow defines one.
- Specify the input, expected output, and handoff for each step. Load and invoke
  the selected skill through the current agent's supported mechanism. Where
  explicit invocation is unavailable, read and follow its `SKILL.md`.
- Keep detailed task instructions in the component skills. A sequence of ordinary
  tool commands alone does not make a skill a `workflow-*` skill.
- State completion checks and where to stop for missing input, a failed step,
  or required approval. Honor existing user authorization and each component's
  applicable constraints; orchestration grants no additional permissions.
- Keep dependencies acyclic: a workflow must not call itself, directly or through
  another skill. Calling a skill does not require spawning a subagent.

The CLI copies skill files; it does not resolve or install dependent skills.
Document dependencies in the workflow instructions. `agents/openai.yaml`
`dependencies.tools` describes tools, not skill-to-skill dependencies.
