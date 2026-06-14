// @ts-check
/**
 * When the Wails Go bridge is missing (e.g. Vite `npm run dev`), use the File System
 * Access API so Open workspace / New project / tree / save still work in Chromium-based browsers.
 */
import * as WailsApp from '../../wailsjs/go/main/App.js'

/** @type {FileSystemDirectoryHandle | null} */
let browserRoot = null

/** @type {FileSystemDirectoryHandle | null} */
let newProjectParentHandle = null

const NEW_MAIN_KODA = `// main.koda — entry point for your project.
// Press Run (F5) to execute.

print("Hello from Koda!");
`

const NEW_README = `# Koda project

- Open **main.koda** and start coding.
- Use **File → Save** (Ctrl+S) to write changes to disk.
- **Run** (F5) compiles and runs the current file.

Docs: see the Koda repo \`examples/\` for larger samples.
`

const SAFE_NAME = /^[a-zA-Z0-9][a-zA-Z0-9._\- ]{0,126}$/

export function wailsAppBinding() {
  if (typeof window === 'undefined') return null
  return window['go']?.['main']?.['App'] ?? null
}

export function wailsNativeBridgeOk() {
  const app = wailsAppBinding()
  return !!(app && typeof app.PickWorkspaceFolder === 'function')
}

export function isBrowserWorkspace() {
  return browserRoot !== null
}

export function clearBrowserWorkspace() {
  browserRoot = null
}

export function clearNewProjectParentHandle() {
  newProjectParentHandle = null
}

/** @returns {boolean} */
export function hasShowDirectoryPicker() {
  return typeof window !== 'undefined' && typeof window.showDirectoryPicker === 'function'
}

/**
 * @param {string} name
 * @returns {string}
 */
function validateProjectName(name) {
  const t = name.trim()
  if (!t) throw new Error('project name is empty')
  if (t === '.' || t === '..' || t.replace(/\./g, '') === '') throw new Error('invalid project name')
  if (!SAFE_NAME.test(t)) throw new Error('use letters, numbers, spaces, dot, dash, or underscore only')
  return t
}

/**
 * @param {FileSystemDirectoryHandle} root
 * @param {string} rel
 * @returns {Promise<FileSystemDirectoryHandle>}
 */
async function walkToDirectory(root, rel) {
  const norm = (rel || '').replace(/\\/g, '/').replace(/^\/+/, '').replace(/\/+$/, '')
  if (!norm) return root
  if (norm.includes('..')) throw new Error('path escapes workspace')
  let dir = root
  for (const p of norm.split('/').filter(Boolean)) {
    dir = await dir.getDirectoryHandle(p)
  }
  return dir
}

/**
 * @param {string} rel
 * @returns {Promise<{ dir: FileSystemDirectoryHandle, fileName: string }>}
 */
async function resolveWritePath(rel) {
  const norm = rel.replace(/\\/g, '/').replace(/^\/+/, '')
  if (!norm || norm.includes('..')) throw new Error('invalid path')
  const parts = norm.split('/').filter(Boolean)
  if (parts.length === 0) throw new Error('invalid path')
  if (!browserRoot) throw new Error('no workspace open')
  let dir = browserRoot
  for (let i = 0; i < parts.length - 1; i++) {
    dir = await dir.getDirectoryHandle(parts[i], { create: true })
  }
  return { dir, fileName: parts[parts.length - 1] }
}

/**
 * @param {string} rel
 * @returns {Promise<{ dir: FileSystemDirectoryHandle, fileName: string }>}
 */
async function resolveReadPath(rel) {
  const norm = rel.replace(/\\/g, '/').replace(/^\/+/, '')
  if (!norm || norm.includes('..')) throw new Error('invalid path')
  const parts = norm.split('/').filter(Boolean)
  if (!browserRoot) throw new Error('no workspace open')
  let dir = browserRoot
  for (let i = 0; i < parts.length - 1; i++) {
    dir = await dir.getDirectoryHandle(parts[i])
  }
  return { dir, fileName: parts[parts.length - 1] }
}

export async function unifiedGetWorkspaceRoot() {
  if (browserRoot) return `Browser · ${browserRoot.name}`
  return WailsApp.GetWorkspaceRoot()
}

/**
 * @param {string} rel
 * @returns {Promise<{ name: string, rel: string, isDir: boolean }[]>}
 */
export async function unifiedListDir(rel) {
  if (browserRoot) {
    const dir = await walkToDirectory(browserRoot, rel)
    /** @type {{ name: string, rel: string, isDir: boolean }[]} */
    const out = []
    for await (const [name, handle] of dir.entries()) {
      if (name.startsWith('.')) continue
      const base = rel ? `${rel.replace(/\/+$/, '')}/${name}` : name
      const relChild = base.replace(/\\/g, '/').replace(/^\/+/, '')
      out.push({
        name,
        rel: relChild,
        isDir: handle.kind === 'directory',
      })
    }
    out.sort((a, b) => {
      if (a.isDir !== b.isDir) return a.isDir ? -1 : 1
      return a.name.localeCompare(b.name)
    })
    return out
  }
  return WailsApp.ListDir(rel)
}

