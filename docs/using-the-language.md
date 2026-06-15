# Getting Started with Koda

Koda is a language for making games and applications. It looks a lot like JavaScript, but it compiles straight to a fast native binary — no VM, no interpreter. This guide walks you from your very first program all the way to writing a game loop.

Keep **[../language.md](../language.md)** open beside this page — it has every feature, every builtin, and every operator with working examples.  
For CLI commands (`koda run`, `koda build`, …) see **[commands.md](commands.md)**.

---

## 1. Running your code

You write `.koda` files and use the `koda` command:

| What you want | Command |
|---------------|---------|
| Run a file | `koda run main.koda` |
| Build a binary to keep | `koda build main.koda -o mygame.exe` |
| Check for mistakes (no build) | `koda check main.koda` |
| Auto-format your code | `koda fmt main.koda` |
| Rebuild whenever you save | `koda watch main.koda` |
| Package everything to ship | `koda bundle main.koda -o dist/MyGame` |

---

## 2. Your first program

Create `hello.koda`:

```koda
print("Hello, Koda!");
```

Run it:

```
koda run hello.koda
```

Output: `Hello, Koda!`

You can also wrap everything in `func main()`:

```koda
func main() {
    print("Hello from main!");
}
```

Both styles work the same way.

---

## 3. Comments

```koda
// This is a line comment.

/*
   This is a block comment.
   It can span multiple lines.
*/
```

---

## 4. Variables

Use `let` for values that change and `const` for values that do not. Every statement ends with `;`.

```koda
let name = "Ada";
let score = 0;
const gravity = 900;
const screenWidth = 800;
let alive = true;
let nothing = null;   // null means "no value"
let later;            // also null until you assign it
```

Change a `let` binding any time (`const` cannot be reassigned):

```koda
score += 10;   // preferred
score++;       // add 1
score--;       // subtract 1
score = score + 10;   // still valid
```

Compound operators work on fields and indices too: `player.x += speed * dt;`, `arr[i] += 1;`, `obj.count++`.

> **Note:** `var` is not allowed in Koda. Always use `let` or `const`.

Optional type annotations (add when you want clarity):

```koda
let lives: int = 3;
let speed: float = 8.0;
let label: string = "Player";
```

Keywords and names are **case-insensitive** — `Let`, `LET`, and `let` all work.

---

## 5. Types

| Type | Example |
|------|---------|
| `int` | `3`, `let n: int = 0` |
| `float` | `3.14`, `let speed: float = 8.0` |
| `bool` | `true`, `false` |
| `string` | `"hello"`, `"Score: {score}"` |
| `array` | `[1, 2, 3]` |
| `map` / object | `{ x: 1, y: 2 }` |
| `func` | `func(x) { return x; }` |
| `null` | `null` |

Types are optional — beginners can omit annotations.

You can ask Koda what type a value is:

```koda
print(type(42));       // number
print(type("hello"));  // string
print(type(true));     // bool
print(type(null));     // null
```

---

## 6. Operators

```koda
// Arithmetic
let a = 10 + 3;   // 13
let b = 10 - 3;   // 7
let c = 10 * 3;   // 30
let d = 10 / 4;   // 2.5
let e = 10 % 3;   // 1  (remainder)
let f = 2 ** 8;   // 256 (power)

// Comparison — result is true or false
1 < 2      // true
1 == 1     // true
1 != 2     // true
2 >= 2     // true

// Logic
true && false   // false
true || false   // true
!true           // false
```

**Nullish coalescing** — use a default when something is null:

```koda
let saved = null;
let score = saved ?? 0;   // 0, because saved is null
```

**Optional chaining** — safely read a property that might not exist:

```koda
let player = null;
let hp = player?.health;   // null, no crash
```

---

## 7. Strings

```koda
let s = "Hello, World!";

print(s.toUpper());               // HELLO, WORLD!
print(s.toLower());               // hello, world!
print(s.includes("World"));       // true
print(s.replace("World", "Koda")); // Hello, Koda!
print(s.slice(0, 5));             // Hello

let words = "one two three".split(" ");
print(words[0]);   // one
print(len(words)); // 3
```

**String interpolation** — embed values in double-quoted strings:

```koda
let name = "player";
let lives = 3;
print("Score: {score}   Lives: {lives}");
```

**Template strings** — backticks with `${}`:

```koda
print(`${name} has ${lives} lives left`);
// player has 3 lives left
```

---

## 8. Arrays

