<script lang="ts">
  import { onMount } from "svelte";
  import { goto } from "$app/navigation";

  let isLoading = true;

  onMount(async () => {
    try {
      await goto("/education");
    } catch (error: any) {
      console.error("Failed to redirect:", error);
      isLoading = false;
    }
  });
</script>

<svelte:head>
  <title>Health Articles | NyumbaniCare</title>
  <meta
    name="description"
    content="Browse health education articles and resources to help you make informed health decisions."
  />
</svelte:head>

<div class="bg-white min-h-screen flex items-center justify-center">
  <div class="text-center p-8">
    <div
      class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-700 mx-auto mb-4"
    ></div>
    <p>Redirecting to education page...</p>
  </div>
</div>
