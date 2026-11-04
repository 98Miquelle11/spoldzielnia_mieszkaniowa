import { test, expect, type Page } from "@playwright/test";

const adminEmail = process.env.E2E_ADMIN_EMAIL;
const adminPassword = process.env.E2E_ADMIN_PASSWORD;
const residentEmail = process.env.E2E_RESIDENT_EMAIL;
const residentPassword = process.env.E2E_RESIDENT_PASSWORD;
const fixture = {
  id: 987, title: "Koncert testowy", description: "Opis\nDrugi akapit", location: "Świetlica",
  starts_at: "2026-11-21T16:00:00Z", ends_at: "2026-11-21T18:00:00Z",
  is_published: true, created_at: "2026-11-01T10:00:00Z", updated_at: "2026-11-01T10:00:00Z",
};
async function login(page: Page, baseURL: string, email: string, password: string) {
  const response = await page.request.post("/api/auth/login", { headers: { Origin: new URL(baseURL).origin }, data: { email, password } });
  expect(response.status()).toBe(200);
}

test("publiczna lista, szczegóły i czas Warszawy niezależny od strefy przeglądarki", async ({ browser, baseURL }) => {
  const context = await browser.newContext({ baseURL, timezoneId: "America/New_York" });
  const page = await context.newPage();
  await page.route("**/api/events?*", route => route.fulfill({ json: { items: [fixture] } }));
  await page.route("**/api/events/987", route => route.fulfill({ json: fixture }));
  await page.goto("/imprezy");
  await expect(page.getByRole("heading", { name: fixture.title })).toBeVisible();
  await expect(page.locator("main")).toContainText("17:00");
  await page.getByRole("link", { name: /^Szczegóły/ }).click();
  await expect(page.getByRole("heading", { name: fixture.title })).toBeVisible();
  await expect(page.locator("main")).toContainText("Drugi akapit");
  await context.close();
});

test("archiwum pozostaje po odświeżeniu, pusta lista i błędny filtr", async ({ page }) => {
  await page.route("**/api/events?*", route => route.fulfill({ json: { items: [] } }));
  await page.goto("/imprezy");
  await page.getByRole("link", { name: "Archiwalne", exact: true }).click();
  await expect(page).toHaveURL(/period=past/);
  await page.reload();
  await expect(page.getByText("Brak archiwalnych wydarzeń.")).toBeVisible();
  await expect(page.getByRole("link", { name: "Archiwalne", exact: true })).toHaveAttribute("aria-current", "page");
  await page.goto("/imprezy?period=wrong");
  await expect(page.locator("main").getByRole("alert")).toContainText("Nieznany filtr");
});

test("ładowanie, błąd połączenia, ponowienie i niedostępne szczegóły", async ({ page }) => {
  let fail = true;
  await page.route("**/api/events?*", async route => {
    await new Promise(resolve => setTimeout(resolve, 400));
    if (fail) await route.abort(); else await route.fulfill({ json: { items: [] } });
  });
  await page.goto("/imprezy");
  await expect(page.getByText("Ładowanie wydarzeń…")).toBeVisible();
  await expect(page.locator("main").getByRole("alert")).toContainText("Nie można połączyć");
  fail = false;
  await page.getByRole("button", { name: "Spróbuj ponownie", exact: true }).click();
  await expect(page.getByText("Brak nadchodzących wydarzeń.")).toBeVisible();
  await page.route("**/api/events/987", route => route.fulfill({ status: 404, json: { error: "Brak" } }));
  await page.goto("/imprezy/987?period=past");
  await expect(page.getByRole("heading", { name: "Wydarzenie niedostępne" })).toBeVisible();
  await page.getByRole("link", { name: "Wróć do archiwum" }).click();
  await expect(page).toHaveURL(/period=past/);
});

