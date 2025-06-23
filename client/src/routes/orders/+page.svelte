<script lang="ts">
  import { onMount } from "svelte";
  import { goto } from "$app/navigation";
  import toast from "svelte-french-toast";
  import { orderService, type TestKitOrder } from "$lib/api";
  import { formatDate } from "$lib/utils/dates";

  let orders: TestKitOrder[] = [];
  let isLoading = true;
  let currentPage = 1;
  let totalItems = 0;
  let itemsPerPage = 10;
  let totalPages = 0;

  type OrderStatus =
    | "pending"
    | "processing"
    | "shipped"
    | "delivered"
    | "cancelled";

  const statusStyles: Record<OrderStatus, string> = {
    pending: "bg-yellow-100 text-yellow-800",
    processing: "bg-blue-100 text-blue-800",
    shipped: "bg-purple-100 text-purple-800",
    delivered: "bg-green-100 text-green-800",
    cancelled: "bg-red-100 text-red-800",
  };

  type PaymentStatus = "pending" | "completed" | "failed";

  const paymentStatusStyles: Record<PaymentStatus, string> = {
    pending: "bg-yellow-100 text-yellow-800",
    completed: "bg-green-100 text-green-800",
    failed: "bg-red-100 text-red-800",
  };

  onMount(async () => {
    loadOrders();
  });

  async function loadOrders() {
    isLoading = true;

    try {
      const response = await orderService.getAll(currentPage, itemsPerPage);
      orders = response.data;
      totalItems = response.total;
      totalPages = Math.ceil(totalItems / itemsPerPage);
    } catch (error) {
      console.error("Failed to load orders:", error);
      toast.error("Failed to load your orders");
    } finally {
      isLoading = false;
    }
  }

  function changePage(newPage: number) {
    if (newPage >= 1 && newPage <= totalPages) {
      currentPage = newPage;
      loadOrders();
    }
  }

  function getStatusClass(status: string): string {
    return (
      statusStyles[status.toLowerCase() as OrderStatus] ||
      "bg-gray-100 text-gray-800"
    );
  }

  function getPaymentStatusClass(status: string): string {
    return (
      paymentStatusStyles[status.toLowerCase() as PaymentStatus] ||
      "bg-gray-100 text-gray-800"
    );
  }

  function viewOrderDetails(orderId: string) {
    goto(`/order-confirmation/${orderId}`);
  }
</script>

<svelte:head>
  <title>My Orders | NyumbaniCare</title>
  <meta
    name="description"
    content="View your order history and track current orders with NyumbaniCare."
  />
</svelte:head>

<div class="bg-white">
  <div class="mx-auto max-w-7xl px-4 py-16 sm:px-6 lg:px-8">
    <div class="sm:flex sm:items-center sm:justify-between">
      <h1 class="text-3xl font-bold tracking-tight text-gray-900">My Orders</h1>
      <div class="mt-4 sm:mt-0">
        <a
          href="/test-kits"
          class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
        >
          Browse Test Kits
        </a>
      </div>
    </div>

    <div class="mt-8">
      {#if isLoading}
        <div class="flex justify-center items-center py-12">
          <div
            class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"
          ></div>
        </div>
      {:else if orders.length === 0}
        <div class="bg-white shadow overflow-hidden sm:rounded-lg">
          <div class="px-4 py-12 text-center sm:px-6">
            <p class="text-lg text-gray-500">
              You haven't placed any orders yet.
            </p>
            <div class="mt-6">
              <a
                href="/test-kits"
                class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
              >
                Browse Our Test Kits
              </a>
            </div>
          </div>
        </div>
      {:else}
        <div class="overflow-hidden bg-white shadow sm:rounded-md">
          <ul role="list" class="divide-y divide-gray-200">
            {#each orders as order}
              <li>
                <div class="block hover:bg-gray-50">
                  <div class="px-4 py-4 sm:px-6">
                    <div class="flex items-center justify-between">
                      <div
                        class="sm:flex sm:items-center space-y-2 sm:space-y-0 sm:space-x-4"
                      >
                        <p
                          class="text-sm font-medium text-primary-600 truncate"
                        >
                          Order #{order.id.substring(0, 8)}
                        </p>
                        <div class="flex">
                          <p class="flex items-center text-sm text-gray-500">
                            <span
                              class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full {getStatusClass(
                                order.status
                              )}"
                            >
                              {order.status}
                            </span>
                          </p>
                          <p
                            class="ml-2 flex items-center text-sm text-gray-500"
                          >
                            <span
                              class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full {getPaymentStatusClass(
                                order.payment_status
                              )}"
                            >
                              Payment: {order.payment_status}
                            </span>
                          </p>
                        </div>
                      </div>
                      <div class="ml-2 flex-shrink-0 flex">
                        <button
                          on:click={() => viewOrderDetails(order.id)}
                          class="font-medium text-primary-600 hover:text-primary-500"
                        >
                          View details
                        </button>
                      </div>
                    </div>
                    <div class="mt-2 sm:flex sm:justify-between">
                      <div class="sm:flex">
                        <p class="flex items-center text-sm text-gray-500">
                          <span class="truncate"
                            >Items: {order.items.length}</span
                          >
                        </p>
                        <p
                          class="mt-2 flex items-center text-sm text-gray-500 sm:mt-0 sm:ml-6"
                        >
                          <span class="font-medium"
                            >KES {order.total_amount.toFixed(2)}</span
                          >
                        </p>
                      </div>
                      <div
                        class="mt-2 flex items-center text-sm text-gray-500 sm:mt-0"
                      >
                        <span>Ordered on {formatDate(order.created_at)}</span>
                      </div>
                    </div>
                  </div>
                </div>
              </li>
            {/each}
          </ul>
        </div>
      {/if}

      <!-- Pagination -->
      {#if totalPages > 1}
        <div class="flex justify-center mt-8">
          <nav class="inline-flex -space-x-px rounded-md shadow-sm">
            <!-- Previous button -->
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

            <!-- Page numbers -->
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

            <!-- Next button -->
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
    </div>
  </div>
</div>
