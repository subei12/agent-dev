type ArchiveInfo = {
  id: string;
  status: string;
  manifestObjectKey?: string;
  bundleObjectKey?: string;
};

export function MissionArchivePanel({ archive }: { archive: ArchiveInfo | null }) {
  return (
    <section className="panel">
      <div className="panel-header">
        <p className="panel-kicker">归档状态</p>
        <span className="badge">{archive?.status ?? "not archived"}</span>
      </div>
      {archive ? (
        <div className="stack">
          <div className="empty-card">
            <strong>{archive.id}</strong>
            <p>{archive.manifestObjectKey ?? "manifest 未写入"}</p>
            <p>{archive.bundleObjectKey ?? "bundle 未写入"}</p>
          </div>
        </div>
      ) : (
        <div className="empty-card">
          <strong>尚未生成归档。</strong>
          <p>当管理员触发归档后，这里会展示归档对象键和当前状态。</p>
        </div>
      )}
    </section>
  );
}
