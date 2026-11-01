"use client";

import Link from "next/link";
import { useCallback, useEffect, useState, type FormEvent } from "react";
import { useAuth } from "./auth";

type Event = {
  id: number;
  title: string;
  description: string;
  starts_at: string;
  ends_at: string;
  location: string;
  is_published: boolean;
};

type EventInput = Omit<Event, "id">;

async function adminRequest<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(`/api/admin/events${path}`, {
    ...init,
    headers: { "Content-Type": "application/json", ...init?.headers },
    credentials: "same-origin",
    cache: "no-store",
    signal: AbortSignal.timeout(15000),
  });
  let data: { error?: string } & Record<string, unknown>;
  try {
    data = await response.json();
  } catch {
    throw new Error("Nie można połączyć się z serwerem. Spróbuj ponownie później.");
  }
  if (!response.ok) throw new Error(data.error ?? "Nie udało się wykonać operacji.");
  return data as T;
}

function localDateTime(value: string) {
  const date = new Date(value);
  return new Date(date.getTime() - date.getTimezoneOffset() * 60000).toISOString().slice(0, 16);
}

function formatDate(value: string) {
  return new Intl.DateTimeFormat("pl-PL", { dateStyle: "medium", timeStyle: "short" }).format(new Date(value));
}

function EventForm({
  event,
  pending,
  onSave,
  onCancel,
}: {
  event?: Event;
  pending: boolean;
  onSave: (input: EventInput) => Promise<void>;
  onCancel: () => void;
}) {
  const [error, setError] = useState("");

  async function submit(formEvent: FormEvent<HTMLFormElement>) {
    formEvent.preventDefault();
    setError("");
    const form = formEvent.currentTarget;
    const fields = new FormData(form);
    try {
      await onSave({
        title: String(fields.get("title")).trim(),
        description: String(fields.get("description")).trim(),
        starts_at: new Date(String(fields.get("starts_at"))).toISOString(),
        ends_at: new Date(String(fields.get("ends_at"))).toISOString(),
        location: String(fields.get("location")).trim(),
        is_published: event?.is_published ?? fields.get("is_published") === "on",
      });
      form.reset();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Nie udało się zapisać wydarzenia.");
    }
  }

  return <form className="card border shadow-sm mb-4" onSubmit={submit} aria-busy={pending}>
    <div className="card-body p-4">
      <h2 className="h4 mb-4">{event ? "Edytuj wydarzenie" : "Nowe wydarzenie"}</h2>
      {error && <div className="alert alert-danger" role="alert">{error}</div>}
      <div className="row g-3">
        <div className="col-md-8">
          <label className="form-label" htmlFor="event-title">Tytuł</label>
          <input className="form-control" id="event-title" name="title" defaultValue={event?.title} maxLength={200} required disabled={pending} />
        </div>
        <div className="col-md-4">
          <label className="form-label" htmlFor="event-location">Miejsce</label>
          <input className="form-control" id="event-location" name="location" defaultValue={event?.location} maxLength={255} required disabled={pending} />
        </div>
        <div className="col-md-6">
          <label className="form-label" htmlFor="event-starts">Rozpoczęcie</label>
          <input className="form-control" type="datetime-local" id="event-starts" name="starts_at" defaultValue={event ? localDateTime(event.starts_at) : undefined} required disabled={pending} />
        </div>
        <div className="col-md-6">
          <label className="form-label" htmlFor="event-ends">Zakończenie</label>
          <input className="form-control" type="datetime-local" id="event-ends" name="ends_at" defaultValue={event ? localDateTime(event.ends_at) : undefined} required disabled={pending} />
        </div>
        <div className="col-12">
          <label className="form-label" htmlFor="event-description">Opis</label>
          <textarea className="form-control" id="event-description" name="description" rows={4} maxLength={65535} defaultValue={event?.description} required disabled={pending} />
        </div>
        {!event && <div className="col-12">
          <div className="form-check">
            <input className="form-check-input" type="checkbox" id="event-published" name="is_published" disabled={pending} />
            <label className="form-check-label" htmlFor="event-published">Opublikuj od razu</label>
          </div>
        </div>}
      </div>
      <div className="d-flex gap-2 mt-4">
        <button type="submit" className="btn btn-primary" disabled={pending}>{pending ? "Zapisywanie…" : "Zapisz wydarzenie"}</button>
        <button type="button" className="btn btn-outline-secondary" onClick={onCancel} disabled={pending}>Anuluj</button>
      </div>
    </div>
  </form>;
}

