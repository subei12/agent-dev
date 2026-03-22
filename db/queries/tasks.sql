-- name: CreateTaskBoard :one
insert into task_boards (
  id,
  mission_id,
  title
) values (
  $1, $2, $3
)
returning id, mission_id, title, created_at, updated_at;

-- name: GetLatestTaskBoardByMission :one
select id, mission_id, title, created_at, updated_at
from task_boards
where mission_id = $1
order by created_at desc
limit 1;

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

-- name: ListTaskBoardsByMission :many
select id, mission_id, title, created_at, updated_at
from task_boards
where mission_id = $1
order by created_at asc;

-- name: GetTaskItem :one
select id, board_id, title, type, status, assigned_agent_id, upstream_task_ids_json, downstream_task_ids_json, input_document_version_ids_json, input_repo_candidate_ids_json, definition_of_done_json, created_at, updated_at
from task_items
where id = $1;

-- name: GetMissionTaskByIdentifier :one
select
  t.id,
  t.title,
  t.status,
  t.assigned_agent_id
from task_items t
join task_boards b on b.id = t.board_id
where b.mission_id = $1
  and (t.id = $2 or t.title = $2)
limit 1;

-- name: ListTaskItemsByBoard :many
select id, board_id, title, type, status, assigned_agent_id, upstream_task_ids_json, downstream_task_ids_json, input_document_version_ids_json, input_repo_candidate_ids_json, definition_of_done_json, created_at, updated_at
from task_items
where board_id = $1
order by created_at asc;

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

-- name: ListActiveTaskClaimIDs :many
select id
from task_claims
where status = 'active'
  and ended_at is null
order by created_at asc;

-- name: CompleteTaskClaim :one
update task_claims
set
  status = 'completed',
  execution_lock_token = null,
  last_error = null,
  ended_at = now()
where id = $1
returning id, task_item_id, agent_id, status, claim_reason, created_at, ended_at;

-- name: AcquireTaskClaimExecution :one
update task_claims
set
  execution_lock_token = $2,
  attempt_count = attempt_count + 1,
  last_heartbeat_at = now()
where id = $1
  and status = 'active'
  and ended_at is null
  and (
    execution_lock_token is null
    or last_heartbeat_at < $3
  )
returning id, task_item_id, agent_id, status, claim_reason, created_at, ended_at;

-- name: ReleaseTaskClaimExecution :one
update task_claims
set
  execution_lock_token = null,
  last_heartbeat_at = now()
where id = $1
returning id, task_item_id, agent_id, status, claim_reason, created_at, ended_at;

-- name: UpdateTaskClaimHeartbeat :one
update task_claims
set
  last_heartbeat_at = now()
where id = $1
  and execution_lock_token = $2
returning id, task_item_id, agent_id, status, claim_reason, created_at, ended_at;

-- name: FailTaskClaimExecution :one
update task_claims
set
  status = case when attempt_count >= $2 then 'failed' else 'active' end,
  execution_lock_token = null,
  last_error = $3,
  last_heartbeat_at = now(),
  ended_at = case when attempt_count >= $2 then now() else null end
where id = $1
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
