# Start here — Koda in 5 minutes

**Koda replaces C and C++ for beginners** who want to make **games and desktop apps** without installing Go, Python, LLVM, Visual Studio, or a package manager.

Download the SDK zip for **your OS**, unzip, and launch **Koda Studio** — or use `koda` from a terminal.

---

## Pick your platform

| Platform | Download | Easiest way to start |
|----------|----------|----------------------|
| **Windows** | `koda-*-sdk-windows-amd64.zip` | Double-click **`Start Koda Studio.bat`** |
| **Linux** (x64 or ARM64) | `koda-*-sdk-linux-*.zip` | Run **`./start-koda-studio.sh`** |
| **macOS** (Intel or Apple Silicon) | `koda-*-sdk-darwin-*.zip` | Double-click **`Start Koda Studio.command`** |

Every zip is self-contained: compiler, stdlib, docs, examples, Raylib (where available), and the IDE. No extra installs.

Press **F1** inside Studio for all help and documentation.

---

## Windows

1. Download **`koda-*-sdk-windows-amd64.zip`** from [GitHub Releases](https://github.com/CharmingBlaze/koda-compiler/releases).
2. Unzip to a folder, e.g. `C:\koda`.
3. Double-click **`Start Koda Studio.bat`** (or **`Koda Studio.exe`**).
4. Pick a template on the welcome screen, edit `src/main.koda`, press **F5** to run.

**Terminal (optional):**

```powershell
.\koda.exe doctor
.\koda.exe new bounce --template graphics
cd bounce
..\koda.exe run
```

**Note:** The IDE needs the **WebView2** runtime (preinstalled on most Windows 10/11). If the window is blank, install [WebView2 Runtime](https://developer.microsoft.com/microsoft-edge/webview2/).

---

## Linux

1. Download **`koda-*-sdk-linux-amd64.zip`** or **`koda-*-sdk-linux-arm64.zip`**.
2. Unzip, e.g. `~/koda`.
3. In a terminal:

```bash
cd ~/koda
chmod +x koda kodawrap start-koda-studio.sh
./start-koda-studio.sh
```

Or open **`koda-studio.desktop`** from your file manager (after `chmod +x start-koda-studio.sh`).

**Terminal-only workflow:**

```bash
chmod +x koda kodawrap
./koda doctor
./koda new bounce --template graphics
cd bounce
../koda run
```

**Note:** Koda Studio on Linux uses GTK + WebKit (standard on Ubuntu, Fedora, etc.). If the IDE fails to start, install `libwebkit2gtk-4.1` / `libwebkit2gtk-4.0` for your distro.

---

## macOS

1. Download **`koda-*-sdk-darwin-arm64.zip`** (Apple Silicon) or **`koda-*-sdk-darwin-amd64.zip`** (Intel).
2. Unzip, e.g. `~/koda`.
3. Double-click **`Start Koda Studio.command`** in Finder.

If macOS blocks the app: **System Settings → Privacy & Security → Open Anyway** (first launch only).

**Terminal (optional):**

```bash
cd ~/koda
chmod +x koda kodawrap
./koda doctor
./koda new bounce --template graphics
cd bounce
../koda run
```

---

## What you get (and what you do not need)

| You get | You do **not** need |
|---------|---------------------|
| Native `.exe` / binary (like C/C++) | Go |
| Koda Studio IDE + full docs (F1) | Python |
| Built-in graphics (`@game` + Raylib in zip) | Node.js or npm |
| Automatic memory (GC) | LLVM, VS, or CMake to **use** Koda |
| `koda run` — edit and play instantly | Rust or Cargo |

Release builds **embed Clang, LLVM, and the runtime inside `koda`**. The first compile unpacks tools to a temp folder. Nothing downloads from the internet at compile time.

---

## Optional: add `koda` to PATH

**Windows** (from SDK folder):

```powershell
powershell -ExecutionPolicy Bypass -File scripts\install-koda.ps1
```

**Linux / macOS:**

```bash
bash scripts/install-koda.sh
```

Keep **`stdlib/`** next to the SDK, or set **`KODA_HOME`** to the unzipped folder.

---

## Your first program

```bash
koda new mytool
cd mytool
koda run
koda build -o mytool
```

Ship the built binary to friends — they do **not** need Koda installed to run it.

---

## Learn more

| Doc | Purpose |
|-----|---------|
| [Beginner's guide](docs/beginners-guide.md) | Full tutorial |
| [Game development](docs/guides/game-dev.md) | Loops, input, shipping |
| [FAQ](docs/faq.md) | Common questions |

When something breaks: **`koda doctor`** or **`./check-sdk.sh`** — fix every line marked **FAIL**.
