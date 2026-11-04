# Spółdzielnia Mieszkaniowa „Mirmiłowo” — dokumentacja projektu

**Nazwa projektu:** Spółdzielnia Mieszkaniowa „Mirmiłowo”  

**Repozytorium:** `spoldzielnia_mieszkaniowa`  

**Rodzaj:** Aplikacja internetowa dla współdzielni mieszkaniowej 

**Środowisko docelowe:** lokalne środowisko deweloperskie — backend Go, frontend Next.js oraz MySQL/MariaDB 

**Technologie:** Go 1.26+, Next.js, React, TypeScript, Node.js 20.9+, MySQL 8 / MariaDB

**Zespół:** 12 osób  

**Stan projektu:** etap 1 wykonany; etap 2 częściowo zaimplementowany; dostępny jest również panel administracyjny wydarzeń

---

## 1. Opis projektu

Spółdzielnia Mieszkaniowa „Mirmiłowo” to projekt aplikacji internetowej obsługującej **jedno fikcyjne osiedle mieszkaniowe**. Ma stanowić wspólne miejsce publikowania informacji o spółdzielni, ogłoszeń, dokumentów i terminów oraz umożliwiać mieszkańcom kontakt ze spółdzielnią i elektroniczną obsługę wybranych spraw.

Aplikacja jest podzielona na dwa główne elementy: backend napisany w języku Go, który udostępnia API i obsługuje m.in. uwierzytelnianie oraz sesje, oraz frontend oparty na Next.js, React i TypeScript. Dane są przechowywane w bazie MySQL/MariaDB. Żądania przeglądarki do API są obsługiwane pod ścieżką `/api` przez Next.js, który przekazuje je do backendu Go.

Docelowo aplikacja będzie przeznaczona dla trzech grup użytkowników: gości, zalogowanych mieszkańców i administratorów. Zakres dostępnych funkcji będzie zależeć od roli użytkownika.

Projekt ma charakter dydaktyczny i jest przeznaczony do lokalnego uruchamiania. Nie zakłada wykorzystywania rzeczywistych danych mieszkańców. Planowany podpis za pomocą Profilu Zaufanego będzie **wyłącznie symulacją**, a nie integracją z usługą państwową.

## 2. Aktualny stan realizacji

### 2.1. Postęp według etapów

| Etap | Zakres | Stan na podstawie aktualnego kodu |
| --- | --- | --- |
| 1 | Repozytorium, struktura projektu, szkielet strony, nawigacja i README | **Wykonany** — elementy są obecne w projekcie |
| 2 | Baza danych, konta użytkowników, logowanie, wylogowanie, role i uprawnienia | **W trakcie** — działa logowanie, odczyt sesji i wylogowanie przez API Go; dostępne jest tworzenie kont oraz trwałe sesje MySQL. |
| 3 | Spółdzielnia: Struktura i Historia | Planowany |
| 4 | Ogłoszenia: Imprezy kulturalne | **Wykonany w aktualnym zakresie** — działa publiczna lista wydarzeń, szczegóły, filtrowanie oraz panel administratora umożliwiający zarządzanie wydarzeniami. |
| 5 | Ogłoszenia: Przetargi | Planowany |
| 6 | Ogłoszenia: Wynajem | Planowany |
| 7 | Kalendarz wydarzeń i terminów | Planowany |
| 8 | Inwestycje | Planowany |
| 9 | Strefa mieszkańca: Druki i wnioski, Statut i regulaminy | Planowany |
| 10 | Czat mieszkańców | Planowany |
| 11 | E-mieszkaniec | Planowany |
| 12 | Kontakt, integracja oraz testy końcowe | **Częściowo wykonany** — przygotowano podstawowe testy backendu i testy przeglądarkowe. Podstrona kontaktowa ma charakter informacyjny. |

### 2.2. Co jest już dostępne

W obecnej wersji projektu znajdują się:

