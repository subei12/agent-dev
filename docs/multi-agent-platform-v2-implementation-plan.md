# Multi-Agent Platform v2 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 从当前方案仓库起步，落地一个可运行的多 Agent 协作治理平台 v2，覆盖 Mission 协作治理、共享文档、任务接力、CLI 执行器、运行态观测、transcript 查看、审批发布与 S3 归档。

**Architecture:** 采用单仓库、双进程后端和单前端工作台架构。后端使用 Go 拆为 `api` 与 `worker` 两个进程，共享 `internal/` 域模块、PostgreSQL、S3/MinIO、SSE 事件流和 CLI 适配层；前端使用 React + Vite 构建 Mission 工作台、任务板、Agent 运行态看板和 Transcript Inspector。

**Tech Stack:** Go, chi, pgx, sqlc, PostgreSQL, MinIO/S3, OpenTelemetry, SSE, `creack/pty`, React, TypeScript, Vite, React Router, Tailwind CSS, shadcn/ui, React Flow, TanStack Query, Monaco Editor, Playwright, pnpm.

---

## Scope And Assumptions

- 当前仓库没有实现代码，计划按“新项目骨架”落地。
- 先做单体 API + 单独 Worker，不做微服务、不做 Temporal、不做 Kafka。
- 先支持 `codex`、`claude code`、`gemini cli` 三种本地 CLI 执行器。
- 先做项目级权限和最小操作人标识，不做企业级 RBAC/SSO。
- 默认 transcript 展示脱敏版；原始未脱敏内容仅允许后端加密保存，不直接在默认 UI 暴露。
- 先以本地 Docker Compose 开发环境为主，生产部署脚本后补。

## Target Repository Layout

```text
backend/
  cmd/
    api/
    worker/
  internal/
    app/
    config/
    http/
    sse/
    telemetry/
    db/
    storage/
    mission/
    discussion/
    document/
    task/
    runtime/
    executor/
    approval/
    archive/
    project/
db/
  migrations/
  queries/
frontend/
  src/
    app/
    components/
    pages/
    features/
      missions/
      runtime/
      discussions/
      documents/
      tasks/
      approvals/
      archive/
  public/
ops/
  docker-compose.yml
openapi/
  openapi.yaml
docs/
  multi-agent-platform-v2-collaboration-governance-plan.md
  multi-agent-platform-v2-implementation-plan.md
Makefile
.env.example
README.md
```

## Module Responsibilities

- `backend/cmd/api`: HTTP API、SSE、用户入口。
- `backend/cmd/worker`: 任务调度、CLI 执行、transcript 捕获、归档执行。
- `backend/internal/app`: 进程组装、依赖注入、生命周期。
- `backend/internal/http`: router、handler、middleware、错误模型。
- `backend/internal/db`: pgx 连接池、事务工具、sqlc 封装。
- `backend/internal/storage`: S3/MinIO 客户端与对象键策略。
- `backend/internal/mission`: Mission、成员、阶段迁移、管理员策略。
- `backend/internal/discussion`: 讨论会话、讨论轮次、收敛结果。
- `backend/internal/document`: 共享文档、版本、采纳、视图。
- `backend/internal/task`: TaskBoard、TaskItem、TaskClaim、TaskHandoff、ReviewCheckpoint。
- `backend/internal/runtime`: AgentPresence、MissionAgentRuntime、ExecutorSession、RuntimeEvent、Transcript、脱敏、访问审计。
- `backend/internal/executor`: CLI 适配接口、PTY 进程管理、事件映射。
- `backend/internal/approval`: Approval、intent snapshot、PublishRecord。
- `backend/internal/archive`: MissionArchive、ArchiveManifest、归档打包。
- `frontend/src/pages`: 页面级路由容器。
- `frontend/src/features/*`: 每个业务域的查询、表单、展示和局部状态。

## Implementation Order

按下面顺序实现，避免前端和执行层在基础模型未定时反复返工：

1. 项目骨架与本地环境
2. 后端平台基础设施
3. 数据库迁移与 sqlc
4. Mission 协作治理 API
5. 文档、讨论、任务接力 API
6. 运行态观测与 transcript 管线
7. CLI 执行器与 Worker
8. 审批、发布、归档
9. 前端工作台与 Agent 运行态界面
10. 端到端测试与 CI

## Cross-Cutting Rules

