import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Badge } from "@/components/ui/badge";
import { AlertCircle } from "lucide-react";
import type { LoginRequest, RegisterRequest } from "@/types/domain";

export default function AuthPage() {
  const [isLogin, setIsLogin] = useState(true);
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);
  const [checkingAuth, setCheckingAuth] = useState(true);
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const navigate = useNavigate();

  // Check if user is already authenticated
  useEffect(() => {
    checkAuthStatus();
  }, []);

  const checkAuthStatus = async () => {
    try {
      const response = await fetch("/api/cart", {
        credentials: "include",
      });

      if (response.ok) {
        setIsAuthenticated(true);
      }
    } catch (err) {
      // User not authenticated, that's fine
      console.warn("Not authenticated:", err);
    } finally {
      setCheckingAuth(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);
    setSuccess(null);

    try {
      const endpoint = isLogin ? "/api/login" : "/api/register";
      const body: LoginRequest | RegisterRequest = isLogin
        ? { email, password }
        : { email, password, role: "user" };

      const response = await fetch(endpoint, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include", // Important: include cookies
        body: JSON.stringify(body),
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.error || "Anmeldung fehlgeschlagen");
      }

      if (isLogin) {
        // Login successful - cookie is set by server, just redirect
        localStorage.setItem("user_email", email);
        setSuccess("Anmeldung erfolgreich! Weiterleitung...");

        setTimeout(() => {
          navigate("/shop");
        }, 1000);
      } else {
        // Registration successful
        setSuccess("Registrierung erfolgreich! Bitte melde dich an.");
        setIsLogin(true);
        setPassword("");
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : "Etwas ist schiefgelaufen");
    } finally {
      setLoading(false);
    }
  };

  const handleLogout = async () => {
    try {
      await fetch("/api/logout", {
        method: "POST",
        credentials: "include",
      });
      localStorage.removeItem("user_email");
      setIsAuthenticated(false);
    } catch (err) {
      console.error("Logout failed:", err);
    }
  };

  if (checkingAuth) {
    return (
      <div className="flex min-h-[400px] items-center justify-center">
        <p>Überprüfe Anmeldung...</p>
      </div>
    );
  }

  if (isAuthenticated) {
    const userEmail = localStorage.getItem("user_email");
    return (
      <div className="flex min-h-[400px] items-center justify-center">
        <Card className="w-full max-w-sm">
          <CardHeader>
            <CardTitle>Bereits angemeldet</CardTitle>
            <CardDescription>Du bist angemeldet als {userEmail}</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <Button onClick={() => navigate("/shop")} className="w-full">
              Zum Menü
            </Button>
            <Button onClick={handleLogout} variant="outline" className="w-full">
              Abmelden
            </Button>
          </CardContent>
        </Card>
      </div>
    );
  }

  return (
    <div className="flex min-h-[600px] items-center justify-center">
      <div className="w-full max-w-sm">
        <Card>
          <CardHeader>
            <CardTitle>
              {isLogin ? "Bei Analytica Restaurant anmelden" : "Konto erstellen"}
            </CardTitle>
            <CardDescription>
              {isLogin
                ? "Gib deine Anmeldedaten ein, um auf dein Konto zuzugreifen"
                : "Registriere dich, um leckeres Essen zu bestellen"}
            </CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleSubmit}>
              <div className="flex flex-col gap-6">
                {error && (
                  <div className="flex items-center gap-2 rounded-md bg-destructive/10 p-3 text-sm text-destructive">
                    <AlertCircle className="h-4 w-4" />
                    {error}
                  </div>
                )}

                {success && (
                  <div className="rounded-md bg-green-50 p-3 text-sm text-green-700">
                    {success}
                  </div>
                )}

                <div className="grid gap-3">
                  <Label htmlFor="email">Email</Label>
                  <Input
                    id="email"
                    type="email"
                    placeholder="dein@email.de"
                    value={email}
                    onChange={(e) => setEmail(e.target.value)}
                    required
                  />
                </div>

                <div className="grid gap-3">
                  <Label htmlFor="password">Passwort</Label>
                  <Input
                    id="password"
                    type="password"
                    placeholder="••••••••"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    required
                  />
                </div>

                <Button type="submit" className="w-full" disabled={loading}>
                  {loading
                    ? "Lädt..."
                    : isLogin
                      ? "Anmelden"
                      : "Konto erstellen"}
                </Button>

                <div className="text-center text-sm">
                  {isLogin
                    ? "Noch kein Konto? "
                    : "Bereits ein Konto? "}
                  <button
                    type="button"
                    onClick={() => {
                      setIsLogin(!isLogin);
                      setError(null);
                      setSuccess(null);
                      setPassword("");
                    }}
                    className="underline underline-offset-4 hover:text-primary"
                  >
                    {isLogin ? "Registrieren" : "Anmelden"}
                  </button>
                </div>
              </div>
            </form>
          </CardContent>
        </Card>

        {/* Demo Accounts */}
        <div className="mt-4 text-center">
          <p className="text-xs text-muted-foreground mb-2">Demo Accounts:</p>
          <div className="flex gap-2 justify-center">
            <Badge variant="outline" className="text-xs">
              user@demo.de / demo123
            </Badge>
          </div>
        </div>
      </div>
    </div>
  );
}
