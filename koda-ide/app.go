package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"

	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
	"koda/api"
)

var reSafeProjectName = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9._\- ]{0,126}$`)

// App is the Wails-bound API for Koda Studio.
type App struct {
	ctx context.Context
	mu  sync.RWMutex

	workspaceRoot string

	docMu   sync.Mutex
	lspDocs map[string]string

	runMu sync.Mutex
}

// NewApp creates a new App application struct.
func NewApp() *App {
	return &App{lspDocs: make(map[string]string)}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) shutdown(ctx context.Context) {
	shutdownTerminals()
}

// PickWorkspaceFolder opens a native folder picker.
func (a *App) PickWorkspaceFolder() (string, error) {
	return wailsruntime.OpenDirectoryDialog(a.ctx, wailsruntime.OpenDialogOptions{
		Title: "Open Koda project folder",
	})
}

// PickParentFolderForNewProject asks where to place a new project folder (parent directory).
func (a *App) PickParentFolderForNewProject() (string, error) {
	return wailsruntime.OpenDirectoryDialog(a.ctx, wailsruntime.OpenDialogOptions{
		Title: "New Koda project — choose parent folder",
	})
}

const newProjectMainKoda = `// main.koda — entry point for your project.
// Press Run (F5) to execute.

print("Hello from Koda!");
`

const newProjectReadme = `# Koda project

- Open **main.koda** and start coding.
- Use **File → Save** (Ctrl+S) to write changes to disk.
- **Run** (F5) compiles and runs the current file via the native (LLVM) pipeline.

