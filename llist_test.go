package llist

import (
	"fmt"
	"testing"

	"github.com/KarelKubat/lnode"
)

func mkList() *LList[int] {
	l := New[int]()
	for i := range 10 {
		l.Append(l.Tail(), lnode.New[int](i))
	}
	return l
}

func TestAppend(t *testing.T) {
	l := mkList()
	if l.Head().Value != 0 {
		t.Errorf("TestAppend: got head value %d, want 0", l.Head().Value)
	}
	if l.Tail().Value != 9 {
		t.Errorf("TestAppend: got tail value %d, want 9", l.Tail().Value)
	}
}

func TestPrepend(t *testing.T) {
	l := mkList()
	for i := -1; i < -10; i-- {
		l.Head().Prepend(lnode.New[int](i))
	}

	want := -0
	for n := l.Head(); n != nil; n = n.Next {
		if n.Value != want {
			t.Errorf("TestAppend: while iterating got value %d, want %d", n.Value, want)
		}
		want++
	}
}

func TestEmptyPrepend(t *testing.T) {
	l := New[int]()
	for i := range 10 {
		l.Prepend(l.Head(), lnode.New[int](i))
	}
	want := 9
	for n := l.Head(); n != nil; n = n.Next {
		if n.Value != want {
			t.Errorf("TestEmptyPrepend: while iterating got value %d, want %d", n.Value, want)
		}
		want--
	}
}

func TestFixHead(t *testing.T) {
	l := mkList()
	hd := l.Head()

	hd.Prepend(lnode.New[int](-1))
	hd.Prev.Prepend(lnode.New[int](-2))
	hd.Prev.Prev.Prepend(lnode.New[int](-3))

	if v := l.Head().Value; v != 0 {
		t.Errorf("TestFixHead: before fixing Head, value is %d, want 0", v)
	}
	l.FixHead()
	if v := l.Head().Value; v != -3 {
		t.Errorf("TestFixHead: after fixing Head, value is %d, want 3", v)
	}
}

func TestFixTail(t *testing.T) {
	l := mkList()
	tl := l.Tail()

	tl.Append(lnode.New[int](10))
	tl.Next.Append(lnode.New[int](11))
	tl.Next.Next.Append(lnode.New[int](12))

	if v := l.Tail().Value; v != 9 {
		t.Errorf("TestFixTail: before fixing Tail, value is %d, want 9", v)
	}
	l.FixTail()
	if v := l.Tail().Value; v != 12 {
		t.Errorf("TestFixTail: after fixing Tail, value is %d, want 12", v)
	}
}

func TestFindNodes(t *testing.T) {
	l := mkList()

	// Simple list 0 - 1 - 2 - ... 9
	for i := range 10 {
		nds := l.FindNodes(i)
		if len(nds) != 1 {
			t.Errorf("FindNodes(%d) = %v, want one result", i, nds)
			continue
		}
		if v := nds[0].Value; v != i {
			t.Errorf("FindNodes(%d)[0].Value = %d, want %d", i, v, i)
		}
	}

	// Add another node with 5
	l.Append(l.Head(), lnode.New[int](5))
	nds := l.FindNodes(5)
	if len(nds) != 2 {
		t.Errorf("FindNodes(5) = %v, want two results after duplication of '5'", nds)
	}
	for _, n := range nds {
		if n.Value != 5 {
			t.Errorf("FindNodes(5) = %v, want two results both with value 5", nds)
		}
	}
}

func TestDelete(t *testing.T) {
	l := New[string]()
	for _, s := range []string{"the", "quick", "brown", "fox"} {
		l.Append(l.Tail(), lnode.New[string](s))
	}

	brown := l.Head().Next.Next
	if brown.Value != "brown" {
		t.Fatalf("TestDelete: node brown initialized to %s, need 'brown'", brown.Value)
	}

	// Add more browns at random locations
	l.Append(l.Head(), lnode.New[string]("brown"))
	l.Append(l.Tail(), lnode.New[string]("brown"))
	l.Prepend(brown, lnode.New[string]("brown"))

	nds := l.FindNodes("brown")
	if len(nds) != 4 {
		t.Fatalf("TestDelete: found %d occurrences of 'brown', need 4", len(nds))
	}

	// Delete all occurrences
	wantCount := 4
	for _, n := range nds {
		l.Delete(n)
		wantCount--
		if gotCount := len(l.FindNodes("brown")); gotCount != wantCount {
			t.Errorf("TestDelete: want count %d after deletion(s), got %d", wantCount, gotCount)
		}
	}

	// Delete beyond
	for range 10 {
		l.Delete(brown)
		if gotCount := len(l.FindNodes("brown")); gotCount != 0 {
			t.Errorf("TestDelete: deleting beyond existence yields count %d, want 0", gotCount)
		}
	}

	l.Head().VisitByNext(func(node *lnode.Node[string]) bool {
		fmt.Println(node.Value)
		return true
	})
}
