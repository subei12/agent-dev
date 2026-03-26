import { Mission, TaskItemView } from "../../lib/api";
import { formatStatusLabel } from "../../lib/display";
import { deriveMissionProgress, summarizeMissionProgress } from "../../lib/mission-progress";
import { MissionStageChart } from "./mission-stage-chart";

/**
 * MissionOverview 渲染或处理当前前端行为。
 */
export function MissionOverview({
  mission,
  tasks
}: {
  mission: Mission | null;
  tasks: TaskItemView[];
}) {
  const steps = deriveMissionProgress(mission, tasks);
  const progressSummary = summarizeMissionProgress(steps);

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
          <div className="compact-card">
            <strong>自动推进状态</strong>
            <p>{progressSummary}</p>
            <p>平台会按收敛方案、开发实现、测试验收、管理员复核的顺序自动推进。</p>
          </div>
          <MissionStageChart steps={steps} />
        </div>
        <div className="mission-overview__facts">
          {steps.map((step) => (
            <div
              key={step.key}
              className={
                step.state === "done"
                  ? "compact-card progress-step progress-step--done"
                  : step.state === "active"
                    ? "compact-card progress-step progress-step--active"
                    : "compact-card progress-step"
              }
            >
              <strong>{step.label}</strong>
              <p>{step.detail}</p>
            </div>
          ))}
        </div>
      </div>
    </article>
  );
}
