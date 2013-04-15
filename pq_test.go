package pq

import (
	"testing"
	"time"
	"sync/atomic"
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

