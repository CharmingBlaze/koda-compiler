

Koda Language - Complete Command Reference for Programmers
Moving to Go + LLVM backend
________________________________________
1. Variable Declarations
Syntax
javascript
let identifier = expression;
let identifier;  // Defaults to null
Examples
javascript
let x = 42;
let name = "Jesse";
let alive = true;
let empty;  // null

// Multiple declarations
let a = 1;
let b = 2;
let c = 3;
Rules
•	Block-scoped (like JavaScript let)
•	Mutable (can be reassigned)
•	Must be declared before use
•	No hoisting
________________________________________
2. Data Types
Primitives
javascript
// Number (IEEE 754 double)
let integer = 42;
let float = 3.14159;
let negative = -100;
let hex = 0xFF;  // 255

// String (UTF-8)
let str1 = "Hello";
let str2 = "World";
let escaped = "Line 1\nLine 2\tTabbed";

// Boolean
let yes = true;
let no = false;

// Null
let nothing = null;
Composite Types
javascript
// Array
let numbers = [1, 2, 3, 4, 5];
let mixed = [1, "hello", true, null];
let nested = [[1, 2], [3, 4]];
let empty = [];

// Object
let person = {
    name: "Alice",
    age: 30,
    active: true
};

let point = {x: 10, y: 20};
let empty_obj = {};
________________________________________
3. Operators
Arithmetic
javascript
let sum = 10 + 5;        // 15
let diff = 10 - 5;       // 5
let product = 10 * 5;    // 50
let quotient = 10 / 5;   // 2
let remainder = 10 % 3;  // 1
let negation = -x;       // Unary minus
Comparison
javascript
10 == 10   // true  (equality)
10 != 5    // true  (inequality)
10 > 5     // true  (greater than)
10 >= 10   // true  (greater or equal)
10 < 20    // true  (less than)
10 <= 10   // true  (less or equal)
Logical
javascript
true && false   // false (AND)
true || false   // true  (OR)
!true          // false (NOT)
Assignment
javascript
x = 10      // Simple assignment
x += 5      // x = x + 5
x -= 3      // x = x - 3
x *= 2      // x = x * 2
x /= 4      // x = x / 4
Increment/Decrement
javascript
x++   // Post-increment
x--   // Post-decrement
++x   // Pre-increment
--x   // Pre-decrement
________________________________________
4. Control Flow
If Statement
javascript
if (condition) {
    // statements
}

if (condition) {
    // then branch
} else {
    // else branch
}

if (condition1) {
    // branch 1
} else if (condition2) {
    // branch 2
} else {
    // default branch
}
Examples:
javascript
if (x > 10) {
    print("Big");
}

if (health <= 0) {
    print("Dead");
} else {
    print("Alive");
}

if (score >= 90) {
    print("A");
} else if (score >= 80) {
    print("B");
} else if (score >= 70) {
    print("C");
} else {
    print("F");
}
________________________________________
Switch Statement
javascript
switch (expression) {
    case value1:
        // statements
        break;
    
    case value2:
        // statements
        break;
    
    default:
        // statements
}
Examples:
javascript
switch (day) {
    case 1:
        print("Monday");
        break;
    
    case 2:
        print("Tuesday");
        break;
    
    default:
        print("Other day");
}

// Fall-through (multiple cases)
switch (grade) {
    case "A":
    case "B":
        print("Pass");
        break;
    
    case "F":
        print("Fail");
        break;
}
________________________________________
While Loop
javascript
while (condition) {
    // statements
}
Examples:
javascript
let i = 0;
while (i < 10) {
    print(i);
    i += 1;
}

// Infinite loop
while (true) {
    if (shouldExit) break;
    update();
}
________________________________________
For Loop
javascript
for (initializer; condition; increment) {
    // statements
}
Examples:
javascript
for (let i = 0; i < 10; i += 1) {
    print(i);
}

for (let i = 10; i > 0; i -= 1) {
    print(i);
}

// Empty sections allowed
let i = 0;
for (; i < 10;) {
    print(i);
    i += 1;
}
________________________________________
For-Of Loop
javascript
// Iterate over array elements
for (let item of array) {
    // statements
}

