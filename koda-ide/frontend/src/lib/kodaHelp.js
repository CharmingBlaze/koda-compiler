/** Koda language reference for completions, hovers, and inline help in Studio. */

export const KODA_KEYWORDS = [
  ['func', 'Define a function.', 'func name(a, b) {\n    return a + b;\n}'],
  ['let', 'Mutable variable binding.', 'let x = 1;\nx = x + 1;'],
  ['const', 'Immutable binding (cannot reassign).', 'const gravity = 900;'],
  ['if', 'Conditional branch.', 'if (cond) {\n    ...\n} else {\n    ...\n}'],
  ['else', 'Alternative branch for if.', 'if (x) { ... } else { ... }'],
  ['while', 'Loop while condition is true.', 'while (i < 10) {\n    i = i + 1;\n}'],
  ['loop', 'Infinite loop — use break to exit.', 'loop {\n    if (done) { break; }\n}'],
  ['for', 'Loop over iterable values.', 'for (let n of items) {\n    print(n);\n}'],
  ['of', 'Used with for…of loops.', 'for (let x of arr) { ... }'],
  ['in', 'Membership / iteration helper.', 'for (let k in obj) { ... }'],
  ['return', 'Return a value from a function.', 'return result;'],
  ['break', 'Exit a loop or switch case.', 'break;'],
  ['continue', 'Skip to next loop iteration.', 'continue;'],
  ['switch', 'Multi-way branch (use break).', 'switch (x) {\n    case 1: break;\n    default: break;\n}'],
  ['case', 'Switch branch label.', 'case 1:\n    print("one");\n    break;'],
  ['default', 'Switch fallback branch.', 'default:\n    print("other");'],
  ['struct', 'Named data type for game/app state.', 'struct Player {\n    x, y, health\n}'],
  ['enum', 'Named constants group.', 'enum State { Idle, Running, Dead }'],
  ['use', 'Import a module (official).', 'use raylib;\nuse koda.game;'],
  ['import', 'Load a stdlib or local module.', 'let math = import "@math";'],
  ['include', 'Include another .koda file (#include, legacy).', '#include "helpers.koda"'],
  ['defer', 'Run cleanup when scope exits.', 'defer cleanup();'],
  ['delete', 'Remove a property from an object.', 'delete obj.key;'],
  ['test', 'Named test block.', 'test "adds numbers" {\n    assert(add(1, 2) == 3);\n}'],
  ['match', 'Pattern match on enum/value.', 'match state {\n    Phase.Playing { update(); }\n    default { drawMenu(); }\n}'],
  ['fallthrough', 'Fall through to next switch case.', 'case 1:\n    print("one");\n    fallthrough;'],
  ['step', 'Range step in for-in.', 'for i in 0..100 step 5 {\n    print(i);\n}'],
  ['and', 'Logical AND (word alias).', 'if (alive and health > 0) { ... }'],
  ['or', 'Logical OR (word alias).', 'let dt = game.delta() or 0.016;'],
  ['not', 'Logical NOT (word alias).', 'while not WindowShouldClose() { }'],
  ['this', 'Current instance in struct methods.', 'this.x += speed * dt;'],
  ['typeof', 'Type name of a value.', 'typeof(x)  // "number", "string", ...'],
  ['true', 'Boolean true literal.', 'let ok = true;'],
  ['false', 'Boolean false literal.', 'let ok = false;'],
  ['null', 'Empty / missing value.', 'let empty = null;'],
  ['do', 'Do-while style loop.', 'do { ... } while (cond);'],
  ['var', 'Reserved — use let instead.', 'let x = 1;  // not var'],
]

