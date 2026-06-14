# Koda language reference

Koda is a dynamically-typed scripting language with JavaScript-like syntax. It compiles to native machine code through LLVM IR and the C runtime (`koda build`, `koda run`).

---

## Variables

```koda
let x = 42;
let name = "Player";
let active = true;
let nothing = null;
let score;           // declared, value is null
```

Assignment:

```koda
x = x + 1;
x += 1;
x -= 1;
x *= 2;
x /= 2;
```

---

## Types

| Type | Examples |
|------|---------|
| Number | `0`, `42`, `3.14`, `-10`, `0xFF` |
| String | `"hello"`, `"line\n"` |
| Boolean | `true`, `false` |
| Null | `null` |
| Array | `[1, 2, 3]` |
| Object | `{x: 10, y: 20}` |
| Function | `func(a, b) { return a + b; }` |

---

## Operators

```koda
// Arithmetic
a + b    a - b    a * b    a / b    a % b

// Comparison
a == b   a != b   a < b   a > b   a <= b   a >= b

// Logical
a && b   a || b   !a

// Prefix / postfix update
++x   --x   x++   x--

// Compound assignment
x += 1   x -= 1   x *= 2   x /= 2
```

---

## Functions

```koda
func add(a, b) {
    return a + b;
}

let result = add(3, 4);   // 7
```

Anonymous function expression:

```koda
let square = func(n) {
    return n * n;
};
```

Closures:

```koda
func makeCounter() {
    let count = 0;
    return func() {
        count += 1;
        return count;
    };
}

let c = makeCounter();
print(c());   // 1
print(c());   // 2
```

---

## Control flow

Braces `{}` are **always** required around `if` / `else` bodies, loops, and `switch` bodies — no single-statement branches without braces.

You may still write a **compact single line** when the body is tiny, or spread it across **multiple lines**:

```koda
if (x > 0) { print("positive"); }

if (x > 0) {
    print("positive");
}

while (running) {
    update();
}

for (let i = 0; i < 10; i += 1) { print(i); }

for (let i = 0; i < 10; i += 1) {
    print(i);
}
```

---

## For loops (classic C form)

Counted loops use `init`, `condition`, and `step`; any part may be omitted (`for (;;)` is valid with `break`).

```koda
for (let i = 0; i < n; i += 1) {
    print(i);
}
```

You can mix this with `while`, `do-while`, `for-in`, and `for-of` in the same codebase—see **`docs/user_guide.md`** (“Choosing a loop style”).

## For-of loops

Single binding: each **value** in **insertion order** (arrays by index, tables by stored slot order).

```koda
let items = ["sword", "shield", "potion"];

for (let item of items) {
    print(item);
}
```

Destructuring **`[indexOrKey, value]`**:

```koda
let tbl = { a: 1, b: 2 };

for (let [k, v] of tbl) {
    print(k, v);
}

let xs = ["x", "y"];

for (let [i, ch] of xs) {
    print(i, ch); // i is numeric index 0, 1 …
}
```

For keys only on objects, **`for-in`** is enough:

```koda
for (let key in tbl) {
    print(key, tbl[key]);
}
```

---

## Arrays

```koda
let arr = [10, 20, 30];

print(arr[0]);          // 10
print(len(arr));        // 3

arr[1] = 99;
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

print(player.name);     // Hero
player.hp = 75;
```

---

## Switch

```koda
switch (state) {
    case "menu":
        renderMenu();
        break;
    case "playing":
        updateGame();
        break;
    default:
        break;
}
```

---

## Modules

```koda
#include "math.koda"
#include "wrappers/raylib/raylib.koda"
```

Resolution order: local path, `KODA_PATH` directories, `KODA_WRAPPERS` directories.

---

## Standard library

| Function | Description |
|----------|-------------|
| `print(...)` | Print values, space-separated |
| `type(v)` | Return type name as string |
| `number(v)` | Convert to number |
| `string(v)` | Convert to string |
| `len(v)` | Length of array or string |
| `time()` | Current time in seconds (float) |
| `sleep(ms)` | Sleep for milliseconds |
| `abs(n)` | Absolute value |
| `sqrt(n)` | Square root |
| `random()` | Random float in [0, 1) |
