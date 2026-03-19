# 多 Agent 自动化开发平台 v1.1 方法文档

## 1. 文档定位

本文是对 [multi-agent-platform-final-plan.md](./multi-agent-platform-final-plan.md) 的补强，不推翻原始产品方向，只补齐 v1 在执行语义、审批固化、代码流转和数据闭环上的缺口。

v1.1 的目标不是“加更多功能”，而是把下面 5 件事定义清楚并能稳定落地：

1. 执行器权限到底是“真实隔离”还是“可信约束”。
2. 多个 `code_task` 之间如何在不写 Git 的前提下继续传递代码结果。
3. 审批通过时，平台到底批准了哪个不可变对象。
4. Run 启动时 pin 住的 Artifact 输入如何持久化并可复现。
5. Artifact 生命周期状态如何统一，避免模型与查询语义冲突。

---

## 2. v1.1 核心结论

v1.1 维持原始方案的三条主轴不变：

- 平台仍然是知识与流程主系统。
- Git 仍然只是代码上下文源和可选发布目标。
- 所有外部副作用仍然必须显式化并进入审批链。

在此基础上，v1.1 增加 5 条执行方法约束：

1. `ExecutorProfile` 在 v1.1 采用“可信执行器模型”，不把本地 CLI 误描述为强沙箱。
2. `code_task` 不直接把“未发布代码结果”塞回 Git，而是产出不可变 `RepoCandidate`。
3. `Run` 在启动时必须显式固化 `RunArtifactInput` 和 `RepoSnapshot`。
4. 审批必须绑定 `intentSnapshot` 和 `intentHash`，审批对象不可变。
5. `Artifact` 的“平台内/已发布/发布失败”属于派生分发状态，不再和核心版本状态混用。

---

## 3. v1.1 要解决的具体问题

### 3.1 执行器权限描述失真

如果 v1 直接在宿主机拉起 Codex CLI / Claude Code CLI / Gemini CLI，那么：

- `allowNetwork`
- `allowRepoWrite`
- `canWriteRepo`

这些字段只能表达“平台调度层的允许策略”，不能天然等于“操作系统级安全隔离”。

因此 v1.1 必须明确：

- 本地 CLI 执行器默认是 `trusted_host`。
- 平台只承诺流程控制、凭据注入控制、发布入口控制。
- 平台不承诺该模式下的强网络隔离和强文件系统隔离。

### 3.2 代码结果无法在节点间连续流转

如果下游节点只能读取：

- `ArtifactVersion`
- `RepoSnapshot(commit_sha)`

那么上游 `code_task` 生成的未发布代码结果只能停留在 `diff` 文本层，无法继续作为真实代码工作区被下游消费。

这会直接影响以下典型链路：

1. 生成代码
2. 跑测试
3. 自动修复
4. 代码审查
5. 决定是否发布

因此 v1.1 必须让“代码中间态”成为平台一等公民。

### 3.3 审批与实际执行对象可能不一致

如果审批只记录 `action=repo_write/open_pr`，而不记录：

- base commit
- target branch
- target path
- payload hash
- source candidate

那么审批通过后，真正执行的对象可能已经变化，平台无法证明“用户批准的就是最终执行的那个对象”。

### 3.4 Run 输入缺少持久化关联

v1 已经要求运行时 pin 住固定 ArtifactVersion，但当前对象模型没有对应关联表。这样会导致：

- 重跑不可复现
- 审计无法准确回答“这次运行读了哪些输入”
- 节点上下文组装依赖隐式逻辑

### 3.5 Artifact 状态语义冲突

当前方案同时存在：

- `Artifact.sourceOfTruth`
- `ArtifactVersion.status`
- `platform_only / published_to_repo / publish_failed`

如果不拆开“版本状态”和“分发状态”，后续查询、筛选、列表展示、发布回写都会出现定义冲突。

---

## 4. v1.1 核心方法

### 4.1 执行器采用可信执行器模型

v1.1 明确定义两种执行器安全模式：

```ts
type ExecutorSecurityMode = "trusted_host" | "container_sandbox"
```

`v1.1` 只强制实现：

```ts
ExecutorProfile {
  id: string
  projectId: string
  name: string
  type: "codex_cli" | "claude_code" | "gemini_cli" | "dummy"
  securityMode: "trusted_host"
  command: string
  args: string[]
  envSecretRefs: string[]
  timeoutSec: number
  maxConcurrency: number
  maxPromptChars: number
  allowRepoRead: boolean
  allowRepoWrite: boolean
  allowNetwork: boolean
}
```

