/**
 * Placeholder OffscreenCanvas loop for a future Raylib / WASM preview.
 * Draws a simple animated frame so the worker + transfer path stays wired.
 */
let ctx
let W = 400
let H = 260
let tick = 0

function frame() {
  tick++
  if (!ctx) return
  ctx.fillStyle = '#11111b'
  ctx.fillRect(0, 0, W, H)
  ctx.strokeStyle = '#89b4fa'
  ctx.lineWidth = 2
  const x = 24 + (tick % 120)
  ctx.strokeRect(x, 48, 100, 72)
  ctx.fillStyle = '#a6adc8'
  ctx.font = '14px ui-monospace, monospace'
  ctx.fillText('preview worker (OffscreenCanvas)', 16, 28)
  setTimeout(frame, 33)
}

self.onmessage = (e) => {
  if (e.data?.type === 'init' && e.data.canvas) {
    const canvas = e.data.canvas
    W = canvas.width
    H = canvas.height
    ctx = canvas.getContext('2d')
    frame()
  }
}