export function AdminEvents() {
  const { user, loading } = useAuth();
  const [events, setEvents] = useState<Event[]>([]);
  const [pageError, setPageError] = useState("");
  const [loadingEvents, setLoadingEvents] = useState(true);
  const [pendingId, setPendingId] = useState<number | null>(null);
  const [creating, setCreating] = useState(false);
  const [editingId, setEditingId] = useState<number | null>(null);

  const loadEvents = useCallback(async () => {
    setLoadingEvents(true);
    setPageError("");
    try {
      const result = await adminRequest<{ events: Event[] }>("", { method: "GET" });
      setEvents(result.events);
    } catch (err) {
      setPageError(err instanceof Error ? err.message : "Nie udało się pobrać wydarzeń.");
    } finally {
      setLoadingEvents(false);
    }
  }, []);

  useEffect(() => {
    if (!loading && user?.role === "admin") void loadEvents();
  }, [loading, user, loadEvents]);

  async function createEvent(input: EventInput) {
    setPendingId(0);
    try {
      await adminRequest<{ event: Event }>("", { method: "POST", body: JSON.stringify(input) });
      setCreating(false);
      await loadEvents();
    } finally {
      setPendingId(null);
    }
  }

  async function updateEvent(id: number, input: EventInput) {
    setPendingId(id);
    try {
      await adminRequest<{ event: Event }>(`/${id}`, { method: "PUT", body: JSON.stringify(input) });
      setEditingId(null);
      await loadEvents();
    } finally {
      setPendingId(null);
    }
  }

  async function setPublication(event: Event) {
    setPendingId(event.id);
    setPageError("");
    try {
      await adminRequest(`/${event.id}/publication`, {
        method: "PATCH",
        body: JSON.stringify({ is_published: !event.is_published }),
      });
      await loadEvents();
    } catch (err) {
      setPageError(err instanceof Error ? err.message : "Nie udało się zmienić statusu publikacji.");
    } finally {
      setPendingId(null);
    }
  }

  async function deleteEvent(event: Event) {
    if (!window.confirm(`Czy na pewno chcesz usunąć wydarzenie „${event.title}”? Tej operacji nie można cofnąć.`)) return;
    setPendingId(event.id);
    setPageError("");
    try {
      await adminRequest(`/${event.id}`, { method: "DELETE" });
      await loadEvents();
    } catch (err) {
      setPageError(err instanceof Error ? err.message : "Nie udało się usunąć wydarzenia.");
    } finally {
      setPendingId(null);
    }
  }

  if (loading) return <div className="container py-5" role="status">Sprawdzanie uprawnień…</div>;
  if (!user) return <section className="container py-5">
    <h1 className="h2">Zarządzanie wydarzeniami</h1>
    <p>Aby zarządzać wydarzeniami, zaloguj się na konto administratora.</p>
    <Link href="/login" className="btn btn-primary">Zaloguj się</Link>
  </section>;
  if (user.role !== "admin") return <section className="container py-5">
    <h1 className="h2">Brak uprawnień</h1>
    <p>Ta sekcja jest dostępna wyłącznie dla administratorów.</p>
    <Link href="/" className="btn btn-outline-primary">Wróć na stronę główną</Link>
  </section>;

  return <section className="container py-5">
    <nav aria-label="Okruszki"><ol className="breadcrumb small">
      <li className="breadcrumb-item"><Link href="/">Strona główna</Link></li>
      <li className="breadcrumb-item active" aria-current="page">Zarządzanie wydarzeniami</li>
    </ol></nav>
    <div className="d-flex flex-column flex-md-row justify-content-between align-items-md-center gap-3 mb-4">
      <div><p className="text-uppercase small text-secondary fw-semibold letter-spacing mb-1">Panel administratora</p><h1 className="h2 mb-0">Wydarzenia</h1></div>
      <button className="btn btn-primary" type="button" onClick={() => { setCreating(!creating); setEditingId(null); }}>
        {creating ? "Anuluj dodawanie" : "Dodaj wydarzenie"}
      </button>
    </div>
    {pageError && <div className="alert alert-danger" role="alert">{pageError}</div>}
    {creating && <EventForm pending={pendingId === 0} onSave={createEvent} onCancel={() => setCreating(false)} />}
    {loadingEvents ? <p role="status">Ładowanie wydarzeń…</p> : events.length === 0 ? <div className="placeholder-card p-4">
      <p className="mb-0">Brak wydarzeń. Dodaj pierwsze wydarzenie, aby pojawiło się na liście.</p>
    </div> : <div className="d-flex flex-column gap-3">
      {events.map(event => editingId === event.id ? <EventForm key={`edit-${event.id}`} event={event} pending={pendingId === event.id}
        onSave={input => updateEvent(event.id, input)} onCancel={() => setEditingId(null)} /> :
        <article className="card border shadow-sm" key={event.id}>
          <div className="card-body p-4">
            <div className="d-flex flex-column flex-lg-row justify-content-between gap-3">
              <div className="flex-grow-1">
                <div className="d-flex flex-wrap align-items-center gap-2 mb-2">
                  <h2 className="h5 mb-0">{event.title}</h2>
                  <span className={`badge ${event.is_published ? "text-bg-success" : "text-bg-secondary"}`}>
                    {event.is_published ? "Opublikowane" : "Ukryte"}
                  </span>
                </div>
                <p className="text-secondary mb-2">{event.description}</p>
                <p className="small mb-0"><strong>{event.location}</strong> · {formatDate(event.starts_at)}–{formatDate(event.ends_at)}</p>
              </div>
              <div className="d-flex flex-wrap align-items-start gap-2">
                <button className="btn btn-sm btn-outline-primary" type="button" disabled={pendingId !== null} onClick={() => { setEditingId(event.id); setCreating(false); }}>Edytuj</button>
                <button className="btn btn-sm btn-outline-primary" type="button" disabled={pendingId !== null} onClick={() => void setPublication(event)}>
                  {pendingId === event.id ? "Zapisywanie…" : event.is_published ? "Ukryj" : "Opublikuj"}
                </button>
                <button className="btn btn-sm btn-outline-danger" type="button" disabled={pendingId !== null} onClick={() => void deleteEvent(event)}>Usuń</button>
              </div>
            </div>
          </div>
        </article>)}
    </div>}
  </section>;
}