- strona główna oraz sześć kafelków prowadzących do sekcji aplikacji;
- wszystkie pozycje nawigacji znane z wcześniejszej wersji projektu;
- podstrony informacyjne „Etap w przygotowaniu” dla modułów, których funkcjonalność nie została jeszcze wdrożona;
- logowanie, odczyt informacji o bieżącej sesji i wylogowanie, obsługiwane przez API backendu Go;
- konta użytkowników z rolami `mieszkaniec` i `admin`, tworzone za pomocą narzędzia konsolowego uruchamianego w PowerShell;
- możliwość ustawienia hasła podczas tworzenia konta; wpisywane hasło nie jest widoczne w konsoli, a w bazie danych przechowywany jest jego skrót wygenerowany za pomocą bcrypt;
- weryfikacja adresu e-mail i hasła podczas logowania; w przypadku nieprawidłowych danych wyświetlany jest komunikat „Nieprawidłowy e-mail lub hasło”;
- trwałe sesje przechowywane w MySQL, które pozostają aktywne po restarcie backendu i wygasają po 8 godzinach;
- panel wydarzeń administratora dostępny pod adresem `/admin/wydarzenia`, umożliwiający wyświetlanie listy wydarzeń, dodawanie, edytowanie, publikowanie, ukrywanie i usuwanie;
- testy backendu obejmujące m.in. logowanie, obsługę sesji, błędne hasła, role i błędy SQL.
- publiczna lista wydarzeń kulturalnych dostępna pod adresem `/imprezy`;
- widok szczegółów pojedynczego wydarzenia dostępny pod adresem `/imprezy/{id}`;
- filtrowanie wydarzeń na nadchodzące i archiwalne;
- wyświetlanie terminów wydarzeń z uwzględnieniem strefy czasowej `Europe/Warsaw`;
- stronicowanie listy wydarzeń, ograniczające liczbę wyświetlanych rekordów do 20 na stronę;
- obsługa pustych wyników wyszukiwania oraz błędów pobierania danych;
- możliwość dodawania, edytowania, publikowania, ukrywania i usuwania wydarzeń przez administratora;
- zabezpieczenie operacji administracyjnych poprzez weryfikację sesji i roli `admin` po stronie backendu Go.

**Ograniczenie zakresu:** pozostałe moduły opisane w dokumentacji nadal są planowane. Czat, wnioski, kalendarz i inne funkcjonalności wymagają dalszej implementacji.

## 3. Funkcjonalności aplikacji — obecnie i docelowo

### 3.1. Strona główna i nawigacja

**Obecnie:** działa szkielet serwisu. Użytkownik może przechodzić pomiędzy zdefiniowanymi podstronami, a większość z nich wyświetla informację o planowanej funkcjonalności.

**Docelowo:** strona główna ma kierować do rzeczywistych treści i usług, w tym bieżących ogłoszeń, istotnych terminów, inwestycji i funkcji przeznaczonych dla mieszkańców.

### 3.2. Spółdzielnia

**Obecnie:** dostępne są pozycje nawigacji „Struktura” i „Historia”, prowadzące do stron zastępczych.

**Docelowo:**

- **Struktura** — informacje o organizacji spółdzielni, zarządzie i podziale odpowiedzialności.
- **Historia** — informacje o powstaniu i rozwoju osiedla Mirmiłowo.

### 3.3. Ogłoszenia

**Obecnie:** menu udostępnia odnośniki do „Imprez kulturalnych”, „Przetargów”, „Wynajmu” i „Kalendarza”. Moduł „Imprezy kulturalne” posiada już działającą obsługę danych, natomiast pozostałe moduły mają obecnie charakter informacyjny.

**Docelowo:**

- **Imprezy kulturalne** — publikowanie wydarzeń i informacji o spotkaniach dla mieszkańców.
- **Przetargi** — prezentacja ogłoszeń o przetargach i związanych z nimi terminów.
- **Wynajem** — przegląd ofert wynajmu lokali lub pomieszczeń.
- **Kalendarz** — wspólny widok wydarzeń, ważnych terminów, spraw pilnych oraz terminów składania wniosków i ankiet. Kategorie powinny być oznaczane kolorami z czytelną legendą.

