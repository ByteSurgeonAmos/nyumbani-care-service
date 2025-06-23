import axios, {
  type AxiosRequestConfig,
  type AxiosResponse,
  type AxiosError,
  type InternalAxiosRequestConfig,
} from "axios";
import { browser } from "$app/environment";

const API_URL = browser
  ? import.meta.env.VITE_API_URL || "http://localhost:8080"
  : "";

const apiClient = axios.create({
  baseURL: `${API_URL}/api/v1`,
  headers: {
    "Content-Type": "application/json",
    Accept: "application/json",
  },
});

apiClient.interceptors.request.use(
  (config: InternalAxiosRequestConfig): InternalAxiosRequestConfig => {
    if (browser) {
      const token = localStorage.getItem("token");
      if (token && config.headers) {
        config.headers.Authorization = `Bearer ${token}`;
      }
    }
    return config;
  },
  (error: AxiosError): Promise<AxiosError> => Promise.reject(error)
);

// Track if a token refresh is already in progress
let isRefreshing = false;
// Store all requests that should be retried after token refresh
let failedQueue: {
  resolve: (value: unknown) => void;
  reject: (reason?: any) => void;
}[] = [];

// Helper to process the queue of failed requests
const processQueue = (error: Error | null, token: string | null = null) => {
  failedQueue.forEach(({ resolve, reject }) => {
    if (error) {
      reject(error);
    } else {
      resolve(token);
    }
  });

  failedQueue = [];
};

apiClient.interceptors.response.use(
  (response: AxiosResponse): AxiosResponse => response,
  async (error: AxiosError) => {
    const originalRequest = error.config as any;

    // If error is 401 Unauthorized and we're not already retrying
    if (error.response?.status === 401 && !originalRequest._retry && browser) {
      // If path is login, just return the error
      if (window.location.pathname.includes("/login")) {
        return Promise.reject(error);
      }

      // If token refresh is in progress, queue this request
      if (isRefreshing) {
        return new Promise((resolve, reject) => {
          failedQueue.push({ resolve, reject });
        })
          .then((token) => {
            originalRequest.headers.Authorization = `Bearer ${token}`;
            return apiClient(originalRequest);
          })
          .catch((err) => {
            return Promise.reject(err);
          });
      }

      originalRequest._retry = true;
      isRefreshing = true;

      // Try to refresh token
      return new Promise((resolve, reject) => {
        // Import dynamically to avoid circular dependencies
        import("./authService").then(({ authService }) => {
          authService
            .refreshToken()
            .then((newToken) => {
              if (newToken) {
                // Update header with new token
                apiClient.defaults.headers.common.Authorization = `Bearer ${newToken}`;
                originalRequest.headers.Authorization = `Bearer ${newToken}`;

                // Process queue with new token
                processQueue(null, newToken);
                resolve(apiClient(originalRequest));
              } else {
                // Refresh failed, redirect to login
                console.log("Token refresh failed, redirecting to login");
                localStorage.removeItem("token");
                localStorage.removeItem("user");
                window.location.href = `/login?redirectTo=${window.location.pathname}`;

                // Process queue with error
                processQueue(new Error("Failed to refresh token"));
                reject(error);
              }
            })
            .catch((refreshError) => {
              console.error("Token refresh error:", refreshError);
              localStorage.removeItem("token");
              localStorage.removeItem("user");
              window.location.href = `/login?redirectTo=${window.location.pathname}`;

              // Process queue with error
              processQueue(refreshError);
              reject(error);
            })
            .finally(() => {
              isRefreshing = false;
            });
        });
      });
    }

    return Promise.reject(error);
  }
);

export default apiClient;
