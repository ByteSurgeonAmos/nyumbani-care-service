<script lang="ts">
  import { testKitService, healthEducationService } from "$lib/api";
  import { onMount } from "svelte";
  import { ArrowRight } from "lucide-svelte";
  import type { TestKit, HealthArticle } from "$lib/api";

  let testKits: TestKit[] = [];
  let articles: HealthArticle[] = [];
  let isTestKitsLoading = true;
  let isArticlesLoading = true;

  onMount(async () => {
    Promise.all([loadTestKits(), loadArticles()]);
  });

  async function loadTestKits() {
    isTestKitsLoading = true;
    try {
      const testKitsResponse = await testKitService.getAll(1, 4);
      testKits = testKitsResponse.data;
    } catch (testKitError) {
      console.error("Failed to fetch test kits:", testKitError);
      testKits = [];
    } finally {
      isTestKitsLoading = false;
    }
  }

  async function loadArticles() {
    isArticlesLoading = true;
    try {
      const articlesResponse = await healthEducationService.getAllArticles(
        1,
        3,
        ""
      );
      articles = articlesResponse.data;
    } catch (articlesError) {
      console.error("Failed to fetch health articles:", articlesError);
      articles = [];
    } finally {
      isArticlesLoading = false;
    }
  }
</script>

<svelte:head>
  <title>NyumbaniCare - Healthcare at Home</title>
  <meta
    name="description"
    content="NyumbaniCare provides accessible healthcare services at home including test kits, consultations, and health education."
  />
</svelte:head>

<!-- Hero Section -->
<div class="relative overflow-hidden bg-white">
  <div class="mx-auto max-w-7xl">
    <div
      class="relative z-10 bg-white pb-8 sm:pb-16 md:pb-20 lg:w-full lg:max-w-2xl lg:pb-28 xl:pb-32"
    >
      <main
        class="mx-auto mt-10 max-w-7xl px-4 sm:mt-12 sm:px-6 md:mt-16 lg:mt-20 lg:px-8 xl:mt-28"
      >
        <div class="sm:text-center lg:text-left">
          <h1
            class="text-4xl font-bold tracking-tight text-gray-900 sm:text-5xl md:text-6xl"
          >
            <span class="block xl:inline">Healthcare delivered</span>
            <span class="block text-primary-600 xl:inline"> at your home</span>
          </h1>
          <p
            class="mt-3 text-base text-gray-500 sm:mx-auto sm:mt-5 sm:max-w-xl sm:text-lg md:mt-5 md:text-xl lg:mx-0"
          >
            Get access to quality healthcare services, test kits, and
            consultations from the comfort of your home. Take control of your
            health with NyumbaniCare.
          </p>
          <div class="mt-5 sm:mt-8 sm:flex sm:justify-center lg:justify-start">
            <div class="rounded-md shadow">
              <a
                href="/test-kits"
                class="flex w-full items-center justify-center rounded-md border border-transparent bg-primary-600 px-8 py-3 text-base font-medium text-white hover:bg-primary-700 md:py-4 md:px-10 md:text-lg"
              >
                Browse Test Kits
              </a>
            </div>
            <div class="mt-3 sm:mt-0 sm:ml-3">
              <a
                href="/education"
                class="flex w-full items-center justify-center rounded-md border border-transparent bg-secondary-100 px-8 py-3 text-base font-medium text-secondary-700 hover:bg-secondary-200 md:py-4 md:px-10 md:text-lg"
              >
                Health Education
              </a>
            </div>
          </div>
        </div>
      </main>
    </div>
  </div>
  <div class="lg:absolute lg:inset-y-0 lg:right-0 lg:w-1/2">
    <img
      class="h-56 w-full object-cover sm:h-72 md:h-96 lg:h-full lg:w-full"
      src="https://images.unsplash.com/photo-1551884170-09fb70a3a2ed?ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&w=1974&q=80"
      alt="Healthcare professional"
    />
  </div>
</div>

