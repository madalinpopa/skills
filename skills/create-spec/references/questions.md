# Questions by section

For each section of `docs/templates/SPEC.md`: where the evidence usually lives,
what to ask, and what a complete answer contains. Ask only what the evidence
does not answer. Keep the user's wording.

## Goal

**Evidence**: README introduction, repository description, first commits, any
issue or note that started the project.

**Ask**:

- What is hard or missing today without this application?
- What should a user be able to do once it exists?
- How would we know it worked? Name things a user can observe, not internal
  metrics.
- Which nearby problems will it not solve, even if asked?

**Complete when**: the problem and outcome are one or two sentences each, every
success criterion is observable, and at least one non-goal is named.

## User

**Evidence**: README audience, install instructions, CI or automation that runs
the application, permissions it asks for.

**Ask**:

- Who runs it, and to get what done?
- Where and how often: a terminal on a laptop, a CI job, a browser, a phone?
- What must they already know or have: accounts, tokens, installed tools?
- Does anything other than a person call it, such as a script or another
  service?

**Complete when**: the primary user and their context are concrete enough to
judge default behavior and error wording.

## Commands / Features

**Evidence**: the command tree or router, help text, existing tests, README
usage examples.

**Ask**:

- List every capability in the agreed scope. For a CLI, the command with its
  arguments and flags; for a web application, the feature and where it lives.
- For each: what need it serves, which inputs are required, which are
  optional and their defaults, what it outputs, and what state it changes.
- Which capabilities are planned but not agreed? Keep them in a separate
  list, not in the table.

**Complete when**: each row can be read as a promise, and the planned list is
clearly outside the agreed scope.

### Contract

Write one contract for each capability whose rules do not fit in the table.
Ask, per capability:

- What must be true before it runs? Configuration, permissions, state,
  dependencies.
- Which inputs are valid? What happens with a missing input, an empty one, or
  two options that conflict?
- Which items does it act on? What happens when the selection is omitted,
  names an unknown item, or matches nothing?
- What does it produce, including when there is nothing to do?
- What happens when it runs twice? Can effects double?
- Which failures are specific to it? Point to the shared failure rules for the
  rest.

Link shared rules instead of repeating them in each contract.

## Usage and Outputs / Outcomes

**Evidence**: README examples, integration tests, golden files, screenshots.

**Ask**, per scenario:

- Who is the user and what do they want? Story form: as, I want, so that.
- What exists before: config, input, permissions, state?
- The exact invocation or the steps in the interface.
- The output: standard output, standard error, exit code, or the visible
  confirmation. Which parts are contractual and which are illustrative?
- What is left behind: files, records, external state, or nothing?

Cover at least one success and one failure per capability that changes state.
Add scenarios for empty results, repeated actions, conflicts, and partial
success when the contracts define them. Name the requirement each scenario
demonstrates.

## Implementation details

Ask these subsection by subsection. Most answers already exist in the
repository; confirm them rather than asking from scratch.

### Stack

**Evidence**: module or package file, Taskfile or Makefile, Dockerfile,
compose file, CI workflow, release configuration.

Confirm: language and version, frameworks and what each is responsible for,
storage and where data lives, how to run locally, how to build and ship, where
it is deployed or installed. Ask for the reason behind each choice when the
repository does not say.

### External libraries / dependencies

**Evidence**: dependency file, external executables the code calls, services
in the compose file or environment variables.

For each: version constraint, why it is needed, and a link to its official
documentation. Do not add dependencies the code does not use.

### Architecture

**Evidence**: package layout, entrypoint, existing diagrams, agent instructions
on package boundaries.

Ask: how does a request or command travel from input to result? Which
packages may depend on which? Who owns persistence and each external
integration? Add a diagram only if the prose is unclear without it.

### State and invariants

**Evidence**: schemas, migrations, lock or state files, validation code.

Ask: what is the source of truth, and what is derived from it? How are items
identified, and what must the application never change? Which states exist and
which transitions are allowed? Which rules must always hold? What happens with
missing, malformed, or incompatible existing state? Can two runs overlap, or
is that a stated limit?

### Failure and recovery

**Evidence**: error handling, transactions, backup or rollback code, exit
codes in tests.

Ask: what is validated before any side effect? Which changes succeed or fail
together? After one failure, does unrelated work continue? What happens on
cancellation or a killed process? How does the report separate done, skipped,
and failed work, and what is the exit status when mixed? What is left behind
and when is a retry safe?

An honest "no rollback; rerun the command" is a valid answer. Record it.

### Configuration and initialization

**Evidence**: config loading code, flags, environment variables, first-run
logic, default paths.

Ask: which sources exist and in what precedence? What is required, and what
happens with missing, empty, invalid, or unknown values? How are relative paths
resolved? What does the first run create or fetch, and which commands trigger
it? Do help, list, or preview commands change state or use the network? What
is left behind when setup fails, and how is it retried?

Drop this subsection only when the application has no configuration and no
setup.

### Project layout

**Evidence**: the tree itself.

Show the intended tree with one short responsibility per entry. Include only
what the project needs.

### Identified domains / modules

**Evidence**: package names, tables, module directories.

Ask: which cohesive responsibilities exist, what data does each own, and which
other modules or external systems does it talk to? A domain is a boundary in
the design, not necessarily a package.

### Constraints

**Evidence**: CI matrix, supported platforms in the README, resource limits in
deployment files, security notes.

Ask: supported platforms and stable interfaces; performance, memory, storage,
or scale limits; network and permission needs; secrets, access, retention, and
data integrity; product or delivery limits.

### Additional details

Ask whether anything cross-cutting is not covered yet: output formats,
logging, migrations, verification strategy. Then list the pending decisions
collected so far, each with its impact and what resolves it.

### Suggested skills

**Evidence**: `.claude/skills/`, `.agents/skills/`, agent instruction files,
the skills store if the CLI is configured.

For each skill that helps with this project: when to use it, what it does, and
where it comes from. State whether it is installed; do not assume it is.

## Research

Agree on the questions before searching. Typical ones: does an existing tool
already do this, which dependency fits the constraints, which conventions
apply. Record the search scope: queries, sources, and limits.

### Documentation and specifications

For each library, runtime, standard, or service in the stack: the official
source, the section that matters, the version checked, the finding, and the
date. Link to sections, do not copy. For a CLI, check the Command Line
Interface Guidelines at https://clig.dev/ for help, output, errors, and
configuration behavior.

### Similar tools and applications

Search GitHub and the web. Shortlist three to five tools. For each: link, who
it serves, what overlaps, where it differs, and what to learn or reuse. Note
license, maintenance status, platforms, and integration options with links.
Mark anything not verified as unverified. If nothing suitable turns up, record
the queries and scope instead of claiming nothing exists.

### Findings and open questions

Ask the user to decide, with the evidence in front of them: what to adopt,
reuse, or avoid; whether to use an existing tool, extend one, or build this
application; and why. List what remains unverified and how to check it.

## Behavioral acceptance

Fill this last. For each capability in the agreed scope, write one or more
rows: a short name, the observable condition given an input or state, and
links to the contract and the scenario that demonstrates it. Every row needs
both links. If a requirement has no scenario yet, go back and add one.
