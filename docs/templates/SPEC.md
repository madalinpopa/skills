# Specification: [Application name]

[Describe what we are building and its main purpose in 2–3 sentences.]

This document defines the application's intended behavior and implementation
boundaries. Use it to agree on what to build and verify the result. Keep it
current when agreed behavior changes; distinguish proposals from decisions.

Replace placeholders with project-specific details. Use the CLI examples for
command-line applications or the web alternatives for web applications. Remove
unused alternatives and mark unresolved decisions as pending rather than
inventing requirements. Keep only applicable prompts; mark relevant unsupported
behavior as a non-goal. A prompt about retries, concurrency, or recovery does
not require adding that capability.

## Goal

- **Problem**: [What is difficult or missing today?]
- **Outcome**: [What should the application enable or improve?]
- **Success criteria**: [Observable user or product outcomes that demonstrate the goal is met]
- **Non-goals**: [Related problems this application will not address]

### Behavioral acceptance

Define verifiable requirements for the agreed scope. Link each to its command
or feature contract and a representative scenario below. Keep execution
checkpoints and test evidence in feature/phase documents when used; link to
them instead of duplicating their progress here.

| Requirement | Observable acceptance condition | Contract / scenario |
| --- | --- | --- |
| [Short name or ID] | [Given an input/state, the required result and side effects] | [Links to relevant sections] |

## User

- **Primary user**: [Who uses the application and for what purpose?]
- **Context**: [Where, when, and how often they use it]
- **Experience and access**: [Technical knowledge, permissions, and prerequisites]
- **Other users**: [Additional roles or automated consumers, if relevant]

## Commands / Features

List the supported capabilities. For a CLI, name commands and relevant arguments
or flags. For a web application, name features and where users access them.
Keep planned capabilities clearly separate from the agreed scope.

| Command / feature | Purpose | Inputs and defaults | Result and side effects |
| --- | --- | --- | --- |
| [Command syntax or feature name] | [User need] | [Required/optional input and defaults] | [Output and state changes] |

### Contract: [Command / feature]

Repeat for capabilities whose rules need more detail than the overview table.
Link to shared rules below rather than repeating them for each capability.

- **Preconditions**: [Required state, configuration, permissions, or dependencies]
- **Input rules**: [Valid values, missing versus empty input, and conflicting options]
- **Selection and scope**: [Which items/targets are affected; omitted selection, unknown items, and no matches]
- **Result**: [Required output and state changes, including when no change is needed]
- **Repeated execution**: [Result when the same action is repeated; whether effects can be duplicated]
- **Failure rules**: [Feature-specific errors and links to shared failure/recovery guarantees]

## Usage and Outputs / Outcomes

Use short stories that demonstrate the commands or features above. Include
representative success and failure cases. State relevant prerequisites and
observable side effects; examples should agree with the defined behavior.

Include empty results, repeated actions, conflicts, and partial success where
they affect agreed behavior. State which output details are contractual (such
as fields, streams, and exit codes) and which are illustrative (such as sample
names or decorative spacing). Do the same for web outcomes where wording or
presentation matters. Link scenarios to the requirements they demonstrate.

### CLI: [User completes a task]

**Story**: As [user], I want to [action] so that [benefit].

**Given**: [Required configuration, input, permissions, or existing state]

**Invocation**:

```text
app command <argument> --option <value>
```

**Standard output**:

```text
[Expected output, or no output]
```

- **Standard error**: [Expected diagnostics, or none]
- **Exit code**: [Expected code and meaning]
- **Outcome**: [Files, records, or external state created/changed; or no changes]

### CLI: [Task cannot complete]

**Given**: [Invalid input, missing prerequisite, or relevant dependency failure]

**Invocation**:

```text
app command <failing-input>
```

**Standard error**:

```text
[Expected error and actionable next step]
```

- **Standard output**: [Any output produced before failure, or none]
- **Exit code**: [Expected code and meaning]
- **Outcome**: [State left behind, including any partial changes]

