import { expect, test } from "@playwright/test";

test("home page lets user create one mission requirement and enter the mission workspace", async ({ page }) => {
  const missions = [
    {
      id: "mission_1",
      title: "交付运行态观测工作台",
      status: "implementation",
      description: "用于演示讨论、文档、任务、运行日志与归档的 Mission。"
    }
  ];

  await page.route("**/api/projects/proj_1/missions", async (route) => {
    if (route.request().method() === "POST") {
      const body = JSON.parse(route.request().postData() ?? "{}");
      const createdMission = {
        id: "mission_2",
        title: body.title,
        status: "discussion",
        description: body.description
      };
      missions.unshift(createdMission);
      await route.fulfill({
        status: 201,
        json: createdMission
      });
      return;
    }

    await route.fulfill({ json: missions });
  });

  await page.route("**/api/projects/proj_1/missions/*/task-board", async (route) => {
    await route.fulfill({
      json: {
        id: "board_demo",
        title: "内部执行看板",
        items: []
      }
    });
  });

  await page.route("**/api/projects/proj_1/missions/mission_2", async (route) => {
    await route.fulfill({
      json: {
        id: "mission_2",
        title: "接入多 Agent 自动拆解流程",
        status: "discussion",
        description: "平台应自动从设计推进到开发和测试。"
      }
    });
  });

  await page.route("**/api/projects/proj_1/missions/mission_2/documents", async (route) => {
    await route.fulfill({ json: [] });
  });

  await page.route("**/api/projects/proj_1/missions/mission_2/discussions", async (route) => {
    await route.fulfill({ json: [] });
  });

  await page.route("**/api/projects/proj_1/missions/mission_2/agent-runtimes", async (route) => {
    await route.fulfill({ json: [] });
  });

  await page.route("**/api/projects/proj_1/agents", async (route) => {
    await route.fulfill({ json: [] });
  });

  await page.route("**/api/projects/proj_1/missions/mission_2/approvals", async (route) => {
    await route.fulfill({ json: [] });
  });

  await page.route("**/api/projects/proj_1/missions/mission_2/archive", async (route) => {
    await route.fulfill({ json: null, status: 404 });
  });

  await page.route("**/api/projects/proj_1/missions/mission_2/task-board", async (route) => {
    await route.fulfill({
      json: {
        id: "board_2",
        title: "内部执行看板",
        items: [
          {
            id: "task_plan",
            title: "管理员 Agent 收敛方案",
            type: "design",
            status: "todo",
            assignedAgentId: "agent_admin"
          }
        ]
      }
    });
  });

  await page.goto("/");

  await page.getByLabel("需求标题").fill("接入多 Agent 自动拆解流程");
  await page.getByLabel("需求说明").fill("平台应自动从设计推进到开发和测试。");
  await page.getByRole("button", { name: "提交需求" }).click();

  await expect(page).toHaveURL(/\/missions\/mission_2$/);
  await expect(page.getByText("自动推进状态", { exact: true })).toBeVisible();
  await expect(page.getByText("内部执行任务", { exact: true })).toBeVisible();
  await expect(page.getByRole("heading", { name: "管理员 Agent 收敛方案" }).first()).toBeVisible();
});
