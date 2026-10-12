# Spółdzielnia Mieszkaniowa „Mirmiłowo”

Aplikacja webowa dla fikcyjnej spółdzielni mieszkaniowej zarządzającej jednym osiedlem — Mirmiłowo. Docelowo umożliwi mieszkańcom przeglądanie informacji o spółdzielni, ogłoszeń i inwestycji, korzystanie z kalendarza wydarzeń, składanie wniosków oraz komunikację ze spółdzielnią. Administrator będzie mógł zarządzać treściami i zgłoszeniami.

**Stan projektu:** po etapie 1 dostępny jest szkielet strony z działającą nawigacją i podstronami zastępczymi. Funkcje bazodanowe, logowanie i formularze będą realizowane w kolejnych etapach.

## Technologie

- PHP 8.1 lub nowszy,
- MySQL 8 / MariaDB,
- HTML5,
- CSS3,
- JavaScript,
- Bootstrap 5.3.8,
- Apache,
- XAMPP,
- Git i GitHub.

## Główne funkcjonalności

*Poniższa lista przedstawia zakres docelowy projektu.*

- Responsywna strona główna i nawigacja po serwisie,
- Sekcja **Spółdzielnia** — struktura i historia,
- Sekcja **Ogłoszenia** — imprezy kulturalne, przetargi i oferty wynajmu,
- Kalendarz wydarzeń, terminów składania wniosków i spraw pilnych, z oznaczeniami kategorii kolorami,
- Sekcja **Inwestycje** — informacje, postęp i status realizowanych inwestycji,
- **Strefa mieszkańca** — druki i wnioski, statut i regulaminy oraz czat mieszkańców,
- Rejestracja, logowanie i wylogowanie użytkowników,
- Role gościa, mieszkańca i administratora,
- **E-mieszkaniec** — elektroniczne składanie wniosków, dodawanie załączników i śledzenie statusu zgłoszeń,
- Demonstracja procesu podpisywania dokumentów Profilem Zaufanym (bez rzeczywistej integracji z usługą państwową),
- Sekcja **Kontakt** — dane spółdzielni i formularz kontaktowy,
- Panel administratora do zarządzania treściami, moderowania czatu i obsługi zgłoszeń.

## Uruchomienie

1. Pobierz lub sklonuj repozytorium projektu.

2. Umieść projekt w katalogu `htdocs` XAMPP.

   Katalog projektu powinien nazywać się:

   ```text
   spoldzielnia_mieszkaniowa
   ```

   Przykładowa lokalizacja:

   ```text
   C:\xampp\htdocs\spoldzielnia_mieszkaniowa
   ```

3. Uruchom **Apache** w panelu XAMPP.

4. Otwórz aplikację w przeglądarce:

   ```text
   http://localhost/spoldzielnia_mieszkaniowa/
   ```

5. Sprawdź działanie strony głównej i menu. Na obecnym etapie część podstron wyświetla komunikat „Funkcjonalność w przygotowaniu”.

**Uwaga:** w etapie 1 nie trzeba uruchamiać MySQL, importować bazy ani konfigurować logowania. Bootstrap jest ładowany z CDN, więc do poprawnego wyświetlania stylów potrzebny jest internet.

### Konfiguracja bazy danych — od etapu 2

*Poniższe kroki dotyczą planowanej wersji z obsługą bazy danych. Pliki i konfiguracja muszą zostać dodane i zweryfikowane w etapie 2.*

1. Uruchom **Apache** i **MySQL** w XAMPP.

2. Otwórz phpMyAdmin:

   ```text
   http://localhost/phpmyadmin
   ```

3. Utwórz bazę danych `mirmilowo` z kodowaniem `utf8mb4`.

4. Zaimportuj pliki bazy danych w podanej kolejności (po ich przygotowaniu):

   ```text
   database/schema.sql
   database/seed.sql
   ```

5. Skopiuj plik `config/database.example.php` i zapisz jego kopię jako `config/database.local.php` (po dodaniu pliku przykładowego).

6. Uzupełnij lokalne dane połączenia z bazą w `config/database.local.php`. Nie publikuj tego pliku ani haseł w repozytorium.
