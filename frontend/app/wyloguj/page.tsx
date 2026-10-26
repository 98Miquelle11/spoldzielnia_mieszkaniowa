import type { Metadata } from "next";
import { LogoutButton } from "@/components/logout-button";

export const metadata: Metadata = { title: "Wylogowanie" };
export default function Logout() {
  return <section className="container py-5"><h1 className="h2">Wylogowanie</h1><p>Zakończ bieżącą sesję przyciskiem poniżej.</p><LogoutButton /></section>;
}
