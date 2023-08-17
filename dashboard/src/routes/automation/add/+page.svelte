<script>
    import AddAutomation from "$lib/AddAutomation.svelte";
    
    export let data;
</script>

<AddAutomation deviceList={data.deviceList} nfcList={data.nfcList} />