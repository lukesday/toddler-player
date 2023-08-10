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
    export let data;

    let isSideNavOpen = false
</script>


<Header platformName="Toddler Player" bind:isSideNavOpen>
    <HeaderNav>
        {#if !data.loggedIn}
            <HeaderNavItem href="/login" text="Login" />
        {:else}
            <HeaderNavItem href="/logout" text="Logout" />
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