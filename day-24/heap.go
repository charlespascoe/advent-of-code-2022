package main

type HeapItem[T any] interface {
	LessThan(other T) bool
}

type Heap[T HeapItem[T]] struct {
	items []T
}

func (h *Heap[T]) Count() int {
	return len(h.items)
}

func (h *Heap[T]) Add(item T) {
	h.items = append(h.items, item)
	i := len(h.items) - 1
	h.heapifyUp(i)
}

func (h *Heap[T]) PopMin() T {
	if len(h.items) == 0 {
		return *new(T)
	}

	item := h.items[0]
	last := len(h.items) - 1

	h.swap(0, last)

	// Remove the old cell from the end of the heap and indexes
	h.items = h.items[:last]

	if len(h.items) > 0 {
		// Move the swapped session to where it should be in the heap
		h.heapifyDown(0)
	}

	return item
}

func (h *Heap[T]) heapifyUp(i int) {
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

func (h *Heap[T]) heapifyDown(i int) {
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

func (h *Heap[T]) swap(i, j int) {
	h.items[i], h.items[j] = h.items[j], h.items[i]
}

func parentIndex(i int) int {
	return (i - 1) / 2
}

func childIndexes(i int) (left, right int) {
	return 2*(i+1) - 1, 2 * (i + 1)
}
