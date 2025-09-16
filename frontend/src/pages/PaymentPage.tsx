import { useState, useEffect } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
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
import {
  AlertCircle,
  ArrowLeft,
  Loader2,
  CheckCircle,
  Wand2,
} from "lucide-react";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faStripe, faPaypal } from "@fortawesome/free-brands-svg-icons";
import { faUniversity } from "@fortawesome/free-solid-svg-icons";
import { fakerDE as faker } from "@faker-js/faker";
import type {
  PaymentRequest,
  PaymentResponse,
  StripePaymentDetails,
  PayPalPaymentDetails,
  BankTransferDetails,
} from "@/types/domain";

export default function PaymentPage() {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();

  const orderId = searchParams.get("order_id") || "";
  const totalAmount = parseFloat(searchParams.get("amount") || "0");
  const currency = searchParams.get("currency") || "EUR";

  const [provider, setProvider] = useState<
    "stripe" | "paypal" | "bank_transfer"
  >("stripe");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);

  const [stripeData, setStripeData] = useState<StripePaymentDetails>({
    card_number: "",
    expiry_month: "",
    expiry_year: "",
    cvv: "",
    cardholder_name: "",
  });

  const [paypalData, setPaypalData] = useState<PayPalPaymentDetails>({
    email: "",
    password: "",
  });

  const [bankData, setBankData] = useState<BankTransferDetails>({
    iban: "",
    account_holder: "",
    bank_name: "",
  });

  useEffect(() => {
    if (!orderId || totalAmount <= 0) {
      navigate("/checkout");
    }
  }, [orderId, totalAmount, navigate]);

  const generateTestPaymentData = () => {
    switch (provider) {
      case "stripe":
        setStripeData({
          cardholder_name: faker.person.fullName(),
          card_number: faker.finance.creditCardNumber("visa"),
          expiry_month: faker.date.month(),
          expiry_year: faker.date.future({ years: 3 }).getFullYear().toString(),
          cvv: faker.finance.creditCardCVV(),
        });
        break;
      case "paypal":
        setPaypalData({
          email: faker.internet.email(),
          password: faker.internet.password({ length: 10, memorable: true }),
        });
        break;
      case "bank_transfer":
        setBankData({
          account_holder: faker.person.fullName(),
          iban: faker.finance.iban({ countryCode: "DE", formatted: true }),
          bank_name: faker.helpers.arrayElement([
            "Deutsche Bank",
            "Commerzbank",
            "Sparkasse",
            "DKB",
            "ING",
          ]),
        });
        break;
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    try {
      let payment_details:
        | StripePaymentDetails
        | PayPalPaymentDetails
        | BankTransferDetails;

      switch (provider) {
        case "stripe":
          if (
            !stripeData.card_number ||
            !stripeData.expiry_month ||
            !stripeData.expiry_year ||
            !stripeData.cvv ||
            !stripeData.cardholder_name
          ) {
            throw new Error("Bitte fülle alle Kreditkarten-Felder aus");
          }
          payment_details = stripeData;
          break;
        case "paypal":
          if (!paypalData.email || !paypalData.password) {
            throw new Error("Bitte fülle alle PayPal-Felder aus");
          }
          payment_details = paypalData;
          break;
        case "bank_transfer":
          if (!bankData.iban || !bankData.account_holder) {
            throw new Error("Bitte fülle alle Überweisungs-Felder aus");
          }
          payment_details = bankData;
          break;
        default:
          throw new Error("Ungültiger Zahlungsanbieter");
      }

      const paymentRequest: PaymentRequest = {
        order_id: orderId,
        user_id: crypto.randomUUID(), // Generate valid UUID for now
        provider,
        amount: totalAmount,
        currency,
        payment_details,
      };

      const response = await fetch("/api/payments", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        credentials: "include",
        body: JSON.stringify(paymentRequest),
      });

      if (!response.ok) {
        const errorData = await response.json();
        throw new Error(errorData.error || "Zahlung fehlgeschlagen");
      }

      const result: PaymentResponse = await response.json();

      setLoading(false);
      setSuccess(true);

      setTimeout(() => {
        navigate(`/tracking?order_id=${orderId}&payment_id=${result.id}`);
      }, 2500);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Zahlung fehlgeschlagen");
    } finally {
      setLoading(false);
    }
  };

  if (success) {
    return (
      <div className="max-w-2xl mx-auto space-y-8 min-h-[400px] flex items-center justify-center">
        <Card className="w-full text-center">
          <CardContent className="pt-6 space-y-6">
            <div className="flex justify-center">
              <CheckCircle className="h-16 w-16 text-green-500" />
            </div>
            <div className="space-y-2">
              <h2 className="text-2xl font-bold text-green-600">
                Zahlung erfolgreich!
              </h2>
              <p className="text-muted-foreground">
                Deine Zahlung wurde verarbeitet. Du wirst zur Bestellverfolgung
                weitergeleitet...
              </p>
            </div>
            <div className="space-y-2">
              <p className="text-sm text-muted-foreground">
                Bestellung #{orderId.slice(-8)}
              </p>
              <p className="text-lg font-semibold">
                €{totalAmount.toFixed(2)} bezahlt
              </p>
            </div>
          </CardContent>
        </Card>
      </div>
    );
  }

  return (
    <div className="max-w-2xl mx-auto space-y-8">
      <div className="flex items-center gap-4">
        <Button
          variant="ghost"
          size="sm"
          onClick={() => navigate("/checkout")}
          className="gap-2"
        >
          <ArrowLeft className="h-4 w-4" />
          Zurück zur Bestellung
        </Button>
      </div>

      <div>
        <h1 className="text-3xl font-bold">Zahlung</h1>
        <p className="text-muted-foreground">
          Wähle deine bevorzugte Zahlungsmethode
        </p>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Bestellübersicht</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex justify-between items-center text-lg font-bold">
            <span>Gesamtbetrag</span>
            <span>€{totalAmount.toFixed(2)}</span>
          </div>
          <p className="text-sm text-muted-foreground mt-1">
            Bestellung #{orderId.slice(-8)}
          </p>
        </CardContent>
      </Card>

      {/* Payment Form */}
      <form onSubmit={handleSubmit} className="space-y-6">
        <Card>
          <CardHeader>
            <CardTitle>Zahlungsmethode</CardTitle>
            <CardDescription>
              Wähle aus unseren sicheren Zahlungsoptionen
            </CardDescription>
          </CardHeader>
          <CardContent>
            <RadioGroup
              value={provider}
              onValueChange={(value) => setProvider(value as typeof provider)}
            >
              <div className="flex items-center space-x-2">
                <RadioGroupItem value="stripe" id="stripe" />
                <Label
                  htmlFor="stripe"
                  className="flex items-center gap-3 cursor-pointer"
                >
                  <FontAwesomeIcon
                    icon={faStripe}
                    className="h-6 w-6 text-[#635bff]"
                  />
                  Kreditkarte (Stripe)
                </Label>
              </div>
              <div className="flex items-center space-x-2">
                <RadioGroupItem value="paypal" id="paypal" />
                <Label
                  htmlFor="paypal"
                  className="flex items-center gap-3 cursor-pointer"
                >
                  <FontAwesomeIcon
                    icon={faPaypal}
                    className="h-6 w-6 text-[#0070ba]"
                  />
                  PayPal Konto
                </Label>
              </div>
              <div className="flex items-center space-x-2">
                <RadioGroupItem value="bank_transfer" id="bank_transfer" />
                <Label
                  htmlFor="bank_transfer"
                  className="flex items-center gap-3 cursor-pointer"
                >
                  <FontAwesomeIcon
                    icon={faUniversity}
                    className="h-5 w-5 text-blue-800"
                  />
                  Banküberweisung (SEPA)
                </Label>
              </div>
            </RadioGroup>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle>Zahlungsdetails</CardTitle>
            <CardDescription>
              {provider === "stripe" &&
                "Gib deine Kreditkarteninformationen ein"}
              {provider === "paypal" && "Melde dich mit deinem PayPal-Konto an"}
              {provider === "bank_transfer" && "Gib deine Bankverbindung ein"}
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex justify-end">
              <Button
                type="button"
                variant="outline"
                size="sm"
                onClick={generateTestPaymentData}
                className="flex items-center gap-2 w-fit"
              >
                <Wand2 className="h-4 w-4" />
                Testdaten generieren
              </Button>
            </div>
            {provider === "stripe" && (
              <>
                <div className="space-y-2">
                  <Label htmlFor="cardholder_name">Karteninhaber *</Label>
                  <Input
                    id="cardholder_name"
                    value={stripeData.cardholder_name}
                    onChange={(e) =>
                      setStripeData({
                        ...stripeData,
                        cardholder_name: e.target.value,
                      })
                    }
                    placeholder="Max Mustermann"
                    required
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="card_number">Kartennummer *</Label>
                  <Input
                    id="card_number"
                    value={stripeData.card_number}
                    onChange={(e) =>
                      setStripeData({
                        ...stripeData,
                        card_number: e.target.value,
                      })
                    }
                    placeholder="1234 5678 9012 3456"
                    maxLength={19}
                    required
                  />
                </div>
                <div className="grid grid-cols-3 gap-4">
                  <div className="space-y-2">
                    <Label htmlFor="expiry_month">Monat *</Label>
                    <Input
                      id="expiry_month"
                      value={stripeData.expiry_month}
                      onChange={(e) =>
                        setStripeData({
                          ...stripeData,
                          expiry_month: e.target.value,
                        })
                      }
                      placeholder="MM"
                      maxLength={2}
                      required
                    />
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="expiry_year">Jahr *</Label>
                    <Input
                      id="expiry_year"
                      value={stripeData.expiry_year}
                      onChange={(e) =>
                        setStripeData({
                          ...stripeData,
                          expiry_year: e.target.value,
                        })
                      }
                      placeholder="YY"
                      maxLength={2}
                      required
                    />
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="cvv">CVV *</Label>
                    <Input
                      id="cvv"
                      value={stripeData.cvv}
                      onChange={(e) =>
                        setStripeData({ ...stripeData, cvv: e.target.value })
                      }
                      placeholder="123"
                      maxLength={3}
                      required
                    />
                  </div>
                </div>
              </>
            )}

            {provider === "paypal" && (
              <>
                <div className="space-y-2">
                  <Label htmlFor="paypal_email">PayPal E-Mail *</Label>
                  <Input
                    id="paypal_email"
                    type="email"
                    value={paypalData.email}
                    onChange={(e) =>
                      setPaypalData({ ...paypalData, email: e.target.value })
                    }
                    placeholder="dein@paypal.de"
                    required
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="paypal_password">PayPal Passwort *</Label>
                  <Input
                    id="paypal_password"
                    type="password"
                    value={paypalData.password}
                    onChange={(e) =>
                      setPaypalData({ ...paypalData, password: e.target.value })
                    }
                    placeholder="••••••••"
                    required
                  />
                </div>
              </>
            )}

            {provider === "bank_transfer" && (
              <>
                <div className="space-y-2">
                  <Label htmlFor="account_holder">Kontoinhaber *</Label>
                  <Input
                    id="account_holder"
                    value={bankData.account_holder}
                    onChange={(e) =>
                      setBankData({
                        ...bankData,
                        account_holder: e.target.value,
                      })
                    }
                    placeholder="Max Mustermann"
                    required
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="iban">IBAN *</Label>
                  <Input
                    id="iban"
                    value={bankData.iban}
                    onChange={(e) =>
                      setBankData({ ...bankData, iban: e.target.value })
                    }
                    placeholder="DE89 3704 0044 0532 0130 00"
                    required
                  />
                </div>
                <div className="space-y-2">
                  <Label htmlFor="bank_name">Bank (Optional)</Label>
                  <Input
                    id="bank_name"
                    value={bankData.bank_name}
                    onChange={(e) =>
                      setBankData({ ...bankData, bank_name: e.target.value })
                    }
                    placeholder="Deutsche Bank"
                  />
                </div>
              </>
            )}
          </CardContent>
        </Card>

        {error && (
          <div className="flex items-center gap-2 rounded-md bg-destructive/10 p-3 text-sm text-destructive">
            <AlertCircle className="h-4 w-4" />
            {error}
          </div>
        )}

        <Button type="submit" className="w-full" size="lg" disabled={loading}>
          {loading ? (
            <>
              <Loader2 className="mr-2 h-4 w-4 animate-spin" />
              Zahlung wird verarbeitet...
            </>
          ) : (
            `Jetzt bezahlen - €${totalAmount.toFixed(2)}`
          )}
        </Button>
      </form>

      <div className="text-center text-sm text-muted-foreground">
        <p>🔒 Deine Zahlungsdaten sind sicher verschlüsselt</p>
      </div>
    </div>
  );
}
