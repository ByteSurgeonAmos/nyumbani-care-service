<script lang="ts">
  import { onMount } from "svelte";
  import { goto } from "$app/navigation";
  import { page } from "$app/stores";
  import { isAuthenticated, testKitResultService } from "$lib/api";
  import type { AnalyzeTestKitResponse } from "$lib/api/extendedServices";
  import toast from "svelte-french-toast";

  let results: AnalyzeTestKitResponse[] = [];
  let isLoading = true;
  let error: string | null = null;

  onMount(async () => {
    if ($isAuthenticated) {
      loadResults();
    } else {
      toast.error("Please login to view your test results");
      goto("/login?redirectTo=" + $page.url.pathname + $page.url.search);
    }
  });

  async function loadResults() {
    isLoading = true;
    error = null;

    try {
      const params = new URLSearchParams($page.url.search);
      const result1Id = params.get("result1");
      const result2Id = params.get("result2");

      if (!result1Id || !result2Id) {
        error = "Missing result IDs for comparison";
        return;
      }

      // Load both results
      const [result1, result2] = await Promise.all([
        testKitResultService.getById(result1Id),
        testKitResultService.getById(result2Id),
      ]);

      results = [result1, result2];
    } catch (err) {
      console.error("Failed to load test results:", err);
      error = "Failed to load test results. Please try again later.";
    } finally {
      isLoading = false;
    }
  }

  function getResultBadgeClass(result: string): string {
    switch (result.toLowerCase()) {
      case "positive":
        return "bg-red-100 text-red-800";
      case "negative":
        return "bg-green-100 text-green-800";
      case "inconclusive":
      default:
        return "bg-yellow-100 text-yellow-800";
    }
  }

  function formatDate(dateString: string): string {
    const date = new Date(dateString);
    return new Intl.DateTimeFormat("en-US", {
      year: "numeric",
      month: "long",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    }).format(date);
  }

  function getDaysApart(date1: string, date2: string): string {
    const d1 = new Date(date1);
    const d2 = new Date(date2);
    const diffTime = Math.abs(d2.getTime() - d1.getTime());
    const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
    return diffDays === 1 ? "1 day" : `${diffDays} days`;
  }
</script>

<svelte:head>
  <title>Compare Test Results | NyumbaniCare</title>
  <meta
    name="description"
    content="Compare multiple test kit analysis results"
  />
</svelte:head>