### Web alternative: [User completes a task]

**Story**: As [user], I want to [action] so that [benefit].

**Given**: [Authentication, permissions, and relevant existing state]

1. [Open the relevant page or feature.]
2. [Enter input and perform an action.]
3. [Observe the expected result.]

- **Outcome**: [Visible confirmation and resulting application state]
- **Failure case**: [Trigger, displayed error, retained input/state, and recovery]

## Implementation details

Record decisions and their reasons at the level needed to implement the
behavior above. Link to detailed references instead of copying them here.

### Stack

| Area | Choice | Reason / operational details |
| --- | --- | --- |
| Language and runtime | [Language and supported version] | [Why it fits] |
| Frameworks | [CLI/web framework, or none] | [Responsibility] |
| Storage | [Database, files, or no persistence] | [Data location and lifecycle] |
| Local execution | [Setup and run commands] | [Required services and configuration] |
| Build and distribution | [Build command and delivery method] | [Supported artifacts/platforms] |
| Deployment | [Hosting/install target, or not applicable] | [Release and upgrade approach] |

### External libraries / dependencies

Include libraries, external executables, and services required by the application.
Record supported versions and link to official documentation for chosen APIs.

| Dependency | Version / compatibility | Purpose | Official documentation |
| --- | --- | --- | --- |
| [Library, executable, or service] | [Version requirement] | [Why it is needed] | [Link] |

### Architecture

- **Execution flow**: [How a command/request moves from input to outcome]
- **Boundaries**: [Responsibilities and allowed dependency directions]
- **State and integrations**: [Who owns persistence and external interactions]

[Add a small diagram only when it clarifies these relationships.]

### State and invariants

Describe the state needed by the agreed behavior, including user-owned data
the application touches. State rules independently of the storage mechanism.

- **Authoritative data**: [Source of truth; how derived or cached data relates to it]
- **Identity and ownership**: [How items are identified; what the application may change and what it must preserve]
- **Lifecycle**: [Valid states and allowed transitions, including creation and deletion]
- **Invariants**: [Rules that must hold across operations, such as uniqueness or target containment]
- **Invalid existing state**: [How missing, malformed, or incompatible state is reported and handled]
- **Concurrent access**: [Supported overlap between operations, or explicit limits and assumptions]

### Failure and recovery

Define shared observable guarantees here; keep exceptions in the relevant
command/feature contract. Distinguish expected validation or conflict failures
from unexpected failures after work begins. State limits explicitly, including
when rollback or recovery is unsupported.

- **Validation before changes**: [What is checked before side effects, and across which selected items]
- **Completion boundary**: [Which changes must succeed or fail together, and where partial progress is possible]
- **Continuation**: [Which failures stop processing and which allow unrelated items to continue]
- **Interruption**: [State and reporting guarantees on cancellation, timeout, or process termination, where applicable]
- **Reporting**: [How completed, skipped, and failed work is distinguished; error/exit status when outcomes are mixed]
- **Recovery and retry**: [State left behind, available recovery steps, and when repeating an action is safe]

### Configuration and initialization

Keep this section when configuration or first-run setup affects behavior.

- **Sources and precedence**: [Supported flags, environment variables, files, and defaults, in precedence order]
- **Value rules**: [Required values; missing, empty, and invalid values; treatment of unknown keys]
- **Path resolution**: [Base for relative paths, supported expansion, and local/global scope where applicable]
- **First run**: [What is created or fetched, when initialization happens, and which operations trigger it]
- **Read-only operations**: [Whether inspection/help/preview operations initialize or change state, and whether they access the network]
- **Initialization failure**: [State left behind and how setup can be retried; reference shared recovery rules]

### Project layout

Show the intended layout with brief responsibilities. Include only directories
and files the project needs; replace this illustrative structure.

```text
<project-root>/
  <entrypoint-file>       Process startup
  <application-dir>/     Commands or request handling and use cases
  <domain-dir>/          Domain behavior
  <integration-dir>/     Persistence and external systems
  <docs-dir>/            Project documentation
```

