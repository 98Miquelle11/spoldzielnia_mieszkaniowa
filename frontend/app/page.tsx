import type { Metadata } from "next";
import { PublicEvents } from "@/components/public-events";
export const metadata: Metadata = { title: "Imprezy kulturalne" };
export default async function EventsPage({ searchParams }: {
  searchParams: Promise<{ period?: string | string[] }>;
}) {
  const { period } = await searchParams;
  const selected = period === undefined || period === "" ? "upcoming"
    : period === "upcoming" || period === "past" ? period : null;
  return <PublicEvents key={selected ?? "invalid"} period={selected} />;
}
