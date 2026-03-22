import { useQueries, useQuery } from "@tanstack/react-query";
import { Link } from "react-router-dom";

import { PageShell } from "../components/page-shell";
import { getMissions, getTaskBoard } from "../lib/api";
import { MissionBoard } from "../features/missions/mission-board";

const projectId = "proj_1";

/**
 * HomePage 渲染当前路由对应的页面级工作区。
 */
export function HomePage() {
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

  return (
    <PageShell
      eyebrow="平台总览"
      title="多 Agent 协同控制台"
      description="用更轻的工作台视图统一查看 Mission 进度、运行日志、共享文档和管理员决策入口。"
      aside={
        <div className="hero-stat-grid">
          <div className="hero-stat">
            <span>任务总数</span>
            <strong>{totalTaskCount}</strong>
          </div>
          <div className="hero-stat">
            <span>活跃 Mission</span>
            <strong>{activeMissionCount}</strong>
          </div>
          <div className="hero-stat hero-stat--muted">
            <span>项目状态</span>
            <strong>{missionCount > 0 ? "协作中" : "等待初始化"}</strong>
          </div>
        </div>
      }
    >
      <section className="home-top-grid">
        <article className="panel panel--feature panel--hero-callout">
          <div className="panel-header panel-header--compact">
            <p className="panel-kicker">今天的任务重心</p>
            <span className="badge">首要入口</span>
          </div>
          <h2>先进入当前 Mission，再围绕任务推进、Agent 执行和文档产物完成协作。</h2>
          <p className="panel-copy">
            这里不再是平均铺开的后台卡片，而是一个强调主路径的工作入口。首页负责判断整体进度，Mission 页负责推进具体任务。
          </p>
          <div className="action-row">
            <Link className="action-link" to="/missions/mission_1">
              进入 Mission 工作台
            </Link>
            <Link className="action-link action-link--ghost" to="/agents">
              配置 Agent
            </Link>
          </div>
        </article>
        <article className="panel home-side-panel">
          <div className="panel-header panel-header--compact">
            <p className="panel-kicker">关键入口</p>
            <span className="badge">运行与检查</span>
          </div>
          <div className="stack">
            <div className="compact-card">
              <strong>运行会话</strong>
              <p>优先看结构化事件，必要时再展开 transcript。</p>
              <Link className="inline-link" to="/runtime/sessions/session_1">
                打开运行会话
              </Link>
            </div>
            <div className="compact-card">
              <strong>管理员决策</strong>
              <p>Agent 配置、审批与归档状态集中放在辅助区，不再抢主任务视觉焦点。</p>
            </div>
          </div>
        </article>
      </section>
      <MissionBoard missions={missionsQuery.data ?? []} boardsByMission={boardsByMission} />
    </PageShell>
  );
}
