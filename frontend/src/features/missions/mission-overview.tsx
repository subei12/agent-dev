import { Mission } from "../../lib/api";

/**
 * MissionOverview 渲染或处理当前前端行为。
 */
export function MissionOverview({ mission }: { mission: Mission | null }) {
  return (
    <article className="panel panel--feature">
      <div className="panel-header">
        <p className="panel-kicker">Mission Snapshot</p>
        <span className="badge">{mission?.status ?? "offline snapshot"}</span>
      </div>
      <h2>{mission?.title ?? "Mission workspace is ready"}</h2>
      <p className="panel-copy">
        {mission?.description ??
          "Backend data is optional during local shell work. The interface keeps rendering and labels the current view when the API is offline."}
      </p>
    </article>
  );
}
