package locklessq

import (
	"sync/atomic"
)

type Q struct {
	Q      []float32
	reader int32
	writer int32
	size   int32
	free   int32
}

func New(size int32) *Q {
	return &Q{Q: make([]float32, size, size), size: size, free: size}
}

func (q *Q) Insert(f float32) bool {
	free := atomic.LoadInt32(&q.free)
	if free == 0 {
		// println("locklessq insert returns false")
		return false
	}
	// println("Insert free", q.free)
	atomic.AddInt32(&q.free, -1)
	q.Q[q.writer] = f
	q.writer++
	q.writer %= q.size
	return true
}

func (q *Q) Pop() float32 {
	free := atomic.LoadInt32(&q.free)
	if free == q.size {
		println("retruning 0 -- free ", free, " size ", q.size)
		//panic("locklessq pop return 0")

		return 0
	}
	// println("Pop free", q.free)

	atomic.AddInt32(&q.free, 1)
	ret := q.Q[q.reader]
	q.reader++
	q.reader %= q.size
	return ret
}

func (q *Q) ReadAvailble() int32 {
	return q.size - atomic.LoadInt32(&q.free)
}

func (q *Q) WriteAvailble() int32 {
	return atomic.LoadInt32(&q.free)
}