- 数据库变更一律先写 SQL migration，再补 sqlc query，再写 handler/service。
- API contract 先写 `openapi/openapi.yaml`，再实现 handler。
- transcript 先落结构化事件，再补原始 transcript 存储，不允许只存 stdout 文本。
- 所有需要审计的查看行为必须写审计表。
- Worker 改动必须配套最少一个后端测试或回放脚本。
- 前端页面必须接入真实 API 类型，不允许长时间停留在硬编码 mock。

---

### Task 1: Bootstrap Repository Skeleton And Local Dev Stack

**Files:**
- Create: `backend/go.mod`
- Create: `backend/cmd/api/main.go`
- Create: `backend/internal/config/config.go`
- Create: `backend/internal/http/router.go`
- Create: `backend/internal/http/health_handler.go`
- Create: `backend/internal/http/health_handler_test.go`
- Create: `frontend/package.json`
- Create: `frontend/tsconfig.json`
- Create: `frontend/vite.config.ts`
- Create: `frontend/src/main.tsx`
- Create: `frontend/src/App.tsx`
- Create: `ops/docker-compose.yml`
- Create: `Makefile`
- Create: `.env.example`
- Modify: `README.md`
- Test: `backend/internal/http/health_handler_test.go`

- [ ] **Step 1: Create base directories and manifests**

```bash
mkdir -p backend/cmd/api backend/internal/{config,http} frontend/src ops openapi db/{migrations,queries}
cd backend && go mod init github.com/your-org/agent-platform
cd ../frontend && pnpm init
```

- [ ] **Step 2: Write the failing backend health test**

```go
func TestHealthHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	router := NewRouter()
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code)
	require.JSONEq(t, `{"status":"ok"}`, rec.Body.String())
}
```

- [ ] **Step 3: Implement minimal API skeleton and health endpoint**

```go
func NewRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		render.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	return r
}
```

- [ ] **Step 4: Add local infrastructure and developer commands**

```yaml
services:
  postgres:
    image: postgres:16
  minio:
    image: minio/minio
```

- [ ] **Step 5: Run smoke checks**

Run: `cd backend && go test ./...`  
Expected: health handler test passes.

Run: `docker compose -f ops/docker-compose.yml config`  
Expected: compose file renders without validation errors.

- [ ] **Step 6: Commit**

```bash
git add backend frontend ops Makefile .env.example README.md
git commit -m "chore: bootstrap repository skeleton"
```

### Task 2: Build Backend App Foundation And Shared Infrastructure

**Files:**
- Create: `backend/cmd/worker/main.go`
- Create: `backend/internal/app/api_app.go`
- Create: `backend/internal/app/worker_app.go`
- Create: `backend/internal/http/middleware/request_id.go`
- Create: `backend/internal/http/middleware/error_response.go`
- Create: `backend/internal/sse/hub.go`
- Create: `backend/internal/telemetry/otel.go`
- Create: `backend/internal/db/pool.go`
- Create: `backend/internal/storage/minio.go`
- Create: `backend/internal/config/config_test.go`
- Test: `backend/internal/config/config_test.go`

- [ ] **Step 1: Write the failing config loader test**

```go
func TestLoadConfigRequiresDatabaseURL(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	t.Setenv("DATABASE_URL", "")

	_, err := Load()
	require.ErrorContains(t, err, "DATABASE_URL")
}
```

- [ ] **Step 2: Implement shared config and app wiring**

```go
type Config struct {
	Env         string
	DatabaseURL string
	S3Endpoint  string
}
```

- [ ] **Step 3: Add API and worker entrypoints**

```go
func main() {
	cfg := config.MustLoad()
	app := app.NewWorker(cfg)
	app.Run(context.Background())
}
```

- [ ] **Step 4: Add middleware, SSE hub, pgx pool, and MinIO client wrappers**

```go
type Hub struct {
	subscribers map[string][]chan Event
}
```

- [ ] **Step 5: Run backend tests**

Run: `cd backend && go test ./...`  
Expected: config tests and existing API tests pass.

- [ ] **Step 6: Commit**

```bash
git add backend
git commit -m "feat: add backend app foundation"
```

### Task 3: Add Database Migrations, sqlc, And Core Project Schema

**Files:**
- Create: `sqlc.yaml`
- Create: `db/migrations/000001_core_projects_agents.sql`
- Create: `db/queries/projects.sql`
- Create: `db/queries/agents.sql`
- Create: `backend/internal/db/queries.go`
- Create: `backend/internal/project/repository.go`
- Create: `backend/internal/project/repository_test.go`
- Test: `backend/internal/project/repository_test.go`