方法约束如下：

1. `trusted_host` 只表示“平台可信地调度本地 CLI”，不表示强隔离。
2. `allowRepoWrite=false` 时，不为该节点挂载可推送凭据，不暴露发布入口，不提供可写目标仓库工作目录。
3. `allowNetwork=false` 时，只能作为平台策略和执行前校验字段，文档和 UI 必须明确这不是宿主机级硬隔离。
4. 如果未来要提供强安全承诺，必须新增 `container_sandbox` 执行器，再单独声明可验证的隔离边界。

结论：v1.1 选择“真实描述能力”，而不是“表面上有权限字段，实际上无法约束”。

### 4.2 引入 RepoCandidate 作为代码中间态主对象

`RepoCandidate` 用来表示“基于某个 repo 快照或上一个 candidate 演化出来的一份不可变代码结果”。

```ts
RepoCandidate {
  id: string
  projectId: string
  repoBindingId: string
  baseCommitSha: string
  parentCandidateId?: string
  createdFrom: "repo_snapshot" | "repo_candidate"
  createdFromRefId: string
  treeHash: string
  patchObjectKey?: string
  workspaceArchiveObjectKey?: string
  summary?: string
  status: "draft" | "approved_for_publish" | "published" | "rejected" | "superseded"
  producedByRunId: string
  producedByNodeRunId: string
  createdAt: string
}
```

设计原则：

1. `RepoCandidate` 不可变。
2. 每个 `code_task` 成功后都产出新的 `RepoCandidate`，而不是覆写旧结果。
3. 下游 `code_task` 可以读取 `RepoSnapshot`，也可以读取 `RepoCandidate`。
4. “代码还没进 Git，但已经可以被后续节点继续使用”这件事由平台负责，不由 Git 负责。

### 4.3 Run 输入显式固化

v1.1 不再把运行输入停留在模糊的 `inputsJson` 语义上，而是拆成显式对象。

```ts
Run {
  id: string
  projectId: string
  workflowVersionId: string
  status: "pending" | "running" | "paused" | "waiting_approval" | "failed" | "succeeded" | "cancelled"
  inputsJson: Json
  repoSnapshotId?: string
  createdBy: string
  createdAt: string
  startedAt?: string
  finishedAt?: string
}

RepoSnapshot {
  id: string
  projectId: string
  repoBindingId: string
  commitSha: string
  refName?: string
  createdAt: string
}

RunArtifactInput {
  runId: string
  artifactId: string
  artifactVersionId: string
  alias?: string
  createdAt: string
}
```

方法约束如下：

1. Run 创建时必须把选中的 `ArtifactVersion` 落库到 `RunArtifactInput`。
2. 如果本次运行依赖代码上下文，必须先创建 `RepoSnapshot`，再把 `repoSnapshotId` 写入 `Run`。
3. Retry 和 Replay 默认复用原始 `RunArtifactInput + RepoSnapshot`，除非用户显式发起新 Run。
4. `ContextAssembler` 只从持久化输入读数据，不从“当前最新 Artifact”做隐式漂移。

### 4.4 审批采用 Intent Snapshot

审批不再只看“动作名称”，而是必须绑定一个不可变审批意图快照。

```ts
Approval {
  id: string
  projectId: string
  runId?: string
  nodeRunId?: string
  action: "artifact_publish" | "repo_candidate_publish" | "open_pr" | "merge" | "release" | "external_write"
  subjectType: "artifact_publish" | "repo_candidate_publish" | "merge_request"
  subjectId: string
  intentSnapshotJson: Json
  intentHash: string
  status: "pending" | "approved" | "rejected"
  comment?: string
  createdBy: string
  decidedBy?: string
  createdAt: string
  decidedAt?: string
}
```

`intentSnapshotJson` 至少包含：

- `repoBindingId`
- `baseCommitSha`
- `targetBranch`
- `targetPath`
- `publishMode`
- `sourceArtifactVersionId` 或 `sourceRepoCandidateId`
- `contentHash` 或 `treeHash`

方法约束如下：

1. 审批通过后，Worker 只能执行与 `intentHash` 完全一致的对象。
2. 如果分支头、内容哈希、目标路径发生变化，原审批自动失效，必须重新申请。
3. 审计日志必须能从 `Approval -> intentSnapshot -> executed publish record` 全链路回放。

### 4.5 统一 Artifact 生命周期语义

