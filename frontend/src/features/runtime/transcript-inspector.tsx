import { TranscriptView } from "../../lib/api";

/**
 * TranscriptInspector 渲染或处理当前前端行为。
 */
export function TranscriptInspector({ transcript }: { transcript: TranscriptView | null }) {
  return (
    <section className="panel panel--transcript">
      <div className="panel-header">
        <p className="panel-kicker">Transcript Inspector</p>
        <span className="badge">{transcript?.transcript.status ?? "redacted view"}</span>
      </div>
      <div className="transcript-box">
        {transcript?.entries?.length ? (
          transcript.entries.map((entry) => (
            <article key={entry.id} className="transcript-entry">
              <div className="transcript-entry__meta">
                <span>{entry.role}</span>
                <span>{entry.entryType}</span>
              </div>
              <p>{entry.redactedText ?? "Structured content only."}</p>
            </article>
          ))
        ) : (
          <div className="empty-card">
            <strong>Transcript is hidden until requested.</strong>
            <p>The default mission page shows structured events first and only expands raw dialogue on demand.</p>
          </div>
        )}
      </div>
    </section>
  );
}
