-- +goose Up
create table if not exists task_boards (
  id text primary key,
  mission_id text not null references missions(id) on delete cascade,
  title text not null,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create table if not exists task_items (
  id text primary key,
  board_id text not null references task_boards(id) on delete cascade,
  title text not null,
  type text not null,
  status text not null,
  assigned_agent_id text,
  upstream_task_ids_json jsonb not null default '[]'::jsonb,
  downstream_task_ids_json jsonb not null default '[]'::jsonb,
  input_document_version_ids_json jsonb not null default '[]'::jsonb,
  input_repo_candidate_ids_json jsonb not null default '[]'::jsonb,
  definition_of_done_json jsonb not null default '{}'::jsonb,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create table if not exists task_claims (
  id text primary key,
  task_item_id text not null references task_items(id) on delete cascade,
  agent_id text not null,
  status text not null,
  claim_reason text,
  created_at timestamptz not null default now(),
  ended_at timestamptz
);

create table if not exists task_handoffs (
  id text primary key,
  task_item_id text not null references task_items(id) on delete cascade,
  from_agent_id text not null,
  to_agent_id text,
  to_admin_agent boolean not null default false,
  summary text not null,
  output_document_version_ids_json jsonb not null default '[]'::jsonb,
  output_repo_candidate_ids_json jsonb not null default '[]'::jsonb,
  output_artifact_version_ids_json jsonb not null default '[]'::jsonb,
  validation_summary text,
  risk_summary text,
  recommended_next_action text,
  status text not null,
  created_at timestamptz not null default now(),
  decided_at timestamptz
);

create table if not exists review_checkpoints (
  id text primary key,
  mission_id text not null references missions(id) on delete cascade,
  task_item_id text,
  kind text not null,
  requested_by_agent_id text not null,
  assigned_agent_id text,
  status text not null,
  summary text,
  linked_document_version_ids_json jsonb not null default '[]'::jsonb,
  linked_repo_candidate_ids_json jsonb not null default '[]'::jsonb,
  created_at timestamptz not null default now(),
  finished_at timestamptz
);
