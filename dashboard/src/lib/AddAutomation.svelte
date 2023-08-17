<script>
  import {
    Form,
    FormGroup,
    Select,
    SelectItem,
    Button,
    TextInput
  } from "carbon-components-svelte";

  export let deviceList
  export let nfcList

  const submitForm = async (e) => {
    e.preventDefault()

		const formData = new FormData(e.target)
		const data = new URLSearchParams()
		for (let field of formData) {
			const [key, value] = field
			data.append(key, value)
		}

    //post data
  }
</script>

<Form on:submit={submitForm}>
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
    <Select id="device" name="device" labelText="Device ID">
      <SelectItem
        disabled
        hidden
        value="placeholder-item"
        text="Choose a Device ID"
      />
      {#each deviceList as device}
        <SelectItem value="{device.id}" text="{device.name}" />
      {/each}
      <SelectItem value="test123" text="test123" />
    </Select>
  </FormGroup>
  <Button type="submit">Submit</Button>
</Form>