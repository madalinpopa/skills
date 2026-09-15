# Context and token optimization

Use this procedure for new skills and substantial edits. For a small correction,
apply only the relevant checks. Optimize the work an agent actually performs.

## Audit the loading path

Measure separately:

- Discovery: name and description offered before selection.
- Invocation: the whole `SKILL.md` and any unconditional reference reads.
- Conditional paths: only references needed for the selected mode or failure.
- Execution/reporting: tool output, repeated calls, clarification turns, and
  generated responses.
- Generated instructions: an asset such as `AGENTS.md` may become recurring
  context in the output project even though the asset itself is not loaded now.

Do not add unrelated conditional references to the ordinary-path baseline.
State whether external skill dependencies and project instructions are included.
Use `wc -w` for words and `wc -c` for bytes; neither is a tokenizer. Compare
measured tokens only when the tokenizer/model is known. Do not promise runtime
savings from document counts alone.

## Apply the patterns used in this catalog

1. Shorten discovery descriptions. Front-load capability and trigger, followed
   by necessary exclusions. Roughly 30–50 words often suffices; it is a useful
   editing heuristic, not a required limit. Retain distinctions needed to avoid
   misrouting. The specification's 1,024-character description limit still applies.
2. Replace mandatory full-reference reads with explicit routing by task and
   section. Keep prerequisite and compatibility checks, defaults, and essential
   operational caveats in the normal path. Architecture internals, maintenance,
   customization, and recovery can be conditional. Add contents to long references.
3. Separate advice from execution. A request to choose a layout should not
   collect scaffold parameters or create files. Gather only required values that
   cannot be inferred; use documented optional defaults. Batch related questions.
4. Delegate repetitive mechanics to existing Bash helpers. Copy assets and
   substitute validated values deterministically. Run a helper by its absolute
   installed path while resolving the target from the user's project, so changing
   into the skill directory cannot redirect generated output there.
5. Report outcomes compactly: destination, meaningful edits, checks passed,
   failed or skipped, and the next action. Show full trees or logs only when
   requested or needed to explain failure. Never hide errors or partial writes.
6. Keep authoring and maintenance instructions off the ordinary execution path.
   Avoid repeatedly loading an upstream skill or copying its general advice into
   every generated skill. Depend only on capabilities needed for the task.
7. Preserve security and operational constraints while shortening text. Essential
   secret setup, config precedence, protected-route scope, existing-file guards,
   and unit versus integration results cannot disappear into an optional reference.
8. Trace scripts before describing recovery. Distinguish preflight refusal from
   failure after writes. A helper that refuses before copying cannot be repaired
   by registration edits alone. Document complete conditional recovery, preserve
   unrelated files, and report what actually completed.

## Keep the structure honest

Keep `SKILL.md` under 500 lines and as short as the job permits. Do not use
abbreviations, compressed prose, arbitrary tiny limits, or missing prerequisites
to reduce counts. A single concise skill does not need a routing table. Prefer
one source of truth for detail, with a short reminder only for essential rules.

When choosing resources, favor a maintained script over repeated agent-authored
code, and an asset over regenerating boilerplate. Avoid adding tool calls solely
to rediscover inputs already present in the request or loaded project context.
Batch independent inspections; keep edits and dependent operations sequential.

## Verify savings without losing behavior

Review the ordinary path, an optional/customized path, missing-input handling,
a relevant failure, and a nearby request that should not trigger the skill.
Confirm each path has all required information without loading unrelated files.

Record before/after discovery and mandatory-read word counts alongside what
changed. For substantial workflows, compare realistic task execution: successful
outcome, correct boundaries, number of calls, unnecessary reads, and verbosity.
A shorter instruction that causes retries or wrong actions is not an improvement.
