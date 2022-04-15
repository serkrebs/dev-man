<script>
    export let deviceId = "foo";
    export let tenantId = "guid";
    import { onMount } from "svelte";

    let device = {
        metadata: {}
    };
    $: extraValueKeys = Object.keys(device.metadata).filter(k => !commonKeys.includes(k));

    let tenant = {
        deviceOwner: tenantId,
        tenant: "DEVTEAM",
        environment: "devteam"
    };
    

    onMount(async () => {
        const tenRes = await fetch(`/api/tenant/${tenantId}`);
        tenant = await tenRes.json();
        const devRes = await fetch(`/api/tenant/${tenantId}/device/${deviceId}`);
        device = await devRes.json();

        console.dir(device);
    });

    let commonKeys = [ "environment", "modelNumber", "displayName", "description", "serialNumber", "manufacturer", 
        "hardwareVersion", "firmwareVersion", "softwareVersion", "dataOwner", "newDataOwner" ];

    let commonLabels = {
        environment: "Tenant Environment",
        modelNumber: "Model Number",
        displayName: "Display Name", 
        description: "Description", 
        serialNumber: "Serial Number", 
        manufacturer: "Manufacturer", 
        hardwareVersion: "Hardware Version", 
        firmwareVersion: "Firmware Version", 
        softwareVersion: "Software Version",
        dataOwner: "Data Owner", 
        newDataOwner: "New Data Owner"
    };

    let addingRow = false;
    let extraKey = null;
    let extraKeyInput;

    function deleteExtraValue(k) {
        delete device.metadata[k];
        device = device;        // force update of extraValueKeys
    }

    function addExtraValue() {
        if (extraKey) {
            device.metadata[extraKey] = null;
        }
        addingRow = false;
        extraKey = null;
    }

    function addExtraKey() {
        addingRow = true;
    }

    function cancelExtraValue() {
        addingRow = false;
        extraKey = null;
    }

    async function updateDevice() {
        const devRes = await fetch(`/api/tenant/${tenantId}/device/${deviceId}`, {
            method: 'PUT',
            mode: 'same-origin',
            cache: 'no-cache',
            credentials: 'same-origin',
            headers: {
            'Content-Type': 'application/json'
            },
            redirect: 'follow',
            referrerPolicy: 'no-referrer',
            body: JSON.stringify(device)
        });
        device = await devRes.json();
        console.dir(device)
    }
</script>

<h4>Device: {deviceId}</h4>
<p>
    Tenant: {tenant.tenant}<br />
    Device Owner: {tenant.deviceOwner}
</p>

<table class="center">
    <tr>
        <th>Device Type: </th>
        <td><input class="data-entry"  bind:value={device.deviceType} /></td>
        <td></td>
    </tr>
    {#each commonKeys as key}
        <tr>
            <th>{commonLabels[key]}: </th>
            <td><input class="data-entry"  bind:value={device.metadata[key]} /></td>
            <td></td>
        </tr>
    {/each}
    {#each extraValueKeys as k}
        <tr>
            <th>{k}: </th>
            <td><input class="data-entry"  bind:value={device.metadata[k]} /></td>
            <td>
                <button on:click={deleteExtraValue(k)}>-</button>
            </td>
        </tr>
    {/each}
    {#if addingRow}
        <tr>
            <th><input class="data-key" bind:this={extraKeyInput} bind:value={extraKey} on:blur={addExtraValue} on:show={() => alert("now")}/></th>
            <td><input class="data-entry" disabled=disabled /></td>
            <td>
                <button on:click={cancelExtraValue}>-</button>
            </td>
        </tr>
    {:else}
        <tr>
            <th></th>
            <td></td>
            <td>
                <button on:click={addExtraKey}>+</button>
            </td>
        </tr>
    {/if}
    <tr></tr>
</table>
<div>
    <button on:click={updateDevice}>Update</button>
</div>

<style>
    .center {
        margin-left: auto;
        margin-right: auto;
    }

    input.data-entry {
        min-width: 20rem;
    }
</style>
