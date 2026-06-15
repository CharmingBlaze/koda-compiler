import { StreamLanguage } from '@codemirror/language'
import { simpleMode } from '@codemirror/legacy-modes/mode/simple-mode'
import { snippetCompletion } from '@codemirror/autocomplete'
import { hoverTooltip } from '@codemirror/view'
import {
  KODA_BUILTINS,
  KODA_KEYWORDS,
  KODA_SNIPPETS,
  KODA_STDLIB_MODULES,
  KODA_TYPES,
  lookupKodaHelp,
  renderHelpDom,
  scanDocumentSymbols,
} from './kodaHelp.js'

const keywordNames = KODA_KEYWORDS.map(([k]) => k).join('|')
const typeNames = KODA_TYPES.map(([k]) => k).join('|')
const builtinNames = KODA_BUILTINS.map(([k]) => k).join('|')

const keywordRe = new RegExp(`(?:${keywordNames})\\b`, 'i')
const typeRe = new RegExp(`(?:${typeNames})\\b`, 'i')
const builtinRe = new RegExp(`(?:${builtinNames})\\b`, 'i')

const sharedExprRules = [
  { regex: /\/\/.*/, token: 'comment' },
  { regex: /\/\*/, token: 'comment', next: 'commentBlock' },
  { regex: /#\s*include\b/i, token: 'meta' },
  { regex: keywordRe, token: 'keyword' },
  { regex: typeRe, token: 'type' },
  { regex: /\b(?:true|false|null)\b/i, token: 'atom' },
  { regex: builtinRe, token: 'builtin' },
  { regex: /0x[a-f\d]+|(?:\d+\.?\d*|\.\d+)(?:e[-+]?\d+)?/i, token: 'number' },
  { regex: /"(?:[^\\"]|\\.)*?(?:"|$)/, token: 'string' },
  { regex: /'(?:[^\\']|\\.)*?(?:'|$)/, token: 'string' },
  { regex: /`/, token: 'string', next: 'template' },
  { regex: /[\{\}\(\)\[\];,.:+\-*/%=<>!&|^~?]+/, token: 'operator' },
  { regex: /[a-zA-Z_]\w*/, token: 'variable' },
]

/** Rich stream-mode highlighting aligned with the Koda lexer. */
export const kodaLanguage = StreamLanguage.define(
  simpleMode({
    start: [
      { regex: /\/\/.*/, token: 'comment' },
      { regex: /\/\*/, token: 'comment', next: 'commentBlock' },
      { regex: /#\s*include\b/i, token: 'meta', next: 'includePath' },
      { regex: /\bfunc\b/i, token: 'keyword', next: 'funcName' },
      { regex: /\bstruct\b/i, token: 'keyword', next: 'typeName' },
      { regex: /\benum\b/i, token: 'keyword', next: 'typeName' },
      { regex: keywordRe, token: 'keyword' },
      { regex: typeRe, token: 'type' },
      { regex: /\b(?:true|false|null)\b/i, token: 'atom' },
      { regex: builtinRe, token: 'builtin' },
      { regex: /0x[a-f\d]+|(?:\d+\.?\d*|\.\d+)(?:e[-+]?\d+)?/i, token: 'number' },
      { regex: /"(?:[^\\"]|\\.)*?(?:"|$)/, token: 'string' },
      { regex: /'(?:[^\\']|\\.)*?(?:'|$)/, token: 'string' },
      { regex: /`/, token: 'string', next: 'template' },
      { regex: /[\{\}\(\)\[\];,.:+\-*/%=<>!&|^~?]+/, token: 'operator' },
      { regex: /[a-zA-Z_]\w*/, token: 'variable' },
    ],
    funcName: [
      { regex: /\s+/, token: null },
      { regex: /[a-zA-Z_]\w*/, token: 'def', next: 'start' },
      { regex: /./, token: null, next: 'start' },
    ],
    typeName: [
      { regex: /\s+/, token: null },
      { regex: /[a-zA-Z_]\w*/, token: 'def', next: 'start' },
      { regex: /./, token: null, next: 'start' },
    ],
    includePath: [
      { regex: /\s+/, token: null },
      { regex: /"(?:[^\\"]|\\.)*?(?:"|$)/, token: 'string', next: 'start' },
      { regex: /[a-zA-Z_@][\w./-]*/, token: 'string', next: 'start' },
      { regex: /./, token: null, next: 'start' },
    ],
    commentBlock: [
      { regex: /.*?\*\//, token: 'comment', next: 'start' },
      { regex: /.*/, token: 'comment' },
    ],
    template: [
      { regex: /\\./, token: 'string' },
      { regex: /\$\{/, token: 'operator', next: 'templateEmbed' },
      { regex: /[^`\\$]+/, token: 'string' },
      { regex: /`/, token: 'string', next: 'start' },
    ],
    templateEmbed: [
      { regex: /\}/, token: 'operator', next: 'template' },
      ...sharedExprRules.map((r) => ({ ...r, next: r.next === 'commentBlock' ? 'commentBlock' : r.next })),
    ],
    languageData: {
      closeBrackets: { brackets: ['(', '[', '{', "'", '"', '`'] },
      commentTokens: { line: '//', block: { open: '/*', close: '*/' } },
    },
  }),
)

