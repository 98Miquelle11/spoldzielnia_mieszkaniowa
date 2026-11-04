export type CulturalEvent = {
  id: number;
  title: string;
  description: string;
  location: string;
  starts_at: string;
  ends_at: string;
  is_published: boolean;
  created_at: string;
  updated_at: string;
};
export type Period = "upcoming" | "past";

export function isCulturalEvent(value: unknown): value is CulturalEvent {
  if (!value || typeof value !== "object") return false;
  const e = value as Record<string, unknown>;
  return typeof e.id === "number" && Number.isSafeInteger(e.id) && e.id > 0 &&
    ["title", "description", "location"].every(k => typeof e[k] === "string") &&
    ["starts_at", "ends_at", "created_at", "updated_at"].every(k =>
      typeof e[k] === "string" && Number.isFinite(Date.parse(e[k] as string))) &&
    typeof e.is_published === "boolean";
}

export function formatEventDate(value: string): string {
  return new Intl.DateTimeFormat("pl-PL", {
    dateStyle: "long", timeStyle: "short", timeZone: "Europe/Warsaw",
  }).format(new Date(value));
}

export class EventsError extends Error {
  constructor(message: string, public readonly status: number) { super(message); }
}

export async function fetchEvents(path: string, signal: AbortSignal): Promise<unknown> {
  let response: Response;
  try {
    response = await fetch(`/api/events${path}`, {
      cache: "no-store", signal: AbortSignal.any([signal, AbortSignal.timeout(15000)]),
    });
  } catch (error) {
    if (signal.aborted) throw error;
    throw new EventsError("Nie można połączyć się z serwerem. Spróbuj ponownie.", 0);
  }
  if (!response.ok) {
    throw new EventsError(response.status === 404
      ? "Wydarzenie jest niedostępne lub nie zostało opublikowane."
      : response.status === 400 ? "Nieprawidłowy adres wydarzenia lub filtr."
      : "Nie udało się pobrać wydarzeń. Spróbuj ponownie później.", response.status);
  }
  try { return await response.json(); }
  catch { throw new EventsError("Serwer zwrócił nieprawidłowe dane.", 0); }
}

export type EventPage = {
  items: CulturalEvent[];
  pagination: { page: number; limit: number; has_more: boolean };
};
export function isEventPage(value: unknown): value is EventPage {
  if (!value || typeof value !== "object") return false;
  const v = value as Record<string, unknown>;
  if (!Array.isArray(v.items) || !v.items.every(isCulturalEvent) ||
      !v.pagination || typeof v.pagination !== "object") return false;
  const p = v.pagination as Record<string, unknown>;
  return typeof p.page === "number" && Number.isSafeInteger(p.page) && p.page >= 1 && p.page <= 1000000 &&
    typeof p.limit === "number" && Number.isInteger(p.limit) && p.limit >= 1 && p.limit <= 50 &&
    typeof p.has_more === "boolean" && v.items.length <= p.limit;
}
