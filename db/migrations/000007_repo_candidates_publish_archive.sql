-- +goose Up
create table if not exists repo_candidates (
  id text primary key,
  mission_id text not null references missions(id) on delete cascade,
  task_item_id text,
  repo_binding_id text not null,
  base_commit_sha text not null,
  parent_candidate_id text,
  created_from text not null,
  created_from_ref_id text not null,
  tree_hash text not null,
  patch_object_key text,
  workspace_archive_object_key text,
  summary text,
  status text not null,
  produced_by_run_id text,
  produced_by_node_run_id text,
  created_at timestamptz not null default now()
);

create table if not exists repo_candidate_publishes (
  id text primary key,
  candidate_id text not null references repo_candidates(id) on delete cascade,
  intent_hash text not null,
  status text not null,
  target_branch text,
  created_at timestamptz not null default now(),
  finished_at timestamptz
);

create table if not exists approvals (
  id text primary key,
  mission_id text references missions(id) on delete cascade,
  run_id text,
  node_run_id text,
  action text not null,
  subject_type text not null,
  subject_id text not null,
  intent_snapshot_json jsonb not null default '{}'::jsonb,
  intent_hash text not null,
  status text not null,
  comment text,
  created_by text not null,
  decided_by text,
  created_at timestamptz not null default now(),
  decided_at timestamptz
);

create table if not exists artifact_publishes (
  id text primary key,
  mission_id text references missions(id) on delete cascade,
  artifact_version_id text not null,
  content_hash text not null,
  base_commit_sha text,
  intent_hash text not null,
  status text not null,
  created_at timestamptz not null default now()
);

create table if not exists mission_archives (
  id text primary key,
  mission_id text not null references missions(id) on delete cascade,
  status text not null,
  manifest_object_key text,
  bundle_object_key text,
  hash text,
  created_at timestamptz not null default now(),
  finished_at timestamptz
);