function helpInfo(name) {
  const entry = lookupKodaHelp(name)
  if (!entry) return undefined
  return () => renderHelpDom(entry)
}

function completionFromKeyword([name, desc, example]) {
  return {
    label: name,
    type: 'keyword',
    detail: desc,
    info: helpInfo(name),
  }
}

function completionFromBuiltin([name, sig, desc]) {
  return {
    label: name,
    type: 'function',
    detail: sig,
    info: helpInfo(name),
  }
}

function completionFromType([name, desc]) {
  return {
    label: name,
    type: 'type',
    detail: desc,
    info: helpInfo(name),
  }
}

function completionFromModule([name, desc]) {
  return {
    label: name,
    type: 'namespace',
    detail: desc,
    info: helpInfo(name),
  }
}

const staticCompletions = [
  ...KODA_SNIPPETS.map((s) =>
    snippetCompletion(s.apply, {
      label: s.label,
      type: s.type,
      detail: s.detail,
      info: s.info,
    }),
  ),
  ...KODA_KEYWORDS.map(completionFromKeyword),
  ...KODA_TYPES.map(completionFromType),
  ...KODA_BUILTINS.map(completionFromBuiltin),
  ...KODA_STDLIB_MODULES.map(completionFromModule),
]

/** Context-aware completions: keywords, builtins, snippets, stdlib, and symbols from the file. */
export function kodaCompletionSource(context) {
  const before = context.matchBefore(/[#@]?[\w.]*/)
  if (!before && !context.explicit) return null
  if (before && before.from === before.to && !context.explicit) return null

  const from = before ? before.from : context.pos
  const prefix = (before?.text || '').toLowerCase()
  const docText = context.state.doc.toString()
  const local = scanDocumentSymbols(docText)
  const options = [...staticCompletions, ...local]

  const filtered = prefix
    ? options.filter((o) => {
        const label = String(o.label || '').toLowerCase()
        return label.startsWith(prefix) || label.includes(prefix)
      })
    : options

  const deduped = []
  const seen = new Set()
  for (const opt of filtered) {
    const key = String(opt.label).toLowerCase()
    if (seen.has(key)) continue
    seen.add(key)
    deduped.push(opt)
  }

  deduped.sort((a, b) => {
    const score = (o) => {
      const l = String(o.label).toLowerCase()
      if (l === prefix) return 0
      if (l.startsWith(prefix)) return 1
      if (o.type === 'snippet') return 2
      if (o.type === 'function') return 3
      return 4
    }
    return score(a) - score(b) || String(a.label).localeCompare(String(b.label))
  })

  return {
    from,
    options: deduped.slice(0, 40),
    validFor: /^[#@]?[\w.]*$/,
  }
}

/** Hover docs for keywords, builtins, types, and local symbols. */
export const kodaHover = hoverTooltip((view, pos) => {
  const word = view.state.wordAt(pos)
  if (!word) return null
  const text = view.state.doc.sliceString(word.from, word.to)
  let entry = lookupKodaHelp(text)

  if (!entry) {
    const symbols = scanDocumentSymbols(view.state.doc.toString())
    const sym = symbols.find((s) => s.label.toLowerCase() === text.toLowerCase())
    if (sym) {
      entry = { kind: sym.type, name: sym.label, signature: sym.label, desc: sym.detail || '' }
    }
  }

  if (!entry) return null

  return {
    pos: word.from,
    end: word.end,
    above: true,
    create() {
      return { dom: renderHelpDom(entry) }
    },
  }
})

/** @deprecated use kodaCompletionSource */
export const kodaCompletions = kodaCompletionSource
