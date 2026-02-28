title: Go Channels and Goroutines: A Practical Guide
date: 2026-01-20
description: How channels and goroutines work in Go, with real patterns for concurrent pipelines, fan-out/fan-in, and graceful cancellation.
tags: Go, Concurrency
---

<h2>The mental model</h2>

<p>Go's concurrency isn't thread-based in the traditional sense. Goroutines are multiplexed onto OS threads by the Go scheduler — you can have millions of goroutines on a few threads. Channels are the synchronisation primitive: typed conduits that goroutines use to send and receive values.</p>

<p>Rob Pike's motto: <em>"Do not communicate by sharing memory; instead, share memory by communicating."</em></p>

<h2>Buffered vs unbuffered channels</h2>

<pre><code>// Unbuffered — sender blocks until receiver is ready
ch := make(chan int)

// Buffered — sender only blocks when buffer is full
ch := make(chan int, 10)
</code></pre>

<p>Unbuffered channels are synchronous rendezvous points. Buffered channels introduce queuing. Start with unbuffered; add buffering only when you have measured a bottleneck.</p>

<h2>Pipeline pattern</h2>

<p>A pipeline connects stages where each stage reads from an input channel and writes to an output channel:</p>

<pre><code>func generate(nums ...int) &lt;-chan int {
    out := make(chan int)
    go func() {
        for _, n := range nums {
            out &lt;- n
        }
        close(out)
    }()
    return out
}

func square(in &lt;-chan int) &lt;-chan int {
    out := make(chan int)
    go func() {
        for n := range in {
            out &lt;- n * n
        }
        close(out)
    }()
    return out
}

// Usage
for result := range square(generate(2, 3, 4)) {
    fmt.Println(result) // 4, 9, 16
}
</code></pre>

<h2>Fan-out / Fan-in</h2>

<p>Distribute work across multiple goroutines, then collect results:</p>

<pre><code>func merge(cs ...&lt;-chan int) &lt;-chan int {
    var wg sync.WaitGroup
    out := make(chan int)
    output := func(c &lt;-chan int) {
        for n := range c { out &lt;- n }
        wg.Done()
    }
    wg.Add(len(cs))
    for _, c := range cs { go output(c) }
    go func() { wg.Wait(); close(out) }()
    return out
}
</code></pre>

<h2>Cancellation with context</h2>

<p>Always wire a <code>context.Context</code> through long-running pipelines. Use <code>ctx.Done()</code> as a select case:</p>

<pre><code>func generator(ctx context.Context, nums ...int) &lt;-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for _, n := range nums {
            select {
            case out &lt;- n:
            case &lt;-ctx.Done():
                return
            }
        }
    }()
    return out
}
</code></pre>

<h2>Common pitfalls</h2>

<ul>
<li><strong>Goroutine leaks:</strong> a goroutine blocked on a channel send/receive that's never consumed will live forever. Always pair a sender with a receiver and close channels when done.</li>
<li><strong>Closing a closed channel panics.</strong> Use a <code>sync.Once</code> or a done channel to signal completion.</li>
<li><strong>nil channel blocks forever.</strong> Useful trick: set a channel to nil inside a select to disable that case.</li>
</ul>

<h2>When not to use channels</h2>

<p>Channels are great for passing ownership and signalling. For protecting shared state that multiple goroutines read and write, a <code>sync.Mutex</code> or <code>sync.RWMutex</code> is often simpler and faster. Don't force channels where a mutex is clearer.</p>
