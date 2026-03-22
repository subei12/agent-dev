import { useEffect, useState } from "react";

import { AgentConfig, MissionRuntime, RuntimeEvent, TranscriptAccessAudit, TranscriptView } from "../../lib/api";
import { formatStatusLabel } from "../../lib/display";
import { RuntimeEventTimeline } from "./runtime-event-timeline";
import { TranscriptInspector } from "./transcript-inspector";
import { TranscriptAccessAuditTable } from "./transcript-access-audit-table";

type TaskRuntimePanelProps = {
  agents: AgentConfig[];
  runtimes: MissionRuntime[];
  selectedSessionId: string | null;
  events: RuntimeEvent[];
  transcript: TranscriptView | null;
  audits: TranscriptAccessAudit[];
  onSelectRuntime: (sessionId: string | null) => void;
};

/**
 * TaskRuntimePanel 在 Mission 页内展示 Agent 执行状态，并支持就地查看日志。
 */
export function TaskRuntimePanel({
  agents,
  runtimes,
  selectedSessionId,
  events,
  transcript,
  audits,
  onSelectRuntime
}: TaskRuntimePanelProps) {
  const [showTranscript, setShowTranscript] = useState(false);
  const agentLookup = Object.fromEntries(agents.map((agent) => [agent.id, agent]));
  const selectedRuntime =
    runtimes.find((runtime) => runtime.currentExecutorSessionId === selectedSessionId) ?? runtimes[0] ?? null;

  useEffect(() => {
    setShowTranscript(false);
  }, [selectedSessionId]);

  return (
    <section className="panel panel--feature">
      <div className="panel-header panel-header--compact">
        <p className="panel-kicker">Agent 执行区</p>
        <span className="badge">{runtimes.length} 个 Agent</span>
      </div>
      {runtimes.length > 0 ? (
        <>
          <div className="panel-subsection">
            <div className="panel-header panel-header--compact">
              <p className="panel-kicker">活动 Agent</p>
              <span className="badge">点击切换日志</span>
            </div>
            <div className="runtime-grid">
              {runtimes.map((runtime) => {
                const agent = agentLookup[runtime.agentId];
                const isSelected = selectedRuntime?.id === runtime.id;

                return (
                  <button
                    key={runtime.id}
                    className={isSelected ? "runtime-card runtime-card--selected" : "runtime-card"}
                    onClick={() => onSelectRuntime(runtime.currentExecutorSessionId ?? null)}
                    type="button"
                  >
                    <div className="runtime-card__meta">
                      <span>{agent?.executorType ?? runtime.agentId}</span>
                      <span>{formatStatusLabel(runtime.status)}</span>
                    </div>
                    <h3>{agent?.name ?? runtime.agentId}</h3>
                    <p>{runtime.statusSummary ?? "暂无运行摘要。"}</p>
                    <code>{runtime.currentExecutorSessionId ?? "暂无会话"}</code>
                  </button>
                );
              })}
            </div>
          </div>

          <div className="empty-card runtime-focus">
            <strong>当前查看：{agentLookup[selectedRuntime?.agentId ?? ""]?.name ?? selectedRuntime?.agentId ?? "未选择"}</strong>
            <p>
              状态：{formatStatusLabel(selectedRuntime?.status)} · 会话：
              {selectedRuntime?.currentExecutorSessionId ?? "暂无会话"}
            </p>
          </div>
        </>
      ) : (
        <div className="empty-card">
          <strong>当前没有正在执行的 Agent。</strong>
          <p>当 Agent 开始执行任务后，这里会显示运行状态、结构化事件和 transcript 入口。</p>
        </div>
      )}
      <RuntimeEventTimeline events={events} />
      <div className="task-runtime-toggle">
        <button
          className="action-button action-button--ghost"
          onClick={() => setShowTranscript((current) => !current)}
          type="button"
        >
          {showTranscript ? "收起原始 transcript" : "展开原始 transcript"}
        </button>
      </div>
      {showTranscript ? (
        <>
          <TranscriptInspector transcript={transcript} />
          <TranscriptAccessAuditTable audits={audits} />
        </>
      ) : (
        <div className="empty-card">
          <strong>默认只展示结构化运行事件。</strong>
          <p>需要查看 Agent 与 codex / claude / gemini 的对话记录时，再展开 transcript。</p>
        </div>
      )}
    </section>
  );
}
