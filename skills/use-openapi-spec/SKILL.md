---
name: use-openapi-spec
description: Creates, updates, and reviews OpenAPI 3.2 HTTP API descriptions with explicit contracts, reusable schemas, security requirements, and tool-aware validation. Use for OpenAPI YAML or JSON authoring and migration; not for implementing handlers or documenting non-HTTP protocols.
status: published
tags: [openapi, api, specification, http]
---

# Use OpenAPI specification

Author the API's observable contract using [OpenAPI 3.2.0](https://swagger.io/specification/v3.2/).
Keep ordinary changes focused; load only the reference sections needed below.
An OpenAPI description documents behavior; it does not implement or enforce it.

## Establish the contract

1. Inspect repository instructions, existing descriptions, relevant routes or
   requirements, and validation/generation commands. Identify the source of truth;
   edit that source rather than generated files. An advice/review request does
   not authorize edits. Preserve unrelated operations and established conventions.
2. Determine consumers, operations, payloads, authentication, and expected errors
   from available evidence. Ask only about missing decisions that change the
   contract. Mark unresolved behavior; do not invent endpoints, scopes, limits,
   retries, or security guarantees to complete a template.
3. Use `openapi: 3.2.0` for new descriptions. `info.version` is the description's
   own version, not the OAS version. Check the actual validator, documentation
   renderer, and generator versions before migration or new feature use. Preserve
   an existing document's version unless migration is in scope. Report a tooling
   blocker rather than silently downgrading or claiming unsupported checks passed.

## Write

For a new document, copy [the starter](assets/openapi.yaml) with an existing
shell tool after checking the destination; never overwrite an existing file.
Resolve asset paths from this installed skill and output paths from the user's
project. Replace the sample metadata and empty paths with the agreed contract.
Reuse repository scaffolding and lint/bundle commands for repeated work.

For each operation, define a stable unique `operationId`, concise purpose,
parameters and their serialization, request media/schema, success responses,
meaningful errors, and applicable security. Match path placeholders to required
path parameters. Use quoted response-code keys and accurate media types.
Document only behavior the API promises; distinguish normative OAS rules from
team design conventions.

Model presence separately from nullability. State constraints and examples that
match the data, and reuse components when they represent the same contract.
Make anonymous access an explicit decision; declaring a security scheme alone
protects no operation. Keep credentials and personal data out of examples.

| When needed | Read |
| --- | --- |
| Endpoints, inputs, responses, pagination, or errors | [HTTP contracts](references/contracts.md) |
| Types, nullability, composition, examples, or references | [Schemas](references/schemas.md) |
| Authentication or operation access | [Security](references/security.md) |
| Migration, streaming, QUERY, or other 3.2 additions | [Version and feature choices](references/openapi-3.2.md) |
| Validator selection, compatibility, and final checks | [Validation](references/validation.md) |

## Verify and report

Run the project's applicable checks with confirmed 3.2 support. Validate syntax,
reference resolution, examples, and changed operation contracts; use the
validation reference for gaps or new tooling. For existing APIs, compare changes
with the previous contract and report consumer-breaking differences.

Inspect the final diff. Report the artifact, OAS version, important decisions,
checks and tool versions, and remaining gaps. Separate lint success from generator
compatibility and runtime conformance. Do not publish, call live APIs, or modify
implementation code unless the user requested that work.
