# Modules and imports

How Koda organizes code across files.

---

## Three mechanisms

| Mechanism | Syntax | Effect |
|-----------|--------|--------|
| **Include** | `#include "path.koda"` | Paste file at compile time |
| **Import** | `let m = import "path.koda"` | Load module; get export object |
| **Stdlib import** | `import "@math"` | Load `stdlib/math.koda` or builtin object |

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
