import { HighlightStyle, syntaxHighlighting } from '@codemirror/language'
import { tags } from '@lezer/highlight'

/** High-contrast syntax colors for comfortable long-form coding. */
export const kodaHighlightStyle = HighlightStyle.define([
  { tag: tags.keyword, color: 'var(--color-syntax-keyword)', fontWeight: '600' },
  { tag: tags.controlKeyword, color: 'var(--color-syntax-keyword)', fontWeight: '600' },
  { tag: tags.definitionKeyword, color: 'var(--color-syntax-keyword)', fontWeight: '600' },
  { tag: tags.moduleKeyword, color: 'var(--color-syntax-keyword)', fontWeight: '600' },
  { tag: tags.self, color: 'var(--color-syntax-type)', fontWeight: '500' },
  { tag: tags.atom, color: 'var(--color-syntax-type)', fontWeight: '500' },
  { tag: tags.bool, color: 'var(--color-syntax-type)' },
  { tag: tags.null, color: 'var(--color-syntax-type)' },
  { tag: tags.comment, color: 'var(--color-syntax-comment)', fontStyle: 'italic' },
  { tag: tags.lineComment, color: 'var(--color-syntax-comment)', fontStyle: 'italic' },
  { tag: tags.blockComment, color: 'var(--color-syntax-comment)', fontStyle: 'italic' },
  { tag: tags.meta, color: 'var(--color-cyan)', fontWeight: '500' },
  { tag: tags.string, color: 'var(--color-syntax-string)' },
  { tag: tags.special(tags.string), color: 'var(--color-syntax-string)' },
  { tag: tags.number, color: 'var(--color-syntax-number)' },
  { tag: tags.integer, color: 'var(--color-syntax-number)' },
  { tag: tags.float, color: 'var(--color-syntax-number)' },
  { tag: tags.operator, color: 'var(--color-syntax-operator)' },
  { tag: tags.punctuation, color: 'var(--color-syntax-punct)' },
  { tag: tags.bracket, color: 'var(--color-syntax-punct)' },
  { tag: tags.variableName, color: 'var(--color-text)' },
  { tag: tags.propertyName, color: 'var(--color-text)' },
  { tag: tags.typeName, color: 'var(--color-syntax-type)', fontWeight: '500' },
  { tag: tags.standard(tags.variableName), color: 'var(--color-syntax-function)', fontWeight: '500' },
  { tag: tags.function(tags.variableName), color: 'var(--color-syntax-function)', fontWeight: '500' },
  { tag: tags.definition(tags.variableName), color: 'var(--color-cyan)', fontWeight: '600' },
  { tag: tags.definition(tags.typeName), color: 'var(--color-cyan)', fontWeight: '600' },
  { tag: tags.definition(tags.function(tags.variableName)), color: 'var(--color-syntax-function)', fontWeight: '600' },
  { tag: tags.namespace, color: 'var(--color-yellow)' },
  { tag: tags.className, color: 'var(--color-syntax-type)', fontWeight: '600' },
])

export const kodaSyntax = syntaxHighlighting(kodaHighlightStyle)
