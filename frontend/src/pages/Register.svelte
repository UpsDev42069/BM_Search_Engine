<script lang="ts">
  import Layout from "../components/Layout.svelte";
  import { writable } from "svelte/store";
  import { login, register } from "../utils/api";
  import { navigate } from "svelte-routing";

  let error = writable<string | null>(null);
  let username = writable("");
  let email = writable("");
  let password = writable("");
  let password2 = writable("");

  async function handleSubmit() {
    if ($password !== $password2) {
      error.set("Passwords do not match");
      return;
    }
    const registerResponse = await register($username, $password, $email);
    console.log(registerResponse);
    if (registerResponse && registerResponse.ok) {
      // TODO: login users after registration
      error.set(null);
      const loginResponse = await login($username, $password);
      if (loginResponse && loginResponse.ok) {
        navigate("/");
      } else {
        error.set("Login failed");
      }
    } else {
      error.set("Registration failed");
    }
  }
</script>

<Layout>
  <div class="flex justify-center items-center min-h-screen">
    <div class="w-11/12 md:w-2/3 lg:w-1/2 bg-white p-6 rounded shadow-md">
      <h2 class="text-xl mb-4 text-center">Sign up</h2>
      {#if $error}
        <div class="my-2 bg-red-100 border border-red-400 rounded p-2 text-xs"><strong>Error:</strong> {$error}</div>
      {/if}
      <form on:submit|preventDefault={handleSubmit} class="flex flex-col items-center">
        <dl class="w-full">
          <div class="flex flex-col items-center mb-4">
            <dt class="w-full md:w-10/12 text-left">Username:</dt>
            <dd class="w-full md:w-10/12">
              <input
                type="text"
                bind:value={$username}
                size="30"
                class="w-full p-2 m-2 border border-gray-300 rounded"
              />
            </dd>
          </div>
          <div class="flex flex-col items-center mb-4">
            <dt class="w-full md:w-10/12 text-left">E-mail:</dt>
            <dd class="w-full md:w-10/12">
              <input type="text" bind:value={$email} size="30" class="w-full p-2 m-2 border border-gray-300 rounded" />
            </dd>
          </div>
          <div class="flex flex-col items-center mb-4">
            <dt class="w-full md:w-10/12 text-left">Password:</dt>
            <dd class="w-full md:w-10/12">
              <input
                type="password"
                bind:value={$password}
                size="30"
                class="w-full p-2 m-2 border border-gray-300 rounded"
              />
            </dd>
          </div>
          <div class="flex flex-col items-center mb-4">
            <dt class="w-full md:w-10/12 text-left">Confirm password:</dt>
            <dd class="w-full md:w-10/12">
              <input
                type="password"
                bind:value={$password2}
                size="30"
                class="w-full p-2 m-2 border border-gray-300 rounded"
              />
            </dd>
          </div>
        </dl>
        <div class="flex justify-center w-full">
          <input
            type="submit"
            value="Sign up"
            class="py-2 px-20 border border-blue-600 rounded bg-blue-600 cursor-pointer ml-2"
          />
        </div>
      </form>
    </div>
  </div>
</Layout>
