-- name: CreateMissionArchive :one
insert into mission_archives (
  id,
  mission_id,
  status,
  manifest_object_key,
  bundle_object_key,
  hash
) values (
  $1, $2, $3, $4, $5, $6
)
returning id, mission_id, status, manifest_object_key, bundle_object_key, hash, created_at, finished_at;

-- name: GetMissionArchiveByMission :one
select id, mission_id, status, manifest_object_key, bundle_object_key, hash, created_at, finished_at
from mission_archives
where mission_id = $1
order by created_at desc
limit 1;

-- name: UpdateMissionArchiveStatus :one
update mission_archives
set
  status = $2,
  manifest_object_key = $3,
  bundle_object_key = $4,
  hash = $5,
  finished_at = case when $2 = 'uploaded' then now() else finished_at end
where id = $1
returning id, mission_id, status, manifest_object_key, bundle_object_key, hash, created_at, finished_at;

-- name: ListAdoptedDocumentVersionIDsByMission :many
select current_adopted_version_id
from shared_documents
where mission_id = $1
  and current_adopted_version_id is not null
order by updated_at desc;

-- name: ListMissionDecisionIDs :many
select id
from mission_decisions
where mission_id = $1
order by created_at desc;

-- name: ListRuntimeSummaryKeysByMission :many
select distinct ('runtime/' || executor_session_id || '.json') as object_key
from agent_runtime_events
where mission_id = $1
order by object_key asc;
