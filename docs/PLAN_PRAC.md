# PLAN PRAC — Spółdzielnia Mieszkaniowa „Mirmiłowo”

**Wersja:** 1.0 • **Repozytorium:** `spoldzielnia_mieszkaniowa` • **Zespół:** 12 osób • **Technologie:** PHP + MySQL/MariaDB + Bootstrap 5.3 • **Uruchamianie:** XAMPP/localhost.

Ten dokument jest instrukcją wykonania projektu od pustego repozytorium do demonstracji końcowej. Każdy etap ma jednego właściciela i zakończony jest działającym fragmentem aplikacji. Nie utożsamiamy liczby podstron z nakładem pracy: trudniejsze moduły mają mniej ekranów, prostsze więcej operacji i testów. Po pierwszych dwóch etapach prace nad większością modułów można prowadzić równolegle.

## A. Reguły obowiązujące całą dwunastkę

### A1. Co każdy uczestnik ma zrobić przy swoim etapie

1. Sprawdzić `README.md`, ten plan i zakres przydzielonego modułu.
2. Zobaczyć aktualny stan gałęzi `main` na GitHub.
3. Pobrać kod i uruchomić stronę lokalnie według README.
4. Utworzyć gałąź `feature/NN-nazwa`, gdzie `NN` oznacza numer etapu.
5. Sprawdzić, czy tabele lub inne pliki, których potrzebuje, już istnieją.
6. Dodać albo zmienić tylko swoje strony, logikę modułu oraz niezbędne pliki wspólne.
7. Użyć istniejącego nagłówka, stopki, stylu i mechanizmów logowania; nie tworzyć drugiego systemu kont.
8. Zaprogramować nie tylko pozytywny przypadek, ale też brak danych, błędne pola i brak uprawnień.
9. Jeśli moduł potrzebuje tabeli: dopisać ją do `database/schema.sql`, zapewnić przykładowe rekordy w `database/seed.sql`, opisać zmiany i uzgodnić wpływ na istniejącą bazę.
10. Wprowadzić 2–5 fikcyjnych rekordów do danych demonstracyjnych, jeśli moduł wyświetla listy.
11. Samodzielnie przetestować na localhost, również przy wąskiej szerokości okna.
12. Uzupełnić README o zmiany instalacyjne, a ten plan o ewentualne decyzje projektowe.
13. Wykonać commit, przesłać gałąź na GitHub i otworzyć Pull Request.
14. Do opisu PR dołączyć: co działa, co zmieniono w bazie, jak przetestować, zrzut ekranu (jeżeli pomocny) i ograniczenia.
15. Uzyskać przegląd przynajmniej jednej innej osoby. Po akceptacji połączyć kod z `main` i zweryfikować uruchomienie tej gałęzi.

### A2. Zasady techniczne

- **PHP 8.1+; bez frameworka PHP** (takiego jak Laravel). Mniej zależności = prostszy lokalny start.
- **Bootstrap przez CDN**, bez obowiązkowego npm i kompilacji CSS. Jeśli uczelnia wymaga działania offline, zespół może później przenieść pliki Bootstrapa do `assets/vendor/`.
- **Jedna nawigacja** i wspólne `includes/header.php`, `includes/footer.php`. Trasy są zdefiniowane jawnie w `config/routes.php`.
- **PDO dla SQL**, zapytania parametryzowane, hashowanie haseł, sesje, kontrola ról i tokeny CSRF dla formularzy zmieniających dane.
- **Brak prawdziwych danych mieszkańców**. Zdjęcia i dokumenty mogą być przykładowe, pod warunkiem posiadania prawa do ich użycia.
- **Modyfikacje bazy są wersjonowane**. Na świeżej instalacji aktualne `schema.sql` + `seed.sql` muszą wystarczyć do odtworzenia aplikacji; zmiany istniejącej bazy wymagają udokumentowanej migracji lub jasnej instrukcji ponownego importu w środowisku testowym.
- **Rozdzielone funkcje**: zawartość stron w `pages/`, operacje na danych w `src/`, konfiguracja w `config/`, wygląd w `assets/`.
- **Wspólne nazwy ról**: `guest` (niezalogowany), `resident`, `admin`; tylko `resident` i `admin` znajdują się w bazie kont.

### A3. Kryteria przyjęcia każdego PR

