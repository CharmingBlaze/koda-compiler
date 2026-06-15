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

## The seven types

| Type | Example | `type()` result |
|------|---------|-----------------|
| Number | `42`, `3.14`, `0xff`, `1e3` | `"number"` |
| String | `"hello"` | `"string"` |
| Bool | `true`, `false` | `"bool"` |
| Null | `null` | `"null"` |
| Array | `[1, 2]` | `"array"` |
| Object | `{ a: 1 }` | `"object"` |
| Function | `func() {}` | `"function"` |

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
```

Supported escapes: `\n`, `\r`, `\t`, `\"`, `\\`, `\'`.

---

## Truthiness

In `if` conditions, these are **falsy**: `false`, `null`, `0`, `""` (empty string).

Everything else is truthy.

---

## Try it yourself

Write a program that:

1. Stores your name and age in variables.
2. Uses `const` for a maximum score.
3. Prints `type()` for each binding.

---

## Next chapter

[Chapter 4 — Control flow](04-control-flow.md)
