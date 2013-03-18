package pq

import (
	"container/heap"
	"fmt"
	"sync"
	"sync/atomic"
	"runtime"
	"log"
)

var ErrQueueNotStarted = fmt.Errorf("Queue not started or closed")
var ErrQueueAlreadyStarted = fmt.Errorf("Queue already started")

type Queue struct {
	numWorkers int
	pq priorityQueue
	work chan *item
	m sync.Mutex
	wg sync.WaitGroup
	working bool
	taskRunning int32
}

// Starts work. You may add tasks only after starting queue
func (q *Queue) Start(numWorkers int) (err error) {
	q.m.Lock()
	defer q.m.Unlock()
	if q.working {
		return ErrQueueAlreadyStarted
	}
	q.numWorkers = numWorkers
	q.pq = make(priorityQueue, 0)
	q.work = make(chan *item)
	if q.numWorkers <= 0 {
		q.numWorkers = runtime.NumCPU()
	}
	q.runWorkers()
	go q.dispatcher()
	return
}

// Add func() to queue
func (q *Queue) AddFunc(f func(), priority int) (err error) {
	task := &funcTask{
		f : f,
		p : priority,
	}
	return q.AddTask(task)
}

// Add func() to queue and wait while tasks will be done
func (q *Queue) WaitFunc(f func(), priority int) (err error) {
	task := &funcTask{
		f : f,
		p : priority,
	}
	return q.WaitTask(task)
}

// Just add group of tasks
func (q *Queue) AddGroup(tasks []Task) (err error) {
	for _, t := range tasks {
		if err = q.AddTask(t); err != nil {
			return
		}
	}
	return
}

// Add group of tasks and waits while all tasks will be done
func (q *Queue) WaitGroup(tasks []Task) (err error) {
	if len(tasks) == 0 {
		return
	}
	done := make(chan bool)
	for _, t := range tasks {
		it := &item{task:t, done:done}
		if err = q.addItem(it); err != nil {
			return
		}
	}
	l := len(tasks)
	for l > 0 {
		<- done
		l--
	}
	return
}

// Add single task to queue
func (q *Queue) AddTask(task Task) (err error) {
	it := &item{task:task}
	return q.addItem(it)
}

// Add single task to queue and waits while task will be done
func (q *Queue) WaitTask(task Task) (err error) {
	// add
	it := &item{task:task, done:make(chan bool)}
	if err = q.addItem(it); err != nil {
		return
	}
	// and wait
	<- it.done
	return
}

// Size of queue
func (q *Queue) Len() int {
	return len(q.pq)
}

// How much workers do work at this moment
func (q *Queue) TaskRunning() int {
	return int(atomic.LoadInt32(&q.taskRunning))
}

func (q *Queue) addItem(it *item) (err error) {
	q.m.Lock()
	defer q.m.Unlock()
	if ! q.working {
		return ErrQueueNotStarted
	}
	heap.Push(&q.pq, it)
	if q.pq.Len() == 1 {
		go q.dispatcher()
	}
	return
}

func (q *Queue) runWorkers() {
	for i := 0; i < q.numWorkers; i++ {
		go q.worker()
	}
	q.working = true
}

func (q *Queue) dispatcher() {
	for {
		q.m.Lock()
		if q.pq.Len() == 0 || ! q.working {
			q.m.Unlock()
			break
		}
		it := heap.Pop(&q.pq)
		q.work <- it.(*item)
		q.m.Unlock()	
	}
}

func (q *Queue) worker() {
	q.wg.Add(1)
	defer q.wg.Done()
	for it := range q.work {
		q.runTask(it)
	}
}

func (q *Queue) runTask(it *item) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("PQ. Panic while executing task: %v", r)
			if it.done != nil {
				it.done <- true
			}
		}
	}()
	atomic.AddInt32(&q.taskRunning, 1)
	defer atomic.AddInt32(&q.taskRunning, -1)
	it.task.Run()
	if it.done != nil {
		it.done <- true
	}
}

// Stopping queue. Wait while all workers finish current tasks
func (q *Queue) Stop() {
	q.m.Lock()
	defer q.m.Unlock()
	close(q.work)
	q.working = false
	q.wg.Wait()
}