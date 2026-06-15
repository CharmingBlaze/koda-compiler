# Koda Studio

Koda Studio is the desktop IDE for the KODA programming language. It is built to keep the first project simple while still giving day-to-day KODA work the essentials: a project tree, tabs, diagnostics, command palette, integrated terminal, run/build actions, and language-aware editing.

## Highlights

- **Welcome screen** with SDK health check — know immediately if the SDK is set up correctly.
- **One-click project templates** — Hello app, text game, or Raylib bouncing ball.
- Create or open KODA workspaces from the desktop app (no terminal required).
- Browse nested project folders and open `.koda` files quickly.
- Edit with CodeMirror, KODA syntax highlighting, snippets, bracket matching, search, folding, and lint squiggles.
- Run the current file with `F5` and build a native executable with `Ctrl+Shift+B`.
- Use the command palette with `Ctrl+P` or `Ctrl+K`.
- Track unsaved files in tabs and the status bar.
- Jump from diagnostics directly to the failing line.
- Use the integrated terminal with `Ctrl+J`.

## Development

**Windows** — builds the frontend, desktop app, and launches Studio:

```powershell
.\run-koda-studio.ps1
```

Or double-click `run-koda-studio.cmd`.

**macOS / Linux:**

```bash
chmod +x run-koda-studio.sh
./run-koda-studio.sh
```

If the window closes immediately on Windows, repair **Microsoft Edge WebView2 Runtime**.
On Linux, install `libwebkit2gtk-4.1-dev` (or 4.0) and GTK3 dev packages to build.

Install dependencies once:

```powershell
cd frontend
npm install
```

Build and run the desktop IDE:

```powershell
.\run-koda-studio.ps1
```

Build a redistributable desktop package:

```powershell
wails build
```

Place the built **`Koda Studio.exe`** next to `koda.exe` and `stdlib/` in the SDK zip so end users can double-click to start coding.

## End-user layout (SDK zip)

```
koda/                          (Windows example)
  koda.exe                     (or `koda` on Linux/macOS)
  Koda Studio.exe              Windows IDE
  Koda Studio                  Linux IDE binary
  Koda Studio.app/             macOS IDE
  Start Koda Studio.bat        Windows launcher
  start-koda-studio.sh         Linux launcher
  Start Koda Studio.command    macOS launcher
  stdlib/
  docs/
```

Koda Studio finds the SDK automatically when it sits in the same folder as `stdlib/`. Set `KODA_HOME` to override.

---

- `app.go` exposes workspace, file, run, build, and diagnostic actions to the UI.
- `lsp.go` provides the lightweight KODA language-server bridge used by the editor.
- `terminal.go` owns integrated terminal sessions.
- `frontend/src/App.svelte` is the main IDE shell.
- `frontend/src/lib` contains the editor, command palette, sidebar, status bar, preview, terminal, and language helpers.

## Current Focus

The IDE is being shaped around a professional KODA workflow: fast startup, clear project navigation, strong diagnostics, safe file handling, and language-specific editor assistance. Future work should keep that bar: every feature should make it easier to understand, write, run, and ship KODA programs.
