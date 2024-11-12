<script lang="ts">
  import { Router, Route } from "svelte-routing";
  import Search from "./pages/Search.svelte";
  import Register from "./pages/Register.svelte";
  import Login from "./pages/Login.svelte";
  import About from "./pages/About.svelte";
  import Weather from "./pages/Weather.svelte";
  import ResetPassword from "./pages/ResetPassword.svelte";
</script>

<Router>
  <Route path="/" component={Search} />
  <Route path="/register" component={Register} />
  <Route path="/login" component={Login} />
  <Route path="/about" component={About} />
  <Route path="/weather" component={Weather} />
  <Route path="/reset-password" component={ResetPassword} />
</Router>
