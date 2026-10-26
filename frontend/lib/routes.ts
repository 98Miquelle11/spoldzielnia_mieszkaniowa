export type Route = { slug: string; title: string; group: string | null; description: string };

export const routes: Route[] = [
  { slug: "start", title: "Strona główna", group: null, description: "Aktualności i najważniejsze sprawy osiedla." },
  { slug: "struktura", title: "Struktura", group: "Spółdzielnia", description: "Zarząd i organizacja spółdzielni." },
  { slug: "historia", title: "Historia", group: "Spółdzielnia", description: "Historia powstania osiedla Mirmiłowo." },
  { slug: "imprezy", title: "Imprezy kulturalne", group: "Ogłoszenia", description: "Wydarzenia organizowane dla mieszkańców." },
  { slug: "przetargi", title: "Przetargi", group: "Ogłoszenia", description: "Ogłoszenia o przetargach i terminach." },
  { slug: "wynajem", title: "Wynajem", group: "Ogłoszenia", description: "Dostępne oferty lokali i pomieszczeń." },
  { slug: "kalendarz", title: "Kalendarz", group: "Ogłoszenia", description: "Terminy, wydarzenia i sprawy pilne." },
  { slug: "inwestycje", title: "Inwestycje", group: null, description: "Postęp inwestycji na osiedlu." },
  { slug: "druki", title: "Druki i wnioski", group: "Strefa mieszkańca", description: "Formularze do pobrania." },
  { slug: "regulaminy", title: "Statut i regulaminy", group: "Strefa mieszkańca", description: "Zasady działania spółdzielni." },
  { slug: "czat", title: "Czat mieszkańców", group: "Strefa mieszkańca", description: "Miejsce rozmów zalogowanych mieszkańców." },
  { slug: "e-mieszkaniec", title: "E-mieszkaniec", group: null, description: "Obsługa spraw i wniosków online." },
  { slug: "kontakt", title: "Kontakt", group: null, description: "Dane kontaktowe i formularz." },
  { slug: "login", title: "Logowanie", group: null, description: "Logowanie do konta mieszkańca lub administratora." },
];

export const pageHref = (slug: string) => slug === "start" ? "/" : `/${slug}`;
export const findRoute = (slug: string) => routes.find(route => route.slug === slug);

export const cards = [
  { title: "Ogłoszenia", desc: "Wydarzenia, przetargi i oferty wynajmu.", slug: "imprezy" },
  { title: "Kalendarz", desc: "Terminy, spotkania i sprawy pilne.", slug: "kalendarz" },
  { title: "Inwestycje", desc: "Zmiany i prace na terenie osiedla.", slug: "inwestycje" },
  { title: "Strefa mieszkańca", desc: "Dokumenty, regulaminy i czat.", slug: "druki" },
  { title: "E-mieszkaniec", desc: "Sprawy i wnioski online — moduł planowany.", slug: "e-mieszkaniec" },
  { title: "Kontakt", desc: "Jak skontaktować się ze spółdzielnią.", slug: "kontakt" },
];
