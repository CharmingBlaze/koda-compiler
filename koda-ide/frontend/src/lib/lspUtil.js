/** Build file:// URI for LSP (Windows drive paths included). */
export function pathToFileURI(absPath) {
  const norm = absPath.replace(/\\/g, '/')
  if (/^[A-Za-z]:\//.test(norm)) {
    return 'file:///' + norm
  }
  return 'file://' + norm
}

export function rpcInitialize(id = 1) {
  return JSON.stringify({
    jsonrpc: '2.0',
    id,
    method: 'initialize',
    params: {
      processId: null,
      clientInfo: { name: 'koda-studio', version: '0.1' },
      rootUri: null,
      capabilities: {},
    },
  })
}

export function notifyInitialized() {
  return JSON.stringify({
    jsonrpc: '2.0',
    method: 'initialized',
    params: {},
  })
}

export function notifyDidOpen(uri, text) {
  return JSON.stringify({
    jsonrpc: '2.0',
    method: 'textDocument/didOpen',
    params: {
      textDocument: { uri, languageId: 'koda', version: 1, text },
    },
  })
}

export function notifyDidChange(uri, text) {
  return JSON.stringify({
    jsonrpc: '2.0',
    method: 'textDocument/didChange',
    params: {
      textDocument: { uri, version: 2 },
      contentChanges: [{ text }],
    },
  })
}
