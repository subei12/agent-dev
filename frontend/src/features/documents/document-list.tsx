import { DocumentItem } from "../../lib/api";

export function DocumentList({ documents }: { documents: DocumentItem[] }) {
  return (
    <section className="panel">
      <div className="panel-header">
        <p className="panel-kicker">Documents</p>
        <span className="badge">{documents.length} tracked</span>
      </div>
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
