package locklessq

import (
	"sync"
	"testing"
)

func BenchmarkPop(b *testing.B) {
	var size int32 = 44100000
	q := New(size)

	for i := 0; i < int(size); i++ {
		q.Insert(1.0)
	}

	for n := 0; n < b.N; n++ {
		q.Pop()
	}
}

func TestFree(t *testing.T) {
	var s int32 = 100
	q := New(s)

	for i := 0; i < int(s); i++ {
		q.Insert(1)
	}

	if q.WriteAvailble() != 0 {
		t.Errorf("expected ")
	}

	if q.ReadAvailble() != s {
		t.Errorf("expected read")
	}

}

func TestThreads(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)

	var size int32 = 44100
	var pushes int32 = 80000
	var pops int32 = 50000

	q := New(size)
	var f float32 = 1.0
	start := false
	go func() {
		defer wg.Done()

		for i := 0; int32(i) < pushes; i++ {
			q.Insert(f)
			start = true
			f++
		}
		println("fff")
	}()

	go func() {
		defer wg.Done()

		for i := 0; int32(i) < pops; i++ {
			if start && q.Pop() == 0 {
				t.Errorf("got 0 q has %d free space", q.ReadAvailble())
				return
			}
		}
		println("ssss")

	}()

	wg.Wait()

}
