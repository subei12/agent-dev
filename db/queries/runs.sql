-- name: CreateRun :one
insert into runs (
  id,
  mission_id,
  task_item_id,
  status,
  inputs_json
) values (
  $1, $2, $3, $4, $5
)
returning id, mission_id, task_item_id, status, inputs_json, created_at, started_at, finished_at;

-- name: CreateNodeRun :one
insert into node_runs (
  id,
  run_id,
  node_key,
  status,
  executor_session_id
) values (
  $1, $2, $3, $4, $5
)
returning id, run_id, node_key, status, executor_session_id, created_at, started_at, finished_at;

-- name: UpdateRunStatus :one
update runs
set
  status = $2,
  finished_at = case when $2 in ('succeeded', 'failed', 'cancelled') then now() else finished_at end
where id = $1
returning id, mission_id, task_item_id, status, inputs_json, created_at, started_at, finished_at;

-- name: GetTaskClaimExecutionContext :one
select
  c.id as claim_id,
  b.mission_id,
  c.task_item_id,
  c.agent_id,
  m.admin_agent_id,
  t.title as task_title,
  t.type as task_type,
  t.downstream_task_ids_json,
  a.executor_profile_id,
  ep.command,
  ep.args_json
from task_claims c
join task_items t on t.id = c.task_item_id
join task_boards b on b.id = t.board_id
join missions m on m.id = b.mission_id
join agents a on a.id = c.agent_id
join executor_profiles ep on ep.id = a.executor_profile_id
where c.id = $1;
