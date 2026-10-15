
# Spółdzielnia Mieszkaniowa „Mirmiłowo”

**Repozytorium:** `spoldzielnia_mieszkaniowa`  
**Wersja dokumentacji:** 1.0  
**Typ projektu:** Aplikacja internetowa — projekt akademicki  
**Technologie:** PHP, MySQL/MariaDB, Bootstrap  
**Środowisko:** Localhost — XAMPP  
**Zespół:** 12 osób  
**Stan realizacji:** Ukończone etapy 1 i 2

---

## 1. Opis projektu

Spółdzielnia Mieszkaniowa „Mirmiłowo” jest projektem aplikacji internetowej obsługującej jedno fikcyjne osiedle mieszkaniowe.

Celem systemu jest ułatwienie dostępu do informacji, dokumentów, wydarzeń i usług spółdzielni oraz umożliwienie mieszkańcom załatwiania wybranych spraw elektronicznie.

Projekt będzie uruchamiany lokalnie. Nie przewiduje się wdrożenia produkcyjnego ani przetwarzania rzeczywistych danych mieszkańców.

---

## 2. Aktualny stan realizacji

Zakończono dwa pierwsze etapy projektu:
- Etap 1 - fundament projektu: struktura katalogów, szkielet aplikacji PHP, interfejs Bootstrap, nawigacja i dokumentacja.
- Etap 2 - baza i użytkownicy: baza MySQL/MariaDB, połączenie PDO, logowanie, wylogowanie oraz obsługa ról mieszkańca i administratora.
Pozostałe moduły będą dołączane do istniejącej aplikacji. Obecność pozycji w menu nie oznacza ukończenia danej funkcjonalności. Niezaimplementowane moduły powinny prezentować strony informacyjne.

| Etap | Zakres | Status |
| --- | --- | --- |
| 1 | Fundament projektu | Ukończony |
| 2 | Baza danych i użytkownicy | Ukończony |
| 3 | Struktura i historia spółdzielni | Planowany |
| 4 | Imprezy kulturalne | Planowany |
| 5 | Przetargi | Planowany |
| 6 | Wynajem | Planowany |
| 7 | Kalendarz | Planowany |
| 8 | Inwestycje | Planowany |
| 9 | Dokumenty | Planowany |
| 10 | Czat mieszkańców | Planowany |
| 11 | E-mieszkaniec | Planowany |
| 12 | Kontakt i wydanie końcowe | Planowany |

---

## 3. Funkcjonalności

### 3.1. Fundament i użytkownicy - wykonane
- Szkielet aplikacji PHP.
- Interfejs wykorzystujący Bootstrap.
- Nawigacja między podstronami.
- Baza danych MySQL/MariaDB.
- Połączenie z bazą przez PDO.
- Logowanie i wylogowanie użytkowników.
- Rozróżnianie roli mieszkańca i administratora.
- Mechanizm kontroli dostępu oparty na rolach.

### 3.2. Ogłoszenia

- Imprezy kulturalne.
- Przetargi.
- Wynajem.
- Kalendarz wydarzeń i terminów.

Kalendarz będzie stosował oznaczenia kolorystyczne oraz legendę kategorii.

### 3.3. Inwestycje

- Lista inwestycji.
- Opisy i szczegóły.
- Statusy oraz postęp prac.

### 3.4. Strefa mieszkańca

- Druki i wnioski.
- Statut i regulaminy.
- Czat mieszkańców.

### 3.5. E-mieszkaniec

- Wypełnianie formularzy online.
- Przesyłanie wniosków.
- Dodawanie załączników.
- Sprawdzanie statusu zgłoszenia.
- Demonstracyjny proces podpisywania dokumentów.

### 3.6. Kontakt

- Dane kontaktowe spółdzielni.
- Formularz wiadomości.

---

## 4. Role użytkowników

**Gość** — przegląda publiczne informacje.

**Mieszkaniec** — po zalogowaniu korzysta z czatu i elektronicznych wniosków.

**Administrator** — zarządza treściami, moderuje czat i obsługuje zgłoszenia.

---

## 5. Instalacja i uruchamianie

### 5.1. Wymagania

- Git
- XAMPP z Apache i PHP
- MySQL lub MariaDB — od etapu 2
- Przeglądarka internetowa

### 5.2. Pobranie projektu

Otwórz terminal i wpisz:

```bash
cd C:\xampp\htdocs
git clone https://github.com/NAZWA_KONTA/spoldzielnia_mieszkaniowa.git
```

