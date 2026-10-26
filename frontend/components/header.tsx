"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { useEffect, useState } from "react";
import { routes, pageHref, type Route } from "@/lib/routes";
import { useAuth } from "./auth";
import { LogoutButton } from "./logout-button";

type Entry = { title: string; children: Route[] } | Route;
const entries: Entry[] = [];
for (const route of routes.filter(r => !["start", "login"].includes(r.slug))) {
  if (!route.group) { entries.push(route); continue; }
  let group = entries.find((e): e is { title: string; children: Route[] } => "children" in e && e.title === route.group);
  if (!group) { group = { title: route.group, children: [] }; entries.push(group); }
  group.children.push(route);
}

export function Header() {
  const pathname = usePathname();
  const [open, setOpen] = useState(false);
  const auth = useAuth();
  useEffect(() => {
    setOpen(false);
    document.querySelectorAll<HTMLDetailsElement>(".nav-group[open]").forEach(el => { el.open = false; });
  }, [pathname]);
  return <header>
    <nav className="navbar navbar-expand-lg bg-white border-bottom" aria-label="Nawigacja główna">
      <div className="container">
        <Link className="navbar-brand fw-bold" href="/">Mirmiłowo <span className="brand-subtitle d-block fw-normal">Spółdzielnia Mieszkaniowa</span></Link>
        <button className="navbar-toggler" type="button" aria-controls="menu" aria-expanded={open} aria-label="Rozwiń menu" onClick={() => setOpen(!open)}><span className="navbar-toggler-icon" /></button>
        <div className={`collapse navbar-collapse${open ? " show" : ""}`} id="menu">
          <ul className="navbar-nav ms-auto gap-lg-1">
            {entries.map(entry => <li key={entry.title} className="nav-item">
              {"children" in entry ? <details className="nav-group">
                <summary className={`nav-link${entry.children.some(r => pageHref(r.slug) === pathname) ? " active" : ""}`}>{entry.title}</summary>
                <ul className="nav-submenu shadow-sm">
                  {entry.children.map(route => <li key={route.slug}><Link className={`dropdown-item${pathname === pageHref(route.slug) ? " active" : ""}`} href={pageHref(route.slug)} aria-current={pathname === pageHref(route.slug) ? "page" : undefined}>{route.title}</Link></li>)}
                </ul>
              </details> : <Link className={`nav-link${pathname === pageHref(entry.slug) ? " active" : ""}`} href={pageHref(entry.slug)} aria-current={pathname === pageHref(entry.slug) ? "page" : undefined}>{entry.title}</Link>}
            </li>)}
            <li className="nav-item ms-lg-3">
              {auth.loading ? <span className="nav-link text-secondary" role="status">Sprawdzanie…</span> : auth.user ? <LogoutButton className="nav-link text-danger fw-bold" /> : <Link className="nav-link text-primary fw-bold" href="/login" aria-current={pathname === "/login" ? "page" : undefined}>Zaloguj</Link>}
            </li>
          </ul>
        </div>
      </div>
    </nav>
    {auth.error && <div className="container small py-2" role="status">{auth.error} <button type="button" className="btn btn-sm btn-outline-secondary" disabled={auth.loading} onClick={() => void auth.reload()}>Spróbuj ponownie</button></div>}
  </header>;
}
