"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { authRequest, useAuth } from "./auth";

export function LogoutButton({ className = "btn btn-primary" }: { className?: string }) {
  const { setUser } = useAuth();
  const router = useRouter();
  const [pending, setPending] = useState(false);
  const [error, setError] = useState("");
  async function logout() {
    setPending(true); setError("");
    try {
      await authRequest("logout"); setUser(null); router.push("/"); router.refresh();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Nie udało się wylogować.");
    } finally { setPending(false); }
  }
  return <>
    <button type="button" className={className} disabled={pending} onClick={logout}>{pending ? "Wylogowywanie…" : "Wyloguj się"}</button>
    {error && <p role="alert" className="small text-danger mb-0">{error}</p>}
  </>;
}
