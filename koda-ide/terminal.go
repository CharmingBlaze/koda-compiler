package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"sync"

	"github.com/creack/pty"
	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

// TermSession wraps a shell with a PTY (Unix) or pipes (Windows fallback when PTY fails).
type TermSession struct {
	id  string
	cmd *exec.Cmd
	tty io.ReadWriteCloser
	mu  sync.Mutex
}

var (
	termMu    sync.Mutex
	terminals = map[string]*TermSession{}
	termSeq   int
)

func shellCommand() (string, []string) {
	if runtime.GOOS == "windows" {
		return "powershell.exe", []string{"-NoLogo"}
	}
	sh := os.Getenv("SHELL")
	if sh == "" {
		sh = "/bin/sh"
	}
	return sh, []string{"-l"}
}

func (a *App) TerminalStart(cols, rows int) (string, error) {
	if a.ctx == nil {
		return "", fmt.Errorf("app not started")
	}
	name, args := shellCommand()
	cmd := exec.Command(name, args...)
	cmd.Env = os.Environ()

	tty, err := pty.Start(cmd)
	if err != nil {
		return "", err
	}
	_ = pty.Setsize(tty, &pty.Winsize{Rows: uint16(rows), Cols: uint16(cols)})

	termMu.Lock()
	termSeq++
	id := fmt.Sprintf("t%d", termSeq)
	terminals[id] = &TermSession{id: id, cmd: cmd, tty: tty}
	termMu.Unlock()

	go a.pumpTerminal(id, tty)
	return id, nil
}

func (a *App) pumpTerminal(id string, r io.Reader) {
	buf := make([]byte, 8192)
	for {
		n, err := r.Read(buf)
		if n > 0 {
			wailsruntime.EventsEmit(a.ctx, "term:data", map[string]any{
				"id":   id,
				"data": string(buf[:n]),
			})
		}
		if err != nil {
			wailsruntime.EventsEmit(a.ctx, "term:exit", map[string]any{"id": id})
			termMu.Lock()
			delete(terminals, id)
			termMu.Unlock()
			return
		}
	}
}

func (a *App) TerminalWrite(id, data string) error {
	termMu.Lock()
	s, ok := terminals[id]
	termMu.Unlock()
	if !ok {
		return fmt.Errorf("unknown terminal session")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	_, err := io.WriteString(s.tty, data)
	return err
}

func (a *App) TerminalResize(id string, cols, rows int) error {
	termMu.Lock()
	s, ok := terminals[id]
	termMu.Unlock()
	if !ok {
		return fmt.Errorf("unknown terminal session")
	}
	f, ok := s.tty.(*os.File)
	if !ok {
		return nil
	}
	return pty.Setsize(f, &pty.Winsize{Rows: uint16(rows), Cols: uint16(cols)})
}

func (a *App) TerminalClose(id string) error {
	termMu.Lock()
	s, ok := terminals[id]
	if ok {
		delete(terminals, id)
	}
	termMu.Unlock()
	if !ok {
		return nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	_ = s.tty.Close()
	if s.cmd.Process != nil {
		_ = s.cmd.Process.Kill()
	}
	return nil
}

func shutdownTerminals() {
	termMu.Lock()
	defer termMu.Unlock()
	for id, s := range terminals {
		_ = s.tty.Close()
		if s.cmd.Process != nil {
			_ = s.cmd.Process.Kill()
		}
		delete(terminals, id)
	}
}
