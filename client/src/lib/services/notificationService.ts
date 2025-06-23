import apiClient from "../api/apiClient";
import { writable } from "svelte/store";

export interface Notification {
  id: string;
  userId: string;
  type: string;
  title: string;
  message: string;
  resourceId?: string;
  resourceType?: string;
  isRead: boolean;
  isSent: boolean;
  sendTo?: string;
  sendMethod: string;
  sentAt?: string;
  createdAt: string;
  updatedAt: string;
}

interface SendNotificationOptions {
  to: string;
  subject: string;
  message: string;
  resultId?: string;
}

// Create stores for notifications
export const notifications = writable<Notification[]>([]);
export const unreadCount = writable<number>(0);

export const notificationService = {
  /**
   * Sends a notification email about test results
   */
  async sendTestResultNotification(
    options: SendNotificationOptions
  ): Promise<any> {
    try {
      const response = await apiClient.post("/notifications/email", {
        to: options.to,
        subject: options.subject,
        message: options.message,
        metadata: {
          type: "test_result",
          resultId: options.resultId,
        },
      });

      return {
        success: true,
        message: "Notification sent successfully",
      };
    } catch (error) {
      console.error("Error sending notification:", error);
      throw error;
    }
  },

  /**
   * Get notifications for the current user
   * @param unreadOnly If true, only fetch unread notifications
   */
  async getNotifications(unreadOnly: boolean = false): Promise<Notification[]> {
    try {
      const params = unreadOnly ? "?unread=true" : "";
      const response = await apiClient.get(`/notifications${params}`);

      // Update the notifications store
      notifications.set(response.data);

      return response.data;
    } catch (error) {
      console.error("Error fetching notifications:", error);
      throw error;
    }
  },

  /**
   * Get count of unread notifications
   */
  async getUnreadCount(): Promise<number> {
    try {
      const response = await apiClient.get("/notifications/unread/count");

      // Update the unread count store
      unreadCount.set(response.data.count);

      return response.data.count;
    } catch (error) {
      console.error("Error fetching unread count:", error);
      throw error;
    }
  },

  /**
   * Marks a notification as read
   */
  async markAsRead(notificationId: string): Promise<any> {
    try {
      const response = await apiClient.put(
        `/notifications/${notificationId}/read`
      );

      // Update the store after marking as read
      notifications.update((items) => {
        return items.map((item) => {
          if (item.id === notificationId) {
            return { ...item, isRead: true };
          }
          return item;
        });
      });

      // Decrease the unread count
      unreadCount.update((count) => Math.max(0, count - 1));

      return response.data;
    } catch (error) {
      console.error("Error marking notification as read:", error);
      throw error;
    }
  },
  /**
   * Marks all notifications as read
   */
  async markAllAsRead(): Promise<any> {
    try {
      // Call the backend to mark all as read
      const response = await apiClient.put("/notifications/read-all");

      // Update local state
      notifications.update((items) => {
        return items.map((item) => ({ ...item, isRead: true }));
      });

      unreadCount.set(0);

      return response.data;
    } catch (error) {
      console.error("Error marking all notifications as read:", error);
      throw error;
    }
  },

  /**
   * Loads initial notifications and sets up refresh interval
   */
  initialize() {
    // Initial load
    this.getNotifications();
    this.getUnreadCount();

    // Set up interval to check for new notifications (every 2 minutes)
    const interval = setInterval(() => {
      this.getUnreadCount();
    }, 2 * 60 * 1000);

    return () => clearInterval(interval); // Return cleanup function
  },
};
