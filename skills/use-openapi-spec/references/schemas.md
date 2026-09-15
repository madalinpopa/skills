# Schemas and examples

Use the [OAS Schema Object](https://swagger.io/specification/v3.2/#schema-object)
and [JSON Schema object guidance](https://json-schema.org/understanding-json-schema/reference/object).
OAS 3.2 uses JSON Schema Draft 2020-12 semantics with the OAS dialect.

## Model the wire data

- `required` controls property presence; a property can be optional and non-null,
  or required and nullable. Use `type: [string, "null"]` for a nullable string;
  do not carry the OpenAPI 3.0 `nullable` keyword into a 3.2 schema.
- Give arrays an `items` schema and maps an appropriate `additionalProperties`
  schema. Set bounds, patterns, formats, and enums only for real constraints.
  `format` validation varies by tool; it is not proof of runtime checking.
- Omitted `additionalProperties` permits extra keys. Close an object only when
  that matches the API. With composition, consider `unevaluatedProperties`
  deliberately and confirm consumer support; `allOf` is intersection, not
  object-oriented inheritance or an override mechanism.
- Use `oneOf` for exactly one matching alternative, `anyOf` for one or more.
  Make variants distinguishable, for example by a required property with `const`.
  A discriminator helps selection but does not replace schema validation.
- Model request/response differences explicitly. Use separate schemas when
  they simplify required fields or prevent writable server-owned properties;
  `readOnly`/`writeOnly` annotations alone do not enforce access control.
  Defaults describe expected behavior; validators need not fill missing values.

```yaml
type: object
required: [id, displayName]
properties:
  id:
    type: string
    format: uuid
  displayName:
    type: [string, "null"]
  labels:
    type: array
    items:
      type: string
```

This example requires `displayName` even when null; `labels` may be absent.
It makes no promise to reject unknown properties.

## References and examples

Keep a small API in one file; split only to aid maintenance. Use stable component
names, correct JSON Pointer escaping (`~0`, `~1`), and resolvable relative paths.
Prefer local references for reproducible checks. When splitting documents, use
OpenAPI or Schema Objects at document roots; arbitrary fragments depend on tools.

Schema `$ref` siblings follow JSON Schema semantics. Other Reference Objects
have their own supported fields; do not attach arbitrary schema keywords to a
response/parameter reference. Check resolution against the full document, including
`$id` or `$self` base URI changes; bundling is not equivalent to full dereferencing.

Prefer schema `examples` arrays for data examples. Media type `examples` is a
map of named Example Objects; `example` and `examples` are mutually exclusive
there. For JSON Example Objects, `value` remains valid; 3.2 `dataValue` can make
data intent explicit. Use `serializedValue` for a non-JSON wire representation,
not as a replacement for validating the underlying data. Respect the Example
Object's mutually exclusive fields.

Check examples against schemas and serialization, including boundary values,
empty collections, and meaningful errors. Keep only examples that explain the
contract. Never include live tokens, customer records, or invented constraints.
