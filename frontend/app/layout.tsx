import type { Metadata } from "next";
import "bootstrap/dist/css/bootstrap.min.css";
import "./globals.css";
import { AuthProvider } from "@/components/auth";
import { Header } from "@/components/header";

export const metadata: Metadata = {
  title: { default: "Strona główna — Spółdzielnia Mieszkaniowa „Mirmiłowo”", template: "%s — Spółdzielnia Mieszkaniowa „Mirmiłowo”" },
  description: "Spółdzielnia Mieszkaniowa Mirmiłowo — informacje dla mieszkańców i aktualności osiedla.",
};

export default function RootLayout({ children }: Readonly<{ children: React.ReactNode }>) {
  return <html lang="pl"><body>
    <a href="#tresc" className="visually-hidden-focusable skip-link">Przejdź do treści</a>
    <AuthProvider><Header /><main id="tresc">{children}</main></AuthProvider>
    <footer className="footer border-top mt-5 py-4"><div className="container d-flex flex-column flex-md-row justify-content-between gap-2 small">
      <span>© {new Date().getFullYear()} Spółdzielnia Mieszkaniowa „Mirmiłowo”.</span>
      <span>Projekt edukacyjny • Wersja startowa • Dane fikcyjne</span>
    </div></footer>
  </body></html>;
}
