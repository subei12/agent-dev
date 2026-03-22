package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"

	platformconfig "github.com/your-org/agent-platform/internal/config"
	platformdb "github.com/your-org/agent-platform/internal/db"
)

// main 启动当前可执行入口。
func main() {
	ctx := context.Background()
	cfg := platformconfig.MustLoad()

	pool, err := platformdb.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	if err := seedDemoData(ctx, pool); err != nil {
		log.Fatal(err)
	}

	fmt.Println("seeded demo data for proj_1 / mission_1")
}

// seedDemoData 重置并重新写入确定性的演示项目数据。
func seedDemoData(ctx context.Context, pool *pgxpool.Pool) error {
	// 1. 按依赖安全的逆序清理旧的种子数据。
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	for _, stmt := range []string{
		`delete from transcript_access_audits where id = 'audit_1'`,
		`delete from executor_transcript_entries where transcript_id = 'transcript_1'`,
		`delete from executor_transcripts where id = 'transcript_1'`,
		`delete from agent_runtime_events where executor_session_id = 'session_1'`,
		`delete from mission_agent_runtimes where id = 'runtime_1'`,
		`delete from agent_presences where agent_id in ('agent_admin','agent_backend')`,
		`delete from review_checkpoints where id = 'checkpoint_1'`,
		`delete from task_handoffs where id = 'handoff_1'`,
		`delete from task_claims where id = 'claim_1'`,
		`delete from task_items where id in ('task_design','task_runtime','task_review')`,
		`delete from task_boards where id = 'board_1'`,
		`delete from shared_document_versions where id in ('docv_1','docv_2')`,
		`delete from shared_documents where id = 'doc_1'`,
		`delete from discussion_rounds where id = 'round_1'`,
		`delete from discussion_sessions where id = 'discussion_1'`,
		`delete from mission_decisions where id in ('decision_1','decision_2')`,
		`delete from node_runs where id = 'node_run_1'`,
		`delete from runs where id = 'run_1'`,
		`delete from mission_archives where id = 'archive_1'`,
		`delete from approvals where id = 'approval_1'`,
		`delete from repo_candidate_publishes where id = 'publish_1'`,
		`delete from repo_candidates where id = 'cand_1'`,
		`delete from mission_members where id in ('member_admin','member_backend')`,
		`delete from runtime_observability_policies where id = 'runtime_policy_1'`,
		`delete from agent_admin_policies where id = 'admin_policy_1'`,
		`delete from missions where id = 'mission_1'`,
		`delete from agent_team_members where id in ('team_member_admin','team_member_backend')`,
		`delete from agent_teams where id = 'team_1'`,
		`delete from project_access_grants where id in ('grant_demo_admin','grant_frontend_demo')`,
		`delete from agents where id in ('agent_admin','agent_backend')`,
		`delete from roles where id in ('role_admin','role_backend')`,
		`delete from executor_profiles where id in ('exec_admin','exec_backend')`,
		`delete from repo_bindings where id = 'repo_1'`,
		`delete from projects where id = 'proj_1'`,
	} {
		if _, err := tx.Exec(ctx, stmt); err != nil {
			return err
		}
	}

	// 2. 写入一致的演示项目、mission、runtime 和 archive 快照，供本地演示使用。
	inserts := []struct {
		sql  string
		args []any
	}{
		{`insert into projects (id, name, description, status) values ($1, $2, $3, $4)`,
			[]any{"proj_1", "演示项目", "用于 Mission 工作台演示的种子项目", "active"}},
		{`insert into repo_bindings (id, project_id, provider, remote_url, default_branch, auth_mode, write_enabled) values ($1,$2,$3,$4,$5,$6,$7)`,
			[]any{"repo_1", "proj_1", "generic", "git@example.com:demo/repo.git", "main", "token", false}},
		{`insert into executor_profiles (id, project_id, name, type, command, args_json, env_secret_refs_json, timeout_sec, max_concurrency, max_prompt_chars, allow_repo_read, allow_repo_write, allow_network)
		  values ($1,$2,$3,$4,$5,$6::jsonb,$7::jsonb,$8,$9,$10,$11,$12,$13)`,
			[]any{"exec_admin", "proj_1", "管理员执行器", "codex_cli", "bash", `["-lc","printf 'admin\n'"]`, `[]`, 3600, 1, 16000, true, false, false}},
		{`insert into executor_profiles (id, project_id, name, type, command, args_json, env_secret_refs_json, timeout_sec, max_concurrency, max_prompt_chars, allow_repo_read, allow_repo_write, allow_network)
		  values ($1,$2,$3,$4,$5,$6::jsonb,$7::jsonb,$8,$9,$10,$11,$12,$13)`,
			[]any{"exec_backend", "proj_1", "后端执行器", "codex_cli", "bash", `["-lc","printf 'backend\n'"]`, `[]`, 3600, 1, 16000, true, true, false}},
		{`insert into roles (id, project_id, key, name, prompts_json, permissions_json) values ($1,$2,$3,$4,$5::jsonb,$6::jsonb)`,
			[]any{"role_admin", "proj_1", "admin", "管理员", `{}`, `{"canWriteArtifacts":true}`}},
		{`insert into roles (id, project_id, key, name, prompts_json, permissions_json) values ($1,$2,$3,$4,$5::jsonb,$6::jsonb)`,
			[]any{"role_backend", "proj_1", "backend", "后端", `{}`, `{"canWriteArtifacts":true}`}},
		{`insert into agents (id, project_id, name, role_id, executor_profile_id, enabled) values ($1,$2,$3,$4,$5,$6)`,
			[]any{"agent_admin", "proj_1", "管理员 Agent", "role_admin", "exec_admin", true}},
		{`insert into agents (id, project_id, name, role_id, executor_profile_id, enabled) values ($1,$2,$3,$4,$5,$6)`,
			[]any{"agent_backend", "proj_1", "后端 Agent", "role_backend", "exec_backend", true}},
		{`insert into agent_teams (id, project_id, name, description) values ($1,$2,$3,$4)`,
			[]any{"team_1", "proj_1", "核心交付团队", "管理员 + 后端执行团队"}},
		{`insert into agent_team_members (id, team_id, agent_id, role_key, capabilities_json) values ($1,$2,$3,$4,$5::jsonb)`,
			[]any{"team_member_admin", "team_1", "agent_admin", "admin", `{"isAdmin":true,"canAdoptDocument":true}`}},
		{`insert into agent_team_members (id, team_id, agent_id, role_key, capabilities_json) values ($1,$2,$3,$4,$5::jsonb)`,
			[]any{"team_member_backend", "team_1", "agent_backend", "backend", `{"canBeAssigned":true}`}},
		{`insert into agent_admin_policies (id, project_id) values ($1,$2)`,
			[]any{"admin_policy_1", "proj_1"}},
		{`insert into runtime_observability_policies (id, project_id) values ($1,$2)`,
			[]any{"runtime_policy_1", "proj_1"}},
		{`insert into project_access_grants (id, project_id, actor_user_id, can_manage_mission, can_view_transcripts, can_export_transcripts, can_manage_archive)
		  values ($1,$2,$3,$4,$5,$6,$7)`,
			[]any{"grant_demo_admin", "proj_1", "demo-admin", true, true, false, true}},
		{`insert into project_access_grants (id, project_id, actor_user_id, can_manage_mission, can_view_transcripts, can_export_transcripts, can_manage_archive)
		  values ($1,$2,$3,$4,$5,$6,$7)`,
			[]any{"grant_frontend_demo", "proj_1", "frontend-demo-user", true, true, false, true}},
		{`insert into missions (id, project_id, title, description, source_type, status, admin_agent_id, team_id, repo_binding_id, created_by)
		  values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
			[]any{"mission_1", "proj_1", "交付运行态观测工作台", "用于演示讨论、文档、任务、运行日志与归档的 Mission。", "project_task", "implementation", "agent_admin", "team_1", "repo_1", "demo-user"}},
		{`insert into mission_members (id, mission_id, agent_id, role_key, is_admin) values ($1,$2,$3,$4,$5)`,
			[]any{"member_admin", "mission_1", "agent_admin", "admin", true}},
		{`insert into mission_members (id, mission_id, agent_id, role_key, is_admin) values ($1,$2,$3,$4,$5)`,
			[]any{"member_backend", "mission_1", "agent_backend", "backend", false}},
		{`insert into mission_decisions (id, mission_id, decided_by_agent_id, decision, summary, related_document_version_ids_json, related_repo_candidate_ids_json)
		  values ($1,$2,$3,$4,$5,$6::jsonb,$7::jsonb)`,
			[]any{"decision_1", "mission_1", "agent_admin", "approve_plan", "架构方案已通过。", `[]`, `[]`}},
		{`insert into mission_decisions (id, mission_id, decided_by_agent_id, decision, summary, related_document_version_ids_json, related_repo_candidate_ids_json)
		  values ($1,$2,$3,$4,$5,$6::jsonb,$7::jsonb)`,
			[]any{"decision_2", "mission_1", "agent_admin", "enter_implementation", "进入开发阶段。", `["docv_2"]`, `[]`}},
		{`insert into discussion_sessions (id, mission_id, topic, status, initiated_by_agent_id) values ($1,$2,$3,$4,$5)`,
			[]any{"discussion_1", "mission_1", "运行态观测范围", "open", "agent_admin"}},
		{`insert into discussion_rounds (id, session_id, round_no, prompt_summary, participant_agent_ids_json, summary, open_questions_json, conflicts_json, conclusion_json, adopted_document_version_ids_json)
		  values ($1,$2,$3,$4,$5::jsonb,$6,$7::jsonb,$8::jsonb,$9::jsonb,$10::jsonb)`,
			[]any{"round_1", "discussion_1", 1, "定义运行态面板范围", `["agent_admin","agent_backend"]`, "确认默认优先展示事件时间线，按需查看 transcript。", `[]`, `{}`, `{"decision":"event-first"}`, `[]`}},
		{`insert into shared_documents (id, mission_id, kind, title, current_adopted_version_id) values ($1,$2,$3,$4,$5)`,
			[]any{"doc_1", "mission_1", "architecture", "架构快照", "docv_2"}},
		{`insert into shared_document_versions (id, document_id, version, parent_version_id, status, content_format, storage_kind, object_key, content_hash, content_text, source_round_id)
		  values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
			[]any{"docv_1", "doc_1", 1, nil, "proposed", "md", "db_text", nil, "hash-docv1", "初始草稿", "round_1"}},
		{`insert into shared_document_versions (id, document_id, version, parent_version_id, status, content_format, storage_kind, object_key, content_hash, content_text, source_round_id)
		  values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
			[]any{"docv_2", "doc_1", 2, "docv_1", "adopted", "md", "db_text", nil, "hash-docv2", "已采纳架构版本", "round_1"}},
		{`insert into task_boards (id, mission_id, title) values ($1,$2,$3)`,
			[]any{"board_1", "mission_1", "交付看板"}},
		{`insert into task_items (id, board_id, title, type, status, assigned_agent_id, upstream_task_ids_json, downstream_task_ids_json, input_document_version_ids_json, input_repo_candidate_ids_json, definition_of_done_json)
		  values ($1,$2,$3,$4,$5,$6,$7::jsonb,$8::jsonb,$9::jsonb,$10::jsonb,$11::jsonb)`,
			[]any{"task_design", "board_1", "收敛架构文档", "design", "done", "agent_admin", `[]`, `["task_runtime"]`, `["docv_2"]`, `[]`, `{"done":"approved doc"}`}},
		{`insert into task_items (id, board_id, title, type, status, assigned_agent_id, upstream_task_ids_json, downstream_task_ids_json, input_document_version_ids_json, input_repo_candidate_ids_json, definition_of_done_json)
		  values ($1,$2,$3,$4,$5,$6,$7::jsonb,$8::jsonb,$9::jsonb,$10::jsonb,$11::jsonb)`,
			[]any{"task_runtime", "board_1", "实现运行态观测", "code", "handoff_pending", "agent_backend", `["task_design"]`, `[]`, `["docv_2"]`, `[]`, `{"done":"session and runtime board"}`}},
		{`insert into task_items (id, board_id, title, type, status, assigned_agent_id, upstream_task_ids_json, downstream_task_ids_json, input_document_version_ids_json, input_repo_candidate_ids_json, definition_of_done_json)
		  values ($1,$2,$3,$4,$5,$6,$7::jsonb,$8::jsonb,$9::jsonb,$10::jsonb,$11::jsonb)`,
			[]any{"task_review", "board_1", "评审 transcript 访问策略", "review", "todo", "agent_admin", `["task_runtime"]`, `[]`, `["docv_2"]`, `[]`, `{"done":"review checkpoint"}`}},
		{`insert into task_claims (id, task_item_id, agent_id, status, claim_reason) values ($1,$2,$3,$4,$5)`,
			[]any{"claim_1", "task_runtime", "agent_backend", "active", "assigned"}},
		{`insert into task_handoffs (id, task_item_id, from_agent_id, to_agent_id, to_admin_agent, summary, output_document_version_ids_json, output_repo_candidate_ids_json, output_artifact_version_ids_json, validation_summary, risk_summary, recommended_next_action, status)
		  values ($1,$2,$3,$4,$5,$6,$7::jsonb,$8::jsonb,$9::jsonb,$10,$11,$12,$13)`,
			[]any{"handoff_1", "task_runtime", "agent_backend", nil, true, "实现完成，等待管理员复核。", `["docv_2"]`, `["cand_1"]`, `[]`, "单元测试已通过", "transcript 保留策略仍待确认", "继续进入评审", "pending"}},
		{`insert into review_checkpoints (id, mission_id, task_item_id, kind, requested_by_agent_id, assigned_agent_id, status, summary, linked_document_version_ids_json, linked_repo_candidate_ids_json)
		  values ($1,$2,$3,$4,$5,$6,$7,$8,$9::jsonb,$10::jsonb)`,
			[]any{"checkpoint_1", "mission_1", "task_review", "doc_review", "agent_admin", "agent_admin", "pending", "等待 transcript 策略评审。", `["docv_2"]`, `["cand_1"]`}},
		{`insert into runs (id, mission_id, task_item_id, status, inputs_json) values ($1,$2,$3,$4,$5::jsonb)`,
			[]any{"run_1", "mission_1", "task_runtime", "succeeded", `{}`}},
		{`insert into node_runs (id, run_id, node_key, status, executor_session_id) values ($1,$2,$3,$4,$5)`,
			[]any{"node_run_1", "run_1", "task_execute", "succeeded", "session_1"}},
		{`insert into repo_candidates (id, mission_id, task_item_id, repo_binding_id, base_commit_sha, parent_candidate_id, created_from, created_from_ref_id, tree_hash, patch_object_key, workspace_archive_object_key, summary, status, produced_by_run_id, produced_by_node_run_id)
		  values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)`,
			[]any{"cand_1", "mission_1", "task_runtime", "repo_1", "abc123", nil, "repo_snapshot", "snapshot_1", "treehash-1", nil, nil, "运行态观测实现结果", "draft", "run_1", "node_run_1"}},
		{`insert into approvals (id, mission_id, run_id, node_run_id, action, subject_type, subject_id, intent_snapshot_json, intent_hash, status, comment, created_by, decided_by)
		  values ($1,$2,$3,$4,$5,$6,$7,$8::jsonb,$9,$10,$11,$12,$13)`,
			[]any{"approval_1", "mission_1", "run_1", "node_run_1", "repo_candidate_publish", "repo_candidate_publish", "cand_1", `{"treeHash":"treehash-1"}`, "intenthash-1", "pending", "等待审批", "demo-user", nil}},
		{`insert into mission_archives (id, mission_id, status, manifest_object_key, bundle_object_key, hash) values ($1,$2,$3,$4,$5,$6)`,
			[]any{"archive_1", "mission_1", "uploaded", "archives/mission_1/manifest.json", "archives/mission_1/bundle.zip", "archive-hash-1"}},
		{`insert into agent_presences (agent_id, availability, current_mission_id, current_task_item_id, current_executor_session_id) values ($1,$2,$3,$4,$5)`,
			[]any{"agent_backend", "busy", "mission_1", "task_runtime", "session_1"}},
		{`insert into mission_agent_runtimes (id, mission_id, agent_id, status, current_task_item_id, current_executor_session_id, status_summary) values ($1,$2,$3,$4,$5,$6,$7)`,
			[]any{"runtime_1", "mission_1", "agent_backend", "working", "task_runtime", "session_1", "正在输出 codex 会话"}},
		{`insert into executor_sessions (id, mission_id, task_item_id, agent_id, executor_profile_id, backend, status) values ($1,$2,$3,$4,$5,$6,$7)`,
			[]any{"session_1", "mission_1", "task_runtime", "agent_backend", "exec_backend", "codex_cli", "completed"}},
		{`insert into executor_transcripts (id, executor_session_id, mission_id, agent_id, storage_kind, status) values ($1,$2,$3,$4,$5,$6)`,
			[]any{"transcript_1", "session_1", "mission_1", "agent_backend", "db_text", "sealed"}},
		{`insert into executor_transcript_entries (id, transcript_id, seq, role, entry_type, redacted_text, content_json, raw_object_key) values ($1,$2,$3,$4,$5,$6,$7::jsonb,$8)`,
			[]any{"entry_1", "transcript_1", 1, "assistant", "message", "已为管理员生成运行摘要。", `{}`, nil}},
		{`insert into executor_transcript_entries (id, transcript_id, seq, role, entry_type, redacted_text, content_json, raw_object_key) values ($1,$2,$3,$4,$5,$6,$7::jsonb,$8)`,
			[]any{"entry_2", "transcript_1", 2, "tool_result", "tool_result", "已持久化事件负载。", `{}`, nil}},
		{`insert into agent_runtime_events (id, mission_id, agent_id, executor_session_id, run_id, node_run_id, task_item_id, level, category, type, title, summary, payload_json)
		  values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13::jsonb)`,
			[]any{"event_1", "mission_1", "agent_backend", "session_1", "run_1", "node_run_1", "task_runtime", "info", "session", "session_started", "会话已启动", "worker 已启动 codex_cli", `{}`}},
		{`insert into agent_runtime_events (id, mission_id, agent_id, executor_session_id, run_id, node_run_id, task_item_id, level, category, type, title, summary, payload_json)
		  values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13::jsonb)`,
			[]any{"event_2", "mission_1", "agent_backend", "session_1", "run_1", "node_run_1", "task_runtime", "info", "handoff", "waiting_admin", "等待管理员", "任务完成但没有下游节点", `{}`}},
		{`insert into transcript_access_audits (id, transcript_id, actor_user_id, access_mode, reason) values ($1,$2,$3,$4,$5)`,
			[]any{"audit_1", "transcript_1", "demo-admin", "redacted_transcript", "审阅转录"}},
	}

	for _, insert := range inserts {
		if _, err := tx.Exec(ctx, insert.sql, insert.args...); err != nil {
			return err
		}
	}

	// 3. 仅在所有语句都成功后再提交整批种子数据。
	return tx.Commit(ctx)
}
