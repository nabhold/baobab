# Context resolution API

`POST /v1/context/resolve` resolves tenant lifecycle and one product
entitlement from a verified workload identity.

```http
POST /v1/context/resolve HTTP/1.1
Authorization: Bearer <short-lived-workload-token>
Content-Type: application/json
X-Correlation-ID: 7c8f131b-d8ba-4d89-b60b-a187d3944074

{"product_id":"baobab-trade"}
```

An authorised, active tenant receives:

```json
{
  "tenant_id": "tn_01k4m7x9q2v6c8r3d5f1h0j4",
  "entity_id": "ZURIBEANS",
  "lifecycle_status": "active",
  "product_id": "baobab-trade",
  "entitled": true,
  "cache_ttl_seconds": 15,
  "resolved_at": "2026-09-01T10:00:00Z",
  "correlation_id": "7c8f131b-d8ba-4d89-b60b-a187d3944074"
}
```

The tenant is never accepted from the body or an identity header. Unknown,
inactive, unreconciled and unentitled contexts all return the same `403`
problem category. An invalid token returns `401`; an unavailable authoritative
store returns retryable `503`. Callers may cache success for 15 seconds only
and must fail closed after expiry.

Production traffic also requires infrastructure-terminated workload mTLS.
Never log the bearer token or place it in examples, traces or support tickets.

