# Koda language reference

Koda compiles to native machine code through LLVM IR and the C runtime (`koda build`, `koda run`).

For the full catalog with every builtin, see **[language.md](../../language.md)**. For a learning path, see **[using-the-language.md](../using-the-language.md)**.

---

## Variables

```koda
let x = 42;              // mutable
const gravity = 900;     // immutable
let name = "Player";
let active = true;
let nothing = null;
let score;               // starts as null

let lives: int = 3;      // optional type annotation
let speed: float = 8.0;
let label: string = "Player";
```

Assignment:

```koda
x = x + 1;
x += 1;
x -= 1;
x *= 2;
x /= 2;
```

> `var` is reserved — use `let` or `const`.

---

## Core types

| Type | Examples |
|------|----------|
| `int` | `3`, `let n: int = 0` |
| `float` | `3.14`, `let speed: float = 8.0` |
| `bool` | `true`, `false` |
| `string` | `"hello"`, `"Score: {score}"` |
| `array` | `[1, 2, 3]` |
| `map` / object | `{ x: 10, y: 20 }` |
| `func` | `func(a, b) { return a + b; }` |
| `null` | `null` |

Type annotations are optional. Runtime `type(x)` returns `"number"`, `"string"`, `"bool"`, etc.

---

## String interpolation

```koda
let score = 42;
drawtext("Score: {score}", 20, 20, 24, colors.white);

let msg = `Hello, ${name}!`;   // backtick templates also work
```

---

## Operators

```koda
// Arithmetic
a + b    a - b    a * b    a / b    a % b    a ** b

// Comparison
a == b   a != b   a < b   a > b   a <= b   a >= b

// Logic
a && b   a || b   !a

// Nullish
a ?? b   a ??= b

// Optional chaining
obj?.field   obj?.[key]

// Prefix / postfix update
++x   --x   x++   x--

// Compound assignment
x += 1   x -= 1   x *= 2   x /= 2
```

---

## Control flow

Braces `{}` are **always** required around `if` / `else` bodies, loops, and `switch` / `match` arms.

### if / while / for

```koda
if (x > 0) { print("positive"); }

while (running) { update(); }

for (let i = 0; i < n; i += 1) { print(i); }

for (let item of items) { print(item); }
for coin in coins { coin.update(); }
for i in 0..10 { print(i); }
```

### switch (C-style)

```koda
switch (state) {
    case State.Playing:
        updateGame();
        break;
    default:
        break;
}
```

Cases do **not** fall through. Use `fallthrough;` at the end of a case when you need C-style chaining.

### match (brace-style)

```koda
enum GameState { Playing, Won, GameOver }

match state {
    GameState.Playing {
        update_game(dt);
    }
    GameState.Won {
        draw.text("STAR GET!", 380, 340, 40, colors.yellow);
    }
    default {
        /* optional */
    }
}
```

No fall-through — each arm is its own block.

---

## Functions

```koda
func add(a, b) {
    return a + b;
}

let square = func(n) {
    return n * n;
};
```

Closures capture outer variables:

```koda
func makeCounter() {
    let count = 0;
    return func() {
        count += 1;
        return count;
    };
}
```

---

## Structs and enums

```koda
struct Player { x, y, speed, health }

let p = Player { x: 100, y: 200, speed: 220, health: 100 };

enum State { Idle, Running, Dead }

let state = State.Running;
```

---

## Arrays

```koda
let arr = [10, 20, 30];
print(arr[0]);          // 10
print(len(arr));        // 3
arr.add(40);            // append
arr.remove_at(0);       // remove at index
arr.clear();            // remove all
print(arr.count);       // element count
```

---

## Objects

```koda
let player = {
    name: "Hero",
    hp: 100,
    x: 0,
    y: 0
};

print(player.name);
player.hp = 75;
```

---

## Modules

```koda
#include "math.koda"
#include "wrappers/raylib_shim/raylib.koda"
#include "@game"

let math = import "@math";
```

Resolution order: local path, `KODA_PATH` directories, `KODA_WRAPPERS` directories.

---

## Keywords

```
break  case  const  continue  default  defer  delete  do  else  enum
false  for   func   if       import in      let  match null  of
return struct switch this     true   typeof  while
```

---

## See also

- [Language reference (full)](../language.md)
- [Beginner's guide](../beginners-guide.md)
- [Stdlib](../stdlib/README.md)
