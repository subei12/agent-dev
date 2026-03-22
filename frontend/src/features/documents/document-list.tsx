import { FormEvent, useState } from "react";

import { DocumentItem, DocumentVersion } from "../../lib/api";
import { formatStatusLabel, formatTypeLabel } from "../../lib/display";

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
        <p className="panel-kicker">文档</p>
        <span className="badge">{documents.length} 份</span>
      </div>
      <form className="inline-form inline-form--wide" onSubmit={handleSubmit}>
        <label className="field">
          <span>文档标题</span>
          <input value={title} onChange={(event) => setTitle(event.target.value)} />
        </label>
        <label className="field field--compact">
          <span>文档类型</span>
          <select value={kind} onChange={(event) => setKind(event.target.value)}>
            <option value="architecture">架构文档</option>
            <option value="implementation_plan">实施计划</option>
            <option value="review_report">评审报告</option>
          </select>
        </label>
        <button className="action-button" disabled={isSubmitting} type="submit">
          新建文档
        </button>
      </form>
      <div className="stack">
        {documents.length === 0 ? (
          <div className="empty-card">
            <strong>当前还没有已采纳文档。</strong>
            <p>架构文档、实施计划和评审报告会在这里统一展示。</p>
          </div>
        ) : (
          documents.map((document) => (
            <article key={document.id} className="list-card">
              <div className="list-card__meta">
                <span>{formatTypeLabel(document.kind)}</span>
                <span>{document.currentAdoptedVersionId ?? "草稿"}</span>
              </div>
              <h3>{document.title}</h3>
              <div className="stack">
                <label className="field">
                  <span>{`${document.title} 的版本内容`}</span>
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
                    保存版本
                  </button>
                </div>
                <div className="stack">
                  {(versionsByDocument[document.id] ?? []).map((version) => (
                    <div key={version.id} className="version-card">
                      <div className="list-card__meta">
                        <span>{`v${version.version}`}</span>
                        <span>{formatStatusLabel(version.status)}</span>
                      </div>
                      <p>{version.contentText}</p>
                      {version.status !== "adopted" ? (
                        <button
                          className="action-button action-button--ghost"
                          disabled={isSubmitting}
                          type="button"
                          onClick={() => onAdoptVersion(document.id, version.id)}
                        >
                          采纳版本
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
