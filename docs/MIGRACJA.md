# Zakres migracji

Przeniesiono istniejący kod dostarczonego archiwum, a nie pełny zakres planowany
w pierwotnych 12 etapach. Ponieważ oryginalne pliki SQL zawierały wyłącznie
komentarze, dodano nowy schemat kont i sesji. Nie było danych do przeniesienia.

| Oryginalny element | Nowy odpowiednik |
| --- | --- |
| `index.php` i `config/routes.php` | Next.js App Router, `frontend/lib/routes.ts` |
| `pages/home.php` | `frontend/app/page.tsx` |
| `pages/placeholder.php` | `frontend/app/[slug]/page.tsx`, `not-found.tsx` |
| `includes/header.php` | `frontend/components/header.tsx` |
| `includes/footer.php` | `frontend/app/layout.tsx` |
| `assets/css/style.css` | `frontend/app/globals.css` |
| `pages/login.php` | formularz React i `POST /api/auth/login` w Go |
| `includes/auth.php` | sesje Go/MySQL i `GET /api/auth/me` |
| `src/logout.php` | `POST /api/auth/logout` |
| `config/database.php` | `backend/internal/database/database.go`, zmienne środowiska |
| puste pliki bazy | `database/schema.sql`, `database/seed.sql` |

## API

| Metoda | Adres | Wynik |
| --- | --- | --- |
| GET | `/api/health` | Stan połączenia z bazą; 503 przy awarii |
| GET | `/api/auth/me` | `{"user":null}` albo użytkownik bez hasła/tokena |
| POST | `/api/auth/login` | JSON `email`, `password`; ustawia cookie i zwraca użytkownika |
| POST | `/api/auth/logout` | Unieważnia sesję w bazie i usuwa cookie |

POST wymaga `Content-Type: application/json` i Origin równego `APP_ORIGIN`.
Brak/błąd danych to 400, błędne dane logowania 401, zły Origin 403,
nieobsługiwany format 415, problem z bazą 503. API nie zwraca hashy haseł.
W tabeli `sessions` znajduje się SHA-256 losowego tokena, nie sam token.
Aktualna rola jest odczytywana z `users`, więc zmiana roli działa przy następnym
żądaniu. Czas wygaśnięcia sesji zapisywany jest w UTC.

## Różnice względem PHP

- PHP i sesje PHP nie są już potrzebne. CSS pozostaje częścią frontendu.
- Hasło nie jest przycinane przez `trim`; spacje należą do hasła.
- Sesja jest rotowana po poprawnym logowaniu i usuwana po wylogowaniu.
- Zamiast wylogowywania metodą GET używany jest POST.
- Komunikaty błędów bazy nie ujawniają danych połączenia użytkownikowi.
- Logowanie pojawia się w nawigacji tylko raz.
- Stare linki PHP kierują do nowych adresów. Nieznana strona zwraca 404.
- Nieistniejące moduły pozostają publicznymi stronami informacyjnymi. Nie
  zawierają żadnych prywatnych danych ani API funkcji mieszkańca.

## Dalsze etapy

1. Organizacja i historia spółdzielni — treści SQL i ich edycja.
2. Wydarzenia kulturalne, przetargi i oferty wynajmu — publiczne widoki i CRUD Go.
3. Kalendarz — integracja terminów, strefa Europe/Warsaw.
4. Inwestycje — statusy i postęp.
5. Dokumenty — kontrolowane przesyłanie i pobieranie.
6. Czat — prywatne endpointy dla mieszkańca i administratora.
7. E-mieszkaniec — własne wnioski, załączniki i statusy.
8. Kontakt i administracja — formularz i panel.

Każdy nowy moduł wymaga własnych migracji SQL, kontroli uprawnień w Go,
walidacji i testów. Wygląd i formularze rozwijaj w Next.js. Ewentualny podpis
Profilem Zaufanym pozostaje edukacyjną symulacją zgodnie z pierwotnym planem.

## Źródła techniczne

- https://nextjs.org/docs/app/getting-started/installation
- https://pkg.go.dev/golang.org/x/crypto/bcrypt
- https://pkg.go.dev/github.com/go-sql-driver/mysql