An array holds a list of values.

```koda
let colors = ["red", "green", "blue"];

print(colors[0]);    // red
print(len(colors));  // 3

colors[1] = "yellow";   // change a value
colors.push("purple");  // add to end
colors.pop();           // remove from end
```

Loop through an array:

```koda
for (let color of colors) {
    print(color);
}
```

Useful array methods:

```koda
let nums = [3, 1, 4, 1, 5];

let big = nums.filter(func(n) { return n > 3; });
print(big);   // [4, 5]

let doubled = nums.map(func(n) { return n * 2; });
print(doubled);   // [6, 2, 8, 2, 10]

print(nums.sort().join(", "));   // 1, 1, 3, 4, 5
```

---

## 9. Objects

An object holds named values (like a record or dictionary).

```koda
let player = {
    name: "Ada",
    health: 100,
    x: 0,
    y: 0
};

print(player.name);    // Ada
print(player.health);  // 100

player.health -= 10;   // change a value
player.speed = 5;      // add a new field
```

**Methods** — functions inside an object can use `this` to refer to the object:

```koda
let player = {
    x: 0,
    speed: 5,
    moveRight() {
        this.x += this.speed;
    }
};

player.moveRight();
print(player.x);   // 5
```

---

## 10. Control flow

All `if`, `while`, `for`, and `switch` bodies **must use `{ }` braces**.

### `if` / `else`

```koda
let score = 75;

if (score >= 90) {
    print("A grade");
} else if (score >= 70) {
    print("B grade");
} else {
    print("C grade");
}
```

### `while` loop

```koda
let i = 0;
while (i < 5) {
    print(i);
    i += 1;
}
// prints 0, 1, 2, 3, 4
```

### `for` loop

```koda
// Count up
for (let i = 0; i < 5; i += 1) {
    print(i);
}

// Loop over a range (half-open: 0, 1, 2, 3, 4)
for i in 0..5 {
    print(i);
}

// Loop over array values
let fruits = ["apple", "banana", "cherry"];
for fruit in fruits {
    print(fruit);
}
```

### `switch`

```koda
let day = "Monday";

switch (day) {
    case "Monday":
        print("Start of the week");
        break;
    case "Friday":
        print("End of the week");
        break;
    default:
        print("Midweek");
}
```

Use `break` to stop at the end of each case. Without `break`, it falls through to the next case.

### `match`

Brace-style dispatch for enums and game states — each arm is its own block, **no fall-through**:

```koda
enum State {
    Start,
    Playing,
    GameOver
}

let current = State.Playing;

match current {
    State.Start {
        print("Press SPACE to begin");
    }
    State.Playing {
        update_game(dt);
    }
    State.GameOver {
        draw.text("GAME OVER", 400, 340, 36, colors.red);
    }
}
```

### `break` and `continue`

```koda
for i in 0..10 {
    if (i == 3) { continue; }   // skip 3
    if (i == 7) { break; }      // stop at 7
    print(i);
}
```

---

## 11. Functions

```koda
func greet(name) {
    print("Hello, " + name + "!");
}

greet("World");   // Hello, World!
```

**Return a value:**

```koda
func add(a, b) {
    return a + b;
}

let result = add(3, 4);
print(result);   // 7
```

**Default parameters:**

```koda
func greet(name = "stranger") {
    print("Hello, " + name);
}

greet();         // Hello, stranger
greet("Ada");    // Hello, Ada
```

**Rest parameters** — collect extra arguments into an array:

```koda
func sum(...numbers) {
    let total = 0;
    for (let n of numbers) { total += n; }
    return total;
}

print(sum(1, 2, 3, 4));   // 10
```

**Closures** — a function that remembers values from where it was created:

```koda
func makeCounter() {
    let n = 0;
    return func() {
        n += 1;
        return n;
    };
}

let count = makeCounter();
print(count());   // 1
print(count());   // 2
print(count());   // 3
```

---

## 12. Struct types

A struct is a named type with fixed fields. Easier to read than plain objects.

```koda
struct Point {
    x,
    y
}

let p = Point { x: 10, y: 20 };
print(p.x);   // 10

p.y = 99;
```

**Case-insensitive names:** `struct Player` and `let player` are the same binding — use `struct Mario` with `let player = Mario { ... }` instead.

```koda
struct Bullet {
    x,
    y,
    vx,
    vy
}

func moveBullet(b) {
    b.x += b.vx;
    b.y += b.vy;
}

let b = Bullet { x: 0, y: 0, vx: 5, vy: -3 };
moveBullet(b);
print(b.x, b.y);   // 5 -3
```

