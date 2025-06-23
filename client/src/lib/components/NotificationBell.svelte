<script lang="ts">
  import { onMount, onDestroy } from "svelte";
  import {
    notificationService,
    notifications,
    unreadCount,
  } from "$lib/services/notificationService";
  import type { Notification } from "$lib/services/notificationService";
  import { goto } from "$app/navigation";
  import { fade, slide } from "svelte/transition";

  let showDropdown = false;
  let cleanupFn: () => void;

  onMount(() => {
    // Initialize notifications and set up polling
    cleanupFn = notificationService.initialize();
  });

  onDestroy(() => {
    if (cleanupFn) cleanupFn();
  });

  function toggleDropdown() {
    showDropdown = !showDropdown;
  }

  function closeDropdown() {
    showDropdown = false;
  }

  async function markAsRead(notification: Notification) {
    await notificationService.markAsRead(notification.id);

    // If it's a resource-linked notification, navigate to it
    if (
      notification.resourceType === "test_result" &&
      notification.resourceId
    ) {
      goto(`/test-kits/analyze/results/${notification.resourceId}`);
    }

    closeDropdown();
  }

  async function markAllAsRead() {
    await notificationService.markAllAsRead();
  }

  function getTimeAgo(dateString: string): string {
    const date = new Date(dateString);
    const now = new Date();
    const diffMs = now.getTime() - date.getTime();

    const diffSecs = Math.floor(diffMs / 1000);
    const diffMins = Math.floor(diffSecs / 60);
    const diffHours = Math.floor(diffMins / 60);
    const diffDays = Math.floor(diffHours / 24);

    if (diffDays > 0) {
      return `${diffDays} day${diffDays !== 1 ? "s" : ""} ago`;
    } else if (diffHours > 0) {
      return `${diffHours} hour${diffHours !== 1 ? "s" : ""} ago`;
    } else if (diffMins > 0) {
      return `${diffMins} minute${diffMins !== 1 ? "s" : ""} ago`;
    } else {
      return "Just now";
    }
  }
</script>

<div class="relative">
  <button
    type="button"
    class="relative p-1 text-gray-700 hover:text-primary-600 focus:outline-none"
    on:click={toggleDropdown}
  >
    <svg
      xmlns="http://www.w3.org/2000/svg"
      class="h-6 w-6"
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
    {#if $unreadCount > 0}
      <span
        class="absolute top-0 right-0 block h-5 w-5 rounded-full bg-red-500 text-xs text-white font-bold flex items-center justify-center"
      >
        {$unreadCount > 9 ? "9+" : $unreadCount}
      </span>
    {/if}
  </button>

  {#if showDropdown}
    <div
      transition:fade={{ duration: 100 }}
      class="origin-top-right absolute right-0 mt-2 w-80 rounded-md shadow-lg bg-white ring-1 ring-black ring-opacity-5 focus:outline-none z-50"
      role="menu"
      aria-orientation="vertical"
      aria-labelledby="notification-button"
    >
      <div class="py-1 divide-y divide-gray-100" role="none">
        <div class="flex items-center justify-between px-4 py-2">
          <h3 class="text-sm font-medium text-gray-900">Notifications</h3>
          {#if $unreadCount > 0}
            <button
              type="button"
              class="text-xs text-primary-600 font-medium hover:underline"
              on:click={markAllAsRead}
            >
              Mark all as read
            </button>
          {/if}
        </div>

        {#if $notifications.length === 0}
          <div class="px-4 py-6 text-sm text-gray-500 text-center">
            No notifications yet
          </div>
        {:else}
          <div class="max-h-80 overflow-y-auto">
            {#each $notifications as notification}
              <div
                transition:slide|local={{ duration: 150 }}
                class="px-4 py-3 hover:bg-gray-50 cursor-pointer {notification.isRead
                  ? 'bg-white'
                  : 'bg-blue-50'}"
                on:click={() => markAsRead(notification)}
                role="menuitem"
              >
                <div class="flex items-start">
                  <div class="flex-shrink-0 mt-0.5">
                    {#if notification.type === "test_result"}
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
                    {:else}
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
                    {/if}
                  </div>
                  <div class="ml-3 w-0 flex-1">
                    <p class="text-sm font-medium text-gray-900 truncate">
                      {notification.title}
                    </p>
                    <p class="text-sm text-gray-500 truncate">
                      {notification.message}
                    </p>
                    <p class="text-xs text-gray-400 mt-1">
                      {getTimeAgo(notification.createdAt)}
                    </p>
                  </div>
                  {#if !notification.isRead}
                    <div class="ml-2 flex-shrink-0">
                      <span
                        class="inline-block h-2 w-2 rounded-full bg-primary-600"
                      ></span>
                    </div>
                  {/if}
                </div>
              </div>
            {/each}
          </div>
          <div class="px-4 py-2 text-center">
            <a
              href="/notifications"
              class="text-sm text-primary-600 font-medium hover:underline"
              on:click={() => {
                closeDropdown();
              }}
            >
              View all notifications
            </a>
          </div>
        {/if}
      </div>
    </div>
  {/if}
</div>
