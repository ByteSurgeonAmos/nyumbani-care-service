// Payment tracking utilities

/**
 * Store payment information in local storage
 * @param paymentId Payment ID to store
 * @param orderId Order ID associated with the payment
 */
export function trackPayment(paymentId: string, orderId: string): void {
  try {
    // Store the payment ID and order ID
    localStorage.setItem("lastPaymentId", paymentId);
    localStorage.setItem("lastOrderId", orderId);

    // Store the timestamp of when the payment was initiated
    localStorage.setItem("paymentInitiatedAt", Date.now().toString());
  } catch (error) {
    console.error("Error storing payment tracking info:", error);
  }
}

/**
 * Get tracked payment information
 */
export function getTrackedPayment(): {
  paymentId: string | null;
  orderId: string | null;
  timestamp: number | null;
} {
  try {
    const paymentId = localStorage.getItem("lastPaymentId");
    const orderId = localStorage.getItem("lastOrderId");
    const timestamp = localStorage.getItem("paymentInitiatedAt")
      ? parseInt(localStorage.getItem("paymentInitiatedAt") || "0", 10)
      : null;

    return { paymentId, orderId, timestamp };
  } catch (error) {
    console.error("Error retrieving payment tracking info:", error);
    return { paymentId: null, orderId: null, timestamp: null };
  }
}

/**
 * Clear payment tracking information
 */
export function clearPaymentTracking(): void {
  try {
    localStorage.removeItem("lastPaymentId");
    localStorage.removeItem("lastOrderId");
    localStorage.removeItem("paymentInitiatedAt");
  } catch (error) {
    console.error("Error clearing payment tracking info:", error);
  }
}

/**
 * Check if payment tracking is expired
 * Default expiry is 1 hour (3600000 ms)
 */
export function isPaymentTrackingExpired(expiryMs: number = 3600000): boolean {
  try {
    const timestamp = localStorage.getItem("paymentInitiatedAt");
    if (!timestamp) return true;

    const initiated = parseInt(timestamp, 10);
    const now = Date.now();

    return now - initiated > expiryMs;
  } catch (error) {
    console.error("Error checking payment tracking expiry:", error);
    return true;
  }
}
