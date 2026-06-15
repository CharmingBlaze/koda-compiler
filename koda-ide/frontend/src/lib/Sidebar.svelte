<script>
  import FileTree from './FileTree.svelte'
  import SettingsPanel from './SettingsPanel.svelte'

  let {
    zen = false,
    workspace = '',
    tree = [],
    activeRel = '',
    drawerPinned = $bindable(false),
    openSettingsToken = 0,
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
    onOpenFile = () => {},
    onOpenWorkspace = () => {},
    onRun = () => {},
    onOpenPalette = () => {},
    onOpenHelp = () => {},
    onResetSettings = () => {},
  } = $props()

  let railHover = $state(false)
  let activeTool = $state('files')

  const drawerOpen = $derived(drawerPinned || railHover)

  function setTool(t) {
    activeTool = t
    if (t === 'palette') {
      onOpenPalette()
      return
    }
    if (t === 'run') {
      onRun()
      return
    }
    if (t === 'help') {
      onOpenHelp('START_HERE.md')
    }
  }

  $effect(() => {
    if (!openSettingsToken) return
    activeTool = 'settings'
    drawerPinned = true
  })
</script>

{#if !zen}
  <div
    class="void-rail relative z-30 flex h-full min-h-0 w-11 shrink-0 flex-col items-center gap-1 border-r border-[var(--color-surface0)] bg-[var(--color-crust)] py-2"
    role="navigation"
    aria-label="Activity rail"
    onmouseenter={() => (railHover = true)}
    onmouseleave={() => (railHover = false)}
  >
    <button
      type="button"
      class="flex h-9 w-9 items-center justify-center rounded-md border text-[var(--color-subtext)] hover:bg-[var(--color-surface0)] hover:text-[var(--color-text)] {activeTool === 'files'
        ? 'border-[var(--color-accent)]'
        : 'border-transparent'}"
      title="Files"
      aria-label="Files"
      onclick={() => setTool('files')}
    >
      <svg class="h-5 w-5" fill="none" stroke="currentColor" stroke-width="1.6" viewBox="0 0 24 24"
        ><path stroke-linecap="round" stroke-linejoin="round" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" /></svg>
    </button>
    <button
      type="button"
      class="flex h-9 w-9 items-center justify-center rounded-md border text-[var(--color-subtext)] hover:bg-[var(--color-surface0)] hover:text-[var(--color-text)] {activeTool === 'search'
        ? 'border-[var(--color-accent)]'
        : 'border-transparent'}"
      title="Search (use ⌃P)"
      aria-label="Search"
      onclick={() => setTool('search')}
    >
      <svg class="h-5 w-5" fill="none" stroke="currentColor" stroke-width="1.6" viewBox="0 0 24 24"
        ><path stroke-linecap="round" d="M21 21l-4.3-4.3M11 18a7 7 0 100-14 7 7 0 000 14z" /></svg>
    </button>
    <button
      type="button"
      class="flex h-9 w-9 items-center justify-center rounded-md text-[var(--color-subtext)] hover:bg-[var(--color-surface0)] hover:text-[var(--color-text)]"
      title="Run"
      aria-label="Run"
      onclick={() => setTool('run')}
    >
      <svg class="h-5 w-5" fill="currentColor" viewBox="0 0 24 24"><path d="M8 5v14l11-7z" /></svg>
    </button>
    <button
      type="button"
      class="flex h-9 w-9 items-center justify-center rounded-md text-[var(--color-subtext)] hover:bg-[var(--color-surface0)] hover:text-[var(--color-text)]"
      title="Command palette"
      aria-label="Command palette"
      onclick={() => setTool('palette')}
    >
      <svg class="h-5 w-5" fill="none" stroke="currentColor" stroke-width="1.6" viewBox="0 0 24 24"
        ><path stroke-linecap="round" d="M4 6h16M4 12h10M4 18h16" /></svg>
    </button>
    <button
      type="button"
      class="flex h-9 w-9 items-center justify-center rounded-md border text-[var(--color-subtext)] hover:bg-[var(--color-surface0)] hover:text-[var(--color-text)] {activeTool === 'help'
        ? 'border-[var(--color-accent)]'
        : 'border-transparent'}"
      title="Help & documentation (F1)"
      aria-label="Help"
      onclick={() => setTool('help')}
    >
      <svg class="h-5 w-5" fill="none" stroke="currentColor" stroke-width="1.6" viewBox="0 0 24 24"
        ><circle cx="12" cy="12" r="9" /><path stroke-linecap="round" d="M9.5 9.5a2.5 2.5 0 115 0c0 2-2.5 2-2.5 4" /><circle cx="12" cy="17" r="0.5" fill="currentColor" /></svg>
    </button>
    <button
      type="button"
      class="mt-auto flex h-9 w-9 items-center justify-center rounded-md border text-[var(--color-subtext)] hover:bg-[var(--color-surface0)] hover:text-[var(--color-text)] {activeTool === 'settings'
        ? 'border-[var(--color-accent)]'
        : 'border-transparent'}"
      title="Settings"
      aria-label="Settings"
      onclick={() => setTool('settings')}
    >
      <svg class="h-5 w-5" fill="none" stroke="currentColor" stroke-width="1.6" viewBox="0 0 24 24"
        ><path
          stroke-linecap="round"
          stroke-linejoin="round"
          d="M10.3 3.3h3.4l1 2.8 2.8 1 2.1 2.1-1 2.8 1 2.8-2.1 2.1-2.8 1-1 2.8h-3.4l-1-2.8-2.8-1-2.1-2.1 1-2.8-1-2.8 2.1-2.1 2.8-1 1-2.8z"
        /><circle cx="12" cy="12" r="3" stroke-width="1.6" /></svg>
    </button>

    {#if drawerOpen}
      <div
        class="glass pointer-events-auto fixed bottom-[24px] left-11 top-8 z-40 w-72 overflow-y-auto border-r border-[var(--color-surface0)] p-3 text-sm shadow-2xl"
        role="region"
        aria-label="Side panel"
        onmouseenter={() => (railHover = true)}
        onmouseleave={() => (railHover = false)}
      >
        {#if activeTool === 'files'}
          <div class="mb-2 text-[11px] font-semibold uppercase tracking-wider text-[var(--color-overlay0)]">Files</div>
          {#if !workspace}
            <p class="text-sm text-[var(--color-subtext)]">Open a folder to browse .koda sources.</p>
            <button
              type="button"
              class="mt-2 w-full rounded-md btn-primary px-2 py-1.5 text-sm font-medium"
              onclick={() => onOpenWorkspace()}>Open workspace…</button>
          {:else}
            <FileTree entries={tree} {activeRel} {onOpenFile} />
          {/if}
        {:else if activeTool === 'search'}
          <p class="text-sm leading-relaxed text-[var(--color-subtext)]">Fuzzy file &amp; command search lives in the palette — press <kbd class="rounded bg-[var(--color-surface0)] px-1">⌃P</kbd>.</p>
        {:else if activeTool === 'help'}
          <div class="help-nav-label" style="font-size:0.6875rem;font-weight:600;letter-spacing:0.06em;text-transform:uppercase;color:var(--color-overlay0);margin-bottom:0.5rem;">Help</div>
          <p class="mb-2 text-sm text-[var(--color-subtext)]">Full SDK docs, tutorials, and language reference — press <kbd class="rounded bg-[var(--color-surface0)] px-1">F1</kbd>.</p>
          <div class="flex flex-col gap-1">
            <button type="button" class="help-quick-btn rounded-md px-2 py-1.5 text-left text-sm text-[var(--color-cyan)] hover:bg-[var(--color-surface0)]" onclick={() => onOpenHelp('START_HERE.md')}>Start here</button>
            <button type="button" class="rounded-md px-2 py-1.5 text-left text-sm text-[var(--color-subtext)] hover:bg-[var(--color-surface0)]" onclick={() => onOpenHelp('docs/beginners-guide.md')}>Beginner's guide</button>
            <button type="button" class="rounded-md px-2 py-1.5 text-left text-sm text-[var(--color-subtext)] hover:bg-[var(--color-surface0)]" onclick={() => onOpenHelp('docs/learn/README.md')}>Learn path</button>
            <button type="button" class="rounded-md px-2 py-1.5 text-left text-sm text-[var(--color-subtext)] hover:bg-[var(--color-surface0)]" onclick={() => onOpenHelp('docs/faq.md')}>FAQ</button>
            <button type="button" class="rounded-md px-2 py-1.5 text-left text-sm text-[var(--color-subtext)] hover:bg-[var(--color-surface0)]" onclick={() => onOpenHelp('docs/troubleshooting.md')}>Troubleshooting</button>
            <button type="button" class="rounded-md px-2 py-1.5 text-left text-sm text-[var(--color-subtext)] hover:bg-[var(--color-surface0)]" onclick={() => onOpenHelp('language.md')}>Language reference</button>
          </div>
        {:else if activeTool === 'settings'}
          <SettingsPanel
            bind:theme
            bind:editorFontSize
            bind:editorLineHeight
            bind:editorTabSize
            bind:terminalFontSize
            bind:panelOpacity
            bind:compactMode
            bind:showHeaderThemePicker
            bind:showTerminal
            bind:showPreview
            bind:outputOpen
            onReset={onResetSettings}
          />
        {/if}
      </div>
    {/if}
  </div>
{/if}
