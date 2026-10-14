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
when their roots differ. Finding a command is not proof that it was run.

- Focused tests: Pending: command and working directory.
- Full default tests: Pending: command and working directory.
- Format: Pending: command and working directory.
- Lint / static analysis: Pending: command and working directory.
- Type checking, when applicable: Pending: command and working directory.
- Build: Pending: command and working directory.
- Integration, when applicable: Pending: command, working directory, prerequisites,
  and any opt-in gate for external services.

## Local setup

- Prerequisites: Pending.
- Development commands and working directories: Pending.
- Configuration names: Pending: reference secret provisioning without recording
  secret values in project docs.
