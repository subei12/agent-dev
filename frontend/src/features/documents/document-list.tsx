import { FormEvent, useState } from "react";

import { DocumentItem } from "../../lib/api";

type DocumentListProps = {
  documents: DocumentItem[];
  onCreateDocument: (title: string, kind: string) => void;
  isSubmitting: boolean;
};

export function DocumentList({ documents, onCreateDocument, isSubmitting }: DocumentListProps) {
  const [title, setTitle] = useState("");
  const [kind, setKind] = useState("architecture");

  const handleSubmit = (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (!title.trim()) {
      return;
    }
    onCreateDocument(title.trim(), kind);
    setTitle("");
  };

  return (
    <section className="panel">
      <div className="panel-header">
        <p className="panel-kicker">Documents</p>
        <span className="badge">{documents.length} tracked</span>
      </div>
      <form className="inline-form inline-form--wide" onSubmit={handleSubmit}>
        <label className="field">
          <span>Document title</span>
          <input value={title} onChange={(event) => setTitle(event.target.value)} />
        </label>
        <label className="field field--compact">
          <span>Document kind</span>
          <select value={kind} onChange={(event) => setKind(event.target.value)}>
            <option value="architecture">architecture</option>
            <option value="implementation_plan">implementation_plan</option>
            <option value="review_report">review_report</option>
          </select>
        </label>
        <button className="action-button" disabled={isSubmitting} type="submit">
          Add Document
        </button>
      </form>
      <div className="stack">
        {documents.length === 0 ? (
          <div className="empty-card">
            <strong>No adopted documents yet.</strong>
            <p>Architecture, implementation plans, and review reports will appear here once the API serves them.</p>
          </div>
        ) : (
          documents.map((document) => (
            <article key={document.id} className="list-card">
              <div className="list-card__meta">
                <span>{document.kind}</span>
                <span>{document.currentAdoptedVersionId ?? "draft"}</span>
              </div>
              <h3>{document.title}</h3>
            </article>
          ))
        )}
      </div>
    </section>
  );
}
