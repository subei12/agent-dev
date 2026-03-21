import { DiscussionSession } from "../../lib/api";

export function DiscussionSessionPanel({ sessions }: { sessions: DiscussionSession[] }) {
  return (
    <section className="panel">
      <div className="panel-header">
        <p className="panel-kicker">Discussion Rounds</p>
        <span className="badge">{sessions.length} sessions</span>
      </div>
      <div className="stack">
        {sessions.length === 0 ? (
          <div className="empty-card">
            <strong>No discussion sessions loaded.</strong>
            <p>The shell still reserves the section so agent review flows have a stable home.</p>
          </div>
        ) : (
          sessions.map((session) => (
            <article key={session.id} className="list-card">
              <h3>{session.topic}</h3>
              <p>{session.status}</p>
            </article>
          ))
        )}
      </div>
    </section>
  );
}
