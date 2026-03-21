import { PropsWithChildren, ReactNode } from "react";

type PageShellProps = PropsWithChildren<{
  eyebrow: string;
  title: string;
  description: string;
  aside?: ReactNode;
}>;

export function PageShell({ eyebrow, title, description, aside, children }: PageShellProps) {
  return (
    <section className="page-shell">
      <header className="hero-panel">
        <div>
          <p className="eyebrow">{eyebrow}</p>
          <h1>{title}</h1>
          <p className="hero-copy">{description}</p>
        </div>
        {aside ? <div className="hero-aside">{aside}</div> : null}
      </header>
      <div className="page-grid">{children}</div>
    </section>
  );
}
