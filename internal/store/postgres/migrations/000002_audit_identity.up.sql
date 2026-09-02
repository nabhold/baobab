ALTER TABLE audit_events
  ADD COLUMN actor_id varchar(255),
  ADD COLUMN actor_type varchar(16),
  ADD COLUMN correlation_id uuid,
  ADD COLUMN idempotency_key varchar(128),
  ADD COLUMN result varchar(32),
  ADD COLUMN policy_decision varchar(64);

ALTER TABLE audit_events
  ADD CONSTRAINT audit_actor_type_check CHECK (actor_type IS NULL OR actor_type IN ('human','workload'));

CREATE INDEX audit_correlation_idx ON audit_events(correlation_id) WHERE correlation_id IS NOT NULL;
