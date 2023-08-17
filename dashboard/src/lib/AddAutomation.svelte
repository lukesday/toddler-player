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
</script>

<Form  method="POST" action="/">
  <FormGroup>
    <TextInput
      id="track"
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
    </Select>
  </FormGroup>
  <Button type="submit">Submit</Button>
</Form>