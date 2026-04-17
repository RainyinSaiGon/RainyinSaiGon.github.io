title: Learning Go (Part 1)
date: 2026-03-19
description: A series about learning Go from zero — starting with what Go is, why it exists, and the basics of the language.
tags: Go
series: Learning Go
series_title: Introduction
---

I started learning Go because of the Kafka pet project I was building. My mentor recommended it, and after a few days with the language I understood why. Go is opinionated, fast, and surprisingly easy to get productive in. This series documents what I learn as I go — starting from the basics and moving into things like concurrency, networking, and building real tools.

## Why Go Exists

Go was created at Google in 2007 by Robert Griesemer, Rob Pike, and Ken Thompson. The motivation was straightforward: Google's infrastructure was mostly written in C++ and Java, and both were becoming painful at scale — slow compile times, verbose code, and awkward concurrency models.

Go was designed to fix these specific problems:

- **Fast compilation** — a large Go codebase compiles in seconds, not minutes
- **Simple concurrency** — goroutines and channels are built into the language
- **Readable code** — minimal syntax, no classes, no inheritance
- **Single binary output** — Go compiles to a standalone executable with no runtime dependency

The language launched publicly in 2009 and reached v1.0 in 2012. Today it powers tools like Docker, Kubernetes, Terraform, and Prometheus — essentially the backbone of modern cloud infrastructure.

## How Go Compiles and Runs

Unlike Java, Go compiles directly to **native machine code** — there is no intermediate bytecode or virtual machine.

```text
.go -> go compiler -> native binary (ELF/PE/Mach-O)
```

This means:

- Startup time is near-instant (no JVM warmup)
- The output is a single self-contained binary — no runtime to install
- Cross-compilation is built in (`GOOS=linux GOARCH=amd64 go build`)

The trade-off compared to Java's JIT: Go does not optimize hot code paths at runtime. The compiler does all optimization up front. In practice, Go's performance is excellent for I/O-bound and network workloads, which is where most backend code lives.

Go does have a garbage collector, so you do not manage memory manually (unlike C or Rust). The GC has improved significantly over the years — pause times are now in the microsecond range for most workloads.

## Setting Up

Install Go from [go.dev](https://go.dev/dl/). After installation:

```bash
go version
# go version go1.22.0 linux/amd64
```

Create a new module (the Go equivalent of a project):

```bash
mkdir hello && cd hello
go mod init hello
```

`go.mod` is the module file — it declares the module name and Go version, similar to `package.json` or `pom.xml`.

## Hello, World

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

```bash
go run main.go
# Hello, World!
```

A few things worth noting immediately:

- Every Go file belongs to a **package**. `package main` is the entry point package.
- `import` pulls in standard library packages. `fmt` handles formatted I/O.
- `func main()` is the entry point — no class wrapper needed.
- There are no semicolons. The formatter (`gofmt`) enforces this.

## Variables and Types

Go is statically typed. Every variable has a type known at compile time.

```go
// explicit type
var name string = "Vu"
var age  int    = 22

// short declaration (most common inside functions)
city := "Ho Chi Minh City"

// multiple assignment
x, y := 10, 20
```

The `:=` short declaration infers the type from the right side. It only works inside functions.

Go's basic types:

| Category | Types |
| --- | --- |
| Integer | `int`, `int8`, `int16`, `int32`, `int64` |
| Unsigned | `uint`, `uint8`, `uint16`, `uint32`, `uint64` |
| Float | `float32`, `float64` |
| String | `string` |
| Boolean | `bool` |
| Byte / Rune | `byte` (alias for `uint8`), `rune` (alias for `int32`) |

`int` is either 32 or 64 bits depending on the platform — on most modern systems it is 64 bits.

## Functions

```go
func add(a int, b int) int {
    return a + b
}

// shorthand when params share a type
func multiply(a, b int) int {
    return a * b
}

// multiple return values
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, fmt.Errorf("cannot divide by zero")
    }
    return a / b, nil
}
```

Multiple return values are idiomatic in Go — the second value is almost always an `error`. This is Go's error handling pattern: no exceptions, just values.

```go
result, err := divide(10, 3)
if err != nil {
    fmt.Println("Error:", err)
    return
}
fmt.Println(result)
```

The `if err != nil` pattern appears everywhere in Go code. It's verbose compared to try/catch, but explicit — you always know where errors can happen.

## Control Flow

```go
// standard if/else
if x > 10 {
    fmt.Println("big")
} else if x > 5 {
    fmt.Println("medium")
} else {
    fmt.Println("small")
}

// for is the only loop in Go
for i := 0; i < 5; i++ {
    fmt.Println(i)
}

// while-style loop
n := 0
for n < 10 {
    n++
}

// infinite loop
for {
    // break when done
    break
}
```

Go has no `while` keyword. The `for` loop covers all cases.

## Arrays, Slices, and Maps

```go
// array — fixed size, rarely used directly
arr := [3]int{1, 2, 3}

// slice — dynamic, the type you actually use
nums := []int{10, 20, 30}
nums = append(nums, 40)     // add element
fmt.Println(nums[1])        // 20
fmt.Println(nums[1:3])      // [20 30]

// map
scores := map[string]int{
    "Alice": 95,
    "Bob":   88,
}
scores["Charlie"] = 72

val, ok := scores["Alice"]  // ok is false if key doesn't exist
if ok {
    fmt.Println(val)
}
```

Slices are the workhorse of Go. They are backed by an underlying array but can grow dynamically via `append`. The three-index slice expression `s[low:high]` gives you a window into the backing array without copying.

## Structs

Go has no classes. Instead it uses **structs** — plain data containers.

```go
type User struct {
    Name  string
    Email string
    Age   int
}

u := User{
    Name:  "Dinh Dai Vu",
    Email: "vudd@i-stech.net",
    Age:   22,
}

fmt.Println(u.Name)
```

You attach behavior to structs using **methods**:

```go
func (u User) Greet() string {
    return "Hello, I'm " + u.Name
}

// pointer receiver — modifies the original
func (u *User) Birthday() {
    u.Age++
}
```

The `(u User)` before the function name is the **receiver**. Use a value receiver when you do not need to modify the struct. Use a pointer receiver (`*User`) when you do.

## Interfaces

Go interfaces are **implicit** — a type satisfies an interface simply by having the right methods. No `implements` keyword needed.

```go
type Greeter interface {
    Greet() string
}

func printGreeting(g Greeter) {
    fmt.Println(g.Greet())
}

// User satisfies Greeter because it has a Greet() method
printGreeting(u)
```

This is one of Go's best features. It makes code easy to compose and test — you can pass any type that has the required methods, without any inheritance hierarchy.

## Packages and Modules

Go organizes code into **packages**. Everything in the same directory belongs to the same package. Exported names start with a capital letter — lowercase is package-private.

```go
// math/calc.go
package math

func Add(a, b int) int { return a + b }   // exported
func helper() int      { return 0 }       // not exported
```

The module system (`go.mod`) manages dependencies. To add an external package:

```bash
go get github.com/some/package
```

Go downloads the dependency and records it in `go.mod` and `go.sum`. The `go.sum` file contains cryptographic checksums — Go verifies every download, which prevents supply-chain attacks.

## What's Next

This post covered the fundamentals: how Go compiles, the basic types, functions, error handling, structs, interfaces, and the module system. That's enough to read most Go code.

In the next part I'll go into what makes Go genuinely different from Java and Python: **goroutines and channels** — Go's concurrency model. This is where Go shines, and it's directly relevant to the Kafka project I'm building alongside this series.
