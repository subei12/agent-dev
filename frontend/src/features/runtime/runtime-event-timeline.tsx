import { RuntimeEvent } from "../../lib/api";
import { formatTypeLabel } from "../../lib/display";

const fallbackEvents: RuntimeEvent[] = [
  { id: "event_1", title: "会话已启动", type: "session_started", category: "session", summary: "worker 已启动 codex_cli" },
  { id: "event_2", title: "模型输出", type: "message_created", category: "model", summary: "assistant 生成了运行摘要" },
  { id: "event_3", title: "等待管理员", type: "waiting_admin", category: "handoff", summary: "任务完成但没有下游节点" }
];

/**
 * RuntimeEventTimeline 渲染或处理当前前端行为。
 */
export function RuntimeEventTimeline({ events }: { events: RuntimeEvent[] }) {
  const items = events.length > 0 ? events : fallbackEvents;

  return (
    <section className="panel">
      <div className="panel-header">
        <p className="panel-kicker">运行事件</p>
        <span className="badge">{items.length} 条</span>
      </div>
      <div className="timeline">
        {items.map((event) => (
          <article key={event.id} className="timeline-item">
            <div className="timeline-item__line" />
            <div>
              <div className="timeline-item__meta">
                <span>{formatTypeLabel(event.category)}</span>
                <span>{formatTypeLabel(event.type)}</span>
              </div>
              <h3>{event.title}</h3>
              <p>{event.summary ?? "当前没有事件摘要。"}</p>
            </div>
          </article>
        ))}
      </div>
    </section>
  );
}
