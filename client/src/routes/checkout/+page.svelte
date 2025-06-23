<script lang="ts">
  import { cart, cartTotal, cartItemCount } from "$lib/stores/cartStore";
  import { isAuthenticated, currentUser } from "$lib/api";
  import { paymentService, checkoutService } from "$lib/api/paymentService";
  import type { CheckoutFormData } from "$lib/api/paymentService";
  import { goto } from "$app/navigation";
  import { onMount } from "svelte";
  import toast from "svelte-french-toast";

  let total = 0;
  let count = 0;
  let isLoading = false;

  // Form data for checkout
  let formData: CheckoutFormData = {
    email: "",
    fullName: "",
    phoneNumber: "",
    address: "",
    city: "",
    state: "",
    zipCode: "",
    country: "Kenya", // Default to Kenya
    paymentMethod: "paystack", // Default to Paystack
  };
  // Subscribe to cart stores
  onMount(() => {
    const unsubscribeTotal = cartTotal.subscribe((value) => {
      total = value;
    });

    const unsubscribeCount = cartItemCount.subscribe((value) => {
      count = value;
    });

    // Fill user data if authenticated
    if ($isAuthenticated && $currentUser) {
      formData.email = $currentUser.email || "";
      formData.fullName =
        `${$currentUser.first_name || ""} ${$currentUser.last_name || ""}`.trim();
      formData.phoneNumber = $currentUser.phone_number || "";
    }

    return () => {
      unsubscribeTotal();
      unsubscribeCount();
    };
  });
  async function handleSubmit() {
    if (count === 0) {
      toast.error("Your cart is empty");
      return;
    }

    // Basic validation
    if (
      !formData.email ||
      !formData.fullName ||
      !formData.address ||
      !formData.phoneNumber
    ) {
      toast.error("Please fill in all required fields");
      return;
    }

    isLoading = true;
    try {
      // Get cart items for order
      const items: { id: string; quantity: number }[] = [];
      let cartItems: any[] = [];

      // Get a snapshot of current cart items
      cart.subscribe((value) => {
        cartItems = value || [];
      })();

      if (!cartItems || cartItems.length === 0) {
        toast.error("Your cart is empty");
        isLoading = false;
        return;
      }

      // Validate inventory before creating order
      const inventoryValidation =
        await checkoutService.validateCartInventory(cartItems);
      if (!inventoryValidation.valid) {
        if (
          inventoryValidation.errors &&
          inventoryValidation.errors.length > 0
        ) {
          // Show specific inventory errors
          const errorMessages = inventoryValidation.errors.map(
            (error: any) =>
              `${error.name}: Only ${error.available} available (you requested ${error.requested})`
          );
          toast.error(
            `Some items are out of stock: ${errorMessages.join(", ")}`
          );
        } else {
          toast.error("Some items in your cart are no longer available");
        }

        // Wait a bit before redirecting to cart
        setTimeout(() => {
          goto("/cart");
        }, 3000);

        isLoading = false;
        return;
      }

      // Format items for API
      cartItems.forEach((item: any) => {
        items.push({
          id: item.testKit.id,
          quantity: item.quantity,
        });
      });

      // Format shipping address
      const shippingAddress = checkoutService.formatShippingAddress(formData);

      // Create order first
      const order = await checkoutService.createOrderWithItems({
        test_kit_ids: items,
        shipping_address: shippingAddress,
        payment_method: formData.paymentMethod,
        contact_number: formData.phoneNumber,
        email: formData.email,
      });

      toast.success("Order created successfully!"); // Then initiate payment
      const paymentResponse = await paymentService.initiatePayment({
        email: formData.email,
        amount: total + 300, // Add shipping
        order_id: order.id,
      });

      // Store payment ID for reference
      localStorage.setItem("lastPaymentId", paymentResponse.payment_id);
      localStorage.setItem("lastOrderId", order.id);

      // Clear cart after successful order
      cart.clearCart();

      // If there's a payment URL in the response (sometimes stored in other fields)
      const paymentUrl =
        paymentResponse.payment_url || paymentResponse.reference;

      if (paymentUrl && paymentUrl.startsWith("http")) {
        // Redirect to payment gateway
        window.location.href = paymentUrl;
      } else {
        // If no payment URL, go to order confirmation
        goto(`/order-confirmation/${order.id}`);
      }
    } catch (error) {
      console.error("Checkout error:", error);
      toast.error("Failed to process checkout");
    } finally {
      isLoading = false;
    }
  }
</script>

<svelte:head>
  <title>Checkout | NyumbaniCare</title>
  <meta name="description" content="Complete your test kit order" />
</svelte:head>

