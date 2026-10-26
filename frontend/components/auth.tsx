"use client";

import { createContext, useContext, useEffect, useState, type ReactNode } from "react";

export type User = { id: number; email: string; role: "admin" | "mieszkaniec" };
type State = { user: User | null; loading: boolean; error: string; setUser: (user: User | null) => void; reload: () => Promise<void> };
const AuthContext = createContext<State | null>(null);

export async function authRequest(path: "login" | "logout", body: object = {}) {
  const response = await fetch(`/api/auth/${path}`, {
    method: "POST", headers: { "Content-Type": "application/json" },
    credentials: "same-origin", body: JSON.stringify(body),
    signal: AbortSignal.timeout(15000),
  });
  let data;
  try { data = await response.json(); }
  catch { throw new Error("Nie można połączyć się z serwerem. Spróbuj ponownie później."); }
  if (!response.ok) throw new Error(data.error ?? "Nie udało się wykonać operacji.");
  return data as { user?: User; ok?: boolean };
}

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  async function reload() {
    setLoading(true);
    try {
      const response = await fetch("/api/auth/me", { cache: "no-store", signal: AbortSignal.timeout(10000) });
      if (!response.ok) throw new Error();
      const data = await response.json();
      setUser(data.user); setError("");
    } catch {
      setError("Nie udało się sprawdzić stanu logowania.");
    } finally { setLoading(false); }
  }

  useEffect(() => { void reload(); }, []);
  function updateUser(value: User | null) { setUser(value); setError(""); }
  return <AuthContext.Provider value={{ user, loading, error, setUser: updateUser, reload }}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const value = useContext(AuthContext);
  if (!value) throw new Error("Brak AuthProvider");
  return value;
}
