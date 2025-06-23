/**
 * Utility functions for displaying notifications in the app
 */

import { browser } from "$app/environment";

/**
 * Shows a test result ready notification
 * @param resultId The ID of the test result
 */
export function showTestResultNotification(resultId: string) {
  if (!browser) return;

  // Dynamically import the notification component to avoid issues with SSR
  import("$lib/components/TestResultReadyNotification.svelte").then(
    (module) => {
      const TestResultReadyNotification = module.default;

      // Create and mount the notification component
      const notification = new TestResultReadyNotification({
        target: document.body,
        props: { resultId },
      });

      // Component will clean itself up with onDestroy
    }
  );
}
