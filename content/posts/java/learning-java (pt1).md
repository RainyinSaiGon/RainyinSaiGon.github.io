title: Learning Java (Part 1)
date: 2026-03-26
description: A series about learning Java — starting from how the JVM works, then diving into concurrency.
tags: Java
series: Learning Java
series_title: Introduction
---

In Vietnam, Java backend development is something everyone wants. I think if I ask some students about what they want to be after graduating, most will answer Java backend. So why is Java so famous, what does it have, and why do people want to learn it? That made me curious, so I decided to research Java, build a big project, and share what I find along the way.

## How Java Runs

Java is a **"Write once, run anywhere"** language — it is **platform-independent** and a Java file can run on **Linux**, **macOS**, and **Windows**. How?

They use a **Java Virtual Machine (JVM)**, which lives inside your **Java Development Kit (JDK)**. The JDK ships with a **compiler**, an **interpreter**, and a **Just-in-Time (JIT) compiler**.

The process looks like this:

```text
.java -> compile to .class (bytecode) -> JVM runs it
```

> **Note:** Bytecode is an intermediate representation of your code. The Java compiler (`javac`) compiles your source file into bytecode, not machine instructions.

When running, the JVM starts with the interpreter. Once it notices that some code blocks are called frequently (called **hot spots**), it uses the JIT compiler to compile those parts into native machine code. This is called **adaptive optimization** — the JVM only optimizes what matters.

That is why Java apps can feel a bit slow at startup but get faster over time. This is one of the key reasons Java dominates long-running backend systems.

```java
public class JitDemo {

    static long work(int n) {
        long s = 0;
        for (int i = 0; i < n; i++) {
            s += i;
        }
        return s;
    }

    public static void main(String[] args) {
        final int n = 10_000;
        final int warmupIterations = 20_000;
        final int measureIterations = 50_000;

        System.out.println("=== Warm-up Phase ===");
        for (int i = 0; i < warmupIterations; i++) {
            work(n);
        }

        System.out.println("=== Measurement Phase ===");
        long t1 = System.nanoTime();
        long result = 0;

        for (int i = 0; i < measureIterations; i++) {
            result += work(n);
        }

        long t2 = System.nanoTime();
        System.out.println("Total time: " + (t2 - t1) / 1_000_000.0 + " ms");
        System.out.println("Average per call: " + (t2 - t1) / (double) measureIterations + " ns");
        System.out.println("Result: " + result);
    }
}
```

```text
=== Warm-up Phase ===
=== Measurement Phase ===
Total time: 195.626171 ms
Average per call: 3912.52342 ns
Result: 2499750000000
```

---

## Concurrency in Java

After understanding how the JVM runs code, the next big topic is **concurrency** — something Java was built with from day one. This is what makes Java genuinely powerful for backend systems, and honestly it's the part most people skip until they hit a production bug.

> *Concurrency is about dealing with lots of things at once. Parallelism is about doing lots of things at once.*
> — Rob Pike

### Java Is Made of Threads

Concurrency was an explicit design goal of Java. When Java was released in 1996, built-in threading support was a key differentiator — it removed developers' dependency on OS-specific features to achieve concurrency.

In Java, a **thread** is the smallest unit of execution — an independent path running within a program. Threads share the same address space (same variables and data structures), but each thread maintains its own **program counter**, **stack**, and **local variables**.

Even the most basic Java program is threaded. This might surprise you:

```java
public class HelloWorld {
    public static void main(String[] args) {
        System.out.println("Hello, World!");
        System.out.println("Executed by thread: " + Thread.currentThread().getName());
    }
}
```

```text
Hello, World!
Executed by thread: main
```

The JVM executes `main()` inside a thread called `main`. Threads are also the backbone of the garbage collector, JIT compiler, and the debugger — things most developers take for granted.

### Parallelism vs Concurrency

These two terms get mixed up a lot:

- **Parallelism** = doing multiple things *at the same time* (requires multiple CPU cores)
- **Concurrency** = *designing* programs so parts can overlap, even if not truly simultaneous

Think of it this way: parallelism is multiple workers building a house side by side. Concurrency is a single chef juggling multiple dishes — chopping onions, then checking the oven, then stirring the soup. The chef isn't doing everything at once, but the kitchen is always moving.

### Creating Threads

Java 1.0 gave us two ways to create threads: extending `Thread` or implementing `Runnable`. Both are still valid today.

