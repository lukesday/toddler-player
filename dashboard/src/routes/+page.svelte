<script>
// @ts-nocheck
	import AutomationList from "$lib/AutomationList.svelte";
    import { getContext } from 'svelte';

    let layoutData = getContext('data')

    export let data;
</script>

<h2>{data.subTitle}</h2>
{#if layoutData.loggedIn}
<AutomationList />
{/if}