- [ ] **Step 1: Write a failing repository integration test**

```go
func TestCreateProject(t *testing.T) {
	repo := newProjectRepo(t)
	project, err := repo.Create(ctx, CreateProjectParams{Name: "Demo"})
	require.NoError(t, err)
	require.Equal(t, "Demo", project.Name)
}
```

- [ ] **Step 2: Create initial schema and sqlc config**

```sql
create table projects (
  id text primary key,
  name text not null,
  description text,
  status text not null,
  created_at timestamptz not null,
  updated_at timestamptz not null
);
```

- [ ] **Step 3: Add core tables required before Mission**

Include at minimum:
- `projects`
- `repo_bindings`
- `executor_profiles`
- `roles`
- `agents`
- `agent_teams`
- `agent_team_members`

- [ ] **Step 4: Generate sqlc code and wrap repositories**

Run: `sqlc generate`  
Expected: generated query package under `backend/internal/db/sqlc`.

- [ ] **Step 5: Run migration and integration tests**

Run: `goose -dir db/migrations postgres "$DATABASE_URL" up`  
Expected: migration `000001` applies cleanly.

Run: `cd backend && go test ./internal/project/... -tags=integration`  
Expected: repository integration test passes.

- [ ] **Step 6: Commit**

```bash
git add db sqlc.yaml backend/internal/db backend/internal/project
git commit -m "feat: add core project schema"
```

### Task 4: Implement Mission Governance Domain And API

**Files:**
- Create: `db/migrations/000002_missions.sql`
- Create: `db/queries/missions.sql`
- Create: `backend/internal/mission/model.go`
- Create: `backend/internal/mission/service.go`
- Create: `backend/internal/mission/repository.go`
- Create: `backend/internal/mission/handler.go`
- Create: `backend/internal/mission/handler_test.go`
- Modify: `openapi/openapi.yaml`
- Modify: `backend/internal/http/router.go`
- Test: `backend/internal/mission/handler_test.go`

- [ ] **Step 1: Write the failing Mission creation handler test**

```go
func TestCreateMission(t *testing.T) {
	body := `{"title":"Build runtime board","teamId":"team_1"}`
	req := httptest.NewRequest(http.MethodPost, "/api/projects/proj_1/missions", strings.NewReader(body))
	rec := httptest.NewRecorder()

	router := testRouter(t)
	router.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
}
```

- [ ] **Step 2: Add Mission governance schema**

Include:
- `missions`
- `mission_members`
- `mission_decisions`
- `agent_admin_policies`

- [ ] **Step 3: Implement service rules**

```go
type Service interface {
	CreateMission(ctx context.Context, cmd CreateMissionCmd) (Mission, error)
	Decide(ctx context.Context, cmd DecideMissionCmd) (MissionDecision, error)
}
```

- [ ] **Step 4: Expose API endpoints and OpenAPI schema**

Endpoints:
- `POST /api/projects/{projectId}/missions`
- `GET /api/projects/{projectId}/missions/{missionId}`
- `POST /api/projects/{projectId}/missions/{missionId}/decisions`

- [ ] **Step 5: Run tests**

Run: `cd backend && go test ./internal/mission/... ./internal/http/...`  
Expected: Mission creation and decision tests pass.

- [ ] **Step 6: Commit**

```bash
git add db/ openapi/ backend/internal/mission backend/internal/http/router.go
git commit -m "feat: add mission governance api"
```

### Task 5: Implement Shared Documents And Discussion Sessions

**Files:**
- Create: `db/migrations/000003_documents_discussions.sql`
- Create: `db/queries/documents.sql`
- Create: `db/queries/discussions.sql`
- Create: `backend/internal/document/model.go`
- Create: `backend/internal/document/service.go`
- Create: `backend/internal/document/handler.go`
- Create: `backend/internal/discussion/model.go`
- Create: `backend/internal/discussion/service.go`
- Create: `backend/internal/discussion/handler.go`
- Create: `backend/internal/document/service_test.go`
- Modify: `openapi/openapi.yaml`
- Modify: `backend/internal/http/router.go`
- Test: `backend/internal/document/service_test.go`

- [ ] **Step 1: Write failing document adoption test**

