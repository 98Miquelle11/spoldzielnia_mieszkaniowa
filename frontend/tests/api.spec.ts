import { test, expect, type APIRequestContext } from "@playwright/test";

const email = process.env.E2E_ADMIN_EMAIL;
const password = process.env.E2E_ADMIN_PASSWORD;
async function login(request: APIRequestContext, origin: string, email: string, password: string) {
  expect((await request.post("/api/auth/login", { headers: { Origin: origin }, data: { email, password } })).status()).toBe(200);
}

test("rzeczywiste publiczne API sortuje i zwraca UTC", async ({ request }) => {
  for (const period of ["upcoming", "past"]) {
    const response = await request.get(`/api/events?period=${period}`);
    expect(response.status()).toBe(200);
    const { items } = await response.json();
    expect(Array.isArray(items)).toBe(true);
    for (const item of items) {
      expect(item.is_published).toBe(true);
      for (const field of ["starts_at", "ends_at", "created_at", "updated_at"]) expect(item[field]).toMatch(/Z$/);
    }
    for (let i = 1; i < items.length; i++) {
      if (period === "upcoming") expect(Date.parse(items[i].starts_at)).toBeGreaterThanOrEqual(Date.parse(items[i-1].starts_at));
      else expect(Date.parse(items[i].ends_at)).toBeLessThanOrEqual(Date.parse(items[i-1].ends_at));
    }
  }
  expect((await request.get("/api/events?period=invalid")).status()).toBe(400);
  expect((await request.get("/api/events/abc")).status()).toBe(400);
  expect((await request.get("/api/events/9223372036854775807")).status()).toBe(404);
});

test("rzeczywiste CRUD, walidacja, ukrywanie i ochrona zapisu", async ({ request, baseURL }) => {
  test.skip(!email || !password, "Wymagane konto admin E2E i testowa baza.");
  const origin = new URL(baseURL!).origin;
  await login(request, origin, email!, password!);
  const headers = { Origin: origin };
  const input = { title: `API E2E ${Date.now()}`, description: "Opis", location: "Sala", starts_at: "2030-11-21T17:00:00+01:00", is_published: false };
  let id: number | undefined;
  try {
    for (const invalid of [{ ...input, title: " " }, { ...input, title: "ą".repeat(201) }, { ...input, ends_at: "2030-11-20T16:00:00Z" }, { ...input, is_published: null }]) {
      expect((await request.post("/api/admin/events", { headers, data: invalid })).status()).toBe(400);
    }
    const created = await request.post("/api/admin/events", { headers, data: input });
    expect(created.status()).toBe(201);
    id = (await created.json()).id;
    expect(created.headers().location).toBe(`/api/admin/events/${id}`);
    const details = await request.get(`/api/admin/events/${id}`);
    expect(details.status()).toBe(200);
    const record = await details.json();
    expect(record.starts_at).toBe("2030-11-21T16:00:00Z");
    expect(record.ends_at).toBe(record.starts_at);
    expect(record.created_at).toMatch(/Z$/);
    expect(record.updated_at).toMatch(/Z$/);
    expect((await request.get(`/api/events/${id}`)).status()).toBe(404);
    const { items } = await (await request.get("/api/admin/events")).json();
    expect(items.some((e: { id: number }) => e.id === id)).toBe(true);
    for (const method of ["POST", "PUT", "DELETE"]) {
      const path = method === "POST" ? "/api/admin/events" : `/api/admin/events/${id}`;
      expect((await request.fetch(path, { method, headers: { Origin: "https://attacker.example" }, data: {} })).status()).toBe(403);
      expect((await request.fetch(path, { method, headers: { Origin: origin, "Content-Type": "text/plain" }, data: "{}" })).status()).toBe(415);
    }
    const published = { ...input, is_published: true };
    const updated = await request.put(`/api/admin/events/${id}`, { headers, data: published });
    expect(updated.status()).toBe(204);
    expect(await updated.text()).toBe("");
    expect((await request.get(`/api/events/${id}`)).status()).toBe(200);
    const hidden = await request.put(`/api/admin/events/${id}`, { headers, data: input });
    expect(hidden.status()).toBe(204);
    expect((await request.get(`/api/events/${id}`)).status()).toBe(404);
    const deleted = await request.delete(`/api/admin/events/${id}`, { headers, data: {} });
    expect(deleted.status()).toBe(204);
    expect(await deleted.text()).toBe("");
    expect((await request.get(`/api/admin/events/${id}`)).status()).toBe(404);
    expect((await request.delete(`/api/admin/events/${id}`, { headers, data: {} })).status()).toBe(404);
    const logout = await request.post("/api/auth/logout", { headers, data: {} });
    expect(logout.status()).toBe(200);
    expect((await request.get("/api/admin/events")).status()).toBe(401);
  } finally {
    if (id) await request.delete(`/api/admin/events/${id}`, { headers, data: {} });
  }
});

test("rzeczywiste API blokuje mieszkańca dla wszystkich operacji", async ({ request, baseURL }) => {
  const residentEmail = process.env.E2E_RESIDENT_EMAIL;
  const residentPassword = process.env.E2E_RESIDENT_PASSWORD;
  test.skip(!residentEmail || !residentPassword, "Wymagane konto mieszkańca E2E.");
  const origin = new URL(baseURL!).origin;
  await login(request, origin, residentEmail!, residentPassword!);
  for (const [method, url] of [["GET", "/api/admin/events"], ["GET", "/api/admin/events/1"], ["POST", "/api/admin/events"], ["PUT", "/api/admin/events/1"], ["DELETE", "/api/admin/events/1"]]) {
    expect((await request.fetch(url, { method, headers: { Origin: origin, "Content-Type": "application/json" }, ...(method === "GET" ? {} : { data: {} }) })).status()).toBe(403);
  }
});
