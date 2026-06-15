<script>
  import { onMount } from 'svelte'
  import { ListDocPages, ReadDocPage } from '../../wailsjs/go/main/App.js'
  import { BrowserOpenURL } from '../../wailsjs/runtime/runtime.js'
  import { renderMarkdown, resolveDocLink } from './markdownLite.js'
  import {
    KODA_KEYWORDS,
    KODA_TYPES,
    KODA_BUILTINS,
    KODA_STDLIB_MODULES,
    lookupKodaHelp,
  } from './kodaHelp.js'

  let {
    open = $bindable(false),
    initialPage = '',
    onClose = () => {},
  } = $props()

  /** @type {'docs' | 'reference'} */
  let tab = $state('docs')
  /** @type {{ rel: string, title: string, category: string, beginner: boolean }[]} */
  let pages = $state([])
  let loadErr = $state('')
  let activeRel = $state('START_HERE.md')
  let docHtml = $state('')
  let docLoading = $state(false)
  let filter = $state('')
  let refFilter = $state('')

  const grouped = $derived.by(() => {
    const q = filter.trim().toLowerCase()
    /** @type {Map<string, typeof pages>} */
    const map = new Map()
    for (const p of pages) {
      if (q && !p.title.toLowerCase().includes(q) && !p.rel.toLowerCase().includes(q)) continue
      const list = map.get(p.category) || []
      list.push(p)
      map.set(p.category, list)
    }
    return [...map.entries()]
  })

  const refItems = $derived.by(() => {
    const q = refFilter.trim().toLowerCase()
    /** @type {{ name: string, kind: string, detail: string }[]} */
    const items = []
    for (const [name, desc] of KODA_KEYWORDS) {
      items.push({ name, kind: 'keyword', detail: desc })
    }
    for (const [name, desc] of KODA_TYPES) {
      items.push({ name, kind: 'type', detail: desc })
    }
    for (const [name, sig, desc] of KODA_BUILTINS) {
      items.push({ name, kind: 'builtin', detail: `${sig} — ${desc}` })
    }
    for (const [name, desc] of KODA_STDLIB_MODULES) {
      items.push({ name, kind: 'module', detail: desc })
    }
    if (!q) return items
    return items.filter(
      (it) => it.name.toLowerCase().includes(q) || it.detail.toLowerCase().includes(q),
    )
  })

  async function loadCatalog() {
    loadErr = ''
    try {
      pages = await ListDocPages()
      if (pages.length === 0) {
        loadErr = 'No documentation found. Unzip the full Koda SDK so docs/ sits next to Koda Studio.'
      } else if (!docHtml && open) {
        await openPage(activeRel)
      }
    } catch (e) {
      loadErr = e instanceof Error ? e.message : String(e)
      pages = []
    }
  }

  /** @param {string} rel */
  async function openPage(rel) {
    if (!rel) return
    activeRel = rel
    tab = 'docs'
    docLoading = true
    try {
      const md = await ReadDocPage(rel)
      docHtml = renderMarkdown(md)
    } catch (e) {
      docHtml = `<p class="help-error">${e instanceof Error ? e.message : String(e)}</p>`
    } finally {
      docLoading = false
    }
  }

  /** @param {MouseEvent} e */
  function onDocClick(e) {
    const t = e.target
    if (!(t instanceof Element)) return
    const a = t.closest('a[data-doc-href]')
    if (!a) return
    e.preventDefault()
    const href = a.getAttribute('data-doc-href')
    if (!href) return
    const resolved = resolveDocLink(activeRel, href)
    if (resolved) void openPage(resolved)
  }

  /** @param {MouseEvent} e */
  function onDocClickCapture(e) {
    const t = e.target
    if (!(t instanceof Element)) return
    const a = t.closest('a[href^="http"]')
    if (!a) return
    e.preventDefault()
    const href = a.getAttribute('href')
    if (href) BrowserOpenURL(href)
  }

  /** @param {string} name */
  function openRef(name) {
    refFilter = name
    const entry = lookupKodaHelp(name)
    if (!entry) return
    /** @type {Record<string, string>} */
    const stdlibDoc = {
      '@math': 'docs/stdlib/math.md',
      '@json': 'docs/stdlib/json.md',
      '@io': 'docs/stdlib/io.md',
      '@array': 'docs/stdlib/array.md',
      '@str': 'docs/stdlib/str.md',
      '@util': 'docs/stdlib/util.md',
      '@game': 'docs/stdlib/game.md',
      '@vec2': 'docs/stdlib/vec2.md',
      '@vec3': 'docs/stdlib/vec3.md',
      '@timer': 'docs/stdlib/timer.md',
      '@noise': 'docs/stdlib/noise.md',
    }
    const doc = stdlibDoc[entry.name]
    if (doc && pages.some((p) => p.rel === doc)) {
      void openPage(doc)
      return
    }
    if (entry.kind === 'builtin' && pages.some((p) => p.rel === 'docs/reference/builtins.md')) {
      void openPage('docs/reference/builtins.md')
    }
  }

  function close() {
    open = false
    onClose()
  }

  $effect(() => {
    if (open && pages.length === 0 && !loadErr) {
      void loadCatalog()
    }
  })

  $effect(() => {
    if (open && initialPage) {
      void openPage(initialPage)
    }
  })

  onMount(() => {
    if (open) void loadCatalog()
  })
