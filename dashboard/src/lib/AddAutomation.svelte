<script>
  import {
    Form,
    FormGroup,
    Select,
    SelectItem,
    Button,
    TextInput
  } from "carbon-components-svelte";
    import { ArrowRight } from "carbon-icons-svelte";

  export let deviceList
  export let nfcList

  const submitForm = async (e) => {
    e.preventDefault()

		const formData = new FormData(e.target)
		let data = {}
		for (let field of formData) {
			const [key, value] = field
			data[key] =  value
		}

    console.log(data)

    const response = await fetch(`/automation/add/submit`, 
    {
      method: "POST",
      mode: "no-cors",
      cache: "no-cache",
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(data)
    })
    
    if (response.status !== 200) {
      // Show failure message
    }
    
    // Show success message etc etc, redirect to home with table
    console.log('automation add success')
  }
</script>

<Form on:submit={submitForm}>
  <FormGroup>
    <TextInput
      id="name"
      name="name"
      invalidText="A valid name is required"
      labelText="Name"
      placeholder="eg My Automation"
    />
  </FormGroup>
  <FormGroup>
    <TextInput
      id="track"
      name="track"
      invalidText="A valid value is required"
      labelText="Track"
      placeholder="eg spotify:track:6rqhFgbbKwnb9MLmUQDhG6"
    />
  </FormGroup>
  <FormGroup>
    <Select id="nfc" name="nfc" labelText="NFC ID">
      <SelectItem
        disabled
        hidden
        value="placeholder-item"
        text="Choose an NFC ID"
      />
      {#each nfcList as nfc}
        <SelectItem value="{nfc.NfcUID}" text="{nfc.NfcUID}" />
      {/each}
    </Select>
  </FormGroup>
  <FormGroup>
    <Select id="device" name="device" labelText="Device">
      <SelectItem
        disabled
        hidden
        value="placeholder-item"
        text="Choose a Device"
      />
      {#each deviceList as device}
        <SelectItem value="{device.Id}" text="{device.Name}" />
      {/each}
      <SelectItem value="test123" text="test123" />
    </Select>
  </FormGroup>
  <Button type="submit" icon={ArrowRight}>Submit</Button>
  <Button kind='secondary' href="/">Back</Button>
</Form>