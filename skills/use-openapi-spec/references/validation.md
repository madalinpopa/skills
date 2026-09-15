# Validation and delivery

## Select supported tools

Prefer the repository's pinned commands and configuration. Record exact versions
and verify OpenAPI 3.2 support across validators, renderers, bundlers, generators,
and any contract-test runner. A parser accepting YAML is not an OAS validator;
a 3.0-only generator is not validated by a 3.2-aware linter.

[Redocly CLI](https://github.com/Redocly/redocly-cli) is one local option with
3.2 support. With a compatible installed release, lint an explicit entry file:

```sh
redocly lint path/to/openapi.yaml
```

Use the project's rules. If no configuration exists, its `--extends spec` ruleset
provides a specification-oriented starting point; design rules are a separate
choice. Consult [lint usage](https://redocly.com/docs/cli/commands/lint) and
[rulesets](https://redocly.com/docs/cli/rules) before changing configuration.
Use existing shell/CI commands for repeated lint and bundle operations, checking
each exit status. Do not replace a failed command with a success from a later one.

If tooling is missing, report it and follow the environment's dependency policy.
Use a pinned version when adding a reproducible tool invocation; do not install
globally, alter lockfiles, or use hosted uploads merely to run a check. Never
silently fetch arbitrary remote references or execute remote examples.

## Check in layers

1. Parse YAML/JSON with duplicate-key detection and JSON-compatible values. Quote
   HTTP status keys and strings that could be misinterpreted as numbers or null.
2. Validate OAS structure, unique operation IDs, parameters, security references,
   schemas, and all local `$ref` targets from the real entry document. Official
   [OAS JSON Schemas](https://spec.openapis.org/oas/) have different scopes:
   structural schemas do not necessarily validate embedded schemas or examples.
3. Validate examples against their schemas/dialects and inspect serialized forms.
   Enable supported example checks deliberately; format assertions and non-JSON
   serialization may require separate checks. A skeleton with `paths: {}` can
   pass format validation while documenting no useful operations.
4. Run affected renderer/generator checks when available and requested. Inspect
   generated changes; source edits must not unexpectedly rename client methods
   or remove supported fields. Keep temporary output outside source directories.
5. Compare with the previous contract and relevant implementation or acceptance
   criteria. Check success, meaningful failure, anonymous/protected access, and
   omitted/null/boundary cases. Use existing contract tests when applicable;
   do not contact a live service without authorization.

## Report limits honestly

Separate specification errors, team style findings, unsupported tool features,
and observed runtime drift. Fix the cause without weakening the documented API
or disabling a rule to obtain a green result. Report commands and versions,
remaining warnings, breaking changes, and anything not checked. Keep successful
logs short; never claim runtime enforcement from an OpenAPI security block.