v1.1 把状态拆成两层：

1. 核心版本状态
2. 分发状态

核心对象保留：

```ts
Artifact {
  id: string
  projectId: string
  kind: "prd" | "arch" | "contract" | "testcase" | "prompt" | "report" | "diff" | "note" | "other"
  title: string
  sourceOfTruth: "platform" | "repo"
  latestVersion: number
  createdAt: string
  updatedAt: string
}

ArtifactVersion {
  id: string
  artifactId: string
  version: number
  contentFormat: "md" | "yaml" | "json" | "txt" | "diff" | "pdf" | "png" | "zip"
  storageKind: "db_text" | "object_store" | "repo_ref"
  status: "draft" | "final" | "archived"
  hash: string
  sizeBytes: number
  producedByRunId?: string
  producedByNodeKey?: string
  createdAt: string
}
```

分发状态改为派生视图：

```ts
type ArtifactDistributionState =
  | "platform_only"
  | "linked_from_repo"
  | "published_to_repo"
  | "publish_failed"
  | "archived"
```

派生规则：

1. `sourceOfTruth=repo` 且无平台发布记录，显示 `linked_from_repo`。
2. 有成功发布记录，显示 `published_to_repo`。
3. 最近一次发布失败，显示 `publish_failed`。
4. `ArtifactVersion.status=archived` 时，在 UI 上显示 `archived`。

结论：`ArtifactVersion.status` 管版本生命周期，`ArtifactDistributionState` 管分发结果，二者不混用。

### 4.6 NodeContext 升级为双源代码输入

```ts
type RepoInputRef =
  | {
      kind: "snapshot"
      repoBindingId: string
      repoSnapshotId: string
      commitSha: string
      paths?: string[]
    }
  | {
      kind: "candidate"
      repoBindingId: string
      repoCandidateId: string
      baseCommitSha: string
      treeHash: string
    }

NodeContext {
  runInputs: Json
  selectedArtifactVersions: ArtifactVersion[]
  repoInput?: RepoInputRef
  rolePrompt: string
  nodePromptOverride?: string
}
```

方法约束如下：

1. 首个 `code_task` 默认读取 `RepoSnapshot`。
2. 后续 `code_task` 可以显式声明读取某个上游 `RepoCandidate`。
3. 普通 `task` 节点仍然主要读取 `ArtifactVersion`，不直接消费工作区目录。

### 4.7 节点输出规则调整

v1.1 对节点输出做以下收敛：

- `task`：产出 `ArtifactVersion`
- `code_task`：产出 `ArtifactVersion`、`test_report`、可选 `RepoCandidate`
- `approval`：只产出审批决定，不产出外部副作用
- `publish`：只消费 `ArtifactVersion` 或 `RepoCandidate`，由平台执行外部发布

关键约束：

1. 代码节点不能自行 push。
2. 代码节点可以形成新的 `RepoCandidate`。
3. 只有 `publish` 节点或显式用户动作，才能把 `ArtifactVersion` / `RepoCandidate` 写向 Git。

### 4.8 重试、重跑与幂等语义

v1.1 采用以下语义：

1. `NodeRun` 重试不会覆写原 `RepoCandidate`，而是产生新的 candidate 或新的失败记录。
2. 任何需要审批的 publish 行为都以 `intentHash` 为幂等键。
3. Retry 只允许复用同一组 `RunArtifactInput + RepoSnapshot + intentSnapshot`。
4. 如果输入或发布意图变化，应创建新的 Run 或新的 Approval，而不是重用旧记录。

---

## 5. v1.1 最小数据模型增量

在 v1 的表基础上，v1.1 最少增加或调整以下结构：

### 5.1 新增表

- `repo_snapshots`
- `run_artifact_inputs`
- `repo_candidates`
- `repo_candidate_publishes`

### 5.2 调整表

- `runs`：增加 `repo_snapshot_id`
- `approvals`：增加 `subject_type`、`subject_id`、`intent_snapshot_json`、`intent_hash`
- `artifact_publishes`：增加 `content_hash`、`base_commit_sha`、`intent_hash`

### 5.3 不新增专门表的对象

`ArtifactDistributionState` 作为查询层派生字段，不落单独表，不单独做主状态字段。

---

## 6. v1.1 API 增量

### 6.1 Run 创建

```http
POST /api/projects/:projectId/runs
```

请求体至少支持：

```json
{
  "workflowVersionId": "wfver_xxx",
  "inputs": {},
  "artifactVersionIds": ["av_1", "av_2"],
  "repoBindingId": "repo_1",
  "repoRef": "main"
}
```

