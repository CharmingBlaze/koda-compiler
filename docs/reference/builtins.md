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

Random: `random()`, `randomint(lo, hi)`, `randomchoice(arr)`, `randomseed(n)`.

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
| `listdir(path)` | array |

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
| `gc()` / `gccollect()` | Full collection |
| `gcdisable()` / `gcenable()` | Toggle GC |
| `gcframestep()` | Incremental step (per frame in games) |
| `gcstats()` | Collector stats |

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
