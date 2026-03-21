import { Link } from "react-router-dom";

import { PageShell } from "../components/page-shell";

/**
 * HomePage 渲染当前路由对应的页面级工作区。
 */
export function HomePage() {
  return (
    <PageShell
      eyebrow="Operations Deck"
      title="Mission Control"
      description="A warm, editorial control room for multi-agent planning, runtime telemetry, transcript review, and final archive decisions."
      aside={
        <div className="hero-stat">
          <span>Current Mode</span>
          <strong>Governance-first</strong>
        </div>
      }
    >
      <article className="panel panel--feature">
        <div className="panel-header">
          <p className="panel-kicker">Launchpad</p>
          <span className="badge">v2 workspace</span>
        </div>
        <h2>Move from discussion to implementation without losing runtime visibility.</h2>
        <p className="panel-copy">
          The shell now has dedicated routes for mission workspace views and per-session transcript inspection.
        </p>
        <div className="action-row">
          <Link className="action-link" to="/missions/mission_1">
            Open Mission Workspace
          </Link>
          <Link className="action-link action-link--ghost" to="/runtime/sessions/session_1">
            Inspect Runtime Session
          </Link>
        </div>
      </article>
    </PageShell>
  );
}
