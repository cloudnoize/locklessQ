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
	this.Q[this.writer] = f
	atomic.AddInt32(&this.space, -1)
	this.writer++
	this.writer %= this.size
	return true
}

func (this *Q) Pop() (float32, bool) {
	space := atomic.LoadInt32(&this.space)
	if space == this.size {
		return 0, false
	}
	ret := this.Q[this.reader]
	atomic.AddInt32(&this.space, 1)
	this.reader++
	this.reader %= this.size
	return ret, true
}

func (q *Q) ReadAvailble() int32 {
	return q.size - atomic.LoadInt32(&q.space)
}

func (q *Q) WriteAvailble() int32 {
	return atomic.LoadInt32(&q.space)
}
