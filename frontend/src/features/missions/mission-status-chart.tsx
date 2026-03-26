import { Mission, TaskBoard } from "../../lib/api";
import { formatStatusLabel } from "../../lib/display";

type MissionStatusChartProps = {
  missions: Mission[];
  boardsByMission: Record<string, TaskBoard | null>;
};

/**
 * MissionStatusChart 展示首页总览里的平台状态趋势图。
 */
export function MissionStatusChart({ missions, boardsByMission }: MissionStatusChartProps) {
  const counts = [
    {
      label: "需求收敛",
      value: missions.filter((mission) => mission.status === "discussion").length
    },
    {
      label: "开发推进",
      value: missions.filter((mission) => mission.status === "implementation").length
    },
    {
      label: "已完成",
      value: missions.filter((mission) => mission.status === "completed").length
    }
  ];
  const maxValue = Math.max(...counts.map((item) => item.value), 1);

  return (
    <section className="panel dashboard-chart-card">
      <div className="panel-header">
        <p className="panel-kicker">平台状态趋势</p>
        <span className="badge">{missions.length} 个 Mission</span>
      </div>
      <div className="dashboard-chart">
        {counts.map((item) => (
          <div key={item.label} className="dashboard-chart__row">
            <div className="dashboard-chart__meta">
              <strong>{item.label}</strong>
              <span>{item.value}</span>
            </div>
            <div className="dashboard-chart__bar">
              <span style={{ width: `${(item.value / maxValue) * 100}%` }} />
            </div>
          </div>
        ))}
      </div>
      <div className="dashboard-chart__legend">
        {missions.slice(0, 3).map((mission) => {
          const taskCount = boardsByMission[mission.id]?.items?.length ?? 0;
          return (
            <div key={mission.id} className="compact-card">
              <strong>{mission.title}</strong>
              <p>{formatStatusLabel(mission.status)} · {taskCount} 个内部任务</p>
            </div>
          );
        })}
      </div>
    </section>
  );
}
