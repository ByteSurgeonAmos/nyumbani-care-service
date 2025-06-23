<script lang="ts">
  import { healthEducationService } from "$lib/api";
  import type { HealthArticle } from "$lib/api";
  import { onMount } from "svelte";
  import toast from "svelte-french-toast";

  export let data: { id: string };
  const { id } = data;

  let article: HealthArticle | null = null;
  let isLoading = true;
  let relatedArticles: HealthArticle[] = [];

  onMount(async () => {
    await loadArticle();
  });

  async function loadArticle() {
    isLoading = true;

    try {
      article = await healthEducationService.getArticleById(id);

      if (article) {
        const response = await healthEducationService.getAllArticles(
          1,
          3,
          article.category
        );
        relatedArticles = response.data
          .filter((a) => a.id !== article?.id)
          .slice(0, 3);
      }
    } catch (error) {
      console.error("Failed to load article:", error);
      toast.error("Failed to load article");
    } finally {
      isLoading = false;
    }
  }
</script>

<svelte:head>
  {#if article}
    <title>{article.title} | NyumbaniCare</title>
    <meta name="description" content={article.content.substring(0, 160)} />
  {:else}
    <title>Article | NyumbaniCare</title>
  {/if}
</svelte:head>

<div class="bg-white">
  {#if isLoading}
    <div class="mx-auto max-w-7xl">
      <div class="pt-6 pb-16 sm:pb-24">
        <!-- Skeleton for breadcrumb -->
        <div class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
          <div class="flex items-center space-x-4">
            <div class="h-4 w-32 bg-gray-200 rounded animate-pulse"></div>
            <div class="h-5 w-5 bg-gray-200 rounded-full"></div>
            <div class="h-4 w-60 bg-gray-200 rounded animate-pulse"></div>
          </div>
        </div>

        <div class="mx-auto mt-8 max-w-3xl px-4 sm:px-6 lg:px-8">
          <!-- Skeleton for article header -->
          <div class="border-b border-gray-200 pb-10">
            <div class="flex items-center gap-3">
              <div class="h-6 w-20 bg-gray-200 rounded animate-pulse"></div>
              <div class="h-4 w-32 bg-gray-200 rounded animate-pulse"></div>
            </div>

            <div
              class="h-10 w-3/4 bg-gray-200 rounded animate-pulse mt-4"
            ></div>
            <div class="h-6 w-1/3 bg-gray-200 rounded animate-pulse mt-4"></div>
          </div>

          <!-- Skeleton for article image -->
          <div class="mt-8 h-64 w-full bg-gray-200 rounded animate-pulse"></div>

          <!-- Skeleton for article content -->
          <div class="mt-8">
            <div
              class="h-6 w-full bg-gray-200 rounded animate-pulse mb-4"
            ></div>
            <div
              class="h-6 w-full bg-gray-200 rounded animate-pulse mb-4"
            ></div>
            <div class="h-6 w-3/4 bg-gray-200 rounded animate-pulse mb-4"></div>
            <div
              class="h-6 w-full bg-gray-200 rounded animate-pulse mb-4"
            ></div>
            <div
              class="h-6 w-full bg-gray-200 rounded animate-pulse mb-4"
            ></div>
            <div class="h-6 w-2/3 bg-gray-200 rounded animate-pulse mb-4"></div>
            <div
              class="h-6 w-full bg-gray-200 rounded animate-pulse mb-4"
            ></div>
            <div
              class="h-6 w-full bg-gray-200 rounded animate-pulse mb-4"
            ></div>
            <div class="h-6 w-5/6 bg-gray-200 rounded animate-pulse mb-4"></div>
          </div>

          <!-- Skeleton for related articles -->
          <div class="mt-16 border-t border-gray-200 pt-10">
            <div class="h-8 w-48 bg-gray-200 rounded animate-pulse mb-6"></div>

            <div
              class="grid grid-cols-1 gap-y-10 gap-x-6 sm:grid-cols-2 lg:grid-cols-3"
            >
              {#each Array(3) as _}
                <div class="group relative">
                  <div
                    class="aspect-w-4 aspect-h-3 overflow-hidden rounded-lg bg-gray-200 animate-pulse"
                  ></div>
                  <div class="mt-4">
                    <div
                      class="h-5 w-full bg-gray-200 rounded animate-pulse mb-2"
                    ></div>
                    <div
                      class="h-4 w-full bg-gray-200 rounded animate-pulse"
                    ></div>
                    <div
                      class="h-4 w-2/3 bg-gray-200 rounded animate-pulse mt-1"
                    ></div>
                  </div>
                </div>
              {/each}
            </div>
          </div>
        </div>
      </div>
    </div>
  {:else if !article}
    <div class="text-center py-12">
      <p class="text-lg text-gray-500">Article not found</p>
      <a
        href="/education"
        class="mt-4 inline-flex items-center text-primary-600 hover:text-primary-700"
      >
        Back to all articles
      </a>
    </div>
  {:else}
    <div class="mx-auto max-w-7xl">
      <div class="pt-6 pb-16 sm:pb-24">
        <nav
          aria-label="Breadcrumb"
          class="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8"
        >
          <ol role="list" class="flex items-center space-x-4">
            <li>
              <div class="flex items-center">
                <a
                  href="/education"
                  class="mr-4 text-sm font-medium text-gray-900"
                  >Health Education</a
                >
                <svg
                  viewBox="0 0 6 20"
                  class="h-5 w-auto text-gray-300"
                  fill="currentColor"
                  aria-hidden="true"
                >
                  <path d="M4.878 4.34H3.551L.27 16.532h1.327l3.281-12.19z" />
                </svg>
              </div>
            </li>
            <li class="text-sm">
              <a
                href={`/education/articles/${article.id}`}
                aria-current="page"
                class="font-medium text-gray-500 hover:text-gray-600"
              >
                {article.title}
              </a>
            </li>
          </ol>
        </nav>

        <div class="mx-auto mt-8 max-w-3xl px-4 sm:px-6 lg:px-8">
          <!-- Article header -->
          <div class="border-b border-gray-200 pb-10">
            <div class="flex items-center gap-3">
              <span
                class="inline-flex items-center rounded-md bg-primary-50 px-2 py-1 text-xs font-medium text-primary-700"
              >
                {article.category}
              </span>
              <span class="text-sm text-gray-500">
                {article.published_at
                  ? new Date(article.published_at).toLocaleDateString()
                  : "No date available"}
              </span>
            </div>

            <h1
              class="mt-4 text-3xl font-extrabold tracking-tight text-gray-900 sm:text-4xl"
            >
              {article.title}
            </h1>
            <div class="mt-4 flex items-center">
              <p class="text-sm text-gray-500">
                By {article.author?.first_name && article.author?.last_name
                  ? `${article.author.first_name} ${article.author.last_name}`
                  : "Unknown Author"}
              </p>
            </div>
          </div>

          <!-- Article image -->
          {#if article.image_url}
            <div class="mt-8">
              <img
                src={article.image_url}
                alt={article.title}
                class="w-full h-auto rounded-lg"
              />
            </div>
          {/if}

          <!-- Article content -->
          <div class="mt-8 prose prose-primary prose-lg max-w-none">
            {#each article.content.split("\n\n") as paragraph}
              <p>{paragraph}</p>
            {/each}
          </div>

          <!-- Related articles -->
          {#if relatedArticles.length > 0}
            <div class="mt-16 border-t border-gray-200 pt-10">
              <h2 class="text-xl font-bold tracking-tight text-gray-900">
                Related Articles
              </h2>

              <div
                class="mt-6 grid grid-cols-1 gap-y-10 gap-x-6 sm:grid-cols-2 lg:grid-cols-3"
              >
                {#each relatedArticles as relatedArticle (relatedArticle.id)}
                  <div class="group relative">
                    <div
                      class="aspect-w-4 aspect-h-3 overflow-hidden rounded-lg bg-gray-100"
                    >
                      <img
                        src={relatedArticle.image_url ||
                          "https://via.placeholder.com/400x300?text=Health+Article"}
                        alt={relatedArticle.title}
                        class="h-full w-full object-cover object-center group-hover:opacity-75"
                      />
                    </div>
                    <div class="mt-4">
                      <h3 class="text-sm font-medium text-gray-900">
                        <a href={`/education/articles/${relatedArticle.id}`}>
                          <span aria-hidden="true" class="absolute inset-0"
                          ></span>
                          {relatedArticle.title}
                        </a>
                      </h3>
                      <p class="mt-1 text-sm text-gray-500">
                        {relatedArticle.content.substring(0, 100)}...
                      </p>
                    </div>
                  </div>
                {/each}
              </div>
            </div>
          {/if}
        </div>
      </div>
    </div>
  {/if}
</div>
