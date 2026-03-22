import { TranscriptAccessAudit } from "../../lib/api";
import { formatAccessModeLabel } from "../../lib/display";

/**
 * TranscriptAccessAuditTable 渲染或处理当前前端行为。
 */
export function TranscriptAccessAuditTable({ audits }: { audits: TranscriptAccessAudit[] }) {
  return (
    <section className="panel">
      <div className="panel-header">
        <p className="panel-kicker">转录访问审计</p>
        <span className="badge">{audits.length} 条</span>
      </div>

      {audits.length === 0 ? (
        <div className="empty-card">
          <strong>当前还没有转录查看记录。</strong>
          <p>当有人打开脱敏 transcript 时，这里会记录访问行为。</p>
        </div>
      ) : (
        <div className="audit-table">
          <div className="audit-row audit-row--head">
            <span>操作人</span>
            <span>模式</span>
            <span>原因</span>
          </div>
          {audits.map((audit) => (
            <div key={audit.id} className="audit-row">
              <span>{audit.actorUserId}</span>
              <span>{formatAccessModeLabel(audit.accessMode)}</span>
              <span>{audit.reason ?? "未填写"}</span>
            </div>
          ))}
        </div>
      )}
    </section>
  );
}
