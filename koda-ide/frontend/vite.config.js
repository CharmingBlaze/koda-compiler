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

/** Wails WebView: no crossorigin, no ES modules — use a classic script tag at end of body. */
function wailsHtmlPlugin() {
  return {
    name: 'wails-html',
    transformIndexHtml(html) {
      let out = html
        .replace(/\s+crossorigin(?:="[^"]*")?/gi, '')
        .replace(/<script type="module" src="\.\/assets\/studio\.js"><\/script>\s*/i, '')
      if (!out.includes('assets/studio.js')) {
        out = out.replace('</body>', '  <script src="./assets/studio.js"></script>\n</body>')
      }
      return out
    },
  }
}

/** Surface script load / runtime failures in the WebView (no devtools needed). */
function bootDiagnosticsPlugin() {
  return {
    name: 'boot-diagnostics',
    transformIndexHtml(html) {
      const snippet = `<script>
(function () {
  function report(msg) {
    var b = document.getElementById('boot');
    if (b) {
      b.style.color = '#ff7b72';
      b.style.whiteSpace = 'pre-wrap';
      b.textContent = msg;
    }
  }
  window.__kodaReport = report;
  window.addEventListener('error', function (e) {
    report('Error: ' + (e.message || e.error || 'unknown') + (e.filename ? '\\n' + e.filename : ''));
  });
  window.addEventListener('unhandledrejection', function (e) {
    var r = e.reason;
    report('Promise rejection: ' + (r && r.message ? r.message : String(r)));
  });
})();
</script>`
      return html.replace('</head>', `${snippet}\n</head>`)
    },
  }
}

export default defineConfig({
  base: './',
  plugins: [faviconIcoProbe(), wailsHtmlPlugin(), bootDiagnosticsPlugin(), svelte(), tailwindcss()],
  build: {
    modulePreload: { polyfill: false },
    cssCodeSplit: false,
    rollupOptions: {
      output: {
        format: 'iife',
        inlineDynamicImports: true,
        entryFileNames: 'assets/studio.js',
        assetFileNames: 'assets/[name][extname]',
      },
    },
  },
  optimizeDeps: {
    noDiscovery: true,
    include: [],
  },
})