### 3.4. Inwestycje

**Obecnie:** odnośnik do strony zastępczej.

**Docelowo:** lista inwestycji realizowanych na osiedlu, ich opisy, szczegóły, statusy oraz postęp prac.

### 3.5. Strefa mieszkańca

**Obecnie:** menu zawiera „Druki i wnioski”, „Statut i regulaminy” oraz „Czat mieszkańców”, ale wyświetlane są jedynie strony zastępcze. Te adresy nie są jeszcze chronione mechanizmem autoryzacji.

**Docelowo:**

- **Druki i wnioski** — udostępnianie dokumentów i formularzy do pobrania.
- **Statut i regulaminy** — prezentacja zasad działania spółdzielni i dokumentów organizacyjnych.
- **Czat mieszkańców** — komunikacja między uprawnionymi, zalogowanymi mieszkańcami; przewidziana jest możliwość moderacji.

Dostęp do funkcji wymagających uwierzytelnienia powinien być sprawdzany po stronie serwera PHP, nie wyłącznie przez ukrywanie pozycji menu.

### 3.6. E-mieszkaniec

**Obecnie:** dostępny jest odnośnik oraz podstrona informacyjna bez obsługi wniosków.

**Docelowo:**

- elektroniczne wypełnianie i przesyłanie wniosków;
- dołączanie załączników;
- podgląd statusu zgłoszenia;
- demonstracja procesu podpisywania dokumentów Profilem Zaufanym.

**Zastrzeżenie:** podpis elektroniczny będzie symulowany, bez rzeczywistego uwierzytelnienia państwowego i bez skutków prawnych.

### 3.7. Kontakt

**Obecnie:** dostępna jest strona zastępcza.

**Docelowo:** dane kontaktowe spółdzielni oraz formularz umożliwiający przesłanie wiadomości.

### 3.8. Konta i administracja

**Obecnie:** dostępny jest formularz logowania, obsługa sesji użytkowników oraz możliwość wylogowania. Konta mieszkańców i administratorów są tworzone za pomocą narzędzia konsolowego. System rozróżnia role `mieszkaniec` i `admin`, a hasła są przechowywane w postaci skrótów bcrypt.

Dostępny jest również panel administracyjny wydarzeń pod adresem `/admin/wydarzenia`, który umożliwia dodawanie, edytowanie, publikowanie, ukrywanie i usuwanie wydarzeń.
Operacje administracyjne są zabezpieczone po stronie backendu Go. Dostęp do endpointów `/api/admin/events` wymaga zalogowania i posiadania roli `admin`. Kontrola uprawnień realizowana jest za pomocą mechanizmu `RequireRole`.

**Ograniczenia:** nie ma jeszcze pełnego panelu zarządzania kontami użytkowników ani wszystkich planowanych funkcji administracyjnych.

**Docelowo:** system kont z rolami i ograniczonym dostępem do funkcji, panel mieszkańca oraz panel administracyjny umożliwiający zarządzanie treściami, obsługę zgłoszeń i moderowanie czatu. Możliwość samodzielnej rejestracji została wymieniona w pierwotnym opisie projektu, ale jej sposób działania nie jest jeszcze określony i nie jest zaimplementowany.

## 4. Role użytkowników

| Rola | Uprawnienia docelowe | Stan obecny |
| --- | --- | --- |
| **Gość** | Przeglądanie publicznej części serwisu | Może korzystać z publicznych stron i stron zastępczych bez logowania |
| **Mieszkaniec** | Po zalogowaniu dostęp do prywatnych funkcji, w tym czatu i E-mieszkańca | Kod rozpoznaje rolę `mieszkaniec` w sesji, ale chronione funkcje nie są jeszcze gotowe |
| **Administrator** | Dostęp do panelu zarządzania treściami, zgłoszeniami i moderacji | Kod rozpoznaje rolę `admin` w sesji, ale panel i kompletne sprawdzanie dostępu nie istnieją |

