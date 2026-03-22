import { NavLink, Outlet } from "react-router-dom";

const navItems = [
  { to: "/", label: "首页" },
  { to: "/missions/mission_1", label: "工作台" },
  { to: "/runtime/sessions/session_1", label: "运行态" }
];

/**
 * AppLayout 渲染或处理当前前端行为。
 */
export function AppLayout() {
  return (
    <div className="app-frame">
      <aside className="sidebar">
        <div className="brand-block">
          <p className="brand-kicker">Agent Platform v2</p>
          <h1>任务指挥台</h1>
          <p className="brand-copy">
            围绕讨论、文档、任务接力、运行态观测和归档决策组织起来的多 Agent 开发平台。
          </p>
        </div>

        <nav aria-label="Primary" className="primary-nav">
          {navItems.map((item) => (
            <NavLink
              key={item.to}
              to={item.to}
              className={({ isActive }) => (isActive ? "nav-link nav-link--active" : "nav-link")}
            >
              {item.label}
            </NavLink>
          ))}
        </nav>

        <div className="status-ribbon">
          <span>在线</span>
          <strong>SSE + 转录记录</strong>
        </div>
      </aside>

      <main className="main-panel">
        <Outlet />
      </main>
    </div>
  );
}
