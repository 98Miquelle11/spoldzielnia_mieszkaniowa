# Weryfikacja dostarczonej wersji

Sprawdzenie wykonano 26.10.2026.

## Zakończone poprawnie

- `go test ./...` — scenariusze logowania, złego hasła, braku konta, rotacji
  sesji, odczytu sesji, wylogowania, wygaśnięcia, kontroli origin, walidacji
  JSON, błędów SQL i uprawnień ról.
- `go vet ./...` i kompilacja obu poleceń Go.
- `npm run typecheck` i `npm run build`.
- Próby HTTP z produkcyjnym serwerem Next.js: strona główna, sześć kart,
  formularz logowania i 12 stron informacyjnych.
- Nieznana podstrona zwraca HTTP 404.
- Stare adresy PHP przekierowują; nieprawidłowy tablicowy parametr strony
  jest odrzucany kodem 404.
- Proxy Next.js przetestowane z lokalną atrapą API: zachowanie kodu 401,
  przekazywanie Origin, JSON, cookie HttpOnly, odczyt sesji i wylogowanie.

Środowisko kompilacji: Go 1.27.3, Node.js 24.22.1, Next.js 16.4.3,
TypeScript 5.9.3. W środowisku testowym brakowało części interfejsu `/proc`;
do sprawdzenia kompilacji zastosowano lokalne obejścia pomiaru pamięci Node
i wyłączono CGO/oznaczanie wersji VCS w Go. Nie są częścią aplikacji i nie są
potrzebne w standardowej instalacji Windows, Linux lub macOS.

## Do sprawdzenia lokalnie

Testy Go korzystają z atrapy SQL. W tym środowisku nie uruchomiono prawdziwego
serwera MySQL/MariaDB, więc nie wykonano importu schematu ani pełnego testu
logowania na prawdziwej bazie. Nie wykonano także testu interakcji i wyglądu w
prawdziwej przeglądarce — pobranie silnika przeglądarki w środowisku testowym
było niedostępne. Struktura i CSS pochodzą z dostarczonego projektu.

Po imporcie `database/schema.sql` i utworzeniu kont wykonaj scenariusz ręczny
z README, w tym menu na małym ekranie. Udane testy z atrapami nie zastępują
testu integracyjnego z własną bazą.
