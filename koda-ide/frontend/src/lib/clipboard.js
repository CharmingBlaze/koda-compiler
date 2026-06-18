export async function copyText(text) {
  const value = text == null ? '' : String(text)
  if (!value) return false
  try {
    await navigator.clipboard.writeText(value)
    return true
  } catch {
    const area = document.createElement('textarea')
    area.value = value
    area.setAttribute('readonly', '')
    area.style.position = 'fixed'
    area.style.inset = '-9999px auto auto -9999px'
    document.body.appendChild(area)
    area.select()
    let ok = false
    try {
      ok = document.execCommand('copy')
    } catch {
      ok = false
    } finally {
      area.remove()
    }
    return ok
  }
}
