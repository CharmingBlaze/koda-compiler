Koda Language - Complete Specification
The Full Master Plan

TABLE OF CONTENTS

Language Overview
Lexical Structure
Type System
Variables & Scope
Functions
Control Flow
Data Structures
Module System
FFI & C Interop
Standard Library
Error Handling
Memory Management
Compiler Architecture
Runtime Systems
Tooling Ecosystem
Complete Grammar
Implementation Roadmap
Implementation Status (this repository)


Language Overview
Design Goals

Simple - Learnable in a day
Powerful - Access to entire C/C++ ecosystem
Fast - Native compilation (LLVM IR + llc + Clang + `runtime/libkoda_runtime.a`)
Practical - Solves real problems
Fun - Joy to write

Key Features

Dynamic typing with optional type hints
JavaScript-like syntax
First-class functions
Closures and lexical scoping
Object literals and arrays
Module system with #include
Automatic C/C++ library wrapping
Garbage collection
Dual compilation: bytecode VM or native binary


Lexical Structure
Character Set

Source files: UTF-8 encoding
Identifiers: Unicode letters, digits, underscore
Strings: UTF-8 text

Comments
javascript// Single-line comment

/* 
   Multi-line comment
   Can span multiple lines
*/

/// Documentation comment (for kodadoc tool)
/// @param x The input value
/// @returns The computed result
func calculate(x) { ... }
Keywords (Reserved)
// Control flow
if          else        for         while       
break       continue    return      switch      case
default

// Declarations
let         func        

// Literals
true        false       null

// Future reserved
try         catch       finally     throw
async       await       yield       
class       struct      enum        interface
import      export      module      namespace
const       static      public      private
Operators
javascript// Arithmetic
+  -  *  /  %  **        // power operator

// Comparison
==  !=  <  <=  >  >=
==  !=                   // equality (=== / !== are not supported; koda fmt rewrites them)

// Logical
&&  ||  !

// Bitwise
&  |  ^  ~  <<  >>  >>>

// Assignment
=  +=  -=  *=  /=  %=  **=
&=  |=  ^=  <<=  >>=

// Unary
++  --  +  -  !  ~

// Ternary
? :

// Member access
.  []

// Other
,  ;  :  ?  
(  )  {  }  [  ]
Literals
javascript// Numbers
42                  // integer
3.14                // float
1.5e-10             // scientific
0xFF                // hex
0o77                // octal
0b1010              // binary
1_000_000           // underscores for readability

// Strings
"hello"             // double quotes
'world'             // single quotes
`template ${x}`     // template strings

// Escape sequences
"\n"  "\t"  "\r"  "\\"  "\""  "\'"
"\x41"              // hex escape
"\u0041"            // unicode escape

// Raw strings (no escapes)
r"C:\path\to\file"

// Multi-line strings
"""
This is a
multi-line string
"""

// Booleans
true
false

// Null
null

Type System
Core Types
javascript// Number (IEEE 754 double precision)
let x = 42;
let y = 3.14;

// String (UTF-8)
let name = "Jesse";

// Boolean
let alive = true;

// Null
let nothing = null;

// Function
let greet = func(name) { print("Hello", name); };

// Array (dynamic, heterogeneous)
let items = [1, "two", true, {x: 10}];

// Object (hash map)
let player = {
    name: "Hero",
    health: 100
};
Type Checking
javascript// Runtime type checking
type(42)              // "number"
type("hello")         // "string"
type(true)            // "bool"
type(null)            // "null"
type([1, 2])          // "array"
type({x: 10})         // "object"
type(func() {})       // "function"

// Type predicates
is_number(x)
is_string(x)
is_bool(x)
is_null(x)
is_array(x)
is_object(x)
is_function(x)
Type Coercion
javascript// Implicit coercion in operations
"5" + 3              // 8 (string to number)
5 + 3                // 8
"5" + "3"            // "53" (concatenation)
true + 1             // 2 (bool: true=1, false=0)

// Explicit conversion
number("42")         // 42
string(42)           // "42"
bool(0)              // false
bool(1)              // true

// Truthiness (like JavaScript)
// Falsy: false, 0, "", null
// Truthy: everything else
if (0)         // false
if ("")        // false
if (null)      // false
if ([])        // true
if ({})        // true
Optional Type Hints (Future)
javascript// For documentation and optional checking
func add(x: number, y: number): number {
    return x + y;
}

let name: string = "Jesse";
let items: array = [];

// Type hints are optional, compiler can ignore them
// But tools like LSP use them for autocomplete

Variables & Scope
Variable Declaration
javascript// Standard declaration
let x = 10;
let name = "Jesse";
let alive = true;

// Multiple declarations
let a = 1, b = 2, c = 3;

// Uninitialized (defaults to null)
let x;               // x = null
Scope Rules
javascript// Global scope
let global_var = 100;

func outer() {
    // Function scope
    let outer_var = 200;
    
    func inner() {
        // Inner function scope (closure)
        let inner_var = 300;
        
        // Can access all outer scopes
        print(global_var);   // 100
        print(outer_var);    // 200
        print(inner_var);    // 300
    }
    
    inner();
}

// Block scope
if (true) {
    let block_var = 50;
    print(block_var);        // 50
}
// print(block_var);         // Error: undefined

// Loop scope
for (let i = 0; i < 10; i++) {
    let loop_var = i * 2;
}
// print(i);                 // Error: undefined
Shadowing
javascriptlet x = 10;              // global

func test() {
    let x = 20;          // shadows global
    print(x);            // 20
    
    if (true) {
        let x = 30;      // shadows function-level
        print(x);        // 30
    }
    
    print(x);            // 20
}

print(x);                // 10
Global Namespace
javascript// All top-level declarations are global
let health = 100;

func damage(amount) {
    health -= amount;    // modifies global
}

// Access from anywhere
#include "other.koda"
// other.koda can see 'health' and 'damage'

Functions
Function Declaration
javascript// Basic function
func greet(name) {
    print("Hello", name);
}

// Function with return
func add(a, b) {
    return a + b;
}

// Function with default parameters
func greet(name, greeting = "Hello") {
    print(greeting, name);
}

// No parameters
func get_random() {
    return random();
}

// Multiple returns (returns array)
func minmax(a, b) {
    if (a < b) return [a, b];
    return [b, a];
}
Function Expressions
javascript// Anonymous function
let greet = func(name) {
    print("Hello", name);
};

// Immediately invoked
let result = func(x) {
    return x * 2;
}(10);               // 20
Arrow Functions (Syntactic Sugar)
javascript// Short form
let square = (x) => x * x;

// With block
let add = (a, b) => {
    return a + b;
};

