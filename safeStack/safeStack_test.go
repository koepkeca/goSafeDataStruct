package safeStack

import (
	"testing"
)

func TestStackCreation(t *testing.T) {
	s := NewSafeStack()
	if s.Len() != 0 {
		t.Errorf("Failed, invalid stack length.")
	}
	s.Destroy()
	return
}

func TestStackPushLength(t *testing.T) {
	s := NewSafeStack()
	s.Push(14)
	s.Push(42)
	s.Push("testing")
	s.Push([]byte("Viper"))
	len := s.Len()
	if len != 4 {
		t.Errorf("Failed, invalid stack length, got %d expected 4", len)
	}
	s.Destroy()
	return
}

func TestStackPushOrder(t *testing.T) {
	var nv = 0
	var ok = false
	s := NewSafeStack()
	s.Push(16)
	s.Push(32)
	s.Push(64)
	if nv, ok = s.Pop().(int); !ok {
		t.Errorf("Failed, Pop got wrong type")
		s.Destroy()
		return
	}
	if nv != 64 {
		t.Errorf("Failed, got incorrect value order")
		s.Destroy()
		return
	}
	s.Destroy()
	return
}

func TestSizeAfterPop(t *testing.T) {
	s := NewSafeStack()
	s.Push(16)
	s.Push("test")
	s.Push("私は笑い男だ")
	_ = s.Pop()
	_ = s.Pop()
	_ = s.Pop()
	if s.Len() != 0 {
		t.Errorf("Failed, poped through entire stack, yet size is non-zero")
	}
	s.Destroy()
	return
}

func TestEmptyPop(t *testing.T) {
	s := NewSafeStack()
	v := s.Pop()
	if v != nil {
		t.Errorf("Empty Pop got non-nil value")
	}
	s.Destroy()
}

func TestEmptyPopWithValues(t *testing.T) {
	s := NewSafeStack()
 	s.Push("Thingy")
	_ = s.Pop()
	v := s.Pop()
	if v != nil {
		t.Errorf("Empty stack with values got non-nil value")
	}
	s.Destroy()
}
