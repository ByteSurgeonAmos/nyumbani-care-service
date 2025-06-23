<script lang="ts">
  import { onMount } from "svelte";
  import { goto } from "$app/navigation";
  import toast from "svelte-french-toast";

  onMount(() => {
    // Notify user
    toast.success("Payment successful!");

    // Get the last order ID from localStorage if available
    const orderId = localStorage.getItem("lastOrderId");

    // Redirect to order confirmation page after a short delay
    setTimeout(() => {
      if (orderId) {
        goto(`/order-confirmation/${orderId}`);
      } else {
        goto("/profile?tab=orders");
      }
    }, 2000);
  });
</script>

<svelte:head>
  <title>Payment Successful | NyumbaniCare</title>
  <meta name="description" content="Your payment was successful" />
</svelte:head>

<div
  class="min-h-screen bg-white px-4 py-16 sm:px-6 sm:py-24 md:grid md:place-items-center lg:px-8"
>
  <div class="mx-auto max-w-max">
    <main class="sm:flex">
      <div class="text-center">
        <div class="flex justify-center">
          <div
            class="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-green-100"
          >
            <svg
              class="h-8 w-8 text-green-600"
              fill="none"
              viewBox="0 0 24 24"
              stroke-width="1.5"
              stroke="currentColor"
              aria-hidden="true"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M4.5 12.75l6 6 9-13.5"
              />
            </svg>
          </div>
        </div>
        <h1
          class="mt-3 text-3xl font-bold tracking-tight text-gray-900 sm:text-4xl"
        >
          Payment Successful!
        </h1>
        <p class="mt-3 text-base text-gray-500">
          Thank you for your purchase. Your order has been processed.
        </p>
        <div class="mt-8 flex items-center justify-center">
          <div class="text-sm font-medium text-gray-500">
            Redirecting you to your order...
          </div>
        </div>
        <div class="mt-6">
          <div class="h-1 w-full bg-gray-200 rounded-full overflow-hidden">
            <div
              class="h-full bg-primary-600 rounded-full animate-progress"
            ></div>
          </div>
        </div>
      </div>
    </main>
  </div>
</div>

<style>
  @keyframes progress {
    0% {
      width: 0%;
    }
    100% {
      width: 100%;
    }
  }

  .animate-progress {
    animation: progress 2s ease-in-out forwards;
  }
</style>
