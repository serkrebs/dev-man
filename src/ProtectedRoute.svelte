<script>
    import { Route, Link } from 'svelte-routing';
    import Restricted from './Restricted.svelte';
    import {loggedInUser} from './login';

    export let path;
    export let component;

    $: isAuthenticated = !!$loggedInUser;
</script>

{#if isAuthenticated}
    <Route path={path} component={component} />
{:else}
    <Route path={path} component={Restricted} />
{/if}