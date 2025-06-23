<script lang="ts">
  import { authService } from "$lib/api";
  import { goto } from "$app/navigation";
  import toast from "svelte-french-toast";

  let email = "";
  let password = "";
  let firstName = "";
  let lastName = "";
  let phoneNumber = "";
  let dateOfBirth = "";
  let gender = "";
  let address = "";
  let isLoading = false;
  let error = "";

  async function handleSubmit() {
    isLoading = true;
    error = "";

    try {
      await authService.register({
        email,
        password,
        first_name: firstName,
        last_name: lastName,
        phone_number: phoneNumber,
        date_of_birth: dateOfBirth,
        gender,
        address,
      });
      toast.success("Registration successful!");
      goto("/");
    } catch (err: any) {
      console.error("Registration failed:", err);
      error =
        err.response?.data?.error || "Registration failed. Please try again.";
      toast.error(error);
    } finally {
      isLoading = false;
    }
  }
</script>

<svelte:head>
  <title>Register | NyumbaniCare</title>
</svelte:head>

<div class="min-h-full flex flex-col justify-center py-12 sm:px-6 lg:px-8">
  <div class="sm:mx-auto sm:w-full sm:max-w-md">
    <h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900">
      Create your account
    </h2>
    <p class="mt-2 text-center text-sm text-gray-600">
      Or
      <a
        href="/login"
        class="font-medium text-primary-600 hover:text-primary-500"
      >
        sign in to your account
      </a>
    </p>
  </div>

  <div class="mt-8 sm:mx-auto sm:w-full sm:max-w-md">
    <div class="bg-white py-8 px-4 shadow sm:rounded-lg sm:px-10">
      <form class="space-y-6" on:submit|preventDefault={handleSubmit}>
        {#if error}
          <div class="rounded-md bg-red-50 p-4">
            <div class="flex">
              <div class="ml-3">
                <h3 class="text-sm font-medium text-red-800">{error}</h3>
              </div>
            </div>
          </div>
        {/if}

        <div class="grid grid-cols-1 gap-y-6 gap-x-4 sm:grid-cols-2">
          <div>
            <label
              for="first-name"
              class="block text-sm font-medium text-gray-700">First name</label
            >
            <div class="mt-1">
              <input
                id="first-name"
                name="first-name"
                type="text"
                autocomplete="given-name"
                required
                bind:value={firstName}
                class="input"
              />
            </div>
          </div>

          <div>
            <label
              for="last-name"
              class="block text-sm font-medium text-gray-700">Last name</label
            >
            <div class="mt-1">
              <input
                id="last-name"
                name="last-name"
                type="text"
                autocomplete="family-name"
                required
                bind:value={lastName}
                class="input"
              />
            </div>
          </div>
        </div>

        <div>
          <label for="email" class="block text-sm font-medium text-gray-700"
            >Email address</label
          >
          <div class="mt-1">
            <input
              id="email"
              name="email"
              type="email"
              autocomplete="email"
              required
              bind:value={email}
              class="input"
            />
          </div>
        </div>

        <div>
          <label for="password" class="block text-sm font-medium text-gray-700"
            >Password</label
          >
          <div class="mt-1">
            <input
              id="password"
              name="password"
              type="password"
              autocomplete="new-password"
              required
              minlength="6"
              bind:value={password}
              class="input"
            />
          </div>
        </div>

        <div>
          <label for="phone" class="block text-sm font-medium text-gray-700"
            >Phone number</label
          >
          <div class="mt-1">
            <input
              id="phone"
              name="phone"
              type="tel"
              autocomplete="tel"
              required
              bind:value={phoneNumber}
              class="input"
            />
          </div>
        </div>

        <div>
          <label for="dob" class="block text-sm font-medium text-gray-700"
            >Date of birth</label
          >
          <div class="mt-1">
            <input
              id="dob"
              name="dob"
              type="date"
              required
              bind:value={dateOfBirth}
              class="input"
            />
          </div>
        </div>

        <div>
          <label for="gender" class="block text-sm font-medium text-gray-700"
            >Gender</label
          >
          <div class="mt-1">
            <select
              id="gender"
              name="gender"
              required
              bind:value={gender}
              class="input"
            >
              <option value="">Select gender</option>
              <option value="male">Male</option>
              <option value="female">Female</option>
              <option value="other">Other</option>
              <option value="prefer_not_to_say">Prefer not to say</option>
            </select>
          </div>
        </div>

        <div>
          <label for="address" class="block text-sm font-medium text-gray-700"
            >Address</label
          >
          <div class="mt-1">
            <textarea
              id="address"
              name="address"
              rows="3"
              required
              bind:value={address}
              class="input"
            ></textarea>
          </div>
        </div>

        <div>
          <button
            type="submit"
            disabled={isLoading}
            class="w-full btn btn-primary py-2 px-4 {isLoading
              ? 'opacity-70 cursor-not-allowed'
              : ''}"
          >
            {isLoading ? "Creating account..." : "Create account"}
          </button>
        </div>
      </form>
    </div>
  </div>
</div>
