/*
Package llist encapsulates a double linked list based on github.com/KarelKubat/lnode and provides handy methods plus lookups by value. All methods run in O(1) time complexity (e.g., finding the head or tail, and looking up nodes by value).
*/
package llist

import (
	"github.com/KarelKubat/lnode"
)

// LList is the receiver.
type LList[V comparable] struct {
	head, tail *lnode.Node[V]         // Head/tail for fast lookup
	nodes      map[V][]*lnode.Node[V] // Nodes keyed by value for lookup
}

/*
New constructs an LList receiver. Example of a linked list for strings:

	l := llist.New[string]()  // Linked list containing strings
*/
func New[V comparable]() *LList[V] {
	return &LList[V]{
		nodes: map[V][]*lnode.Node[V]{},
	}
}

/*
Head returns the "start" of the chain, or nil when (a) the chain is empty, or (b) the chain is circular. Head runs in O(1) time (in contrast to github.com/KarelKubat/lnode).

Example:

	l := New[string]()
	for _, s := range []string{"the", "quick", "brown", "fox"} {
		l.Append(l.Tail(), lnode.New[string](s))
	}
	fmt.Println(l.Head().Value)  // "the"

The head node is also a typical start to traverse the list from "left" to "right":

	// see llnode.VisitByNext()
	l.Head().VisitByNext(func (node *llnode.Node[string]) {
		fmt.Printf("%s ", node.Value)
	})
	// Output: the quick brown fox
*/
func (l *LList[V]) Head() *lnode.Node[V] {
	return l.head
}

/*
Tail returns the "end" of the chain, or nil when (a) the chain is empty, or (b) the chain is circlular. Tail runs in O(1) time (in contrast to github.com/KarelKubat/lnode).

Example:

	l := New[string]()
	for _, s := range []string{"the", "quick", "brown", "fox"} {
		l.Append(l.Tail(), lnode.New[string](s))
	}
	fmt.Println(l.Tail().Value)  // "fox"

The tail node is also a typical start to travers the list in reverse order, "right to left":

	l.Tail().VisitByPrev(...) // Details of the visiting function not shown
*/
func (l *LList[V]) Tail() *lnode.Node[V] {
	return l.tail
}

/*
Append inserts the named node "after" (or "to the right") of the anchor and, if necessary adjusts the tail value. Either Append() or Prepend() can be used to initialize an empty LList.

Note that using the node's Append() (lnode.Append()) bypasses the internal administration of llist. The head or tail position, and internal pointers (used in llist.FindNodes()) will not be updated. Re-fixing these can be done using FixHead(), FixTail() and FixCounts().

Append() and Tail() can be used to build a list from left to right:

	l := New[string]()
	for _, s := range []string{"the", "quick", "brown", "fox"} {
		l.Append(l.Tail(), lnode.New[string](s))
	}
	fmt.Println(l.Head().Value)  // "the"
	fmt.Println(l.Tail().Value)  // "fox"
*/
func (l *LList[V]) Append(anchor, n *lnode.Node[V]) {
	if l.head == nil {
		l.head = n
		l.tail = n
	} else if anchor == l.tail {
		l.tail.Append(n)
		l.tail = n
	} else {
		anchor.Append(n)
	}
	l.addCount(n)
}

/*
Prepend inserts the named node "before" (or "to the left") of the anchor, and, if necessary adjust the head value. Either Append() or Prepend() can be used to initialize an empty LList.

Note that using the node's Append() (lnode.Append()) bypasses the internal administration of llist. The head or tail position, and internal pointers (used in llist.FindNodes()) will not be updated. Re-fixing these can be done using FixHead(), FixTail() and FixCounts().

Prepend() and Head() can be used to build a list from right to left:

	l := New[string]()
	for _, s := range []string{"the", "quick", "brown", "fox"} {
		l.Prepend(l.Head(), lnode.New[string](s))
	}
	fmt.Println(l.Head().Value)  // "fox"
	fmt.Println(l.Tail().Value)  // "the"
*/
func (l *LList[V]) Prepend(anchor, n *lnode.Node[V]) {
	switch {
	case l.head == nil && l.tail == nil:
		l.head = n
		l.tail = n
	case anchor == l.head:
		l.head.Prepend(n)
		l.head = n
	default:
		anchor.Prepend(n)
	}
	l.addCount(n)
}