Gość nie musi posiadać konta. Mieszkaniec i administrator docelowo mają być użytkownikami zapisanymi w tabeli `users`, z hasłami przechowywanymi jako hashe, a nie tekst jawny.

**Obecnie nie należy zakładać, że uprawnienia są skutecznie egzekwowane**, ponieważ w projekcie nie ma chronionych paneli ani kompletnego mechanizmu kontroli dostępu do tras.

## 5. Technologie i struktura projektu

### 5.1. Technologie

| Technologia | Zastosowanie |
| --- | --- |
| Go 1.26+ | Backend aplikacji: serwer HTTP, API, logika uwierzytelniania, obsługa sesji i kontrola ról |
| Next.js | Framework frontendowy; obsługa stron aplikacji oraz przekazywanie żądań przeglądarki do API backendu |
| React | Budowanie komponentów i interfejsu użytkownika |
| TypeScript | Typowanie kodu frontendowego |
| Node.js 20.9+ | Środowisko uruchomieniowe i narzędzia wymagane przez frontend Next.js |
| MySQL 8 / MariaDB | Przechowywanie danych aplikacji, kont użytkowników i sesji |
| XAMPP (opcjonalnie) | Może dostarczać lokalny serwer MySQL/MariaDB; Apache służy w tym przypadku do uruchomienia phpMyAdmin, a nie do obsługi aplikacji |
| phpMyAdmin (opcjonalnie) | Import schematu i lokalna administracja bazą danych |

W przeglądarce API jest dostępne pod ścieżką `/api`; Next.js przekazuje żądania do backendu Go. Nie ma potrzeby włączania CORS ani przechowywania tokena w `localStorage`. Cookie sesji ma ustawienia `HttpOnly` i `SameSite=Lax`, a dla połączeń HTTPS również `Secure`.

### 5.2. Aktualna struktura plików

Poniższe zestawienie przedstawia rzeczywistą strukturę projektu na podstawie dostarczonego archiwum ZIP. Projekt został podzielony na backend napisany w Go, frontend oparty na Next.js oraz katalogi zawierające bazę danych i dokumentację.

```text
spoldzielnia_mieszkaniowa/
├── .gitignore
├── README.md
├── DOKUMENTACJA.md
├── LICENSE
│
├── backend/
│   ├── .env.example
│   ├── go.mod
│   ├── go.sum
│   │
│   ├── cmd/
│   │   ├── server/
│   │   │   └── main.go
│   │   └── create-user/
│   │       └── main.go
│   │
│   └── internal/
│       ├── app/
│       │   ├── admin_events.go
│       │   ├── app.go
│       │   ├── app_test.go
│       │   ├── auth.go
│       │   ├── public_events.go
│       │   └── public_events_test.go
│       │
│       ├── database/
│       │   └── database.go
│       │
│       └── events/
│           ├── event.go
│           ├── input.go
│           ├── repository.go
│           ├── repository_public_test.go
│           └── repository_test.go
│
├── database/
│   ├── schema.sql
│   └── seed.sql
│
├── docs/
│   ├── MIGRACJA.md
│   ├── PLAN_PRAC.md
│   └── WERYFIKACJA.md
│
└── frontend/
    ├── .env.example
    ├── next-env.d.ts
    ├── next.config.ts
    ├── package.json
    ├── package-lock.json
    ├── tsconfig.json
    │
    ├── app/
    │   ├── [slug]/
    │   │   └── page.tsx
    │   ├── admin/
    │   │   └── wydarzenia/
    │   │       └── page.tsx
    │   ├── index.php/
    │   │   └── page.tsx
    │   ├── login/
    │   │   └── page.tsx
    │   ├── src/
    │   │   └── logout.php/
    │   │       └── page.tsx
    │   ├── wyloguj/
    │   │   └── page.tsx
    │   ├── globals.css
    │   ├── layout.tsx
    │   ├── not-found.tsx
    │   └── page.tsx
    │
    ├── components/
    │   ├── admin-events.tsx
    │   ├── auth.tsx
    │   ├── header.tsx
    │   ├── login-form.tsx
    │   └── logout-button.tsx
    │
    └── lib/
        └── routes.ts
```

