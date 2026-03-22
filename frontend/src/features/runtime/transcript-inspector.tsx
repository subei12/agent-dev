import { TranscriptView } from "../../lib/api";

/**
 * TranscriptInspector 渲染或处理当前前端行为。
 */
export function TranscriptInspector({ transcript }: { transcript: TranscriptView | null }) {
  return (
    <section className="panel panel--transcript">
      <div className="panel-header">
        <p className="panel-kicker">转录查看器</p>
        <span className="badge">{transcript?.transcript.status ?? "脱敏视图"}</span>
      </div>
      <div className="transcript-box">
        {transcript?.entries?.length ? (
          transcript.entries.map((entry) => (
            <article key={entry.id} className="transcript-entry">
              <div className="transcript-entry__meta">
                <span>{entry.role}</span>
                <span>{entry.entryType}</span>
              </div>
              <p>{entry.redactedText ?? "当前仅返回结构化内容。"}</p>
            </article>
          ))
        ) : (
          <div className="empty-card">
            <strong>当前未加载 transcript。</strong>
            <p>默认优先展示结构化事件，只有在需要时才展开脱敏后的转录内容。</p>
          </div>
        )}
      </div>
    </section>
  );
}
