"use client";

import Link from "next/link";
import {
  useCallback,
  useEffect,
  useState,
  type FormEvent,
} from "react";
import { useAuth } from "./auth";

type CulturalEvent = {
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

type EventInput = Pick<
  CulturalEvent,
  | "title"
  | "description"
  | "location"
  | "starts_at"
  | "ends_at"
  | "is_published"
>;

type Draft = {
  id: number | null;
  title: string;
  description: string;
  location: string;
  starts_at: string;  //daty w formacie wymaganym przez input datetime-local.
  ends_at: string;
  is_published: boolean;
};

type FieldErrors = Record<string, string>;
const DRAFT_KEY = "admin-events-draft-v1";  //docelowo można użyć klucza zawierającego id użytkownika

class ApiError extends Error {
  constructor(
    message: string,
    public readonly status: number,
    public readonly fields: FieldErrors = {},  //kod odpowiedzi HTTP + błędy poszczególnych pól
  ) {
    super(message);
    this.name = "ApiError";
  }
}

async function adminRequest<T = void>(
  path: string,
  init: RequestInit = {},
): Promise<T> {
  const headers = new Headers(init.headers);
  headers.set("Content-Type", "application/json");

  let response: Response;

  try {
    response = await fetch(`/api/admin/events${path}`, {
      ...init,
      headers,
      credentials: "same-origin",
      cache: "no-store",
      signal: AbortSignal.timeout(15000),
    });
  } catch {
    throw new ApiError(
      "Nie można połączyć się z serwerem. Spróbuj ponownie.",
      0,
    );
  }

  if (response.status === 204) {    //204 response
    return undefined as T;
  }

  let data: unknown;

  try {
    data = await response.json();
  } catch {
    throw new ApiError(
      response.ok
        ? "Serwer zwrócił nieprawidłową odpowiedź."
        : `Nie udało się wykonać operacji (HTTP ${response.status}).`,
      response.status,
    );
  }

  if (!response.ok) {
    const payload =
      typeof data === "object" && data !== null
        ? (data as Record<string, unknown>)
        : {};

    const fields: FieldErrors = {};

    if (
      typeof payload.fields === "object" &&
      payload.fields !== null
    ) {
      for (const [name, message] of Object.entries(payload.fields)) {
        if (typeof message === "string") {
          fields[name] = message;
        }
      }
    }

    const message =
      response.status === 401
        ? "Sesja wygasła. Zaloguj się ponownie. Szkic formularza został zachowany."
        : typeof payload.error === "string"
          ? payload.error
          : response.status === 403
            ? "Serwer odrzucił dostęp do tej operacji."
            : "Nie udało się wykonać operacji.";

    throw new ApiError(message, response.status, fields);
  }

  return data as T;
}

function emptyDraft(): Draft {
  return {
    id: null,
    title: "",
    description: "",
    location: "",
    starts_at: "",
    ends_at: "",
    is_published: false,
  };
}

function toLocalInput(value: string): string {
  const date = new Date(value);

  if (Number.isNaN(date.getTime())) {
    throw new Error("Wydarzenie zawiera nieprawidłową datę.");
  }

  const pad = (number: number) => String(number).padStart(2, "0");

  return (
    `${String(date.getFullYear()).padStart(4, "0")}-` +
    `${pad(date.getMonth() + 1)}-${pad(date.getDate())}` +
    `T${pad(date.getHours())}:${pad(date.getMinutes())}:` +
    `${pad(date.getSeconds())}`
  );
}

function eventToDraft(event: CulturalEvent): Draft {
  return {
    id: event.id,
    title: event.title,
    description: event.description,
    location: event.location,
    starts_at: toLocalInput(event.starts_at),
    ends_at: toLocalInput(event.ends_at),
    is_published: event.is_published,
  };
}

//bez id, created_at i updated_at
function eventToInput(event: CulturalEvent): EventInput {
  return {
    title: event.title,
    description: event.description,
    location: event.location,
    starts_at: event.starts_at,
    ends_at: event.ends_at,
    is_published: event.is_published,
  };
}

function formatDate(value: string): string {
  const date = new Date(value);

  if (Number.isNaN(date.getTime())) {
    return "Nieprawidłowa data";
  }

  return new Intl.DateTimeFormat("pl-PL", {
    dateStyle: "medium",
    timeStyle: "short",
  }).format(date);
}

// sprawdzenie danych z sessionstorage, bo mogą pochodzić ze starszej wersji formularza lub zostać zmodyfikowane
function isDraft(value: unknown): value is Draft {
  if (typeof value !== "object" || value === null) return false;

  const draft = value as Record<string, unknown>;

  return (
    (draft.id === null ||
      (typeof draft.id === "number" &&
        Number.isSafeInteger(draft.id) &&
        draft.id > 0)) &&
    typeof draft.title === "string" &&
    typeof draft.description === "string" &&
    typeof draft.location === "string" &&
    typeof draft.starts_at === "string" &&
    typeof draft.ends_at === "string" &&
    typeof draft.is_published === "boolean"
  );
}

function validateDraft(draft: Draft): FieldErrors {
  const errors: FieldErrors = {};

  const textFields = [
    ["title", "Tytuł", 200],
    ["location", "Miejsce", 200],
    ["description", "Opis", 5000],
  ] as const;

  for (const [field, label, limit] of textFields) {
    const value = draft[field].trim();

    if (!value) {
      errors[field] = `${label} jest wymagany.`;
    } else if (Array.from(value).length > limit) {
      errors[field] = `${label}: maksymalnie ${limit} znaków.`;
    }
  }

  const startsAt = draft.starts_at
    ? new Date(draft.starts_at)
    : null;

  const endsAt = draft.ends_at
    ? new Date(draft.ends_at)
    : null;

  if (!startsAt || Number.isNaN(startsAt.getTime())) {
    errors.starts_at = "Podaj poprawny termin rozpoczęcia.";
  }

  if (!endsAt || Number.isNaN(endsAt.getTime())) {
    errors.ends_at = "Podaj poprawny termin zakończenia.";
  }

  // błąd kolejności terminów pokazujemy przed wysłaniem żądania równe rozpoczęcie i zakończenie dozwolone
  if (
    startsAt &&
    endsAt &&
    !errors.starts_at &&
    !errors.ends_at &&
    endsAt.getTime() < startsAt.getTime()
  ) {
    errors.ends_at =
      "Zakończenie nie może być wcześniejsze niż rozpoczęcie.";
  }

  return errors;
}

export function AdminEvents() {
  const { user, loading: authLoading, setUser } = useAuth();

  const [events, setEvents] = useState<CulturalEvent[]>([]);

  const [listStatus, setListStatus] = useState<
    "loading" | "success" | "error"
  >("loading");

  const [draft, setDraft] = useState<Draft | null>(null);

  // chroni przed nadpisaniem zapisanego szkicu zanim go odczytamy
  const [draftReady, setDraftReady] = useState(false);
  const [draftWarning, setDraftWarning] = useState("");

  const [fieldErrors, setFieldErrors] = useState<FieldErrors>({});
  const [formError, setFormError] = useState("");
  const [pageError, setPageError] = useState("");
  const [notice, setNotice] = useState("");
  const [needsLogin, setNeedsLogin] = useState(false);

  // blokowanie przycisków w trakcie operacji
  const [busy, setBusy] = useState(false);
  const [timeZone, setTimeZone] = useState("");

  useEffect(() => {
    setTimeZone(Intl.DateTimeFormat().resolvedOptions().timeZone);

    try {
      const saved = sessionStorage.getItem(DRAFT_KEY);

      if (saved) {
        const parsed: unknown = JSON.parse(saved);

        if (isDraft(parsed)) {
          setDraft(parsed);
        } else {
          sessionStorage.removeItem(DRAFT_KEY);
        }
      }
    } catch {
      setDraftWarning(
        "Nie można odczytać zapisanego szkicu formularza.",
      );
    }

    setDraftReady(true);
  }, []);

  useEffect(() => {
    if (!draftReady) return;

    // po zapisie lub anulowaniu draft staje się null
    try {
      if (draft) {
        sessionStorage.setItem(DRAFT_KEY, JSON.stringify(draft));
      } else {
        sessionStorage.removeItem(DRAFT_KEY);
      }

      setDraftWarning("");
    } catch {
      setDraftWarning(
        "Przeglądarka nie pozwala zapisać szkicu. " +
          "Przed opuszczeniem strony skopiuj treść formularza.",
      );
    }
  }, [draft, draftReady]);

  const loadEvents = useCallback(async () => {
    setListStatus("loading");
    setPageError("");

    try {
      const result = await adminRequest<{
        items: CulturalEvent[];
      }>("", { method: "GET" });

      if (!Array.isArray(result.items)) {
        throw new Error("Serwer zwrócił nieprawidłową listę wydarzeń.");
      }

      setEvents(result.items);
      setListStatus("success");
      setNeedsLogin(false);
    } catch (error) {
      setListStatus("error");
      setPageError(
        error instanceof Error
          ? error.message
          : "Nie udało się pobrać wydarzeń.",
      );

      if (error instanceof ApiError && error.status === 401) {
        setNeedsLogin(true);
        setUser(null);
      }
    }
  }, [setUser]);

  useEffect(() => {
    if (!authLoading && user?.role === "admin") {
      void loadEvents();
    }
  }, [authLoading, user?.id, user?.role, loadEvents]);

  function showError(error: unknown, forForm = false) {
    const message =
      error instanceof Error
        ? error.message
        : "Nie udało się wykonać operacji.";

    if (forForm) {
      setFormError(message);
      setFieldErrors(error instanceof ApiError ? error.fields : {});
    } else {
      setPageError(message);
    }

    if (error instanceof ApiError && error.status === 401) {
      setNeedsLogin(true);
        setUser(null);
    }
  }

  function beginCreate() {
    if (
      draft &&
      !window.confirm("Odrzucić aktualny szkic formularza?")
    ) {
      return;
    }

    setDraft(emptyDraft());
    setFieldErrors({});
    setFormError("");
    setNotice("");
  }

  function cancelDraft() {
    if (!window.confirm("Odrzucić szkic formularza?")) return;

    setDraft(null);
    setFieldErrors({});
    setFormError("");
  }

  async function beginEdit(id: number) {
    if (
      draft &&
      !window.confirm("Odrzucić aktualny szkic formularza?")
    ) {
      return;
    }

    setBusy(true);
    setPageError("");  //***
    setNotice("");

    try {
      const event = await adminRequest<CulturalEvent>(`/${id}`, {
        method: "GET",
      });
      setDraft(eventToDraft(event));
      setFieldErrors({});
      setFormError("");
    } catch (error) {
      if (error instanceof ApiError && error.status === 404) {
        await loadEvents();
      }

      showError(error);
    } finally {
      setBusy(false);
    }
  }

  async function saveEvent(formEvent: FormEvent<HTMLFormElement>) {
    formEvent.preventDefault();
    if (!draft || busy) return;

    setFormError("");
    setNotice("");

    const errors = validateDraft(draft);
    setFieldErrors(errors);

    if (Object.keys(errors).length > 0) {
      setFormError("Popraw zaznaczone pola.");
      return;
    }

    const input: EventInput = {
      title: draft.title.trim(),
      description: draft.description.trim(),
      location: draft.location.trim(),
      starts_at: new Date(draft.starts_at).toISOString(),
      ends_at: new Date(draft.ends_at).toISOString(),
      is_published: draft.is_published,
    };

    setBusy(true);

    try {
      if (draft.id === null) {
        await adminRequest<{ id: number }>("", {
          method: "POST",
          body: JSON.stringify(input),
        });
      } else {
        await adminRequest<void>(`/${draft.id}`, {
          method: "PUT",
          body: JSON.stringify(input),
        });
      }

      setDraft(null);
      setFieldErrors({});
      setNeedsLogin(false);
      setNotice("Wydarzenie zostało zapisane.");

      await loadEvents();
    } catch (error) {
      showError(error, true);
    } finally {
      setBusy(false);
    }
  }

  async function togglePublication(event: CulturalEvent) {
    setBusy(true);
    setPageError("");
    setNotice("");

    try {
      const current = await adminRequest<CulturalEvent>(
        `/${event.id}`,
        { method: "GET" },
      );

      const published = !current.is_published;

      await adminRequest<void>(`/${current.id}`, {
        method: "PUT",
        body: JSON.stringify({
          ...eventToInput(current),
          is_published: published,
        }),
      });

      setNotice(
        published
          ? "Wydarzenie zostało opublikowane."
          : "Wydarzenie zostało ukryte.",
      );

      await loadEvents();
    } catch (error) {
      if (error instanceof ApiError && error.status === 404) {
        await loadEvents();
      }

      showError(error);
    } finally {
      setBusy(false);
    }
  }

  async function deleteEvent(event: CulturalEvent) {
    if (
      !window.confirm(
        `Czy na pewno usunąć wydarzenie „${event.title}”?`,
      )
    ) {
      return;
    }

    setBusy(true);
    setPageError("");
    setNotice("");

    try {
      await adminRequest<void>(`/${event.id}`, {
        method: "DELETE",
        body: JSON.stringify({}),
      });

      setNotice("Wydarzenie zostało usunięte.");
      await loadEvents();
    } catch (error) {
      if (error instanceof ApiError && error.status === 404) {
        await loadEvents();
      }

      showError(error);
    } finally {
      setBusy(false);
    }
  }

  function updateDraft(field: keyof EventInput, value: string | boolean) {
    setDraft((current) =>
      current ? { ...current, [field]: value } : current,
    );
    setFieldErrors((current) => {
      const next = { ...current };
      delete next[field];
      return next;
    });
  }

  function fieldError(name: string) {
    return fieldErrors[name] ? (
      <div id={`error-${name}`} className="invalid-feedback">
        {fieldErrors[name]}
      </div>
    ) : null;
  }

  function fieldProps(name: string) {
    return {
      className: `form-control${fieldErrors[name] ? " is-invalid" : ""}`,
      "aria-invalid": Boolean(fieldErrors[name]),
      "aria-describedby": fieldErrors[name]
        ? `error-${name}`
        : undefined,
    };
  }

  if (authLoading || !draftReady) {
    return (
      <div className="container py-5" role="status">
        Sprawdzanie uprawnień…
      </div>
    );
  }

  //prezentacja
  if (!user) {
    return (
      <section className="container py-5">
        <h1 className="h2">Zarządzanie wydarzeniami</h1>
        <p>Zaloguj się na konto administratora.</p>

        {draft && (
          <p>Twój szkic formularza został zachowany na tej karcie.</p>
        )}

        {draftWarning && (
          <p className="text-danger" role="alert">
            {draftWarning}
          </p>
        )}

        <Link href="/login" className="btn btn-primary">
          Zaloguj się
        </Link>
      </section>
    );
  }

  if (user.role !== "admin") {
    return (
      <section className="container py-5">
        <h1 className="h2">Brak uprawnień</h1>
        <p>Ta sekcja jest dostępna wyłącznie dla administratorów.</p>

        <Link href="/" className="btn btn-outline-primary">
          Wróć na stronę główną
        </Link>
      </section>
    );
  }

  return (
    <section className="container py-5">
      <div className="d-flex flex-wrap justify-content-between align-items-center gap-3 mb-4">
        <h1 className="h2 mb-0">Zarządzanie wydarzeniami</h1>

        <button
          type="button"
          className="btn btn-primary"
          disabled={busy}
          onClick={beginCreate}
        >
          Dodaj wydarzenie
        </button>
      </div>

      {notice && (
        <div className="alert alert-success" role="status">
          {notice}
        </div>
      )}

      {pageError && (
        <div className="alert alert-danger" role="alert">
          {pageError}
        </div>
      )}

      {/* !!wygaśnięcie sesji ma osobny komunikat i link logowania!! */}
      {needsLogin && (
        <div className="alert alert-warning" role="alert">
          <p>
            Zaloguj się ponownie. Po zalogowaniu wróć do zarządzania
            wydarzeniami na tej samej karcie.
          </p>

          <Link href="/login" className="btn btn-outline-primary">
            Przejdź do logowania
          </Link>
        </div>
      )}

      {draftWarning && (
        <div className="alert alert-warning" role="alert">
          {draftWarning}
        </div>
      )}

      {draft && (
        <form
          className="card shadow-sm mb-4"
          onSubmit={saveEvent}
          noValidate
          aria-busy={busy}
        >
          {/* noValidate wyłącza domyślne komunikaty przeglądarki
              Wymagane pola i limity sprawdza validateDraft */}
          <div className="card-body p-4">
            <h2 className="h4">
              {draft.id === null ? "Nowe wydarzenie" : "Edytuj wydarzenie"}
            </h2>

            {formError && (
              <div className="alert alert-danger" role="alert">
                {formError}
              </div>
            )}

            {/* Blokujemy cały formularz podczas operacji */}
            <fieldset disabled={busy}>
              <div className="row g-3">
                <div className="col-md-8">
                  <label htmlFor="event-title" className="form-label">
                    Tytuł
                  </label>

                  <input
                    {...fieldProps("title")}
                    id="event-title"
                    name="title"
                    value={draft.title}
                    required
                    onChange={(e) => updateDraft("title", e.target.value)}
                  />

                  {/* Limity sprawdzamy przez Array.from zamiast maxLength */}
                  <div className="form-text">Maksymalnie 200 znaków.</div>
                  {fieldError("title")}
                </div>

                <div className="col-md-4">
                  <label htmlFor="event-location" className="form-label">
                    Miejsce
                  </label>

                  <input
                    {...fieldProps("location")}
                    id="event-location"
                    name="location"
                    value={draft.location}
                    required
                    onChange={(e) => updateDraft("location", e.target.value)}
                  />

                  <div className="form-text">Maksymalnie 200 znaków.</div>
                  {fieldError("location")}
                </div>

                <div className="col-12">
                  <p className="small text-secondary mb-0">
                    Daty w strefie przeglądarki: {timeZone}.
                  </p>
                </div>

                <div className="col-md-6">
                  <label htmlFor="event-starts" className="form-label">
                    Rozpoczęcie
                  </label>

                  <input
                    {...fieldProps("starts_at")}
                    id="event-starts"
                    name="starts_at"
                    type="datetime-local"
                    step={1}
                    value={draft.starts_at}
                    required
                    onChange={(e) => updateDraft("starts_at", e.target.value)}
                  />

                  {fieldError("starts_at")}
                </div>

                <div className="col-md-6">
                  <label htmlFor="event-ends" className="form-label">
                    Zakończenie
                  </label>

                  <input
                    {...fieldProps("ends_at")}
                    id="event-ends"
                    name="ends_at"
                    type="datetime-local"
                    step={1}
                    value={draft.ends_at}
                    required
                    onChange={(e) => updateDraft("ends_at", e.target.value)}
                  />

                  {fieldError("ends_at")}
                </div>

                <div className="col-12">
                  <label htmlFor="event-description" className="form-label">
                    Opis
                  </label>

                  <textarea
                    {...fieldProps("description")}
                    id="event-description"
                    name="description"
                    rows={5}
                    value={draft.description}
                    required
                    onChange={(e) =>
                      updateDraft("description", e.target.value)
                    }
                  />

                  <div className="form-text">Maksymalnie 5000 znaków.</div>
                  {fieldError("description")}
                </div>

                <div className="col-12">
                  <div className="form-check">
                    <input
                      className="form-check-input"
                      id="event-published"
                      name="is_published"
                      type="checkbox"
                      checked={draft.is_published}
                      onChange={(e) =>
                        updateDraft("is_published", e.target.checked)
                      }
                    />

                    <label
                      htmlFor="event-published"
                      className="form-check-label"
                    >
                      Wydarzenie opublikowane
                    </label>
                  </div>

                  {fieldErrors.is_published && (
                    <p className="text-danger small">
                      {fieldErrors.is_published}
                    </p>
                  )}
                </div>
              </div>

              <div className="d-flex gap-2 mt-4">
                <button type="submit" className="btn btn-primary">
                  {busy ? "Zapisywanie…" : "Zapisz wydarzenie"}
                </button>

                <button
                  type="button"
                  className="btn btn-outline-secondary"
                  onClick={cancelDraft}
                >
                  Anuluj
                </button>
              </div>
            </fieldset>
          </div>
        </form>
      )}

      {listStatus === "loading" && (
        <p role="status">Ładowanie wydarzeń…</p>
      )}
      {listStatus === "error" && (
        <button
          type="button"
          className="btn btn-outline-primary"
          disabled={busy}
          onClick={() => void loadEvents()}
        >
          Spróbuj ponownie pobrać listę
        </button>
      )}

      {listStatus === "success" && events.length === 0 && (
        <p>Brak wydarzeń. Dodaj pierwsze wydarzenie.</p>
      )}

      {listStatus === "success" && events.length > 0 && (
        <div className="d-flex flex-column gap-3">
          {events.map((event) => (
            <article className="card shadow-sm" key={event.id}>
              <div className="card-body p-4">
                <div className="d-flex flex-column flex-lg-row justify-content-between gap-3">
                  <div className="flex-grow-1">
                    <div className="d-flex flex-wrap align-items-center gap-2 mb-2">
                      <h2 className="h5 mb-0">{event.title}</h2>

                      <span
                        className={`badge ${
                          event.is_published
                            ? "text-bg-success"
                            : "text-bg-secondary"
                        }`}
                      >
                        {event.is_published ? "Opublikowane" : "Ukryte"}
                      </span>
                    </div>

                    <p
                      style={{
                        whiteSpace: "pre-wrap",
                        overflowWrap: "anywhere",
                      }}
                    >
                      {event.description}
                    </p>

                    <p className="small mb-0">
                      <strong>{event.location}</strong>
                      {" · "}
                      {formatDate(event.starts_at)}
                      {" – "}
                      {formatDate(event.ends_at)}
                      {timeZone && ` (${timeZone})`}
                    </p>
                  </div>

                  <div className="d-flex flex-wrap align-items-start gap-2">
                    <button
                      type="button"
                      className="btn btn-sm btn-outline-primary"
                      disabled={busy}
                      onClick={() => void beginEdit(event.id)}
                    >
                      Edytuj
                    </button>

                    <button
                      type="button"
                      className="btn btn-sm btn-outline-primary"
                      disabled={busy || draft !== null}
                      onClick={() => void togglePublication(event)}
                    >
                      {event.is_published ? "Ukryj" : "Opublikuj"}
                    </button>

                    <button
                      type="button"
                      className="btn btn-sm btn-outline-danger"
                      disabled={busy || draft !== null}
                      onClick={() => void deleteEvent(event)}
                    >
                      Usuń
                    </button>
                  </div>
                </div>
              </div>
            </article>
          ))}
        </div>
      )}
    </section>
  );
}
