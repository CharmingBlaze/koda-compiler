<script>
  import FileTree from './FileTree.svelte'

  let {
    entries = [],
    activeRel = '',
    depth = 0,
    onOpenFile = () => {},
  } = $props()

  function norm(p) {
    return (p || '').replace(/\\/g, '/').toLowerCase()
  }
</script>

<ul class={depth === 0 ? 'space-y-0.5' : 'mt-0.5 space-y-0.5'}>
  {#each entries as entry (entry.rel)}
    <li>
      {#if entry.isDir}
        <details open class="group">
          <summary
            class="flex cursor-pointer select-none items-center gap-1 rounded px-1 py-0.5 text-sm text-[var(--color-subtext)] hover:bg-[var(--color-surface0)] hover:text-[var(--color-text)]"
            style={`padding-left: ${depth * 10 + 4}px`}
          >
            <span class="text-[10px] transition group-open:rotate-90">▶</span>
            <span class="truncate">{entry.name}</span>
          </summary>
          {#if entry.children?.length}
            <FileTree entries={entry.children} {activeRel} depth={depth + 1} {onOpenFile} />
          {/if}
        </details>
      {:else}
        <button
          type="button"
          class="flex w-full items-center gap-1 truncate rounded px-1 py-0.5 text-left text-sm hover:bg-[var(--color-surface0)]"
          class:text-[var(--color-accent)]={norm(activeRel) === norm(entry.rel)}
          class:text-[var(--color-text)]={norm(activeRel) !== norm(entry.rel)}
          style={`padding-left: ${depth * 10 + 16}px`}
          title={entry.rel}
          onclick={() => onOpenFile(entry.rel)}
        >
          <span class="text-[11px] font-semibold text-[var(--color-accent)]">K</span>
          <span class="truncate">{entry.name}</span>
        </button>
      {/if}
    </li>
  {/each}
</ul>
