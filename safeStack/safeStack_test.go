package safeStack

import (
	"testing"
)

func TestStackCreation(t *testing.T) {
	s := New()
	if s.Len() != 0 {
		t.Errorf("Failed, invalid stack length.")
	}
	s.Destroy()
	return
}

func TestStackPushLength(t *testing.T) {
	s := New()
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
	s := New()
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
	s := New()
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
	s := New()
	v := s.Pop()
	if v != nil {
		t.Errorf("Empty Pop got non-nil value")
	}
	s.Destroy()
}

func TestEmptyPopWithValues(t *testing.T) {
	s := New()
	s.Push("Thingy")
	_ = s.Pop()
	v := s.Pop()
	if v != nil {
		t.Errorf("Empty stack with values got non-nil value")
	}
	s.Destroy()
}

func BenchmarkEqualRWWithInt(b *testing.B) {
	s := New()
	write := false
	for i := 0; i < b.N; i++ {
		if s.Pop() == nil || write {
			s.Push(i)
		} else {
			s.Pop()
			write = true
		}
	}
	s.Destroy()
}

func BenchmarkROnlyWithInt(b *testing.B) {
	s := New()
	nbr := b.N
	for i := 0; i < nbr; i++ {
		s.Push(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = s.Pop()
	}
	s.Destroy()
}

func BenchmarkWOnlyWithInt(b *testing.B) {
	s := New()
	for i := 0; i < b.N; i++ {
		s.Push(i)
	}
	s.Destroy()
}

/*
func BenchConcurrRW(b *testing.B) {
	s := New()
	b.RunParallel(func(pb *testing.PB) {
*/
