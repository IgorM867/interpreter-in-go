# MyGoInterpreter

MyGoInterpreter is a custom programming language interpreter written in Go. The interpreter supports various expressions, control structures, and native functions, providing a foundation for learning about language design and implementation.

## Features

- **Basic Expressions**: Supports arithmetic operations, boolean expressions, and variable assignments.
- **Control Structures**: Includes if and else statements
- **Functions:** Allows defining and invoking both user-defined and native functions.
- **Native Functions:** Provides built-in functions for common operations like printing

## Language Syntax

### Variables

```
let x = 10;
const y = x + 5;
const arr = [1,2,3];
```

### Arithmetic Operations

```
let result = (3 + 4) * 2 - 1;
```

### Boolean Expressions

```
let isTrue = !false;
let comparison = 10 >= 5;
```

### Control Structures

```
if (x > 10) {
    print("x is greater than 10");
} else {
    print("x is less than 10");
}
```

### Loops

```
let i = 0;
while (i < 10) {
    println(i)
    i = i + 1
}
```

### Functions

```
fn add(a, b) {
    a + b
}
```

### Native Functions

```
print()
println()
```
