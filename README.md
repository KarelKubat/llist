# llist

`llist` is a doubly linked list, with handy features.

<!-- toc -->
- [Purpose](#purpose)
- [Synopsis](#synopsis)
- [Description](#description)
  - [Don't break the encapsulation: use <code>llist</code>'s methods](#dont-break-the-encapsulation-use-llists-methods)
<!-- /toc -->

## Purpose

Package `llist` encapsulates the double linked list of [github.com/KarelKubat/lnode](https://github.com/KarelKubat/lnode) but provides O(1) speed in:

- Determining the head and tail of the contained list
- Fetching nodes by value

Inserting or deleting nodes, and changing the value of a given node, has already O(1) time complexity in the internally contained `lnode`s.

Having O(1) time complexity for all operations has drawbacks:

- `llist` has to duplicate a linked list, but rearranged as a map. This effectively duplicates the memory requirement.
- Inserting nodes, or deleting, or changing the values must occur through `llist`'s methods, so that the internal administration is kept up to date. These actions must not be performed by directly manipulating the `lnode`'s fields `Next`, `Prev` and `Value` (or by using the corresponding `lnode`'s methods).

> My application for this is that I needed an LRU (Least Recently Used) memory cache of fixed size, where all operations would run in O(1) time complexity. The LRU cache can implement this using `llist`:  Most recent entries can be kept at the head of the list, least recent at the tail. Lookups and deletes (evictions) are lightning fast given that `llist` runs at O(1).

**Package `llist` is not thread-safe. The caller must ensure that concurrent updates are mutex-protected.**

## Synopsis

Not all methods are shown in this listing.

```go
package main

import (
    "fmt"

    "github.com/KarelKubat/llist"
    "github.com/KarelKubat/lnode"
)

var words = []string{"the", "quick", "brown", "fox", "jumped", "over", "the", "lazy", "dog"}

func print(title string, l *llist.LList[string]) {
    fmt.Printf("\n%s: ", title)

    // Walk all nodes, starting at the head.
    l.Head().VisitByNext(func(node *lnode.Node[string]) bool {
        fmt.Print(node.Value, " ") // Print contained string
        return true                // Continue visiting
    })
    fmt.Println()

    // Print how many times each word occurs.
    for _, word := range words {
        nds := l.FindNodes(word)
        fmt.Println("  Word", word, "occurs", len(nds), "times")
    }
}

func main() {
    // Initialize the list.
    l := llist.New[string]()
    for _, w := range words {
        l.Append(l.Tail(), lnode.New[string](w))
    }
    // Structure:
    // the - quick - brown - fox - jumped - over - the - lazy - dog
    // ^head                                                    ^tail
    print("Vanilla list", l)

    // Prepend to the start of the list.
    l.Prepend(l.Head(), lnode.New[string]("yesterday"))
    // Structure:
    // yesterday - the - quick - brown - fox - jumped - over - the - lazy - dog
    // ^head                                                                ^tail
    print("After prepending 'yesterday'", l)

    // Set fx to point to 5th node, or 4x Next from head ("fox")
    fx := l.Head().Next.Next.Next.Next
    // Insert beyond fx more "fox" nodes
    l.Append(fx, lnode.New[string]("fox"))
    l.Append(fx, lnode.New[string]("fox"))
    l.Append(fx, lnode.New[string]("fox"))

    // Structure:
    // yesterday - the - quick - brown - fox - fox - fox - fox - jumped (etc.)
    // ^head
    print("Three more foxes", l)

    // Remove 3x the 5th node.
    for range 3 {
        l.Delete(l.Head().Next.Next.Next.Next)
    }
    // Structure:
    // yesterday - the - quick - brown - fox - jumped - over - the - lazy - dog
    // ^head                                                                ^tail
    print("Three deletes of the 5th node", l)

    // Change all nodes with "lazy" to "quick" (in this example just one)
    for _, node := range l.FindNodes("lazy") {
        node.SetValue(node, "quick")
    }
    print("Lazy is now quick", l)
}
```

## Description

- `Head()` and `Tail()` return the start resp. end of the list.
- `Append()` and `Prepend()` insert a new node "left" or "right" of a given anchor.
- `Delete()` removes a given node.
- `FindNodes()` returns a list of nodes that match a value.
- `SetValue()` modifies the value of a given node.

### Don't break the encapsulation: use `llist`'s methods

Package `llist` doesn't hide that it internally uses `lnode`s. But. In the event that some of `lnode`s methods are used (or fields are directly modified) instead of `llist`'s alternatives, then the internal administration of `llist` can be rebuilt:

- `FixHead()` and `FixTail()` recompute the head and tail nodes
- `FixCounts()` recomputes the internal map that's used by `FindNodes()`.

```go
// Example: changing a node's value directly
for _, node := range l.FindNodes("lazy") {
    // Directly change the value, not using llist's method l.SetValue(node, "quick")
    node.Value = "quick"
}
// l.FindNodes("quick") won't return correct results
l.FixCounts() // Rebuild the internal map
// l.FindNodes("quick") will now return correct results

// Example: inserting at head without telling llist
// Instead, it would be better to: l.Head().Prepend(lnode.New[string]("yesterday"))
node := lnode.New[string]("yesterday")
l.Head().Prev = node
node.Next = l.Head()
// l.Head() will still return the old head, not the "yesterday" node
l.FixHead() // Recompute the head
// l.Head() will now correctly return the chain start
```
