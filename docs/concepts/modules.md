# Modules and imports

How Koda organizes code across files.

---

## Three mechanisms

| Mechanism | Syntax | Effect |
|-----------|--------|--------|
| **Use** | `use raylib;` / `use koda.math;` | Official module import (expanded at compile time) |
| **Include** | `#include "path.koda"` | Paste file at compile time (legacy / shims) |
| **Import** | `let m = import "path.koda"` | Load module; get export object |
| **Stdlib import** | `import "@math"` or `use koda.math` | Load `stdlib/math.koda` or builtin object |

---

## Use (recommended)

Official style for wrappers and stdlib:

```koda
use koda.math;
use koda.input;

func main() {
    let a = abs(-3);
    print(a);
}
```

```koda
use raylib;   // full wrapper when installed under wrappers/raylib/

func main() {
    InitWindow(800, 600, "Hello");
    defer CloseWindow();
    …
}
```

Resolution order for `use name`:

1. `koda.*` → stdlib (`use koda.math` → `@math`)
2. `wrappers/` (e.g. `raylib` → `wrappers/raylib/raylib.koda`)
3. Stdlib shorthand (`use math` → `@math`)

Unknown modules report searched paths. Legacy `#include` and `import "koda.game"` still work.

---

## Include

Best for:

- Small helpers used in one project
- Stdlib `.koda` files (`timer.koda`, `array.koda`)
- No separate namespace — functions become globals in the including file

```koda
#include "stdlib/timer.koda"
let cd = cooldown(0.5);
```

---

## Import

Best for:

- Clear boundaries between game systems
- Returning a single object of exports

```koda
let physics = import "src/physics.koda";
physics.step(world, deltatime());
```

The loaded file's exports become properties on the returned object.

---

## `@` stdlib modules

| Import | Source |
|--------|--------|
| `@math`, `@json`, `@io` | Compiler-built export objects + `stdlib/*.koda` |
| `@array`, `@timer`, … | Primarily `stdlib/*.koda` |

Dot calls like `math.sin(x)` route to native functions when the binding is named `math`, `json`, or `io`.

---

## Related

- [Learn — Modules](../learn/08-modules-and-imports.md)
- [Stdlib overview](../stdlib/README.md)
