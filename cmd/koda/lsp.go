package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"koda/api"
)

func runLSP() error {
	reader := bufio.NewReader(os.Stdin)
	for {
		msg, err := readLSPMessage(reader)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		resp, notify, err := handleLSPMessage(msg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "lsp: %v\n", err)
			continue
		}
		if notify != "" {
			if err := writeLSPMessage(os.Stdout, notify); err != nil {
				return err
			}
		}
		if resp != "" {
			if err := writeLSPMessage(os.Stdout, resp); err != nil {
				return err
			}
		}
	}
}

func readLSPMessage(r *bufio.Reader) (string, error) {
	var contentLen int
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			return "", err
		}
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		if strings.HasPrefix(strings.ToLower(line), "content-length:") {
			fmt.Sscanf(line, "Content-Length: %d", &contentLen)
		}
	}
	if contentLen <= 0 {
		return "", fmt.Errorf("missing Content-Length")
	}
	buf := make([]byte, contentLen)
	_, err := io.ReadFull(r, buf)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

func writeLSPMessage(w io.Writer, body string) error {
	header := fmt.Sprintf("Content-Length: %d\r\n\r\n", len(body))
	if _, err := w.Write([]byte(header)); err != nil {
		return err
	}
	_, err := w.Write([]byte(body))
	return err
}

func handleLSPMessage(raw string) (response, notification string, err error) {
	var envelope struct {
		JSONRPC string           `json:"jsonrpc"`
		ID      *json.RawMessage `json:"id"`
		Method  string           `json:"method"`
		Params  json.RawMessage  `json:"params"`
	}
	if err := json.Unmarshal([]byte(raw), &envelope); err != nil {
		return `{"jsonrpc":"2.0","id":null,"error":{"code":-32700,"message":"parse error"}}`, "", nil
	}
	if envelope.ID == nil {
		handleLSPNotify(envelope.Method, envelope.Params)
		return "", "", nil
	}
	switch envelope.Method {
	case "initialize":
		res := map[string]any{
			"capabilities": map[string]any{
				"textDocumentSync": 1,
			},
			"serverInfo": map[string]string{"name": "koda-lsp", "version": version},
		}
		return jsonRPCResult(*envelope.ID, res), "", nil
	case "shutdown":
		return jsonRPCResult(*envelope.ID, nil), "", nil
	default:
		return jsonRPCError(*envelope.ID, -32601, "method not found: "+envelope.Method), "", nil
	}
}

func handleLSPNotify(method string, params json.RawMessage) {
	switch method {
	case "initialized", "$/cancelRequest":
		return
	case "textDocument/didOpen", "textDocument/didChange":
		var p struct {
			TextDocument struct {
				URI  string `json:"uri"`
			} `json:"textDocument"`
			ContentChanges []struct {
				Text string `json:"text"`
			} `json:"contentChanges"`
			Text string `json:"text"`
		}
		if json.Unmarshal(params, &p) != nil {
			return
		}
		text := p.Text
		if len(p.ContentChanges) > 0 {
			text = p.ContentChanges[len(p.ContentChanges)-1].Text
		}
		path, err := lspURIToPath(p.TextDocument.URI)
		if err != nil {
			return
		}
		publishLSPDiagnostics(path, text)
	}
}

func publishLSPDiagnostics(path, text string) {
	diags := api.Diagnose(path, text)
	items := make([]map[string]any, 0, len(diags))
	for _, d := range diags {
		sev := 1
		if d.Severity == "warning" {
			sev = 2
		}
		items = append(items, map[string]any{
			"range": map[string]any{
				"start": map[string]int{"line": d.Line - 1, "character": d.Col - 1},
				"end":   map[string]int{"line": d.Line - 1, "character": d.Col},
			},
			"severity": sev,
			"message":  d.Message,
		})
	}
	uri, _ := pathToLSPURI(path)
	notify := map[string]any{
		"jsonrpc": "2.0",
		"method":  "textDocument/publishDiagnostics",
		"params": map[string]any{
			"uri":         uri,
			"diagnostics": items,
		},
	}
	b, _ := json.Marshal(notify)
	_ = writeLSPMessage(os.Stdout, string(b))
}

func jsonRPCResult(id json.RawMessage, result any) string {
	b, _ := json.Marshal(map[string]any{"jsonrpc": "2.0", "id": id, "result": result})
	return string(b)
}

func jsonRPCError(id json.RawMessage, code int, msg string) string {
	b, _ := json.Marshal(map[string]any{
		"jsonrpc": "2.0",
		"id":      id,
		"error":   map[string]any{"code": code, "message": msg},
	})
	return string(b)
}

func lspURIToPath(uri string) (string, error) {
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

func pathToLSPURI(abs string) (string, error) {
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