</script>

{#if open}
  <div class="help-overlay" role="dialog" aria-modal="true" aria-label="Koda help and documentation">
    <header class="help-header">
      <div class="help-header-left">
        <h2 class="help-title">Help &amp; documentation</h2>
        <div class="help-tabs">
          <button type="button" class:active={tab === 'docs'} onclick={() => (tab = 'docs')}>Documentation</button>
          <button type="button" class:active={tab === 'reference'} onclick={() => (tab = 'reference')}>Language quick reference</button>
        </div>
      </div>
      <button type="button" class="help-close" onclick={close} aria-label="Close help">✕</button>
    </header>

    <div class="help-body">
      {#if tab === 'docs'}
        <aside class="help-nav">
          <div class="help-quick">
            <div class="help-nav-label">New to Koda?</div>
            <button type="button" class="help-quick-btn" onclick={() => openPage('START_HERE.md')}>Start here</button>
            <button type="button" class="help-quick-btn" onclick={() => openPage('docs/beginners-guide.md')}>Beginner's guide</button>
            <button type="button" class="help-quick-btn" onclick={() => openPage('docs/learn/README.md')}>Learn path (chapters)</button>
            <button type="button" class="help-quick-btn" onclick={() => openPage('docs/faq.md')}>FAQ</button>
            <button type="button" class="help-quick-btn" onclick={() => openPage('docs/troubleshooting.md')}>Troubleshooting</button>
            <button type="button" class="help-quick-btn" onclick={() => openPage('docs/guides/game-dev.md')}>Game development</button>
            <button type="button" class="help-quick-btn" onclick={() => openPage('language.md')}>Language reference</button>
          </div>

          <input class="help-search" type="search" placeholder="Search all docs…" bind:value={filter} />

          {#if loadErr}
            <p class="help-error">{loadErr}</p>
          {:else if pages.length === 0}
            <p class="help-muted">Loading documentation…</p>
          {:else}
            {#each grouped as [category, items] (category)}
              <div class="help-group">
                <div class="help-nav-label">{category}</div>
                <ul class="help-list">
                  {#each items as page (page.rel)}
                    <li>
                      <button
                        type="button"
                        class="help-link"
                        class:active={page.rel === activeRel}
                        onclick={() => openPage(page.rel)}>{page.title}</button>
                    </li>
                  {/each}
                </ul>
              </div>
            {/each}
          {/if}
        </aside>

        <div
          class="help-content prose-koda"
          role="document"
          onkeydown={(e) => e.key === 'Enter' && onDocClick(e)}
          onclick={onDocClick}
          onclickcapture={onDocClickCapture}
        >
          {#if docLoading}
            <p class="help-muted">Loading…</p>
          {:else}
            {@html docHtml}
          {/if}
        </div>
      {:else}
        <aside class="help-nav help-nav-wide">
          <input class="help-search" type="search" placeholder="Search keywords, builtins, @modules…" bind:value={refFilter} />
          <p class="help-muted help-ref-hint">Hover symbols in the editor for inline help. Click an item for related docs.</p>
          <ul class="help-ref-list">
            {#each refItems as item (item.name + item.kind)}
              <li>
                <button type="button" class="help-ref-item" onclick={() => openRef(item.name)}>
                  <span class="help-ref-kind">{item.kind}</span>
                  <span class="help-ref-name">{item.name}</span>
                  <span class="help-ref-detail">{item.detail}</span>
                </button>
              </li>
            {/each}
          </ul>
        </aside>
        <article class="help-content help-content-ref">
          <h3>Language quick reference</h3>
          <p>Type in the editor and use autocomplete for snippets. Hover any word for a short explanation.</p>
          <p>Studio shortcuts: <kbd>F5</kbd> run · <kbd>Ctrl</kbd>+<kbd>S</kbd> save · <kbd>Ctrl</kbd>+<kbd>P</kbd> palette · <kbd>F1</kbd> help</p>
          {#if refFilter}
            {@const entry = lookupKodaHelp(refFilter)}
            {#if entry}
              <div class="help-ref-card">
                <div class="help-ref-card-title">{entry.signature || entry.name}</div>
                <p>{entry.desc}</p>
                {#if entry.example}
                  <pre><code>{entry.example}</code></pre>
                {/if}
              </div>
            {/if}
          {/if}
        </article>
      {/if}
    </div>
  </div>
{/if}

<style>
  .help-overlay {
    position: fixed;
    inset: 2rem 1rem 2.5rem 3rem;
    z-index: 70;
    display: flex;
    flex-direction: column;
    border-radius: 12px;
    border: 1px solid var(--color-surface1);
    background: var(--color-base);
    box-shadow: 0 24px 80px rgba(0, 0, 0, 0.55);
  }
  .help-header {
    display: flex;
    align-items: flex-start;
    justify-content: space-between;
    gap: 1rem;
    padding: 0.75rem 1rem;
    border-bottom: 1px solid var(--color-surface0);
    flex-shrink: 0;
  }
  .help-header-left { display: flex; flex-wrap: wrap; align-items: center; gap: 1rem; }
  .help-title { margin: 0; font-size: 1rem; font-weight: 600; color: var(--color-text); }
  .help-tabs { display: flex; gap: 0.25rem; }
  .help-tabs button {
    border-radius: 6px; border: 1px solid transparent; padding: 0.25rem 0.65rem;
    font-size: 0.8125rem; color: var(--color-subtext); background: transparent; cursor: pointer;
  }
  .help-tabs button.active {
    border-color: var(--color-accent); color: var(--color-text); background: var(--color-surface0);
  }
  .help-close {
    border: none; background: transparent; color: var(--color-overlay0);
    font-size: 1.1rem; cursor: pointer; padding: 0.25rem 0.5rem; border-radius: 6px;
  }
  .help-close:hover { background: var(--color-surface0); color: var(--color-text); }
  .help-body { display: flex; min-height: 0; flex: 1; }
  .help-nav {
    width: min(280px, 34vw); flex-shrink: 0; overflow-y: auto;
    border-right: 1px solid var(--color-surface0); padding: 0.75rem;
  }
  .help-nav-wide { width: min(360px, 40vw); }
  .help-nav-label {
    font-size: 0.6875rem; font-weight: 600; letter-spacing: 0.06em; text-transform: uppercase;
    color: var(--color-overlay0); margin: 0.75rem 0 0.35rem;
  }
  .help-nav-label:first-child { margin-top: 0; }
  .help-quick { margin-bottom: 0.75rem; }
  .help-quick-btn {
    display: block; width: 100%; text-align: left; border: none; border-radius: 6px;
    padding: 0.35rem 0.5rem; margin-bottom: 0.15rem; font-size: 0.8125rem; color: var(--color-cyan);
    background: color-mix(in srgb, var(--color-cyan) 8%, transparent); cursor: pointer;
  }
  .help-quick-btn:hover { background: color-mix(in srgb, var(--color-cyan) 16%, transparent); }
  .help-search {
    width: 100%; box-sizing: border-box; margin-bottom: 0.5rem; border-radius: 6px;
    border: 1px solid var(--color-surface0); background: var(--color-crust);
    padding: 0.4rem 0.55rem; font-size: 0.8125rem; color: var(--color-text); outline: none;
  }
  .help-search:focus { border-color: var(--color-accent); }
  .help-list { list-style: none; margin: 0; padding: 0; }
  .help-link {
    display: block; width: 100%; text-align: left; border: none; border-radius: 4px;
    padding: 0.2rem 0.45rem; font-size: 0.8125rem; color: var(--color-subtext);
    background: transparent; cursor: pointer;
  }
  .help-link:hover, .help-link.active { color: var(--color-text); background: var(--color-surface0); }
  .help-link.active { color: var(--color-accent); }
  .help-content {
    flex: 1; overflow-y: auto; padding: 1rem 1.25rem 2rem;
    font-size: 0.9375rem; line-height: 1.65; color: var(--color-subtext);
  }
  .help-muted { color: var(--color-overlay0); font-size: 0.8125rem; }
  .help-error { color: var(--color-red); font-size: 0.8125rem; }
  .help-ref-list { list-style: none; margin: 0; padding: 0; }
  .help-ref-item {
    display: grid; grid-template-columns: 4.5rem 1fr; gap: 0.15rem 0.5rem; width: 100%;
    text-align: left; border: none; border-radius: 6px; padding: 0.35rem 0.45rem;
    margin-bottom: 0.15rem; background: transparent; cursor: pointer;
  }
  .help-ref-item:hover { background: var(--color-surface0); }
  .help-ref-kind {
    grid-row: span 2; font-size: 0.625rem; font-weight: 600; text-transform: uppercase;
    letter-spacing: 0.04em; color: var(--color-overlay0); align-self: center;
  }
  .help-ref-name { font-family: var(--font-mono); font-size: 0.8125rem; color: var(--color-cyan); }
  .help-ref-detail { grid-column: 2; font-size: 0.75rem; color: var(--color-overlay0); line-height: 1.35; }
  .help-ref-hint { margin: 0 0 0.75rem; }
  .help-ref-card {
    margin-top: 1.5rem; padding: 1rem; border-radius: 8px;
    border: 1px solid var(--color-surface0); background: var(--color-mantle);
  }
  .help-ref-card-title {
    font-family: var(--font-mono); font-weight: 600; color: var(--color-cyan); margin-bottom: 0.5rem;
  }
  .help-content-ref h3 { color: var(--color-text); margin-top: 0; }
  .help-content-ref kbd {
    border-radius: 4px; padding: 0.1rem 0.35rem; font-size: 0.75rem;
    background: var(--color-surface0); color: var(--color-text);
  }
  .help-content :global(h1), .help-content :global(h2), .help-content :global(h3) {
    color: var(--color-text); margin: 1.25rem 0 0.5rem; line-height: 1.3;
  }
  .help-content :global(h1) { font-size: 1.35rem; }
  .help-content :global(h2) { font-size: 1.1rem; }
  .help-content :global(h3) { font-size: 1rem; }
  .help-content :global(p) { margin: 0.65rem 0; }
  .help-content :global(a) { color: var(--color-cyan); }
  .help-content :global(code) {
    font-family: var(--font-mono); font-size: 0.85em; background: var(--color-surface0);
    padding: 0.1rem 0.3rem; border-radius: 4px;
  }
  .help-content :global(pre) {
    overflow-x: auto; border-radius: 8px; border: 1px solid var(--color-surface0);
    background: var(--color-crust); padding: 0.75rem 1rem; margin: 0.75rem 0;
  }
  .help-content :global(pre code) { background: none; padding: 0; }
  .help-content :global(ul), .help-content :global(ol) { margin: 0.5rem 0; padding-left: 1.35rem; }
  .help-content :global(blockquote) {
    margin: 0.75rem 0; padding-left: 0.75rem; border-left: 3px solid var(--color-accent);
    color: var(--color-overlay0);
  }
  .help-content :global(table) {
    width: 100%; border-collapse: collapse; margin: 0.75rem 0; font-size: 0.875rem;
  }
  .help-content :global(th), .help-content :global(td) {
    border: 1px solid var(--color-surface0); padding: 0.35rem 0.5rem; text-align: left;
  }
  .help-content :global(th) { background: var(--color-mantle); color: var(--color-text); }
</style>
