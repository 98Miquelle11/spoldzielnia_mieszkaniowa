import Link from "next/link";
import { cards, pageHref } from "@/lib/routes";

export default function HomePage() {
  return <>
    <section className="hero py-5">
      <div className="container py-md-4">
        <p className="small text-uppercase letter-spacing fw-semibold">Spółdzielnia Mieszkaniowa</p>
        <h1 className="display-4 fw-bold">Witaj w Mirmiłowie</h1>
        <p className="lead hero-lead">Aktualności i najważniejsze sprawy osiedla w jednym miejscu. Sprawdź wydarzenia, znajdź potrzebne dokumenty i skontaktuj się ze spółdzielnią.</p>
        <div className="d-flex flex-wrap gap-2 mt-4">
          <Link href="/imprezy" className="btn btn-primary">Imprezy kulturalne</Link>
          <Link href="/kontakt" className="btn btn-outline-primary">Kontakt</Link>
        </div>
      </div>
    </section>
    <section className="container py-5" aria-labelledby="sections-title">
      <h2 id="sections-title" className="h3 mb-4">Sprawy mieszkańców</h2>
      <div className="row g-4">
        {cards.map(card => <div key={card.slug} className="col-md-6 col-lg-4">
          <article className="card h-100 shadow-sm"><div className="card-body p-4 d-flex flex-column">
            <h3 className="h5">{card.title}</h3>
            <p className="text-secondary">{card.desc}</p>
            <Link href={pageHref(card.slug)} className="link-primary mt-auto" aria-label={`Przejdź do sekcji ${card.title}`}>Przejdź do sekcji →</Link>
          </div></article>
        </div>)}
      </div>
    </section>
  </>;
}
