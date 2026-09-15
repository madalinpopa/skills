# OpenAPI 3.2 choices

Target `openapi: 3.2.0`, the version on the user's selected
[specification page](https://swagger.io/specification/v3.2/). Check newer patch
clarifications when needed; do not automatically change the feature family.

## Migration

Inspect the existing version and consumers before editing. Follow the official
[3.1 to 3.2 guide](https://learn.openapis.org/upgrading/v3.1-to-v3.2.html), then
validate with the exact consuming tools. For 3.0 sources, first account for the
[3.0 to 3.1 schema changes](https://learn.openapis.org/upgrading/v3.0-to-v3.1.html),
including nullability and JSON Schema keyword semantics. Changing the version
string alone is not evidence of compatibility.

Keep migration separate from changing endpoint behavior. If a generator only
accepts older versions, report the exact limitation. Do not relabel a document,
discard keywords, or introduce a second derived contract without an agreed plan.

## Adopt features only for actual requirements

| Need | 3.2 feature and check |
| --- | --- |
| Complex safe queries with a body | Path Item `query` describes QUERY; verify actual server/proxy/client support. |
| Other HTTP methods | `additionalOperations` keys use wire-method capitalization; do not duplicate a fixed method such as POST there. |
| Whole-query encoding | `in: querystring` uses `content`; only one is allowed and it cannot coexist with `in: query` parameters on the operation/path item. |
| Streaming sequences | Media Type `itemSchema` describes individual items. Match the actual sequential media type and framing; do not model an unbounded stream as an ordinary JSON array. |
| Richer examples | `dataValue` represents data; `serializedValue` represents wire text. Follow mutual-exclusion and validation rules. |
| Navigation hierarchy | Tag `summary`, `parent`, and `kind` organize existing operations; verify renderer support. |
| Document identity | Root `$self` changes the reference base. Add only with a deliberate URI/resolution strategy. |
| OAuth device flow | `deviceAuthorization` needs the specified authorization/token endpoints and scopes. |

For server-sent events, `itemSchema` models parsed event fields, not just the
application payload. The `data` field is a string; describe embedded JSON using
`contentMediaType`/`contentSchema` when applicable. Verify framing and field types
against the specification's sequential-media rules.

Read the specific normative object before using a less-common feature, including
multipart encodings, XML `nodeType`, discriminator `defaultMapping`, or external
security-scheme URIs. Keep these details off an ordinary CRUD authoring path.
