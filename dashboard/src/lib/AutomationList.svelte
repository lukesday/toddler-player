<script>
    import { DataTable, Toolbar, ToolbarContent, Button } from "carbon-components-svelte";
    import { TrashCan, Add, Play } from "carbon-icons-svelte";
    import { showDeleteSuccess, showSuccess } from "$lib/toast.js";

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
      showDeleteSuccess('Automation', `"${automationList.find((automation) => automation.ID === id).Name}" deleted successfully`)
      automationList = automationList.filter((automation) => automation.ID !== id)
    }

    const handleTrigger = async (id) => {
      const response = await fetch(`/automation/trigger`, {
        method: "POST",
        mode: "no-cors",
        cache: "no-cache",
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ id: id })
      })

      if (response.status !== 200) {
        console.log('automation trigger failure', id, response.status)
        return
      }
      
      console.log('automation trigger success', id)
      showSuccess('Automation', "Automation", `"${automationList.find((automation) => automation.ID === id).Name}" triggered successfully`)
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
      <Button
        icon={TrashCan}
        kind="danger-ghost"
        on:click={handleDelete(cell.value)}
        iconDescription="Delete"
      />
      <Button
        icon={Play}
        kind="ghost"
        on:click={handleTrigger(cell.value)}
        iconDescription="Trigger"
      />
    {:else}
      {cell.value}
    {/if}
  </svelte:fragment>
  </DataTable>