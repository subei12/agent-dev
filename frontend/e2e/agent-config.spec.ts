import { expect, test } from "@playwright/test";

test("agent config page lists agents and saves local cli settings", async ({ page }) => {
  const agents = [
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
  ];

  await page.route("**/api/projects/proj_1/agents", async (route) => {
    if (route.request().method() === "PATCH") {
      const body = JSON.parse(route.request().postData() ?? "{}");
      agents[0] = { ...agents[0], ...body };
      await route.fulfill({ status: 200, json: agents[0] });
      return;
    }

    await route.fulfill({ json: agents });
  });

  await page.goto("/agents");

  await expect(page.getByText("后端 Agent")).toBeVisible();
  const commandInput = page.getByLabel("本地命令");
  await commandInput.fill("claude");
  await page.getByRole("button", { name: "保存 Agent 配置" }).click();
  await expect(commandInput).toHaveValue("claude");
});
