# Koda — Syntax Reference

> **Full specification** — for a compact single-page reference, see [`language.md`](../language.md) in the root.

This is the **compact syntax reference** for Koda. For full explanations with examples, see the root **`language.md`**. For a step-by-step beginner guide, see **`using-the-language.md`**.

> The compiler sources are the ground truth: **`internal/lexer/token.go`**, **`internal/parser/`**, **`internal/codegen/builtin_register.go`**.

---

## Program shape

```koda
// Top-level statements (script style)
print("hi");

// Or define func main() as the entry point
func main() {
    print("hi");
}
```

Statements end with **`;`**. All keywords and builtin names are **case-insensitive**.

---

## Variables

```koda
let x = 10;          // declare and assign (mutable)
const gravity = -22; // immutable — cannot reassign
let y;               // declare as null
x = x + 1;          // reassign let only
x += 1;  x++;       // shorthand

let lives: int = 3;  // optional type annotation

let { a, b } = obj;  // object destructuring
```

> `var` is reserved — always use `let`.

---

## Types

| Type | Literal |
|------|---------|
| Number | `42`, `3.14`, `1e3`, `0xff`, `0b1010` |
| String | `"hello"`, escapes: `\n \t \" \\`, unicode: `\u{1F600}` |
| Bool | `true`, `false` |
| Null | `null` |
| Array | `[1, 2, 3]`, spread: `[0, ...arr, 4]` |
| Object | `{ x: 1, y: 2 }`, method shorthand: `name() { … }` |
| Function | `func(x) { return x; }` |
| Template string | `` `value = ${ expr }` `` |

---

## Operators

**Arithmetic:** `+` `-` `*` `/` `%` `**`

**Bitwise:** `&` `|` `^` `~` `<<` `>>` `>>>`

**Comparison:** `<` `<=` `>` `>=` `==` `!=`

> **Note:** `===` and `!==` are **not supported**. Use `==` / `!=` only; run `koda fmt` to migrate old sources.

**Logic:** `&&` `||` `!`

**Unary:** `+x` `-x` `!x` `typeof x` `++x` `--x` `x++` `x--`

**Assignment:** `=` `+=` `-=` `*=` `/=` `%=` `&=` `|=` `^=` `<<=` `>>=` `??=`

**Nullish coalescing:** `a ?? b` — uses `b` only when `a` is `null`

**Optional chaining:** `obj?.prop`, `obj?.[expr]` — yields `null` if receiver is `null`

**Range:** `lo..hi` — integer sequence; used in `for (let i of lo..hi)`

**Spread:** `...arr` inside array literals; `...name` as last function param (rest)

---

## Control flow

### `if`

```koda
if (cond) {
    // ...
} else if (other) {
    // ...
} else {
    // ...
}

// if as an expression
let x = if (n > 0) { 1 } else { -1 };
```

### `while` and `do`…`while`

```koda
while (cond) { /* ... */ }

do { /* ... */ } while (cond);
```

### `for` (C-style)

```koda
for (let i = 0; i < 10; i += 1) { /* ... */ }
for (;;) { break; }                             // infinite
```

### `for`…`in` (keys)

```koda
for (let key in obj) { /* key is string */ }
for (let i in arr)   { /* i is numeric index */ }
```

### `for`…`in` (values and ranges)

```koda
for coin in coins { /* each element */ }
for i in 0..10    { /* 0,1,…,9 half-open */ }
```

Legacy parenthesized form still works: `for (let v of arr)`.

### `for`…`of` with pairs

```koda
for (let [k, v] of obj) { /* k=key, v=value */ }
for (let [i, v] of arr) { /* i=index, v=element */ }
```

### `switch` (statement)

```koda
switch (x) {
    case 1:
        /* ... */
        break;
    case 2:
        /* ... */
        break;
    default:
        /* ... */
}
```

Falls through unless `break` is used.

### `match`

```koda
match state {
    GameState.Playing {
        update_game(dt);
    }
    GameState.Won {
        drawtext("STAR GET!", 380, 340, 40, colors.yellow);
    }
    default {
        /* ... */
    }
}
```

