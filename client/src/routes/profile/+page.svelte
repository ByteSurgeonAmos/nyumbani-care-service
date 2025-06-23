<script lang="ts">
  import {
    currentUser,
    authService,
    orderService,
    testKitResultService,
    type TestKitOrder,
  } from "$lib/api";
  import type { AnalyzeTestKitResponse } from "$lib/api/extendedServices";
  import { onMount } from "svelte";
  import toast from "svelte-french-toast";
  import { formatDate } from "$lib/utils/dates";
  import { goto } from "$app/navigation";

  let isLoading = false;
  let firstName = $currentUser?.first_name || "";
  let lastName = $currentUser?.last_name || "";
  let email = $currentUser?.email || "";
  let phoneNumber = $currentUser?.phone_number || "";
  let address = $currentUser?.address || "";

  let recentOrders: TestKitOrder[] = [];
  let isLoadingOrders = true;

  let recentResults: AnalyzeTestKitResponse[] = [];
  let isLoadingResults = true;

  onMount(async () => {
    try {
      await authService.getCurrentUser();

      // Update local variables with the latest user data
      firstName = $currentUser?.first_name || "";
      lastName = $currentUser?.last_name || "";
      email = $currentUser?.email || "";
      phoneNumber = $currentUser?.phone_number || "";
      address = $currentUser?.address || "";

      // Load recent orders and test results
      loadRecentOrders();
      loadRecentResults();
    } catch (error) {
      console.error("Failed to fetch user:", error);
    }
  });

  async function loadRecentOrders() {
    isLoadingOrders = true;
    try {
      // Fetch just the most recent 3 orders
      const response = await orderService.getAll(1, 3);
      recentOrders = response.data;
    } catch (error) {
      console.error("Failed to load recent orders:", error);
    } finally {
      isLoadingOrders = false;
    }
  }

  async function loadRecentResults() {
    isLoadingResults = true;
    try {
      // Fetch just the most recent 3 test results
      const response = await testKitResultService.getAll(1, 3);
      recentResults = response.data;
    } catch (error) {
      console.error("Failed to load recent test results:", error);
    } finally {
      isLoadingResults = false;
    }
  }

  function viewOrderDetails(orderId: string) {
    goto(`/order-confirmation/${orderId}`);
  }

  function viewTestResultDetails(resultId: string) {
    goto(`/test-result/${resultId}`);
  }

  async function handleUpdateProfile() {
    isLoading = true;

    try {
      await authService.updateProfile({
        first_name: firstName,
        last_name: lastName,
        phone_number: phoneNumber,
        address: address,
      });

      toast.success("Profile updated successfully");
    } catch (error) {
      console.error("Failed to update profile:", error);
      toast.error("Failed to update profile");
    } finally {
      isLoading = false;
    }
  }

  function getResultBadgeClass(result: string): string {
    switch (result.toLowerCase()) {
      case "positive":
        return "bg-red-100 text-red-800";
      case "negative":
        return "bg-green-100 text-green-800";
      case "inconclusive":
      default:
        return "bg-yellow-100 text-yellow-800";
    }
  }
</script>

<svelte:head>
  <title>My Profile | NyumbaniCare</title>
</svelte:head>

