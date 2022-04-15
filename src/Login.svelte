<script>
    import {loggedInUser, azureLogin, azureLogout} from './login';
</script>

{#if !!$loggedInUser}
    <span>Welcome {$loggedInUser.name}</span>
    <button on:click={azureLogout}>Log out</button>
{:else}
    <button on:click={azureLogin}>Log in</button>
{/if}
