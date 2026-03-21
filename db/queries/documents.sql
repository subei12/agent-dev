-- name: CreateSharedDocument :one
insert into shared_documents (
  id,
  mission_id,
  kind,
  title,
  current_adopted_version_id
) values (
  $1, $2, $3, $4, $5
)
returning id, mission_id, kind, title, current_adopted_version_id, created_at, updated_at;

-- name: GetSharedDocument :one
select id, mission_id, kind, title, current_adopted_version_id, created_at, updated_at
from shared_documents
where id = $1 and mission_id = $2;

-- name: ListSharedDocumentsByMission :many
select id, mission_id, kind, title, current_adopted_version_id, created_at, updated_at
from shared_documents
where mission_id = $1
order by updated_at desc;

-- name: GetNextSharedDocumentVersionNumber :one
select coalesce(max(version), 0)::int4 + 1
from shared_document_versions
where document_id = $1;

-- name: CreateSharedDocumentVersion :one
insert into shared_document_versions (
  id,
  document_id,
  version,
  parent_version_id,
  status,
  content_format,
  storage_kind,
  object_key,
  content_hash,
  content_text,
  produced_by_agent_id,
  produced_by_run_id,
  source_round_id
) values (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
)
returning id, document_id, version, parent_version_id, status, content_format, storage_kind, object_key, content_hash, content_text, produced_by_agent_id, produced_by_run_id, source_round_id, created_at;

-- name: GetSharedDocumentVersion :one
select id, document_id, version, parent_version_id, status, content_format, storage_kind, object_key, content_hash, content_text, produced_by_agent_id, produced_by_run_id, source_round_id, created_at
from shared_document_versions
where id = $1;

-- name: ListSharedDocumentVersions :many
select id, document_id, version, parent_version_id, status, content_format, storage_kind, object_key, content_hash, content_text, produced_by_agent_id, produced_by_run_id, source_round_id, created_at
from shared_document_versions
where document_id = $1
order by version asc;

-- name: UpdateSharedDocumentVersionStatus :one
update shared_document_versions
set status = $2
where id = $1
returning id, document_id, version, parent_version_id, status, content_format, storage_kind, object_key, content_hash, content_text, produced_by_agent_id, produced_by_run_id, source_round_id, created_at;

-- name: ClearAdoptedSharedDocumentVersions :exec
update shared_document_versions
set status = 'proposed'
where document_id = $1 and status = 'adopted';

-- name: SetSharedDocumentCurrentAdoptedVersion :one
update shared_documents
set
  current_adopted_version_id = $2,
  updated_at = now()
where id = $1
returning id, mission_id, kind, title, current_adopted_version_id, created_at, updated_at;
