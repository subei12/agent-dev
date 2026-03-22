import { useState } from "react";

import { MissionRuntime, TranscriptAccessAudit, TranscriptView, RuntimeEvent } from "../../lib/api";
import { formatStatusLabel } from "../../lib/display";
import { RuntimeEventTimeline } from "./runtime-event-timeline";
import { TranscriptInspector } from "./transcript-inspector";
import { TranscriptAccessAuditTable } from "./transcript-access-audit-table";

type TaskRuntimePanelProps = {
  runtimes: MissionRuntime[];
  events: RuntimeEvent[];
  transcript: TranscriptView | null;
  audits: TranscriptAccessAudit[];
};

/**
 * TaskRuntimePanel 在 Mission 页内展示 Agent 执行状态，并支持就地查看日志。
 */
export function TaskRuntimePanel({ runtimes, events, transcript, audits }: TaskRuntimePanelProps) {
  const [selectedAgentId, setSelectedAgentId] = useState<string | null>(runtimes[0]?.agentId ?? null);

  return (
    <section className="panel panel--feature">
      <div className="panel-header">
        <p className="panel-kicker">Agent 执行区</p>
        <span className="badge">{runtimes.length} 个 Agent</span>
      </div>
      <div className="runtime-grid">
        {runtimes.map((runtime) => (
          <button
            key={runtime.id}
            className={selectedAgentId === runtime.agentId ? "runtime-card runtime-card--selected" : "runtime-card"}
            onClick={() => setSelectedAgentId(runtime.agentId)}
            type="button"
          >
            <div className="runtime-card__meta">
              <span>{runtime.agentId}</span>
              <span>{formatStatusLabel(runtime.status)}</span>
            </div>
            <p>{runtime.statusSummary ?? "暂无运行摘要。"}</p>
            <code>{runtime.currentExecutorSessionId ?? "暂无会话"}</code>
          </button>
        ))}
      </div>
      <RuntimeEventTimeline events={events} />
      <TranscriptInspector transcript={transcript} />
      <TranscriptAccessAuditTable audits={audits} />
    </section>
  );
}