```go
func TestAdoptDocumentVersionSetsSingleAdoptedVersion(t *testing.T) {
	svc := newDocumentService(t)
	doc := seedDocument(t, svc)
	v1 := seedVersion(t, svc, doc.ID)
	v2 := seedVersion(t, svc, doc.ID)

	err := svc.AdoptVersion(ctx, doc.ID, v2.ID, "admin_agent")
	require.NoError(t, err)
	requireVersionAdopted(t, doc.ID, v2.ID)
	requireVersionNotAdopted(t, doc.ID, v1.ID)
}
```

- [ ] **Step 2: Add document and discussion schema**

Include:
- `discussion_sessions`
- `discussion_rounds`
- `shared_documents`
- `shared_document_versions`

- [ ] **Step 3: Implement discussion -> document linkage**

```go
type CreateRoundCmd struct {
	SessionID   string
	Summary     string
	OpenQuestions []string
}
```

- [ ] **Step 4: Add handlers and OpenAPI definitions**

Endpoints:
- `POST /api/projects/{projectId}/missions/{missionId}/discussions`
- `POST /api/projects/{projectId}/missions/{missionId}/discussions/{sessionId}/rounds`
- `POST /api/projects/{projectId}/missions/{missionId}/documents`
- `POST /api/projects/{projectId}/missions/{missionId}/documents/{documentId}/adopt`

- [ ] **Step 5: Run tests**

Run: `cd backend && go test ./internal/document/... ./internal/discussion/...`  
Expected: adoption rule and discussion creation tests pass.

- [ ] **Step 6: Commit**

```bash
git add db/ openapi/ backend/internal/document backend/internal/discussion
git commit -m "feat: add shared documents and discussions"
```

### Task 6: Implement Task Board, Claims, Handoffs, And Review Checkpoints

**Files:**
- Create: `db/migrations/000004_tasks_reviews.sql`
- Create: `db/queries/tasks.sql`
- Create: `backend/internal/task/model.go`
- Create: `backend/internal/task/service.go`
- Create: `backend/internal/task/handler.go`
- Create: `backend/internal/task/service_test.go`
- Modify: `openapi/openapi.yaml`
- Modify: `backend/internal/http/router.go`
- Test: `backend/internal/task/service_test.go`

- [ ] **Step 1: Write failing handoff rule test**

```go
func TestCreateHandoffRequiresAdminFallbackWhenNoDownstream(t *testing.T) {
	svc := newTaskService(t)
	task := seedStandaloneTask(t, svc)

	_, err := svc.CreateHandoff(ctx, CreateHandoffCmd{
		TaskItemID: task.ID,
		FromAgentID: "agent_dev",
	})

	require.ErrorContains(t, err, "toAdminAgent")
}
```

- [ ] **Step 2: Add task schema**

Include:
- `task_boards`
- `task_items`
- `task_claims`
- `task_handoffs`
- `review_checkpoints`

- [ ] **Step 3: Implement assignment, claim, handoff, and checkpoint services**

```go
type TaskService interface {
	Claim(ctx context.Context, cmd ClaimTaskCmd) (TaskClaim, error)
	CreateHandoff(ctx context.Context, cmd CreateHandoffCmd) (TaskHandoff, error)
	RequestCheckpoint(ctx context.Context, cmd RequestCheckpointCmd) (ReviewCheckpoint, error)
}
```

- [ ] **Step 4: Expose task board endpoints**

Endpoints:
- `POST /api/projects/{projectId}/missions/{missionId}/task-board`
- `POST /api/projects/{projectId}/missions/{missionId}/tasks/{taskId}/claim`
- `POST /api/projects/{projectId}/missions/{missionId}/tasks/{taskId}/handoffs`
- `POST /api/projects/{projectId}/missions/{missionId}/tasks/{taskId}/review-checkpoints`

- [ ] **Step 5: Run tests**

Run: `cd backend && go test ./internal/task/...`  
Expected: claim, handoff fallback, and checkpoint tests pass.

- [ ] **Step 6: Commit**

```bash
git add db/ openapi/ backend/internal/task
git commit -m "feat: add task board and handoffs"
```

### Task 7: Implement Runtime Observability Schema And Event Pipeline

**Files:**
- Create: `db/migrations/000005_runtime_observability.sql`
- Create: `db/queries/runtime.sql`
- Create: `backend/internal/runtime/model.go`
- Create: `backend/internal/runtime/service.go`
- Create: `backend/internal/runtime/redactor.go`
- Create: `backend/internal/runtime/handler.go`
- Create: `backend/internal/runtime/service_test.go`
- Modify: `openapi/openapi.yaml`
- Modify: `backend/internal/http/router.go`
- Test: `backend/internal/runtime/service_test.go`

