package locklessq

import (
	"sync/atomic"
)

type Q struct {
	Q      []float32
	reader int32
	writer int32
	size   int32
	space  int32
}

//Allocates slice with length and capacity of size.
func New(size int32) *Q {
	return &Q{Q: make([]float32, size, size), size: size, space: size}
}

//Inserts item to Q.
//Tests if Q has space, since its single writer, test is safe.
func (this *Q) Insert(f float32) bool {
	space := atomic.LoadInt32(&this.space)
	if space == 0 {
		return false
	}
	atomic.AddInt32(&this.space, -1)
	this.Q[this.writer] = f
	this.writer++
	this.writer %= this.size
	return true
}

func (q *Q) Pop() (float32, bool) {
	space := atomic.LoadInt32(&q.space)
	if space == q.size {
		return 0, false
	}
	atomic.AddInt32(&q.space, 1)
	ret := q.Q[q.reader]
	q.reader++
	q.reader %= q.size
	return ret, true
}

func (q *Q) ReadAvailble() int32 {
	return q.size - atomic.LoadInt32(&q.space)
}

func (q *Q) WriteAvailble() int32 {
	return atomic.LoadInt32(&q.space)
}
