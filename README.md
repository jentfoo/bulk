# go-analyze/bulk

[![license](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/go-analyze/charts/blob/master/LICENSE)
[![Build Status](https://github.com/go-analyze/bulk/actions/workflows/tests-main.yml/badge.svg)](https://github.com/go-analyze/bulk/actions/workflows/tests-main.yml)

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

### Slice Operations

#### Filtering

##### `SliceFilter[T any](slice []T, predicate func(v T) bool) []T`
Filters elements that pass the predicate function. **Zero allocations** when all elements pass or fail.

```go
numbers := []int{1, 2, 3, 4, 5, 6}
evens := bulk.SliceFilter(numbers, func(n int) bool { return n%2 == 0 })
// Result: [2, 4, 6]
// Original slice unchanged: [1, 2, 3, 4, 5, 6]
```

##### `SliceFilterInPlace[T any](slice []T, predicate func(v T) bool) []T`
**Zero-allocation** filtering by reusing input slice memory.

```go
numbers := []int{1, 2, 3, 4, 5, 6} // This slice will be modified!
evens := bulk.SliceFilterInPlace(numbers, func(n int) bool { return n%2 == 0 })
// Result: [2, 4, 6]
// Warning: numbers slice is now corrupted and must be discarded
```

#### Partitioning Operations

##### `SliceSplit[T any](slice []T, predicate func(v T) bool) ([]T, []T)`
Partitions elements into two slices based on predicate.

```go
numbers := []int{1, 2, 3, 4, 5, 6}
evens, odds := bulk.SliceSplit(numbers, func(n int) bool { return n%2 == 0 })
// evens: [2, 4, 6]
// odds: [1, 3, 5]
```

##### `SliceSplitInPlace[T any](slice []T, predicate func(v T) bool) ([]T, []T)`
**Memory-efficient** partitioning that reuses the input slice for one partition.

```go
numbers := []int{1, 2, 3, 4, 5, 6} // This slice will be modified!
evens, odds := bulk.SliceSplitInPlace(numbers, func(n int) bool { return n%2 == 0 })
// evens: [2, 4, 6] (reuses original slice memory)
// odds: [1, 3, 5] (new allocation only when needed)
```

##### `SliceSplitInPlaceUnstable[T any](slice []T, predicate func(v T) bool) ([]T, []T)`
**Fastest** partitioning using two-pointer technique. Element order may change.

```go
numbers := []int{1, 2, 3, 4, 5, 6} // This slice will be modified!
evens, odds := bulk.SliceSplitInPlaceUnstable(numbers, func(n int) bool { return n%2 == 0 })
// evens: [2, 4, 6] (order may differ from input)
// odds: [1, 3, 5] (order may differ from input)
```

#### Transformation Operations

##### `SliceTransform[I any, R any](input []I, conversion func(I) R) []R`
Converts each element using the provided conversion function.

```go
strings := []string{"1", "2", "3", "4"}
numbers := bulk.SliceTransform(func(s string) int {
    n, _ := strconv.Atoi(s)
    return n
}, strings)
// Result: [1, 2, 3, 4]
```

#### Set Operations

##### `SliceToMap[T comparable](slices ...[]T) map[T]struct{}`
Converts slices to a map for fast lookup and deduplication. Accepts multiple slices for union operations.

```go
slice1 := []string{"a", "b", "c", "b"}
slice2 := []string{"c", "d", "e"}
set := bulk.SliceToMap(slice1, slice2)
// Result: map[string]struct{}{"a": {}, "b": {}, "c": {}, "d": {}, "e": {}} (order may vary)
```

**Deduplication Pattern:**
```go
import (
    "maps"
    "slices"
    "github.com/go-analyze/bulk"
)

duplicates := []string{"apple", "banana", "apple", "cherry", "banana"}
unique := slices.Collect(maps.Keys(bulk.SliceToMap(duplicates)))
// Result: ["apple", "banana", "cherry"] (order may vary)
```

---

## When to use `bulk`

* **Large data analyses** where minimizing memory pressure is critical
* **Performance-sensitive loops** processing millions of elements
* Scenarios where **in-place mutations** are safe and desired

If you require copy-on-write semantics as your primary workflow, consider other collection utilities like [Pie](https://github.com/elliotchance/pie) instead, which always returns independent slices but incurs more and larger memory allocations.