> **Uwaga:** `NAZWA_KONTA` należy zastąpić nazwą użytkownika lub organizacji GitHub, w której znajduje się repozytorium.

Następnie uruchom Apache w XAMPP.

### 5.3. Otwarcie aplikacji

W przeglądarce otwórz:

`http://localhost/spoldzielnia_mieszkaniowa/`

Po pierwszym etapie aplikacja działa bez bazy danych.

### 5.4. Konfiguracja bazy danych

Od etapu 2:

1. Uruchom Apache i MySQL w XAMPP.
2. Otwórz phpMyAdmin.
3. Utwórz bazę `mirmilowo`.
4. Zaimportuj `database/schema.sql`.
5. Zaimportuj `database/seed.sql`.
6. Skopiuj przykładową konfigurację bazy do pliku lokalnego.
7. Wprowadź parametry połączenia.
8. Utwórz konto administracyjne zgodnie z instrukcją dołączoną w etapie 2.

> **Ważne:** Lokalny plik konfiguracyjny z hasłem nie powinien trafiać do GitHub.

---

## 6. Struktura projektu

```text
spoldzielnia_mieszkaniowa/
│
├── index.php
├── assets/
├── config/
├── includes/
├── pages/
├── src/
├── database/
├── storage/
└── docs/
```

Opis katalogów i plików:

| Katalog / plik | Przeznaczenie |
|---|---|
| `index.php` | Punkt wejścia aplikacji |
| `assets/` | Style i pliki interfejsu |
| `config/` | Konfiguracja i trasy |
| `includes/` | Wspólne elementy strony |
| `pages/` | Widoki podstron |
| `src/` | Logika aplikacji |
| `database/` | Struktura i dane bazy |
| `storage/` | Przechowywane pliki |
| `docs/` | Dokumentacja techniczna |

---

## 7. Organizacja pracy

Projekt podzielono na **12 etapów odpowiadających 12 osobom**.

Współpraca jest realizowana przez repozytorium GitHub.

Każdy uczestnik:

1. Wykonuje przydzielone zadanie na osobnej gałęzi (*branch*).
2. Testuje wprowadzone zmiany.
3. Zgłasza *Pull Request* do gałęzi `main`.

Połączenie zmian następuje po ich sprawdzeniu przez innego członka zespołu.

---

## 8. Bezpieczeństwo

Projekt powinien stosować:

- Weryfikację uprawnień po stronie PHP.
- Bezpieczne przechowywanie haseł.
- Parametryzowane zapytania SQL.
- Walidację formularzy.
- Ochronę przed CSRF i XSS.
- Kontrolowane przesyłanie i udostępnianie plików.
- Wyłącznie fikcyjne dane testowe.

**Uwaga dotycząca podpisu elektronicznego:**

Funkcja podpisu Profilem Zaufanym pozostaje demonstracją. Nie potwierdza rzeczywistego złożenia ani zweryfikowania podpisu elektronicznego.

---

## 9. Dokumentowanie i testowanie modułów

Każdy autor powinien opisać:
- zakres wykonanej funkcji;
- dodane i zmienione pliki;
- zmiany bazy danych;
- reguły walidacji;
- zastosowane zabezpieczenia;
- dodatkową konfigurację;
- instrukcję użycia;
- wyniki testów;
- znane ograniczenia.
Testy powinny obejmować:
- poprawne wykonanie operacji;
- błędne i niekompletne dane;
- brak wymaganych uprawnień;
- nieistniejący rekord lub pustą listę;
- działanie na urządzeniach mobilnych;
- wpływ zmian na istniejące funkcje.
Wyniki należy zapisywać wraz z datą, osobą testującą i wersją kodu

---

## 10. Warunki ukończenia

Projekt zostanie uznany za ukończony, kiedy:

- [ ] Można go pobrać i uruchomić zgodnie z dokumentacją.
- [ ] Wszystkie główne funkcje działają na localhost.
- [ ] Baza danych zawiera komplet wymaganych tabel.
- [ ] Każda rola ma właściwe uprawnienia.
- [ ] Formularze zapisują i odczytują dane.
- [ ] Kalendarz prawidłowo prezentuje terminy.
- [ ] Aplikacja poprawnie działa na urządzeniach mobilnych.
- [ ] Wszystkie funkcjonalności przejdą testy końcowe.

---

**Status projektu:** rozwój aplikacji po ukończeniu etapów 1 i 2.
**Status dokumentacji:** wersja rozwojowa - wymaga uzupełnienia szczegółów konfiguracji oraz aktualizacji podczas realizacji kolejnych etapów.