<!-- Features Section -->
<div class="py-12 bg-white">
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
    <div class="lg:text-center">
      <h2
        class="text-base text-primary-600 font-semibold tracking-wide uppercase"
      >
        Our Services
      </h2>
      <p
        class="mt-2 text-3xl leading-8 font-extrabold tracking-tight text-gray-900 sm:text-4xl"
      >
        A better way to manage your health
      </p>
      <p class="mt-4 max-w-2xl text-xl text-gray-500 lg:mx-auto">
        NyumbaniCare provides a comprehensive set of services designed to make
        healthcare accessible, convenient, and affordable.
      </p>
    </div>

    <div class="mt-10">
      <dl
        class="space-y-10 md:grid md:grid-cols-2 md:gap-x-8 md:gap-y-10 md:space-y-0"
      >
        <div class="relative">
          <dt>
            <div
              class="absolute flex h-12 w-12 items-center justify-center rounded-md bg-primary-500 text-white"
            >
              <svg
                class="h-6 w-6"
                fill="none"
                viewBox="0 0 24 24"
                stroke-width="1.5"
                stroke="currentColor"
                aria-hidden="true"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M9.75 3.104v5.714a2.25 2.25 0 01-.659 1.591L5 14.5M9.75 3.104c-.251.023-.501.05-.75.082m.75-.082a24.301 24.301 0 014.5 0m0 0v5.714a2.25 2.25 0 01-.659 1.591L9.5 14.5M15 3.104c.251.023.501.05.75.082M15 3.104a24.301 24.301 0 01-4.5 0m0 0v5.714a2.25 2.25 0 01-.659 1.591L5 14.5M9.75 3.104V1.75h4.5v1.354a24.301 24.301 0 01-4.5 0z"
                />
              </svg>
            </div>
            <p class="ml-16 text-lg leading-6 font-medium text-gray-900">
              Test Kits
            </p>
          </dt>
          <dd class="mt-2 ml-16 text-base text-gray-500">
            Order diagnostic test kits delivered directly to your home. Get
            results quickly and securely.
          </dd>
        </div>

        <div class="relative">
          <dt>
            <div
              class="absolute flex h-12 w-12 items-center justify-center rounded-md bg-primary-500 text-white"
            >
              <svg
                class="h-6 w-6"
                fill="none"
                viewBox="0 0 24 24"
                stroke-width="1.5"
                stroke="currentColor"
                aria-hidden="true"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M15.75 6a3.75 3.75 0 11-7.5 0 3.75 3.75 0 017.5 0zM4.501 20.118a7.5 7.5 0 0114.998 0A17.933 17.933 0 0112 21.75c-2.676 0-5.216-.584-7.499-1.632z"
                />
              </svg>
            </div>
            <p class="ml-16 text-lg leading-6 font-medium text-gray-900">
              Virtual Consultations
            </p>
          </dt>
          <dd class="mt-2 ml-16 text-base text-gray-500">
            Connect with healthcare professionals via telehealth for
            consultations, follow-ups, and medical advice.
          </dd>
        </div>

        <div class="relative">
          <dt>
            <div
              class="absolute flex h-12 w-12 items-center justify-center rounded-md bg-primary-500 text-white"
            >
              <svg
                class="h-6 w-6"
                fill="none"
                viewBox="0 0 24 24"
                stroke-width="1.5"
                stroke="currentColor"
                aria-hidden="true"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M12 6.042A8.967 8.967 0 006 3.75c-1.052 0-2.062.18-3 .512v14.25A8.987 8.987 0 016 18c2.305 0 4.408.867 6 2.292m0-14.25a8.966 8.966 0 016-2.292c1.052 0 2.062.18 3 .512v14.25A8.987 8.987 0 0018 18a8.967 8.967 0 00-6 2.292m0-14.25v14.25"
                />
              </svg>
            </div>
            <p class="ml-16 text-lg leading-6 font-medium text-gray-900">
              Health Education
            </p>
          </dt>
          <dd class="mt-2 ml-16 text-base text-gray-500">
            Access a wealth of health education resources, articles, and guides
            to help you make informed health decisions.
          </dd>
        </div>

        <div class="relative">
          <dt>
            <div
              class="absolute flex h-12 w-12 items-center justify-center rounded-md bg-primary-500 text-white"
            >
              <svg
                class="h-6 w-6"
                fill="none"
                viewBox="0 0 24 24"
                stroke-width="1.5"
                stroke="currentColor"
                aria-hidden="true"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M9 12h3.75M9 15h3.75M9 18h3.75m3 .75H18a2.25 2.25 0 002.25-2.25V6.108c0-1.135-.845-2.098-1.976-2.192a48.424 48.424 0 00-1.123-.08m-5.801 0c-.065.21-.1.433-.1.664 0 .414.336.75.75.75h4.5a.75.75 0 00.75-.75 2.25 2.25 0 00-.1-.664m-5.8 0A2.251 2.251 0 0113.5 2.25H15c1.012 0 1.867.668 2.15 1.586m-5.8 0c-.376.023-.75.05-1.124.08C9.095 4.01 8.25 4.973 8.25 6.108V8.25m0 0H4.875c-.621 0-1.125.504-1.125 1.125v11.25c0 .621.504 1.125 1.125 1.125h9.75c.621 0 1.125-.504 1.125-1.125V9.375c0-.621-.504-1.125-1.125-1.125H8.25zM6.75 12h.008v.008H6.75V12zm0 3h.008v.008H6.75V15zm0 3h.008v.008H6.75V18z"
                />
              </svg>
            </div>
            <p class="ml-16 text-lg leading-6 font-medium text-gray-900">
              Health Records
            </p>
          </dt>
          <dd class="mt-2 ml-16 text-base text-gray-500">
            Store and access your medical records securely. Share them with
            healthcare providers as needed.
          </dd>
        </div>
      </dl>
    </div>
  </div>
