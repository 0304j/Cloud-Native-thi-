import { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group";
import { Separator } from "@/components/ui/separator";
import { Badge } from "@/components/ui/badge";
import { AlertCircle, ShoppingCart, Truck, MapPin } from "lucide-react";

interface CartItem {
  product_id: string;
  product_name: string;
  quantity: number;
  unit_price: number;
  total_price: number;
}

interface Cart {
  UserID: string;
  Items: Array<{
    ProductID: string;
    Qty: number;
  }>;
}

interface DeliveryInfo {
  customer_name: string;
  customer_phone: string;
  street: string;
  house_number: string;
  postal_code: string;
  city: string;
  floor?: string;
  instructions?: string;
}

export default function CheckoutPage() {
  const [cart, setCart] = useState<CartItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [orderType, setOrderType] = useState<"delivery" | "pickup">("delivery");
  const [deliveryInfo, setDeliveryInfo] = useState<DeliveryInfo>({
    customer_name: "",
    customer_phone: "",
    street: "",
    house_number: "",
    postal_code: "",
    city: "",
    floor: "",
    instructions: ""
  });

  const navigate = useNavigate();

  useEffect(() => {
    fetchCart();
  }, []);

  const fetchCart = async () => {
    try {
      setLoading(true);
      const response = await fetch('/api/cart', {
        credentials: 'include'
      });

      if (!response.ok) {
        if (response.status === 401) {
          navigate('/auth');
          return;
        }
        throw new Error('Failed to fetch cart');
      }

      const cartData: Cart = await response.json();

      // Convert cart items to display format
      const cartItems: CartItem[] = [];
      for (const item of cartData.Items) {
        // Fetch product details for each cart item
        try {
          const productResponse = await fetch('/api/products');
          const products = await productResponse.json();
          const product = products.find((p: any) => p.id === item.ProductID);

          if (product) {
            cartItems.push({
              product_id: item.ProductID,
              product_name: product.name,
              quantity: item.Qty,
              unit_price: product.price,
              total_price: product.price * item.Qty
            });
          }
        } catch (err) {
          console.error('Failed to fetch product details:', err);
        }
      }

      setCart(cartItems);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load cart');
    } finally {
      setLoading(false);
    }
  };

  const getTotalAmount = () => {
    return cart.reduce((sum, item) => sum + item.total_price, 0);
  };

  const handleSubmitOrder = async (e: React.FormEvent) => {
    e.preventDefault();

    if (cart.length === 0) {
      setError('Your cart is empty');
      return;
    }

    // Validate delivery info if delivery is selected
    if (orderType === 'delivery') {
      const required = ['customer_name', 'customer_phone', 'street', 'house_number', 'postal_code', 'city'];
      for (const field of required) {
        if (!deliveryInfo[field as keyof DeliveryInfo]) {
          setError(`Please fill in ${field.replace('_', ' ')}`);
          return;
        }
      }
    }

    setSubmitting(true);
    setError(null);

    try {
      const checkoutRequest = {
        user_id: "mock-user-id", // TODO: Get from JWT token
        items: cart.map(item => ({
          product_id: item.product_id,
          product_name: item.product_name,
          quantity: item.quantity,
          unit_price: item.unit_price,
          total_price: item.total_price
        })),
        total_amount: getTotalAmount(),
        currency: "EUR",
        order_type: orderType,
        delivery_info: orderType === 'delivery' ? deliveryInfo : undefined
      };

      const response = await fetch('/api/checkout', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        credentials: 'include',
        body: JSON.stringify(checkoutRequest)
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || 'Order failed');
      }

      const result = await response.json();

      // Order successful - redirect to tracking
      navigate(`/kitchen?order_id=${result.order_id}`);

    } catch (err) {
      setError(err instanceof Error ? err.message : 'Order failed');
    } finally {
      setSubmitting(false);
    }
  };


  if (loading) {
    return (
      <div className="max-w-4xl mx-auto space-y-8">
        <div className="text-center">
          <h1 className="text-3xl font-bold">Loading Cart...</h1>
        </div>
      </div>
    );
  }

  if (cart.length === 0) {
    return (
      <div className="max-w-4xl mx-auto space-y-8">
        <div className="text-center space-y-4">
          <ShoppingCart className="h-16 w-16 mx-auto text-muted-foreground" />
          <h1 className="text-3xl font-bold">Your Cart is Empty</h1>
          <p className="text-muted-foreground">Add some delicious items from our menu!</p>
          <Button onClick={() => navigate('/shop')} size="lg">
            Browse Menu
          </Button>
        </div>
      </div>
    );
  }

  return (
    <div className="max-w-4xl mx-auto space-y-8">
      <div>
        <h1 className="text-3xl font-bold">Checkout</h1>
        <p className="text-muted-foreground">Review your order and complete your purchase</p>
      </div>

      <form onSubmit={handleSubmitOrder} className="grid gap-8 lg:grid-cols-2">
        {/* Left Column - Order Details */}
        <div className="space-y-6">
          {/* Cart Summary */}
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <ShoppingCart className="h-5 w-5" />
                Your Order
              </CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              {cart.map((item) => (
                <div key={item.product_id} className="flex justify-between items-center">
                  <div className="flex-1">
                    <p className="font-medium">{item.product_name}</p>
                    <p className="text-sm text-muted-foreground">
                      €{item.unit_price.toFixed(2)} × {item.quantity}
                    </p>
                  </div>
                  <p className="font-medium">€{item.total_price.toFixed(2)}</p>
                </div>
              ))}

              <Separator />

              <div className="flex justify-between items-center text-lg font-bold">
                <span>Total</span>
                <span>€{getTotalAmount().toFixed(2)}</span>
              </div>
            </CardContent>
          </Card>

          {/* Order Type Selection */}
          <Card>
            <CardHeader>
              <CardTitle>Order Type</CardTitle>
              <CardDescription>How would you like to receive your order?</CardDescription>
            </CardHeader>
            <CardContent>
              <RadioGroup value={orderType} onValueChange={(value) => setOrderType(value as "delivery" | "pickup")}>
                <div className="flex items-center space-x-2">
                  <RadioGroupItem value="delivery" id="delivery" />
                  <Label htmlFor="delivery" className="flex items-center gap-2 cursor-pointer">
                    <Truck className="h-4 w-4" />
                    Delivery
                    <Badge variant="secondary">€2.50</Badge>
                  </Label>
                </div>
                <div className="flex items-center space-x-2">
                  <RadioGroupItem value="pickup" id="pickup" />
                  <Label htmlFor="pickup" className="flex items-center gap-2 cursor-pointer">
                    <MapPin className="h-4 w-4" />
                    Pickup
                    <Badge variant="outline">Free</Badge>
                  </Label>
                </div>
              </RadioGroup>
            </CardContent>
          </Card>
        </div>

        {/* Right Column - Delivery Information */}
        <div className="space-y-6">
          {orderType === 'delivery' && (
            <Card>
              <CardHeader>
                <CardTitle>Delivery Information</CardTitle>
                <CardDescription>Where should we deliver your order?</CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="grid gap-4 sm:grid-cols-2">
                  <div className="space-y-2">
                    <Label htmlFor="customer_name">Full Name *</Label>
                    <Input
                      id="customer_name"
                      value={deliveryInfo.customer_name}
                      onChange={(e) => setDeliveryInfo({...deliveryInfo, customer_name: e.target.value})}
                      placeholder="Max Mustermann"
                      required={orderType === 'delivery'}
                    />
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="customer_phone">Phone Number *</Label>
                    <Input
                      id="customer_phone"
                      type="tel"
                      value={deliveryInfo.customer_phone}
                      onChange={(e) => setDeliveryInfo({...deliveryInfo, customer_phone: e.target.value})}
                      placeholder="+49 123 456789"
                      required={orderType === 'delivery'}
                    />
                  </div>
                </div>

                <div className="grid gap-4 sm:grid-cols-3">
                  <div className="space-y-2 sm:col-span-2">
                    <Label htmlFor="street">Street *</Label>
                    <Input
                      id="street"
                      value={deliveryInfo.street}
                      onChange={(e) => setDeliveryInfo({...deliveryInfo, street: e.target.value})}
                      placeholder="Hauptstraße"
                      required={orderType === 'delivery'}
                    />
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="house_number">House Number *</Label>
                    <Input
                      id="house_number"
                      value={deliveryInfo.house_number}
                      onChange={(e) => setDeliveryInfo({...deliveryInfo, house_number: e.target.value})}
                      placeholder="123"
                      required={orderType === 'delivery'}
                    />
                  </div>
                </div>

                <div className="grid gap-4 sm:grid-cols-2">
                  <div className="space-y-2">
                    <Label htmlFor="postal_code">Postal Code *</Label>
                    <Input
                      id="postal_code"
                      value={deliveryInfo.postal_code}
                      onChange={(e) => setDeliveryInfo({...deliveryInfo, postal_code: e.target.value})}
                      placeholder="12345"
                      required={orderType === 'delivery'}
                    />
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="city">City *</Label>
                    <Input
                      id="city"
                      value={deliveryInfo.city}
                      onChange={(e) => setDeliveryInfo({...deliveryInfo, city: e.target.value})}
                      placeholder="München"
                      required={orderType === 'delivery'}
                    />
                  </div>
                </div>

                <div className="space-y-2">
                  <Label htmlFor="floor">Floor/Apartment (Optional)</Label>
                  <Input
                    id="floor"
                    value={deliveryInfo.floor}
                    onChange={(e) => setDeliveryInfo({...deliveryInfo, floor: e.target.value})}
                    placeholder="2nd Floor, Apt 4"
                  />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="instructions">Delivery Instructions (Optional)</Label>
                  <Input
                    id="instructions"
                    value={deliveryInfo.instructions}
                    onChange={(e) => setDeliveryInfo({...deliveryInfo, instructions: e.target.value})}
                    placeholder="Ring doorbell twice"
                  />
                </div>
              </CardContent>
            </Card>
          )}

          {orderType === 'pickup' && (
            <Card>
              <CardHeader>
                <CardTitle>Pickup Information</CardTitle>
                <CardDescription>Pickup details</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-2">
                  <p className="text-sm text-muted-foreground">
                    <strong>Pickup Address:</strong><br />
                    Analytica Restaurant<br />
                    Musterstraße 123<br />
                    80333 München
                  </p>
                  <p className="text-sm text-muted-foreground">
                    <strong>Pickup Time:</strong> Ready in ~30 minutes
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
            {submitting ? 'Processing Order...' : `Place Order - €${getTotalAmount().toFixed(2)}`}
          </Button>
        </div>
      </form>
    </div>
  );
}