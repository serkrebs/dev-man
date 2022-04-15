<script>
    export let tenantId;

    import { onMount } from "svelte";
	import { Link } from "svelte-routing";
    let devices = [];

    onMount(async () => {
        const res = await fetch(`/api/tenant/${tenantId}/device`);
        devices = await res.json();
    });

    function goback() {
        window.history.back()
    }

    function notYetImplemented(evt) {
        alert("Feature is not implemented (yet).");
    }
</script>

<button on:click="{goback}">Back</button>
<p>Showing devices for {tenantId}</p>

<div>
    <table class="center">
        <thead>
            <tr>
                <th> Device ID </th>
                <th> Device Type </th>
                <th> Model No. </th>
                <th> Serial No. </th>
                <th> Tenant </th>
                <th></th>
                <th></th>
            </tr>
        </thead>
        <tbody>
            {#each devices as row}
                <tr>
                    <td>
                        {row.deviceId}
                    </td>
                    <td>
                        {row.deviceType}
                    </td>
                    <td>
                        {row.metadata.modelNumber}
                    </td>
                    <td>
                        {row.metadata.serialNumber}
                    </td>
                    <td>
                        {row.tenant}
                    </td>
                    <td>
                        <!-- svelte-ignore a11y-invalid-attribute -->
                        <a href="#" on:click="{notYetImplemented}">Delete</a>
                    </td>
                    <td>
                        <!-- svelte-ignore a11y-invalid-attribute -->
                        <Link to="/tenant/{tenantId}/device/{row.deviceId}">Edit</Link>
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

    td {
        text-align: left;
        padding: 0.25em;
    }

    tr:nth-child(even) {
        background-color: #f2f2f2;
    }
</style>