- Nowy moduł działa na komputerze autora i recenzenta.
- Zmiana nie powoduje błędu na stronie głównej ani w menu.
- Niedozwolony użytkownik nie może wykonać chronionej operacji, także po bezpośrednim wpisaniu URL.
- Formularze mają etykiety, kontrolę poprawności danych i czytelny komunikat o wyniku.
- Nowy moduł opisany jest w README, a niezbędne pliki SQL są w repozytorium.
- Brak haseł, danych osobowych ani lokalnych załączników w commicie.

---

## ETAP 1 — GitHub, uruchamialna aplikacja i nawigacja

**Odpowiada:** osoba 1. **Zależności:** brak. **Efekt:** recenzent pobiera repozytorium i od razu uruchamia witrynę bez bazy danych.

### Instrukcja od zera

1. Założyć lub wybrać konto GitHub, które będzie właścicielem wspólnego repozytorium.
2. Na https://github.com/new wybrać tworzenie nowego repozytorium.
3. W polu **Repository name** wpisać dokładnie `spoldzielnia_mieszkaniowa`.
4. Jako opis wpisać np. `Projekt strony Spółdzielni Mieszkaniowej Mirmiłowo — PHP, MySQL, Bootstrap`.
5. Ustawić repozytorium jako **Public** (łatwy dostęp dla oceniającego) albo **Private**, jeśli uczelnia tak wymaga; przy prywatnym zapewnić oceniającemu dostęp.
6. Utworzyć repozytorium z gałęzią `main`. Można zaznaczyć początkowy README — zostanie zastąpiony przygotowaną dokumentacją.
7. Dodać wszystkich 11 współpracowników przez ustawienia repozytorium (**Collaborators / Access**). Wymagane jest zaakceptowanie zaproszeń.
8. Zainstalować na komputerze Git i XAMPP, następnie uruchomić Apache.
9. Przejść w terminalu do katalogu `C:\xampp\htdocs`.
10. Uruchomić `git clone https://github.com/NAZWA_KONTA/spoldzielnia_mieszkaniowa.git`.
11. Otworzyć katalog sklonowanego repozytorium w edytorze.
12. Skopiować do niego **zawartość** dostarczonego pakietu startowego; uważać, aby nie powstał zagnieżdżony katalog `spoldzielnia_mieszkaniowa/spoldzielnia_mieszkaniowa`.
13. Sprawdzić obecność plików `index.php`, `README.md`, `config/routes.php`, `includes/header.php`, `includes/footer.php`, `assets/css/style.css` i `docs/PLAN_PRAC.md`.
14. W przeglądarce otworzyć `http://localhost/spoldzielnia_mieszkaniowa/`.
15. Sprawdzić wyświetlanie tytułu **Spółdzielnia Mieszkaniowa „Mirmiłowo”**, menu, sześciu kafelków strony głównej i stopki.
16. Otworzyć każdą pozycję menu, również z poziomu telefonu lub węższego okna przeglądarki. Wszystkie podstrony mają informować o planowanej funkcji, a nie pokazywać błąd 404.
17. Sprawdzić adres nieistniejący, np. `?page=xyz`: powinien pokazywać stronę 404 z przyciskiem powrotu.
18. Sprawdzić w `README.md`, czy instrukcja klonowania zawiera **placeholder właściciela repozytorium**, a następnie po utworzeniu repozytorium zastąpić go prawdziwym adresem w README. Bez tej poprawki instrukcja nie będzie gotowa dla oceniającego.
19. Wykonać `git status`, sprawdzić pliki, a następnie `git add .`, `git commit -m "Dodaj uruchamialny szkielet Mirmilowa"` i `git push origin main` (dla pierwszego zatwierdzenia, zanim zespół rozpocznie równoległe prace).
20. Otworzyć repozytorium na GitHub, upewnić się, że wszystkie pliki są widoczne, a w README znajduje się instrukcja uruchomienia.
21. Poprosić osobę 2 o **test na czysto**: sklonowanie projektu na innym komputerze i otwarcie strony wyłącznie według README. Poprawić ujawnione problemy.
22. Włączyć ochronę gałęzi `main` (np. wymóg co najmniej jednej recenzji PR), jeśli pozwalają na to ustawienia repozytorium.
23. Utworzyć GitHub Issues dla kolejnych 11 etapów i przydzielić je osobom.
24. Oznaczyć zakończenie etapu na GitHub (np. tag `v0.1.0` lub GitHub Release) — oceniający może odtworzyć pierwszy kamień milowy.

**Co musi działać:** publiczna strona główna, wszystkie docelowe pozycje nawigacji, podstrony zastępcze, responsywny Bootstrap, czytelny README. **Czego nie wymagamy:** MySQL, logowania, panelu administracyjnego i aktywnych formularzy.

