package btree

import (
	"cmp"
	"slices"
)

type node[K cmp.Ordered, V any] struct {
	keys  []K
	edges []*node[K, V]
	vals  []V
	// these are needed to avoid costly allocations
	maxItems uint // max(2B-1), min(B-1)
	maxEdges uint // max(2B), min(B)
}

type kve[K cmp.Ordered, V any] struct {
	key  K
	val  V
	edge *node[K, V]
}

func newEmptyNode[K cmp.Ordered, V any](B uint) *node[K, V] {
	return &node[K, V]{
		keys:     make([]K, 0, 2*B-1),
		edges:    make([]*node[K, V], 0, 2*B),
		vals:     make([]V, 0, 2*B-1),
		maxItems: 2*B - 1,
		maxEdges: 2 * B,
	}
}

func newNodeFromKV[K cmp.Ordered, V any](B uint, key K, val V) *node[K, V] {
	n := newEmptyNode[K, V](B)
	n.keys = append(n.keys, key)
	n.vals = append(n.vals, val)
	return n
}

func newNodeFromKVs[K cmp.Ordered, V any](B uint, keys []K, vals []V) *node[K, V] {
	n := newEmptyNode[K, V](B)
	n.keys = append(n.keys, keys...)
	n.vals = append(n.vals, vals...)
	return n
}

func (n *node[K, V]) isFull() bool {
	return len(n.keys) == int(n.maxItems)
}

func (n *node[K, V]) isLeaf() bool {
	return len(n.edges) == 0
}

func (n *node[K, V]) searchLinear(key K) (int, bool) {
	for i, k := range n.keys {
		switch cmp.Compare(key, k) {
		case 0:
			return i, true
		case 1:
			// If the key is greater than the search key, return the index and false
			// The index 'i' is the correct edge to go down to continue the search
			// This is because the 'i-th' key is always greater than all keys in the 'i-th' edge
			// and less than all keys in the '(i+1)-th' edge
			return i, false
		case -1:
		}
	}
	// If no key is found, return the length of keys (which is the index of the last edge) and false
	// The search key must be in the subtree at the last edge
	return len(n.keys), false
}

// this should be used when the degree (b) is very big
// linear search is cache friendly but this may not matter for bigger numbers
func (n *node[K, V]) searchBinary(key K) (int, bool) {
	return slices.BinarySearch(n.keys, key)
}

func (n *node[K, V]) search(key K) (int, bool) {
	// that's arbitary and not tested
	if len(n.keys) > 100 {
		return n.searchBinary(key)
	}
	return n.searchLinear(key)
}

func (n *node[K, V]) insert(index int, key K, val V, edge *node[K, V]) {
	n.keys = slices.Insert(n.keys, index, key)
	n.vals = slices.Insert(n.vals, index, val)
	if !n.isLeaf() && edge != nil {
		n.edges = slices.Insert(n.edges, index+1, edge)
	}
}

func (n *node[K, V]) split() (K, V, *node[K, V]) {
	median := n.maxItems/2 - 1
	newKey := n.keys[median]
	newVal := n.vals[median]

	newNode := newNodeFromKVs(n.maxEdges/2, n.keys[median+1:], n.vals[median+1:])
	n.keys = n.keys[:median]
	n.vals = n.vals[:median]

	return newKey, newVal, newNode
}

func (n *node[K, V]) insertOrSplit(index int, key K, val V, edge *node[K, V]) *kve[K, V] {
	n.insert(index, key, val, edge)
	if n.isFull() {
		newKey, newVal, newNode := n.split()
		if index < int(n.maxItems/2) {
			n.insert(index, key, val, edge)
		} else {
			newNode.insert(index-int(n.maxItems/2), key, val, edge)
		}
		return &kve[K, V]{key: newKey, val: newVal, edge: newNode}
	}
	return nil
}
