ALTER TABLE audit_events
  ADD COLUMN target varchar(255),
  ADD COLUMN client_id varchar(255),
  ADD COLUMN token_id varchar(255);

CREATE INDEX audit_target_idx ON audit_events(target,occurred_at);
