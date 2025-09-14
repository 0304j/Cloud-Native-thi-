import { BrowserRouter as Router, Routes, Route } from "react-router-dom";

// Pages
import HomePage from "@/pages/HomePage";
import AuthPage from "@/pages/AuthPage";
import ShopPage from "@/pages/ShopPage";
import CheckoutPage from "@/pages/CheckoutPage";
import KitchenPage from "@/pages/KitchenPage";

function App() {
  return (
    <Router>
      <div className="min-h-screen bg-gray-50">
        <Routes>
          <Route path="/" element={<HomePage />} />
          <Route path="/auth" element={<AuthPage />} />
          <Route path="/shop" element={<ShopPage />} />
          <Route path="/checkout" element={<CheckoutPage />} />
          <Route path="/kitchen" element={<KitchenPage />} />
        </Routes>
      </div>
    </Router>
  );
}

export default App;