- [ ] **Step 1: Write failing transcript redaction test**

```go
func TestRedactorMasksKnownSecrets(t *testing.T) {
	got := Redact("token=abc123")
	require.NotContains(t, got, "abc123")
	require.Contains(t, got, "[REDACTED]")
}
```

- [ ] **Step 2: Add runtime schema**

Include:
- `agent_presences`
- `mission_agent_runtimes`
- `executor_sessions`
- `agent_runtime_events`
- `executor_transcripts`
- `executor_transcript_entries`
- `transcript_access_audits`
- `runtime_observability_policies`

- [ ] **Step 3: Implement runtime service**

```go
type RuntimeService interface {
	StartSession(ctx context.Context, cmd StartSessionCmd) (ExecutorSession, error)
	AppendEvent(ctx context.Context, cmd AppendEventCmd) error
	AppendTranscriptEntry(ctx context.Context, cmd AppendTranscriptEntryCmd) error
	SealSession(ctx context.Context, sessionID string) error
}
```

- [ ] **Step 4: Add query endpoints for runtime board and transcript viewer**

Endpoints:
- `GET /api/projects/{projectId}/missions/{missionId}/agent-runtimes`
- `GET /api/projects/{projectId}/executor-sessions/{sessionId}`
- `GET /api/projects/{projectId}/executor-sessions/{sessionId}/events`
- `GET /api/projects/{projectId}/executor-sessions/{sessionId}/transcript`

- [ ] **Step 5: Run tests**

Run: `cd backend && go test ./internal/runtime/...`  
Expected: redaction and session lifecycle tests pass.

- [ ] **Step 6: Commit**

```bash
git add db/ openapi/ backend/internal/runtime
git commit -m "feat: add runtime observability domain"
```

### Task 8: Implement CLI Executor Framework And Session Capture

**Files:**
- Create: `backend/internal/executor/adapter.go`
- Create: `backend/internal/executor/manager.go`
- Create: `backend/internal/executor/pty_runner.go`
- Create: `backend/internal/executor/codex/adapter.go`
- Create: `backend/internal/executor/claude/adapter.go`
- Create: `backend/internal/executor/gemini/adapter.go`
- Create: `backend/internal/executor/manager_test.go`
- Modify: `backend/internal/runtime/service.go`
- Test: `backend/internal/executor/manager_test.go`

- [ ] **Step 1: Write failing executor manager test**

```go
func TestManagerEmitsSessionStartedEvent(t *testing.T) {
	manager := newManager(t, fakeAdapter{})
	err := manager.RunTask(ctx, seedTaskCommand())
	require.NoError(t, err)
	requireEventRecorded(t, "session_started")
}
```

- [ ] **Step 2: Define adapter contract**

```go
type Adapter interface {
	Name() string
	BuildCommand(ctx context.Context, input TaskExecutionInput) (*exec.Cmd, error)
	MapOutput(line []byte) ([]runtime.EventDraft, []runtime.TranscriptDraft)
}
```

- [ ] **Step 3: Implement PTY-backed process runner**

Use:
- `github.com/creack/pty`
- line chunking for stdout/stderr
- heartbeat updates
- cancellation on context timeout

- [ ] **Step 4: Implement first-pass event mappers**

Support:
- `codex`
- `claude code`
- `gemini cli`

The first version may map only:
- session start
- prompt sent
- message created
- tool call started/finished
- session completed/failed

- [ ] **Step 5: Run tests**

Run: `cd backend && go test ./internal/executor/... ./internal/runtime/...`  
Expected: manager test and output mapping tests pass.

- [ ] **Step 6: Commit**

```bash
git add backend/internal/executor backend/internal/runtime
git commit -m "feat: add cli executor framework"
```

### Task 9: Implement Worker Scheduling And Run Orchestration

**Files:**
- Create: `db/migrations/000006_runs_workflow.sql`
- Create: `db/queries/runs.sql`
- Create: `backend/internal/worker/scheduler.go`
- Create: `backend/internal/worker/dispatcher.go`
- Create: `backend/internal/worker/run_service.go`
- Create: `backend/internal/worker/run_service_test.go`
- Modify: `backend/cmd/worker/main.go`
- Modify: `backend/internal/task/service.go`
- Test: `backend/internal/worker/run_service_test.go`

- [ ] **Step 1: Write failing worker orchestration test**

