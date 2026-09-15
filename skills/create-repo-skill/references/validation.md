# Validation and evaluation

## Format and executable checks

Run the local `scripts/validate-skill.sh` on the source skill. It validates the
repository's extended frontmatter directly; no external authoring skill or
manual removal of store fields is needed. Preserve source files while checking
shared and Claude-specific metadata rules. Supported optional/unknown fields
must not be rejected merely because an old validator used a small allowlist.

Check names, descriptions, explicit status, tag types, unique YAML keys, and
invalid `x-claude` collisions. Reject unfinished scaffold instructions. Check
relative links, reference anchors, packaging, and changed assets during review.
Tool validation is bounded: valid YAML and a valid name do not prove useful
instructions or correct agent behavior. Report unsupported checks honestly.

Exercise new or changed scripts against their public command interface in
isolated temporary directories. Test observable behavior: expected files,
correct defaults, preserved existing content, malformed input, and realistic
failures. Avoid tests that only match headings, copy implementation logic, or
exercise libraries independently. Follow the repository's test-first review
checkpoints for behavioral helper changes.

Use Bash for orchestration and a proper YAML parser where needed. Report a
missing Python/PyYAML runtime instead of claiming validation passed or installing
packages silently. See the scripts reference for one-call environment setup.
Keep test output concise while retaining the actual failure reason.

## Behavioral review

Test representative user requests against the description, including a nearby
request that should not trigger it. Verify normal execution, the relevant
conditional mode, failed prerequisites, and user authorization boundaries.
Check that scripts and commands exist, references are discoverable, required
inputs/defaults are clear, and output claims match checks actually run.

For workflow skills, verify component availability and every input/output
handoff. Missing dependencies or failed steps stop dependent work. Keep skill
dependencies acyclic; tool metadata is not a skill dependency installer.

Preserve default automatic invocation unless explicitly asked for explicit-only
use. Preserve existing sidecar policy and dependencies when editing UI metadata.
Do not equate invocation permission with permission for external actions.

## Independent forward-testing

Use an independent subagent evaluation when complexity or risk makes realistic
behavioral validation worthwhile and delegation is available and authorized.
Ordinary edits do not automatically require subagents.

Give the evaluator the actual skill, a realistic user request, and the minimum
raw artifacts required. Do not include the intended answer, suspected defect,
proposed fix, or prior conclusions unless the task requires that information.
For example: ask it to use the skill to initialize an example inside a supplied
temporary repository, then inspect the files and actions it actually produced.

Constrain evaluation to permitted resources and side effects. Use isolated
workspaces; do not let generated files enter the working tree or another run's
fixtures. Obtain additional authorization only when the evaluation exceeds the
existing scope, affects a live service, or imposes material cost. Respect a
missing dependency or approval boundary rather than bypassing it.

Review outcomes and artifacts. Correct demonstrated weaknesses narrowly, then
repeat only the affected checks. Record limitations, including what was not
executed; successful format validation is not an end-to-end evaluation.
