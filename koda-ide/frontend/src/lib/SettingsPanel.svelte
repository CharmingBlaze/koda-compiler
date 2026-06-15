<script>
  import { themes } from './themes.js'

  let {
    theme = $bindable('koda-dark'),
    editorFontSize = $bindable(15),
    editorLineHeight = $bindable(1.65),
    editorTabSize = $bindable(4),
    terminalFontSize = $bindable(13),
    panelOpacity = $bindable(92),
    compactMode = $bindable(false),
    showHeaderThemePicker = $bindable(true),
    showTerminal = $bindable(false),
    showPreview = $bindable(false),
    outputOpen = $bindable(false),
    onReset = () => {},
  } = $props()
</script>

<div class="settings-panel">
  <div class="settings-heading">
    <div>
      <div class="settings-kicker">Settings</div>
      <h2>Customize Studio</h2>
    </div>
    <button type="button" class="settings-reset" onclick={() => onReset()}>Reset</button>
  </div>

  <section class="settings-section">
    <label class="settings-field">
      <span>Theme</span>
      <select bind:value={theme}>
        {#each themes as item (item.id)}
          <option value={item.id}>{item.name}</option>
        {/each}
      </select>
    </label>
    <label class="settings-check">
      <input type="checkbox" bind:checked={showHeaderThemePicker} />
      <span>Show theme picker in toolbar</span>
    </label>
  </section>

  <section class="settings-section">
    <div class="settings-section-title">Editor</div>
    <label class="settings-range">
      <span>Font size <strong>{editorFontSize}px</strong></span>
      <input type="range" min="12" max="22" step="1" bind:value={editorFontSize} />
    </label>
    <label class="settings-range">
      <span>Line height <strong>{editorLineHeight}</strong></span>
      <input type="range" min="1.35" max="2" step="0.05" bind:value={editorLineHeight} />
    </label>
    <label class="settings-field">
      <span>Tab width</span>
      <select bind:value={editorTabSize}>
        <option value={2}>2 spaces</option>
        <option value={4}>4 spaces</option>
        <option value={8}>8 spaces</option>
      </select>
    </label>
  </section>

  <section class="settings-section">
    <div class="settings-section-title">Workspace</div>
    <label class="settings-check">
      <input type="checkbox" bind:checked={showTerminal} />
      <span>Integrated terminal</span>
    </label>
    <label class="settings-check">
      <input type="checkbox" bind:checked={outputOpen} />
      <span>Output panel</span>
    </label>
    <label class="settings-check">
      <input type="checkbox" bind:checked={showPreview} />
      <span>3D preview pane</span>
    </label>
  </section>

  <section class="settings-section">
    <div class="settings-section-title">Interface</div>
    <label class="settings-check">
      <input type="checkbox" bind:checked={compactMode} />
      <span>Compact layout</span>
    </label>
    <label class="settings-range">
      <span>Panel opacity <strong>{panelOpacity}%</strong></span>
      <input type="range" min="76" max="100" step="1" bind:value={panelOpacity} />
    </label>
    <label class="settings-range">
      <span>Terminal size <strong>{terminalFontSize}px</strong></span>
      <input type="range" min="11" max="20" step="1" bind:value={terminalFontSize} />
    </label>
  </section>
</div>

<style>
  .settings-panel {
    display: flex;
    flex-direction: column;
    gap: 0.85rem;
    color: var(--color-text);
  }

  .settings-heading {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    gap: 0.75rem;
  }

  .settings-kicker,
  .settings-section-title {
    color: var(--color-overlay0);
    font-size: 0.6875rem;
    font-weight: 700;
    letter-spacing: 0.06em;
    text-transform: uppercase;
  }

  h2 {
    margin: 0.15rem 0 0;
    font-size: 1rem;
    font-weight: 650;
    color: var(--color-text);
  }

  .settings-reset {
    border: 1px solid var(--color-surface1);
    border-radius: 0.375rem;
    background: var(--color-surface0);
    color: var(--color-subtext);
    cursor: pointer;
    font: inherit;
    font-size: 0.75rem;
    padding: 0.25rem 0.45rem;
  }

  .settings-reset:hover {
    color: var(--color-text);
    border-color: var(--color-accent);
  }

  .settings-section {
    display: grid;
    gap: 0.65rem;
    border-top: 1px solid var(--color-surface0);
    padding-top: 0.8rem;
  }

  .settings-field,
  .settings-range,
  .settings-check {
    display: grid;
    gap: 0.35rem;
    font-size: 0.8125rem;
    color: var(--color-subtext);
  }

  .settings-range span {
    display: flex;
    justify-content: space-between;
    gap: 0.75rem;
  }

  strong {
    color: var(--color-cyan);
    font-weight: 650;
  }

  select,
  input[type='range'] {
    width: 100%;
  }

  select {
    border: 1px solid var(--color-surface1);
    border-radius: 0.375rem;
    background: var(--color-crust);
    color: var(--color-text);
    font: inherit;
    min-height: 2rem;
    padding: 0.3rem 0.45rem;
    outline: none;
  }

  select:focus {
    border-color: var(--color-accent);
  }

  input[type='range'] {
    accent-color: var(--color-accent);
  }

  .settings-check {
    grid-template-columns: 1rem 1fr;
    align-items: center;
    gap: 0.55rem;
  }

  .settings-check input {
    accent-color: var(--color-accent);
  }
</style>
