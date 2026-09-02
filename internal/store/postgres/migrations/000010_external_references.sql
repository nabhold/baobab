CREATE TABLE IF NOT EXISTS registry.external_reference (
    external_reference_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    canonical_entity_id uuid NOT NULL REFERENCES registry.canonical_entity(canonical_entity_id) ON DELETE CASCADE,
    provider text NOT NULL,
    provider_key text NOT NULL,
    external_url text,
    created_at timestamptz NOT NULL DEFAULT now(),
    UNIQUE (canonical_entity_id, provider, provider_key)
);

CREATE INDEX IF NOT EXISTS external_reference_provider_idx
    ON registry.external_reference(provider, provider_key);
