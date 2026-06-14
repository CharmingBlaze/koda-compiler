import { HighlightStyle, syntaxHighlighting } from '@codemirror/language'
import { tags } from '@lezer/highlight'

/** Syntax colors aligned with Catppuccin Macchiato CSS variables. */
export const kodaHighlightStyle = HighlightStyle.define([
  { tag: tags.keyword, color: 'var(--color-mauve)' },
  { tag: tags.comment, color: 'var(--color-overlay0)', fontStyle: 'italic' },
  { tag: tags.string, color: 'var(--color-green)' },
  { tag: tags.number, color: 'var(--color-yellow)' },
  { tag: tags.operator, color: 'var(--color-accent)' },
  { tag: tags.variableName, color: 'var(--color-text)' },
  { tag: tags.typeName, color: 'var(--color-peach)' },
])

export const kodaSyntax = syntaxHighlighting(kodaHighlightStyle)
