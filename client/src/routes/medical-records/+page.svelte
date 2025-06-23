<script lang="ts">
  import { onMount } from "svelte";
  import { goto } from "$app/navigation";
  import { isAuthenticated } from "$lib/api";
  import toast from "svelte-french-toast";

  interface MedicalRecord {
    id: string;
    title: string;
    date: string;
    recordType: string;
    doctor?: string;
    notes: string;
  }

  let records: MedicalRecord[] = []; // Will be populated from API
  let isLoading = true;
  let error: string | null = null;

  // Mock data for the prototype
  const mockRecords = [
    {
      id: "1",
      title: "Annual Checkup",
      date: "2025-05-10",
      recordType: "doctor_visit",
      doctor: "Dr. Smith",
      notes: "Regular annual checkup. All vitals normal.",
    },
    {
      id: "2",
      title: "COVID-19 Test Result",
      date: "2025-05-01",
      recordType: "test_result",
      notes: "Associated test result ID: abc12345",
    },
    {
      id: "3",
      title: "Flu Vaccination",
      date: "2025-04-15",
      recordType: "vaccination",
      notes: "Annual flu vaccine administered.",
    },
  ];

  onMount(async () => {
    if ($isAuthenticated) {
      loadRecords();
    } else {
      toast.error("Please login to view your medical records");
      goto("/login?redirectTo=/medical-records");
    }
  });

  async function loadRecords() {
    isLoading = true;
    error = null;

    try {
      // Simulate API call with mock data
      await new Promise((resolve) => setTimeout(resolve, 800));
      records = mockRecords;
    } catch (err) {
      console.error("Failed to load medical records:", err);
      error = "Failed to load medical records. Please try again later.";
    } finally {
      isLoading = false;
    }
  }

  function getRecordTypeIcon(recordType: string) {
    switch (recordType) {
      case "test_result":
        return `<svg class="h-5 w-5 text-primary-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
        </svg>`;
      case "doctor_visit":
        return `<svg class="h-5 w-5 text-green-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
        </svg>`;
      case "vaccination":
        return `<svg class="h-5 w-5 text-blue-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19.428 15.428a2 2 0 00-1.022-.547l-2.387-.477a6 6 0 00-3.86.517l-.318.158a6 6 0 01-3.86.517L6.05 15.21a2 2 0 00-1.806.547M8 4h8l-1 1v5.172a2 2 0 00.586 1.414l5 5c1.26 1.26.367 3.414-1.415 3.414H4.828c-1.782 0-2.674-2.154-1.414-3.414l5-5A2 2 0 009 10.172V5L8 4z" />
        </svg>`;
      case "prescription":
        return `<svg class="h-5 w-5 text-red-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
        </svg>`;
      default:
        return `<svg class="h-5 w-5 text-gray-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
        </svg>`;
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

  function viewRecord(id: string) {
    // Navigate to record details
    goto(`/medical-records/${id}`);
  }
</script>

<svelte:head>
  <title>Medical Records | NyumbaniCare</title>
  <meta name="description" content="View and manage your medical records" />
</svelte:head>

<div class="bg-white">
  <div class="max-w-7xl mx-auto px-4 py-16 sm:px-6 lg:px-8">
    <div class="sm:flex sm:items-center sm:justify-between">
      <h1 class="text-3xl font-bold tracking-tight text-gray-900">
        My Medical Records
      </h1>
      <div class="mt-4 sm:mt-0">
        <a
          href="/medical-records/add"
          class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
        >
          Add New Record
        </a>
      </div>
    </div>

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
      {:else if records.length === 0}
        <div class="bg-white shadow overflow-hidden sm:rounded-lg">
          <div class="px-4 py-12 text-center sm:px-6">
            <p class="text-lg text-gray-500">
              You haven't added any medical records yet.
            </p>
            <div class="mt-6">
              <a
                href="/medical-records/add"
                class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-primary-600 hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
              >
                Add Your First Record
              </a>
            </div>
          </div>
        </div>
      {:else}
        <div class="overflow-hidden bg-white shadow sm:rounded-md">
          <ul role="list" class="divide-y divide-gray-200">
            {#each records as record}
              <li>
                <button
                  class="block w-full text-left hover:bg-gray-50"
                  on:click={() => viewRecord(record.id)}
                >
                  <div class="px-4 py-4 sm:px-6">
                    <div class="flex items-center justify-between">
                      <div class="flex items-center">
                        <div class="flex-shrink-0 mr-3">
                          {@html getRecordTypeIcon(record.recordType)}
                        </div>
                        <p
                          class="text-sm font-medium text-primary-600 truncate"
                        >
                          {record.title}
                        </p>
                      </div>
                      <div class="ml-2 flex-shrink-0 flex">
                        <p
                          class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-green-100 text-green-800"
                        >
                          {record.recordType.replace("_", " ")}
                        </p>
                      </div>
                    </div>
                    <div class="mt-2 sm:flex sm:justify-between">
                      <div class="sm:flex">
                        {#if record.doctor}
                          <p class="flex items-center text-sm text-gray-500">
                            <svg
                              class="flex-shrink-0 mr-1.5 h-5 w-5 text-gray-400"
                              xmlns="http://www.w3.org/2000/svg"
                              viewBox="0 0 20 20"
                              fill="currentColor"
                            >
                              <path
                                fill-rule="evenodd"
                                d="M10 9a3 3 0 100-6 3 3 0 000 6zm-7 9a7 7 0 1114 0H3z"
                                clip-rule="evenodd"
                              />
                            </svg>
                            <span>{record.doctor}</span>
                          </p>
                        {/if}
                        {#if record.notes}
                          <p
                            class="mt-2 flex items-center text-sm text-gray-500 sm:mt-0 sm:ml-6"
                          >
                            <svg
                              class="flex-shrink-0 mr-1.5 h-5 w-5 text-gray-400"
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
                            <span>
                              {record.notes.length > 60
                                ? record.notes.substring(0, 60) + "..."
                                : record.notes}
                            </span>
                          </p>
                        {/if}
                      </div>
                      <div
                        class="mt-2 flex items-center text-sm text-gray-500 sm:mt-0"
                      >
                        <svg
                          class="flex-shrink-0 mr-1.5 h-5 w-5 text-gray-400"
                          xmlns="http://www.w3.org/2000/svg"
                          viewBox="0 0 20 20"
                          fill="currentColor"
                        >
                          <path
                            fill-rule="evenodd"
                            d="M6 2a1 1 0 00-1 1v1H4a2 2 0 00-2 2v10a2 2 0 002 2h12a2 2 0 002-2V6a2 2 0 00-2-2h-1V3a1 1 0 10-2 0v1H7V3a1 1 0 00-1-1zm0 5a1 1 0 000 2h8a1 1 0 100-2H6z"
                            clip-rule="evenodd"
                          />
                        </svg>
                        <span>{formatDate(record.date)}</span>
                      </div>
                    </div>
                  </div>
                </button>
              </li>
            {/each}
          </ul>
        </div>
      {/if}
    </div>
  </div>
</div>
