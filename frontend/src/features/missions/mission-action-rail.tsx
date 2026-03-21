type MissionActionRailProps = {
  archiveStatus: string;
  onArchive: () => void;
  isArchiving: boolean;
};

export function MissionActionRail({ archiveStatus, onArchive, isArchiving }: MissionActionRailProps) {
  return (
    <section className="panel">
      <div className="panel-header">
        <p className="panel-kicker">Mission Actions</p>
        <span className="badge">admin controls</span>
      </div>
      <div className="stack">
        <div className="empty-card">
          <strong>Archive status: {archiveStatus}</strong>
          <p>Archive the current mission package once the admin agent is satisfied with the latest artifacts.</p>
        </div>
        <button className="action-button" onClick={onArchive} disabled={isArchiving} type="button">
          Archive Mission
        </button>
      </div>
    </section>
  );
}
