<script lang="ts">
  import { onMount, onDestroy } from "svelte";
  import toast from "svelte-french-toast";
  import {
    notificationService,
    unreadCount,
  } from "$lib/services/notificationService";
  import { goto } from "$app/navigation";

  export let resultId: string;

  let toastId: string | null = null;
  let styleElement: HTMLStyleElement | null = null;

  onMount(() => {
    // Mark any existing notifications for this result as read
    notificationService.getNotifications().then((notifications) => {
      const resultNotification = notifications.find(
        (n) =>
          n.resourceType === "test_result" &&
          n.resourceId === resultId &&
          !n.isRead
      );

      if (resultNotification) {
        notificationService.markAsRead(resultNotification.id);
      }
    }); // Show toast notification
    toastId = toast(
      `Your test result is ready! <button 
        class="text-primary-700 font-bold underline ml-2" 
        onclick="window.location.href='/test-kits/analyze/results/${resultId}'">
        View Result
      </button>`,
      {
        duration: 10000,
        position: "top-center",
        className: "test-result-notification",
      }
    );

    // Add some custom CSS to style the toast notification
    styleElement = document.createElement("style");
    styleElement.innerHTML = `
      .test-result-notification {
        border-left: 4px solid #4f46e5 !important;
      }
      .test-result-notification button {
        background: none;
        border: none;
        padding: 0;
        cursor: pointer;
      }
    `;
    document.head.appendChild(styleElement);
  });
  onDestroy(() => {
    if (toastId) {
      toast.dismiss(toastId);
    }

    // Clean up the style element
    if (styleElement && document.head.contains(styleElement)) {
      document.head.removeChild(styleElement);
    }
  });
</script>

<!-- This component has no UI, it just shows a toast notification -->
```
