import { FormEvent, useState } from "react";

import { DiscussionSession } from "../../lib/api";
import { formatStatusLabel } from "../../lib/display";

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
        <p className="panel-kicker">讨论轮次</p>
        <span className="badge">{sessions.length} 条会话</span>
      </div>
      <form className="inline-form" onSubmit={handleSubmit}>
        <label className="field">
          <span>讨论主题</span>
          <input value={topic} onChange={(event) => setTopic(event.target.value)} />
        </label>
        <button className="action-button" disabled={isSubmitting} type="submit">
          发起讨论
        </button>
      </form>
      <div className="stack">
        {sessions.length === 0 ? (
          <div className="empty-card">
            <strong>当前还没有讨论会话。</strong>
            <p>创建第一轮讨论后，这里会展示各轮讨论的主题和状态。</p>
          </div>
        ) : (
          sessions.map((session) => (
            <article key={session.id} className="list-card">
              <h3>{session.topic}</h3>
              <p>{formatStatusLabel(session.status)}</p>
            </article>
          ))
        )}
      </div>
    </section>
  );
}