No fall-through — each arm is an isolated block. See [String interpolation](#string-interpolation) and [Match](#match) below for full examples.

### `switch` (expression)

```koda
let label = switch (x) {
    case 1 => "one"
    case 2 => "two"
    default => "other"
};
```

### `break` / `continue` / `return`

```koda
break;          // exit loop or switch
continue;       // next loop iteration
return;         // return null from function
return expr;    // return value
```

### `defer`

```koda
func f() {
    defer cleanup();  // runs when f() exits, LIFO order
}
```

### `delete`

```koda
delete obj.key;      // remove own property from object
delete obj["key"];   // bracket form also works
```

---

## Functions

```koda
// Named declaration
func add(a, b) {
    return a + b;
}

// Function expression
let mul = func(a, b) { return a * b; };

// Default parameter
func greet(name = "world") { print(name); }

// Rest parameter (must be last)
func sum(...nums) { /* nums is an array */ }
```

`this` is bound to the receiver for `obj.method()` calls.

---

## Struct types

```koda
struct Coin {
    x = 0.0;
    z = 0.0;
    on = true;

    func new(x, z) {
        this.x = x;
        this.z = z;
    }

    func draw() {
        if on {
            draw.cube(x, 1.5, z, 0.45, 0.45, 0.12, colors.gold);
        }
    }
}

let c = Coin { x: 7, z: -5 };   // on defaults to true
let c2 = Coin(7, -5);           // constructor via func new

for coin in coins {
    coin.draw();
}
```

Inside struct methods, bare field names (`x`, `on`) refer to `this.x`, `this.on`. Use `this.field` when you need to be explicit.

---

## Enum types

```koda
enum Dir {
    Up,       // 0
    Down,     // 1
    Left,     // 2
    Right     // 3
}

let d = Dir.Up;   // 0
```

---

## Objects

```koda
let o = { x: 1, y: 2 };
o.x;          // dot access
o["x"];       // bracket access (identical)
o.z = 3;      // add/update property
len(o);       // number of keys

// Method shorthand
let obj = {
    value: 10,
    get() { return this.value; }
};
```

---

## Arrays

```koda
let a = [1, 2, 3];
a[0];              // read
a[0] = 99;         // write
len(a);            // count
a.count;           // element count (property)
a.add(4);          // append (alias for push)
a.push(4);         // append
a.pop();           // remove last
a.remove_at(0);    // remove at index
a.clear();         // remove all elements

for x in a { /* each element */ }
a.each(func(x) { /* callback per element */ });
```

---

## String interpolation

Embed values in double-quoted strings with `{expression}`:

```koda
drawtext("Score: {score}", 20, 20, 24, colors.white);
drawtext("Lives: {lives}", 20, 50, 24, colors.white);
```

Backtick templates also work: `` `Hello, ${name}!` ``.

---

## Core types

| Type | Example |
|------|---------|
| `int` | `3`, `let n: int = 0` |
| `float` | `3.14`, `let speed: float = 8.0` |
| `bool` | `true`, `false` |
| `string` | `"hello"`, `"Score: {score}"` |
| `array` | `[1, 2, 3]` |
| `map` | `{ x: 1, y: 2 }` |
| `func` | `func(x) { return x; }` |
| `object` | plain `{ key: val }` tables |

Types are optional — beginners can omit them. Use `let` for mutable values, `const` for fixed ones.

---

## Constants

```koda
const sw = 1280;
const sh = 720;
const gravity = -22.0;

let score = 0;   // changes during play
```

---

## Enums

```koda
enum GameState {
    Playing,
    Won,
    GameOver
}

let state = GameState.Playing;
```

Members are integers from `0` upward. Use `Type.Member` (e.g. `GameState.Won`).

---

## Match

Cleaner than long `if` / `else if` chains for states:

```koda
match state {
    GameState.Playing {
        update_game(dt);
    }
    GameState.Won {
        drawtext("STAR GET!", 380, 340, 40, colors.yellow);
    }
    GameState.GameOver {
        drawtext("GAME OVER - press R", 400, 340, 36, colors.red);
    }
}
```

Classic `switch (x) { case …: … }` still works.

---

## Template strings

```koda
let s = `Hello, ${ name }! Score: ${ score * 2 }`;
```

---

## Includes and imports

```koda
#include "relative/path.koda"    // textual include (most common)
let m = import("./module.koda"); // expression form
```

---

## Native FFI hint

```koda
// koda: extern bindingName c_symbol arity
let bindingName;
```

---

## Truthy / falsy

**Falsy:** `false`, `null`, `0`, `""`  
**Truthy:** everything else (including `[]` and `{}`)

---

## Precedence (high → low)

`call .member [index] ++ --` → `typeof ! + - ++ --` (prefix) → `**` → `* / %` → `+ -` → `<< >> >>>` → `< <= > >=` → `== !=` → `&` → `^` → `|` → `&&` → `|| ??` → `= += …` (assignment)

---

## Keywords

```
break  case  const  continue  default  defer  delete  do  else  enum
false  for   func      if       import in      let match null  of
return struct switch   this     true   typeof  while
```

`#include` is a directive. `var` is reserved (use `let`).

---

## Built-in names (all case-insensitive)

See **`../language.md`** section 23 for the full grouped list, or **`internal/codegen/builtin_register.go`** for the definitive source.