</div>

<!-- Featured Test Kits Section -->
<div class="bg-gray-50 py-12">
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
    <div class="lg:text-center mb-10">
      <h2
        class="text-base text-primary-600 font-semibold tracking-wide uppercase"
      >
        Popular Test Kits
      </h2>
      <p
        class="mt-2 text-3xl leading-8 font-extrabold tracking-tight text-gray-900 sm:text-4xl"
      >
        Take control of your health
      </p>
    </div>

    {#if isTestKitsLoading}
      <div
        class="grid grid-cols-1 gap-y-10 gap-x-6 sm:grid-cols-2 lg:grid-cols-4 xl:gap-x-8"
      >
        {#each Array(4) as _ (Math.random())}
          <div
            class="group relative bg-white rounded-lg shadow-md overflow-hidden"
          >
            <!-- Skeleton for test kit image -->
            <div
              class="w-full min-h-80 aspect-w-1 aspect-h-1 overflow-hidden bg-gray-200 animate-pulse lg:aspect-none lg:h-48"
            ></div>

            <div class="p-4">
              <!-- Skeleton for test kit title -->
              <div class="h-6 w-3/4 bg-gray-200 rounded animate-pulse"></div>

              <!-- Skeleton for description -->
              <div
                class="mt-1 h-4 w-full bg-gray-200 rounded animate-pulse"
              ></div>

              <!-- Skeleton for price and category -->
              <div class="mt-2 flex justify-between items-center">
                <div class="h-5 w-16 bg-gray-200 rounded animate-pulse"></div>
                <div class="h-4 w-20 bg-gray-200 rounded animate-pulse"></div>
              </div>
            </div>
          </div>
        {/each}
      </div>
    {:else if testKits.length === 0}
      <div class="text-center py-8">
        <p class="text-gray-500">No test kits available at the moment.</p>
        <p class="text-gray-500 text-sm mt-2">
          Please check back later for our available test kits!
        </p>
      </div>
    {:else}
      <div
        class="grid grid-cols-1 gap-y-10 gap-x-6 sm:grid-cols-2 lg:grid-cols-4 xl:gap-x-8"
      >
        {#each testKits as kit (kit.id)}
          <div
            class="group relative bg-white rounded-lg shadow-md overflow-hidden"
          >
            <div
              class="w-full min-h-80 aspect-w-1 aspect-h-1 overflow-hidden group-hover:opacity-75 lg:aspect-none lg:h-48"
            >
              <img
                src={kit.image_url ||
                  "https://via.placeholder.com/300x200?text=Test+Kit"}
                alt={kit.name}
                class="w-full h-full object-center object-cover lg:w-full lg:h-full"
              />
            </div>
            <div class="p-4">
              <h3 class="text-lg font-medium text-gray-900">
                <a href={`/test-kits/${kit.id}`}>
                  <span aria-hidden="true" class="absolute inset-0"></span>
                  {kit.name}
                </a>
              </h3>
              <p class="mt-1 text-sm text-gray-500">
                {kit.description
                  ? kit.description.substring(0, 60) + "..."
                  : "No description available"}
              </p>
              <div class="mt-2 flex justify-between items-center">
                <p class="text-lg font-medium text-primary-600">
                  KES {kit.price ? kit.price.toFixed(2) : "0.00"}
                </p>
                <p class="text-sm font-medium text-gray-500">
                  {kit.category || "General"}
                </p>
              </div>
            </div>
          </div>
        {/each}
      </div>

      <div class="mt-8 text-center">
        <a
          href="/test-kits"
          class="inline-flex items-center px-4 py-2 border border-transparent text-base font-medium rounded-md shadow-sm text-white bg-primary-600 hover:bg-primary-700"
        >
          View All Test Kits
          <ArrowRight class="ml-2 h-4 w-4" />
        </a>
      </div>
    {/if}
  </div>
</div>

<!-- Health Education Articles -->
<div class="bg-white py-12">
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
    <div class="lg:text-center mb-10">
      <h2
        class="text-base text-primary-600 font-semibold tracking-wide uppercase"
      >
        Health Education
      </h2>
      <p
        class="mt-2 text-3xl leading-8 font-extrabold tracking-tight text-gray-900 sm:text-4xl"
      >
        Latest Health Articles
      </p>
    </div>

    {#if isArticlesLoading}
      <div class="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
        {#each Array(3) as _ (Math.random())}
          <div
            class="group relative bg-white rounded-lg shadow-md overflow-hidden"
          >
            <!-- Skeleton for article image -->
            <div
              class="w-full h-48 overflow-hidden bg-gray-200 animate-pulse"
            ></div>

            <div class="p-4">
              <!-- Skeleton for category -->
              <div
                class="h-4 w-16 bg-gray-200 rounded animate-pulse mb-1"
              ></div>

              <!-- Skeleton for title -->
              <div
                class="h-6 w-full bg-gray-200 rounded animate-pulse mb-2"
              ></div>

              <!-- Skeleton for content -->
              <div
                class="h-4 w-full bg-gray-200 rounded animate-pulse mb-1"
              ></div>
              <div
                class="h-4 w-full bg-gray-200 rounded animate-pulse mb-1"
              ></div>
              <div class="h-4 w-3/4 bg-gray-200 rounded animate-pulse"></div>

              <!-- Skeleton for author and date -->
              <div class="mt-3 flex justify-between">
                <div class="h-3 w-24 bg-gray-200 rounded animate-pulse"></div>
                <div class="h-3 w-24 bg-gray-200 rounded animate-pulse"></div>
              </div>
            </div>
          </div>
        {/each}
      </div>
    {:else if articles.length === 0}
      <div class="text-center py-8">
        <p class="text-gray-500">No health articles available at the moment.</p>
        <p class="text-gray-500 text-sm mt-2">
          Check back soon for new content!
        </p>
      </div>
    {:else}
      <div class="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
        {#each articles as article (article.id)}
          <div
            class="group relative bg-white rounded-lg shadow-md overflow-hidden"
          >
            <div class="w-full h-48 overflow-hidden">
              <img
                src={article.image_url ||
                  "https://via.placeholder.com/400x200?text=Health+Article"}
                alt={article.title}
                class="w-full h-full object-center object-cover group-hover:opacity-75"
              />
            </div>
            <div class="p-4">
              <p class="text-sm text-primary-600 mb-1">
                {article.category || "General"}
              </p>
              <h3 class="text-lg font-semibold text-gray-900 mb-2">
                <a href={`/education/articles/${article.id}`}>
                  <span aria-hidden="true" class="absolute inset-0"></span>
                  {article.title || "Untitled Article"}
                </a>
              </h3>
              <p class="text-sm text-gray-500">
                {article.content
                  ? article.content.substring(0, 100) + "..."
                  : "No content available"}
              </p>
              <div class="mt-3 flex justify-between">
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

      <div class="mt-8 text-center">
        <a
          href="/education/articles"
          class="inline-flex items-center px-4 py-2 border border-transparent text-base font-medium rounded-md shadow-sm text-white bg-primary-600 hover:bg-primary-700"
        >
          Read More Articles
          <ArrowRight class="ml-2 h-4 w-4" />
        </a>
      </div>
    {/if}
  </div>
</div>

<!-- Test Analysis Section -->
<div class="py-12 bg-white">
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
    <div class="lg:text-center">
      <h2
        class="text-base text-primary-600 font-semibold tracking-wide uppercase"
      >
        Test Kit Analysis
      </h2>
      <p
        class="mt-2 text-3xl leading-8 font-extrabold tracking-tight text-gray-900 sm:text-4xl"
      >
        Get immediate results from your test kits
      </p>
      <p class="mt-4 max-w-2xl text-xl text-gray-500 lg:mx-auto">
        Our AI-powered analysis provides quick, accurate interpretations of your
        test kit results.
      </p>
    </div>

    <div class="mt-10 max-w-lg mx-auto grid gap-6 lg:grid-cols-2 lg:max-w-none">
      <div class="flex flex-col rounded-lg shadow-lg overflow-hidden">
        <div class="flex-1 bg-white p-6 flex flex-col justify-between">
          <div class="flex-1">
            <div class="block mt-2">
              <p class="text-xl font-semibold text-gray-900">
                Analyze a Test Kit
              </p>
              <p class="mt-3 text-base text-gray-500">
                Upload an image of your test kit result and get an immediate
                analysis with AI-powered technology.
              </p>
            </div>
          </div>
          <div class="mt-6">
            <a
              href="/test-kits/analyze"
              class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-primary-600 hover:bg-primary-700"
            >
              Analyze Now
            </a>
          </div>
        </div>
      </div>

      <div class="flex flex-col rounded-lg shadow-lg overflow-hidden">
        <div class="flex-1 bg-white p-6 flex flex-col justify-between">
          <div class="flex-1">
            <div class="block mt-2">
              <p class="text-xl font-semibold text-gray-900">
                View Test Results
              </p>
              <p class="mt-3 text-base text-gray-500">
                Access your past test kit analyses, track your health over time,
                and share results with healthcare providers.
              </p>
            </div>
          </div>
          <div class="mt-6">
            <a
              href="/test-kits/analyze/results"
              class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-primary-600 hover:bg-primary-700"
            >
              View Results
            </a>
          </div>
        </div>
      </div>
    </div>
  </div>
</div>

<!-- Call to Action -->
<div class="bg-primary-700">
  <div
    class="max-w-2xl mx-auto text-center py-16 px-4 sm:py-20 sm:px-6 lg:px-8"
  >
    <h2 class="text-3xl font-extrabold text-white sm:text-4xl">
      <span class="block">Ready to take control of your health?</span>
      <span class="block">Start with NyumbaniCare today.</span>
    </h2>
    <p class="mt-4 text-lg leading-6 text-primary-100">
      Sign up now and get access to our complete suite of healthcare services
      from the comfort of your home.
    </p>
    <a
      href="/register"
      class="mt-8 w-full inline-flex items-center justify-center px-5 py-3 border border-transparent text-base font-medium rounded-md text-primary-700 bg-white hover:bg-primary-50 sm:w-auto"
    >
      Create Your Account
    </a>
  </div>
</div>

<style>
  /* Global style overrides can be added here if needed */
  h1 {
    width: 100%;
  }
</style>
