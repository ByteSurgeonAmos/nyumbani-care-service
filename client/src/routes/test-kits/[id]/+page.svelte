<script lang="ts">
  import { testKitService, isAuthenticated } from "$lib/api";
  import type { TestKit } from "$lib/api";
  import { onMount } from "svelte";
  import toast from "svelte-french-toast";
  import { goto } from "$app/navigation";
  import { cart } from "$lib/stores/cartStore";
  import TestKitOrderCallout from "$lib/components/TestKitOrderCallout.svelte";

  export let data: { id: string };
  const { id } = data;

  let testKit: TestKit | null = null;
  let isLoading = true;
  let quantity = 1;
  let isOrderLoading = false;

  onMount(async () => {
    await loadTestKit();
  });
  async function loadTestKit() {
    isLoading = true;

    try {
      testKit = await testKitService.getById(id);
      if (!testKit) {
        throw new Error("Test kit not found");
      }
    } catch (err) {
      console.error("Failed to load test kit:", err);
      toast.error("Failed to load test kit details");
      testKit = null;
    } finally {
      isLoading = false;
    }
  }
  function incrementQuantity() {
    if (testKit && testKit.in_stock) {
      quantity++;
    }
  }

  function decrementQuantity() {
    if (quantity > 1) {
      quantity--;
    }
  }
  function addToCart() {
    if (!testKit || !testKit.in_stock) {
      toast.error("This test kit is currently out of stock");
      return;
    }

    isOrderLoading = true;

    try {
      // Add item to cart store
      cart.addItem(testKit, quantity);
      toast.success("Added to cart successfully");

      // Show the success message for a moment before going to cart
      setTimeout(() => {
        goto("/cart");
      }, 1000);
    } catch (err) {
      console.error("Failed to add to cart:", err);
      toast.error("Failed to add to cart");
    } finally {
      isOrderLoading = false;
    }
  }
</script>

