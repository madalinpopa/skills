---
name: create-project-docs
description: Initializes a project's AGENTS.md, CLAUDE.md, and documentation for conventions, product specs, decisions, and features. Use when starting a project or requesting this documentation foundation, especially for Go and TypeScript. Not for application scaffolding or completing an existing product specification.
compatibility: The copy helper requires Bash and standard Unix tools.
status: published
tags: [setup, docs, conventions, go, typescript]
---

# Create project docs

Create a documentation foundation, then tailor new files from the user's
project intent and available repository evidence. Keep recurring agent
instructions small; route to ordinary documents only when the task needs them.

## Inspect and initialize

Resolve the selected project root and this skill's installed absolute directory
as `<skill-dir>`. The root must already exist. Inspect current instructions,
README, docs, manifests, task scripts, and the diff when available. In an empty
project, use the user's description; ask only for facts needed to proceed.
Do not infer that the skill's installation directory is the target project.
If an existing spec or convention document already owns the same information
under another path, resolve the proposed routing with the user before copying
a competing owner. Initialization does not authorize migrating existing docs.

Run the bundled helper by its absolute path:

```sh
bash "<skill-dir>/scripts/init-project-docs.sh" "<project-root>"
```

It copies missing files from [assets](assets/), reports created and kept paths,
and preserves existing regular files, including empty files. It refuses
symlink destinations and file/directory collisions before writing. A later
write failure can leave directories or partial files; inspect reported paths
and fix the cause before retrying. Never delete or replace existing content
merely to make a rerun succeed.

## Tailor the new documents

Use the helper's created list to identify files available for tailoring. Keep
existing documents unchanged unless their update is explicitly authorized.
If an existing instruction file needs new routes or a CLAUDE import, show the
targeted proposal and obtain authorization for that edit. Preserve Claude-only
additions. Reuse the agreed owner when existing docs use different paths.

| File | Owner / tailoring |
| --- | --- |
| [AGENTS.md](assets/AGENTS.md) | Project identity and task-based document routes |
| [CLAUDE.md](assets/CLAUDE.md) | One `@AGENTS.md` import; no reverse import |
| [CODE.md](assets/docs/conventions/CODE.md) | Shared code rules; TypeScript/frontend rules only when applicable |
| [COMMITS.md](assets/docs/conventions/COMMITS.md) | Commit proposals and delivery authorization |
| [GOLANG.md](assets/docs/conventions/GOLANG.md) | Go-specific conventions and module toolchain |
| [TDD.md](assets/docs/conventions/TDD.md) | Tests, review gates, and verification |
| [MISSION.md](assets/docs/specs/MISSION.md) | Problem, audience, outcome, and non-goals |
| [FEATURES.md](assets/docs/specs/FEATURES.md) | Capability index linking to feature requirements |
| [ROADMAP.md](assets/docs/specs/ROADMAP.md) | Delivery order and links to feature work |
| [STACK.md](assets/docs/specs/STACK.md) | Technology choices, supported versions, and project commands |
| [Decision guide](assets/docs/decisions/README.md) | Numbered decision records and their template |
| [Feature guide](assets/docs/features/README.md) | Numbered folders with plan.md, requirements.md, and validation.md templates |

Record known facts and explicit decisions. Gather missing information with the
available ask-user-question tool, or ask in chat if none is available. Start
with the minimum needed for the next document slice, using repository evidence
and previous answers first. Leave non-blocking unknowns as `Pending:` questions;
do not turn initialization into an exhaustive interview. Preserve this rule in
generated AGENTS.md so it also applies when later filling specs and feature records.
Go backend and TypeScript frontend are candidate defaults, not requirements.
Confirm them from the request or manifests; do not choose React, storage,
architecture, versions, dependencies, or deployment merely to fill a table.
For Go-only projects, mark frontend choices not applicable and remove irrelevant
frontend rules. For non-Go projects, mark Go guidance not applicable.

Use real commands and working directories from project scripts/configuration;
record missing commands as pending. Research current official docs only for
technology choices actually being documented, and record source links and
versions in STACK.md. Adapt engineering defaults to existing project policy.
Preserve the user's explicit policies during tailoring: external libraries need
user approval before adoption, installation, or version changes, and their latest
released versions must always be verified. Keep periodic `go fix` guidance in
GOLANG.md for Go projects, scoped to the supported toolchain and reviewed changes.
Preserve the adversarial-review gate in AGENTS.md and TDD.md: use an independent
subagent after each code slice, resolve findings within approved scope, and
re-review the final diff before user approval. This documents a reviewer role;
it does not enable delegation tools or install platform-specific agents.
Keep AGENTS.md around 300 words where practical, with essential constraints
and ordinary Markdown links; do not auto-import all linked documents.

Keep detailed numbered records for real work: `0001-<decision>.md` under
`docs/decisions/` and `0001-<feature>/` under `docs/features/`, containing
`plan.md`, `requirements.md`, and `validation.md`. Requirements own the behavior
contract; the plan owns implementation steps and the current checkpoint;
validation owns evidence and acceptance. FEATURES.md indexes these records.
Use each directory's guide when creating a record; do not create example
decisions or features during setup. Proposals remain proposals until approved.

For instructions-only setup, use `create-agents-setup` when available. Existing
`docs/SPEC.md` authoring belongs to `create-spec` when available; do not invoke
it to rewrite this split-spec layout. Application scaffolding is a separate
request. This skill has no mandatory dependency on those skills.

## Verify and report

Inspect the resulting files and changes. Check relative links, task routes,
CLAUDE imports, real commands, document ownership, and explicit pending facts.
Distinguish copied starter templates from agreed product requirements. Check
that existing bytes were preserved and no competing spec owner was introduced.
Do not claim application tests passed merely because documentation was copied.

Report changes and verification in short `-` bullets; avoid repeated details or
full successful logs. Include **User must be aware** with only material limitations,
partial writes, and pending questions/decisions, or "None". Record unresolved
items in the relevant project docs and ask the user to review them using the
available question tool, falling back to chat. Include a proposed commit message.
Preserve these response rules in generated AGENTS.md. Leave changes uncommitted unless
the user explicitly requests a commit. Published skill metadata grants no
permission to commit, push, or publish the target project.

The small instruction file and conditional routes follow
[OpenAI instruction guidance](https://learn.chatgpt.com/docs/agent-configuration/agents-md)
and [context guidance](https://developers.openai.com/blog/rethinking-skills-and-prompts-for-gpt-6-astra).
The single import follows [Claude memory guidance](https://code.claude.com/docs/en/memory).
