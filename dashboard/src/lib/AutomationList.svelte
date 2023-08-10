<script>
    import { onMount } from 'svelte';
  
    let automations = [];
  
    onMount(async () => {
      const res = await fetch(`/api/automation`);
      automations = await res.json();
    });
  </script>

<p>test</p>