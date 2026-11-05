import { test, expect } from "@playwright/test";
import { routes, pageHref } from "../lib/routes";

test.beforeEach(async ({ page }) => {
  await page.route("**/api/auth/me", route => route.fulfill({ json: { user: null } }));
  await page.route("**/api/events?*", route => route.fulfill({ json: {
    items: [], pagination: { page: 1, limit: 20, has_more: false },
  } }));
});

for (const route of routes) {
  test(`dymny test trasy ${pageHref(route.slug)}`, async ({ page }) => {
    const response = await page.goto(pageHref(route.slug));
    expect(response?.status()).toBe(200);
    const heading = route.slug === "start" ? "Witaj w Mirmiłowie"
      : route.title;
    await expect(page.locator("main h1")).toHaveText(heading);
    await expect(page.locator("main")).not.toContainText("Wydarzenie niedostępne");
    if (route.slug === "start") await expect(page.locator("main article")).toHaveCount(6);
  });
}

test("powrót z kontaktu prowadzi na stronę główną", async ({ page }) => {
  await page.goto("/kontakt");
  await page.getByRole("link", { name: "Wróć na stronę główną" }).click();
  await expect(page).toHaveURL(/\/$/);
  await expect(page.getByRole("heading", { name: "Witaj w Mirmiłowie" })).toBeVisible();
});

test("menu administratora jest rozwinięte na laptopach i nie przepełnia ekranu", async ({ page }) => {
  await page.route("**/api/auth/me", route => route.fulfill({ json: { user: { id: 1, email: "admin@example.test", role: "admin" } } }));
  await page.goto("/");
  for (const width of [992, 1024, 1199, 1200, 1280, 1399, 1400]) {
    await page.setViewportSize({ width, height: 900 });
    await expect(page.getByRole("button", { name: "Rozwiń menu" })).toBeHidden();
    await expect(page.getByRole("link", { name: "Panel wydarzeń" })).toBeVisible();
    const logout = page.getByRole("button", { name: "Wyloguj się" });
    await expect(logout).toBeVisible();
    expect(await page.evaluate(() => document.documentElement.scrollWidth)).toBeLessThanOrEqual(width);
    const box = await logout.boundingBox();
    expect(box!.x + box!.width).toBeLessThanOrEqual(width);
  }
});
