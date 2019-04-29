package locklessq

import (
	"sync/atomic"
)

type Qfloat32 struct {
	Q      []float32
	reader int32
	writer int32
	size   int32
	space  int32
}

//Allocates slice with length and capacity of size.
func NewQfloat32(size int32) *Qfloat32 {
	return &Qfloat32{Q: make([]float32, size, size), size: size, space: size}
}

//Inserts item to Q.
//Tests if Q has space, since its single writer, test is safe.
func (this *Qfloat32) Insert(f float32) bool {
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

func (this *Qfloat32) Pop() (float32, bool) {
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

func (q *Qfloat32) ReadAvailble() int32 {
	return q.size - atomic.LoadInt32(&q.space)
}

func (q *Qfloat32) WriteAvailble() int32 {
	return atomic.LoadInt32(&q.space)
}

func (q *Qfloat32) Size() int32 {
	return int32(len(q.Q)) - atomic.LoadInt32(&q.space)
}

type Qint16 struct {
	Q      []int16
	reader int32
	writer int32
	size   int32
	space  int32
}

//Allocates slice with length and capacity of size.
func NewQint16(size int32) *Qint16 {
	return &Qint16{Q: make([]int16, size, size), size: size, space: size}
}

//Inserts item to Q.
//Tests if Q has space, since its single writer, test is safe.
func (this *Qint16) Insert(s int16) bool {
	space := atomic.LoadInt32(&this.space)
	if space == 0 {
		return false
	}
	this.Q[this.writer] = s
	atomic.AddInt32(&this.space, -1)
	this.writer++
	this.writer %= this.size
	return true
}

func (this *Qint16) Pop() (int16, bool) {
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

func (q *Qint16) ReadAvailble() int32 {
	return q.size - atomic.LoadInt32(&q.space)
}

func (q *Qint16) WriteAvailble() int32 {
	return atomic.LoadInt32(&q.space)
}

func (q *Qint16) Size() int32 {
	return int32(len(q.Q)) - atomic.LoadInt32(&q.space)
}
