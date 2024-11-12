<script lang="ts">
import Layout from "../components/Layout.svelte";
import { writable } from "svelte/store";
import { navigate } from "svelte-routing";
import { session } from "../stores/session";
import { resetPassword } from "../utils/api";



let username = $session.username;
let oldPassword = writable("");
let newPassword = writable("");
let confirmPassword = writable("");



async function handleResetPassword(){
    if (username) {

        if ($newPassword !== $confirmPassword) {
            console.error("Passwords do not match");
            return;
        }
        const res = await resetPassword(username, $oldPassword, $newPassword)
        if (res && res.ok) {
            await session.checkSession();
            navigate("/");
        } else {
            console.error("Reset password failed: ", res);
        }
    }
}
</script>





<Layout>
    <div class="flex justify-center items-center min-h-screen">
      <div class="w-11/12 md:w-2/3 lg:w-1/2 bg-white p-6 rounded shadow-md">
        <h2 class="text-xl mb-4 text-center">Log in</h2>
        <form on:submit|preventDefault={handleResetPassword} class="flex flex-col items-center">
          <dl class="w-full">
            <div class="flex flex-col items-center mb-4">
              <dt class="w-full md:w-10/12 text-left">Old Password:</dt>
              <dd class="w-full md:w-10/12">
                <input
                  type="password"
                  bind:value={$oldPassword}
                  size="30"
                  class="w-full p-2 m-2 border border-gray-300 rounded"
                />
              </dd>
            </div>
            <div class="flex flex-col items-center mb-4">
              <dt class="w-full md:w-10/12 text-left">New Password:</dt>
              <dd class="w-full md:w-10/12">
                <input
                  type="password"
                  bind:value={$newPassword}
                  size="30"
                  class="w-full p-2 m-2 border border-gray-300 rounded"
                />
              </dd>
            </div>
            <div class="flex flex-col items-center mb-4">
                <dt class="w-full md:w-10/12 text-left">New Password:</dt>
                <dd class="w-full md:w-10/12">
                  <input
                    type="password"
                    bind:value={$confirmPassword}
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
  

