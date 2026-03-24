import { expect, test } from "@playwright/test";

test("app shell shows mission navigation", async ({ page }) => {
  await page.goto("/");

  const nav = page.getByRole("navigation");
  await expect(nav).toBeVisible();
  await expect(nav.getByRole("link", { name: "首页", exact: true })).toBeVisible();
  await expect(nav.getByRole("link", { name: "工作台", exact: true })).toBeVisible();
  await expect(nav.getByRole("link", { name: "运行态", exact: true })).toBeVisible();
  await expect(page.getByText("Agent Platform v2")).toBeVisible();
  await expect(page.getByRole("main").getByRole("heading", { name: "任务指挥台" })).toBeVisible();
  await expect(page.getByText("新建需求", { exact: true })).toBeVisible();
  await expect(page.getByText("任务脉冲", { exact: true })).toBeVisible();
  await expect(page.getByText("最近动态", { exact: true })).toBeVisible();
  await expect(page.getByText("项目任务列表", { exact: true })).toBeVisible();
});

test("mission page shows documents and task board", async ({ page }) => {
  await page.goto("/missions/mission_1");

  await expect(page.getByText("内部执行任务", { exact: true })).toBeVisible();
  await expect(page.getByText("自动推进状态", { exact: true })).toBeVisible();
  await expect(page.getByText("需求收敛", { exact: true })).toBeVisible();
  await expect(page.getByText("待开始", { exact: true }).first()).toBeVisible();
  await expect(page.getByText("文档", { exact: true })).toBeVisible();
  await expect(page.getByText("任务详情", { exact: true })).toBeVisible();
});

test("runtime board shows active agents and event timeline", async ({ page }) => {
  await page.goto("/missions/mission_1");

  await expect(page.getByText("Agent 执行区", { exact: true })).toBeVisible();
  await expect(page.getByText("运行事件", { exact: true })).toBeVisible();
});
