-- name: CreateDiscussionSession :one
insert into discussion_sessions (
  id,
  mission_id,
  topic,
  status,
  initiated_by_agent_id
) values (
  $1, $2, $3, $4, $5
)
returning id, mission_id, topic, status, initiated_by_agent_id, created_at, closed_at;

-- name: ListDiscussionSessionsByMission :many
select id, mission_id, topic, status, initiated_by_agent_id, created_at, closed_at
from discussion_sessions
where mission_id = $1
order by created_at desc;

-- name: GetNextDiscussionRoundNumber :one
select coalesce(max(round_no), 0)::int4 + 1
from discussion_rounds
where session_id = $1;

-- name: CreateDiscussionRound :one
insert into discussion_rounds (
  id,
  session_id,
  round_no,
  prompt_summary,
  participant_agent_ids_json,
  summary,
  open_questions_json,
  conflicts_json,
  conclusion_json,
  adopted_document_version_ids_json
) values (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
)
returning id, session_id, round_no, prompt_summary, participant_agent_ids_json, summary, open_questions_json, conflicts_json, conclusion_json, adopted_document_version_ids_json, created_at;
