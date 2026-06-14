# Chapter 7 — Structs and enums

**You will learn:** typed structs for game entities and enums for states.

**Time:** ~10 minutes.

---

## Structs

Declare fields without repeating types (all fields share struct layout rules):

```koda
struct Enemy {
    x, y, hp, speed
}

let e = Enemy { x: 0, y: 0, hp: 50, speed: 80 };
e.x = e.x + e.speed * 0.016;
```

Structs are ideal for **hot game data** — fixed fields, compile-time field names.

---

## Enums

```koda
enum AIState {
    Idle, Patrol, Chase, Attack
}

let state = AIState.Patrol;

if (state == AIState.Chase) {
    print("chasing");
}
```

Enum members are compared by identity like named constants.

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
let team = Team.Enemy;

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
| Fixed set of fields | Dynamic keys |
| Game entities, math records | Config maps, JSON-like data |
| Compile-time field checks | Schema varies at runtime |

---

## Try it yourself

Define a `Player` struct and an enum `GamePhase` with `Menu`, `Play`, `Pause`. Switch on phase in `func main()`.

---

## Next chapter

[Chapter 8 — Modules and imports](08-modules-and-imports.md)