// No parameters
let greet = () => print("Hello");

// Single parameter (no parens)
let double = x => x * 2;
First-Class Functions
javascript// Functions are values
func add(a, b) { return a + b; }
func sub(a, b) { return a - b; }

let operations = [add, sub];

print(operations[0](10, 5));    // 15
print(operations[1](10, 5));    // 5

// Pass as arguments
func apply(f, x, y) {
    return f(x, y);
}

print(apply(add, 3, 4));        // 7
Closures
javascriptfunc make_counter() {
    let count = 0;
    
    return func() {
        count += 1;
        return count;
    };
}

let counter = make_counter();
print(counter());    // 1
print(counter());    // 2
print(counter());    // 3

// Each counter has its own state
let counter2 = make_counter();
print(counter2());   // 1
Variadic Functions
javascript// Rest parameters
func sum(...numbers) {
    let total = 0;
    for (let i = 0; i < numbers.length(); i++) {
        total += numbers[i];
    }
    return total;
}

print(sum(1, 2, 3));           // 6
print(sum(1, 2, 3, 4, 5));     // 15

// Mix regular and rest
func greet(greeting, ...names) {
    for (let i = 0; i < names.length(); i++) {
        print(greeting, names[i]);
    }
}

greet("Hello", "Alice", "Bob", "Charlie");
Named Parameters (Object Destructuring)
javascriptfunc create_player(options) {
    let name = options.name or "Player";
    let health = options.health or 100;
    let x = options.x or 0;
    let y = options.y or 0;
    
    return {
        name: name,
        health: health,
        x: x,
        y: y
    };
}

let player = create_player({
    name: "Hero",
    health: 150,
    x: 10
});
Recursion
javascriptfunc factorial(n) {
    if (n <= 1) return 1;
    return n * factorial(n - 1);
}

print(factorial(5));    // 120

// Tail call optimization (future)
func factorial_tail(n, acc = 1) {
    if (n <= 1) return acc;
    return factorial_tail(n - 1, n * acc);
}

Control Flow
If/Else
javascript// Basic if
if (health < 50) {
    print("Low health!");
}

// If-else
if (health > 80) {
    print("Healthy");
} else {
    print("Damaged");
}

// If-else-if chain
if (score > 1000) {
    print("S Rank");
} else if (score > 500) {
    print("A Rank");
} else if (score > 100) {
    print("B Rank");
} else {
    print("C Rank");
}

// Ternary operator
let status = health > 50 ? "alive" : "critical";

// Expression form (everything returns value)
let max = if (a > b) a else b;
Switch Statement
javascriptlet color = "red";

switch (color) {
    case "red":
        print("Fire!");
        break;
    
    case "blue":
        print("Water!");
        break;
    
    case "green":
        print("Earth!");
        break;
    
    default:
        print("Unknown color");
}

// Fall-through (no break)
switch (x) {
    case 1:
    case 2:
    case 3:
        print("Small");
        break;
    
    case 4:
    case 5:
        print("Medium");
        break;
}

// Switch as expression
let message = switch (score) {
    case 100: "Perfect!";
    case 90..99: "Excellent!";
    default: "Good try!";
};
For Loop
javascript// C-style for
for (let i = 0; i < 10; i++) {
    print(i);
}

// Multiple variables
for (let i = 0, j = 10; i < j; i++, j--) {
    print(i, j);
}

// Infinite loop (omit condition)
for (;;) {
    if (should_exit) break;
}

// For-in (array iteration)
let items = ["sword", "shield", "potion"];

for (let item in items) {
    print(item);
}

// For-in with index
for (let i, item in items) {
    print(i, ":", item);
}

// For-in (object iteration)
let player = {name: "Hero", health: 100, level: 5};

for (let key in player) {
    print(key, "=", player[key]);
}

// For-in with key and value
for (let key, value in player) {
    print(key, "=", value);
}
While Loop
javascript// Basic while
let x = 0;
while (x < 10) {
    print(x);
    x++;
}

// Condition-only
while (is_running()) {
    update();
}

// Infinite loop
while (true) {
    if (should_exit) break;
}
Do-While Loop
javascript// Execute at least once
let x = 0;
do {
    print(x);
    x++;
} while (x < 10);

// Input validation example
let input;
do {
    input = read_line("Enter positive number: ");
} while (number(input) <= 0);
Break & Continue
javascript// Break exits loop
for (let i = 0; i < 100; i++) {
    if (i == 50) break;
    print(i);
}

// Continue skips iteration
for (let i = 0; i < 10; i++) {
    if (i % 2 == 0) continue;
    print(i);    // Only odd numbers
}

// Labeled break (future)
outer: for (let i = 0; i < 10; i++) {
    for (let j = 0; j < 10; j++) {
        if (i * j > 50) break outer;
        print(i, j);
    }
}
Match Expression (Future Pattern Matching)
javascript// Enhanced switch with patterns
let result = match (value) {
    0 => "zero",
    1..10 => "small",
    11..100 => "medium",
    _ if value > 100 => "large",
    _ => "unknown"
};

// Type matching
let result = match (value) {
    x: number => x * 2,
    s: string => s.upper(),
    a: array => a.length(),
    _ => null
};

Data Structures
Arrays
javascript// Array literals
let empty = [];
let numbers = [1, 2, 3, 4, 5];
let mixed = [1, "two", true, {x: 10}];

// Nested arrays
let matrix = [
    [1, 2, 3],
    [4, 5, 6],
    [7, 8, 9]
];

// Array access
print(numbers[0]);        // 1
print(matrix[1][2]);      // 6

// Negative indexing (from end)
print(numbers[-1]);       // 5 (last element)
print(numbers[-2]);       // 4 (second to last)

// Array slicing
let slice = numbers[1:4];      // [2, 3, 4]
let slice2 = numbers[:3];      // [1, 2, 3]
let slice3 = numbers[2:];      // [3, 4, 5]
let slice4 = numbers[-2:];     // [4, 5]

// Array methods
numbers.push(6);              // Add to end
let last = numbers.pop();     // Remove from end
numbers.insert(0, 0);         // Insert at index
numbers.remove(2);            // Remove at index
let len = numbers.length();   // Get length
numbers.clear();              // Remove all elements

// Array iteration
for (let num in numbers) {
    print(num);
}

// Array methods (functional)
let doubled = numbers.map((x) => x * 2);
let evens = numbers.filter((x) => x % 2 == 0);
let sum = numbers.reduce((acc, x) => acc + x, 0);

// Array comprehension (future)
let squares = [x * x for x in numbers if x % 2 == 0];
Objects
javascript// Object literals
let player = {
    name: "Hero",
    health: 100,
    position: {x: 10, y: 20}
};

