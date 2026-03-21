type TaskCard = {
  id: string;
  title: string;
  type: string;
  status: string;
  assignedAgentId?: string;
};

type TaskBoardProps = {
  tasks?: TaskCard[];
  onClaimTask: (taskId: string) => void;
  onSendToAdmin: (taskId: string) => void;
  onRequestReview: (taskId: string) => void;
  taskActionMessage?: string;
};

const fallbackTasks: TaskCard[] = [
  { id: "task-design", title: "Converge architecture notes", type: "design", status: "done", assignedAgentId: "agent_admin" },
  { id: "task-code", title: "Implement runtime observability service", type: "code", status: "in_progress", assignedAgentId: "agent_backend" },
  { id: "task-review", title: "Review transcript access policy", type: "review", status: "todo", assignedAgentId: "agent_reviewer" }
];

/**
 * TaskBoard 渲染或处理当前前端行为。
 */
export function TaskBoard({
  tasks,
  onClaimTask,
  onSendToAdmin,
  onRequestReview,
  taskActionMessage
}: TaskBoardProps) {
  const items = tasks && tasks.length > 0 ? tasks : fallbackTasks;

  return (
    <section className="panel">
      <div className="panel-header">
        <p className="panel-kicker">Task Board</p>
        <span className="badge">{items.length} active lanes</span>
      </div>
      {taskActionMessage ? (
        <div className="empty-card">
          <strong>{taskActionMessage}</strong>
        </div>
      ) : null}
      <div className="task-grid">
        {items.map((task) => (
          <article key={task.id} className="task-card">
            <div className="task-card__meta">
              <span>{task.type}</span>
              <span>{task.status}</span>
            </div>
            <h3>{task.title}</h3>
            <p>{task.assignedAgentId ?? "unassigned"}</p>
            <div className="task-actions">
              <button className="action-button" onClick={() => onClaimTask(task.id)} type="button">
                Claim Task
              </button>
              <button className="action-button action-button--ghost" onClick={() => onSendToAdmin(task.id)} type="button">
                Send To Admin
              </button>
              <button className="action-button action-button--ghost" onClick={() => onRequestReview(task.id)} type="button">
                Request Review
              </button>
            </div>
          </article>
        ))}
      </div>
    </section>
  );
}
