<script>
  import { onMount, onDestroy, tick } from 'svelte'
  import {
    PickWorkspaceFolder,
    PickParentFolderForNewProject,
    CreateProjectInParent,
    OpenWorkspace,
    GetWorkspaceRoot,
    ListDir,
    ReadFile,
    WriteFile,
    AbsFromWorkspace,
    RunProgram,
    BuildProgram,
    DefaultBuildOutput,
    LSPMessage,
  } from '../wailsjs/go/main/App.js'
  import { EventsOn } from '../wailsjs/runtime/runtime.js'
  import CodeEditor from './lib/CodeEditor.svelte'
  import CommandPalette from './lib/CommandPalette.svelte'
  import TerminalPane from './lib/TerminalPane.svelte'
  import PreviewViewport from './lib/PreviewViewport.svelte'
  import Sidebar from './lib/Sidebar.svelte'
  import StatusBar from './lib/StatusBar.svelte'
  import DiagnosticToast from './lib/DiagnosticToast.svelte'
  import FileMenuBar from './lib/FileMenuBar.svelte'
  import { pathToFileURI, rpcInitialize, notifyInitialized, notifyDidOpen, notifyDidChange } from './lib/lspUtil.js'

  let workspace = $state('')
  let tree = $state([])
  let tabs = $state([])
  let activeIndex = $state(-1)
  let zen = $state(false)
  let outputOpen = $state(false)
  let outputText = $state('')
  let paletteOpen = $state(false)
  let showTerminal = $state(false)
  let showPreview = $state(false)
  let themeNudge = $state(0)
  let sidebarPinned = $state(false)
  let activeAbsPath = $state('')
  let diagErrors = $state(0)
  let diagWarnings = $state(0)
  /** @type { { line: number, col: number, token: number } | null } */
  let jumpTarget = $state(null)
  let newProjectOpen = $state(false)
  let newProjectParent = $state('')
  let newProjectName = $state('my-koda-project')
  let newProjectErr = $state('')
  let newFileOpen = $state(false)
  let newFileRel = $state('hello.koda')
  let newFileErr = $state('')
  let unsubOut = []
  let lspTimer = 0

  const active = $derived(activeIndex >= 0 ? tabs[activeIndex] : null)

  function norm(p) {
    return (p || '').replace(/\\/g, '/').toLowerCase()
  }

  function absToRel(abs) {
    const root = workspace.replace(/\\/g, '/').replace(/\/$/, '')
    const a = abs.replace(/\\/g, '/')
    if (!root) return ''
    const prefix = root.toLowerCase() + '/'
    if (!a.toLowerCase().startsWith(prefix)) return ''
    return a.slice(root.length + 1).replace(/\\/g, '/')
  }

  $effect(() => {
    const a = active
    if (!a) {
      activeAbsPath = ''
      return
    }
    void AbsFromWorkspace(a.rel)
      .then((p) => {
        activeAbsPath = p
      })
      .catch(() => {
        activeAbsPath = ''
      })
  })

  /** @param {any} payload */
  function updateDiagCounts(payload) {
    if (!payload?.path || !activeAbsPath) {
      diagErrors = 0
      diagWarnings = 0
      return
    }
    if (norm(payload.path) !== norm(activeAbsPath)) {
      diagErrors = 0
      diagWarnings = 0
      return
    }
    const raw = payload.diagnosticsRaw
    if (!Array.isArray(raw)) {
      diagErrors = 0
      diagWarnings = 0
      return
    }
    let e = 0
    let w = 0
    for (const d of raw) {
      const s = (d.severity || 'error').toLowerCase()
      if (s === 'warning') w++
      else e++
    }
    diagErrors = e
    diagWarnings = w
  }

  function wailsAppBinding() {
    if (typeof window === 'undefined') return null
    return window['go']?.['main']?.['App'] ?? null
  }

  async function refreshTree() {
    if (!workspace) {
      tree = []
      return
    }
    try {
      tree = await ListDir('')
    } catch {
      tree = []
    }
  }

  async function openWorkspaceFlow() {
    const app = wailsAppBinding()
    if (!app || typeof app.PickWorkspaceFolder !== 'function') {
      outputOpen = true
      outputText =
        'Open workspace needs the desktop app (Wails). The browser preview has no folder picker.\n\nRun: wails dev\nfrom the koda-ide folder.'
      return
    }
    let picked = ''
    try {
      picked = await PickWorkspaceFolder()
    } catch (e) {
      outputOpen = true
      outputText = e instanceof Error ? e.message : String(e)
      return
    }
    if (!picked) return
    try {
      await OpenWorkspace(picked)
      workspace = await GetWorkspaceRoot()
      tabs = []
      activeIndex = -1
      outputText = ''
      await refreshTree()
    } catch (e) {
      outputOpen = true
      outputText = e instanceof Error ? e.message : String(e)
    }
  }

  async function startNewProjectWizard() {
    newProjectErr = ''
    const app = wailsAppBinding()
    if (!app || typeof app.PickParentFolderForNewProject !== 'function') {
      newProjectParent = ''
      newProjectErr =
        'New project needs the desktop app: the browser preview (npm run dev) does not load the Go runtime, so the folder dialog cannot run.\n\nFrom koda-ide run: wails dev'
      newProjectOpen = true
      return
    }
    let parent = ''
    try {
      parent = await PickParentFolderForNewProject()
    } catch (e) {
      newProjectParent = ''
      newProjectErr = e instanceof Error ? e.message : String(e)
      newProjectOpen = true
      return
    }
    if (!parent) return
    newProjectParent = parent
    newProjectName = 'my-koda-project'
    newProjectOpen = true
  }

  async function confirmNewProject() {
    newProjectErr = ''
    if (!newProjectParent.trim()) return
    try {
      const root = await CreateProjectInParent(newProjectParent, newProjectName.trim())
      await OpenWorkspace(root)
      workspace = await GetWorkspaceRoot()
      tabs = []
      activeIndex = -1
      outputText = ''
      await refreshTree()
      newProjectOpen = false
      await openRel('main.koda')
    } catch (err) {
      newProjectErr = err instanceof Error ? err.message : String(err)
    }
  }

  function openNewFileDialog() {
    newFileErr = ''
    newFileRel = 'hello.koda'
    newFileOpen = true
  }

  async function confirmNewFile() {
    newFileErr = ''
    let rel = newFileRel.trim().replace(/\\/g, '/').replace(/^\/+/, '')
    if (!rel) {
      newFileErr = 'Enter a file path.'
      return
    }
    if (rel.includes('..')) {
      newFileErr = 'Path cannot contain "..".'
      return
    }
    if (!rel.endsWith('.koda')) rel = `${rel}.koda`
    const starter = `// ${rel}\n\nprint("hello");\n`
    try {
      await WriteFile(rel, starter)
      newFileOpen = false
      await refreshTree()
      await openRel(rel)
    } catch (err) {
      newFileErr = err instanceof Error ? err.message : String(err)
    }
  }

  function closeActiveTab() {
    if (activeIndex < 0 || tabs.length === 0) return
    jumpTarget = null
    const next = tabs.filter((_, i) => i !== activeIndex)
    const nextIndex = next.length === 0 ? -1 : activeIndex >= next.length ? next.length - 1 : activeIndex
    tabs = next
    activeIndex = nextIndex
  }

  async function openRel(rel) {
    jumpTarget = null
    const existing = tabs.findIndex((t) => t.rel === rel)
    if (existing >= 0) {
      activeIndex = existing
      return
    }
    const text = await ReadFile(rel)
    tabs = [...tabs, { rel, name: rel.split('/').pop() || rel, text }]
    activeIndex = tabs.length - 1
    const abs = await AbsFromWorkspace(rel)
    void LSPMessage(notifyDidOpen(pathToFileURI(abs), text))
  }

  async function saveActive() {
    if (!active) return
    await WriteFile(active.rel, active.text)
    outputText = outputText + `Saved ${active.rel}\n`
  }

  async function runActive() {
    if (!active) return
    const abs = await AbsFromWorkspace(active.rel)
    RunProgram(abs, active.text)
  }

  async function buildActive() {
    if (!active) return
    const abs = await AbsFromWorkspace(active.rel)
    const out = await DefaultBuildOutput(abs)
    BuildProgram(abs, active.text, out)
  }

  function pushLsp(path, text) {
    if (!path) return
    clearTimeout(lspTimer)
    lspTimer = window.setTimeout(() => {
      const uri = pathToFileURI(path)
      void LSPMessage(notifyDidChange(uri, text))
    }, 400)
  }

  async function tryAmbientContrast() {
    try {
      // @ts-ignore — optional in WebView2
      const Sensor = window.AmbientLightSensor
      if (!Sensor) return
      const s = new Sensor()
      s.addEventListener('reading', () => {
        const lux = s.illuminance
        themeNudge = lux > 50 ? 1 : 0
      })
      s.start()
    } catch {
      /* ignore */
    }
  }

  async function collectPaletteItems() {
    const cmds = [
      { id: 'c:open', label: 'Open workspace folder…', hint: 'workspace', shortcut: '' },
      { id: 'c:new-project', label: 'New project…', hint: 'workspace', shortcut: '⌃⇧N' },
      { id: 'c:new-file', label: 'New file…', hint: 'file', shortcut: '' },
      { id: 'c:close-tab', label: 'Close tab', hint: 'file', shortcut: '' },
      { id: 'c:zen', label: 'Toggle zen mode', hint: 'view', shortcut: '' },
      { id: 'c:out', label: 'Toggle output panel', hint: 'view', shortcut: '' },
      { id: 'c:term', label: 'Toggle integrated terminal', hint: 'view', shortcut: '⌃J' },
      { id: 'c:pv', label: 'Toggle 3D preview pane', hint: 'view', shortcut: '' },
      { id: 'c:run', label: 'Run current file (VM)', hint: 'koda', shortcut: 'F5' },
      { id: 'c:build', label: 'Build native executable', hint: 'koda', shortcut: '⌃⇧B' },
      { id: 'c:sidebar', label: 'Pin / toggle file drawer (rail hover)', hint: 'view', shortcut: '⌃B' },
      { id: 'c:save', label: 'Save current file', hint: 'file', shortcut: '⌃S' },
      { id: 'c:settings', label: 'Settings (placeholder)', hint: 'settings', shortcut: '⌃,' },
    ]
    const files = []
    async function walk(rel) {
      let entries = []
      try {
        entries = await ListDir(rel)
      } catch {
        return
      }
      for (const e of entries) {
        if (e.isDir) await walk(e.rel)
        else if (e.rel.endsWith('.koda')) {
          files.push({ id: 'f:' + e.rel, label: e.rel, hint: 'file', rel: e.rel })
        }
      }
    }
    if (workspace) await walk('')
    return [...cmds, ...files]
  }

  let paletteItems = $state([])

  async function openPalette() {
    paletteItems = await collectPaletteItems()
    paletteOpen = true
  }

  function onPalettePick(it) {
    if (it.id === 'c:open') void openWorkspaceFlow()
    else if (it.id === 'c:new-project') void startNewProjectWizard()
    else if (it.id === 'c:new-file') {
      if (workspace) openNewFileDialog()
    } else if (it.id === 'c:close-tab') closeActiveTab()
    else if (it.id === 'c:zen') zen = !zen
    else if (it.id === 'c:out') outputOpen = !outputOpen
    else if (it.id === 'c:term') showTerminal = !showTerminal
    else if (it.id === 'c:pv') showPreview = !showPreview
    else if (it.id === 'c:run') void runActive()
    else if (it.id === 'c:build') void buildActive()
    else if (it.id === 'c:sidebar') sidebarPinned = !sidebarPinned
    else if (it.id === 'c:save') void saveActive()
    else if (it.id === 'c:settings') {
      /* placeholder — wire to settings route later */
    } else if (it.rel) void openRel(it.rel)
  }

  async function onDiagnosticJump(absPath, line, col) {
    const rel = absToRel(absPath)
    if (rel) await openRel(rel)
    await tick()
    jumpTarget = { line, col, token: Date.now() }
  }

  function onJumpConsumed() {
    jumpTarget = null
  }

  /** @param {KeyboardEvent} e */
  function onGlobalKey(e) {
    if ((e.ctrlKey || e.metaKey) && e.key.toLowerCase() === 'b') {
      e.preventDefault()
      sidebarPinned = !sidebarPinned
      return
    }
    if ((e.ctrlKey || e.metaKey) && e.key.toLowerCase() === 'p') {
      e.preventDefault()
      void openPalette()
      return
    }
    if ((e.ctrlKey || e.metaKey) && e.key.toLowerCase() === 'k') {
      e.preventDefault()
      void openPalette()
      return
    }
    if ((e.ctrlKey || e.metaKey) && e.key.toLowerCase() === 'j' && !e.shiftKey) {
      e.preventDefault()
      showTerminal = !showTerminal
      return
    }
    if ((e.ctrlKey || e.metaKey) && e.key === ',') {
      e.preventDefault()
      /* settings placeholder */
      return
    }
    if (e.key === 'F5' && !e.ctrlKey && !e.metaKey) {
      e.preventDefault()
      void runActive()
    }
    if ((e.ctrlKey || e.metaKey) && e.shiftKey && e.key.toLowerCase() === 'b') {
      e.preventDefault()
      void buildActive()
    }
    if ((e.ctrlKey || e.metaKey) && e.shiftKey && e.key.toLowerCase() === 'n') {
      e.preventDefault()
      void startNewProjectWizard()
    }
  }

  onMount(async () => {
    workspace = await GetWorkspaceRoot()
    await refreshTree()
    const init = await LSPMessage(rpcInitialize(1))
    void init
    void LSPMessage(notifyInitialized())
    void tryAmbientContrast()

    unsubOut.push(
      EventsOn('koda:stdout', (line) => {
        outputText = outputText + line
      }),
    )
    unsubOut.push(
      EventsOn('koda:stderr', (line) => {
        outputText = outputText + line
      }),
    )
    unsubOut.push(
      EventsOn('koda:runDone', () => {
        outputText = outputText + '\n[run finished]\n'
      }),
    )
    unsubOut.push(
      EventsOn('koda:buildDone', () => {
        outputText = outputText + '\n[build finished]\n'
      }),
    )
    unsubOut.push(
      EventsOn('lsp:publishDiagnostics', (payload) => {
        updateDiagCounts(payload)
      }),
    )

    window.addEventListener('keydown', onGlobalKey)
  })

  onDestroy(() => {
    window.removeEventListener('keydown', onGlobalKey)
    for (const u of unsubOut) {
      if (typeof u === 'function') u()
    }
    unsubOut = []
    clearTimeout(lspTimer)
  })
