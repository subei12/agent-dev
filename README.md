# agent-dev

多 Agent 自动化开发平台方案与实现仓库。

当前仓库包含：

- [v1 原始方案](./docs/multi-agent-platform-final-plan.md)
- [v1.1 方法补强文档](./docs/multi-agent-platform-v1.1-method.md)
- `backend/`：Go API 与 Worker 代码
- `frontend/`：React 工作台
- `ops/docker-compose.yml`：本地 PostgreSQL 与 MinIO

当前开发建议：

1. 先看 v1 原始方案，理解产品边界和技术主张。
2. 再看 v1.1 方法文档，理解执行器边界、代码中间态、审批快照和 Run 输入固化等落地约束。
3. 使用 `make dev-up` 启动本地依赖。
4. 使用 `make test-backend` 运行后端测试。
