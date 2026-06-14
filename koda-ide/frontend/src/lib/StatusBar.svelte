<script>
  let {
    zen = false,
    scope = '',
    diagErrors = 0,
    diagWarnings = 0,
  } = $props()
</script>

{#if !zen}
  <footer
    class="void-status flex h-[22px] shrink-0 items-center gap-3 border-t border-[var(--color-surface0)] bg-[var(--color-crust)] px-2 text-[10px] text-[var(--color-overlay0)]"
  >
    <span class="rounded-full bg-[var(--color-green)]/15 px-2 py-0.5 font-medium text-[var(--color-green)]">Koda v0.1.0</span>
    <span class="min-w-0 flex-1 truncate text-[var(--color-subtext)]" title={scope}>{scope || '—'}</span>
    <span class="flex shrink-0 items-center gap-2">
      {#if diagErrors > 0}
        <span class="flex items-center gap-0.5 text-[var(--color-diag-error)]" title="Errors"
          ><span class="inline-block h-1.5 w-1.5 rounded-full bg-[var(--color-diag-error)]"></span>{diagErrors}</span>
      {/if}
      {#if diagWarnings > 0}
        <span class="flex items-center gap-0.5 text-[var(--color-diag-warning)]" title="Warnings"
          ><span class="inline-block h-1.5 w-1.5 rounded-full bg-[var(--color-diag-warning)]"></span>{diagWarnings}</span>
      {/if}
      <span class="text-[var(--color-overlay0)]">⌃P</span>
    </span>
  </footer>
{/if}