// Dot notation
print(player.name);           // "Hero"
player.health = 90;

// Bracket notation
print(player["name"]);        // "Hero"
let key = "health";
print(player[key]);           // 90

// Computed keys
let prop = "score";
let obj = {
    name: "Test",
    [prop]: 100               // score: 100
};

// Shorthand properties
let name = "Hero";
let health = 100;
let player = {name, health};  // {name: "Hero", health: 100}

// Methods
let player = {
    name: "Hero",
    health: 100,
    
    damage: func(amount) {
        this.health -= amount;
    },
    
    // Method shorthand
    heal(amount) {
        this.health += amount;
    }
};

player.damage(10);
player.heal(5);

// Dynamic properties
player.level = 5;             // Add new property
delete player.level;          // Remove property

// Object iteration
for (let key in player) {
    print(key, ":", player[key]);
}

// Object methods
let keys = player.keys();         // ["name", "health"]
let values = player.values();     // ["Hero", 100]
let entries = player.entries();   // [["name", "Hero"], ["health", 100]]

// Object spread (future)
let newPlayer = {...player, level: 5};

// Object destructuring (future)
let {name, health} = player;
Strings
javascript// String operations
let s = "hello world";

s.length()              // 11
s.upper()               // "HELLO WORLD"
s.lower()               // "hello world"
s.trim()                // Remove whitespace
s.split(" ")            // ["hello", "world"]
s.replace("world", "koda")  // "hello koda"
s.starts_with("hello")  // true
s.ends_with("world")    // true
s.contains("lo")        // true
s.index_of("world")     // 6
s.substring(0, 5)       // "hello"
s.repeat(3)             // "hello worldhello worldhello world"

// String interpolation
let name = "Jesse";
let age = 25;
let msg = `Hello ${name}, you are ${age} years old`;

// Multi-line strings
let text = """
This is a
multi-line
string
""";

// Raw strings (no escape sequences)
let path = r"C:\Users\Jesse\file.txt";

// Character access
print(s[0]);            // "h"
print(s[-1]);           // "d"

// String slicing
print(s[0:5]);          // "hello"

// String iteration
for (let char in s) {
    print(char);
}
Maps (Future Built-in)
javascript// Explicit map type (hash map with any key type)
let map = Map();

map.set("key", "value");
map.set(42, "number key");
map.set(obj, "object key");

print(map.get("key"));       // "value"
print(map.has("key"));       // true
map.delete("key");
print(map.size());           // 2

// Iteration
for (let key, value in map) {
    print(key, "=>", value);
}
Sets (Future Built-in)
javascript// Unordered collection of unique values
let set = Set();

set.add(1);
set.add(2);
set.add(1);              // Duplicate ignored

print(set.has(1));       // true
print(set.size());       // 2

set.remove(1);

// Iteration
for (let value in set) {
    print(value);
}

// Set operations
let a = Set([1, 2, 3]);
let b = Set([2, 3, 4]);

let union = a.union(b);         // [1, 2, 3, 4]
let intersection = a.intersect(b);  // [2, 3]
let difference = a.difference(b);   // [1]

Module System
File Inclusion
javascript// Include another Koda file
#include "player.koda"
#include "enemy.koda"
#include "utils.koda"

// Include from subdirectory
#include "libs/json.koda"
#include "game/entities/boss.koda"

// Include system library (from KODA_PATH)
#include <raylib>
#include <sqlite>
Module Structure
File: math_utils.koda
javascript// All top-level declarations are exported

func square(x) {
    return x * x;
}

func cube(x) {
    return x * x * x;
}

let PI = 3.14159;

// Private (not exported, use underscore prefix)
func _helper(x) {
    return x * 2;
}
File: main.koda
javascript#include "math_utils.koda"

print(math_utils.square(5));     // 25
print(math_utils.PI);             // 3.14159
// print(math_utils._helper(5));  // Error: private
Namespace Control
javascript// Import with custom namespace
#include "math_utils.koda" as math

print(math.square(5));

// Import specific symbols into current scope
#include "math_utils.koda" {square, cube}

print(square(5));        // No prefix needed
print(cube(3));          // No prefix needed

// Import all into current scope (use sparingly)
#include "math_utils.koda" *

print(square(5));
print(PI);
Module Resolution
Search order:
1. Relative to current file
   #include "utils.koda"  → ./utils.koda

2. Relative with path
   #include "libs/json.koda"  → ./libs/json.koda

3. System libraries (KODA_PATH environment variable)
   #include <raylib>  → $KODA_PATH/raylib/raylib.koda
   
4. Standard library (built-in)
   #include <math>    → Built-in math module
Circular Dependencies
javascript// a.koda
#include "b.koda"
func a_func() { b.b_func(); }

// b.koda
#include "a.koda"
func b_func() { a.a_func(); }

// Error: Circular dependency detected: a.koda <-> b.koda
Solution: Introduce third file
javascript// shared.koda
func shared_func() { ... }

// a.koda
#include "shared.koda"
func a_func() { shared.shared_func(); }

// b.koda
#include "shared.koda"
func b_func() { shared.shared_func(); }
Module Initialization
javascript// player.koda
print("Loading player module...");

let player_count = 0;

func create_player(name) {
    player_count++;
    return {name: name, id: player_count};
}

// main.koda
#include "player.koda"  // Prints "Loading player module..."

let p1 = player.create_player("Alice");
let p2 = player.create_player("Bob");

FFI & C Interop
Calling C Functions
javascript// Manual FFI declaration
#native "libm"              // Link with libm

#ffi func sqrt(number) -> number;
#ffi func pow(number, number) -> number;
#ffi func sin(number) -> number;

// Now use them
let result = sqrt(16);      // 4.0
let power = pow(2, 8);      // 256.0
Type Mapping
Koda TypeC TypeNotesnumberdoubleIEEE 754 doublestringchar*UTF-8, null-terminatedboolint0 or 1nullNULLNull pointerarrayvoid*Opaque array handleobjectvoid*Opaque object handlefunctionvoid*Function pointer wrapperptrvoid*Raw pointer
Struct Mapping
c// C header: point.h
typedef struct {
    int x;
    int y;
} Point;

Point create_point(int x, int y);
void print_point(Point p);
javascript// Koda wrapper (auto-generated by kodawrap)
#native "point"

#ffi func create_point(number, number) -> object;
#ffi func print_point(object) -> void;

// Helper constructor
func Point(x, y) {
    return create_point(x, y);
}

// Usage
let p = Point(10, 20);
print_point(p);
Pointer Handling
javascript#native "mylib"

