<script>
// @ts-nocheck
	import AutomationList from "$lib/AutomationList.svelte";
    import { getContext } from 'svelte';
    import {Grid, Row, Column } from "carbon-components-svelte";

    let layoutData = getContext('data')

    export let data;
</script>


    <Row padding>
        <Column>
        <h2>{data.subTitle}</h2>
        </Column>
    </Row>
    <Row>
        <Column>
        {#if layoutData.loggedIn}
        <AutomationList automationList={data.automationList} />
        {/if}
        </Column>
    </Row>