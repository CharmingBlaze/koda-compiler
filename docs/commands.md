# Koda CLI — every command

Use the **`koda`** binary from [GitHub Releases](https://github.com/CharmingBlaze/koda-compiler/releases). For full inline help:

```bash
koda help
```

---

## `koda new`

Create a new project directory with `koda.json`, `src/main.koda`, `assets/`, and a template README.

```bash
koda new <name> [--template hello|game|graphics]
```

| Template | What you get |
|----------|----------------|
| **hello** (default) | Minimal print program — runs immediately |
| **game** | Text-based lunar lander — no native libraries |
| **graphics** | Bouncing-ball window via bundled Raylib shim — requires Raylib link flags |

**Examples**

```bash
koda new myapp
koda new lander --template game
koda new bounce --template graphics
cd bounce
# Windows: $env:KODA_LINKFLAGS = "-lraylib -lopengl32 -lgdi32 -lwinmm"
koda run
```

---

## `koda.json` (project manifest)

When `koda.json` sits in the project root (or a parent directory), you can omit the entry file on `run`, `build`, `watch`, `check`, and `bundle`.

```json
{
  "name": "mygame",
  "version": "0.1.0",
  "entry": "src/main.koda",
  "bundle": {
    "assets": ["assets"],
    "extra": ["LICENSE"]
  },
  "native": {
    "sources": ["wrappers/mylib/wrapper.c"],
    "linkflags": "-lmylib"
  }
}
```

| Field | Purpose |
|-------|---------|
| `name` | Project name (default executable / bundle name) |
| `entry` | Main `.koda` file (required) |
| `bundle.assets` | Directories copied into `koda bundle` output |
| `bundle.extra` | Additional files or folders to copy |
| `native.sources` | C/C++ files linked on every build (sets `KODA_NATIVE_SOURCES` when unset) |
| `native.linkflags` | Linker flags (sets `KODA_LINKFLAGS` when unset) |

Environment variables always override manifest values.

**Project-aware commands**

```bash
cd mygame
koda run
koda build -o mygame
koda check
koda bundle -o dist/mygame
```

---

## `koda run` / `koda native`

Compile the entry **`.koda`** file to a native executable (temporary), run it, then delete the temp binary.

```bash
koda run [--no-opt] [<file.koda>]
koda native [--no-opt] [<file.koda>]   # same as run
```

- **`--no-opt`** — skips LLVM IR optimisation and uses a less aggressive native compile. Use if a very large program hits a flaky Clang optimiser on your machine.
- Omit the file when `koda.json` defines `entry`.

**Example**

```bash
koda run src/main.koda
koda run              # uses koda.json entry
koda run --no-opt src/main.koda
```

---

## `koda watch`

Rebuild and rerun whenever **`.koda`** files under the entry file’s directory change.

```bash
koda watch [--no-opt] [<file.koda>]
```

**Example**

```bash
koda watch src/main.koda
koda watch            # project entry
```

---

## `koda run`

Compile to a native binary and run (temporary executable).

```bash
koda run [--no-opt] [--debug] [<file.koda>] [-- <program args...>]
```

Arguments after `--` are passed to your program.

---

## `koda test`

```bash
koda test [--no-opt] [-v] [--failfast] [-run <pattern>] [<files...>]
```

---

## `koda lint`

Run `check` plus `fmt --check` on paths (default `./...`).

```bash
koda lint [paths...]
```

---

## `koda bench` / `profile`

Time repeated compile+run cycles.

```bash
koda bench [--count N] [--warmup N] [--no-opt] [--debug] <file.koda> [-- <args...>]
```

---

## `koda debug`

Run with debug symbols (`koda run --debug`).

---

## `koda eval` / `repl`

```bash
koda eval 'print(1 + 2)'
koda repl
```

---

## `koda init`

Alias for `koda new`.

---

## `koda env` / `completions` / `update` / `doc` / `lsp`

```bash
koda env [--export]
koda completions bash|zsh|fish
koda update [--check-only]
koda doc [stdlib | module @name]
koda lsp
```

See **[reference/cli.md](reference/cli.md)** for full details.

---

## `koda clean`

```bash
koda clean [<dir>] [--cache]
```

---

## `koda check`

Parse, resolve imports, and run semantic analysis only — **no** native compile.

```bash
koda check [<file.koda>]
koda check ./...
```

Prints **`OK`** if the program is valid.

**Example**

```bash
koda check src/main.koda
koda check            # project entry
```

---

## `koda fmt`

Format **`.koda`** sources with the canonical formatter.

```bash
koda fmt [--check] <file.koda> [more files...]
koda fmt [--check] ./...
```

- **`--check`** — do not write files; exit with an error if any file would change (for CI).

**Examples**

```bash
koda fmt src/main.koda
koda fmt ./...
koda fmt --check .
```

---

## `koda disasm`

Print **LLVM IR** for the program (after parse, sema, and codegen). Useful for debugging the compiler pipeline, not for everyday game dev.

```bash
koda disasm <file.koda>
```

---

## `koda build`

Produce a **native executable** you keep.

```bash
koda build [--no-opt] [--debug] <file.koda> [-o <exe>]
```

- **`--no-opt`** — same meaning as for `run`.
- **`--debug`** — emit debug symbols and favour debuggable builds (implies **`--no-opt`** for the native step).
- **`-o`** — output path; if omitted, a default name next to the source is used.

**Examples**

```bash
koda build game.koda -o dist/game.exe
koda build --debug game.koda -o dist/game_debug.exe
```

---

## `koda bundle`

Build the program and write a **folder** you can zip or ship (executable, helper scripts, README, etc.).

```bash
koda bundle <file.koda> [-o <dir>]
```

Default output directory is **`dist`** if you omit **`-o`**.

**Example**

```bash
koda bundle game.koda -o dist/MyGame
```

Copy extra files and directories (DLLs, assets folders) with **`KODA_BUNDLE_FILES`**. Use your OS path-list separator (`;` on Windows, `:` on Linux/macOS), or quote paths that contain spaces.

---

## `koda wrap`

Forward arguments to **`kodawrap`** (C/C++ header → **`.koda`** bindings + **`wrapper.c`**). **`koda`** looks for **`kodawrap`** next to itself, then **`wrapgen`**, **`kujiwrap`**, then **`PATH`**.

```bash
koda wrap --help
koda wrap -name mylib -headers ./include/mylib.h -out ./wrappers/mylib
```

Full workflow: **[wrappers.md](wrappers.md)**.

---

## `koda paths`

Print **machine-readable** toolchain paths (for scripts and CI).

```bash
koda paths
```

---

## `koda doctor`

Human-readable **health check**: clang resolution, runtime library, stdlib, install directory, etc.

```bash
koda doctor
```

---

## `koda version`

Print version and platform.

```bash
koda version
koda --version
```

---

## `koda help`

Print the full help screen (commands, environment variables, examples).

```bash
koda help
koda --help
```

---

## Environment variables (quick reference)

| Variable | Purpose |
|----------|---------|
| **`KODA_CLANG`** / **`CC`** | Clang driver for native builds (if not using embedded release toolchain) |
| **`KODA_PATH`** | Extra **`@module`** search directories |
| **`KODA_WRAPPERS`** | Pre-built wrapper **`.koda`** trees |
| **`KODA_NATIVE_SOURCES`** | C/C++ sources linked into your app (e.g. **`wrapper.c`**) |
| **`KODA_LINKFLAGS`** | Extra linker flags (**`-l`**, **`-L`**, frameworks, …) |
| **`KODA_BUNDLE_FILES`** | Extra files or directories copied into **`koda bundle`** output |

See **`koda help`** for the complete list and notes.
 Essential Commands for Game/App Development
These are critical quality-of-life features that should be in v1.0

Critical Missing Features Analysis
Looking at your Breakout game code, you're manually doing:
javascript// Current pain points:
let ft = rl.GetFrameTime()  // Delta time from raylib
this.x = this.x + this.vx * ft  // Manual delta time math
ball.vx = 220  // Magic numbers everywhere
if (ball.x < ball.radius) { ... }  // Manual bounds checking
These should be language-level features, not library calls.

Tier 1: Must-Have Built-ins (Add to v1.0)
1. Time & Delta Time
javascript// Built-in time functions
deltaTime()     // Seconds since last frame (automatic)
time()          // Seconds since program start
timestamp()     // Unix timestamp
sleep(ms)       // Sleep for milliseconds
setTimeout(fn, ms)    // Call function after delay
setInterval(fn, ms)   // Call function repeatedly

// Usage in game loop:
func update() {
    let dt = deltaTime();  // Automatic frame timing
    player.x += player.vx * dt;
    player.y += player.vy * dt;
}
Implementation:
go// internal/runtime/time.go

var (
    programStart time.Time
    lastFrame    time.Time
)

func init() {
    programStart = time.Now()
    lastFrame = programStart
}

func native_deltaTime(args []Value) Value {
    now := time.Now()
    dt := now.Sub(lastFrame).Seconds()
    lastFrame = now
    return NumberVal(dt)
}

func native_time(args []Value) Value {
    elapsed := time.Since(programStart).Seconds()
    return NumberVal(elapsed)
}

func native_timestamp(args []Value) Value {
    return NumberVal(float64(time.Now().Unix()))
}

func native_sleep(args []Value) Value {
    if len(args) != 1 || !IsNumber(args[0]) {
        return NullVal
    }
    ms := AsNumber(args[0])
    time.Sleep(time.Duration(ms) * time.Millisecond)
    return NullVal
}

2. Random Numbers (Essential for Games)
javascript// Random functions
random()              // [0, 1)
random(max)           // [0, max)
random(min, max)      // [min, max)
randomInt(max)        // Integer [0, max)
randomInt(min, max)   // Integer [min, max)
randomChoice(array)   // Random element from array
randomSeed(seed)      // Set seed for reproducibility

// Game usage:
let spawnX = random(0, screenWidth);
let spawnY = random(0, screenHeight);
let enemyType = randomChoice(["goblin", "orc", "dragon"]);
let damage = randomInt(10, 20);

// Procedural generation:
randomSeed(12345);  // Same seed = same level
let level = generateLevel();
Implementation:
go// internal/runtime/random.go

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

func native_random(args []Value) Value {
    switch len(args) {
    case 0:
        // random() → [0, 1)
        return NumberVal(rng.Float64())
    
    case 1:
        // random(max) → [0, max)
        if !IsNumber(args[0]) {
            return NullVal
        }
        max := AsNumber(args[0])
        return NumberVal(rng.Float64() * max)
    
    case 2:
        // random(min, max) → [min, max)
        if !IsNumber(args[0]) || !IsNumber(args[1]) {
            return NullVal
        }
        min := AsNumber(args[0])
        max := AsNumber(args[1])
        return NumberVal(min + rng.Float64()*(max-min))
    
    default:
        return NullVal
    }
}

func native_randomInt(args []Value) Value {
    switch len(args) {
    case 1:
        // randomInt(max) → [0, max)
        if !IsNumber(args[0]) {
            return NullVal
        }
        max := int(AsNumber(args[0]))
        return NumberVal(float64(rng.Intn(max)))
    
    case 2:
        // randomInt(min, max) → [min, max)
        if !IsNumber(args[0]) || !IsNumber(args[1]) {
            return NullVal
        }
        min := int(AsNumber(args[0]))
        max := int(AsNumber(args[1]))
        return NumberVal(float64(min + rng.Intn(max-min)))
    
    default:
        return NullVal
    }
}

func native_randomChoice(args []Value) Value {
    if len(args) != 1 || !IsArray(args[0]) {
        return NullVal
    }
    
    arr := AsArray(args[0])
    if arr.Count == 0 {
        return NullVal
    }
    
    idx := rng.Intn(arr.Count)
    return arr.Values[idx]
}

func native_randomSeed(args []Value) Value {
    if len(args) != 1 || !IsNumber(args[0]) {
        return NullVal
    }
    
    seed := int64(AsNumber(args[0]))
    rng = rand.New(rand.NewSource(seed))
    return NullVal
}

3. Math Utilities (Game Math)
javascript// Core math
PI          // 3.14159...
E           // 2.71828...

// Trigonometry
sin(x)
cos(x)
tan(x)
asin(x)
acos(x)
atan(x)
atan2(y, x)

// Exponential
pow(base, exp)
sqrt(x)
exp(x)
log(x)
log10(x)

// Rounding
floor(x)
ceil(x)
round(x)
trunc(x)

// Utilities
abs(x)
sign(x)        // -1, 0, or 1
min(a, b, ...)
max(a, b, ...)
clamp(val, min, max)
lerp(a, b, t)      // Linear interpolation
smoothstep(a, b, t) // Smooth interpolation
map(val, inMin, inMax, outMin, outMax)  // Remap range

// Distance & angles
distance(x1, y1, x2, y2)
distanceSq(x1, y1, x2, y2)  // Faster (no sqrt)
angleBetween(x1, y1, x2, y2)
normalize(x, y)  // Returns {x, y} normalized

// Game usage:
let angle = atan2(target.y - player.y, target.x - player.x);
player.vx = cos(angle) * speed;
player.vy = sin(angle) * speed;

let dist = distance(player.x, player.y, enemy.x, enemy.y);
if (dist < 50) {
    takeDamage();
}

let health = clamp(player.health, 0, 100);
let alpha = lerp(0, 255, health / 100);  // Fade based on health
Implementation:
go// internal/runtime/math.go

func native_lerp(args []Value) Value {
    if len(args) != 3 {
        return NullVal
    }
    a := AsNumber(args[0])
    b := AsNumber(args[1])
    t := AsNumber(args[2])
    return NumberVal(a + (b-a)*t)
}

func native_clamp(args []Value) Value {
    if len(args) != 3 {
        return NullVal
    }
    val := AsNumber(args[0])
    min := AsNumber(args[1])
    max := AsNumber(args[2])
    
    if val < min {
        return NumberVal(min)
    }
    if val > max {
        return NumberVal(max)
    }
    return NumberVal(val)
}

func native_distance(args []Value) Value {
    if len(args) != 4 {
        return NullVal
    }
    x1 := AsNumber(args[0])
    y1 := AsNumber(args[1])
    x2 := AsNumber(args[2])
    y2 := AsNumber(args[3])
    
    dx := x2 - x1
    dy := y2 - y1
    return NumberVal(math.Sqrt(dx*dx + dy*dy))
}

func native_angleBetween(args []Value) Value {
    if len(args) != 4 {
        return NullVal
    }
    x1 := AsNumber(args[0])
    y1 := AsNumber(args[1])
    x2 := AsNumber(args[2])
    y2 := AsNumber(args[3])
    
    return NumberVal(math.Atan2(y2-y1, x2-x1))
}

func native_map(args []Value) Value {
    if len(args) != 5 {
        return NullVal
    }
    val := AsNumber(args[0])
    inMin := AsNumber(args[1])
    inMax := AsNumber(args[2])
    outMin := AsNumber(args[3])
    outMax := AsNumber(args[4])
    
    normalized := (val - inMin) / (inMax - inMin)
    return NumberVal(outMin + normalized*(outMax-outMin))
}

4. Array Utilities (Critical)
javascript// Array methods (your code desperately needs these)
arr.map(fn)
arr.filter(fn)
arr.forEach(fn)
arr.find(fn)
arr.findIndex(fn)
arr.some(fn)
arr.every(fn)
arr.reduce(fn, initial)
arr.sort(fn)
arr.reverse()
arr.slice(start, end)
arr.concat(other)
arr.indexOf(item)
arr.includes(item)

// Your Breakout code would become:
let aliveBricks = this.bricks.filter(b => b.alive());
let allDead = this.bricks.every(b => !b.alive());

// Instead of:
let i = 0;
while (i < this.bricks.length()) {
    if (this.bricks[i].alive()) { ... }
    i = i + 1;
}

5. String Utilities
javascript// String methods
str.split(sep)
str.trim()
str.upper()
str.lower()
str.startsWith(prefix)
str.endsWith(suffix)
str.indexOf(substr)
str.slice(start, end)
str.replace(old, new)
str.replaceAll(old, new)

// String formatting
format("Player: {}, Score: {}", name, score)
// Or template literals (better):
`Player: ${name}, Score: ${score}`

Tier 2: Nice-to-Have (v1.1)
6. Wait/Async Utilities
javascript// Coroutine-style waiting
wait(seconds)           // Pause execution
waitUntil(condition)    // Wait for condition to be true
waitForFrame()          // Wait one frame

// Game usage:
func enemyAI() {
    while (true) {
        patrol();
        wait(2);  // Patrol for 2 seconds
        
        if (playerVisible()) {
            chase();
            wait(5);
        }
    }
}

func cutscene() {
    showText("Once upon a time...");
    wait(3);
    showText("A hero arose...");
    wait(3);
    startGame();
}
Note: This requires coroutine support or async/await, which is complex. Consider for v1.1+.

7. Input Helpers
javascript// Keyboard (if not using raylib)
isKeyDown(key)
isKeyPressed(key)
isKeyReleased(key)

// Mouse
getMouseX()
getMouseY()
getMousePosition()  // Returns {x, y}
isMouseButtonDown(button)
isMouseButtonPressed(button)

// Gamepad
getGamepadAxis(gamepad, axis)
isGamepadButtonDown(gamepad, button)

8. Debug Utilities
javascript// Debugging
assert(condition, message)
trace(...args)           // Debug print
profile(name, fn)        // Profile function execution
benchmark(fn, iterations)

// Usage:
assert(player.health >= 0, "Health cannot be negative!");

profile("update", func() {
    updatePhysics();
    updateAI();
});

Tier 3: Advanced (v2.0+)
9. Collections
javascript// Set
let enemies = Set();
enemies.add(goblin);
enemies.has(goblin);
enemies.delete(goblin);

// Map
let scores = Map();
scores.set("Alice", 100);
scores.get("Alice");

10. File I/O
javascript// File operations
readFile(path)
writeFile(path, content)
appendFile(path, content)
fileExists(path)
deleteFile(path)

// JSON
parseJSON(str)
toJSON(obj)

// Save game example:
let saveData = {
    level: currentLevel,
    score: score,
    inventory: items
};
writeFile("save.json", toJSON(saveData));

// Load game:
let data = parseJSON(readFile("save.json"));
currentLevel = data.level;

Updated Standard Library (v1.0)
Complete Built-ins List
javascript// ============================================================================
// TIME & TIMING
// ============================================================================
deltaTime()              // Seconds since last frame
time()                   // Seconds since program start
timestamp()              // Unix timestamp
sleep(ms)                // Sleep milliseconds

// ============================================================================
// RANDOM
// ============================================================================
random()                 // [0, 1)
random(max)              // [0, max)
random(min, max)         // [min, max)
randomInt(max)           // Integer [0, max)
randomInt(min, max)      // Integer [min, max)
randomChoice(array)      // Random element
randomSeed(seed)         // Set seed

// ============================================================================
// MATH CONSTANTS
// ============================================================================
PI                       // 3.14159265358979
E                        // 2.71828182845905

// ============================================================================
// MATH - TRIGONOMETRY
// ============================================================================
sin(x)
cos(x)
tan(x)
asin(x)
acos(x)
atan(x)
atan2(y, x)

// ============================================================================
// MATH - EXPONENTIAL
// ============================================================================
pow(base, exp)
sqrt(x)
exp(x)
log(x)
log10(x)

// ============================================================================
// MATH - ROUNDING
// ============================================================================
floor(x)
ceil(x)
round(x)
trunc(x)

// ============================================================================
// MATH - UTILITIES
// ============================================================================
abs(x)
sign(x)
min(...values)
max(...values)
clamp(val, min, max)
lerp(a, b, t)
smoothstep(a, b, t)
map(val, inMin, inMax, outMin, outMax)

// ============================================================================
// MATH - GEOMETRY
// ============================================================================
distance(x1, y1, x2, y2)
distanceSq(x1, y1, x2, y2)
angleBetween(x1, y1, x2, y2)
normalize(x, y)          // Returns {x, y}

// ============================================================================
// ARRAY METHODS
// ============================================================================
arr.length()
arr.push(item)
arr.pop()
arr.map(fn)
arr.filter(fn)
arr.forEach(fn)
arr.find(fn)
arr.findIndex(fn)
arr.some(fn)
arr.every(fn)
arr.reduce(fn, initial)
arr.sort(compareFn)
arr.reverse()
arr.slice(start, end)
arr.concat(other)
arr.indexOf(item)
arr.includes(item)

// ============================================================================
// STRING METHODS
// ============================================================================
str.length()
str.split(sep)
str.trim()
str.upper()
str.lower()
str.startsWith(prefix)
str.endsWith(suffix)
str.indexOf(substr)
str.slice(start, end)
str.replace(old, new)
str.replaceAll(old, new)

// ============================================================================
// TYPE CHECKING
// ============================================================================
type(value)              // Returns "number", "string", etc.
isNumber(value)
isString(value)
isBool(value)
isNull(value)
isArray(value)
isObject(value)
isFunction(value)

// ============================================================================
// CONVERSION
// ============================================================================
number(value)
string(value)
bool(value)

// ============================================================================
// I/O
// ============================================================================
print(...values)
format(template, ...values)

// ============================================================================
// DEBUG
// ============================================================================
assert(condition, message)
trace(...values)

Your Breakout Code - Before/After
Before (Current)
javascriptfunc update(ft) {
    this.paddle.update(this.rl, ft, this.sw, this.KEY_LEFT, this.KEY_RIGHT);
    this.ball.step(ft);
    this.ball.bounceWalls(this.sw, this.sh);
    this.ball.bouncePaddle(this.paddle);
    this.ball.resolveBricks(this.bricks);
    if (this.ball.pastBottom(this.sh)) {
        this.ball.reset(this.sw, this.sh);
    }
}

func resolveBricks(bricks) {
    let i = 0;
    while (i < bricks.length()) {
        let b = bricks[i];
        if (b.alive()) {
            if (this.hitsRect(b.x, b.y, b.w, b.h)) {
                b.hit();
                this.vy = -this.vy;
            }
        }
        i = i + 1;
    }
}
After (With New Built-ins)
javascriptfunc update() {
    let dt = deltaTime();  // Automatic!
    
    this.paddle.update(dt);
    this.ball.step(dt);
    this.ball.bounceWalls();
    this.ball.bouncePaddle(this.paddle);
    this.ball.resolveBricks(this.bricks);
    
    if (this.ball.pastBottom()) {
        this.ball.reset();
    }
    
    // Check win condition
    if (this.bricks.every(b => !b.alive())) {
        this.nextLevel();
    }
}

func resolveBricks(bricks) {
    // Much cleaner!
    bricks
        .filter(b => b.alive())
        .find(b => this.hitsRect(b.x, b.y, b.w, b.h))
        ?.hit();  // Optional chaining (v1.1)
    
    // Or explicit:
    let hit = bricks.find(b => 
        b.alive() && this.hitsRect(b.x, b.y, b.w, b.h)
    );
    
    if (hit) {
        hit.hit();
        this.vy = -this.vy;
    }
}

func spawnPowerup() {
    let x = random(20, this.sw - 20);
    let y = random(20, this.sh - 20);
    let type = randomChoice(["extra_ball", "laser", "expand_paddle"]);
    
    return createPowerup(x, y, type);
}

Implementation Priority
Week 1: Core Time/Random/Math
go// Add to internal/runtime/

// time.go
- deltaTime()
- time()
- timestamp()
- sleep()

// random.go
- random() variants
- randomInt() variants
- randomChoice()
- randomSeed()

// math.go
- PI, E constants
- All trig functions
- Rounding functions
- Utility functions (clamp, lerp, distance, etc.)
Tests:
javascript// tests/stdlib/time_test.koda
let start = time();
sleep(100);
let end = time();
assert(end - start >= 0.1, "Sleep failed");

// tests/stdlib/random_test.koda
randomSeed(12345);
let a = random();
randomSeed(12345);
let b = random();
assert(a == b, "Seed failed");

// tests/stdlib/math_test.koda
assert(abs(-5) == 5);
assert(clamp(15, 0, 10) == 10);
assert(distance(0, 0, 3, 4) == 5);

Week 2: Array/String Methods
go// Add to internal/runtime/

// array.go
- All array methods as methods on array objects

// string.go
- All string methods as methods on string objects

Week 3: Debug/Format
go// debug.go
- assert()
- trace()

// format.go
- format() function

Final Recommendation
Add to v1.0 (Essential):

✅ Time functions (deltaTime, time, sleep)
✅ Random functions (all variants)
✅ Math functions (complete set)
✅ Array methods (map, filter, etc.)
✅ String methods (split, trim, etc.)
✅ format() function

Add to v1.1 (Quality of Life):

🔵 Template literals
🔵 Arrow functions
🔵 Optional chaining
🔵 Coroutines/wait

Add to v2.0 (Advanced):

🔵 Async/await
🔵 File I/O
🔵 Collections (Set, Map)
🔵 Modules (import/export)


## Standard library wishlist (notes)

Possible directions for builtins and stdlib (for planning only):

## Standard Library - Built-in Functions

Implement these essential functions for game/app development:

### Time & Timing
- deltaTime() - Automatic frame delta
- time() - Program runtime
- timestamp() - Unix timestamp
- sleep(ms) - Pause execution

### Random
- random(), random(max), random(min, max)
- randomInt(max), randomInt(min, max)
- randomChoice(array)
- randomSeed(seed)

### Math
- Constants: PI, E
- Trig: sin, cos, tan, asin, acos, atan, atan2
- Exp: pow, sqrt, exp, log, log10
- Rounding: floor, ceil, round, trunc
- Utils: abs, sign, min, max, clamp, lerp, map
- Geometry: distance, distanceSq, angleBetween

### Arrays
Implement as methods: map, filter, forEach, find, some, every, reduce, etc.

### Strings
Implement as methods: split, trim, upper, lower, startsWith, endsWith, etc.

These are critical for game development. Implement in internal/runtime/ package
`r`n---`r`n`r`n## Legacy Notes (migrated from root "more commands.md")`r`n
 Essential Commands for Game/App Development
These are critical quality-of-life features that should be in v1.0

Critical Missing Features Analysis
Looking at your Breakout game code, you're manually doing:
javascript// Current pain points:
let ft = rl.GetFrameTime()  // Delta time from raylib
this.x = this.x + this.vx * ft  // Manual delta time math
ball.vx = 220  // Magic numbers everywhere
if (ball.x < ball.radius) { ... }  // Manual bounds checking
These should be language-level features, not library calls.

Tier 1: Must-Have Built-ins (Add to v1.0)
1. Time & Delta Time
javascript// Built-in time functions
deltaTime()     // Seconds since last frame (automatic)
time()          // Seconds since program start
timestamp()     // Unix timestamp
sleep(ms)       // Sleep for milliseconds
setTimeout(fn, ms)    // Call function after delay
setInterval(fn, ms)   // Call function repeatedly

// Usage in game loop:
func update() {
    let dt = deltaTime();  // Automatic frame timing
    player.x += player.vx * dt;
    player.y += player.vy * dt;
}
Implementation:
go// internal/runtime/time.go

var (
    programStart time.Time
    lastFrame    time.Time
)

func init() {
    programStart = time.Now()
    lastFrame = programStart
}

func native_deltaTime(args []Value) Value {
    now := time.Now()
    dt := now.Sub(lastFrame).Seconds()
    lastFrame = now
    return NumberVal(dt)
}

func native_time(args []Value) Value {
    elapsed := time.Since(programStart).Seconds()
    return NumberVal(elapsed)
}

func native_timestamp(args []Value) Value {
    return NumberVal(float64(time.Now().Unix()))
}

func native_sleep(args []Value) Value {
    if len(args) != 1 || !IsNumber(args[0]) {
        return NullVal
    }
    ms := AsNumber(args[0])
    time.Sleep(time.Duration(ms) * time.Millisecond)
    return NullVal
}

2. Random Numbers (Essential for Games)
javascript// Random functions
random()              // [0, 1)
random(max)           // [0, max)
random(min, max)      // [min, max)
randomInt(max)        // Integer [0, max)
randomInt(min, max)   // Integer [min, max)
randomChoice(array)   // Random element from array
randomSeed(seed)      // Set seed for reproducibility

// Game usage:
let spawnX = random(0, screenWidth);
let spawnY = random(0, screenHeight);
let enemyType = randomChoice(["goblin", "orc", "dragon"]);
let damage = randomInt(10, 20);

// Procedural generation:
randomSeed(12345);  // Same seed = same level
let level = generateLevel();
Implementation:
go// internal/runtime/random.go

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

func native_random(args []Value) Value {
    switch len(args) {
    case 0:
        // random() → [0, 1)
        return NumberVal(rng.Float64())
    
    case 1:
        // random(max) → [0, max)
        if !IsNumber(args[0]) {
            return NullVal
        }
        max := AsNumber(args[0])
        return NumberVal(rng.Float64() * max)
    
    case 2:
        // random(min, max) → [min, max)
        if !IsNumber(args[0]) || !IsNumber(args[1]) {
            return NullVal
        }
        min := AsNumber(args[0])
        max := AsNumber(args[1])
        return NumberVal(min + rng.Float64()*(max-min))
    
    default:
        return NullVal
    }
}

func native_randomInt(args []Value) Value {
    switch len(args) {
    case 1:
        // randomInt(max) → [0, max)
        if !IsNumber(args[0]) {
            return NullVal
        }
        max := int(AsNumber(args[0]))
        return NumberVal(float64(rng.Intn(max)))
    
    case 2:
        // randomInt(min, max) → [min, max)
        if !IsNumber(args[0]) || !IsNumber(args[1]) {
            return NullVal
        }
        min := int(AsNumber(args[0]))
        max := int(AsNumber(args[1]))
        return NumberVal(float64(min + rng.Intn(max-min)))
    
    default:
        return NullVal
    }
}

func native_randomChoice(args []Value) Value {
    if len(args) != 1 || !IsArray(args[0]) {
        return NullVal
    }
    
    arr := AsArray(args[0])
    if arr.Count == 0 {
        return NullVal
    }
    
    idx := rng.Intn(arr.Count)
    return arr.Values[idx]
}

func native_randomSeed(args []Value) Value {
    if len(args) != 1 || !IsNumber(args[0]) {
        return NullVal
    }
    
    seed := int64(AsNumber(args[0]))
    rng = rand.New(rand.NewSource(seed))
    return NullVal
}

3. Math Utilities (Game Math)
javascript// Core math
PI          // 3.14159...
E           // 2.71828...

// Trigonometry
sin(x)
cos(x)
tan(x)
asin(x)
acos(x)
atan(x)
atan2(y, x)

// Exponential
pow(base, exp)
sqrt(x)
exp(x)
log(x)
log10(x)

// Rounding
floor(x)
ceil(x)
round(x)
trunc(x)

// Utilities
abs(x)
sign(x)        // -1, 0, or 1
min(a, b, ...)
max(a, b, ...)
clamp(val, min, max)
lerp(a, b, t)      // Linear interpolation
smoothstep(a, b, t) // Smooth interpolation
map(val, inMin, inMax, outMin, outMax)  // Remap range

// Distance & angles
distance(x1, y1, x2, y2)
distanceSq(x1, y1, x2, y2)  // Faster (no sqrt)
angleBetween(x1, y1, x2, y2)
normalize(x, y)  // Returns {x, y} normalized

// Game usage:
let angle = atan2(target.y - player.y, target.x - player.x);
player.vx = cos(angle) * speed;
player.vy = sin(angle) * speed;

let dist = distance(player.x, player.y, enemy.x, enemy.y);
if (dist < 50) {
    takeDamage();
}

let health = clamp(player.health, 0, 100);
let alpha = lerp(0, 255, health / 100);  // Fade based on health
Implementation:
go// internal/runtime/math.go

func native_lerp(args []Value) Value {
    if len(args) != 3 {
        return NullVal
    }
    a := AsNumber(args[0])
    b := AsNumber(args[1])
    t := AsNumber(args[2])
    return NumberVal(a + (b-a)*t)
}

func native_clamp(args []Value) Value {
    if len(args) != 3 {
        return NullVal
    }
    val := AsNumber(args[0])
    min := AsNumber(args[1])
    max := AsNumber(args[2])
    
    if val < min {
        return NumberVal(min)
    }
    if val > max {
        return NumberVal(max)
    }
    return NumberVal(val)
}

func native_distance(args []Value) Value {
    if len(args) != 4 {
        return NullVal
    }
    x1 := AsNumber(args[0])
    y1 := AsNumber(args[1])
    x2 := AsNumber(args[2])
    y2 := AsNumber(args[3])
    
    dx := x2 - x1
    dy := y2 - y1
    return NumberVal(math.Sqrt(dx*dx + dy*dy))
}

func native_angleBetween(args []Value) Value {
    if len(args) != 4 {
        return NullVal
    }
    x1 := AsNumber(args[0])
    y1 := AsNumber(args[1])
    x2 := AsNumber(args[2])
    y2 := AsNumber(args[3])
    
    return NumberVal(math.Atan2(y2-y1, x2-x1))
}

func native_map(args []Value) Value {
    if len(args) != 5 {
        return NullVal
    }
    val := AsNumber(args[0])
    inMin := AsNumber(args[1])
    inMax := AsNumber(args[2])
    outMin := AsNumber(args[3])
    outMax := AsNumber(args[4])
    
    normalized := (val - inMin) / (inMax - inMin)
    return NumberVal(outMin + normalized*(outMax-outMin))
}

4. Array Utilities (Critical)
javascript// Array methods (your code desperately needs these)
arr.map(fn)
arr.filter(fn)
arr.forEach(fn)
arr.find(fn)
arr.findIndex(fn)
arr.some(fn)
arr.every(fn)
arr.reduce(fn, initial)
arr.sort(fn)
arr.reverse()
arr.slice(start, end)
arr.concat(other)
arr.indexOf(item)
arr.includes(item)

// Your Breakout code would become:
let aliveBricks = this.bricks.filter(b => b.alive());
let allDead = this.bricks.every(b => !b.alive());

// Instead of:
let i = 0;
while (i < this.bricks.length()) {
    if (this.bricks[i].alive()) { ... }
    i = i + 1;
}

5. String Utilities
javascript// String methods
str.split(sep)
str.trim()
str.upper()
str.lower()
str.startsWith(prefix)
str.endsWith(suffix)
str.indexOf(substr)
str.slice(start, end)
str.replace(old, new)
str.replaceAll(old, new)

// String formatting
format("Player: {}, Score: {}", name, score)
// Or template literals (better):
`Player: ${name}, Score: ${score}`

Tier 2: Nice-to-Have (v1.1)
6. Wait/Async Utilities
javascript// Coroutine-style waiting
wait(seconds)           // Pause execution
waitUntil(condition)    // Wait for condition to be true
waitForFrame()          // Wait one frame

// Game usage:
func enemyAI() {
    while (true) {
        patrol();
        wait(2);  // Patrol for 2 seconds
        
        if (playerVisible()) {
            chase();
            wait(5);
        }
    }
}

func cutscene() {
    showText("Once upon a time...");
    wait(3);
    showText("A hero arose...");
    wait(3);
    startGame();
}
Note: This requires coroutine support or async/await, which is complex. Consider for v1.1+.

7. Input Helpers
javascript// Keyboard (if not using raylib)
isKeyDown(key)
isKeyPressed(key)
isKeyReleased(key)

// Mouse
getMouseX()
getMouseY()
getMousePosition()  // Returns {x, y}
isMouseButtonDown(button)
isMouseButtonPressed(button)

// Gamepad
getGamepadAxis(gamepad, axis)
isGamepadButtonDown(gamepad, button)

8. Debug Utilities
javascript// Debugging
assert(condition, message)
trace(...args)           // Debug print
profile(name, fn)        // Profile function execution
benchmark(fn, iterations)

// Usage:
assert(player.health >= 0, "Health cannot be negative!");

profile("update", func() {
    updatePhysics();
    updateAI();
});

Tier 3: Advanced (v2.0+)
9. Collections
javascript// Set
let enemies = Set();
enemies.add(goblin);
enemies.has(goblin);
enemies.delete(goblin);

// Map
let scores = Map();
scores.set("Alice", 100);
scores.get("Alice");

10. File I/O
javascript// File operations
readFile(path)
writeFile(path, content)
appendFile(path, content)
fileExists(path)
deleteFile(path)

// JSON
parseJSON(str)
toJSON(obj)

// Save game example:
let saveData = {
    level: currentLevel,
    score: score,
    inventory: items
};
writeFile("save.json", toJSON(saveData));

// Load game:
let data = parseJSON(readFile("save.json"));
currentLevel = data.level;

Updated Standard Library (v1.0)
Complete Built-ins List
javascript// ============================================================================
// TIME & TIMING
// ============================================================================
deltaTime()              // Seconds since last frame
time()                   // Seconds since program start
timestamp()              // Unix timestamp
sleep(ms)                // Sleep milliseconds

// ============================================================================
// RANDOM
// ============================================================================
random()                 // [0, 1)
random(max)              // [0, max)
random(min, max)         // [min, max)
randomInt(max)           // Integer [0, max)
randomInt(min, max)      // Integer [min, max)
randomChoice(array)      // Random element
randomSeed(seed)         // Set seed

// ============================================================================
// MATH CONSTANTS
// ============================================================================
PI                       // 3.14159265358979
E                        // 2.71828182845905

// ============================================================================
// MATH - TRIGONOMETRY
// ============================================================================
sin(x)
cos(x)
tan(x)
asin(x)
acos(x)
atan(x)
atan2(y, x)

// ============================================================================
// MATH - EXPONENTIAL
// ============================================================================
pow(base, exp)
sqrt(x)
exp(x)
log(x)
log10(x)

// ============================================================================
// MATH - ROUNDING
// ============================================================================
floor(x)
ceil(x)
round(x)
trunc(x)

// ============================================================================
// MATH - UTILITIES
// ============================================================================
abs(x)
sign(x)
min(...values)
max(...values)
clamp(val, min, max)
lerp(a, b, t)
smoothstep(a, b, t)
map(val, inMin, inMax, outMin, outMax)

// ============================================================================
// MATH - GEOMETRY
// ============================================================================
distance(x1, y1, x2, y2)
distanceSq(x1, y1, x2, y2)
angleBetween(x1, y1, x2, y2)
normalize(x, y)          // Returns {x, y}

// ============================================================================
// ARRAY METHODS
// ============================================================================
arr.length()
arr.push(item)
arr.pop()
arr.map(fn)
arr.filter(fn)
arr.forEach(fn)
arr.find(fn)
arr.findIndex(fn)
arr.some(fn)
arr.every(fn)
arr.reduce(fn, initial)
arr.sort(compareFn)
arr.reverse()
arr.slice(start, end)
arr.concat(other)
arr.indexOf(item)
arr.includes(item)

// ============================================================================
// STRING METHODS
// ============================================================================
str.length()
str.split(sep)
str.trim()
str.upper()
str.lower()
str.startsWith(prefix)
str.endsWith(suffix)
str.indexOf(substr)
str.slice(start, end)
str.replace(old, new)
str.replaceAll(old, new)

// ============================================================================
// TYPE CHECKING
// ============================================================================
type(value)              // Returns "number", "string", etc.
isNumber(value)
isString(value)
isBool(value)
isNull(value)
isArray(value)
isObject(value)
isFunction(value)

// ============================================================================
// CONVERSION
// ============================================================================
number(value)
string(value)
bool(value)

// ============================================================================
// I/O
// ============================================================================
print(...values)
format(template, ...values)

// ============================================================================
// DEBUG
// ============================================================================
assert(condition, message)
trace(...values)

Your Breakout Code - Before/After
Before (Current)
javascriptfunc update(ft) {
    this.paddle.update(this.rl, ft, this.sw, this.KEY_LEFT, this.KEY_RIGHT);
    this.ball.step(ft);
    this.ball.bounceWalls(this.sw, this.sh);
    this.ball.bouncePaddle(this.paddle);
    this.ball.resolveBricks(this.bricks);
    if (this.ball.pastBottom(this.sh)) {
        this.ball.reset(this.sw, this.sh);
    }
}

func resolveBricks(bricks) {
    let i = 0;
    while (i < bricks.length()) {
        let b = bricks[i];
        if (b.alive()) {
            if (this.hitsRect(b.x, b.y, b.w, b.h)) {
                b.hit();
                this.vy = -this.vy;
            }
        }
        i = i + 1;
    }
}
After (With New Built-ins)
javascriptfunc update() {
    let dt = deltaTime();  // Automatic!
    
    this.paddle.update(dt);
    this.ball.step(dt);
    this.ball.bounceWalls();
    this.ball.bouncePaddle(this.paddle);
    this.ball.resolveBricks(this.bricks);
    
    if (this.ball.pastBottom()) {
        this.ball.reset();
    }
    
    // Check win condition
    if (this.bricks.every(b => !b.alive())) {
        this.nextLevel();
    }
}

func resolveBricks(bricks) {
    // Much cleaner!
    bricks
        .filter(b => b.alive())
        .find(b => this.hitsRect(b.x, b.y, b.w, b.h))
        ?.hit();  // Optional chaining (v1.1)
    
    // Or explicit:
    let hit = bricks.find(b => 
        b.alive() && this.hitsRect(b.x, b.y, b.w, b.h)
    );
    
    if (hit) {
        hit.hit();
        this.vy = -this.vy;
    }
}

func spawnPowerup() {
    let x = random(20, this.sw - 20);
    let y = random(20, this.sh - 20);
    let type = randomChoice(["extra_ball", "laser", "expand_paddle"]);
    
    return createPowerup(x, y, type);
}

Implementation Priority
Week 1: Core Time/Random/Math
go// Add to internal/runtime/

// time.go
- deltaTime()
- time()
- timestamp()
- sleep()

// random.go
- random() variants
- randomInt() variants
- randomChoice()
- randomSeed()

// math.go
- PI, E constants
- All trig functions
- Rounding functions
- Utility functions (clamp, lerp, distance, etc.)
Tests:
javascript// tests/stdlib/time_test.koda
let start = time();
sleep(100);
let end = time();
assert(end - start >= 0.1, "Sleep failed");

// tests/stdlib/random_test.koda
randomSeed(12345);
let a = random();
randomSeed(12345);
let b = random();
assert(a == b, "Seed failed");

// tests/stdlib/math_test.koda
assert(abs(-5) == 5);
assert(clamp(15, 0, 10) == 10);
assert(distance(0, 0, 3, 4) == 5);

Week 2: Array/String Methods
go// Add to internal/runtime/

// array.go
- All array methods as methods on array objects

// string.go
- All string methods as methods on string objects

Week 3: Debug/Format
go// debug.go
- assert()
- trace()

// format.go
- format() function

Final Recommendation
Add to v1.0 (Essential):

✅ Time functions (deltaTime, time, sleep)
✅ Random functions (all variants)
✅ Math functions (complete set)
✅ Array methods (map, filter, etc.)
✅ String methods (split, trim, etc.)
✅ format() function

Add to v1.1 (Quality of Life):

🔵 Template literals
🔵 Arrow functions
🔵 Optional chaining
🔵 Coroutines/wait

Add to v2.0 (Advanced):

🔵 Async/await
🔵 File I/O
🔵 Collections (Set, Map)
🔵 Modules (import/export)


## Standard library wishlist (notes)

Possible directions for builtins and stdlib (for planning only):

## Standard Library - Built-in Functions

Implement these essential functions for game/app development:

### Time & Timing
- deltaTime() - Automatic frame delta
- time() - Program runtime
- timestamp() - Unix timestamp
- sleep(ms) - Pause execution

### Random
- random(), random(max), random(min, max)
- randomInt(max), randomInt(min, max)
- randomChoice(array)
- randomSeed(seed)

### Math
- Constants: PI, E
- Trig: sin, cos, tan, asin, acos, atan, atan2
- Exp: pow, sqrt, exp, log, log10
- Rounding: floor, ceil, round, trunc
- Utils: abs, sign, min, max, clamp, lerp, map
- Geometry: distance, distanceSq, angleBetween

### Arrays
Implement as methods: map, filter, forEach, find, some, every, reduce, etc.

### Strings
Implement as methods: split, trim, upper, lower, startsWith, endsWith, etc.

These are critical for game development. Implement in internal/runtime/ package
