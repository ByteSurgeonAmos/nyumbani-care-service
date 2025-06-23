<script lang="ts">
  import { page } from "$app/stores";
  import { isAuthenticated, currentUser, authService } from "$lib/api";
  import { LogOut, Menu, X } from "lucide-svelte";
  import { onMount } from "svelte";
  import NotificationBell from "$lib/components/NotificationBell.svelte";

  let isMobileMenuOpen = false;

  function toggleMobileMenu() {
    isMobileMenuOpen = !isMobileMenuOpen;
  }

  function handleLogout() {
    authService.logout();
  }

  onMount(() => {
    // Update the current user information if authenticated
    if ($isAuthenticated) {
      authService.getCurrentUser();
    }
  });
</script>

<header class="bg-white shadow">
  <nav class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
    <div class="flex justify-between h-16">
      <div class="flex">
        <div class="flex-shrink-0 flex items-center">
          <a href="/" class="text-2xl font-bold text-primary-700"
            >NyumbaniCare</a
          >
        </div>
        <div class="hidden sm:ml-6 sm:flex sm:space-x-8">
          <a
            href="/"
            class="{$page.url.pathname === '/'
              ? 'border-primary-500 text-gray-900'
              : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'} 
            inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium"
          >
            Home
          </a>
          <a
            href="/test-kits"
            class="{$page.url.pathname.startsWith('/test-kits')
              ? 'border-primary-500 text-gray-900'
              : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'} 
            inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium"
          >
            Test Kits
          </a>
          <a
            href="/education"
            class="{$page.url.pathname.startsWith('/education')
              ? 'border-primary-500 text-gray-900'
              : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'} 
            inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium"
          >
            Health Education
          </a>
          {#if $isAuthenticated}
            <a
              href="/medical-records"
              class="{$page.url.pathname.startsWith('/medical-records')
                ? 'border-primary-500 text-gray-900'
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'} 
              inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium"
            >
              Medical Records
            </a>
            <a
              href="/orders"
              class="{$page.url.pathname.startsWith('/orders')
                ? 'border-primary-500 text-gray-900'
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'} 
              inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium"
            >
              Orders
            </a>
            <a
              href="/api-demo"
              class="{$page.url.pathname.startsWith('/api-demo')
                ? 'border-primary-500 text-gray-900'
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'} 
              inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium"
            >
              API Demo
            </a>
            <a
              href="/test-kits/analyze/results"
              class="{$page.url.pathname.startsWith(
                '/test-kits/analyze/results'
              )
                ? 'border-primary-500 text-gray-900'
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'} 
              inline-flex items-center px-1 pt-1 border-b-2 text-sm font-medium"
            >
              Test Results
            </a>
          {/if}
        </div>
      </div>
      <div class="hidden sm:ml-6 sm:flex sm:items-center">
        {#if $isAuthenticated}
          <div class="ml-3 relative flex items-center gap-4">
            <NotificationBell />
            <span class="text-gray-700"
              >Hi, {$currentUser?.first_name || "User"}</span
            >
            <button
              class="flex items-center text-gray-700 hover:text-primary-700"
              on:click={handleLogout}
            >
              <LogOut class="h-5 w-5 mr-1" />
              <span>Logout</span>
            </button>
          </div>
        {:else}
          <div class="ml-3 relative flex items-center gap-2">
            <a
              href="/login"
              class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-primary-700 hover:bg-primary-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
            >
              Login
            </a>
            <a href="/register" class="btn btn-primary"> Register </a>
          </div>
        {/if}
      </div>
      <div class="-mr-2 flex items-center sm:hidden">
        <button
          on:click={toggleMobileMenu}
          type="button"
          class="bg-white inline-flex items-center justify-center p-2 rounded-md text-gray-400 hover:text-gray-500 hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-primary-500"
          aria-expanded="false"
        >
          <span class="sr-only"
            >{isMobileMenuOpen ? "Close main menu" : "Open main menu"}</span
          >
          {#if isMobileMenuOpen}
            <X class="h-6 w-6" />
          {:else}
            <Menu class="h-6 w-6" />
          {/if}
        </button>
      </div>
    </div>
  </nav>

  {#if isMobileMenuOpen}
    <div class="sm:hidden">
      <div class="pt-2 pb-3 space-y-1">
        <a
          href="/"
          class="{$page.url.pathname === '/'
            ? 'bg-primary-50 border-primary-500 text-primary-700'
            : 'border-transparent text-gray-600 hover:bg-gray-50 hover:border-gray-300 hover:text-gray-800'} 
          block pl-3 pr-4 py-2 border-l-4 text-base font-medium"
        >
          Home
        </a>
        <a
          href="/test-kits"
          class="{$page.url.pathname.startsWith('/test-kits')
            ? 'bg-primary-50 border-primary-500 text-primary-700'
            : 'border-transparent text-gray-600 hover:bg-gray-50 hover:border-gray-300 hover:text-gray-800'} 
          block pl-3 pr-4 py-2 border-l-4 text-base font-medium"
        >
          Test Kits
        </a>
        <a
          href="/education"
          class="{$page.url.pathname.startsWith('/education')
            ? 'bg-primary-50 border-primary-500 text-primary-700'
            : 'border-transparent text-gray-600 hover:bg-gray-50 hover:border-gray-300 hover:text-gray-800'} 
          block pl-3 pr-4 py-2 border-l-4 text-base font-medium"
        >
          Health Education
        </a>
        {#if $isAuthenticated}
          <a
            href="/medical-records"
            class="{$page.url.pathname.startsWith('/medical-records')
              ? 'bg-primary-50 border-primary-500 text-primary-700'
              : 'border-transparent text-gray-600 hover:bg-gray-50 hover:border-gray-300 hover:text-gray-800'} 
            block pl-3 pr-4 py-2 border-l-4 text-base font-medium"
          >
            Medical Records
          </a>
          <a
            href="/orders"
            class="{$page.url.pathname.startsWith('/orders')
              ? 'bg-primary-50 border-primary-500 text-primary-700'
              : 'border-transparent text-gray-600 hover:bg-gray-50 hover:border-gray-300 hover:text-gray-800'} 
            block pl-3 pr-4 py-2 border-l-4 text-base font-medium"
          >
            Orders
          </a>
          <a
            href="/api-demo"
            class="{$page.url.pathname.startsWith('/api-demo')
              ? 'bg-primary-50 border-primary-500 text-primary-700'
              : 'border-transparent text-gray-600 hover:bg-gray-50 hover:border-gray-300 hover:text-gray-800'} 
            block pl-3 pr-4 py-2 border-l-4 text-base font-medium"
          >
            API Demo
          </a>
          <a
            href="/test-kits/analyze/results"
            class="{$page.url.pathname.startsWith('/test-kits/analyze/results')
              ? 'bg-primary-50 border-primary-500 text-primary-700'
              : 'border-transparent text-gray-600 hover:bg-gray-50 hover:border-gray-300 hover:text-gray-800'} 
            block pl-3 pr-4 py-2 border-l-4 text-base font-medium"
          >
            Test Results
          </a>
        {/if}
      </div>
      <div class="pt-4 pb-3 border-t border-gray-200">
        {#if $isAuthenticated}
          <div class="flex items-center px-4">
            <div class="flex-shrink-0">
              <div
                class="h-10 w-10 rounded-full bg-primary-100 flex items-center justify-center text-primary-700 font-bold"
              >
                {$currentUser?.first_name.charAt(0) || "U"}
              </div>
            </div>
            <div class="ml-3 flex-grow">
              <div class="text-base font-medium text-gray-800">
                {$currentUser?.first_name}
                {$currentUser?.last_name}
              </div>
              <div class="text-sm font-medium text-gray-500">
                {$currentUser?.email}
              </div>
            </div>
            <div>
              <NotificationBell />
            </div>
          </div>
          <div class="mt-3 space-y-1">
            <a
              href="/profile"
              class="block px-4 py-2 text-base font-medium text-gray-500 hover:text-gray-800 hover:bg-gray-100"
            >
              Your Profile
            </a>
            <button
              on:click={handleLogout}
              class="block w-full text-left px-4 py-2 text-base font-medium text-gray-500 hover:text-gray-800 hover:bg-gray-100"
            >
              Sign out
            </button>
          </div>
        {:else}
          <div class="mt-3 space-y-1">
            <a
              href="/login"
              class="block px-4 py-2 text-base font-medium text-gray-500 hover:text-gray-800 hover:bg-gray-100"
            >
              Login
            </a>
            <a
              href="/register"
              class="block px-4 py-2 text-base font-medium text-gray-500 hover:text-gray-800 hover:bg-gray-100"
            >
              Register
            </a>
          </div>
        {/if}
      </div>
    </div>
  {/if}
</header>
