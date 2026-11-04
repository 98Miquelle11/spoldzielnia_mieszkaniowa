import type { Metadata } from "next";
import { PublicEventDetails } from "@/components/public-events";
export const metadata: Metadata = { title: "Szczegóły wydarzenia" };
export default async function EventPage({ params, searchParams }: {
  params: Promise<{ id: string }>;
  searchParams: Promise<{ period?: string | string[] }>;
}) {
  const { id } = await params;
  const { period } = await searchParams;
  return <PublicEventDetails key={id} id={id} period={period === "past" ? "past" : "upcoming"} />;
}
