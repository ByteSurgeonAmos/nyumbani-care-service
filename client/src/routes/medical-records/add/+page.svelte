<script lang="ts">
  import { onMount } from "svelte";
  import { goto, beforeNavigate } from "$app/navigation";
  import { isAuthenticated } from "$lib/api";
  import { page } from "$app/stores";
  import toast from "svelte-french-toast";

  let testResultId: string | null = null;
  let formData = {
    title: "",
    date: new Date().toISOString().substring(0, 10),
    recordType: "test_result",
    notes: "",
    symptoms: "",
    doctor: "",
    attachments: [],
  };

  let isSubmitting = false;
  let hasChanges = false;

  onMount(() => {
    if (!$isAuthenticated) {
      toast.error("Please login to add medical records");
      goto("/login?redirectTo=/medical-records/add");
    }

    // Get test result ID from query params if available
    const params = new URLSearchParams($page.url.search);
    testResultId = params.get("test_result_id");

    if (testResultId) {
      formData.title = "Test Result Analysis";
      formData.recordType = "test_result";
      formData.notes = `Associated test result ID: ${testResultId}`;
    }
  });

  // Confirm before leaving if there are unsaved changes
  beforeNavigate(({ cancel }) => {
    if (
      hasChanges &&
      !window.confirm(
        "You have unsaved changes. Are you sure you want to leave?"
      )
    ) {
      cancel();
    }
  });

  function updateField() {
    hasChanges = true;
  }

  async function handleSubmit() {
    if (!formData.title) {
      toast.error("Please enter a record title");
      return;
    }

    isSubmitting = true;

    try {
      // Add code to submit to API when endpoints are available
      // For now, simulate a successful submission
      await new Promise((resolve) => setTimeout(resolve, 1000));

      toast.success("Medical record saved successfully");
      hasChanges = false;
      goto("/medical-records");
    } catch (err) {
      console.error("Failed to save medical record:", err);
      toast.error("Failed to save medical record");
    } finally {
      isSubmitting = false;
    }
  }
</script>

<svelte:head>
  <title>Add Medical Record | NyumbaniCare</title>
  <meta
    name="description"
    content="Add a new medical record to your health history"
  />
</svelte:head>

