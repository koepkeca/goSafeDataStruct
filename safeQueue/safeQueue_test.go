package safeQueue

import (
	"testing"
)

func TestQueueCreation(t *testing.T) {
	q := New()
	if q.Len() != 0 {
		t.Errorf("Failed, invalid queue length.")
	}
	q.Destroy()
	return
}

func TestQueueLength(t *testing.T) {
	q := New()
	q.Enqueue(14)
	q.Enqueue(42)
	q.Enqueue("testing")
	q.Enqueue([]byte("Viper"))
	len := q.Len()
	if len != 4 {
		t.Errorf("Failed, invalid stack length, got %d expected 4", len)
	}
	q.Destroy()
	return
}

func TestQueueOrder(t *testing.T) {
	q := New()
	q.Enqueue(16)
	q.Enqueue(32)
	q.Enqueue(64)
	nv, ok := q.Dequeue().(int); 
	if !ok {
		t.Errorf("Failed, Dequeue got wrong type")
		q.Destroy()
		return
	}
	if nv != 16 {
		t.Errorf("Failed, got incorrect value order")
		q.Destroy()
		return
	}
	q.Destroy()
	return
}

func TestSizeAfterDequeue(t *testing.T) {
	q := New()
	q.Enqueue(16)
	q.Enqueue("test")
	q.Enqueue("私は笑い男だ")
	_ = q.Dequeue()
	_ = q.Dequeue()
	_ = q.Dequeue()
	if q.Len() != 0 {
		t.Errorf("Failed, poped through entire stack, yet size is non-zero")
	}
	q.Destroy()
	return
}

func TestEmptyDequeue(t *testing.T) {
	d := New()
	v := d.Dequeue()
	if v != nil {
		t.Errorf("Empty Pop got non-nil value")
	}
	d.Destroy()
}

func TestEmptyDequeueWithValues(t *testing.T) {
	q := New()
	q.Enqueue("Thingy")
	_ = q.Dequeue()
	v := q.Dequeue()
	if v != nil {
		t.Errorf("Empty stack with values got non-nil value")
	}
	q.Destroy()
}

func BenchmarkEqualRWWithInt(b *testing.B) {
	q := New()
	write := false
	for i := 0; i < b.N; i++ {
		if q.Dequeue() == nil || write {
			q.Enqueue(i)
		} else {
			q.Dequeue()
			write = true
		}
	}
	q.Destroy()
}

func BenchmarkROnlyWithInt(b *testing.B) {
	q := New()
	nbr := b.N
	for i := 0; i < nbr; i++ {
		q.Enqueue(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = q.Dequeue()
	}
	q.Destroy()
}

func BenchmarkWOnlyWithInt(b *testing.B) {
	q := New()
	for i := 0; i < b.N; i++ {
		q.Enqueue(i)
	}
	q.Destroy()
}
