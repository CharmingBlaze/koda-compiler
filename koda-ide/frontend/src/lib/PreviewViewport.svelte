<script>
  import { onMount, onDestroy } from 'svelte'

  let wrap = $state(null)
  let worker = null

  onMount(() => {
    if (!wrap) return
    const canvas = document.createElement('canvas')
    canvas.width = 420
    canvas.height = 280
    canvas.className = 'h-auto max-w-full rounded-md border border-[var(--color-surface0)] bg-[var(--color-crust)]'
    wrap.appendChild(canvas)
    const off = canvas.transferControlToOffscreen()
    worker = new Worker(new URL('./preview.worker.js', import.meta.url), { type: 'module' })
    worker.postMessage({ type: 'init', canvas: off }, [off])
  })

  onDestroy(() => {
    worker?.terminate()
    worker = null
    wrap?.replaceChildren()
  })
</script>

<div class="flex flex-col gap-2 p-2">
  <p class="text-xs text-[var(--color-subtext)]">
    OffscreenCanvas preview (worker). Swap this for Raylib child-process frames or WASM when you wire gameplay.
  </p>
  <div bind:this={wrap} class="flex justify-center"></div>
</div>
