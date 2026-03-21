-- +goose Up
create table if not exists agent_presences (
  agent_id text primary key,
  availability text not null,
  current_mission_id text,
  current_task_item_id text,
  current_executor_session_id text,
  last_heartbeat_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create table if not exists mission_agent_runtimes (
  id text primary key,
  mission_id text not null references missions(id) on delete cascade,
  agent_id text not null,
  status text not null,
  current_task_item_id text,
  current_run_id text,
  current_node_run_id text,
  current_executor_session_id text,
  status_summary text,
  started_at timestamptz,
  updated_at timestamptz not null default now(),
  unique(mission_id, agent_id)
);

create table if not exists executor_sessions (
  id text primary key,
  mission_id text not null references missions(id) on delete cascade,
  task_item_id text,
  agent_id text not null,
  executor_profile_id text not null,
  backend text not null,
  status text not null,
  started_at timestamptz not null default now(),
  ended_at timestamptz
);

create table if not exists agent_runtime_events (
  id text primary key,
  mission_id text not null references missions(id) on delete cascade,
  agent_id text not null,
  executor_session_id text not null references executor_sessions(id) on delete cascade,
  run_id text,
  node_run_id text,
  task_item_id text,
  level text not null,
  category text not null,
  type text not null,
  title text not null,
  summary text,
  payload_json jsonb not null default '{}'::jsonb,
  occurred_at timestamptz not null default now()
);

create table if not exists executor_transcripts (
  id text primary key,
  executor_session_id text not null unique references executor_sessions(id) on delete cascade,
  mission_id text not null references missions(id) on delete cascade,
  agent_id text not null,
  storage_kind text not null,
  redacted_object_key text,
  raw_encrypted_object_key text,
  status text not null,
  started_at timestamptz not null default now(),
  ended_at timestamptz
);

create table if not exists executor_transcript_entries (
  id text primary key,
  transcript_id text not null references executor_transcripts(id) on delete cascade,
  seq integer not null,
  role text not null,
  entry_type text not null,
  redacted_text text,
  content_json jsonb not null default '{}'::jsonb,
  raw_object_key text,
  created_at timestamptz not null default now(),
  unique(transcript_id, seq)
);

create table if not exists transcript_access_audits (
  id text primary key,
  transcript_id text not null references executor_transcripts(id) on delete cascade,
  actor_user_id text not null,
  access_mode text not null,
  reason text,
  created_at timestamptz not null default now()
);

create table if not exists runtime_observability_policies (
  id text primary key,
  project_id text not null references projects(id) on delete cascade,
  allow_transcript_view boolean not null default true,
  redact_secrets_by_default boolean not null default true,
  allow_transcript_export boolean not null default false,
  retain_runtime_events_days integer not null default 30,
  retain_transcript_days integer not null default 30,
  archive_full_transcript boolean not null default false,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);
