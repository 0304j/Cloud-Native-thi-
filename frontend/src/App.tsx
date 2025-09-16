import { BrowserRouter as Router, Routes, Route } from "react-router-dom";

// Layout
import Layout from "@/components/Layout";

// Pages
import HomePage from "@/pages/HomePage";
import AuthPage from "@/pages/AuthPage";
import ShopPage from "@/pages/ShopPage";
import CheckoutPage from "@/pages/CheckoutPage";
import PaymentPage from "@/pages/PaymentPage";
import KitchenPage from "@/pages/KitchenPage";

function App() {
  return (
    <Router>
      <Layout>
        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route path="/auth" element={<AuthPage />} />
          <Route path="/shop" element={<ShopPage />} />
          <Route path="/checkout" element={<CheckoutPage />} />
          <Route path="/payment" element={<PaymentPage />} />
          <Route path="/kitchen" element={<KitchenPage />} />
        </Routes>
      </Layout>
    </Router>
  );
}

export default App;