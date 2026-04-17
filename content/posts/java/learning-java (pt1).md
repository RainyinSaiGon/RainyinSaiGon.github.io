title: Learning Java (Part 1)
date: 2026-03-26
description: A series about learning Java — starting from how the JVM works, then diving into concurrency, the Java Memory Model, and low-level CPU concepts like false sharing.
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

### Thread Lifecycle

A Java thread is never just "running" or "not running." It moves through six distinct states during its lifetime:

```text
NEW → RUNNABLE → (BLOCKED | WAITING | TIMED_WAITING) → TERMINATED
```

| State | What it means |
|---|---|
| `NEW` | Thread object created but `start()` not called yet |
| `RUNNABLE` | Eligible to run — may be executing or queued for a CPU slice |
| `BLOCKED` | Waiting to acquire a monitor lock held by another thread |
| `WAITING` | Indefinitely waiting for another thread to signal (`Object.wait()`, `Thread.join()`) |
| `TIMED_WAITING` | Same as WAITING but with a timeout (`Thread.sleep(n)`, `wait(n)`) |
| `TERMINATED` | `run()` finished or an uncaught exception killed it |

The distinction between `BLOCKED` and `WAITING` matters in practice. A `BLOCKED` thread is competing for a lock — it will automatically become `RUNNABLE` once the lock is released. A `WAITING` thread voluntarily gave up control and will only resume when explicitly notified with `notify()` or `notifyAll()`.

You can inspect thread state programmatically:

```java
Thread t = new Thread(() -> {
    try { Thread.sleep(5000); } catch (InterruptedException e) {}
});
System.out.println(t.getState()); // NEW
t.start();
Thread.sleep(100);
System.out.println(t.getState()); // TIMED_WAITING
t.join();
System.out.println(t.getState()); // TERMINATED
```

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

### Context Switching: What Actually Happens

When the OS decides to stop running thread A and start running thread B, it performs a **context switch**. This is not free — here is what happens at the CPU level:

1. **Save the register file** — The CPU has a fixed set of registers (program counter, stack pointer, general-purpose registers). The OS saves all of them for thread A into a kernel data structure called the **Thread Control Block (TCB)**.
2. **Flush the instruction pipeline** — Modern CPUs speculatively decode and execute future instructions. When switching threads, all in-flight speculative work must be discarded.
3. **Restore the register file** — The TCB for thread B is loaded back into the CPU registers.
4. **TLB flush (sometimes)** — The **Translation Lookaside Buffer (TLB)** is a hardware cache that maps virtual memory addresses to physical page addresses. Switching between different *processes* (not threads within the same process) invalidates the TLB, forcing expensive page-table walks on the next memory accesses.
5. **Cache warming** — Thread B may have been working on entirely different data. The L1/L2 caches are still populated with thread A's data, which is useless for thread B. Thread B will suffer cache misses until its working set is loaded into cache.

A single context switch costs roughly **1–10 microseconds** of pure CPU overhead. On a busy web server doing tens of thousands of context switches per second, this is measurable latency — not a theoretical concern.

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

The answer is **work stealing**: each thread in the pool has its own local deque (double-ended queue). Threads push and pop their own tasks from one end. When a thread runs out of work, it "steals" tasks from the *other end* of another thread's deque. This design minimizes contention — the owner and the thief work at opposite ends, so they rarely collide.

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

---

## The Java Memory Model

Concurrency is not just about threads — it is also about **memory**. In a multi-core system, each CPU core has its own private caches. When two threads on different cores read and write the same variable, you cannot assume they see each other's changes without explicit synchronization. The **Java Memory Model (JMM)** defines the rules.

### The Visibility Problem

Consider this program:

```java
public class VisibilityDemo {
    private static boolean stop = false;

    public static void main(String[] args) throws InterruptedException {
        Thread worker = new Thread(() -> {
            int i = 0;
            while (!stop) {   // reads stop from register or local cache
                i++;
            }
            System.out.println("Stopped at: " + i);
        });

        worker.start();
        Thread.sleep(1000);
        stop = true;   // write may never reach the worker's cache
        System.out.println("Main set stop = true");
    }
}
```

