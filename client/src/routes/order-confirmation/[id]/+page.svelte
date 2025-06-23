<script lang="ts">
  import { onMount } from "svelte";
  import { goto } from "$app/navigation";
  import { isAuthenticated, orderService } from "$lib/api";
  import { paymentService, checkoutService } from "$lib/api/paymentService";
  import type { ExtendedTestKitOrder } from "$lib/api/types";
  import type { PaymentStatus } from "$lib/api/paymentService";
  import {
    trackPayment,
    getTrackedPayment,
    clearPaymentTracking,
  } from "$lib/utils/payment";
  import toast from "svelte-french-toast";
  export let data: { id: string };
  const { id } = data;

  let order: ExtendedTestKitOrder | null = null;
  let payment: PaymentStatus | null = null;
  let isLoading = true;
  let pollingInterval: any = null;
  onMount(() => {
    // Check if user is authenticated
    if (!$isAuthenticated) {
      toast.error("Please login to view order details");
      goto("/login?redirectTo=/order-confirmation/" + id);
      return;
    }

    // Store user email for payment processing
    if ($isAuthenticated && localStorage.getItem("user")) {
      try {
        const user = JSON.parse(localStorage.getItem("user") || "{}");
        if (user.email) {
          localStorage.setItem("userEmail", user.email);
        }
      } catch (e) {
        console.error("Failed to parse user data:", e);
      }
    } // Check if we came back from a payment gateway
    const trackedPayment = getTrackedPayment();
    if (trackedPayment.orderId === id && trackedPayment.paymentId) {
      toast("Checking payment status...", {
        duration: 4000,
        position: "top-center",
      });

      // Clear payment tracking since we're handling it now
      clearPaymentTracking();
    }

    // Start loading order details
    loadOrderDetails().then(() => {
      // Start polling for payment status updates if payment is pending
      if (payment && payment.status === "pending") {
        pollingInterval = setInterval(checkPaymentStatus, 5000);
      }
    });

    // Return the cleanup function directly
    return () => {
      if (pollingInterval) clearInterval(pollingInterval);
    };
  });
  async function loadOrderDetails() {
    isLoading = true;

    try {
      // Load order details
      order = await orderService.getById(id);

      // Ensure order total price is properly set
      if (order && order.total_amount && !order.total_price) {
        order.total_price = order.total_amount;
      }

      // Ensure order items are properly initialized
      if (order && !order.items) {
        order.items = [];
      } // Try to find the associated payment
      try {
        const payments = await paymentService.listPayments();
        payment = payments.find((p: any) => p.order_id === id) || null;
      } catch (err) {
        console.error("Failed to load payment details:", err);
        payment = null;
      }
    } catch (err) {
      console.error("Failed to load order details:", err);
      toast.error("Failed to load order details");
      order = null;
    } finally {
      isLoading = false;
    }
  }
  async function checkPaymentStatus() {
    if (!payment) return;

    try {
      const updatedPayment = await paymentService.getPaymentStatus(payment.id);

      // Preserve the original order_id and transaction_id if they exist
      payment = {
        ...updatedPayment,
        order_id: payment.order_id || updatedPayment.order_id,
        transaction_id: payment.transaction_id || updatedPayment.transaction_id,
      };

      if (updatedPayment.status === "completed") {
        if (pollingInterval) clearInterval(pollingInterval);
        await loadOrderDetails();
        toast.success("Payment completed successfully!");
      } else if (updatedPayment.status === "failed") {
        if (pollingInterval) clearInterval(pollingInterval);
        toast.error("Payment failed. Please try again.");
      }
    } catch (err) {
      console.error("Failed to check payment status:", err);
    }
  }
  async function retryPayment() {
    if (!order) return;

    try {
      // Check if we have order items
      if (
        !order.items ||
        !Array.isArray(order.items) ||
        order.items.length === 0
      ) {
        toast.error("Cannot retry payment: Order details are incomplete");
        return;
      }

      // Check inventory before allowing retry
      const inventoryValidation = await checkoutService.validateCartInventory(
        order.items.map((item) => ({
          testKit: item.test_kit,
          quantity: item.quantity,
        }))
      );

      if (!inventoryValidation.valid) {
        if (
          inventoryValidation.errors &&
          inventoryValidation.errors.length > 0
        ) {
          // Show specific inventory errors
          const errorMessages = inventoryValidation.errors.map(
            (error) =>
              `${error.name}: Only ${error.available} available (you requested ${error.requested})`
          );
          toast.error(`Cannot retry payment: ${errorMessages.join(", ")}`);
        } else {
          toast.error(
            "Cannot retry payment: Some items are no longer available"
          );
        }
        return;
      }

      // Make sure we have access to the user email
      const userEmail =
        typeof order.user_id === "string" && order.user_id.includes("@")
          ? order.user_id
          : order.user_email || localStorage.getItem("userEmail");

      if (!userEmail) {
        toast.error(
          "Missing user email information. Please try again or contact support."
        );
        return;
      }

      // Initiate new payment
      const paymentResponse = await paymentService.initiatePayment({
        email: userEmail,
        amount:
          order.total_amount || order.total_price || calculateTotal(order),
        order_id: order.id,
      }); // Store payment tracking information
      trackPayment(paymentResponse.payment_id, order.id);

      toast.success("Redirecting to payment...");

      // If there's a payment URL in the response (sometimes stored in other fields)
      const paymentUrl =
        paymentResponse.payment_url || paymentResponse.reference;

      if (paymentUrl && paymentUrl.startsWith("http")) {
        // Redirect to payment gateway
        window.location.href = paymentUrl;
      } else {
        // If no payment URL, refresh the page
        await loadOrderDetails();
      }
    } catch (error) {
      console.error("Payment retry error:", error);
      toast.error("Failed to process payment");
    }
  }

  function getStatusColor(status: string): string {
    switch (status.toLowerCase()) {
      case "completed":
      case "paid":
      case "confirmed":
      case "delivered":
        return "bg-green-100 text-green-800";
      case "pending":
      case "processing":
      case "shipped":
        return "bg-yellow-100 text-yellow-800";
      case "failed":
      case "cancelled":
        return "bg-red-100 text-red-800";
      default:
        return "bg-gray-100 text-gray-800";
    }
  }
  function formatDate(dateString: string): string {
    const date = new Date(dateString);
    return new Intl.DateTimeFormat("en-US", {
      year: "numeric",
      month: "long",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    }).format(date);
  }
  // Calculate total price including shipping
  function calculateTotal(order: ExtendedTestKitOrder | null): number {
    if (!order) return 0;

    // Check if the total price already includes shipping
    // This is a heuristic - if the order has total_amount and total_price, and they differ by around 300,
    // assume shipping is already included
    const basePrice = order.total_price || order.total_amount || 0;

    // Check if shipping cost is already included in the total
    if (
      order.total_amount &&
      order.total_price &&
      Math.abs(order.total_price - order.total_amount) >= 290 &&
      Math.abs(order.total_price - order.total_amount) <= 310
    ) {
      return basePrice;
    }

    const shippingCost = 300; // Standard shipping cost
    return basePrice + shippingCost;
  }

  // Check if shipping is already included in the price
  function isShippingIncluded(order: ExtendedTestKitOrder | null): boolean {
    if (!order) return false;

    // If we have both fields and they differ by around 300, shipping is likely included
    if (
      order.total_amount &&
      order.total_price &&
      Math.abs(order.total_price - order.total_amount) >= 290 &&
      Math.abs(order.total_price - order.total_amount) <= 310
    ) {
      return true;
    }

    return false;
  }
  async function handlePayNow() {
    if (!order) return;

    try {
      // Make sure we have access to the user email
      const userEmail =
        typeof order.user_id === "string" && order.user_id.includes("@")
          ? order.user_id
          : order.user_email || localStorage.getItem("userEmail");

      if (!userEmail) {
        toast.error(
          "Missing user email information. Please try again or contact support."
        );
        return;
      }

      // Generate a new payment for this order
      const paymentResponse = await paymentService.initiatePayment({
        email: userEmail,
        amount:
          order.total_amount || order.total_price || calculateTotal(order),
        order_id: order.id,
      });

      // Store payment tracking information
      trackPayment(paymentResponse.payment_id, order.id);

      toast.success("Redirecting to payment...");

      // If there's a payment URL in the response
      const paymentUrl =
        paymentResponse.payment_url || paymentResponse.reference;

      if (paymentUrl && paymentUrl.startsWith("http")) {
        // Redirect to payment gateway
        window.location.href = paymentUrl;
      } else {
        // If no payment URL, refresh the page
        await loadOrderDetails();
      }
    } catch (error) {
      console.error("Payment initiation error:", error);
      toast.error("Failed to initiate payment");
    }
  }