// With index
for (let index, item of array) {
    // statements
}

// Iterate over object keys
for (let key of object) {
    // statements
}

// With key and value
for (let key, value of object) {
    // statements
}
Examples:
javascript
let items = ["sword", "shield", "potion"];

// Just values
for (let item of items) {
    print(item);
}

// Index and value
for (let i, item of items) {
    print(i, ":", item);
}

// Object iteration
let player = {name: "Hero", level: 5};

// Just keys
for (let key of player) {
    print(key);
}

// Keys and values
for (let key, value of player) {
    print(key, "=", value);
}

// String iteration
for (let char of "hello") {
    print(char);
}
________________________________________
Break and Continue
javascript
// Break - exit loop immediately
for (let i = 0; i < 10; i += 1) {
    if (i == 5) break;
    print(i);  // 0, 1, 2, 3, 4
}

// Continue - skip to next iteration
for (let i = 0; i < 10; i += 1) {
    if (i % 2 == 0) continue;
    print(i);  // 1, 3, 5, 7, 9
}
________________________________________
5. Functions
Function Declaration
javascript
func name(parameters) {
    // statements
    return value;
}
Examples:
javascript
// No parameters
func greet() {
    print("Hello!");
}

// With parameters
func add(a, b) {
    return a + b;
}

// Multiple statements
func factorial(n) {
    if (n <= 1) {
        return 1;
    }
    return n * factorial(n - 1);
}
________________________________________
Function Calls
javascript
greet();
let result = add(10, 20);
let fact = factorial(5);
________________________________________
Function Expressions
javascript
// Assign to variable
let greet = func(name) {
    print("Hello, " + name);
};

// Pass as argument
array.map(func(x) {
    return x * 2;
});

// Return from function
func makeAdder(x) {
    return func(y) {
        return x + y;
    };
}

let add5 = makeAdder(5);
print(add5(10));  // 15
________________________________________
Return Statement
javascript
return;           // Return null
return value;     // Return value
return a + b;     // Return expression
________________________________________
6. Closures
javascript
func makeCounter() {
    let count = 0;
    
    return func() {
        count += 1;
        return count;
    };
}

let counter1 = makeCounter();
let counter2 = makeCounter();

print(counter1());  // 1
print(counter1());  // 2
print(counter2());  // 1  (separate closure)
________________________________________
7. Objects
Object Literals
javascript
let obj = {
    property: value,
    key: "value",
    nested: {
        inner: 42
    }
};
Property Access
javascript
// Dot notation
obj.property
obj.nested.inner

// Assignment
obj.property = newValue;
obj.newProp = 123;
Methods
javascript
let obj = {
    value: 42,
    
    // Method shorthand
    getValue() {
        return this.value;
    },
    
    setValue(v) {
        this.value = v;
    }
};

obj.getValue();
obj.setValue(100);
________________________________________
The this Keyword
javascript
let player = {
    name: "Hero",
    health: 100,
    
    damage(amount) {
        this.health -= amount;  // 'this' refers to player
        if (this.health <= 0) {
            this.die();
        }
    },
    
    die() {
        print(this.name + " died!");
    }
};

player.damage(150);
Rules:
•	Method call: obj.method() → this = obj
•	Regular call: func() → this = null
________________________________________
8. Arrays
Array Literals
javascript
let empty = [];
let numbers = [1, 2, 3, 4, 5];
let mixed = [1, "hello", true, null];
Array Access
javascript
// Indexing (0-based)
let first = array[0];
let second = array[1];

// Assignment
array[0] = newValue;
array[5] = 42;
Array Methods
javascript
// Length
let count = array.length();

// Push (add to end)
array.push(item);

// Pop (remove from end)
let last = array.pop();
________________________________________
9. Modules
Include Directive
javascript
#include "path/to/file.koda"
Examples:
javascript
// main.koda
#include "utils.koda"
#include "player.koda"

let p = createPlayer("Hero");
utils.greet(p.name);

// utils.koda
func greet(name) {
    print("Hello, " + name);
}

