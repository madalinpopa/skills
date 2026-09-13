---
name: demo
description: Demonstrates Agent Skills packaging with a bundled script, reference, and output template. Use when exploring this demo, learning the skill folder layout, or smoke-testing a skills CLI installation.
license: MIT. See LICENSE.txt for terms.
compatibility: Requires Python 3 for the bundled script. No third-party packages or network access are needed.
allowed-tools: Read
metadata:
  purpose: installation-example
status: published
tags: [demo, example, testing]
x-claude:
  disable-model-invocation: true
  user-invocable: true
---

# Demo

Use this working example to explain skill resources or exercise an installed
copy. Resolve the paths below relative to the directory containing this file.

## Explore the layout

```text
demo/
├── SKILL.md
├── LICENSE.txt
├── scripts/show-template.py
├── references/guidelines.md
└── assets/report-template.md
```

Read [references/guidelines.md](references/guidelines.md) when explaining the
format or adapting this example into a new skill. Keep only the resources the
new task uses.

## Exercise an installation

1. Run [scripts/show-template.py](scripts/show-template.py) with Python 3:

   ```sh
   python3 scripts/show-template.py
   ```

   Run this command from the skill directory, or pass the script's absolute
   path from another directory. Use `--help` for usage.
2. Confirm that stdout matches [assets/report-template.md](assets/report-template.md).
   This demonstrates that the script can find its bundled asset. It does not
   validate metadata or prove that every installed file is intact.
3. Fill the template in your response with the path checked, observed result,
   and any failure. Mark checks you did not run as unverified. Save a report
   only if requested, using the user's output location.

For example, “Check the demo installation” should produce a short report with
the actual script result. “Explain the optional folders” only needs an
explanation; it does not require running the script.

If Python is unavailable, report that the script check could not run. If a
resource is missing, report its path rather than inventing replacement content.

## Store metadata

`status`, `tags`, and `x-claude` are this repository's optional source
extensions. A skill without them is published. Installation removes `status`
and `tags`. For Claude it lifts `x-claude` fields into frontmatter; for every
other agent it drops that block. The demo is explicitly invoked in Claude
because `disable-model-invocation` is true.

The other fields demonstrate the Agent Skills specification. `allowed-tools`
pre-approves reading where supported; script execution uses the client's normal
permission checks. Validate the shared `.agents/skills` copy against the
standard and check the Claude extensions against Claude's documentation.
