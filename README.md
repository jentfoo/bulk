# go-analyze/bulk

[![license](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/go-analyze/bulk/blob/master/LICENSE)
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
* **Multi-slice operations** enabling intelligent optimizations across all slices
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

##### `SliceFilter[T any](predicate func(v T) bool, slices ...[]T) []T`
Filters elements that pass the predicate function from one or multiple slices. **Zero allocations** when all true elements are consecutive.

```go
numbers := []int{1, 2, 3, 4, 5, 6}
evens := bulk.SliceFilter(func(n int) bool { return n%2 == 0 }, numbers)
// Result: [2, 4, 6]
// Original slice unchanged: [1, 2, 3, 4, 5, 6]
```

##### `SliceFilterInPlace[T any](predicate func(v T) bool, slices ...[]T) []T`
**Zero-allocation** filtering by reusing input slice memory (zero allocation only guaranteed with a single slice, but will always have less allocations than SliceFilter).

```go
numbers := []int{1, 2, 3, 4, 5, 6} // This slice will be modified!
evens := bulk.SliceFilterInPlace(func(n int) bool { return n%2 == 0 }, numbers)
// Result: [2, 4, 6]
// Warning: numbers slice is now corrupted and must be discarded
```

#### Partitioning Operations

##### `SliceSplit[T any](predicate func(v T) bool, slices ...[]T) ([]T, []T)`
Partitions elements from one or multiple slices into two slices based on predicate.

```go
numbers := []int{1, 2, 3, 4, 5, 6}
evens, odds := bulk.SliceSplit(func(n int) bool { return n%2 == 0 }, numbers)
// evens: [2, 4, 6]
// odds: [1, 3, 5]
```

##### `SliceSplitInPlace[T any](predicate func(v T) bool, slice []T) ([]T, []T)`
**Memory-efficient** partitioning that reuses the input slice for one partition.

```go
numbers := []int{1, 2, 3, 4, 5, 6} // This slice will be modified!
evens, odds := bulk.SliceSplitInPlace(func(n int) bool { return n%2 == 0 }, numbers)
// evens: [2, 4, 6] (reuses original slice memory)
// odds: [1, 3, 5] (new allocation only when needed)
```

##### `SliceSplitInPlaceUnstable[T any](predicate func(v T) bool, slice []T) ([]T, []T)`
**Fastest** partitioning using two-pointer technique. Element order may change.

```go
numbers := []int{1, 2, 3, 4, 5, 6} // This slice will be modified!
evens, odds := bulk.SliceSplitInPlaceUnstable(func(n int) bool { return n%2 == 0 }, numbers)
// evens: [2, 4, 6] (order may differ from input)
// odds: [1, 3, 5] (order may differ from input)
```

#### Transformation Operations

##### `SliceTransform[I any, R any](conversion func(I) R, inputs ...[]I) []R`
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

##### `SliceToSet[T comparable](slices ...[]T) map[T]struct{}`
Converts slices to a set for fast lookup and deduplication. Accepts multiple slices for union operations.

```go
slice1 := []string{"a", "b", "c", "b"}
slice2 := []string{"c", "d", "e"}
set := bulk.SliceToSet(slice1, slice2)
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
unique := slices.Collect(maps.Keys(bulk.SliceToSet(duplicates)))
// Result: ["apple", "banana", "cherry"] (order may vary)
```

##### `SliceToSetBy[I any, R comparable](keyfunc func(I) R, slices ...[]I) map[R]struct{}`
Creates a set using a key function to transform elements into comparable keys. Typically used to provide a field from within the structs within the slice.

##### `SliceIntersect[T comparable](a, b []T) []T`
Returns elements that exist in both slices, preserving order from slice a. Duplicates are automatically removed from the result.

```go
numbers1 := []int{1, 2, 3, 4, 5}
numbers2 := []int{3, 4, 5, 6, 7}
intersection := bulk.SliceIntersect(numbers1, numbers2)
// Result: [3, 4, 5]
```

##### `SliceDifference[T comparable](a, b []T) []T`
Returns elements that exist in slice a but not in slice b, preserving order from slice a. Duplicates are automatically removed from the result.

```go
numbers1 := []int{1, 2, 3, 4, 5}
numbers2 := []int{3, 4, 5, 6, 7}
difference := bulk.SliceDifference(numbers1, numbers2)
// Result: [1, 2]
```

**Practical Set Operation Examples:**
```go
// Find common interests between users
userA := []string{"music", "sports", "reading", "cooking"}
userB := []string{"sports", "movies", "cooking", "travel"}
common := bulk.SliceIntersect(userA, userB)
// Result: ["sports", "cooking"]

// Find unique preferences for user A
unique := bulk.SliceDifference(userA, userB)
// Result: ["music", "reading"]

// Union can be achieved with SliceToSet for deduplication
union := slices.Collect(maps.Keys(bulk.SliceToSet(userA, userB)))
// Result: all unique elements from both slices
```

#### Counting Operations

##### `SliceToCounts[T comparable](slices ...[]T) map[T]int`
Counts occurrences of each element across multiple slices.

```go
slice1 := []string{"a", "b", "c", "b"}
slice2 := []string{"c", "d", "a"}
counts := bulk.SliceToCounts(slice1, slice2)
// Result: map[string]int{"a": 2, "b": 2, "c": 2, "d": 1}
```

##### `SliceToCountsBy[T any, K comparable](keyfunc func(T) K, slices ...[]T) map[K]int`
Counts occurrences using a key function to group elements.

#### Indexing Operations

##### `SliceToIndexBy[T any, K comparable](keyfunc func(T) K, slices ...[]T) map[K]T`
Creates an index map where each key maps to the **last** value encountered.

```go
type Person struct{ ID int; Name string }
people := []Person{{1, "Alice"}, {2, "Bob"}, {1, "Alice_Updated"}}
index := bulk.SliceToIndexBy(func(p Person) int { return p.ID }, people)
// Result: map[int]Person{1: {1, "Alice_Updated"}, 2: {2, "Bob"}} (last wins)
```

#### Grouping Operations

##### `SliceToGroupsBy[T any, K comparable](keyfunc func(T) K, slices ...[]T) map[K][]T`
Groups elements by key derived by each entry, preserving all values for each key.

```go
type Person struct{ Dept, Name string }
people := []Person{
    {"eng", "Alice"}, {"sales", "Bob"}, 
    {"eng", "Charlie"}, {"sales", "Dave"},
}
groups := bulk.SliceToGroupsBy(func(p Person) string { return p.Dept }, people)
// Result: map[string][]Person{
//   "eng": [{"eng", "Alice"}, {"eng", "Charlie"}],
//   "sales": [{"sales", "Bob"}, {"sales", "Dave"}]
// }
```

---

## When to use `bulk`

* **Large data analyses** where minimizing memory pressure is critical
* **Performance-sensitive loops** processing millions of elements
* Scenarios where **in-place mutations** are safe and desired

If you require copy-on-write semantics as your primary workflow, consider other collection utilities like [Pie](https://github.com/elliotchance/pie) instead, which always returns independent slices but incurs more and larger memory allocations.
