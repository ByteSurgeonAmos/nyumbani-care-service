<script lang="ts">
  import { onMount } from "svelte";
  import toast from "svelte-french-toast";

  let email = "";
  let isLoading = false;
  let isSubmitted = false;

  async function handleSubmit() {
    isLoading = true;

    try {
      // This is a placeholder - in a real app you would:
      // await authService.resetPassword(email);
      isSubmitted = true;
      toast.success("Password reset instructions sent to your email");
    } catch (error) {
      console.error("Failed to request password reset:", error);
      toast.error("Failed to request password reset. Please try again.");
    } finally {
      isLoading = false;
    }
  }
</script>

<svelte:head>
  <title>Reset Password | NyumbaniCare</title>
</svelte:head>

<div class="min-h-full flex flex-col justify-center py-12 sm:px-6 lg:px-8">
  <div class="sm:mx-auto sm:w-full sm:max-w-md">
    <h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900">
      Reset your password
    </h2>
    <p class="mt-2 text-center text-sm text-gray-600">
      Enter your email address and we'll send you instructions to reset your
      password
    </p>
  </div>

  <div class="mt-8 sm:mx-auto sm:w-full sm:max-w-md">
    <div class="bg-white py-8 px-4 shadow sm:rounded-lg sm:px-10">
      {#if isSubmitted}
        <div class="rounded-md bg-green-50 p-4">
          <div class="flex">
            <div class="ml-3">
              <h3 class="text-sm font-medium text-green-800">
                Password reset email sent
              </h3>
              <p class="mt-2 text-sm text-green-700">
                If there's an account associated with {email}, you'll receive an
                email with password reset instructions.
              </p>
              <div class="mt-4">
                <a
                  href="/login"
                  class="text-sm font-medium text-primary-600 hover:text-primary-500"
                >
                  Return to login
                </a>
              </div>
            </div>
          </div>
        </div>
      {:else}
        <form class="space-y-6" on:submit|preventDefault={handleSubmit}>
          <div>
            <label for="email" class="block text-sm font-medium text-gray-700">
              Email address
            </label>
            <div class="mt-1">
              <input
                id="email"
                name="email"
                type="email"
                autocomplete="email"
                required
                bind:value={email}
                class="appearance-none block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-primary-500 focus:border-primary-500 sm:text-sm"
              />
            </div>
          </div>

          <div>
            <button
              type="submit"
              disabled={isLoading}
              class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 {isLoading
                ? 'opacity-70 cursor-not-allowed'
                : ''}"
            >
              {isLoading ? "Sending..." : "Send reset instructions"}
            </button>
          </div>

          <div class="text-center mt-4">
            <a
              href="/login"
              class="text-sm font-medium text-primary-600 hover:text-primary-500"
            >
              Back to login
            </a>
          </div>
        </form>
      {/if}
    </div>
  </div>
</div>
