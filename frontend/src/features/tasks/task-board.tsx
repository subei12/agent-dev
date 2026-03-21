type TaskCard = {
  id: string;
  title: string;
  type: string;
  status: string;
  assignedAgentId?: string;
};

const fallbackTasks: TaskCard[] = [
  { id: "task-design", title: "Converge architecture notes", type: "design", status: "done", assignedAgentId: "agent_admin" },
  { id: "task-code", title: "Implement runtime observability service", type: "code", status: "in_progress", assignedAgentId: "agent_backend" },
  { id: "task-review", title: "Review transcript access policy", type: "review", status: "todo", assignedAgentId: "agent_reviewer" }
];

export function TaskBoard({ tasks = fallbackTasks }: { tasks?: TaskCard[] }) {
  return (
    <section className="panel">
      <div className="panel-header">
        <p className="panel-kicker">Task Board</p>
        <span className="badge">{tasks.length} active lanes</span>
      </div>
      <div className="task-grid">
        {tasks.map((task) => (
          <article key={task.id} className="task-card">
            <div className="task-card__meta">
              <span>{task.type}</span>
              <span>{task.status}</span>
            </div>
            <h3>{task.title}</h3>
            <p>{task.assignedAgentId ?? "unassigned"}</p>
          </article>
        ))}
      </div>
    </section>
  );
}
