import { expect, test } from "@playwright/test";

test("mission workspace renders seeded api data and runtime transcript", async ({ page }) => {
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

  await page.route("http://127.0.0.1:8080/api/projects/proj_1/missions/mission_1/documents", async (route) => {
    await route.fulfill({
      json: [
        {
          id: "doc_1",
          title: "Architecture Snapshot",
          kind: "architecture",
          currentAdoptedVersionId: "docv_2"
        }
      ]
    });
  });

  await page.route("http://127.0.0.1:8080/api/projects/proj_1/missions/mission_1/discussions", async (route) => {
    await route.fulfill({
      json: [
        {
          id: "session_1",
          topic: "Runtime observability scope",
          status: "open"
        }
      ]
    });
  });

  await page.route("http://127.0.0.1:8080/api/projects/proj_1/missions/mission_1/agent-runtimes", async (route) => {
    await route.fulfill({
      json: [
        {
          id: "runtime_1",
          agentId: "agent_backend",
          status: "working",
          statusSummary: "streaming codex session",
          currentExecutorSessionId: "session_telemetry"
        }
      ]
    });
  });

  await page.route("http://127.0.0.1:8080/api/projects/proj_1/executor-sessions/session_1", async (route) => {
    await route.fulfill({
      json: {
        id: "session_1",
        backend: "codex_cli",
        status: "completed"
      }
    });
  });

  await page.route("http://127.0.0.1:8080/api/projects/proj_1/executor-sessions/session_1/events", async (route) => {
    await route.fulfill({
      json: [
        {
          id: "event_1",
          title: "Session started",
          type: "session_started",
          category: "session",
          summary: "worker launched codex_cli"
        }
      ]
    });
  });

  await page.route("http://127.0.0.1:8080/api/projects/proj_1/executor-sessions/session_1/transcript?view=redacted", async (route) => {
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
            redactedText: "Drafted runtime summary for admin review."
          }
        ]
      }
    });
  });

  await page.route("http://127.0.0.1:8080/api/projects/proj_1/executor-sessions/session_1/access-audits", async (route) => {
    await route.fulfill({
      json: [
        {
          id: "audit_1",
          actorUserId: "demo-admin",
          accessMode: "redacted_transcript",
          reason: "review transcript"
        }
      ]
    });
  });

  await page.goto("/missions/mission_1");

  await expect(page.getByText("Seeded mission")).toBeVisible();
  await expect(page.getByText("Architecture Snapshot")).toBeVisible();
  await expect(page.getByText("Runtime observability scope")).toBeVisible();
  await expect(page.getByText("streaming codex session")).toBeVisible();

  await page.goto("/runtime/sessions/session_1");

  await expect(page.getByText("Runtime Session · codex_cli")).toBeVisible();
  await expect(page.getByText("Drafted runtime summary for admin review.")).toBeVisible();
  await expect(page.getByText("demo-admin")).toBeVisible();
});
