import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Badge } from "@/components/ui/badge";
import { AlertCircle } from "lucide-react";

export default function AuthPage() {
  const [isLogin, setIsLogin] = useState(true);
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);
    setSuccess(null);

    try {
      const endpoint = isLogin ? '/api/login' : '/api/register';
      const body = isLogin
        ? { email, password }
        : { email, password, role: 'user' };

      const response = await fetch(endpoint, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(body),
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.error || 'Authentication failed');
      }

      if (isLogin) {
        // Login successful - store token and redirect
        localStorage.setItem('jwt_token', data.token);
        localStorage.setItem('user_email', email);
        setSuccess('Login successful! Redirecting...');

        setTimeout(() => {
          navigate('/shop');
        }, 1000);
      } else {
        // Registration successful
        setSuccess('Registration successful! Please login.');
        setIsLogin(true);
        setPassword('');
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Something went wrong');
    } finally {
      setLoading(false);
    }
  };

  const isLoggedIn = localStorage.getItem('jwt_token');

  if (isLoggedIn) {
    const userEmail = localStorage.getItem('user_email');
    return (
      <div className="flex min-h-[400px] items-center justify-center">
        <Card className="w-full max-w-sm">
          <CardHeader>
            <CardTitle>Already Logged In</CardTitle>
            <CardDescription>
              You are logged in as {userEmail}
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <Button
              onClick={() => navigate('/shop')}
              className="w-full"
            >
              Go to Menu
            </Button>
            <Button
              onClick={() => {
                localStorage.removeItem('jwt_token');
                localStorage.removeItem('user_email');
                window.location.reload();
              }}
              variant="outline"
              className="w-full"
            >
              Logout
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
              {isLogin ? 'Login to Analytica Restaurant' : 'Create Account'}
            </CardTitle>
            <CardDescription>
              {isLogin
                ? 'Enter your credentials to access your account'
                : 'Sign up to start ordering delicious food'
              }
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
                  <Label htmlFor="password">Password</Label>
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
                  {loading ? 'Loading...' : (isLogin ? 'Login' : 'Create Account')}
                </Button>

                <div className="text-center text-sm">
                  {isLogin ? "Don't have an account? " : "Already have an account? "}
                  <button
                    type="button"
                    onClick={() => {
                      setIsLogin(!isLogin);
                      setError(null);
                      setSuccess(null);
                      setPassword('');
                    }}
                    className="underline underline-offset-4 hover:text-primary"
                  >
                    {isLogin ? 'Sign up' : 'Login'}
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
            <Badge variant="outline" className="text-xs">user@demo.de / demo123</Badge>
          </div>
        </div>
      </div>
    </div>
  );
}