```go
func TestClaimedTaskCreatesRunAndExecutorSession(t *testing.T) {
	svc := newRunService(t)
	err := svc.ExecuteClaimedTask(ctx, "task_claim_1")
	require.NoError(t, err)
	requireRunCreated(t)
	requireExecutorSessionCreated(t)
}
```

- [ ] **Step 2: Add run schema**

Include:
- `runs`
- `node_runs`
- `repo_snapshots`
- `run_artifact_inputs`

- [ ] **Step 3: Implement claim -> run -> executor session flow**

```go
func (s *RunService) ExecuteClaimedTask(ctx context.Context, claimID string) error {
	// load task + inputs
	// create run
	// create executor session
	// delegate to adapter manager
}
```

- [ ] **Step 4: Emit SSE updates from worker**

Publish:
- runtime status changes
- task handoff updates
- review checkpoint updates

- [ ] **Step 5: Run tests**

Run: `cd backend && go test ./internal/worker/...`  
Expected: worker orchestration tests pass.

- [ ] **Step 6: Commit**

```bash
git add db/ backend/internal/worker backend/cmd/worker/main.go
git commit -m "feat: orchestrate worker runs"
```

### Task 10: Implement RepoCandidate, Approval, Publish, And Archive Services

**Files:**
- Create: `db/migrations/000007_repo_candidates_publish_archive.sql`
- Create: `db/queries/repo_candidates.sql`
- Create: `db/queries/approvals.sql`
- Create: `db/queries/archives.sql`
- Create: `backend/internal/approval/service.go`
- Create: `backend/internal/approval/handler.go`
- Create: `backend/internal/archive/service.go`
- Create: `backend/internal/archive/manifest.go`
- Create: `backend/internal/archive/service_test.go`
- Modify: `openapi/openapi.yaml`
- Modify: `backend/internal/http/router.go`
- Test: `backend/internal/archive/service_test.go`

- [ ] **Step 1: Write failing archive manifest test**

```go
func TestBuildArchiveManifestIncludesRuntimeSummary(t *testing.T) {
	svc := newArchiveService(t)
	manifest, err := svc.BuildManifest(ctx, "mission_1")
	require.NoError(t, err)
	require.NotEmpty(t, manifest.ObjectKeys)
	require.NotEmpty(t, manifest.IncludedDecisionIDs)
}
```

- [ ] **Step 2: Add publish and archive schema**

Include:
- `repo_candidates`
- `repo_candidate_publishes`
- `approvals`
- `artifact_publishes`
- `mission_archives`

- [ ] **Step 3: Implement intent snapshot validation**

```go
type IntentSnapshot struct {
	RepoBindingID string `json:"repoBindingId"`
	BaseCommitSHA string `json:"baseCommitSha"`
	TreeHash      string `json:"treeHash"`
}
```

- [ ] **Step 4: Implement archive assembly**

Bundle at minimum:
- adopted documents
- task board snapshot
- runtime event summary
- key transcript references
- approvals
- final summary

- [ ] **Step 5: Run tests**

Run: `cd backend && go test ./internal/approval/... ./internal/archive/...`  
Expected: intent validation and archive manifest tests pass.

- [ ] **Step 6: Commit**

```bash
git add db/ openapi/ backend/internal/approval backend/internal/archive
git commit -m "feat: add publish approval and archive services"
```

### Task 11: Build Frontend Shell, Routing, And Shared Data Layer

**Files:**
- Create: `frontend/src/app/router.tsx`
- Create: `frontend/src/app/query-client.ts`
- Create: `frontend/src/app/layout.tsx`
- Create: `frontend/src/components/page-shell.tsx`
- Create: `frontend/src/lib/api.ts`
- Create: `frontend/src/lib/sse.ts`
- Create: `frontend/src/pages/home.tsx`
- Create: `frontend/src/pages/mission.tsx`
- Create: `frontend/src/pages/runtime-session.tsx`
- Create: `frontend/src/styles.css`
- Modify: `frontend/src/App.tsx`
- Test: `frontend/e2e/app-shell.spec.ts`

- [ ] **Step 1: Write failing Playwright smoke test**

```ts
test("app shell shows mission navigation", async ({ page }) => {
  await page.goto("/");
  await expect(page.getByRole("navigation")).toBeVisible();
});
```

- [ ] **Step 2: Build app shell, router, and API client**

```ts
export const router = createBrowserRouter([
  { path: "/", element: <HomePage /> },
  { path: "/missions/:missionId", element: <MissionPage /> },
]);
```

