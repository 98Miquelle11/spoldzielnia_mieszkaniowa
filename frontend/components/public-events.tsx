"use client";

import Link from "next/link";
import { useEffect, useState } from "react";
import { EventsError, fetchEvents, formatEventDate, isCulturalEvent, isEventPage, type EventPage, type CulturalEvent, type Period } from "@/lib/events";

type Result<T> = { state: "loading" } | { state: "success"; data: T } | { state: "error"; message: string };

function EventDates({ event }: { event: CulturalEvent }) {
  return <p className="text-secondary">
    <time dateTime={event.starts_at}>{formatEventDate(event.starts_at)}</time>
    {event.ends_at !== event.starts_at && <> – <time dateTime={event.ends_at}>{formatEventDate(event.ends_at)}</time></>}
    <span className="d-block small">Czas lokalny: Europe/Warsaw</span>
  </p>;
}

export function PublicEvents({ period }: { period: Period | null }) {
  const [result, setResult] = useState<Result<EventPage>>({ state: "loading" });
  const [retry, setRetry] = useState(0);
  const [page, setPage] = useState(1);
  useEffect(() => {
    if (!period) return;
    const controller = new AbortController();
    setResult({ state: "loading" });
    void fetchEvents(`?period=${period}&page=${page}&limit=20`, controller.signal).then(data => {
      if (!isEventPage(data) || data.pagination.page !== page) {
        throw new Error("Serwer zwrócił nieprawidłową listę wydarzeń.");
      }
      if (!controller.signal.aborted) setResult({ state: "success", data: data });
    }).catch(error => {
      if (!controller.signal.aborted) setResult({ state: "error", message: error instanceof Error ? error.message : "Nie udało się pobrać wydarzeń." });
    });
    return () => controller.abort();
  }, [period, page, retry]);
  return <section className="container py-5">
    <h1 className="h2">Imprezy kulturalne</h1>
    <nav aria-label="Filtr wydarzeń" className="nav nav-pills gap-2 my-4">
      <Link className={`nav-link${period === "upcoming" ? " active" : ""}`} aria-current={period === "upcoming" ? "page" : undefined} href="/imprezy?period=upcoming">Nadchodzące</Link>
      <Link className={`nav-link${period === "past" ? " active" : ""}`} aria-current={period === "past" ? "page" : undefined} href="/imprezy?period=past">Archiwalne</Link>
    </nav>
    {!period ? <p className="alert alert-danger" role="alert">Nieznany filtr. Wybierz „Nadchodzące” lub „Archiwalne”.</p>
      : result.state === "loading" ? <p role="status">Ładowanie wydarzeń…</p>
      : result.state === "error" ? <div role="alert" className="alert alert-danger"><p>{result.message}</p><button className="btn btn-outline-primary" onClick={() => setRetry(n => n + 1)}>Spróbuj ponownie</button></div>
      : <>
        {result.data.items.length === 0 && <p role="status">Brak {period === "past" ? "archiwalnych" : "nadchodzących"} wydarzeń.</p>}
        <div className="row g-4">{result.data.items.map(event => <div className="col-12 col-md-6 col-xl-4" key={event.id}>
          <article className="card h-100 shadow-sm"><div className="card-body d-flex flex-column">
            <h2 className="h5 text-break">{event.title}</h2>
            <EventDates event={event} />
            <p className="text-break"><strong>Miejsce:</strong> {event.location}</p>
            <p className="text-break" style={{ whiteSpace: "pre-wrap" }}>{Array.from(event.description).slice(0, 240).join("")}{Array.from(event.description).length > 240 ? "…" : ""}</p>
            <Link className="btn btn-outline-primary mt-auto align-self-start" href={`/imprezy/${event.id}?period=${period}`}>Szczegóły<span className="visually-hidden">: {event.title}</span></Link>
          </div></article>
        </div>)}</div>
        <nav aria-label="Strony wydarzeń" className="d-flex align-items-center gap-3 mt-4">
          <button className="btn btn-outline-primary" disabled={page === 1} onClick={() => setPage(n => n - 1)}>Poprzednia strona</button>
          <span aria-live="polite">Strona {page}</span>
          <button className="btn btn-primary" disabled={!result.data.pagination.has_more} onClick={() => setPage(n => n + 1)}>Następna strona</button>
        </nav>
      </>}
  </section>;
}

export function PublicEventDetails({ id, period }: { id: string; period: Period }) {
  const [result, setResult] = useState<Result<CulturalEvent>>({ state: "loading" });
  const [retry, setRetry] = useState(0);
  useEffect(() => {
    const controller = new AbortController();
    setResult({ state: "loading" });
    if (!/^\d+$/.test(id) || !Number.isSafeInteger(Number(id)) || Number(id) <= 0) {
      setResult({ state: "error", message: "Nieprawidłowy identyfikator wydarzenia." });
      return;
    }
    void fetchEvents(`/${id}`, controller.signal).then(data => {
      if (!isCulturalEvent(data)) throw new EventsError("Serwer zwrócił nieprawidłowe dane wydarzenia.", 0);
      if (!controller.signal.aborted) setResult({ state: "success", data });
    }).catch(error => {
      if (!controller.signal.aborted) setResult({ state: "error", message: error instanceof Error ? error.message : "Nie udało się pobrać wydarzenia." });
    });
    return () => controller.abort();
  }, [id, retry]);
  return <section className="container py-5">
    <Link href={`/imprezy?period=${period}`} className="btn btn-outline-primary mb-4">Wróć do {period === "past" ? "archiwum" : "nadchodzących wydarzeń"}</Link>
    {result.state === "loading" ? <p role="status">Ładowanie wydarzenia…</p>
      : result.state === "error" ? <><h1 className="h2">Wydarzenie niedostępne</h1><p role="alert">{result.message}</p><button className="btn btn-primary" onClick={() => setRetry(n => n + 1)}>Spróbuj ponownie</button></>
      : <article className="card shadow-sm"><div className="card-body p-4 p-md-5">
        <h1 className="h2 text-break">{result.data.title}</h1>
        <EventDates event={result.data} />
        <p className="text-break"><strong>Miejsce:</strong> {result.data.location}</p>
        <p className="text-break" style={{ whiteSpace: "pre-wrap" }}>{result.data.description}</p>
      </div></article>}
  </section>;
}
