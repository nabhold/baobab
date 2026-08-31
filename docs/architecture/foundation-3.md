# Foundation 3: identity and security boundary

Foundation 3 removes the bootstrap administrative secret. Human management
requests now require an asymmetrically signed OIDC token with the configured
issuer and `baobab-control-plane` audience, a `human` actor type, and the
explicit `tenant:write` scope.

The reusable verifier also accepts the Shared workload-identity claim shape for
future `context:resolve` use. Workload mTLS terminates at the APISIX boundary
owned by `nabhold/infrastructure`; a valid token remains mandatory after that
network check, and identity headers are never authoritative.

Every accepted registration, compatible replay, and conflicting idempotent
replay records authenticated actor and correlation provenance. The Context
Resolution contract is now defined in `nabhold/shared`, but the endpoint is not
implemented in this foundation; persistence and entitlement semantics land in
the next vertical slice.
