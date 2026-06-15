# Pong (Koda)

Classic two-player Pong written in **Koda** — structs, enums, functions, and the `@game` stdlib module.

## Controls

| Player | Keys |
|--------|------|
| Left | **W** / **S** |
| Right | **↑** / **↓** |
| Serve / restart | **Space** |

First to **11** wins.

## Run

From this folder:

```bash
koda run
```

Or create a new project from the template:

```bash
koda new mypong --template pong
cd mypong
koda run
```

In **Koda Studio**: pick the **Pong** template, then press **F5**.

## Source

All game logic is in `src/main.koda` — paddles and ball as `struct`, game flow as `enum Phase`, physics in plain Koda functions.
