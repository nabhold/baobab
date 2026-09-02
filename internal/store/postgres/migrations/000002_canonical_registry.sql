CREATE TABLE IF NOT EXISTS registry.canonical_entity (
    canonical_entity_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id text,
    legal_entity_id text,
    entity_type text NOT NULL,
    external_key text,
    status text NOT NULL DEFAULT 'active',
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS canonical_entity_status_idx
    ON registry.canonical_entity(status);

CREATE INDEX IF NOT EXISTS canonical_entity_tenant_idx
    ON registry.canonical_entity(tenant_id);
