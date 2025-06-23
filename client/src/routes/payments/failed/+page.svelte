<script lang="ts">
  import { onMount } from "svelte";
  import { goto } from "$app/navigation";
  import { paymentService } from "$lib/api/paymentService";
  import type { PaymentStatus } from "$lib/api/paymentService";
  import toast from "svelte-french-toast";

  let orderId = "";
  let paymentId = "";
  let error = "";
  let isLoading = true;
  let paymentStatus: PaymentStatus | null = null;
  let redirectCountdown = 15; // countdown before auto-redirect
  let countdownInterval: ReturnType<typeof setInterval>;

  onMount(async () => {
    // Notify user
    toast.error("Payment failed");

    // Get the last order ID and payment ID from localStorage if available
    orderId = localStorage.getItem("lastOrderId") || "";
    paymentId = localStorage.getItem("lastPaymentId") || "";

    // Get query parameters
    const urlParams = new URLSearchParams(window.location.search);
    const errorParam = urlParams.get("error");
    const errorMessage = urlParams.get("error_message");

    if (errorParam) {
      error = errorParam;
      if (errorMessage) {
        error += `: ${errorMessage}`;
      }
    } // Try to get payment details if we have a payment ID
    if (paymentId) {
      try {
        paymentStatus = await paymentService.getPaymentStatus(paymentId);
      } catch (err) {
        console.error("Error fetching payment status:", err);
      }
    }

    isLoading = false;

    // Start countdown for auto-redirect
    countdownInterval = setInterval(() => {
      redirectCountdown--;
      if (redirectCountdown <= 0) {
        clearInterval(countdownInterval);
        handleRedirect();
      }
    }, 1000);
  });

  function handleRedirect() {
    clearInterval(countdownInterval);

    if (orderId) {
      goto(`/order-confirmation/${orderId}`);
    } else {
      goto("/cart");
    }
  }

  // Clean up on component destruction
  function onDestroy() {
    if (countdownInterval) {
      clearInterval(countdownInterval);
    }
  }
</script>

<svelte:head>
  <title>Payment Failed | NyumbaniCare</title>
  <meta name="description" content="Your payment was not successful" />
</svelte:head>

<div
  class="min-h-screen bg-white px-4 py-16 sm:px-6 sm:py-24 md:grid md:place-items-center lg:px-8"
>
  <div class="mx-auto max-w-md">
    <main>
      <div class="text-center">
        <div class="flex justify-center">
          <div
            class="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-red-100"
          >
            <svg
              class="h-8 w-8 text-red-600"
              fill="none"
              viewBox="0 0 24 24"
              stroke-width="1.5"
              stroke="currentColor"
              aria-hidden="true"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M6 18L18 6M6 6l12 12"
              />
            </svg>
          </div>
        </div>
        <h1
          class="mt-3 text-3xl font-bold tracking-tight text-gray-900 sm:text-4xl"
        >
          Payment Failed
        </h1>

        {#if isLoading}
          <div class="mt-6 flex justify-center">
            <div
              class="animate-spin rounded-full h-10 w-10 border-b-2 border-primary-600"
            ></div>
          </div>
        {:else}
          <p class="mt-3 text-base text-gray-500">
            We were unable to process your payment for your order.
          </p>

          {#if error}
            <div class="mt-4 p-4 bg-red-50 rounded-md">
              <p class="text-sm text-red-700">{error}</p>
            </div>
          {/if}

          {#if paymentStatus}
            <div class="mt-6 border border-gray-200 rounded-md p-4">
              <h3 class="text-lg font-medium text-gray-900">Payment Details</h3>
              <dl class="mt-2 text-sm text-gray-500">
                <div class="mt-1 flex justify-between">
                  <dt>Status:</dt>
                  <dd class="text-red-600 font-medium">
                    {paymentStatus.status}
                  </dd>
                </div>
                <div class="mt-1 flex justify-between">
                  <dt>Reference ID:</dt>
                  <dd>{paymentStatus.transaction_id || "N/A"}</dd>
                </div>
                <div class="mt-1 flex justify-between">
                  <dt>Amount:</dt>
                  <dd>KES {paymentStatus.amount.toFixed(2)}</dd>
                </div>
              </dl>
            </div>
          {/if}

          <div class="mt-6 space-y-4">
            <p class="text-gray-500">What would you like to do?</p>
            <div class="flex justify-center space-x-4">
              <button
                class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-primary-600 hover:bg-primary-700"
                on:click={() => goto("/checkout")}
              >
                Try Again
              </button>
              <button
                class="inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md shadow-sm text-gray-700 bg-white hover:bg-gray-50"
                on:click={handleRedirect}
              >
                View Order Status
              </button>
            </div>
          </div>

          <div class="mt-8">
            <p class="text-sm text-gray-500">
              You'll be redirected to order status in {redirectCountdown} seconds
            </p>
            <div
              class="mt-2 h-1 w-full bg-gray-200 rounded-full overflow-hidden"
            >
              <div
                class="h-full bg-primary-600 rounded-full"
                style="width: {(redirectCountdown / 15) * 100}%;"
              ></div>
            </div>
          </div>

          {#if orderId}
            <div
              class="mt-6 text-sm text-gray-500 border-t border-gray-200 pt-4"
            >
              <p>Your order (#{orderId.substring(0, 8)}) is still saved.</p>
              <p class="mt-1">
                If you continue to experience issues, please contact our support
                team.
              </p>
            </div>
          {/if}
        {/if}
      </div>
    </main>
  </div>
</div>
