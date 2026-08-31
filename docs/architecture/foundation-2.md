# Foundation 2: executable control-plane core

Foundation 2 replaces the legacy polyglot monolith scaffold with a focused Go service. The first vertical slice implements the canonical `POST /v1/tenants` contract from [`nabhold/shared`](https://github.com/nabhold/shared), PostgreSQL desired state, idempotent operations, an append-only audit record, and a transactional outbox record.

Infrastructure manifests remain in [`nabhold/infrastructure`](https://github.com/nabhold/infrastructure). Cross-repository schemas remain in `nabhold/shared`. This repository owns only control-plane runtime logic and its own database migrations.

The administrative bearer token is a bootstrap authentication boundary. OIDC validation and workload mTLS are Foundation 3 security work; production deployment must not proceed with the bootstrap token.
