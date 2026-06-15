<script>
  let {
    zen = false,
    scope = '',
    diagErrors = 0,
    diagWarnings = 0,
    dirtyCount = 0,
    line = 1,
    col = 1,
  } = $props()
</script>

{#if !zen}
  <footer
    class="void-status flex h-[26px] shrink-0 items-center gap-3 border-t border-[var(--color-surface0)] bg-[var(--color-crust)] px-3 text-[13px] text-[var(--color-subtext)]"
  >
    <span class="rounded px-2 py-0.5 font-medium text-[var(--color-cyan)]" style="background: color-mix(in srgb, var(--color-cyan) 12%, transparent);">Koda Studio v0.4.0</span>
    <span class="min-w-0 flex-1 truncate text-[var(--color-subtext)]" title={scope}>{scope || '—'}</span>
    <span class="flex shrink-0 items-center gap-2">
      {#if dirtyCount > 0}
        <span class="text-[var(--color-yellow)]" title="Unsaved files">{dirtyCount} unsaved</span>
      {/if}
      {#if diagErrors > 0}
        <span class="flex items-center gap-0.5 text-[var(--color-diag-error)]" title="Errors"
          ><span class="inline-block h-1.5 w-1.5 rounded-full bg-[var(--color-diag-error)]"></span>{diagErrors}</span>
      {/if}
      {#if diagWarnings > 0}
        <span class="flex items-center gap-0.5 text-[var(--color-diag-warning)]" title="Warnings"
          ><span class="inline-block h-1.5 w-1.5 rounded-full bg-[var(--color-diag-warning)]"></span>{diagWarnings}</span>
      {/if}
      <span class="text-[var(--color-overlay0)]">Ln {line}, Col {col}</span>
      <span class="text-[var(--color-overlay0)]">⌃P</span>
    </span>
  </footer>
{/if}
