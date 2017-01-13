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
		lChan <- curr.size
	}
	return <-lChan
}

//Dequeue removes the next item in the queue and returns it's value
func (s *SafeQueue) Dequeue() interface{} {
	vChan := make(chan interface{})
	s.op <- func(curr *queue) {
		if curr.size == 0 {
			vChan <- nil
			return
		}
		val := curr.head.value
		if curr.size == 1 {
			curr.tail = nil
		}
		curr.head = curr.head.next
		curr.size--
		vChan <- val
		return
	}
	return <-vChan
}

//Enqueue places a new item at the end of the queue
func (s *SafeQueue) Enqueue(v interface{}) {
	s.op <- func(curr *queue) {
		newElem := &elem{v, nil}
		if curr.size == 0 {
			curr.head = newElem
		} else {
			curr.tail.next = newElem
		}
		curr.tail = newElem
		curr.size++
		return
	}
}

//Front returns the value from the front of the queue.
//It does not remove it from the queue.
func (s *SafeQueue) Front() interface{} {
	vChan := make(chan interface{})
	s.op <- func(curr *queue) {
		if curr.size == 0 {
			vChan <- nil
			return
		}
		vChan <- curr.head.value
		return
	}
	return <-vChan
}

//Back returns the value from the back of the queue.
//It does not remove it from the queue.
func (s *SafeQueue) Back() interface{} {
	vChan := make(chan interface{})
	s.op <- func(curr *queue) {
		if curr.size == 0 {
			vChan <- nil
			return
		}
		vChan <- curr.tail.value
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

//stack is the basic container for the queue
type queue struct {
	head *elem
	tail *elem
	size int64
}

//elem is the element structure
type elem struct {
	value interface{}
	next  *elem
}

//loop creates the guarded data structure and listens for
//methods on the op channel. loop terminates when the op
//channel is closed.
func (s *SafeQueue) loop() {
	st := &queue{}
	for op := range s.op {
		op(st)
	}
}
