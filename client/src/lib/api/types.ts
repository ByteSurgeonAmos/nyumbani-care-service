// Common types used across the application

import type { TestKit } from "./apiServices";

// Extended TestKitOrder interface with optional fields for flexibility
export interface ExtendedTestKitOrder {
  id: string;
  user_id: string;
  user_email?: string;
  status: string;
  total_amount?: number;
  total_price?: number;
  shipping_address: string;
  payment_status: string;
  tracking_number?: string;
  quantity?: number;
  test_kit?: TestKit;
  items?: TestKitOrderItem[];
  created_at: string;
  updated_at: string;
}

export interface TestKitOrderItem {
  id: string;
  order_id: string;
  test_kit_id: string;
  test_kit: TestKit;
  quantity: number;
  unit_price: number;
  subtotal: number;
}
