import { expect, test } from "@playwright/test";

test("mission workspace renders seeded api data and runtime transcript", async ({ page }) => {
  await page.route("**/api/projects/proj_1/missions/mission_1", async (route) => {
    await route.fulfill({
      json: {
        id: "mission_1",
        title: "演示任务",
        status: "implementation",
        description: "交付运行态观测工作台。"
      }
    });
  });

  await page.route("**/api/projects/proj_1/missions/mission_1/documents", async (route) => {
    await route.fulfill({
      json: [
        {
          id: "doc_1",
          title: "架构快照",
          kind: "architecture",
          currentAdoptedVersionId: "docv_2"
        }
      ]
    });
  });

  await page.route("**/api/projects/proj_1/missions/mission_1/discussions", async (route) => {
    await route.fulfill({
      json: [
        {
          id: "session_1",
          topic: "运行态观测范围",
          status: "open"
        }
      ]
    });
  });

  await page.route("**/api/projects/proj_1/missions/mission_1/agent-runtimes", async (route) => {
    await route.fulfill({
      json: [
        {
          id: "runtime_1",
          agentId: "agent_backend",
          status: "working",
          statusSummary: "正在输出 codex 会话",
          currentExecutorSessionId: "session_telemetry"
        }
      ]
    });
  });

  await page.route("**/api/projects/proj_1/missions/mission_1/task-board", async (route) => {
    await route.fulfill({
      json: {
        id: "board_1",
        title: "Delivery Board",
        items: [
          {
            id: "task_runtime",
            title: "实现运行态观测",
            type: "code",
            status: "handoff_pending",
            assignedAgentId: "agent_backend"
          }
        ]
      }
    });
  });

  await page.route("**/api/projects/proj_1/executor-sessions/session_1", async (route) => {
    await route.fulfill({
      json: {
        id: "session_1",
        backend: "codex_cli",
        status: "completed"
      }
    });
  });

  await page.route("**/api/projects/proj_1/executor-sessions/session_1/events", async (route) => {
    await route.fulfill({
      json: [
        {
          id: "event_1",
          title: "会话已启动",
          type: "session_started",
          category: "session",
          summary: "worker 已启动 codex_cli"
        }
      ]
    });
  });

  await page.route("**/api/projects/proj_1/executor-sessions/session_1/transcript?view=redacted", async (route) => {
    await route.fulfill({
      json: {
        transcript: {
          id: "transcript_1",
          status: "sealed",
          executorSessionId: "session_1"
        },
        entries: [
          {
            id: "entry_1",
            role: "assistant",
            entryType: "message",
            redactedText: "已为管理员生成运行摘要。"
          }
        ]
      }
    });
  });

  await page.route("**/api/projects/proj_1/executor-sessions/session_1/access-audits", async (route) => {
    await route.fulfill({
      json: [
        {
          id: "audit_1",
          actorUserId: "demo-admin",
          accessMode: "redacted_transcript",
          reason: "审阅转录"
        }
      ]
    });
  });

  await page.goto("/missions/mission_1");

  await expect(page.getByText("演示任务")).toBeVisible();
  await expect(page.getByRole("heading", { name: "架构快照" })).toBeVisible();
  await expect(page.getByText("运行态观测范围")).toBeVisible();
  await expect(page.getByText("正在输出 codex 会话")).toBeVisible();
  await expect(page.getByText("实现运行态观测")).toBeVisible();

  await page.goto("/runtime/sessions/session_1");

  await expect(page.getByText("运行会话 · codex_cli")).toBeVisible();
  await expect(page.getByText("已为管理员生成运行摘要。")).toBeVisible();
  await expect(page.getByText("demo-admin")).toBeVisible();
});
