import { defineConfig } from "@playwright/test";
export default defineConfig({
  testDir: "./tests", workers: 1, timeout: 30000,
  use: {
    baseURL: process.env.E2E_BASE_URL ?? "http://localhost:3000",
    timezoneId: "Europe/Warsaw",
    launchOptions: process.env.E2E_CHROMIUM_PATH ? { executablePath: process.env.E2E_CHROMIUM_PATH, args: ["--no-sandbox"] } : {},
  },
});
