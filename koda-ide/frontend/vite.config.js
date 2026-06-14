import fs from 'node:fs'
import path from 'node:path'
import { fileURLToPath } from 'node:url'
import { defineConfig } from 'vite'
import { svelte } from '@sveltejs/vite-plugin-svelte'
import tailwindcss from '@tailwindcss/vite'

const __dirname = path.dirname(fileURLToPath(import.meta.url))

/** Browsers still probe /favicon.ico; serve the SVG with the correct MIME type. */
function faviconIcoProbe() {
  return {
    name: 'favicon-ico-probe',
    configureServer(server) {
      server.middlewares.use((req, res, next) => {
        const url = req.url?.split('?')[0]
        if (url !== '/favicon.ico') return next()
        const fp = path.join(__dirname, 'public', 'favicon.svg')
        res.setHeader('Content-Type', 'image/svg+xml; charset=utf-8')
        res.end(fs.readFileSync(fp))
      })
    },
    configurePreviewServer(server) {
      server.middlewares.use((req, res, next) => {
        const url = req.url?.split('?')[0]
        if (url !== '/favicon.ico') return next()
        const fp = path.join(__dirname, 'dist', 'favicon.svg')
        res.setHeader('Content-Type', 'image/svg+xml; charset=utf-8')
        res.end(fs.readFileSync(fp))
      })
    },
    closeBundle() {
      const dist = path.join(__dirname, 'dist')
      const pub = path.join(__dirname, 'public', 'favicon.svg')
      if (!fs.existsSync(dist) || !fs.existsSync(pub)) return
      fs.copyFileSync(pub, path.join(dist, 'favicon.ico'))
    },
  }
}

export default defineConfig({
  plugins: [faviconIcoProbe(), svelte(), tailwindcss()],
})
