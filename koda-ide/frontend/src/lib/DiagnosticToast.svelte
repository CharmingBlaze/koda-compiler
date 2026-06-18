<script>
  import { onMount, onDestroy } from 'svelte'
  import { fly } from 'svelte/transition'
  import { EventsOn } from '../../wailsjs/runtime/runtime.js'
  import CopyButton from './CopyButton.svelte'

  let {
    activeAbsPath = '',
    onJump = async () => {},
  } = $props()

  /** @type { { id: string, path: string, line: number, col: number, message: string }[] } */
  let queue = $state([])
  /** @type { Map<string, ReturnType<typeof setTimeout>> } */
  let timers = new Map()
  /** @type { null | (() => void) } */
  let unsub = null

  function norm(p) {
    return (p || '').replace(/\\/g, '/').toLowerCase()
  }

  function pathMatchesActive(path) {
    if (!activeAbsPath || !path) return false
    return norm(path) === norm(activeAbsPath)
  }

  function pushToast(entry) {
    const id = `${entry.path}:${entry.line}:${entry.col}:${Date.now()}`
    queue = [...queue.slice(-2), { ...entry, id }]
    const t = window.setTimeout(() => {
      queue = queue.filter((x) => x.id !== id)
      timers.delete(id)
    }, 6000)
    timers.set(id, t)
  }

  /** @param {any} payload */
  function handlePayload(payload) {
    if (!pathMatchesActive(payload?.path)) return
    const raw = payload?.diagnosticsRaw
    if (!Array.isArray(raw)) return
    for (const d of raw) {
      const sev = (d.severity || 'error').toLowerCase()
      if (sev === 'warning') continue
      pushToast({
        path: d.path || payload.path,
        line: d.line || 1,
        col: d.col || 1,
        message: d.message || 'Error',
      })
    }
  }

  onMount(() => {
    unsub = EventsOn('lsp:publishDiagnostics', handlePayload)
  })

  onDestroy(() => {
    if (typeof unsub === 'function') unsub()
    for (const t of timers.values()) clearTimeout(t)
    timers.clear()
  })

  async function clickEntry(entry) {
    await onJump(entry.path, entry.line, entry.col)
    queue = queue.filter((x) => x.id !== entry.id)
    const t = timers.get(entry.id)
    if (t) clearTimeout(t)
    timers.delete(entry.id)
  }
</script>

<div
  class="pointer-events-none fixed bottom-8 right-3 z-50 flex w-[min(380px,94vw)] flex-col gap-2"
  aria-live="polite"
>
  {#each queue as entry (entry.id)}
    <div
      class="pointer-events-auto glass w-full overflow-hidden rounded-lg border border-[var(--color-surface0)] px-3 py-2 text-xs text-[var(--color-subtext)] shadow-xl transition hover:border-[var(--color-accent)]/50"
      in:fly={{ x: 28, duration: 200 }}
    >
      <div class="mb-0.5 font-mono text-[10px] text-[var(--color-diag-error)]">L{entry.line}:{entry.col}</div>
      <div class="line-clamp-3 text-[var(--color-text)]">{entry.message}</div>
      <div class="mt-2 flex items-center justify-end gap-1.5">
        <CopyButton text={`${entry.path}:${entry.line}:${entry.col}\n${entry.message}`} label="Copy" copiedLabel="Copied" title="Copy diagnostic" compact />
        <button
          type="button"
          class="diag-jump-button"
          onclick={() => void clickEntry(entry)}
        >Jump</button>
      </div>
    </div>
  {/each}
</div>

<style>
  .diag-jump-button {
    border: 1px solid var(--color-surface1);
    border-radius: 0.35rem;
    background: var(--color-surface0);
    color: var(--color-text);
    cursor: pointer;
    font: inherit;
    font-size: 0.75rem;
    padding: 0.2rem 0.45rem;
  }

  .diag-jump-button:hover {
    border-color: var(--color-accent);
  }
</style>
