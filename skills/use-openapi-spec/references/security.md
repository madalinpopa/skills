# Security contracts

Follow the [Security Scheme](https://swagger.io/specification/v3.2/#security-scheme-object)
and [Security Requirement](https://swagger.io/specification/v3.2/#security-requirement-object)
objects. Describe the actual mechanism and access policy; the application must
implement authentication and authorization separately.

## Declare and apply

Define reusable schemes in `components.securitySchemes`, then apply them through
root or operation `security`. Declaring a scheme without a requirement enables
no protection. Use HTTP bearer for bearer tokens, OAuth2/OpenID Connect for those
protocols, and the correct name/location for an API key. `bearerFormat: JWT` is
an informational hint, not token verification.

```yaml
components:
  securitySchemes:
    bearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT
security:
  - bearerAuth: []
```

An empty array after a scheme name does not mean public access. It leaves the
scheme requirement in place without listed scopes/roles. OAuth/OpenID arrays
name required scopes; other schemes can list roles in 3.2 when the contract and
tools support them. Do not invent a role or scope model.

## Review effective access per operation

- An operation inherits root security when its own `security` is absent.
- Operation `security` replaces the root requirement; it does not merge with it.
- `security: []` explicitly removes inherited requirements for that operation.
- Separate objects in the array mean alternatives (OR); schemes within one
  object are all required (AND).
- An empty object `{}` among alternatives permits anonymous access. Use it only
  when optional authentication is intentional.

Keep public exceptions explicit and review them alongside protected operations.
Document the actual authentication/permission failure responses and relevant
challenge headers. Object ownership, tenant isolation, and field-level access
usually need clear prose and runtime checks beyond security-scheme declarations.
Use placeholder credentials and approved server URLs; never upload private API
descriptions or contact identity providers merely to validate syntax.
