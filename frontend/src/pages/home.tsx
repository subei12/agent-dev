import { useMutation, useQueries, useQuery, useQueryClient } from "@tanstack/react-query";
import { Link, useNavigate } from "react-router-dom";

import { PageShell } from "../components/page-shell";
import { createMission, getMissions, getTaskBoard } from "../lib/api";
import { formatStatusLabel, formatTypeLabel } from "../lib/display";
import { MissionBoard } from "../features/missions/mission-board";
import { MissionRequestComposer } from "../features/missions/mission-request-composer";
import { MissionStatusChart } from "../features/missions/mission-status-chart";
import { MissionWorkbenchPreview } from "../features/missions/mission-workbench-preview";

const projectId = "proj_1";

/**
 * HomePage 渲染当前路由对应的页面级工作区。
 */
export function HomePage() {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const missionsQuery = useQuery({
    queryKey: ["missions"],
    queryFn: () => getMissions(projectId)
  });
  const boardQueries = useQueries({
    queries: (missionsQuery.data ?? []).map((mission) => ({
      queryKey: ["home-task-board", mission.id],
      queryFn: () => getTaskBoard(projectId, mission.id)
    }))
  });
  const boardsByMission = Object.fromEntries(
    (missionsQuery.data ?? []).map((mission, index) => [mission.id, boardQueries[index]?.data ?? null])
  );
  const missionCount = missionsQuery.data?.length ?? 0;
  const totalTaskCount = Object.values(boardsByMission).reduce((count, board) => count + (board?.items?.length ?? 0), 0);
  const activeMissionCount = (missionsQuery.data ?? []).filter((mission) => mission.status !== "done").length;
  const allTasks = (missionsQuery.data ?? []).flatMap((mission) =>
    (boardsByMission[mission.id]?.items ?? []).map((task) => ({
      ...task,
      missionTitle: mission.title
    }))
  );
  const reviewTaskCount = allTasks.filter((task) => task.status === "review" || task.status === "handoff_pending").length;
  const spotlightMission = missionsQuery.data?.[0] ?? null;
  const spotlightTasks = spotlightMission ? boardsByMission[spotlightMission.id]?.items ?? [] : [];

  const missionMutation = useMutation({
    mutationFn: (payload: { title: string; description: string }) => createMission(projectId, payload),
    onSuccess: async (mission) => {
      await queryClient.invalidateQueries({ queryKey: ["missions"] });
      navigate(`/missions/${mission.id}`);
    }
  });

  return (
    <PageShell
      eyebrow="总任务看板"
      title="任务指挥台"
      description="用户只提交一个需求，管理员 Agent 再组织设计、开发、测试和评审等内部执行任务。"
      aside={
        <div className="hero-stat">
          <span>任务总数</span>
          <strong>{totalTaskCount}</strong>
        </div>
      }
    >
      <section className="dashboard-hero-grid">
        <MissionRequestComposer
          isSubmitting={missionMutation.isPending}
          onCreateMission={(payload) => missionMutation.mutate(payload)}
        />
        <MissionWorkbenchPreview missionTitle={spotlightMission?.title ?? ""} board={boardsByMission[spotlightMission?.id ?? ""] ?? null} />
      </section>
      <section className="dashboard-summary-grid">
        <article className="panel summary-card">
          <p className="panel-kicker">活跃 Mission</p>
          <h2>{activeMissionCount}</h2>
          <p className="panel-copy">仍在推进中的项目任务。</p>
        </article>
        <article className="panel summary-card">
          <p className="panel-kicker">等待管理员</p>
          <h2>{reviewTaskCount}</h2>
          <p className="panel-copy">待评审或等待交接确认的任务。</p>
        </article>
        <article className="panel summary-card">
          <p className="panel-kicker">项目状态</p>
          <h2>{missionCount > 0 ? "协作中" : "待初始化"}</h2>
          <p className="panel-copy">演示项目已接入真实任务、文档和运行日志。</p>
        </article>
      </section>
      <section className="dashboard-main-grid">
        <MissionStatusChart missions={missionsQuery.data ?? []} boardsByMission={boardsByMission} />
        <article className="panel dashboard-pulse">
          <div className="panel-header">
            <p className="panel-kicker">任务脉冲</p>
            <span className="badge">{spotlightMission?.title ?? "暂无 Mission"}</span>
          </div>
          <h2>{spotlightMission?.title ?? "当前没有可展示的 Mission"}</h2>
          <p className="panel-copy">
            {spotlightMission?.description ?? "接入更多 Mission 后，这里会展示主任务的推进节奏。"}
          </p>
          <div className="dashboard-progress-list">
            {(missionsQuery.data ?? []).map((mission) => {
              const missionTasks = boardsByMission[mission.id]?.items ?? [];
              const completedCount = missionTasks.filter((task) => task.status === "done").length;
              const progress = missionTasks.length > 0 ? Math.max(12, (completedCount / missionTasks.length) * 100) : 12;

              return (
                <div key={mission.id} className="dashboard-progress-row">
                  <div>
                    <strong>{mission.title}</strong>
                    <p>{missionTasks.length} 个任务</p>
                  </div>
                  <div className="dashboard-progress-bar">
                    <span style={{ width: `${progress}%` }} />
                  </div>
                </div>
              );
            })}
          </div>
          <div className="action-row">
            <Link className="action-link" to="/missions/mission_1">
              查看当前 Mission
            </Link>
            <Link className="action-link action-link--ghost" to="/runtime/sessions/session_1">
              查看运行日志
            </Link>
          </div>
        </article>

        <article className="panel dashboard-side-card">
          <div className="panel-header">
            <p className="panel-kicker">最近动态</p>
            <span className="badge">{allTasks.slice(0, 4).length} 条</span>
          </div>
          <div className="stack">
            {allTasks.slice(0, 4).map((task) => (
              <div key={task.id} className="compact-card">
                <strong>{task.title}</strong>
                <p>{task.missionTitle}</p>
                <p>{task.assignedAgentId ?? "未分配"} · {formatStatusLabel(task.status)}</p>
              </div>
            ))}
          </div>
        </article>

        <article className="panel dashboard-side-card">
          <div className="panel-header">
            <p className="panel-kicker">近期任务列表</p>
            <span className="badge">{spotlightTasks.length} 个</span>
          </div>
          <div className="stack">
            {spotlightTasks.slice(0, 4).map((task) => (
              <div key={task.id} className="compact-card">
                <strong>{task.title}</strong>
                <p>{formatTypeLabel(task.type)} · {formatStatusLabel(task.status)}</p>
              </div>
            ))}
          </div>
        </article>
      </section>
      <MissionBoard missions={missionsQuery.data ?? []} boardsByMission={boardsByMission} />
    </PageShell>
  );
}
