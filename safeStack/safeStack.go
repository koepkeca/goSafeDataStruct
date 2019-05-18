//Package safeStack implements a thread safe stack in go.
package safeStack

//SafeStack is the structure that contains the channel used to
//communicate with the stack.
type SafeStack struct {
	op chan (func(*stack))
}

//Len will get return the number of items in the stack.
func (s *SafeStack) Len() (i int64) {
	lChan := make(chan int64)
	s.op <- func(curr *stack) {
		lChan <- int64(len(*curr))
	}
	return <-lChan
}

//Pop will perform a pop on the stack, removing the first item
//and returning it's value.
func (s *SafeStack) Pop() (v interface{}) {
	vChan := make(chan interface{})
	s.op <- func(curr *stack) {
		old := *curr
		n := len(old)
		if n == 0 {
			vChan <- nil
			return
		}
		item := old[0]
		*curr = old[1:n]
		vChan <- item
		return
	}
	return <-vChan
}

//Push will push the value v onto the stack.
func (s *SafeStack) Push(v interface{}) {
	s.op <- func(curr *stack) {
		*curr = append([]interface{}{v}, *curr...)
		return
	}
	return
}

//Destroy closes the primary channel thus stopping
//the running go-routine.
func (s *SafeStack) Destroy() {
	close(s.op)
}

//New creates a new Safe Stack, this also starts the go-routine
//so once this is called, you need to clean up after yourself
//by using the Destroy method.
func New() (s *SafeStack) {
	s = &SafeStack{make(chan func(*stack))}
	go s.loop()
	return
}

//We emulate a stack using an interface slice to reduce memory overhead
type stack []interface{}

//loop creates the guarded data structure and listens for
//methods on the op channel. loop terminates when the op
//channel is closed.
func (s *SafeStack) loop() {
	st := &stack{}
	for op := range s.op {
		op(st)
	}
}
