<script>
// @ts-nocheck

    export let data;
</script>

<h1>Welcome to the Toddler Player Dashboard</h1>
{#if data.loggedIn}
    <p>Welcome back {data.userData.display_name}</p>
{:else}
    <p><a href="/login">Login</a></p>
{/if}
<slot />