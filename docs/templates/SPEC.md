# Specification: [Application name]

[Describe what we are building and its main purpose in 2–3 sentences.]

This document defines the application's intended behavior and implementation
boundaries. Use it to agree on what to build and verify the result. Keep it
current when agreed behavior changes; distinguish proposals from decisions.

Replace placeholders with project-specific details. Use the CLI examples for
command-line applications or the web alternatives for web applications. Remove
unused alternatives and mark unresolved decisions as pending rather than
inventing requirements.

## Goal

- **Problem**: [What is difficult or missing today?]
- **Outcome**: [What should the application enable or improve?]
- **Success criteria**: [Observable conditions that demonstrate the goal is met]
- **Non-goals**: [Related problems this application will not address]

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

## Usage and Outputs / Outcomes

Use short stories that demonstrate the commands or features above. Include
representative success and failure cases. State relevant prerequisites and
observable side effects; examples should agree with the defined behavior.

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
- **Failure handling**: [How errors, cancellation, and partial work are handled]

[Add a small diagram only when it clarifies these relationships.]

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
configuration sources and precedence, output formats, logging, migrations,
verification strategy, or other cross-cutting behavior.

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

### Documentation and specifications

Research the latest official documentation for relevant specifications,
standards, libraries, frameworks, runtimes, and services. Check applicability
to the versions selected for this project. Link to specific sections rather
than copying documentation; reference existing dependency entries where useful.

| Topic | Official source / section | Version / revision | Relevance and findings | Last checked |
| --- | --- | --- | --- | --- |
| [Specification, library, or other topic] | [Direct link] | [Applicable version] | [Supported behavior, constraint, or decision] | [YYYY-MM-DD] |

### Similar tools and applications

Research existing tools serving the same users or solving a similar problem.
Use official product documentation or repositories to verify comparisons.
Summarize what already exists, where it differs from this application's goal,
and what we can learn or reuse without expanding the agreed scope.

| Tool / application | Official website / repository | Users and overlapping capabilities | Relevant differences and lessons | Last checked |
| --- | --- | --- | --- | --- |
| [Name] | [Direct link] | [Who it serves and what it solves] | [Strengths, gaps for our use case, or useful approaches] | [YYYY-MM-DD] |

### Findings and open questions

- **Implications**: [How the evidence informs the goal, scope, or implementation]
- **Proposed decisions**: [What to adopt, reuse, or avoid, with source references]
- **Unresolved questions**: [What remains unverified and how to investigate it]
