-- +goose Up
alter table task_claims
  add column if not exists execution_lock_token text,
  add column if not exists attempt_count integer not null default 0,
  add column if not exists last_error text,
  add column if not exists last_heartbeat_at timestamptz;
