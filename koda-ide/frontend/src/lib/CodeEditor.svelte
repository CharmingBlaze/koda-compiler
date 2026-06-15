<script>
  import { onMount, onDestroy, tick } from 'svelte'
  import { EditorState, Prec } from '@codemirror/state'
  import {
    EditorView,
    drawSelection,
    highlightActiveLine,
    highlightActiveLineGutter,
    keymap,
    lineNumbers,
  } from '@codemirror/view'
  import { defaultKeymap, history, historyKeymap, indentWithTab } from '@codemirror/commands'
  import { bracketMatching, foldGutter, indentOnInput, indentUnit } from '@codemirror/language'
  import { closeBrackets, autocompletion, closeBracketsKeymap, completionKeymap } from '@codemirror/autocomplete'
  import { searchKeymap, highlightSelectionMatches } from '@codemirror/search'
  import { linter } from '@codemirror/lint'
  import { DiagnoseFile } from '../../wailsjs/go/main/App.js'
  import { kodaCompletionSource, kodaHover, kodaLanguage } from './kodaLang.js'
  import { kodaSyntax } from './cmKodaTheme.js'
  import { breadcrumbsPanel } from './cmBreadcrumbsPanel.js'

  let {
    absPath = '',
    relPath = '',
    seed = '',
    jumpTarget = null,
    onTextChange = () => {},
    onCursorChange = () => {},
    onLspNotify = () => {},
    onSave = () => {},
    onJumpConsumed = () => {},
    tabSize = 4,
  } = $props()

  let host = $state(null)
  let view = $state(null)

  function posForLineCol(doc, line, col) {
    const lineN = Number(line)
    const colN = Number(col)
    const safeLine = Number.isFinite(lineN) && lineN >= 1 ? lineN : 1
    const safeCol = Number.isFinite(colN) && colN >= 1 ? colN : 1
    const ln = Math.min(Math.max(safeLine, 1), doc.lines)
    const L = doc.line(ln)
    const ch = Math.min(Math.max(safeCol, 1), L.length + 1) - 1
    return L.from + ch
  }

  const voidCmTheme = EditorView.theme({
    '&': {
      backgroundColor: 'var(--color-editor-bg)',
      color: 'var(--color-text)',
    },
    '.cm-content': {
      caretColor: 'var(--color-accent)',
      padding: '12px 8px',
      fontSize: 'var(--editor-font-size)',
      lineHeight: 'var(--editor-line-height)',
    },
    '.cm-activeLine': {
      backgroundColor: 'var(--color-editor-active-line) !important',
    },
    '.cm-selectionBackground, &.cm-focused .cm-selectionBackground': {
      backgroundColor: 'var(--color-editor-selection) !important',
    },
    '&.cm-focused .cm-cursor': {
      borderLeftColor: 'var(--color-accent) !important',
    },
    '.cm-cursor, .cm-dropCursor': {
      borderLeft: '2.5px solid var(--color-accent) !important',
      marginLeft: '-1px',
    },
    '.cm-lintRange.cm-lintRange-error': {
      background: 'transparent',
      textDecoration: 'underline wavy var(--color-diag-error)',
      textDecorationThickness: '1.2px',
      textUnderlineOffset: '2px',
    },
    '.cm-lintRange.cm-lintRange-warning': {
      background: 'transparent',
      textDecoration: 'underline wavy var(--color-diag-warning)',
      textDecorationThickness: '1.2px',
      textUnderlineOffset: '2px',
    },
    '.cm-lintRange.cm-lintRange-info': {
      background: 'transparent',
      textDecoration: 'underline wavy var(--color-diag-info)',
      textDecorationThickness: '1.2px',
      textUnderlineOffset: '2px',
    },
    '.cm-tooltip.cm-tooltip-lint': {
      maxWidth: '28rem',
      background: 'var(--glass-bg)',
      backdropFilter: 'blur(var(--glass-blur))',
      border: '1px solid var(--glass-border)',
      borderRadius: '8px',
      color: 'var(--color-text)',
      boxShadow: '0 12px 40px color-mix(in srgb, var(--color-crust) 55%, transparent)',
    },
    '.cm-diagnosticAction': {
      color: 'var(--color-accent)',
    },
  })

  function makeExtensions(pathForLint, relForCrumb) {
    const safeTabSize = [2, 4, 8].includes(Number(tabSize)) ? Number(tabSize) : 4
    return [
      EditorState.tabSize.of(safeTabSize),
      indentUnit.of(' '.repeat(safeTabSize)),
      breadcrumbsPanel(relForCrumb),
      lineNumbers(),
      highlightActiveLineGutter(),
      foldGutter(),
      drawSelection(),
      EditorState.allowMultipleSelections.of(true),
      indentOnInput(),
      bracketMatching(),
      closeBrackets(),
      autocompletion({
        override: [kodaCompletionSource],
        activateOnTyping: true,
        maxRenderedOptions: 25,
        defaultKeymap: true,
        icons: true,
        closeOnBlur: true,
      }),
      kodaHover,
      highlightActiveLine(),
      highlightSelectionMatches(),
      history(),
      kodaLanguage,
      kodaSyntax,
      voidCmTheme,
      Prec.highest(
        keymap.of([
          {
            key: 'Mod-s',
            run: () => {
              onSave()
              return true
            },
          },
        ]),
      ),
      keymap.of([
        ...closeBracketsKeymap,
        ...defaultKeymap,
        ...searchKeymap,
        ...historyKeymap,
        ...completionKeymap,
        indentWithTab,
      ]),
      linter(
        async (v) => {
          if (!pathForLint) return []
          let rows = []
          try {
            const raw = await DiagnoseFile(pathForLint, v.state.doc.toString())
            rows = Array.isArray(raw) ? raw : []
          } catch {
            return []
          }
          if (!rows.length) return []
          const out = []
          for (const d of rows) {
            try {
              const from0 = posForLineCol(v.state.doc, d.line, d.col)
              if (!Number.isFinite(from0) || from0 < 0 || from0 > v.state.doc.length) continue
              const w = v.state.wordAt(from0)
              const from = w ? w.from : from0
              const to = w ? w.to : Math.min(from0 + 1, v.state.doc.length)
              const sev = d.severity === 'warning' ? 'warning' : d.severity === 'info' ? 'info' : 'error'
              const msg = d.message != null ? String(d.message) : ''
              out.push({
                from,
                to,
                message: msg,
                severity: sev,
                actions: [
                  {
                    name: 'Copy message',
                    apply() {
                      void navigator.clipboard?.writeText(msg)
                    },
                  },
                ],
              })
            } catch {
              /* skip malformed diagnostic */
            }
          }
          return out
        },
        { delay: 150 },
      ),
      EditorView.updateListener.of((u) => {
        if (u.docChanged) {
          onTextChange(u.state.doc.toString())
          onLspNotify(pathForLint, u.state.doc.toString())
        }
        if (u.selectionSet || u.docChanged) {
          const head = u.state.selection.main.head
          const line = u.state.doc.lineAt(head)
          onCursorChange({ line: line.number, col: head - line.from + 1 })
        }
      }),
    ]
  }

  function mountEditor() {
    if (!host) return
    if (view) {
      view.destroy()
      view = null
    }
    const state = EditorState.create({
      doc: seed,
      extensions: makeExtensions(absPath, relPath),
    })
    view = new EditorView({ state, parent: host })
    onCursorChange({ line: 1, col: 1 })
  }

  $effect(() => {
    if (!view || !jumpTarget) return
    const { line, col } = jumpTarget
    queueMicrotask(() => {
      if (!view) return
      const pos = posForLineCol(view.state.doc, line, col)
      view.dispatch({
        selection: { anchor: pos, head: pos },
        scrollIntoView: true,
      })
      onJumpConsumed()
    })
  })

  onMount(() => {
    void tick().then(() => mountEditor())
  })

  onDestroy(() => {
    if (view) {
      view.destroy()
      view = null
    }
  })
</script>

<div
  class="h-full min-h-[200px] w-full overflow-hidden rounded-lg border border-[var(--color-surface0)] bg-[var(--color-editor-bg)] shadow-inner"
  bind:this={host}
></div>
