<script lang="ts">
  import { Router, Link, Route } from "svelte-routing";
  import Logo from "./Logo.svelte";
  import TenantList from "./TenantList.svelte";
  import DeviceList from "./DeviceList.svelte";
  import Device from "./Device.svelte";
  import MissingFeatures from "./MissingFeatures.svelte";
  import Login from "./Login.svelte";
  import ProtectedRoute from "./ProtectedRoute.svelte";

  export let url = "";
</script>

<main>
  <Router {url}>
    <header id="top-navbar">
      <div>
        <Logo />
        <div id="logo">
          <h3 class="no-vert-margin">Development Device Manager</h3>
        </div>

        <nav>
          <Link to="/">Tenants</Link>
          <Link to="/missing">Work Remaining</Link>
        </nav>
        <Login />
      </div>
    </header>
    <div>
      <ProtectedRoute path="tenant/:tenantId/device/:deviceId" component={Device} />
      <ProtectedRoute path="tenant/:tenantId/device" component={DeviceList} />
      <Route path="missing" component={MissingFeatures} />
      <ProtectedRoute path="/" component={TenantList} />
    </div>
  </Router>
</main>

<style>
  main {
    text-align: center;
    padding: 0;
    margin: 0 auto;
  }

  header > div {
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    height: 42px;
    width: 100%;
    background-color: #e7e7e7;
  }

  #logo {
    padding-top: 5px;
  }

  nav {
    padding-top: 8px;
	padding-right: 2em;
  }

  .no-vert-margin {
    margin-block-start: 0;
    margin-block-end: 0;
  }

  :global(nav a) {
    padding-left: 30px;
  }
</style>
