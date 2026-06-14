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
let x = 10;          // declare and assign
let y;               // declare as null
x = x + 1;          // reassign
x += 1;  x++;       // shorthand

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

**Comparison:** `<` `<=` `>` `>=` `==` `!=` `===` `!==`

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

### `for`…`of` (values)

```koda
for (let v of arr)      { /* each element */ }
for (let i of 0..10)    { /* 0,1,…,9 */ }
```

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
struct Point {
    x,
    y
}

let p = Point { x: 3, y: 4 };
p.y = 10;
```

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
a[0];           // read
a[0] = 99;      // write
len(a);         // count
a.push(4);      // append
a.pop();        // remove last
```

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

`call .member [index] ++ --` → `typeof ! + - ++ --` (prefix) → `**` → `* / %` → `+ -` → `<< >> >>>` → `< <= > >=` → `== != === !==` → `&` → `^` → `|` → `&&` → `|| ??` → `= += …` (assignment)

---

## Keywords

```
break  case  continue  default  defer  delete  do  else  enum
false  for   func      if       import in      let null  of
return struct switch   this     true   typeof  while
```

`#include` is a directive. `var` is reserved (use `let`).

---

## Built-in names (all case-insensitive)

See **`../language.md`** section 23 for the full grouped list, or **`internal/codegen/builtin_register.go`** for the definitive source.