### Identified domains / modules

Identify cohesive responsibilities, their owned data, and interactions. A domain
does not automatically require a separate package, service, or database.

| Domain / module | Responsibility and owned data | Interactions |
| --- | --- | --- |
| [Name] | [Behavior and state it owns] | [Modules or external systems it uses] |

### Constraints

- **Compatibility**: [Supported platforms, runtimes, and stable interfaces]
- **Resources**: [Relevant performance, memory, storage, or scale limits]
- **Environment**: [Network/offline requirements and process permissions]
- **Data and security**: [Access, secrets, retention, and data-integrity requirements]
- **Product boundaries**: [Scope, delivery, or operational limits]

### Additional details

Keep only topics relevant to this application that are not covered above:
output formats, logging, migrations, verification strategy, or other
cross-cutting behavior.

- **[Topic]**: [Required behavior and rationale]
- **Pending decisions**: [Question, impact, and what is needed to resolve it]

### Suggested skills

List agent skills that would help with this project's work. Explain when each
applies and where to find it; do not assume it is installed or enabled.

| Skill | When to use it | Purpose | Source / availability |
| --- | --- | --- | --- |
| [Skill name] | [Specific task or trigger] | [How it helps] | [Path/link and availability] |

## Research

Keep concise findings and direct source links. Verify sources when filling out
this section; record versions and review dates so later readers can assess
freshness. Distinguish verified facts from assumptions and proposed decisions.

Define the questions before searching. Stop when the important questions have
evidence-backed answers; record remaining uncertainty. If searches find no
suitable alternative, describe that result and the search scope rather than
claiming no alternatives exist.

- **Research questions**: [Decisions to resolve, such as whether an existing tool
  meets the goal, which dependency fits the constraints, or which conventions apply]
- **Search scope**: [Relevant queries, sources searched, and any limitations]

### Documentation and specifications

Research the latest official documentation for relevant specifications,
standards, libraries, frameworks, runtimes, and services. Check applicability
to the versions selected for this project. Link to specific sections rather
than copying documentation; reference existing dependency entries where useful.

For CLI projects, consider the [Command Line Interface Guidelines](https://clig.dev/)
for help, output, errors, configuration, and scripting behavior.

| Topic | Official source / section | Version / revision | Relevance and findings | Last checked |
| --- | --- | --- | --- | --- |
| [Specification, library, or other topic] | [Direct link] | [Applicable version] | [Supported behavior, constraint, or decision] | [YYYY-MM-DD] |

### Similar tools and applications

Research existing tools serving the same users or solving a similar problem.
Use both GitHub search for public repositories and regular web search to find
relevant tools and applications.
Use official product documentation or repositories to verify comparisons.
Start with a shortlist of 3–5 relevant tools, or fewer if fewer are found.
Explain why each matters; expand only when a research question remains unresolved.
Summarize what already exists, where it differs from this application's goal,
and what we can learn or reuse without expanding the agreed scope.

| Tool / application | Official website / repository | Users and overlapping capabilities | Relevant differences and lessons | Last checked |
| --- | --- | --- | --- | --- |
| [Name] | [Direct link] | [Who it serves and what it solves] | [Strengths, gaps for our use case, or useful approaches] | [YYYY-MM-DD] |

For each shortlisted repository, record reuse suitability in concise notes:
declared license, maintenance status, supported platforms, and integration
options, with supporting links. Distinguish documented claims from behavior
actually tested; mark unknown details as unverified.

### Findings and open questions

- **Implications**: [How the evidence informs the goal, scope, or implementation]
- **Proposed decisions**: [What to adopt, reuse, or avoid, with source references]
- **Recommendation**: [Use an existing tool / extend or integrate one / build this application]
- **Reason**: [Evidence, unmet requirements, and tradeoffs supporting the recommendation]
- **Unresolved questions**: [What remains unverified and how to investigate it]