</script>

<div class="relative h-full text-[var(--color-text)]" style={themeNudge ? 'filter: contrast(1.06);' : ''}>
<div
  class="void-app h-full"
  class:void-app--zen={zen}
  class:void-app--preview={showPreview && !zen}
>
  {#if !zen}
    <header
      class="void-header glass flex h-8 shrink-0 items-center gap-2 border-b border-[var(--color-surface0)] px-3 text-sm"
    >
      <span class="font-semibold tracking-tight text-[var(--color-accent)]">Koda</span>
      <FileMenuBar
        {zen}
        hasWorkspace={!!workspace}
        onNewProject={() => void startNewProjectWizard()}
        onOpenWorkspace={() => void openWorkspaceFlow()}
        onNewFile={() => openNewFileDialog()}
        onSave={() => void saveActive()}
        onCloseTab={() => closeActiveTab()}
        canSave={!!active}
        canCloseTab={!!active}
      />
      <button
        type="button"
        class="rounded-md px-2 py-0.5 text-[var(--color-subtext)] hover:bg-[var(--color-surface0)]"
        onclick={() => void openPalette()}>Palette <kbd class="font-sans text-[10px] text-[var(--color-overlay0)]">⌃P</kbd></button>
      <button
        type="button"
        class="rounded-md px-2 py-0.5 text-[var(--color-subtext)] hover:bg-[var(--color-surface0)]"
        onclick={() => (zen = !zen)}>Zen</button>
      <button
        type="button"
        class="rounded-md px-2 py-0.5 text-[var(--color-subtext)] hover:bg-[var(--color-surface0)]"
        onclick={() => (outputOpen = !outputOpen)}>Output</button>
      <button
        type="button"
        class="rounded-md px-2 py-0.5 text-[var(--color-subtext)] hover:bg-[var(--color-surface0)]"
        onclick={() => void runActive()}
        disabled={!active}>Run</button>
      <button
        type="button"
        class="rounded-md px-2 py-0.5 text-[var(--color-subtext)] hover:bg-[var(--color-surface0)]"
        onclick={() => void buildActive()}
        disabled={!active}>Build</button>
      <span class="ml-auto truncate text-xs text-[var(--color-overlay0)]" title={workspace || ''}
        >{workspace || 'No folder open'}</span>
    </header>
  {/if}

  <Sidebar
    {zen}
    {workspace}
    {tree}
    activeRel={active?.rel ?? ''}
    bind:drawerPinned={sidebarPinned}
    onOpenFile={(rel) => void openRel(rel)}
    onOpenWorkspace={() => void openWorkspaceFlow()}
    onRun={() => void runActive()}
    onOpenPalette={() => void openPalette()}
  />

  <main class="void-main flex min-h-0 flex-col">
    {#if !zen}
      <div
        class="flex shrink-0 gap-0.5 border-b border-[var(--color-surface0)] bg-[var(--color-mantle)] px-2 py-0.5 text-xs"
      >
        {#each tabs as t, i (t.rel)}
          <button
            type="button"
            class="max-w-[160px] truncate rounded-md border-b-2 border-transparent px-2 py-1 transition-colors"
            class:border-[var(--color-accent)]={i === activeIndex}
            class:bg-[var(--color-surface0)]={i === activeIndex}
            class:text-[var(--color-text)]={i === activeIndex}
            class:text-[var(--color-subtext)]={i !== activeIndex}
            onclick={() => {
              jumpTarget = null
              activeIndex = i
            }}>{t.name}</button>
        {/each}
      </div>
    {/if}

    <div class="flex min-h-0 flex-1 flex-col">
      <section class="flex min-h-0 min-w-0 flex-1 flex-col p-2">
        {#if active}
          {#key active.rel}
            {#await AbsFromWorkspace(active.rel)}
              <p class="p-4 text-sm text-[var(--color-subtext)]">Loading editor…</p>
            {:then abs}
              <CodeEditor
                absPath={abs}
                relPath={active.rel}
                seed={active.text}
                {jumpTarget}
                onTextChange={(t) => {
                  const idx = activeIndex
                  tabs = tabs.map((tab, i) => (i === idx ? { ...tab, text: t } : tab))
                }}
                onLspNotify={(p, txt) => pushLsp(p, txt)}
                onSave={() => void saveActive()}
                onJumpConsumed={onJumpConsumed}
              />
            {:catch}
              <p class="p-4 text-sm text-[var(--color-red)]">Could not resolve file path.</p>
            {/await}
          {/key}
        {:else}
          <div
            class="flex flex-1 flex-col items-center justify-center gap-3 rounded-lg border border-dashed border-[var(--color-surface1)] bg-[var(--color-mantle)]/50 p-8 text-center text-[var(--color-subtext)]"
          >
            <p class="max-w-md text-sm">
              Start a <strong class="text-[var(--color-text)]">new project</strong> (starter <code class="text-[var(--color-accent)]">main.koda</code>) or open an existing folder. Then edit, save with
              <kbd class="rounded bg-[var(--color-surface0)] px-1">⌃S</kbd>, run with <kbd class="rounded bg-[var(--color-surface0)] px-1">F5</kbd>.
            </p>
            <div class="flex flex-wrap items-center justify-center gap-2">
              <button
                type="button"
                class="rounded-lg bg-[var(--color-accent)] px-4 py-2 font-medium text-[var(--color-crust)]"
                onclick={() => void startNewProjectWizard()}>New project…</button>
              <button
                type="button"
                class="rounded-lg border border-[var(--color-surface1)] px-4 py-2 font-medium text-[var(--color-text)] hover:bg-[var(--color-surface0)]"
                onclick={() => void openWorkspaceFlow()}>Open workspace…</button>
            </div>
          </div>
        {/if}
      </section>

      {#if showTerminal && !zen}
        <div class="h-56 shrink-0 border-t border-[var(--color-surface0)] p-2">
          <TerminalPane active={showTerminal} />
        </div>
      {/if}
    </div>
  </main>

  {#if showPreview && !zen}
    <aside class="void-preview glass min-h-0 overflow-y-auto border-l border-[var(--color-surface0)] p-2">
      <div class="mb-2 text-xs font-medium text-[var(--color-overlay0)]">Preview</div>
      <PreviewViewport />
    </aside>
  {/if}

  <StatusBar {zen} scope={active ? active.rel : ''} diagErrors={diagErrors} diagWarnings={diagWarnings} />
</div>

  <DiagnosticToast activeAbsPath={activeAbsPath} onJump={onDiagnosticJump} />

  {#if outputOpen}
    <aside
      class="glass fixed bottom-[22px] right-0 top-8 z-40 flex w-[min(480px,92vw)] min-h-0 flex-col border-l border-[var(--color-surface0)] shadow-2xl"
    >
      <div class="flex h-8 shrink-0 items-center justify-between border-b border-[var(--color-surface0)] px-3 text-sm">
        <span>Output</span>
        <button
          type="button"
          class="text-[var(--color-overlay0)] hover:text-[var(--color-text)]"
          onclick={() => (outputOpen = false)}>Close</button>
      </div>
      <pre
        class="h-[calc(100vh-2rem-22px-2rem)] min-h-0 flex-1 overflow-auto p-3 font-mono text-xs leading-relaxed whitespace-pre-wrap text-[var(--color-subtext)]">{outputText}</pre>
    </aside>
  {/if}

  {#if newProjectOpen}
    <div
      class="fixed inset-0 z-[60] flex items-center justify-center bg-black/50 p-4 backdrop-blur-sm"
      role="presentation"
      onclick={(e) => {
        if (e.target === e.currentTarget) newProjectOpen = false
      }}
      onkeydown={(e) => {
        if (e.key === 'Escape') newProjectOpen = false
      }}
    >
      <div
        class="glass w-full max-w-md rounded-xl border border-[var(--color-surface0)] p-5 shadow-2xl"
        role="dialog"
        aria-modal="true"
        aria-labelledby="np-title"
        tabindex="-1"
        onclick={(e) => e.stopPropagation()}
        onkeydown={(e) => e.stopPropagation()}
      >
        <h2 id="np-title" class="mb-1 text-base font-semibold text-[var(--color-text)]">New project</h2>
        {#if newProjectParent}
          <p class="mb-4 text-xs text-[var(--color-subtext)]">
            Folder will be created inside:
            <span class="break-all font-mono text-[var(--color-overlay0)]">{newProjectParent}</span>
          </p>
          <label class="mb-1 block text-xs font-medium text-[var(--color-subtext)]" for="np-name">Project folder name</label>
          <input
            id="np-name"
            class="mb-3 w-full rounded-md border border-[var(--color-surface0)] bg-[var(--color-crust)] px-3 py-2 text-sm text-[var(--color-text)] outline-none focus:border-[var(--color-accent)]"
            bind:value={newProjectName}
            autocomplete="off"
          />
        {:else}
          <p class="mb-4 whitespace-pre-wrap text-xs text-[var(--color-red)]">{newProjectErr}</p>
        {/if}
        {#if newProjectParent && newProjectErr}
          <p class="mb-3 whitespace-pre-wrap text-xs text-[var(--color-red)]">{newProjectErr}</p>
        {/if}
        <div class="flex justify-end gap-2">
          <button
            type="button"
            class="rounded-md px-3 py-1.5 text-sm text-[var(--color-subtext)] hover:bg-[var(--color-surface0)]"
            onclick={() => (newProjectOpen = false)}>{newProjectParent ? 'Cancel' : 'Close'}</button>
          {#if newProjectParent}
            <button
              type="button"
              class="rounded-md bg-[var(--color-accent)] px-3 py-1.5 text-sm font-medium text-[var(--color-crust)]"
              onclick={() => void confirmNewProject()}>Create &amp; open</button>
          {/if}
        </div>
      </div>
    </div>
  {/if}

  {#if newFileOpen}
    <div
      class="fixed inset-0 z-[60] flex items-center justify-center bg-black/50 p-4 backdrop-blur-sm"
      role="presentation"
      onclick={(e) => {
        if (e.target === e.currentTarget) newFileOpen = false
      }}
      onkeydown={(e) => {
        if (e.key === 'Escape') newFileOpen = false
      }}
    >
      <div
        class="glass w-full max-w-md rounded-xl border border-[var(--color-surface0)] p-5 shadow-2xl"
        role="dialog"
        aria-modal="true"
        aria-labelledby="nf-title"
        tabindex="-1"
        onclick={(e) => e.stopPropagation()}
        onkeydown={(e) => e.stopPropagation()}
      >
        <h2 id="nf-title" class="mb-1 text-base font-semibold text-[var(--color-text)]">New file</h2>
        <p class="mb-4 text-xs text-[var(--color-subtext)]">Path relative to workspace (use <code class="text-[var(--color-accent)]">/</code>). Example: <code class="text-[var(--color-accent)]">src/app.koda</code></p>
        <label class="mb-1 block text-xs font-medium text-[var(--color-subtext)]" for="nf-rel">File path</label>
        <input
          id="nf-rel"
          class="mb-3 w-full rounded-md border border-[var(--color-surface0)] bg-[var(--color-crust)] px-3 py-2 font-mono text-sm text-[var(--color-text)] outline-none focus:border-[var(--color-accent)]"
          bind:value={newFileRel}
          autocomplete="off"
        />
        {#if newFileErr}
          <p class="mb-3 text-xs text-[var(--color-red)]">{newFileErr}</p>
        {/if}
        <div class="flex justify-end gap-2">
          <button
            type="button"
            class="rounded-md px-3 py-1.5 text-sm text-[var(--color-subtext)] hover:bg-[var(--color-surface0)]"
            onclick={() => (newFileOpen = false)}>Cancel</button>
          <button
            type="button"
            class="rounded-md bg-[var(--color-accent)] px-3 py-1.5 text-sm font-medium text-[var(--color-crust)]"
            onclick={() => void confirmNewFile()}>Create &amp; open</button>
        </div>
      </div>
    </div>
  {/if}

  <CommandPalette bind:open={paletteOpen} items={paletteItems} onPick={onPalettePick} />
</div>
