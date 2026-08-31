# Control-plane security model

The canonical trust-boundary contract is
`nabhold/shared/docs/security/control-plane-trust-boundaries.md`. This document
records the runtime obligations of `baobab-cp`.

## Authentication and authorisation

- Validate OIDC issuer, audience, expiry, not-before time, signature algorithm,
  and current JWKS key.
- Authorise every management operation using explicit scopes and tenant or
  platform policy; authentication alone never grants administration authority.
- Use workload identity between deployed services. Static shared API keys are
  not the production default.

## Data handling

- Encrypt transport across every trust boundary.
- Store only opaque secret references in control-plane metadata.
- Parameterise SQL and restrict migration/provisioning privileges by role.
- Redact authorisation headers, tokens, connection strings, credentials, and
  personal information from logs and traces.
- Keep audit events append-only under a separately controlled retention policy.

## Provisioner credentials

Each adapter uses a distinct identity with only its required operations. A
database provisioner cannot edit APISIX; a gateway reconciler cannot create
database schemas; a digital estate has neither permission.

Failures never trigger fallback to a broader identity.
