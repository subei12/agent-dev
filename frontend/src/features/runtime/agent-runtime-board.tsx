import { MissionRuntime } from "../../lib/api";

const fallbackRuntimes: MissionRuntime[] = [
  { id: "runtime_1", agentId: "agent_backend", status: "working", statusSummary: "streaming codex session", currentExecutorSessionId: "session_1" },
  { id: "runtime_2", agentId: "agent_reviewer", status: "waiting_admin", statusSummary: "awaiting handoff decision", currentExecutorSessionId: "session_2" }
];

export function AgentRuntimeBoard({ runtimes }: { runtimes: MissionRuntime[] }) {
  const items = runtimes.length > 0 ? runtimes : fallbackRuntimes;

  return (
    <section className="panel">
      <div className="panel-header">
        <p className="panel-kicker">Active Agents</p>
        <span className="badge">{items.length} visible</span>
      </div>
      <div className="runtime-grid">
        {items.map((runtime) => (
          <article key={runtime.id} className="runtime-card">
            <div className="runtime-card__meta">
              <span>{runtime.agentId}</span>
              <span>{runtime.status}</span>
            </div>
            <p>{runtime.statusSummary ?? "No runtime summary yet."}</p>
            <code>{runtime.currentExecutorSessionId ?? "session unavailable"}</code>
          </article>
        ))}
      </div>
    </section>
  );
}