---

## 13. Enum types

An enum gives names to a set of numbers. Great for game states, directions, etc.

```koda
enum State {
    Start,    // 0
    Playing,  // 1
    Paused,   // 2
    GameOver  // 3
}

let current = State.Playing;

if (current == State.Playing) {
    print("Game is running!");
}
```

Use them in `switch` or `match`:

```koda
switch (current) {
    case State.Start:
        print("Press SPACE to begin");
        break;
    case State.Playing:
        print("Playing...");
        break;
    case State.GameOver:
        print("Game Over!");
        break;
}

match current {
    State.Playing {
        update_game(dt);
    }
    State.GameOver {
        draw.text("GAME OVER", 400, 340, 36, colors.red);
    }
}
```

---

## 14. Math

The `math` object (or global functions) gives you everything you need:

```koda
print(math.floor(3.9));          // 3
print(math.ceil(3.1));           // 4
print(math.round(3.5));          // 4
print(math.abs(-5));             // 5
print(math.sqrt(16));            // 4
print(math.min(3, 7));           // 3
print(math.max(3, 7));           // 7
print(math.clamp(150, 0, 100));  // 100

// Lerp: smooth value between two points
let pos = math.lerp(0, 100, 0.25);   // 25

// Distance between two points
let d = math.distance(0, 0, 3, 4);   // 5

// Trig (angles in radians)
let angle = math.atan2(1, 0);   // π/2
print(math.degrees(angle));     // 90
```

---

## 15. Useful builtins

```koda
// Output
print("hello", 42, true);         // space-separated
let s = format("x={} y={}", x, y); // build a string

// Assertions (crash with message if false)
assert(lives > 0, "player is dead");

// File I/O
let text = readFile("level.txt");
writeFile("save.dat", data);
if (fileExists("save.dat")) { /* ... */ }

// Random numbers
randomSeed(42);
let roll = randomInt(1, 7);    // 1-6 inclusive
let chance = random();         // 0.0 to 1.0

// JSON
let obj = parseJSON("{\"hp\": 100}");
let str = toJSON({ hp: 100, name: "Ada" });
```

---

## 16. Game loop pattern

Here is a minimal game loop using Raylib (Koda's built-in graphics library):

```koda
#include "../wrappers/raylib_shim/raylib.koda"

let screenW = 800;
let screenH = 600;
let x = 400;
let y = 300;
let speed = 4;

InitWindow(screenW, screenH, "My Game");
SetTargetFPS(60);

while (!WindowShouldClose()) {
    // Update
    if (IsKeyDown(263)) { x -= speed; }   // left arrow
    if (IsKeyDown(262)) { x += speed; }   // right arrow
    if (IsKeyDown(265)) { y -= speed; }   // up arrow
    if (IsKeyDown(264)) { y += speed; }   // down arrow

    // Draw
    BeginDrawing();
    ClearBackground(colors.dark);
    DrawCircle(x, y, 20, colors.red);
    EndDrawing();
}

CloseWindow();
```

Run it with:

```
koda run mygame.koda
```

---

## 17. Splitting code across files

Use `#include` to pull in another `.koda` file:

```koda
// main.koda
#include "player.koda"
#include "enemies.koda"

func main() {
    // player.koda and enemies.koda code is available here
}
```

The included file is merged in as if you typed it directly. Paths are relative to the file doing the including.

---

## 18. `defer` — cleanup made easy

`defer` schedules a call to run when the function exits. Multiple defers run in **reverse order** (last in, first out):

```koda
func loadLevel(path) {
    let f = openFile(path);
    defer closeFile(f);   // always runs, even if we return early

    if (!f) { return; }   // closeFile still runs here

    // ... process file ...
}
```

---

## 19. Where to go next

| Document | What it covers |
|----------|----------------|
| **[../language.md](../language.md)** | Complete reference — every feature, every builtin, with examples |
| **[language.md](language.md)** | Compact syntax cheat sheet |
| **[commands.md](commands.md)** | All `koda` CLI commands |
| **[wrappers.md](wrappers.md)** | Using C/C++ libraries from Koda |
| **`tests/*.koda`** | Small runnable examples for every feature |
| **`examples/games/`** | Complete game examples including brick breaker |

When something behaves unexpectedly, run `koda check yourfile.koda` — it checks your code without compiling and gives clear error messages.
