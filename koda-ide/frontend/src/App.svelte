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
    CheckSDK,
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
  import ThemePicker from './lib/ThemePicker.svelte'
  import WelcomeScreen from './lib/WelcomeScreen.svelte'
  import { defaultTheme, safeTheme } from './lib/themes.js'
  import HelpPanel from './lib/HelpPanel.svelte'
  import { pathToFileURI, rpcInitialize, notifyInitialized, notifyDidOpen, notifyDidChange } from './lib/lspUtil.js'

  const defaultStudioSettings = {
    theme: defaultTheme,
    editorFontSize: 15,
    editorLineHeight: 1.65,
    editorTabSize: 4,
    terminalFontSize: 13,
    panelOpacity: 92,
    compactMode: false,
    showHeaderThemePicker: true,
  }

  let workspace = $state('')
  let tree = $state([])
  let tabs = $state([])
  let activeIndex = $state(-1)
  let zen = $state(false)
  let outputOpen = $state(false)
  let outputText = $state('')
  let outputEl = $state(null)
  let paletteOpen = $state(false)
  let showTerminal = $state(false)
  let showPreview = $state(false)
  let theme = $state(defaultTheme)
  let themeReady = $state(false)
  let themeNudge = $state(0)
  let editorFontSize = $state(defaultStudioSettings.editorFontSize)
  let editorLineHeight = $state(defaultStudioSettings.editorLineHeight)
  let editorTabSize = $state(defaultStudioSettings.editorTabSize)
  let terminalFontSize = $state(defaultStudioSettings.terminalFontSize)
  let panelOpacity = $state(defaultStudioSettings.panelOpacity)
  let compactMode = $state(defaultStudioSettings.compactMode)
  let showHeaderThemePicker = $state(defaultStudioSettings.showHeaderThemePicker)
  let openSettingsToken = $state(0)
  let sidebarPinned = $state(false)
  let activeAbsPath = $state('')
  let diagErrors = $state(0)
  let diagWarnings = $state(0)
  let cursorLine = $state(1)
  let cursorCol = $state(1)
  /** @type { { line: number, col: number, token: number } | null } */
  let jumpTarget = $state(null)
  let newProjectOpen = $state(false)
  let newProjectParent = $state('')
  let newProjectName = $state('my-koda-project')
  let newProjectTemplate = $state('hello')
  let newProjectErr = $state('')
  /** @type {{ ok: boolean, version?: string, installDir?: string, stdlibDir?: string, lines?: { ok: boolean, label: string, detail: string, fix?: string }[] }} */
  let sdkStatus = $state({ ok: false, lines: [] })
  let startupError = $state('')
  let newFileOpen = $state(false)
  let newFileRel = $state('hello.koda')
  let newFileErr = $state('')
  let helpOpen = $state(false)
  let helpPage = $state('START_HERE.md')
  let unsubOut = []
  let lspTimer = 0

  const active = $derived(activeIndex >= 0 ? tabs[activeIndex] : null)
  const dirtyCount = $derived(tabs.filter((t) => t.dirty).length)
  const shellStyle = $derived(
    `${themeNudge ? 'filter: contrast(1.06);' : ''}` +
      `--editor-font-size:${clampNumber(editorFontSize, 12, 22, defaultStudioSettings.editorFontSize)}px;` +
      `--editor-line-height:${clampNumber(editorLineHeight, 1.35, 2, defaultStudioSettings.editorLineHeight)};` +
      `--terminal-font-size:${clampNumber(terminalFontSize, 11, 20, defaultStudioSettings.terminalFontSize)}px;` +
      `--panel-opacity:${clampNumber(panelOpacity, 76, 100, defaultStudioSettings.panelOpacity)}%;`,
  )

  function clampNumber(value, min, max, fallback) {
    const n = Number(value)
    if (!Number.isFinite(n)) return fallback
    return Math.min(max, Math.max(min, n))
  }

  function boolValue(value, fallback) {
    return typeof value === 'boolean' ? value : fallback
  }

  function readStudioSettings() {
    let raw = null
    try {
      raw = JSON.parse(localStorage.getItem('koda-studio-settings') || 'null')
    } catch {
      raw = null
    }
    return {
      theme: safeTheme(raw?.theme || localStorage.getItem('koda-studio-theme') || defaultStudioSettings.theme),
      editorFontSize: clampNumber(raw?.editorFontSize, 12, 22, defaultStudioSettings.editorFontSize),
      editorLineHeight: clampNumber(raw?.editorLineHeight, 1.35, 2, defaultStudioSettings.editorLineHeight),
      editorTabSize: [2, 4, 8].includes(Number(raw?.editorTabSize)) ? Number(raw.editorTabSize) : defaultStudioSettings.editorTabSize,
      terminalFontSize: clampNumber(raw?.terminalFontSize, 11, 20, defaultStudioSettings.terminalFontSize),
      panelOpacity: clampNumber(raw?.panelOpacity, 76, 100, defaultStudioSettings.panelOpacity),
      compactMode: boolValue(raw?.compactMode, defaultStudioSettings.compactMode),
      showHeaderThemePicker: boolValue(raw?.showHeaderThemePicker, defaultStudioSettings.showHeaderThemePicker),
    }
  }

  function applyStudioSettings(settings) {
    theme = safeTheme(settings.theme)
    editorFontSize = settings.editorFontSize
    editorLineHeight = settings.editorLineHeight
    editorTabSize = settings.editorTabSize
    terminalFontSize = settings.terminalFontSize
    panelOpacity = settings.panelOpacity
    compactMode = settings.compactMode
    showHeaderThemePicker = settings.showHeaderThemePicker
  }

  function resetStudioSettings() {
    applyStudioSettings(defaultStudioSettings)
  }

  function openSettings() {
    openSettingsToken += 1
    sidebarPinned = true
  }

  $effect(() => {
    if (!outputEl || !outputOpen) return
    outputText
    queueMicrotask(() => {
      if (outputEl) outputEl.scrollTop = outputEl.scrollHeight
    })
  })

  $effect(() => {
    if (!themeReady || typeof localStorage === 'undefined') return
    localStorage.setItem('koda-studio-theme', theme)
    localStorage.setItem(
      'koda-studio-settings',
      JSON.stringify({
        theme,
        editorFontSize: clampNumber(editorFontSize, 12, 22, defaultStudioSettings.editorFontSize),
        editorLineHeight: clampNumber(editorLineHeight, 1.35, 2, defaultStudioSettings.editorLineHeight),
        editorTabSize: [2, 4, 8].includes(Number(editorTabSize)) ? Number(editorTabSize) : defaultStudioSettings.editorTabSize,
        terminalFontSize: clampNumber(terminalFontSize, 11, 20, defaultStudioSettings.terminalFontSize),
        panelOpacity: clampNumber(panelOpacity, 76, 100, defaultStudioSettings.panelOpacity),
        compactMode,
        showHeaderThemePicker,
      }),
    )
  })

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

  function confirmDiscardDirty(message = 'You have unsaved files. Continue without saving?') {
    if (dirtyCount === 0) return true
    return window.confirm(message)
  }

  async function refreshTree() {
    if (!workspace) {
      tree = []
      return
    }
    try {
      tree = await loadTree('')
    } catch {
      tree = []
    }
  }

  async function loadTree(rel) {
    const entries = await ListDir(rel)
    return Promise.all(
      entries.map(async (entry) => {
        if (!entry.isDir) return entry
        try {
          return { ...entry, children: await loadTree(entry.rel) }
        } catch {
          return { ...entry, children: [] }
        }
      }),
    )
  }

  async function openWorkspaceFlow() {
    if (!confirmDiscardDirty('Open another workspace and discard unsaved tab changes?')) return
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

  async function startNewProjectWizard(template = 'hello') {
    if (!confirmDiscardDirty('Create a new project and discard unsaved tab changes?')) return
    newProjectErr = ''
    newProjectTemplate = template || 'hello'
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

  async function openProjectEntry() {
    for (const rel of ['src/main.koda', 'main.koda']) {
      try {
        await openRel(rel)
        return
      } catch {
        /* try next */
      }
    }
  }

  async function confirmNewProject() {
    newProjectErr = ''
    if (!newProjectParent.trim()) return
    try {
      const root = await CreateProjectInParent(
        newProjectParent,
        newProjectName.trim(),
        newProjectTemplate || 'hello',
      )
      await OpenWorkspace(root)
      workspace = await GetWorkspaceRoot()
      tabs = []
      activeIndex = -1
      outputText = ''
      await refreshTree()
      newProjectOpen = false
      await openProjectEntry()
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
    if (tabs[activeIndex]?.dirty && !window.confirm(`Close ${tabs[activeIndex].name} without saving?`)) return
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
    tabs = [...tabs, { rel, name: rel.split('/').pop() || rel, text, savedText: text, dirty: false }]
    activeIndex = tabs.length - 1
    const abs = await AbsFromWorkspace(rel)
    void LSPMessage(notifyDidOpen(pathToFileURI(abs), text))
  }

  async function saveActive() {
    if (!active) return
    await WriteFile(active.rel, active.text)
    const idx = activeIndex
    tabs = tabs.map((tab, i) => (i === idx ? { ...tab, savedText: tab.text, dirty: false } : tab))
    outputText = outputText + `Saved ${active.rel}\n`
  }

  async function runActive() {
    if (!active) return
    const abs = await AbsFromWorkspace(active.rel)
    outputOpen = true
    outputText = outputText + `\n> koda run ${active.rel}\n`
    RunProgram(abs, active.text)
  }

  async function buildActive() {
    if (!active) return
    const abs = await AbsFromWorkspace(active.rel)
    const out = await DefaultBuildOutput(abs)
    outputOpen = true
    outputText = outputText + `\n> koda build ${active.rel}\n`
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

  function openHelp(page = 'START_HERE.md') {
    helpPage = page
    helpOpen = true
  }

  async function collectPaletteItems() {
    const cmds = [
      { id: 'c:help', label: 'Help & documentation', hint: 'help', shortcut: 'F1' },
      { id: 'c:help-begin', label: 'Beginner\'s guide', hint: 'help', shortcut: '' },
      { id: 'c:help-faq', label: 'FAQ', hint: 'help', shortcut: '' },
      { id: 'c:open', label: 'Open workspace folder…', hint: 'workspace', shortcut: '' },
      { id: 'c:new-project', label: 'New project…', hint: 'workspace', shortcut: '⌃⇧N' },
      { id: 'c:new-file', label: 'New file…', hint: 'file', shortcut: '' },
      { id: 'c:close-tab', label: 'Close tab', hint: 'file', shortcut: '' },
      { id: 'c:zen', label: 'Toggle zen mode', hint: 'view', shortcut: '' },
      { id: 'c:out', label: 'Toggle output panel', hint: 'view', shortcut: '' },
      { id: 'c:term', label: 'Toggle integrated terminal', hint: 'view', shortcut: '⌃J' },
      { id: 'c:pv', label: 'Toggle 3D preview pane', hint: 'view', shortcut: '' },
      { id: 'c:run', label: 'Run current file', hint: 'koda', shortcut: 'F5' },
      { id: 'c:build', label: 'Build native executable', hint: 'koda', shortcut: '⌃⇧B' },
      { id: 'c:sidebar', label: 'Pin / toggle file drawer (rail hover)', hint: 'view', shortcut: '⌃B' },
      { id: 'c:save', label: 'Save current file', hint: 'file', shortcut: '⌃S' },
      { id: 'c:settings', label: 'Open settings', hint: 'settings', shortcut: '⌃,' },
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
    if (it.id === 'c:help') openHelp('START_HERE.md')
    else if (it.id === 'c:help-begin') openHelp('docs/beginners-guide.md')
    else if (it.id === 'c:help-faq') openHelp('docs/faq.md')
    else if (it.id === 'c:open') void openWorkspaceFlow()
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
    else if (it.id === 'c:settings') openSettings()
    else if (it.rel) void openRel(it.rel)
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
    if ((e.ctrlKey || e.metaKey) && e.key.toLowerCase() === 's') {
      e.preventDefault()
      void saveActive()
      return
    }
    if ((e.ctrlKey || e.metaKey) && e.key.toLowerCase() === 'w') {
      e.preventDefault()
      closeActiveTab()
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
      openSettings()
      return
    }
    if (e.key === 'F1') {
      e.preventDefault()
      openHelp('START_HERE.md')
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
    try {
      applyStudioSettings(readStudioSettings())
    } catch {
      applyStudioSettings(defaultStudioSettings)
    }
    themeReady = true
    try {
      workspace = await GetWorkspaceRoot()
      await refreshTree()
      if (workspace) await openProjectEntry()
      try {
        sdkStatus = await CheckSDK()
      } catch {
        sdkStatus = { ok: false, lines: [{ ok: false, label: 'SDK', detail: 'Could not check SDK', fix: 'Run Koda Studio from the SDK folder.' }] }
      }
      const init = await LSPMessage(rpcInitialize(1))
      void init
      void LSPMessage(notifyInitialized())
      void tryAmbientContrast()
    } catch (err) {
      startupError = err instanceof Error ? err.message : String(err)
    }

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
    window.addEventListener('beforeunload', onBeforeUnload)
  })

  onDestroy(() => {
    window.removeEventListener('keydown', onGlobalKey)
    window.removeEventListener('beforeunload', onBeforeUnload)
    for (const u of unsubOut) {
      if (typeof u === 'function') u()
    }
    unsubOut = []
    clearTimeout(lspTimer)
  })

  function onBeforeUnload(e) {
    if (dirtyCount === 0) return
    e.preventDefault()
    e.returnValue = ''
  }
</script>

<div class="studio-shell" data-theme={theme} class:studio-shell--compact={compactMode} style={shellStyle}>
  <div
    class="void-app"
    class:void-app--zen={zen}
    class:void-app--preview={showPreview && !zen}
  >
    {#if !zen}
      <header class="void-header">
        <span class="studio-title">Koda Studio</span>
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
        onclick={() => openHelp()}>Help <kbd class="font-sans text-[10px] text-[var(--color-overlay0)]">F1</kbd></button>
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
      {#if showHeaderThemePicker}
        <ThemePicker bind:value={theme} />
      {/if}
      <span class="studio-spacer" title={workspace || ''}>{workspace || 'No folder open'}</span>
    </header>
  {/if}

  <Sidebar
    {zen}
    {workspace}
    {tree}
    activeRel={active?.rel ?? ''}
    bind:drawerPinned={sidebarPinned}
    {openSettingsToken}
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
    onOpenFile={(rel) => void openRel(rel)}
    onOpenWorkspace={() => void openWorkspaceFlow()}
    onRun={() => void runActive()}
    onOpenPalette={() => void openPalette()}
    onOpenHelp={(page) => openHelp(page)}
    onResetSettings={resetStudioSettings}
  />

  <main class="void-main">
    {#if !zen}
      <div
        class="flex shrink-0 gap-0.5 border-b border-[var(--color-surface0)] bg-[var(--color-mantle)] px-2 py-1 text-sm"
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
            }}>
              <span class="mr-1 text-[var(--color-accent)]">{t.dirty ? '●' : ''}</span>{t.name}
            </button>
        {/each}
      </div>
    {/if}

    <div class="void-main-body">
      <section class="void-main-editor p-2">
        {#if startupError}
          <div class="welcome">
            <h1>Startup error</h1>
            <p class="welcome-lead text-[var(--color-red)]">{startupError}</p>
            <div class="welcome-actions">
              <button type="button" onclick={() => void openWorkspaceFlow()}>Open workspace…</button>
            </div>
          </div>
        {:else if active}
          {#key `${active.rel}:${editorTabSize}`}
            {#await AbsFromWorkspace(active.rel)}
              <p class="p-4 text-sm text-[var(--color-subtext)]">Loading editor…</p>
            {:then abs}
              <CodeEditor
                absPath={abs}
                relPath={active.rel}
                seed={active.text}
                {jumpTarget}
                tabSize={editorTabSize}
                onTextChange={(t) => {
                  const idx = activeIndex
                  tabs = tabs.map((tab, i) =>
                    i === idx ? { ...tab, text: t, dirty: t !== tab.savedText } : tab,
                  )
                }}
                onCursorChange={(pos) => {
                  cursorLine = pos.line
                  cursorCol = pos.col
                }}
                onLspNotify={(p, txt) => pushLsp(p, txt)}
                onSave={() => void saveActive()}
                onJumpConsumed={onJumpConsumed}
              />
            {:catch}
              <p class="p-4 text-sm text-[var(--color-red)]">Could not resolve file path.</p>
            {/await}
          {/key}
        {:else if !workspace}
          <WelcomeScreen
            sdk={sdkStatus}
            onNewProject={(template) => void startNewProjectWizard(template)}
            onOpenWorkspace={() => void openWorkspaceFlow()}
            onOpenHelp={(page) => openHelp(page)}
          />
        {:else}
          <div
            class="flex flex-1 flex-col items-center justify-center gap-3 rounded-lg border border-dashed border-[var(--color-surface1)] bg-[var(--color-mantle)]/50 p-8 text-center text-[var(--color-subtext)]"
          >
            <p class="max-w-md text-sm">
              Open a <code class="text-[var(--color-cyan)]">.koda</code> file from the sidebar, or start with
              <code class="text-[var(--color-cyan)]">src/main.koda</code>. Save with
              <kbd class="rounded bg-[var(--color-surface0)] px-1">⌃S</kbd>, run with <kbd class="rounded bg-[var(--color-surface0)] px-1">F5</kbd>.
            </p>
            <button
              type="button"
              class="rounded-lg border border-[var(--color-surface1)] px-4 py-2 font-medium text-[var(--color-text)] hover:bg-[var(--color-surface0)]"
              onclick={() => void openProjectEntry()}>Open main.koda</button>
          </div>
        {/if}
      </section>

      {#if showTerminal && !zen}
        <div class="h-56 shrink-0 border-t border-[var(--color-surface0)] p-2">
          {#key `${theme}:${terminalFontSize}`}
            <TerminalPane active={showTerminal} fontSize={terminalFontSize} />
          {/key}
        </div>
      {/if}

      {#if outputOpen && !zen}
        <aside class="output-dock">
          <div class="output-dock-header">
            <span>Output</span>
            <button type="button" onclick={() => (outputText = '')}>Clear</button>
            <button type="button" onclick={() => (outputOpen = false)}>Close</button>
          </div>
          <pre class="output-panel" bind:this={outputEl}>{outputText || 'No output yet.'}</pre>
        </aside>
      {/if}
    </div>
  </main>

  {#if showPreview && !zen}
    <aside class="void-preview glass min-h-0 overflow-y-auto border-l border-[var(--color-surface0)] p-2">
      <div class="mb-2 text-sm font-medium text-[var(--color-overlay0)]">Preview</div>
      <PreviewViewport />
    </aside>
  {/if}

  <StatusBar
    {zen}
    scope={active ? active.rel : ''}
    {diagErrors}
    {diagWarnings}
    {dirtyCount}
    line={cursorLine}
    col={cursorCol}
  />
  </div>

  <div class="studio-overlays">
  <DiagnosticToast activeAbsPath={activeAbsPath} onJump={onDiagnosticJump} />

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
          <p class="mb-4 text-sm text-[var(--color-subtext)]">
            Folder will be created inside:
            <span class="break-all font-mono text-[var(--color-overlay0)]">{newProjectParent}</span>
          </p>
          <label class="mb-1 block text-sm font-medium text-[var(--color-subtext)]" for="np-name">Project folder name</label>
          <input
            id="np-name"
            class="mb-3 w-full rounded-md border border-[var(--color-surface0)] bg-[var(--color-crust)] px-3 py-2 text-sm text-[var(--color-text)] outline-none focus:border-[var(--color-accent)]"
            bind:value={newProjectName}
            autocomplete="off"
          />
          <fieldset class="mb-3">
            <legend class="mb-2 text-sm font-medium text-[var(--color-subtext)]">Template</legend>
            <div class="flex flex-col gap-1.5 text-sm">
              <label class="flex cursor-pointer items-center gap-2">
                <input type="radio" bind:group={newProjectTemplate} value="hello" />
                <span>Hello app — console output</span>
              </label>
              <label class="flex cursor-pointer items-center gap-2">
                <input type="radio" bind:group={newProjectTemplate} value="game" />
                <span>Text game — lunar lander in the terminal</span>
              </label>
              <label class="flex cursor-pointer items-center gap-2">
                <input type="radio" bind:group={newProjectTemplate} value="graphics" />
                <span>Bouncing ball — Raylib window</span>
              </label>
              <label class="flex cursor-pointer items-center gap-2">
                <input type="radio" bind:group={newProjectTemplate} value="pong" />
                <span>Pong — two-player paddle game</span>
              </label>
            </div>
          </fieldset>
        {:else}
          <p class="mb-4 whitespace-pre-wrap text-sm text-[var(--color-red)]">{newProjectErr}</p>
        {/if}
        {#if newProjectParent && newProjectErr}
          <p class="mb-3 whitespace-pre-wrap text-sm text-[var(--color-red)]">{newProjectErr}</p>
        {/if}
        <div class="flex justify-end gap-2">
          <button
            type="button"
            class="rounded-md px-3 py-1.5 text-sm text-[var(--color-subtext)] hover:bg-[var(--color-surface0)]"
            onclick={() => (newProjectOpen = false)}>{newProjectParent ? 'Cancel' : 'Close'}</button>
          {#if newProjectParent}
            <button
              type="button"
              class="rounded-md btn-primary px-3 py-1.5 text-sm font-medium"
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
        <p class="mb-4 text-sm text-[var(--color-subtext)]">Path relative to workspace (use <code class="text-[var(--color-accent)]">/</code>). Example: <code class="text-[var(--color-accent)]">src/app.koda</code></p>
        <label class="mb-1 block text-sm font-medium text-[var(--color-subtext)]" for="nf-rel">File path</label>
        <input
          id="nf-rel"
          class="mb-3 w-full rounded-md border border-[var(--color-surface0)] bg-[var(--color-crust)] px-3 py-2 font-mono text-sm text-[var(--color-text)] outline-none focus:border-[var(--color-accent)]"
          bind:value={newFileRel}
          autocomplete="off"
        />
        {#if newFileErr}
          <p class="mb-3 text-sm text-[var(--color-red)]">{newFileErr}</p>
        {/if}
        <div class="flex justify-end gap-2">
          <button
            type="button"
            class="rounded-md px-3 py-1.5 text-sm text-[var(--color-subtext)] hover:bg-[var(--color-surface0)]"
            onclick={() => (newFileOpen = false)}>Cancel</button>
          <button
            type="button"
            class="rounded-md btn-primary px-3 py-1.5 text-sm font-medium"
            onclick={() => void confirmNewFile()}>Create &amp; open</button>
        </div>
      </div>
    </div>
  {/if}

  <CommandPalette bind:open={paletteOpen} items={paletteItems} onPick={onPalettePick} />

  <HelpPanel bind:open={helpOpen} initialPage={helpPage} />
  </div>
</div>
