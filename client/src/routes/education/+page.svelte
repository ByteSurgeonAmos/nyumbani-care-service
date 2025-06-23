<script lang="ts">
  import { healthEducationService } from "$lib/api";
  import type { HealthArticle } from "$lib/api";
  import { onMount } from "svelte";
  import toast from "svelte-french-toast";

  let articles: HealthArticle[] = [];
  let isLoading = true;
  let currentPage = 1;
  let totalArticles = 0;
  let articlesPerPage = 9;
  let totalPages = 0;
  let categories: string[] = [];
  let selectedCategory = "";

  onMount(async () => {
    loadArticles();
  });

  async function loadArticles() {
    isLoading = true;

    try {
      const response = await healthEducationService.getAllArticles(
        currentPage,
        articlesPerPage,
        selectedCategory
      );
      articles = response.data;
      totalArticles = response.total;
      totalPages = Math.ceil(totalArticles / articlesPerPage);

      if (categories.length === 0) {
        const uniqueCategories = new Set(
          articles.map((article) => article.category)
        );
        categories = Array.from(uniqueCategories);
      }
    } catch (error) {
      console.error("Failed to load articles:", error);
      toast.error("Failed to load health articles");
    } finally {
      isLoading = false;
    }
  }
  function changePage(newPage: number): void {
    if (newPage >= 1 && newPage <= totalPages && newPage !== currentPage) {
      currentPage = newPage;
      loadArticles();
    }
  }

  function filterByCategory(category: string): void {
    selectedCategory = category;
    currentPage = 1;
    loadArticles();
  }
</script>

<svelte:head>
  <title>Health Education | NyumbaniCare</title>
  <meta
    name="description"
    content="Browse health education articles and resources to help you make informed health decisions."
  />
</svelte:head>