**Test odbiorowy:** osoba nieznająca kodu klonuje repozytorium, uruchamia tylko Apache i widzi stronę oraz menu. Nie musi tworzyć bazy danych.

## ETAP 2 — Baza danych, konta i uprawnienia

**Odpowiada:** osoba 2. **Zależność:** ukończony etap 1. **Efekt:** spójna baza i mechanizm logowania dla wszystkich przyszłych modułów.

1. Zainstalować/uruchomić w XAMPP Apache i MySQL/MariaDB.
2. Otworzyć `http://localhost/phpmyadmin` i utworzyć bazę `mirmilowo` (utf8mb4).
3. Zaprojektować tabelę `users`: `id`, `name`, `email` (unikalny), `password_hash`, `role`, `created_at`, `updated_at`.
4. Zdefiniować SQL tworzący tabelę w `database/schema.sql` w sposób możliwy do wykonania na czystej bazie.
5. Utworzyć `config/database.example.php` z nazwami pól połączenia i `config/database.local.php` lokalnie (ten drugi plik ignorowany przez Git).
6. Napisać w `src/` wspólną funkcję tworzenia połączenia PDO, z bezpieczną obsługą błędów.
7. Napisać formularz logowania: adres e-mail i hasło; przy niepowodzeniu pokazać ogólny błąd bez ujawniania, czy konto istnieje.
8. Implementować `password_hash`, `password_verify`, regenerację identyfikatora sesji przy logowaniu oraz wylogowanie.
9. Utworzyć funkcje `currentUser()`, `isLoggedIn()`, `requireLogin()` i `requireRole('admin')` albo spójne odpowiedniki.
10. Zapewnić sposób **lokalnego stworzenia pierwszego administratora** (np. skrypt CLI), bez hasła zapisanego w kodzie czy README.
11. Zapewnić utworzenie testowego mieszkańca z wybranym lokalnie hasłem (np. przez skrypt CLI lub bezpieczny formularz administracyjny).
12. Dodać strony logowania i panelu konta do nawigacji z czytelnym podziałem widoku gościa i zalogowanego użytkownika.
13. Sprawdzić bezpośredni URL strony dla mieszkańca przy braku logowania — system ma przekierować do logowania albo zwrócić odmowę dostępu.
14. Sprawdzić niedozwoloną próbę wejścia na panel admina z konta mieszkańca.
15. Przygotować wspólne helpery: kodowanie HTML, walidacja i token CSRF dla formularzy POST.
16. Opisać pełne kroki importu SQL, konfiguracji lokalnej i utworzenia kont w README.
17. Przetestować od zera: skasować **testową** bazę (nie dane innych!), utworzyć ponownie, zaimportować SQL i zalogować się.
18. Uzgodnić z resztą zespołu schemat nazewnictwa tabel oraz zasady dodawania nowych definicji bazy.

**Test odbiorowy:** da się stworzyć bazę według README; admin loguje się do panelu, mieszkaniec do strefy mieszkańca, gość nie widzi prywatnych treści. Uwaga: bezpieczeństwo jest wymagane również na localhost.

## ETAP 3 — Spółdzielnia: Struktura i Historia

**Odpowiada:** osoba 3. **Zależność:** etap 2. **Efekt:** dwie działające sekcje informacyjne z edycją przez administratora.

1. Otworzyć placeholdery „Struktura” i „Historia”, rozpoznać ich identyfikatory routingu.
2. Zaplanować zawartość struktury: zarząd, rada, stanowiska, opis kompetencji (wszystko fikcyjne).
3. Zaplanować historię w formie chronologii: data/rok, tytuł zdarzenia, opis.
4. Dodać tabelę/tabele dla informacji spółdzielni (np. `organization_pages` i opcjonalnie `history_events`) do SQL.
5. Dodać kilka fikcyjnych wpisów demonstracyjnych do seedów.
6. Napisać stronę Struktura w Bootstrap: karty ról, sekcje organizacyjne i dane publiczne.
7. Napisać stronę Historia w Bootstrap: uporządkowana oś czasu lub lista chronologiczna.
8. Dodać widok admina, w którym można dodać lub zmienić treść i kolejność wpisów.
9. Dodać usuwanie z potwierdzeniem i walidację dat/nazw.
10. Zadbać, aby gość widział informacje, ale nie miał prawa ich edycji.
11. Przetestować pustą bazę, więcej wpisów, długie opisy i układ mobilny.
12. Udokumentować strukturę tabel i testy, wysłać PR.

