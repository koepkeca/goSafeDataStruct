package main

import (
	"log"
	//	"math/rand"
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
	//Example #1, Implement a new heap, push something onto it, then iterate through it.
	basic := safeHeap.New(&IntHeap{9, 3, 17, 6, 1})
	basic.Push(147)
	for basic.Len() != 0 {
		log.Printf("The next value is: %d", basic.Pop().(int))
	}
	//Example #2, Non-pooled worker concurrency test creates a number of workers that wait a random number
	//of milliseconds to create a random distribution. Once these workers wait
	//concludes, they add themselves to the min_head and return.
	np_start := time.Now()
	log.Printf("Beginning concurrency test...")
	concurr_heap := safeHeap.New(&IntHeap{})
	wg := sync.WaitGroup{}
	for n := 0; n < 4096; n++ {
		wg.Add(1)
		go func(i int, wg *sync.WaitGroup) {
			select {
			//			case <- time.After(time.Duration(rand.Intn(1000)) * time.Millisecond):
			case <-time.After(time.Duration(20) * time.Millisecond):
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
	np_end := time.Since(np_start)
	p_start := time.Now()
	worker_pool_heap := safeHeap.New(&IntHeap{})
	const num_jobs = 4096
	const workers = 128
	jobs := make(chan int, num_jobs)
	rslt := make(chan int, num_jobs)
	for w := 1; w <= workers; w++ {
		go worker(w, jobs, rslt)
	}
	for j := 1; j <= num_jobs; j++ {
		jobs <- j
	}
	close(jobs)
	for r := 1; r <= num_jobs; r++ {
		worker_pool_heap.Push(<-rslt)
	}
	close(rslt)
	for worker_pool_heap.Len() != 0 {
		log.Printf("The next value in heap is : %d", worker_pool_heap.Pop().(int))
	}
	p_end := time.Since(p_start)
	log.Printf("Sequential execution time: %s, Concurrent execution time: %s", np_end, p_end)
	worker_pool_heap.Destroy()
	basic.Destroy()
	return
}

func worker(id int, jobs <-chan int, rslt chan<- int) {
	for next := range jobs {
		log.Printf("Worker %d processing job %d", id, next)
		time.Sleep(time.Duration(20) * time.Millisecond)
		//		time.Sleep(time.Duration(rand.Intn(20)) * time.Millisecond)
		log.Printf("Inserting value %d to heap", next)
		rslt <- next
	}
	return
}
