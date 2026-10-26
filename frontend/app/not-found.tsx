import Link from "next/link";

export default function NotFound() {
  return <section className="container py-5"><div className="placeholder-card p-4 p-md-5">
    <span className="badge text-bg-secondary mb-3">404</span>
    <h1 className="display-6 fw-semibold">Nie znaleziono strony</h1>
    <p className="lead text-secondary">Sprawdź adres lub wróć na stronę główną.</p>
    <Link href="/" className="btn btn-primary mt-3">Wróć na stronę główną</Link>
  </div></section>;
}