Struktura nie obejmuje katalogów generowanych automatycznie podczas instalacji zależności i uruchamiania aplikacji, takich jak node_modules czy .next.

### 5.3. Najważniejsze pliki i sposób działania

W projekcie wyróżnić można kilka najważniejszych plików i katalogów odpowiedzialnych za działanie aplikacji. Backend napisany w Go odpowiada za logikę aplikacji, komunikację z bazą danych i obsługę użytkowników, natomiast frontend Next.js odpowiada za interfejs oraz interakcję z użytkownikiem.

| Plik | Obecne zastosowanie |
| --- | --- |
| `backend/cmd/server/main.go` | Główny punkt uruchomienia backendu. Odpowiada za uruchomienie serwera HTTP i konfigurację aplikacji. |
| `backend/cmd/create-user/main.go` | Narzędzie konsolowe umożliwiające tworzenie kont użytkowników z rolą mieszkańca lub administratora. |
| `backend/internal/app/app.go` | Główna konfiguracja aplikacji backendowej oraz rejestracja endpointów API. |
| `backend/internal/app/auth.go` | Obsługa logowania, wylogowania, sesji użytkowników i kontroli uprawnień. |
| `backend/internal/app/admin_events.go` | Obsługa operacji administracyjnych związanych z wydarzeniami, takich jak tworzenie, edycja, publikowanie, ukrywanie i usuwanie. |
| `backend/internal/app/public_events.go` | Obsługa pobierania opublikowanych wydarzeń dostępnych dla użytkowników. |
| `backend/internal/database/database.go` | Konfiguracja i nawiązywanie połączenia z bazą MySQL/MariaDB. |
| `backend/internal/events/repository.go` | Obsługa operacji na wydarzeniach w bazie danych. |
| `frontend/app/page.tsx` | Strona główna aplikacji zawierająca sekcję powitalną i kafelki prowadzące do poszczególnych sekcji serwisu. |
| `frontend/app/layout.tsx` | Główny układ frontendowy aplikacji, obejmujący wspólne elementy interfejsu. |
| `frontend/app/login/page.tsx` | Strona logowania użytkowników. |
| `frontend/app/admin/wydarzenia/page.tsx` | Strona panelu administracyjnego służącego do zarządzania wydarzeniami. |
| `frontend/components/header.tsx` | Komponent nagłówka i menu nawigacyjnego aplikacji. |
| `frontend/components/auth.tsx` | Obsługa stanu uwierzytelnienia użytkownika po stronie frontendu. |
| `frontend/components/login-form.tsx` | Formularz logowania komunikujący się z backendem. |
| `frontend/components/admin-events.tsx` | Komponent umożliwiający administratorowi zarządzanie wydarzeniami. |
| `frontend/lib/routes.ts` | Definicje podstron, ich tytułów, opisów oraz grup nawigacyjnych. |
| `database/schema.sql` | Definicja struktury bazy danych `mirmilowo`, obejmująca tabele użytkowników, sesji i wydarzeń kulturalnych. |
| `database/seed.sql` | Plik zawierający przykładowe wydarzenia kulturalne. |
| `frontend/next.config.ts` | Konfiguracja Next.js, w tym przekazywanie żądań API do backendu Go. |
| `frontend/package.json` | Definicja zależności frontendu oraz skryptów służących do uruchamiania, budowania i sprawdzania aplikacji. |

**Sposób działania aplikacji:**

Aplikacja działa w architekturze klient–serwer. Użytkownik korzysta z interfejsu frontendowego dostępnego pod adresem `http://localhost:3000`. Frontend wysyła żądania do API, które są przekazywane przez Next.js do backendu Go. Backend przetwarza żądania, wykonuje odpowiednie operacje na bazie MySQL/MariaDB i zwraca odpowiedzi w formacie JSON.

W przypadku logowania backend weryfikuje dane użytkownika, tworzy sesję i zapisuje ją w bazie danych. W przypadku panelu wydarzeń administrator może zarządzać wydarzeniami, a backend sprawdza jego uprawnienia przed wykonaniem operacji.

