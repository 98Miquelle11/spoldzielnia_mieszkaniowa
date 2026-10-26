import type { Metadata } from "next";
import Link from "next/link";
import { notFound, redirect } from "next/navigation";
import { findRoute } from "@/lib/routes";

type Props = { params: Promise<{ slug: string }> };
export async function generateMetadata({ params }: Props): Promise<Metadata> {
  const { slug } = await params;
  return { title: findRoute(slug)?.title ?? "Nie znaleziono strony" };
}

export default async function Section({ params }: Props) {
  const { slug } = await params;
  const route = findRoute(slug);
  if (!route) notFound();
  if (slug === "start") redirect("/");
  return <section className="container py-5">
    <nav aria-label="Okruszki"><ol className="breadcrumb small"><li className="breadcrumb-item"><Link href="/">Strona główna</Link></li><li className="breadcrumb-item active" aria-current="page">{route.title}</li></ol></nav>
    <div className="placeholder-card p-4 p-md-5 mt-3">
      <span className="badge text-bg-secondary mb-3">Etap w przygotowaniu</span>
      <h1 className="display-6 fw-semibold">{route.title}</h1>
      <p className="lead text-secondary">{route.description}</p>
      <p>Ta podstrona jest już dostępna w nawigacji. Jej właściwa funkcjonalność zostanie dodana w jednym z kolejnych etapów projektu.</p>
      <Link href="/" className="btn btn-primary mt-3">Wróć na stronę główną</Link>
    </div>
  </section>;
}
