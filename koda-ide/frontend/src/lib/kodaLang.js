import { StreamLanguage } from '@codemirror/language'
import { simpleMode } from '@codemirror/legacy-modes/mode/simple-mode'

/** Phase-A Koda highlighting via legacy stream mode (keywords, strings, comments). */
export const kodaLanguage = StreamLanguage.define(
  simpleMode({
    start: [
      { regex: /\/\/.*/, token: 'comment' },
      { regex: /\/\*/, token: 'comment', next: 'commentBlock' },
      {
        regex:
          /(?:func|let|const|var|if|else|while|for|return|import|from|export|true|false|null|break|continue|switch|case|default|class|try|catch|finally|throw|typeof|in|of|do|this|include)\b/,
        token: 'keyword',
      },
      { regex: /0x[a-f\d]+|(?:\d+\.?\d*|\.\d+)(?:e[-+]?\d+)?/i, token: 'number' },
      { regex: /"(?:[^\\"]|\\.)*?(?:"|$)/, token: 'string' },
      { regex: /'(?:[^\\']|\\.)*?(?:'|$)/, token: 'string' },
      { regex: /[\{\}\(\)\[\];,.:+\-*/%=<>!&|^~?]+/, token: 'operator' },
      { regex: /[a-zA-Z_]\w*/, token: 'variable' },
    ],
    commentBlock: [
      { regex: /.*?\*\//, token: 'comment', next: 'start' },
      { regex: /.*/, token: 'comment' },
    ],
    languageData: { closeBrackets: { brackets: ['(', '[', '{', "'", '"'] } },
  }),
)
