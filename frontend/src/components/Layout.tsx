import { Link, useLocation } from "react-router-dom";
import { useState, useEffect } from "react";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import type { Cart } from "@/types/domain";
import {
  NavigationMenu,
  NavigationMenuItem,
  NavigationMenuList,
  navigationMenuTriggerStyle,
} from "@/components/ui/navigation-menu";
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from "@/components/ui/sheet";
import { ShoppingCart, User, Menu, Clock } from "lucide-react";
import { cn } from "@/lib/utils";

interface LayoutProps {
  children: React.ReactNode;
}

export default function Layout({ children }: LayoutProps) {
  const location = useLocation();
  const isHomePage = location.pathname === "/";
  const [cartCount, setCartCount] = useState(0);
  const [cartAnimating, setCartAnimating] = useState(false);

  useEffect(() => {
    fetchCartCount();

    // Listen for cart updates
    const handleCartUpdate = () => {
      setCartAnimating(true);
      fetchCartCount();

      // Clear animation after a short delay
      setTimeout(() => {
        setCartAnimating(false);
      }, 600);
    };

    window.addEventListener('cartUpdated', handleCartUpdate);

    return () => {
      window.removeEventListener('cartUpdated', handleCartUpdate);
    };
  }, []);

  const fetchCartCount = async () => {
    try {
      const response = await fetch("/api/cart", {
        credentials: "include",
      });

      if (response.ok) {
        const cartData: Cart = await response.json();
        const totalItems = cartData.Items.reduce(
          (sum, item) => sum + item.Qty,
          0
        );
        setCartCount(totalItems);
      }
    } catch (err) {
      console.error("Failed to fetch cart count:", err);
    }
  };

  return (
    <div className="min-h-screen bg-background flex flex-col">
      <header className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
        <div className="w-full max-w-none px-4 sm:px-6 lg:px-8 flex h-16 items-center relative">
          {/* Logo - Left */}
          <div className="flex items-center space-x-2">
            <Link to="/" className="flex items-center space-x-2">
              <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-primary text-primary-foreground">
                <span className="text-sm font-bold">A</span>
              </div>
              <span className="text-xl font-bold text-foreground">
                Analytica Restaurant
              </span>
            </Link>
          </div>

          {/* Navigation - Absolute Center */}
          <div className="absolute left-1/2 transform -translate-x-1/2">
            <NavigationMenu className="hidden md:flex">
              <NavigationMenuList>
                <NavigationMenuItem>
                  <Link to="/shop" className={navigationMenuTriggerStyle()}>
                    Menü
                  </Link>
                </NavigationMenuItem>
                <NavigationMenuItem>
                  <Link to="/checkout" className={navigationMenuTriggerStyle()}>
                    Bestellung verfolgen
                  </Link>
                </NavigationMenuItem>
              </NavigationMenuList>
            </NavigationMenu>
          </div>

          {/* Actions - Right */}
          <div className="ml-auto flex items-center space-x-3">
            <Button variant="ghost" size="icon" asChild className="relative">
              <Link to="/checkout">
                <ShoppingCart className="h-5 w-5" />
                {cartCount > 0 && (
                  <Badge
                    variant="destructive"
                    className={cn(
                      "absolute -top-2 -right-2 h-5 w-5 rounded-full p-0 flex items-center justify-center text-xs transition-transform",
                      cartAnimating && "animate-bounce"
                    )}
                  >
                    {cartCount}
                  </Badge>
                )}
              </Link>
            </Button>

            <Button variant="ghost" size="icon" asChild>
              <Link to="/auth">
                <User className="h-5 w-5" />
              </Link>
            </Button>

            <Sheet>
              <SheetTrigger asChild>
                <Button variant="ghost" size="icon" className="md:hidden">
                  <Menu className="h-5 w-5" />
                </Button>
              </SheetTrigger>
              <SheetContent>
                <SheetHeader>
                  <SheetTitle>Navigation</SheetTitle>
                  <SheetDescription>
                    Navigiere durch unser Restaurant
                  </SheetDescription>
                </SheetHeader>
                <nav className="flex flex-col space-y-4 mt-6">
                  <Link
                    to="/shop"
                    className="flex items-center space-x-2 text-lg font-medium"
                  >
                    <ShoppingCart className="h-5 w-5" />
                    <span>Menü</span>
                  </Link>
                  <Link
                    to="/checkout"
                    className="flex items-center space-x-2 text-lg font-medium"
                  >
                    <Clock className="h-5 w-5" />
                    <span>Bestellung verfolgen</span>
                  </Link>
                  <Link
                    to="/auth"
                    className="flex items-center space-x-2 text-lg font-medium"
                  >
                    <User className="h-5 w-5" />
                    <span>Mein Konto</span>
                  </Link>
                </nav>
              </SheetContent>
            </Sheet>
          </div>
        </div>
      </header>

      <main
        className={cn(
          "flex-1", // Takes available space, pushing footer to bottom
          isHomePage
            ? ""
            : "container mx-auto px-4 sm:px-6 lg:px-8 py-8 lg:py-12"
        )}
      >
        {children}
      </main>

      {!isHomePage && (
        <footer className="border-t bg-muted/30 mt-auto">
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
      )}
    </div>
  );
}