On many JVMs with JIT enabled, this program loops forever. The worker thread caches `stop` in a CPU register or L1 cache. The main thread's write to `stop` goes to *its own* cache, which is never flushed to the worker's cache. The worker never sees the change.

This is not a bug — it is legal under the JMM. The JMM allows the JIT and CPU to aggressively reorder and cache memory operations unless you use **synchronization constructs** that establish a happens-before relationship.

### Happens-Before

The central concept of the JMM is **happens-before**: if action A *happens-before* action B, then every write A performed is guaranteed visible to B.

The JMM guarantees happens-before in these situations:

| Rule | Guarantee |
|---|---|
| **Program order** | Within a single thread, each action happens-before the next one in source order |
| **Monitor unlock** | Releasing a `synchronized` lock happens-before any subsequent acquisition of that same lock |
| **Volatile write** | A write to a `volatile` field happens-before every subsequent read of that same field |
| **Thread start** | `Thread.start()` happens-before any action inside the new thread |
| **Thread join** | All actions inside a thread happen-before `Thread.join()` returns |
| **Transitivity** | If A happens-before B and B happens-before C, then A happens-before C |

Without happens-before, shared state between threads is undefined — even something as simple as a boolean flag.

### `volatile`: Visibility Without Locking

Declaring a field `volatile` establishes a happens-before edge between every write and every subsequent read of that field. Under the hood, the JVM emits a **memory fence** (also called a memory barrier) — a CPU instruction that:

- Forces all pending writes to be flushed from the CPU's store buffer to main memory (or at least to a cache level visible to other cores).
- Prevents the CPU and compiler from reordering instructions across the fence.

```java
public class VisibilityFixed {
    private static volatile boolean stop = false;

    public static void main(String[] args) throws InterruptedException {
        Thread worker = new Thread(() -> {
            int i = 0;
            while (!stop) { i++; }  // volatile read — always fetches fresh value
            System.out.println("Stopped at: " + i);
        });

        worker.start();
        Thread.sleep(1000);
        stop = true;   // volatile write — immediately visible to all threads
    }
}
```

`volatile` is cheap — it does not acquire a lock. But it only gives you **visibility**, not **atomicity**. The compound operation `counter++` is three steps (read, add 1, write), and `volatile` does not protect it from race conditions between those steps.

### `synchronized`: Atomicity and Visibility

`synchronized` on a method or block acquires a **monitor lock** before entering and releases it on exit. Every thread that acquires the same monitor is guaranteed to see all writes made by the previous holder of that lock.

```java
public class SafeCounter {
    private int count = 0;  // no need for volatile — synchronized covers visibility

    public synchronized void increment() {
        count++;  // read-modify-write is now atomic
    }

    public synchronized int get() {
        return count;
    }
}
```

`synchronized` is heavier than `volatile` because it involves lock acquisition (which may put the thread in `BLOCKED` state if another thread holds the lock). It is the right tool when you need to protect a multi-step read-modify-write sequence.

### `ReentrantLock`: Explicit Locking

`java.util.concurrent.locks.ReentrantLock` is the explicit alternative to `synchronized`. It provides the same mutual exclusion and visibility guarantees, but with more control:

```java
import java.util.concurrent.locks.ReentrantLock;

public class BetterCounter {
    private final ReentrantLock lock = new ReentrantLock();
    private int count = 0;

    public void increment() {
        lock.lock();
        try {
            count++;
        } finally {
            lock.unlock();  // always in finally — never skip the unlock
        }
    }

    // Non-blocking: returns false immediately if lock is held
    public boolean tryIncrement() {
        if (lock.tryLock()) {
            try {
                count++;
                return true;
            } finally {
                lock.unlock();
            }
        }
        return false;
    }
}
```

Key advantages over `synchronized`:

