package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"path/filepath"
	"runtime"
	"strings"

	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
	"koda/api"
)

func fileURIToPath(uri string) (string, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return "", err
	}
	if u.Scheme != "file" {
		return "", fmt.Errorf("expected file URI")
	}
	p := u.Path
	if runtime.GOOS == "windows" && len(p) >= 3 && p[0] == '/' && p[2] == ':' {
		p = strings.TrimPrefix(p, "/")
	}
	return filepath.Clean(filepath.FromSlash(p)), nil
}

func pathToFileURI(abs string) (string, error) {
	abs, err := filepath.Abs(abs)
	if err != nil {
		return "", err
	}
	abs = filepath.Clean(abs)
	if runtime.GOOS == "windows" {
		return "file:///" + filepath.ToSlash(abs), nil
	}
	return "file://" + filepath.ToSlash(abs), nil
}

// LSPMessage handles a single JSON-RPC 2.0 message from the client (subset for Koda IDE).
func (a *App) LSPMessage(raw string) (string, error) {
	var envelope struct {
		JSONRPC string           `json:"jsonrpc"`
		ID      *json.RawMessage `json:"id"`
		Method  string           `json:"method"`
		Params  json.RawMessage  `json:"params"`
	}
	if err := json.Unmarshal([]byte(raw), &envelope); err != nil {
		return `{"jsonrpc":"2.0","id":null,"error":{"code":-32700,"message":"parse error"}}`, nil
	}
	if envelope.ID == nil {
		a.handleLSPNotify(envelope.Method, envelope.Params)
		return "", nil
	}
	switch envelope.Method {
	case "initialize":
		res := map[string]any{
			"capabilities": map[string]any{
				"textDocumentSync":   1,
				"completionProvider": map[string]any{"resolveProvider": false},
			},
			"serverInfo": map[string]string{"name": "koda-lsp", "version": "0.1"},
		}
		return jsonRPCResult(*envelope.ID, res)
	case "shutdown":
		return jsonRPCResult(*envelope.ID, nil)
	case "textDocument/completion":
		return jsonRPCResult(*envelope.ID, map[string]any{"isIncomplete": true, "items": []any{}})
	default:
		return jsonRPCError(*envelope.ID, -32601, "method not found: "+envelope.Method)
	}
}

func (a *App) handleLSPNotify(method string, params json.RawMessage) {
	switch method {
	case "initialized", "$/cancelRequest":
		return
	case "textDocument/didOpen":
		var p struct {
			TextDocument struct {
				URI  string `json:"uri"`
				Text string `json:"text"`
			} `json:"textDocument"`
		}
		if json.Unmarshal(params, &p) != nil {
			return
		}
		path, err := fileURIToPath(p.TextDocument.URI)
		if err != nil {
			return
		}
		a.docMu.Lock()
		a.lspDocs[path] = p.TextDocument.Text
		a.docMu.Unlock()
		a.publishDiagnostics(path, p.TextDocument.Text)
	case "textDocument/didChange":
		var p struct {
			TextDocument struct {
				URI string `json:"uri"`
			} `json:"textDocument"`
			ContentChanges []struct {
				Text string `json:"text"`
			} `json:"contentChanges"`
		}
		if json.Unmarshal(params, &p) != nil || len(p.ContentChanges) == 0 {
			return
		}
		path, err := fileURIToPath(p.TextDocument.URI)
		if err != nil {
			return
		}
		text := p.ContentChanges[len(p.ContentChanges)-1].Text
		a.docMu.Lock()
		a.lspDocs[path] = text
		a.docMu.Unlock()
		a.publishDiagnostics(path, text)
	}
}

func (a *App) publishDiagnostics(path, text string) {
	diags := api.Diagnose(path, text)
	uri, err := pathToFileURI(path)
	if err != nil {
		return
	}
	items := make([]map[string]any, 0, len(diags))
	for _, d := range diags {
		line := d.Line
		if line < 1 {
			line = 1
		}
		col := d.Col
		if col < 1 {
			col = 1
		}
		severity := 1
		if d.Severity == "warning" {
			severity = 2
		}
		items = append(items, map[string]any{
			"range": map[string]any{
				"start": map[string]int{"line": line - 1, "character": col - 1},
				"end":   map[string]int{"line": line - 1, "character": col},
			},
			"message":  d.Message,
			"severity": severity,
			"source":   "koda",
		})
	}
	wailsruntime.EventsEmit(a.ctx, "lsp:publishDiagnostics", map[string]any{
		"uri":           uri,
		"path":          path,
		"diagnostics":   items,
		"diagnosticsRaw": diags,
	})
}

func jsonRPCResult(id json.RawMessage, result any) (string, error) {
	out := struct {
		JSONRPC string          `json:"jsonrpc"`
		ID      json.RawMessage `json:"id"`
		Result  any             `json:"result"`
	}{JSONRPC: "2.0", ID: id, Result: result}
	b, err := json.Marshal(out)
	return string(b), err
}

func jsonRPCError(id json.RawMessage, code int, message string) (string, error) {
	out := struct {
		JSONRPC string          `json:"jsonrpc"`
		ID      json.RawMessage `json:"id"`
		Error   map[string]any  `json:"error"`
	}{JSONRPC: "2.0", ID: id, Error: map[string]any{"code": code, "message": message}}
	b, err := json.Marshal(out)
	return string(b), err
}
