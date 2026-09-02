CREATE TABLE IF NOT EXISTS audit.audit_record (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id uuid,
    actor_id uuid,
    action text NOT NULL,
    resource_type text NOT NULL,
    resource_id uuid NOT NULL,
    previous_version jsonb,
    resulting_version jsonb,
    reason text,
    request_id uuid,
    correlation_id uuid,
    source_ip_hash text,
    occurred_at timestamptz NOT NULL DEFAULT now(),
    metadata jsonb NOT NULL DEFAULT '{}'::jsonb
);
