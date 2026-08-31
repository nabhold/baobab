# Baobab Control Plane (`baobab-cp`)

> The authoritative source of tenant lifecycle, entitlement, and desired-state truth for the Baobab ecosystem.

**Status:** Foundation 2 — executable Go control-plane core under active development (see [ADR-0001](docs/adr/0001-go-control-plane-runtime.md)).
**Architecture:** [ADR-0003 — Multi-Tenant, Production-Ready Control Plane Architecture](docs/adr/0003-multi-tenant-control-plane-architecture.md)

---

## What this repository is

`baobab-cp` decides *who a tenant is, what state they're in, and what they're entitled to use* — and reconciles that decision against reality. It does not process commerce, ERP, or research-intelligence business logic; those live in their own product engines and consume this repository's decisions over the network.

If you are looking for:
- **Commerce logic** → [`nabhold/baobab-trade`](https://github.com/nabhold/baobab-trade)
- **ERP logic** → [`nabhold/baobab-erp`](https://github.com/nabhold/baobab-erp)
- **Research intelligence** → [`nabhold/baobab-pulse`](https://github.com/nabhold/baobab-pulse)
- **Canonical contracts** (schemas this repo implements against) → [`nabhold/shared`](https://github.com/nabhold/shared)
- **Infrastructure provisioning** (Terraform, APISIX bootstrap, RabbitMQ/Postgres/Redis topology) → [`nabhold/infrastructure`](https://github.com/nabhold/infrastructure)
- **Local dev container image** → [`nabhold/baobab-dev`](https://github.com/nabhold/baobab-dev)

...you want one of those repositories instead. This one is intentionally narrow.

## Responsibilities

| Owns | Does not own |
|---|---|
| Tenant lifecycle (provision, suspend, reinstate, decommission) | Business data of any kind |
| Tenancy hierarchy metadata (Tenant Group → Tenant → Business Unit → Function → Team) | UI / end-user surfaces |
| Product entitlements per tenant | Infrastructure provisioning mechanics (that's `nabhold/infrastructure`) |
| Desired-state reconciliation (APISIX routes, per-tenant Postgres boundaries) | Canonical contract *definitions* (that's `nabhold/shared` — this repo implements against them) |
| Auditable provisioning history | Product-specific integrations |
| Lifecycle/entitlement event publication | — |
| Authenticated tenant-context resolution for product engines | — |

## Architecture at a glance

```
                     ┌────────────────────────────┐
                     │        nabhold/shared        │
                     │  canonical contracts (OpenAPI,│
                     │  AsyncAPI, JSON Schema)        │
                     └───────────────┬────────────┘
                                     │ implements
┌────────────────────────────────────▼────────────────────────────────────┐
│                              baobab-cp                                   │
│                                                                          │
│   API layer (REST)         Reconciler loop         Event outbox         │
│   /v1/context/resolve      desired vs actual        → RabbitMQ          │
│   /v1/tenants              → APISIX / Postgres                          │
│   /v1/entitlements                                                      │
│                                                                          │
│                    PostgreSQL (authoritative)                           │
└───────────┬──────────────────────────────────────────┬──────────────────┘
            │ resolves context for                     │ provisioned by
            ▼                                            ▼
┌───────────────────────────┐                 ┌────────────────────────────┐
│ baobab-trade / baobab-erp  │                 │   nabhold/infrastructure    │
│ baobab-pulse (consumers)   │                 │  Terraform · APISIX · RMQ   │
└───────────────────────────┘                 └────────────────────────────┘
```

See [ADR-0003](docs/adr/0003-multi-tenant-control-plane-architecture.md) for the full rationale, including why REST over gRPC for v1, why fail-closed context resolution is a contract not an implementation detail, and why a bespoke reconciler rather than a full Kubernetes operator at this stage.

## Tech stack

| Concern | Choice | Why |
|---|---|---|
| Language / runtime | Go | [ADR-0001](docs/adr/0001-go-control-plane-runtime.md) |
| HTTP | `net/http` + `chi` router | minimal, idiomatic, no framework lock-in |
| Database | PostgreSQL 17 via `pgx` | authoritative store; matches `nabhold/infrastructure`'s provisioned topology |
| Migrations | `golang-migrate` | plain SQL, reviewable diffs |
| Messaging | RabbitMQ via `amqp091-go` | matches provisioned topology; DLQ support |
| Gateway integration | APISIX Admin API client | control plane reconciles routes it owns |
| AuthN | mTLS (service-to-service) + OIDC (admin API) | see ADR-0003 §7 |
| Observability | OpenTelemetry (traces, metrics, logs) | org-wide observability contract (`nabhold/shared`) |
| Config | environment variables, validated at startup | 12-factor, container-friendly |
| Testing | standard `testing` + Testcontainers (Postgres, RabbitMQ) | real dependencies in CI, not mocks-only |
| CI/CD | reusable workflows from `nabhold/shared` | org-wide standardisation |
| Container | multi-stage Dockerfile, distroless final stage | minimal attack surface |

## Repository structure

```
baobab-cp/
├── cmd/
│   └── controlplane/          # main entrypoint
├── internal/
│   ├── api/                   # HTTP layer: router, handlers, middleware
│   │   ├── handlers/
│   │   └── middleware/
│   ├── domain/                 # core domain model, framework-free
│   │   ├── tenant/
│   │   └── entitlement/
│   ├── reconcile/               # desired-state reconciliation loop
│   ├── events/                  # outbox + RabbitMQ publisher
│   ├── store/
│   │   ├── postgres/            # repository implementations
│   │   └── migrations/          # SQL migrations
│   ├── gateway/                 # APISIX admin API client
│   └── config/                  # startup configuration + validation
├── pkg/
│   └── contracts/                # generated clients from nabhold/shared
├── docs/
│   ├── adr/                      # this repo's local ADR register
│   └── architecture/
├── test/
│   └── integration/
├── .github/workflows/            # calls nabhold/shared reusable workflows
├── Dockerfile
├── Makefile
├── go.mod
├── .env.example
└── README.md
```

`internal/domain` contains no framework or infrastructure imports — it is the part of this codebase that should be easiest to test and hardest to accidentally couple to a specific database or transport.

## Getting started

```bash
git clone https://github.com/nabhold/baobab-cp.git
cd baobab-cp
cp .env.example .env
make dev-up      # starts local Postgres + RabbitMQ via nabhold/infrastructure compose topology
make migrate
make run
```

Run tests:

```bash
make test              # unit tests
make test-integration  # Testcontainers-backed integration tests
```

## Relationship with other repositories

| Repository | Relationship |
|---|---|
| `nabhold/shared` | Contract source of truth. `baobab-cp` implements the OpenAPI/AsyncAPI schemas defined there; it never redefines them locally. |
| `nabhold/infrastructure` | Provisions the Postgres, RabbitMQ, and APISIX instances `baobab-cp` depends on and reconciles against. `baobab-cp` never provisions its own infrastructure. |
| `nabhold/baobab-trade` | Calls `POST /v1/context/resolve` on every commerce request; fails closed if unresolved. Consumer, not a dependency of this repo. |
| `nabhold/baobab-erp` | Contract-level consumer of canonical identifiers; does not yet call this repo's context-resolution API directly (open question — see ADR-0003 §14). |
| `nabhold/baobab-pulse` | Consumer of tenant/entitlement context (integration not yet established — repository is pre-Foundation). |
| `nabhold/baobab-dev` | Provides this repository's local/CI development container image. |
| Digital-estate frontends (`nabhold/nabhold`, `zuribeans`, `thamani`, `equator-estate`) | Consume Baobab exclusively through product-engine APIs; do not call `baobab-cp` directly in the general case. |

## Security

- Every request to `/v1/context/resolve` is authenticated via mutual TLS + workload identity; responses are signed.
- Secrets are never committed. `.env` is for local development only and must never contain production credentials — see [SECURITY.md](SECURITY.md).
- Every provisioning and lifecycle-transition action is written to an append-only audit log before being acknowledged.

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md). Changes to `/v1/context/resolve`'s contract (including its fail-closed behaviour) require an ADR, not just a PR — it is depended upon by production consumers.

## License

Apache-2.0. See [LICENSE](LICENSE).
