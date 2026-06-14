# The Koda Language Guide

> **Start here for an accurate, native-only guide:** [using-the-language.md](using-the-language.md) and the exhaustive catalog [language.md](../language.md).  
> This file is a long-form tutorial; some early sections still mention older VM-era names (`koda`, `repl`) — treat **`koda run`** / **`koda build`** and the docs above as the source of truth.

A complete, teachable reference for the Koda programming language — from first
principles through advanced features, standard library, and native compilation.

---

## Table of Contents

1. [What is Koda?](#1-what-is-koda)
2. [Getting Started](#2-getting-started)
3. [Values and Types](#3-values-and-types)
4. [Variables](#4-variables)
5. [Operators](#5-operators)
6. [Strings](#6-strings)
7. [Arrays](#7-arrays)
8. [Objects](#8-objects)
9. [Control Flow](#9-control-flow)
10. [Functions](#10-functions)
11. [Closures](#11-closures)
12. [Modules and Imports](#12-modules-and-imports)
13. [Built-in Functions](#13-built-in-functions)
14. [Standard Library](#14-standard-library)
15. [Methods on Values](#15-methods-on-values)
16. [Template Literals](#16-template-literals)
17. [Ranges](#17-ranges)
18. [Destructuring and Spread](#18-destructuring-and-spread)
19. [Tail-Call Optimisation](#19-tail-call-optimisation)
20. [Graphics and Windowing](#20-graphics-and-windowing)
21. [LLVM Native Compilation](#21-llvm-native-compilation)
22. [Error Reference](#22-error-reference)
23. [Complete Examples](#23-complete-examples)

---

## 1. What is Koda?

Koda is a statically-scoped, dynamically-typed scripting language designed for
performance and expressiveness. It compiles to bytecode that runs in a fast VM,
and can also be lowered to LLVM IR for native binary generation.

**Key design goals:**

- Simple, familiar C-family syntax with no surprises
- First-class functions and closures
- NaN-boxed values for compact, fast representation
- A module system with path-based and stdlib imports
- Optional ahead-of-time compilation via LLVM
- Built-in 2D/3D graphics through Raylib

**Everything in Koda is a value.** Numbers, strings, booleans, null, arrays,
objects, and functions are all first-class values that can be assigned to
variables, passed to functions, and returned from functions.

---

## 2. Getting Started

### Running a script

```
koda run hello.koda
```

### REPL

```
koda repl
```

### Compiling to native binary

```
koda compile myapp.koda -o myapp
```

### Your first program

```koda
print("Hello, Koda!");
```

Output:
```
Hello, Koda!
```

`print` accepts any number of arguments and prints them separated by spaces,
followed by a newline. This is the primary output function.

---

## 3. Values and Types

Koda has six value types:

| Type       | Example             | Notes                            |
|------------|---------------------|----------------------------------|
| `number`   | `42`, `3.14`, `-7`  | All numbers are 64-bit floats    |
| `string`   | `"hello"`, `'hi'`   | UTF-8, immutable                 |
| `bool`     | `true`, `false`     | Logical values                   |
| `null`     | `null`              | The absence of a value           |
| `array`    | `[1, 2, 3]`         | Ordered, mutable, any element    |
| `object`   | `{x: 1, y: 2}`      | Key-value map, string keys       |
| `function` | `func(x) { ... }`   | First-class, can close over vars |

### Checking the type of a value

```koda
print(type(42));        // number
print(type("hello"));   // string
print(type(true));      // bool
print(type(null));      // null
print(type([1, 2]));    // array
print(type({a: 1}));    // object
print(type(print));     // function
```

`type(v)` always returns a **string** describing the type. Use it to branch on
unknown inputs.

### Truthiness

Every value has a truth value in a boolean context:

- `false` and `null` are **falsy**
- Everything else — including `0`, `""`, `[]`, `{}` — is **truthy**

This differs from JavaScript. In Koda, `0` is truthy. Design your guards
accordingly.

---

## 4. Variables

### Declaring a variable

```koda
let x = 10;
let name = "Alice";
let active = true;
let nothing = null;
```

`let` declares a variable and assigns an initial value. Every variable must be
declared before use. Variables are lexically scoped — they are visible from the
point of declaration to the end of the enclosing block.

### Reassigning a variable

```koda
let score = 0;
score = score + 1;   // reassignment
score += 10;         // compound assignment
score++;             // increment
score--;             // decrement
```

**Compound assignment operators:** `+=`, `-=`, `*=`, `/=`, `%=`, `**=`,
`&=`, `|=`, `^=`, `<<=`, `>>=`, `>>>=`

### Scope

```koda
let x = 1;

{
    let x = 2;      // shadows outer x inside this block
    print(x);       // 2
}

print(x);           // 1 — outer x unchanged
```

Blocks `{ ... }` create a new scope. Inner declarations shadow outer ones
without modifying them.

---

## 5. Operators

### Arithmetic

| Operator | Meaning          | Example          |
|----------|------------------|------------------|
| `+`      | Add / concatenate| `3 + 4` → `7`    |
| `-`      | Subtract         | `10 - 3` → `7`   |
| `*`      | Multiply         | `4 * 5` → `20`   |
| `/`      | Divide           | `9 / 2` → `4.5`  |
| `%`      | Remainder        | `10 % 3` → `1`   |
| `**`     | Exponentiation   | `2 ** 10` → `1024`|
| `-x`     | Negate           | `-5`             |

Division always produces a float. There is no integer division operator;
use `floor(a / b)` when you need integer semantics.

### Comparison

```koda
10 == 10       // true
10 != 5        // true
10 > 5         // true
10 >= 10       // true
5  < 10        // true
5  <= 5        // true
```

`==` checks **value equality** for numbers, strings, and booleans.
Arrays and objects compare by **reference**.

```koda
[1,2] == [1,2]   // false — different objects in memory
let a = [1,2];
let b = a;
b == a           // true — same reference
```

For strict equality (no type coercion — Koda has none anyway), `===` is also
available and behaves identically to `==`.

### Logical

```koda
true && false    // false  — logical AND, short-circuits
true || false    // true   — logical OR, short-circuits
!true            // false  — logical NOT
```

`&&` evaluates the right side only when the left is truthy.
`||` evaluates the right side only when the left is falsy.
Both return one of the two original values (not necessarily a bool).

```koda
let name = input() || "Guest";   // default if input is falsy
let safe = obj && obj.value;     // guard null access
```

### Bitwise

```koda
a & b     // AND
a | b     // OR
a ^ b     // XOR
~a        // NOT (bitwise complement)
a << n    // left shift
a >> n    // signed right shift
a >>> n   // unsigned right shift
```

### Ternary

```koda
let label = score > 50 ? "pass" : "fail";
```

The ternary operator `condition ? thenValue : elseValue` is an **expression**,
not a statement. It can appear anywhere a value is expected.

### Operator precedence (high to low)

1. Unary: `!`, `-`, `~`
2. Exponentiation: `**`
3. Multiplicative: `*`, `/`, `%`
4. Additive: `+`, `-`
5. Shift: `<<`, `>>`, `>>>`
6. Bitwise AND: `&`
7. Bitwise XOR: `^`
8. Bitwise OR: `|`
9. Comparison: `<`, `<=`, `>`, `>=`
10. Equality: `==`, `!=`, `===`, `!==`
11. Logical AND: `&&`
12. Logical OR: `||`
13. Ternary: `?:`
14. Assignment: `=`, `+=`, `-=`, ...

---

## 6. Strings

Strings are immutable sequences of UTF-8 characters. They support both `"` and
`'` delimiters.

```koda
let s = "Hello, World!";
let t = 'Single quotes work too';
```

### Concatenation

`+` with at least one string operand concatenates. The other operand is
automatically converted to a string.

```koda
"Hello, " + "World!"         // "Hello, World!"
"Score: " + 100              // "Score: 100"
"Pi is " + 3.14              // "Pi is 3.14"
```

### Template literals (string interpolation)

Use backticks and `${}` to embed expressions:

```koda
let name = "Alice";
let age  = 30;
let msg  = `Hello, ${name}! You are ${age} years old.`;
print(msg);   // Hello, Alice! You are 30 years old.
```

Template literals can contain any expression inside `${}`, including function
calls and arithmetic.

### String indexing

```koda
let s = "hello";
print(s[0]);    // h
print(s[4]);    // o
```

Strings are zero-indexed. Indexing returns a single-character string.

### String length

```koda
print(len("hello"));   // 5
```

### String methods

Methods are called with dot notation:

```koda
"hello".upper()         // "HELLO"
"HELLO".lower()         // "hello"
"  hi  ".trim()         // "hi"
"a,b,c".split(",")      // ["a", "b", "c"]
"hello".contains("ell") // true
"hello".startsWith("he")// true
"hello".endsWith("lo")  // true
"hello".indexOf("l")    // 2
"hello".slice(1, 3)     // "el"
"hi".repeat(3)          // "hihihi"
"abc".replace("b","X")  // "aXc"
```

---

## 7. Arrays

Arrays are ordered, mutable collections of any values.

### Creating arrays

```koda
let empty  = [];
let nums   = [1, 2, 3, 4, 5];
let mixed  = [1, "two", true, null, [3, 4]];
```

### Accessing elements

Arrays are zero-indexed:

```koda
let arr = [10, 20, 30];
print(arr[0]);    // 10
print(arr[2]);    // 30
```

### Modifying elements

```koda
arr[1] = 99;
print(arr);    // [10, 99, 30]
```

### Array length

```koda
print(len([1, 2, 3]));   // 3
```

### Slicing

```koda
let arr = [0, 1, 2, 3, 4];
print(arr[1..3]);    // [1, 2]   — from index 1 up to (not including) 3
print(arr[2..]);     // [2, 3, 4]
print(arr[..3]);     // [0, 1, 2]
```

### Array methods

```koda
let a = [1, 2, 3];

a.push(4);          // [1, 2, 3, 4]  — appends, returns new array
a.pop();            // returns 4, modifies in place
a.includes(2);      // true
a.indexOf(2);       // 1
a.reverse();        // [3, 2, 1]
a.slice(0, 2);      // [1, 2]
a.concat([4, 5]);   // [1, 2, 3, 4, 5]
a.join(", ");       // "1, 2, 3"
```

### Iterating

```koda
let fruits = ["apple", "banana", "cherry"];

for (let i = 0; i < len(fruits); i++) {
    print(fruits[i]);
}

// for-in loop iterates over indices
for (let i in fruits) {
    print(i, fruits[i]);
}
```

---

## 8. Objects

Objects are unordered key-value maps. Keys are always strings.

### Creating objects

```koda
let point = { x: 10, y: 20 };

let user = {
    name: "Alice",
    age:  30,
    active: true
};
```

### Accessing fields

```koda
print(user.name);       // Alice  — dot access
print(user["age"]);     // 30     — bracket access (dynamic keys)

let field = "name";
print(user[field]);     // Alice
```

### Modifying and adding fields

```koda
user.age = 31;
user.city = "London";    // new field created dynamically
```

### Deleting fields

```koda
delete user.city;
```

### Checking field existence

```koda
print(keys(user));    // ["name", "age", "active"]
```

`keys(obj)` returns an array of the object's own key names.

### Nested objects

```koda
let config = {
    server: {
        host: "localhost",
        port: 8080
    },
    debug: true
};

print(config.server.host);    // localhost
print(config.server.port);    // 8080
```

### Objects as namespaces

Objects are how Koda organises related data and behaviour:

```koda
let Vec2 = {
    new: func(x, y) { return {x: x, y: y}; },
    add: func(a, b) { return {x: a.x + b.x, y: a.y + b.y}; },
    len: func(v)    { return sqrt(v.x*v.x + v.y*v.y); }
};

let a = Vec2.new(3, 4);
let b = Vec2.new(1, 1);
let c = Vec2.add(a, b);
print(Vec2.len(a));    // 5
```

---

## 9. Control Flow

Every `if`, `else`, loop, and `switch` body uses **`{ … }` braces** — Koda does not allow braceless one-statement branches.

Within that rule you can format however you like: **one short line** inside the braces, or **several lines** when the body is longer or you want clearer structure.

```koda
// Single-line body (still braced)
if (x > 0) { print("ok"); }

// Multi-line body — same meaning, easier to extend
if (x > 0) {
    print("ok");
}

for (let i = 0; i < 3; i += 1) { print(i); }

for (let i = 0; i < 3; i += 1) {
    print(i);
}
```

`koda fmt` will reflow layout; choose a style that reads well for your team.

### if / else

```koda
let x = 10;

if (x > 5) {
    print("big");
} else if (x == 5) {
    print("five");
} else {
    print("small");
}
```

The condition must be in parentheses. Braces are required on every branch — there are no
braceless one-liners.

### if as an expression

```koda
let label = if (score > 50) { "pass" } else { "fail" };
```

When `if` is used as an expression, each branch must produce a value (the last
expression in the block is the result).

### switch statement

```koda
let day = 3;

switch (day) {
    case 1: print("Monday");
    case 2: print("Tuesday");
    case 3: print("Wednesday");
    default: print("Other");
}
```

Each `case` runs its body and **does not fall through** — no `break` needed.
The `default` branch runs when no case matches.

### switch as an expression

```koda
let name = switch (day) {
    case 1 => "Monday"
    case 2 => "Tuesday"
    case 3 => "Wednesday"
    default => "Other"
};
```

Arrow syntax `=>` makes each arm a value expression.

### while loop

```koda
let i = 0;
while (i < 5) {
    print(i);
    i++;
}
```

Runs the body repeatedly as long as the condition is truthy.

### do-while loop

```koda
let i = 0;
do {
    print(i);
    i++;
} while (i < 5);
```

The body runs **at least once**, then checks the condition.

### Choosing a loop style

You can mix **classic C-style loops** with Koda’s **JavaScript-flavored** forms in the same program:

| Style | Good when |
|-------|-----------|
| **`for (init; cond; step)`** | You want initialization, test, and step written together (counted loops). |
| **`while` / `do-while`** | The condition is simpler than a three-part header, or the step is uneven. |
| **`for-in`** | You iterate **keys** (object fields or array indices as values). |
| **`for-of`** | You iterate **array elements** in order (half-open `[0, len)` indexing). |

Nothing forces one style: pick whichever reads best.

### for loop

```koda
for (let i = 0; i < 10; i++) {
    print(i);
}
```

The classic C-style `for` loop. All three parts are optional:

```koda
for (;;) {    // infinite loop
    // ...
    break;
}
```

Multiple `let` bindings and comma-separated steps are allowed:

```koda
for (let i = 0, let j = 10; i < j; i += 1, j -= 1) {
    print(i, j);
}
```

### for-in loop (iterate over keys)

```koda
let obj = {a: 1, b: 2, c: 3};

for (let key in obj) {
    print(key, obj[key]);
}
```

`for-in` iterates over the **keys** of an object (or indices of an array).

### for-of loop (iterate values)

```koda
let items = ["sword", "shield", "potion"];

for (let item of items) {
    print(item);
}
```

`for-of` walks stored slots **`0 .. len(iterable)-1`** and binds each **value** (arrays: elements; objects/tables: values in insertion order). Order matches how entries were stored in the runtime table.

Bind **key and value** together:

```koda
let tbl = { a: 1, b: 2 };

for (let [k, v] of tbl) {
    print(k, ":", v);
}

let xs = ["x", "y"];

for (let [i, ch] of xs) {
    print(i, ch); // i is the numeric index
}
```

You can still use **`for-in`** when you only need keys, or **`for (let k of ["a","b"])`** when keys are fixed.

Need only indices without destructuring? Use a classic **`for`** over **`len(items)`**.

### break and continue

```koda
for (let i = 0; i < 10; i++) {
    if (i == 3) { continue; }   // skip 3
    if (i == 7) { break; }      // stop at 7
    print(i);
}
// prints: 0 1 2 4 5 6
```

`break` exits the innermost loop immediately.
`continue` skips to the next iteration of the innermost loop.

---

## 10. Functions

### Declaring a function

```koda
func greet(name) {
    print("Hello, " + name + "!");
}

greet("Alice");    // Hello, Alice!
```

`func name(params) { body }` declares a named function at the current scope.

### Returning values

```koda
func add(a, b) {
    return a + b;
}

let result = add(3, 4);
print(result);    // 7
```

`return` exits the function immediately with a value. A function without a
`return` statement (or with bare `return;`) returns `null`.

### Function expressions

Functions are values. They can be assigned to variables and passed around:

```koda
let square = func(x) {
    return x * x;
};

print(square(5));    // 25
```

### Default parameters

Parameters can have default values for when the caller omits them:

```koda
func greet(name, greeting = "Hello") {
    print(greeting + ", " + name + "!");
}

greet("Alice");            // Hello, Alice!
greet("Bob", "Hi");        // Hi, Bob!
```

### Rest parameters

A rest parameter collects all remaining arguments into an array:

```koda
func sum(...nums) {
    let total = 0;
    for (let i = 0; i < len(nums); i++) {
        total += nums[i];
    }
    return total;
}

print(sum(1, 2, 3, 4));    // 10
```

The rest parameter must be last and is prefixed with `...`.

### Arrow functions

A shorthand for single-expression functions:

```koda
let double = (x) => x * 2;
let add    = (a, b) => a + b;

print(double(5));      // 10
print(add(3, 4));      // 7
```

Arrow functions use `=> expr` for the body — no braces, no `return`.

### Passing functions as arguments (higher-order)

```koda
func apply(f, x) {
    return f(x);
}

func triple(n) { return n * 3; }

print(apply(triple, 5));    // 15
print(apply(func(x) { return x + 10; }, 5));    // 15
```

### Recursive functions

```koda
func factorial(n) {
    if (n <= 1) { return 1; }
    return n * factorial(n - 1);
}

print(factorial(10));    // 3628800
```

---

## 11. Closures

A **closure** is a function that remembers the variables from the scope where it
was defined, even after that scope has exited.

### Basic closure

```koda
func makeCounter() {
    let count = 0;
    return func() {
        count = count + 1;
        return count;
    };
}

let counter = makeCounter();
print(counter());    // 1
print(counter());    // 2
print(counter());    // 3
```

Each call to `makeCounter()` creates an independent closure with its own `count`.
The returned function "closes over" the `count` variable.

### Closure over parameters

```koda
func makeAdder(n) {
    return func(x) { return x + n; };
}

let add5  = makeAdder(5);
let add10 = makeAdder(10);

print(add5(3));     // 8
print(add10(3));    // 13
```

### Why closures matter

Closures let you create **factories** — functions that produce customised
behaviour based on their creation-time arguments. They are the foundation of
callbacks, event handlers, and stateful objects in Koda.

---

## 12. Modules and Imports

Koda has a first-class module system. Every `.koda` file is a module — all its
top-level `let` declarations and `func` declarations become exported values.

### Importing a module

```koda
let math = import "@math";
print(math.pi);             // 3.141592653589793
print(math.sin(math.pi));   // ~0
```

The import expression loads the module, runs it, and returns an **object**
whose fields are the module's top-level names.

### Import by path

```koda
let utils = import "utils.koda";        // relative path
let lib   = import "./lib/helpers.koda";
```

### Stdlib imports

Standard library modules use the `@` prefix:

```koda
let math  = import "@math";
let io    = import "@io";
let json  = import "@json";
let str   = import "@str";
let array = import "@array";
```

Set the `KODA_PATH` environment variable to point to the stdlib directory:

```
set KODA_PATH=C:\koda\stdlib    (Windows)
export KODA_PATH=/usr/local/koda/stdlib    (Unix)
```

### Writing a module

Any `.koda` file is a module. Whatever you declare at top level is exported:

```koda
// utils.koda

let VERSION = "1.0";

func clamp(v, lo, hi) {
    if (v < lo) { return lo; }
    if (v > hi) { return hi; }
    return v;
}

func lerp(a, b, t) {
    return a + (b - a) * t;
}
```

Importing:

```koda
let u = import "utils.koda";
print(u.VERSION);          // 1.0
print(u.clamp(15, 0, 10)); // 10
print(u.lerp(0, 100, 0.5));// 50
```

### Include directive

For simple code inclusion (no module object), use `#include`:

```koda
#include <stdlib/math.koda>
#include "helpers.koda"
```

This inlines the file as if you had typed it directly. Unlike `import`, there
is no module object returned.

---

## 13. Built-in Functions

These are always available without any import.

### Output

```koda
print(v1, v2, ...)
```
Prints all arguments separated by spaces, then a newline.

```koda
print("x =", 42, "active:", true);
// x = 42 active: true
```

### Type

```koda
type(value) -> string
```
Returns `"number"`, `"string"`, `"bool"`, `"null"`, `"array"`, `"object"`,
or `"function"`.

### Type predicates

```koda
is_number(v)    // true if v is a number
is_string(v)    // true if v is a string
is_bool(v)      // true if v is a boolean
is_null(v)      // true if v is null
is_array(v)     // true if v is an array
is_object(v)    // true if v is an object
is_function(v)  // true if v is a function
```

These are faster and more idiomatic than `type(v) == "number"`.

### Conversion

```koda
number("42")        // 42         — parse string to number
number("3.14")      // 3.14
number(true)        // 1
number(false)       // 0
string(42)          // "42"
string(true)        // "true"
string(null)        // "null"
```

### Collections

```koda
len(v)     // length of string, array, or object (key count)
keys(obj)  // array of object's own keys
```

### Math

```koda
abs(x)      // absolute value
sqrt(x)     // square root
sin(x)      // sine (radians)
cos(x)      // cosine (radians)
floor(x)    // round down
ceil(x)     // round up
min(a, b)   // smaller of two numbers
max(a, b)   // larger of two numbers
pow(a, b)   // a raised to the power b
random()    // random float in [0, 1)
```

### Time

```koda
clock()     // CPU time used (seconds, high precision)
time()      // wall-clock Unix timestamp (seconds)
sleep(ms)   // pause execution for ms milliseconds
```

### I/O

```koda
input()          // read a line from stdin, returns string
input("Prompt: ")// print prompt then read line

readFile(path)             // read file, returns string
writeFile(path, content)   // write string to file
appendFile(path, content)  // append string to file
fileExists(path)           // returns true/false
createDirectory(path)      // create directory
```

### Garbage collection

```koda
gc()    // trigger garbage collection manually
```

---

## 14. Standard Library

### @math

```koda
let math = import "@math";
```

**Constants:**

```koda
math.pi    // 3.141592653589793
math.e     // 2.718281828459045
math.tau   // 6.283185307179586  (2 * pi)
math.phi   // 1.618033988749895  (golden ratio)
```

**Trigonometry:**

```koda
math.sin(x)         // sine of x (radians)
math.cos(x)         // cosine of x (radians)
math.tan(x)         // tangent of x (radians)
math.asin(x)        // arc sine, result in [-pi/2, pi/2]
math.acos(x)        // arc cosine, result in [0, pi]
math.atan(x)        // arc tangent, result in [-pi/2, pi/2]
math.atan2(y, x)    // two-argument arc tangent (full quadrant)
math.hypot(a, b)    // sqrt(a*a + b*b), numerically stable
```

**Exponential and logarithm:**

```koda
math.exp(x)     // e^x
math.log(x)     // natural logarithm
math.log2(x)    // base-2 logarithm
math.log10(x)   // base-10 logarithm
math.pow(a, b)  // a^b (same as built-in pow)
math.sqrt(x)    // square root
math.abs(x)     // absolute value
```

**Rounding:**

```koda
math.floor(x)   // round toward -infinity
math.ceil(x)    // round toward +infinity
math.round(x)   // round to nearest integer
math.trunc(x)   // truncate toward zero
```

**Utilities:**

```koda
math.min(a, b)           // smaller value
math.max(a, b)           // larger value
math.clamp(v, lo, hi)    // clamp v to [lo, hi]
math.sign(x)             // -1, 0, or 1
math.lerp(a, b, t)       // linear interpolation
math.deg(radians)        // convert radians to degrees
math.rad(degrees)        // convert degrees to radians
math.map(v, a1, b1, a2, b2)  // re-map value from one range to another
```

**Example:**

```koda
let math = import "@math";

let angle = math.rad(45);           // 45 degrees → radians
print(math.sin(angle));             // 0.7071...
print(math.clamp(150, 0, 100));     // 100
print(math.lerp(0, 255, 0.5));      // 127.5
print(math.map(5, 0, 10, 0, 100));  // 50
```

---

### @io

```koda
let io = import "@io";
```

**File operations:**

```koda
io.read(path)              // read entire file, returns string
io.write(path, content)    // write string to file (overwrites)
io.append(path, content)   // append string to file
io.exists(path)            // true if file or directory exists
io.remove(path)            // delete file
io.mkdir(path)             // create directory (and parents)
io.ls(path)                // list directory, returns array of names
```

**Path utilities:**

```koda
io.joinPath("dir", "file.txt")   // "dir/file.txt"  (OS separator)
io.baseName("/path/to/file.txt") // "file.txt"
io.dirName("/path/to/file.txt")  // "/path/to"
io.extName("file.txt")           // ".txt"
```

**Environment:**

```koda
io.env("HOME")    // get environment variable value, or null
```

**Example:**

```koda
let io = import "@io";

io.write("notes.txt", "Hello\nWorld\n");
let content = io.read("notes.txt");
print(content);

let files = io.ls(".");
for (let i = 0; i < len(files); i++) {
    print(files[i]);
}
```

---

### @json

```koda
let json = import "@json";
```

**Stringify** — convert any Koda value to a JSON string:

```koda
json.stringify({name: "Koda", version: 1.2})
// '{"name":"Koda","version":1.2}'

json.stringify([1, true, null, "hi"])
// '[1,true,null,"hi"]'

json.stringify({name: "Koda"}, 2)   // second arg = indent spaces
// {
//   "name": "Koda"
// }
```

**Parse** — convert a JSON string into Koda values:

```koda
let obj = json.parse('{"x": 1, "y": [2, 3]}');
print(obj.x);       // 1
print(obj.y[0]);    // 2
```

**tryParse** — parse without throwing on bad input:

```koda
let result = json.tryParse("{bad json}");
if (result.error != null) {
    print("parse failed:", result.error);
} else {
    print(result.value);
}
```

**Round-trip example:**

```koda
let json = import "@json";

let data = {
    name: "Koda",
    version: 1.2,
    features: ["fast", "native", "clean"],
    active: true,
    meta: null
};

let s = json.stringify(data);
let p = json.parse(s);
print(p.name);         // Koda
print(p.features[0]);  // fast
```

---

### @str

```koda
let str = import "@str";
```

```koda
str.upper("hello")          // "HELLO"
str.lower("HELLO")          // "hello"
str.trim("  hi  ")          // "hi"
str.split("a,b,c", ",")     // ["a", "b", "c"]
str.join(["a","b","c"], "-") // "a-b-c"
str.charCode("A")           // 65
```

---

### @array

```koda
let array = import "@array";
```

```koda
array.sort([3, 1, 2])    // [1, 2, 3]  — returns sorted copy
```

---

## 15. Methods on Values

Koda supports dot-method calls directly on values. These delegate to built-in
implementations based on the type.

### String methods

```koda
"hello".upper()              // "HELLO"
"HELLO".lower()              // "hello"
"  trim me  ".trim()         // "trim me"
"a,b,c".split(",")           // ["a", "b", "c"]
"hello".contains("ell")      // true
"hello".startsWith("hel")    // true
"hello".endsWith("llo")      // true
"hello".indexOf("l")         // 2
"hello".slice(1, 4)          // "ell"
"abc".replace("b", "B")      // "aBc"
"ha".repeat(3)               // "hahaha"
```

### Array methods

```koda
let a = [1, 2, 3];
a.push(4)           // returns [1, 2, 3, 4]
a.pop()             // removes and returns last element
a.includes(2)       // true
a.indexOf(3)        // 2
a.reverse()         // [3, 2, 1]
a.slice(0, 2)       // [1, 2]
a.concat([4, 5])    // [1, 2, 3, 4, 5]
a.join(", ")        // "1, 2, 3"
```

Methods can be chained:

```koda
"abc".upper().split("")   // ["A", "B", "C"]
```

---

## 16. Template Literals

Template literals use backtick syntax and support embedded expressions:

```koda
let x = 42;
let s = `The answer is ${x}`;
print(s);    // The answer is 42
```

Any expression works inside `${}`:

```koda
let a = 3;
let b = 4;
print(`Hypotenuse: ${sqrt(a*a + b*b)}`);   // Hypotenuse: 5

let items = ["apple", "banana"];
print(`First item: ${items[0]}`);
print(`Count: ${len(items)}`);
```

Multi-word expressions:

```koda
let user = {name: "Alice", score: 95};
print(`${user.name} scored ${user.score > 90 ? "A" : "B"}`);
// Alice scored A
```

---

## 17. Ranges

The `..` operator creates a range, which Koda automatically expands to an array:

```koda
let r = 0..5;
print(r);    // [0, 1, 2, 3, 4]   — exclusive end
```

Ranges are useful with for loops:

```koda
for (let i in 0..10) {
    print(i);
}
```

### Slice syntax

Use ranges to slice arrays and strings:

```koda
let arr = [10, 20, 30, 40, 50];
print(arr[1..3]);    // [20, 30]
print(arr[..2]);     // [10, 20]
print(arr[3..]);     // [40, 50]
```

```koda
let s = "Hello, World!";
print(s[7..12]);   // World
```

---

## 18. Destructuring and Spread

### Spread in function calls

The `...` prefix spreads an array as individual arguments:

```koda
func add(a, b, c) { return a + b + c; }

let nums = [1, 2, 3];
print(add(...nums));    // 6
```

### Rest in function parameters

```koda
func first(head, ...rest) {
    print("head:", head);
    print("rest:", rest);
}

first(1, 2, 3, 4);
// head: 1
// rest: [2, 3, 4]
```

---

## 19. Tail-Call Optimisation

Koda optimises **tail calls** — recursive calls that are the last operation in
a function. This allows unbounded recursion without stack overflow:

```koda
func loop(n) {
    if (n == 0) { return "done"; }
    return loop(n - 1);    // tail call — no stack growth
}

print(loop(1000000));    // "done" — no stack overflow
```

A call is a tail call when its result is returned directly with no further
computation. Both direct recursion and mutual recursion are optimised.

```koda
// Accumulator pattern for tail-recursive factorial
func factTail(n, acc = 1) {
    if (n <= 1) { return acc; }
    return factTail(n - 1, n * acc);    // tail call
}

print(factTail(20));    // 2432902008176640000
```

---

## 20. Graphics and Windowing

Koda includes native Raylib bindings for 2D and 3D graphics. These are
available as global built-in functions — no import required.

### Window lifecycle

```koda
initWindow(800, 600, "My App");
setTargetFPS(60);

while (!windowShouldClose()) {
    beginDrawing();
    clearBackground(0x181818FF);

    // draw here

    endDrawing();
}

closeWindow();
```

### 2D drawing

```koda
drawText("Hello!", 100, 100, 24, 0xFFFFFFFF);
drawRectangle(x, y, width, height, color);
drawCircle(cx, cy, radius, color);
```

Colors are 32-bit RGBA hex integers: `0xRRGGBBAA`.

### 3D drawing

```koda
let camera = {
    px: 0, py: 5, pz: 10,   // position
    tx: 0, ty: 0, tz: 0,    // target
    ux: 0, uy: 1, uz: 0,    // up vector
    fov: 45
};

beginMode3D(camera);

drawGrid(20, 1.0);
drawCube(0, 1, 0, 2, 2, 2, 0x00AAFFFF);
drawCubeWires(0, 1, 0, 2.1, 2.1, 2.1, 0xFFFFFFFF);
drawLine3D(0, 0, 0, 5, 0, 0, 0xFF0000FF);

endMode3D();
```

### Input

```koda
isKeyPressed(KEY_SPACE)    // true on the frame the key is pressed
isKeyDown(KEY_LEFT)        // true while key is held
```

### Complete 3D example

```koda
initWindow(1024, 768, "Koda 3D");
setTargetFPS(60);

let angle = 0;
let camera = {px:0, py:5, pz:10, tx:0, ty:0, tz:0, ux:0, uy:1, uz:0, fov:45};

while (!windowShouldClose()) {
    angle = angle + 0.02;
    camera.px = 10 * sin(angle);
    camera.pz = 10 * cos(angle);

    beginDrawing();
    clearBackground(0x181818FF);
    beginMode3D(camera);

    drawGrid(20, 1.0);
    drawCube(0, 1, 0, 2, 2, 2, 0x00AAFFFF);

    endMode3D();

    drawText("Koda 3D", 10, 10, 20, 0xFFFFFFFF);
    endDrawing();
}

closeWindow();
```

---

## 21. LLVM Native Compilation

Koda can compile scripts to native machine code via LLVM IR.

### Compiling

```
koda compile myapp.koda -o myapp
```

This produces a native executable. The output runs without the Koda runtime.

### What gets compiled

The LLVM backend handles:

- All arithmetic and logic
- Variables, closures, upvalues
- Function calls and tail calls
- Arrays and objects
- All control flow (if, while, for, switch)
- String operations
- Math intrinsics (`sin`, `cos`, `sqrt`, `floor`, `ceil`, `abs`, `pow`, `min`,
  `max`, `tan`, `asin`, `acos`, `atan`, `atan2`, `log`, `log2`, `log10`,
  `exp`, `round`, `trunc`, `hypot`)
- I/O (`readFile`, `writeFile`)
- Graphics (all Raylib bindings)

### Extern declarations for FFI

Use `#extern` to call native C functions from compiled code:

```koda
// #extern double my_c_func(double x);

let result = my_c_func(3.14);
```

The extern directive tells the compiler the function signature so it can
generate the correct call.

---

## 22. Error Reference

### Runtime errors

| Error message                   | Cause                                  |
|---------------------------------|----------------------------------------|
| `undefined global variable`     | Reading an undeclared variable         |
| `expected number, got X`        | Math operation on non-number           |
| `expected string, got X`        | String operation on non-string         |
| `index out of range`            | Array index beyond bounds              |
| `cannot index type X`           | Using `[]` on a non-indexable type     |
| `not callable`                  | Calling something that is not function |
| `wrong number of arguments`     | Arity mismatch on native function      |

### Lexer errors

| Error message                         | Cause                            |
|---------------------------------------|----------------------------------|
| `unexpected character`                | Unknown character in source      |
| `unterminated string`                 | String literal missing closing `"`|

### Compiler errors

| Error message                         | Cause                               |
|---------------------------------------|-------------------------------------|
| `variable already declared`           | Redeclaring `let` in same scope     |
| `return outside function`             | `return` at top level               |
| `break outside loop`                  | `break` outside a loop              |
| `continue outside loop`               | `continue` outside a loop           |

---

## 23. Complete Examples

### Fibonacci with memoisation

```koda
let cache = {};

func fib(n) {
    let k = string(n);
    if (cache[k] != null) { return cache[k]; }
    if (n <= 1) { return n; }
    let result = fib(n - 1) + fib(n - 2);
    cache[k] = result;
    return result;
}

for (let i = 0; i <= 20; i++) {
    print(`fib(${i}) = ${fib(i)}`);
}
```

### Linked list

```koda
func node(value, next = null) {
    return {value: value, next: next};
}

func prepend(list, value) {
    return node(value, list);
}

func printList(list) {
    let cur = list;
    while (cur != null) {
        print(cur.value);
        cur = cur.next;
    }
}

let list = null;
list = prepend(list, 3);
list = prepend(list, 2);
list = prepend(list, 1);
printList(list);   // 1, 2, 3
```

### JSON config loader

```koda
let io   = import "@io";
let json = import "@json";

func loadConfig(path) {
    if (!fileExists(path)) {
        return {debug: false, port: 8080};
    }
    let content = io.read(path);
    return json.parse(content);
}

let cfg = loadConfig("config.json");
print(`Running on port ${cfg.port}`);
```

### Functional pipeline

```koda
func map(arr, f) {
    let result = [];
    for (let i = 0; i < len(arr); i++) {
        result = result.push(f(arr[i]));
    }
    return result;
}

func filter(arr, pred) {
    let result = [];
    for (let i = 0; i < len(arr); i++) {
        if (pred(arr[i])) {
            result = result.push(arr[i]);
        }
    }
    return result;
}

func reduce(arr, f, init) {
    let acc = init;
    for (let i = 0; i < len(arr); i++) {
        acc = f(acc, arr[i]);
    }
    return acc;
}

let nums = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10];

let result = reduce(
    map(
        filter(nums, (n) => n % 2 == 0),
        (n) => n * n
    ),
    (acc, n) => acc + n,
    0
);

print(result);   // 4 + 16 + 36 + 64 + 100 = 220
```

### State machine

```koda
func makeStateMachine(initial) {
    let state       = initial;
    let transitions = {};

    return {
        on: func(from, event, to, action) {
            let key = from + ":" + event;
            transitions[key] = {to: to, action: action};
        },
        send: func(event) {
            let key   = state + ":" + event;
            let trans = transitions[key];
            if (trans == null) {
                print(`No transition from ${state} on ${event}`);
                return;
            }
            if (trans.action != null) { trans.action(); }
            state = trans.to;
        },
        getState: func() { return state; }
    };
}

let door = makeStateMachine("closed");
door.on("closed", "open",  "open",   func() { print("Opening door"); });
door.on("open",   "close", "closed", func() { print("Closing door"); });
door.on("closed", "lock",  "locked", func() { print("Locking door"); });
door.on("locked", "unlock","closed", func() { print("Unlocking door"); });

door.send("open");      // Opening door
door.send("close");     // Closing door
door.send("lock");      // Locking door
door.send("unlock");    // Unlocking door
print(door.getState()); // closed
```

---

*This guide covers the complete Koda language as of the current release.*
*For the latest changes, see the changelog and test suite in the repository.*
