# HTTP contracts

Read only affected topics. Use the [OAS operation and parameter rules](https://swagger.io/specification/v3.2/#operation-object)
as the format authority; resource naming and pagination choices are API design
conventions, not OAS requirements.

## Operations and inputs

- Match actual method semantics and route names. Keep operation IDs unique and
  stable because generators and links may depend on them. Reuse existing tags.
- Every `{placeholder}` needs a matching `in: path` parameter with
  `required: true`. Identify parameters by both name and location; do not repeat
  the same pair. Define limits, defaults, units, and serialization only when known.
- Use either parameter `schema` or `content`; parameter `content` has one media
  type. For ordinary query arrays/objects, choose supported `style`/`explode`
  deliberately and show a serialized example when ambiguity matters. Do not
  assume nested `deepObject` structures have portable encoding.
- Request-body presence and property presence are separate: set
  `requestBody.required` for a mandatory body and schema `required` for mandatory
  fields. Describe PATCH semantics, including omitted fields and explicit nulls.
- Avoid introducing GET/HEAD request bodies. For existing unusual behavior,
  describe the interoperability limits rather than silently redesigning it.

## Responses and lifecycle

Describe actual successful statuses and foreseeable failure responses. Include
body schemas, content types, and meaningful headers; omit response content when
no body is sent. Do not assign every endpoint the same generic success/error map.
Use `default` for otherwise unspecified responses, not to hide known contracts.

For collection endpoints, specify the implemented cursor/offset rules, bounds,
ordering, continuation indicators, and empty results. For asynchronous work,
state how clients observe completion. Include idempotency, rate-limit, caching,
conditional-request, and retry details only where the service promises them.

For provider-initiated requests, use an operation's `callbacks` when they depend
on an initiating API call, or root `webhooks` for independently registered events.
Describe the receiver's request contract and responses; do not invent delivery
guarantees, signature schemes, or retry schedules.

Use a consistent error shape. Preserve an established format; for new APIs,
consider [RFC 9457 Problem Details](https://www.rfc-editor.org/rfc/rfc9457.html)
with `application/problem+json`. Its `type` identifies the problem, `title`
summarizes it, and optional `detail`/`instance` describe an occurrence. If `status`
is included, keep it consistent with the actual HTTP status. Do not expose stack
traces, secrets, or internal diagnostics as public error contracts.

## Compatibility review

Assess changes from the consumer's perspective: removed/renamed operations or
fields, new required inputs, narrowed accepted values, changed serialization,
response types/statuses, and stronger authentication can break clients. Adding
response enum values can also break exhaustive clients. A new `info.version`
does not make a change safe; report the specific affected consumers/contracts.
