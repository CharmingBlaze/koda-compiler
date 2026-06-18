<script>

  /** @type {{ ok: boolean, version?: string, installDir?: string, stdlibDir?: string, lines?: { ok: boolean, label: string, detail: string, fix?: string }[] }} */

  let { sdk = { ok: false, lines: [] }, examples = [], onNewProject, onOpenWorkspace, onOpenExample = () => {}, onOpenHelp = () => {} } = $props()



  const templates = [

    {

      id: 'hello',

      title: 'Hello app',

      blurb: 'Print to the console — your first native program.',

    },

    {

      id: 'game',

      title: 'Text game',

      blurb: 'Lunar lander in the terminal — no graphics setup.',

    },

    {

      id: 'graphics',

      title: 'Bouncing ball',

      blurb: 'Raylib window with use raylib + koda.game — press F5 to run.',

    },

    {

      id: 'pong',

      title: 'Pong',

      blurb: 'Classic two-player paddle game — W/S vs arrow keys.',

    },

  ]



  const categoryLabels = {

    game: 'Games',

    graphics: 'Graphics demos',

    language: 'Language',

    console: 'Console',

  }



  /** @type {Record<string, typeof examples>} */

  let grouped = $derived.by(() => {

    /** @type {Record<string, typeof examples>} */

    const g = {}

    for (const ex of examples) {

      const cat = ex.category || 'other'

      if (!g[cat]) g[cat] = []

      g[cat].push(ex)

    }

    return g

  })



  const categoryOrder = ['game', 'graphics', 'console', 'language']



  /** @param {string} template */

  function start(template) {

    onNewProject?.(template)

  }

</script>



<div class="welcome">

  <header>

    <h1>Welcome to Koda Studio</h1>

    <p class="welcome-lead">

      Unzip the Koda SDK for your OS, pick a template or open an example, then press

      <kbd>F5</kbd> to run. Press <kbd>F1</kbd> for help and all documentation.

    </p>

  </header>



  <section class="welcome-card" style="margin-top: 2rem;">

    <div style="font-weight: 600; font-size: 0.875rem; margin-bottom: 0.625rem; color: var(--color-text);">SDK status</div>

    {#if sdk.lines?.length}

      <ul style="list-style: none; padding: 0; margin: 0; font-size: 0.875rem; line-height: 1.5;">

        {#each sdk.lines as line}

          <li style="display: flex; gap: 0.625rem; margin-bottom: 0.5rem; align-items: flex-start;">

            <span

              style="font-family: var(--font-mono); font-size: 0.75rem; font-weight: 600; min-width: 2.25rem; color: {line.ok

                ? 'var(--color-green)'

                : 'var(--color-red)'};"

            >

              {line.ok ? 'OK' : 'FAIL'}

            </span>

            <span style="color: var(--color-subtext);">

              <strong style="color: var(--color-text); font-weight: 500;">{line.label}</strong>

              — {line.detail}

              {#if !line.ok && line.fix}

                <span style="display: block; font-size: 0.8125rem; color: var(--color-overlay0); margin-top: 0.125rem;">{line.fix}</span>

              {/if}

            </span>

          </li>

        {/each}

      </ul>

    {:else}

      <p style="color: var(--color-overlay0); margin: 0; font-size: 0.875rem;">Checking SDK…</p>

    {/if}

  </section>



  <section style="margin-top: 2rem;">

    <h2>Start a new project</h2>

    <div class="template-grid" style="margin-top: 1rem;">

      {#each templates as t}

        <button type="button" class="template-btn" onclick={() => start(t.id)}>

          <strong>{t.title}</strong>

          <span>{t.blurb}</span>

        </button>

      {/each}

    </div>

  </section>



  <section style="margin-top: 2rem;">

    <h2>Open an example</h2>

    <p style="margin: 0.5rem 0 0; font-size: 0.875rem; color: var(--color-subtext);">

      Runnable samples from the SDK — opens the project folder; press <kbd>F5</kbd> to compile and run.

    </p>

    {#if examples.length === 0}

      <p style="margin-top: 1rem; font-size: 0.875rem; color: var(--color-overlay0);">Loading examples…</p>

    {:else}

      {#each categoryOrder as cat}

        {#if grouped[cat]?.length}

          <h3 style="margin: 1.25rem 0 0.5rem; font-size: 0.9375rem; color: var(--color-subtext);">

            {categoryLabels[cat] || cat}

          </h3>

          <div class="template-grid" style="margin-top: 0.5rem;">

            {#each grouped[cat] as d}

              <button type="button" class="template-btn" onclick={() => onOpenExample?.(d.id)}>

                <strong>{d.title}</strong>

                <span>{d.blurb}</span>

              </button>

            {/each}

          </div>

        {/if}

      {/each}

    {/if}

  </section>



  <section class="welcome-card" style="margin-top: 2rem;">

    <h2 style="margin: 0 0 0.75rem; font-size: 1rem;">Need help?</h2>

    <p style="margin: 0 0 0.75rem; font-size: 0.875rem; color: var(--color-subtext);">

      All tutorials, guides, and reference docs are built in — press <kbd>F1</kbd> anytime.

    </p>

    <div class="template-grid" style="grid-template-columns: repeat(auto-fill, minmax(10rem, 1fr));">

      <button type="button" class="template-btn" onclick={() => onOpenHelp('START_HERE.md')}>

        <strong>Start here</strong>

        <span>5-minute overview</span>

      </button>

      <button type="button" class="template-btn" onclick={() => onOpenHelp('docs/beginners-guide.md')}>

        <strong>Beginner's guide</strong>

        <span>Full walkthrough</span>

      </button>

      <button type="button" class="template-btn" onclick={() => onOpenHelp('docs/learn/README.md')}>

        <strong>Learn path</strong>

        <span>Step-by-step chapters</span>

      </button>

      <button type="button" class="template-btn" onclick={() => onOpenHelp('docs/guides/game-dev.md')}>

        <strong>Game dev</strong>

        <span>Loops, input, shipping</span>

      </button>

      <button type="button" class="template-btn" onclick={() => onOpenHelp('docs/faq.md')}>

        <strong>FAQ</strong>

        <span>Common questions</span>

      </button>

      <button type="button" class="template-btn" onclick={() => onOpenHelp()}>

        <strong>All documentation</strong>

        <span>Browse everything</span>

      </button>

    </div>

  </section>



  <div class="welcome-actions" style="margin-top: 2rem;">

    <button type="button" onclick={() => onOpenWorkspace?.()}>Open existing folder…</button>

  </div>

</div>

