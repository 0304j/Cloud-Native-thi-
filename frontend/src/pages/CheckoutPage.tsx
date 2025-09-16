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
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group";
import { Separator } from "@/components/ui/separator";
import { Badge } from "@/components/ui/badge";
import { AlertCircle, ShoppingCart, Truck, MapPin, Wand2, Plus, Minus, X } from "lucide-react";
import { fakerDE as faker } from "@faker-js/faker";
import type {
  Cart,
  CheckoutItem,
  DeliveryInfo,
  CheckoutRequest,
  OrderResponse,
  Product,
} from "@/types/domain";

export default function CheckoutPage() {
  const [cart, setCart] = useState<CheckoutItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [orderType, setOrderType] = useState<"delivery" | "pickup">("delivery");
  const [updatingItem, setUpdatingItem] = useState<string | null>(null); // Track which item is being updated
  const [deliveryInfo, setDeliveryInfo] = useState<DeliveryInfo>({
    customer_name: "",
    customer_phone: "",
    street: "",
    house_number: "",
    postal_code: "",
    city: "",
    floor: "",
    instructions: "",
  });

  const navigate = useNavigate();

  useEffect(() => {
    fetchCart();
  }, []);

  const fetchCart = async () => {
    try {
      setLoading(true);
      const response = await fetch("/api/cart", {
        credentials: "include",
      });

      if (!response.ok) {
        if (response.status === 401) {
          // Don't navigate during cart updates, just show error
          setError("Sitzung abgelaufen. Bitte lade die Seite neu und melde dich erneut an.");
          return;
        }
        throw new Error("Warenkorb konnte nicht geladen werden");
      }

      const cartData: Cart = await response.json();

      // Convert cart items to display format
      const cartItems: CheckoutItem[] = [];
      for (const item of cartData.Items) {
        // Fetch product details for each cart item
        try {
          const productResponse = await fetch("/api/products");
          const products: Product[] = await productResponse.json();
          const product = products.find(
            (p: Product) => p.id === item.ProductID
          );

          if (product) {
            cartItems.push({
              product_id: item.ProductID,
              product_name: product.name,
              quantity: item.Qty,
              unit_price: product.price,
              total_price: product.price * item.Qty,
            });
          }
        } catch (err) {
          console.error("Failed to fetch product details:", err);
        }
      }

      setCart(cartItems);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Warenkorb konnte nicht geladen werden");
    } finally {
      setLoading(false);
    }
  };

  const getTotalAmount = () => {
    return cart.reduce((sum, item) => sum + item.total_price, 0);
  };

  const updateCartItemQuantity = async (productId: string, newQuantity: number) => {
    setUpdatingItem(productId);
    try {
      const response = await fetch('/api/cart', {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify({
          product_id: productId,
          qty: newQuantity
        })
      });

      if (!response.ok) {
        throw new Error('Failed to update cart item');
      }

      // Refresh cart data
      await fetchCart();

      // Trigger cart count refresh in header
      window.dispatchEvent(new CustomEvent('cartUpdated'));
    } catch (err) {
      console.error('Failed to update cart item:', err);
      setError('Failed to update item quantity');
    } finally {
      setUpdatingItem(null);
    }
  };

  const removeCartItem = async (productId: string) => {
    setUpdatingItem(productId);
    try {
      const response = await fetch(`/api/cart/${productId}`, {
        method: 'DELETE',
        credentials: 'include',
      });

      if (!response.ok) {
        throw new Error('Failed to remove cart item');
      }

      // Refresh cart data
      await fetchCart();

      // Trigger cart count refresh in header
      window.dispatchEvent(new CustomEvent('cartUpdated'));
    } catch (err) {
      console.error('Failed to remove cart item:', err);
      setError('Failed to remove item');
    } finally {
      setUpdatingItem(null);
    }
  };

  const generateTestAddress = () => {
    // Set German locale for realistic German data

    setDeliveryInfo({
      customer_name: faker.person.fullName(),
      customer_phone: faker.phone.number({ style: "international" }),
      street: faker.location.street(),
      house_number: faker.number.int({ min: 1, max: 299 }).toString(),
      postal_code: faker.location.zipCode("#####"),
      city: faker.helpers.arrayElement([
        "München",
        "Berlin",
        "Hamburg",
        "Köln",
        "Frankfurt",
      ]),
      floor: faker.helpers.maybe(
        () => faker.number.int({ min: 1, max: 5 }) + ". OG",
        { probability: 0.6 }
      ),
      instructions: faker.helpers.maybe(
        () =>
          faker.helpers.arrayElement([
            "Klingel 2x drücken",
            "Hintereingang nutzen",
            "Bei Nachbar abgeben falls nicht da",
            "Vor der Tür abstellen",
            "Paket beim Pförtner abgeben",
          ]),
        { probability: 0.4 }
      ),
    });
  };

  const handleSubmitOrder = async (e: React.FormEvent) => {
    e.preventDefault();

    if (cart.length === 0) {
      setError("Dein Warenkorb ist leer");
      return;
    }

    // Validate delivery info if delivery is selected
    if (orderType === "delivery") {
      const required = [
        "customer_name",
        "customer_phone",
        "street",
        "house_number",
        "postal_code",
        "city",
      ];
      for (const field of required) {
        if (!deliveryInfo[field as keyof DeliveryInfo]) {
          setError(`Please fill in ${field.replace("_", " ")}`);
          return;
        }
      }
    }

    setSubmitting(true);
    setError(null);

    try {
      const checkoutRequest: CheckoutRequest = {
        user_id: "mock-user-id", // TODO: Get from JWT token
        items: cart,
        total_amount: getTotalAmount(),
        currency: "EUR",
        order_type: orderType,
        delivery_info: orderType === "delivery" ? deliveryInfo : undefined,
      };

      const response = await fetch("/api/checkout", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify(checkoutRequest),
      });

      // Debug: Log response details
      console.log("Checkout response status:", response.status);
      console.log("Checkout response headers:", Object.fromEntries(response.headers.entries()));

      if (!response.ok) {
        const responseText = await response.text();
        console.error("Checkout error response:", responseText);

        // Try to parse as JSON, fallback to text
        let errorMessage = "Bestellung fehlgeschlagen";
        try {
          const errorData = JSON.parse(responseText);
          errorMessage = errorData.error || errorMessage;
        } catch {
          // Response is not JSON, use the text directly
          errorMessage = responseText || errorMessage;
        }

        throw new Error(errorMessage);
      }

      const responseText = await response.text();
      console.log("Checkout success response:", responseText);

      // Parse JSON response
      const result: OrderResponse = JSON.parse(responseText);

      // Order successful - redirect to payment
      navigate(`/payment?order_id=${result.order_id}&amount=${getTotalAmount()}&currency=EUR`);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Bestellung fehlgeschlagen");
    } finally {
      setSubmitting(false);
    }
  };

  if (loading) {
    return (
      <div className="max-w-4xl mx-auto space-y-8">
        <div className="text-center">
          <h1 className="text-3xl font-bold">Warenkorb wird geladen...</h1>
        </div>
      </div>
    );
  }

  if (cart.length === 0) {
    return (
      <div className="max-w-4xl mx-auto space-y-8">
        <div className="text-center space-y-4">
          <ShoppingCart className="h-16 w-16 mx-auto text-muted-foreground" />
          <h1 className="text-3xl font-bold">Dein Warenkorb ist leer</h1>
          <p className="text-muted-foreground">
            Füge leckere Gerichte aus unserem Menü hinzu!
          </p>
          <Button onClick={() => navigate("/shop")} size="lg">
            Menü durchstöbern
          </Button>
        </div>
      </div>
    );
  }

  return (
    <div className="max-w-4xl mx-auto space-y-8">
      <div>
        <h1 className="text-3xl font-bold">Bestellung</h1>
        <p className="text-muted-foreground">
          Überprüfe deine Bestellung und schließe deinen Kauf ab
        </p>
      </div>

      <form onSubmit={handleSubmitOrder} className="grid gap-8 lg:grid-cols-2">
        {/* Left Column - Order Details */}
        <div className="space-y-6">
          {/* Cart Summary */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <ShoppingCart className="h-5 w-5" />
                Deine Bestellung
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              {cart.map((item) => (
                <div
                  key={item.product_id}
                  className="flex items-center gap-4 py-2"
                >
                  <div className="flex-1">
                    <p className="font-medium">{item.product_name}</p>
                    <p className="text-sm text-muted-foreground">
                      €{item.unit_price.toFixed(2)} each
                    </p>
                  </div>

                  {/* Quantity Controls */}
                  <div className="flex items-center gap-2">
                    <Button
                      type="button"
                      variant="outline"
                      size="sm"
                      onClick={() => updateCartItemQuantity(item.product_id, Math.max(1, item.quantity - 1))}
                      disabled={updatingItem === item.product_id || item.quantity <= 1}
                      className="h-8 w-8 p-0"
                    >
                      <Minus className="h-3 w-3" />
                    </Button>

                    <span className="w-8 text-center font-medium">
                      {item.quantity}
                    </span>

                    <Button
                      type="button"
                      variant="outline"
                      size="sm"
                      onClick={() => updateCartItemQuantity(item.product_id, item.quantity + 1)}
                      disabled={updatingItem === item.product_id}
                      className="h-8 w-8 p-0"
                    >
                      <Plus className="h-3 w-3" />
                    </Button>
                  </div>

                  {/* Price */}
                  <div className="text-right">
                    <p className="font-medium">€{item.total_price.toFixed(2)}</p>
                  </div>

                  {/* Remove Button */}
                  <Button
                    type="button"
                    variant="ghost"
                    size="sm"
                    onClick={() => removeCartItem(item.product_id)}
                    disabled={updatingItem === item.product_id}
                    className="h-8 w-8 p-0 text-muted-foreground hover:text-destructive"
                  >
                    <X className="h-4 w-4" />
                  </Button>
                </div>
              ))}

              <Separator />

              <div className="flex justify-between items-center text-lg font-bold">
                <span>Gesamt</span>
                <span>€{getTotalAmount().toFixed(2)}</span>
              </div>
            </CardContent>
          </Card>

          {/* Order Type Selection */}
          <Card>
            <CardHeader>
              <CardTitle>Bestellart</CardTitle>
              <CardDescription>
                Wie möchtest du deine Bestellung erhalten?
              </CardDescription>
            </CardHeader>
            <CardContent>
              <RadioGroup
                value={orderType}
                onValueChange={(value) =>
                  setOrderType(value as "delivery" | "pickup")
                }
              >
                <div className="flex items-center space-x-2">
                  <RadioGroupItem value="delivery" id="delivery" />
                  <Label
                    htmlFor="delivery"
                    className="flex items-center gap-2 cursor-pointer"
                  >
                    <Truck className="h-4 w-4" />
                    Lieferung
                    <Badge variant="secondary">€2.50</Badge>
                  </Label>
                </div>
                <div className="flex items-center space-x-2">
                  <RadioGroupItem value="pickup" id="pickup" />
                  <Label
                    htmlFor="pickup"
                    className="flex items-center gap-2 cursor-pointer"
                  >
                    <MapPin className="h-4 w-4" />
                    Abholung
                    <Badge variant="outline">Kostenlos</Badge>
                  </Label>
                </div>
              </RadioGroup>
            </CardContent>
          </Card>
        </div>

        {/* Right Column - Delivery Information */}
        <div className="space-y-6">
          {orderType === "delivery" && (
            <Card>
              <CardHeader>
                <CardTitle>Lieferinformationen</CardTitle>
                <CardDescription>
                  Wohin sollen wir deine Bestellung liefern?
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="grid gap-4 sm:grid-cols-2">
                  <div className="space-y-2">
                    <Label htmlFor="customer_name">Vollständiger Name *</Label>
                    <Input
                      id="customer_name"
                      value={deliveryInfo.customer_name}
                      onChange={(e) =>
                        setDeliveryInfo({
                          ...deliveryInfo,
                          customer_name: e.target.value,
                        })
                      }
                      placeholder="Max Mustermann"
                      required={orderType === "delivery"}
                    />
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="customer_phone">Telefonnummer *</Label>
                    <Input
                      id="customer_phone"
                      type="tel"
                      value={deliveryInfo.customer_phone}
                      onChange={(e) =>
                        setDeliveryInfo({
                          ...deliveryInfo,
                          customer_phone: e.target.value,
                        })
                      }
                      placeholder="+49 123 456789"
                      required={orderType === "delivery"}
                    />
                  </div>
                </div>

                <div className="grid gap-4 sm:grid-cols-3">
                  <div className="space-y-2 sm:col-span-2">
                    <Label htmlFor="street">Straße *</Label>
                    <Input
                      id="street"
                      value={deliveryInfo.street}
                      onChange={(e) =>
                        setDeliveryInfo({
                          ...deliveryInfo,
                          street: e.target.value,
                        })
                      }
                      placeholder="Hauptstraße"
                      required={orderType === "delivery"}
                    />
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="house_number">Hausnummer *</Label>
                    <Input
                      id="house_number"
                      value={deliveryInfo.house_number}
                      onChange={(e) =>
                        setDeliveryInfo({
                          ...deliveryInfo,
                          house_number: e.target.value,
                        })
                      }
                      placeholder="123"
                      required={orderType === "delivery"}
                    />
                  </div>
                </div>

                <div className="grid gap-4 sm:grid-cols-2">
                  <div className="space-y-2">
                    <Label htmlFor="postal_code">Postleitzahl *</Label>
                    <Input
                      id="postal_code"
                      value={deliveryInfo.postal_code}
                      onChange={(e) =>
                        setDeliveryInfo({
                          ...deliveryInfo,
                          postal_code: e.target.value,
                        })
                      }
                      placeholder="12345"
                      required={orderType === "delivery"}
                    />
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="city">Stadt *</Label>
                    <Input
                      id="city"
                      value={deliveryInfo.city}
                      onChange={(e) =>
                        setDeliveryInfo({
                          ...deliveryInfo,
                          city: e.target.value,
                        })
                      }
                      placeholder="München"
                      required={orderType === "delivery"}
                    />
                  </div>
                </div>

                <div className="space-y-2">
                  <Label htmlFor="floor">Stockwerk/Wohnung (Optional)</Label>
                  <Input
                    id="floor"
                    value={deliveryInfo.floor}
                    onChange={(e) =>
                      setDeliveryInfo({
                        ...deliveryInfo,
                        floor: e.target.value,
                      })
                    }
                    placeholder="2. Stock, Wohnung 4"
                  />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="instructions">
                    Lieferhinweise (Optional)
                  </Label>
                  <Input
                    id="instructions"
                    value={deliveryInfo.instructions}
                    onChange={(e) =>
                      setDeliveryInfo({
                        ...deliveryInfo,
                        instructions: e.target.value,
                      })
                    }
                    placeholder="Zweimal klingeln"
                  />
                </div>

                <Button
                  type="button"
                  variant="outline"
                  size="sm"
                  onClick={generateTestAddress}
                  className="flex items-center gap-2 w-fit"
                >
                  <Wand2 className="h-4 w-4" />
                  Testdaten generieren
                </Button>
              </CardContent>
            </Card>
          )}

          {orderType === "pickup" && (
            <Card>
              <CardHeader>
                <CardTitle>Abholinformationen</CardTitle>
                <CardDescription>Details zur Abholung</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-2">
                  <p className="text-sm text-muted-foreground">
                    <strong>Abholadresse:</strong>
                    <br />
                    Analytica Restaurant
                    <br />
                    Musterstraße 123
                    <br />
                    80333 München
                  </p>
                  <p className="text-sm text-muted-foreground">
                    <strong>Abholzeit:</strong> Bereit in ~30 Minuten
                  </p>
                </div>
              </CardContent>
            </Card>
          )}

          {error && (
            <div className="flex items-center gap-2 rounded-md bg-destructive/10 p-3 text-sm text-destructive">
              <AlertCircle className="h-4 w-4" />
              {error}
            </div>
          )}

          <Button
            type="submit"
            className="w-full"
            size="lg"
            disabled={submitting}
          >
            {submitting
              ? "Bestellung wird verarbeitet..."
              : `Bestellung aufgeben - €${getTotalAmount().toFixed(2)}`}
          </Button>
        </div>
      </form>
    </div>
  );
}