/**
 * @param {string} rel
 * @returns {Promise<string>}
 */
export async function unifiedReadFile(rel) {
  if (browserRoot) {
    const { dir, fileName } = await resolveReadPath(rel)
    const fh = await dir.getFileHandle(fileName)
    const file = await fh.getFile()
    return await file.text()
  }
  return WailsApp.ReadFile(rel)
}

/**
 * @param {string} rel
 * @param {string} content
 */
export async function unifiedWriteFile(rel, content) {
  if (browserRoot) {
    const { dir, fileName } = await resolveWritePath(rel)
    const fh = await dir.getFileHandle(fileName, { create: true })
    const w = await fh.createWritable()
    await w.write(content)
    await w.close()
    return
  }
  return WailsApp.WriteFile(rel, content)
}

/**
 * @param {string} rel
 * @returns {Promise<string>}
 */
export async function unifiedAbsFromWorkspace(rel) {
  if (browserRoot) {
    const norm = rel.replace(/\\/g, '/').replace(/^\/+/, '')
    return `file:///koda-browser/${norm.split('/').map(encodeURIComponent).join('/')}`
  }
  return WailsApp.AbsFromWorkspace(rel)
}

/** Open folder picker and set browser workspace (no Wails). */
export async function browserOnlyOpenWorkspace() {
  if (!hasShowDirectoryPicker()) {
    throw new Error('This browser does not support folder access. Try Chrome or Edge, or run wails dev.')
  }
  // @ts-ignore — Chromium
  browserRoot = await window.showDirectoryPicker({ mode: 'readwrite' })
}

/**
 * Pick parent folder for new project (browser). Sets internal handle; returns a short label for the UI.
 * @returns {Promise<string>}
 */
export async function browserOnlyPickParentForNewProject() {
  if (!hasShowDirectoryPicker()) {
    throw new Error('This browser does not support folder access. Try Chrome or Edge, or run wails dev.')
  }
  // @ts-ignore
  newProjectParentHandle = await window.showDirectoryPicker({ mode: 'readwrite' })
  return newProjectParentHandle.name
}

/**
 * Create project folder under picked parent, write starter files, set as workspace.
 * @param {string} projectName
 */
export async function browserOnlyFinishNewProject(projectName) {
  if (!newProjectParentHandle) throw new Error('No parent folder selected')
  const name = validateProjectName(projectName)
  /** @type {FileSystemDirectoryHandle} */
  let child
  try {
    await newProjectParentHandle.getDirectoryHandle(name)
    throw new Error(`folder already exists: ${name}`)
  } catch (e) {
    // @ts-ignore
    if (e?.name === 'NotFoundError') {
      child = await newProjectParentHandle.getDirectoryHandle(name, { create: true })
    } else {
      throw e
    }
  }
  const mainH = await child.getFileHandle('main.koda', { create: true })
  const w1 = await mainH.createWritable()
  await w1.write(NEW_MAIN_KODA)
  await w1.close()
  const readmeH = await child.getFileHandle('README.md', { create: true })
  const w2 = await readmeH.createWritable()
  await w2.write(NEW_README)
  await w2.close()
  newProjectParentHandle = null
  browserRoot = child
}

/**
 * Wails: pick folder + OpenWorkspace. Browser: pick folder + set handle.
 * @returns {Promise<boolean>} false if user cancelled (Wails empty pick)
 */
export async function unifiedOpenWorkspaceFlow() {
  if (wailsNativeBridgeOk()) {
    clearBrowserWorkspace()
    clearNewProjectParentHandle()
    const picked = await WailsApp.PickWorkspaceFolder()
    if (!picked) return false
    await WailsApp.OpenWorkspace(picked)
    return true
  }
  await browserOnlyOpenWorkspace()
  return true
}

/**
 * @returns {Promise<string>} parent path (Wails) or display label (browser)
 */
export async function unifiedPickParentForNewProject() {
  if (wailsNativeBridgeOk()) {
    clearNewProjectParentHandle()
    return await WailsApp.PickParentFolderForNewProject()
  }
  return await browserOnlyPickParentForNewProject()
}

/**
 * @param {string} parentDir
 * @param {string} projectName
 */
export async function unifiedCreateProjectAndOpen(parentDir, projectName) {
  if (wailsNativeBridgeOk()) {
    clearBrowserWorkspace()
    const root = await WailsApp.CreateProjectInParent(parentDir, projectName)
    await WailsApp.OpenWorkspace(root)
    return
  }
  await browserOnlyFinishNewProject(projectName)
}
