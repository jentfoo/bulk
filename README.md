# go-analyze/bulk

**Performance-first collection operations for Go**

`bulk` provides optimized utilities for working with slices and maps in Go, designed to **minimize memory allocations** and **maximize performance**. Unlike typical collection libraries that always copy, copies are avoided when possible (instead using views or truncated slices). And for uses where the input view can be discarded, `InPlace` functionality is offered to prevent copying entirely.

---

## Operation Types

**Default Operations** (recommended):
- Safe: never corrupt input data
- Smart: avoid allocations when possible (return views, truncate slices, etc.)
- Fast: better performance than standard "always copy" designs

**InPlace Operations** (maximum performance):
- Fastest: zero allocations by reusing input slice memory
- Unsafe: input slice is modified and must be discarded after use
- Critical: for performance-sensitive code where memory pressure matters

---

## Features

* **Zero-allocation InPlace variants** for maximum performance
* **Conditional optimizations** that avoid copies when safe
* Generic support using Go 1.18+ type parameters
* Simple, consistent API in a single package

---

## Installation

```bash
go get github.com/go-analyze/bulk
```

Import in your code:

```go
import "github.com/go-analyze/bulk"
```

---

## Usage

TODO - UNDER CONSTRUCTION

---

## When to use `bulk`

* **Large data analyses** where minimizing memory pressure is critical
* **Performance-sensitive loops** processing millions of elements
* Scenarios where **in-place mutations** are safe and desired

If you require copy-on-write semantics as your primary workflow, consider other collection utilities like [Pie](https://github.com/elliotchance/pie) instead, which always returns independent slices but incurs more and larger memory allocations.