func square(x) {
    return x * x;
}
Rules:
•	Files loaded once (no duplicates)
•	Relative paths supported
•	All declarations become available
•	Circular includes handled
________________________________________
10. Comments
javascript
// Single-line comment

/* 
   Multi-line
   comment
*/

let x = 42;  // Inline comment
________________________________________
11. Standard Library
I/O
javascript
print(value);              // Print to stdout
print(a, b, c);            // Multiple values
print("Result:", result);  // Mixed types
Type Checking
javascript
type(value)       // Returns: "number", "string", "bool", 
                  //          "null", "array", "object", "function"

// Examples:
type(42)          // "number"
type("hello")     // "string"
type(true)        // "bool"
type(null)        // "null"
type([1, 2])      // "array"
type({x: 10})     // "object"
type(greet)       // "function"
Type Conversion
javascript
number(value)     // Convert to number
string(value)     // Convert to string

// Examples:
number("42")      // 42
number("3.14")    // 3.14
string(42)        // "42"
string(true)      // "true"
Math Functions
javascript
abs(x)            // Absolute value
sqrt(x)           // Square root
random()          // Random [0, 1)

// Examples:
abs(-10)          // 10
sqrt(16)          // 4
random()          // 0.7834...
Array/String Functions
javascript
len(value)        // Length of array or string

// Examples:
len([1, 2, 3])    // 3
len("hello")      // 5
Time Functions
javascript
time()            // Seconds since epoch
sleep(ms)         // Sleep milliseconds

// Examples:
let now = time();           // 1234567890
sleep(1000);                // Sleep 1 second
________________________________________
12. Reserved Keywords
let        - Variable declaration
func       - Function declaration
if         - Conditional
else       - Conditional alternative
for        - Loop
while      - Loop
switch     - Multi-way branch
case       - Switch case
default    - Switch default
break      - Exit loop/switch
continue   - Skip iteration
return     - Return from function
true       - Boolean literal
false      - Boolean literal
null       - Null literal
this       - Object self-reference
of         - For-of loop
________________________________________
13. Operators Precedence (Highest to Lowest)
1.  ()  []  .              (Grouping, indexing, property)
2.  ++  --  !  -  +        (Unary)
3.  *  /  %                (Multiplication, division, modulo)
4.  +  -                   (Addition, subtraction)
5.  <  <=  >  >=           (Comparison)
6.  ==  !=                 (Equality)
7.  &&                     (Logical AND)
8.  ||                     (Logical OR)
9.  =  +=  -=  *=  /=      (Assignment)
________________________________________
14. Truthy/Falsy Values
Falsy:
javascript
false
null
0
""  (empty string)
Truthy:
javascript
true
42 (any non-zero number)
"hello" (any non-empty string)
[] (arrays)
{} (objects)
________________________________________
15. Complete Language Grammar (EBNF)
ebnf
program          = declaration* EOF

declaration      = varDecl | funcDecl | includeDecl | statement

varDecl          = "let" IDENTIFIER ("=" expression)? ";"?
funcDecl         = "func" IDENTIFIER "(" parameters? ")" block
parameters       = IDENTIFIER ("," IDENTIFIER)*
includeDecl      = "#include" STRING

statement        = exprStmt | ifStmt | switchStmt | whileStmt 
                 | forStmt | forOfStmt | returnStmt 
                 | breakStmt | continueStmt | block

exprStmt         = expression ";"?
ifStmt           = "if" "(" expression ")" statement ("else" statement)?
switchStmt       = "switch" "(" expression ")" "{" switchCase* "}"
switchCase       = "case" expression ":" statement* | "default" ":" statement*
whileStmt        = "while" "(" expression ")" statement
forStmt          = "for" "(" (varDecl | exprStmt | ";") expression? ";" expression? ")" statement
forOfStmt        = "for" "(" "let" IDENTIFIER ("," IDENTIFIER)? "of" expression ")" statement
returnStmt       = "return" expression? ";"?
breakStmt        = "break" ";"?
continueStmt     = "continue" ";"?
block            = "{" declaration* "}"

