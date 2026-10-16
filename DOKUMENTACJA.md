# Spółdzielnia Mieszkaniowa „Mirmiłowo” — dokumentacja projektu

**Nazwa projektu:** Spółdzielnia Mieszkaniowa „Mirmiłowo”  

**Repozytorium:** `spoldzielnia_mieszkaniowa`  

**Rodzaj:** Aplikacja internetowa dla współdzielni mieszkaniowej 

**Środowisko docelowe:** lokalne (XAMPP: Apache, PHP, MySQL/MariaDB)  

**Technologie:** PHP 8.1+, HTML5, CSS3, JavaScript, Bootstrap 5.3.8, MySQL/MariaDB

**Zespół:** 12 osób  

**Stan projektu:** etap 1 wykonany, etap 2 rozpoczęty

---

## 1. Opis projektu

Spółdzielnia Mieszkaniowa „Mirmiłowo” to projekt strony internetowej obsługującej **jedno fikcyjne osiedle mieszkaniowe**. Ma stanowić wspólne miejsce publikowania informacji o spółdzielni, ogłoszeń, dokumentów i terminów oraz umożliwiać mieszkańcom kontakt ze spółdzielnią i elektroniczną obsługę wybranych spraw.

Docelowo aplikacja będzie przeznaczona dla trzech grup użytkowników: gości, zalogowanych mieszkańców i administratorów. Zakres dostępnych funkcji będzie zależeć od roli użytkownika.

Projekt ma charakter dydaktyczny. Zakłada uruchamianie na własnym komputerze za pomocą XAMPP, bez wdrożenia produkcyjnego i bez wykorzystywania rzeczywistych danych mieszkańców. Planowany podpis za pomocą Profilu Zaufanego będzie **wyłącznie symulacją**, a nie integracją z usługą państwową.

## 2. Aktualny stan realizacji

### 2.1. Postęp według etapów

| Etap | Zakres | Stan na podstawie aktualnego kodu |
| --- | --- | --- |
| 1 | Repozytorium, struktura projektu, szkielet strony, nawigacja i README | **Wykonany** — elementy są obecne w projekcie |
| 2 | Baza danych, konta użytkowników, logowanie, wylogowanie, role i uprawnienia | **W trakcie** — część kodu istnieje, lecz etap nie spełnia wszystkich założeń |
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

- strona główna z opisem serwisu, przyciskami prowadzącymi do kalendarza i E-mieszkańca oraz kafelkami najważniejszych sekcji;
- responsywny układ zbudowany z użyciem Bootstrap 5.3.8 i własnego arkusza CSS;
- wspólny nagłówek, menu nawigacyjne i stopka;
- lista dozwolonych adresów podstron w `config/routes.php`;
- podstrony informacyjne („Etap w przygotowaniu”) dla modułów, których jeszcze nie wdrożono;
- obsługa nieznanej podstrony — komunikat 404;
- formularz logowania z polami adresu e-mail i hasła;
- kod łączący się z MySQL/MariaDB przez PDO oraz zapytanie odczytujące użytkownika według adresu e-mail;
- sprawdzanie hasła za pomocą `password_verify()` i zapis identyfikatora oraz roli użytkownika w sesji po pomyślnym logowaniu;
- prosty skrypt wylogowania i podstawowe funkcje sprawdzające stan sesji oraz rolę;
- katalogi przewidziane pod przyszłe moduły, dokumenty i pliki aplikacji.

**Ograniczenie:** istnienie kodu logowania nie oznacza, że da się od razu zalogować po wypakowaniu ZIP-a. Projekt nie dostarcza jeszcze kompletnej struktury bazy ani kont testowych.

## 3. Funkcjonalności aplikacji — obecnie i docelowo

### 3.1. Strona główna i nawigacja

**Obecnie:** działa szkielet serwisu. Użytkownik może przechodzić pomiędzy zdefiniowanymi podstronami, a większość z nich wyświetla informację o planowanej funkcjonalności. Dostępna jest również wersja menu na mniejszych ekranach.

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