<div class="bg-white">
  <div class="max-w-3xl mx-auto px-4 py-16 sm:px-6 lg:px-8">
    <div class="mb-8">
      <a
        href="/medical-records"
        class="text-primary-600 hover:text-primary-500"
      >
        &larr; Back to all records
      </a>
    </div>

    <h1 class="text-3xl font-bold tracking-tight text-gray-900">
      Add Medical Record
    </h1>
    <p class="mt-2 text-lg text-gray-500">
      {#if testResultId}
        Add test result to your medical records
      {:else}
        Track your health information for future reference
      {/if}
    </p>

    <form class="mt-10 space-y-8" on:submit|preventDefault={handleSubmit}>
      <div class="grid grid-cols-1 gap-y-6 gap-x-4 sm:grid-cols-6">
        <div class="sm:col-span-4">
          <label for="title" class="block text-sm font-medium text-gray-700">
            Record Title
          </label>
          <div class="mt-1">
            <input
              type="text"
              id="title"
              name="title"
              bind:value={formData.title}
              on:input={updateField}
              class="shadow-sm focus:ring-primary-500 focus:border-primary-500 block w-full sm:text-sm border-gray-300 rounded-md"
              required
            />
          </div>
        </div>

        <div class="sm:col-span-3">
          <label for="date" class="block text-sm font-medium text-gray-700">
            Date
          </label>
          <div class="mt-1">
            <input
              type="date"
              id="date"
              name="date"
              bind:value={formData.date}
              on:input={updateField}
              class="shadow-sm focus:ring-primary-500 focus:border-primary-500 block w-full sm:text-sm border-gray-300 rounded-md"
              required
            />
          </div>
        </div>

        <div class="sm:col-span-3">
          <label
            for="recordType"
            class="block text-sm font-medium text-gray-700"
          >
            Record Type
          </label>
          <div class="mt-1">
            <select
              id="recordType"
              name="recordType"
              bind:value={formData.recordType}
              on:change={updateField}
              class="shadow-sm focus:ring-primary-500 focus:border-primary-500 block w-full sm:text-sm border-gray-300 rounded-md"
              required
            >
              <option value="test_result">Test Result</option>
              <option value="doctor_visit">Doctor Visit</option>
              <option value="prescription">Prescription</option>
              <option value="vaccination">Vaccination</option>
              <option value="allergy">Allergy</option>
              <option value="other">Other</option>
            </select>
          </div>
        </div>

        <div class="sm:col-span-6">
          <label for="notes" class="block text-sm font-medium text-gray-700">
            Notes
          </label>
          <div class="mt-1">
            <textarea
              id="notes"
              name="notes"
              rows="4"
              bind:value={formData.notes}
              on:input={updateField}
              class="shadow-sm focus:ring-primary-500 focus:border-primary-500 block w-full sm:text-sm border-gray-300 rounded-md"
            ></textarea>
          </div>
          <p class="mt-2 text-sm text-gray-500">
            Add any additional notes or information about this record.
          </p>
        </div>

        {#if formData.recordType === "doctor_visit"}
          <div class="sm:col-span-4">
            <label for="doctor" class="block text-sm font-medium text-gray-700">
              Doctor/Healthcare Provider
            </label>
            <div class="mt-1">
              <input
                type="text"
                id="doctor"
                name="doctor"
                bind:value={formData.doctor}
                on:input={updateField}
                class="shadow-sm focus:ring-primary-500 focus:border-primary-500 block w-full sm:text-sm border-gray-300 rounded-md"
              />
            </div>
          </div>

          <div class="sm:col-span-6">
            <label
              for="symptoms"
              class="block text-sm font-medium text-gray-700"
            >
              Symptoms/Reason for Visit
            </label>
            <div class="mt-1">
              <textarea
                id="symptoms"
                name="symptoms"
                rows="3"
                bind:value={formData.symptoms}
                on:input={updateField}
                class="shadow-sm focus:ring-primary-500 focus:border-primary-500 block w-full sm:text-sm border-gray-300 rounded-md"
              ></textarea>
            </div>
          </div>
        {/if}

        {#if testResultId}
          <div class="sm:col-span-6">
            <div class="rounded-md bg-blue-50 p-4">
              <div class="flex">
                <div class="flex-shrink-0">
                  <svg
                    class="h-5 w-5 text-blue-400"
                    xmlns="http://www.w3.org/2000/svg"
                    viewBox="0 0 20 20"
                    fill="currentColor"
                  >
                    <path
                      fill-rule="evenodd"
                      d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z"
                      clip-rule="evenodd"
                    />
                  </svg>
                </div>
                <div class="ml-3 flex-1 md:flex md:justify-between">
                  <p class="text-sm text-blue-700">
                    This record will be linked to your test result #{testResultId.substring(
                      0,
                      8
                    )}
                  </p>
                  <p class="mt-3 text-sm md:mt-0 md:ml-6">
                    <a
                      href={`/test-kits/analyze/results/${testResultId}`}
                      class="whitespace-nowrap font-medium text-blue-700 hover:text-blue-600"
                    >
                      View result <span aria-hidden="true">&rarr;</span>
                    </a>
                  </p>
                </div>
              </div>
            </div>
          </div>
        {/if}
      </div>

      <div class="pt-5">
        <div class="flex justify-end">
          <button
            type="button"
            on:click={() => goto("/medical-records")}
            class="bg-white py-2 px-4 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
          >
            Cancel
          </button>
          <button
            type="submit"
            disabled={isSubmitting}
            class="ml-3 inline-flex justify-center py-2 px-4 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500 disabled:opacity-50"
          >
            {isSubmitting ? "Saving..." : "Save Record"}
          </button>
        </div>
      </div>
    </form>
  </div>
</div>
