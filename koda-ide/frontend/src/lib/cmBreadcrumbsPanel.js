import { showPanel } from '@codemirror/view'

/**
 * Top panel with workspace-relative breadcrumb path.
 * @param {string} relPath e.g. "internal/foo.koda"
 */
export function breadcrumbsPanel(relPath) {
  const path = relPath || ''
  const text = path ? path.replace(/\\/g, '/').split('/').join(' › ') : '—'
  return showPanel.of(() => ({
    top: true,
    dom: (() => {
      const dom = document.createElement('div')
      dom.className = 'cm-koda-breadcrumbs'
      dom.textContent = text
      return dom
    })(),
  }))
}
