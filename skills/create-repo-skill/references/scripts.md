# Local authoring scripts

## Runtime and use

Use Bash for file creation and command orchestration. The YAML operations use
Python 3 and PyYAML to preserve structured data reliably. Helpers must report
missing dependencies before changing files; they do not install packages.
An agent with `uv` can supply the dependency for one command with
`uv run --with pyyaml bash <script> ...`, or use an existing Python environment.

Run helpers by absolute installed path. Project/output paths are explicit;
never derive the authoring checkout from the installed skill's location or
configured store. Quote paths and data, including values with spaces, dollar
signs, quotes, and newlines. Pass subprocess arguments directly; do not use eval.

## Initialize

```sh
scripts/init-skill.sh <name> --repo <checkout> --description "<capability and trigger>" [--status published|draft] [--resources scripts,references,assets]
```

Creates only `skills/<name>/SKILL.md` and requested resource directories inside
the specified repository. Require its `AGENTS.md`, `docs/SPEC.md`, and `skills/`.
Default to explicit `published`; `--status draft` is opt-in. Validate names,
options, and output paths before writes. Refuse an existing skill or symlink;
do not overwrite, follow an escaping path, or silently normalize a supplied name.

Do not create UI metadata or example files by default. Add meaningful examples
when the real workflow needs them. Finish scaffold placeholders before validating
or publishing; a published status does not make unfinished content acceptable.

## Validate

```sh
scripts/validate-skill.sh <skill-folder>
```

Validate source metadata, naming, and unfinished scaffold text without rewriting
source files. Accept supported optional and unknown fields; enforce repository
store types, unique YAML keys, and `x-claude` exclusions/collisions. The review
still checks links, script behavior, and realistic requests separately.

## Write agent metadata

```sh
scripts/write-agent-metadata.sh <skill-folder> --interface 'display_name=Example' --interface 'short_description=Create useful repository examples'
```

Use repeatable `--interface key=value` for supported UI fields. Read the local
agent-metadata reference for constraints. Update only supplied fields, preserving
existing interface fields, policy, dependencies, and other supported data.
Validate values before creating directories or changing a file. Add optional
fields only when provided or requested; do not infer tool dependencies or policy.

## Design other helpers

Prefer executable Bash helpers for repeated scaffold creation, asset copying,
validated substitutions, file inventories, and verification chains. Keep them
portable to the target environment, including macOS Bash 3.2 when applicable.
Use `set -euo pipefail`, explicit arguments, meaningful exit codes, and short
summaries. Run deterministic automation without printing full source into context.

Parse options and check dependencies, inputs, destinations, and collisions before
writes. Preserve user files. Use safe temporary files for replacements and clean
up only temporary content owned by the helper. Report partial work if a later
step fails; do not claim full rollback unless implemented and tested.

Use Python or another established tool for structured formats and logic that
would be brittle in shell. Bash should reduce agent calls and repetitive file
work, not replace proper parsers or hide complex code in enormous heredocs.
Do not add a helper for a one-off trivial edit that ordinary tools already handle.