## 6. Instalacja i uruchamianie

Poniższa instrukcja dotyczy **systemu Windows i pobrania gotowego archiwum ZIP**, bez klonowania repozytorium i bez używania Git. Oddzielnie opisano uruchomienie obecnej strony oraz wymagania do działania logowania.

### 6.1. Wymagania systemowe

Przed rozpoczęciem instalacji należy przygotować następujące oprogramowanie:

1. **Go 1.26 lub nowszy** – środowisko wymagane do uruchomienia backendu aplikacji. Instalator dostępny jest na stronie https://go.dev/dl/.
2. **Node.js 20.9 lub nowszy** – wymagany do uruchomienia frontendu. Zalecana jest aktualna wersja LTS. Instalator dostępny jest na stronie https://nodejs.org/.
3. **MySQL 8 lub MariaDB** – system zarządzania bazą danych wykorzystywaną przez aplikację.
4. **XAMPP** – może zostać wykorzystany do uruchomienia MySQL/MariaDB oraz phpMyAdmin.
5. **Przeglądarka internetowa** – np. Google Chrome, Mozilla Firefox lub Microsoft Edge.
6. **Program do rozpakowywania plików ZIP** – w systemie Windows można skorzystać z wbudowanej funkcji „Wyodrębnij wszystkie”.

Do uruchomienia aplikacji nie jest wymagany Git, GitHub Desktop ani Composer.

### 6.2. Przygotowanie bazy danych

Aplikacja wykorzystuje bazę danych MySQL/MariaDB o nazwie `mirmilowo`.

W przypadku korzystania z XAMPP należy wykonać następujące czynności:

1. Uruchomić XAMPP Control Panel.
2. Włączyć usługi MySQL oraz Apache.
3. Otworzyć przeglądarkę internetową i przejść pod adres:

   `http://localhost/phpmyadmin`

4. Wybrać zakładkę **Import** na poziomie serwera.
5. Wskazać plik `database/schema.sql` znajdujący się w katalogu projektu.
6. Wykonać import pliku.

Plik `schema.sql` automatycznie tworzy bazę danych `mirmilowo` oraz wymagane tabele:

- `users` – przechowuje konta użytkowników;
- `sessions` – przechowuje informacje o sesjach użytkowników;
- `cultural_events` – przechowuje dane dotyczące wydarzeń kulturalnych.

Alternatywnie można wykonać import za pomocą terminala klienta MySQL:

```sql
SOURCE /pelna/sciezka/do/mirmilowo/database/schema.sql;
```

Po zaimportowaniu struktury bazy należy zaimportować również plik `database/seed.sql`. Zawiera on trzy przykładowe wydarzenia: opublikowane przyszłe, opublikowane archiwalne oraz ukryte.

Dane przyszłych wydarzeń zostały ustalone na 21 listopada 2026 r. Po upływie tego terminu wydarzenia naturalnie trafią do archiwum.

W przypadku istniejącej bazy danych zawierającej już tabele kont użytkowników i sesji można zamiast pełnego schematu wykonać migrację:

`database/migrations/001_cultural_events.sql`

Migrację należy wykonać po wybraniu właściwej bazy danych. Dodaje ona brakującą tabelę wydarzeń i nie zmienia struktury już istniejących tabel.

Plik `seed.sql` nie tworzy kont użytkowników. Konta administratora i mieszkańca zostaną utworzone w kolejnym etapie instalacji.

Ponowny import schematu nie usuwa istniejących danych. Należy jednak pamiętać, że instrukcje `IF NOT EXISTS` nie aktualizują struktury już istniejących tabel. W przypadku przyszłych zmian struktury bazy należy przygotowywać osobne migracje SQL.

### 6.3. Konfiguracja i uruchomienie backendu

Backend aplikacji został napisany w języku Go i odpowiada za obsługę komunikacji z bazą danych, logowania, sesji użytkowników oraz funkcjonalności dostępnych przez API.