**Test odbiorowy:** gość otwiera obie sekcje, admin może zmienić treści, a zmiany są widoczne po odświeżeniu.

## ETAP 4 — Ogłoszenia: Imprezy kulturalne

**Odpowiada:** osoba 4. **Zależność:** etap 2. **Efekt:** publiczne ogłoszenia o wydarzeniach, zarządzane przez administratora.

1. Ustalić dane wydarzenia: tytuł, opis, data i godzina rozpoczęcia, miejsce, status publikacji.
2. Dodać tabelę `cultural_events` do `schema.sql`.
3. Dodać 3 fikcyjne wydarzenia do seedów (przyszłe i archiwalne).
4. Zbudować publiczną listę wydarzeń, posortowaną po dacie.
5. Zbudować stronę szczegółów wydarzenia.
6. Oznaczyć datę i miejsce w czytelnej formie.
7. Dodać prosty filtr „nadchodzące / archiwalne”.
8. Dodać formularz admina do tworzenia wydarzenia.
9. Dodać edycję, publikację/ukrywanie i usuwanie.
10. Zweryfikować poprawność dat (np. data zakończenia nie przed rozpoczęciem, jeśli jest używana).
11. Zablokować akcje admina osobom niezalogowanym.
12. Przygotować sposób udostępniania wydarzeń kalendarzowi w etapie 7 (np. poprzez czytelną funkcję pobierania danych).
13. Przetestować pustą listę, wiele wydarzeń i widok mobilny.
14. Zaktualizować dokumentację i wysłać PR.

**Test odbiorowy:** admin dodaje wydarzenie, gość widzi je w sekcji Imprezy kulturalne; nieopublikowane wydarzenie nie jest pokazywane publicznie.

## ETAP 5 — Ogłoszenia: Przetargi

**Odpowiada:** osoba 5. **Zależność:** etap 2. **Efekt:** lista przetargów z terminami i statusem.

1. Uzgodnić pola: tytuł, opis, numer/referencja, data ogłoszenia, termin składania ofert, status, ewentualny plik.
2. Utworzyć tabelę `tenders` w SQL.
3. Dodać fikcyjne przetargi z różnymi terminami i statusami.
4. Zbudować listę przetargów: tytuł, termin, status, przycisk „Szczegóły”.
5. Zbudować stronę szczegółów z opisem i informacją o postępowaniu.
6. Dodać filtrowanie po statusie (np. otwarty/zakończony).
7. Dodać sortowanie po zbliżającym się terminie.
8. Dodać formularz administratora tworzący przetarg.
9. Dodać edycję i zamknięcie przetargu; przy usuwaniu wymagane potwierdzenie.
10. Zweryfikować poprawność dat i długości pól.
11. Jeśli są załączniki, przechowywać je w kontrolowany sposób (bez bezpośredniego wykonywania plików).
12. Przygotować pole terminu dla przyszłego kalendarza (etap 7).
13. Sprawdzić, że gość nie może zmienić danych w bazie.
14. Przetestować listę pustą, zakończone postępowanie i termin w przeszłości; dopisać dokumentację.

**Test odbiorowy:** da się dodać, edytować i zamknąć przetarg; gość widzi aktualną listę i daty.

## ETAP 6 — Ogłoszenia: Wynajem

**Odpowiada:** osoba 6. **Zależność:** etap 2. **Efekt:** katalog ofert wynajmu dla osiedla.

1. Uzgodnić dane oferty: tytuł, typ lokalu, metraż, orientacyjna cena, opis, dostępność, status.
2. Utworzyć tabelę `rentals`; opcjonalny obraz przechowywać jako bezpiecznie obsługiwany załącznik.
3. Dodać fikcyjne oferty: np. lokal usługowy, miejsce parkingowe, pomieszczenie użytkowe.
4. Zbudować publiczne karty ofert w Bootstrap.
5. Zbudować szczegóły oferty.
6. Dodać filtr po typie lub statusie dostępności.
7. Dodać sortowanie po cenie lub powierzchni (jedno jest wystarczające, jeśli zakres rośnie).
8. Dodać formularz administracyjny tworzenia oferty.
9. Dodać edycję oferty i zmianę statusu na „niedostępna”.
10. Dodać walidację ceny i powierzchni (liczby nieujemne).
11. Przy ewentualnych zdjęciach: walidować format i rozmiar, nie ufać rozszerzeniu pliku.
12. Dodać komunikat dla pustej listy i poprawić wygląd mobilny.
13. Sprawdzić, że operacje administracyjne wymagają uprawnień.
14. Udokumentować dane demonstracyjne, testy i wysłać PR.

