-- +goose Up
create table if not exists discussion_sessions (
  id text primary key,
  mission_id text not null references missions(id) on delete cascade,
  topic text not null,
  status text not null,
  initiated_by_agent_id text not null,
  created_at timestamptz not null default now(),
  closed_at timestamptz
);

create table if not exists discussion_rounds (
  id text primary key,
  session_id text not null references discussion_sessions(id) on delete cascade,
  round_no integer not null,
  prompt_summary text not null,
  participant_agent_ids_json jsonb not null default '[]'::jsonb,
  summary text not null,
  open_questions_json jsonb not null default '[]'::jsonb,
  conflicts_json jsonb not null default '{}'::jsonb,
  conclusion_json jsonb not null default '{}'::jsonb,
  adopted_document_version_ids_json jsonb not null default '[]'::jsonb,
  created_at timestamptz not null default now(),
  unique(session_id, round_no)
);

create table if not exists shared_documents (
  id text primary key,
  mission_id text not null references missions(id) on delete cascade,
  kind text not null,
  title text not null,
  current_adopted_version_id text,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create table if not exists shared_document_versions (
  id text primary key,
  document_id text not null references shared_documents(id) on delete cascade,
  version integer not null,
  parent_version_id text,
  status text not null,
  content_format text not null,
  storage_kind text not null,
  object_key text,
  content_hash text not null,
  content_text text not null default '',
  produced_by_agent_id text,
  produced_by_run_id text,
  source_round_id text,
  created_at timestamptz not null default now(),
  unique(document_id, version)
);
