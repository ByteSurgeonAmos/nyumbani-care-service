<script lang="ts">
  import { testKitService } from "$lib/api";
  import type { TestKit } from "$lib/api";
  import { onMount } from "svelte";
  import toast from "svelte-french-toast";
  import { cart } from "$lib/stores/cartStore";

  let testKits: TestKit[] = [];
  let isLoading = true;
  let currentPage = 1;
  let totalKits = 0;
  let kitsPerPage = 12;
  let totalPages = 0;
  let addingToCart: Record<string, boolean> = {};

  onMount(async () => {
    loadTestKits();
  });

  async function loadTestKits() {
    isLoading = true;

    try {
      const response = await testKitService.getAll(currentPage, kitsPerPage);
      testKits = response.data;
      console.log("Loaded test kits:", testKits);

      totalKits = response.total;
      totalPages = Math.ceil(totalKits / kitsPerPage);
    } catch (error) {
      console.error("Failed to load test kits:", error);
      toast.error("Failed to load test kits");
    } finally {
      isLoading = false;
    }
  }

  function addToCart(testKit: TestKit) {
    if (!testKit.in_stock) {
      toast.error("This test kit is currently out of stock");
      return;
    }

    addingToCart[testKit.id] = true;

    try {
      cart.addItem(testKit, 1);
      toast.success(`Added ${testKit.name} to cart`);
    } catch (err) {
      console.error("Failed to add to cart:", err);
      toast.error("Failed to add to cart");
    } finally {
      // Remove loading state after a short delay to show feedback
      setTimeout(() => {
        addingToCart[testKit.id] = false;
      }, 500);
    }
  }

  function changePage(newPage: number): void {
    if (newPage >= 1 && newPage <= totalPages && newPage !== currentPage) {
      currentPage = newPage;
      loadTestKits();
    }
  }
</script>

<svelte:head>
  <title>Test Kits | NyumbaniCare</title>
  <meta
    name="description"
    content="Browse our selection of at-home health test kits available for delivery."
  />
</svelte:head>

