# Chapter 7 — Structs and enums

**You will learn:** structs as the main way to model game and app state, plus enums for phases.

**Time:** ~15 minutes.

> **Read this before chapter 6.** Structs are for game/app data. Objects are for JSON and config.

---

## Structs

Declare named fields — the compiler checks field names at compile time:

```koda
struct Mario {
    x, y,
    speed,
    health
}

let player = Mario {
    x: 400,
    y: 300,
    speed: 220,
    health: 100
};
```

Identifiers are **case-insensitive** — `struct Player` and `let player` are the same name and cause a duplicate-binding error. Use a different type name (`Mario`, `Hero`) or variable name.

Update fields in your game loop:

```koda
func updatePlayer(player, dt) {
    if (isKeyDown(KEY_RIGHT)) {
        player.x = player.x + player.speed * dt;
    }
}
```

Structs are ideal for **hot game data** — fixed fields, fast access, helpful errors if you typo a field name.

### Field defaults

Omit fields in a literal when they have defaults:

```koda
struct Coin {
    x = 0.0;
    z = 0.0;
    on = true;
}

let coin = Coin { x: 7, z: -5 };   // on is true automatically
```

### Constructors and methods

Define `func new(...)` inside the struct, then call `Coin(7, -5)`:

```koda
struct Coin {
    x = 0;
    z = 0;
    on = true;

    func new(x, z) {
        this.x = x;
        this.z = z;
    }

    func collect(player) {
        if on && player.distance_xz(x, z) < 1.4 {
            on = false;
        }
    }
}

let coins = [Coin(7, -5), Coin(-8, 3)];

for coin in coins {
    coin.collect(player);
}
```

Inside methods, bare names like `x` and `on` mean the struct's own fields.

---

## Constants

Use `const` for values that never change:

```koda
const gravity = 900;
const screenWidth = 800;
```

| Keyword | Meaning |
|---------|---------|
| `let` | Mutable binding |
| `const` | Immutable — cannot reassign |

---

## Optional type annotations

Beginners can omit types. Add them when you want clarity or integer math:

```koda
let score = 0;              // inferred number
let lives: int = 3;         // optional
let dt: float = 0.016;      // optional
let name: string = "Jesse"; // optional
```

Core beginner types: `int`, `float`, `bool`, `string`, `array`, `map`, `func`, `object`, `byte`.

---

## Enums

```koda
enum GamePhase {
    Menu, Play, Pause
}

let phase = GamePhase.Play;

match phase {
    GamePhase.Menu {
        draw_menu();
    }
    GamePhase.Play {
        update_game(dt);
    }
    GamePhase.Pause {
        draw_pause_overlay();
    }
}
```

Members are integers from `0` upward. Use `Type.Member` (e.g. `GamePhase.Play`).

---

## Combining structs and enums

```koda
struct Entity {
    x, y, hp
}

enum Team {
    Player, Enemy, Neutral
}

let bot = Entity { x: 10, y: 10, hp: 30 };

func updateEntity(ent, dt) {
    if (ent.hp <= 0) {
        return false;
    }
    ent.x = ent.x + 50 * dt;
    return true;
}
```

---

## Struct vs object

| Use struct when… | Use object when… |
|------------------|------------------|
| Player, enemy, bullet state | JSON from a file |
| Fixed set of fields | Config key/value maps |
| Game entities | Parsed API responses |

---

## Try it yourself

Define a `Player` struct and a `GamePhase` enum. In `func main()`, create a player and move them when phase is `Play`.

---

## Next chapter

[Chapter 8 — Modules and imports](08-modules-and-imports.md)
