CREATE TABLE IF NOT EXISTS system.idempotency_record (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    idempotency_key text NOT NULL,
    tenant_id uuid,
    operation text NOT NULL,
    request_hash text NOT NULL,
    response_status integer,
    response_headers jsonb,
    response_body jsonb,
    resource_id uuid,
    state text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    completed_at timestamptz,
    expires_at timestamptz NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idempotency_operation_key_uq
    ON system.idempotency_record(
        coalesce(tenant_id, '00000000-0000-0000-0000-000000000000'::uuid),
        operation,
        idempotency_key
    );

CREATE TABLE IF NOT EXISTS system.registry_revision (
    registry_name text PRIMARY KEY,
    revision bigint NOT NULL DEFAULT 0,
    updated_at timestamptz NOT NULL DEFAULT now()
);

INSERT INTO system.registry_revision (registry_name, revision, updated_at)
VALUES
    ('canonical', 0, now()),
    ('mapping', 0, now()),
    ('market', 0, now()),
    ('estate', 0, now()),
    ('topology', 0, now()),
    ('capability', 0, now()),
    ('isolation', 0, now())
ON CONFLICT (registry_name) DO NOTHING;