// Pass pointer to C
#ffi func process_data(ptr) -> void;

// Get pointer from C
#ffi func get_buffer() -> ptr;

// Usage
let buffer = get_buffer();
process_data(buffer);

// Free when done (if not auto-freed)
free(buffer);
Callbacks
javascript// C function that takes callback
#ffi func iterate(array, function) -> void;

// Koda callback
let my_callback = func(item) {
    print("Item:", item);
};

let items = [1, 2, 3, 4, 5];
iterate(items, my_callback);
Memory Management with C
javascript// Automatic cleanup with finalizers
#native "sqlite"

#ffi func sqlite_open(string) -> ptr [finalizer: sqlite_close];
#ffi func sqlite_close(ptr) -> void;

// Usage - no manual cleanup needed
func use_database() {
    let db = sqlite_open("data.db");
    // ... use db
    // Automatically closed when db goes out of scope
}

Standard Library
Core Module (Always Available)
javascript// I/O
print(...)              // Print to stdout
warn(...)               // Print to stderr
input(prompt)           // Read line from stdin
read_file(path)         // Read entire file as string
write_file(path, text)  // Write string to file

// Type checking
type(value)             // Returns type name as string
is_number(x)
is_string(x)
is_bool(x)
is_null(x)
is_array(x)
is_object(x)
is_function(x)

// Conversion
number(value)           // Convert to number
string(value)           // Convert to string
bool(value)             // Convert to bool
array(value)            // Convert to array

// Math (basic)
abs(x)
min(a, b)
max(a, b)
random()                // 0.0 to 1.0
random_int(min, max)    // Inclusive

// String
len(str)                // Length
split(str, delim)       // Split string
join(array, delim)      // Join array to string

// Array
len(array)              // Length

// Time
time()                  // Seconds since epoch
sleep(ms)               // Sleep milliseconds
Math Module
javascript#include <math>

// Constants
math.PI                 // 3.14159265359
math.E                  // 2.71828182846
math.INF                // Infinity
math.NAN                // Not a Number

// Trigonometry
math.sin(x)
math.cos(x)
math.tan(x)
math.asin(x)
math.acos(x)
math.atan(x)
math.atan2(y, x)

// Exponential & logarithmic
math.exp(x)
math.log(x)             // Natural log
math.log10(x)
math.log2(x)
math.pow(x, y)
math.sqrt(x)
math.cbrt(x)            // Cube root

// Rounding
math.floor(x)
math.ceil(x)
math.round(x)
math.trunc(x)

// Other
math.abs(x)
math.sign(x)            // -1, 0, or 1
math.min(a, b)
math.max(a, b)
math.clamp(x, min, max)
math.lerp(a, b, t)      // Linear interpolation
math.random()           // 0.0 to 1.0
math.random_range(min, max)
math.random_int(min, max)
String Module
javascript#include <string>

string.len(s)
string.upper(s)
string.lower(s)
string.trim(s)
string.trim_left(s)
string.trim_right(s)
string.split(s, delim)
string.join(array, delim)
string.replace(s, old, new)
string.replace_all(s, old, new)
string.starts_with(s, prefix)
string.ends_with(s, suffix)
string.contains(s, substr)
string.index_of(s, substr)
string.last_index_of(s, substr)
string.substring(s, start, end)
string.repeat(s, count)
string.reverse(s)
string.pad_left(s, width, char)
string.pad_right(s, width, char)

// Format strings
string.format("Hello {}, you are {} years old", name, age)
Array Module
javascript#include <array>

array.len(arr)
array.push(arr, item)
array.pop(arr)
array.insert(arr, index, item)
array.remove(arr, index)
array.clear(arr)
array.contains(arr, item)
array.index_of(arr, item)
array.last_index_of(arr, item)
array.reverse(arr)
array.sort(arr)
array.sort(arr, compare_func)
array.slice(arr, start, end)
array.concat(arr1, arr2)

// Functional
array.map(arr, func)
array.filter(arr, func)
array.reduce(arr, func, initial)
array.foreach(arr, func)
array.any(arr, predicate)
array.all(arr, predicate)
array.find(arr, predicate)
array.find_index(arr, predicate)
Object Module
javascript#include <object>

object.keys(obj)
object.values(obj)
object.entries(obj)
object.has(obj, key)
object.get(obj, key, default)
object.set(obj, key, value)
object.delete(obj, key)
object.clear(obj)
object.size(obj)
object.merge(obj1, obj2)
object.clone(obj)
File I/O Module
javascript#include <file>

// Reading
file.read(path)                  // Read entire file as string
file.read_bytes(path)            // Read as byte array
file.read_lines(path)            // Read as array of lines

// Writing
file.write(path, text)
file.write_bytes(path, bytes)
file.append(path, text)

// Checks
file.exists(path)
file.is_file(path)
file.is_dir(path)
file.size(path)

// Directory operations
file.list_dir(path)
file.create_dir(path)
file.remove(path)
file.rename(old, new)
file.copy(src, dst)
file.move(src, dst)

// Path operations
file.join(parts...)
file.dirname(path)
file.basename(path)
file.extension(path)
file.absolute(path)
JSON Module
javascript#include <json>

// Parse JSON string to object
let data = json.parse('{"name": "Jesse", "age": 25}');
print(data.name);        // "Jesse"

// Stringify object to JSON
let text = json.stringify({x: 10, y: 20});
print(text);             // {"x":10,"y":20}

// Pretty print
let pretty = json.stringify(data, indent: 2);

// Handle errors
let result = json.try_parse(text);
if (result.error) {
    print("Parse error:", result.error);
} else {
    print("Data:", result.value);
}
HTTP Module
javascript#include <http>

// Simple GET request
let response = http.get("https://api.example.com/data");
print(response.status);    // 200
print(response.body);      // Response text

// POST request with JSON
let response = http.post("https://api.example.com/users", {
    headers: {"Content-Type": "application/json"},
    body: json.stringify({name: "Jesse", age: 25})
});

// Full control
let response = http.request({
    method: "PUT",
    url: "https://api.example.com/user/123",
    headers: {
        "Authorization": "Bearer token123",
        "Content-Type": "application/json"
    },
    body: json.stringify(data),
    timeout: 5000
});

// Handle errors
if (response.error) {
    warn("HTTP error:", response.error);
} else {
    print("Success:", response.body);
}
OS Module
javascript#include <os>

// Environment
os.getenv("PATH")
os.setenv("MY_VAR", "value")

// Process
os.exit(code)
os.args()                // Command line arguments
os.exec("ls -la")        // Run command, return output
os.spawn("python", ["script.py"])

