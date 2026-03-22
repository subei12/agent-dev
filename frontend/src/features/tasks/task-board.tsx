import { formatStatusLabel, formatTypeLabel } from "../../lib/display";

type TaskCard = {
  id: string;
  title: string;
  type: string;
  status: string;
  assignedAgentId?: string;
};

type TaskBoardProps = {
  tasks?: TaskCard[];
  selectedTaskId?: string | null;
  onSelectTask: (taskId: string) => void;
  onClaimTask: (taskId: string) => void;
  onSendToAdmin: (taskId: string) => void;
  onRequestReview: (taskId: string) => void;
  taskActionMessage?: string;
};

const fallbackTasks: TaskCard[] = [
  { id: "task-design", title: "收敛架构文档", type: "design", status: "done", assignedAgentId: "agent_admin" },
  { id: "task-code", title: "实现运行态观测服务", type: "code", status: "in_progress", assignedAgentId: "agent_backend" },
  { id: "task-review", title: "评审 transcript 访问策略", type: "review", status: "todo", assignedAgentId: "agent_reviewer" }
];

/**
 * TaskBoard 渲染或处理当前前端行为。
 */
export function TaskBoard({
  tasks,
  selectedTaskId,
  onSelectTask,
  onClaimTask,
  onSendToAdmin,
  onRequestReview,
  taskActionMessage
}: TaskBoardProps) {
  const items = tasks && tasks.length > 0 ? tasks : fallbackTasks;
  const groups = [
    { key: "todo", label: "待开始", statuses: ["todo"] },
    { key: "progress", label: "进行中", statuses: ["claimed", "in_progress", "handoff_pending"] },
    { key: "review", label: "评审", statuses: ["review"] },
    { key: "done", label: "已完成", statuses: ["done"] }
  ];

  return (
    <div className="task-board-shell">
      {taskActionMessage ? (
        <div className="empty-card">
          <strong>{taskActionMessage}</strong>
        </div>
      ) : null}
      <div className="task-board-columns">
        {groups.map((group) => {
          const groupItems = items.filter((task) => group.statuses.includes(task.status));

          return (
            <section key={group.key} className="task-board-column">
              <div className="task-column-header">
                <strong>{group.label}</strong>
                <span className="badge">{groupItems.length} 个</span>
              </div>
              <div className="stack">
                {groupItems.length === 0 ? (
                  <div className="empty-card task-column-empty">
                    <strong>当前无任务</strong>
                  </div>
                ) : (
                  groupItems.map((task) => (
                    <article
                      key={task.id}
                      className={selectedTaskId === task.id ? "task-list-item task-list-item--selected" : "task-list-item"}
                      onClick={() => onSelectTask(task.id)}
                    >
                      <div className="task-card__meta">
                        <span>{formatTypeLabel(task.type)}</span>
                        <span>{formatStatusLabel(task.status)}</span>
                      </div>
                      <h3>{task.title}</h3>
                      <p>{task.assignedAgentId ?? "未分配 Agent"}</p>
                      <div className="task-actions task-actions--compact">
                        <button
                          className="action-button"
                          onClick={(event) => {
                            event.stopPropagation();
                            onClaimTask(task.id);
                          }}
                          type="button"
                        >
                          领取任务
                        </button>
                        <button
                          className="action-button action-button--ghost"
                          onClick={(event) => {
                            event.stopPropagation();
                            onSendToAdmin(task.id);
                          }}
                          type="button"
                        >
                          提交管理员
                        </button>
                        <button
                          className="action-button action-button--ghost"
                          onClick={(event) => {
                            event.stopPropagation();
                            onRequestReview(task.id);
                          }}
                          type="button"
                        >
                          请求评审
                        </button>
                      </div>
                    </article>
                  ))
                )}
              </div>
            </section>
          );
        })}
      </div>
    </div>
  );
}
