# Stack and project commands

Record chosen technologies, supported versions, responsibilities, and official
documentation sources. Choices need project evidence or user agreement;
architecture tradeoffs belong in [decision records](../decisions/README.md).

## Technology choices

- Backend: Pending: confirm Go or another language and its supported version.
- Frontend: Pending: confirm whether a frontend is needed and whether it uses TypeScript.
- Framework: Pending: choose only if required; React is not assumed.
- Storage and external services: Pending: record actual needs and choices.
- Architecture and source layout: Pending: record selected boundaries and paths.
- Package/build tools: Pending: record the tools used by this project.
- Deployment: Pending: record the chosen environment when relevant.

For each resolved entry, include the chosen version, evidence or decision link,
official documentation URL, and date verified. Mark unused components not
applicable. Do not install dependencies merely to complete this document.
Always verify the latest released library version from an official source and
record it alongside the selected version. Obtain user approval before dependency
additions, installations, replacements, or upgrades, following
[CODE.md](../conventions/CODE.md). An existing version pin is not permission to
upgrade to the latest release.

## Verification commands

Replace pending entries with exact commands and their working directories from
the project's manifests and scripts. Record Go and frontend commands separately
when their roots differ. Command discovery is not proof a command was executed.

| Check | Command | Working directory |
| --- | --- | --- |
| Focused tests | Pending | Pending |
| Full default tests | Pending | Pending |
| Format | Pending | Pending |
| Lint / static analysis | Pending | Pending |
| Type checking, when applicable | Pending | Pending |
| Build | Pending | Pending |
| Integration, when applicable | Pending | Pending |

## Local setup

Pending: record prerequisites, development commands, and configuration names.
Reference secret provisioning without recording secret values in project docs.
