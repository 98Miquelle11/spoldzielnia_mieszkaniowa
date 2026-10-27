# Mirmiłowo — Go + Next.js + MySQL

Backend jest napisany w Go, frontend w Next.js/React/TypeScript. Baza to
MySQL/MariaDB.

## Co działa

- Strona główna i sześć kafelków sekcji.
- Wszystkie pozycje nawigacji z oryginalnego `config/routes.php`.
- Podstrony „Etap w przygotowaniu”.
- Logowanie, odczyt bieżącej sesji i wylogowanie przez API Go.
- Konta `mieszkaniec` i `admin`; narzędzie do ich tworzenia.
- Trwałe sesje MySQL: pozostają aktywne po restarcie backendu, wygasają po 8 h.
- Przekierowania starych adresów `index.php?page=...`.

Pozostałe moduły w pierwotnym repozytorium były jedynie planowane. Czat,
wnioski, kalendarz, zarządzanie ogłoszeniami i panel administratora nadal mają
strony informacyjne. Migracja nie dodaje ich docelowej funkcjonalności.

## Wymagania

- Go 1.26 lub nowszy — https://go.dev/dl/
- Node.js 20.9 lub nowszy (zalecany aktualny LTS) — https://nodejs.org/
- MySQL 8 lub MariaDB; można użyć MySQL/MariaDB z istniejącego XAMPP.

Apache i PHP nie obsługują nowej aplikacji. W XAMPP możesz korzystać z Apache
wyłącznie po to, aby otworzyć phpMyAdmin i zaimportować bazę.

## 1. Utwórz bazę

W XAMPP uruchom MySQL oraz Apache. Otwórz http://localhost/phpmyadmin,
wybierz zakładkę **Import** na poziomie serwera i wskaż `database/schema.sql`.
Plik sam tworzy bazę `mirmilowo` oraz tabele `users` i `sessions`.

Alternatywnie w terminalu klienta MySQL:

```sql
SOURCE /pelna/sciezka/do/mirmilowo/database/schema.sql;
```

`database/seed.sql` nie dodaje kont ze wspólnymi hasłami. Konta utworzysz w kroku 3.
Ponowny import schematu nie usuwa danych. `IF NOT EXISTS` nie aktualizuje jednak
struktury już istniejących tabel — dla przyszłych zmian przygotuj migracje SQL.

## 2. Uruchom backend Go

Otwórz PowerShell w folderze projektu, a następnie:

```powershell
cd backend
$env:DB_ADDR = "127.0.0.1:3306"
$env:DB_NAME = "mirmilowo"
$env:DB_USER = "root"
$env:DB_PASSWORD = ""
$env:APP_ORIGIN = "http://localhost:3000"
go mod download
go run ./cmd/server
```

Wartości powyżej odpowiadają typowej lokalnej instalacji XAMPP. Jeśli konto root
ma hasło lub serwer używa innego portu, wpisz własne wartości. Dostępne ustawienia
opisuje `backend/.env.example`. Go czyta zmienne procesu; nie wczytuje tego pliku
automatycznie.

Na macOS/Linux odpowiednikiem ustawienia jest np.:

```bash
export DB_ADDR=127.0.0.1:3306 DB_NAME=mirmilowo DB_USER=root DB_PASSWORD=''
export APP_ORIGIN=http://localhost:3000
go mod download
go run ./cmd/server
```

Backend domyślnie nasłuchuje lokalnie na `127.0.0.1:8080`.

## 3. Utwórz konta

W drugim terminalu (koniecznie drugim, nie wyłączaj pierwszego) przejdź do `backend` i ustaw te same zmienne bazy co w kroku 2.
Zmienne PowerShell obowiązują tylko w terminalu, w którym je ustawiono.

```powershell
go run ./cmd/create-user -email admin@example.test -role admin
go run ./cmd/create-user -email mieszkaniec@example.test -role mieszkaniec
```

Każde polecenie poprosi o hasło, które nie będzie widoczne podczas wpisywania.
Wybierz hasło o długości 12–72 bajtów (polskie znaki mogą zajmować więcej niż bajt).
W bazie zostanie zapisany wyłącznie hash bcrypt. Powtórzenie polecenia dla tego
samego adresu zgłosi błąd, zamiast nadpisać konto.

## 4. Uruchom frontend Next.js

