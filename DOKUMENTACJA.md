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
| 4 | Ogłoszenia: Imprezy kulturalne | Planowany |
| 5 | Ogłoszenia: Przetargi | Planowany |
| 6 | Ogłoszenia: Wynajem | Planowany |
| 7 | Kalendarz wydarzeń i terminów | Planowany |
| 8 | Inwestycje | Planowany |
| 9 | Strefa mieszkańca: Druki i wnioski, Statut i regulaminy | Planowany |
| 10 | Czat mieszkańców | Planowany |
| 11 | E-mieszkaniec | Planowany |
| 12 | Kontakt, integracja oraz testy końcowe | Planowany |

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

**Obecnie:** menu udostępnia odnośniki do „Imprez kulturalnych”, „Przetargów”, „Wynajmu” i „Kalendarza”. Żaden z tych modułów nie zawiera jeszcze docelowej obsługi danych.

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

### 6.1. Co trzeba zainstalować na komputerze

1. **XAMPP** zawierający serwer Apache, PHP 8.1 lub nowszy oraz MySQL/MariaDB (wraz z phpMyAdmin).
2. **Przeglądarkę internetową**, np. Chrome, Firefox lub Edge.
3. **Program do rozpakowywania plików ZIP** — w Windows wystarczy wbudowana funkcja „Wyodrębnij wszystkie”.
4. **Dostęp do internetu**, jeśli Bootstrap ma zostać pobrany z CDN.

Do samego uruchomienia pobranej aplikacji **nie potrzeba instalować Git ani GitHub Desktop**. Nie jest również wymagany Composer ani Node.js.

### 6.2. Pobranie i rozpakowanie aplikacji ZIP

1. Pobierz archiwum ZIP z aktualną wersją projektu.
2. Kliknij pobrany plik prawym przyciskiem myszy i wybierz **Wyodrębnij wszystkie**.
3. Otwórz wypakowany folder i znajdź katalog zawierający bezpośrednio plik `index.php` oraz katalogi `assets`, `config`, `includes` i `pages`.
4. Skopiuj **ten właśnie katalog** do folderu `C:\xampp\htdocs\`.
5. Nadaj mu nazwę `spoldzielnia_mieszkaniowa`.

Docelowa lokalizacja pliku startowego powinna być następująca:

```text
C:\xampp\htdocs\spoldzielnia_mieszkaniowa\index.php
```

### 6.3. Uruchomienie strony głównej

1. Otwórz **XAMPP Control Panel**.
2. Kliknij **Start** przy **Apache** i sprawdź, czy serwer się uruchomił.
3. W przeglądarce wpisz:

   ```text
   http://localhost/spoldzielnia_mieszkaniowa/
   ```
