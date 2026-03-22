import { Mission } from "../../lib/api";
import { formatStatusLabel } from "../../lib/display";

/**
 * MissionOverview 渲染或处理当前前端行为。
 */
export function MissionOverview({ mission }: { mission: Mission | null }) {
  return (
    <article className="panel panel--feature">
      <div className="panel-header">
        <p className="panel-kicker">任务概览</p>
        <span className="badge">{mission ? formatStatusLabel(mission.status) : "离线快照"}</span>
      </div>
      <h2>{mission?.title ?? "Mission 工作台已就绪"}</h2>
      <p className="panel-copy">
        {mission?.description ??
          "当前页面在后端不可用时也会保留结构，方便联调布局和交互。"}
      </p>
    </article>
  );
}
