<script>
  let {
    zen = false,
    hasWorkspace = false,
    onNewProject = () => {},
    onOpenWorkspace = () => {},
    onNewFile = () => {},
    onSave = () => {},
    onCloseTab = () => {},
    canSave = false,
    canCloseTab = false,
  } = $props()

  let menuOpen = $state(false)
  let wrapEl = $state(null)

  function closeMenu() {
    menuOpen = false
  }

  /** @param {MouseEvent} e */
  function onDocClick(e) {
    if (!menuOpen || !wrapEl) return
    if (e.target instanceof Node && !wrapEl.contains(e.target)) closeMenu()
  }
</script>

<svelte:window onclick={onDocClick} />

{#if !zen}
  <div class="relative flex items-center" bind:this={wrapEl}>
    <button
      type="button"
      class="rounded-md px-2.5 py-0.5 text-[var(--color-subtext)] hover:bg-[var(--color-surface0)] hover:text-[var(--color-text)]"
      class:bg-[var(--color-surface0)]={menuOpen}
      aria-expanded={menuOpen}
      aria-haspopup="true"
      onclick={() => (menuOpen = !menuOpen)}>File</button>
    {#if menuOpen}
      <div
        class="glass absolute left-0 top-full z-50 mt-0.5 min-w-[14rem] overflow-hidden rounded-lg border border-[var(--color-surface0)] py-1 text-sm shadow-xl"
        role="menu"
      >
        <button
          type="button"
          class="flex w-full items-center justify-between gap-6 px-3 py-1.5 text-left text-[var(--color-text)] hover:bg-[var(--color-surface0)]"
          role="menuitem"
          onclick={() => {
            closeMenu()
            onNewProject()
          }}>New project…</button>
        <button
          type="button"
          class="flex w-full items-center justify-between gap-6 px-3 py-1.5 text-left text-[var(--color-text)] hover:bg-[var(--color-surface0)]"
          role="menuitem"
          onclick={() => {
            closeMenu()
            onOpenWorkspace()
          }}>Open workspace…</button>
        <div class="my-1 h-px bg-[var(--color-surface0)]"></div>
        <button
          type="button"
          class="flex w-full items-center justify-between gap-6 px-3 py-1.5 text-left text-[var(--color-text)] hover:bg-[var(--color-surface0)] disabled:opacity-40"
          role="menuitem"
          disabled={!hasWorkspace}
          title={!hasWorkspace ? 'Open or create a workspace first' : ''}
          onclick={() => {
            if (!hasWorkspace) return
            closeMenu()
            onNewFile()
          }}>New file…</button>
        <button
          type="button"
          class="flex w-full items-center justify-between gap-6 px-3 py-1.5 text-left text-[var(--color-text)] hover:bg-[var(--color-surface0)] disabled:opacity-40"
          role="menuitem"
          disabled={!canSave}
          onclick={() => {
            if (!canSave) return
            closeMenu()
            onSave()
          }}
        >
          <span>Save</span>
          <kbd class="font-sans text-[10px] text-[var(--color-overlay0)]">⌃S</kbd>
        </button>
        <div class="my-1 h-px bg-[var(--color-surface0)]"></div>
        <button
          type="button"
          class="flex w-full items-center justify-between gap-6 px-3 py-1.5 text-left text-[var(--color-text)] hover:bg-[var(--color-surface0)] disabled:opacity-40"
          role="menuitem"
          disabled={!canCloseTab}
          onclick={() => {
            if (!canCloseTab) return
            closeMenu()
            onCloseTab()
          }}>Close tab</button>
      </div>
    {/if}
  </div>
{/if}