服务端动作：

1. 固化 `RepoSnapshot`
2. 写入 `Run`
3. 写入 `RunArtifactInput`

### 6.2 RepoCandidate 查询与发布

```http
GET  /api/projects/:projectId/repo-candidates/:candidateId
POST /api/projects/:projectId/repo-candidates/:candidateId/publish
```

### 6.3 Run 输入回放

```http
GET /api/projects/:projectId/runs/:runId/inputs
```

返回：

- `inputsJson`
- `repoSnapshot`
- `artifactInputs`

### 6.4 Approval 查看

```http
GET /api/projects/:projectId/approvals/:approvalId
```

返回：

- `subjectType`
- `subjectId`
- `intentSnapshotJson`
- `intentHash`

---

## 7. v1.1 关键执行流程

### 7.1 文档型工作流

1. 用户发起 Run。
2. 平台固定 `RunArtifactInput`。
3. `task` 节点读取固定 ArtifactVersion。
4. 节点输出写入新的 `ArtifactVersion`。
5. 如需外部发布，创建 `ArtifactPublish + Approval(intentSnapshot)`。
6. 审批通过后，由 `publish` 节点执行 Git 写入。

### 7.2 代码型工作流

1. 用户发起 Run 并指定仓库。
2. 平台先创建 `RepoSnapshot(commitSha)`。
3. 第一个 `code_task` 基于 `RepoSnapshot` 物化工作区。
4. 节点完成后产出 `RepoCandidate#1`、`diff`、`test_report`。
5. 第二个 `code_task` 读取 `RepoCandidate#1`，继续修改后产出 `RepoCandidate#2`。
6. review 节点读取 `RepoCandidate#2` 对应的 diff、报告和摘要。
7. 用户或工作流触发 `publish` 节点。
8. 平台生成 `RepoCandidatePublish + Approval(intentSnapshot)`。
9. 审批通过后，由 Worker 将该 candidate 写入目标 branch 或 PR。

### 7.3 审批执行流程

1. 平台创建待审批对象。
2. 同步写入 `intentSnapshotJson` 和 `intentHash`。
3. 审批人看到的是固定快照，而不是动态对象。
4. Worker 执行前再次校验 `intentHash`。
5. 校验失败则拒绝执行并要求重新审批。

---

## 8. v1.1 实施顺序

建议按以下顺序落地，避免大面积返工：

### Phase 1

- 增加 `RepoSnapshot`
- 增加 `RunArtifactInput`
- 调整 `ContextAssembler`
- 跑通“固定输入可复现”

### Phase 2

- 增加 `RepoCandidate`
- 改造 `code_task` 输入输出
- 跑通“多代码节点串联但不写 Git”

### Phase 3

- 改造 `Approval`
- 引入 `intentSnapshotJson + intentHash`
- 跑通“审批对象不可变”

### Phase 4

- 补 `RepoCandidatePublish`
- 补 Artifact 分发状态派生逻辑
- 完成 Runs / Approvals / Publish UI 对应展示

---

## 9. v1.1 验收标准

满足以下条件即可认为 v1.1 方法闭环成立：

1. 同一个 Run 可以明确查到全部 `ArtifactVersion` 输入和固定的 `RepoSnapshot`。
2. 两个以上 `code_task` 可以串联执行，且中间结果不需要先写 Git。
3. 任意一次发布审批都能回放出“批准的具体对象、具体 commit 基线、具体内容哈希”。
4. Worker 执行外部写操作前，可以校验审批快照未漂移。
5. UI 和 API 对 Artifact 的显示不再混淆“版本状态”和“分发状态”。
6. 平台文档明确说明 `trusted_host` 不是强安全沙箱。

---

## 10. 与原始方案的关系

这份 v1.1 方法文档不改变原方案的产品定义，只把原方案从“方向正确”推进到“可以稳定实现并减少返工”。

原始方案继续成立的部分：

- 平台优先，Git 次级
- Artifact 是长期主数据
- 外部副作用显式审批
- 工作流版本化运行
- 单体 API + 内置 Scheduler + Worker

v1.1 新增的只是实现层闭环：

- 可信执行器边界
- 代码中间态对象
- Run 输入固化
- 审批快照化
- Artifact 状态语义统一

结论：v1 定义产品，v1.1 定义方法；两者应同时存在，并以前者为产品总纲、以后者为落地约束。