**Test odbiorowy:** gość przegląda i filtruje oferty; administrator może je dodawać i zmieniać dostępność.

## ETAP 7 — Kalendarz spraw, wydarzeń i terminów

**Odpowiada:** osoba 7. **Zależność:** etap 2, a do pełnej integracji etapy 4–6. **Efekt:** kalendarz z kategoriami i kolorami.

1. Ustalić stałe kategorie: `urgent`, `deadline`, `culture`, `tender`, `other` (nazwy techniczne), odpowiednio czerwony, pomarańczowy, zielony, niebieski, szary.
2. Dodać tabelę `calendar_events` dla własnych terminów: id, tytuł, data/godzina, kategoria, opis, data utworzenia.
3. Przygotować dane demonstracyjne różnych kategorii i dat.
4. Zbudować widok bieżącego miesiąca z przechodzeniem na miesiąc poprzedni/następny.
5. W każdym dniu pokazać wydarzenia przypisane do właściwej daty.
6. Dodać legendę kategorii z kolorem **i tekstową nazwą**, aby nie polegać tylko na kolorze.
7. Po kliknięciu wpisu wyświetlić szczegóły wydarzenia i termin.
8. Dodać filtrowanie według kategorii, bez usuwania wpisów z bazy.
9. Dodać formularz admina tworzenia terminu, edycję i usuwanie.
10. Zadbać o poprawność miesięcy i lat na przełomie roku; nie zgubić wpisów z różnych stref czasowych (daty osiedlowe traktować konsekwentnie jako `Europe/Warsaw`).
11. Zintegrować dane o imprezach z etapu 4 oraz terminach przetargów z etapu 5. Jeśli moduł nie został jeszcze połączony, w PR wyraźnie oznaczyć niezrealizowane połączenie jako zadanie integracyjne.
12. Gdy jest potrzeba, uwzględnić daty związane z wynajmem/wnioskami w sposób uzgodniony z autorami modułów (bez kopiowania wpisów do wielu tabel).
13. Sprawdzić widok mobilny; w razie potrzeby pokazać listę wydarzeń pod kalendarzem.
14. Przetestować daty skrajne, puste miesiące, filtry, brak uprawnień i brak powielonych rekordów; zaktualizować dokumentację.

**Test odbiorowy:** widoczne są wydarzenia z co najmniej czterech kategorii, kolory odpowiadają legendzie, a kliknięcie pokazuje termin i opis.

## ETAP 8 — Inwestycje na osiedlu

**Odpowiada:** osoba 8. **Zależność:** etap 2. **Efekt:** publiczny rejestr inwestycji i ich postępów.

1. Ustalić model: nazwa, opis, data startu, przewidywane zakończenie, status, procent realizacji, ewentualny obraz.
2. Dodać tabelę `investments` i dane demonstracyjne.
3. Utworzyć listę inwestycji z kartami Bootstrap i statusem.
4. Wprowadzić trzy statusy: planowana, w trakcie, zakończona.
5. Dodać pasek postępu (`progress`) z widoczną wartością procentową.
6. Zbudować podstronę szczegółów inwestycji.
7. Dodać filtr po statusie.
8. Dodać formularz admina tworzenia inwestycji.
9. Dodać edycję opisu i postępu oraz zmianę statusu.
10. Walidować procent w zakresie 0–100, daty i wymagane pola.
11. Rozwiązać sytuację, gdy inwestycja nie ma jeszcze daty zakończenia.
12. Sprawdzić brak danych, długie opisy, mobile i próby edycji przez gościa.
13. Dodać sekcję/odnośnik do inwestycji ze strony głównej (bez zmiany wspólnego layoutu bez uzgodnienia).
14. Uzupełnić README i wysłać PR.

**Test odbiorowy:** zmiana postępu w panelu administracyjnym aktualizuje publiczny opis inwestycji.

## ETAP 9 — Strefa mieszkańca: Druki i wnioski, Statut i regulaminy

**Odpowiada:** osoba 9. **Zależność:** etap 2. **Efekt:** repozytorium dokumentów do pobrania.

