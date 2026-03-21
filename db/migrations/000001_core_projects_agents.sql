-- +goose Up
create table if not exists projects (
  id text primary key,
  name text not null,
  description text,
  status text not null,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create table if not exists repo_bindings (
  id text primary key,
  project_id text not null references projects(id) on delete cascade,
  provider text not null,
  remote_url text not null,
  default_branch text not null,
  auth_mode text not null,
  write_enabled boolean not null default false,
  created_at timestamptz not null default now()
);

create table if not exists executor_profiles (
  id text primary key,
  project_id text not null references projects(id) on delete cascade,
  name text not null,
  type text not null,
  command text not null,
  args_json jsonb not null default '[]'::jsonb,
  env_secret_refs_json jsonb not null default '[]'::jsonb,
  timeout_sec integer not null default 3600,
  max_concurrency integer not null default 1,
  max_prompt_chars integer not null default 16000,
  allow_repo_read boolean not null default true,
  allow_repo_write boolean not null default false,
  allow_network boolean not null default false,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create table if not exists roles (
  id text primary key,
  project_id text not null references projects(id) on delete cascade,
  key text not null,
  name text not null,
  prompts_json jsonb not null default '{}'::jsonb,
  permissions_json jsonb not null default '{}'::jsonb,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  unique(project_id, key)
);

create table if not exists agents (
  id text primary key,
  project_id text not null references projects(id) on delete cascade,
  name text not null,
  role_id text not null references roles(id),
  executor_profile_id text not null references executor_profiles(id),
  enabled boolean not null default true,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create table if not exists agent_teams (
  id text primary key,
  project_id text not null references projects(id) on delete cascade,
  name text not null,
  description text,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create table if not exists agent_team_members (
  id text primary key,
  team_id text not null references agent_teams(id) on delete cascade,
  agent_id text not null references agents(id),
  role_key text not null,
  capabilities_json jsonb not null default '{}'::jsonb,
  created_at timestamptz not null default now()
);
