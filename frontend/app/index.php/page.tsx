import { redirect, notFound } from "next/navigation";
import { findRoute, pageHref } from "@/lib/routes";

// Compatibility with links/bookmarks from the PHP application.
export default async function Legacy({ searchParams }: { searchParams: Promise<Record<string, string | string[] | undefined>> }) {
  const query = await searchParams;
  if (Object.keys(query).some(key => key.startsWith("page["))) notFound();
  const slug = query.page ?? "start";
  if (typeof slug !== "string" || !findRoute(slug)) notFound();
  redirect(pageHref(slug));
}
