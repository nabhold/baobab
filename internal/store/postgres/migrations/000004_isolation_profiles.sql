CREATE TABLE IF NOT EXISTS policy.isolation_profile (
    isolation_profile_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    name text NOT NULL,
    strategy text NOT NULL CHECK (strategy IN ('schema_per_tenant', 'row_level_security')),
    tenant_scope text NOT NULL DEFAULT 'tenant',
    data_partitioning text NOT NULL DEFAULT 'per_tenant',
    is_default boolean NOT NULL DEFAULT false,
    created_at timestamptz NOT NULL DEFAULT now(),
    UNIQUE (name)
);

CREATE TABLE IF NOT EXISTS policy.tenant_isolation_profile (
    tenant_id text NOT NULL,
    isolation_profile_id uuid NOT NULL REFERENCES policy.isolation_profile(isolation_profile_id) ON DELETE CASCADE,
    effective_from timestamptz NOT NULL DEFAULT now(),
    effective_to timestamptz,
    PRIMARY KEY (tenant_id, isolation_profile_id, effective_from)
);

CREATE INDEX IF NOT EXISTS tenant_isolation_profile_tenant_idx
    ON policy.tenant_isolation_profile(tenant_id);
