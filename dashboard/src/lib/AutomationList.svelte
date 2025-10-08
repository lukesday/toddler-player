<script>
    import { DataTable, Link, Toolbar, ToolbarContent, Button } from "carbon-components-svelte";
    import { TrashCan, Add } from "carbon-icons-svelte";
    export let automationList

    const handleDelete = async (id) => {

      const response = await fetch(`/automation/delete`, 
      {
        method: "POST",
        mode: "no-cors",
        cache: "no-cache",
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ id: id })
      })
      
      if (response.status !== 200) {
        // Show failure message
        console.log('automation delete failure', id)
        return
      }
      
      // Show success message, remove from list
      console.log('automation delete success', id)
    }
</script>
  
<DataTable
    headers={[
        { key: "Name", value: "Name" },
        { key: "NfcUID", value: "NFC UID" },
        { key: "MediaId", value: "Track" },
        { key: "Action", value: "Actions"},
    ]}
    rows={automationList.map((automation) => ({
        Name: automation.Name,
        NfcUID: automation.NfcTag.NfcUID,
        MediaId: automation.MediaId,
        Action: automation.ID,
    }))}
  title="Automations"
  description="All active automations">

  <Toolbar>
    <ToolbarContent>
        <Button 
        href="/automation/add"
        icon={Add}>
        New Automation
        </Button>
    </ToolbarContent>
  </Toolbar>
  <svelte:fragment slot="cell" let:row let:cell>
    {#if cell.key === "Action"}
      <Link
        icon={TrashCan}
        on:click={handleDelete(cell.value)}
        target="_blank">Delete</Link>
    {:else}
      {cell.value}
    {/if}
  </svelte:fragment>
  </DataTable>