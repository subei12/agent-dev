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
  const totalTaskCount = Object.values(boardsByMission).reduce((count, board) => count + (board?.items?.length ?? 0), 0);

  return (
    <PageShell
      eyebrow="总任务看板"
      title="Mission 总览"
      description="首页先展示所有 Mission 的任务进度，再进入单个 Mission 的作业空间。"
      aside={
        <div className="hero-stat">
          <span>任务总数</span>
          <strong>{totalTaskCount}</strong>
        </div>
      }
    >
      <MissionBoard missions={missionsQuery.data ?? []} boardsByMission={boardsByMission} />
      <section className="home-top-grid">
        <article className="panel">
          <div className="panel-header">
            <p className="panel-kicker">快速入口</p>
            <span className="badge">工作台</span>
          </div>
          <h2>从首页直接进入任务，而不是先看一堆对象卡片。</h2>
          <p className="panel-copy">首页现在优先展示 Mission 任务总览和快捷入口，方便直接切入具体作业。</p>
          <div className="action-row">
            <Link className="action-link" to="/missions/mission_1">
              打开任务页
            </Link>
            <Link className="action-link action-link--ghost" to="/agents">
              管理 Agent
            </Link>
          </div>
        </article>
        <article className="panel">
          <div className="panel-header">
            <p className="panel-kicker">运行入口</p>
            <span className="badge">日志</span>
          </div>
          <h2>运行日志和 transcript 入口放在固定位置。</h2>
          <p className="panel-copy">如果要排查 Agent 执行过程，可以直接进入运行会话查看事件流、转录和访问审计。</p>
          <div className="action-row">
            <Link className="action-link action-link--ghost" to="/runtime/sessions/session_1">
              查看运行会话
            </Link>
          </div>
        </article>
      </section>
    </PageShell>
  );
}
