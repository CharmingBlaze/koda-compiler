export const themes = [
  { id: 'koda-dark', name: 'Koda Dark' },
  { id: 'aurora', name: 'Aurora' },
  { id: 'midnight', name: 'Midnight' },
  { id: 'graphite', name: 'Graphite' },
  { id: 'ocean', name: 'Ocean' },
  { id: 'forest', name: 'Forest' },
  { id: 'ember', name: 'Ember' },
  { id: 'rose', name: 'Rose' },
  { id: 'violet', name: 'Violet' },
  { id: 'solar', name: 'Solar Dark' },
  { id: 'dracula', name: 'Dracula' },
  { id: 'nord', name: 'Nord' },
  { id: 'tokyo', name: 'Tokyo Night' },
  { id: 'matrix', name: 'Matrix' },
  { id: 'cyberpunk', name: 'Cyberpunk' },
  { id: 'paper', name: 'Paper Light' },
  { id: 'github-light', name: 'GitHub Light' },
  { id: 'latte', name: 'Latte' },
  { id: 'sepia', name: 'Sepia' },
  { id: 'high-contrast', name: 'High Contrast' },
]

export const defaultTheme = 'koda-dark'

export function safeTheme(id) {
  return themes.some((theme) => theme.id === id) ? id : defaultTheme
}
