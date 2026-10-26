# Plan zespołowy po migracji — 12 etapów

Zachowano podział odpowiedzialności z oryginalnego projektu. Implementację
backendu rozwijaj w Go, widoki w Next.js, a strukturę danych nadal w SQL.

| Etap | Zakres | Stan tej paczki |
| --- | --- | --- |
| 1 | Struktura, strona główna, nawigacja, dokumentacja | Przeniesione do Next.js |
| 2 | Baza, konta, logowanie, wylogowanie, role | Kod Go i SQL gotowy; wymaga importu i testu z lokalnym MySQL |
| 3 | Struktura i historia spółdzielni | Strony informacyjne; treści i edycja do wykonania |
| 4 | Imprezy kulturalne | Moduł planowany |
| 5 | Przetargi | Moduł planowany |
| 6 | Wynajem | Moduł planowany |
| 7 | Kalendarz | Moduł planowany |
| 8 | Inwestycje | Moduł planowany |
| 9 | Druki, statut i regulaminy | Moduł planowany |
| 10 | Czat mieszkańców | Moduł planowany |
| 11 | E-mieszkaniec | Moduł planowany |
| 12 | Kontakt, integracja i testy końcowe | Moduł planowany |

## Zasady dalszej pracy

- Nie twórz drugiego systemu kont. Korzystaj z tabel `users`, `sessions` i z
  middleware `RequireRole` w Go; ukrycie przycisku w React nie chroni API.
- Dla nowych tabel dodawaj wersjonowane skrypty migracji SQL oraz fikcyjne dane
  demonstracyjne. Nie zastępuj MySQL bazą innego typu.
- Publiczne listy i szczegóły udostępniaj przez GET; operacje zapisu wymagają
  autoryzacji, walidacji i ochrony CSRF po stronie Go.
- Prywatne zasoby, np. wnioski mieszkańca, wymagają także sprawdzenia właściciela
  rekordu. Sama zgodność roli nie wystarczy.
- Formularze i pobieranie plików powinny obsługiwać brak danych, błędy i brak
  uprawnień. Nie wystawiaj katalogu z prywatnymi załącznikami jako publicznego.
- Kalendarz prezentuj w strefie Europe/Warsaw; czas sesji pozostaje w UTC.
- Zachowaj wspólny layout, menu, polskie treści i responsywność.
- Wnioskowanie i podpis Profilem Zaufanym pozostają planem edukacyjnym, bez
  deklarowania prawdziwej integracji lub skutku prawnego.
- Każdy etap wykonuj na osobnej gałęzi Git, z opisem instalacji i testów w PR.
  Po połączeniu zmian sprawdź instalację od zera według README.
