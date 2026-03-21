-- name: UpsertAgentPresence :one
insert into agent_presences (
  agent_id,
  availability,
  current_mission_id,
  current_task_item_id,
  current_executor_session_id,
  last_heartbeat_at,
  updated_at
) values (
  $1, $2, $3, $4, $5, now(), now()
)
on conflict (agent_id) do update
set
  availability = excluded.availability,
  current_mission_id = excluded.current_mission_id,
  current_task_item_id = excluded.current_task_item_id,
  current_executor_session_id = excluded.current_executor_session_id,
  last_heartbeat_at = now(),
  updated_at = now()
returning agent_id, availability, current_mission_id, current_task_item_id, current_executor_session_id, last_heartbeat_at, updated_at;

-- name: UpsertMissionAgentRuntime :one
insert into mission_agent_runtimes (
  id,
  mission_id,
  agent_id,
  status,
  current_task_item_id,
  current_run_id,
  current_node_run_id,
  current_executor_session_id,
  status_summary,
  started_at,
  updated_at
) values (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, now()
)
on conflict (mission_id, agent_id) do update
set
  status = excluded.status,
  current_task_item_id = excluded.current_task_item_id,
  current_run_id = excluded.current_run_id,
  current_node_run_id = excluded.current_node_run_id,
  current_executor_session_id = excluded.current_executor_session_id,
  status_summary = excluded.status_summary,
  started_at = coalesce(mission_agent_runtimes.started_at, excluded.started_at),
  updated_at = now()
returning id, mission_id, agent_id, status, current_task_item_id, current_run_id, current_node_run_id, current_executor_session_id, status_summary, started_at, updated_at;

-- name: CreateExecutorSession :one
insert into executor_sessions (
  id,
  mission_id,
  task_item_id,
  agent_id,
  executor_profile_id,
  backend,
  status
) values (
  $1, $2, $3, $4, $5, $6, $7
)
returning id, mission_id, task_item_id, agent_id, executor_profile_id, backend, status, started_at, ended_at;

-- name: GetExecutorSession :one
select id, mission_id, task_item_id, agent_id, executor_profile_id, backend, status, started_at, ended_at
from executor_sessions
where id = $1;

-- name: UpdateExecutorSessionStatus :one
update executor_sessions
set
  status = $2,
  ended_at = case when $2 in ('completed', 'failed', 'cancelled') then now() else ended_at end
where id = $1
returning id, mission_id, task_item_id, agent_id, executor_profile_id, backend, status, started_at, ended_at;

-- name: ListMissionAgentRuntimes :many
select id, mission_id, agent_id, status, current_task_item_id, current_run_id, current_node_run_id, current_executor_session_id, status_summary, started_at, updated_at
from mission_agent_runtimes
where mission_id = $1
order by updated_at desc;

-- name: CreateAgentRuntimeEvent :one
insert into agent_runtime_events (
  id,
  mission_id,
  agent_id,
  executor_session_id,
  run_id,
  node_run_id,
  task_item_id,
  level,
  category,
  type,
  title,
  summary,
  payload_json
) values (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
)
returning id, mission_id, agent_id, executor_session_id, run_id, node_run_id, task_item_id, level, category, type, title, summary, payload_json, occurred_at;

-- name: ListAgentRuntimeEventsBySession :many
select id, mission_id, agent_id, executor_session_id, run_id, node_run_id, task_item_id, level, category, type, title, summary, payload_json, occurred_at
from agent_runtime_events
where executor_session_id = $1
order by occurred_at asc;

-- name: CreateExecutorTranscript :one
insert into executor_transcripts (
  id,
  executor_session_id,
  mission_id,
  agent_id,
  storage_kind,
  redacted_object_key,
  raw_encrypted_object_key,
  status
) values (
  $1, $2, $3, $4, $5, $6, $7, $8
)
returning id, executor_session_id, mission_id, agent_id, storage_kind, redacted_object_key, raw_encrypted_object_key, status, started_at, ended_at;

-- name: GetExecutorTranscriptBySession :one
select id, executor_session_id, mission_id, agent_id, storage_kind, redacted_object_key, raw_encrypted_object_key, status, started_at, ended_at
from executor_transcripts
where executor_session_id = $1;

-- name: UpdateExecutorTranscriptStatus :one
update executor_transcripts
set
  status = $2,
  ended_at = case when $2 = 'sealed' then now() else ended_at end
where id = $1
returning id, executor_session_id, mission_id, agent_id, storage_kind, redacted_object_key, raw_encrypted_object_key, status, started_at, ended_at;

-- name: GetNextTranscriptEntrySeq :one
select coalesce(max(seq), 0)::int4 + 1
from executor_transcript_entries
where transcript_id = $1;

-- name: CreateExecutorTranscriptEntry :one
insert into executor_transcript_entries (
  id,
  transcript_id,
  seq,
  role,
  entry_type,
  redacted_text,
  content_json,
  raw_object_key
) values (
  $1, $2, $3, $4, $5, $6, $7, $8
)
returning id, transcript_id, seq, role, entry_type, redacted_text, content_json, raw_object_key, created_at;

-- name: ListExecutorTranscriptEntriesBySession :many
select e.id, e.transcript_id, e.seq, e.role, e.entry_type, e.redacted_text, e.content_json, e.raw_object_key, e.created_at
from executor_transcript_entries e
join executor_transcripts t on t.id = e.transcript_id
where t.executor_session_id = $1
order by e.seq asc;

-- name: CreateTranscriptAccessAudit :one
insert into transcript_access_audits (
  id,
  transcript_id,
  actor_user_id,
  access_mode,
  reason
) values (
  $1, $2, $3, $4, $5
)
returning id, transcript_id, actor_user_id, access_mode, reason, created_at;

-- name: ListTranscriptAccessAuditsBySession :many
select a.id, a.transcript_id, a.actor_user_id, a.access_mode, a.reason, a.created_at
from transcript_access_audits a
join executor_transcripts t on t.id = a.transcript_id
where t.executor_session_id = $1
order by a.created_at desc;
