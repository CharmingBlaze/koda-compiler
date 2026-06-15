<script>
  import { themes } from './themes.js'

  let {
    value = $bindable('koda-dark'),
  } = $props()
</script>

<label class="theme-picker" title="Theme">
  <span class="theme-picker-label">Theme</span>
  <select bind:value aria-label="Theme">
    {#each themes as theme (theme.id)}
      <option value={theme.id}>{theme.name}</option>
    {/each}
  </select>
</label>
