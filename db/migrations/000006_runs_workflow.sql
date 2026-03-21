-- +goose Up
create table if not exists runs (
  id text primary key,
  mission_id text not null references missions(id) on delete cascade,
  task_item_id text,
  status text not null,
  inputs_json jsonb not null default '{}'::jsonb,
  created_at timestamptz not null default now(),
  started_at timestamptz not null default now(),
  finished_at timestamptz
);

create table if not exists node_runs (
  id text primary key,
  run_id text not null references runs(id) on delete cascade,
  node_key text not null,
  status text not null,
  executor_session_id text,
  created_at timestamptz not null default now(),
  started_at timestamptz not null default now(),
  finished_at timestamptz
);

create table if not exists repo_snapshots (
  id text primary key,
  mission_id text not null references missions(id) on delete cascade,
  repo_binding_id text,
  commit_sha text not null,
  ref_name text,
  created_at timestamptz not null default now()
);

create table if not exists run_artifact_inputs (
  id text primary key,
  run_id text not null references runs(id) on delete cascade,
  artifact_id text not null,
  artifact_version_id text not null,
  alias text,
  created_at timestamptz not null default now()
);
