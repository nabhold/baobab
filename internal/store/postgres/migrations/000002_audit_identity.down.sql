DROP INDEX IF EXISTS audit_correlation_idx;
ALTER TABLE audit_events DROP CONSTRAINT IF EXISTS audit_actor_type_check;
ALTER TABLE audit_events
  DROP COLUMN IF EXISTS policy_decision,
  DROP COLUMN IF EXISTS result,
  DROP COLUMN IF EXISTS idempotency_key,
  DROP COLUMN IF EXISTS correlation_id,
  DROP COLUMN IF EXISTS actor_type,
  DROP COLUMN IF EXISTS actor_id;
