-- +goose Up
create table if not exists project_access_grants (
  id text primary key,
  project_id text not null references projects(id) on delete cascade,
  actor_user_id text not null,
  can_manage_mission boolean not null default false,
  can_view_transcripts boolean not null default false,
  can_export_transcripts boolean not null default false,
  can_manage_archive boolean not null default false,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  unique(project_id, actor_user_id)
);
