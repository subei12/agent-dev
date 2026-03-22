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
    { key: "todo", label: "待开始" },
    { key: "claimed", label: "已领取" },
    { key: "in_progress", label: "进行中" },
    { key: "handoff_pending", label: "等待交接" },
    { key: "done", label: "已完成" },
    { key: "review", label: "评审中" }
  ];

  return (
    <section className="panel">
      <div className="panel-header">
        <p className="panel-kicker">任务看板</p>
        <span className="badge">{items.length} 个任务</span>
      </div>
      {taskActionMessage ? (
        <div className="empty-card">
          <strong>{taskActionMessage}</strong>
        </div>
      ) : null}
      <div className="task-list-groups">
        {groups.map((group) => {
          const groupItems = items.filter((task) => task.status === group.key);
          if (groupItems.length === 0) return null;

          return (
            <section key={group.key} className="task-group">
              <div className="panel-header">
                <p className="panel-kicker">{group.label}</p>
                <span className="badge">{groupItems.length} 个</span>
              </div>
              <div className="stack">
                {groupItems.map((task) => (
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
                    <p>{task.assignedAgentId ?? "未分配"}</p>
                    <div className="task-actions">
                      <button className="action-button" onClick={() => onClaimTask(task.id)} type="button">
                        领取任务
                      </button>
                      <button className="action-button action-button--ghost" onClick={() => onSendToAdmin(task.id)} type="button">
                        提交管理员
                      </button>
                      <button className="action-button action-button--ghost" onClick={() => onRequestReview(task.id)} type="button">
                        请求评审
                      </button>
                    </div>
                  </article>
                ))}
              </div>
            </section>
          );
        })}
      </div>
    </section>
  );
}
