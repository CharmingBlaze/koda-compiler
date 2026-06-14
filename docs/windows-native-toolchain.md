# Windows 11 native toolchain (LLVM 14, MSVC, go-llvm)

End-to-end setup for **Koda on Windows** without WSL: Go, LLVM 14 (clang + llc), MSVC (for MSVC-linked CGo / system libs), GNU Make, and environment variables for **go-llvm** (CGo).

Use **Developer Command Prompt for VS 2022** for compiler work so MSVC’s environment is loaded.

---

## Step 1 — Install LLVM 14

1. Open: [LLVM 14.0.6 release — llvmorg-14.0.6](https://github.com/llvm/llvm-project/releases/tag/llvmorg-14.0.6)
2. Under **Assets**, download **`LLVM-14.0.6-win64.exe`**.
3. Run the installer. On the options screen, set:
   - **Add LLVM to the system PATH for all users** (default is “Do not add to PATH”; change it or `clang` / `llc` will not be found).
4. Default install path: **`C:\Program Files\LLVM`**.

Open a **new** Command Prompt and verify:

```bat
clang --version
llc --version
```

Both should report **14.0.6**. If the command is not recognized, add to **System** `Path`:

`C:\Program Files\LLVM\bin`

---

## Step 2 — Install MSVC Build Tools

The C runtime and MSVC-linked tools expect the Microsoft linker and libraries when using the LLVM/MSVC stack.

1. [Visual Studio Build Tools](https://visualstudio.microsoft.com/visual-cpp-build-tools/) → **Download Build Tools**.
2. In the installer workload list, enable only:
   - **Desktop development with C++**
3. After install, use **Start → Developer Command Prompt for VS 2022** for builds (not plain `cmd` / PowerShell unless you’ve loaded the MSVC environment yourself).

Verify:

```bat
cl.exe /?
```

---

## Step 3 — Install Make

The runtime ships a **GNU Makefile** (`runtime/Makefile`). On Windows 11:

```bat
winget install ezwinports.make
```

Open a **new** Developer Command Prompt and verify:

```bat
make --version
```

If `winget` is unavailable, download [make without Guile (ezwinports)](https://sourceforge.net/projects/ezwinports/files/make-4.4.1-without-guile-w32-bin.zip), extract **`make.exe`**, and copy it to a directory already on `PATH` (for example `C:\Program Files\LLVM\bin`).

**Note:** `runtime/Makefile` sets `SHELL := sh.exe`. **Git for Windows** provides `sh.exe` on `PATH` for many setups. If `make` fails looking for `sh`, install Git for Windows or run the runtime build via **`scripts\build-runtime.ps1`** (see [CONTRIBUTING.md](../CONTRIBUTING.md)).

---

## Step 4 — Build the Koda C runtime on Windows

`runtime/Makefile` defaults to **`CC=gcc`** and **`AR=ar`**. If you only installed **LLVM** (no MinGW `gcc`), build with **clang** and **llvm-ar**:

```bat
set CC=clang
set AR=llvm-ar
make -C runtime
```

Alternatively, from the repo root in **PowerShell**:

```powershell
$env:CC = "clang"
$env:AR = "llvm-ar"
.\scripts\build-runtime.ps1
```

That produces **`runtime/libkoda_runtime.a`**.

If `make` fails on Unix-only commands (`rm`, `mkdir -p`), use the manual object + archive sequence from [CONTRIBUTING.md](../CONTRIBUTING.md) or the PowerShell script above.

---

## Step 5 — CGo environment variables (go-llvm)

go-llvm calls the LLVM C API via CGo. Set **system** environment variables (then restart the Developer Command Prompt):

| Variable        | Example value |
|-----------------|---------------|
| `CGO_CPPFLAGS`  | `-I"C:\Program Files\LLVM\include"` |
| `CGO_LDFLAGS`   | `-L"C:\Program Files\LLVM\lib"` plus the LLVM import library you actually have |

A typical pattern is to add **`-lLLVM-14`** (or whatever matches the `.lib` name under `C:\Program Files\LLVM\lib`). If linking fails with **cannot open input file `LLVM-14.lib`**, list the directory and point `CGO_LDFLAGS` at the real name:

```bat
dir "C:\Program Files\LLVM\lib\*.lib"
```

Ensure **`C:\Program Files\LLVM\bin`** is on `PATH` (Step 1).

---

## Step 6 — Verify the stack

In **Developer Command Prompt for VS 2022**:

```bat
go version
clang --version
llc --version
make --version
echo %CGO_CPPFLAGS%
echo %CGO_LDFLAGS%
```

From the Koda repo root (with `CC`/`AR` set if using clang for the runtime):

```bat
make -C runtime
go build ./cmd/koda/...
go test ./...
koda build tests\hello.koda
hello.exe
```

---

## Step 7 — Add go-llvm and verify CGo

```bat
cd path\to\koda-main
go get github.com/tinygo-org/go-llvm
go build github.com/tinygo-org/go-llvm
```

If **`go build github.com/tinygo-org/go-llvm`** succeeds, LLVM headers/libs and MSVC’s link environment are visible to CGo.

Pin the module version to the LLVM line you installed (for example the **llvm14** tagged revision) when you wire go-llvm into Koda’s `go.mod`.

---

## Troubleshooting

| Symptom | What to check |
|--------|----------------|
| `clang` / `llc` not found | LLVM `bin` on **system** `PATH`; new shell after install. |
| `make` not found | `winget` install or copy `make.exe` into a `PATH` directory. |
| go-llvm: missing LLVM headers | `CGO_CPPFLAGS` set as **system** variable; new Developer Prompt. |
| `LNK1181` cannot open `LLVM-14.lib` | `dir` the lib folder; adjust `CGO_LDFLAGS` to the actual `.lib` name. |
| `make -C runtime`: no `sh.exe` | Install Git for Windows, or use **`build-runtime.ps1`** with `CC`/`AR` set. |

---

## Self-contained release `koda` (optional; no LLVM on PATH for end users)

The **Release** GitHub Actions workflow (`.github/workflows/release.yml`, runs on **`v*`** tags) builds **`koda-windows-amd64.exe`** with **embedded** **`clang.exe`**, **`lld.exe`**, and **`libkoda_runtime.a`** under **`internal/embed/windows/amd64/`**, then compiles Koda with:

`go build -trimpath -tags release -ldflags="-s -w" -o koda-windows-amd64.exe ./cmd/koda`

That artifact is intended for **players and authors who only download one file**: it extracts the embedded tools to a temp directory on first **`koda build`** / **`koda run`** (see **`kodahome.FindToolchain`**). It does **not** replace the developer setup above for **contributing** from a clean clone (normal **`go build ./cmd/koda/...`** without **`-tags release`** still expects **LLVM on `PATH`** or a tarball next to the install, as in the rest of this document).

---

## What you have when this is done

- **Go** — Koda compiler and tests  
- **LLVM 14** — `clang`, `llc`, headers/libs for CGo and the current `.ll` pipeline  
- **MSVC** — linker/libs for typical Windows CGo and native links  
- **Make (+ sh)** — `runtime` archive via Makefile, or PowerShell fallback  
- **go-llvm** — ready for **OptimiseModule**, **EmitObjectFile**, and a **`--no-opt`** fallback as in the Phase 3 go-llvm integration plan  

No Linux, WSL, or VM required for this path.
