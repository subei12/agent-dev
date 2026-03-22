-- name: ListAgentsByProject :many
select id, project_id, name, role_id, executor_profile_id, enabled, created_at, updated_at
from agents
where project_id = $1
order by created_at desc;

-- name: ListExecutorProfilesByProject :many
select id, project_id, name, type, command, args_json, env_secret_refs_json, timeout_sec, max_concurrency, max_prompt_chars, allow_repo_read, allow_repo_write, allow_network, created_at, updated_at
from executor_profiles
where project_id = $1
order by created_at desc;

-- name: UpdateAgent :one
update agents
set
  name = $2,
  enabled = $3,
  updated_at = now()
where id = $1
returning id, project_id, name, role_id, executor_profile_id, enabled, created_at, updated_at;

-- name: UpdateExecutorProfile :one
update executor_profiles
set
  type = $2,
  command = $3,
  args_json = $4,
  timeout_sec = $5,
  max_concurrency = $6,
  allow_repo_read = $7,
  allow_repo_write = $8,
  allow_network = $9,
  updated_at = now()
where id = $1
returning id, project_id, name, type, command, args_json, env_secret_refs_json, timeout_sec, max_concurrency, max_prompt_chars, allow_repo_read, allow_repo_write, allow_network, created_at, updated_at;
