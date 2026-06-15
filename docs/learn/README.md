# Learn Koda

A **chapter-by-chapter** path from zero to shipping a small program. Each chapter is short (~10 minutes). Read in order or jump to what you need.

> **Prefer one file?** See the [Beginner's guide](../beginners-guide.md).

| Chapter | Topic | You will |
|---------|-------|----------|
| [01 — Welcome](01-welcome.md) | What Koda is | Understand goals and workflow |
| [02 — Install](02-install-and-first-run.md) | Setup | Install SDK and run hello world |
| [03 — Variables](03-variables-and-types.md) | Types | Use `let`, `const`, types, string interpolation |
| [04 — Control flow](04-control-flow.md) | Logic | Write if, while, for, switch, match |
| [05 — Functions](05-functions.md) | Functions | Define, return, closures |
| [07 — Structs & enums](07-structs-and-enums.md) | Game data | **Structs first** — players, enemies, phases |
| [06 — Objects & arrays](06-objects-and-arrays.md) | Lists & JSON | Arrays and config objects |
| [08 — Modules](08-modules-and-imports.md) | Code organization | `import "@math"`, local modules |
| [09 — Files & JSON](09-files-and-json.md) | Persistence | Read/write files, parse JSON |
| [10 — Ship it](10-building-and-shipping.md) | CLI | build, bundle, test, clean |

> **Note:** Read chapter **7 before chapter 6**. Structs model game state; objects are for JSON and config.

---

## After this path

| Goal | Guide |
|------|-------|
| Games | [Game development](../guides/game-dev.md) |
| Desktop / CLI apps | [Applications](../guides/applications.md) |
| C background | [From C](../guides/from-c.md) |
| Full syntax lookup | [Language reference](../language.md) |
| Every CLI flag | [CLI reference](../reference/cli.md) |

---

## How to practice

1. Type examples yourself — do not only read.
2. Change one line and predict the output before running.
3. Run `koda check` when something fails to parse.
4. Run `koda doctor` if builds fail on your machine.
5. Keep a `playground.koda` file for experiments.
