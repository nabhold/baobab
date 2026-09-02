CREATE TABLE IF NOT EXISTS mapping.canonical_mapping (
    canonical_mapping_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    source_entity_id uuid NOT NULL REFERENCES registry.canonical_entity(canonical_entity_id) ON DELETE CASCADE,
    target_entity_id uuid NOT NULL REFERENCES registry.canonical_entity(canonical_entity_id) ON DELETE CASCADE,
    mapping_type text NOT NULL,
    status text NOT NULL DEFAULT 'active',
    created_at timestamptz NOT NULL DEFAULT now(),
    UNIQUE (source_entity_id, target_entity_id, mapping_type)
);

CREATE INDEX IF NOT EXISTS canonical_mapping_status_idx
    ON mapping.canonical_mapping(status);
