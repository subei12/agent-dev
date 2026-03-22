import { Mission, TaskItemView } from "./api";

export type MissionProgressStep = {
  key: string;
  label: string;
  state: "pending" | "active" | "done";
  detail: string;
};

/**
 * deriveMissionProgress 根据 Mission 和内部任务状态推导自动推进阶段轨迹。
 */
export function deriveMissionProgress(mission: Mission | null, tasks: TaskItemView[] = []): MissionProgressStep[] {
  const designTask = tasks.find((task) => task.type === "design");
  const codeTask = tasks.find((task) => task.type === "code");
  const testTask = tasks.find((task) => task.type === "test");
  const reviewTask = tasks.find((task) => task.type === "review");

  const designState = deriveTaskState(designTask);
  const codeState = deriveTaskState(codeTask);
  const testState = deriveTaskState(testTask);
  const reviewState = deriveTaskState(reviewTask);
  const missionState = mission?.status === "completed" ? "done" : reviewState === "done" ? "active" : "pending";

  return [
    {
      key: "design",
      label: "需求收敛",
      state: designState,
      detail: designTask ? `管理员 Agent 收敛统一方案` : "等待系统生成方案任务"
    },
    {
      key: "code",
      label: "开发实现",
      state: codeState,
      detail: codeTask ? "执行 Agent 正在实现与自测" : "等待进入开发阶段"
    },
    {
      key: "test",
      label: "测试验收",
      state: testState,
      detail: testTask ? "自动进入测试与验收" : "等待进入测试阶段"
    },
    {
      key: "review",
      label: "管理员复核",
      state: reviewState,
      detail: reviewTask ? "管理员 Agent 复核并决定是否结项" : "等待进入管理员复核"
    },
    {
      key: "complete",
      label: "Mission 完成",
      state: missionState,
      detail: missionState === "done" ? "平台已自动完成结项" : "等待最终结项决策"
    }
  ];
}

/**
 * summarizeMissionProgress 返回当前阶段摘要，方便页面展示自动推进状态。
 */
export function summarizeMissionProgress(steps: MissionProgressStep[]): string {
  const activeStep = steps.find((step) => step.state === "active");
  if (activeStep) {
    return `当前阶段：${activeStep.label}`;
  }

  const pendingStep = steps.find((step) => step.state === "pending");
  if (pendingStep) {
    return `下一阶段：${pendingStep.label}`;
  }

  return "当前阶段：Mission 完成";
}

/**
 * deriveTaskState 把内部任务状态映射为阶段轨迹状态。
 */
function deriveTaskState(task: TaskItemView | undefined): "pending" | "active" | "done" {
  if (!task) {
    return "pending";
  }
  if (task.status === "done") {
    return "done";
  }
  if (task.status === "claimed" || task.status === "in_progress" || task.status === "handoff_pending" || task.status === "review") {
    return "active";
  }
  return "pending";
}