// Platform info
os.platform()            // "linux", "macos", "windows"
os.arch()                // "x64", "arm64"
os.hostname()
os.username()
os.homedir()
os.tmpdir()

// Current directory
os.getcwd()
os.chdir(path)
Time Module
javascript#include <time>

// Current time
time.now()               // Seconds since epoch
time.now_ms()            // Milliseconds since epoch

// Sleep
time.sleep(1000)         // Sleep 1 second (milliseconds)
time.sleep_s(1)          // Sleep 1 second

// Date/time
let dt = time.datetime();
print(dt.year);          // 2024
print(dt.month);         // 5
print(dt.day);           // 1
print(dt.hour);          // 14
print(dt.minute);        // 30
print(dt.second);        // 45

// Formatting
time.format(dt, "%Y-%m-%d %H:%M:%S")
time.parse("2024-05-01", "%Y-%m-%d")

// Timer
let timer = time.timer();
// ... do work
print("Elapsed:", timer.elapsed(), "ms");
Regex Module (Future)
javascript#include <regex>

let pattern = regex.compile(r"\d+");

// Test match
regex.test(pattern, "abc123");     // true

// Find match
let match = regex.match(pattern, "abc123def");
print(match[0]);                    // "123"

// Find all matches
let matches = regex.find_all(pattern, "12 abc 34 def 56");
// ["12", "34", "56"]

// Replace
regex.replace(pattern, "abc123def", "XXX");  // "abcXXXdef"

Error Handling
Runtime Errors
javascript// Division by zero
let x = 10 / 0;          // Error: division by zero

// Null access
let obj = null;
print(obj.name);         // Error: null access

// Array out of bounds
let arr = [1, 2, 3];
print(arr[10]);          // Error: index out of bounds

// Type error
let x = "hello";
x + 5;                   // Error: cannot add string and number
Try-Catch (Future)
javascripttry {
    let data = json.parse(text);
    process(data);
} catch (e) {
    warn("Error:", e.message);
    warn("Stack trace:", e.stack);
} finally {
    cleanup();
}

// Specific error types
try {
    file.read("missing.txt");
} catch (e: FileNotFound) {
    print("File not found");
} catch (e: PermissionError) {
    print("Permission denied");
} catch (e) {
    print("Other error:", e);
}
Error Propagation
javascript// Return error objects
func read_config(path) {
    if (!file.exists(path)) {
        return {error: "Config file not found"};
    }
    
    let text = file.read(path);
    let data = json.parse(text);
    
    if (data.error) {
        return data;  // Propagate error
    }
    
    return {value: data};
}

// Check for errors
let result = read_config("config.json");
if (result.error) {
    warn("Failed to read config:", result.error);
} else {
    let config = result.value;
    // ... use config
}
Assert
javascript// Runtime assertion
assert(x > 0, "x must be positive");
assert(items.length() > 0);

// Disabled in release builds
// koda build --release (asserts removed)

Memory Management
Garbage Collection
javascript// Automatic memory management
let player = {health: 100};
// ... use player
// player automatically freed when no longer referenced
GC Algorithm

Mark-and-Sweep collector
Runs when allocation threshold reached
Can be manually triggered: gc.collect()

Reference Counting for C Objects
javascript// C objects use reference counting
#include <raylib>

func test() {
    let texture = raylib.LoadTexture("player.png");
    // ... use texture
    // Automatically unloaded when function returns
}
Manual Memory Control (Advanced)
javascript// For performance-critical code
let buffer = alloc(1024);        // Allocate 1KB
// ... use buffer
free(buffer);                     // Manual free

// With statement (auto-free at end of block)
with (let buffer = alloc(1024)) {
    // ... use buffer
}  // Automatically freed

Compiler Architecture
Frontend Pipeline
Source Code (.koda)
    ↓
┌─────────────────────┐
│   Lexer             │ - Tokenization
│   (hand-written)    │ - Character → Tokens
└─────────────────────┘
    ↓
┌─────────────────────┐
│   Parser            │ - Syntax analysis
│   (recursive descent)│ - Tokens → AST
└─────────────────────┘
    ↓
┌─────────────────────┐
│   Semantic Analysis │ - Type checking
│                     │ - Name resolution
│                     │ - Scope analysis
└─────────────────────┘
    ↓
┌─────────────────────┐
│   IR Generator      │ - High-level IR
│                     │ - SSA form
└─────────────────────┘
    ↓
┌─────────────────────┐
│   Optimizer         │ - Constant folding
│                     │ - Dead code elimination
│                     │ - Inlining
└─────────────────────┘
    ↓
┌─────────────┬───────┐
│             │       │
Backend 1     Backend 2
Backend 1: Bytecode VM
Optimized IR
    ↓
┌─────────────────────┐
│ Bytecode Compiler   │ - IR → Bytecode
└─────────────────────┘
    ↓
┌─────────────────────┐
│ Bytecode VM         │ - Stack-based VM
│ (like Lua/CPython)  │ - Direct threading
└─────────────────────┘
    ↓
Execution (5-10x slower than C)
Backend 2: Native Compilation
Optimized IR
    ↓
┌─────────────────────┐
│ LLVM IR (llir)      │ - bundle → .ll
└─────────────────────┘
    ↓
┌─────────────────────┐
│ LLVM llc + Clang    │ - .ll → .o + libkoda_runtime.a → native
└─────────────────────┘
    ↓
Native Binary (AOT)

Runtime Systems
Bytecode VM
c// VM opcode set (30-40 opcodes)
typedef enum {
    OP_CONSTANT,        // Load constant
    OP_NULL,
    OP_TRUE,
    OP_FALSE,
    
    OP_ADD,
    OP_SUBTRACT,
    OP_MULTIPLY,
    OP_DIVIDE,
    OP_MODULO,
    OP_POWER,
    
    OP_NEGATE,
    OP_NOT,
    
    OP_EQUAL,
    OP_NOT_EQUAL,
    OP_LESS,
    OP_LESS_EQUAL,
    OP_GREATER,
    OP_GREATER_EQUAL,
    
    OP_GET_GLOBAL,
    OP_SET_GLOBAL,
    OP_GET_LOCAL,
    OP_SET_LOCAL,
    OP_GET_UPVALUE,
    OP_SET_UPVALUE,
    
    OP_GET_PROPERTY,
    OP_SET_PROPERTY,
    OP_GET_INDEX,
    OP_SET_INDEX,
    
    OP_ARRAY,
    OP_OBJECT,
    
    OP_JUMP,
    OP_JUMP_IF_FALSE,
    OP_LOOP,
    
    OP_CALL,
    OP_RETURN,
    
    OP_CLOSURE,
    OP_CLOSE_UPVALUE,
    
    OP_POP,
    OP_PRINT,
} OpCode;