- [ ] **Step 3: Add SSE client and shared fetch wrapper**

```ts
export function subscribe(path: string, onMessage: (event: MessageEvent) => void) {
  const source = new EventSource(path, { withCredentials: true });
  source.onmessage = onMessage;
  return () => source.close();
}
```

- [ ] **Step 4: Run typecheck and Playwright**

Run: `pnpm --dir frontend exec tsc --noEmit`  
Expected: TypeScript check passes.

Run: `pnpm --dir frontend exec playwright test frontend/e2e/app-shell.spec.ts`  
Expected: smoke test passes.

- [ ] **Step 5: Commit**

```bash
git add frontend
git commit -m "feat: add frontend shell and routing"
```

### Task 12: Build Mission Workspace Pages

**Files:**
- Create: `frontend/src/features/missions/mission-overview.tsx`
- Create: `frontend/src/features/discussions/discussion-session-panel.tsx`
- Create: `frontend/src/features/documents/document-list.tsx`
- Create: `frontend/src/features/documents/document-editor.tsx`
- Create: `frontend/src/features/tasks/task-board.tsx`
- Create: `frontend/src/features/tasks/task-detail-drawer.tsx`
- Create: `frontend/src/features/tasks/task-handoff-form.tsx`
- Create: `frontend/e2e/mission-workspace.spec.ts`
- Test: `frontend/e2e/mission-workspace.spec.ts`

- [ ] **Step 1: Write failing Mission workspace e2e**

```ts
test("mission page shows documents and task board", async ({ page }) => {
  await page.goto("/missions/mission_1");
  await expect(page.getByText("Task Board")).toBeVisible();
  await expect(page.getByText("Documents")).toBeVisible();
});
```

- [ ] **Step 2: Implement Mission page sections**

Sections:
- Mission header
- discussion rounds
- adopted documents
- task board with claims and handoffs
- review checkpoints

- [ ] **Step 3: Wire mutations**

Support:
- create discussion round
- adopt document
- claim task
- create handoff
- request review checkpoint

- [ ] **Step 4: Run tests**

Run: `pnpm --dir frontend exec playwright test frontend/e2e/mission-workspace.spec.ts`  
Expected: Mission workspace interactions pass.

- [ ] **Step 5: Commit**

```bash
git add frontend
git commit -m "feat: add mission workspace ui"
```

### Task 13: Build Agent Runtime Board And Transcript Inspector

**Files:**
- Create: `frontend/src/features/runtime/agent-runtime-board.tsx`
- Create: `frontend/src/features/runtime/agent-runtime-card.tsx`
- Create: `frontend/src/features/runtime/runtime-event-timeline.tsx`
- Create: `frontend/src/features/runtime/transcript-inspector.tsx`
- Create: `frontend/src/features/runtime/transcript-access-audit-table.tsx`
- Create: `frontend/e2e/runtime-board.spec.ts`
- Test: `frontend/e2e/runtime-board.spec.ts`

- [ ] **Step 1: Write failing runtime board e2e**

```ts
test("runtime board shows active agents and event timeline", async ({ page }) => {
  await page.goto("/missions/mission_1");
  await expect(page.getByText("Active Agents")).toBeVisible();
  await expect(page.getByText("Runtime Events")).toBeVisible();
});
```

- [ ] **Step 2: Implement runtime board**

Show:
- agent availability
- mission runtime status
- current executor backend
- last heartbeat
- latest events

- [ ] **Step 3: Implement transcript inspector**

Show:
- event timeline by default
- transcript collapsed by default
- redacted transcript on expand
- access audit table for privileged viewers

- [ ] **Step 4: Run tests**

Run: `pnpm --dir frontend exec playwright test frontend/e2e/runtime-board.spec.ts`  
Expected: runtime board and transcript inspector test passes.

- [ ] **Step 5: Commit**

```bash
git add frontend
git commit -m "feat: add runtime board and transcript inspector"
```

### Task 14: Add End-To-End Flow, CI, And Delivery Guardrails

**Files:**
- Create: `.github/workflows/ci.yml`
- Create: `frontend/e2e/mission-end-to-end.spec.ts`
- Create: `backend/scripts/seed_demo_data.go`
- Create: `backend/scripts/replay_executor_session.go`
- Modify: `README.md`
- Test: `.github/workflows/ci.yml`

- [ ] **Step 1: Write failing end-to-end scenario**

Scenario must cover:
- create Mission
- add discussion round
- adopt document
- claim task
- create executor session
- view runtime events
- complete Mission
- build archive

