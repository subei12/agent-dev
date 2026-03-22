import { NavLink, Outlet } from "react-router-dom";

const navItems = [
  { to: "/", label: "首页" },
  { to: "/missions/mission_1", label: "Mission" },
  { to: "/runtime/sessions/session_1", label: "运行会话" },
  { to: "/agents", label: "Agent 配置" }
];

/**
 * AppLayout 渲染或处理当前前端行为。
 */
export function AppLayout() {
  return (
    <div className="app-shell">
      <header className="topbar">
        <div className="topbar__brand">
          <p className="brand-kicker">Agent Platform v2</p>
          <strong>多 Agent 协同开发工作台</strong>
        </div>

        <nav aria-label="Primary" className="topbar__nav">
          {navItems.map((item) => (
            <NavLink
              key={item.to}
              to={item.to}
              className={({ isActive }) => (isActive ? "topbar__link topbar__link--active" : "topbar__link")}
            >
              {item.label}
            </NavLink>
          ))}
        </nav>

        <div className="topbar__meta">
          <span className="topbar__meta-dot" />
          <div>
            <span>在线协作</span>
            <strong>SSE + transcript</strong>
          </div>
        </div>
      </header>

      <main className="main-panel">
        <Outlet />
      </main>
    </div>
  );
}
