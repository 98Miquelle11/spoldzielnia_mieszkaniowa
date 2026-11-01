import type { Metadata } from "next";
import { AdminEvents } from "@/components/admin-events";

export const metadata: Metadata = { title: "Zarządzanie wydarzeniami" };

export default function AdminEventsPage() {
  return <AdminEvents />;
}