expression       = assignment
assignment       = (call "." IDENTIFIER | call "[" expression "]")? 
                   ("=" | "+=" | "-=" | "*=" | "/=") assignment | logicalOr
logicalOr        = logicalAnd ("||" logicalAnd)*
logicalAnd       = equality ("&&" equality)*
equality         = comparison (("==" | "!=") comparison)*
comparison       = addition ((">" | ">=" | "<" | "<=") addition)*
addition         = multiplication (("+" | "-") multiplication)*
multiplication   = unary (("*" | "/" | "%") unary)*
unary            = ("!" | "-" | "++" | "--") unary | postfix
postfix          = call ("++" | "--")?
call             = primary ("(" arguments? ")" | "." IDENTIFIER | "[" expression "]")*
arguments        = expression ("," expression)*

primary          = "true" | "false" | "null" | "this"
                 | NUMBER | STRING | IDENTIFIER
                 | "(" expression ")"
                 | arrayLiteral | objectLiteral | funcExpr

arrayLiteral     = "[" (expression ("," expression)*)? "]"
objectLiteral    = "{" (objectPair ("," objectPair)*)? "}"
objectPair       = IDENTIFIER ":" expression
                 | IDENTIFIER "(" parameters? ")" block
funcExpr         = "func" "(" parameters? ")" block
________________________________________
16. Example Programs
Hello World
javascript
print("Hello, World!");
________________________________________
Fibonacci
javascript
func fib(n) {
    if (n <= 1) return n;
    return fib(n - 1) + fib(n - 2);
}

for (let i = 0; i < 10; i += 1) {
    print(fib(i));
}
________________________________________
Factorial
javascript
func factorial(n) {
    if (n <= 1) return 1;
    return n * factorial(n - 1);
}

print(factorial(5));  // 120
________________________________________
Object-Oriented Pattern
javascript
func createPlayer(name, health) {
    return {
        name: name,
        health: health,
        x: 0,
        y: 0,
        
        move(dx, dy) {
            this.x += dx;
            this.y += dy;
            print(this.name + " moved to " + this.x + "," + this.y);
        },
        
        damage(amount) {
            this.health -= amount;
            if (this.health <= 0) {
                this.die();
            }
        },
        
        die() {
            print(this.name + " died!");
        }
    };
}

let player = createPlayer("Hero", 100);
player.move(10, 5);
player.damage(150);
________________________________________
Closures
javascript
func makeCounter() {
    let count = 0;
    
    return {
        increment() {
            count += 1;
            return count;
        },
        
        decrement() {
            count -= 1;
            return count;
        },
        
        get() {
            return count;
        }
    };
}

let counter = makeCounter();
print(counter.increment());  // 1
print(counter.increment());  // 2
print(counter.get());        // 2
________________________________________
Array Processing
javascript
let numbers = [1, 2, 3, 4, 5];

// Sum
let total = 0;
for (let num of numbers) {
    total += num;
}
print("Sum:", total);

// Filter evens
let evens = [];
for (let num of numbers) {
    if (num % 2 == 0) {
        evens.push(num);
    }
}
print("Evens:", evens);
________________________________________
17. Type Reference
Type	Literal	Example	Size
Number	42, 3.14	let x = 42;	8 bytes
String	"hello"	let s = "hi";	Variable
Boolean	true, false	let b = true;	8 bytes
Null	null	let n = null;	8 bytes
Array	[1, 2, 3]	let a = [];	Variable
Object	{x: 10}	let o = {};	Variable
Function	func() {}	let f = func(){}	Variable
________________________________________
18. Compilation Targets
Go + LLVM Backend
Source (.koda)
    ↓
Go Lexer → Tokens
    ↓
Go Parser → AST
    ↓
Go Compiler → LLVM IR
    ↓
LLVM → Native Code
    ↓
Executable
________________________________________
This is the complete Koda language specification.
Total keywords: 17
Total operators: ~25
Standard library functions: ~15
Designed for:
•	✅ Easy to learn (JavaScript-like)
•	✅ Fast execution (LLVM backend)
•	✅ Simple implementation (~5,000 lines Go)
•	✅ Desktop scripting focus