<div class="bg-white">
  <div class="mx-auto max-w-7xl px-4 py-16 sm:px-6 lg:px-8">
    <h1 class="text-3xl font-bold tracking-tight text-gray-900">
      Health Education
    </h1>
    <p class="mt-4 max-w-3xl text-base text-gray-500">
      Browse our collection of health articles, guides, and resources to help
      you make informed decisions about your health.
    </p>

    <!-- Category filter -->
    {#if categories.length > 0}
      <div class="mt-8">
        <div class="flex flex-wrap gap-2">
          <button
            class="px-3 py-1 rounded-full text-sm font-medium {selectedCategory ===
            ''
              ? 'bg-primary-100 text-primary-800'
              : 'bg-gray-100 text-gray-800 hover:bg-gray-200'}"
            on:click={() => filterByCategory("")}
          >
            All Categories
          </button>
          {#each categories as category}
            <button
              class="px-3 py-1 rounded-full text-sm font-medium {selectedCategory ===
              category
                ? 'bg-primary-100 text-primary-800'
                : 'bg-gray-100 text-gray-800 hover:bg-gray-200'}"
              on:click={() => filterByCategory(category)}
            >
              {category}
            </button>
          {/each}
        </div>
      </div>
    {/if}

    {#if isLoading}
      <div
        class="mt-8 grid grid-cols-1 gap-y-10 gap-x-6 sm:grid-cols-2 lg:grid-cols-3"
      >
        {#each Array(6) as _, i}
          <div class="group relative">
            <!-- Skeleton for image -->
            <div
              class="aspect-w-4 aspect-h-3 overflow-hidden rounded-lg bg-gray-200 animate-pulse"
            ></div>
            <!-- Skeleton for category -->
            <div class="mt-4">
              <div class="h-4 w-20 bg-gray-200 rounded animate-pulse"></div>
              <!-- Skeleton for title -->
              <div
                class="h-6 w-full bg-gray-200 rounded animate-pulse mt-2"
              ></div>
              <!-- Skeleton for content preview -->
              <div
                class="h-4 w-full bg-gray-200 rounded animate-pulse mt-3"
              ></div>
              <div
                class="h-4 w-full bg-gray-200 rounded animate-pulse mt-1"
              ></div>
              <!-- Skeleton for author and date -->
              <div class="mt-3 flex justify-between items-center">
                <div class="h-3 w-24 bg-gray-200 rounded animate-pulse"></div>
                <div class="h-3 w-32 bg-gray-200 rounded animate-pulse"></div>
              </div>
            </div>
          </div>
        {/each}
      </div>
    {:else if articles.length === 0}
      <div class="text-center py-12">
        <p class="text-lg text-gray-500">
          No articles available at the moment.
        </p>
      </div>
    {:else}
      <div
        class="mt-8 grid grid-cols-1 gap-y-10 gap-x-6 sm:grid-cols-2 lg:grid-cols-3"
      >
        {#each articles as article (article.id)}
          <div class="group relative">
            <div
              class="aspect-w-4 aspect-h-3 overflow-hidden rounded-lg bg-gray-100"
            >
              <img
                src={article.image_url ||
                  "https://via.placeholder.com/600x400?text=Health+Article"}
                alt={article.title}
                class="h-full w-full object-cover object-center group-hover:opacity-75"
              />
            </div>
            <div class="mt-4">
              <p class="text-sm text-primary-600">{article.category}</p>
              <h3 class="mt-1 text-lg font-medium text-gray-900">
                <a href={`/education/articles/${article.id}`}>
                  <span aria-hidden="true" class="absolute inset-0"></span>
                  {article.title}
                </a>
              </h3>
              <p class="mt-2 text-sm text-gray-500">
                {article.content.substring(0, 120)}...
              </p>
              <div class="mt-3 flex justify-between items-center">
                <p class="text-xs text-gray-500">
                  By {article.author?.first_name && article.author?.last_name
                    ? `${article.author.first_name} ${article.author.last_name}`
                    : "Unknown Author"}
                </p>
                <p class="text-xs text-gray-500">
                  {article.created_at
                    ? new Date(article.created_at).toLocaleDateString()
                    : "Unknown date"}
                </p>
              </div>
            </div>
          </div>
        {/each}
      </div>

      <!-- Pagination -->
      {#if totalPages > 1}
        <div class="flex justify-center mt-10">
          <nav
            class="inline-flex -space-x-px rounded-md shadow-sm"
            aria-label="Pagination"
          >
            <button
              class="inline-flex items-center rounded-l-md px-2 py-2 text-gray-400 ring-1 ring-inset ring-gray-300 hover:bg-gray-50 focus:z-20 focus:outline-offset-0 {currentPage ===
              1
                ? 'cursor-not-allowed opacity-50'
                : ''}"
              on:click={() => changePage(currentPage - 1)}
              disabled={currentPage === 1}
            >
              <span class="sr-only">Previous</span>
              <svg
                class="h-5 w-5"
                viewBox="0 0 20 20"
                fill="currentColor"
                aria-hidden="true"
              >
                <path
                  fill-rule="evenodd"
                  d="M12.79 5.23a.75.75 0 01-.02 1.06L8.832 10l3.938 3.71a.75.75 0 11-1.04 1.08l-4.5-4.25a.75.75 0 010-1.08l4.5-4.25a.75.75 0 011.06.02z"
                  clip-rule="evenodd"
                />
              </svg>
            </button>

            {#each Array(totalPages) as _, i}
              <button
                class="inline-flex items-center px-4 py-2 text-sm font-semibold {currentPage ===
                i + 1
                  ? 'bg-primary-600 text-white focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-primary-600'
                  : 'text-gray-900 ring-1 ring-inset ring-gray-300 hover:bg-gray-50 focus:z-20 focus:outline-offset-0'}"
                on:click={() => changePage(i + 1)}
              >
                {i + 1}
              </button>
            {/each}

            <button
              class="inline-flex items-center rounded-r-md px-2 py-2 text-gray-400 ring-1 ring-inset ring-gray-300 hover:bg-gray-50 focus:z-20 focus:outline-offset-0 {currentPage ===
              totalPages
                ? 'cursor-not-allowed opacity-50'
                : ''}"
              on:click={() => changePage(currentPage + 1)}
              disabled={currentPage === totalPages}
            >
              <span class="sr-only">Next</span>
              <svg
                class="h-5 w-5"
                viewBox="0 0 20 20"
                fill="currentColor"
                aria-hidden="true"
              >
                <path
                  fill-rule="evenodd"
                  d="M7.21 14.77a.75.75 0 01.02-1.06L11.168 10 7.23 6.29a.75.75 0 111.04-1.08l4.5 4.25a.75.75 0 010 1.08l-4.5 4.25a.75.75 0 01-1.06-.02z"
                  clip-rule="evenodd"
                />
              </svg>
            </button>
          </nav>
        </div>
      {/if}
    {/if}
  </div>
</div>
