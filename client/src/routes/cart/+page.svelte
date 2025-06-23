<script lang="ts">
  import {
    cart,
    cartTotal,
    cartItemCount,
    type CartItem,
  } from "$lib/stores/cartStore";
  import { isAuthenticated } from "$lib/api";
  import { goto } from "$app/navigation";
  import { fly, fade } from "svelte/transition";
  import { onMount } from "svelte";
  import toast from "svelte-french-toast";

  let items: CartItem[] = [];
  let total = 0;
  let count = 0;

  // Subscribe to cart store
  onMount(() => {
    const unsubscribe = cart.subscribe((value) => {
      items = value;
    });

    const unsubscribeTotal = cartTotal.subscribe((value) => {
      total = value;
    });

    const unsubscribeCount = cartItemCount.subscribe((value) => {
      count = value;
    });

    return () => {
      unsubscribe();
      unsubscribeTotal();
      unsubscribeCount();
    };
  });

  function incrementQuantity(testKitId: string) {
    const item = items.find((item) => item.testKit.id === testKitId);
    if (item && item.testKit.in_stock) {
      cart.updateQuantity(testKitId, item.quantity + 1);
    }
  }

  function decrementQuantity(testKitId: string) {
    const item = items.find((item) => item.testKit.id === testKitId);
    if (item && item.quantity > 1) {
      cart.updateQuantity(testKitId, item.quantity - 1);
    }
  }

  function removeItem(testKitId: string) {
    cart.removeItem(testKitId);
    toast.success("Item removed from cart");
  }

  function proceedToCheckout() {
    if (!$isAuthenticated) {
      toast.error("Please login to continue with checkout");
      goto("/login?redirectTo=/checkout");
      return;
    }

    if (count === 0) {
      toast.error("Your cart is empty");
      return;
    }

    goto("/checkout");
  }
</script>

<svelte:head>
  <title>Shopping Cart | NyumbaniCare</title>
  <meta
    name="description"
    content="Your shopping cart with selected test kits"
  />
</svelte:head>

