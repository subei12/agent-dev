-- name: CreateRepoCandidate :one
insert into repo_candidates (
  id,
  mission_id,
  task_item_id,
  repo_binding_id,
  base_commit_sha,
  parent_candidate_id,
  created_from,
  created_from_ref_id,
  tree_hash,
  patch_object_key,
  workspace_archive_object_key,
  summary,
  status,
  produced_by_run_id,
  produced_by_node_run_id
) values (
  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
)
returning id, mission_id, task_item_id, repo_binding_id, base_commit_sha, parent_candidate_id, created_from, created_from_ref_id, tree_hash, patch_object_key, workspace_archive_object_key, summary, status, produced_by_run_id, produced_by_node_run_id, created_at;

-- name: GetRepoCandidate :one
select id, mission_id, task_item_id, repo_binding_id, base_commit_sha, parent_candidate_id, created_from, created_from_ref_id, tree_hash, patch_object_key, workspace_archive_object_key, summary, status, produced_by_run_id, produced_by_node_run_id, created_at
from repo_candidates
where id = $1;
