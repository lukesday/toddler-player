<script>
// @ts-nocheck
    import "carbon-components-svelte/css/g80.css";
    import {
        Header,
        HeaderNav,
        HeaderNavItem,
        Content,
        Grid,
        Row,
        Column,
    } from "carbon-components-svelte";
    import { setContext } from 'svelte';
    import Toast from '../lib/Toast.svelte';

    export let data;

    setContext('data', data)
    
    let isSideNavOpen = false
</script>

<Header platformName="Toddler Player" bind:isSideNavOpen>
    <HeaderNav>
        {#if !data.loggedIn}
            <HeaderNavItem href="/login" text="Login" />
        {:else}
            <HeaderNavItem href="/" text="Home" />
            <HeaderNavItem data-sveltekit-preload-data="off" href="/logout" text="Logout" />
        {/if}
    </HeaderNav>
  </Header>
  
  <Content>
    <Grid>
      <Row padding>
        <Column>
            <h1>Dashboard</h1>
            {#if data.loggedIn}
                <p>Welcome back {data.userData.display_name}</p>
            {/if}
        </Column>
      </Row>
      <Row>
        <Column>
            <slot />
        </Column>
      </Row>
    </Grid>
  </Content>

  <!-- Toast notifications -->
  <Toast />