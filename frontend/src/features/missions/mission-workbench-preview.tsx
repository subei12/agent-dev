import { TaskBoard } from "../../lib/api";
import { summarizeMissionProgress, deriveMissionProgress } from "../../lib/mission-progress";

type MissionWorkbenchPreviewProps = {
  missionTitle: string;
  board: TaskBoard | null;
};

/**
 * MissionWorkbenchPreview 用简化后的真实工作台结构做首页主视觉演示。
 */
export function MissionWorkbenchPreview({ missionTitle, board }: MissionWorkbenchPreviewProps) {
  const tasks = board?.items ?? [];
  const summary = summarizeMissionProgress(
    deriveMissionProgress(
      missionTitle
        ? {
            id: "preview",
            title: missionTitle,
            status: tasks.every((task) => task.status === "done") ? "completed" : "implementation"
          }
        : null,
      tasks
    )
  );

  return (
    <section className="panel dashboard-preview-card">
      <div className="panel-header">
        <p className="panel-kicker">工作台演示</p>
        <span className="badge">真实产品视图</span>
      </div>
      <div className="preview-shell">
        <aside className="preview-sidebar">
          <strong>Mission</strong>
          <p>{missionTitle || "等待需求进入"}</p>
          <span>{summary}</span>
        </aside>
        <div className="preview-main">
          <div className="preview-main__header">
            <strong>内部执行任务</strong>
            <span>{tasks.length} 个</span>
          </div>
          <div className="preview-task-list">
            {tasks.slice(0, 4).map((task) => (
              <article key={task.id} className="preview-task-item">
                <strong>{task.title}</strong>
                <p>{task.assignedAgentId ?? "未分配 Agent"}</p>
              </article>
            ))}
          </div>
        </div>
      </div>
    </section>
  );
}
