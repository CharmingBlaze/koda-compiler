# Building applications with Koda

Koda is not only for games. The same **native compile → single binary** pipeline works for:

- Command-line tools and scripts
- File processors and small utilities
- Desktop apps with a window (via Raylib or other C UI libs)
- Data helpers using JSON and the stdlib

If you would have written a **C program** or a **Python script** but want a **compiled executable** with no runtime install for end users, Koda is a good fit.

---

## When Koda beats C for apps

| Task | Why Koda helps |
|------|----------------|
| Read/write files | `readFile`, `writeFile`, `appendFile` builtins |
| Parse JSON | `parseJSON`, `toJSON` |
| Quick CLI | `print`, `assert`, `len`, loops — no `printf` formats |
| Ship to users | `koda bundle` → folder with exe + assets |
| Call existing C libs | `kodawrap` without writing all logic in C |

---

## Project setup

```bash
koda new mytool
cd mytool
```

`koda.json` defines the entry and optional bundle assets:

```json
{
  "name": "mytool",
  "version": "0.1.0",
  "entry": "src/main.koda",
  "bundle": {
    "assets": ["assets"],
    "extra": ["README.md"]
  }
}
```

```bash
koda run
koda build -o mytool
koda bundle -o dist/mytool
```

---

## Example: file batch processor

```koda
func main() {
    let path = "data/input.txt";
    if (!fileExists(path)) {
        print("Missing:", path);
        return;
    }
    let text = readFile(path);
    let lines = text.split("\n");
    let out = [];
    for (let line of lines) {
        let trimmed = line.trim();
        if (len(trimmed) > 0) {
            out.push(trimmed.toUpper());
        }
    }
    writeFile("data/output.txt", out.join("\n"));
    print("Processed", len(out), "lines");
}
```

Run with `koda run` from the project root (paths relative to cwd).

---

## Example: JSON config tool

```koda
func main() {
    let raw = readFile("config.json");
    let cfg = parseJSON(raw);
    print("Server:", cfg.host, "Port:", cfg.port);
    cfg.debug = true;
    writeFile("config.json", toJSON(cfg));
}
```

See [reference.md](../reference.md) for `stdlib/json.koda` helpers.

---

## Example: CLI with arguments

`args()` returns an array of command-line strings. Index `0` is the program name; arguments after `--` follow:

```koda
func main() {
    let a = args();
    if (len(a) < 3) {
        print("Usage: mytool <input> <output>");
        return;
    }
    let inputPath = a[1];
    let outputPath = a[2];
    if (!fileExists(inputPath)) {
        print("Missing:", inputPath);
        return;
    }
    writeFile(outputPath, readFile(inputPath).toUpper());
    print("Wrote", outputPath);
}
```

```bash
koda run src/main.koda -- data/input.txt data/output.txt
```

Read environment variables with `env("VAR_NAME")` — returns `null` when unset.

---

## Example: batch mode CLI

Process a config file when no arguments are passed:

```koda
func main() {
    let a = args();
    if (len(a) > 1) {
        print("Processing", a[1]);
        writeFile("output.log", readFile(a[1]));
        return;
    }
    let cfg = parseJSON(readFile("config.json"));
    print("Running", cfg.name, "version", cfg.version);
    writeFile("output.log", "done\n");
}
```

---

## Desktop apps with a window

Use the **graphics** template or Raylib:

```bash
koda new myapp --template graphics
```

Set Raylib link flags (see [Raylib guide](raylib.md)), then `koda run`. The same binary can be a settings panel, visualizer, or simple editor UI.

For pure terminal UIs, stay with `print` / `input` — no window required.

---

## Structs for application data

Use structs when C would use a `struct` — configs, entities, records:

```koda
struct Config {
    host, port, debug
}

let cfg = Config { host: "localhost", port: 8080, debug: false };
if (cfg.debug) {
    print("Debug mode on");
}
```

---

## Shipping

1. **`koda build -o MyApp`** — executable only.
2. **`koda bundle -o dist/MyApp`** — exe + launcher script + README for end users.
3. List data files in `koda.json` → `bundle.assets` so they copy automatically.

Players and coworkers **only run the folder** — no Koda install required.

Full guide: [Distribution](../distribution.md).

---

## Next steps

| Topic | Guide |
|-------|--------|
| Language tutorial | [Using the language](../using-the-language.md) |
| From C mindset | [Coming from C](from-c.md) |
| Games | [Game development](game-dev.md) |
| All commands | [commands.md](../commands.md) |
