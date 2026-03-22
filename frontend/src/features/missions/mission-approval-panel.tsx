type ApprovalItem = {
  id: string;
  action: string;
  status: string;
  subjectType: string;
};

export function MissionApprovalPanel({ approvals }: { approvals: ApprovalItem[] }) {
  return (
    <section className="panel">
      <div className="panel-header">
        <p className="panel-kicker">审批列表</p>
        <span className="badge">{approvals.length} 条</span>
      </div>
      <div className="stack">
        {approvals.length === 0 ? (
          <div className="empty-card">
            <strong>当前 Mission 没有审批记录。</strong>
            <p>创建发布审批后，这里会展示动作、对象和状态。</p>
          </div>
        ) : (
          approvals.map((approval) => (
            <article key={approval.id} className="list-card">
              <div className="list-card__meta">
                <span>{approval.action}</span>
                <span>{approval.status}</span>
              </div>
              <h3>{approval.id}</h3>
              <p>{approval.subjectType}</p>
            </article>
          ))
        )}
      </div>
    </section>
  );
}
