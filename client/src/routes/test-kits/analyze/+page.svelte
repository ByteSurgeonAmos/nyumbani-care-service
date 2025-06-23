<script lang="ts">
  import { onMount } from "svelte";
  import { testKitService, testKitResultService } from "$lib/api";
  import { isAuthenticated } from "$lib/api";
  import type { TestKit } from "$lib/api";
  import toast from "svelte-french-toast";

  let testKits: TestKit[] = [];
  let selectedTestKit: string | null = null;
  let imageFile: File | null = null;
  let previewUrl: string | null = null;
  let isAnalyzing = false;
  let analysisResult: any = null;
  let isLoading = true;
  let error: string | null = null;

  onMount(async () => {
    if ($isAuthenticated) {
      await loadTestKits();
    }
  });

  async function loadTestKits() {
    isLoading = true;
    error = null;

    try {
      const response = await testKitService.getAll(1, 100);
      testKits = response.data;
    } catch (err) {
      console.error("Failed to load test kits:", err);
      error = "Failed to load test kits. Please try again later.";
    } finally {
      isLoading = false;
    }
  }

  function handleFileChange(event: Event) {
    const input = event.target as HTMLInputElement;
    if (!input.files || input.files.length === 0) {
      imageFile = null;
      previewUrl = null;
      return;
    }

    imageFile = input.files[0];

    // Create a preview URL for the selected image
    if (previewUrl) {
      URL.revokeObjectURL(previewUrl);
    }
    previewUrl = URL.createObjectURL(imageFile);
  }

  async function analyzeImage() {
    if (!selectedTestKit || !imageFile) {
      toast.error("Please select a test kit and upload an image");
      return;
    }

    isAnalyzing = true;
    analysisResult = null;
    error = null;

    try {
      // Convert image to base64
      const base64Image = await fileToBase64(imageFile);

      // Send to API
      const result = await testKitResultService.analyzeImage({
        test_kit_id: selectedTestKit,
        image_data: base64Image,
      });
      analysisResult = result;
      toast.success(
        "Analysis completed successfully. A notification has been sent."
      );
    } catch (err) {
      console.error("Failed to analyze image:", err);
      error = "Failed to analyze image. Please try again.";
      toast.error("Failed to analyze image");
    } finally {
      isAnalyzing = false;
    }
  }

  function fileToBase64(file: File): Promise<string> {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();
      reader.readAsDataURL(file);
      reader.onload = () => {
        let encoded = reader.result?.toString() || "";
        // Remove data:image/jpeg;base64, prefix
        if (encoded.includes(",")) {
          encoded = encoded.split(",")[1];
        }
        resolve(encoded);
      };
      reader.onerror = (error) => reject(error);
    });
  }
</script>

<svelte:head>
  <title>Analyze Test Kit | NyumbaniCare</title>
  <meta name="description" content="Analyze your test kit results" />
</svelte:head>

