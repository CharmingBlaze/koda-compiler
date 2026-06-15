# @pool — object pool

Reuse object instances to reduce GC pressure in hot loops. Ships in `stdlib/pool.koda`.

```koda
#include "@pool"

let bullets = pool(64, func() { return { x: 0, y: 0, active: false }; });
let b = get(bullets);
release(bullets, b);
```

Functions: `pool(size, makeItem)`, `get(p)`, `release(p, obj)`, `reset(p)`.
