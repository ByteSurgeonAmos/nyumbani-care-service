<script lang="ts">
  import { authService } from "$lib/api";
  import { goto } from "$app/navigation";
  import toast from "svelte-french-toast";
  import { page } from "$app/stores";

  let email = "";
  let password = "";
  let isLoading = false;
  let error = "";

  // Get the redirect URL from the query string if it exists
  $: redirectTo = $page.url.searchParams.get("redirectTo") || "/";

  async function handleSubmit() {
    isLoading = true;
    error = "";

    try {
      await authService.login({ email, password });
      toast.success("Login successful!");
      goto(redirectTo);
    } catch (err: any) {
      console.error("Login failed:", err);
      error =
        err.response?.data?.error ||
        "Login failed. Please check your credentials.";
      toast.error(error);
    } finally {
      isLoading = false;
    }
  }
</script>

<svelte:head>
  <title>Login | NyumbaniCare</title>
</svelte:head>

<div class="min-h-full flex flex-col justify-center py-12 sm:px-6 lg:px-8">
  <div class="sm:mx-auto sm:w-full sm:max-w-md">
    <h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900">
      Sign in to your account
    </h2>
    <p class="mt-2 text-center text-sm text-gray-600">
      Or
      <a
        href="/register"
        class="font-medium text-primary-600 hover:text-primary-500"
      >
        create a new account
      </a>
    </p>
  </div>

  <div class="mt-8 sm:mx-auto sm:w-full sm:max-w-md">
    <div class="bg-white py-8 px-4 shadow sm:rounded-lg sm:px-10">
      <form class="space-y-6" on:submit|preventDefault={handleSubmit}>
        {#if error}
          <div class="rounded-md bg-red-50 p-4">
            <div class="flex">
              <div class="ml-3">
                <h3 class="text-sm font-medium text-red-800">{error}</h3>
              </div>
            </div>
          </div>
        {/if}

        <div>
          <label for="email" class="block text-sm font-medium text-gray-700"
            >Email address</label
          >
          <div class="mt-1">
            <input
              id="email"
              name="email"
              type="email"
              autocomplete="email"
              required
              bind:value={email}
              class="input"
            />
          </div>
        </div>

        <div>
          <label for="password" class="block text-sm font-medium text-gray-700"
            >Password</label
          >
          <div class="mt-1">
            <input
              id="password"
              name="password"
              type="password"
              autocomplete="current-password"
              required
              bind:value={password}
              class="input"
            />
          </div>
        </div>

        <div class="flex items-center justify-between">
          <div class="flex items-center">
            <input
              id="remember-me"
              name="remember-me"
              type="checkbox"
              class="h-4 w-4 text-primary-600 focus:ring-primary-500 border-gray-300 rounded"
            />
            <label for="remember-me" class="ml-2 block text-sm text-gray-900">
              Remember me
            </label>
          </div>
          <div class="text-sm">
            <a
              href="/reset-password"
              class="font-medium text-primary-600 hover:text-primary-500"
            >
              Forgot your password?
            </a>
          </div>
        </div>

        <div>
          <button
            type="submit"
            disabled={isLoading}
            class="w-full btn btn-primary py-2 px-4 {isLoading
              ? 'opacity-70 cursor-not-allowed'
              : ''}"
          >
            {isLoading ? "Signing in..." : "Sign in"}
          </button>
        </div>
      </form>
    </div>
  </div>
</div>
