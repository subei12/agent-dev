import { MissionProgressStep } from "../../lib/mission-progress";

type MissionStageChartProps = {
  steps: MissionProgressStep[];
};

/**
 * MissionStageChart 用可视化轨迹展示 Mission 当前自动推进到哪一步。
 */
export function MissionStageChart({ steps }: MissionStageChartProps) {
  return (
    <div className="stage-chart">
      <div className="panel-header panel-header--compact">
        <p className="panel-kicker">阶段推进图</p>
      </div>
      <div className="stage-chart__rail">
        {steps.map((step) => (
          <div key={step.key} className="stage-chart__step">
            <div
              className={
                step.state === "done"
                  ? "stage-chart__dot stage-chart__dot--done"
                  : step.state === "active"
                    ? "stage-chart__dot stage-chart__dot--active"
                    : "stage-chart__dot"
              }
            />
            <strong>{step.label}</strong>
            <p>{step.detail}</p>
          </div>
        ))}
      </div>
    </div>
  );
}
