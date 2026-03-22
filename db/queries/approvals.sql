-- name: CreateApproval :one
insert into approvals (
  id,
  mission_id,
  run_id,
  node_run_id,
  action,
  subject_type,
  subject_id,
  intent_snapshot_json,
  intent_hash,
  status,
  comment,
  created_by,
  decided_by
) values (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
)
returning id, mission_id, run_id, node_run_id, action, subject_type, subject_id, intent_snapshot_json, intent_hash, status, comment, created_by, decided_by, created_at, decided_at;

-- name: GetApproval :one
select id, mission_id, run_id, node_run_id, action, subject_type, subject_id, intent_snapshot_json, intent_hash, status, comment, created_by, decided_by, created_at, decided_at
from approvals
where id = $1;

-- name: ListApprovalIDsByMission :many
select id
from approvals
where mission_id = $1
order by created_at desc;

-- name: ListApprovalsByMission :many
select id, mission_id, run_id, node_run_id, action, subject_type, subject_id, intent_snapshot_json, intent_hash, status, comment, created_by, decided_by, created_at, decided_at
from approvals
where mission_id = $1
order by created_at desc;