1. Uzgodnić kategorie plików: `form` (druk/wniosek), `statute` (statut), `regulation` (regulamin).
2. Zaprojektować tabelę `documents` z tytułem, kategorią, opisem, nazwą zapisanego pliku, datą dodania i widocznością.
3. Dodać do projektu fikcyjne, nieszkodliwe pliki PDF w ramach danych przykładowych (lub dokumenty generowane na potrzeby demonstracji).
4. Zbudować publiczną stronę Druki i wnioski z listą i pobieraniem.
5. Zbudować publiczną stronę Statut i regulaminy z podziałem na rodzaje dokumentów.
6. Dodać wyszukiwanie po tytule lub filtrowanie po kategorii.
7. Dodać formularz administratora dodawania dokumentu.
8. Wprowadzić walidację rozmiaru, typu MIME, rozszerzenia i nazw; przechowywać pliki w niepublicznym `storage/uploads`.
9. Dodać kontrolowany skrypt pobierania, który sprawdza identyfikator dokumentu zamiast przyjmować niezweryfikowaną ścieżkę.
10. Dodać możliwość zmiany tytułu/kategorii oraz wyłączenia publikacji dokumentu.
11. Unikać uruchamiania przesłanych plików jako kodu PHP; nie ufać nazwom plików użytkownika.
12. Zapewnić komunikat przy braku dokumentów i przy próbie pobrania nieistniejącego pliku.
13. Sprawdzić pobieranie, brak uprawnień do dodawania i działanie na telefonie.
14. Opisać ograniczenia rozmiaru pliku i strukturę katalogu w README.

**Test odbiorowy:** admin publikuje plik, gość pobiera go z listy; nie da się przeglądać całego katalogu uploadów przez URL.

## ETAP 10 — Czat mieszkańców

**Odpowiada:** osoba 10. **Zależność:** etap 2. **Efekt:** prosty czat dla zalogowanych mieszkańców i administratorów.

1. Potwierdzić, że czat nie jest publiczny; wymaga aktywnej sesji mieszkańca lub administratora.
2. Utworzyć tabelę `chat_messages`: id, autor (`user_id`), treść, data wysłania, status moderacji.
3. Utworzyć stronę czatu z listą wiadomości i polem wpisywania tekstu.
4. Dodać obsługę wysłania wiadomości metodą POST.
5. Wymusić maksymalną długość wiadomości i blokadę pustych wpisów.
6. Wyświetlać autora i godzinę, escapować HTML, nie interpretować skryptów w wiadomościach.
7. Pokazywać wiadomości chronologicznie; ograniczyć liczbę pobieranych wpisów (np. ostatnie 50).
8. Po wysłaniu przekierować do czatu, aby odświeżenie strony nie wysłało tej samej wiadomości drugi raz.
9. Dodać proste automatyczne odświeżanie co kilkanaście sekund **lub** zwykły przycisk „Odśwież”; WebSocket nie jest wymagany.
10. Dodać możliwość usuwania/ukrywania wiadomości przez administratora, z kontrolą serwerową.
11. Wprowadzić proste ograniczenie liczby wysyłanych wiadomości (ochrona przed spamem).
12. Zabezpieczyć formularze przed CSRF i dane przed XSS.
13. Sprawdzić rozmowę na dwóch różnych kontach i odmowę dostępu dla gościa.
14. Udokumentować ograniczenie: to czat demonstracyjny, nie komunikator czasu rzeczywistego.

**Test odbiorowy:** dwóch mieszkańców widzi wspólne wiadomości po odświeżeniu; gość nie odczyta czatu, admin może ukryć wpis.

## ETAP 11 — E-mieszkaniec: wnioski online i demonstracja podpisu

**Odpowiada:** osoba 11. **Zależności:** etap 2 i zasady załączników z etapu 9. **Efekt:** formularz i obsługa wniosków bez fałszywych deklaracji o Profilem Zaufanym.

