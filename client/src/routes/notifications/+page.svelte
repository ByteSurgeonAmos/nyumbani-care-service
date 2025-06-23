<script lang="ts">
  import { onMount } from "svelte";
  import { goto } from "$app/navigation";
  import {
    notificationService,
    notifications,
  } from "$lib/services/notificationService";
  import type { Notification } from "$lib/services/notificationService";
  import { isAuthenticated } from "$lib/api";
  import toast from "svelte-french-toast";
  import { showTestResultNotification } from "$lib/utils/notifications";

  let isLoading = true;
  let error: string | null = null;
  let filters = {
    all: true,
    unread: false,
    testResults: false,
    system: false,
  };

  onMount(async () => {
    if ($isAuthenticated) {
      await loadNotifications();
    } else {
      toast.error("Please log in to view your notifications");
      goto("/login?redirectTo=/notifications");
    }
  });

  async function loadNotifications() {
    isLoading = true;
    try {
      await notificationService.getNotifications(filters.unread);
      error = null;
    } catch (err) {
      console.error("Failed to load notifications:", err);
      error = "Failed to load notifications. Please try again later.";
    } finally {
      isLoading = false;
    }
  }

  async function markAsRead(notification: Notification) {
    try {
      await notificationService.markAsRead(notification.id);

      // Navigate if it's a linked resource
      if (
        notification.resourceType === "test_result" &&
        notification.resourceId
      ) {
        goto(`/test-kits/analyze/results/${notification.resourceId}`);
      }
    } catch (err) {
      console.error("Failed to mark notification as read:", err);
      toast.error("Failed to mark notification as read");
    }
  }

  async function markAllAsRead() {
    try {
      await notificationService.markAllAsRead();
      toast.success("All notifications marked as read");
    } catch (err) {
      console.error("Failed to mark all notifications as read:", err);
      toast.error("Failed to mark all notifications as read");
    }
  }

  function formatDate(dateString: string): string {
    const date = new Date(dateString);
    return date.toLocaleString("en-US", {
      weekday: "short",
      month: "short",
      day: "numeric",
      year: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  }

  function handleFilterChange() {
    loadNotifications();
  }

  $: filteredNotifications = $notifications.filter((notification) => {
    if (!filters.all) {
      if (filters.testResults && notification.type !== "test_result")
        return false;
      if (filters.system && notification.type !== "system") return false;
      if (filters.unread && notification.isRead) return false;
    }
    return true;
  });
</script>

<svelte:head>
  <title>Notifications | NyumbaniCare</title>
  <meta name="description" content="Your notifications" />
</svelte:head>

<div class="bg-white">
  <div class="max-w-7xl mx-auto px-4 py-16 sm:px-6 lg:px-8">
    <div class="flex items-center justify-between">
      <h1 class="text-3xl font-bold tracking-tight text-gray-900">
        Notifications
      </h1>
      {#if filteredNotifications.length > 0}
        <button
          on:click={markAllAsRead}
          class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
        >
          Mark all as read
        </button>
      {/if}
    </div>

    <div class="mt-8">
      <!-- Filter Options -->
      <div class="bg-white shadow rounded-md mb-6">
        <div class="px-4 py-5 border-b border-gray-200 sm:px-6">
          <h3 class="text-lg leading-6 font-medium text-gray-900">
            Filter Notifications
          </h3>
        </div>
        <div class="px-4 py-4 flex flex-wrap gap-4">
          <label class="inline-flex items-center">
            <input
              type="radio"
              bind:group={filters.all}
              value={true}
              on:change={() => {
                filters = {
                  all: true,
                  unread: false,
                  testResults: false,
                  system: false,
                };
                handleFilterChange();
              }}
              class="form-radio h-4 w-4 text-primary-600"
            />
            <span class="ml-2 text-sm text-gray-700">All</span>
          </label>
          <label class="inline-flex items-center">
            <input
              type="checkbox"
              bind:checked={filters.unread}
              on:change={() => {
                filters.all = false;
                handleFilterChange();
              }}
              class="form-checkbox h-4 w-4 text-primary-600"
            />
            <span class="ml-2 text-sm text-gray-700">Unread only</span>
          </label>
          <label class="inline-flex items-center">
            <input
              type="checkbox"
              bind:checked={filters.testResults}
              on:change={() => {
                filters.all = false;
                handleFilterChange();
              }}
              class="form-checkbox h-4 w-4 text-primary-600"
            />
            <span class="ml-2 text-sm text-gray-700">Test Results</span>
          </label>
          <label class="inline-flex items-center">
            <input
              type="checkbox"
              bind:checked={filters.system}
              on:change={() => {
                filters.all = false;
                handleFilterChange();
              }}
              class="form-checkbox h-4 w-4 text-primary-600"
            />
            <span class="ml-2 text-sm text-gray-700">System</span>
          </label>
        </div>
      </div>

      {#if isLoading}
        <div class="flex justify-center items-center py-12">
          <div
            class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"
          ></div>
        </div>
      {:else if error}
        <div class="rounded-md bg-red-50 p-4">
          <div class="flex">
            <div class="ml-3">
              <p class="text-sm font-medium text-red-800">{error}</p>
            </div>
          </div>
        </div>
      {:else if filteredNotifications.length === 0}
        <div
          class="text-center py-12 bg-white shadow overflow-hidden sm:rounded-lg"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            class="mx-auto h-12 w-12 text-gray-400"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9"
            />
          </svg>
          <h3 class="mt-2 text-lg font-medium text-gray-900">
            No notifications found
          </h3>
          <p class="mt-1 text-sm text-gray-500">
            You don't have any notifications matching your filters.
          </p>
        </div>
      {:else}
        <div class="overflow-hidden shadow rounded-md">
          <ul role="list" class="divide-y divide-gray-200">
            {#each filteredNotifications as notification}
              <li
                class="{notification.isRead
                  ? 'bg-white'
                  : 'bg-blue-50'} hover:bg-gray-50 cursor-pointer"
                on:click={() => markAsRead(notification)}
              >
                <div class="px-4 py-5 sm:px-6">
                  <div class="flex items-start">
                    <div class="flex-shrink-0 mt-0.5">
                      {#if notification.type === "test_result"}
                        <div class="rounded-full bg-primary-100 p-2">
                          <svg
                            xmlns="http://www.w3.org/2000/svg"
                            class="h-5 w-5 text-primary-600"
                            fill="none"
                            viewBox="0 0 24 24"
                            stroke="currentColor"
                          >
                            <path
                              stroke-linecap="round"
                              stroke-linejoin="round"
                              stroke-width="2"
                              d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
                            />
                          </svg>
                        </div>
                      {:else}
                        <div class="rounded-full bg-gray-100 p-2">
                          <svg
                            xmlns="http://www.w3.org/2000/svg"
                            class="h-5 w-5 text-gray-500"
                            fill="none"
                            viewBox="0 0 24 24"
                            stroke="currentColor"
                          >
                            <path
                              stroke-linecap="round"
                              stroke-linejoin="round"
                              stroke-width="2"
                              d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                            />
                          </svg>
                        </div>
                      {/if}
                    </div>
                    <div class="ml-4 flex-1">
                      <div class="flex justify-between">
                        <h3 class="text-sm font-medium text-gray-900">
                          {notification.title}
                        </h3>
                        <p class="text-xs text-gray-500">
                          {formatDate(notification.createdAt)}
                        </p>
                      </div>
                      <p class="mt-1 text-sm text-gray-600">
                        {notification.message}
                      </p>
                    </div>
                    {#if !notification.isRead}
                      <div class="ml-3 flex-shrink-0">
                        <span
                          class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-primary-100 text-primary-800"
                        >
                          New
                        </span>
                      </div>
                    {/if}
                  </div>
                </div>
              </li>
            {/each}
          </ul>
        </div>
      {/if}
    </div>
  </div>
</div>