W celu uruchomienia backendu należy:

1. Otworzyć folder główny projektu.
2. Uruchomić PowerShell.
3. Przejść do katalogu `backend`:

```powershell
cd backend
```

4. Ustawić zmienne środowiskowe odpowiedzialne za połączenie z bazą danych:

```powershell
$env:DB_ADDR = "127.0.0.1:3306"
$env:DB_NAME = "mirmilowo"
$env:DB_USER = "root"
$env:DB_PASSWORD = ""
$env:APP_ORIGIN = "http://localhost:3000"
```

Powyższe wartości odpowiadają typowej lokalnej konfiguracji XAMPP. Jeżeli konto `root` posiada hasło lub serwer MySQL/MariaDB działa na innym porcie, należy odpowiednio zmienić wartości zmiennych.

Dostępne ustawienia opisuje plik `backend/.env.example`. Backend odczytuje zmienne ze środowiska procesu i nie wczytuje automatycznie pliku `.env`.

5. Pobrać wymagane zależności:

```powershell
go mod download
```

6. Uruchomić serwer backendowy:

```powershell
go run ./cmd/server
```

Po poprawnym uruchomieniu backend domyślnie nasłuchuje lokalnie pod adresem:

`http://127.0.0.1:8080`

Terminal, w którym uruchomiono backend, musi pozostać otwarty podczas korzystania z aplikacji.

W systemach macOS/Linux odpowiednikiem konfiguracji zmiennych środowiskowych jest:

```bash
export DB_ADDR=127.0.0.1:3306 DB_NAME=mirmilowo DB_USER=root DB_PASSWORD=''
export APP_ORIGIN=http://localhost:3000
go mod download
go run ./cmd/server
```

### 6.4. Tworzenie kont użytkowników

Po uruchomieniu backendu należy utworzyć konta użytkowników, które umożliwią logowanie do aplikacji.

Projekt udostępnia narzędzie konsolowe `create-user`, pozwalające tworzyć konta z odpowiednimi rolami.

W tym celu należy:

1. Otworzyć drugi terminal PowerShell. Nie należy zamykać pierwszego terminala, w którym działa backend.
2. Przejść do katalogu `backend`.
3. Ustawić ponownie te same zmienne środowiskowe połączenia z bazą danych, które zostały wykorzystane podczas uruchamiania backendu:

```powershell
$env:DB_ADDR = "127.0.0.1:3306"
$env:DB_NAME = "mirmilowo"
$env:DB_USER = "root"
$env:DB_PASSWORD = ""
$env:APP_ORIGIN = "http://localhost:3000"
```

Zmienne środowiskowe PowerShell obowiązują wyłącznie w terminalu, w którym zostały ustawione.

4. Utworzyć konto administratora:

```powershell
go run ./cmd/create-user -email admin@example.test -role admin
```

5. Utworzyć konto mieszkańca:

```powershell
go run ./cmd/create-user -email mieszkaniec@example.test -role mieszkaniec
```

Każde polecenie poprosi o podanie hasła, które nie będzie widoczne podczas wpisywania.

Hasło powinno mieć długość od 12 do 72 bajtów. Polskie znaki mogą zajmować więcej niż jeden bajt.

W bazie danych przechowywany jest wyłącznie skrót hasła wygenerowany przy użyciu algorytmu bcrypt.

Ponowne wykonanie polecenia dla tego samego adresu e-mail spowoduje wyświetlenie błędu zamiast nadpisania istniejącego konta.

Podane adresy e-mail są przykładowe i mogą zostać zastąpione własnymi danymi.

### 6.5. Uruchomienie frontendu Next.js

Frontend został przygotowany z wykorzystaniem Next.js, React oraz TypeScript. Odpowiada za wyświetlanie interfejsu użytkownika i komunikację z backendem.

W celu uruchomienia frontendu należy:

1. Otworzyć kolejne okno PowerShell.
2. Przejść do głównego katalogu projektu, a następnie do folderu `frontend`:

```powershell
cd frontend
```

3. Zainstalować wymagane zależności:

