# Chapter 4 — Control flow

**You will learn:** `if`, loops, `switch`, `match`, `break`, and `continue`.

**Time:** ~10 minutes.

---

## If / else

```koda
let hp = 35;

if (hp <= 0) {
    print("dead");
} else if (hp < 50) {
    print("hurt");
} else {
    print("healthy");
}
```

Parentheses around the condition are required.

---

## While and do-while

```koda
let i = 0;
while (i < 5) {
    print(i);
    i = i + 1;
}

let j = 0;
do {
    print(j);
    j = j + 1;
} while (j < 3);
```

---

## Classic for

```koda
for (let i = 0; i < 3; i = i + 1) {
    print(i);
}
```

---

## For-of iteration

```koda
for (let item of ["a", "b", "c"]) {
    print(item);
}

let cfg = { width: 800, height: 600 };
for (let key in cfg) {
    print(key, cfg[key]);
}

for (let k, v of cfg) {
    print(k, v);
}
```

---

## Switch (C-style)

Cases do **not** fall through. Use `fallthrough;` when you need the next case to run.

```koda
let weapon = 2;
switch (weapon) {
    case 1:
        print("sword");
        break;
    case 2:
        print("bow");
        break;
    default:
        print("fists");
}
```

---

## Match (game states)

Brace-style dispatch — cleaner than long `if` / `else if` chains:

```koda
enum GameState { Playing, Won, GameOver }

let state = GameState.Playing;

match state {
    GameState.Playing {
        update_game(dt);
    }
    GameState.Won {
        draw.text("STAR GET!", 380, 340, 40, colors.yellow);
    }
    GameState.GameOver {
        draw.text("GAME OVER - press R", 400, 340, 36, colors.red);
    }
}
```

Each arm is its own block — no fall-through. Classic `switch (x) { case …: … }` still works.

---

## Break and continue

```koda
for (let n of range(0, 10)) {
    if (n == 5) {
        break;
    }
    if (n % 2 == 0) {
        continue;
    }
    print(n);
}
```

(`range` from `import "@array"` or `#include "stdlib/array.koda"`.)

---

## Try it yourself

Write a loop that prints numbers 1–10 but skips 7 and stops at 9.

---

## Next chapter

[Chapter 5 — Functions](05-functions.md)
