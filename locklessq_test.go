package locklessq

import (
	"testing"
	"time"
)

func TestInsert(t *testing.T) {
	size := 100
	q := New(int32(size))
	for i := 0; i < size; i++ {
		ok := q.Insert(float32(i))
		if !ok {
			t.Error("Failed to insert")
		}
	}
	ok := q.Insert(float32(1.1))
	if ok {
		t.Error("Should have Failed to insert")
	}
}

func TestPop(t *testing.T) {
	size := 2
	q := New(int32(size))
	_, ok := q.Pop()
	if ok {
		t.Error("Should have Failed pop")
	}

	for i := 0; i < size; i++ {
		ok := q.Insert(float32(i))
		if !ok {
			t.Error("Failed to pop")
		}
	}

	v, ok := q.Pop()
	v, ok = q.Pop()
	if !ok {
		t.Error("Failed pop")
	}
	if v != float32(1) {
		t.Error("Expected ", float32(1), "got ", v)
	}

	v, ok = q.Pop()
	if ok {
		t.Error("Should have Failed pop")
	}
}

func TestStress(t *testing.T) {
	size := 100
	q := New(int32(size))
	cont := true
	go func() {
		for cont {
			q.Pop()
		}
	}()

	go func() {
		var v float32
		for cont {
			q.Insert(v)
			v++
		}
	}()

	time.Sleep(10 * time.Second)
	cont = false
}

func TestCorrectness(t *testing.T) {
	size := 100
	q := New(int32(size))
	q2 := New(int32(size))
	cont := true
	go func() {
		for cont {
			v, ok := q.Pop()
			if ok {
				ok = q2.Insert(v)
				if v >= float32(size) && ok {
					t.Errorf("shouldnt have happened %v", v)
					return
				}
			}
		}
		println("reader exit")
	}()

	go func() {
		v := float32(0)
		for cont {
			q.Insert(v)
			v++
		}
		println("writer exit")
	}()

	time.Sleep(3 * time.Second)
	cont = false
	time.Sleep(1 * time.Second)
	for i := 0; i < size; i++ {
		v, ok := q2.Pop()
		if !ok {
			t.Error("shouldnt have failed")
		}
		if i != int(v) {
			t.Errorf("left side %v != %v \nshould be incremental %v", i, v, q2.Q)
		}
	}
}