// VM execution with direct threading
void vm_run() {
    static void* dispatch[] = {
        &&OP_CONSTANT,
        &&OP_ADD,
        // ... all opcodes
    };
    
    register uint8_t* ip = vm.ip;
    register Value* sp = vm.sp;
    
    DISPATCH();
    
OP_CONSTANT: {
    PUSH(CONSTANT);
    DISPATCH();
}

OP_ADD: {
    Value b = POP();
    Value a = POP();
    PUSH(value_add(a, b));
    DISPATCH();
}
    
    // ... all operations
}
Value Representation
c// NaN-boxing for compact values (8 bytes each)
typedef uint64_t Value;

#define SIGN_BIT     ((uint64_t)0x8000000000000000)
#define QNAN         ((uint64_t)0x7ffc000000000000)

#define TAG_NULL     1
#define TAG_FALSE    2
#define TAG_TRUE     3

#define IS_NUMBER(v) (((v) & QNAN) != QNAN)
#define IS_NULL(v)   ((v) == (QNAN | TAG_NULL))
#define IS_BOOL(v)   (((v) | 1) == (QNAN | TAG_TRUE))

#define AS_NUMBER(v) value_to_num(v)
#define AS_BOOL(v)   ((v) == (QNAN | TAG_TRUE))

#define NUMBER_VAL(n) num_to_value(n)
#define BOOL_VAL(b)   ((b) ? (QNAN | TAG_TRUE) : (QNAN | TAG_FALSE))
#define NULL_VAL      (QNAN | TAG_NULL)

// Heap objects (strings, arrays, objects, functions)
#define IS_OBJ(v)    (((v) & (QNAN | SIGN_BIT)) == (QNAN | SIGN_BIT))
#define AS_OBJ(v)    ((Obj*)(uintptr_t)((v) & ~(SIGN_BIT | QNAN)))
#define OBJ_VAL(obj) (Value)(SIGN_BIT | QNAN | (uintptr_t)(obj))
Garbage Collector
c// Mark-and-sweep GC
void gc_collect() {
    // Mark phase
    mark_roots();           // Mark globals, stack, upvalues
    mark_compiler_roots();  // Mark compiler temporaries
    trace_references();     // Trace all reachable objects
    
    // Sweep phase
    sweep_strings();        // Sweep interned strings
    sweep_objects();        // Sweep heap objects
    
    // Update stats
    vm.next_gc = vm.bytes_allocated * GC_HEAP_GROW_FACTOR;
}

// Tri-color marking
typedef enum {
    OBJ_WHITE,    // Not reached
    OBJ_GRAY,     // Reached but not traced
    OBJ_BLACK,    // Reached and traced
} ObjColor;

Tooling Ecosystem
Koda Compiler (koda)
bash# Run script in bytecode VM (fast iteration)
koda run script.koda

# Build native binary (production)
koda build script.koda -o program

# Build with optimizations
koda build script.koda -o program --release

# Watch mode (auto-reload on change)
koda watch script.koda

# REPL
koda

# Show bytecode
koda disasm script.koda

# Format code
koda fmt script.koda

# Check syntax without running
koda check script.koda
Kodawrap (kodawrap)
bash# Generate wrapper from C header
kodawrap header.h -o wrapper.koda

# With config file
kodawrap --config kodawrap.toml

# Multiple headers
kodawrap *.h -o library.koda

# With documentation
kodawrap --docs header.h -o wrapper.koda

# Generate markdown docs
kodawrap --markdown header.h
Package Manager (koda pkg)
bash# Install package
koda pkg install raylib
koda pkg install sqlite

# Search packages
koda pkg search graphics

# List installed
koda pkg list

# Update package
koda pkg update raylib

# Remove package
koda pkg remove raylib

# Publish package
koda pkg publish
Language Server (koda-lsp)
bash# Start LSP server
koda-lsp

# Features:
# - Autocomplete
# - Go to definition
# - Find references
# - Hover documentation
# - Error diagnostics
# - Code formatting
Debugger (koda-debug)
bash# Run with debugger
koda debug script.koda

# Commands:
# break file.koda:10    - Set breakpoint
# continue              - Continue execution
# step                  - Step to next line
# next                  - Step over function
# print var             - Print variable
# backtrace             - Show call stack
# quit                  - Exit debugger
Testing Framework
javascript// test.koda
#include <test>

test.describe("Math operations", func() {
    test.it("should add numbers", func() {
        let result = add(2, 3);
        test.assert_equal(result, 5);
    });
    
    test.it("should multiply numbers", func() {
        let result = multiply(4, 5);
        test.assert_equal(result, 20);
    });
});

test.run();
bash# Run tests
koda test test.koda

# Run all tests in directory
koda test tests/

# Watch mode
koda test --watch
Profiler
bash# Profile execution
koda profile script.koda

# Output:
# Function         Calls    Time (ms)    % Total
# main                 1      150.2      75.1%
# update             100       30.5      15.3%
# render             100       19.3       9.6%
Documentation Generator (kodadoc)
bash# Generate docs from source
kodadoc src/ -o docs/

# Generates HTML documentation
# with API reference

Complete Grammar
ebnfprogram          = declaration* EOF

declaration      = varDecl
                 | funcDecl
                 | includeDecl
                 | nativeDecl
                 | ffiDecl
                 | statement

varDecl          = "let" IDENTIFIER ("=" expression)? ("," IDENTIFIER ("=" expression)?)* ";"?

funcDecl         = "func" IDENTIFIER "(" parameters? ")" block

parameters       = IDENTIFIER ("=" expression)? ("," IDENTIFIER ("=" expression)?)*
                 | "..." IDENTIFIER

includeDecl      = "#include" (STRING | "<" IDENTIFIER ">") 
                   ("as" IDENTIFIER)? 
                   ("{" IDENTIFIER ("," IDENTIFIER)* "}")? 
                   ("*")?

nativeDecl       = "#native" STRING

ffiDecl          = "#ffi" "func" IDENTIFIER "(" ffiParams? ")" "->" ffiType ";"

ffiParams        = ffiType ("," ffiType)*
ffiType          = "number" | "string" | "bool" | "ptr" | "void" | "object" | "array" | "function"

statement        = exprStmt
                 | ifStmt
                 | switchStmt
                 | whileStmt
                 | doWhileStmt
                 | forStmt
                 | breakStmt
                 | continueStmt
                 | returnStmt
                 | block

exprStmt         = expression ";"?

ifStmt           = "if" "(" expression ")" statement ("else" statement)?

