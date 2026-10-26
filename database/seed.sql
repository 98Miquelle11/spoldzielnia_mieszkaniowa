USE mirmilowo;
-- Brak wspólnych haseł i automatycznego nadawania uprawnień.
-- Utwórz konta narzędziem Go po imporcie schema.sql (patrz README):
-- go run ./cmd/create-user -email admin@example.test -role admin
-- go run ./cmd/create-user -email mieszkaniec@example.test -role mieszkaniec
-- Hasło podajesz interaktywnie; w bazie zapisany zostanie bcrypt.
