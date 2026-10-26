package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/mail"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"
	"mirmilowo/backend/internal/database"
)

func main() {
	email := flag.String("email", "", "Adres e-mail konta")
	role := flag.String("role", "mieszkaniec", "mieszkaniec albo admin")
	flag.Parse()
	*email = strings.TrimSpace(*email)
	parsed, err := mail.ParseAddress(*email)
	if err != nil || parsed.Address != *email || len(*email) > 254 {
		log.Fatal("Podaj poprawny adres e-mail przez -email")
	}
	if *role != "admin" && *role != "mieszkaniec" {
		log.Fatal("Rola: admin albo mieszkaniec")
	}
	fmt.Fprint(os.Stderr, "Hasło (12–72 bajty, wpisywane niewidocznie): ")
	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Fprintln(os.Stderr)
	if err != nil {
		log.Fatal("Uruchom polecenie w interaktywnym terminalu: ", err)
	}
	if len(password) < 12 || len(password) > 72 {
		log.Fatal("Hasło musi zawierać 12–72 bajty")
	}
	hash, err := bcrypt.GenerateFromPassword(password, 12)
	clear(password)
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db, err := database.Open(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// Deliberately no upsert: rerunning cannot silently replace a password or role.
	if _, err = db.ExecContext(ctx, "INSERT INTO users (email, password_hash, role) VALUES (?, ?, ?)", *email, string(hash), *role); err != nil {
		log.Fatal("Nie utworzono konta (sprawdź m.in. czy e-mail już istnieje): ", err)
	}
	fmt.Printf("Utworzono konto %s (%s).\n", *email, *role)
}
