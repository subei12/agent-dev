-- +goose Up
create table if not exists missions (
  id text primary key,
  project_id text not null references projects(id) on delete cascade,
  title text not null,
  description text,
  source_type text not null,
  status text not null,
  admin_agent_id text not null default '',
  team_id text,
  repo_binding_id text,
  created_by text not null,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  completed_at timestamptz
);

create table if not exists mission_members (
  id text primary key,
  mission_id text not null references missions(id) on delete cascade,
  agent_id text not null,
  role_key text not null,
  is_admin boolean not null default false,
  joined_at timestamptz not null default now()
);

create table if not exists mission_decisions (
  id text primary key,
  mission_id text not null references missions(id) on delete cascade,
  decided_by_agent_id text not null,
  decision text not null,
  summary text not null,
  related_task_item_id text,
  related_document_version_ids_json jsonb not null default '[]'::jsonb,
  related_repo_candidate_ids_json jsonb not null default '[]'::jsonb,
  created_at timestamptz not null default now()
);

create table if not exists agent_admin_policies (
  id text primary key,
  project_id text not null references projects(id) on delete cascade,
  require_plan_adoption_before_implementation boolean not null default true,
  require_test_checkpoint_before_complete boolean not null default true,
  require_review_checkpoint_before_complete boolean not null default true,
  allow_auto_create_tasks boolean not null default false,
  allow_auto_assign_tasks boolean not null default false,
  allow_auto_trigger_downstream boolean not null default false,
  allow_complete_with_open_risks boolean not null default false,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);
