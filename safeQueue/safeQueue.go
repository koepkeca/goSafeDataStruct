//Package safeQueue implements a thread safe queue in go.
package safeQueue

//SafeQueue is the structure that contains the channel used to
//communicate with the queue.
type SafeQueue struct {
	op chan (func(*queue))
}

//Len will get return the number of items in the queue.
func (s *SafeQueue) Len() int64 {
	lChan := make(chan int64)
	s.op <- func(curr *queue) {
		lChan <- int64(len(*curr))
	}
	return <-lChan
}

//Dequeue removes the next item in the queue and returns it's value
func (s *SafeQueue) Dequeue() interface{} {
	vChan := make(chan interface{})
	s.op <- func(curr *queue) {
		if len(*curr) == 0 {
			vChan <- nil
			return
		}
		old := *curr
		n := len(old)
		val := old[0]
		*curr = old[1:n]
		vChan <- val
		return
	}
	return <-vChan
}

//Enqueue places a new item at the end of the queue
func (s *SafeQueue) Enqueue(v interface{}) {
	s.op <- func(curr *queue) {
		*curr = append(*curr, v)
		return
	}
}

//Front returns the value from the front of the queue.
//It does not remove it from the queue.
func (s *SafeQueue) Front() interface{} {
	vChan := make(chan interface{})
	s.op <- func(curr *queue) {
		if len(*curr) == 0 {
			vChan <- nil
			return
		}
		tmp := *curr
		vChan <- tmp[0]
		return
	}
	return <-vChan
}

//Back returns the value from the back of the queue.
//It does not remove it from the queue.
func (s *SafeQueue) Back() interface{} {
	vChan := make(chan interface{})
	s.op <- func(curr *queue) {
		if len(*curr) == 0 {
			vChan <- nil
			return
		}
		tmp := *curr
		vChan <- tmp[len(tmp)]
		return
	}
	return <-vChan
}

//Destroy closes the primary channel thus stopping
//the running go-routine.
func (s *SafeQueue) Destroy() {
	close(s.op)
}

//New creates a new Safe Stack, this also starts the go-routine
//so once this is called, you need to clean up after yourself
//by using the Destroy method.
func New() (s *SafeQueue) {
	s = &SafeQueue{make(chan func(*queue))}
	go s.loop()
	return
}

//queue contains the queue represented as a slice of interfaces.
type queue []interface{}

//loop creates the guarded data structure and listens for
//methods on the op channel. loop terminates when the op
//channel is closed.
func (s *SafeQueue) loop() {
	st := &queue{}
	for op := range s.op {
		op(st)
	}
}
