package safeHeap

import (
	"container/heap"
)

//SafeHeap creates a co-routine safe Heap comprised a data type that implements heap.Interface
type SafeHeap struct {
	op chan (func(heap.Interface))
}

//Len provides the current length of the Heap
func (s *SafeHeap) Len() (l int) {
	lch := make(chan int)
	s.op <- func(curr heap.Interface) {
		lch <- curr.Len()
		return
	}
	return <-lch
}

//Push will insert the item into the heap 
func (s *SafeHeap) Push(x interface{}) bool {
	bch := make(chan bool)
	s.op <- func(curr heap.Interface) {
		heap.Push(curr, x)
		bch <- true
		return
	}
	return <- bch
}

//Pop removes the next value fom the heap
func (s *SafeHeap) Pop() interface{} {
	ich := make(chan interface{})
	s.op <- func(curr heap.Interface) {
		ich <- heap.Pop(curr)
		return
	}
	return <-ich
}

//Remove will remove the item and idx from the heap
func (s *SafeHeap) Remove(idx int) interface{} {
	ich := make(chan interface{})
	s.op <- func(curr heap.Interface) {
		ich <- heap.Remove(curr, idx)
		return
	}
	return <-ich
}

//New creates a new heap that is safe for concurrent use from h
func New(h heap.Interface) (s *SafeHeap) {
	s = &SafeHeap{op: make(chan func(heap.Interface))}
	go s.loop(h)
	return
}

func (s *SafeHeap) loop(h heap.Interface) {
	heap.Init(h)
	for op := range s.op {
		op(h)
	}
	return
}
