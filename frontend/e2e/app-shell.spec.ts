import { expect, test } from "@playwright/test";

test("app shell shows mission navigation", async ({ page }) => {
  await page.goto("/");

  const nav = page.getByRole("navigation");
  await expect(nav).toBeVisible();
  await expect(nav.getByRole("link", { name: "Missions", exact: true })).toBeVisible();
  await expect(nav.getByRole("link", { name: "Runtime", exact: true })).toBeVisible();
  await expect(page.getByRole("main").getByRole("heading", { name: "Mission Control" })).toBeVisible();
});

test("mission page shows documents and task board", async ({ page }) => {
  await page.goto("/missions/mission_1");

  await expect(page.getByText("Task Board", { exact: true })).toBeVisible();
  await expect(page.getByText("Documents", { exact: true })).toBeVisible();
  await expect(page.getByText("Discussion Rounds", { exact: true })).toBeVisible();
});

test("runtime board shows active agents and event timeline", async ({ page }) => {
  await page.goto("/missions/mission_1");

  await expect(page.getByText("Active Agents", { exact: true })).toBeVisible();
  await expect(page.getByText("Runtime Events", { exact: true })).toBeVisible();
});
