// Domain Models - Shared types based on Go backend domain models
// Keep these in sync with the corresponding Go structs

// Auth Service Domain
export interface User {
  id: string;
  email: string;
  password_hash: string;
  role: 'admin' | 'user';
  created_at: string;
  updated_at: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  email: string;
  password: string;
  role: 'admin' | 'user';
}

// Shopping Service Domain
export interface Product {
  id: string;
  name: string;
  description: string;
  price: number;
  user_id: string;
}

export interface CartItem {
  ProductID: string;
  Qty: number;
}

export interface Cart {
  UserID: string;
  Items: CartItem[];
}

export interface AddToCartRequest {
  product_id: string;
  qty: number;
}

// Checkout Service Domain
export interface CheckoutItem {
  product_id: string;
  product_name: string;
  quantity: number;
  unit_price: number;
  total_price: number;
}

export interface DeliveryInfo {
  customer_name: string;
  customer_phone: string;
  street: string;
  house_number: string;
  postal_code: string;
  city: string;
  floor?: string;
  instructions?: string;
}

export interface CheckoutRequest {
  user_id: string;
  items: CheckoutItem[];
  total_amount: number;
  currency: string;
  order_type: 'delivery' | 'pickup';
  delivery_info?: DeliveryInfo;
  event_type?: string;
  status?: string;
  timestamp?: string;
}

export interface Order {
  id: string;
  user_id: string;
  items: CheckoutItem[];
  total_amount: number;
  currency: string;
  order_type: 'delivery' | 'pickup';
  delivery_info?: DeliveryInfo;
  status: string;
  created_at: string;
  updated_at: string;
}

export interface OrderResponse {
  status: string;
  order_id: string;
  order_type: 'delivery' | 'pickup';
  total_amount: number;
  currency: string;
  created_at: string;
}

// Kitchen Service Domain (for future TrackingPage)
export interface KitchenOrder {
  id: string;
  user_id: string;
  items: CheckoutItem[];
  total_amount: number;
  currency: string;
  order_type: 'delivery' | 'pickup';
  delivery_info?: DeliveryInfo;
  kitchen_status: 'received' | 'preparing' | 'ready' | 'completed';
  estimated_ready_time?: string;
  created_at: string;
  updated_at: string;
}

// API Response wrappers
export interface ApiError {
  error: string;
  details?: string;
}

export interface ApiResponse<T> {
  data?: T;
  error?: string;
  message?: string;
}