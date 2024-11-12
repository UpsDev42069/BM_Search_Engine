<script lang="ts">
  import { writable } from "svelte/store";
  import Layout from "../components/Layout.svelte";
  import { login } from "../utils/api";
  import { navigate } from "svelte-routing";
  import { session } from "../stores/session";
  let username = writable("");
  let password = writable("");

  async function handleLogin() {
    const res = await login($username, $password);
    if (res && res.ok) {
      await session.checkSession();
      const data = await res.json();
      if (data.resetPassword){
        navigate("/reset-password");
      } else {
        navigate("/");
      }
    } else {
      console.error("Login failed: ", res);
    }
  }
</script>

<Layout>
  <div class="flex justify-center items-center min-h-screen">
    <div class="w-11/12 md:w-2/3 lg:w-1/2 bg-white p-6 rounded shadow-md">
      <h2 class="text-xl mb-4 text-center">Log in</h2>
      <form on:submit|preventDefault={handleLogin} class="flex flex-col items-center">
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
        </dl>
        <div class="flex justify-center w-full">
          <input
            type="submit"
            value="Log in"
            class="py-2 px-20 border border-blue-600 rounded bg-blue-600 cursor-pointer ml-2"
          />
        </div>
      </form>
    </div>
  </div>
</Layout>
