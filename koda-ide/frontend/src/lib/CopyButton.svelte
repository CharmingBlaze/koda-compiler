<script>
  import { copyText } from './clipboard.js'

  let {
    text = '',
    label = 'Copy',
    copiedLabel = 'Copied',
    title = 'Copy to clipboard',
    compact = false,
  } = $props()

  let copied = $state(false)
  let timer = 0

  async function copy() {
    if (!(await copyText(text))) return
    copied = true
    clearTimeout(timer)
    timer = window.setTimeout(() => {
      copied = false
    }, 1400)
  }
</script>

<button
  type="button"
  class="copy-button"
  class:copy-button-compact={compact}
  {title}
  aria-label={title}
  onclick={copy}
>
  <span aria-hidden="true">{copied ? '✓' : '⧉'}</span>
  <span>{copied ? copiedLabel : label}</span>
</button>