/*
FixHead recomputes the stored value for the head of the list. This may be necessary when the list is modified using llnode.Append(), llnode.Prepend() etc., instead of the corresponding llist functions.

FixHead() doesn't work on circular lists.

Example:

	l := New[string]()
	for _, s := range []string{"the", "quick", "brown", "fox"} {
		l.Append(l.Tail(), lnode.New[string](s))
	}

	fmt.Println(l.Head().Value) 					// "the"

	hd := l.Head()
	hd.Prepend(lnode.New[string]("today")			// inserted via node, not list
	fmt.Println(l.Head().Value) 					// still "the"
	hd.FixHead()                					// recompute
	fmt.Println(l.Head().Value) 					// now it's "today"

	l.Prepend(l.Head(),
		lnode.New[string]("yesterday and"))			// inserted via list, not node
	fmt.Println(l.Head().Value)						// "yesterday and", which is correct
*/
func (l *LList[V]) FixHead() {
	l.head = l.head.Head()
}

/*
FixTail recomputes the stored value for the tail of the list. This may be necessary when the list is modified using llnode.Append(), llnode.Prepend() etc., instead of the corresponding llist functions. For an example see FixHead().

FixTail() doesn't work on circular lists.
*/
func (l *LList[V]) FixTail() {
	l.tail = l.tail.Tail()
}

/*
FixCounts recomputes the stored pointers to kept nodes. This may be necessary when the list is modified using llnode's functions instead of the corresponding llist functions.

FixCounts() doesn't work on circular lists.
*/
func (l *LList[V]) FixCounts() {
	hd := l.Head()
	if hd == nil {
		return
	}
	l.nodes = map[V][]*lnode.Node[V]{}
	hd.VisitByNext(func(node *lnode.Node[V]) bool {
		l.addCount(node)
		return true
	})
}

/*
FindNodes returns a slice of nodes matching the stated value. The length of the slice may be 0 when the value doesn't occur in the linked list, or 1 when it is unique, and so on.

Example:

	l := New[string]()
	for _, s := range []string{"the", "quick", "brown", "fox"} {
		l.Append(l.Tail(), lnode.New[string](s))
	}
	nds := l.FindNodes("dog")						// search for "dog"
	fmt.Println(length(nds))						// occurs 0 times

	nds = l.FindNodes("fox")						// search for "fox"
	fmt.Println(length(nds))						// occurs 1 time

	l.Prepend(l.Head(), lnode.New[string]("fox"))	// insert another fox
	nds = l.FindNodes("fox")						// search again for "fox"
	fmt.Println(length(nds))						// occurs 2 times
	for _, n := range nds {
		fmt.Println(n.Value)						// "fox", obviously
	}
*/
func (l *LList[V]) FindNodes(v V) []*lnode.Node[V] {
	nodes, ok := l.nodes[v]
	if !ok {
		return nil
	}
	return nodes
}

/*
Delete removes nodes from the list, adjusting head/tail as appropriate. Example:

	l := New[string]()
	for _, s := range []string{"the", "quick", "brown", "fox"} {
		l.Append(l.Tail(), lnode.New[string](s))
	}

	brown := l.Head().Next.Next 	// Or in this case, l.FindNodes("brown")[0]
	l.Delete(brown)
	l.Head().VisitByNext(func(node *lnode.Node[string]) bool {
		fmt.Printf("%s ", node.Value)
	}
	// Output: the quick fox
*/
func (l *LList[V]) Delete(node *lnode.Node[V]) {
	switch node {
	case l.head:
		l.head = node.Next
	case l.tail:
		l.tail = node.Prev
	}
	node.Delete()
	l.subCount(node)
}

/*
SetValue changes the value of a node and updates the internal structure so that FindNodes() will properly return nodes that contain a value.

Example:

	l := New[string]()
	for _, s := range []string{"the", "quick", "brown", "fox"} {
		l.Append(l.Tail(), lnode.New[string](s))
	}
	// Head -> next -> next points to the node containing "brown"
	l.SetValue(l.Head().Next.Next, "red")
*/
func (l *LList[V]) SetValue(node *lnode.Node[V], value V) {
	l.subCount(node)
	node.Value = value
	l.addCount(node)
}

// Helper
func (l *LList[V]) addCount(n *lnode.Node[V]) {
	if _, ok := l.nodes[n.Value]; !ok {
		l.nodes[n.Value] = []*lnode.Node[V]{}
	}
	l.nodes[n.Value] = append(l.nodes[n.Value], n)
}

// Helper
func (l *LList[V]) subCount(n *lnode.Node[V]) {
	nds, ok := l.nodes[n.Value]
	if !ok {
		return
	}
	newNds := []*lnode.Node[V]{}
	for _, nd := range nds {
		if nd != n {
			newNds = append(newNds, nd)
		}
	}
	l.nodes[n.Value] = newNds
}