<svelte:head>
  {#if testKit}
    <title>{testKit.name} | NyumbaniCare</title>
    <meta
      name="description"
      content={testKit.description
        ? testKit.description.substring(0, 160)
        : "Test Kit Details"}
    />
  {:else}
    <title>Test Kit Details | NyumbaniCare</title>
  {/if}
</svelte:head>

<div class="bg-white">
  {#if isLoading}
    <div class="mx-auto max-w-7xl">
      <div class="pt-6 pb-16 sm:pb-24">
        <!-- Skeleton Loader for Navigation -->
        <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
          <div class="h-4 w-32 bg-gray-200 rounded animate-pulse mb-8"></div>
        </div>

        <div class="mx-auto mt-8 max-w-2xl px-4 sm:px-6 lg:max-w-7xl lg:px-8">
          <div class="lg:grid lg:grid-cols-2 lg:items-start lg:gap-x-8">
            <!-- Skeleton for image -->
            <div class="aspect-h-1 aspect-w-1 w-full">
              <div
                class="overflow-hidden rounded-lg bg-gray-200 animate-pulse h-96"
              ></div>
            </div>

            <!-- Skeleton for product info -->
            <div class="mt-10 px-4 sm:mt-16 sm:px-0 lg:mt-0">
              <!-- Title -->
              <div
                class="h-8 w-3/4 bg-gray-200 rounded animate-pulse mb-3"
              ></div>

              <!-- Price -->
              <div
                class="h-6 w-24 bg-gray-200 rounded animate-pulse mt-4"
              ></div>

              <!-- Availability -->
              <div
                class="h-4 w-32 bg-gray-200 rounded animate-pulse mt-4"
              ></div>

              <!-- Description -->
              <div class="mt-6">
                <div
                  class="h-4 w-full bg-gray-200 rounded animate-pulse mb-2"
                ></div>
                <div
                  class="h-4 w-full bg-gray-200 rounded animate-pulse mb-2"
                ></div>
                <div class="h-4 w-3/4 bg-gray-200 rounded animate-pulse"></div>
              </div>

              <!-- Quantity and add to cart -->
              <div class="mt-10">
                <div
                  class="h-10 w-full bg-gray-200 rounded animate-pulse"
                ></div>
              </div>

              <!-- Details section -->
              <div class="mt-10">
                <div
                  class="h-6 w-32 bg-gray-200 rounded animate-pulse mb-4"
                ></div>

                <div class="border-t border-gray-200 pt-4">
                  <div
                    class="h-4 w-32 bg-gray-200 rounded animate-pulse mb-4"
                  ></div>
                  <div
                    class="h-4 w-full bg-gray-200 rounded animate-pulse mb-2"
                  ></div>
                  <div
                    class="h-4 w-full bg-gray-200 rounded animate-pulse mb-2"
                  ></div>
                </div>

                <div class="border-t border-gray-200 pt-4">
                  <div
                    class="h-4 w-24 bg-gray-200 rounded animate-pulse mb-2"
                  ></div>
                  <div class="h-4 w-32 bg-gray-200 rounded animate-pulse"></div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  {:else if !testKit}
    <div class="text-center py-12">
      <p class="text-lg text-gray-500">Test kit not found</p>
      <a
        href="/test-kits"
        class="mt-4 inline-flex items-center text-primary-600 hover:text-primary-700"
      >
        Back to all test kits
      </a>
    </div>
  {:else}
    <div class="mx-auto max-w-7xl">
      <div class="pt-6 pb-16 sm:pb-24">
        <nav
          aria-label="Breadcrumb"
          class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8"
        >
          <ol role="list" class="flex items-center space-x-4">
            <li>
              <div class="flex items-center">
                <a
                  href="/test-kits"
                  class="mr-4 text-sm font-medium text-gray-900">Test Kits</a
                >
                <svg
                  viewBox="0 0 6 20"
                  class="h-5 w-auto text-gray-300"
                  fill="currentColor"
                  aria-hidden="true"
                >
                  <path d="M4.878 4.34H3.551L.27 16.532h1.327l3.281-12.19z" />
                </svg>
              </div>
            </li>
            <li class="text-sm">
              <a
                href={`/test-kits/${testKit.id}`}
                aria-current="page"
                class="font-medium text-gray-500 hover:text-gray-600"
              >
                {testKit.name}
              </a>
            </li>
          </ol>
        </nav>

        <div class="mx-auto mt-8 max-w-2xl px-4 sm:px-6 lg:max-w-7xl lg:px-8">
          <div class="lg:grid lg:grid-cols-2 lg:items-start lg:gap-x-8">
            <!-- Image gallery -->
            <div class="aspect-h-1 aspect-w-1 w-full">
              <div class="overflow-hidden rounded-lg">
                <img
                  src={testKit.image_url ||
                    "https://via.placeholder.com/600x600?text=Test+Kit"}
                  alt={testKit.name}
                  class="h-full w-full object-cover object-center"
                />
              </div>
            </div>

            <!-- Product info -->
            <div class="mt-10 px-4 sm:mt-16 sm:px-0 lg:mt-0">
              <h1 class="text-3xl font-bold tracking-tight text-gray-900">
                {testKit.name}
              </h1>
              <div class="mt-3">
                <h2 class="sr-only">Product information</h2>
                <p class="text-3xl tracking-tight text-gray-900">
                  KES {typeof testKit.price === "number"
                    ? testKit.price.toFixed(2)
                    : "0.00"}
                </p>
              </div>

              <div class="mt-2">
                <div class="flex items-center">
                  {#if testKit.in_stock}
                    <div class="flex items-center">
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
                      <p class="ml-2 text-sm text-green-600">
                        In stock and ready to ship
                      </p>
                    </div>
                  {:else}
                    <div class="flex items-center">
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
                      <p class="ml-2 text-sm text-red-600">
                        Currently out of stock
                      </p>
                    </div>
                  {/if}
                </div>
              </div>

              <div class="mt-6">
                <h3 class="sr-only">Description</h3>
                <div class="space-y-6 text-base text-gray-700">
                  <p>{testKit.description}</p>
                </div>
              </div>

              <div class="mt-6">
                <div class="mt-10 flex">
                  <div class="mr-4">
                    <label
                      for="quantity"
                      class="block text-sm font-medium text-gray-700"
                      >Quantity</label
                    >
                    <div class="mt-1 flex rounded-md shadow-sm">
                      <button
                        type="button"
                        class="relative inline-flex items-center rounded-l-md bg-gray-50 px-3 py-2 text-gray-500 ring-1 ring-inset ring-gray-300 hover:bg-gray-100"
                        on:click={decrementQuantity}
                        disabled={quantity <= 1}
                      >
                        <span class="text-gray-500">-</span>
                      </button>
                      <div
                        class="w-16 border-y border-gray-300 text-center py-2"
                      >
                        {quantity}
                      </div>
                      <button
                        type="button"
                        class="relative inline-flex items-center rounded-r-md bg-gray-50 px-3 py-2 text-gray-500 ring-1 ring-inset ring-gray-300 hover:bg-gray-100"
                        on:click={incrementQuantity}
                        disabled={!testKit.in_stock}
                      >
                        <span class="text-gray-500">+</span>
                      </button>
                    </div>
                  </div>
                  <button
                    type="button"
                    class="flex-1 rounded-md bg-primary-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-primary-700 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-600 disabled:opacity-50 disabled:cursor-not-allowed"
                    on:click={addToCart}
                    disabled={isOrderLoading || !testKit.in_stock}
                  >
                    {#if isOrderLoading}
                      <span>Processing...</span>
                    {:else}
                      <span>Add to Cart</span>
                    {/if}
                  </button>
                </div>
              </div>

              <!-- Test kit details -->
              <div class="mt-10">
                <h3 class="text-lg font-medium text-gray-900">Details</h3>

                <div class="mt-4 space-y-6">
                  <div class="border-t border-gray-200 pt-4">
                    <h4 class="font-medium text-gray-900">Instructions</h4>
                    <div class="prose prose-sm mt-2 text-gray-500">
                      <p>{testKit.instructions}</p>
                    </div>
                  </div>
                  <div class="border-t border-gray-200 pt-4">
                    <h4 class="font-medium text-gray-900">Category</h4>
                    <p class="mt-2 text-sm text-gray-500">{testKit.category}</p>
                  </div>

                  {#if testKit.reference_range}
                    <div class="border-t border-gray-200 pt-4">
                      <h4 class="font-medium text-gray-900">Reference Range</h4>
                      <p class="mt-2 text-sm text-gray-500">
                        {testKit.reference_range}
                      </p>
                    </div>
                  {/if}

                  {#if testKit.manufacturer}
                    <div class="border-t border-gray-200 pt-4">
                      <h4 class="font-medium text-gray-900">Manufacturer</h4>
                      <p class="mt-2 text-sm text-gray-500">
                        {testKit.manufacturer}
                      </p>
                    </div>
                  {/if}

                  {#if testKit.expiry_date}
                    <div class="border-t border-gray-200 pt-4">
                      <h4 class="font-medium text-gray-900">Expiry Date</h4>
                      <p class="mt-2 text-sm text-gray-500">
                        {new Date(testKit.expiry_date).toLocaleDateString()}
                      </p>
                    </div>
                  {/if}

                  {#if testKit.requires_analysis}
                    <div class="border-t border-gray-200 pt-4">
                      <h4 class="font-medium text-gray-900">Result Analysis</h4>
                      <p class="mt-2 text-sm text-gray-500">
                        This test kit requires analysis of results.
                        <a
                          href="/test-kits/analyze"
                          class="text-primary-600 hover:text-primary-700"
                        >
                          Upload your results for AI analysis
                        </a>.
                      </p>
                    </div>
                  {/if}
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
    {#if testKit && !isLoading}
      <TestKitOrderCallout />
    {/if}
  {/if}
</div>