export const KODA_TYPES = [
  ['int', 'Integer type annotation.', 'let lives: int = 3;'],
  ['float', 'Floating-point annotation.', 'let dt: float = 0.016;'],
  ['string', 'String type annotation.', 'let name: string = "Ada";'],
  ['bool', 'Boolean type annotation.', 'let alive: bool = true;'],
  ['number', 'Numeric type (64-bit float).', 'let score: number = 42;'],
  ['void', 'No return value.', 'func log(msg): void { print(msg); }'],
  ['i32', '32-bit signed integer annotation.', 'let id: i32 = 0;'],
  ['u8', '8-bit unsigned integer annotation.', 'let channel: u8 = 255;'],
]

export const KODA_STDLIB_MODULES = [
  ['raylib', 'Full Raylib API (548 functions).', 'use raylib;'],
  ['koda.game', 'Game loop helpers over Raylib.', 'use koda.game;'],
  ['koda.input', 'Keyboard/mouse helpers.', 'use koda.input;'],
  ['koda.camera', '3D orbit / FPS camera.', 'use koda.camera;'],
  ['koda.color', 'asRaylib() color conversion.', 'use koda.color;'],
  ['koda.math', 'Math helpers (via @math).', 'use koda.math;'],
  ['@math', 'Math helpers (sqrt, sin, lerp, …).', 'use koda.math;'],
  ['@io', 'File read/write helpers.', 'use koda.io;'],
  ['@json', 'JSON parse/stringify.', 'use koda.json;'],
  ['@array', 'Array utilities (range, sum, shuffle).', 'use koda.array;'],
  ['@str', 'String utilities.', 'use koda.str;'],
  ['@util', 'General utilities.', 'use koda.util;'],
  ['@game', 'Alias for koda.game (legacy @ import).', 'use koda.game;'],
  ['@color', 'Color helpers + asRaylib.', 'use koda.color;'],
  ['@vec2', '2D vector math.', 'use koda.vec2;'],
  ['@vec3', '3D vector math.', 'use koda.vec3;'],
  ['@timer', 'Timers and scheduling.', 'use koda.timer;'],
  ['@easing', 'Easing functions for animation.', 'use koda.easing;'],
  ['@noise', 'Procedural noise.', 'use koda.noise;'],
  ['@input', 'Input helpers (legacy @ import).', 'use koda.input;'],
  ['@pool', 'Object pooling.', 'use koda.pool;'],
]

