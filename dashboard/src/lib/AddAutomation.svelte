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
    import { goto } from "$app/navigation";

  export let nfcList

  // Form data
  let name = '';
  let track = '';
  let nfc = 'placeholder-item';

  // Track which fields have been touched (blurred)
  let touched = {};

  // Validation errors
  let errors = {};

  // Reactive validation - only show errors for touched fields
  $: errors.name = touched.name && name.trim().length < 2 ? 'Name must be at least 2 characters' : '';
  $: errors.track = touched.track && track.trim().length < 1 ? 'Track ID is required' : '';
  $: errors.nfc = touched.nfc && (!nfc || nfc === '' || nfc === 'placeholder-item') ? 'Please select an NFC ID' : '';

  // Check if form is valid (based on actual values, not just errors)
  $: isFormValid = name.trim().length >= 2 && 
                   track.trim().length >= 1 && 
                   nfc && nfc !== '' && nfc !== 'placeholder-item';

  // Debug logging
  $: console.log('NFC value:', nfc, 'Type:', typeof nfc, 'Is valid:', nfc && nfc !== '' && nfc !== 'placeholder-item');

  // Function to mark field as touched
  function markTouched(field) {
    touched[field] = true;
  }

  const submitForm = async (e) => {
    e.preventDefault()

    if (!isFormValid) {
      console.log('Form is not valid');
      return;
    }

    const data = {
      name: name.trim(),
      track: track.trim(),
      nfc: nfc
    };

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
      console.log('automation add failure')
      return
    }
    
    // Show success message etc etc, redirect to home with table
    console.log('automation add success')
    goto('/')
  }
</script>

<Form on:submit={submitForm}>
  <FormGroup>
    <TextInput
      id="name"
      name="name"
      bind:value={name}
      on:blur={() => markTouched('name')}
      invalid={errors.name !== ''}
      invalidText={errors.name || "A valid name is required"}
      labelText="Name"
      placeholder="eg My Automation"
    />
  </FormGroup>
  <FormGroup>
    <TextInput
      id="track"
      name="track"
      bind:value={track}
      on:blur={() => markTouched('track')}
      invalid={errors.track !== ''}
      invalidText={errors.track || "A valid value is required"}
      labelText="Track"
      placeholder="eg 6rqhFgbbKwnb9MLmUQDhG6"
    />
  </FormGroup>
  <FormGroup>
    <Select 
      id="nfc" 
      name="nfc" 
      labelText="NFC ID"
      bind:selected={nfc}
      on:blur={() => markTouched('nfc')}
      invalid={errors.nfc !== ''}
      invalidText={errors.nfc || "Please select an NFC ID"}
    >
      <SelectItem
        disabled
        hidden
        value="placeholder-item"
        text="Choose an NFC ID"
        selected
      />
      {#each nfcList as nfcItem}
        <SelectItem value="{nfcItem.NfcUID}" text="{nfcItem.NfcUID}" />
      {/each}
    </Select>
  </FormGroup>
  <Button type="submit" icon={ArrowRight} disabled={!isFormValid}>Submit</Button>
  <Button kind='secondary' href="/">Back</Button>
</Form>