- **`tryLock()`** — attempt to acquire without blocking; useful for avoiding deadlocks.
- **`tryLock(time, unit)`** — wait up to a timeout before giving up.
- **`lockInterruptibly()`** — allow a waiting thread to be interrupted.
- **Fairness** — `new ReentrantLock(true)` grants the lock to the longest-waiting thread, preventing starvation.

### Compare-And-Swap: Lock-Free Atomicity

For simple counters and references, the `java.util.concurrent.atomic` package provides **lock-free thread safety** using the hardware **Compare-And-Swap (CAS)** instruction.

**What CAS does:** CAS takes three operands — a memory address, an expected value, and a new value. It reads the current value at the address, compares it to the expected value, and only writes the new value if they match. All three steps happen atomically at the hardware level. No lock, no OS call.

```text
CAS(address, expected, new_value):
    if *address == expected:
        *address = new_value
        return SUCCESS
    else:
        return FAILURE  (current value unchanged)
```

Java exposes CAS through `AtomicInteger`, `AtomicLong`, `AtomicReference`, and others:

```java
import java.util.concurrent.atomic.AtomicInteger;

AtomicInteger counter = new AtomicInteger(0);

// Internally: CAS loop until the increment succeeds
counter.incrementAndGet();

// Explicit CAS — only set to 10 if current value is 5
boolean updated = counter.compareAndSet(5, 10);
```

The internal implementation of `incrementAndGet` is conceptually:

```java
int current, next;
do {
    current = get();       // read current value
    next = current + 1;
} while (!compareAndSet(current, next));  // retry if someone else changed it
return next;
```

This is called a **CAS loop** or **spin loop**. Under low contention it is very fast — usually just one iteration. Under high contention (many threads all trying to increment simultaneously), the loop retries frequently, burning CPU. In those cases, `LongAdder` (also in `java.util.concurrent.atomic`) uses per-thread cells to spread contention across multiple counters and sum them on read.

---

## CPU Caches and Cache Lines

To understand the next topic (false sharing), you need to understand how CPU caches work at the hardware level.

### The Memory Hierarchy

Accessing main memory (RAM) takes roughly **50–100 nanoseconds**. A modern CPU core executes an instruction roughly every **0.3 nanoseconds**. If the CPU waited for RAM on every read, it would be idle more than 99% of the time.

To close this gap, CPUs have multiple levels of fast, on-chip memory called **caches**:

| Level | Location | Typical Size | Latency |
|---|---|---|---|
| **Registers** | Inside the execution unit | Bytes | 0 ns (immediate) |
| **L1 cache** | Per core, on die | 32–64 KB | ~1 ns (4 cycles) |
| **L2 cache** | Per core, on die | 256 KB – 1 MB | ~4 ns (12 cycles) |
| **L3 cache** | Shared by all cores | 4–64 MB | ~10–40 ns |
| **RAM** | Off-chip DIMM | GBs | ~50–100 ns |
| **NVMe SSD** | PCIe slot | TBs | ~100,000 ns |

When the CPU needs a value, it checks L1 first. If found, that is an **L1 cache hit** — nearly free. If not found (**L1 miss**), it checks L2, then L3, then finally fetches from RAM (**LLC miss** — Last Level Cache miss). Each miss step is 3–10× slower than the previous.

### Cache Lines: The Granularity of Transfer

Caches do not store individual bytes or words. They move memory in fixed-size chunks called **cache lines** — **64 bytes** on virtually every modern x86 and ARM processor.

When the CPU reads a single `int` (4 bytes) from RAM, it actually fetches the entire 64-byte aligned block containing that `int` into the L1 cache. This exploits **spatial locality**: if you need one value, you will likely need the nearby values soon.

This is why iterating over an array in order is dramatically faster than jumping around randomly:

```java
int[] arr = new int[16_000_000];

// Fast: sequential — each 64-byte cache line holds 16 ints
// Load one cache line, use 16 ints, load next cache line
long sum = 0;
for (int i = 0; i < arr.length; i++) {
    sum += arr[i];
}

// Slow: stride-64 — every access is a guaranteed cache miss
// Load a cache line, use 1 int, load next cache line
for (int i = 0; i < arr.length; i += 16) {
    sum += arr[i];
}
```