/** Built-in runtime functions (aligned with internal/sema/builtin_globals.go). */
export const KODA_BUILTINS = [
  ['print', 'print(value, …)', 'Write values to stdout (auto newline).', 'print("Hello", 42);'],
  ['len', 'len(value)', 'Length of string, array, or object keys.', 'len(items)'],
  ['type', 'type(value)', 'Type name as string.', 'type(name)  // "string"'],
  ['typeof', 'typeof(value)', 'Same as type().', 'typeof(42)'],
  ['assert', 'assert(cond, msg?)', 'Abort if condition is false.', 'assert(x > 0, "positive");'],
  ['expect', 'expect(cond, msg?)', 'Test assertion (non-fatal in tests).', 'expect(add(1,1) == 2);'],
  ['panic', 'panic(msg?)', 'Abort with error message.', 'panic("unexpected");'],
  ['warn', 'warn(msg)', 'Print warning to stderr.', 'warn("deprecated");'],
  ['ok', 'ok(value)', 'Wrap success result.', 'ok(result)'],
  ['err', 'err(message)', 'Wrap error result.', 'err("failed")'],
  ['format', 'format(fmt, …)', 'Format string with placeholders.', 'format("x=%d", x)'],
  ['string', 'string(value)', 'Convert to string.', 'string(42)'],
  ['number', 'number(value)', 'Convert to number.', 'number("3.14")'],
  ['bool', 'bool(value)', 'Convert to boolean.', 'bool(1)'],
  ['abs', 'abs(n)', 'Absolute value.', 'abs(-5)'],
  ['sqrt', 'sqrt(n)', 'Square root.', 'sqrt(16)'],
  ['cbrt', 'cbrt(n)', 'Cube root.', 'cbrt(27)'],
  ['min', 'min(a, b)', 'Smaller of two numbers.', 'min(a, b)'],
  ['max', 'max(a, b)', 'Larger of two numbers.', 'max(a, b)'],
  ['floor', 'floor(n)', 'Round down.', 'floor(3.7)'],
  ['ceil', 'ceil(n)', 'Round up.', 'ceil(3.2)'],
  ['round', 'round(n)', 'Round to nearest integer.', 'round(3.5)'],
  ['trunc', 'trunc(n)', 'Truncate toward zero.', 'trunc(-3.7)'],
  ['sign', 'sign(n)', 'Sign of number (-1, 0, 1).', 'sign(-3)'],
  ['clamp', 'clamp(v, lo, hi)', 'Constrain value to range.', 'clamp(x, 0, 100)'],
  ['lerp', 'lerp(a, b, t)', 'Linear interpolation.', 'lerp(0, 100, 0.5)'],
  ['map', 'map(v, inLo, inHi, outLo, outHi)', 'Remap value between ranges.', 'map(x, 0, 1, 0, 255)'],
  ['sin', 'sin(rad)', 'Sine (radians).', 'sin(pi / 2)'],
  ['cos', 'cos(rad)', 'Cosine (radians).', 'cos(0)'],
  ['tan', 'tan(rad)', 'Tangent (radians).', 'tan(pi / 4)'],
  ['asin', 'asin(n)', 'Arc sine.', 'asin(1)'],
  ['acos', 'acos(n)', 'Arc cosine.', 'acos(0)'],
  ['atan', 'atan(n)', 'Arc tangent.', 'atan(1)'],
  ['atan2', 'atan2(y, x)', 'Arc tangent of y/x.', 'atan2(1, 1)'],
  ['pow', 'pow(base, exp)', 'Power.', 'pow(2, 8)'],
  ['exp', 'exp(n)', 'e^n.', 'exp(1)'],
  ['log', 'log(n)', 'Natural logarithm.', 'log(e)'],
  ['log10', 'log10(n)', 'Base-10 logarithm.', 'log10(100)'],
  ['log2', 'log2(n)', 'Base-2 logarithm.', 'log2(8)'],
  ['pi', 'pi', 'Constant π.', 'let r = pi * 2;'],
  ['e', 'e', "Constant Euler's number.", 'exp(1)'],
  ['degrees', 'degrees(rad)', 'Radians to degrees.', 'degrees(pi)'],
  ['radians', 'radians(deg)', 'Degrees to radians.', 'radians(180)'],
  ['hypot', 'hypot(x, y)', 'Hypotenuse length.', 'hypot(3, 4)'],
  ['fmod', 'fmod(x, y)', 'Floating remainder.', 'fmod(7, 3)'],
  ['wrap', 'wrap(v, lo, hi)', 'Wrap value in range.', 'wrap(angle, 0, 360)'],
  ['approach', 'approach(cur, target, step)', 'Move toward target by step.', 'approach(x, 100, 5)'],
  ['smoothstep', 'smoothstep(edge0, edge1, x)', 'Smooth Hermite interpolation.', 'smoothstep(0, 1, t)'],
  ['smoothdamp', 'smoothdamp(cur, target, vel, smooth, dt)', 'Smooth damping.', 'smoothdamp(x, target, v, 0.3, dt)'],
  ['normalize', 'normalize(x, y)', 'Unit vector from components.', 'normalize(dx, dy)'],
  ['distance', 'distance(x1, y1, x2, y2)', 'Distance between points.', 'distance(0, 0, 3, 4)'],
  ['distancesq', 'distancesq(x1, y1, x2, y2)', 'Squared distance (faster).', 'distancesq(0, 0, 3, 4)'],
  ['anglebetween', 'anglebetween(x1, y1, x2, y2)', 'Angle between vectors.', 'anglebetween(1, 0, 0, 1)'],
  ['random', 'random()', 'Random float in [0, 1).', 'random()'],
  ['randomint', 'randomint(lo, hi)', 'Random integer in range.', 'randomint(1, 6)'],
  ['randomchoice', 'randomchoice(arr)', 'Pick random array element.', 'randomchoice(items)'],
  ['randomseed', 'randomseed(n)', 'Seed RNG.', 'randomseed(42)'],
  ['time', 'time()', 'Wall-clock seconds.', 'time()'],
  ['clock', 'clock()', 'High-resolution timer.', 'clock()'],
  ['timestamp', 'timestamp()', 'Unix timestamp.', 'timestamp()'],
  ['programtime', 'programtime()', 'Seconds since program start.', 'programtime()'],
  ['deltatime', 'deltatime()', 'Frame delta (game loop).', 'let dt = deltatime();'],
  ['sleep', 'sleep(seconds)', 'Pause execution.', 'sleep(0.5)'],
  ['readfile', 'readfile(path)', 'Read entire file as string.', 'readfile("data.txt")'],
  ['writefile', 'writefile(path, text)', 'Write string to file.', 'writefile("out.txt", data)'],
  ['appendfile', 'appendfile(path, text)', 'Append to file.', 'appendfile("log.txt", line)'],
  ['fileexists', 'fileexists(path)', 'True if path exists.', 'fileexists("save.dat")'],
  ['deletefile', 'deletefile(path)', 'Delete a file.', 'deletefile("tmp.txt")'],
  ['isfile', 'isfile(path)', 'True if path is a file.', 'isfile(p)'],
  ['isdir', 'isdir(path)', 'True if path is a directory.', 'isdir(p)'],
  ['filesize', 'filesize(path)', 'File size in bytes.', 'filesize("data.bin")'],
  ['listdir', 'listdir(path)', 'List directory entries.', 'listdir(".")'],
  ['parsejson', 'parsejson(text)', 'Parse JSON string to object.', 'parsejson(raw)'],
  ['tojson', 'tojson(value)', 'Serialize value to JSON.', 'tojson(obj)'],
  ['keys', 'keys(obj)', 'Object key names as array.', 'keys(config)'],
  ['matches', 'matches(value, pattern)', 'Pattern / type match helper.', 'matches(x, "number")'],
  ['replace', 'replace(s, old, new)', 'Replace first substring.', 'replace(s, "a", "b")'],
  ['replaceall', 'replaceall(s, old, new)', 'Replace all substrings.', 'replaceall(s, " ", "")'],
  ['trace', 'trace(…)', 'Debug print with location.', 'trace("here", x)'],
  ['isnumber', 'isnumber(v)', 'Type check.', 'isnumber(x)'],
  ['isstring', 'isstring(v)', 'Type check.', 'isstring(x)'],
  ['isbool', 'isbool(v)', 'Type check.', 'isbool(x)'],
  ['isnull', 'isnull(v)', 'Type check.', 'isnull(x)'],
  ['isarray', 'isarray(v)', 'Type check.', 'isarray(x)'],
  ['isobject', 'isobject(v)', 'Type check.', 'isobject(x)'],
  ['isfunction', 'isfunction(v)', 'Type check.', 'isfunction(x)'],
  ['arraypush', 'arraypush(arr, item)', 'Push onto array.', 'arraypush(items, "sword")'],
  ['arraypop', 'arraypop(arr)', 'Pop from array.', 'arraypop(stack)'],
  ['arrayslice', 'arrayslice(arr, start, end?)', 'Slice array.', 'arrayslice(a, 1, 3)'],
  ['arraysort', 'arraysort(arr)', 'Sort array in place.', 'arraysort(nums)'],
  ['arrayreverse', 'arrayreverse(arr)', 'Reverse array in place.', 'arrayreverse(a)'],
  ['arrayincludes', 'arrayincludes(arr, item)', 'Contains check.', 'arrayincludes(a, 2)'],
  ['arrayindexof', 'arrayindexof(arr, item)', 'Index of item.', 'arrayindexof(a, "x")'],
  ['arrayconcat', 'arrayconcat(a, b)', 'Concatenate arrays.', 'arrayconcat(a, b)'],
  ['assetpath', 'assetpath(rel)', 'Resolve asset path from project.', 'assetpath("sprites/hero.png")'],
  ['gc', 'gc()', 'Run garbage collection.', 'gc()'],
  ['gccollect', 'gccollect()', 'Force GC collection.', 'gccollect()'],
  ['arena', 'arena()', 'Create memory arena.', 'let a = arena();'],
  ['arenareset', 'arenareset(a)', 'Reset arena.', 'arenareset(a)'],
]

