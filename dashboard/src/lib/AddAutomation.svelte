<script>
  import {
    Form,
    FormGroup,
    Checkbox,
    RadioButtonGroup,
    RadioButton,
    Select,
    SelectItem,
    Button,
  } from "carbon-components-svelte";

  export let deviceList
  export let nfcList
</script>

<Form  method="POST" action="/api/automation">
  <input name="track" id="track" type="text">
  <FormGroup>
    <Select id="nfc" name="nfc" labelText="NFC ID">
      <SelectItem
        disabled
        hidden
        value="placeholder-item"
        text="Choose an NFC ID"
      />
      <SelectItem value="option-1" text="Option 1" />
      {#each nfcList as nfc}
        <SelectItem value="{nfc.nfc_uid}" text="{nfc.nfc_uid}" />
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