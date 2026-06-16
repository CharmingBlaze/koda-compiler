# Chapter 3 — Variables and types

**You will learn:** Koda's types, `let`, `const`, reassignment, and `type()`.

**Time:** ~10 minutes.

---

## Declaring variables

Use `let` for values that change:

```koda
let score = 100;
let name = "Player";
let ready = true;
let nothing = null;
```

Use `const` for values that never change:

```koda
const gravity = 900;
const screenWidth = 800;
```

Reassign `let` with `=` (not `let` again). You cannot reassign `const`.

```koda
score = score + 10;
```

---

## Optional type annotations

Types are usually inferred. Add them when you want clarity:

```koda
let score = 0;              // inferred number
let lives: int = 3;         // optional
let dt: float = 0.016;      // optional
let name: string = "Jesse"; // optional
```

---

## The core types

| Type | Example | Notes |
|------|---------|-------|
| `int` | `3`, `let n: int = 0` | Integer semantics |
| `float` | `3.14`, `let speed: float = 8.0` | Default for number literals |
| `bool` | `true`, `false` | |
| `string` | `"hello"`, `"Score: {score}"` | See interpolation below |
| `array` | `[1, 2, 3]` | |
| `map` / object | `{ x: 1, y: 2 }` | Plain tables |
| `func` | `func(x) { return x; }` | |
| `null` | `null` | |

Use `type(x)` at runtime — it returns `"number"`, `"string"`, `"bool"`, etc.

Numbers are 64-bit floats by default. Use `int` annotations when you need integer semantics.

---

## Equality

Use `==` and `!=` to compare values:

```koda
if (score == 100) {
    print("perfect");
}
```

There is one equality operator — no `===` / `!==` confusion.

---

## Strings and escapes

```koda
let line = "Line one\nLine two";
let quote = "She said \"hi\"";
let msg = "Score: {score}   Lives: {lives}";
```

Embed values in double-quoted strings with `{expression}` — no manual `"Score: " + score` building.

Supported escapes: `\n`, `\r`, `\t`, `\"`, `\\`, `\'`.

---

## Truthiness

In `if` conditions, these are **falsy**: `false`, `null`, `0`, `""` (empty string).

Everything else is truthy — including **`[]` and `{}`**. This differs from JavaScript, where empty arrays and objects are falsy. Use `len(arr) > 0` instead of `if (arr)` when you mean “non-empty array”.

---

## Try it yourself

Write a program that:

1. Stores your name and age in variables.
2. Uses `const` for a maximum score.
3. Prints `type()` for each binding.

---

## Next chapter

[Chapter 4 — Control flow](04-control-flow.md)