const keywordMap = new Map()
const builtinMap = new Map()
const moduleMap = new Map()
const typeMap = new Map()

for (const [name, desc, example] of KODA_KEYWORDS) {
  keywordMap.set(name.toLowerCase(), { kind: 'keyword', name, signature: name, desc, example })
}
for (const [name, desc, example] of KODA_TYPES) {
  typeMap.set(name.toLowerCase(), { kind: 'type', name, signature: name, desc, example })
}
for (const [name, sig, desc, example] of KODA_BUILTINS) {
  builtinMap.set(name.toLowerCase(), { kind: 'builtin', name, signature: sig, desc, example })
}
for (const [name, desc, example] of KODA_STDLIB_MODULES) {
  moduleMap.set(name.toLowerCase(), { kind: 'module', name, signature: name, desc, example })
}

/** @param {string} word */
export function lookupKodaHelp(word) {
  const key = (word || '').toLowerCase()
  return keywordMap.get(key) || builtinMap.get(key) || typeMap.get(key) || moduleMap.get(key) || null
}

/** @param {{ kind: string, name: string, signature?: string, desc: string, example?: string }} entry */
export function renderHelpDom(entry) {
  const root = document.createElement('div')
  root.className = 'koda-help-tooltip'
  root.style.cssText =
    'max-width:22rem;padding:0.35rem 0.15rem;font-family:var(--font-sans);font-size:13px;line-height:1.5;color:var(--color-text);'

  const title = document.createElement('div')
  title.style.cssText = 'font-family:var(--font-mono);font-weight:600;color:var(--color-cyan);margin-bottom:0.25rem;'
  title.textContent = entry.signature || entry.name
  root.appendChild(title)

  const desc = document.createElement('div')
  desc.style.cssText = 'color:var(--color-subtext);margin-bottom:0.35rem;'
  desc.textContent = entry.desc
  root.appendChild(desc)

  if (entry.example) {
    const ex = document.createElement('pre')
    ex.style.cssText =
      'margin:0;padding:0.5rem 0.625rem;border-radius:6px;background:var(--color-surface0);font-family:var(--font-mono);font-size:12px;color:var(--color-text);white-space:pre-wrap;'
    ex.textContent = entry.example
    root.appendChild(ex)
  }

  const kind = document.createElement('div')
  kind.style.cssText = 'margin-top:0.35rem;font-size:11px;color:var(--color-overlay0);text-transform:uppercase;letter-spacing:0.04em;'
  kind.textContent = entry.kind
  root.appendChild(kind)

  return root
}

