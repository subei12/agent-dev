import { NavLink, Outlet } from "react-router-dom";

const navItems = [
  { to: "/", label: "Missions" },
  { to: "/missions/mission_1", label: "Workspace" },
  { to: "/runtime/sessions/session_1", label: "Runtime" }
];

export function AppLayout() {
  return (
    <div className="app-frame">
      <aside className="sidebar">
        <div className="brand-block">
          <p className="brand-kicker">Agent Platform v2</p>
          <h1>Mission Control</h1>
          <p className="brand-copy">
            Editorial command deck for discussions, code handoffs, runtime telemetry, and final archive
            decisions.
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
          <span>Live</span>
          <strong>SSE + Transcript</strong>
        </div>
      </aside>

      <main className="main-panel">
        <Outlet />
      </main>
    </div>
  );
}
