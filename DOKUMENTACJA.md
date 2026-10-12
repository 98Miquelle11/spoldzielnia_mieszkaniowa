
# Spółdzielnia Mieszkaniowa „Mirmiłowo”

**Repozytorium:** `spoldzielnia_mieszkaniowa`  
**Wersja dokumentacji:** 1.0  
**Typ projektu:** Aplikacja internetowa — projekt akademicki  
**Technologie:** PHP, MySQL/MariaDB, Bootstrap  
**Środowisko:** Localhost — XAMPP  
**Zespół:** 12 osób  

---

## 1. Opis projektu

Spółdzielnia Mieszkaniowa „Mirmiłowo” jest projektem aplikacji internetowej obsługującej jedno fikcyjne osiedle mieszkaniowe.

Celem systemu jest ułatwienie dostępu do informacji, dokumentów, wydarzeń i usług spółdzielni oraz umożliwienie mieszkańcom załatwiania wybranych spraw elektronicznie.

Projekt będzie uruchamiany lokalnie. Nie przewiduje się wdrożenia produkcyjnego ani przetwarzania rzeczywistych danych mieszkańców.

---

## 2. Funkcjonalności

### 2.1. Spółdzielnia

- Struktura organizacyjna spółdzielni.
- Historia osiedla.

### 2.2. Ogłoszenia

- Imprezy kulturalne.
- Przetargi.
- Wynajem.
- Kalendarz wydarzeń i terminów.

Kalendarz będzie stosował oznaczenia kolorystyczne oraz legendę kategorii.

### 2.3. Inwestycje

- Lista inwestycji.
- Opisy i szczegóły.
- Statusy oraz postęp prac.

### 2.4. Strefa mieszkańca

- Druki i wnioski.
- Statut i regulaminy.
- Czat mieszkańców.

### 2.5. E-mieszkaniec

- Wypełnianie formularzy online.
- Przesyłanie wniosków.
- Dodawanie załączników.
- Sprawdzanie statusu zgłoszenia.
- Demonstracyjny proces podpisywania dokumentów.

### 2.6. Kontakt

- Dane kontaktowe spółdzielni.
- Formularz wiadomości.

---

## 3. Role użytkowników

**Gość** — przegląda publiczne informacje.

**Mieszkaniec** — po zalogowaniu korzysta z czatu i elektronicznych wniosków.

**Administrator** — zarządza treściami, moderuje czat i obsługuje zgłoszenia.

---

## 4. Instalacja i uruchamianie

### 4.1. Wymagania

- Git
- XAMPP z Apache i PHP
- MySQL lub MariaDB — od etapu 2
- Przeglądarka internetowa

### 4.2. Pobranie projektu

Otwórz terminal i wpisz:

```bash
cd C:\xampp\htdocs
git clone https://github.com/NAZWA_KONTA/spoldzielnia_mieszkaniowa.git
```

> **Uwaga:** `NAZWA_KONTA` należy zastąpić nazwą użytkownika lub organizacji GitHub, w której znajduje się repozytorium.

Następnie uruchom Apache w XAMPP.

### 4.3. Otwarcie aplikacji

W przeglądarce otwórz:

`http://localhost/spoldzielnia_mieszkaniowa/`

Po pierwszym etapie aplikacja działa bez bazy danych.

### 4.4. Konfiguracja bazy danych

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

## 5. Struktura projektu

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

## 6. Organizacja pracy

Projekt podzielono na **12 etapów odpowiadających 12 osobom**.

Współpraca jest realizowana przez repozytorium GitHub.

Każdy uczestnik:

1. Wykonuje przydzielone zadanie na osobnej gałęzi (*branch*).
2. Testuje wprowadzone zmiany.
3. Zgłasza *Pull Request* do gałęzi `main`.

Połączenie zmian następuje po ich sprawdzeniu przez innego członka zespołu.

---

## 7. Bezpieczeństwo

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

## 8. Warunki ukończenia

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

**Status projektu:** Dokumentacja wstępna — podlega aktualizacji wraz z rozwojem aplikacji.