1. Wybrać 2–3 kategorie spraw demonstracyjnych: zgłoszenie usterki, wniosek administracyjny, prośba o zaświadczenie.
2. Zaprojektować tabelę `applications`: id, `user_id`, rodzaj, tytuł, opis, data złożenia, status, informacja o pliku i trybie podpisu.
3. Zdefiniować statusy: szkic / wysłany / w trakcie / zakończony / odrzucony (można uprościć do czterech).
4. Przygotować dla mieszkańca formularz wyboru rodzaju sprawy i jej opisu.
5. Umożliwić dodanie jednego załącznika z kontrolą typu i wielkości.
6. Dodać podsumowanie wniosku przed potwierdzeniem.
7. Dodać osobny krok z **widocznym ostrzeżeniem „Symulacja podpisu — bez skutku prawnego”**; żadnego pola na login i hasło do Profilem Zaufanego.
8. Opcjonalnie podać link do oficjalnej usługi podpisania pliku na gov.pl i pozwolić ponownie załączyć dokument; bez weryfikacji kryptograficznej taki plik traktować jako „przesłany plik”, nie „zweryfikowany podpis”.
9. Zapisać nowy wniosek do bazy z identyfikatorem konta zalogowanego mieszkańca.
10. Dodać „Moje wnioski”: listę własnych zgłoszeń ze statusem i datą.
11. Zablokować możliwość odczytu cudzego wniosku przez zmianę ID w adresie.
12. Zbudować dla administratora listę wszystkich wniosków oraz zmianę statusu i opcjonalną notatkę.
13. Przy zmianie statusu stosować walidację i ochronę CSRF.
14. Przetestować dwóch mieszkańców, administratora, brak załącznika, niepoprawny plik, dostęp bez logowania i widoki mobilne.
15. Wyjaśnić w README i interfejsie ograniczenia prawne i techniczne funkcji demonstracyjnej.

**Test odbiorowy:** mieszkaniec wysyła wniosek, widzi własny numer i status; admin zmienia status; mieszkaniec nie zobaczy spraw innych kont. Nigdzie nie pojawia się fałszywe potwierdzenie państwowego podpisu.

## ETAP 12 — Kontakt, integracja i wydanie końcowe

**Odpowiada:** osoba 12. **Zależność:** etap 2; finał po PR modułów 3–11. **Efekt:** działający formularz kontaktowy i spójna wersja demonstracyjna.

1. Zaprojektować sekcję Kontakt: fikcyjny adres biura, telefon, godziny pracy, e-mail i informacja o osiedlu.
2. Zbudować czytelną stronę z informacjami i formularzem: imię, e-mail, temat, treść.
3. Dodać tabelę `contact_messages` do SQL i fikcyjne rekordy demonstracyjne.
4. Zaprogramować zapis wiadomości do bazy — nie jest wymagane wysyłanie prawdziwych e-maili.
5. Dodać walidację adresu e-mail, wymaganych pól, długości wiadomości i token CSRF.
6. Dodać czytelne potwierdzenie „Wiadomość zapisana w systemie demonstracyjnym”, bez sugerowania doręczenia wiadomości na rzeczywistą skrzynkę.
7. Zbudować widok administratora do odczytu zgłoszeń i oznaczania ich jako obsłużonych.
8. Przygotować krótką checklistę integracyjną dla wszystkich 12 osób i rozdzielić testy modułów między ich autorów (osoba 12 **nie naprawia samodzielnie cudzego kodu**).
9. Po połączeniu wszystkich PR sprawdzić od nowa wszystkie linki menu, logowanie, nadawanie ról i formularze.
10. Sprawdzić, czy dane z modułów wydarzeń i terminów są pokazywane w kalendarzu, a nie jedynie w odrębnych sekcjach.
11. Zweryfikować, że każda tabela jest opisana w `database/schema.sql`, a dane przykładowe wystarczają do demonstracji.
12. Zweryfikować uprawnienia: gość nie używa prywatnych funkcji, mieszkaniec nie edytuje treści administracyjnych, administrator obsługuje zgłoszenia.
13. Poprosić **inną osobę niż autor instalacji** o świeże sklonowanie repozytorium na nowym komputerze i uruchomienie całości wyłącznie według README.
14. Sprawdzić podstawową dostępność: klawiatura, etykiety formularzy, kontrast, legenda kalendarza, działanie telefonu.
15. Uzupełnić dokumentację o ostateczne wymagania i proces instalacji SQL; usunąć z README nieaktualne instrukcje „po etapie 2”.
16. Przygotować 5–10 minut scenariusza prezentacji: strona główna → przetargi → kalendarz → logowanie mieszkańca → czat → wysłanie wniosku → zmiana statusu przez admina.
17. Wygenerować końcowy tag/release (np. `v1.0.0`) i upewnić się, że `main` jest stabilny.

**Test odbiorowy:** instalacja na czystym komputerze od zera, wszystkie moduły dostępne, demonstracja przechodzi bez błędów, a README zgadza się z rzeczywistą procedurą.

---

## B. Sugerowany harmonogram zależności (można pracować równolegle)