</script>

<svelte:head>
  <title>Order Confirmation | NyumbaniCare</title>
  <meta
    name="description"
    content="Your order confirmation and tracking details"
  />
</svelte:head>

<div class="bg-white">
  <div class="mx-auto max-w-3xl px-4 py-16 sm:px-6 sm:py-24 lg:px-8">
    {#if isLoading}
      <div class="flex justify-center items-center h-64">
        <div class="animate-pulse text-center">
          <div class="h-8 w-3/4 bg-gray-200 rounded mx-auto mb-4"></div>
          <div class="h-4 w-1/2 bg-gray-200 rounded mx-auto"></div>
          <div class="mt-8 space-y-2">
            <div class="h-4 w-full bg-gray-200 rounded"></div>
            <div class="h-4 w-5/6 bg-gray-200 rounded"></div>
            <div class="h-4 w-4/6 bg-gray-200 rounded"></div>
          </div>
        </div>
      </div>
    {:else if !order}
      <div class="text-center">
        <h1 class="text-3xl font-bold tracking-tight text-gray-900 sm:text-4xl">
          Order Not Found
        </h1>
        <p class="mt-4 text-base text-gray-500">
          We couldn't find the order you're looking for.
        </p>
        <div class="mt-6">
          <a href="/test-kits" class="text-primary-600 hover:text-primary-500">
            Continue Shopping
            <span aria-hidden="true"> &rarr;</span>
          </a>
        </div>
      </div>
    {:else}
      <div>
        <h1 class="text-3xl font-bold tracking-tight text-gray-900 sm:text-4xl">
          Order Confirmation
        </h1>
        <p class="mt-2 text-base text-gray-500">Thank you for your order.</p>

        <div class="mt-6">
          <dl class="grid grid-cols-1 gap-x-6 gap-y-4 sm:grid-cols-2">
            <div class="border-t border-gray-200 pt-4">
              <dt class="text-sm font-medium text-gray-500">Order number</dt>
              <dd class="mt-1 text-sm font-medium text-gray-900">{order.id}</dd>
            </div>
            <div class="border-t border-gray-200 pt-4">
              <dt class="text-sm font-medium text-gray-500">Order date</dt>
              <dd class="mt-1 text-sm font-medium text-gray-900">
                {formatDate(order.created_at)}
              </dd>
            </div>
            <div class="border-t border-gray-200 pt-4 sm:col-span-2">
              <dt class="text-sm font-medium text-gray-500">Order status</dt>
              <dd class="mt-1 text-sm font-medium">
                <span
                  class={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${getStatusColor(order.status)}`}
                >
                  {order.status.charAt(0).toUpperCase() + order.status.slice(1)}
                </span>
              </dd>
            </div>
            <div class="border-t border-gray-200 pt-4 sm:col-span-2">
              <dt class="text-sm font-medium text-gray-500">Payment status</dt>
              <dd class="mt-1 text-sm font-medium">
                <span
                  class={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${getStatusColor(order.payment_status)}`}
                >
                  {order.payment_status.charAt(0).toUpperCase() +
                    order.payment_status.slice(1)}
                </span>
              </dd>
            </div>
            {#if order.tracking_number}
              <div class="border-t border-gray-200 pt-4 sm:col-span-2">
                <dt class="text-sm font-medium text-gray-500">
                  Tracking number
                </dt>
                <dd class="mt-1 text-sm font-medium text-gray-900">
                  {order.tracking_number}
                </dd>
              </div>
            {/if}
            <div class="border-t border-gray-200 pt-4 sm:col-span-2">
              <dt class="text-sm font-medium text-gray-500">
                Shipping address
              </dt>
              <dd class="mt-1 text-sm text-gray-900 whitespace-pre-line">
                {order.shipping_address}
              </dd>
            </div>
          </dl>
          {#if order.payment_status === "pending" || order.payment_status === "failed"}
            <div class="mt-8 border-t border-gray-200 pt-8">
              <div
                class="rounded-md {order.payment_status === 'failed'
                  ? 'bg-red-50'
                  : 'bg-yellow-50'} p-4"
              >
                <div class="flex">
                  <div class="flex-shrink-0">
                    <svg
                      class="h-5 w-5 {order.payment_status === 'failed'
                        ? 'text-red-400'
                        : 'text-yellow-400'}"
                      viewBox="0 0 20 20"
                      fill="currentColor"
                      aria-hidden="true"
                    >
                      <path
                        fill-rule="evenodd"
                        d="M8.485 2.495c.673-1.167 2.357-1.167 3.03 0l6.28 10.875c.673 1.167-.17 2.625-1.516 2.625H3.72c-1.347 0-2.189-1.458-1.515-2.625L8.485 2.495zM10 5a.75.75 0 01.75.75v3.5a.75.75 0 01-1.5 0v-3.5A.75.75 0 0110 5zm0 9a1 1 0 100-2 1 1 0 000 2z"
                        clip-rule="evenodd"
                      />
                    </svg>
                  </div>
                  <div class="ml-3">
                    <h3
                      class="text-sm font-medium {order.payment_status ===
                      'failed'
                        ? 'text-red-800'
                        : 'text-yellow-800'}"
                    >
                      {order.payment_status === "failed"
                        ? "Payment failed"
                        : "Payment pending"}
                    </h3>
                    <div
                      class="mt-2 text-sm {order.payment_status === 'failed'
                        ? 'text-red-700'
                        : 'text-yellow-700'}"
                    >
                      {#if order.payment_status === "failed"}
                        <p>
                          Your payment was not successful. You can try again
                          with a different payment method or contact customer
                          support.
                        </p>
                      {:else}
                        <p>
                          Your order has been placed but payment is still
                          pending. Please complete the payment to process your
                          order.
                        </p>
                      {/if}
                    </div>
                    <div class="mt-4">
                      <button
                        type="button"
                        class="rounded-md bg-primary-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-primary-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-600"
                        on:click={order.payment_status === "failed"
                          ? retryPayment
                          : handlePayNow}
                      >
                        Pay Now
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          {:else if order.payment_status === "paid"}
            <div class="mt-8 border-t border-gray-200 pt-8">
              <div class="rounded-md bg-green-50 p-4">
                <div class="flex">
                  <div class="flex-shrink-0">
                    <svg
                      class="h-5 w-5 text-green-400"
                      viewBox="0 0 20 20"
                      fill="currentColor"
                      aria-hidden="true"
                    >
                      <path
                        fill-rule="evenodd"
                        d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.857-9.809a.75.75 0 00-1.214-.882l-3.483 4.79-1.88-1.88a.75.75 0 10-1.06 1.061l2.5 2.5a.75.75 0 001.137-.089l4-5.5z"
                        clip-rule="evenodd"
                      />
                    </svg>
                  </div>
                  <div class="ml-3">
                    <h3 class="text-sm font-medium text-green-800">
                      Payment successful
                    </h3>
                    <div class="mt-2 text-sm text-green-700">
                      <p>
                        Your payment has been processed successfully. We'll ship
                        your order soon.
                      </p>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          {/if}

          <table class="mt-8 w-full text-gray-500">
            <caption class="sr-only">Products</caption>
            <thead
              class="sr-only text-left text-sm text-gray-500 sm:not-sr-only"
            >
              <tr>
                <th scope="col" class="py-3 pr-8 font-normal sm:w-2/5 lg:w-1/3"
                  >Product</th
                >
                <th
                  scope="col"
                  class="hidden w-1/5 py-3 pr-8 font-normal sm:table-cell"
                  >Price</th
                >
                <th
                  scope="col"
                  class="hidden py-3 pr-8 font-normal sm:table-cell"
                  >Quantity</th
                >
                <th scope="col" class="w-0 py-3 text-right font-normal"
                  >Total</th
                >
              </tr>
            </thead>
            <tbody
              class="divide-y divide-gray-200 border-b border-gray-200 text-sm sm:border-t"
            >
              {#if order.items && order.items.length > 0}
                {#each order.items as item}
                  <tr>
                    <td class="py-6 pr-8">
                      <div class="flex items-center">
                        <img
                          src={item.test_kit?.image_url ||
                            "https://via.placeholder.com/100x100?text=Test+Kit"}
                          alt={item.test_kit?.name}
                          class="mr-4 h-16 w-16 rounded object-cover object-center"
                        />
                        <div>
                          <div class="font-medium text-gray-900">
                            {item.test_kit?.name || "Test Kit"}
                          </div>
                          <div class="mt-1 sm:hidden">
                            KES {item.unit_price.toFixed(2)}
                          </div>
                        </div>
                      </div>
                    </td>
                    <td class="hidden py-6 pr-8 sm:table-cell"
                      >KES {item.unit_price.toFixed(2)}</td
                    >
                    <td class="hidden py-6 pr-8 sm:table-cell"
                      >{item.quantity}</td
                    >
                    <td class="py-6 text-right font-medium"
                      >KES {item.subtotal.toFixed(2)}</td
                    >
                  </tr>
                {/each}
              {:else}
                <tr>
                  <td class="py-6 pr-8" colspan="4">
                    <div class="flex items-center">
                      <img
                        src="https://via.placeholder.com/100x100?text=Test+Kit"
                        alt="Test Kit"
                        class="mr-4 h-16 w-16 rounded object-cover object-center"
                      />
                      <div>
                        <div class="font-medium text-gray-900">
                          {order.test_kit?.name || "Test Kit"}
                        </div>
                        <div class="mt-1 sm:hidden">
                          KES {order.total_price?.toFixed(2)}
                        </div>
                      </div>
                    </div>
                  </td>
                  <td class="hidden py-6 pr-8 sm:table-cell"
                    >KES {order.total_price?.toFixed(2)}</td
                  >
                  <td class="hidden py-6 pr-8 sm:table-cell"
                    >{order.quantity || 1}</td
                  >
                  <td class="py-6 text-right font-medium"
                    >KES {order.total_price?.toFixed(2)}</td
                  >
                </tr>
              {/if}
              {#if !isShippingIncluded(order)}
                <!-- Shipping row (only shown when not already included) -->
                <tr>
                  <td class="py-6 pr-8" colspan="3">
                    <div class="font-medium text-gray-900">Shipping</div>
                  </td>
                  <td class="py-6 text-right font-medium">KES 300.00</td>
                </tr>
              {/if}
              <!-- Total row -->
              <tr>
                <td class="py-6 pr-8" colspan="3">
                  <div class="font-medium text-gray-900">Total</div>
                </td>
                <td class="py-6 text-right font-medium"
                  >KES {calculateTotal(order).toFixed(2)}</td
                >
              </tr>
            </tbody>
          </table>

          <div class="mt-8 border-t border-gray-200 pt-8">
            <h2 class="text-lg font-medium text-gray-900">Need help?</h2>
            <div class="mt-4 flex space-x-4">
              <a
                href="/support"
                class="text-sm font-medium text-primary-600 hover:text-primary-500"
                >Contact Support</a
              >
              <a
                href="/about"
                class="text-sm font-medium text-primary-600 hover:text-primary-500"
                >About NyumbaniCare</a
              >
            </div>
          </div>
        </div>
      </div>
    {/if}

    <div class="mt-16 border-t border-gray-200 pt-8">
      <a
        href="/test-kits"
        class="text-primary-600 hover:text-primary-500 flex items-center"
      >
        <svg
          class="h-5 w-5 mr-2"
          viewBox="0 0 20 20"
          fill="currentColor"
          aria-hidden="true"
        >
          <path
            fill-rule="evenodd"
            d="M17 10a.75.75 0 01-.75.75H5.612l4.158 3.96a.75.75 0 11-1.04 1.08l-5.5-5.25a.75.75 0 010-1.08l5.5-5.25a.75.75 0 111.04 1.08L5.612 9.25H16.25A.75.75 0 0117 10z"
            clip-rule="evenodd"
          />
        </svg>
        Continue Shopping
      </a>
    </div>
  </div>
</div>
