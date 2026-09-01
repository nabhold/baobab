DROP INDEX IF EXISTS audit_target_idx;
ALTER TABLE audit_events
  DROP COLUMN IF EXISTS token_id,
  DROP COLUMN IF EXISTS client_id,
  DROP COLUMN IF EXISTS target;
