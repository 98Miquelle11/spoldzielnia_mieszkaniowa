import Link from "next/link";
import { cards, pageHref } from "@/lib/routes";

export default function Home() {
  return <>
    <section className="hero py-5"><div className="container py-lg-4">
      <p className="text-uppercase fw-semibold mb-2 text-muted small letter-spacing">Jedno osiedle. Wspólna przestrzeń.</p>
      <h1 className="display-5 fw-bold">Spółdzielnia Mieszkaniowa „Mirmiłowo”</h1>
      <p className="lead mt-3 mb-4 hero-lead">Informacje o życiu osiedla, terminach, inwestycjach i sprawach mieszkańców w jednym miejscu.</p>
      <Link href="/kalendarz" className="btn btn-primary btn-lg">Sprawdź kalendarz</Link>
      <Link href="/e-mieszkaniec" className="btn btn-outline-primary btn-lg ms-2">E-mieszkaniec</Link>
    </div></section>
    <section className="container py-5" aria-labelledby="sekcje-tytul">
      <h2 className="h3 mb-2" id="sekcje-tytul">Najważniejsze sekcje</h2>
      <p className="text-secondary mb-4">Dostępna jest nawigacja i strony informacyjne. Funkcje będą wdrażane w kolejnych etapach.</p>
      <div className="row g-4">{cards.map(card => <div className="col-md-6 col-lg-4" key={card.slug}>
        <article className="card h-100 border-0 shadow-sm"><div className="card-body p-4">
          <h3 className="h5 card-title">{card.title}</h3><p className="card-text text-secondary">{card.desc}</p>
          <Link href={pageHref(card.slug)} className="stretched-link link-primary">Przejdź do sekcji <span aria-hidden="true">→</span></Link>
        </div></article>
      </div>)}</div>
    </section>
  </>;
}
