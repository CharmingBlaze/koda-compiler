# Chapter 5 — Functions

**You will learn:** defining functions, return values, anonymous functions, and closures.

**Time:** ~10 minutes.

---

## Defining and calling

```koda
func add(a, b) {
    return a + b;
}

print(add(2, 3));

func greet(name) {
  print("Hello, " + name);
}

greet("Ada");
```

`func main()` is a common entry point:

```koda
func main() {
    print("Starting...");
}
```

---

## Anonymous functions

```koda
let square = func(x) {
    return x * x;
};

print(square(4));
```

Functions are values — pass them to other functions:

```koda
func apply(fn, x) {
    return fn(x);
}

print(apply(square, 5));
```

---

## Closures

Inner functions can read and update outer locals:

```koda
func makeAccumulator(start) {
    let total = start;
    return func(n) {
        total = total + n;
        return total;
    };
}

let acc = makeAccumulator(10);
print(acc(5));   // 15
print(acc(3));   // 18
```

---

## Default and rest parameters

```koda
func power(base, exp) {
    if (exp == null) {
        exp = 2;
    }
    return pow(base, exp);
}
```

See [Language reference](../language.md) for `...rest` parameters.

---

## Try it yourself

Write `func clamp(value, low, high)` that returns `value` bounded between `low` and `high`. (Or use `math.clamp` after `import "@math"`.)

---

## Next chapter

[Chapter 6 — Objects and arrays](06-objects-and-arrays.md)