<div class="bg-gray-50 py-8">
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
    <h1 class="text-3xl font-bold text-gray-900 mb-6">My Profile</h1>

    <div class="bg-white rounded-lg shadow overflow-hidden">
      <div class="p-6">
        <form on:submit|preventDefault={handleUpdateProfile} class="space-y-6">
          <div class="grid grid-cols-1 gap-y-6 gap-x-4 sm:grid-cols-2">
            <div>
              <label
                for="first-name"
                class="block text-sm font-medium text-gray-700"
                >First name</label
              >
              <div class="mt-1">
                <input
                  id="first-name"
                  name="first-name"
                  type="text"
                  bind:value={firstName}
                  class="input"
                />
              </div>
            </div>

            <div>
              <label
                for="last-name"
                class="block text-sm font-medium text-gray-700">Last name</label
              >
              <div class="mt-1">
                <input
                  id="last-name"
                  name="last-name"
                  type="text"
                  bind:value={lastName}
                  class="input"
                />
              </div>
            </div>

            <div>
              <label for="email" class="block text-sm font-medium text-gray-700"
                >Email</label
              >
              <div class="mt-1">
                <input
                  id="email"
                  name="email"
                  type="email"
                  disabled
                  value={email}
                  class="input bg-gray-50"
                />
              </div>
              <p class="mt-1 text-sm text-gray-500">Email cannot be changed</p>
            </div>

            <div>
              <label
                for="phone-number"
                class="block text-sm font-medium text-gray-700"
                >Phone number</label
              >
              <div class="mt-1">
                <input
                  id="phone-number"
                  name="phone-number"
                  type="tel"
                  bind:value={phoneNumber}
                  class="input"
                />
              </div>
            </div>

            <div class="sm:col-span-2">
              <label
                for="address"
                class="block text-sm font-medium text-gray-700">Address</label
              >
              <div class="mt-1">
                <textarea
                  id="address"
                  name="address"
                  rows="3"
                  bind:value={address}
                  class="input"
                ></textarea>
              </div>
            </div>
          </div>

          <div class="flex justify-end">
            <button
              type="submit"
              disabled={isLoading}
              class="btn btn-primary py-2 px-4 {isLoading
                ? 'opacity-70 cursor-not-allowed'
                : ''}"
            >
              {isLoading ? "Saving..." : "Save Changes"}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- Recent Orders Section -->
    <div class="mt-8 bg-white rounded-lg shadow overflow-hidden">
      <div
        class="px-6 py-5 border-b border-gray-200 flex justify-between items-center"
      >
        <h2 class="text-lg font-medium text-gray-900">Recent Orders</h2>
        <a
          href="/orders"
          class="text-sm font-medium text-primary-600 hover:text-primary-500"
        >
          View all orders
        </a>
      </div>
      <div class="px-6 py-5">
        {#if isLoadingOrders}
          <div class="flex justify-center items-center py-8">
            <div
              class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"
            ></div>
          </div>
        {:else if recentOrders.length === 0}
          <p class="text-sm text-gray-500 py-4">
            You haven't placed any orders yet.
            <a href="/test-kits" class="text-primary-600 hover:text-primary-500"
              >Browse our test kits</a
            >
            to place your first order.
          </p>
        {:else}
          <div class="divide-y divide-gray-200">
            {#each recentOrders as order}
              <div class="py-4 flex justify-between items-center">
                <div>
                  <p class="text-sm font-medium text-gray-900">
                    Order #{order.id.substring(0, 8)}
                  </p>
                  <p class="text-sm text-gray-500">
                    {formatDate(order.created_at)} · {order.items.length} items ·
                    KES {order.total_amount.toFixed(2)}
                  </p>
                  <div class="mt-1">
                    <span
                      class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium
                      {order.status === 'delivered'
                        ? 'bg-green-100 text-green-800'
                        : order.status === 'cancelled'
                          ? 'bg-red-100 text-red-800'
                          : 'bg-yellow-100 text-yellow-800'}"
                    >
                      {order.status.charAt(0).toUpperCase() +
                        order.status.slice(1)}
                    </span>
                    <span
                      class="ml-2 inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium
                      {order.payment_status === 'completed'
                        ? 'bg-green-100 text-green-800'
                        : order.payment_status === 'failed'
                          ? 'bg-red-100 text-red-800'
                          : 'bg-yellow-100 text-yellow-800'}"
                    >
                      Payment: {order.payment_status.charAt(0).toUpperCase() +
                        order.payment_status.slice(1)}
                    </span>
                  </div>
                </div>
                <button
                  on:click={() => viewOrderDetails(order.id)}
                  class="text-sm font-medium text-primary-600 hover:text-primary-500"
                >
                  View details
                </button>
              </div>
            {/each}
          </div>
        {/if}
      </div>
    </div>

    <!-- Recent Test Results Section -->
    <div class="mt-8 bg-white rounded-lg shadow overflow-hidden">
      <div
        class="px-6 py-5 border-b border-gray-200 flex justify-between items-center"
      >
        <h2 class="text-lg font-medium text-gray-900">Recent Test Results</h2>
        <a
          href="/test-kits/analyze/results"
          class="text-sm font-medium text-primary-600 hover:text-primary-500"
          >View all</a
        >
      </div>
      <div class="px-6 py-5">
        {#if isLoadingResults}
          <div class="flex justify-center items-center py-8">
            <div
              class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"
            ></div>
          </div>
        {:else if recentResults.length === 0}
          <p class="text-sm text-gray-500 py-4">
            You haven't analyzed any test kits yet.
            <a
              href="/test-kits/analyze"
              class="text-primary-600 hover:text-primary-500"
              >Analyze a test kit</a
            >
            to get your first result.
          </p>
        {:else}
          <div class="divide-y divide-gray-200">
            {#each recentResults as result}
              <div class="py-4 flex justify-between items-center">
                <div>
                  <p class="text-sm font-medium text-gray-900">
                    Test Result #{result.id.substring(0, 8)}
                  </p>
                  <p class="text-sm text-gray-500">
                    {result.created_at
                      ? formatDate(result.created_at)
                      : "Unknown date"}
                  </p>
                  <div class="mt-1">
                    <span
                      class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium
                      {getResultBadgeClass(result.result)}"
                    >
                      {result.result.toUpperCase()}
                    </span>
                    <span class="ml-2 text-xs text-gray-500">
                      Confidence: {(result.confidence * 100).toFixed(0)}%
                    </span>
                  </div>
                </div>
                <div>
                  <a
                    href={`/test-kits/analyze/results/${result.id}`}
                    class="text-sm font-medium text-primary-600 hover:text-primary-500"
                    >View details</a
                  >
                </div>
              </div>
            {/each}
          </div>
        {/if}
      </div>
    </div>

    <div class="mt-8 bg-white rounded-lg shadow overflow-hidden">
      <div class="px-6 py-5 border-b border-gray-200">
        <h2 class="text-lg font-medium text-gray-900">Account Security</h2>
      </div>
      <div class="px-6 py-5">
        <div class="space-y-4">
          <div>
            <h3 class="text-base font-medium text-gray-900">Change Password</h3>
            <p class="mt-1 text-sm text-gray-500">
              Update your password to maintain account security.
            </p>
            <button class="mt-3 btn btn-outline">Change Password</button>
          </div>

          <div class="pt-4 border-t border-gray-200">
            <h3 class="text-base font-medium text-gray-900">
              Two-Factor Authentication
            </h3>
            <p class="mt-1 text-sm text-gray-500">
              Add an extra layer of security to your account.
            </p>
            <button class="mt-3 btn btn-outline">Enable 2FA</button>
          </div>
        </div>
      </div>
    </div>

    <div class="mt-8 bg-white rounded-lg shadow overflow-hidden">
      <div class="px-6 py-5 border-b border-gray-200">
        <h2 class="text-lg font-medium text-gray-900">Privacy Settings</h2>
      </div>
      <div class="px-6 py-5">
        <div class="space-y-4">
          <div class="flex items-start">
            <div class="flex items-center h-5">
              <input
                id="marketing"
                name="marketing"
                type="checkbox"
                class="h-4 w-4 text-primary-600 focus:ring-primary-500 border-gray-300 rounded"
              />
            </div>
            <div class="ml-3 text-sm">
              <label for="marketing" class="font-medium text-gray-700"
                >Marketing emails</label
              >
              <p class="text-gray-500">
                Receive emails about new products, features, and more.
              </p>
            </div>
          </div>

          <div class="flex items-start">
            <div class="flex items-center h-5">
              <input
                id="newsletter"
                name="newsletter"
                type="checkbox"
                class="h-4 w-4 text-primary-600 focus:ring-primary-500 border-gray-300 rounded"
              />
            </div>
            <div class="ml-3 text-sm">
              <label for="newsletter" class="font-medium text-gray-700"
                >Health newsletter</label
              >
              <p class="text-gray-500">
                Receive our weekly health newsletter with tips and articles.
              </p>
            </div>
          </div>
        </div>

        <div class="mt-5 flex justify-end">
          <button class="btn btn-primary">Save Preferences</button>
        </div>
      </div>
    </div>
  </div>
</div>
