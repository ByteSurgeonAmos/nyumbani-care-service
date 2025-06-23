<script lang="ts">
  import { cart } from "$lib/stores/cartStore";
  import { isAuthenticated } from "$lib/api";
  import { type TestKit } from "$lib/api/apiServices";
  import toast from "svelte-french-toast";
  import { goto } from "$app/navigation";

  function addToCart(testKit: TestKit) {
    cart.addItem(testKit, 1);
    toast.success(`${testKit.name} added to cart`);
    goto("/cart");
  }
</script>

<div class="sticky bottom-0 w-full z-40">
  <div class="bg-white border-t border-gray-200 py-3 px-4 sm:px-6">
    <div class="flex justify-center md:justify-between items-center space-x-4">
      <div class="hidden md:block">
        <p class="text-sm font-medium text-gray-900">Complete your order</p>
        <p class="text-xs text-gray-500">
          Buy test kits and manage your health with ease
        </p>
      </div>

      <div class="flex space-x-2">
        <a
          href="/test-kits"
          class="rounded-md bg-white px-3 py-2 text-sm font-semibold text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 hover:bg-gray-50"
        >
          Browse Test Kits
        </a>

        <a
          href="/cart"
          class="rounded-md bg-primary-600 px-3 py-2 text-sm font-semibold text-white shadow-sm hover:bg-primary-700"
        >
          View Cart
        </a>
      </div>
    </div>
  </div>
</div>
