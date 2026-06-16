# Koda for VS Code

Syntax highlighting for [Koda](https://github.com/CharmingBlaze/koda-compiler) (`.koda` files).

## Install (local / unpublished)

`.koda` files map to the **Koda** language via `extensions: [".koda"]` in `contributes.languages` — no manual language mode selection after install.

1. Open this folder in VS Code: `vscode-koda/`
2. Press **F5** to launch an Extension Development Host with Koda highlighting enabled.

Or package and install:

```bash
npm install -g @vscode/vsce
cd vscode-koda
vsce package
code --install-extension koda-0.5.0.vsix
```

## What you get

- Keywords (`func`, `let`, `struct`, `match`, `fallthrough`, …)
- Comments (`//`, `/* */`)
- Strings (double, single, backtick templates with `${…}`)
- String interpolation in `"Score: {score}"`
- `#include` directives and `@game` / `@math` stdlib names
- Builtin calls (`print`, `assert`, `gc`, …)

Full language server and completions live in **Koda Studio** (`koda-ide/`). This extension is the lightweight path for VS Code / Cursor users.
