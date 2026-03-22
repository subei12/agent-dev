import { DocumentItem, TaskItemView } from "../../lib/api";
import { formatStatusLabel, formatTypeLabel } from "../../lib/display";

type TaskDetailPanelProps = {
  task: TaskItemView | null;
  documents: DocumentItem[];
};

/**
 * TaskDetailPanel 展示当前选中任务的详情和关联产物。
 */
export function TaskDetailPanel({ task, documents }: TaskDetailPanelProps) {
  if (!task) {
    return (
      <section className="panel">
        <div className="panel-header">
          <p className="panel-kicker">任务详情</p>
        </div>
        <div className="empty-card">
          <strong>请先选择一个任务。</strong>
          <p>选中任务后，这里会展示任务状态、负责人和关联产物。</p>
        </div>
      </section>
    );
  }

  const upstreamCount = Array.isArray(task.upstreamTaskIds) ? task.upstreamTaskIds.length : 0;
  const downstreamCount = Array.isArray(task.downstreamTaskIds) ? task.downstreamTaskIds.length : 0;

  return (
    <section className="panel panel--feature">
      <div className="panel-header panel-header--compact">
        <p className="panel-kicker">任务详情</p>
        <span className="badge">{formatStatusLabel(task.status)}</span>
      </div>
      <h2>{task.title}</h2>
      <p className="panel-copy">
        {formatTypeLabel(task.type)} · {task.assignedAgentId ?? "未分配 Agent"}
      </p>
      <div className="detail-metrics">
        <div className="compact-card">
          <strong>任务类型</strong>
          <p>{formatTypeLabel(task.type)}</p>
        </div>
        <div className="compact-card">
          <strong>当前负责人</strong>
          <p>{task.assignedAgentId ?? "尚未分配"}</p>
        </div>
      </div>
      <div className="stack">
        <div className="empty-card">
          <strong>上下游依赖</strong>
          <p>上游：{upstreamCount} 个 · 下游：{downstreamCount} 个</p>
        </div>
        <div className="empty-card">
          <strong>关联产物</strong>
          <p>{documents.length ? documents.map((doc) => doc.title).join("、") : "当前任务还没有关联文档。"}</p>
        </div>
      </div>
    </section>
  );
}
