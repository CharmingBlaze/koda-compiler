<script>
  import { onDestroy, tick } from 'svelte'
  import { Terminal } from '@xterm/xterm'
  import { FitAddon } from '@xterm/addon-fit'
  import '@xterm/xterm/css/xterm.css'
  import { TerminalStart, TerminalWrite, TerminalResize, TerminalClose } from '../../wailsjs/go/main/App.js'
  import { EventsOn } from '../../wailsjs/runtime/runtime.js'

  let { active = false } = $props()

  let el = $state(null)
  let term = null
  let fit = null
  let sid = ''
  let offData = null
  let offExit = null
  let resizeHandler = null

  async function start() {
    if (!el || sid) return
    await tick()
    if (!el) return
    term = new Terminal({
      theme: {
        background: '#181926',
        foreground: '#cad3f5',
        cursor: '#c6a0f6',
      },
      fontFamily: "'JetBrains Mono Variable', ui-monospace, Cascadia Code, Consolas, monospace",
      fontSize: 13,
    })
    fit = new FitAddon()
    term.loadAddon(fit)
    term.open(el)
    fit.fit()
    const id = await TerminalStart(term.cols, term.rows)
    sid = id
    term.onData((d) => {
      void TerminalWrite(sid, d)
    })
    offData = EventsOn('term:data', (payload) => {
      if (payload?.id === sid && payload?.data) term.write(payload.data)
    })
    offExit = EventsOn('term:exit', (payload) => {
      if (payload?.id === sid) {
        term.writeln('\r\n\x1b[33m[session ended]\x1b[0m')
        cleanupSession()
      }
    })
    resizeHandler = () => {
      if (!fit || !term || !sid) return
      fit.fit()
      void TerminalResize(sid, term.cols, term.rows)
    }
    window.addEventListener('resize', resizeHandler)
  }

  function cleanupSession() {
    if (resizeHandler) {
      window.removeEventListener('resize', resizeHandler)
      resizeHandler = null
    }
    if (offData) {
      offData()
      offData = null
    }
    if (offExit) {
      offExit()
      offExit = null
    }
    if (sid) {
      void TerminalClose(sid)
      sid = ''
    }
    if (term) {
      term.dispose()
      term = null
    }
    fit = null
  }

  $effect(() => {
    if (active) {
      void start()
    } else {
      cleanupSession()
    }
  })

  onDestroy(() => cleanupSession())
</script>

<div class="h-full min-h-[160px] w-full overflow-hidden rounded-md border border-[var(--color-surface0)] bg-[var(--color-crust)]">
  <div class="h-full w-full" bind:this={el}></div>
</div>
