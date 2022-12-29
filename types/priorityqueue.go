package types

import "container/heap"

/** copied entirely from https://pkg.go.dev/container/heap */

// An Item is something we manage in a priority queue.
type Item[T any] struct {
	value *T
	// The priority of the item in the queue.
	priority int
	// The index is needed by update and is maintained by the heap.Interface methods.
	// The index of the item in the heap.
	index int
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue[T any] []*Item[T]

// sets index manually
func (pq *PriorityQueue[T]) NewItem(value *T, priority, index int) {
	(*pq)[index] = &Item[T]{
		value,
		priority,
		index,
	}
}

// sets index automatically
func (pq *PriorityQueue[T]) PushValue(value *T, priority int) {
	newItem := &Item[T]{
		value,
		priority,
		0,
	}

	heap.Push(pq, newItem)
}

func (pq PriorityQueue[T]) Len() int { return len(pq) }

func (pq PriorityQueue[T]) Less(i, j int) bool {
	// We want Pop to give us the LOWEST priority so we use less than here.
	return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue[T]) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue[T]) Push(x interface{}) {
	n := len(*pq)
	item := x.(*Item[T])
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue[T]) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
func (pq *PriorityQueue[T]) Update(value *T, priority int) {
	// find item
	for _, item := range *pq {
		if item.value == value {
			item.priority = priority
			heap.Fix(pq, item.index)

			return
		}
		// else: uh oh
	}
}

func (pq *PriorityQueue[T]) Get() *T {
	return heap.Pop(pq).(*Item[T]).value
}

func (pq *PriorityQueue[T]) Init() {
	heap.Init(pq)
}
