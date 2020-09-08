package main

import (
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/koepkeca/goSafeDataStruct/safeHeap"
)

//Create a InitHeap (minimum integer heap) which implements container/heap.Interface
//this will allow us to create a concurrent-ready Min Heap
type IntHeap []int

func (h IntHeap) Len() int { return len(h) }

//opposite for max heap
func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *IntHeap) Pop() (x interface{}) {
	old := *h
	n := len(old)
	x = old[n-1]
	*h = old[0 : n-1]
	return

}

func (h *IntHeap) Push(x interface{}) {
	*h = append(*h, x.(int))
	return
}

func main() {
	basic := safeHeap.New(&IntHeap{9, 3, 17, 6, 1})
	basic.Push(147)
	for basic.Len() != 0 {
		log.Printf("The next value is: %d", basic.Pop().(int))
	}
	//Concurrency test creates a number of workers that wait a random number
	//of milliseconds to create a random distribution. Once these workers wait
	//concludes, they add themselves to the min_head and return.
	log.Printf("Beginning concurrency test...")
	concurr_heap := safeHeap.New(&IntHeap{})
	wg := sync.WaitGroup{}
	for n := 0; n < 32; n++ {
		wg.Add(1)
		go func(i int, wg *sync.WaitGroup) {
			select {
			case <- time.After(time.Duration(rand.Intn(1000)) * time.Millisecond):
				defer wg.Done()
				concurr_heap.Push(i)
				log.Printf("Added value %d", i)
				return
			}
		}(n, &wg)
	}
	wg.Wait()
	log.Printf("Heap length: %d", concurr_heap.Len())
	log.Printf("Heap concurrent heap contents: ")
	for concurr_heap.Len() != 0 {
		log.Printf("The next value is: %d", concurr_heap.Pop().(int))
	}
	return
}
