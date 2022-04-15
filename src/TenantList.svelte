<script>
	import { Link } from "svelte-routing";
    import { onMount } from "svelte";
    let tenants = [];

    onMount(async () => {
        const res = await fetch(`/api/tenant`);
        tenants = await res.json();
    });

    function notYetImplemented(evt) {
        alert("Feature is not implemented (yet).");
    }
</script>

<div>
    <table class="center">
        <thead>
            <tr>
                <th> Environment </th>
                <th> Tenant Name </th>
                <th> Data Owner ID </th>
                <th></th>
                <th></th>
            </tr>
        </thead>
        <tbody>
            {#each tenants as row}
                <tr>
                    <td>
                        {row.environment}
                    </td>
                    <td>
                        {row.tenant}
                    </td>
                    <td>
                        {row.deviceOwner}
                    </td>
                    <td>
                        <!-- svelte-ignore a11y-invalid-attribute -->
                        <a href="#" on:click="{notYetImplemented}">Delete</a>
                    </td>
                    <td>
                        <Link to="/tenant/{row.deviceOwner}/device">Devices</Link>
                    </td>
                </tr>
            {/each}
        </tbody>
    </table>
</div>

<style>
    .center {
        margin-left: auto;
        margin-right: auto;
    }

    table {
        min-width: 36em;
    }

    td {
        text-align: left;
        padding: 0.25em;
    }

    tr:nth-child(even) {
        background-color: #f2f2f2;
    }
</style>
