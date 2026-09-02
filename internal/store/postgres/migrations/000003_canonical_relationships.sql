CREATE TABLE IF NOT EXISTS registry.canonical_relationship (
    relationship_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    parent_entity_id uuid NOT NULL REFERENCES registry.canonical_entity(canonical_entity_id) ON DELETE CASCADE,
    child_entity_id uuid NOT NULL REFERENCES registry.canonical_entity(canonical_entity_id) ON DELETE CASCADE,
    relationship_type text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    UNIQUE (parent_entity_id, child_entity_id, relationship_type)
);

CREATE INDEX IF NOT EXISTS canonical_relationship_parent_idx
    ON registry.canonical_relationship(parent_entity_id);

CREATE INDEX IF NOT EXISTS canonical_relationship_child_idx
    ON registry.canonical_relationship(child_entity_id);
