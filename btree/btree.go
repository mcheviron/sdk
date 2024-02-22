package btree

import (
	"cmp"
)

type BTree[K cmp.Ordered, V any] struct {
	root *node[K, V]
	b    uint
	// len   uint
	// depth uint
}

func (t *BTree[K, V]) Find(key K) (*V, bool) {
	curr := t.root
	for {
		pos, found := curr.search(key)
		if found {
			return &curr.vals[pos], true
		}
		if curr.edges[pos] == nil {
			return nil, false
		}
		curr = curr.edges[pos]
	}
}

func (t *BTree[K, V]) Insert(key K, val V) {

}