- [ ] **Step 2: Add demo seed and replay scripts**

```go
func main() {
	// insert demo project, agents, mission, documents, tasks, runtime events
}
```

- [ ] **Step 3: Add CI workflow**

Stages:
- backend unit tests
- backend integration tests
- frontend typecheck
- frontend Playwright

- [ ] **Step 4: Run the full verification suite**

Run: `cd backend && go test ./...`  
Expected: all backend tests pass.

Run: `pnpm --dir frontend exec tsc --noEmit`  
Expected: no type errors.

Run: `pnpm --dir frontend exec playwright test`  
Expected: e2e suite passes.

- [ ] **Step 5: Commit**

```bash
git add .github backend/scripts frontend/e2e README.md
git commit -m "chore: add ci and end-to-end coverage"
```

---

## API Surface Checklist

The implementation is not done until these endpoint groups exist and are backed by real data:

- Mission
  - `POST /api/projects/:projectId/missions`
  - `GET /api/projects/:projectId/missions/:missionId`
  - `POST /api/projects/:projectId/missions/:missionId/decisions`
- Discussion
  - `POST /api/projects/:projectId/missions/:missionId/discussions`
  - `POST /api/projects/:projectId/missions/:missionId/discussions/:sessionId/rounds`
- Documents
  - `POST /api/projects/:projectId/missions/:missionId/documents`
  - `POST /api/projects/:projectId/missions/:missionId/documents/:documentId/versions`
  - `POST /api/projects/:projectId/missions/:missionId/documents/:documentId/adopt`
- Tasks
  - `POST /api/projects/:projectId/missions/:missionId/task-board`
  - `POST /api/projects/:projectId/missions/:missionId/tasks/:taskId/claim`
  - `POST /api/projects/:projectId/missions/:missionId/tasks/:taskId/handoffs`
  - `POST /api/projects/:projectId/missions/:missionId/tasks/:taskId/review-checkpoints`
- Runtime
  - `GET /api/projects/:projectId/missions/:missionId/agent-runtimes`
  - `GET /api/projects/:projectId/executor-sessions/:sessionId`
  - `GET /api/projects/:projectId/executor-sessions/:sessionId/events`
  - `GET /api/projects/:projectId/executor-sessions/:sessionId/transcript`
  - `GET /api/projects/:projectId/executor-sessions/:sessionId/access-audits`
- Approval / Publish / Archive
  - `GET /api/projects/:projectId/approvals/:approvalId`
  - `POST /api/projects/:projectId/repo-candidates/:candidateId/publish`
  - `POST /api/projects/:projectId/missions/:missionId/archive`
  - `GET /api/projects/:projectId/missions/:missionId/archive`

## Verification Checklist

- `docker compose -f ops/docker-compose.yml up -d` can bring up PostgreSQL and MinIO.
- `goose up` applies all migrations on a clean database.
- `sqlc generate` succeeds after each schema change.
- `go test ./...` passes under `backend/`.
- `pnpm --dir frontend exec tsc --noEmit` passes.
- `pnpm --dir frontend exec playwright test` passes.
- A demo Mission can progress from discussion to archive with a visible Agent runtime session and transcript.
- At least one `codex` executor session can emit structured runtime events and a redacted transcript.
- Transcript access creates audit records.
- Archive manifest includes runtime event summary and key execution session references.

## Risks To Watch During Execution

- CLI output formats may differ more than expected; keep adapter parsing minimal and event mapping tolerant.
- PTY buffering can duplicate or split lines; build transcript chunking tests early.
- SSE fan-out can leak connections; cap subscribers and clean up on disconnect.
- sqlc query sprawl can make repositories noisy; keep one domain package per query group.
- Frontend can drift into mock data; connect real API responses by Task 11, not later.
- Archive bundling can get large; store large transcript blobs in object storage instead of DB text columns.

## Out Of Scope For This Plan

- Enterprise SSO / SCIM / ABAC
- Kubernetes-native sandboxed executors
- Multi-region deployment
- Rich plugin marketplace
- Full GitHub/GitLab app installation flow

## Handoff Notes

- Start with Task 1 and do not parallelize backend schema work until Task 3 is committed.
- Frontend implementation should not start before Task 6 APIs exist, except for app shell scaffolding.
- Runtime board UI should wait for Task 7 and Task 8 contracts to stabilize.
- If transcript privacy requirements change, update `RuntimeObservabilityPolicy` before shipping UI export features.
