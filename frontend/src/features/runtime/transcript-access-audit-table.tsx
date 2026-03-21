import { TranscriptAccessAudit } from "../../lib/api";

/**
 * TranscriptAccessAuditTable 渲染或处理当前前端行为。
 */
export function TranscriptAccessAuditTable({ audits }: { audits: TranscriptAccessAudit[] }) {
  return (
    <section className="panel">
      <div className="panel-header">
        <p className="panel-kicker">Transcript Access Audit</p>
        <span className="badge">{audits.length} records</span>
      </div>

      {audits.length === 0 ? (
        <div className="empty-card">
          <strong>No transcript views have been recorded yet.</strong>
          <p>Audit rows appear once viewers open the redacted transcript view.</p>
        </div>
      ) : (
        <div className="audit-table">
          <div className="audit-row audit-row--head">
            <span>Actor</span>
            <span>Mode</span>
            <span>Reason</span>
          </div>
          {audits.map((audit) => (
            <div key={audit.id} className="audit-row">
              <span>{audit.actorUserId}</span>
              <span>{audit.accessMode}</span>
              <span>{audit.reason ?? "not provided"}</span>
            </div>
          ))}
        </div>
      )}
    </section>
  );
}
