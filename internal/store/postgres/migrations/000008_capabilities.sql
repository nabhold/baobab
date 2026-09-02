CREATE TABLE IF NOT EXISTS capability.capability (
    capability_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    code text NOT NULL UNIQUE,
    name text NOT NULL,
    description text,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS capability.tenant_capability (
    tenant_id text NOT NULL,
    capability_id uuid NOT NULL REFERENCES capability.capability(capability_id) ON DELETE CASCADE,
    status text NOT NULL DEFAULT 'enabled',
    granted_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (tenant_id, capability_id)
);

CREATE INDEX IF NOT EXISTS tenant_capability_tenant_idx
    ON capability.tenant_capability(tenant_id);
