<script lang="ts">
  import { onMount } from "svelte";
  import { goto } from "$app/navigation";
  import { isAuthenticated } from "$lib/api";
  import toast from "svelte-french-toast";

  export let data: { id: string };
  const { id } = data;

  let record: any = null;
  let isLoading = true;
  let error: string | null = null;

  // Mock data for the prototype
  const mockRecords = {
    "1": {
      id: "1",
      title: "Annual Checkup",
      date: "2025-05-10",
      recordType: "doctor_visit",
      doctor: "Dr. Smith",
      notes: "Regular annual checkup. All vitals normal.",
      symptoms: "None - routine visit",
      created_at: "2025-05-10T15:30:00Z",
    },
    "2": {
      id: "2",
      title: "COVID-19 Test Result",
      date: "2025-05-01",
      recordType: "test_result",
      notes: "Associated test result ID: abc12345",
      result: "negative",
      created_at: "2025-05-01T09:15:00Z",
    },
    "3": {
      id: "3",
      title: "Flu Vaccination",
      date: "2025-04-15",
      recordType: "vaccination",
      notes: "Annual flu vaccine administered.",
      doctor: "Dr. Johnson",
      created_at: "2025-04-15T11:45:00Z",
    },
  };

  onMount(async () => {
    if ($isAuthenticated) {
      loadRecord();
    } else {
      toast.error("Please login to view medical records");
      goto(`/login?redirectTo=/medical-records/${id}`);
    }
  });

  async function loadRecord() {
    isLoading = true;
    error = null;

    try {
      // Simulate API call with mock data
      await new Promise((resolve) => setTimeout(resolve, 800));
      record = mockRecords[id as keyof typeof mockRecords];

      if (!record) {
        error = "Medical record not found";
      }
    } catch (err) {
      console.error("Failed to load medical record:", err);
      error = "Failed to load medical record. Please try again later.";
    } finally {
      isLoading = false;
    }
  }

  function formatDate(dateString: string): string {
    const date = new Date(dateString);
    return new Intl.DateTimeFormat("en-US", {
      year: "numeric",
      month: "long",
      day: "numeric",
    }).format(date);
  }

  function formatDateTime(dateString: string): string {
    const date = new Date(dateString);
    return new Intl.DateTimeFormat("en-US", {
      year: "numeric",
      month: "long",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    }).format(date);
  }

  function getRecordTypeLabel(recordType: string): string {
    return recordType
      .replace(/_/g, " ")
      .replace(/\b\w/g, (l) => l.toUpperCase());
  }
</script>

<svelte:head>
  <title>Medical Record | NyumbaniCare</title>
  <meta name="description" content="View medical record details" />
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
    {:else if !record}
      <div class="text-center py-12">
        <p class="text-lg text-gray-500">Medical record not found.</p>
        <div class="mt-6">
          <a
            href="/medical-records"
            class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
          >
            View All Records
          </a>
        </div>
      </div>
    {:else}
      <div class="bg-white shadow overflow-hidden sm:rounded-lg">
        <div class="px-4 py-5 sm:px-6">
          <div class="flex justify-between items-center">
            <h3 class="text-lg leading-6 font-medium text-gray-900">
              {record.title}
            </h3>
            <span
              class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-green-100 text-green-800"
            >
              {getRecordTypeLabel(record.recordType)}
            </span>
          </div>
          <p class="mt-1 max-w-2xl text-sm text-gray-500">
            {formatDate(record.date)}
          </p>
        </div>

        <div class="border-t border-gray-200">
          {#if record.doctor}
            <div
              class="bg-gray-50 px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6"
            >
              <dt class="text-sm font-medium text-gray-500">
                Healthcare Provider
              </dt>
              <dd class="mt-1 text-sm text-gray-900 sm:mt-0 sm:col-span-2">
                {record.doctor}
              </dd>
            </div>
          {/if}

          {#if record.symptoms}
            <div
              class="bg-white px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6"
            >
              <dt class="text-sm font-medium text-gray-500">Symptoms</dt>
              <dd class="mt-1 text-sm text-gray-900 sm:mt-0 sm:col-span-2">
                {record.symptoms}
              </dd>
            </div>
          {/if}

          {#if record.result}
            <div
              class="bg-gray-50 px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6"
            >
              <dt class="text-sm font-medium text-gray-500">Result</dt>
              <dd class="mt-1 text-sm text-gray-900 sm:mt-0 sm:col-span-2">
                {record.result}
              </dd>
            </div>
          {/if}

          <div
            class="bg-white px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6"
          >
            <dt class="text-sm font-medium text-gray-500">Notes</dt>
            <dd class="mt-1 text-sm text-gray-900 sm:mt-0 sm:col-span-2">
              {record.notes || "No additional notes"}

              {#if record.notes && record.notes.includes("test result ID:")}
                {#if record.notes.match(/test result ID: ([a-zA-Z0-9]+)/)}
                  {@const testResultId = record.notes.match(
                    /test result ID: ([a-zA-Z0-9]+)/
                  )[1]}
                  <div class="mt-3">
                    <a
                      href={`/test-kits/analyze/results/${testResultId}`}
                      class="inline-flex items-center px-2.5 py-1.5 border border-transparent text-xs font-medium rounded text-primary-700 bg-primary-100 hover:bg-primary-200 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
                    >
                      View test result
                    </a>
                  </div>
                {/if}
              {/if}
            </dd>
          </div>

          <div
            class="bg-gray-50 px-4 py-5 sm:grid sm:grid-cols-3 sm:gap-4 sm:px-6"
          >
            <dt class="text-sm font-medium text-gray-500">Created</dt>
            <dd class="mt-1 text-sm text-gray-900 sm:mt-0 sm:col-span-2">
              {formatDateTime(record.created_at)}
            </dd>
          </div>
        </div>

        <div class="border-t border-gray-200 px-4 py-5 sm:px-6">
          <div class="flex justify-between">
            <a
              href="/medical-records"
              class="inline-flex items-center px-4 py-2 border border-gray-300 shadow-sm text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
            >
              Back to Records
            </a>

            <button
              type="button"
              class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
            >
              Edit Record
            </button>
          </div>
        </div>
      </div>
    {/if}
  </div>
</div>