<div class="bg-white">
  <div class="mx-auto max-w-7xl px-4 pt-16 pb-24 sm:px-6 lg:px-8">
    <div class="mx-auto max-w-2xl lg:max-w-none">
      <h1 class="text-3xl font-bold tracking-tight text-gray-900 sm:text-4xl">
        Checkout
      </h1>

      {#if count === 0}
        <div class="mt-12 text-center">
          <p class="text-lg text-gray-500">Your cart is empty</p>
          <a
            href="/test-kits"
            class="mt-6 inline-flex items-center rounded-md border border-transparent bg-primary-600 px-4 py-2 text-sm font-medium text-white shadow-sm hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2"
          >
            Browse Test Kits
          </a>
        </div>
      {:else}
        <div
          class="mt-12 lg:grid lg:grid-cols-12 lg:items-start lg:gap-x-12 xl:gap-x-16"
        >
          <!-- Order summary -->
          <section
            aria-labelledby="order-heading"
            class="bg-gray-50 rounded-lg px-4 py-6 sm:p-6 lg:col-span-5 lg:p-8 lg:mt-0 mt-10"
          >
            <h2 id="order-heading" class="text-lg font-medium text-gray-900">
              Order summary
            </h2>

            <div class="mt-6 flow-root">
              <ul role="list" class="-my-6 divide-y divide-gray-200">
                {#each $cart as item}
                  <li class="flex py-6">
                    <div
                      class="h-24 w-24 flex-shrink-0 overflow-hidden rounded-md border border-gray-200"
                    >
                      <img
                        src={item.testKit.image_url ||
                          "https://via.placeholder.com/200x200?text=Test+Kit"}
                        alt={item.testKit.name}
                        class="h-full w-full object-cover object-center"
                      />
                    </div>

                    <div class="ml-4 flex flex-1 flex-col">
                      <div>
                        <div
                          class="flex justify-between text-base font-medium text-gray-900"
                        >
                          <h3>
                            <a href={`/test-kits/${item.testKit.id}`}
                              >{item.testKit.name}</a
                            >
                          </h3>
                          <p class="ml-4">
                            KES {item.testKit.price.toFixed(2)}
                          </p>
                        </div>
                        <p class="mt-1 text-sm text-gray-500">
                          {item.testKit.category}
                        </p>
                      </div>
                      <div
                        class="flex flex-1 items-end justify-between text-sm"
                      >
                        <p class="text-gray-500">Qty {item.quantity}</p>
                      </div>
                    </div>
                  </li>
                {/each}
              </ul>
            </div>

            <dl class="mt-8 space-y-4">
              <div class="flex items-center justify-between">
                <dt class="text-sm text-gray-600">Subtotal</dt>
                <dd class="text-sm font-medium text-gray-900">
                  KES {total.toFixed(2)}
                </dd>
              </div>
              <div
                class="flex items-center justify-between border-t border-gray-200 pt-4"
              >
                <dt class="flex items-center text-sm text-gray-600">
                  <span>Shipping estimate</span>
                </dt>
                <dd class="text-sm font-medium text-gray-900">KES 300.00</dd>
              </div>
              <div
                class="flex items-center justify-between border-t border-gray-200 pt-4"
              >
                <dt class="text-base font-medium text-gray-900">Order total</dt>
                <dd class="text-base font-medium text-gray-900">
                  KES {(total + 300).toFixed(2)}
                </dd>
              </div>
            </dl>
          </section>

          <!-- Contact information -->
          <div class="lg:col-span-7">
            <form
              class="space-y-8 divide-y divide-gray-200"
              on:submit|preventDefault={handleSubmit}
            >
              <div class="space-y-8 divide-y divide-gray-200">
                <div>
                  <h3 class="text-lg font-medium leading-6 text-gray-900">
                    Contact Information
                  </h3>
                  <div
                    class="mt-6 grid grid-cols-1 gap-y-6 gap-x-4 sm:grid-cols-6"
                  >
                    <div class="sm:col-span-6">
                      <label
                        for="email"
                        class="block text-sm font-medium text-gray-700"
                        >Email address</label
                      >
                      <div class="mt-1">
                        <input
                          type="email"
                          name="email"
                          id="email"
                          autocomplete="email"
                          required
                          class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500"
                          bind:value={formData.email}
                        />
                      </div>
                    </div>

                    <div class="sm:col-span-3">
                      <label
                        for="fullName"
                        class="block text-sm font-medium text-gray-700"
                        >Full name</label
                      >
                      <div class="mt-1">
                        <input
                          type="text"
                          name="fullName"
                          id="fullName"
                          autocomplete="given-name"
                          required
                          class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500"
                          bind:value={formData.fullName}
                        />
                      </div>
                    </div>

                    <div class="sm:col-span-3">
                      <label
                        for="phoneNumber"
                        class="block text-sm font-medium text-gray-700"
                        >Phone number</label
                      >
                      <div class="mt-1">
                        <input
                          type="tel"
                          name="phoneNumber"
                          id="phoneNumber"
                          autocomplete="tel"
                          required
                          class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500"
                          bind:value={formData.phoneNumber}
                        />
                      </div>
                    </div>
                  </div>
                </div>

                <div class="pt-8">
                  <h3 class="text-lg font-medium leading-6 text-gray-900">
                    Shipping address
                  </h3>

                  <div
                    class="mt-6 grid grid-cols-1 gap-y-6 gap-x-4 sm:grid-cols-6"
                  >
                    <div class="sm:col-span-6">
                      <label
                        for="address"
                        class="block text-sm font-medium text-gray-700"
                        >Street address</label
                      >
                      <div class="mt-1">
                        <input
                          type="text"
                          name="address"
                          id="address"
                          autocomplete="street-address"
                          required
                          class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500"
                          bind:value={formData.address}
                        />
                      </div>
                    </div>

                    <div class="sm:col-span-3">
                      <label
                        for="city"
                        class="block text-sm font-medium text-gray-700"
                        >City</label
                      >
                      <div class="mt-1">
                        <input
                          type="text"
                          name="city"
                          id="city"
                          autocomplete="address-level2"
                          required
                          class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500"
                          bind:value={formData.city}
                        />
                      </div>
                    </div>

                    <div class="sm:col-span-3">
                      <label
                        for="state"
                        class="block text-sm font-medium text-gray-700"
                        >State / Province</label
                      >
                      <div class="mt-1">
                        <input
                          type="text"
                          name="state"
                          id="state"
                          autocomplete="address-level1"
                          class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500"
                          bind:value={formData.state}
                        />
                      </div>
                    </div>

                    <div class="sm:col-span-3">
                      <label
                        for="zipCode"
                        class="block text-sm font-medium text-gray-700"
                        >ZIP / Postal code</label
                      >
                      <div class="mt-1">
                        <input
                          type="text"
                          name="zipCode"
                          id="zipCode"
                          autocomplete="postal-code"
                          class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500"
                          bind:value={formData.zipCode}
                        />
                      </div>
                    </div>

                    <div class="sm:col-span-3">
                      <label
                        for="country"
                        class="block text-sm font-medium text-gray-700"
                        >Country</label
                      >
                      <div class="mt-1">
                        <select
                          id="country"
                          name="country"
                          autocomplete="country-name"
                          class="block w-full rounded-md border-gray-300 shadow-sm focus:border-primary-500 focus:ring-primary-500"
                          bind:value={formData.country}
                        >
                          <option value="Kenya">Kenya</option>
                          <option value="Tanzania">Tanzania</option>
                          <option value="Uganda">Uganda</option>
                          <option value="Rwanda">Rwanda</option>
                          <option value="Other">Other</option>
                        </select>
                      </div>
                    </div>
                  </div>
                </div>

                <div class="pt-8">
                  <h3 class="text-lg font-medium leading-6 text-gray-900">
                    Payment method
                  </h3>
                  <div class="mt-6">
                    <div class="space-y-4">
                      <div class="flex items-center">
                        <input
                          id="paystack"
                          name="payment-method"
                          type="radio"
                          checked
                          class="h-4 w-4 border-gray-300 text-primary-600 focus:ring-primary-500"
                          bind:group={formData.paymentMethod}
                          value="paystack"
                        />
                        <label
                          for="paystack"
                          class="ml-3 block text-sm font-medium text-gray-700"
                        >
                          Paystack (Credit Card / Mobile Money)
                        </label>
                      </div>
                    </div>
                  </div>
                </div>
              </div>

              <div class="pt-5">
                <div class="flex justify-end">
                  <a
                    href="/cart"
                    class="rounded-md border border-gray-300 bg-white py-2 px-4 text-sm font-medium text-gray-700 shadow-sm hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2"
                  >
                    Back to cart
                  </a>
                  <button
                    type="submit"
                    class="ml-3 inline-flex justify-center rounded-md border border-transparent bg-primary-600 py-2 px-4 text-sm font-medium text-white shadow-sm hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 disabled:opacity-70 disabled:cursor-not-allowed"
                    disabled={isLoading}
                  >
                    {#if isLoading}
                      Processing...
                    {:else}
                      Complete Order
                    {/if}
                  </button>
                </div>
              </div>
            </form>
          </div>
        </div>
      {/if}
    </div>
  </div>
</div>
