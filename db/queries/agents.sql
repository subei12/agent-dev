-- name: ListAgentsByProject :many
select id, project_id, name, role_id, executor_profile_id, enabled, created_at, updated_at
from agents
where project_id = $1
order by created_at desc;
