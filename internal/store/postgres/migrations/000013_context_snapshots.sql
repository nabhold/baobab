CREATE TABLE IF NOT EXISTS audit.context_snapshot (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    context_version text NOT NULL,
    context_hash text NOT NULL UNIQUE,
    tenant_id uuid,
    legal_entity_id uuid,
    market_id uuid,
    digital_estate_id uuid,
    context_data jsonb NOT NULL,
    resolved_at timestamptz NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now()
);
