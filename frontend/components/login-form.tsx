"use client";

import { useState, type FormEvent } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { authRequest, useAuth } from "./auth";

export function LoginForm() {
  const { user, loading, setUser } = useAuth();
  const router = useRouter();
  const [pending, setPending] = useState(false);
  const [error, setError] = useState("");
  async function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault(); setError(""); setPending(true);
    const form = event.currentTarget;
    const fields = new FormData(form);
    try {
      const result = await authRequest("login", { email: String(fields.get("email")).trim(), password: String(fields.get("password")) });
      if (!result.user) throw new Error("Nieprawidłowa odpowiedź serwera.");
      form.reset(); setUser(result.user); router.push("/"); router.refresh();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Nie udało się zalogować.");
    } finally { setPending(false); }
  }
  if (user) return <><p>Jesteś zalogowany jako {user.email}.</p><Link className="btn btn-primary" href="/">Przejdź na stronę główną</Link></>;
  return <form onSubmit={submit} aria-busy={pending}>
    {error && <div className="alert alert-danger" role="alert">{error}</div>}
    <div className="mb-3"><label htmlFor="email" className="form-label">Adres e-mail</label><input type="email" className="form-control" id="email" name="email" autoComplete="username" maxLength={254} required disabled={pending} /></div>
    <div className="mb-4"><label htmlFor="password" className="form-label">Hasło</label><input type="password" className="form-control" id="password" name="password" autoComplete="current-password" required disabled={pending} /></div>
    <button type="submit" className="btn btn-primary w-100" disabled={pending || loading}>{pending ? "Logowanie…" : "Zaloguj się"}</button>
  </form>;
}