W kolejnym terminalu, z folderu głównego projektu:

```powershell
cd frontend
npm ci
npm run dev
```

Otwórz **http://localhost:3000**. Formularz logowania znajduje się pod `/login`.
Backend Go i MySQL muszą nadal działać. Nie otwieraj pliku HTML bezpośrednio ani
adresu starej aplikacji w `htdocs`.

Domyślne ustawienia działają bez `.env.local`. Jeśli zmienisz adres backendu,
skopiuj `frontend/.env.example` do `frontend/.env.local` i ustaw `API_INTERNAL_URL`.
Po zmianie uruchom Next.js ponownie; dla wersji produkcyjnej wykonaj nowy build.

Używaj dokładnie hosta z `APP_ORIGIN`: `localhost` i `127.0.0.1` to różne originy.
Jeśli Next.js proponuje port 3001, zwolnij port 3000 albo odpowiednio zmień
`APP_ORIGIN` w backendzie i uruchom go ponownie.

## Testy i kompilacja

Backend (w `backend`):

```powershell
go test ./...
go vet ./...
go build ./cmd/server ./cmd/create-user
```

Testy używają atrapy SQL i nie wymagają działającej bazy. Sprawdzają m.in.
logowanie, rotację i usuwanie sesji, błędne hasła, origin/CSRF, role i awarie SQL.

Frontend (w `frontend`):

```powershell
npm run typecheck
npm run build
npm start
```

`npm start` uruchamia wcześniej zbudowaną wersję. Nie uruchamiaj go równocześnie
z `npm run dev` na tym samym porcie.

Test ręczny z prawdziwą bazą: zaloguj oba konta; sprawdź komunikat dla błędnego
hasła; odśwież stronę po zalogowaniu; zrestartuj backend i odśwież ponownie;
wyloguj się i sprawdź, że `/api/auth/me` zwraca `{"user":null}`. Otwórz też menu na
telefonie i nieistniejący adres. Szczegóły migracji: `docs/MIGRACJA.md`.

## Struktura

```text
backend/cmd/server/          serwer HTTP Go
backend/cmd/create-user/     tworzenie kont
backend/internal/app/        API, sesje, kontrola ról, testy
backend/internal/database/   połączenie z MySQL
frontend/app/               strony Next.js i CSS
frontend/components/        menu, formularz i stan logowania
frontend/lib/routes.ts      polskie tytuły, grupy i opisy
database/schema.sql         struktura MySQL
database/seed.sql           wskazówki tworzenia kont
docs/                       opis migracji i dalszych etapów
```

## Konfiguracja i utrzymanie

W przeglądarce API działa pod `/api`; Next.js przekazuje żądania do Go. Nie trzeba
włączać CORS ani przechowywać tokena w localStorage. Cookie sesji ma `HttpOnly`,
`SameSite=Lax`, a dla HTTPS także `Secure`. Zmiany sesji wymagają JSON i dokładnie
zgodnego nagłówka Origin. Wylogowanie używa POST, więc samo wejście na adres
nie kończy sesji. Stary `/src/logout.php` kieruje do strony z przyciskiem.

Wygasłe sesje są odrzucane przy odczycie. Okresowo usuń je poleceniem SQL:

```sql
DELETE FROM mirmilowo.sessions WHERE expires_at <= UTC_TIMESTAMP();
```

Przyszłe prywatne endpointy zabezpieczaj przez `RequireRole` po stronie Go.
Obecny frontend rozpoznaje rolę, ale nie zawiera panelu administratora.
Przed publicznym wdrożeniem skonfiguruj HTTPS, własne konto bazy z ograniczonymi
uprawnieniami i ograniczanie prób logowania na reverse proxy. Nie udostępniaj
publicznie bazy danych. Hasła i lokalne pliki środowiska pozostaw poza Git.

## Praca z GitHub

Rozpakuj paczkę poza oryginalnym projektem i sprawdź nową wersję lokalnie.
Następnie na osobnej gałęzi oryginalnego repozytorium dodaj strukturę z tej paczki
(bez `node_modules`, `.next` i własnych sekretów), zastępując stary punkt wejścia
PHP oraz nieaktualne instrukcje. Zachowaj historię Git i plik LICENSE. Paczka
nie zawiera metadanych `.git`; niczego nie opublikowano automatycznie na GitHub.
