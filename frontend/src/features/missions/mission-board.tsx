import { Link } from "react-router-dom";

import { Mission, TaskBoard } from "../../lib/api";
import { formatStatusLabel, formatTypeLabel } from "../../lib/display";

type MissionBoardProps = {
  missions: Mission[];
  boardsByMission: Record<string, TaskBoard | null>;
};

/**
 * MissionBoard 展示首页总任务看板，帮助用户快速看到每个 Mission 的进度和入口。
 */
export function MissionBoard({ missions, boardsByMission }: MissionBoardProps) {
  return (
    <section className="panel panel--feature">
      <div className="panel-header panel-header--compact">
        <p className="panel-kicker">总任务看板</p>
        <span className="badge">{missions.length} 个 Mission</span>
      </div>
      <div className="mission-board-grid">
        {missions.map((mission, index) => {
          const board = boardsByMission[mission.id];
          const tasks = board?.items ?? [];
          const pendingCount = tasks.filter((task) => task.status !== "done").length;
          return (
            <article key={mission.id} className={index === 0 ? "mission-card mission-card--featured" : "mission-card"}>
              <div className="list-card__meta">
                <span>{formatStatusLabel(mission.status)}</span>
                <span>{pendingCount}/{tasks.length || 0} 进行中</span>
              </div>
              <h3>{mission.title}</h3>
              <p>{mission.description ?? "暂无任务说明。"}</p>
              <div className="stack">
                {tasks.slice(0, 3).map((task) => (
                  <div key={task.id} className="mini-row">
                    <span>{task.title}</span>
                    <span>{formatTypeLabel(task.type)} · {formatStatusLabel(task.status)}</span>
                  </div>
                ))}
              </div>
              <div className="mission-card__footer">
                <span className="mission-card__hint">管理员可在 Mission 内决定是否继续分派或归档。</span>
                <Link className="action-link" to={`/missions/${mission.id}`}>
                  打开任务
                </Link>
              </div>
            </article>
          );
        })}
      </div>
    </section>
  );
}
