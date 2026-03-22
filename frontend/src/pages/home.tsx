import { Link } from "react-router-dom";

import { PageShell } from "../components/page-shell";

/**
 * HomePage 渲染当前路由对应的页面级工作区。
 */
export function HomePage() {
  return (
    <PageShell
      eyebrow="运行看板"
      title="任务指挥台"
      description="面向多 Agent 讨论、执行、运行日志和归档决策的统一工作入口。"
      aside={
        <div className="hero-stat">
          <span>当前模式</span>
          <strong>治理优先</strong>
        </div>
      }
    >
      <article className="panel panel--feature">
        <div className="panel-header">
          <p className="panel-kicker">启动入口</p>
          <span className="badge">v2 工作台</span>
        </div>
        <h2>从讨论、拆解到执行与归档，都在同一个工作台里完成。</h2>
        <p className="panel-copy">
          当前页面已经接入 Mission 工作台和运行会话检视页，方便直接联调核心流程。
        </p>
        <div className="action-row">
          <Link className="action-link" to="/missions/mission_1">
            打开 Mission 工作台
          </Link>
          <Link className="action-link action-link--ghost" to="/runtime/sessions/session_1">
            查看运行会话
          </Link>
        </div>
      </article>
    </PageShell>
  );
}
