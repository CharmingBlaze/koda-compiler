# Chapter 3 — Variables and types

**You will learn:** Koda's types, `let`, reassignment, and `type()`.

**Time:** ~10 minutes.

---

## Declaring variables

Use `let` at the start of a statement (semicolon required after `let` lines):

```koda
let score = 100;
let name = "Player";
let ready = true;
let nothing = null;
```

Reassign with `=` (not `let` again):

```koda
score = score + 10;
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

Numbers are 64-bit floats internally — integers and fractions use the same type.

---

## Strings and escapes

```koda
let line = "Line one\nLine two";
let quote = "She said \"hi\"";
let path = "C:\\games\\save.dat";
```

Supported escapes: `\n`, `\r`, `\t`, `\"`, `\\`, `\'`.

---

## Truthiness

In `if` conditions, these are **falsy**: `false`, `null`, `0`, `""` (empty string).

Everything else is truthy, including `"0"` and empty arrays.

---

## Operators (preview)

```koda
let a = 10 + 3 * 2;     // 16
let ok = score >= 100;
let both = ready && ok;
let msg = "Score: " + string(score);
```

Full operator list: [Language reference — Operators](../language.md#5-operators).

---

## Try it yourself

Write a program that:

1. Stores your name and age in variables.
2. Prints `type()` for each.
3. Prints whether age is 18 or older using `if`.

---

## Next chapter

[Chapter 4 — Control flow](04-control-flow.md)