Both loops do the same number of cache line loads, but the sequential loop uses all 16 ints from each line. The stride-64 loop wastes 15 of every 16 slots — effectively paying cache miss cost for every element.

---

## False Sharing: The Silent Performance Killer

With cache lines understood, we can address one of the most subtle and destructive performance problems in concurrent Java: **false sharing**.

### Definition

**False sharing** is a performance degradation that occurs when two or more threads on different CPU cores repeatedly write to *different variables* that happen to reside within the *same 64-byte cache line*.

The "false" in false sharing is critical: the threads have no logical dependency on each other's data. They are not sharing any variable. But because their variables are packed together within the same cache line, the hardware treats the entire line as a single unit of sharing — and coherence traffic is generated as if they were writing to the same memory location.

### The MESI Cache Coherence Protocol

Modern multi-core CPUs maintain a consistent view of memory across all cores using a **cache coherence protocol**. The most widely used is **MESI**, named after the four states a cache line can be in at any given core:

| State | Symbol | Meaning |
|---|---|---|
| **Modified** | M | This core has the **only** valid copy. It has been written to and differs from RAM. All other cores have this line in Invalid state. |
| **Exclusive** | E | This core has the **only** copy, but it has **not** been modified. Matches RAM. Can transition to Modified on write without notifying others. |
| **Shared** | S | **Multiple** cores hold valid copies, all matching RAM. A write by any core must first invalidate all others. |
| **Invalid** | I | This core's copy is **stale**. Must be fetched before use. |

The protocol enforces one invariant: at most one core can hold a line in M or E state; multiple cores can hold a line in S state; a line can be Invalid at any number of cores.

### Why False Sharing Hurts: The RFO Cycle

**RFO** stands for **Request For Ownership** — the message a core sends when it needs to transition a cache line from Shared or Invalid to Exclusive (in preparation for writing).

Here is the false sharing sequence when Thread 0 (on Core 0) writes variable `A` and Thread 1 (on Core 1) writes variable `B`, where `A` and `B` share a cache line:

```text
Initial state: both cores have the line in Shared (S) state

Step 1 — Core 0 wants to write A:
  Core 0 sends RFO to all other cores.
  Core 1 transitions its copy to Invalid (I).
  Core 0 gets the line in Modified (M) state.
  Core 0 writes to A.

Step 2 — Core 1 wants to write B (different variable, same cache line):
  Core 1 sends RFO to all other cores.
  Core 0 must write-back its Modified copy to L3/RAM.
  Core 0 transitions its copy to Invalid (I).
  Core 1 fetches the updated line from L3/RAM.
  Core 1 gets the line in Modified (M) state.
  Core 1 writes to B.

Step 3 — Core 0 wants to write A again:
  [repeat from Step 1]
```

This back-and-forth is called **cache line ping-pong**. Each full cycle involves:

- An RFO message on the inter-core bus
- A write-back from the modifying core to L3 or RAM
- A cache line fetch by the requesting core from L3 or RAM
- Transition overhead through the coherence state machine

Each cycle costs roughly **100–300 CPU cycles** — far more expensive than an L1 cache hit (~4 cycles). When two threads are running a tight write loop, this happens millions of times per second.

### Demonstrating the Problem

```java
public class FalseSharingDemo {

    static long counterA = 0L;
    static long counterB = 0L;
    // These two longs are 8 bytes each, 16 bytes total — same cache line

    static final long ITERATIONS = 500_000_000L;

    public static void main(String[] args) throws InterruptedException {
        Thread t0 = new Thread(() -> {
            for (long i = 0; i < ITERATIONS; i++) counterA++;
        });
        Thread t1 = new Thread(() -> {
            for (long i = 0; i < ITERATIONS; i++) counterB++;
        });

        long start = System.nanoTime();
        t0.start(); t1.start();
        t0.join();  t1.join();

        System.out.printf("False sharing:  %d ms%n",
            (System.nanoTime() - start) / 1_000_000);
    }
}
```