```java
// Option 1: extend Thread
class MyThread extends Thread {
    @Override
    public void run() {
        System.out.println("Running in: " + getName());
    }
}

// Option 2: implement Runnable (preferred)
class MyRunnable implements Runnable {
    @Override
    public void run() {
        System.out.println("Running in: " + Thread.currentThread().getName());
    }
}

// Option 3: lambda (Java 8+)
Thread t = new Thread(() -> System.out.println("Lambda thread"));
t.start();
```

The key rule: always call `start()`, never `run()` directly. Calling `run()` executes it on the *current* thread, not a new one — a classic mistake.

`Runnable` became the preferred approach because it separates task logic from thread management and keeps the single inheritance slot open. Lambda expressions (Java 8) made this even more concise since `Runnable` is a functional interface.

### The Hidden Cost of Threads

Here's the part most tutorials skip: threads are expensive.

Each Java thread wraps a native OS thread and consumes around **2 MiB of memory** outside the heap by default (on most Linux environments). That sounds small, but if you're running a web server with thousands of concurrent connections, the memory bill adds up fast.

Beyond memory, there is **context switching** — every time the CPU switches between threads, it has to save the current thread's state and restore the next one. Under high load, this CPU overhead becomes significant.

You can test the limit on your machine:

```java
import java.util.concurrent.atomic.AtomicInteger;
import java.util.concurrent.locks.LockSupport;

public class ThreadLimitTest {
    public static void main(String[] args) {
        var threadCount = new AtomicInteger(0);
        try {
            while (true) {
                new Thread(() -> {
                    threadCount.incrementAndGet();
                    LockSupport.park();
                }).start();
            }
        } catch (OutOfMemoryError e) {
            System.out.println("Reached thread limit: " + threadCount);
        }
    }
}
```

On a typical machine you will hit `OutOfMemoryError` around 10,000–20,000 threads. The OS imposes a hard ceiling, so **"just add more threads"** is not a real scaling strategy.

### The Real Problem: Blocking

The deeper issue is not just the number of threads — it's that threads **block**. Consider a method that calculates a credit score:

```java
public Credit calculateCredit(Long personId) {
    var person = getPerson(personId);       // database call — blocks thread
    var assets = getAssets(person);         // API call    — blocks thread
    var liabilities = getLiabilities(person); // database call — blocks thread
    importantWork();
    return calculateCredits(assets, liabilities);
}
```

If each call takes 200 ms, the total is 1,000 ms (1 second). During every database or API call, the thread is just sitting there waiting — it cannot be used for anything else. In a high-throughput application, this is how you run out of threads under load.

### Making It Parallel

Notice that `getAssets()`, `getLiabilities()`, and `importantWork()` have no dependency on each other. We can run them in parallel:

```java
public Credit calculateCreditWithExecutor(Long personId)
        throws ExecutionException, InterruptedException {

    try (ExecutorService executor = Executors.newFixedThreadPool(5)) {
        var person = getPerson(personId);

        var assetsFuture      = executor.submit(() -> getAssets(person));
        var liabilitiesFuture = executor.submit(() -> getLiabilities(person));
        executor.submit(() -> importantWork());

        return calculateCredits(assetsFuture.get(), liabilitiesFuture.get());
    }
}
```

With an `ExecutorService`, we submit tasks to a thread pool. The pool reuses pre-created threads, avoiding the overhead of creating and destroying threads on every request. The result: what took 1,000 ms sequentially now takes around **600 ms**.

The `ExecutorService` also gives us lifecycle management — the `try-with-resources` block automatically shuts the pool down when we are done.

### The Executor Framework's Limitations

The Executor framework is better, but not perfect:

- **`Future.get()` still blocks** — we shifted the blocking from the I/O call to the result retrieval, but blocking is still there.
- **No composability** — chaining multiple async operations is awkward with `Future`.

### CompletableFuture: Composable Async

Java 8 introduced `CompletableFuture` to address the composability problem. Instead of blocking on `get()`, you chain operations:

```java
public Credit calculateCreditWithCompletableFuture(Long personId)
        throws InterruptedException, ExecutionException {

    return CompletableFuture
        .runAsync(() -> importantWork())
        .thenCompose(v -> CompletableFuture.supplyAsync(() -> getPerson(personId)))
        .thenCombineAsync(
            CompletableFuture.supplyAsync(() -> getAssets(getPerson(personId))),
            (person, assets) -> calculateCredits(assets, getLiabilities(person))
        )
        .get();
}
```

