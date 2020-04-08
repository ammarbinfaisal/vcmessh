<script>
  import Router from "svelte-navaid/Router.svelte";
  import Route from "svelte-navaid/Route.svelte";
  import Link from "svelte-navaid/Link.svelte";
  import navaidHash from 'svelte-navaid/navaid-hash';
  import CreateRoom from "./components/CreateRoom.svelte";
  import Room from "./components/Room.svelte";
</script>

<style global>
  body {
    font-family: "Fira Code", monospace;
    font-weight: 500;
  }
  @tailwind base;
  @tailwind components;
  @tailwind utilities;
</style>

<main class="h-screen w-screen">
  <Router library={navaidHash}>
    <Route path="/">
      <CreateRoom />
    </Route>
    <Route path="/r/:roomname" let:params>
      <Room roomname={params.roomname}/>
    </Route>
    <Route>
      <h2>Page Not Found</h2>
    </Route>
  </Router>
</main>
