import { useState, useEffect } from "react";
import { Button } from "@/components/ui/button";
import { Card, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Skeleton } from "@/components/ui/skeleton";
import { Plus, AlertCircle, Check, Loader2 } from "lucide-react";
import { Product, AddToCartRequest } from "@/types/domain";

export default function ShopPage() {
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [addingToCart, setAddingToCart] = useState<string | null>(null); // Track which product is being added
  const [cartSuccess, setCartSuccess] = useState<string | null>(null); // Track successful additions

  // Fetch products from Shopping Service
  useEffect(() => {
    const fetchProducts = async () => {
      try {
        setLoading(true);
        const response = await fetch('/api/products');

        if (!response.ok) {
          throw new Error(`HTTP ${response.status}: ${response.statusText}`);
        }

        const data = await response.json();
        setProducts(data);
        setError(null);
      } catch (err) {
        console.error('Failed to fetch products:', err);
        setError(err instanceof Error ? err.message : 'Produkte konnten nicht geladen werden');
      } finally {
        setLoading(false);
      }
    };

    fetchProducts();
  }, []);

  const handleAddToCart = async (product: Product) => {
    setAddingToCart(product.id);

    try {
      const response = await fetch('/api/cart', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify({
          product_id: product.id,
          qty: 1
        } as AddToCartRequest)
      });

      if (!response.ok) {
        throw new Error('Artikel konnte nicht zum Warenkorb hinzugefügt werden');
      }

      // Show success state
      setAddingToCart(null);
      setCartSuccess(product.id);

      // Clear success state after animation
      setTimeout(() => {
        setCartSuccess(null);
      }, 1500);

      console.log('Added to cart:', product.name);

      // Trigger cart count refresh in header
      window.dispatchEvent(new CustomEvent('cartUpdated'));
    } catch (err) {
      console.error('Cart error:', err);
      setAddingToCart(null);
      alert('Artikel konnte nicht zum Warenkorb hinzugefügt werden. Bitte versuche es erneut.');
    }
  };

  // Loading skeleton component
  const ProductSkeleton = () => (
    <Card className="overflow-hidden">
      <div className="aspect-[4/3] bg-muted">
        <Skeleton className="h-full w-full" />
      </div>
      <CardHeader className="pb-3">
        <Skeleton className="h-6 w-3/4" />
        <Skeleton className="h-4 w-full" />
        <Skeleton className="h-4 w-2/3" />
      </CardHeader>
      <CardFooter className="pt-0">
        <div className="flex items-center justify-between w-full">
          <Skeleton className="h-8 w-20" />
          <Skeleton className="h-9 w-24" />
        </div>
      </CardFooter>
    </Card>
  );

  // Product card component
  const ProductCard = ({ product }: { product: Product }) => (
    <Card className="group overflow-hidden transition-all hover:shadow-lg">
      <div className="relative overflow-hidden">
        <div className="aspect-[4/3] bg-muted flex items-center justify-center">
          <span className="text-4xl">🍽️</span>
        </div>
        {/* Mark popular items */}
        {(product.name.includes('Big Data') || product.name.includes('Pizza')) && (
          <Badge className="absolute top-3 left-3 bg-orange-500 hover:bg-orange-600">
            Beliebt
          </Badge>
        )}
      </div>
      <CardHeader className="pb-3">
        <CardTitle className="text-lg leading-tight">{product.name}</CardTitle>
        <CardDescription className="text-sm leading-relaxed">
          {product.description}
        </CardDescription>
      </CardHeader>
      <CardFooter className="pt-0">
        <div className="flex items-center justify-between w-full">
          <span className="text-2xl font-bold text-primary">
            €{product.price.toFixed(2)}
          </span>
          <Button
            onClick={() => handleAddToCart(product)}
            disabled={addingToCart === product.id || cartSuccess === product.id}
            className="gap-2 cursor-pointer disabled:cursor-not-allowed"
            variant={cartSuccess === product.id ? "default" : "default"}
          >
            {addingToCart === product.id ? (
              <>
                <Loader2 className="h-4 w-4 animate-spin" />
                Wird hinzugefügt...
              </>
            ) : cartSuccess === product.id ? (
              <>
                <Check className="h-4 w-4" />
                Hinzugefügt!
              </>
            ) : (
              <>
                <Plus className="h-4 w-4" />
                Hinzufügen
              </>
            )}
          </Button>
        </div>
      </CardFooter>
    </Card>
  );


  if (error) {
    return (
      <div className="space-y-8">
        <div className="space-y-4">
          <h1 className="text-4xl font-bold tracking-tight">Unser Menü</h1>
          <p className="text-xl text-muted-foreground">
            Entdecke unsere datengetriebene Küche
          </p>
        </div>

        <div className="flex flex-col items-center justify-center py-12 space-y-4">
          <AlertCircle className="h-12 w-12 text-muted-foreground" />
          <h3 className="text-lg font-semibold">Fehler beim Laden des Menüs</h3>
          <p className="text-sm text-muted-foreground text-center max-w-md">
            {error}
          </p>
          <Button onClick={() => window.location.reload()} variant="outline">
            Erneut versuchen
          </Button>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-8">
      {/* Header */}
      <div className="space-y-4">
        <h1 className="text-4xl font-bold tracking-tight">Unser Menü</h1>
        <p className="text-xl text-muted-foreground">
          Entdecke unsere datengetriebene Küche
        </p>
      </div>

      {/* Products Grid */}
      {loading ? (
        <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
          {Array.from({ length: 6 }).map((_, i) => (
            <ProductSkeleton key={i} />
          ))}
        </div>
      ) : products.length === 0 ? (
        <div className="text-center py-12">
          <h3 className="text-lg font-semibold mb-2">Keine Produkte verfügbar</h3>
          <p className="text-muted-foreground">
            Unser Menü wird gerade zusammengestellt. Bitte versuchen Sie es später noch einmal.
          </p>
        </div>
      ) : (
        <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
          {products.map((product) => (
            <ProductCard key={product.id} product={product} />
          ))}
        </div>
      )}
    </div>
  );
}