# Koda Language Reference

> **Quick reference** — compact catalog of every syntax form. For learning path and guides, see **[docs/README.md](docs/README.md)**.

Koda is a **modern language for games and applications** that compiles to a **native binary** — a practical alternative to C for gameplay logic, tools, and desktop apps when you want one executable and optional C library interop without a VM.

It combines **C-style** structs, enums, and control flow with **JavaScript-style** objects and closures. This page covers **everything you can write** with examples.

> **New to Koda?** [Beginner's guide](docs/beginners-guide.md) · [Learn path](docs/learn/README.md) · [Using the language](docs/using-the-language.md) · this page as lookup.  
> **From C?** [Coming from C](docs/guides/from-c.md).  
> **CLI:** [CLI reference](docs/reference/cli.md) · [commands.md](docs/commands.md).

**Case-insensitive:** all keywords and builtin names are case-insensitive. `Print`, `PRINT`, and `print` all work.

---

## Table of contents

1. [Your first program](#1-your-first-program)
2. [Comments](#2-comments)
3. [Values and types](#3-values-and-types)
4. [Variables](#4-variables)
5. [Operators](#5-operators)
6. [Control flow](#6-control-flow)
7. [Functions](#7-functions)
8. [Closures](#8-closures)
9. [Objects](#9-objects)
10. [Arrays](#10-arrays)
11. [Struct types](#11-struct-types)
12. [Enum types](#12-enum-types)
13. [String methods](#13-string-methods)
14. [Array methods](#14-array-methods)
15. [The math object](#15-the-math-object)
16. [Built-in functions](#16-built-in-functions)
17. [Template strings](#17-template-strings)
18. [Includes and imports](#18-includes-and-imports)
19. [Native FFI hint](#19-native-ffi-hint)
20. [Truthy and falsy](#20-truthy-and-falsy)
21. [Operator precedence](#21-operator-precedence)
22. [All keywords](#22-all-keywords)
23. [All built-in names](#23-all-built-in-names)

---

## 1. Your first program

Save this as `hello.koda` and run it with `koda run hello.koda`:

```koda
print("Hello, Koda!");
```

You can also wrap code in a `func main()`:

```koda
func main() {
    print("Hello from main!");
}
```

Both styles work. `func main()` is useful when you want a clear entry point.

---

## 2. Comments

```koda
// This is a line comment.

/*
   This is a
   block comment.
*/

let x = 1; // comment at end of line
```

---

## 3. Values and types

Koda has seven types. Everything is a value — you can pass any of these to functions, store them in variables, and put them in arrays.

| Type | Example | Notes |
|------|---------|-------|
| **Number** | `42`, `3.14`, `0xff`, `0b1010`, `1e3` | All numbers are 64-bit floats |
| **String** | `"hello"` | UTF-8; use `\n`, `\t`, `\"`, `\\` for escapes |
| **Bool** | `true`, `false` | |
| **Null** | `null` | Means "no value" |
| **Array** | `[1, 2, 3]` | Ordered list of any values |
| **Object** | `{ x: 1, y: 2 }` | Key/value pairs |
| **Function** | `func(x) { return x; }` | First-class |

Check the type of a value at runtime:

```koda
print(type(42));       // number
print(type("hi"));     // string
print(type(true));     // bool
print(type(null));     // null
print(type([1, 2]));   // array
print(type({ a: 1 })); // object
print(type(print));    // function
```

---

## 4. Variables

Declare with `let`. Variables must be declared before use.

```koda
let name = "Koda";
let score = 0;
let active = true;
let nothing = null;
let uninit;           // starts as null
```

**Reassign** at any time:

```koda
score = score + 10;
score += 10;   // same thing
score++;       // increment by 1
score--;       // decrement by 1
```

**Destructuring** — pull fields out of an object into individual variables:

```koda
let player = { x: 10, y: 20, name: "Ada" };
let { x, y } = player;   // x = 10, y = 20
```

> `var` is **reserved** — always use `let`.

---

## 5. Operators

### Arithmetic

```koda
let a = 10 + 3;   // 13
let b = 10 - 3;   // 7
let c = 10 * 3;   // 30
let d = 10 / 3;   // 3.333...
let e = 10 % 3;   // 1  (remainder)
let f = 2 ** 8;   // 256 (power)
```

### Comparison

```koda
1 < 2     // true
1 <= 1    // true
2 > 1     // true
1 >= 1    // true
1 == 1    // true   (loose equal)
1 != 2    // true
1 === 1   // true   (strict equal)
1 !== 2   // true
```

### Logic

```koda
true && false   // false  (and)
true || false   // true   (or)
!true           // false  (not)
```

### Bitwise

```koda
5 & 3    // 1   (AND)
5 | 3    // 7   (OR)
5 ^ 3    // 6   (XOR)
~5       // -6  (NOT)
1 << 2   // 4   (left shift)
8 >> 1   // 4   (right shift)
8 >>> 1  // 4   (unsigned right shift)
```

### Compound assignment

```koda
x += 5;   x -= 5;   x *= 5;   x /= 5;   x %= 5;
x &= 3;   x |= 3;   x ^= 3;   x <<= 1;  x >>= 1;
```

### Nullish coalescing — `??`

Returns the right side **only** when the left side is `null`. (It does **not** trigger for `0` or `""`.)

```koda
let saved = null;
let score = saved ?? 0;   // score = 0

let name = "";
let display = name ?? "Guest";   // display = "" (NOT "Guest" — "" is not null)
```

### Nullish assign — `??=`

Assigns only when the variable is currently `null`:

```koda
let x = null;
x ??= 42;   // x is now 42

let y = 10;
y ??= 99;   // y is still 10
```

### Optional chaining — `?.`

Access a property without crashing when the receiver is `null`:

```koda
let player = null;
let hp = player?.health;   // null, no crash

let obj = { pos: { x: 5 } };
let x = obj?.pos?.x;   // 5
```

### `typeof`

```koda
let t = typeof 42;      // "number"
let s = typeof "hi";    // "string"
```

### Range — `..`

Used in `for`-`of` loops to generate integer sequences:

```koda
for (let i of 0..5) {
    print(i);   // 0, 1, 2, 3, 4
}
```

---

## 6. Control flow

> All `if`, `else`, `while`, `for`, and `switch` bodies **must** use `{ }` braces.

### `if` / `else if` / `else`

```koda
let score = 85;

if (score >= 90) {
    print("A");
} else if (score >= 70) {
    print("B");
} else {
    print("C");
}
```

**`if` as an expression** — assign the result directly:

```koda
let label = if (score > 50) { "pass" } else { "fail" };
```

### `while`

```koda
let n = 0;
while (n < 5) {
    print(n);
    n += 1;
}
```

### `do`…`while`

Runs the body **at least once**, then checks the condition:

```koda
let n = 0;
do {
    print(n);
    n += 1;
} while (n < 3);
```

### `for` (C-style)

```koda
for (let i = 0; i < 5; i += 1) {
    print(i);
}

// Multiple variables in the init
for (let i = 0, let j = 10; i < j; i += 1) {
    print(i, j);
}

// Infinite loop
for (;;) {
    break;
}
```

### `for`…`in` — iterate keys

Arrays yield numeric indices; objects yield string keys in insertion order:

```koda
let scores = [10, 20, 30];
for (let i in scores) {
    print(i);   // 0, 1, 2
}

let obj = { a: 1, b: 2 };
for (let key in obj) {
    print(key);   // "a", "b"
}
```

### `for`…`of` — iterate values

```koda
let colors = ["red", "green", "blue"];
for (let c of colors) {
    print(c);
}
```

### `for`…`of` with range

```koda
for (let i of 0..5) {
    print(i);   // 0, 1, 2, 3, 4
}
```

### `for`…`of` with pairs `[k, v]`

Gets both the key **and** value at once:

```koda
let obj = { x: 10, y: 20 };
for (let [key, val] of obj) {
    print(key, val);   // "x" 10, then "y" 20
}

let arr = ["a", "b"];
for (let [idx, item] of arr) {
    print(idx, item);  // 0 "a", then 1 "b"
}
```

### `switch`

```koda
let direction = "up";

switch (direction) {
    case "up":
        print("moving up");
        break;
    case "down":
        print("moving down");
        break;
    default:
        print("standing still");
}
```

Cases **fall through** unless you use `break`. Use `break` at the end of each branch you want to stop at.

**`switch` as an expression** — arms use `=>`:

```koda
let msg = switch (direction) {
    case "up"   => "going up"
    case "down" => "going down"
    default     => "stopped"
};
```

### `break` and `continue`

```koda
for (let i of 0..10) {
    if (i == 3) { continue; }   // skip 3
    if (i == 7) { break; }      // stop at 7
    print(i);
}
```

### `defer`

Runs a call when the **enclosing function exits**. Multiple defers run in **last-in, first-out** order:

```koda
func work() {
    defer print("done last");
    defer print("done second");
    defer print("done first");
    print("working...");
}
// prints: working... → done first → done second → done last
```

Useful for cleanup:

```koda
func loadLevel(path) {
    let file = openFile(path);
    defer closeFile(file);
    // ... rest of function; closeFile always runs
}
```

### `delete`

Removes an own property from an object:

```koda
let config = { volume: 80, muted: false };
delete config.muted;
print(len(config));   // 1
```

---

## 7. Functions

### Declaring a function

```koda
func greet(name) {
    print("Hello, " + name + "!");
}

greet("world");   // Hello, world!
```

### Returning a value

```koda
func add(a, b) {
    return a + b;
}

let result = add(3, 4);   // 7
```

`return;` with no value returns `null`.

### Default parameters

```koda
func greet(name = "world") {
    print("Hello, " + name);
}

greet();          // Hello, world
greet("Koda");    // Hello, Koda
```

### Rest parameters

`...name` collects all extra arguments into an array:

```koda
func sum(...numbers) {
    let total = 0;
    for (let n of numbers) {
        total += n;
    }
    return total;
}

print(sum(1, 2, 3, 4));   // 10
```

### Function values

Functions are values — assign them to variables, pass them around:

```koda
let double = func(x) {
    return x * 2;
};

print(double(5));   // 10
```

### `this` in object methods

When you call `obj.method()`, `this` inside the method is the object:

```koda
let player = {
    health: 100,
    heal(amount) {
        this.health = this.health + amount;
    }
};

player.heal(20);
print(player.health);   // 120
```

---

## 8. Closures

A function can capture variables from the scope where it was created:

```koda
func makeCounter() {
    let count = 0;
    return func() {
        count += 1;
        return count;
    };
}

let counter = makeCounter();
print(counter());   // 1
print(counter());   // 2
print(counter());   // 3
```

Each call to `makeCounter()` creates a separate `count` — they don't interfere with each other.

---

## 9. Objects

Objects are key/value stores. Keys are strings.

```koda
let pos = { x: 10, y: 20 };

// Read a field
print(pos.x);        // 10
print(pos["y"]);     // 20  (bracket access is the same)

// Write a field
pos.x = 99;
pos["y"] = 0;

// Add a new field
pos.z = 5;
```

### Computed keys

```koda
let key = "speed";
let obj = {};
obj[key] = 100;   // same as obj.speed = 100
```

### Methods (shorthand syntax)

```koda
let rect = {
    w: 100,
    h: 50,
    area() {
        return this.w * this.h;
    }
};

print(rect.area());   // 5000
```

### Dot notation on values

Property access uses `.field` (or `?.field` for optional chaining when the receiver may be `null`):

```koda
let math = import "@math";
print(math.pi);
print(math.sin(math.pi / 2));

let json = import "@json";
let text = json.stringify({ score: 10 });
let data = json.parse(text);
```

Namespace calls like `math.lerp(a, b, t)` and `json.parse(text)` work without an import when the name is a known stdlib namespace.

### `len` on objects

```koda
let config = { a: 1, b: 2, c: 3 };
print(len(config));   // 3  (number of keys)
```

---

## 10. Arrays

```koda
let items = [10, 20, 30];

print(items[0]);    // 10
print(items[2]);    // 30
print(len(items));  // 3

items[1] = 99;      // set a value
```

### Building arrays

```koda
let list = [];
list.push("apple");
list.push("banana");
list.push("cherry");
print(len(list));   // 3
print(list.pop());  // "cherry"
```

### Spread — `...`

Splice one array into another:

```koda
let a = [1, 2, 3];
let b = [0, ...a, 4];   // [0, 1, 2, 3, 4]
```

### Iterating arrays

```koda
let scores = [10, 20, 30];

// by value
for (let s of scores) {
    print(s);
}

// by index
for (let i in scores) {
    print(i, scores[i]);
}

// by index+value pair
for (let [i, s] of scores) {
    print(i, s);
}
```

---

## 11. Struct types

Structs give you named, ordered fields. Construct them with `TypeName { field: value, … }`.

```koda
struct Point {
    x,
    y
}

let p = Point { x: 3, y: 4 };
print(p.x);   // 3
print(p.y);   // 4

p.y = 10;
print(p.y);   // 10
```

A more complete example:

```koda
struct Rect {
    x,
    y,
    w,
    h
}

func area(r) {
    return r.w * r.h;
}

let box = Rect { x: 0, y: 0, w: 100, h: 50 };
print(area(box));   // 5000
```

See `tests/struct_test.koda` for more examples.

---

## 12. Enum types

Enums declare a set of named constants. Members are numbered `0`, `1`, `2`, … in the order you list them. Access them as `EnumName.Member`.

```koda
enum Dir {
    Up,
    Down,
    Left,
    Right
}

let d = Dir.Up;     // 0
let r = Dir.Right;  // 3

print(d == Dir.Up);    // true
print(r == 3);         // true
```

Use enums in `switch`:

```koda
enum State {
    Playing,
    Paused,
    GameOver
}

let current = State.Playing;

switch (current) {
    case State.Playing:
        print("game is running");
        break;
    case State.Paused:
        print("paused");
        break;
    case State.GameOver:
        print("game over");
        break;
}
```

See `tests/enum_test.koda` for more examples.

---

## 13. String methods

Call these on any string value. Names are case-insensitive.

| Method | What it does | Example |
|--------|-------------|---------|
| `split(delim)` | Split into an array. Empty `""` delimiter splits every character. | `"a,b,c".split(",")` → `["a","b","c"]` |
| `trim()` | Remove leading/trailing whitespace | `"  hi  ".trim()` → `"hi"` |
| `toUpper()` | Uppercase | `"hello".toUpper()` → `"HELLO"` |
| `toLower()` | Lowercase | `"HELLO".toLower()` → `"hello"` |
| `replace(a, b)` | Replace first occurrence of `a` with `b` | `"aabbaa".replace("b", "x")` → `"aaxbaa"` |
| `replaceAll(a, b)` | Replace every occurrence | `"aabbaa".replaceAll("a", "x")` → `"xxbbxx"` |
| `indexOf(needle)` | First position of needle, or `-1` | `"hello".indexOf("ll")` → `2` |
| `includes(needle)` | `true` if needle is found | `"hello".includes("ell")` → `true` |
| `slice(start, end)` | Substring from `start` up to (not including) `end` | `"hello".slice(1, 3)` → `"el"` |
| `startsWith(prefix)` | `true` if string starts with prefix | `"hello".startsWith("he")` → `true` |
| `endsWith(suffix)` | `true` if string ends with suffix | `"hello".endsWith("lo")` → `true` |

```koda
let s = "  Hello, World!  ";
print(s.trim());                         // "Hello, World!"
print(s.trim().toLower());               // "hello, world!"
print(s.trim().includes("World"));       // true
print(s.trim().replace("World", "Koda")); // "Hello, Koda!"

let parts = "one,two,three".split(",");
for (let p of parts) {
    print(p);   // one, two, three
}
```

---

## 14. Array methods

| Method | What it does |
|--------|-------------|
| `push(x)` | Add `x` to the end |
| `pop()` | Remove and return the last element (or `null` if empty) |
| `length()` | Number of elements (same as `len(arr)`) |
| `concat(...values)` | Return a new array with `values` appended |
| `join(sep)` | Join elements into a string with `sep` between them |
| `slice(start, end)` | Copy a range of elements (half-open) |
| `sort()` | Sort in place, returns the array |
| `reverse()` | Reverse in place, returns the array |
| `indexOf(item)` | First position of `item`, or `-1` |
| `includes(item)` | `true` if `item` is in the array |
| `map(callback)` | Return new array with `callback(element)` applied to each item |
| `filter(callback)` | Return new array with only elements where `callback(element)` is truthy |
| `find(callback)` | Return first element where `callback(element)` is truthy, or `null` |
| `reduce(callback)` or `reduce(callback, initial)` | Reduce to a single value; `callback(accumulator, element)` |

```koda
let nums = [3, 1, 4, 1, 5, 9];

// map — double every number
let doubled = nums.map(func(n) { return n * 2; });
print(doubled);   // [6, 2, 8, 2, 10, 18]

// filter — keep only numbers > 3
let big = nums.filter(func(n) { return n > 3; });
print(big);   // [4, 5, 9]

// find — first number > 4
let first = nums.find(func(n) { return n > 4; });
print(first);   // 5

// reduce — sum all
let total = nums.reduce(func(acc, n) { return acc + n; }, 0);
print(total);   // 23

// sort and join
print(nums.sort().join(", "));   // 1, 1, 3, 4, 5, 9
```

---

## 15. The math object

The compiler provides a built-in `math` object. You can use either `math.xxx(...)` or call the function directly as a global.

```koda
// Both of these are identical:
let a = math.floor(3.9);
let b = floor(3.9);
```

| Function | What it does |
|----------|-------------|
| `floor(x)` | Round down |
| `ceil(x)` | Round up |
| `round(x)` | Round to nearest |
| `trunc(x)` | Remove fractional part |
| `abs(x)` | Absolute value |
| `sqrt(x)` | Square root |
| `pow(base, exp)` | Power |
| `sign(x)` | `-1`, `0`, or `1` |
| `min(a, b)` | Smaller of two values |
| `max(a, b)` | Larger of two values |
| `clamp(val, lo, hi)` | Keep `val` between `lo` and `hi` |
| `lerp(a, b, t)` | Linear interpolation from `a` to `b` by `t` (0..1) |
| `sin(x)`, `cos(x)`, `tan(x)` | Trig (radians) |
| `asin(x)`, `acos(x)`, `atan(x)`, `atan2(y, x)` | Inverse trig |
| `log(x)`, `log10(x)`, `exp(x)` | Logarithm / exponent |
| `distance(x1, y1, x2, y2)` | 2D Euclidean distance |
| `distanceSq(x1, y1, x2, y2)` | Distance squared (faster, no sqrt) |
| `hypot(a, b)` | `sqrt(a*a + b*b)` |
| `degrees(r)` | Radians → degrees |
| `radians(d)` | Degrees → radians |
| `fmod(x, y)` | Floating-point remainder |
| `wrap(val, lo, hi)` | Wrap `val` around the range `[lo, hi)` |
| `smoothstep(a, b, t)` | Smooth interpolation |
| `smoothdamp(cur, target, ...)` | Spring-like smoothing |
| `approach(cur, target, step)` | Move `cur` toward `target` by at most `step` |
| `angleBetween(x1, y1, x2, y2)` | Angle between two 2D points |
| `normalize(x, y)` | Normalize a 2D vector |
| `pi` | π ≈ 3.14159… |
| `e` | Euler's number ≈ 2.71828… |

```koda
let angle = math.atan2(1, 1);        // ~0.785 radians (45°)
let deg   = math.degrees(angle);     // 45
let pos   = math.lerp(0, 100, 0.25); // 25
let safe  = math.clamp(150, 0, 100); // 100
```

---

## 16. Built-in functions

These are globally available (no import needed).

### Output and debugging

```koda
print("score:", score);          // prints to stdout, space-separated, newline at end
trace("debug value:", x);        // like print but for debug output
let s = format("x={} y={}", x, y);  // format a string without printing
```

### Values and types

```koda
let n = len([1, 2, 3]);   // 3 — works on arrays, strings, objects
let t = type(42);         // "number"
let t2 = typeof "hi";     // "string" — same as type()

// Type checks
print(isNumber(3));    // true
print(isString("x"));  // true
print(isBool(true));   // true
print(isNull(null));   // true
print(isArray([]));    // true
print(isObject({}));   // true
print(isFunction(print)); // true

// Coercions
let n = number("3.14");  // 3.14
let s = string(42);      // "42"
let b = bool(1);         // true
```

### Assertions and errors

```koda
assert(score >= 0, "score must be non-negative");
panic("something went very wrong");   // prints message and exits
```

### Files

```koda
let text = readFile("data.txt");
writeFile("out.txt", "hello\n");
appendFile("log.txt", "another line\n");

if (fileExists("save.dat")) {
    let data = readFile("save.dat");
}

deleteFile("temp.tmp");
```

### JSON

```koda
let obj  = parseJSON("{\"x\": 1}");
let text = toJSON({ x: 1, y: 2 });
print(text);   // {"x":1,"y":2}
```

### Time

```koda
let t  = time();         // current wall time (seconds)
let c  = clock();        // CPU clock
let ts = timestamp();    // Unix timestamp
let pt = programTime();  // time since program started (seconds)
sleep(100);              // sleep ~100 milliseconds
```

### Random numbers

```koda
randomSeed(42);               // seed for reproducible results

let u = random();             // float in [0, 1)
let i = randomInt(10);        // integer in [0, 10)
let j = randomInt(5, 15);     // integer in [5, 15)
```

### Results — `ok` / `err`

Convention for functions that may fail:

```koda
func divide(a, b) {
    if (b == 0) { return err("division by zero"); }
    return ok(a / b);
}
```

### Garbage collection

Koda manages memory automatically. For games you can fine-tune GC timing:

```koda
gcFrameStep(0.25);   // do a small GC step each frame (recommended in game loops)
gc();                // force a full collection now
gcDisable();         // pause GC (use carefully)
gcEnable();          // resume GC
let stats = gcStats(); // get GC statistics
```

### Substring check

```koda
if (matches(body, "error")) {
    print("found the word error");
}
```

`matches` is a **substring** check, not a regex.

---

## 17. Template strings

Use backticks. Embed any expression with `${ }`:

```koda
let name = "Koda";
let score = 42;

print(`Hello, ${name}!`);             // Hello, Koda!
print(`Score: ${score * 2}`);         // Score: 84
print(`Type: ${type(name)}`);         // Type: string
print(`Pi is about ${math.round(math.pi * 100) / 100}`);
```

---

## 18. Includes and imports

### `#include` — merge another file

The most common way to split code across files:

```koda
#include "lib/utils.koda"
#include "../stdlib/vec2.koda"
```

The included file's declarations become available as if they were written in the current file. Paths are relative to the file doing the including.

### `import()` — expression form

Load a module as an object. Exported top-level `let` and `func` names become properties:

```koda
let math = import "@math";
print(math.sqrt(16));

let utils = import("./lib/utils.koda");
utils.helper();
```

**Stdlib `@` modules** (shipped next to `koda`, or under `KODA_PATH`):

| Module | Properties (examples) |
|--------|------------------------|
| `@math` | `pi`, `sin`, `cos`, `sqrt`, `lerp`, `clamp`, … |
| `@json` | `parse`, `stringify`, `tryparse` |
| `@vec2` | `vec2`, `add`, `dot`, `normalize`, … |

The builtin `parseJSON(text)` returns `ok(value)` or `err(message)`; `json.parse` returns the value directly (or `null` on failure).

---

## 19. Native FFI hint

To bind a Koda name to a C symbol, use a special comment before the declaration:

```koda
// koda: extern myFunc my_c_function 2

let myFunc;
```

`2` is the argument count (arity). This wires `myFunc` to the C symbol `my_c_function` in the linked binary.

For wrapping C libraries in a friendlier way, use `koda wrap` — see **`docs/wrappers.md`**.

---

## 20. Truthy and falsy

In conditions (`if`, `while`, `&&`, `||`):

| Value | Truthy/Falsy |
|-------|-------------|
| `false` | **falsy** |
| `null` | **falsy** |
| `0` | **falsy** |
| `""` (empty string) | **falsy** |
| everything else | **truthy** |

> This differs from JavaScript — in Koda, **non-empty arrays** and **non-empty objects** are truthy.

```koda
if (0) { print("won't run"); }
if ("") { print("won't run"); }
if ([]) { print("WILL run — empty array is truthy"); }
if ({}) { print("WILL run — empty object is truthy"); }
```

---

## 21. Operator precedence

From highest (evaluated first) to lowest:

| Level | Operators |
|-------|-----------|
| Highest | `()` call, `.` member, `[]` index, `++` `--` postfix |
| | `++` `--` prefix, `!` `+` `-` unary, `typeof` |
| | `**` power |
| | `*` `/` `%` |
| | `+` `-` |
| | `<<` `>>` `>>>` |
| | `<` `<=` `>` `>=` |
| | `==` `!=` `===` `!==` |
| | `&` |
| | `^` |
| | `\|` |
| | `&&` |
| | `\|\|` `??` |
| Lowest | `=` `+=` `-=` etc. (assignment) |

When in doubt, use parentheses.

---

## 22. All keywords

All keywords are **case-insensitive**.

```
break    case     continue  default   defer
delete   do       else      enum      false
for      func     if        import    in
let      null     of        return    struct
switch   this     true      typeof    while
```

Directive: `#include "path.koda"`

> `var` is **reserved** — you will get an error if you use it. Use `let` instead.

---

## 23. All built-in names

These are all registered in `internal/codegen/builtin_register.go`. All case-insensitive at call sites.

**Control / errors:**
`ok`, `err`, `panic`, `assert`

**Output:**
`print`, `trace`, `format`

**Types and values:**
`len`, `type`, `typeof`, `number`, `string`, `bool`, `isNumber`, `isString`, `isBool`, `isNull`, `isArray`, `isObject`, `isFunction`

**Files:**
`readFile`, `writeFile`, `appendFile`, `fileExists`, `deleteFile`

**JSON:**
`parseJSON`, `toJSON`

**Time:**
`time`, `clock`, `timestamp`, `programTime`, `sleep`, `deltaTime`

**Random:**
`random`, `randomInt`, `randomChoice`, `randomSeed`

**Math globals:**
`abs`, `sqrt`, `pow`, `exp`, `log`, `log10`,
`floor`, `ceil`, `round`, `trunc`, `sign`,
`min`, `max`, `clamp`, `lerp`, `smoothstep`, `smoothdamp`, `approach`,
`sin`, `cos`, `tan`, `asin`, `acos`, `atan`, `atan2`,
`distance`, `distanceSq`, `hypot`, `normalize`, `angleBetween`,
`degrees`, `radians`, `fmod`, `wrap`,
`pi`, `e`

**Substring:**
`matches`

**GC:**
`gc`, `gcCollect`, `gcDisable`, `gcEnable`, `gcFrameStep`, `gcStats`

---

## See also

| File | What it's for |
|------|---------------|
| `docs/using-the-language.md` | Beginner walkthrough |
| `docs/commands.md` | `koda run`, `koda build`, and all CLI commands |
| `docs/wrappers.md` | Wrapping C/C++ libraries |
| `tests/*.koda` | Working examples for every feature |
