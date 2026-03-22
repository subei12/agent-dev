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
      <div className="panel-header">
        <p className="panel-kicker">项目任务列表</p>
        <span className="badge">{missions.length} 个 Mission</span>
      </div>
      <div className="mission-board-list">
        {missions.map((mission) => {
          const board = boardsByMission[mission.id];
          const tasks = board?.items ?? [];
          const pendingCount = tasks.filter((task) => task.status !== "done").length;

          return (
            <article key={mission.id} className="mission-row-card">
              <div className="mission-row-card__main">
                <div className="list-card__meta">
                  <span>{formatStatusLabel(mission.status)}</span>
                  <span>{pendingCount}/{tasks.length || 0} 处理中</span>
                </div>
                <h3>{mission.title}</h3>
                <p>{mission.description ?? "暂无任务说明。"}</p>
              </div>
              <div className="mission-row-card__tasks">
                {tasks.slice(0, 3).map((task) => (
                  <div key={task.id} className="mini-row">
                    <span>{task.title}</span>
                    <span>{formatTypeLabel(task.type)} · {formatStatusLabel(task.status)}</span>
                  </div>
                ))}
              </div>
              <div className="mission-row-card__action">
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
