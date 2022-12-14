package main

type HeapItem[T any, K comparable] interface {
	Key() K
	LessThan(other T) bool
}

type Heap[T HeapItem[T, K], K comparable] struct {
	items   []T
	indexes map[K]int
}

func (h *Heap[T, K]) Count() int {
	return len(h.items)
}

func (h *Heap[T, K]) Add(item T) {
	if h.indexes == nil {
		h.indexes = make(map[K]int)
	}

	i, exists := h.indexes[item.Key()]

	if exists {
		// Move it to the right place if it's changed up or down
		h.heapifyUp(i)
		h.heapifyDown(i)
	} else {
		h.items = append(h.items, item)
		i = len(h.items) - 1
		h.indexes[item.Key()] = i
		h.heapifyUp(i)
	}
}

func (h *Heap[T, K]) PopMin() T {
	if len(h.items) == 0 {
		return *new(T)
	}

	c := h.items[0]
	last := len(h.items) - 1

	h.swap(0, last)

	// Remove the old cell from the end of the heap and indexes
	h.items = h.items[:last]
	delete(h.indexes, c.Key())

	if len(h.items) > 0 {
		// Move the swapped session to where it should be in the heap
		h.heapifyDown(0)
	}

	return c
}

func (h *Heap[T, K]) heapifyUp(i int) {
	for i > 0 {
		parent := parentIndex(i)

		if !h.items[i].LessThan(h.items[parent]) {
			// Min-heap condition restored; stop here
			break
		}

		h.swap(parent, i)
		i = parent
	}
}

func (h *Heap[T, K]) heapifyDown(i int) {
	for {
		left, right := childIndexes(i)
		min := i

		if left < len(h.items) && h.items[left].LessThan(h.items[min]) {
			min = left
		}

		if right < len(h.items) && h.items[right].LessThan(h.items[min]) {
			min = right
		}

		if min == i {
			break
		}

		h.swap(i, min)
		i = min
	}
}

func (h *Heap[T, K]) swap(i, j int) {
	// Get current items
	ci := h.items[i]
	cj := h.items[j]

	// Swap items
	h.items[i] = cj
	h.items[j] = ci

	// Update indexes (cj now at i, ci now at j)
	h.indexes[cj.Key()] = i
	h.indexes[ci.Key()] = j
}

func parentIndex(i int) int {
	return (i - 1) / 2
}

func childIndexes(i int) (left, right int) {
	return 2*(i+1) - 1, 2 * (i + 1)
}
