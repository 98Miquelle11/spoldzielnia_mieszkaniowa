import type { Metadata } from "next";
import { LoginForm } from "@/components/login-form";

export const metadata: Metadata = { title: "Logowanie" };
export default function Login() {
  return <div className="container py-5"><div className="row justify-content-center"><div className="col-md-6"><div className="card shadow-sm"><div className="card-body p-4">
    <h1 className="h2 card-title mb-4">Logowanie</h1><LoginForm />
  </div></div></div></div></div>;
}
