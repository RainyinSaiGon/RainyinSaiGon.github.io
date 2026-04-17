title: Learning Java (Part 1)
date: 2026-03-26
description: A series about learning Java.
tags: Java
series: Learning Java
series_title: Introduction
---

In Vietnam, Java backend development is something everyone wants. I think if i ask some students about what they want to be after graduating, most will answer Java backend. Yeah, so why Java is so famous, what it haves and why people want to learn it. That wonder me a lot so i decided to research something about Java, make a big project and share something interesting with everyone.

## How Java run

Java is a **"Write once, run anywhere"** language, it means that Java is **platform-independent** and a Java file can run on both **Linux**, **macOS**, **Windows**. How can Java do it?

They use a **Java Virtual Machine (JVM)**, a virtual machine inside your **Java Developement Kit (JDK)** and they have both **compiler**, **interpreter** and **Just-in-time (JIT) compiler** in their language.

Their process look like this:

```text
.java -> compile to .class (bytecode) -> JVM run it
```
> **Note:** Bytecode is an intermediate representation of you code, it means that first Java compiler (javac keyword) will compile your file into bytecode file, not machine instructions. 

When running, first JVM still use Interpreter but the time they realize some blocks are reused frequently, they will use JIT to compile it to machine code. This technique ...is called adaptive optimization. The JVM watches which methods are executed many times (often called hot spots) and then compiles only those parts into optimized machine code. Because of that, Java does not need to compile everything at maximum optimization level from the beginning.

That is why some Java applications may feel a bit slower at startup, but after running for a while, performance improves significantly. This is one of the key reasons Java is widely used for long-running backend systems, where stable and optimized performance over time is very important.

We can run a simple program to check

``` java
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

        // Warm-up phase: JVM recognizes this as hot code and JIT-compiles it
        System.out.println("=== Warm-up Phase ===");
        for (int i = 0; i < warmupIterations; i++) {
            work(n);
        }

        // Measurement phase: now the code is JIT-compiled and optimized
        System.out.println("=== Measurement Phase ===");
        long t1 = System.nanoTime();
        long result = 0;
        
        for (int i = 0; i < measureIterations; i++) {
            result += work(n);
        }
        
        long t2 = System.nanoTime();

        double elapsedMs = (t2 - t1) / 1_000_000.0;
        double avgNs = (t2 - t1) / (double) measureIterations;

        System.out.println("Total time: " + elapsedMs + " ms");
        System.out.println("Average per call: " + avgNs + " ns");
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