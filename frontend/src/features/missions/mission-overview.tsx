import { Mission } from "../../lib/api";
import { formatStatusLabel } from "../../lib/display";

/**
 * MissionOverview 渲染或处理当前前端行为。
 */
export function MissionOverview({ mission }: { mission: Mission | null }) {
  return (
    <article className="panel panel--feature mission-overview-card">
      <div className="panel-header panel-header--compact">
        <p className="panel-kicker">当前 Mission</p>
        <span className="badge">{mission ? formatStatusLabel(mission.status) : "离线快照"}</span>
      </div>
      <div className="mission-overview__grid">
        <div>
          <h2>任务推进摘要</h2>
          <p className="panel-copy">
            {mission?.description ?? "当前页面在后端不可用时也会保留结构，方便联调布局和交互。"}
          </p>
        </div>
        <div className="mission-overview__facts">
          <div className="compact-card">
            <strong>任务推进优先</strong>
            <p>第一屏优先暴露新增任务、看板选择和关键操作。</p>
          </div>
          <div className="compact-card">
            <strong>运行事件优先</strong>
            <p>默认先看结构化事件，只有在排查时才展开 transcript。</p>
          </div>
        </div>
      </div>
    </article>
  );
}
