import apiClient from "./apiClient";
import axios from "axios";
import { goto } from "$app/navigation";
import { browser } from "$app/environment";
import { writable } from "svelte/store";

// Types
export interface User {
  id: string;
  email: string;
  first_name: string;
  last_name: string;
  role: string;
  phone_number?: string;
  date_of_birth?: string;
  gender?: string;
  address?: string;
  created_at?: string;
  updated_at?: string;
}

interface LoginRequest {
  email: string;
  password: string;
}

interface RegisterRequest {
  email: string;
  password: string;
  first_name: string;
  last_name: string;
  phone_number: string;
  date_of_birth: string;
  gender: string;
  address: string;
}

interface AuthResponse {
  token: string;
  user: User;
}

export const currentUser = writable<User | null>(null);
export const isAuthenticated = writable<boolean>(false);

if (browser) {
  const storedUser = localStorage.getItem("user");
  const storedToken = localStorage.getItem("token");

  if (storedUser && storedToken) {
    currentUser.set(JSON.parse(storedUser));
    isAuthenticated.set(true);
  }
}

export const authService = {
  async login(credentials: LoginRequest): Promise<User> {
    try {
      const response = await apiClient.post<
        AuthResponse & { refresh_token?: string }
      >("/auth/login", credentials);
      const { token, user, refresh_token } = response.data;

      if (browser) {
        localStorage.setItem("token", token);
        localStorage.setItem("user", JSON.stringify(user));

        // Store refresh token if available
        if (refresh_token) {
          localStorage.setItem("refreshToken", refresh_token);
        }
      }

      currentUser.set(user);
      isAuthenticated.set(true);

      return user;
    } catch (error) {
      console.error("Login failed:", error);
      throw error;
    }
  },
  async register(userData: RegisterRequest): Promise<User> {
    try {
      const response = await apiClient.post<
        AuthResponse & { refresh_token?: string }
      >("/auth/register", userData);
      const { token, user, refresh_token } = response.data;

      if (browser) {
        localStorage.setItem("token", token);
        localStorage.setItem("user", JSON.stringify(user));

        // Store refresh token if available
        if (refresh_token) {
          localStorage.setItem("refreshToken", refresh_token);
        }
      }

      currentUser.set(user);
      isAuthenticated.set(true);

      return user;
    } catch (error) {
      console.error("Registration failed:", error);
      throw error;
    }
  },
  logout(): void {
    if (browser) {
      localStorage.removeItem("token");
      localStorage.removeItem("refreshToken");
      localStorage.removeItem("user");
    }

    currentUser.set(null);
    isAuthenticated.set(false);
    goto("/login");
  },

  async getCurrentUser(): Promise<User | null> {
    try {
      const response = await apiClient.get<User>("/users/me");
      const user = response.data;

      if (browser) {
        localStorage.setItem("user", JSON.stringify(user));
      }

      currentUser.set(user);
      isAuthenticated.set(true);

      return user;
    } catch (error) {
      console.error("Failed to get current user:", error);
      return null;
    }
  },
  async refreshToken(): Promise<string | null> {
    try {
      // Use stored refresh token if available
      const refreshToken = browser
        ? localStorage.getItem("refreshToken")
        : null;

      if (!refreshToken) {
        console.log("No refresh token available");
        return null;
      }

      const response = await apiClient.post<{
        token: string;
        refresh_token: string;
      }>("/auth/refresh", {
        refresh_token: refreshToken,
      });

      const { token, refresh_token } = response.data;

      if (browser) {
        localStorage.setItem("token", token);
        localStorage.setItem("refreshToken", refresh_token);
      }

      return token;
    } catch (error) {
      console.error("Failed to refresh token:", error);
      return null;
    }
  },

  async updateProfile(profileData: {
    first_name?: string;
    last_name?: string;
    phone_number?: string;
    address?: string;
  }): Promise<User> {
    try {
      const response = await apiClient.put<User>("/users/me", profileData);
      const updatedUser = response.data;

      if (browser) {
        localStorage.setItem("user", JSON.stringify(updatedUser));
      }

      currentUser.set(updatedUser);

      return updatedUser;
    } catch (error) {
      console.error("Profile update failed:", error);
      throw error;
    }
  },
};

export function requireAuth(
  url: URL
): { redirect: string; status: number } | void {
  if (browser && !localStorage.getItem("token")) {
    return {
      redirect: `/login?redirectTo=${url.pathname}`,
      status: 302,
    };
  }
}