switchStmt       = "switch" "(" expression ")" "{" switchCase* "}"
switchCase       = "case" expression ":" statement* ("break" ";"?)?
                 | "default" ":" statement*

whileStmt        = "while" "(" expression ")" statement

doWhileStmt      = "do" statement "while" "(" expression ")" ";"?

forStmt          = "for" "(" (varDecl | exprStmt | ";") expression? ";" expression? ")" statement
                 | "for" "(" IDENTIFIER ("," IDENTIFIER)? "in" expression ")" statement

breakStmt        = "break" ";"?
continueStmt     = "continue" ";"?
returnStmt       = "return" expression? ";"?

block            = "{" declaration* "}"

expression       = assignment

assignment       = (call "." IDENTIFIER | call "[" expression "]")? 
                   ("=" | "+=" | "-=" | "*=" | "/=" | "%=" | "**=" | "&=" | "|=" | "^=" | "<<=" | ">>=") 
                   assignment
                 | ternary

ternary          = logicalOr ("?" expression ":" ternary)?

logicalOr        = logicalAnd ("||" logicalAnd)*

logicalAnd       = bitwiseOr ("&&" bitwiseOr)*

bitwiseOr        = bitwiseXor ("|" bitwiseXor)*

bitwiseXor       = bitwiseAnd ("^" bitwiseAnd)*

bitwiseAnd       = equality ("&" equality)*

equality         = comparison (("!=" | "==") comparison)*

comparison       = bitwiseShift ((">" | ">=" | "<" | "<=") bitwiseShift)*

bitwiseShift     = range (("<<" | ">>" | ">>>") range)*

range            = addition (".." addition)?

addition         = multiplication (("-" | "+") multiplication)*

multiplication   = exponentiation (("/" | "*" | "%") exponentiation)*

exponentiation   = unary ("**" unary)*

unary            = ("!" | "-" | "+" | "~" | "++" | "--") unary
                 | postfix

postfix          = call ("++" | "--")?

call             = primary ("(" arguments? ")" | "." IDENTIFIER | "[" expression "]")*

arguments        = expression ("," expression)*

primary          = "true" | "false" | "null"
                 | NUMBER | STRING
                 | IDENTIFIER
                 | "(" expression ")"
                 | arrayLiteral
                 | objectLiteral
                 | funcExpr
                 | arrowFunc
                 | templateString

arrayLiteral     = "[" (expression ("," expression)*)? "]"

objectLiteral    = "{" (objectPair ("," objectPair)*)? "}"
objectPair       = (IDENTIFIER | STRING | "[" expression "]") ":" expression
                 | IDENTIFIER  // shorthand

funcExpr         = "func" "(" parameters? ")" block

arrowFunc        = (IDENTIFIER | "(" parameters? ")") "=>" (expression | block)

templateString   = "`" (TEXT | "${" expression "}")* "`"

NUMBER           = [0-9]+ ("." [0-9]+)? ([eE] [+-]? [0-9]+)?
                 | "0x" [0-9a-fA-F]+
                 | "0o" [0-7]+
                 | "0b" [01]+

STRING           = '"' (CHAR | ESCAPE)* '"'
                 | "'" (CHAR | ESCAPE)* "'"
                 | '"""' .*? '"""'
                 | 'r"' .*? '"'

IDENTIFIER       = [a-zA-Z_] [a-zA-Z0-9_]*

COMMENT          = "//" .* NEWLINE
                 | "/*" .*? "*/"

Implementation Status (this repository)

This section describes the Go compiler/VM in this repo (`cmd/koda`, `internal/parser`, `internal/runtime`). It is the source of truth for what is implemented today. The **Implementation Roadmap** below remains a **target vision** for the language project as a whole.

