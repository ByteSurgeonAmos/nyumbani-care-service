import apiClient from "./apiClient";
import type { TestKitOrder } from "./apiServices";
import type { ExtendedTestKitOrder } from "./types";

// Types for payment processing
export interface PaymentInitiateRequest {
  email: string;
  amount: number;
  order_id: string;
}

export interface PaymentInitiateResponse {
  payment_id: string;
  reference: string;
  status: string;
  message: string;
  payment_url?: string;
}

export interface PaymentStatus {
  id: string;
  order_id: string;
  user_id: string;
  amount: number;
  currency: string;
  status: string; // pending, completed, failed
  transaction_id: string;
  payment_method: string;
  created_at: string;
  updated_at: string;
}

export const paymentService = {
  async initiatePayment(
    data: PaymentInitiateRequest
  ): Promise<PaymentInitiateResponse> {
    const response = await apiClient.post<PaymentInitiateResponse>(
      "/payments/paystack",
      data
    );
    return response.data;
  },

  async getPaymentStatus(paymentId: string): Promise<PaymentStatus> {
    const response = await apiClient.get<PaymentStatus>(
      `/payments/${paymentId}`
    );
    return response.data;
  },

  async listPayments(page = 1, limit = 10): Promise<PaymentStatus[]> {
    const response = await apiClient.get<PaymentStatus[]>("/payments", {
      params: { page, limit },
    });
    return response.data;
  },
};

export interface CheckoutFormData {
  email: string;
  fullName: string;
  phoneNumber: string;
  address: string;
  city: string;
  state: string;
  zipCode: string;
  country: string;
  paymentMethod: string;
}

export interface CreateOrderWithItemsRequest {
  test_kit_ids: { id: string; quantity: number }[];
  shipping_address: string;
  payment_method: string;
  contact_number?: string;
  email?: string;
}

export interface InventoryError {
  id: string;
  name: string;
  requested: number;
  available: number;
}

export const checkoutService = {
  async createOrderWithItems(
    data: CreateOrderWithItemsRequest
  ): Promise<TestKitOrder> {
    try {
      const response = await apiClient.post<TestKitOrder>("/orders", data);
      return response.data;
    } catch (error: any) {
      // Handle specific inventory errors from API
      if (error?.response?.data?.errors?.inventory) {
        const inventoryErrors: InventoryError[] =
          error.response.data.errors.inventory;
        const errorItems = inventoryErrors.map(
          (item) =>
            `${item.name} (requested: ${item.requested}, available: ${item.available})`
        );
        throw new Error(
          `Inventory insufficient for some items: ${errorItems.join(", ")}`
        );
      }
      throw error;
    }
  },

  // Format address from form data
  formatShippingAddress(formData: CheckoutFormData): string {
    return `${formData.fullName}, ${formData.address}, ${formData.city}, ${formData.state} ${formData.zipCode}, ${formData.country}`;
  },

  // Validate if all items in cart are in stock
  async validateCartInventory(
    cartItems: { testKit: { id: string }; quantity: number }[]
  ): Promise<{ valid: boolean; errors?: InventoryError[] }> {
    try {
      const response = await apiClient.post<{
        valid: boolean;
        errors?: InventoryError[];
      }>("/orders/validate-inventory", {
        items: cartItems.map((item) => ({
          id: item.testKit.id,
          quantity: item.quantity,
        })),
      });
      return response.data;
    } catch (error) {
      console.error("Inventory validation error:", error);
      return { valid: false };
    }
  },
};
