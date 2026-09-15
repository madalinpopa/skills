# Skill design

Adapted and reorganized for this repository from the Apache-2.0 material
identified in `../NOTICE`. Local authoring defaults and commands take precedence.

## Contents

- [Useful instructions](#useful-instructions)
- [Choose the structure](#choose-the-structure)
- [Resource roles](#resource-roles)
- [Create and iterate](#create-and-iterate)

## Useful instructions

Assume the agent is capable. Include non-obvious constraints, local knowledge,
fragile operations, and decisions that change the result. Remove generic advice,
repetition, speculative edge cases, and examples that do not clarify a task.

Preserve the user's intent and scope. Do not replace their product choice,
modify unrelated configuration, or turn a past failure or personal example
into a universal requirement. Existing authorization applies within its scope;
a skill or orchestration step grants no additional permissions. Define clear
stopping conditions for failed prerequisites, partial work, and external writes.

Match specificity to risk. Open-ended work needs outcomes and decision criteria;
fragile operations may need exact commands, deterministic scripts, and narrow
parameters. Distinguish requirements from recommendations. Avoid rigid sequences,
fixed numbers of steps, and detailed abstractions without a concrete reason.

Keep skills self-contained. Require another skill or tool only when the workflow
needs it and it is available in the target environment. Specialized review,
hardening, and audit procedures apply when requested or needed, not merely
because a routine task touches the topic. Calling a skill is not delegation.

## Choose the structure

A skill is a folder containing YAML frontmatter plus a Markdown body in
`SKILL.md`. `name` and `description` are required; preserve supported optional
fields. Detailed repository constraints live in the authoring reference.

Information loads in three stages: name/description during selection, the
entrypoint when invoked, and supporting resources only when relevant. Make
each stage independently economical. A high size limit is not a target.

A short single-purpose skill may need only `SKILL.md`. Multiple substantial
modes justify routing. For example, a deployment skill can keep selection and
shared constraints in the entrypoint and separate AWS, Azure, and GCP details;
choosing AWS should not require reading the other providers. Similarly, ordinary
document edits need not load redlining or file-format internals.

Do not split files merely to hit a line count. Moving content into a reference
saves nothing if every run is still told to read it. Keep enough context in the
entrypoint to choose the correct mode and execute the ordinary path safely.

## Resource roles

| Resource | Purpose | Load behavior |
| --- | --- | --- |
| `scripts/` | Deterministic or repeated operations, such as scaffold creation, PDF rotation, transformations, or data processing | Execute without reading implementation unless adapting or debugging it |
| `references/` | Conditional procedures, schemas, policies, format details, and substantial examples | Link directly from the entrypoint and name the triggering task |
| `assets/` | Templates, images, icons, fonts, slides, and boilerplate copied into output | Treat as output material, not instructions |
| `agents/openai.yaml` | UI metadata, tool dependencies, and invocation settings for supported hosts | Create or change only when useful or requested |

Use descriptive file names and forward slashes. Keep maintained task-specific
knowledge in one place. Do not copy entire external manuals or catalogs when
a current authoritative source is enough. Inspect callers and purpose before
removing existing resources. Add a contents section or useful search terms to
long references so agents can retrieve relevant sections.

Avoid unnecessary README files, installation guides, changelogs, duplicated
quick references, empty directories, and placeholder examples. Add them only
for an actual packaging or task requirement. A useful template belongs in
assets; a naming convention or schema belongs in references.

## Create and iterate

1. Identify the outcome, intended requests, inputs, defaults, scope boundaries,
   and nearby requests that should use another skill. Inspect the catalog for
   overlap when adding or expanding capability.
2. Choose a clear action-oriented name using the local naming reference. Do
   not rename installed identities as an incidental cleanup.
3. Select only resources that make execution more reliable or cheaper. A
   recurring transformation may justify a Bash helper; a website starter may
   justify template assets; an analysis task may need a schema reference.
4. Initialize a new skill with the local Bash command. Request only necessary
   resource directories. Do not initialize an existing skill again. Replace
   unfinished scaffold text before validation and publication. Add examples
   only when they represent a real supported workflow.
5. Write imperative, task-specific instructions. State outcomes, constraints,
   commands, and conditional reference routes. Keep detailed workflows and
   tool inventories out of the discovery description.
6. Validate format, exercise scripts, and review representative agent behavior.
   Make narrow corrections based on observed failures. Do not accumulate a
   universal instruction for every incidental example.

New skills here are `published` by default; users may explicitly request a
draft. Existing metadata and invocation policy remain unless a change is
requested. Automatic discovery remains enabled by default; sensitive actions
need appropriate authorization, not automatic exclusion from discovery.
