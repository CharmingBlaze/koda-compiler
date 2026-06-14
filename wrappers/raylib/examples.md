# raylib - Examples

Copy any snippet into your Koda program and adjust the arguments.

---

## Include the library

```koda
#include "raylib.koda"
```

---

## Functions

### void

```koda
let result = void();
print(result);
```

### bool

```koda
let result = bool();
print(result);
```

### bool

```koda
let result = bool();
print(result);
```

### InitWindow

```koda
InitWindow(width, height, title);
```

### CloseWindow

```koda
CloseWindow();
```

### WindowShouldClose

```koda
let result = WindowShouldClose();
print(result);
```

### IsWindowReady

```koda
let result = IsWindowReady();
print(result);
```

### IsWindowFullscreen

```koda
let result = IsWindowFullscreen();
print(result);
```

*See [api_reference.md](api_reference.md) for all 554 functions.*

---

## Structs

Structs are passed as Koda objects with matching field names:

```koda
let obj = { x: 0,  y: 0 };
```

---

## Enum values

```koda
let bool_false = 0;
let bool_true = 1;
```