**Obecnie:** dostępny jest formularz logowania, podstawowa obsługa sesji i wylogowanie. Brak gotowych kont, paneli użytkowników i kompletnej ochrony poszczególnych funkcji.

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
| PHP 8.1+ | Generowanie stron, obsługa formularzy i logika aplikacji |
| MySQL / MariaDB | Planowane przechowywanie kont oraz danych kolejnych modułów |
| PDO (`pdo_mysql`) | Połączenie PHP z bazą i wykonywanie zapytań SQL |
| HTML5, CSS3, JavaScript | Struktura, wygląd i zachowanie interfejsu |
| Bootstrap 5.3.8 | Układ responsywny, menu i komponenty interfejsu |
| Apache (XAMPP) | Lokalna obsługa aplikacji PHP |
| phpMyAdmin | Lokalna administracja bazą danych |

Bootstrap jest pobierany z CDN, dlatego poprawne wyświetlanie jego stylów i działanie menu rozwijanego wymagają dostępu do internetu lub późniejszego udostępnienia plików biblioteki lokalnie.

### 5.2. Aktualna struktura plików

Poniższe zestawienie opisuje **rzeczywistą zawartość dostarczonego ZIP-a**, a nie strukturę planowaną:

```text
spoldzielnia_mieszkaniowa/
├── index.php
├── .gitignore
├── .htaccess
├── README.md
├── DOKUMENTACJA.md
├── LICENSE
├── assets/
│   └── css/
│       └── style.css
├── config/
│   ├── .htaccess
│   ├── database.php
│   └── routes.php
├── database/
│   ├── .htaccess
│   ├── schema.sql
│   └── seed.sql
├── docs/
│   └── PLAN_PRAC.md
├── includes/
│   ├── .htaccess
│   ├── auth.php
│   ├── header.php
│   └── footer.php
├── pages/
│   ├── .htaccess
│   ├── home.php
│   ├── login.php
│   └── placeholder.php
├── src/
│   ├── .gitkeep
│   └── logout.php
└── storage/
    ├── .htaccess
    └── uploads/
        └── .gitkeep
```

Archiwum zawiera też katalog `.git/` z lokalnymi metadanymi repozytorium. Nie jest on wymagany do uruchomienia aplikacji z ZIP-a.

### 5.3. Najważniejsze pliki i sposób działania

| Plik | Obecne zastosowanie |
| --- | --- |
| `index.php` | Główny punkt wejścia; wybiera podstronę z parametru `page` i ładuje widok |
| `config/routes.php` | Lista dozwolonych podstron oraz ich tytułów i opisów |
| `pages/home.php` | Treść strony głównej |
| `pages/placeholder.php` | Widoki niewdrożonych modułów i komunikat 404 |
| `pages/login.php` | Formularz logowania oraz próba pobrania użytkownika z bazy |
| `includes/auth.php` | Rozpoczęcie sesji i podstawowe funkcje rozpoznające zalogowanie i rolę |
| `src/logout.php` | Wywołanie zakończenia sesji i przekierowanie do strony głównej |
| `config/database.php` | Bezpośrednia konfiguracja i nawiązanie połączenia PDO |
| `includes/header.php`, `includes/footer.php` | Wspólna nawigacja i stopka |
| `assets/css/style.css` | Dodatkowe style aplikacji |
| `database/schema.sql`, `database/seed.sql` | Pliki przygotowane pod bazę, obecnie bez wykonywalnych definicji SQL |
| `docs/PLAN_PRAC.md` | Szczegółowy plan projektu |

**Przykładowy adres podstrony:** `http://localhost/spoldzielnia_mieszkaniowa/index.php?page=imprezy`.

**Przepływ w aplikacji:** przeglądarka otwiera `index.php` → PHP sprawdza parametr `page` w `config/routes.php` → wczytuje nagłówek → wyświetla stronę główną, istniejący widok albo stronę zastępczą → dołącza stopkę.

**Przepływ logowania (w kodzie):** użytkownik wysyła e-mail i hasło → `pages/login.php` inicjuje połączenie PDO → PHP próbuje pobrać rekord z `users` → `password_verify()` porównuje hasło z hashem → w przypadku poprawnych danych kod ustawia zmienne sesji i przekierowuje na stronę główną. Ten proces jest zależny od ręcznie przygotowanej bazy i kont, których obecne archiwum nie dostarcza.

---

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
