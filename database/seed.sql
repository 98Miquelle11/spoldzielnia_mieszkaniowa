USE mirmilowo;
SET NAMES utf8mb4;

START TRANSACTION;

INSERT INTO cultural_events
  (title, description, starts_at, ends_at, location, is_published)
SELECT
  'Wspólne wyjście do GOK na wystawę sztuki',
  'Zapraszamy mieszkańców Mirmiłowa na wspólne wyjście do Gminnego Ośrodka Kultury. Obejrzymy wystawę sztuki lokalnych twórców i porozmawiamy o prezentowanych pracach. Zbiórka przy wejściu do GOK o godzinie 17:00. Wydarzenie fikcyjne, przygotowane na potrzeby projektu akademickiego.',
  '2026-11-21 16:00:00',
  '2026-11-21 18:00:00',
  'Gminny Ośrodek Kultury, sala wystawowa (lokalizacja fikcyjna)',
  1
WHERE NOT EXISTS (
  SELECT 1 FROM cultural_events
  WHERE title = 'Wspólne wyjście do GOK na wystawę sztuki'
    AND starts_at = '2026-11-21 16:00:00'
);

INSERT INTO cultural_events
  (title, description, starts_at, ends_at, location, is_published)
SELECT
  'Sylwester 2024 dla mieszkańców',
  'Archiwalne ogłoszenie o spotkaniu sylwestrowym mieszkańców Mirmiłowa, podczas którego wspólnie żegnaliśmy rok 2023 i witaliśmy rok 2024. W programie przewidziano muzykę, wspólne rozmowy i symboliczne odliczanie do północy. Wydarzenie fikcyjne, przygotowane na potrzeby projektu akademickiego.',
  '2023-12-31 19:00:00',
  '2024-01-01 01:00:00',
  'Świetlica osiedlowa Mirmiłowo (lokalizacja fikcyjna)',
  1
WHERE NOT EXISTS (
  SELECT 1 FROM cultural_events
  WHERE title = 'Sylwester 2024 dla mieszkańców'
    AND starts_at = '2023-12-31 19:00:00'
);

INSERT INTO cultural_events
  (title, description, starts_at, ends_at, location, is_published)
SELECT
  'Sylwester 2027 dla mieszkańców',
  'Robocze ogłoszenie o planowanym spotkaniu sylwestrowym mieszkańców Mirmiłowa na zakończenie roku 2026. Wstępnie przewidziano muzykę, poczęstunek i wspólne powitanie roku 2027. Ogłoszenie jest ukryte do czasu zatwierdzenia programu przez administratora. Wydarzenie fikcyjne, przygotowane na potrzeby projektu akademickiego.',
  '2026-12-31 19:00:00',
  '2027-01-01 01:00:00',
  'Świetlica osiedlowa Mirmiłowo (lokalizacja fikcyjna)',
  0
WHERE NOT EXISTS (
  SELECT 1 FROM cultural_events
  WHERE title = 'Sylwester 2027 dla mieszkańców'
    AND starts_at = '2026-12-31 19:00:00'
);
COMMIT;
