import { writable } from "svelte/store";
import { browser } from "$app/environment";
import type { TestKit } from "$lib/api";

export interface CartItem {
  testKit: TestKit;
  quantity: number;
}

// Initialize the cart store from localStorage if available
const storedCart = browser ? localStorage.getItem("cart") : null;
const initialCart: CartItem[] = storedCart ? JSON.parse(storedCart) : [];

// Create the cart store
const createCartStore = () => {
  const { subscribe, set, update } = writable<CartItem[]>(initialCart);

  // Save to localStorage whenever the store changes
  if (browser) {
    subscribe((value) => {
      localStorage.setItem("cart", JSON.stringify(value));
    });
  }

  return {
    subscribe,

    // Add item to cart
    addItem: (testKit: TestKit, quantity: number = 1) => {
      update((items) => {
        // Check if the item already exists in the cart
        const existingItemIndex = items.findIndex(
          (item) => item.testKit.id === testKit.id
        );

        if (existingItemIndex !== -1) {
          // Update quantity if item exists
          const updatedItems = [...items];
          updatedItems[existingItemIndex].quantity += quantity;
          return updatedItems;
        } else {
          // Add new item if it doesn't exist
          return [...items, { testKit, quantity }];
        }
      });
    },

    // Update item quantity
    updateQuantity: (testKitId: string, quantity: number) => {
      update((items) =>
        items.map((item) =>
          item.testKit.id === testKitId
            ? { ...item, quantity: Math.max(1, quantity) }
            : item
        )
      );
    },

    // Remove item from cart
    removeItem: (testKitId: string) => {
      update((items) => items.filter((item) => item.testKit.id !== testKitId));
    },

    // Clear the entire cart
    clearCart: () => {
      set([]);
    },

    // Calculate total cost
    getTotal: (items: CartItem[]): number => {
      return items.reduce(
        (total, item) => total + item.testKit.price * item.quantity,
        0
      );
    },
  };
};

export const cart = createCartStore();

// Derived stores for cart information
export const cartItemCount = writable(0);
export const cartTotal = writable(0);

// Update derived stores when cart changes
if (browser) {
  cart.subscribe((items) => {
    cartItemCount.set(items.reduce((total, item) => total + item.quantity, 0));
    cartTotal.set(cart.getTotal(items));
  });
}
