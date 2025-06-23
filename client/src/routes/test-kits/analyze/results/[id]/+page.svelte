<script lang="ts">
  import { onMount } from "svelte";
  import { goto } from "$app/navigation";
  import { isAuthenticated, testKitResultService } from "$lib/api";
  import type { AnalyzeTestKitResponse } from "$lib/api/extendedServices";
  import { medicalRecordService } from "$lib/services/medicalRecordService";
  import { notificationService } from "$lib/services/notificationService";
  import { showTestResultNotification } from "$lib/utils/notifications";
  import toast from "svelte-french-toast";

  export let data: { id: string };
  const { id } = data;

  let result: AnalyzeTestKitResponse | null = null;
  let isLoading = true;
  let error: string | null = null;
  let otherResults: AnalyzeTestKitResponse[] = [];
  let showCompareModal = false;
  let selectedResultId: string | null = null;
  let showShareModal = false;
  let shareEmail = "";
  let isSharing = false;
  let isSavingToMedicalRecord = false;
  onMount(async () => {
    if ($isAuthenticated) {
      loadResult();

      // Check if there are any notifications related to this result and mark them as read
      try {
        const allNotifications = await notificationService.getNotifications();
        const resultNotification = allNotifications.find(
          (n) =>
            n.resourceType === "test_result" && n.resourceId === id && !n.isRead
        );

        if (resultNotification) {
          await notificationService.markAsRead(resultNotification.id);
        }
      } catch (err) {
        // Don't show errors for this background operation
        console.log("Could not process notifications:", err);
      }
    } else {
      toast.error("Please login to view test results");
      goto(`/login?redirectTo=/test-kits/analyze/results/${id}`);
    }
  });

  async function loadResult() {
    isLoading = true;
    error = null;

    try {
      // Load the current result
      result = await testKitResultService.getById(id);

      // Load other results for comparison (except current)
      if (result && result.test_kit_id) {
        const response = await testKitResultService.getAll();
        otherResults = response.data
          .filter((r) => r.id !== id && r.test_kit_id === result?.test_kit_id)
          .slice(0, 5); // Get up to 5 other results of the same test kit type
      }
    } catch (err) {
      console.error("Failed to load test result:", err);
      error = "Failed to load test result. Please try again later.";
      toast.error("Failed to load test result");
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

  function handleCompareClick(resultId: string) {
    selectedResultId = resultId;
    showCompareModal = true;
  }

  function openCompareModal() {
    if (otherResults.length > 0) {
      showCompareModal = true;
    } else {
      toast.error("No other results available for comparison");
    }
  }

  function closeCompareModal() {
    showCompareModal = false;
    selectedResultId = null;
  }

  function compareResults() {
    if (!selectedResultId) {
      toast.error("Please select a result to compare with");
      return;
    }

    goto(
      `/test-kits/analyze/results/compare?result1=${id}&result2=${selectedResultId}`
    );
  }

  function openShareModal() {
    showShareModal = true;
  }

  function closeShareModal() {
    showShareModal = false;
    shareEmail = "";
  }
  async function shareResult() {
    if (!shareEmail) {
      toast.error("Please enter an email address");
      return;
    }

    if (!validateEmail(shareEmail)) {
      toast.error("Please enter a valid email address");
      return;
    }

    isSharing = true;

    try {
      const message = `
        A test result has been shared with you from Nyumbani Care.
        
        Test Type: ${result?.test_kit?.name || "Home Test Kit"}
        Test Date: ${result?.created_at ? formatDate(result.created_at) : "Unknown date"}
        Result: ${result?.result?.toUpperCase()}
        
        Interpretation: ${result?.interpretation || "No interpretation available"}
        
        Recommendations: ${result?.advice || "No recommendations available"}
        
        You can view the complete results by logging into your Nyumbani Care account or creating one if you don't have one yet.
      `;

      await notificationService.sendTestResultNotification({
        to: shareEmail,
        subject: "Test Results Shared with You",
        message: message,
        resultId: id,
      });

      toast.success(`Result shared with ${shareEmail}`);
      closeShareModal();
    } catch (err) {
      console.error("Failed to share result:", err);
      toast.error("Failed to share result");
    } finally {
      isSharing = false;
    }
  }

  async function saveToMedicalRecords() {
    if (!result) return;

    isSavingToMedicalRecord = true;

    try {
      const response = await medicalRecordService.createFromTestResult(result);
      toast.success("Test result saved to medical records");

      // Optionally redirect to the medical record
      if (response.recordId) {
        goto(`/medical-records/${response.recordId}`);
      }
    } catch (err) {
      console.error("Failed to save to medical records:", err);
      toast.error("Failed to save to medical records");
    } finally {
      isSavingToMedicalRecord = false;
    }
  }

  function validateEmail(email: string) {
    const re = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return re.test(email);
  }
</script>

<svelte:head>
  <title>Test Result Details | NyumbaniCare</title>
  <meta
    name="description"
    content="View details of your test kit analysis results"
  />
</svelte:head>

<div class="bg-white">
  <div class="max-w-3xl mx-auto px-4 py-16 sm:px-6 lg:px-8">
    <div class="mb-8">
      <a
        href="/test-kits/analyze/results"
        class="text-primary-600 hover:text-primary-500"
      >
        &larr; Back to all results
      </a>
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
    {:else if !result}
      <div class="text-center py-12">
        <p class="text-lg text-gray-500">Test result not found.</p>
        <div class="mt-6">
          <a
            href="/test-kits/analyze"
            class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
          >
            Analyze New Test
          </a>
        </div>
      </div>
    {:else}
      <div class="bg-white shadow overflow-hidden sm:rounded-lg">
        <div class="px-4 py-5 sm:px-6">
          <div class="flex justify-between items-center">
            <h3 class="text-lg leading-6 font-medium text-gray-900">
              Test Result Details
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
            {result.created_at ? formatDate(result.created_at) : "Unknown date"}
          </p>
        </div>

        <div class="border-t border-gray-200">
          {#if result.test_kit}
            <div
              class="bg-gray-50 px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6"
            >
              <dt class="text-sm font-medium text-gray-500">Test kit</dt>
              <dd class="mt-1 text-sm text-gray-900 sm:mt-0 sm:col-span-2">
                {result.test_kit.name}
              </dd>
            </div>
          {/if}

          <div
            class="bg-white px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6"
          >
            <dt class="text-sm font-medium text-gray-500">Result</dt>
            <dd class="mt-1 text-sm text-gray-900 sm:mt-0 sm:col-span-2">
              {result.result}
            </dd>
          </div>

          <div
            class="bg-gray-50 px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6"
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
                <span class="ml-2">{(result.confidence * 100).toFixed(0)}%</span
                >
              </div>
            </dd>
          </div>

          {#if result.interpretation}
            <div
              class="bg-white px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6"
            >
              <dt class="text-sm font-medium text-gray-500">Interpretation</dt>
              <dd class="mt-1 text-sm text-gray-900 sm:mt-0 sm:col-span-2">
                {result.interpretation}
              </dd>
            </div>
          {/if}

          {#if result.advice}
            <div
              class="bg-gray-50 px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6"
            >
              <dt class="text-sm font-medium text-gray-500">
                Recommended steps
              </dt>
              <dd class="mt-1 text-sm text-gray-900 sm:mt-0 sm:col-span-2">
                <ul class="list-disc pl-5 space-y-1">
                  {#each result.advice
                    .split("\n")
                    .filter((item) => item.trim()) as item}
                    <li>{item}</li>
                  {/each}
                </ul>
              </dd>
            </div>
          {/if}

          {#if result.detected_markers && result.detected_markers.length > 0}
            <div
              class="bg-white px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6"
            >
              <dt class="text-sm font-medium text-gray-500">
                Detected markers
              </dt>
              <dd class="mt-1 text-sm text-gray-900 sm:mt-0 sm:col-span-2">
                <div class="flex flex-wrap gap-1">
                  {#each result.detected_markers as marker}
                    <span
                      class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800"
                    >
                      {marker}
                    </span>
                  {/each}
                </div>
              </dd>
            </div>
          {/if}
        </div>

        {#if result.image_url}
          <div class="border-t border-gray-200">
            <div class="px-4 py-5 sm:px-6">
              <h4 class="text-base font-medium text-gray-900">Test Image</h4>
              <div class="mt-2 max-w-lg mx-auto">
                <img
                  src={result.image_url}
                  alt="Test result"
                  class="h-auto w-full object-contain rounded-lg shadow"
                />
              </div>
            </div>
          </div>
        {/if}

        <div class="border-t border-gray-200 px-4 py-5 sm:px-6">
          <div class="flex justify-between">
            <a
              href="/test-kits/analyze"
              class="inline-flex items-center px-4 py-2 border border-gray-300 shadow-sm text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
            >
              Analyze Another Test
            </a>

            <div class="flex space-x-3">
              <button
                type="button"
                on:click={openShareModal}
                class="inline-flex items-center px-4 py-2 border border-gray-300 shadow-sm text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
              >
                Share with Healthcare Provider
              </button>

              {#if otherResults.length > 0}
                <button
                  type="button"
                  on:click={openCompareModal}
                  class="inline-flex items-center px-4 py-2 border border-gray-300 shadow-sm text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
                >
                  Compare Results
                </button>
              {/if}
              <button
                type="button"
                on:click={saveToMedicalRecords}
                disabled={isSavingToMedicalRecord}
                class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 disabled:opacity-50"
              >
                {#if isSavingToMedicalRecord}
                  <svg
                    class="animate-spin -ml-1 mr-2 h-4 w-4 text-white"
                    xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                  >
                    <circle
                      class="opacity-25"
                      cx="12"
                      cy="12"
                      r="10"
                      stroke="currentColor"
                      stroke-width="4"
                    ></circle>
                    <path
                      class="opacity-75"
                      fill="currentColor"
                      d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                    ></path>
                  </svg>
                  Saving...
                {:else}
                  Save to Medical Records
                {/if}
              </button>
            </div>
          </div>
        </div>
      </div>
    {/if}

    {#if showCompareModal}
      <div
        class="fixed z-10 inset-0 overflow-y-auto"
        aria-labelledby="modal-title"
        role="dialog"
        aria-modal="true"
      >
        <div
          class="flex items-end justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0"
        >
          <!-- Background overlay -->
          <div
            class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity"
            aria-hidden="true"
          ></div>

          <!-- Modal panel -->
          <div
            class="inline-block align-bottom bg-white rounded-lg px-4 pt-5 pb-4 text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full sm:p-6"
          >
            <div>
              <div class="mt-3 text-center sm:mt-5">
                <h3
                  class="text-lg leading-6 font-medium text-gray-900"
                  id="modal-title"
                >
                  Compare Test Results
                </h3>
                <div class="mt-4">
                  <p class="text-sm text-gray-500 mb-4">
                    Select a result to compare with the current one.
                  </p>

                  <div class="mt-2">
                    {#each otherResults as otherResult}
                      <div class="relative flex items-start py-2">
                        <div class="min-w-0 flex-1 text-sm">
                          <label
                            for={`result-${otherResult.id}`}
                            class="font-medium text-gray-700 flex items-center"
                          >
                            <input
                              id={`result-${otherResult.id}`}
                              name="result-option"
                              type="radio"
                              bind:group={selectedResultId}
                              value={otherResult.id}
                              class="focus:ring-primary-500 h-4 w-4 text-primary-600 border-gray-300"
                            />
                            <div class="ml-3 text-left">
                              <span class="block">
                                Test Result #{otherResult.id.substring(0, 8)}
                              </span>
                              <span class="block text-xs text-gray-500">
                                {otherResult.created_at
                                  ? formatDate(otherResult.created_at)
                                  : "Unknown date"}
                              </span>
                              <span
                                class="inline-flex mt-1 items-center px-2.5 py-0.5 rounded-full text-xs font-medium
                                {getResultBadgeClass(otherResult.result)}"
                              >
                                {otherResult.result.toUpperCase()}
                              </span>
                            </div>
                          </label>
                        </div>
                      </div>
                    {/each}
                  </div>
                </div>
              </div>
            </div>
            <div
              class="mt-5 sm:mt-6 sm:grid sm:grid-cols-2 sm:gap-3 sm:grid-flow-row-dense"
            >
              <button
                type="button"
                on:click={compareResults}
                disabled={!selectedResultId}
                class="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-primary-600 text-base font-medium text-white hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 sm:col-start-2 sm:text-sm disabled:opacity-50 disabled:cursor-not-allowed"
              >
                Compare
              </button>
              <button
                type="button"
                on:click={closeCompareModal}
                class="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 sm:mt-0 sm:col-start-1 sm:text-sm"
              >
                Cancel
              </button>
            </div>
          </div>
        </div>
      </div>
    {/if}

    {#if showShareModal}
      <div
        class="fixed z-10 inset-0 overflow-y-auto"
        aria-labelledby="modal-title"
        role="dialog"
        aria-modal="true"
      >
        <div
          class="flex items-end justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0"
        >
          <!-- Background overlay -->
          <div
            class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity"
            aria-hidden="true"
          ></div>

          <!-- Modal panel -->
          <div
            class="inline-block align-bottom bg-white rounded-lg px-4 pt-5 pb-4 text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full sm:p-6"
          >
            <div>
              <div class="mt-3 text-center sm:mt-5">
                <h3
                  class="text-lg leading-6 font-medium text-gray-900"
                  id="modal-title"
                >
                  Share with Healthcare Provider
                </h3>
                <div class="mt-4">
                  <p class="text-sm text-gray-500 mb-4">
                    Enter the email address of your healthcare provider to share
                    this test result.
                  </p>

                  <div class="mt-2">
                    <input
                      type="email"
                      bind:value={shareEmail}
                      placeholder="healthcare.provider@example.com"
                      class="shadow-sm focus:ring-primary-500 focus:border-primary-500 block w-full sm:text-sm border-gray-300 rounded-md"
                    />
                  </div>
                </div>
              </div>
            </div>
            <div
              class="mt-5 sm:mt-6 sm:grid sm:grid-cols-2 sm:gap-3 sm:grid-flow-row-dense"
            >
              <button
                type="button"
                on:click={shareResult}
                disabled={isSharing || !shareEmail}
                class="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-primary-600 text-base font-medium text-white hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 sm:col-start-2 sm:text-sm disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {isSharing ? "Sharing..." : "Share Result"}
              </button>
              <button
                type="button"
                on:click={closeShareModal}
                class="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 sm:mt-0 sm:col-start-1 sm:text-sm"
              >
                Cancel
              </button>
            </div>
          </div>
        </div>
      </div>
    {/if}
  </div>
</div>
