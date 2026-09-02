CREATE TABLE IF NOT EXISTS market.market (
    market_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    code text NOT NULL UNIQUE,
    name text NOT NULL,
    currency text NOT NULL,
    region text NOT NULL,
    is_active boolean NOT NULL DEFAULT true,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS market.market_assignment (
    market_assignment_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id text NOT NULL,
    market_id uuid NOT NULL REFERENCES market.market(market_id) ON DELETE CASCADE,
    effective_from timestamptz NOT NULL DEFAULT now(),
    effective_to timestamptz,
    UNIQUE (tenant_id, market_id, effective_from)
);

CREATE INDEX IF NOT EXISTS market_assignment_tenant_idx
    ON market.market_assignment(tenant_id);
