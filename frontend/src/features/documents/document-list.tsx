import { FormEvent, useState } from "react";

import { DocumentItem, DocumentVersion } from "../../lib/api";

type DocumentListProps = {
  documents: DocumentItem[];
  versionsByDocument: Record<string, DocumentVersion[]>;
  onCreateDocument: (title: string, kind: string) => void;
  onCreateVersion: (documentId: string, contentText: string) => void;
  onAdoptVersion: (documentId: string, versionId: string) => void;
  isSubmitting: boolean;
};

/**
 * DocumentList 渲染或处理当前前端行为。
 */
export function DocumentList({
  documents,
  versionsByDocument,
  onCreateDocument,
  onCreateVersion,
  onAdoptVersion,
  isSubmitting
}: DocumentListProps) {
  const [title, setTitle] = useState("");
  const [kind, setKind] = useState("architecture");
  const [versionDrafts, setVersionDrafts] = useState<Record<string, string>>({});

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
              <div className="stack">
                <label className="field">
                  <span>{`Version content for ${document.title}`}</span>
                  <textarea
                    className="field-textarea"
                    value={versionDrafts[document.id] ?? ""}
                    onChange={(event) =>
                      setVersionDrafts((current) => ({
                        ...current,
                        [document.id]: event.target.value
                      }))
                    }
                  />
                </label>
                <div className="task-actions">
                  <button
                    className="action-button"
                    disabled={isSubmitting}
                    type="button"
                    onClick={() => {
                      const contentText = versionDrafts[document.id]?.trim();
                      if (!contentText) {
                        return;
                      }
                      onCreateVersion(document.id, contentText);
                      setVersionDrafts((current) => ({
                        ...current,
                        [document.id]: ""
                      }));
                    }}
                  >
                    Save Version
                  </button>
                </div>
                <div className="stack">
                  {(versionsByDocument[document.id] ?? []).map((version) => (
                    <div key={version.id} className="version-card">
                      <div className="list-card__meta">
                        <span>{`v${version.version}`}</span>
                        <span>{version.status}</span>
                      </div>
                      <p>{version.contentText}</p>
                      {version.status !== "adopted" ? (
                        <button
                          className="action-button action-button--ghost"
                          disabled={isSubmitting}
                          type="button"
                          onClick={() => onAdoptVersion(document.id, version.id)}
                        >
                          Adopt Version
                        </button>
                      ) : null}
                    </div>
                  ))}
                </div>
              </div>
            </article>
          ))
        )}
      </div>
    </section>
  );
}
