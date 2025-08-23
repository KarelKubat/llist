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
		l.SetValue(node, "quick")
	}
	print("Lazy is now quick", l)
}
