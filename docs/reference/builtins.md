# Built-in globals reference

Global functions and constants available without import. Names are **case-insensitive**.

For grouped APIs see [stdlib](../stdlib/README.md).

---

## Output and debugging

| Name | Purpose |
|------|---------|
| `print(…)` | Print to stdout |
| `warn(…)` | Warning to stderr |
| `assert(cond, msg?)` | Panic if false |
| `panic(msg)` | Stop with message |
| `trace()` | Stack trace |
| `type(x)` / `typeof(x)` | Type name string |

---

## Process environment

| Name | Purpose |
|------|---------|
| `args()` | Array of command-line arguments (includes program name at index 0) |
| `env(name)` | Environment variable value, or `null` if unset |
| `rgb(r, g, b)` | Raylib color from 0–255 channels (alpha 255) |
| `rgba(r, g, b, a)` | Raylib color from 0–255 RGBA channels |
| `vec2(x, y)` | 2D vector with `+`, `-`, `*`, `+=` |
| `vec3(x, y, z)` | 3D vector with `+`, `-`, `*`, `+=` |
| `color(r, g, b)` | RGBA color object (`.packed` for Raylib) |
| `rect(x, y, w, h)` | Rectangle bounds |
| `box(center, size)` | 3D AABB from two `vec3` values |

Named palette colors are on the global **`colors`** object: `colors.white`, `colors.sky`, `colors.grass`, etc. See [game-types.md](../stdlib/game-types.md).

```koda
game.clear(colors.sky);
game.rect(0, 0, 40, 40, rgb(34, 139, 34));
let custom = rgba(255, 216, 168, 255);
```

```koda
func main() {
    for (let i = 0; i < len(args()); i = i + 1) {
        print("arg", i, ":", args()[i]);
    }
    let home = env("USERPROFILE");  // Windows; use HOME on Unix
    if (home != null) {
        print("home:", home);
    }
}
```

Run with extra arguments: `koda run tool.koda -- input.txt output.txt`

---

## Time

| Name | Purpose |
|------|---------|
| `deltatime()` | Seconds since last frame |
| `programtime()` | Seconds since start |
| `time()` | Wall clock |
| `clock()` | CPU/monotonic time |
| `timestamp()` | Unix timestamp |
| `sleep(ms)` | Block for milliseconds |

---

## Math (globals)

`sin`, `cos`, `tan`, `asin`, `acos`, `atan`, `atan2`, `pow`, `exp`, `log`, `log2`, `log10`, `sqrt`, `cbrt`, `abs`, `floor`, `ceil`, `round`, `trunc`, `sign`, `min`, `max`, `hypot`, `fmod`, `lerp`, `clamp`, `wrap`, `approach`, `smoothstep`, `distance`, `distancesq`, `normalize`, `map`, `degrees`, `radians`, `smoothdamp`, `pi`, `e`

Random: `random()`, `randomint(lo, hi)`, `randomchoice(arr)` (pick a random element), `randomseed(n)`.

Prefer `import "@math"` for namespace grouping.

---

## Strings and conversion

| Name | Purpose |
|------|---------|
| `string(x)` | To string |
| `number(x)` | To number |
| `len(x)` | Length of string/array/object |
| `keys(obj)` | Key list |
| `format(…)` | Formatted print helper |

String **methods**: `.trim()`, `.toupper()`, `.split()`, etc. — see [str stdlib](../stdlib/str.md).

---

## Files

| Name | Returns |
|------|---------|
| `readfile(path)` | `{ ok, value, error }` |
| `writefile(path, text)` | `{ ok, value, error }` |
| `appendfile(path, text)` | bool |
| `fileexists(path)` | bool |
| `deletefile(path)` | bool |
| `isfile(path)` | bool |
| `isdir(path)` | bool |
| `filesize(path)` | number |
| `listdir(path)` | array of **entry names** (not full paths; like Node.js `readdir`) |
| `readdir(path)` | alias of `listDir` |

Or `import "@io"`.

---

## JSON

| Name | Purpose |
|------|---------|
| `parsejson(text)` | Parse with ok/err wrapper |
| `tojson(value)` | Stringify |
| `json.parse` via import | Raw parse |

See [json stdlib](../stdlib/json.md).

---

## Type predicates

`isnumber`, `isstring`, `isbool`, `isnull`, `isarray`, `isobject`, `isfunction`, `bool(x)`

---

## Garbage collection

| Name | Purpose |
|------|---------|
| `gc()` | Full collection (canonical) |
| `gcCollect()` | deprecated alias of `gc()` |
| `gcDisable()` / `gcEnable()` | Toggle GC |
| `gcFrameStep(ms)` | Incremental step (per frame in games; budget in ms) |
| `gcStats()` | Collector stats |
| `arena(bytes)` | Create bump allocator |
| `arenaReset(arena)` | Reset arena for next frame |
| `arenaAllocArray(arena, cap)` | Array inside arena |
| `arenaAllocStruct(arena, fields)` | Struct inside arena |

---

## Result helpers

| Name | Purpose |
|------|---------|
| `ok(value)` | `{ ok: true, value, error: null }` |
| `err(message)` | `{ ok: false, value: null, error }` |

---

## Related

- [Language reference](../../language.md)
- [stdlib overview](../stdlib/README.md)
