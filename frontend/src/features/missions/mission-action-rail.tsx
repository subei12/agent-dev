import { formatStatusLabel } from "../../lib/display";

type MissionActionRailProps = {
  archiveStatus: string;
  onArchive: () => void;
  isArchiving: boolean;
};

/**
 * MissionActionRail 渲染或处理当前前端行为。
 */
export function MissionActionRail({ archiveStatus, onArchive, isArchiving }: MissionActionRailProps) {
  return (
    <section className="panel">
      <div className="panel-header">
        <p className="panel-kicker">任务操作</p>
        <span className="badge">管理员操作</span>
      </div>
      <div className="stack">
        <div className="empty-card">
          <strong>归档状态：{formatStatusLabel(archiveStatus)}</strong>
          <p>当管理员确认当前文档、任务和运行结果满足要求后，可在这里触发归档。</p>
        </div>
        <button className="action-button" onClick={onArchive} disabled={isArchiving} type="button">
          归档任务
        </button>
      </div>
    </section>
  );
}
