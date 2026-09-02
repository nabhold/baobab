CREATE TABLE IF NOT EXISTS mapping.mapping_scope (
    mapping_scope_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id text,
    market_id uuid,
    entity_type text NOT NULL,
    scope_type text NOT NULL DEFAULT 'tenant',
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS mapping_scope_tenant_market_idx
    ON mapping.mapping_scope(tenant_id, market_id);
