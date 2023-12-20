package pq

import (
	"container/heap"
)

type Item[T any] struct {
	value    T
	priority int
}

type PriorityQueue[T any] []Item[T]

func (pq PriorityQueue[_]) Len() int {
	return len(pq)
}

func (pq PriorityQueue[_]) Less(i, j int) bool {
	return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue[_]) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue[T]) Push(x any) {
	*pq = append(*pq, x.(Item[T]))
}

func (pq *PriorityQueue[_]) Pop() (x any) {
	x, *pq = (*pq)[len(*pq)-1], (*pq)[:len(*pq)-1]
	return x
}

func (pq *PriorityQueue[T]) Insert(v T, p int) {
	heap.Push(pq, Item[T]{v, p})
}

func (pq *PriorityQueue[T]) DeleteMin() (T, int) {
	x := heap.Pop(pq).(Item[T])
	return x.value, x.priority
}

func (pq PriorityQueue[_]) IsEmpty() bool {
	return pq.Len() == 0
}
