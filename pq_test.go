package pq

import (
	"testing"
	"time"
	"sync/atomic"
	"runtime"
	"math/rand"
)

var inc int64

type testtask struct {
	p int
}

func (tt testtask) Run() {
	sleep()
}

func (tt testtask) Priority() int {
	return tt.p
}

func sleep() {
	time.Sleep(time.Second / 2)
	atomic.AddInt64(&inc, 1)
}

func finc() {
	atomic.AddInt64(&inc, 1)
}

func TestTask(t *testing.T) {
	n := 10
	q := new(Queue)
	if err := q.Start(n); err != nil {
		t.Errorf("Can't start queue: %v", err)
		return
	}
	
	defer q.Stop()
	
	st := time.Now()
	
	for i := 0; i < n; i++ {
		if err := q.AddFunc(sleep, 10); err != nil {
			t.Errorf("Can't add task: %v", err)
		}
	}
	// add last and wait
	if err := q.WaitFunc(sleep, 0); err != nil {
		t.Errorf("Can't add task: %v", err)
	}	
	stop := time.Since(st)
	
	if int(stop / time.Second) != 1 {
		t.Errorf("Unexpected stop time: %v", stop)
	} 
	
	if int(inc) != n+1 {
		t.Errorf("Unexpected inc: %v", inc)
	}
}

func TestGroup(t *testing.T) {
	inc = 0
	n := 10
	q := new(Queue)
	if err := q.Start(n); err != nil {
		t.Errorf("Can't start queue: %v", err)
		return
	}
	defer q.Stop()
	
	st := time.Now()
	tasks := make([]Task, n*2)
	for i := 0; i < n*2; i++ {
		tasks[i] = testtask{1}
	}
	// add last and wait
	if err := q.WaitGroup(tasks); err != nil {
		t.Errorf("Can't add tasks: %v", err)
	}	
	
	stop := time.Since(st)	
	
	if int(stop / time.Second) != 1 {
		t.Errorf("Unexpected stop time: %v", stop)
	} 	
	
	if int(inc) != n*2 {
		t.Errorf("Unexpected inc: %v", inc)
	}
}

func BenchmarkTaskSingleThread1(b *testing.B) {	
	benchTask(b, 1, 1, false)
}
func BenchmarkTaskSingleThread10(b *testing.B) {	
	benchTask(b, 10, 1, false)
}
func BenchmarkTaskMultiThread1(b *testing.B) {	
	benchTask(b, 1, runtime.NumCPU(), false)
}
func BenchmarkTaskMultiThread10(b *testing.B) {	
	benchTask(b, 10, runtime.NumCPU(), false)
}
func BenchmarkTaskMultiThread30(b *testing.B) {	
	benchTask(b, 30, runtime.NumCPU(), false)
}
func BenchmarkTaskMultiThread30RandPriority(b *testing.B) {	
	benchTask(b, 30, runtime.NumCPU(), true)
}

func benchTask(b *testing.B, w, c int, r bool) {
	runtime.GOMAXPROCS(c)
	q := new(Queue)
	q.Start(w)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p := 0
		if r {
			p = rand.Intn(10)
		}
		q.AddFunc(finc, p)
	}
	q.WaitFunc(finc, -1)
	defer q.Stop()
}