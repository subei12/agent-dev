import { expect, test } from "@playwright/test";

test("app shell shows mission navigation", async ({ page }) => {
  await page.goto("/");

  const nav = page.getByRole("navigation");
  await expect(nav).toBeVisible();
  await expect(nav.getByRole("link", { name: "首页", exact: true })).toBeVisible();
  await expect(nav.getByRole("link", { name: "运行态", exact: true })).toBeVisible();
  await expect(page.getByRole("main").getByRole("heading", { name: "Mission 总览" })).toBeVisible();
  await expect(page.getByRole("main").getByText("总任务看板", { exact: true }).first()).toBeVisible();
});

test("mission page shows documents and task board", async ({ page }) => {
  await page.goto("/missions/mission_1");

  await expect(page.getByText("任务看板", { exact: true })).toBeVisible();
  await expect(page.getByText("文档", { exact: true })).toBeVisible();
  await expect(page.getByText("讨论轮次", { exact: true })).toBeVisible();
});

test("runtime board shows active agents and event timeline", async ({ page }) => {
  await page.goto("/missions/mission_1");

  await expect(page.getByText("活动 Agent", { exact: true })).toBeVisible();
  await expect(page.getByText("运行事件", { exact: true })).toBeVisible();
});