The two threads touch logically independent variables. Yet on a typical machine this runs in **~4,000 ms** because the cores constantly invalidate each other's cache line.

For comparison, running both increments on a single thread (no contention) takes **~1,200 ms**. The false sharing overhead is **~3×** slower than sequential execution — a genuine regression from adding parallelism.

### Fix 1: Manual Padding

The solution is to push the two variables onto separate cache lines by inserting enough padding bytes between them:

```java
public class FalseSharingPadded {

    static long counterA = 0L;
    // 7 longs × 8 bytes = 56 bytes padding + 8 bytes for counterA = 64 bytes (one cache line)
    static long p1, p2, p3, p4, p5, p6, p7;

    static long counterB = 0L;
    static long p8, p9, p10, p11, p12, p13, p14;

    // same benchmark...
}
```

Now `counterA` and `counterB` occupy different cache lines. Core 0 can write to `counterA` without affecting Core 1's copy of `counterB`.

Result: **~700 ms** — roughly **5–6× faster** than the false-sharing version, with zero logic changes.

### Fix 2: `@Contended` Annotation

Manual padding is brittle. The JVM can reorder fields during class loading, and the JIT can eliminate unused fields entirely. Java 8 introduced `@jdk.internal.vm.annotation.Contended` to solve this properly:

```java
// Available as sun.misc.Contended in Java 8, or via --add-opens in Java 9+
import jdk.internal.vm.annotation.Contended;

public class FalseSharingAnnotated {

    @Contended
    static long counterA = 0L;

    @Contended
    static long counterB = 0L;
}
```

Run with:

```bash
java -XX:-RestrictContended FalseSharingAnnotated
```

The JVM inserts **128 bytes** of padding around each `@Contended` field (not 64, because hardware prefetchers can read ahead by a full cache line, so 128 bytes provides a safety margin). This is the exact technique used inside `java.util.concurrent.atomic.LongAdder`, `ConcurrentHashMap`, and `ForkJoinPool` in the JDK itself.

### When False Sharing Appears in Practice

False sharing is easy to create accidentally:

- **Per-thread counters in an array** — `long[] counts = new long[numThreads]`; if each thread writes to `counts[threadId]`, adjacent elements share cache lines.
- **Ring buffer head and tail pointers** — if head and tail are adjacent fields in a class, the producer and consumer constantly invalidate each other's cache line. LMAX Disruptor explicitly pads both pointers.
- **Node-based data structures** — fields of a `Node` written by different threads (e.g., next pointer and value) can coexist in one cache line.

The fix is always: identify which variables are written by different threads, and ensure they are on separate cache lines.

---

## Where We Are

Java's concurrency story has evolved over 30 years:

| Era | Tool | Problem Solved | Remaining Issue |
|-----|------|----------------|-----------------|
| Java 1.0 | `Thread` | Built-in threading | Expensive, unmanaged |
| Java 5 | `ExecutorService` | Thread pooling, lifecycle | `Future.get()` still blocks |
| Java 7 | Fork/Join Pool | Work stealing, cache affinity | Complex task decomposition |
| Java 8 | `CompletableFuture` | Composable async | Steep learning curve |
| Java 9+ | Reactive (WebFlux) | Non-blocking I/O at scale | Very steep learning curve |
| Java 21 | Virtual Threads | Blocking code without cost | Still learning the limits |

And the memory/synchronization tools:

| Tool | Visibility | Atomicity | Cost |
|---|---|---|---|
| `volatile` | Yes | No | Low — memory fence, no lock |
| `synchronized` | Yes | Yes | Medium — lock acquisition, possible blocking |
| `ReentrantLock` | Yes | Yes | Medium — same as synchronized, more flexible |
| `AtomicInteger` (CAS) | Yes | Yes | Low under low contention, degrades under high |
| Padding / `@Contended` | N/A | N/A | Memory overhead (~128 bytes per field) |

In the next part I'll go deeper into the Spring Boot ecosystem and how these primitives show up in real backend applications — thread pools, connection pools, request handling, and where developers most commonly get burned.
