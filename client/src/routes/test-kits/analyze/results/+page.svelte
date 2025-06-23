<script lang="ts">
  import { onMount } from "svelte";
  import { goto } from "$app/navigation";
  import { isAuthenticated, testKitResultService } from "$lib/api";
  import type { AnalyzeTestKitResponse } from "$lib/api/extendedServices";
  import toast from "svelte-french-toast";

  let results: AnalyzeTestKitResponse[] = [];
  let filteredResults: AnalyzeTestKitResponse[] = [];
  let isLoading = true;
  let error: string | null = null;

  // Filtering options
  let filterStatus = "all";
  let sortOrder = "newest";
  let searchQuery = "";

  onMount(async () => {
    if ($isAuthenticated) {
      loadResults();
    } else {
      toast.error("Please login to view your test results");
      goto("/login?redirectTo=/test-kits/analyze/results");
    }
  });

  async function loadResults() {
    isLoading = true;
    error = null;

    try {
      const response = await testKitResultService.getAll();
      results = response.data;
      applyFilters();
    } catch (err) {
      console.error("Failed to load test results:", err);
      error = "Failed to load test results. Please try again later.";
    } finally {
      isLoading = false;
    }
  }

  function applyFilters() {
    // Apply status filter
    let filtered = [...results];
    if (filterStatus !== "all") {
      filtered = filtered.filter(
        (r) => r.result.toLowerCase() === filterStatus.toLowerCase()
      );
    }

    // Apply search query if present
    if (searchQuery) {
      const query = searchQuery.toLowerCase();
      filtered = filtered.filter((r) => {
        return (
          r.test_kit?.name.toLowerCase().includes(query) ||
          r.interpretation?.toLowerCase().includes(query) ||
          r.result.toLowerCase().includes(query)
        );
      });
    }

    // Apply sort order
    filtered.sort((a, b) => {
      const dateA = new Date(a.created_at || 0).getTime();
      const dateB = new Date(b.created_at || 0).getTime();

      if (sortOrder === "newest") {
        return dateB - dateA;
      } else {
        return dateA - dateB;
      }
    });

    filteredResults = filtered;
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

  function viewResultDetails(id: string) {
    goto(`/test-kits/analyze/results/${id}`);
  }

  function handleFilterChange() {
    applyFilters();
  }
</script>

<svelte:head>
  <title>Test Kit Results | NyumbaniCare</title>
  <meta name="description" content="View your test kit analysis results" />
</svelte:head>

<div class="bg-white">
  <div class="max-w-7xl mx-auto px-4 py-16 sm:px-6 lg:px-8">
    <div class="sm:flex sm:items-center sm:justify-between">
      <h1 class="text-3xl font-bold tracking-tight text-gray-900">
        My Test Results
      </h1>
      <div class="mt-4 sm:mt-0">
        <a
          href="/test-kits/analyze"
          class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
        >
          Analyze New Test
        </a>
      </div>
    </div>

    <div class="mt-8">
      <!-- Filter and Sort Options -->
      <div class="mb-4">
        <div
          class="flex flex-col sm:flex-row sm:items-center sm:justify-between"
        >
          <!-- Status Filter -->
          <div class="flex items-center mb-2 sm:mb-0">
            <label
              for="statusFilter"
              class="mr-2 text-sm font-medium text-gray-700">Status:</label
            >
            <select
              id="statusFilter"
              bind:value={filterStatus}
              on:change={handleFilterChange}
              class="block w-full sm:w-auto border border-gray-300 rounded-md shadow-sm focus:ring-primary-500 focus:border-primary-500"
            >
              <option value="all">All</option>
              <option value="positive">Positive</option>
              <option value="negative">Negative</option>
              <option value="inconclusive">Inconclusive</option>
            </select>
          </div>

          <!-- Sort Order -->
          <div class="flex items-center mb-2 sm:mb-0">
            <label
              for="sortOrder"
              class="mr-2 text-sm font-medium text-gray-700">Sort by:</label
            >
            <select
              id="sortOrder"
              bind:value={sortOrder}
              on:change={handleFilterChange}
              class="block w-full sm:w-auto border border-gray-300 rounded-md shadow-sm focus:ring-primary-500 focus:border-primary-500"
            >
              <option value="newest">Newest First</option>
              <option value="oldest">Oldest First</option>
            </select>
          </div>

          <!-- Search -->
          <div class="flex items-center">
            <label for="searchQuery" class="sr-only">Search</label>
            <input
              id="searchQuery"
              type="text"
              bind:value={searchQuery}
              on:input={handleFilterChange}
              placeholder="Search results..."
              class="block w-full sm:w-auto border border-gray-300 rounded-md shadow-sm focus:ring-primary-500 focus:border-primary-500 px-3 py-2"
            />
          </div>
        </div>
      </div>

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
      {:else if filteredResults.length === 0}
        <div class="bg-white shadow overflow-hidden sm:rounded-lg">
          <div class="px-4 py-12 text-center sm:px-6">
            <p class="text-lg text-gray-500">
              No test results found matching your criteria.
            </p>
            <div class="mt-6">
              <a
                href="/test-kits/analyze"
                class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
              >
                Analyze a Test Kit
              </a>
            </div>
          </div>
        </div>
      {:else}
        <div class="overflow-hidden bg-white shadow sm:rounded-md">
          <ul role="list" class="divide-y divide-gray-200">
            {#each filteredResults as result}
              <li>
                <div class="block hover:bg-gray-50">
                  <div class="px-4 py-4 sm:px-6">
                    <div class="flex items-center justify-between">
                      <div class="flex items-center">
                        <div class="flex flex-col">
                          <p
                            class="text-sm font-medium text-primary-600 truncate"
                          >
                            Test Result #{result.id.substring(0, 8)}
                          </p>
                          <div class="flex mt-1">
                            <span
                              class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full {getResultBadgeClass(
                                result.result
                              )}"
                            >
                              {result.result.toUpperCase()}
                            </span>
                            <span class="ml-2 text-xs text-gray-500">
                              Confidence: {(result.confidence * 100).toFixed(
                                0
                              )}%
                            </span>
                          </div>
                        </div>
                      </div>
                      <div class="ml-2 flex-shrink-0 flex">
                        <button
                          on:click={() => viewResultDetails(result.id)}
                          class="font-medium text-primary-600 hover:text-primary-500"
                        >
                          View details
                        </button>
                      </div>
                    </div>
                    <div class="mt-2 sm:flex sm:justify-between">
                      <div class="sm:flex">
                        <p class="flex items-center text-sm text-gray-500">
                          <svg
                            class="flex-shrink-0 mr-1.5 h-5 w-5 text-gray-400"
                            fill="none"
                            stroke="currentColor"
                            viewBox="0 0 24 24"
                          >
                            <path
                              stroke-linecap="round"
                              stroke-linejoin="round"
                              stroke-width="2"
                              d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
                            />
                          </svg>
                          <span>
                            {result.interpretation
                              ? result.interpretation.substring(0, 60) + "..."
                              : "No interpretation available"}
                          </span>
                        </p>
                      </div>
                      <div
                        class="mt-2 flex items-center text-sm text-gray-500 sm:mt-0"
                      >
                        <svg
                          class="flex-shrink-0 mr-1.5 h-5 w-5 text-gray-400"
                          fill="none"
                          stroke="currentColor"
                          viewBox="0 0 24 24"
                        >
                          <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"
                          />
                        </svg>
                        <span
                          >{result.created_at
                            ? formatDate(result.created_at)
                            : "Unknown date"}</span
                        >
                      </div>
                    </div>
                  </div>
                </div>
              </li>
            {/each}
          </ul>
        </div>
      {/if}
    </div>
  </div>
</div>