<div class="bg-white">
  <div class="mx-auto max-w-2xl px-4 pb-24 pt-16 sm:px-6 lg:max-w-7xl lg:px-8">
    <h1 class="text-3xl font-bold tracking-tight text-gray-900 sm:text-4xl">
      Shopping Cart
    </h1>

    {#if count === 0}
      <div
        class="mt-12 lg:grid lg:grid-cols-12 lg:items-start lg:gap-x-12 xl:gap-x-16"
      >
        <div class="lg:col-span-12" in:fade={{ duration: 300 }}>
          <div class="flex flex-col items-center py-12">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              class="h-24 w-24 text-gray-400"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="1"
                d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 11-4 0 2 2 0 014 0z"
              />
            </svg>
            <p class="mt-4 text-lg font-medium text-gray-900">
              Your cart is empty
            </p>
            <p class="mt-1 text-sm text-gray-500">
              Add some test kits to get started!
            </p>
            <a
              href="/test-kits"
              class="mt-6 rounded-md bg-primary-600 px-3.5 py-2.5 text-sm font-semibold text-white shadow-sm hover:bg-primary-500"
            >
              Browse Test Kits
            </a>
          </div>
        </div>
      </div>
    {:else}
      <div
        class="mt-12 lg:grid lg:grid-cols-12 lg:items-start lg:gap-x-12 xl:gap-x-16"
      >
        <section aria-labelledby="cart-heading" class="lg:col-span-7">
          <h2 id="cart-heading" class="sr-only">Items in your shopping cart</h2>

          <ul
            role="list"
            class="divide-y divide-gray-200 border-b border-t border-gray-200"
          >
            {#each items as item (item.testKit.id)}
              <li class="flex py-6 sm:py-10" in:fly={{ y: 20, duration: 300 }}>
                <div class="flex-shrink-0">
                  <img
                    src={item.testKit.image_url ||
                      "https://via.placeholder.com/200x200?text=Test+Kit"}
                    alt={item.testKit.name}
                    class="h-24 w-24 rounded-md object-cover object-center sm:h-48 sm:w-48"
                  />
                </div>

                <div class="ml-4 flex flex-1 flex-col justify-between sm:ml-6">
                  <div
                    class="relative pr-9 sm:grid sm:grid-cols-2 sm:gap-x-6 sm:pr-0"
                  >
                    <div>
                      <div class="flex justify-between">
                        <h3 class="text-sm">
                          <a
                            href={`/test-kits/${item.testKit.id}`}
                            class="font-medium text-gray-700 hover:text-gray-800"
                          >
                            {item.testKit.name}
                          </a>
                        </h3>
                      </div>
                      <div class="mt-1 flex text-sm">
                        <p class="text-gray-500">
                          Category: {item.testKit.category}
                        </p>
                      </div>
                      <p class="mt-1 text-sm font-medium text-gray-900">
                        KES {item.testKit.price.toFixed(2)}
                      </p>
                    </div>

                    <div class="mt-4 sm:mt-0 sm:pr-9">
                      <label
                        for={`quantity-${item.testKit.id}`}
                        class="sr-only"
                      >
                        Quantity, {item.testKit.name}
                      </label>
                      <div class="flex items-center">
                        <button
                          type="button"
                          class="text-gray-500 hover:text-gray-600"
                          on:click={() => decrementQuantity(item.testKit.id)}
                          disabled={item.quantity <= 1}
                        >
                          <svg
                            class="h-5 w-5"
                            viewBox="0 0 20 20"
                            fill="currentColor"
                            aria-hidden="true"
                          >
                            <path
                              fill-rule="evenodd"
                              d="M3 10a.75.75 0 01.75-.75h10.5a.75.75 0 010 1.5H3.75A.75.75 0 013 10z"
                              clip-rule="evenodd"
                            />
                          </svg>
                        </button>
                        <input
                          id={`quantity-${item.testKit.id}`}
                          name={`quantity-${item.testKit.id}`}
                          value={item.quantity}
                          disabled
                          class="mx-2 w-10 rounded-md border-0 py-1.5 text-center text-gray-900 ring-1 ring-inset ring-gray-300"
                        />
                        <button
                          type="button"
                          class="text-gray-500 hover:text-gray-600"
                          on:click={() => incrementQuantity(item.testKit.id)}
                          disabled={!item.testKit.in_stock}
                        >
                          <svg
                            class="h-5 w-5"
                            viewBox="0 0 20 20"
                            fill="currentColor"
                            aria-hidden="true"
                          >
                            <path
                              fill-rule="evenodd"
                              d="M10 3a.75.75 0 01.75.75v5.5h5.5a.75.75 0 010 1.5h-5.5v5.5a.75.75 0 01-1.5 0v-5.5h-5.5a.75.75 0 010-1.5h5.5v-5.5A.75.75 0 0110 3z"
                              clip-rule="evenodd"
                            />
                          </svg>
                        </button>
                      </div>

                      <div class="absolute right-0 top-0">
                        <button
                          type="button"
                          class="-m-2 inline-flex p-2 text-gray-400 hover:text-gray-500"
                          on:click={() => removeItem(item.testKit.id)}
                        >
                          <span class="sr-only">Remove</span>
                          <svg
                            class="h-5 w-5"
                            viewBox="0 0 20 20"
                            fill="currentColor"
                            aria-hidden="true"
                          >
                            <path
                              fill-rule="evenodd"
                              d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z"
                              clip-rule="evenodd"
                            />
                          </svg>
                        </button>
                      </div>
                    </div>
                  </div>

                  <p class="mt-4 flex space-x-2 text-sm text-gray-700">
                    {#if item.testKit.in_stock}
                      <svg
                        class="h-5 w-5 flex-shrink-0 text-green-500"
                        viewBox="0 0 20 20"
                        fill="currentColor"
                        aria-hidden="true"
                      >
                        <path
                          fill-rule="evenodd"
                          d="M16.704 4.153a.75.75 0 01.143 1.052l-8 10.5a.75.75 0 01-1.127.075l-4.5-4.5a.75.75 0 011.06-1.06l3.894 3.893 7.48-9.817a.75.75 0 011.05-.143z"
                          clip-rule="evenodd"
                        />
                      </svg>
                      <span>In stock</span>
                    {:else}
                      <svg
                        class="h-5 w-5 flex-shrink-0 text-red-500"
                        viewBox="0 0 20 20"
                        fill="currentColor"
                        aria-hidden="true"
                      >
                        <path
                          fill-rule="evenodd"
                          d="M4 10a.75.75 0 01.75-.75h10.5a.75.75 0 010 1.5H4.75A.75.75 0 014 10z"
                          clip-rule="evenodd"
                        />
                      </svg>
                      <span>Out of stock</span>
                    {/if}
                  </p>
                </div>
              </li>
            {/each}
          </ul>
        </section>

        <!-- Order summary -->
        <section
          aria-labelledby="summary-heading"
          class="mt-16 rounded-lg bg-gray-50 px-4 py-6 sm:p-6 lg:col-span-5 lg:mt-0 lg:p-8"
          in:fly={{ y: 20, duration: 300, delay: 200 }}
        >
          <h2 id="summary-heading" class="text-lg font-medium text-gray-900">
            Order summary
          </h2>

          <dl class="mt-6 space-y-4">
            <div class="flex items-center justify-between">
              <dt class="text-sm text-gray-600">Subtotal</dt>
              <dd class="text-sm font-medium text-gray-900">
                KES {total.toFixed(2)}
              </dd>
            </div>
            <div
              class="flex items-center justify-between border-t border-gray-200 pt-4"
            >
              <dt class="flex text-sm text-gray-600">
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

          <div class="mt-6">
            <button
              type="button"
              class="w-full rounded-md border border-transparent bg-primary-600 px-4 py-3 text-base font-medium text-white shadow-sm hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 focus:ring-offset-gray-50"
              on:click={proceedToCheckout}
            >
              Checkout
            </button>
          </div>

          <div class="mt-6 text-center text-sm text-gray-500">
            <p>
              or
              <a
                href="/test-kits"
                class="font-medium text-primary-600 hover:text-primary-500"
              >
                Continue Shopping
                <span aria-hidden="true"> &rarr;</span>
              </a>
            </p>
          </div>
        </section>
      </div>
    {/if}
  </div>
</div>
