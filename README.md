# go-analyze/bulk

**High-performance, large collection operations for Go**

`bulk` provides a suite of utilities for working with large in-memory data structures (slices and maps) in Go, designed to **minimize memory allocations** and **maximize performance**. Copies are generally avoided, and where reasonable `InPlace` functionality is offered to prevent copying entirely.

---

## Features

* **Zero-allocation, in-place operations** for filters, partitions, removals, and more
* Generic support using Go 1.18+ type parameters (any type `T`)
* Simple, consistent API: single package, `bulk`

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

If you require copy-on-write semantics as your primary workflow, consider other collection utilities like [Pie](https://github.com/elliotchance/pie) instead, which always returns independent slices but incurs more and larger memory allocations.