```powershell
npm ci
```

4. Uruchomić aplikację w trybie deweloperskim:

```powershell
npm run dev
```

W przypadku problemów z wykonywaniem skryptu `npm.ps1` w PowerShell można wykorzystać alternatywne polecenia:

```powershell
npm.cmd ci
npm.cmd run dev
```

Po poprawnym uruchomieniu frontendu należy otworzyć przeglądarkę i wpisać adres:

`http://localhost:3000`

Formularz logowania znajduje się pod adresem:

`http://localhost:3000/login`

Podczas korzystania z aplikacji backend Go oraz serwer MySQL/MariaDB muszą pozostawać uruchomione.

Domyślna konfiguracja działa bez konieczności tworzenia pliku `.env.local`.

W przypadku zmiany adresu backendu należy skopiować plik `frontend/.env.example` do `frontend/.env.local` i ustawić odpowiednią wartość `API_INTERNAL_URL`.

Po zmianie konfiguracji należy ponownie uruchomić Next.js. W przypadku wersji produkcyjnej wymagane jest ponowne wykonanie kompilacji.

Należy używać dokładnie tego samego hosta, który został ustawiony w zmiennej `APP_ORIGIN`. Adresy `localhost` i `127.0.0.1` są traktowane jako różne originy.

Jeżeli Next.js zaproponuje uruchomienie aplikacji na porcie 3001 z powodu zajętego portu 3000, należy zwolnić port 3000 albo odpowiednio zmienić wartość `APP_ORIGIN` w backendzie i ponownie go uruchomić.

### 6.6. Testowanie i kompilacja aplikacji

Po uruchomieniu wszystkich elementów systemu można sprawdzić poprawność działania backendu i frontendu.

**Testowanie backendu**

W katalogu `backend` należy wykonać następujące polecenia:

```powershell
go test ./...
go vet ./...
go build ./cmd/server ./cmd/create-user
```

Poszczególne polecenia odpowiadają za:

- `go test ./...` – uruchomienie testów jednostkowych backendu;
- `go vet ./...` – sprawdzenie potencjalnych błędów i nieprawidłowości w kodzie;
- `go build ./cmd/server ./cmd/create-user` – kompilację programów backendowych.

Testy Go wykorzystują atrapę SQL i nie wymagają uruchomionej bazy danych.

Sprawdzają między innymi:

- proces logowania;
- rotację i usuwanie sesji;
- obsługę nieprawidłowych haseł;
- kontrolę origin i zabezpieczenia CSRF;
- kontrolę ról użytkowników;
- obsługę awarii SQL.

**Testowanie i kompilacja frontendu**

W katalogu `frontend` należy wykonać:

```powershell
npm run typecheck
npm run build
```

Polecenie `npm run typecheck` sprawdza poprawność typów TypeScript, natomiast `npm run build` przygotowuje produkcyjną wersję aplikacji.

Uruchomienie wcześniej zbudowanej wersji produkcyjnej odbywa się za pomocą polecenia:

```powershell
npm start
```

Polecenia `npm start` nie należy uruchamiać równocześnie z `npm run dev` na tym samym porcie.

**Testowanie ręczne z wykorzystaniem rzeczywistej bazy danych**

Po uruchomieniu MySQL/MariaDB, backendu oraz frontendu należy wykonać podstawowe testy:

1. Zalogować się na utworzone konto administratora.
2. Zalogować się na utworzone konto mieszkańca.
3. Sprawdzić komunikat pojawiający się po podaniu nieprawidłowego hasła.
4. Odświeżyć stronę po zalogowaniu i sprawdzić, czy sesja nadal jest aktywna.
5. Zrestartować backend i ponownie odświeżyć stronę.
6. Wylogować użytkownika.
7. Sprawdzić, czy endpoint `/api/auth/me` zwraca odpowiedź:

```json
{"user":null}
```

8. Otworzyć menu aplikacji na urządzeniu mobilnym.
9. Sprawdzić zachowanie aplikacji po wejściu na nieistniejący adres.