| Moment | Kto pracuje | Dlaczego |
| --- | --- | --- |
| Start | Osoba 1 | Zakłada GitHub i publikuje działający szkielet |
| Po zatwierdzeniu etapu 1 | Osoba 2 | Buduje wspólne logowanie, role i bazę |
| Po zatwierdzeniu etapu 2 | Osoby 3, 4, 5, 6, 8, 9, 10, 11, 12 | Mogą wykonywać niezależne moduły we własnych gałęziach |
| Po istnieniu wspólnych tabel zdarzeń | Osoba 7 | Finalizuje połączenie kalendarza z imprezami i terminami |
| Końcówka | Wszyscy; osoba 12 koordynuje checklistę | Testy integracyjne, naprawy przez autorów modułów i wydanie |

**Istotne:** „12 etapów” oznacza tutaj 12 **pakietów odpowiedzialności**, nie 12 tygodni pracy jeden po drugim. Dzięki temu zespół 12 osób nie czeka bezczynnie na swoją kolej. Liczba godzin powinna być monitorowana w GitHub Issues i wyrównywana pomocniczymi zadaniami (testy, dokumentacja, przeglądy PR).

## C. Gotowy proces w GitHub dla każdej osoby

**Pobranie projektu** (pierwszy raz):

```powershell
cd C:\xampp\htdocs
git clone https://github.com/NAZWA_KONTA/spoldzielnia_mieszkaniowa.git
cd spoldzielnia_mieszkaniowa
```

**Rozpoczęcie etapu na nowej gałęzi:**

```bash
git switch main
git pull origin main
git switch -c feature/08-inwestycje
```

**Publikacja zmian:**

```bash
git status
git add .
git commit -m "Dodaj liste inwestycji i edycje w panelu admina"
git push -u origin feature/08-inwestycje
```

**Na stronie GitHub:** otwórz PR z gałęzi `feature/08-inwestycje` do `main`; w opisie napisz, jak uruchomić/testować, jakie pliki SQL zmieniono i kto powinien to zweryfikować.

**Po scaleniu swojej pracy:**

```bash
git switch main
git pull origin main
```

W razie konfliktu w Git nie nadpisywać cudzej pracy. Poprosić autora kolidującego modułu o wspólne rozwiązanie, następnie ponownie wykonać testy.

## D. Minimalne scenariusze testowe na ocenę

1. **Pierwsze uruchomienie:** pobranie `main`, uruchomienie Apache, import SQL i otwarcie strony.
2. **Nawigacja:** każda sekcja z menu otwiera właściwą zawartość.
3. **Autoryzacja:** gość nie czyta czatu ani cudzych wniosków; mieszkaniec nie zarządza treściami.
4. **CRUD administracyjny:** admin tworzy ogłoszenie, edytuje je, a publiczna lista natychmiast odzwierciedla zmianę.
5. **Kalendarz:** pokazuje wydarzenie kulturalne, termin przetargu, deadline i pilną sprawę z poprawną legendą.
6. **Dokumenty:** publikacja i pobranie bezpośrednio z sekcji Druki i wnioski.
7. **Czat:** zapis i odczyt wiadomości z dwóch różnych kont.
8. **E-mieszkaniec:** złożenie własnego wniosku, aktualizacja statusu przez admina i odmowa odczytu cudzej sprawy.
9. **Kontakt:** walidacja nieprawidłowego formularza i zapis poprawnej wiadomości.
10. **Błędy:** nieistniejąca podstrona zwraca 404; strona działa przy pustych tabelach.
11. **RWD:** menu oraz formularze są obsługiwane na małym ekranie.
12. **Bezpieczeństwo demonstracyjne:** nieprzesłane hasła i nieuprawnione modyfikacje, dane testowe zamiast prawdziwych, żadnego pozorowanego podpisu PZ.

## E. Definicja końca projektu (Definition of Done)

Projekt jest ukończony dopiero, kiedy **inna osoba, bez wiedzy o kodzie**, na komputerze z nowo zainstalowanym XAMPP:

- sklonuje repozytorium,
- uruchomi Apache i MySQL/MariaDB,
- zaimportuje przygotowane skrypty SQL,
- utworzy konta lokalne według dokumentacji,
- przejdzie wszystkie scenariusze A/B z prezentacji bez ręcznego poprawiania plików,
- zobaczy wyraźne ograniczenie podpisu Profilem Zaufanym w wersji edukacyjnej.

Jeżeli trzeba wysłać autorowi prywatną wiadomość w stylu „zmień sobie ręcznie wartość w kodzie” albo „na moim komputerze działa”, etap nie jest zamknięty.
