import { Link } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { ChefHat, ShoppingCart, Clock, MapPin, Truck } from "lucide-react";

export default function HomePage() {
  return (
    <>
      {/* Hero Section */}
      <div className="relative overflow-hidden bg-gradient-to-br from-background via-background to-muted/20">
        <div className="container mx-auto px-4 sm:px-6 lg:px-8 py-20 sm:py-32 lg:py-40">
          <div className="mx-auto max-w-2xl text-center">
            <h1 className="text-4xl font-bold tracking-tight text-foreground sm:text-6xl">
              Analytica Restaurant
            </h1>
            <p className="mt-6 text-lg leading-8 text-muted-foreground">
              Datengetriebene Küche trifft auf authentischen Geschmack. Frische
              Zutaten, optimierte Rezepte - direkt zu dir nach Hause geliefert.
            </p>
            <div className="mt-10 flex items-center justify-center gap-x-6">
              <Button asChild size="lg">
                <Link to="/shop">
                  <ShoppingCart className="mr-2 h-4 w-4" />
                  Jetzt Bestellen
                </Link>
              </Button>
              <Button variant="outline" asChild size="lg">
                <Link to="/auth">Mein Konto</Link>
              </Button>
            </div>
          </div>
        </div>
      </div>

      {/* Features Section */}
      <div className="container mx-auto px-4 sm:px-6 lg:px-8 py-16 sm:py-24 lg:py-32">
        <div className="mx-auto max-w-3xl text-center">
          <h2 className="text-3xl font-bold tracking-tight text-foreground sm:text-4xl">
            Warum Analytica Restaurant?
          </h2>
          <p className="mt-4 text-lg text-muted-foreground">
            Qualität, Schnelligkeit und Geschmack in perfekter Balance
          </p>
        </div>

        <div className="mx-auto mt-16 grid max-w-5xl grid-cols-1 gap-8 sm:grid-cols-2 lg:grid-cols-3">
          {/* Fresh Ingredients */}
          <div className="group relative overflow-hidden rounded-lg border bg-card p-8 shadow-sm transition-all hover:shadow-md">
            <div className="flex h-12 w-12 items-center justify-center rounded-lg bg-green-500/10">
              <ChefHat className="h-6 w-6 text-green-500" />
            </div>
            <h3 className="mt-4 text-xl font-semibold text-card-foreground">
              Frische Zutaten
            </h3>
            <p className="mt-2 text-sm text-muted-foreground">
              Täglich frische, regionale Zutaten für optimalen Geschmack und
              Qualität.
            </p>
          </div>

          {/* Fast Delivery */}
          <div className="group relative overflow-hidden rounded-lg border bg-card p-8 shadow-sm transition-all hover:shadow-md">
            <div className="flex h-12 w-12 items-center justify-center rounded-lg bg-orange-500/10">
              <Truck className="h-6 w-6 text-orange-500" />
            </div>
            <h3 className="mt-4 text-xl font-semibold text-card-foreground">
              Schnelle Lieferung
            </h3>
            <p className="mt-2 text-sm text-muted-foreground">
              Durchschnittlich 30 Minuten - heiß und frisch direkt vor deine
              Haustür.
            </p>
          </div>

          {/* Real-time Tracking */}
          <div className="group relative overflow-hidden rounded-lg border bg-card p-8 shadow-sm transition-all hover:shadow-md">
            <div className="flex h-12 w-12 items-center justify-center rounded-lg bg-blue-500/10">
              <MapPin className="h-6 w-6 text-blue-500" />
            </div>
            <h3 className="mt-4 text-xl font-semibold text-card-foreground">
              Live-Tracking
            </h3>
            <p className="mt-2 text-sm text-muted-foreground">
              Verfolge deine Bestellung in Echtzeit - von der Küche bis vor
              deine Tür.
            </p>
          </div>
        </div>
      </div>

      {/* Action Section */}
      <div className="container mx-auto px-4 sm:px-6 lg:px-8 py-16 sm:py-24">
        <div className="mx-auto max-w-3xl text-center">
          <h2 className="text-3xl font-bold tracking-tight text-foreground sm:text-4xl">
            Bereit für deine Bestellung?
          </h2>
          <p className="mt-4 text-lg text-muted-foreground">
            Entdecke unser vielfältiges Menü und lass dich verwöhnen
          </p>
          <div className="mt-8 flex items-center justify-center gap-x-6">
            <Button asChild size="lg">
              <Link to="/shop">
                <ShoppingCart className="mr-2 h-4 w-4" />
                Unser Menü entdecken
              </Link>
            </Button>
            <Button variant="ghost" asChild size="lg">
              <Link to="/checkout" className="flex items-center">
                <Clock className="mr-2 h-4 w-4" />
                Bestellung verfolgen
              </Link>
            </Button>
          </div>
        </div>
      </div>

      {/* Footer für HomePage */}
      <footer className="border-t bg-muted/30 mt-16">
        <div className="container mx-auto px-4 sm:px-6 lg:px-8 py-8 lg:py-12">
          <div className="text-center">
            <p className="text-sm text-muted-foreground">
              © 2025 Analytica Restaurant - Cloud-Native Delivery Service
            </p>
            <p className="mt-1 text-xs text-muted-foreground">
              Ein THI Cloud-Native Projekt
            </p>
          </div>
        </div>
      </footer>
    </>
  );
}
