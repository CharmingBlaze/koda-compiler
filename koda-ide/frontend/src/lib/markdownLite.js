/** Minimal markdown → HTML for in-IDE documentation (no external deps). */

function escapeHtml(s) {
  return s
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
}

function escapeAttr(s) {
  return escapeHtml(s).replace(/'/g, '&#39;')
}

/** @param {string} text */
function inlineFormat(text) {
  let s = escapeHtml(text)
  s = s.replace(/`([^`]+)`/g, '<code>$1</code>')
  s = s.replace(/\*\*([^*]+)\*\*/g, '<strong>$1</strong>')
  s = s.replace(/\*([^*]+)\*/g, '<em>$1</em>')
  s = s.replace(/\[([^\]]+)\]\(([^)]+)\)/g, (_, label, href) => {
    const h = href.trim()
    if (/^https?:\/\//i.test(h) || h.startsWith('#')) {
      return `<a href="${escapeHtml(h)}"${/^https?:/i.test(h) ? ' target="_blank" rel="noopener noreferrer"' : ''}>${escapeHtml(label)}</a>`
    }
    return `<a href="#" data-doc-href="${escapeHtml(h)}">${escapeHtml(label)}</a>`
  })
  return s
}

/**
 * @param {string} md
 * @returns {string}
 */
export function renderMarkdown(md) {
  if (!md) return ''
  const lines = md.replace(/\r\n/g, '\n').split('\n')
  /** @type {string[]} */
  const out = []
  let i = 0

  while (i < lines.length) {
    const line = lines[i]

    if (/^```/.test(line)) {
      const lang = line.slice(3).trim()
      i++
      const code = []
      while (i < lines.length && !/^```/.test(lines[i])) {
        code.push(lines[i])
        i++
      }
      i++
      const cls = lang ? ` class="language-${escapeHtml(lang)}"` : ''
      const rawCode = code.join('\n')
      out.push(
        `<div class="code-copy-wrap"><button type="button" class="copy-button copy-button-compact" data-copy-text="${escapeAttr(rawCode)}" title="Copy code" aria-label="Copy code"><span aria-hidden="true">⧉</span><span>Copy</span></button><pre><code${cls}>${escapeHtml(rawCode)}</code></pre></div>`,
      )
      continue
    }

    if (/^#{1,3}\s/.test(line)) {
      const level = line.match(/^#+/)?.[0].length ?? 1
      const tag = `h${Math.min(level, 3)}`
      out.push(`<${tag}>${inlineFormat(line.replace(/^#+\s*/, ''))}</${tag}>`)
      i++
      continue
    }

    if (/^---+$/.test(line.trim()) || /^\*\*\*+$/.test(line.trim())) {
      out.push('<hr />')
      i++
      continue
    }

    if (/^>\s?/.test(line)) {
      const quote = []
      while (i < lines.length && /^>\s?/.test(lines[i])) {
        quote.push(lines[i].replace(/^>\s?/, ''))
        i++
      }
      out.push(`<blockquote><p>${inlineFormat(quote.join(' '))}</p></blockquote>`)
      continue
    }

    if (/^[-*]\s/.test(line)) {
      out.push('<ul>')
      while (i < lines.length && /^[-*]\s/.test(lines[i])) {
        out.push(`<li>${inlineFormat(lines[i].replace(/^[-*]\s+/, ''))}</li>`)
        i++
      }
      out.push('</ul>')
      continue
    }

    if (/^\d+\.\s/.test(line)) {
      out.push('<ol>')
      while (i < lines.length && /^\d+\.\s/.test(lines[i])) {
        out.push(`<li>${inlineFormat(lines[i].replace(/^\d+\.\s+/, ''))}</li>`)
        i++
      }
      out.push('</ol>')
      continue
    }

    if (/^\|.+\|$/.test(line.trim())) {
      const rows = []
      while (i < lines.length && /^\|.+\|$/.test(lines[i].trim())) {
        rows.push(
          lines[i]
            .trim()
            .slice(1, -1)
            .split('|')
            .map((c) => c.trim()),
        )
        i++
      }
      if (rows.length >= 2 && /^[-:\s|]+$/.test(rows[1].join(''))) {
        rows.splice(1, 1)
      }
      if (rows.length) {
        out.push('<table><thead><tr>')
        for (const cell of rows[0]) {
          out.push(`<th>${inlineFormat(cell)}</th>`)
        }
        out.push('</tr></thead><tbody>')
        for (let r = 1; r < rows.length; r++) {
          out.push('<tr>')
          for (const cell of rows[r]) {
            out.push(`<td>${inlineFormat(cell)}</td>`)
          }
          out.push('</tr>')
        }
        out.push('</tbody></table>')
      }
      continue
    }

    if (line.trim() === '') {
      i++
      continue
    }

    const para = []
    while (i < lines.length && lines[i].trim() !== '' && !/^#{1,3}\s/.test(lines[i]) && !/^```/.test(lines[i])) {
      para.push(lines[i])
      i++
    }
    out.push(`<p>${inlineFormat(para.join(' '))}</p>`)
  }

  return out.join('\n')
}

/** Resolve relative markdown links against current doc path. */
export function resolveDocLink(currentRel, href) {
  const h = (href || '').trim().replace(/\\/g, '/')
  if (!h || /^https?:\/\//i.test(h) || h.startsWith('#')) return null
  if (h.startsWith('docs/') || h === 'language.md' || h === 'START_HERE.md' || h === 'README.md') {
    return h.endsWith('.md') ? h : `${h}.md`
  }
  const base = (currentRel || 'docs/README.md').replace(/\\/g, '/')
  const dir = base.includes('/') ? base.slice(0, base.lastIndexOf('/') + 1) : ''
  let target = dir + h
  if (target.startsWith('../')) {
    const parts = base.split('/')
    parts.pop()
    let t = h
    while (t.startsWith('../')) {
      t = t.slice(3)
      parts.pop()
    }
    target = [...parts, t].join('/')
  }
  if (!target.endsWith('.md')) {
    if (target.endsWith('/')) target += 'README.md'
    else target += '.md'
  }
  return target.replace(/^\.\//, '')
}