Docs: see the Koda repo examples/ folder for larger samples.
`

// CreateProjectInParent creates parent/name with starter files and returns the new workspace path.
func (a *App) CreateProjectInParent(parentDir, projectName string) (string, error) {
	name := strings.TrimSpace(projectName)
	if name == "" {
		return "", fmt.Errorf("project name is empty")
	}
	if strings.TrimRight(name, ".") == "" || name == "." || name == ".." {
		return "", fmt.Errorf("invalid project name")
	}
	if !reSafeProjectName.MatchString(name) {
		return "", fmt.Errorf("use letters, numbers, spaces, dot, dash, or underscore only")
	}
	parent, err := filepath.Abs(parentDir)
	if err != nil {
		return "", err
	}
	if fi, err := os.Stat(parent); err != nil || !fi.IsDir() {
		return "", fmt.Errorf("parent is not a directory")
	}
	root := filepath.Join(parent, name)
	if _, err := os.Stat(root); err == nil {
		return "", fmt.Errorf("folder already exists: %s", name)
	} else if !os.IsNotExist(err) {
		return "", err
	}
	if err := os.MkdirAll(root, 0755); err != nil {
		return "", err
	}
	if err := os.WriteFile(filepath.Join(root, "main.koda"), []byte(newProjectMainKoda), 0644); err != nil {
		return "", err
	}
	if err := os.WriteFile(filepath.Join(root, "README.md"), []byte(newProjectReadme), 0644); err != nil {
		return "", err
	}
	return filepath.Abs(root)
}

// OpenWorkspace sets the active workspace root (must be an existing directory).
func (a *App) OpenWorkspace(root string) error {
	fi, err := os.Stat(root)
	if err != nil {
		return err
	}
	if !fi.IsDir() {
		return fmt.Errorf("not a directory: %s", root)
	}
	abs, err := filepath.Abs(root)
	if err != nil {
		return err
	}
	a.mu.Lock()
	a.workspaceRoot = abs
	a.mu.Unlock()
	return nil
}

// GetWorkspaceRoot returns the current workspace or empty string.
func (a *App) GetWorkspaceRoot() string {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.workspaceRoot
}

func workspacePath(root, rel string) (string, error) {
	if root == "" {
		return "", fmt.Errorf("no workspace open")
	}
	cleanRoot, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}
	full := filepath.Join(cleanRoot, filepath.FromSlash(rel))
	cleanFull, err := filepath.Abs(full)
	if err != nil {
		return "", err
	}
	if !strings.EqualFold(cleanFull, cleanRoot) {
		prefix := cleanRoot + string(os.PathSeparator)
		if !strings.HasPrefix(strings.ToLower(cleanFull), strings.ToLower(prefix)) {
			return "", fmt.Errorf("path escapes workspace")
		}
	}
	return cleanFull, nil
}

// DirEntry is one row in the file tree.
type DirEntry struct {
	Name  string `json:"name"`
	Rel   string `json:"rel"`
	IsDir bool   `json:"isDir"`
}

// ListDir lists immediate children of rel (use "" for workspace root).
func (a *App) ListDir(rel string) ([]DirEntry, error) {
	a.mu.RLock()
	root := a.workspaceRoot
	a.mu.RUnlock()
	full, err := workspacePath(root, rel)
	if err != nil {
		return nil, err
	}
	f, err := os.Open(full)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	names, err := f.Readdirnames(-1)
	if err != nil {
		return nil, err
	}
	var out []DirEntry
	for _, name := range names {
		if strings.HasPrefix(name, ".") {
			continue
		}
		child := filepath.Join(full, name)
		relChild, _ := filepath.Rel(root, child)
		st, err := os.Stat(child)
		if err != nil {
			continue
		}
		out = append(out, DirEntry{Name: name, Rel: filepath.ToSlash(relChild), IsDir: st.IsDir()})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].IsDir != out[j].IsDir {
			return out[i].IsDir
		}
		return strings.ToLower(out[i].Name) < strings.ToLower(out[j].Name)
	})
	return out, nil
}

// ReadFile reads a file relative to workspace.
func (a *App) ReadFile(rel string) (string, error) {
	a.mu.RLock()
	root := a.workspaceRoot
	a.mu.RUnlock()
	full, err := workspacePath(root, rel)
	if err != nil {
		return "", err
	}
	b, err := os.ReadFile(full)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// WriteFile writes a file relative to workspace (creates parent dirs).
func (a *App) WriteFile(rel, content string) error {
	a.mu.RLock()
	root := a.workspaceRoot
	a.mu.RUnlock()
	full, err := workspacePath(root, rel)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(full), 0755); err != nil {
		return err
	}
	return os.WriteFile(full, []byte(content), 0644)
}

// AbsFromWorkspace joins workspace with a slash rel path and returns absolute host path.
func (a *App) AbsFromWorkspace(rel string) (string, error) {
	a.mu.RLock()
	root := a.workspaceRoot
	a.mu.RUnlock()
	return workspacePath(root, rel)
}

// DiagnoseFile returns load+compile diagnostics; overlay replaces on-disk content for that file when non-empty.
func (a *App) DiagnoseFile(absPath, overlay string) []api.Diagnostic {
	return api.Diagnose(absPath, overlay)
}

func (a *App) runProgramImpl(path, overlay string) {
	a.runMu.Lock()
	defer a.runMu.Unlock()
	out := &ideLineWriter{ctx: a.ctx, event: "koda:stdout"}
	errW := &ideLineWriter{ctx: a.ctx, event: "koda:stderr"}
	if err := api.RunWithWriters(path, overlay, out, errW); err != nil {
		wailsruntime.EventsEmit(a.ctx, "koda:stderr", err.Error()+"\n")
		wailsruntime.EventsEmit(a.ctx, "koda:runDone", map[string]any{"ok": false})
		return
	}
	wailsruntime.EventsEmit(a.ctx, "koda:runDone", map[string]any{"ok": true})
}

// ideLineWriter buffers stdout/stderr and emits one Wails event per line.
type ideLineWriter struct {
	ctx   context.Context
	event string
	buf   []byte
}

func (w *ideLineWriter) Write(p []byte) (n int, err error) {
	w.buf = append(w.buf, p...)
	for {
		i := bytes.IndexByte(w.buf, '\n')
		if i < 0 {
			break
		}
		line := string(w.buf[:i])
		w.buf = w.buf[i+1:]
		wailsruntime.EventsEmit(w.ctx, w.event, line+"\n")
	}
	return len(p), nil
}

// RunProgram compiles and runs path in a background goroutine; overlay may be "" to use disk.
func (a *App) RunProgram(path, overlay string) {
	go a.runProgramImpl(path, overlay)
}

func (a *App) buildProgramImpl(path, overlay, output string) {
	log := func(s string) { wailsruntime.EventsEmit(a.ctx, "koda:stdout", s) }
	err := api.BuildNativeHost(path, overlay, output, log)
	if err != nil {
		wailsruntime.EventsEmit(a.ctx, "koda:stderr", err.Error()+"\n")
		wailsruntime.EventsEmit(a.ctx, "koda:buildDone", map[string]any{"ok": false})
		return
	}
	wailsruntime.EventsEmit(a.ctx, "koda:buildDone", map[string]any{"ok": true, "output": output})
}

// BuildProgram runs the native Clang pipeline like `koda build`.
func (a *App) BuildProgram(path, overlay, output string) {
	go a.buildProgramImpl(path, overlay, output)
}

// DefaultBuildOutput suggests an output path next to the source file.
func (a *App) DefaultBuildOutput(entryPath string) string {
	return filepath.Join(filepath.Dir(entryPath), api.DefaultExeName(entryPath))
}
