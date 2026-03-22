import { expect, test } from "@playwright/test";

test("mission workspace supports creating discussion, document, claiming task, and archiving", async ({
  page
}) => {
  const sessions = [
    {
      id: "session_1",
      topic: "运行态观测范围",
      status: "open"
    }
  ];

  const documents = [
    {
      id: "doc_1",
      title: "架构快照",
      kind: "architecture",
      currentAdoptedVersionId: "docv_2"
    }
  ];
  const versionsByDocument: Record<string, Array<{
    id: string;
    documentId: string;
    version: number;
    status: string;
    contentFormat: string;
    contentHash: string;
    contentText: string;
  }>> = {
    doc_1: [
      {
        id: "docv_2",
        documentId: "doc_1",
        version: 2,
        status: "adopted",
        contentFormat: "md",
        contentHash: "hash-docv2",
        contentText: "已采纳架构版本"
      }
    ]
  };

  const taskBoard = {
    id: "board_1",
    title: "Delivery Board",
    items: [
      {
        id: "task_runtime",
        title: "实现运行态观测",
        type: "code",
        status: "todo",
        assignedAgentId: "agent_backend"
      }
    ]
  };

  let archiveStatus = "未归档";

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

  await page.route("**/api/projects/proj_1/missions/mission_1/discussions", async (route) => {
    if (route.request().method() === "POST") {
      const body = JSON.parse(route.request().postData() ?? "{}");
      sessions.push({
        id: "session_2",
        topic: body.topic,
        status: "open"
      });
      await route.fulfill({
        status: 201,
        json: sessions[sessions.length - 1]
      });
      return;
    }

    await route.fulfill({ json: sessions });
  });

  await page.route("**/api/projects/proj_1/missions/mission_1/documents", async (route) => {
    if (route.request().method() === "POST") {
      const body = JSON.parse(route.request().postData() ?? "{}");
      documents.push({
        id: "doc_2",
        title: body.title,
        kind: body.kind,
        currentAdoptedVersionId: undefined
      });
      await route.fulfill({
        status: 201,
        json: documents[documents.length - 1]
      });
      return;
    }

    await route.fulfill({ json: documents });
  });

  await page.route(
    new RegExp(".*/api/projects/proj_1/missions/mission_1/documents/([^/]+)/versions"),
    async (route) => {
      const documentId = route.request().url().split("/documents/")[1].split("/versions")[0];

      if (route.request().method() === "POST") {
        const body = JSON.parse(route.request().postData() ?? "{}");
        const nextVersion = (versionsByDocument[documentId]?.length ?? 0) + 1;
        const created = {
          id: `${documentId}_version_${nextVersion}`,
          documentId,
          version: nextVersion,
          status: "proposed",
          contentFormat: body.contentFormat,
          contentHash: body.contentHash,
          contentText: body.contentText
        };
        versionsByDocument[documentId] = [...(versionsByDocument[documentId] ?? []), created];
        await route.fulfill({
          status: 201,
          json: created
        });
        return;
      }

      await route.fulfill({
        json: versionsByDocument[documentId] ?? []
      });
    }
  );

  await page.route(
    new RegExp(".*/api/projects/proj_1/missions/mission_1/documents/([^/]+)/adopt"),
    async (route) => {
      const documentId = route.request().url().split("/documents/")[1].split("/adopt")[0];
      const body = JSON.parse(route.request().postData() ?? "{}");

      versionsByDocument[documentId] = (versionsByDocument[documentId] ?? []).map((version) => ({
        ...version,
        status: version.id === body.versionId ? "adopted" : "proposed"
      }));
      const adopted = versionsByDocument[documentId].find((version) => version.id === body.versionId)!;
      const document = documents.find((item) => item.id === documentId);
      if (document) {
        document.currentAdoptedVersionId = adopted.id;
      }
      await route.fulfill({
        status: 201,
        json: adopted
      });
    }
  );

  await page.route("**/api/projects/proj_1/missions/mission_1/task-board", async (route) => {
    await route.fulfill({ json: taskBoard });
  });

  await page.route("**/api/projects/proj_1/missions/mission_1/tasks/task_runtime/claim", async (route) => {
    taskBoard.items[0].status = "claimed";
    await route.fulfill({
      status: 201,
      json: {
        id: "claim_1",
        taskItemId: "task_runtime",
        agentId: "agent_backend",
        status: "active"
      }
    });
  });

  await page.route("**/api/projects/proj_1/missions/mission_1/tasks/task_runtime/handoffs", async (route) => {
    taskBoard.items[0].status = "handoff_pending";
    await route.fulfill({
      status: 201,
      json: {
        id: "handoff_1",
        taskItemId: "task_runtime",
        toAdminAgent: true,
        status: "pending"
      }
    });
  });

  await page.route(
    "**/api/projects/proj_1/missions/mission_1/tasks/task_runtime/review-checkpoints",
    async (route) => {
      await route.fulfill({
        status: 201,
        json: {
          id: "checkpoint_1",
          missionId: "mission_1",
          taskItemId: "task_runtime"
        }
      });
    }
  );

  await page.route("**/api/projects/proj_1/missions/mission_1/archive", async (route) => {
    if (route.request().method() === "POST") {
      archiveStatus = "uploaded";
      await route.fulfill({
        status: 201,
        json: {
          archive: {
            id: "archive_1",
            missionId: "mission_1",
            status: "uploaded"
          },
          manifest: {
            missionId: "mission_1",
            objectKeys: ["runtime/session_1.json"]
          }
        }
      });
      return;
    }

    await route.fulfill({
      json: {
        id: "archive_1",
        missionId: "mission_1",
        status: archiveStatus
      }
    });
  });

  await page.route("**/api/projects/proj_1/missions/mission_1/agent-runtimes", async (route) => {
    await route.fulfill({ json: [] });
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
          args: ["--approval-mode", "full-auto"],
          timeoutSec: 3600,
          maxConcurrency: 1,
          allowRepoRead: true,
          allowRepoWrite: true,
          allowNetwork: false
        }
      ]
    });
  });

  await page.route("**/api/projects/proj_1/missions/mission_1/approvals", async (route) => {
    await route.fulfill({ json: [] });
  });

  await page.goto("/missions/mission_1");

  await page.getByLabel("讨论主题").fill("管理员复核流程");
  await page.getByRole("button", { name: "发起讨论" }).click();
  await expect(page.getByText("管理员复核流程")).toBeVisible();

  await page.getByLabel("文档标题").fill("测试策略");
  await page.getByRole("button", { name: "新建文档" }).click();
  await expect(page.getByRole("heading", { name: "测试策略" })).toBeVisible();

  await page.getByLabel("架构快照 的版本内容").fill("版本 3 草稿");
  await page.getByRole("button", { name: "保存版本" }).first().click();
  await expect(page.getByText("版本 3 草稿")).toBeVisible();
  await page.getByRole("button", { name: "采纳版本" }).last().click();
  await expect(page.getByText("已采纳", { exact: true })).toBeVisible();

  await page.getByRole("button", { name: "领取任务" }).click();
  await expect(page.getByText("任务已领取", { exact: true })).toBeVisible();
  await expect(page.getByText("已领取", { exact: true }).first()).toBeVisible();

  await page.getByRole("button", { name: "提交管理员" }).click();
  await expect(page.getByText("任务已提交给管理员", { exact: true })).toBeVisible();
  await expect(page.getByText("等待交接", { exact: true }).first()).toBeVisible();

  await page.getByRole("button", { name: "请求评审" }).click();
  await expect(page.getByText("已发起检查点")).toBeVisible();

  await page.getByRole("button", { name: "归档任务" }).click();
  await expect(page.getByText("归档状态：已上传")).toBeVisible();
});