test("duża lista i układ mobilny", async ({ page }) => {
  await page.setViewportSize({ width: 375, height: 812 });
  await page.route("**/api/events?*", route => route.fulfill({ json: { items: Array.from({ length: 45 }, (_, i) => ({ ...fixture, id: i + 1, title: `Wydarzenie ${i + 1} ${"x".repeat(200)}` })) } }));
  await page.goto("/imprezy");
  await expect(page.locator("main article")).toHaveCount(20);
  await page.getByRole("button", { name: /Pokaż kolejne/ }).click();
  await expect(page.locator("main article")).toHaveCount(40);
  expect(await page.evaluate(() => document.documentElement.scrollWidth <= window.innerWidth)).toBe(true);
});

test("API blokuje gościa na wszystkich trasach administracyjnych", async ({ request, baseURL }) => {
  for (const [method, url] of [["GET", "/api/admin/events"], ["GET", "/api/admin/events/1"], ["POST", "/api/admin/events"], ["PUT", "/api/admin/events/1"], ["DELETE", "/api/admin/events/1"]]) {
    const response = await request.fetch(url, { method, headers: { Origin: new URL(baseURL!).origin, "Content-Type": "application/json" }, ...(method === "GET" ? {} : { data: {} }) });
    expect(response.status()).toBe(401);
  }
});

test("mieszkaniec nie ma dostępu administracyjnego", async ({ page, baseURL }) => {
  test.skip(!residentEmail || !residentPassword, "Ustaw E2E_RESIDENT_EMAIL i E2E_RESIDENT_PASSWORD.");
  await login(page, baseURL!, residentEmail!, residentPassword!);
  for (const method of ["GET", "POST", "PUT", "DELETE"]) {
    const response = await page.request.fetch(`/api/admin/events${["PUT", "DELETE"].includes(method) ? "/1" : ""}`, { method, headers: { Origin: new URL(baseURL!).origin, "Content-Type": "application/json" }, ...(method === "GET" ? {} : { data: {} }) });
    expect(response.status()).toBe(403);
  }
  await page.goto("/admin/wydarzenia");
  await expect(page.getByRole("heading", { name: "Brak uprawnień" })).toBeVisible();
});

