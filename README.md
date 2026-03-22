# agent-dev

多 Agent 自动化开发平台方案与实现仓库。

当前仓库包含：

- [v1 原始方案](./docs/multi-agent-platform-final-plan.md)
- [v1.1 方法补强文档](./docs/multi-agent-platform-v1.1-method.md)
- [v2 协作治理方案](./docs/multi-agent-platform-v2-collaboration-governance-plan.md)
- [v2 落地实施文档](./docs/multi-agent-platform-v2-implementation-plan.md)
- `backend/`：Go API 与 Worker 代码
- `frontend/`：React 工作台
- `ops/docker-compose.yml`：本地 PostgreSQL 与 MinIO

当前开发建议：

1. 先看 v1 原始方案，理解产品边界和技术主张。
2. 再看 v1.1 方法文档，理解执行器边界、代码中间态、审批快照和 Run 输入固化等落地约束。
3. 使用 `make dev-up` 启动本地依赖。
4. 使用 `make test-backend` 运行后端测试。
5. 进入 `frontend/` 后使用 `corepack pnpm typecheck` 和 `corepack pnpm test:e2e` 验证前端。

## 快速启动

### 1. 启动本地依赖

在仓库根目录执行：

```bash
make dev-up
```

这会启动：

- PostgreSQL：`localhost:5432`
- MinIO：
  - API：`http://localhost:9000`
  - Console：`http://localhost:9001`

### 2. 配置环境变量

后端默认依赖下面这些环境变量：

```bash
export APP_ENV=development
export APP_ADDR=:8080
export DATABASE_URL=postgres://postgres:postgres@localhost:5432/agent_platform?sslmode=disable
export S3_ENDPOINT=http://localhost:9000
export S3_ACCESS_KEY=minio
export S3_SECRET_KEY=minio123
export S3_BUCKET=agent-platform-dev
```

也可以直接参考仓库根目录的 `.env.example`。

### 3. 启动后端 API

在一个终端里执行：

```bash
cd backend
go run ./cmd/api
```

默认监听：

- API：`http://127.0.0.1:8080`
- Health：`http://127.0.0.1:8080/healthz`

### 4. 启动 Worker

在另一个终端里执行：

```bash
cd backend
go run ./cmd/worker
```

Worker 会轮询 active claim，并触发任务执行链路。

### 5. 启动前端

在第三个终端里执行：

```bash
cd frontend
corepack pnpm install
corepack pnpm dev
```

默认访问：

- 前端：`http://127.0.0.1:5173`

### 6. 写入演示数据

在任意终端里执行：

```bash
cd backend
DATABASE_URL=postgres://postgres:postgres@localhost:5432/agent_platform?sslmode=disable go run ./scripts/seed_demo_data
```

这会写入一套演示数据：

- Project：`proj_1`
- Mission：`mission_1`

### 7. 生成一条新的 runtime session 回放

```bash
cd backend
DATABASE_URL=postgres://postgres:postgres@localhost:5432/agent_platform?sslmode=disable go run ./scripts/replay_executor_session
```

这会生成一条新的 executor session，方便你验证 runtime board、transcript 和 SSE 刷新。

## 启动后建议验证

### 后端验证

```bash
cd backend
go test ./...
```

### 前端类型检查

```bash
cd frontend
corepack pnpm typecheck
```

### 前端端到端

```bash
cd frontend
corepack pnpm test:e2e
```

### 手工访问建议

启动前后端并 seed 数据后，可以先看：

- `http://127.0.0.1:5173/`
- `http://127.0.0.1:5173/missions/mission_1`
- `http://127.0.0.1:5173/runtime/sessions/session_1`

本地演示脚本：

1. `cd backend && DATABASE_URL=postgres://postgres:postgres@localhost:5432/agent_platform?sslmode=disable go run ./scripts/seed_demo_data`
2. `cd backend && DATABASE_URL=postgres://postgres:postgres@localhost:5432/agent_platform?sslmode=disable go run ./scripts/replay_executor_session`

这两个脚本会写入一套 `proj_1 / mission_1` 的演示数据，并生成一条新的 runtime executor session 供调试和 UI 联调。
