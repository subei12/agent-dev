import { FormEvent, useState } from "react";

import { DiscussionSession } from "../../lib/api";

type DiscussionSessionPanelProps = {
  sessions: DiscussionSession[];
  onCreateSession: (topic: string) => void;
  isSubmitting: boolean;
};

/**
 * DiscussionSessionPanel 渲染或处理当前前端行为。
 */
export function DiscussionSessionPanel({ sessions, onCreateSession, isSubmitting }: DiscussionSessionPanelProps) {
  const [topic, setTopic] = useState("");

  const handleSubmit = (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (!topic.trim()) {
      return;
    }
    onCreateSession(topic.trim());
    setTopic("");
  };

  return (
    <section className="panel">
      <div className="panel-header">
        <p className="panel-kicker">Discussion Rounds</p>
        <span className="badge">{sessions.length} sessions</span>
      </div>
      <form className="inline-form" onSubmit={handleSubmit}>
        <label className="field">
          <span>Discussion topic</span>
          <input value={topic} onChange={(event) => setTopic(event.target.value)} />
        </label>
        <button className="action-button" disabled={isSubmitting} type="submit">
          Start Discussion
        </button>
      </form>
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