`CompletableFuture` lets you define workflows as a chain of transformations. Built-in methods like `exceptionally`, `handle`, and `whenComplete` handle errors cleanly. Under the hood it uses the Fork/Join pool, so task distribution is efficient.

The downside: the API has a steep learning curve. The chained style makes debugging harder — setting a breakpoint and stepping through the code doesn't work the same way when execution jumps between threads and callbacks.

### The Fork/Join Pool

Between the basic `ExecutorService` and `CompletableFuture` sits the **Fork/Join Pool** (introduced in Java 7). It solves a specific problem: when tasks create subtasks, how do you keep all threads busy?

The answer is **work stealing**: each thread in the pool has its own local queue. When a thread finishes its queue, it "steals" tasks from another thread's queue. This keeps all cores busy and avoids idle time.

```java
ForkJoinPool pool = new ForkJoinPool();
pool.submit(() -> {
    // your parallel tasks here
}).join();
```

The Fork/Join pool also promotes **cache affinity** — related tasks tend to run on the same CPU core, so data loaded into that core's L1/L2 cache stays warm and available to subsequent tasks.

### Reactive Programming

For applications with very high I/O throughput, frameworks like **Spring WebFlux**, **RxJava**, and **Project Reactor** take a different approach: instead of threads blocking on I/O, you register callbacks that fire when data is ready. No thread is ever blocked.

```java
public Mono<Credit> calculateCreditReactive(Long personId) {
    Mono<Person> personMono    = Mono.fromSupplier(() -> getPerson(personId));
    Mono<List<Asset>> assetsMono = personMono.map(p -> getAssets(p));
    Mono<List<Liability>> liabilitiesMono = personMono.map(p -> getLiabilities(p));

    return Mono.zip(assetsMono, liabilitiesMono)
               .map(tuple -> calculateCredits(tuple.getT1(), tuple.getT2()));
}
```

The performance ceiling is very high, but the learning investment is steep. You need to understand Publishers, Subscribers, backpressure, and schedulers. Debugging a reactive pipeline with dozens of operators and multiple async boundaries is genuinely hard.

### Virtual Threads: The Modern Answer

Java 21 shipped **Project Loom** and with it **virtual threads** — lightweight threads managed by the JVM, not the OS. You can create millions of them without running out of memory.

```java
public Credit calculateCreditWithVirtualThread(Long personId)
        throws ExecutionException, InterruptedException {

    try (var executor = Executors.newVirtualThreadPerTaskExecutor()) {
        var person            = getPerson(personId);
        var assetsFuture      = executor.submit(() -> getAssets(person));
        var liabilitiesFuture = executor.submit(() -> getLiabilities(person));
        executor.submit(this::importantWork);

        return calculateCredits(assetsFuture.get(), liabilitiesFuture.get());
    }
}
```

The only change from the `ExecutorService` version is `Executors.newVirtualThreadPerTaskExecutor()`. That's it.

When a virtual thread hits a blocking operation (database call, network I/O, `Thread.sleep()`), it **yields** control back to the underlying platform thread — which then picks up another virtual thread and keeps running. Once the blocking call completes, the virtual thread resumes. From the developer's perspective, the code looks sequential and imperative. From the runtime's perspective, the platform threads are always busy.

This is what makes virtual threads exciting: you get the performance of reactive programming without giving up readable, debuggable, sequential code.

### Where We Are

Java's concurrency story has evolved over 30 years:

| Era | Tool | Problem Solved | Remaining Issue |
|-----|------|----------------|-----------------|
| Java 1.0 | `Thread` | Built-in threading | Expensive, unmanaged |
| Java 5 | `ExecutorService` | Thread pooling, lifecycle | `Future.get()` still blocks |
| Java 7 | Fork/Join Pool | Work stealing, cache affinity | Complex task decomposition |
| Java 8 | `CompletableFuture` | Composable async | Steep learning curve |
| Java 9+ | Reactive (WebFlux) | Non-blocking I/O at scale | Very steep learning curve |
| Java 21 | Virtual Threads | Blocking code without cost | Still learning the limits |

In the next part I'll go deeper into the Virtual Thread  and how these concurrency primitives show up in real backend applications.
