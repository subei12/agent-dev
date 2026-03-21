-- name: UpsertProjectAccessGrant :one
insert into project_access_grants (
  id,
  project_id,
  actor_user_id,
  can_manage_mission,
  can_view_transcripts,
  can_export_transcripts,
  can_manage_archive
) values (
  $1, $2, $3, $4, $5, $6, $7
)
on conflict (project_id, actor_user_id) do update
set
  can_manage_mission = excluded.can_manage_mission,
  can_view_transcripts = excluded.can_view_transcripts,
  can_export_transcripts = excluded.can_export_transcripts,
  can_manage_archive = excluded.can_manage_archive,
  updated_at = now()
returning id, project_id, actor_user_id, can_manage_mission, can_view_transcripts, can_export_transcripts, can_manage_archive, created_at, updated_at;

-- name: GetProjectAccessGrant :one
select id, project_id, actor_user_id, can_manage_mission, can_view_transcripts, can_export_transcripts, can_manage_archive, created_at, updated_at
from project_access_grants
where project_id = $1 and actor_user_id = $2;
