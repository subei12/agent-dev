import { expect, test } from "@playwright/test";

test("mission workspace shows internal execution tasks and expands transcript only on demand", async ({
  page
}) => {
  const taskBoard = {
    id: "board_1",
    title: "交付看板",
    items: [
      {
        id: "task_runtime",
        title: "实现运行态观测",
        type: "code",
        status: "in_progress",
        assignedAgentId: "agent_backend"
      }
    ]
  };

  const runtimes = [
    {
      id: "runtime_backend",
      agentId: "agent_backend",
      status: "working",
      statusSummary: "正在实现运行态观测",
      currentExecutorSessionId: "session_backend"
    },
    {
      id: "runtime_reviewer",
      agentId: "agent_reviewer",
      status: "waiting_admin",
      statusSummary: "等待管理员确认评审结论",
      currentExecutorSessionId: "session_reviewer"
    }
  ];

  const sessionEvents = {
    session_backend: [
      {
        id: "event_backend_1",
        title: "开始编码",
        type: "message_created",
        category: "model",
        summary: "后端 Agent 正在写运行态接口。"
      }
    ],
    session_reviewer: [
      {
        id: "event_reviewer_1",
        title: "等待管理员",
        type: "waiting_admin",
        category: "handoff",
        summary: "评审 Agent 已提交检查意见。"
      }
    ]
  } as const;

  const sessionTranscripts = {
    session_backend: {
      transcript: {
        id: "transcript_backend",
        status: "sealed",
        executorSessionId: "session_backend"
      },
      entries: [
        {
          id: "entry_backend_1",
          role: "assistant",
          entryType: "message",
          redactedText: "已生成运行态接口草稿。"
        }
      ]
    },
    session_reviewer: {
      transcript: {
        id: "transcript_reviewer",
        status: "sealed",
        executorSessionId: "session_reviewer"
      },
      entries: [
        {
          id: "entry_reviewer_1",
          role: "assistant",
          entryType: "message",
          redactedText: "建议管理员复核 transcript 脱敏策略。"
        }
      ]
    }
  } as const;

  await page.route("**/api/projects/proj_1/missions/mission_1", async (route) => {
    await route.fulfill({
      json: {
        id: "mission_1",
        title: "运行态工作台联调",
        status: "implementation",
        description: "验证任务新增、Agent 执行区和 transcript 展开流程。"
      }
    });
  });

  await page.route("**/api/projects/proj_1/missions/mission_1/documents", async (route) => {
    await route.fulfill({
      json: [
        {
          id: "doc_1",
          title: "运行态接口说明",
          kind: "architecture",
          currentAdoptedVersionId: "docv_1"
        }
      ]
    });
  });

  await page.route("**/api/projects/proj_1/missions/mission_1/documents/doc_1/versions", async (route) => {
    await route.fulfill({
      json: [
        {
          id: "docv_1",
          documentId: "doc_1",
          version: 1,
          status: "adopted",
          contentFormat: "md",
          contentHash: "hash-docv1",
          contentText: "当前采用事件优先展示。"
        }
      ]
    });
  });

  await page.route("**/api/projects/proj_1/missions/mission_1/discussions", async (route) => {
    await route.fulfill({
      json: [
        {
          id: "discussion_1",
          topic: "事件优先还是 transcript 优先",
          status: "open"
        }
      ]
    });
  });

  await page.route("**/api/projects/proj_1/missions/mission_1/approvals", async (route) => {
    await route.fulfill({ json: [] });
  });

  await page.route("**/api/projects/proj_1/missions/mission_1/archive", async (route) => {
    await route.fulfill({
      json: {
        id: "archive_1",
        missionId: "mission_1",
        status: "uploaded"
      }
    });
  });

  await page.route("**/api/projects/proj_1/missions/mission_1/task-board", async (route) => {
    await route.fulfill({ json: taskBoard });
  });

  await page.route("**/api/projects/proj_1/missions/mission_1/agent-runtimes", async (route) => {
    await route.fulfill({ json: runtimes });
  });

  await page.route("**/api/projects/proj_1/agents", async (route) => {
    await route.fulfill({
      json: [
        {
          id: "agent_backend",
          name: "后端 Agent",
          enabled: true,
          executorProfileId: "exec_backend",
          executorType: "codex_cli",
          command: "codex",
          args: ["exec"],
          timeoutSec: 3600,
          maxConcurrency: 1,
          allowRepoRead: true,
          allowRepoWrite: true,
          allowNetwork: true
        },
        {
          id: "agent_reviewer",
          name: "评审 Agent",
          enabled: true,
          executorProfileId: "exec_reviewer",
          executorType: "claude_code",
          command: "claude",
          args: ["--continue"],
          timeoutSec: 3600,
          maxConcurrency: 1,
          allowRepoRead: true,
          allowRepoWrite: false,
          allowNetwork: true
        }
      ]
    });
  });

  await page.route("**/api/projects/proj_1/executor-sessions/*/events", async (route) => {
    const sessionId = route.request().url().split("/executor-sessions/")[1].split("/events")[0] as keyof typeof sessionEvents;
    await route.fulfill({ json: sessionEvents[sessionId] ?? [] });
  });

  await page.route("**/api/projects/proj_1/executor-sessions/*/transcript?view=redacted", async (route) => {
    const sessionId = route.request().url().split("/executor-sessions/")[1].split("/transcript")[0] as keyof typeof sessionTranscripts;
    await route.fulfill({ json: sessionTranscripts[sessionId] });
  });

  await page.route("**/api/projects/proj_1/executor-sessions/*/access-audits", async (route) => {
    await route.fulfill({
      json: [
        {
          id: "audit_1",
          actorUserId: "demo-admin",
          accessMode: "redacted_transcript",
          reason: "排查运行日志"
        }
      ]
    });
  });

  await page.goto("/missions/mission_1");

  await expect(page.getByText("内部执行任务", { exact: true })).toBeVisible();
  await expect(page.getByText("这些任务由管理员 Agent 拆解后分派给不同执行 Agent。")).toBeVisible();
  await expect(page.getByRole("heading", { name: "实现运行态观测" }).first()).toBeVisible();

  await page.getByRole("button", { name: /评审 Agent/ }).click();

  await expect(page.getByText("当前查看：评审 Agent")).toBeVisible();
  await expect(page.getByText("状态：等待管理员 · 会话：session_reviewer")).toBeVisible();
  await expect(page.getByText("建议管理员复核 transcript 脱敏策略。")).toHaveCount(0);

  await page.getByRole("button", { name: "展开原始 transcript" }).click();
  await expect(page.getByText("建议管理员复核 transcript 脱敏策略。")).toBeVisible();
});
