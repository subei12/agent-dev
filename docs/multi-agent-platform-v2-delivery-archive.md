# 多 Agent 平台 v2 交付归档

## 分支信息

- 分支：`v2-bootstrap`
- 远端：`origin/v2-bootstrap`

## 已完成能力

- 已完成 Go 后端、React 前端、Docker Compose 本地环境和 CI 的项目骨架。
- 已完成包含 discussion、document、task board、runtime、approval、archive 的 Mission 治理 API。
- 已完成共享文档的版本创建与采纳流程。
- 已完成任务 claim、handoff 和 review checkpoint 流程。
- 已完成运行时观测，包括 session 生命周期、结构化事件、transcript 捕获、审计和脱敏存储。
- 已完成 `codex_cli`、`claude_code`、`gemini_cli` 的执行器框架。
- 已完成包含 claim 轮询、执行锁、过期锁回收、心跳更新和失败重试的 Worker 编排。
- 已完成 SSE 事件流和前端基于查询失效的实时刷新。
- 已完成 manifest 与 bundle 的对象存储写入。
- 已完成前端 Mission 工作台的写操作，包括 discussions、documents、task claims、handoffs、checkpoints 和 archive。
- 已完成演示 seed 脚本和 executor replay 脚本。
- 已完成针对 Mission 壳、Runtime 壳、seeded 数据和 Mission 操作的 Playwright 端到端覆盖。

## 仍未完成

- 真实的用户认证与会话管理。
- 更完整的 RBAC，或可在 UI 中编辑角色策略的能力。
- 超出当前共享适配基线的、更细粒度的执行器专属事件解析。
- 更生产级的 Worker 队列语义，例如死信处理和并发池。
- 超出当前 zip/json 基线的、更完整的大对象归档打包格式。
- 面向审批审核和 mission/archive 历史浏览的前端页面。

## 验证快照

- `cd backend && go test ./...`
- `cd frontend && corepack pnpm typecheck`
- `cd frontend && corepack pnpm test:e2e`

在更新这份归档时，上述命令均已通过。
