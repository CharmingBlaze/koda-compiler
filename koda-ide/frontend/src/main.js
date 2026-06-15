import '@fontsource-variable/inter/index.css'
import '@fontsource-variable/jetbrains-mono/index.css'
import { mount } from 'svelte'
import './app.css'
import App from './App.svelte'

const root = document.getElementById('app')
if (root) {
  root.style.height = '100%'
  root.style.width = '100%'
  try {
    mount(App, { target: root })
    document.getElementById('boot')?.remove()
  } catch (err) {
    const msg = err instanceof Error ? err.message : String(err)
    if (typeof window.__kodaReport === 'function') {
      window.__kodaReport('Koda Studio failed to start:\n' + msg)
    } else {
      root.innerHTML = `<div style="padding:2rem;font-family:system-ui;color:#f3f6fa;background:#131920;max-width:40rem"><h1 style="color:#ff7b72;margin:0 0 1rem">Koda Studio failed to start</h1><pre style="white-space:pre-wrap;color:#a8b8c8">${msg}</pre></div>`
    }
    console.error(err)
  }
}