**Distribution model:** Ship a **prebuilt `koda` executable** (and optionally **kodawrap** next to it). **`koda run` / `check` / `disasm`** need no C toolchain. **`koda build`** / **`koda bundle`** emit **LLVM IR** (`.ll` via [llir/llvm](https://github.com/llir/llvm)), write **`.KODA_build/`**, run **llc** to an object file, then drive **LLVM `clang`** to link **`runtime/libkoda_runtime.a`** and headers under **`runtime/src/`** — see [DISTRIBUTION_GUIDE.md](DISTRIBUTION_GUIDE.md), [HANDOFF.md](HANDOFF.md), **`KODA_CLANG`**, **`CC`**, **`KODA_USE_LLD`**. **`bundle`** writes a tidy folder + README for sharing with end users (no Go/Python on their side).

**Tracking:** Use **[list.md](list.md)** in this repository as a checkbox audit vs the rest of this document.

Implemented (high level)

- Lexer, parser, AST, bytecode compiler + stack VM with GC (`internal/parser`).
- CLI: `koda run`, `check`, `disasm`, `build`, **`bundle`**, **`help`**, **`version`**, **`wrap`** (forwards to **kodawrap** if present beside `koda`) — [cmd/koda/main.go](cmd/koda/main.go).
- **kodawrap** (sources under **`cmd/wrapgen`**, same Go module): `go build -o kodawrap ./cmd/wrapgen` generates readable `.koda` + `wrapper.c` + Markdown from C headers.
- **C interop today:** `// koda:extern` lines parsed in [internal/codegen.go](internal/codegen.go), not `#native` / `#ffi` tokens from the grammar below.
- Modules: `import "path"` expression form, `import("@scope/...")`, and `#include "path"` / `#include <name>` ([internal/parser/loader.go](internal/parser/loader.go)). Optional `#include` forms `as`, `{names}`, `*` from the grammar are **not** implemented.
- Core language: control flow (`if`, `while`, `for`, `for (let x in/of …)`, `do`/`while`, `switch` / `match` — **no fall-through** unless you use `fallthrough;`; LLVM backend), operators including bitwise/shift, ternary, objects/arrays, closures, template literals, raw strings, numeric radix forms, `...rest` / defaults (see tests and [KODA_PROGRAMMER_REFERENCE.md](KODA_PROGRAMMER_REFERENCE.md)).
- Builtins: `type`, conversions, `is_*` predicates, core math/time/file helpers per [internal/sema.go](internal/sema.go).

Release-safe additions implemented after release hardening

- `assert(condition[, message])` is available as a core VM builtin for runtime invariants and tests.
- `@math` includes `cbrt`, `inf`, `nan`, `random`, `randomRange`, and `randomInt` in addition to the existing trig/log/rounding helpers.
- `@io` includes safe metadata helpers: `isFile(path)`, `isDir(path)`, `size(path)`, and `list(path)`.
- `json.stringify(value, indent)` supports pretty JSON indentation; `json.try_parse(text)` returns `{error, value}` without throwing.
- `koda disasm <file.koda>` exists for bytecode inspection.
- Current native/LLVM backend release gate supports core/native-lowered language features. Imported stdlib module objects are primarily VM-supported until module-object lowering is expanded in LLVM.
Not implemented yet (defer)

- **`#native` / `#ffi`** as first-class syntax (use **wrapgen** + `// koda:extern` + `wrapper.c` instead).
- **Rich stdlib remaining work:** `@math`, `@json`, `@io`, `@str`, `@array`, `os`, and `path` exist with production-tested core coverage. Larger modules from the vision spec (HTTP, regex, full object/file/string/array APIs, package registry integration) remain deferred.
- **Tooling**: no REPL, `koda fmt`, `koda test`, LSP, package registry, debugger, or profiler in this repo.
- **Grammar gaps vs spec examples:** arrow-function expressions, `for (let a, b in x)` two-variable forms, `switch`/`if` as value expressions, object shorthand/computed keys, many array/string conveniences from the spec — see [list.md](list.md).

Implementation Roadmap (target vision)
Phase 1: Core Language (3 months)
Month 1: Frontend

- Lexer (tokenization)
- Parser (AST generation)
- Basic AST nodes
- Error reporting

Month 2: Bytecode VM

- Value representation
- Bytecode opcodes (30-40)
- Stack-based VM
- Basic garbage collector
- Standard library (core functions)

Month 3: Language Features

- Variables & scope
- Functions & closures
- Objects & arrays
- Control flow
- String operations

Deliverable: koda run command works, can run basic programs

Phase 2: Native Compilation (2 months)
Month 4: LLVM IR + runtime

- LLVM IR emission (llir) for user code
- Embedded C runtime ([internal/runtime/data](internal/runtime/data)) linked by **LLVM clang**
- Optional **LLD** via **`KODA_USE_LLD=1`** (`-fuse-ld=lld`); driver selection via **`KODA_CLANG`** / **`CC`**

Month 5: Optimization

- Constant folding
- Dead code elimination
- Function inlining
- Build system

Deliverable: koda build command works, produces native binaries

Phase 3: Module System (1 month)
Month 6: Includes & Packages

- #include directive (full selective forms)
- Module resolution
- Namespace management
- Package manager basics

Deliverable: Can split code into modules

Phase 4: C Interop — **partially delivered**: `// koda:extern` + **wrapgen** (`cmd/wrapgen`); `#native` / `#ffi` syntax still deferred. See Implementation Status above and [list.md](list.md).
Month 7: FFI Foundation

- #native directive
- #ffi declarations
- Type marshalling
- Callback support

Month 8: Kodawrap Tool

- libclang integration
- Header parsing
- Wrapper generation
- Documentation extraction

Deliverable: Can wrap C libraries automatically

Phase 5: Standard Library (1 month)
Month 9: Built-in Modules

- Math module
- String module
- Array module
- File I/O
- JSON
- HTTP
- Time

Deliverable: Rich standard library

Phase 6: Tooling (2 months)
Month 10: Developer Tools

- Code formatter
- Syntax checker
- REPL
- Testing framework

Month 11: IDE Support

- Language server (LSP)
- VS Code extension
- Syntax highlighting
- Autocomplete

Deliverable: Professional development experience

Phase 7: Polish & Launch (1 month)
Month 12: Final Push

- Documentation
- Tutorial
- Example projects
- Website
- Package registry
- Community setup

Deliverable: Public release!

Total Timeline

Minimum Viable Product: 6 months (Phases 1-3)
Feature Complete: 9 months (Phases 1-5)
Production Ready: 12 months (All phases)

Team Size:

Solo developer: 12 months
2 developers: 6-8 months
Small team (3-4): 4-6 months


Success Metrics
Technical

- Bytecode VM: 5-10x faster than Python
- Native compilation: Within 2x of hand-written C
- Compile time: <100ms for small projects
- Binary size: <5MB for compiler
- Memory efficient: GC overhead <20%

User Experience

- Learn basics in <1 day
- Zero-config setup
- Fast edit-run cycle (<1s)
- Access to 100+ C libraries via Kodawrap
- Active community


Competitive Positioning
FeatureKodaLuaPythonJavaScriptEase of Learning⭐⭐⭐⭐⭐⭐⭐⭐⭐⭐⭐⭐⭐⭐⭐⭐⭐⭐⭐Speed (VM)⭐⭐⭐⭐⭐⭐⭐⭐⭐⭐⭐⭐⭐⭐⭐⭐Speed (Native)⭐⭐⭐⭐⭐⭐⭐⭐⭐⭐C Interop⭐⭐⭐⭐⭐⭐⭐⭐⭐⭐⭐⭐⭐Auto C Wrapping⭐⭐⭐⭐⭐⭐⭐⭐⭐Zero Config⭐⭐⭐⭐⭐⭐⭐⭐⭐⭐⭐⭐⭐⭐Modern Syntax⭐⭐⭐⭐⭐⭐⭐⭐⭐⭐⭐⭐⭐⭐⭐⭐⭐
Koda's Unique Value: JavaScript-easy syntax + C-level speed + effortless C library access

Example Projects
1. Game Engine
javascript#include <raylib>

let player = {x: 400, y: 300, speed: 5};

func update() { /* game logic */ }
func draw() { /* rendering */ }
func main() { /* game loop */ }
2. CLI Tool
javascript#include <os>
#include <file>

let args = os.args();
// Process files, do work
3. Web Server
javascript#include <http>

func handle_request(req) {
    return {
        status: 200,
        body: "Hello from Koda!"
    };
}

http.serve(":8080", handle_request);
4. Data Processing
javascript#include <file>
#include <json>

let data = json.parse(file.read("data.json"));
let processed = data.map((item) => transform(item));
file.write("output.json", json.stringify(processed));

Community & Ecosystem
Website: github.com/CharmingBlaze/koda-compiler

Documentation
Tutorial
Playground (WASM)
Package registry
Blog

GitHub: github.com/CharmingBlaze/koda-compiler

Compiler source
Standard library
Community packages
Issue tracker

Discord/Forum

Community support
Show & tell
Package announcements

Package Registry

Central repository
Version management
Documentation hosting
Download statistics


Marketing Tagline

Koda: Fortune Favors the Bold
Draw your ideas into code.
JavaScript-easy syntax.
C-powered performance.
Effortless C library access.
One command to rule them all: koda

bashcurl github.com/CharmingBlaze/koda-compiler/install | sh
echo 'print("Hello, World!")' > hello.koda
koda run hello.koda
Write it. Draw it. Ship it. 🎴