<div class="bg-white">
  <div class="max-w-7xl mx-auto px-4 py-16 sm:px-6 lg:px-8">
    <div class="mb-8">
      <a
        href="/test-kits/analyze/results"
        class="text-primary-600 hover:text-primary-500"
      >
        &larr; Back to all results
      </a>
    </div>

    <h1 class="text-3xl font-bold tracking-tight text-gray-900">
      Compare Test Results
    </h1>

    <div class="mt-8">
      {#if isLoading}
        <div class="flex justify-center items-center py-12">
          <div
            class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"
          ></div>
        </div>
      {:else if error}
        <div class="rounded-md bg-red-50 p-4">
          <div class="flex">
            <div class="ml-3">
              <p class="text-sm font-medium text-red-800">{error}</p>
            </div>
          </div>
        </div>
      {:else if results.length !== 2}
        <div class="text-center py-12">
          <p class="text-lg text-gray-500">
            Unable to load both test results for comparison.
          </p>
          <div class="mt-6">
            <a
              href="/test-kits/analyze/results"
              class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
            >
              Return to Results
            </a>
          </div>
        </div>
      {:else}
        {#if results[0].test_kit && results[1].test_kit && results[0].test_kit.id === results[1].test_kit.id}
          <div class="mb-8 bg-blue-50 p-4 rounded-md">
            <p class="text-blue-700">
              Comparing two {results[0].test_kit.name} tests taken {results[0]
                .created_at && results[1].created_at
                ? getDaysApart(results[0].created_at, results[1].created_at)
                : "some time"} apart.
            </p>
          </div>
        {:else}
          <div class="mb-8 bg-yellow-50 p-4 rounded-md">
            <p class="text-yellow-700">
              Note: You are comparing results from different types of test kits,
              which may not provide a meaningful comparison.
            </p>
          </div>
        {/if}

        <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
          {#each results as result, index}
            <div class="bg-white shadow overflow-hidden sm:rounded-lg">
              <div class="px-4 py-5 sm:px-6">
                <div class="flex justify-between items-center">
                  <h3 class="text-lg leading-6 font-medium text-gray-900">
                    Test Result #{result.id.substring(0, 8)}
                  </h3>
                  <span
                    class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full {getResultBadgeClass(
                      result.result
                    )}"
                  >
                    {result.result.toUpperCase()}
                  </span>
                </div>
                <p class="mt-1 max-w-2xl text-sm text-gray-500">
                  {result.created_at
                    ? formatDate(result.created_at)
                    : "Unknown date"}
                </p>
              </div>

              <div class="border-t border-gray-200">
                {#if result.test_kit}
                  <div
                    class="bg-gray-50 px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6"
                  >
                    <dt class="text-sm font-medium text-gray-500">Test kit</dt>
                    <dd
                      class="mt-1 text-sm text-gray-900 sm:mt-0 sm:col-span-2"
                    >
                      {result.test_kit.name}
                    </dd>
                  </div>
                {/if}

                <div
                  class="bg-white px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6"
                >
                  <dt class="text-sm font-medium text-gray-500">Confidence</dt>
                  <dd class="mt-1 text-sm text-gray-900 sm:mt-0 sm:col-span-2">
                    <div class="flex items-center">
                      <div class="w-full bg-gray-200 rounded-full h-2.5">
                        <div
                          class="bg-primary-600 h-2.5 rounded-full"
                          style="width: {result.confidence * 100}%"
                        ></div>
                      </div>
                      <span class="ml-2"
                        >{(result.confidence * 100).toFixed(0)}%</span
                      >
                    </div>
                  </dd>
                </div>

                {#if result.interpretation}
                  <div
                    class="bg-gray-50 px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6"
                  >
                    <dt class="text-sm font-medium text-gray-500">
                      Interpretation
                    </dt>
                    <dd
                      class="mt-1 text-sm text-gray-900 sm:mt-0 sm:col-span-2"
                    >
                      {result.interpretation}
                    </dd>
                  </div>
                {/if}
              </div>

              <div class="border-t border-gray-200 px-4 py-5 sm:px-6">
                <a
                  href={`/test-kits/analyze/results/${result.id}`}
                  class="inline-flex items-center px-4 py-2 border border-gray-300 shadow-sm text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
                >
                  View Details
                </a>
              </div>
            </div>
          {/each}
        </div>

        <div class="mt-8 bg-white shadow overflow-hidden sm:rounded-lg">
          <div class="px-4 py-5 sm:px-6">
            <h3 class="text-lg leading-6 font-medium text-gray-900">
              Comparison Summary
            </h3>
          </div>

          <div class="border-t border-gray-200">
            <div
              class="bg-white px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6"
            >
              <dt class="text-sm font-medium text-gray-500">Result Status</dt>
              <dd class="mt-1 text-sm text-gray-900 sm:mt-0 sm:col-span-2">
                {#if results[0].result === results[1].result}
                  <div class="rounded-md bg-green-50 p-2">
                    <span class="text-green-700"
                      >Both tests show <strong>{results[0].result}</strong> results</span
                    >
                  </div>
                {:else}
                  <div class="rounded-md bg-yellow-50 p-2">
                    <span class="text-yellow-700"
                      >Results differ: <strong>{results[0].result}</strong> vs
                      <strong>{results[1].result}</strong></span
                    >
                  </div>
                {/if}
              </dd>
            </div>

            <div
              class="bg-gray-50 px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6"
            >
              <dt class="text-sm font-medium text-gray-500">
                Confidence Difference
              </dt>
              <dd class="mt-1 text-sm text-gray-900 sm:mt-0 sm:col-span-2">
                {Math.abs(
                  (results[0].confidence - results[1].confidence) * 100
                ).toFixed(1)}% difference in confidence levels
              </dd>
            </div>

            <div
              class="bg-white px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6"
            >
              <dt class="text-sm font-medium text-gray-500">Interpretation</dt>
              <dd class="mt-1 text-sm text-gray-900 sm:mt-0 sm:col-span-2">
                {#if results[0].interpretation === results[1].interpretation}
                  <div class="rounded-md bg-green-50 p-2">
                    <span class="text-green-700"
                      >The interpretations are consistent</span
                    >
                  </div>
                {:else}
                  <div class="rounded-md bg-yellow-50 p-2">
                    <span class="text-yellow-700"
                      >The interpretations differ</span
                    >
                  </div>
                {/if}
              </dd>
            </div>

            <div
              class="bg-gray-50 px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6"
            >
              <dt class="text-sm font-medium text-gray-500">Recommendation</dt>
              <dd class="mt-1 text-sm text-gray-900 sm:mt-0 sm:col-span-2">
                {#if results[0].result !== results[1].result}
                  Consider consulting with a healthcare professional to review
                  these different test results.
                {:else if results[0].result === "positive"}
                  Both tests show positive results. Please follow health
                  guidelines and consider consulting with a healthcare
                  professional.
                {:else if results[0].result === "negative"}
                  Both tests show negative results, which is encouraging.
                  Continue monitoring your health.
                {:else}
                  Both tests are inconclusive. Consider retaking the test or
                  consulting with a healthcare professional.
                {/if}
              </dd>
            </div>
          </div>
        </div>
      {/if}
    </div>
  </div>
</div>
