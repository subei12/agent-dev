import { MissionRuntime } from "../../lib/api";
import { formatStatusLabel } from "../../lib/display";

const fallbackRuntimes: MissionRuntime[] = [
  { id: "runtime_1", agentId: "agent_backend", status: "working", statusSummary: "正在输出 codex 会话", currentExecutorSessionId: "session_1" },
  { id: "runtime_2", agentId: "agent_reviewer", status: "waiting_admin", statusSummary: "等待管理员处理交接", currentExecutorSessionId: "session_2" }
];

/**
 * AgentRuntimeBoard 渲染或处理当前前端行为。
 */
export function AgentRuntimeBoard({
  runtimes,
  onSelectRuntime
}: {
  runtimes: MissionRuntime[];
  onSelectRuntime: (sessionId: string | null) => void;
}) {
  const items = runtimes.length > 0 ? runtimes : fallbackRuntimes;

  return (
    <section className="panel">
      <div className="panel-header">
        <p className="panel-kicker">活动 Agent</p>
        <span className="badge">{items.length} 个</span>
      </div>
      <div className="runtime-grid">
        {items.map((runtime) => (
          <button
            key={runtime.id}
            className="runtime-card"
            onClick={() => onSelectRuntime(runtime.currentExecutorSessionId ?? null)}
            type="button"
          >
            <div className="runtime-card__meta">
              <span>{runtime.agentId}</span>
              <span>{formatStatusLabel(runtime.status)}</span>
            </div>
            <p>{runtime.statusSummary ?? "当前还没有运行摘要。"}</p>
            <code>{runtime.currentExecutorSessionId ?? "当前没有 session"}</code>
          </button>
        ))}
      </div>
    </section>
  );
}
