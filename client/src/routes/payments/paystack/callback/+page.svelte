<script lang="ts">
  import { onMount } from "svelte";
  import { goto } from "$app/navigation";
  import toast from "svelte-french-toast";

  // Check URL parameters for payment information
  onMount(() => {
    const url = new URL(window.location.href);
    const reference = url.searchParams.get("reference");
    const status = url.searchParams.get("status");
    const trxref = url.searchParams.get("trxref");
    const orderId = localStorage.getItem("lastOrderId");

    // If there's a reference, it means we're returning from Paystack
    if (reference && status) {
      if (status === "success") {
        toast.success("Payment successful!");

        // If we have an order ID stored, redirect to the order confirmation page
        if (orderId) {
          goto(`/order-confirmation/${orderId}`);
        } else {
          // Otherwise, go to the orders page
          goto("/profile?tab=orders");
        }
      } else {
        toast.error("Payment was not successful. Please try again.");

        // If we have an order ID stored, redirect to the order confirmation page
        if (orderId) {
          goto(`/order-confirmation/${orderId}`);
        }
      }
    }
  });
</script>

<svelte:head>
  <title>Payment Processing | NyumbaniCare</title>
  <meta name="description" content="Processing your payment" />
</svelte:head>

<div
  class="min-h-screen bg-white px-4 py-16 sm:px-6 sm:py-24 md:grid md:place-items-center lg:px-8"
>
  <div class="mx-auto max-w-max">
    <main class="sm:flex">
      <div class="text-center">
        <h1
          class="text-4xl font-bold tracking-tight text-primary-600 sm:text-5xl"
        >
          Processing Payment
        </h1>
        <p class="mt-4 text-base text-gray-500">
          We are processing your payment. Please wait...
        </p>
        <div class="mt-10 flex justify-center space-x-3 text-center">
          <div class="flex items-center justify-center">
            <div
              class="h-8 w-8 animate-spin rounded-full border-4 border-primary-200 border-t-primary-600"
            ></div>
          </div>
        </div>
        <p class="mt-4 text-sm text-gray-500">
          You will be redirected automatically once the payment is processed.
        </p>
      </div>
    </main>
  </div>
</div>