<div class="bg-white">
  <div
    class="mx-auto max-w-2xl px-4 py-16 sm:px-6 sm:py-24 lg:max-w-7xl lg:px-8"
  >
    <div class="flex justify-between items-start">
      <div>
        <h1 class="text-3xl font-bold tracking-tight text-gray-900">
          Test Kits
        </h1>
        <p class="mt-4 max-w-3xl text-base text-gray-500">
          Our diagnostic test kits allow you to test for various health
          conditions in the comfort of your home. All kits include detailed
          instructions and support from our healthcare team.
        </p>
      </div>
      <div class="flex gap-4">
        <a
          href="/cart"
          class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-primary-600 hover:bg-primary-700"
        >
          View Cart
        </a>
        <a
          href="/test-kits/analyze"
          class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-primary-600 hover:bg-primary-700"
        >
          Analyze Test Results
        </a>
      </div>
    </div>

    <!-- Filtering and sorting options could be added here -->

    {#if isLoading}
      <div
        class="mt-8 grid grid-cols-1 gap-y-12 sm:grid-cols-2 sm:gap-x-6 lg:grid-cols-4 xl:gap-x-8"
      >
        {#each Array(8) as _, i}
          <div class="group relative">
            <!-- Skeleton for image -->
            <div
              class="aspect-w-1 aspect-h-1 w-full overflow-hidden rounded-lg bg-gray-200 animate-pulse"
            ></div>
            <!-- Skeleton for product info -->
            <div class="mt-4 flex justify-between">
              <div>
                <div
                  class="h-5 w-36 bg-gray-200 rounded animate-pulse mb-2"
                ></div>
                <div class="h-4 w-24 bg-gray-200 rounded animate-pulse"></div>
              </div>
              <div class="h-5 w-16 bg-gray-200 rounded animate-pulse"></div>
            </div>
            <!-- Skeleton for description -->
            <div
              class="h-4 w-full bg-gray-200 rounded animate-pulse mt-2"
            ></div>
            <div class="h-4 w-3/4 bg-gray-200 rounded animate-pulse mt-1"></div>
            <!-- Skeleton for stock status -->
            <div class="h-5 w-16 bg-gray-200 rounded animate-pulse mt-2"></div>
          </div>
        {/each}
      </div>
    {:else if testKits.length === 0}
      <div class="text-center py-12">
        <p class="text-lg text-gray-500">
          No test kits available at the moment.
        </p>
      </div>
    {:else}
      <div
        class="mt-8 grid grid-cols-1 gap-y-12 sm:grid-cols-2 sm:gap-x-6 lg:grid-cols-4 xl:gap-x-8"
      >
        {#each testKits as kit (kit.id)}
          <div class="group relative">
            <div
              class="aspect-w-1 aspect-h-1 w-full overflow-hidden rounded-lg bg-gray-100"
            >
              <img
                src={kit.image_url ||
                  "https://via.placeholder.com/300x300?text=Test+Kit"}
                alt={kit.name}
                class="h-full w-full object-cover object-center group-hover:opacity-75"
              />
            </div>
            <div class="mt-4 flex justify-between">
              <div>
                <h3 class="text-lg font-medium text-gray-900">
                  <a href={`/test-kits/${kit.id}`}>
                    <span aria-hidden="true" class="absolute inset-0"></span>
                    {kit.name}
                  </a>
                </h3>
                <p class="mt-1 text-sm text-gray-500">{kit.category}</p>
              </div>
              <p class="text-lg font-medium text-primary-600">
                KES {kit.price.toFixed(2)}
              </p>
            </div>
            <p class="mt-1 text-sm text-gray-500">
              {kit.description.substring(0, 100)}...
            </p>
            <div class="mt-2 flex items-center justify-between">
              {#if kit.in_stock}
                <span
                  class="inline-flex items-center rounded-full bg-green-100 px-2.5 py-0.5 text-xs font-medium text-green-800"
                >
                  In Stock
                </span>
              {:else}
                <span
                  class="inline-flex items-center rounded-full bg-red-100 px-2.5 py-0.5 text-xs font-medium text-red-800"
                >
                  Out of Stock
                </span>
              {/if}

              <button
                on:click={() => addToCart(kit)}
                disabled={!kit.in_stock || addingToCart[kit.id]}
                class="text-sm font-medium {!kit.in_stock
                  ? 'text-gray-400 cursor-not-allowed'
                  : 'text-primary-600 hover:text-primary-800'} flex items-center"
              >
                {#if addingToCart[kit.id]}
                  <svg
                    class="animate-spin -ml-1 mr-2 h-4 w-4 text-primary-600"
                    xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                  >
                    <circle
                      class="opacity-25"
                      cx="12"
                      cy="12"
                      r="10"
                      stroke="currentColor"
                      stroke-width="4"
                    ></circle>
                    <path
                      class="opacity-75"
                      fill="currentColor"
                      d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                    ></path>
                  </svg>
                  Adding...
                {:else}
                  <svg
                    class="w-4 h-4 mr-1"
                    fill="currentColor"
                    viewBox="0 0 20 20"
                    xmlns="http://www.w3.org/2000/svg"
                  >
                    <path
                      d="M3 1a1 1 0 000 2h1.22l.305 1.222a.997.997 0 00.01.042l1.358 5.43-.893.892C3.74 11.846 4.632 14 6.414 14H15a1 1 0 000-2H6.414l1-1H14a1 1 0 00.894-.553l3-6A1 1 0 0017 3H6.28l-.31-1.243A1 1 0 005 1H3zM16 16.5a1.5 1.5 0 11-3 0 1.5 1.5 0 013 0zM6.5 18a1.5 1.5 0 100-3 1.5 1.5 0 000 3z"
                    ></path>
                  </svg>
                  Add to Cart
                {/if}
              </button>
            </div>
          </div>
        {/each}
      </div>

      <!-- Pagination -->
      {#if totalPages > 1}
        <div class="flex justify-center mt-10">
          <nav
            class="inline-flex -space-x-px rounded-md shadow-sm"
            aria-label="Pagination"
          >
            <button
              class="inline-flex items-center rounded-l-md px-2 py-2 text-gray-400 ring-1 ring-inset ring-gray-300 hover:bg-gray-50 focus:z-20 focus:outline-offset-0 {currentPage ===
              1
                ? 'cursor-not-allowed opacity-50'
                : ''}"
              on:click={() => changePage(currentPage - 1)}
              disabled={currentPage === 1}
            >
              <span class="sr-only">Previous</span>
              <svg
                class="h-5 w-5"
                viewBox="0 0 20 20"
                fill="currentColor"
                aria-hidden="true"
              >
                <path
                  fill-rule="evenodd"
                  d="M12.79 5.23a.75.75 0 01-.02 1.06L8.832 10l3.938 3.71a.75.75 0 11-1.04 1.08l-4.5-4.25a.75.75 0 010-1.08l4.5-4.25a.75.75 0 011.06.02z"
                  clip-rule="evenodd"
                />
              </svg>
            </button>

            {#each Array(totalPages) as _, i}
              <button
                class="inline-flex items-center px-4 py-2 text-sm font-semibold {currentPage ===
                i + 1
                  ? 'bg-primary-600 text-white'
                  : 'text-gray-900 ring-1 ring-inset ring-gray-300 hover:bg-gray-50'}"
                on:click={() => changePage(i + 1)}
              >
                {i + 1}
              </button>
            {/each}

            <button
              class="inline-flex items-center rounded-r-md px-2 py-2 text-gray-400 ring-1 ring-inset ring-gray-300 hover:bg-gray-50 focus:z-20 focus:outline-offset-0 {currentPage ===
              totalPages
                ? 'cursor-not-allowed opacity-50'
                : ''}"
              on:click={() => changePage(currentPage + 1)}
              disabled={currentPage === totalPages}
            >
              <span class="sr-only">Next</span>
              <svg
                class="h-5 w-5"
                viewBox="0 0 20 20"
                fill="currentColor"
                aria-hidden="true"
              >
                <path
                  fill-rule="evenodd"
                  d="M7.21 14.77a.75.75 0 01.02-1.06L11.168 10 7.23 6.29a.75.75 0 111.04-1.08l4.5 4.25a.75.75 0 010 1.08l-4.5 4.25a.75.75 0 01-1.06-.02z"
                  clip-rule="evenodd"
                />
              </svg>
            </button>
          </nav>
        </div>
      {/if}
    {/if}
  </div>
</div>
