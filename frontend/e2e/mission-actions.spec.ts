import { expect, test } from "@playwright/test";

test("mission workspace supports creating discussion, document, claiming task, and archiving", async ({
  page
}) => {
  const sessions = [
    {
      id: "session_1",
      topic: "Runtime observability scope",
      status: "open"
    }
  ];

  const documents = [
    {
      id: "doc_1",
      title: "Architecture Snapshot",
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
        contentText: "Adopted architecture"
      }
    ]
  };

  const taskBoard = {
    id: "board_1",
    title: "Delivery Board",
    items: [
      {
        id: "task_runtime",
        title: "Implement runtime observability",
        type: "code",
        status: "todo",
        assignedAgentId: "agent_backend"
      }
    ]
  };

  let archiveStatus = "none";

  await page.route("http://127.0.0.1:8080/api/projects/proj_1/missions/mission_1", async (route) => {
    await route.fulfill({
      json: {
        id: "mission_1",
        title: "Seeded mission",
        status: "implementation",
        description: "Deliver the runtime telemetry dashboard."
      }
    });
  });

  await page.route("http://127.0.0.1:8080/api/projects/proj_1/missions/mission_1/discussions", async (route) => {
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

  await page.route("http://127.0.0.1:8080/api/projects/proj_1/missions/mission_1/documents", async (route) => {
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
    /http:\/\/127\.0\.0\.1:8080\/api\/projects\/proj_1\/missions\/mission_1\/documents\/([^/]+)\/versions/,
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
    /http:\/\/127\.0\.0\.1:8080\/api\/projects\/proj_1\/missions\/mission_1\/documents\/([^/]+)\/adopt/,
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

  await page.route("http://127.0.0.1:8080/api/projects/proj_1/missions/mission_1/task-board", async (route) => {
    await route.fulfill({ json: taskBoard });
  });

  await page.route("http://127.0.0.1:8080/api/projects/proj_1/missions/mission_1/tasks/task_runtime/claim", async (route) => {
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

  await page.route("http://127.0.0.1:8080/api/projects/proj_1/missions/mission_1/tasks/task_runtime/handoffs", async (route) => {
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
    "http://127.0.0.1:8080/api/projects/proj_1/missions/mission_1/tasks/task_runtime/review-checkpoints",
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

  await page.route("http://127.0.0.1:8080/api/projects/proj_1/missions/mission_1/archive", async (route) => {
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

  await page.route("http://127.0.0.1:8080/api/projects/proj_1/missions/mission_1/agent-runtimes", async (route) => {
    await route.fulfill({ json: [] });
  });

  await page.goto("/missions/mission_1");

  await page.getByLabel("Discussion topic").fill("Admin review loop");
  await page.getByRole("button", { name: "Start Discussion" }).click();
  await expect(page.getByText("Admin review loop")).toBeVisible();

  await page.getByLabel("Document title").fill("Test Strategy");
  await page.getByRole("button", { name: "Add Document" }).click();
  await expect(page.getByRole("heading", { name: "Test Strategy" })).toBeVisible();

  await page.getByLabel("Version content for Architecture Snapshot").fill("Version 3 draft");
  await page.getByRole("button", { name: "Save Version" }).first().click();
  await expect(page.getByText("Version 3 draft")).toBeVisible();
  await page.getByRole("button", { name: "Adopt Version" }).last().click();
  await expect(page.getByText("adopted", { exact: true })).toBeVisible();

  await page.getByRole("button", { name: "Claim Task" }).click();
  await expect(page.getByText("Task claimed", { exact: true })).toBeVisible();
  await expect(page.getByText("claimed", { exact: true })).toBeVisible();

  await page.getByRole("button", { name: "Send To Admin" }).click();
  await expect(page.getByText("Task handed to admin", { exact: true })).toBeVisible();
  await expect(page.getByText("handoff_pending", { exact: true })).toBeVisible();

  await page.getByRole("button", { name: "Request Review" }).click();
  await expect(page.getByText("Checkpoint requested")).toBeVisible();

  await page.getByRole("button", { name: "Archive Mission" }).click();
  await expect(page.getByText("Archive status: uploaded")).toBeVisible();
});
