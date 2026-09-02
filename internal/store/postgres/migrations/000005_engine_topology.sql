CREATE TABLE IF NOT EXISTS topology.engine (
    engine_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    code text NOT NULL UNIQUE,
    name text NOT NULL,
    description text,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS topology.engine_instance (
    engine_instance_id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    engine_id uuid NOT NULL REFERENCES topology.engine(engine_id) ON DELETE CASCADE,
    region text NOT NULL,
    environment text NOT NULL,
    status text NOT NULL DEFAULT 'active',
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS engine_instance_engine_idx
    ON topology.engine_instance(engine_id);
