-- name: CreateMission :one
insert into missions (
  id,
  project_id,
  title,
  description,
  source_type,
  status,
  admin_agent_id,
  team_id,
  repo_binding_id,
  created_by
) values (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
)
returning id, project_id, title, description, source_type, status, admin_agent_id, team_id, repo_binding_id, created_by, created_at, updated_at, completed_at;

-- name: ListMissionsByProject :many
select id, project_id, title, description, source_type, status, admin_agent_id, team_id, repo_binding_id, created_by, created_at, updated_at, completed_at
from missions
where project_id = $1
order by updated_at desc;

-- name: GetMission :one
select id, project_id, title, description, source_type, status, admin_agent_id, team_id, repo_binding_id, created_by, created_at, updated_at, completed_at
from missions
where id = $1 and project_id = $2;

-- name: CreateMissionDecision :one
insert into mission_decisions (
  id,
  mission_id,
  decided_by_agent_id,
  decision,
  summary,
  related_task_item_id,
  related_document_version_ids_json,
  related_repo_candidate_ids_json
) values (
  $1, $2, $3, $4, $5, $6, $7, $8
)
returning id, mission_id, decided_by_agent_id, decision, summary, related_task_item_id, related_document_version_ids_json, related_repo_candidate_ids_json, created_at;

-- name: UpdateMissionStatus :one
update missions
set
  status = $2,
  updated_at = now(),
  completed_at = case when $2 = 'completed' then now() else completed_at end
where id = $1
returning id, project_id, title, description, source_type, status, admin_agent_id, team_id, repo_binding_id, created_by, created_at, updated_at, completed_at;