test("pełny CRUD, publikacja, CSRF i zachowanie szkicu po wygaśnięciu sesji", async ({ page, baseURL }) => {
  test.skip(!adminEmail || !adminPassword, "Ustaw E2E_ADMIN_EMAIL i E2E_ADMIN_PASSWORD; użyj testowej bazy.");
  await login(page, baseURL!, adminEmail!, adminPassword!);
  await page.goto("/admin/wydarzenia");
  await page.getByRole("button", { name: "Dodaj wydarzenie" }).click();
  const title = `E2E wydarzenie ${Date.now()}`;
  let id: number | undefined;
  try {
    await page.getByLabel("Tytuł", { exact: true }).fill(title);
    await page.getByLabel("Miejsce", { exact: true }).fill("Sala testowa");
    await page.getByLabel("Opis", { exact: true }).fill("Opis z testu integracyjnego");
    await page.getByLabel("Rozpoczęcie", { exact: true }).fill("2030-12-01T17:00");
    const created = page.waitForResponse(r => r.url().endsWith("/api/admin/events") && r.request().method() === "POST");
    await page.getByRole("button", { name: "Zapisz wydarzenie" }).click();
    const response = await created;
    expect(response.status()).toBe(201);
    id = (await response.json()).id;
    expect(response.headers().location).toBe(`/api/admin/events/${id}`);
    expect((await page.request.get(`/api/events/${id}`)).status()).toBe(404);
    let card = page.locator("article").filter({ has: page.getByRole("heading", { name: title, exact: true }) });
    await expect(card.getByText("Ukryte", { exact: true })).toBeVisible();
    await card.getByRole("button", { name: "Opublikuj", exact: true }).click();
    await expect(card.getByText("Opublikowane", { exact: true })).toBeVisible();
    expect((await page.request.get(`/api/events/${id}`)).status()).toBe(200);
    const csrf = await page.request.delete(`/api/admin/events/${id}`, { headers: { Origin: "https://attacker.example", "Content-Type": "application/json" }, data: {} });
    expect(csrf.status()).toBe(403);
    await card.getByRole("button", { name: "Edytuj", exact: true }).click();
    await page.getByLabel("Opis", { exact: true }).fill("Szkic po wygaśnięciu sesji");
    await page.context().clearCookies();
    await page.getByRole("button", { name: "Zapisz wydarzenie" }).click();
    await expect(page.getByText("Sesja wygasła. Zaloguj się ponownie na konto administratora.")).toBeVisible();
    await page.getByRole("link", { name: "Zaloguj się", exact: true }).click();
    await page.getByLabel("Adres e-mail").fill(adminEmail!);
    await page.getByLabel("Hasło", { exact: true }).fill(adminPassword!);
    await page.getByRole("button", { name: "Zaloguj się", exact: true }).click();
    await expect(page.getByLabel("Opis", { exact: true })).toHaveValue("Szkic po wygaśnięciu sesji");
    await page.getByRole("button", { name: "Zapisz wydarzenie" }).click();
    await expect(page.getByText("Wydarzenie zostało zapisane.")).toBeVisible();
    card = page.locator("article").filter({ has: page.getByRole("heading", { name: title, exact: true }) });
    await card.getByRole("button", { name: "Ukryj", exact: true }).click();
    await expect(card.getByText("Ukryte", { exact: true })).toBeVisible();
    expect((await page.request.get(`/api/events/${id}`)).status()).toBe(404);
    page.once("dialog", async dialog => { expect(dialog.message()).toContain(title); await dialog.accept(); });
    await card.getByRole("button", { name: "Usuń", exact: true }).click();
    await expect(page.getByText("Wydarzenie zostało usunięte.")).toBeVisible();
    expect((await page.request.get(`/api/admin/events/${id}`)).status()).toBe(404);
  } finally {
    if (id) await page.request.delete(`/api/admin/events/${id}`, { headers: { Origin: new URL(baseURL!).origin, "Content-Type": "application/json" }, data: {} });
  }
});

test("edycja zachowuje drugą godzinę jesienną, a formularz odrzuca godzinę nieistniejącą wiosną", async ({ page }) => {
  const event = { ...fixture, starts_at: "2026-10-25T01:30:00Z", ends_at: "2026-10-25T02:00:00Z" };
  await page.route("**/api/auth/me", route => route.fulfill({ json: { user: { id: 1, email: "admin@example.test", role: "admin" } } }));
  await page.route("**/api/admin/events", route => route.fulfill({ json: { items: [event] } }));
  await page.route("**/api/admin/events/987", route => route.request().method() === "PUT"
    ? route.fulfill({ status: 204 }) : route.fulfill({ json: event }));
  await page.goto("/admin/wydarzenia");
  await page.getByRole("button", { name: "Edytuj", exact: true }).click();
  await page.getByLabel("Tytuł", { exact: true }).fill("Zmieniony tytuł");
  const sent = page.waitForRequest(r => r.url().endsWith("/api/admin/events/987") && r.method() === "PUT");
  await page.getByRole("button", { name: "Zapisz wydarzenie" }).click();
  expect((await sent).postDataJSON().starts_at).toBe("2026-10-25T01:30:00Z");
  await expect(page.getByText("Wydarzenie zostało zapisane.")).toBeVisible();
  await page.getByRole("button", { name: "Dodaj wydarzenie" }).click();
  await page.getByLabel("Tytuł", { exact: true }).fill("Wiosenny termin");
  await page.getByLabel("Miejsce", { exact: true }).fill("Sala");
  await page.getByLabel("Opis", { exact: true }).fill("Opis");
  await page.getByLabel("Rozpoczęcie", { exact: true }).fill("2027-03-28T02:30");
  await page.getByRole("button", { name: "Zapisz wydarzenie" }).click();
  await expect(page.getByText(/Ta godzina nie istnieje/)).toBeVisible();
});
