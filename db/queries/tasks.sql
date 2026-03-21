-- name: CreateTaskBoard :one
insert into task_boards (
  id,
  mission_id,
  title
) values (
  $1, $2, $3
)
returning id, mission_id, title, created_at, updated_at;

-- name: CreateTaskItem :one
insert into task_items (
  id,
  board_id,
  title,
  type,
  status,
  assigned_agent_id,
  upstream_task_ids_json,
  downstream_task_ids_json,
  input_document_version_ids_json,
  input_repo_candidate_ids_json,
  definition_of_done_json
) values (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
)
returning id, board_id, title, type, status, assigned_agent_id, upstream_task_ids_json, downstream_task_ids_json, input_document_version_ids_json, input_repo_candidate_ids_json, definition_of_done_json, created_at, updated_at;

-- name: GetTaskItem :one
select id, board_id, title, type, status, assigned_agent_id, upstream_task_ids_json, downstream_task_ids_json, input_document_version_ids_json, input_repo_candidate_ids_json, definition_of_done_json, created_at, updated_at
from task_items
where id = $1;

-- name: UpdateTaskItemStatus :one
update task_items
set
  status = $2,
  updated_at = now()
where id = $1
returning id, board_id, title, type, status, assigned_agent_id, upstream_task_ids_json, downstream_task_ids_json, input_document_version_ids_json, input_repo_candidate_ids_json, definition_of_done_json, created_at, updated_at;

-- name: CreateTaskClaim :one
insert into task_claims (
  id,
  task_item_id,
  agent_id,
  status,
  claim_reason
) values (
  $1, $2, $3, $4, $5
)
returning id, task_item_id, agent_id, status, claim_reason, created_at, ended_at;

-- name: CreateTaskHandoff :one
insert into task_handoffs (
  id,
  task_item_id,
  from_agent_id,
  to_agent_id,
  to_admin_agent,
  summary,
  output_document_version_ids_json,
  output_repo_candidate_ids_json,
  output_artifact_version_ids_json,
  validation_summary,
  risk_summary,
  recommended_next_action,
  status
) values (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
)
returning id, task_item_id, from_agent_id, to_agent_id, to_admin_agent, summary, output_document_version_ids_json, output_repo_candidate_ids_json, output_artifact_version_ids_json, validation_summary, risk_summary, recommended_next_action, status, created_at, decided_at;

-- name: CreateReviewCheckpoint :one
insert into review_checkpoints (
  id,
  mission_id,
  task_item_id,
  kind,
  requested_by_agent_id,
  assigned_agent_id,
  status,
  summary,
  linked_document_version_ids_json,
  linked_repo_candidate_ids_json
) values (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
)
returning id, mission_id, task_item_id, kind, requested_by_agent_id, assigned_agent_id, status, summary, linked_document_version_ids_json, linked_repo_candidate_ids_json, created_at, finished_at;
