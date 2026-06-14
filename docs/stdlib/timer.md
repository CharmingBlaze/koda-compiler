# `@timer` — cooldowns, intervals, countdowns

**Include:** `#include "stdlib/timer.koda"`  
**Import:** `import "@timer"` (when exposed in your SDK layout)

---

## Countdown timer

Object-based elapsed timer toward a duration.

| Function | Purpose |
|----------|---------|
| `create(duration)` | New timer object |
| `update(t, dt)` | Add delta seconds |
| `done(t)` | True when elapsed ≥ duration |
| `reset(t)` | Clear elapsed |
| `progress(t)` | 0..1 ratio |
| `remaining(t)` | Seconds left |

```koda
#include "stdlib/timer.koda"

let intro = create(3.0);
intro = update(intro, deltatime());
if (done(intro)) {
    start_game();
}
```

---

## Cooldown (one-shot gate)

Uses `programtime()` — good for fire rates.

| Function | Purpose |
|----------|---------|
| `cooldown(seconds)` | Create cooldown state |
| `cooldown_try(c)` | True once per interval when called |
| `cooldown_reset(c)` | Reset gate |

```koda
let fire = cooldown(0.25);
if (cooldown_try(fire)) {
    shoot();
}
```

---

## Interval (repeating)

| Function | Purpose |
|----------|---------|
| `interval(seconds)` | Create interval state |
| `interval_tick(t)` | True every `seconds` |
| `interval_reset(t)` | Reset schedule |

```koda
let spawn = interval(2.0);
if (interval_tick(spawn)) {
    spawn_enemy();
}
```

---

## `wait_ms(ms)`

Wrapper around `sleep(ms)` for scripted delays (blocks the thread).

---

## Related

- [Game dev guide](../guides/game-dev.md)