/** Scan source for user-defined symbols. */
export function scanDocumentSymbols(text) {
  /** @type {{ label: string, type: string, detail: string, info?: () => HTMLElement }[]} */
  const out = []
  const seen = new Set()

  const add = (label, type, detail) => {
    const key = label.toLowerCase()
    if (seen.has(key)) return
    seen.add(key)
    out.push({
      label,
      type,
      detail,
      info: () =>
        renderHelpDom({
          kind: type,
          name: label,
          signature: label,
          desc: detail,
        }),
    })
  }

  let m
  const funcRe = /\bfunc\s+([A-Za-z_]\w*)\s*\(/g
  while ((m = funcRe.exec(text))) add(m[1], 'function', 'Function in this file')

  const structRe = /\bstruct\s+([A-Za-z_]\w*)\b/g
  while ((m = structRe.exec(text))) add(m[1], 'class', 'Struct defined in this file')

  const enumRe = /\benum\s+([A-Za-z_]\w*)\b/g
  while ((m = enumRe.exec(text))) add(m[1], 'enum', 'Enum defined in this file')

  const letRe = /\blet\s+([A-Za-z_]\w*)/g
  while ((m = letRe.exec(text))) add(m[1], 'variable', 'Local variable')

  const constRe = /\bconst\s+([A-Za-z_]\w*)/g
  while ((m = constRe.exec(text))) add(m[1], 'constant', 'Constant binding')

  return out
}

export const KODA_SNIPPETS = [
  {
    label: 'func',
    type: 'snippet',
    detail: 'function',
    apply: 'func ${name}(${args}) {\n\t${}\n}',
    info: () => renderHelpDom({ kind: 'snippet', name: 'func', desc: 'New function block', example: 'func update(dt) {\n    ...\n}' }),
  },
  {
    label: 'main',
    type: 'snippet',
    detail: 'entry point',
    apply: 'func main() {\n\t${}\n}',
    info: () => renderHelpDom({ kind: 'snippet', name: 'main', desc: 'Program entry function', example: 'func main() {\n    print("Hello!");\n}' }),
  },
  {
    label: 'if',
    type: 'snippet',
    detail: 'conditional',
    apply: 'if (${condition}) {\n\t${}\n}',
  },
  {
    label: 'ifelse',
    type: 'snippet',
    detail: 'if / else',
    apply: 'if (${condition}) {\n\t${}\n} else {\n\t${}\n}',
  },
  {
    label: 'while',
    type: 'snippet',
    detail: 'loop',
    apply: 'while (${condition}) {\n\t${}\n}',
  },
  {
    label: 'for',
    type: 'snippet',
    detail: 'for…of',
    apply: 'for (let ${item} of ${items}) {\n\t${}\n}',
  },
  {
    label: 'switch',
    type: 'snippet',
    detail: 'switch',
    apply: 'switch (${expr}) {\n\tcase ${1}:\n\t\t${}\n\t\tbreak;\n\tdefault:\n\t\t${}\n}',
  },
  {
    label: 'struct',
    type: 'snippet',
    detail: 'struct',
    apply: 'struct ${Name} {\n\t${fields}\n}',
  },
  {
    label: 'enum',
    type: 'snippet',
    detail: 'enum',
    apply: 'enum ${Name} {\n\t${Member1}, ${Member2}\n}',
  },
  {
    label: 'test',
    type: 'snippet',
    detail: 'test block',
    apply: 'test "${name}" {\n\texpect(${condition});\n}',
  },
  {
    label: 'import',
    type: 'snippet',
    detail: 'import module',
    apply: 'let ${alias} = import "${module}";',
  },
  {
    label: 'include',
    type: 'snippet',
    detail: '#include file',
    apply: '#include "${path}"',
  },
  {
    label: 'print',
    type: 'snippet',
    detail: 'print line',
    apply: 'print("${message}");',
  },
  {
    label: 'use',
    type: 'snippet',
    detail: 'module import',
    apply: 'use ${module};',
  },
  {
    label: 'useraylib',
    type: 'snippet',
    detail: 'full Raylib',
    apply: 'use raylib;',
  },
  {
    label: 'game',
    type: 'snippet',
    detail: 'game loop',
    apply:
      'use raylib;\nuse koda.game;\n\nfunc main() {\n\tgame.open(800, 600, "My Game");\n\tdefer game.close();\n\tgame.fps(60);\n\twhile (game.running()) {\n\t\tlet dt = game.delta();\n\t\tgame.begin();\n\t\tgame.clear(colors.dark);\n\t\t${}\n\t\tgame.end();\n\t}\n}',
  },
  {
    label: 'defer',
    type: 'snippet',
    detail: 'defer cleanup',
    apply: 'defer ${fn}();',
  },
]
