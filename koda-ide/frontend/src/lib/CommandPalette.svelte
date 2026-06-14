<script>
  import Fuse from 'fuse.js'

  let {
    open = $bindable(false),
    items = [],
    onPick = () => {},
  } = $props()

  let q = $state('')
  let inputEl = $state(null)
  let cursor = $state(0)
  let listEl = $state(null)

  /** @param {readonly [number, number][]} indices */
  function mergeInclusiveRanges(indices) {
    const arr = indices.map(([a, b]) => [Math.min(a, b), Math.max(a, b)])
    arr.sort((x, y) => x[0] - y[0])
    /** @type [number, number][] */
    const out = []
    for (const [s, e] of arr) {
      if (!out.length || s > out[out.length - 1][1] + 1) out.push([s, e])
      else out[out.length - 1][1] = Math.max(out[out.length - 1][1], e)
    }
    return out
  }

  /**
   * @param {string} label
   * @param {readonly import('fuse.js').FuseResultMatch[] | undefined} matches
   */
  function highlightParts(label, matches) {
    const labelMatch = matches?.find((m) => m.key === 'label')
    if (!labelMatch?.indices?.length) return [{ text: label, hi: false }]
    const merged = mergeInclusiveRanges([...labelMatch.indices])
    /** @type { { text: string, hi: boolean }[] } */
    const parts = []
    let i = 0
    for (const [s, e] of merged) {
      if (i < s) parts.push({ text: label.slice(i, s), hi: false })
      parts.push({ text: label.slice(s, e + 1), hi: true })
      i = e + 1
    }
    if (i < label.length) parts.push({ text: label.slice(i), hi: false })
    return parts
  }

  const filtered = $derived.by(() => {
    const query = q.trim()
    if (!items.length) return []
    const fuse = new Fuse(items, {
      keys: ['label', 'hint'],
      threshold: 0.35,
      includeMatches: true,
      ignoreLocation: true,
    })
    if (!query) return items.slice(0, 40).map((item) => ({ item, matches: [] }))
    return fuse.search(query).slice(0, 40)
  })

  $effect(() => {
    if (open && inputEl) {
      queueMicrotask(() => inputEl?.focus())
      cursor = 0
    }
  })

  $effect(() => {
    if (cursor >= filtered.length) cursor = Math.max(0, filtered.length - 1)
  })

  $effect(() => {
    if (!open || !listEl) return
    const row = listEl.querySelector(`[data-pal-idx="${cursor}"]`)
    row?.scrollIntoView({ block: 'nearest' })
  })

  function pick(it) {
    onPick(it)
    open = false
    q = ''
    cursor = 0
  }

  /** @param {KeyboardEvent} e */
  function onKeydown(e) {
    if (!open) return
    if (e.key === 'Escape') {
      open = false
      q = ''
      cursor = 0
      return
    }
    if (e.key === 'ArrowDown') {
      e.preventDefault()
      if (!filtered.length) return
      cursor = Math.min(filtered.length - 1, cursor + 1)
      return
    }
    if (e.key === 'ArrowUp') {
      e.preventDefault()
      if (!filtered.length) return
      cursor = Math.max(0, cursor - 1)
      return
    }
    if (e.key === 'Enter') {
      e.preventDefault()
      const row = filtered[cursor]
      if (row && row.item) pick(row.item)
    }
  }
</script>

<svelte:window onkeydown={onKeydown} />

{#if open}
  <div
    class="fixed inset-0 z-50 flex items-start justify-center bg-black/40 pt-[12vh] backdrop-blur-sm"
    role="presentation"
    onclick={(e) => {
      if (e.target === e.currentTarget) {
        open = false
        q = ''
        cursor = 0
      }
    }}
  >
    <div
      class="glass w-[min(560px,92vw)] overflow-hidden rounded-xl shadow-2xl"
      role="dialog"
      aria-modal="true"
      aria-label="K-Bar command palette"
    >
      <input
        bind:this={inputEl}
        class="w-full border-0 bg-transparent px-4 py-3 text-sm text-[var(--color-text)] outline-none placeholder:text-[var(--color-overlay0)]"
        placeholder="Files, commands…"
        bind:value={q}
      />
      <ul bind:this={listEl} class="max-h-[50vh] overflow-y-auto border-t border-[var(--color-surface0)] text-sm">
        {#each filtered as row, i (row.item.id)}
          {@const parts = highlightParts(row.item.label, row.matches)}
          <li>
            <button
              type="button"
              data-pal-idx={i}
              class="flex w-full cursor-pointer items-center gap-2 px-4 py-2 text-left text-[var(--color-text)] hover:bg-[var(--color-surface0)]"
              class:bg-[var(--color-surface0)]={i === cursor}
              onclick={() => pick(row.item)}
            >
              <span class="min-w-0 flex-1 truncate">
                {#each parts as seg, si (si)}
                  <span class:palette-match={seg.hi}>{seg.text}</span>
                {/each}
              </span>
              {#if row.item.hint}
                <span class="shrink-0 text-xs text-[var(--color-overlay0)]">{row.item.hint}</span>
              {/if}
              {#if row.item.shortcut}
                <kbd class="shrink-0 font-sans text-[10px] tracking-wide text-[var(--color-overlay0)] opacity-80">{row.item.shortcut}</kbd>
              {/if}
            </button>
          </li>
        {/each}
      </ul>
    </div>
  </div>
{/if}
