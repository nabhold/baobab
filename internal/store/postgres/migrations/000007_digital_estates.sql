CREATE TABLE IF NOT EXISTS estate.digital_estate (
    digital_estate_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id text NOT NULL,
    name text NOT NULL,
    domain text NOT NULL,
    status text NOT NULL DEFAULT 'active',
    created_at timestamptz NOT NULL DEFAULT now(),
    UNIQUE (domain)
);

CREATE INDEX IF NOT EXISTS digital_estate_tenant_idx
    ON estate.digital_estate(tenant_id);