<div class="bg-white">
  <div class="max-w-3xl mx-auto px-4 py-16 sm:px-6 sm:py-24">
    <div class="text-center">
      <h1 class="text-3xl font-bold tracking-tight text-gray-900 sm:text-4xl">
        Test Kit Result Analysis
      </h1>
      <p class="mt-4 text-lg text-gray-500">
        Upload a photo of your test kit result for AI analysis
      </p>
      {#if $isAuthenticated}
        <div class="mt-4">
          <a
            href="/test-kits/analyze/results"
            class="inline-flex items-center text-primary-600 hover:text-primary-700"
          >
            <span>View my test results</span>
            <svg class="ml-1 h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
              <path
                fill-rule="evenodd"
                d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z"
                clip-rule="evenodd"
              />
            </svg>
          </a>
        </div>
      {/if}
    </div>

    {#if !$isAuthenticated}
      <div class="mt-10 py-8 bg-gray-50 rounded-lg text-center">
        <p class="text-lg text-gray-700 mb-4">
          You need to be logged in to use this feature.
        </p>
        <a
          href="/login?redirectTo=/test-kits/analyze"
          class="inline-flex items-center px-4 py-2 border border-transparent text-base font-medium rounded-md shadow-sm text-white bg-primary-600 hover:bg-primary-700"
        >
          Log In
        </a>
      </div>
    {:else if isLoading}
      <div class="mt-10 flex justify-center">
        <div
          class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-700"
        ></div>
      </div>
    {:else}
      <div class="mt-10">
        {#if error}
          <div class="mb-6 rounded-md bg-red-50 p-4">
            <div class="flex">
              <div class="ml-3">
                <p class="text-sm font-medium text-red-800">{error}</p>
              </div>
            </div>
          </div>
        {/if}

        <form class="space-y-6" on:submit|preventDefault={analyzeImage}>
          <div>
            <label
              for="testKit"
              class="block text-sm font-medium text-gray-700"
            >
              Select Test Kit Type
            </label>
            <select
              id="testKit"
              bind:value={selectedTestKit}
              class="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-primary-500 focus:border-primary-500 sm:text-sm rounded-md"
              required
            >
              <option value="">Select a test kit</option>
              {#each testKits as kit}
                <option value={kit.id}>{kit.name}</option>
              {/each}
            </select>
          </div>

          <div>
            <label for="image" class="block text-sm font-medium text-gray-700">
              Upload Test Kit Image
            </label>
            <div
              class="mt-1 flex justify-center px-6 pt-5 pb-6 border-2 border-gray-300 border-dashed rounded-md"
            >
              <div class="space-y-1 text-center">
                {#if previewUrl}
                  <img
                    src={previewUrl}
                    alt="Test kit preview"
                    class="mx-auto h-32 w-auto object-contain mb-4"
                  />
                {:else}
                  <svg
                    class="mx-auto h-12 w-12 text-gray-400"
                    stroke="currentColor"
                    fill="none"
                    viewBox="0 0 48 48"
                    aria-hidden="true"
                  >
                    <path
                      d="M28 8H12a4 4 0 00-4 4v20m32-12v8m0 0v8a4 4 0 01-4 4H8m36-12h-4m4 0H20"
                      stroke-width="2"
                      stroke-linecap="round"
                      stroke-linejoin="round"
                    />
                  </svg>
                {/if}
                <div class="flex text-sm text-gray-600">
                  <label
                    for="file-upload"
                    class="relative cursor-pointer bg-white rounded-md font-medium text-primary-600 hover:text-primary-500 focus-within:outline-none focus-within:ring-2 focus-within:ring-offset-2 focus-within:ring-primary-500"
                  >
                    <span>Upload a file</span>
                    <input
                      id="file-upload"
                      name="file-upload"
                      type="file"
                      class="sr-only"
                      accept="image/*"
                      on:change={handleFileChange}
                    />
                  </label>
                  <p class="pl-1">or drag and drop</p>
                </div>
                <p class="text-xs text-gray-500">PNG, JPG, GIF up to 10MB</p>
              </div>
            </div>
          </div>

          <div>
            <button
              type="submit"
              class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 disabled:opacity-50 disabled:cursor-not-allowed"
              disabled={isAnalyzing || !selectedTestKit || !imageFile}
            >
              {isAnalyzing ? "Analyzing..." : "Analyze Result"}
            </button>
          </div>
        </form>

        {#if analysisResult}
          <div class="mt-10 bg-gray-50 shadow overflow-hidden sm:rounded-lg">
            <div class="px-4 py-5 sm:px-6">
              <h3 class="text-lg leading-6 font-medium text-gray-900">
                Analysis Results
              </h3>
              <p class="mt-1 max-w-2xl text-sm text-gray-500">
                AI analysis of your test kit image
              </p>
            </div>
            <div class="border-t border-gray-200">
              <dl>
                <div
                  class="bg-white px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6"
                >
                  <dt class="text-sm font-medium text-gray-500">Result</dt>
                  <dd class="mt-1 text-sm text-gray-900 sm:mt-0 sm:col-span-2">
                    {analysisResult.result}
                  </dd>
                </div>
                <div
                  class="bg-gray-50 px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6"
                >
                  <dt class="text-sm font-medium text-gray-500">Confidence</dt>
                  <dd class="mt-1 text-sm text-gray-900 sm:mt-0 sm:col-span-2">
                    {(analysisResult.confidence * 100).toFixed(2)}%
                  </dd>
                </div>
                <div
                  class="bg-white px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6"
                >
                  <dt class="text-sm font-medium text-gray-500">
                    Interpretation
                  </dt>
                  <dd class="mt-1 text-sm text-gray-900 sm:mt-0 sm:col-span-2">
                    {analysisResult.interpretation}
                  </dd>
                </div>
                <div
                  class="bg-gray-50 px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6"
                >
                  <dt class="text-sm font-medium text-gray-500">
                    Recommendations
                  </dt>
                  <dd class="mt-1 text-sm text-gray-900 sm:mt-0 sm:col-span-2">
                    {analysisResult.advice}
                  </dd>
                </div>
              </dl>
            </div>
            <div class="border-t border-gray-200 px-4 py-5 sm:px-6">
              <div class="flex justify-between">
                <button
                  type="button"
                  on:click={() => {
                    selectedTestKit = null;
                    imageFile = null;
                    previewUrl = null;
                    analysisResult = null;
                  }}
                  class="inline-flex items-center px-4 py-2 border border-gray-300 shadow-sm text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
                >
                  Analyze Another Test
                </button>

                <a
                  href="/test-kits/analyze/results"
                  class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
                >
                  View All Results
                </a>
              </div>
            </div>
          </div>
        {/if}
      </div>
    {/if}
  </div>
</div>
